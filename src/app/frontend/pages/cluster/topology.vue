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
					<div class="col-12">
						<div class="card">
							<div class="card-body m-0 p-2">
								<div id="wrapGraph" class="p-0" style="min-height: calc(100vh - 210px - 60px)"></div>
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
			viewModel:{}
		}
	},
	layout: "default",
	methods: {
		query(d) {
			this.$data.start = (new Date()).getTime();
			const ns = (d && d.namespace) ? d.namespace: this.selectNamespace();

			let g = new TopologyGraph("#wrapGraph", {
				global: {
					scale: {ratio:1},
					toolbar: {
						align: { horizontal: "left" }
					}
				},
				extends: {
					topology:{
						simulation: {
							alphaDecay:0.3
						}
					}
				}
			})
			.on("nodeclick", (e,d)=> { 
				if (d.kind == "pod" && d.namespace)  {
					this.isShowSidebar = true;
					if(this.isShowSidebar) this.viewModel = this.getViewLink("", "pods", d.namespace, d.name);
				} else if (d.kind == "node")  {
					this.isShowSidebar = true;
					if(this.isShowSidebar) this.viewModel = this.getViewLink("", "nodes", "", d.name);
				} else {
					this.isShowSidebar = false;
				}
			})

			let url = `/api/clusters/${this.currentContext()}/graph/topology`;
			if (ns) url += `/namespaces/${this.selectNamespace()}`;
			this.$axios.get(url)
				.then( resp => {
					this.$data.nodeLen = resp.data.nodes.length;
					g.data(resp.data).render();
				})
				.catch(e => { this.msghttp(e);})
		}
	}
}
</script>
