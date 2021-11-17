<template>
  <div class="content-wrapper">
		<section class="content-header">
			<div class="container-fluid">
				<c-navigator group="Cluster"></c-navigator>
				<div class="row mb-2">
					<div class="col-sm"><h1 class="m-0 text-dark"><span class="badge badge-info mr-2">N</span>Nodes</h1></div>
				</div>
			</div>
		</section>
		<section class="content">
			<div class="container-fluid">
				<!-- search & total count & items per page  -->
				<div class="row pb-2">
					<div class="col-sm-10"><c-search-form no-namespace no-label-selector @input="query_All" @keyword="(k)=>{keyword=k}"/></div>
					<div class="col-sm-2">
						<b-form inline class="float-right">
							<c-colums-selector name="grdSheet1" v-model="fields" :fields="fieldsAll" ></c-colums-selector>
							<i class="text-secondary ml-2 mr-2">|</i>
							<span class="text-sm align-middle">Total : {{ totalItems }}</span>
						</b-form>
					</div>
				</div>
				<!-- GRID-->
				<div class="row">
					<div class="col-12">
						<div class="card">
							<div class="card-body table-responsive p-0">
								<b-table hover selectable show-empty select-mode="single" @row-selected="onRowSelected" @sort-changed="currentPage=1" ref="grdSheet1" :items="items" :fields="fields" :filter="keyword" :filter-included-fields="filterOn" @filtered="onFiltered" :busy="isBusy" class="text-sm">
									<template #table-busy>
										<div class="text-center text-success lh-vh-50">
											<b-spinner type="grow" variant="success" class="align-middle mr-2"></b-spinner>
											<span class="text-lg align-middle">Loading...</span>
										</div>
									</template>
									<template v-slot:cell(cpu)="data">
										<b-progress :value="data.value" :max="100" variant="info" show-value class="mb-3"></b-progress>
									</template>
									<template v-slot:cell(memory)="data">
										<b-progress :value="data.value" :max="100" variant="info" show-value class="mb-3"></b-progress>
									</template>
									<template v-slot:cell(storage)="data">
										<b-progress :value="data.value" :max="100" variant="info" show-value class="mb-3"></b-progress>
									</template>
								</b-table>
							</div>
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
				{ key: "roles", label: "Roles", sortable: true },
				{ key: "cpu", label: "CPU", sortable: true  },
				{ key: "memory", label: "Memory", sortable: true  },
				{ key: "storage", label: "Storage", sortable: true },
				{ key: "pods", label: "Pods", sortable: true },
				{ key: "version", label: "Version", sortable: true },
				{ key: "creationTimestamp", label: "Age", sortable: true, formatter: this.getElapsedTime },
				{ key: "status", label: "Status", sortable: true },
			],
			isBusy: false,
			items: [],
			currentItems:[],
			selectIndex: 0,
			totalItems: 0,
			metrics: [],
			isShowSidebar: false,
			viewModel:{},
		}
	},
	layout: "default",
	methods: {
		onRowSelected(items) {
			this.isShowSidebar = (items && items.length > 0)
			if (this.isShowSidebar) this.viewModel = this.getViewLink('', 'nodes', items[0].namespace, items[0].name)
		},
		// 조회
		query_All() {
			this.isBusy = true;
			this.$axios.get(`/api/clusters/${this.currentContext()}/nodes`)
					.then((resp) => {
						this.items = [];
						if(resp.data.nodes) {
							for(let n in resp.data.nodes) {
								this.items.push({
									name: n,
									cpu: resp.data.nodes[n].metrics.percent.cpu,
									memory: resp.data.nodes[n].metrics.percent.memory,
									storage: resp.data.nodes[n].metrics.percent.storage,
									pods: `${resp.data.nodes[n].metrics.usage.pods}/${resp.data.nodes[n].metrics.allocatable.pods}`,
									version: resp.data.nodes[n].version,
									roles: resp.data.nodes[n].role,
									creationTimestamp: resp.data.nodes[n].creationTimestamp,
									status: resp.data.nodes[n].status
								});
							}
						}
						this.onFiltered(this.items);
					})
					.catch(e => { this.msghttp(e);})
					.finally(()=> { this.isBusy = false;});
		},
		onFiltered(filteredItems) {
			this.totalItems = filteredItems.length;
		}
	}
}
</script>
