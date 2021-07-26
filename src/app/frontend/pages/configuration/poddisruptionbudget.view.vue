<template>
<div>
	<!-- 1. metadata -->
	<c-metadata v-model="metadata" dtCols="3" ddCols="9">
		<dt v-if="metadata.ownerReferences" class="col-sm-3">Controlled By</dt>
		<dd v-if="metadata.ownerReferences" class="col-sm-9">{{ controller.kind }} <a href="#" @click="$emit('navigate', getViewLink(controller.group, controller.resource, metadata.namespace, controller.name))">{{ controller.name }}</a></dd>
		<dt v-if="selector" class="col-sm-3">Selector</dt>
		<dd v-if="selector" class="col-sm-9">
			<span v-for="(value, name) in selector" v-bind:key="name" class="border-box background">{{ value }}</span>
		</dd>
		<dt class="col-sm-3">Min Available</dt><dd class="col-sm-9">{{ info.minAvailable }}</dd>
		<dt class="col-sm-3">max Unavailable</dt><dd class="col-sm-9">{{ info.maxUnavailable }}</dd>
		<dt class="col-sm-3">Current Healthy</dt><dd class="col-sm-9">{{ info.currentHealthy }}</dd>
		<dt class="col-sm-3">Desired Healthy</dt><dd class="col-sm-9">{{ info.desiredHealthy }}</dd>
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
			info: {},
			selector: {},
			controller: {}
		}
	},
	mounted() {
		this.$nuxt.$on("onReadCompleted", (data) => {
			if(!data) return
			this.metadata = data.metadata;
			this.info =  {
				minAvailable: data.spec.minAvailable || "N/A",
				maxUnavailable: data.spec.maxUnavailable || "N/A",
				currentHealthy: data.status.currentHealthy,
				desiredHealthy: data.status.desiredHealthy,
			}
			this.controller = data.metadata.ownerReferences? this.getResource(data.metadata.ownerReferences[0]): {};
			this.selector = this.stringifyLabels(data.spec.selector.matchLabels);
		});
		this.$nuxt.$emit("onCreated",'')
	},
	beforeDestroy(){
		this.$nuxt.$off("onReadCompleted");
	},
}
</script>
