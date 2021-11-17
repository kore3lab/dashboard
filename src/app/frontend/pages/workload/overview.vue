<template>
  <div class="content-wrapper">
		<section class="content-header">
			<div class="container-fluid">
				<c-navigator group="Workload"></c-navigator>
				<div class="row mb-2">
					<!-- title & search -->
					<div class="col-sm"><h1 class="m-0 text-dark"><span class="badge badge-info mr-2">O</span>Overview</h1><b-icon-plus-square class="ml-2"/></div>
				</div>
			</div>
		</section>

		<section class="content">
			<div class="container-fluid">
				<!-- total count & items per page  -->
				<div class="row mb-2">
					<div class="col-12">
						<b-form inline>
							<b-form-select v-model="selectedNamespace" :options="namespaces()" size="sm" @input="query_All(); selectNamespace(selectedNamespace);"></b-form-select>
							<span class="text-sm mr-2 ml-2">Group by</span>
							<b-form-input size="sm" list="listGroupBy" value="app"></b-form-input>
							<datalist id="listGroupBy">
								<option v-for="d in groups" :key="d">{{ d }}</option>
							</datalist>
							<b-form-tags class="form-labels-filter" size="sm" v-model="labelFilters" tag-pills tag-variant="info" separator=" ," placeholder="Enter new label filter"></b-form-tags>
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
									<template v-slot:cell(jobStatus)="data">
										<div v-bind:key="data.value" v-bind:class="{'text-success':data.value=='Complete', 'text-danger':data.value=='Failed'}">{{ data.value }}</div>
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
			groups: ["app","app.kubernetes.io/name", "app.kubernetes.io/part-of"],
			labelFilters: ["app=kore-board"],
			filterOn: ["name"],
			fields: [],
			fieldsAll: [
				{ key: "name", label: "Name", sortable: true },
				{ key: "namespace", label: "Namespace", sortable: true  },
				{ key: "completions", label: "Completions", sortable: true },
				{ key: "creationTimestamp", label: "Age", sortable: true, formatter: this.getElapsedTime },
				{ key: "jobStatus", label: "Status", sortable: true  },
			],
			isBusy: true,
			items: [],
			itemsPerPage: this.$storage.global.get("itemsPerPage",10),
			currentPage: 1,
			totalItems: 0,
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
		this.$nuxt.$on("context-selected", (ctx) => {
			this.selectedNamespace = this.selectNamespace()
			this.query_All()
		});
		if(this.currentContext()) this.$nuxt.$emit("context-selected");
	},
	methods: {
		onRowSelected(items) {
			this.isShowSidebar = (items && items.length > 0)
			if (this.isShowSidebar) this.viewModel = this.getViewLink('batch', 'jobs', items[0].namespace, items[0].name)
		},
		// 조회
		query_All() {
			this.isBusy = true;
			this.$axios.get(this.getApiUrl("batch","jobs",this.selectedNamespace))
					.then((resp) => {
						this.items = [];
						resp.data.items.forEach(el => {
							this.items.push({
								name: el.metadata.name,
								namespace: el.metadata.namespace,
								completions: `${el.status.succeeded?el.status.succeeded:"0"} / ${el.spec.completions?el.spec.completions:"1"}`,
								jobStatus: el.status.conditions? el.status.conditions[0].type : "Failed",
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
		}
	},
	beforeDestroy(){
		this.$nuxt.$off('context-selected')
	}
}
</script>

				<!-- search & total count & items per page  -->
<!--
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


import VueSearchForm		from "@/components/list/searchForm"
		"c-search-form": { extends: VueSearchForm},


(d && d.namespace) ? d.namespace: this.selectNamespace(), "", d && d.labelSelector? `labelSelector=${d.labelSelector}`: ""
"", "", d && d.labelSelector? `labelSelector=${d.labelSelector}`: ""
-->
