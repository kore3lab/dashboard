<template>
<div>
	<!-- 1. charts -->
	<c-charts v-model="value" class="row"></c-charts>
	<!-- 2. metadata -->
	<c-metadata v-model="value" dtCols="2" ddCols="10" @navigate="$emit('navigate', arguments[0])">
		<dt class="col-sm-2">Replicas</dt>
		<dd class="col-sm-10">{{ replicas }}</dd>
	</c-metadata>
	<!-- 3. pods -->
	<c-podlist v-model="value" class="row" @navigate="$emit('navigate',arguments[0])"></c-podlist>
	<!-- 4. events -->
	<c-events v-model="value" class="row"></c-events>
</div>
</template>
<script>
import VueMetadataView	from "@/components/view/metadataView.vue";
import VueEventsView	from "@/components/view/eventsView.vue";
import VueChartsView	from "@/components/view/metricsChartsView.vue";
import VuePodListView	from "@/components/view/podListView.vue";

export default {
	props:["value"],
	components: {
		"c-metadata": { extends: VueMetadataView },
		"c-events": { extends: VueEventsView },
		"c-charts": { extends: VueChartsView },
		"c-podlist": { extends: VuePodListView }
	},
	data() {
		return {
			replicas: ""
		}
	},
	watch: {
		value(d) { this.onSync(d) }
	},
	methods: {
		onSync(data) {
			if(!data) return
			this.replicas = `${data.status.availableReplicas || 0} current / ${data.status.replicas || 0} desired`;
		}
	}
}
</script>
