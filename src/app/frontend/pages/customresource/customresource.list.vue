<template>
  <div class="content-wrapper">
		<section class="content-header">
			<div class="container-fluid">
				<c-navigator :group="'Custom Resource / '+ crdQuery.group"></c-navigator>
				<div class="row mb-2">
					<div class="col-sm">
						<h1 class="m-0 text-dark"><span class="badge badge-info mr-2">{{crdQuery.name.charAt(0)}}</span><span>{{ crdQuery.name }}</span></h1>
					</div>
					<div class="col-sm-1 text-right">
						<b-button variant="primary" size="sm" @click="$router.push(`customresource.create?group=${crdQuery.group}&crd=${crdQuery.crd}&version=${crdQuery.version}&name=${crdQuery.name}`)">Create</b-button>
					</div>
				</div>
			</div>
		</section>

		<section class="content">
			<div class="container-fluid">
				<!-- search & total count  -->
				<div class="row pb-2">
					<div class="col-sm-10"><c-search-form no-label-selector :no-namespace="!isNamespaced" @input="query_All" @keyword="(k)=>{keyword=k}"/></div>
					<div class="col-sm-2">
						<b-form inline class="float-right">
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
										<div class="text-center text-success" style="margin:150px 0">
											<b-spinner type="grow" variant="success" class="align-middle mr-2"></b-spinner>
											<span class="align-middle text-lg">Loading...</span>
										</div>
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
import VueNavigator 	from "@/components/navigator"
import VueView			from "@/pages/view";
import VueSearchForm	from "@/components/list/searchForm"
import jsonpath			from "jsonpath"

export default {
	components: {
		"c-navigator": { extends: VueNavigator },
		"c-search-form": { extends: VueSearchForm},
		"c-view": { extends: VueView }
	},
	data() {
		return {
			keyword: "",
			filterOn: [".metadata.name"],
			fields: [],
			isBusy: false,
			origin: [],
			items: [],
			itemsPerPage: this.$storage.global.get("itemsPerPage",10),
			currentPage: 1,
			totalItems: 0,
			isShowSidebar: false,
			viewModel:{},
			isNamespaced: false
		}
	},
	computed: {
		crdQuery: {
			get () {
				return {
					group: this.$route.query.group,		//ex. networking.istio.io
					crd: this.$route.query.crd,			//ex. virtualservices
					name: this.$route.query.name,		//ex. VirtualService
					version: this.$route.query.version	//ex. v1beta1
				}
			},
		},
	},
	layout: "default",
	watch: {
		itemsPerPage(n) {
			this.$storage.global.set("itemsPerPage",n)
		},
		crdQuery(d) {
			if(this.currentContext()) this.$nuxt.$emit("context-selected");
		}
	},
	watchQuery: ["group","crd", "version", "name"],
	created() {
		this.$nuxt.$on("context-selected", () => {
			//crd 정보조회
			this.$axios.get(`${this.getApiUrl("apiextensions.k8s.io","customresourcedefinitions")}/${this.crdQuery.crd}.${this.crdQuery.group}`)
				.then((resp) => {
					//column - name
					this.fields.push({ key: ".metadata.name", label: "Name", sortable: true })

					//column - namespace
					this.isNamespaced = (resp.data.spec.scope=="Namespaced")
					if (this.isNamespaced) this.fields.push({key: ".metadata.namespace", label: "Namespace", sortable: true})

					// custom columns (https://kubernetes.io/docs/reference/using-api/deprecation-guide/#customresourcedefinition-v122)
					let additionalPrinterColumns;
					if( resp.data.apiVersion == "apiextensions.k8s.io/v1") {
						const version = resp.data.spec.versions.find(el => el.name= this.crdQuery.version);
						if(version) additionalPrinterColumns = version.additionalPrinterColumns;
					} else {
						//v1beta1
						additionalPrinterColumns = resp.data.spec.additionalPrinterColumns;
					}
					if(additionalPrinterColumns) {
						additionalPrinterColumns.forEach(c => {
							this.fields.push({key: c.jsonPath || c.JSONPath, label: c.name, sortable: true})
						});
					}
					//column - age
					if(! this.fields.find(e => e.key == ".metadata.creationTimestamp") ) {
						this.fields.push({key: ".metadata.creationTimestamp", label: "Age", sortable: true, formatter: this.getElapsedTime})
					}
					this.query_All();
				})
				.catch(e => { 
					if (e.response && e.response.status == 404) {
						return this.$nuxt.error({
							statusCode: 404, 
							redirect: "/customresource/customresourcedefinitions.list",
							message: `${this.crdQuery.crd} ${this.crdQuery.group}/${this.crdQuery.version} is not found` });
					} else this.msghttp(e);
				})
		});
		if(this.currentContext()) this.$nuxt.$emit("context-selected");

	},
	methods: {
		onRowSelected(items) {
			this.isShowSidebar = (items && items.length > 0)
			if (this.isShowSidebar) this.viewModel = this.getViewLink(this.crdQuery.group, this.crdQuery.crd, items[0][".metadata.namespace"], items[0][".metadata.name"])
		},
		// 조회
		query_All(d) {
			this.isBusy = true;
			this.$axios.get(this.getApiUrl(this.crdQuery.group, this.crdQuery.crd, (d && d.namespace) ? d.namespace: this.selectNamespace()))
				.then((resp) => {
					this.items = [];
					resp.data.items.forEach(el => {
						let it = {}
						this.fields.forEach(e => {
							if(e.key==".metadata.creationTimestamp") {
								it[e.key] = el.metadata.creationTimestamp
							} else {
								try {
									let v = jsonpath.query(el,`$.${e.key}`)
									it[e.key] = (v && v.length > 0) ? v[0]: ""
								} catch(_) {}
							}
						});
						this.items.push(it);
					});
					this.origin = this.items;
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
		this.$nuxt.$off("context-selected")
	}
}
</script>
