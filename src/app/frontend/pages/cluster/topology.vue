<template>
	<div class="content-wrapper">
		<div class="content-header">
			<div class="container-fluid">
				<c-navigator group="Cluster"></c-navigator>
				<div class="row mb-2">
					<!-- title & search -->
					<div class="col-sm"><h1 class="m-0 text-dark"><span class="badge badge-info mr-2">T</span>Topology</h1></div>
					<div class="col-sm-2"><b-form-select v-model="selectedNamespace" :options="namespaces()" size="sm" @input="query(); selectNamespace(selectedNamespace);"></b-form-select></div>
					<div class="col-sm-1 text-right"><b-button variant="primary" size="sm" @click="$nuxt.$emit('navbar-context-selected',currentContext)">Reload</b-button></div>
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
</template>
<script>
import * as graph	from "../../static/kore3lab.graph/kore3lab.graph.topology"
import VueNavigator from "@/components/navigator"

export default {
	components: {
		"c-navigator": { extends: VueNavigator }
	},
	data() {
		let ns = this.$route.query.namespace;
		return {
			selectedNamespace: ns ? ns: "",
			nodeLen: 0
		}
	},
	layout: "default",
	created() {
		this.$nuxt.$on("navbar-context-selected", (ctx) => {
			this.selectedNamespace = this.selectNamespace()
			this.query()
		});
		if(this.currentContext()) this.$nuxt.$emit("navbar-context-selected");
	},
	beforeDestroy(){
		this.$nuxt.$off("navbar-context-selected")
	},
	methods: {
		query() {
			this.$data.start = (new Date()).getTime();
			let url = `/api/clusters/${this.currentContext()}/topology`;
			if (this.$data.selectedNamespace !== "") url += `/namespaces/${this.$data.selectedNamespace}`;

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
