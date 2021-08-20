<template>
<div>
	<!-- 1. metadata -->
	<c-metadata dtCols="2" ddCols="10"></c-metadata>
	<!-- 2. rules -->
	<div class="row">
		<div class="col-md-12">
			<div class="card card-secondary card-outline">
				<div class="card-header p-2"><h3 class="card-title">Rules</h3></div>
				<div class="card-body group">
					<ul>
						<li v-for="(val, idx) in rules" v-bind:key="idx">
							<dl class="row">
								<dt v-if="val.resources" class="col-sm-2">Resources</dt><dd v-if="val.resources" class="col-sm-9">{{ val.resources.join(", ") }}</dd>
								<dt v-if="val.verbs" class="col-sm-2">Verbs</dt><dd v-if="val.verbs" class="col-sm-9">{{ val.verbs.join(", ") }}</dd>
								<dt v-if="val.apiGroups" class="col-sm-2">Api Groups</dt><dd v-if="val.apiGroups" class="col-sm-9">{{ val.apiGroups.join(", ") }}</dd>
								<dt v-if="val.resourceNames" class="col-sm-2">Resource Names</dt><dd v-if="val.resourceNames" class="col-sm-9">{{ val.resourceNames.join(", ") }}</dd>
							</dl>
						</li>
					</ul>
				</div>
			</div>
		</div>
	</div>
	<!-- 3. events -->
	<c-events class="row"></c-events>

</div>
</template>
<script>
import VueMetadataView	from "@/components/view/metadataView.vue";
import VueEventsView	from "@/components/view/eventsView.vue";

export default {
	components: {
		"c-metadata": { extends: VueMetadataView },
		"c-events": { extends: VueEventsView }
	},
	data() {
		return {
			rules: []
		}
	},
	mounted() {
		this.$nuxt.$on("view-data-read-completed", (data) => {
			if(!data) return
			this.rules = data.rules?data.rules:{};
		});
	},
	methods: {},
	beforeDestroy(){
		this.$nuxt.$off("view-data-read-completed");
	},
}
</script>
