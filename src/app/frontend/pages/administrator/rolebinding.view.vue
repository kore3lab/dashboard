<template>
<div>
	<!-- 1. metadata -->
	<c-metadata v-model="metadata" dtCols="2" ddCols="10" @navigate="$emit('navigate', arguments[0])"></c-metadata>
	<!-- 2. reference -->
	<div class="row">
		<div class="col-md-12">
			<div class="card card-secondary card-outline">
				<div class="card-header p-2"><h3 class="card-title">Reference</h3></div>
				<div class="card-body p-2">
					<b-table-lite small :items="roleRefs" :fields="fields"></b-table-lite>
				</div>
			</div>
		</div>
	</div>
	<!-- 3. bindings -->
	<div class="row">
		<div class="col-md-12">
			<div class="card card-secondary card-outline">
				<div class="card-header p-2"><h3 class="card-title">Bindings</h3></div>
				<div class="card-body p-2">
					<div v-if="bindings">
						<b-table-lite small :items="bindings" :fields="bindFields">
							<template v-slot:cell(namespace)="data">
								{{ data.value? data.value : '-'}}
							</template>
						</b-table-lite>
					</div>
				</div>
			</div>
		</div>
	</div>
	<!-- 4. events -->
	<c-events class="row" v-model="metadata.uid"></c-events>
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
			metadata: {},
			roleRefs: [],
			bindings: [],
			fields: [
				{ key: "kind", label: "Kind" },
				{ key: "name", label: "Name" },
				{ key: "apiGroup", label: "Api Group" },
			],
			bindFields: [
				{ key: "name", label: "Binding" },
				{ key: "kind", label: "Type" },
				{ key: "namespace", label: "Namespace" },
			]
		}
	},
	mounted() {
		this.$nuxt.$on("onReadCompleted", (data) => {
			if(!data) return
			this.metadata = data.metadata;
			this.roleRefs = [data.roleRef];
			this.bindings = data.subjects;
		});
		this.$nuxt.$emit("onCreated",'')
	},
	methods: {},
	beforeDestroy(){
		this.$nuxt.$off("onReadCompleted");
	},
}
</script>
