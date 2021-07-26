<template>
<div>
	<!-- 1. metadata -->
	<c-metadata v-model="metadata" dtCols="3" ddCols="9">
		<dt class="col-sm-3">Reference</dt>
		<dd class="col-sm-9">{{ ref.kind }} / <a href="#" @click="$emit('navigate', getViewLink(ref.group, ref.resource, metadata.namespace, ref.name))">{{ ref.name }}</a></dd>
		<dt class="col-sm-3">Min Pods</dt><dd class="col-sm-9">{{ info.minPods }}</dd>
		<dt class="col-sm-3">Max Pods</dt><dd class="col-sm-9">{{ info.maxPods }}</dd>
		<dt class="col-sm-3">Replicas</dt><dd class="col-sm-9">{{ info.replicas }}</dd>
	</c-metadata>
	<!-- 2. metrics-->
	<div class="row" v-show="metrics">
		<div class="col-md-12">
			<div class="card card-secondary card-outline">
				<div class="card-header p-2"><h3 class="card-title">Metrics</h3></div>
				<div class="card-body group">
					<b-table-lite :items="metrics" :fields="fields" class="subset">
					</b-table-lite>
				</div>
			</div>
		</div>
	</div>
	<!-- 3. evnets -->
	<c-events class="row" v-model="metadata.uid"></c-events>
</div>
</template>
<script>
import VueMetadataView	from "@/components/view/metadataView.vue";
import VueEventsView		from "@/components/view/eventsView.vue";
import VueJsonTree			from "@/components/jsontree";

export default {
	components: {
		"c-metadata": { extends: VueMetadataView },
		"c-jsontree": { extends: VueJsonTree },
		"c-events": { extends: VueEventsView }
	},
	data() {
		return {
			metadata: {},
			info: [],
			ref: {},
			fields: [
				{ key: "name", label: "Name" },
				{ key: "metric", label: "Current/Target" },
			],
			metrics: []
		}
	},
	mounted() {
		this.$nuxt.$on("onReadCompleted", (data) => {
			if(!data) return
			this.metadata = data.metadata;
			this.info = {
				minPods: data.spec.minReplicas || 0,
				maxPods: data.spec.maxReplicas || 0,
				replicas: data.status.currentReplicas
			},
			this.metrics = [
				{name: "Resource cpu on Pods (as a percentage of request)", metric: `${data.status.currentCPUUtilizationPercentage ? data.status.currentCPUUtilizationPercentage + '%': "<unknown>"} / ${data.spec.targetCPUUtilizationPercentage?data.spec.targetCPUUtilizationPercentage:"0"}%`}
			];
			this.ref = this.getResource(data.spec.scaleTargetRef)
		});
		this.$nuxt.$emit("onCreated",'')
	},
	methods: {},
	beforeDestroy(){
		this.$nuxt.$off("onReadCompleted");
	},
}
</script>
