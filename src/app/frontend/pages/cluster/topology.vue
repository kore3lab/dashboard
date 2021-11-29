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
	</div>
</template>
<script>
import * as graph		from "../../static/kore3lab.graph/kore3lab.graph.topology"
import VueNavigator		from "@/components/navigator"
import VueSearchForm	from "@/components/list/searchForm"

export default {
	components: {
		"c-navigator": { extends: VueNavigator },
		"c-search-form": { extends: VueSearchForm}
	},
	data() {
		return {
		}
	},
	layout: "default",
	methods: {
		query(d) {
			this.$data.start = (new Date()).getTime();
			const ns = (d && d.namespace) ? d.namespace: this.selectNamespace();

			let url = `/api/clusters/${this.currentContext()}/topology`;
			if (ns) url += `/namespaces/${this.selectNamespace()}`;

			let g = new graph.TopologyGraph("#wrapGraph");
			this.$axios.get(url)
				.then( resp => {
					this.$data.nodeLen = resp.data.nodes.length;
					g.config({
						topology:{
							simulation: {
								alphaDecay:0.3
							}
						}
					}).data(resp.data).render();
				})
				.catch(e => { this.msghttp(e);})
		}
	}
}
</script>
