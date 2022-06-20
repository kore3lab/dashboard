<template>
	<div class="content-wrapper">
		<section class="content-header">
			<div class="container-fluid">
				<c-navigator group="Cluster"></c-navigator>
				<div class="row mb-2">
					<!-- title & search -->
					<div class="col-sm"><h1 class="m-0 text-dark"><span class="badge badge-info mr-2">T</span>Topology</h1></div>
					<div class="col-sm-1 text-right"><b-button variant="light" @click="query"><b-icon-arrow-repeat></b-icon-arrow-repeat></b-button></div>
				</div>
			</div>
		</section>
		<!-- Main content -->
		<section class="content">
			<div class="container-fluid">
				<!-- search  -->
				<c-search-form class="mb-2" no-label-selector no-keyword @input="query" />
				<!-- graph  -->
				<div class="row">
					<b-nav tabs class="col-12">
						<b-nav-item class="p-0" link-classes="pl-3 pr-3" :active="selctedTab==1" @click="selctedTab=1;query()">Topology</b-nav-item>
						<b-nav-item class="p-0" link-classes="pl-3 pr-3" :active="selctedTab==2" @click="selctedTab=2;query()">Workloads</b-nav-item>
					</b-nav>
				</div>
				<div class="row">
					<div class="col-12">
						<div class="card">
							<div class="card-body m-0 p-2">
								<b-overlay :show="isShowOverlay" rounded="sm" spinner-variant="primary">
								<div id="wrapTopologyGraph" class="p-0" style="min-height: calc(100vh - 210px - 60px)"></div>
								</b-overlay>
							</div>
						</div>
					</div>
				</div>
			</div>
		</section>
		<b-sidebar v-model="isShowSidebar" width="50em" right shadow no-header>
			<c-view v-model="viewModel"  @close="isShowSidebar=false"/>
		</b-sidebar>
	</div>
</template>
<script>
import VueNavigator		from "@/components/navigator"
import VueSearchForm	from "@/components/list/searchForm"
import TopologyGraph	from "@/components/graph/graph.topology"
import HierarchyGraph	from "@/components/graph/graph.hierarchy";
import VueView			from "@/pages/view";

export default {
	components: {
		"c-navigator": { extends: VueNavigator },
		"c-search-form": { extends: VueSearchForm},
		"c-view": { extends: VueView }
	},
	data() {
		return {
			isShowSidebar: false,
			isShowOverlay: false,
			viewModel:{},
			selctedTab: 1
		}
	},
	layout: "default",
	methods: {
		query(d) {
			const ns = (d && d.namespace) ? d.namespace: this.selectNamespace();
			const g = this.selctedTab == 1 ? new TopologyGraph("#wrapTopologyGraph"): new HierarchyGraph("#wrapTopologyGraph");
			// node-click()
			g.on("nodeclick", (e,data)=> { 
				const d = data.data ? data.data: data;
				if ("Container,Cluster".includes(d.kind)) {
					this.isShowSidebar = false;
				} else {
					const group = ("DaemonSet,ReplicaSet,StatefulSet,Deployment".includes(d.kind)) ? "apps": "";
					this.isShowSidebar = true;
					if(this.isShowSidebar) this.viewModel = this.getViewLink(group, `${d.kind.toLowerCase()}s`, d.namespace, d.name);
				}
			});

			this.isShowOverlay = true;
			let url = `/api/clusters/${this.currentContext()}/graph/${this.selctedTab==1?"topology":"workloads"}`;
			if (ns) url += `/namespaces/${this.selectNamespace()}`;
			this.$axios.get(url)
				.then( resp => {
					g.config({
						global: {
							scale: {ratio:1},
							toolbar: {
								align: { horizontal: "right" }
							}
						},
						extends: {
							topology:{
								simulation: {
									alphaDecay:0.3
								}
							},
							hierarchy: {
								group: {
									divide: ns? true: false,
									title: { display: "none" },
									box: {
										border: { width: 0 },
										background: { fill: "none" }
									}
								},
								node: {
									forEach: ns? (d)=> { if(d.kind == "Pod") d.depth = 2; }: (d)=> { if(d.kind == "Pod") d.depth = 3; }
								}
							}
						}
					}).data(resp.data).render();
				})
				.catch(e => { this.msghttp(e);})
				.finally(()=>{this.isShowOverlay = false;})
			

		}
	}
}
</script>
