<template>
<div>
	<!-- 1. charts -->
	<c-charts v-model="value" class="row"></c-charts>
	<!-- 2. metadata -->
	<c-metadata v-model="value" dtCols="2" ddCols="10" @navigate="$emit('navigate', arguments[0])">
		<dt class="col-sm-2">Replicas</dt><dd class="col-sm-10">{{ info.replicas }}</dd>
		<dt class="col-sm-2">Strategy Type</dt><dd class="col-sm-10">{{ info.strategyType }}</dd>
		<dt class="col-sm-2">Conditions</dt>
		<dd class="col-sm-10">
			<span v-for="(d, idx) in info.conditions" v-bind:key="idx" v-bind:class="{'badge-primary':d.type=='Progressing', 'badge-success':d.type=='Available'}" class="badge font-weight-light text-sm mb-1 mr-1">{{ d.type }}</span>
		</dd>
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
			spec: {},
			info: {}
		}
	},
	watch: {
		value(d) { this.onSync(d) }
	},
	methods: {
		onSync(data) {
			if(!data) return
			this.spec = data.spec;
			this.info = {
				replicas: `${data.spec.replicas} desired, ${data.status.updatedReplicas || 0} updated, ${data.status.replicas || 0} total, ${data.status.availableReplicas || 0} available, ${data.status.unavailableReplicas || 0} unavailable`,
				strategyType: data.spec.strategy.type,
				conditions: data.status.conditions? data.status.conditions.sort( (a,b)=> {return a.type < b.type ? -1 : a.type > b.type ? 1 : 0;} ): []
			}
		}
	}
}
</script>
