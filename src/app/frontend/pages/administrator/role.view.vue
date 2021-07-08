<template>
<div>
	<!-- 1. metadata -->
	<div class="row">
		<div class="col-md-12">
			<div class="card card-secondary card-outline">
				<div class="card-body p-2">
					<dl class="row mb-0">
						<dt class="col-sm-2">Create at</dt><dd class="col-sm-10">{{ this.getTimestampString(metadata.creationTimestamp)}} ago ({{ metadata.creationTimestamp }})</dd>
						<dt class="col-sm-2">Name</dt><dd class="col-sm-10">{{ metadata.name }}</dd>
						<dt class="col-sm-2">Namespace</dt><dd class="col-sm-10">{{ metadata.namespace }}</dd>
						<dt class="col-sm-2">Annotations</dt>
						<dd class="col-sm-10">
							<ul class="list-unstyled mb-0">
								<li v-for="(value, name) in metadata.annotations" v-bind:key="name">{{ name }}=<span class="font-weight-light">{{ value }}</span></li>
							</ul>
						</dd>
						<dt class="col-sm-2">Labels</dt>
						<dd class="col-sm-10">
							<span v-for="(value, name) in metadata.labels" v-bind:key="name" class="label">{{ name }}={{ value }}</span>
						</dd>
					</dl>
				</div>
			</div>
		</div>
	</div>
	<!-- 2. rules -->
	<div class="row">
		<div class="col-md-12">
			<div class="card card-secondary card-outline">
				<div class="card-header p-2"><h3 class="card-title">Rules</h3></div>
				<div class="card-body group">
					<ul>
						<li v-for="(val, idx) in rules" v-bind:key="idx">
							<dl class="row">
								<dt v-if="val.resources" class="col-sm-2">Resources</dt><dd v-if="val.resources" class="col-sm-9">{{ val.resources }}</dd>
								<dt v-if="val.verbs" class="col-sm-2">Verbs</dt><dd v-if="val.verbs" class="col-sm-9">{{ val.verbs }}</dd>
								<dt v-if="val.apiGroups" class="col-sm-2">Api Groups</dt><dd v-if="val.apiGroups" class="col-sm-9">{{ val.apiGroups }}</dd>
								<dt v-if="val.resourceNames" class="col-sm-2">Resource Names</dt><dd v-if="val.resourceNames" class="col-sm-9">{{ val.resourceNames }}</dd>
							</dl>
						</li>
					</ul>
				</div>
			</div>
		</div>
	</div>
	<!-- 3. events -->
	<c-events class="row" v-model="metadata.uid"></c-events>

</div>
</template>
<script>
import VueEventsView	from "@/components/view/eventsView.vue";

export default {
	components: {
		"c-events": { extends: VueEventsView }
	},
	data() {
		return {
			metadata: {},
			rules: [],
		}
	},
	mounted() {
		this.$nuxt.$on("onReadCompleted", (data) => {
			if(!data) return
			this.origin = data;
			this.metadata = data.metadata;
			this.rules = this.getRules(data.rules)
		});
		this.$nuxt.$emit("onCreated",'')
	},
	methods: {
		getRules(rules) {
			if(!rules) return

			let list = []
			rules.map(({ resourceNames, apiGroups, resources, verbs }, index) => {
				list.push({
					resourceNames: resourceNames? resourceNames.join(", ") : null,
					apiGroups: apiGroups? apiGroups.join(", ") ? apiGroups.join(", ") : '*'  : null,
					resources: resources? resources.join(", ") : null,
					verbs: verbs? verbs.join(", ") : null,
				})
			})
			return list
		},
	},
	beforeDestroy(){
		this.$nuxt.$off("onReadCompleted");
	},
}
</script>
