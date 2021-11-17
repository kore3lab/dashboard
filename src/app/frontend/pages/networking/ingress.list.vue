<template>
  <div class="content-wrapper">
		<section class="content-header">
			<div class="container-fluid">
				<c-navigator group="Networking"></c-navigator>
				<div class="row mb-2">
					<div class="col-sm"><h1 class="m-0 text-dark"><span class="badge badge-info mr-2">I</span>Ingresses</h1></div>
					<div class="col-sm-1 text-right">
						<b-button variant="primary" size="sm" @click="$router.push(`/create?context=${currentContext()}&group=Networking&crd=Ingress`)">Create</b-button>
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
									<template v-slot:cell(loadBalancers)="data">
										<ul class="list-unstyled mb-0">
											<li v-for="(value, idx) in data.value" v-bind:key="idx">{{ value.hostname || value.ip }} </li>
										</ul>
									</template>
									<template v-slot:cell(rules)="data">
										<ul class="list-unstyled mb-0">
											<li v-for="(value, idx) in data.value" v-bind:key="idx"><b-link @click="onLink(value.fr)">{{ value.fr }}</b-link> ⇢ {{ value.to }}</li>
										</ul>
									</template>
									<template v-slot:cell(endpoints)="data">
										<ul class="list-unstyled mb-0">
											<li v-for="(d, idx) in data.item.endpoints" v-bind:key="idx">{{ d }}</li>
										</ul>
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
import VueView				from "@/pages/view"

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
				{key: "name", label: "Name", sortable: true},
				{key: "namespace", label: "Namespace", sortable: true},
				{key: "loadBalancers", label: "LoadBalancers"},
				{key: "rules", label: "Rules", formatter: this.formatRules },
				{key: "creationTimestamp", label: "Age", sortable: true, formatter: this.getElapsedTime},
			],
			isBusy: false,
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
	methods: {
		onLink(val) {
			window.open(val,'')
		},
		onRowSelected(items) {
			this.isShowSidebar = (items && items.length > 0)
			if (this.isShowSidebar) this.viewModel = this.getViewLink('networking.k8s.io', 'ingresses', items[0].namespace, items[0].name)
		},
		// 조회
		query_All(d) {
			this.isBusy = true;
			this.$axios.get(this.getApiUrl("networking.k8s.io","ingresses",(d && d.namespace) ? d.namespace: this.selectNamespace(), "", d && d.labelSelector? `labelSelector=${d.labelSelector}`: ""))
				.then((resp) => {
					this.items = [];
					resp.data.items.forEach(el => {
						this.items.push({
							name: el.metadata.name,
							namespace: el.metadata.namespace,
							loadBalancers: el.status.loadBalancer || el.status.loadBalancer.ingress ? el.status.loadBalancer.ingress: [],
							rules: el.spec.rules,
							tls: el.spec.tls,
							creationTimestamp: el.metadata.creationTimestamp,
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
		formatRules(value, key, item) {
			let rules = item.rules
			let tls = item.tls
			if(!rules) return [];

			let protocol = 'http';
			const routes = [];

			if (tls && tls.length > 0) {
				protocol += "s";
			}
			rules.map(rule => {
				const host = rule.host ? rule.host : "*";

				if (rule.http && rule.http.paths) {
					rule.http.paths.forEach(path => {
						const serviceName = "service" in path.backend ? path.backend.service.name : path.backend.serviceName;
						const servicePort = "service" in path.backend ? path.backend.service.port.number ?? path.backend.service.port.name : path.backend.servicePort;
						routes.push({
								fr: `${protocol}://${host}${path.path || "/"}`,
								to: `${serviceName}:${servicePort}`
						})
					});
				}
			});
			return routes;
		}
	}
}
</script>
