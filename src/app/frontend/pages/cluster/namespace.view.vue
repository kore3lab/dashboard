<template>
<div>
	<!-- 1. metadata -->
	<c-metadata v-model="value" dtCols="3" ddCols="9"  @navigate="$emit('navigate', arguments[0])">
		<dt class="col-sm-3">Status</dt><dd class="col-sm-9" v-bind:class="{ 'text-success': status.phase=='Active' }">{{ status.phase }}</dd>
		<dt class="col-sm-3">Resource Quotas</dt>
		<dd class="col-sm-9"><span v-if="quotas.length==0">-</span><span v-for="(val, idx) in quotas" v-bind:key="idx" class="mr-1"><a href="#" @click="$emit('navigate', getViewLink('', 'resourcequotas', metadata.name,val))">{{ val }} </a></span></dd>
		<dt class="col-sm-3">Limit Ranges</dt>
		<dd class="col-sm-9"><span v-if="limits.length==0">-</span><span v-for="(val, idx) in limits" v-bind:key="idx" class="mr-1"><a href="#" @click="$emit('navigate', getViewLink('', 'limitranges', metadata.name,val))">{{ val }} </a></span></dd>
	</c-metadata>
	<!-- 2. graph -->
	<div class="row" v-show="isWorkload">
		<div class="col-md-12">
			<div class="card card-secondary card-outline m-0">
				<div class="card-header p-2"><h3 class="card-title">Workloads</h3></div>
				<div class="card-body mw-100" id="wrapWorkloadsGraph"  v-show="isWorkload"></div>
			</div>
		</div>
	</div>

</div>
</template>
<style scoped>
#wrapWorkloadsGraph {min-height: 40em;}
</style>
<script>
import VueMetadataView	from "@/components/view/metadataView.vue";
import HierarchyGraph	from "@/components/graph/graph.hierarchy";

export default {
	props:["value"],
	components: {
		"c-metadata": { extends: VueMetadataView }
	},
	data() {
		return {
			status: {},
			quotas: [],
			limits: [],
			isWorkload: true
		}
	},
	watch: {
		value(d) { this.onSync(d) }
	},
	methods: {
		onSync(data) {
			if(!data) return
			this.status = data.status;
			this.quotas = [];
			this.$axios.get(this.getApiUrl("","resourcequotas",data.metadata.name))
			.then(resp => {
				resp.data.items.forEach(el =>{
					this.quotas.push(el.metadata.name)
				})
			});

			this.limits = [];
			this.$axios.get(this.getApiUrl("","limitranges",data.metadata.name))
			.then(resp => {
				resp.data.items.forEach(el =>{
					this.limits.push(el.metadata.name)
				})
			});
			//workloads graph
			let g = new HierarchyGraph("#wrapWorkloadsGraph", {
				global: {
					toolbar: { visible:false }
				},
				extends: {
					hierarchy: {
						group: {
							title: { display: "none" }
						},
						node: {
							forEach: (d)=> { if(d.kind == "Pod") d.depth = 2; }
						}
					}
				}
			})
			.on("nodeclick", (e,d)=> {
				if (d.data.namespace && d.data.name) {
					// (core) pod, (apps) daemonset, replicaset, deployment
					const model = this.getViewLink(d.data.kind=="Pod"?"":"apps", `${d.data.kind.toLowerCase()}s`, d.data.namespace, d.data.name);
					this.$emit("navigate", model)
				}
			})
			this.isWorkload = true;
			this.$axios.get(`/api/clusters/${this.currentContext()}/graph/workloads/namespaces/${data.metadata.name}`)
				.then( resp => {
					this.isWorkload = resp.data[data.metadata.name].length > 0;
					g.data(resp.data).render();
				})
				.catch(e => { this.msghttp(e);})

		}
	}
}
</script>
