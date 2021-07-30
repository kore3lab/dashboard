<template>
<div>
	<!-- 1. metadata -->
	<c-metadata v-model="metadata" dtCols="3" ddCols="9"  @navigate="$emit('navigate', arguments[0])">
		<dt class="col-sm-3">Status</dt><dd class="col-sm-9" v-bind:class="{ 'text-success': status.phase=='Active' }">{{ status.phase }}</dd>
		<dt class="col-sm-3">Resource Quotas</dt>
		<dd class="col-sm-9"><span v-if="quotas.length==0">-</span><span v-for="(val, idx) in quotas" v-bind:key="idx" class="mr-1"><a href="#" @click="$emit('navigate', getViewLink('', 'resourcequotas', metadata.name,val))">{{ val }} </a></span></dd>
		<dt class="col-sm-3">Limit Ranges</dt>
		<dd class="col-sm-9"><span v-if="limits.length==0">-</span><span v-for="(val, idx) in limits" v-bind:key="idx" class="mr-1"><a href="#" @click="$emit('navigate', getViewLink('', 'limitranges', metadata.name,val))">{{ val }} </a></span></dd>
	</c-metadata>
</div>
</template>
<script>
import VueMetadataView	from "@/components/view/metadataView.vue";

export default {
	components: {
		"c-metadata": { extends: VueMetadataView }
	},
	data() {
		return {
			metadata: {},
			status: {},
			quotas: [],
			limits: []
		}
	},
	mounted() {
		this.$nuxt.$on("onReadCompleted", (data) => {
			if(!data) return
			this.metadata = data.metadata;
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

		});
		this.$nuxt.$emit("onCreated",'')
	},
	methods: {},
	beforeDestroy(){
		this.$nuxt.$off("onReadCompleted");
	},
}
</script>
