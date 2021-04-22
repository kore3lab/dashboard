<template>
	<div class="content-wrapper">
		<div class="content-header">
			<div class="container-fluid">
				<c-navigator group="Workload"></c-navigator>
				<div class="row mb-2">
					<!-- title & search -->
					<div class="col-sm"><h1 class="m-0 text-dark"><span class="badge badge-info mr-2">D</span>Deployments</h1></div>
					<div class="col-sm-2"><b-form-select v-model="selectedNamespace" :options="namespaces()" size="sm" @input="query_All"></b-form-select></div>
					<div class="col-sm-2 float-left">
						<div class="input-group input-group-sm" >
							<b-form-input id="txtKeyword" v-model="keyword" class="form-control float-right" placeholder="Search"></b-form-input>
							<div class="input-group-append"><button type="submit" class="btn btn-default" @click="query_All"><i class="fas fa-search"></i></button></div>
						</div>
					</div>
					<!-- button -->
					<div class="col-sm-1 text-right">
						<b-button variant="primary" size="sm" @click="$router.push(`/create?context=${currentContext()}&group=Workload&crd=Deployment`)">Create</b-button>
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
									<template v-slot:cell(status)="data">
										<div class="list-unstyled mb-0" v-if="data.item.status">
											<span v-for="(val, idx) in data.item.status" v-bind:key="idx" v-bind:class="val.style" class=" text-sm ml-1 ">{{ val.type }}</span>
										</div>
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
import axios		from "axios"
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
				{ key: "pods", label: "Pods", sortable: true  },
				{ key: "replicas", label: "Replicas", sortable: true  },
				{ key: "creationTimestamp", label: "Age", sortable: true  },
				{ key: "status", label: "Status", sortable: true  },
			],
			isBusy: false,
			items: [],
			currentPage: 1,
			totalItems: 0,
			isShowSidebar: false,
			viewModel:{},
		}
	},
	layout: "default",
	created() {
		this.$nuxt.$on("navbar-context-selected", (ctx) => this.query_All() );
		if(this.currentContext()) this.$nuxt.$emit("navbar-context-selected");
	},
	methods: {
		onSortChanged() {
			this.currentPage = 1
		},
		onRowSelected(items) {
			if(items) {
				if(items.length) {
					this.viewModel = this.getViewLink('apps', 'deployments', items[0].namespace, items[0].name)
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
			axios.get(this.getApiUrl("apps","deployments",this.selectedNamespace))
					.then((resp) => {
						this.items = [];
						resp.data.items.forEach(el => {
							this.items.push({
								name: el.metadata.name,
								namespace: el.metadata.namespace,
								pods: this.toPods(el.status),
								replicas: el.status.replicas ? el.status.replicas : 0,
								creationTimestamp: this.getElapsedTime(el.metadata.creationTimestamp),
								status: this.toStatus(el.status.conditions),
							});
						});
						this.onFiltered(this.items);
					})
					.catch(e => { this.msghttp(e);})
					.finally(()=> { this.isBusy = false;});
		},
		toPods(status) {
			var podsReady = 0
			var podsLength = 0
			if ( status.readyReplicas ) podsReady = status.readyReplicas
			if ( status.availableReplicas ) podsLength = status.availableReplicas
			return `${podsReady}/${podsLength}`
		},
		toStatus(conditions) {
			if(!conditions) return
			let condition = []
			conditions.forEach(el => {
				condition.push({
					type: el.type,
					style: this.checkStyle(el.type),
				})
			})
			condition.sort(function(a,b) {
				return a.type < b.type ? -1 : a.type > b.type ? 1 : 0;
			})
			return condition
		},
		checkStyle(t) {
			if(t === 'Progressing') return 'text-primary'
			if(t === 'Available') return 'text-success'
			else return 'text-danger'
		},
		onFiltered(filteredItems) {
			this.totalItems = filteredItems.length;
			this.currentPage = 1
		}
	},
	beforeDestroy(){
		this.$nuxt.$off('navbar-context-selected')
	}
}
</script>
<style scoped>
</style>
