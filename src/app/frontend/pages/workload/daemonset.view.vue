<template>
<div>
	<!-- 1. chart -->
	<c-charts class="row" v-model="selectUrl"></c-charts>
	<!-- 2. metadata -->
	<c-metadata v-model="metadata" :workload="spec" dtCols="2" ddCols="10" @navigate="$emit('navigate', arguments[0])">
		<dt class="col-sm-2">Strategy Type</dt><dd class="col-sm-10">{{ strategyType }}</dd>
	</c-metadata>
	<!-- 3. pods -->
	<c-podlist class="row" v-model="selectUrl" @navigate="$emit('navigate',arguments[0])"></c-podlist>
	<!-- 4. events -->
	<c-events class="row" v-model="metadata.uid"></c-events>
</div>
</template>
<script>
import VueMetadataView	from "@/components/view/metadataView.vue";
import VueEventsView	from "@/components/view/eventsView.vue";;
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
			metadata: {},
			selectUrl: "",
			spec: {},
			strategyType: ""
		}
	},
	mounted() {
		this.$nuxt.$on("onReadCompleted", (data) => {
			if(!data) return
			this.metadata = data.metadata;
			this.selectUrl = `namespaces/${data.metadata.namespace}/daemonsets/${data.metadata.name}`;
			this.strategyType = data.spec.updateStrategy.type;
		});
		this.$nuxt.$emit("onCreated",'')
	},
	methods: {},
	beforeDestroy(){
		this.$nuxt.$off("onReadCompleted");
	},
}
</script>

