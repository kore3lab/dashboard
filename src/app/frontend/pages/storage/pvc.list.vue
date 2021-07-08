<template>
	<div class="content-wrapper">
		<section class="content-header">
			<div class="container-fluid">
				<c-navigator group="Storage"></c-navigator>
				<div class="row mb-2">
					<!-- title & search -->
					<div class="col-sm"><h1 class="m-0 text-dark"><span class="badge badge-info mr-2">P</span>Persistent Volume Claims</h1></div>
					<div class="col-sm-2"><b-form-select v-model="selectedNamespace" :options="namespaces()" size="sm" @input="query_All(); selectNamespace(selectedNamespace);"></b-form-select></div>
					<div class="col-sm-2 float-left">
						<div class="input-group input-group-sm" >
							<b-form-input id="txtKeyword" v-model="keyword" class="form-control float-right" placeholder="Search"></b-form-input>
							<div class="input-group-append"><button type="submit" class="btn btn-default" @click="query_All"><i class="fas fa-search"></i></button></div>
						</div>
					</div>
					<!-- button -->
					<div class="col-sm-1 text-right">
						<b-button variant="primary" size="sm" @click="$router.push(`/create?context=${currentContext()}&group=Storage&crd=PersistentVolumeClaim`)">Create</b-button>
					</div>
				</div>
			</div>
		</section>

		<section class="content">
			<div class="container-fluid">
				<!-- total count & items per page  -->
				<div class="d-flex flex-row-reverse">
					<div class="p-2">
						<b-form inline>
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
										<a href="#" @click="viewModel=getViewLink(data.value.group,data.value.rs,data.item.namespace, data.value.name); isShowSidebar=true;">{{ data.value.name }}</a>
									</template>
									<template v-slot:cell(pods)="data">
										<a href="#" v-for="(d, idx) in data.value" v-bind:key="idx" @click="viewModel=getViewLink('','pods',d.podNamespace, d.podName); isShowSidebar=true;">{{ d.podName }} </a>
									</template>
									<template v-slot:cell(accessModes)="data">
										<ul class="list-unstyled mb-0">
											<li v-for="d in data.item.accessModes" v-bind:key="d">{{ d }}</li>
										</ul>
									</template>
									<template v-slot:cell(status)="data">
										<div v-bind:class="data.item.status.style" class="text-sm">{{ data.item.status.value }}</div>
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
import VueColumsSelector	from "@/components/columnsSelector"
import VueView				from "@/pages/view";

export default {
	components: {
		"c-navigator": { extends: VueNavigator },
		"c-colums-selector": { extends: VueColumsSelector},
		"c-view": { extends: VueView }
	},
	data() {
		return {
			selectedNamespace: "",
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
	created() {
		this.$nuxt.$on("navbar-context-selected", (ctx) => {
			this.selectedNamespace = this.selectNamespace()
			this.getPods()
		});
		if(this.currentContext()) this.$nuxt.$emit("navbar-context-selected");
	},
	methods: {
		onRowSelected(items) {
			this.isShowSidebar = (items && items.length > 0)
			if (this.isShowSidebar) this.viewModel = this.getViewLink('', 'persistentvolumeclaims', items[0].namespace, items[0].name)
		},
		// 조회
		query_All() {
			this.isBusy = true;
			this.$axios.get(this.getApiUrl("","persistentvolumeclaims",this.selectedNamespace))
					.then((resp) => {
						this.items = [];
						resp.data.items.forEach(el => {
							this.items.push({
								name: el.metadata.name,
								namespace: el.metadata.namespace,
								labels: el.metadata.labels,
								status: this.getStatus(el.status.phase),
								capacity: el.spec.resources.requests.storage ? el.spec.resources.requests.storage: "",
								pods: this.getPvc(el.metadata.name),
								accessModes: el.spec.accessModes,
								storageClass: this.getStorageClass(el.spec.storageClassName),
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
		getStorageClass(name) {
			return {
				"name" : name,
				"group" : "storage.k8s.io",
				"rs" : "storageclasses"
			}
		},
		// pod List 조회 이후 query_All 실행
		getPods() {
			this.isBusy = true;
			this.$axios.get(this.getApiUrl("","pods",this.selectedNamespace))
					.then((resp) => {
						resp.data.items.forEach(el => {
							this.getPodname(el)
						});
					})
					.catch(e => { this.msghttp(e);})
					.finally(()=> this.query_All());
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
		},
		// pvc 상태 체크
		getStatus(status) {
			if (status === "Bound") {
				return {
					"value": "Bound",
					"style": "text-success",
				}
			} else if (status === "Pending") {
				return {
					"value": "Pending",
					"style": "text-warning",
				}
			} else {
				return {
					"value": "Lost",
					"style": "text-secondary",
				}
			}
		}
	},
	beforeDestroy(){
		this.$nuxt.$off('navbar-context-selected')
	}
}
</script>
<style scoped>label {font-weight: 500;}</style>
