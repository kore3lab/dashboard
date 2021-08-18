<template>
<div>
	<!-- 1. metadata -->
	<c-metadata dtCols="2" ddCols="10"  @navigate="$emit('navigate', arguments[0])">
		<dt class="col-sm-2">Conditions</dt>
		<dd class="col-sm-10">
			<span v-for="(value, idx) in info.conditions" v-bind:key="idx" v-bind:class="{'badge-success':value.type=='Complete', 'badge-danger':value.type=='Failed'}" class="badge font-weight-light text-sm mb-1 mr-1"> {{ value.type }} </span>
		</dd>
		<dt class="col-sm-2">Completions</dt><dd class="col-sm-10">{{ info.completions }}</dd>
		<dt class="col-sm-2">Parallelism</dt><dd class="col-sm-10">{{ info.parallelism }}</dd>
	</c-metadata>
	<!-- 2. pods -->
	<c-podlist class="row" @navigate="$emit('navigate',arguments[0])"></c-podlist>
	<!-- 3. events -->
	<c-events class="row"></c-events>
</div>
</template>
<script>
import VueMetadataView	from "@/components/view/metadataView.vue";
import VueEventsView	from "@/components/view/eventsView.vue";
import VuePodListView	from "@/components/view/podListView.vue";

export default {
	components: {
		"c-metadata": { extends: VueMetadataView },
		"c-events": { extends: VueEventsView },
		"c-podlist": { extends: VuePodListView }
	},
	data() {
		return {
			info: {}
		}
	},
	mounted() {
		this.$nuxt.$on("view-data-read-completed", (data) => {
			if(!data) return
			this.info = {
				conditions: data.status.conditions?data.status.conditions.filter(el=>{return el.status === 'True'}):[],
				completions: data.spec.completions,
				parallelism: data.spec.parallelism,
			};
		});
	},
	beforeDestroy(){
		this.$nuxt.$off("view-data-read-completed");
	}
}
</script>
