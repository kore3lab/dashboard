<template>
	<div class="content-wrapper">
		<section class="content-header">
			<div class="container-fluid">
				<c-navigator group="Configuration"></c-navigator>
				<div class="row mb-2">
					<!-- title & search -->
					<div class="col-sm"><h1 class="m-0 text-dark"><span class="badge badge-info mr-2">H</span>Horizontal Pod Autoscalers</h1></div>
					<div class="col-sm-2"><b-form-select v-model="selectedNamespace" :options="namespaces()" size="sm" @input="query_All(); selectNamespace(selectedNamespace);"></b-form-select></div>
					<div class="col-sm-2 float-left">
						<div class="input-group input-group-sm" >
							<b-form-input id="txtKeyword" v-model="keyword" class="form-control float-right" placeholder="Search"></b-form-input>
							<div class="input-group-append"><button type="submit" class="btn btn-default" @click="query_All"><i class="fas fa-search"></i></button></div>
						</div>
					</div>
					<!-- button -->
					<div class="col-sm-1 text-right">
						<b-button variant="primary" size="sm" @click="$router.push(`/create?context=${currentContext()}&group=Configuration&crd=HPA`)">Create</b-button>
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
									<template v-slot:cell(target)="data">
										<span>{{ data.value.kind }} / <a href="#" @click="viewModel=getViewLink(data.value.group,data.value.resource,data.item.namespace, data.value.name); isShowSidebar=true;">{{ data.value.name }}</a></span>
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
				{ key: "metrics", label: "Metrics", sortable: true  },
				{ key: "minpods", label: "Min Pods", sortable: true },
				{ key: "maxpods", label: "Max Pods", sortable: true },
				{ key: "replicas", label: "Replicas", sortable: true },
				{ key: "target", label: "Target", sortable: true, formatter: this.formatTarget },
				{ key: "creationTimestamp", label: "Age", sortable: true, formatter: this.getElapsedTime },
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
	created() {
		this.$nuxt.$on("navbar-context-selected", (ctx) => {
			this.selectedNamespace = this.selectNamespace()
			this.query_All()
		});
		if(this.currentContext()) this.$nuxt.$emit("navbar-context-selected");
	},
	methods: {
		onRowSelected(items) {
			this.isShowSidebar = (items && items.length > 0)
			if (this.isShowSidebar) this.viewModel = this.getViewLink('autoscaling', 'horizontalpodautoscalers', items[0].namespace, items[0].name)
		},
		// 조회
		query_All() {
			this.isBusy = true;
			this.$axios.get(this.getApiUrl("autoscaling","horizontalpodautoscalers",this.selectedNamespace))
					.then((resp) => {
						this.items = [];
						resp.data.items.forEach(el => {
							this.items.push({
								name: el.metadata.name,
								namespace: el.metadata.namespace,
								metrics: `${el.status.currentCPUUtilizationPercentage ? el.status.currentCPUUtilizationPercentage +'%': "<unknown>"} / ${el.spec.targetCPUUtilizationPercentage?el.spec.targetCPUUtilizationPercentage:"0"}%`,
								minpods: el.spec.minReplicas,
								maxpods: el.spec.maxReplicas,
								replicas: el.status.currentReplicas,
								target: el.spec.scaleTargetRef,
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
		formatTarget(ref) {
			let group = ""
			let version = ref.apiVersion.split('/')
			if (version.length>1) {
				group = version[0]
			}
			return {
				"group": group,
				"kind": ref.kind,
				"name": ref.name,
				"resource": ref.kind.toLowerCase()+'s'
			}
		},
	},
	beforeDestroy(){
		this.$nuxt.$off('navbar-context-selected')
	}
}
</script>