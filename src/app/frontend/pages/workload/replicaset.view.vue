<template>
<div>
	<!-- 1. charts -->
	<c-charts class="row"></c-charts>
	<!-- 2. metadata -->
	<c-metadata dtCols="2" ddCols="10" @navigate="$emit('navigate', arguments[0])">
		<dt class="col-sm-2">Replicas</dt>
		<dd class="col-sm-10">{{ replicas }}</dd>
	</c-metadata>
	<!-- 3. pods -->
	<c-podlist class="row" @navigate="$emit('navigate',arguments[0])"></c-podlist>
	<!-- 4. events -->
	<c-events class="row"></c-events>
</div>
</template>
<script>
import VueMetadataView	from "@/components/view/metadataView.vue";
import VueEventsView	from "@/components/view/eventsView.vue";
import VueChartsView	from "@/components/view/metricsChartsView.vue";
import VuePodListView	from "@/components/view/podListView.vue";

export default {
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
	mounted() {
		this.$nuxt.$on("view-data-read-completed", (data) => {
			if(!data) return
			this.replicas = `${data.status.availableReplicas || 0} current / ${data.status.replicas || 0} desired`;
		});
	},
	beforeDestroy(){
		this.$nuxt.$off("view-data-read-completed");
	}
}
</script>
