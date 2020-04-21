<template>
<!-- content-wrapper -->
<div class="content-wrapper">

	<!-- Content Header (Page header) -->
	<div class="content-header">
		<div class="container-fluid">
			<c-navigator group="Cluster"></c-navigator>
			<div class="row mb-2">
				<div class="col-sm-2">
					<h1 class="m-0 text-dark">Topology</h1>
				</div>
				<div class="col-sm-2">
					<b-form-select v-model="selectedNamespace" :options="namespaces()" size="sm" @input="query"></b-form-select>
				</div>
				<!-- KCL -->
				<div class="col-sm-3 text-lg">
					nodes : <span class="text-danger">{{ nodeLen }}</span> <span class="text-sm">ea</span>, elapsed : <span class="text-danger">{{ elapsed }}</span> <span class="text-sm">ms</span>
				</div>
				<div class="col-sm-3 text-lg">
					node/sec : <span class="text-danger">{{ Math.round(nodeLen/(elapsed/1000)*100)/100 }}</span> <span class="text-sm">ea</span>
				</div>
				<!-- KCL -->
				<div class="col-sm-2 text-right">
					<b-button variant="primary" size="sm" @click="$nuxt.$emit('navbar-context-selected',currentContext)">Reload</b-button>
				</div>
			</div>
		</div>
	</div>
	<!-- Main content -->
	<section class="content">
		<div class="container-fluid">
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
<!-- /.content-wrapper -->
</template>
<script>
import * as graph	from "../../static/acorn-graph/acorn.graph.topology"
import axios		from "axios"
import VueNavigator from "@/components/navigator"

export default {
	components: {
		"c-navigator": { extends: VueNavigator }
	},
	data() {
		let ns = this.$route.query.namespace;
		return {
			selectedNamespace: ns ? ns: " ",
			nodeLen: 0,
			start: 0,		//KCL
			elapsed: 0		//KCL
		}
	},
	layout: "default",
	created() {
		this.$nuxt.$on("navbar-context-selected", () => this.query());
		if(this.currentContext()) this.$nuxt.$emit("navbar-context-selected");
	},
	beforeDestroy(){
		this.$nuxt.$off("navbar-context-selected")
	},
	methods: {
		onEnd() {
			this.$data.elapsed = (new Date()).getTime() - this.$data.start; // KCL
		},
		query() {
			this.$data.start = (new Date()).getTime();
			let url = `${this.backendUrl()}/api/clusters/${this.currentContext()}/topology`;
			if (this.$data.selectedNamespace != " ") url += `/namespaces/${this.$data.selectedNamespace}`;
			

			let g = new graph.TopologyGraph("#wrapGraph");
			axios.get(url)
				.then( resp => {
					// KCL
					this.$data.nodeLen = resp.data.nodes.length;
					g.config({
						topology:{
							simulation: {
								alphaDecay:0.3,
								onEnd: this.onEnd
							}
						}
					}).data(resp.data).render();
					//-- KCL
					// g.data(resp.data).render();

				})
				.catch((error) => {
					this.toast(error.message);
					console.log(error);
				});
		}
	}
}
</script>
