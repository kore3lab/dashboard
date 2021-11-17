<template>
  <div class="content-wrapper">
		<section class="content-header">
			<div class="container-fluid">
				<c-navigator group="Storage"></c-navigator>
				<div class="row mb-2">
					<div class="col-sm"><h1 class="m-0 text-dark"><span class="badge badge-info mr-2">P</span>Persistent Volume Claims</h1></div>
					<div class="col-sm-1 text-right">
						<b-button variant="primary" size="sm" @click="$router.push(`/create?context=${currentContext()}&group=Storage&crd=PersistentVolumeClaim`)">Create</b-button>
					</div>
				</div>
			</div>
		</section>

		<section class="content">
			<div class="container-fluid">
				<!-- search & total count & items per page  -->
				<div class="row pb-2">
					<div class="col-sm-10"><c-search-form  @input="query_All" @keyword="(k)=>{keyword=k}"/></div>
					<div class="col-sm-2">
						<b-form inline class="float-right">
							<c-colums-selector name="grdSheet1" v-model="fields" :fields="fieldsAll" ></c-colums-selector>
							<i class="text-secondary ml-2 mr-2">|</i>
							<b-form-select size="sm" :options="this.var('ITEMS_PER_PAGE')" v-model="itemsPerPage"></b-form-select>
							<span class="text-sm align-middle ml-2">Total : {{ totalItems }}</span>
						</b-form>
					</div>
				</div>
				<!-- GRID-->
				<div class="row">
					<div class="col-12">
						<div class="card">
							<div class="card-body table-responsive p-0">
								<b-table hover selectable show-empty select-mode="single" @row-selected="onRowSelected" @sort-changed="currentPage=1" ref="grdSheet1" :items="items" :fields="fields" :filter="keyword" :filter-included-fields="filterOn" @filtered="onFiltered" :current-page="currentPage" :per-page="itemsPerPage" :busy="isBusy" class="text-sm">
									<template #table-busy>
										<div class="text-center text-success lh-vh-50">
											<b-spinner type="grow" variant="success" class="align-middle mr-2"></b-spinner>
											<span class="text-lg align-middle">Loading...</span>
										</div>
									</template>
									<template v-slot:cell(storageClass)="data">
										<a href="#" @click="viewModel=getViewLink('storage.k8s.io','storageclasses','', data.value); isShowSidebar=true;">{{ data.value }}</a>
									</template>
									<template v-slot:cell(pods)="data">
										<ul class="list-unstyled m-0 p-0">
											<li v-for="(d, idx) in data.value" v-bind:key="idx"><a href="#" @click="viewModel=getViewLink('','pods',d.podNamespace, d.podName); isShowSidebar=true;">{{ d.podName }} </a></li>
										</ul>
									</template>
									<template v-slot:cell(accessModes)="data">
										<ul class="list-unstyled mb-0">
											<li v-for="d in data.item.accessModes" v-bind:key="d">{{ d }}</li>
										</ul>
									</template>
									<template v-slot:cell(status)="data">
										<div v-bind:class="{'text-success': data.value=='Bound'}" class="text-sm">{{ data.value }}</div>
									</template>
								</b-table>
							</div>
							<b-pagination v-model="currentPage" :per-page="itemsPerPage" :total-rows="totalItems" size="sm" align="center"></b-pagination>
						</div>
					</div>
				</div><!-- //GRID-->
			</div>
		</section>
		<b-sidebar v-model="isShowSidebar" width="50em" @hidden="$refs.grdSheet1.clearSelected()" right shadow no-header>
			<c-view v-model="viewModel" @delete="query_All()" @close="isShowSidebar=false"/>
		</b-sidebar>
	</div>
</template>
<script>
import VueNavigator			from "@/components/navigator"
import VueColumsSelector	from "@/components/list/columnsSelector"
import VueSearchForm		from "@/components/list/searchForm"
import VueView				from "@/pages/view";

export default {
	components: {
		"c-navigator": { extends: VueNavigator },
		"c-colums-selector": { extends: VueColumsSelector},
		"c-search-form": { extends: VueSearchForm},
		"c-view": { extends: VueView }
	},
	data() {
		return {
			keyword: "",
			filterOn: ["name"],
			fields: [],
			fieldsAll: [
				{ key: "name", label: "Name", sortable: true },
				{ key: "namespace", label: "Namespace", sortable: true  },
				{ key: "storageClass", label: "Storage Class", sortable: true  },
				{ key: "capacity", label: "Capacity", sortable: true  },
				{ key: "pods", label: "Pods", sortable: true  },
				{ key: "accessModes", label: "Access Mode", sortable: true  },
				{ key: "creationTimestamp", label: "Age", sortable: true, formatter: this.getElapsedTime },
				{ key: "status", label: "Status", sortable: true  },
			],
			isBusy: false,
			items: [],
			itemsPerPage: this.$storage.global.get("itemsPerPage",10),
			currentPage: 1,
			totalItems: 0,
			pvcPod: [],
			podVersion: [],
			isShowSidebar: false,
			viewModel:{},
		}
	},
	watch: {
		itemsPerPage(n) {
			this.$storage.global.set("itemsPerPage",n)
		}
	},
	layout: "default",
	methods: {
		onRowSelected(items) {
			this.isShowSidebar = (items && items.length > 0)
			if (this.isShowSidebar) this.viewModel = this.getViewLink('', 'persistentvolumeclaims', items[0].namespace, items[0].name)
		},
		// 조회
		async query_All(d) {
			this.isBusy = true;

			const ns = (d && d.namespace) ? d.namespace: this.selectNamespace();
			const labels = d && d.labelSelector? `labelSelector=${d.labelSelector}`: "";

			// 개선필요
			try {
				let resp = await this.$axios.get(this.getApiUrl("","pods", ns, "", labels))
				resp.data.items.forEach(el => {
					this.getPodname(el)
				});
			} catch (e) {
				this.msghttp(e);
			}

			this.$axios.get(this.getApiUrl("","persistentvolumeclaims", ns, "", labels))
					.then((resp) => {
						this.items = [];
						resp.data.items.forEach(el => {
							this.items.push({
								name: el.metadata.name,
								namespace: el.metadata.namespace,
								labels: el.metadata.labels,
								status: el.status.phase,
								capacity: el.spec.resources.requests.storage ? el.spec.resources.requests.storage: "",
								pods: this.getPvc(el.metadata.name),
								accessModes: el.spec.accessModes,
								storageClass: el.spec.storageClassName,
								creationTimestamp: el.metadata.creationTimestamp
							});
						});
						this.onFiltered(this.items);
					})
					.catch(e => { this.msghttp(e);})
					.finally(()=> { this.isBusy = false;});
		},
		onFiltered(filteredItems) {
			this.totalItems = filteredItems.length;
			this.currentPage = 1
		},
		formatterPods(value, key, item) {
			return {
				name: item.name,
				namespace: item.namepsace
			}
		},
		getPvc(el) {
			let list = []
			for(let i=0;i<this.pvcPod.length;i++) {
				if (el === this.pvcPod[i].claimName) {
					list.push(this.pvcPod[i])
				}
			}
			return list
		},
		getPodname(el) {
			let pvclist = []
			if(el.spec.volumes) {
				el.spec.volumes.filter(volume => {
					if(volume.persistentVolumeClaim && volume.persistentVolumeClaim.claimName) {
						pvclist = {
							claimName : volume.persistentVolumeClaim.claimName,
							podName : el.metadata.name,
							podNamespace: el.metadata.namespace,
						}
					}
				})
				if (pvclist.length !== 0) {
					this.pvcPod.push(pvclist)
				}
			}
		}
	}
}
</script>
