<template>
<div class="content-wrapper">

	<div class="content-header">
		<div class="container-fluid">
			<c-navigator group="Cluster"></c-navigator>
			<div class="row mb-2">
				<div class="col-sm-10">
					<h1 class="m-0 text-dark">ServiceMesh</h1>
				</div>
				<div class="col-sm-2 text-right">
					<b-button variant="primary" size="sm" @click="$nuxt.$emit('navbar-context-selected',currentContext)">Reload</b-button>
				</div>
			</div>
		</div>
	</div>

	<section class="content">
	<div class="container-fluid">
		<div class="row">
			<div class="col-12">
				<div class="card">
					<iframe :src="source" class="embed-responsive-item card-body m-0 p-2 border-0" style="min-height: calc((100vh - 210px) - 60px);"></iframe>
				</div>
			</div>
		</div>
		
	</div>
	</section>
</div>
</template>

<script>
import VueNavigator from "@/components/navigator"

export default {
	components: {
		"c-navigator": { extends: VueNavigator }
	},
	data() {
		return {
			selectedNamespaces: [" "],
			source: ""
		}
	},
	layout: "default",
	created() {
		this.$nuxt.$on("navbar-context-selected", () => this.query_All());
		if(this.currentContext()) this.$nuxt.$emit("navbar-context-selected");
	},
	beforeDestroy(){
		this.$nuxt.$off("navbar-context-selected")
	},
	methods: {
		query_All() {
			let namespaces = ""
			this.namespaces().forEach(el => {
				namespaces += el.value + ",";
			});
			this.$data.source = `${this.kialiRootUrl()}/kiali/console/graph/namespaces/?edges=noEdgeLabels&graphType=versionedApp&unusedNodes=true&operationNodes=false&injectServiceNodes=true&duration=60&refresh=15000&namespaces=${namespaces}&layout=dagre&kiosk=true`;
		}
	}
}
</script>
