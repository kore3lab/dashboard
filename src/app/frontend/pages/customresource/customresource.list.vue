<template>
  <div class="content-wrapper" v-bind:style="{paddingBottom: contentPadding()}">
		<section class="content-header">
			<div class="container-fluid">
				<c-navigator :group="'Custom Resource / '+ apiQuery.group"></c-navigator>
				<div class="row mb-2">
					<!-- title & search -->
					<div class="col-sm"><h1 class="m-0 text-dark"><span class="badge badge-info mr-2 text-capitalize">{{apiQuery.name.charAt(0)}}</span><span class="text-capitalize">{{apiQuery.name}}</span></h1></div>
					<div v-if="$route.query.scope === 'Namespaced' || $route.query.ns === 'true'" class="col-sm-2"><b-form-select v-model="selectedNamespace" :options="namespaces()" size="sm" @input="query_All(); selectNamespace(selectedNamespace);"></b-form-select></div>
					<div class="col-sm-2 float-left">
						<div class="input-group input-group-sm" >
							<b-form-input id="txtKeyword" v-model="keyword" class="form-control float-right" placeholder="Search"></b-form-input>
							<div class="input-group-append"><button type="submit" class="btn btn-default" @click="query_All"><i class="fas fa-search"></i></button></div>
						</div>
					</div>
					<!-- button -->
					<div class="col-sm-1 text-right">
						<b-button variant="primary" size="sm" @click="$router.push(`customresource.create?group=${apiQuery.group}&name=${apiQuery.name}&crd=${crdQuery.name}&version=${crdQuery.version}`)">Create</b-button>
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
import VueNavigator from "@/components/navigator"
import VueView		from "@/pages/view";
import jsonpath		from "jsonpath"

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
			fields: [],
			isBusy: false,
			origin: [],
			items: [],
			itemsPerPage: this.$storage.global.get("itemsPerPage",10),
			currentPage: 1,
			totalItems: 0,
			isShowSidebar: false,
			viewModel:{},
			crdQuery: {
				name: this.$route.query.crd,
				version: this.$route.query.version
			},
			apiQuery: {
				group: "",
				name: ""
			}
		}
	},
	layout: "default",
	watch: {
		itemsPerPage(n) {
			this.$storage.global.set("itemsPerPage",n)
		},
		crdQuery(d) {
			if(this.currentContext()) this.$nuxt.$emit("navbar-context-selected");
		}
	},
	watchQuery: ["crd","version"],
	async asyncData({ query }) {
		return {crdQuery : {name: query.crd, version: query.version}}
	},
	created() {
		this.$nuxt.$on("navbar-context-selected", () => {

			//crd 정보조회
			this.$axios.get(`${this.getApiUrl("apiextensions.k8s.io","customresourcedefinitions")}/${this.crdQuery.name}`)
				.then((resp) => {
					//컬럼 정의
					this.fields.push({ key: ".metadata.name", label: "Name", sortable: true })
					if (resp.data.spec.scope=="Namespaced") this.fields.push({key: ".metadata.namespace", label: "Namespace", sortable: true})
					let version = resp.data.spec.versions.find(el => el.name= this.crdQuery.version);

					if(version.additionalPrinterColumns) {
						version.additionalPrinterColumns.forEach(c => {
							this.fields.push({key: c.jsonPath, label: c.name, sortable: true})
						});
					}
					if(! this.fields.find(e => e.key == ".metadata.creationTimestamp") ) {
						this.fields.push({key: ".metadata.creationTimestamp", label: "Age", sortable: true, formatter: this.getElapsedTime})
					}
					// crd 정보 (group, name)
					this.apiQuery = {group: resp.data.spec["group"], name: resp.data.spec.names["plural"]}
					this.query_All();
				})
				//.catch(e => { this.msghttp(e);})
		});
		if(this.currentContext()) this.$nuxt.$emit("navbar-context-selected");

	},
	methods: {
		onRowSelected(items) {
			this.isShowSidebar = (items && items.length > 0)
			if (this.isShowSidebar) this.viewModel = this.getViewLink(this.apiQuery.group, this.apiQuery.name, items[0][".metadata.namespace"], items[0][".metadata.name"])
		},
		// 조회
		query_All() {
			this.isBusy = true;
			this.$axios.get(this.getApiUrl(this.apiQuery.group, this.apiQuery.name,this.selectedNamespace))
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
				//.catch(e => { this.msghttp(e);})
				.finally(()=> { this.isBusy = false;});
		},
		onFiltered(filteredItems) {
			this.totalItems = filteredItems.length;
			this.currentPage = 1
		},
	},
	beforeDestroy(){
		this.$nuxt.$off("navbar-context-selected")
	}
}
</script>
