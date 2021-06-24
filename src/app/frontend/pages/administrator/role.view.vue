<template>
	<div class="card-body p-2">
		<div class="row">
			<div class="col-md-12">
				<div class="card card-secondary card-outline">
					<div class="card-body p-2">
						<dl class="row mb-0">
							<dt class="col-sm-2 text-truncate">Create at</dt><dd class="col-sm-10">{{ this.getTimestampString(metadata.creationTimestamp)}} ago ({{ metadata.creationTimestamp }})</dd>
							<dt class="col-sm-2">Name</dt><dd class="col-sm-10">{{ metadata.name }}</dd>
							<dt class="col-sm-2">Namespace</dt><dd class="col-sm-10">{{ metadata.namespace }}</dd>
							<dt class="col-sm-2">Annotations</dt>
							<dd class="col-sm-10 text-truncate">
								<ul class="list-unstyled mb-0">
									<li v-for="(value, name) in metadata.annotations" v-bind:key="name"><span class="badge badge-secondary font-weight-light text-sm mb-1">{{ name }}={{ value }}</span></li>
								</ul>
							</dd>
							<dt class="col-sm-2">Labels</dt>
							<dd class="col-sm-10 text-truncate">
								<ul class="list-unstyled mb-0">
									<li v-for="(value, name) in metadata.labels" v-bind:key="name"><span class="badge badge-secondary font-weight-light text-sm mb-1">{{ name }}={{ value }}</span></li>
								</ul>
							</dd>
						</dl>
					</div>
				</div>
			</div>
		</div>

		<div class="row">
			<div class="col-md-12">
				<div class="card card-secondary card-outline">
					<div class="card-header p-2"><h3 class="card-title text-md">Rules</h3></div>
					<div class="card-body p-2">
						<div class="card card-secondary m-1 mb-2" v-for="(val, idx) in rules" v-bind:key="idx">
							<dl class="row m-1 ">
								<dt v-if="val.resources" class="col-sm-2">Resources</dt><dd v-if="val.resources" class="col-sm-9">{{ val.resources }}</dd>
								<dt v-if="val.verbs" class="col-sm-2">Verbs</dt><dd v-if="val.verbs" class="col-sm-9">{{ val.verbs }}</dd>
								<dt v-if="val.apiGroups" class="col-sm-2">Api Groups</dt><dd v-if="val.apiGroups" class="col-sm-9">{{ val.apiGroups }}</dd>
								<dt v-if="val.resourceNames" class="col-sm-2">Resource Names</dt><dd v-if="val.resourceNames" class="col-sm-9">{{ val.resourceNames }}</dd>
							</dl>
						</div>
					</div>
				</div>
			</div>
		</div>

		<div class="row" v-show="event.length>0">
			<div class="col-md-12">
				<div class="card card-secondary card-outline m-0">
					<div class="card-header p-2"><h3 class="card-title text-md">Events</h3></div>
					<div class="card-body p-2">
						<dl v-for="(val, idx) in event" v-bind:key="idx" class="row mb-0 card-body p-2 border-bottom">
							<dt class="col-sm-12"><p v-bind:class="val.type" class="mb-1">{{ val.name }}</p></dt>
							<dt class="col-sm-2 text-truncate">Source</dt><dd class="col-sm-10">{{ val.source }}</dd>
							<dt class="col-sm-2 text-truncate">Count</dt><dd class="col-sm-10">{{ val.count }}</dd>
							<dt class="col-sm-2 text-truncate">Sub-object</dt><dd class="col-sm-10">{{ val.subObject }}</dd>
							<dt class="col-sm-2 text-truncate">Last seen</dt><dd class="col-sm-10">{{ val.lastSeen }}</dd>
						</dl>
					</div>
				</div>
			</div>
		</div>

	</div>
</template>
<script>

export default {
	data() {
		return {
			metadata: {},
			event: [],
			rules: [],
		}
	},
	mounted() {
		this.$nuxt.$on("onReadCompleted", (data) => {
			if(!data) return
			this.origin = data;
			this.metadata = data.metadata;
			this.onSync(data)
		});
		this.$nuxt.$emit("onCreated",'')
	},
	methods: {
		onSync(data) {
			this.event = this.getEvents(data.metadata.uid,'fieldSelector=involvedObject.name='+data.metadata.name);
			this.rules = this.getRules(data.rules)
		},
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
