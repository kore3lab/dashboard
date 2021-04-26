<template>
	<div class="content-wrapper">
		<div class="content-header">
			<div class="container-fluid">
				<c-navigator group="Storage"></c-navigator>
				<div class="row mb-2">
					<!-- title & search -->
					<div class="col-sm"><h1 class="m-0 text-dark"><span class="badge badge-info mr-2">P</span>Persistent Volume Claims</h1></div>
					<div class="col-sm-2"><b-form-select v-model="selectedNamespace" :options="namespaces()" size="sm" @input="query_All"></b-form-select></div>
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
		</div>

		<section class="content">
			<div class="container-fluid">
				<!-- count -->
				<div class="row mb-2">
					<div class="col-12 text-right "><span class="text-sm align-middle">Total : {{ totalItems }}</span></div>
				</div>
				<!-- GRID-->
				<div class="row">
					<div class="col-12">
						<div class="card">
							<div class="card-body table-responsive p-0">
								<b-table id="list" hover selectable show-empty select-mode="single" @row-selected="onRowSelected" @sort-changed="onSortChanged()" ref="selectableTable" :items="items" :fields="fields" :filter="keyword" :filter-included-fields="filterOn" @filtered="onFiltered" :current-page="currentPage" :per-page="$config.itemsPerPage" :busy="isBusy" class="text-sm">
									<template #table-busy>
										<div class="text-center text-success lh-vh-50">
											<b-spinner type="grow" variant="success" class="align-middle mr-2"></b-spinner>
											<span class="text-lg align-middle">Loading...</span>
										</div>
									</template>
									<template #empty="scope">
										<h4 class="text-center">does not exist.</h4>
									</template>
									<template v-slot:cell(name)="data">
										{{ data.value }}
									</template>
									<template v-slot:cell(storageClass)="data">
										<a href="#" @click="viewModel=getViewLink(data.value.group,data.value.rs,data.item.namespace, data.value.name); isShowSidebar=true;">{{ data.value.name }}</a>
									</template>
									<template v-slot:cell(pods)="data">
										<a href="#" v-for="(d, idx) in data.value" v-bind:key="idx" @click="viewModel=getViewLink('','pods',d[2], d[1]); isShowSidebar=true;">{{ d[1] }} </a>
									</template>
									<template v-slot:cell(accessModes)="data">
										<ul class="list-unstyled mb-0">
											<li v-for="d in data.item.accessModes" v-bind:key="d">{{ d }}</li>
										</ul>
									</template>
									<template v-slot:cell(status)="data">
										<div v-bind:class="data.item.status.style" class="text-sm">{{ data.item.status.value }}</div>
									</template>
									<template v-slot:cell(creationTimestamp)="data">
										{{ data.value.str }}
									</template>
								</b-table>
							</div>
							<b-pagination v-model="currentPage" :per-page="$config.itemsPerPage" :total-rows="totalItems" size="sm" align="center"></b-pagination>
						</div>
					</div>
				</div><!-- //GRID-->
			</div>
		</section>
		<b-sidebar v-model="isShowSidebar" width="50em" right shadow no-header>
			<c-view v-model="viewModel" @delete="query_All()" @close="onRowSelected"/>
		</b-sidebar>
	</div>
</template>
<script>
import VueNavigator from "@/components/navigator"
import VueView from "@/pages/view";

export default {
	components: {
		"c-navigator": { extends: VueNavigator },
		"c-view": { extends: VueView }
	},
	data() {
		return {
			selectedNamespace: "",
			keyword: "",
			filterOn: ["name"],
			fields: [
				{ key: "name", label: "Name", sortable: true },
				{ key: "namespace", label: "Namespace", sortable: true  },
				{ key: "storageClass", label: "Storage Class", sortable: true  },
				{ key: "capacity", label: "Capacity", sortable: true  },
				{ key: "pods", label: "Pods", sortable: true  },
				{ key: "accessModes", label: "Access Mode", sortable: true  },
				{ key: "creationTimestamp", label: "Age", sortable: true },
				{ key: "status", label: "Status", sortable: true  },
			],
			isBusy: false,
			items: [],
			currentPage: 1,
			totalItems: 0,
			pvcPod: [],
			podVersion: [],
			isShowSidebar: false,
			viewModel:{},
		}
	},
	layout: "default",
	created() {
		this.$nuxt.$on("navbar-context-selected", (ctx) => this.getPods() );
		if(this.currentContext()) this.$nuxt.$emit("navbar-context-selected");
	},
	methods: {
		onSortChanged() {
			this.currentPage = 1
		},
		onRowSelected(items) {
			if(items) {
				if(items.length) {
					this.viewModel = this.getViewLink('', 'persistentvolumeclaims', items[0].namespace, items[0].name)
					this.isShowSidebar = true
				} else {
					this.isShowSidebar = false
					this.$refs.selectableTable.clearSelected()
				}
			} else {
				this.isShowSidebar = false
				this.$refs.selectableTable.clearSelected()
			}
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
								creationTimestamp: this.getElapsedTime(el.metadata.creationTimestamp)
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
				if (el === this.pvcPod[i][0]) {
					list.push(this.pvcPod[i])
				}
			}
			return list
		},
		getPodname(el) {
			let pvclist = []
			if (el.spec.volumes[0].persistentVolumeClaim !== undefined)
			{
				pvclist.push(el.spec.volumes[0].persistentVolumeClaim.claimName)
				pvclist.push(el.metadata.name)
				pvclist.push(el.metadata.namespace)
			}
			if (pvclist.length !== 0) {
				this.pvcPod.push(pvclist)
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
