<template>
<div>
	<!-- 1. metadata -->
	<div class="row">
		<div class="col-md-12">
			<div class="card card-secondary card-outline">
				<div class="card-body p-2">
					<dl class="row mb-0">
						<dt class="col-sm-3">Create at</dt><dd class="col-sm-9">{{ this.getTimestampString(metadata.creationTimestamp)}} ago ({{ metadata.creationTimestamp }})</dd>
						<dt class="col-sm-3">Name</dt><dd class="col-sm-9">{{ metadata.name }}</dd>
						<dt class="col-sm-3">Namespace</dt><dd class="col-sm-9">{{ metadata.namespace }}</dd>
						<dt class="col-sm-3">Message</dt><dd class="col-sm-9">{{ info.message }}</dd>
						<dt class="col-sm-3">Reason</dt><dd class="col-sm-9">{{ info.reason }}</dd>
						<dt class="col-sm-3">Source</dt><dd class="col-sm-9"><span v-for="(val, idx) in info.source" v-bind:key="idx">{{ val }} </span></dd>
						<dt class="col-sm-3">First seen</dt><dd class="col-sm-9"><span v-if="info.firstSeen">{{ this.getTimestampString(info.firstSeen)}} ago ({{ info.firstSeen }})</span><span v-else>-</span></dd>
						<dt class="col-sm-3">Last seen</dt><dd class="col-sm-9"><span v-if="info.lastSeen">{{ this.getTimestampString(info.lastSeen)}} ago ({{ info.lastSeen }})</span><span v-else>-</span></dd>
						<dt class="col-sm-3">Count</dt><dd class="col-sm-9">{{ info.count ? info.count : '-' }}</dd>
						<dt class="col-sm-3">Type</dt><dd class="col-sm-9"><span v-bind:class="{'text-warning':info.type=='Warning'}">{{ info.type }}</span></dd>
					</dl>
				</div>
			</div>
		</div>
	</div>
	<!-- 2. Involved object -->
	<div class="row">
		<div class="col-md-12">
			<div class="card card-secondary card-outline m-0">
				<div class="card-header p-2"><h3 class="card-title">Involved object</h3></div>
				<div class="card-body p-2">
					<b-table-lite striped hover small :items="involvedObject" :fields="fields">
						<template v-slot:cell(name)="data">
							<a href="#" @click="$emit('navigate', getViewLink(data.item.group, data.item.resource, data.item.namespace, data.value))">{{ data.value }}</a>
						</template>
					</b-table-lite>
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
			info: {},
			involvedObject: [],
			fields: [
				{ key: "name", label: "Name" },
				{ key: "namespace", label: "Namespace"  },
				{ key: "kind", label: "Kind" },
				{ key: "fieldPath", label: "Field Path" },
			],
		}
	},
	mounted() {
		this.$nuxt.$on("view-data-read-completed", (data) => {
			if(!data) return
			this.metadata = data.metadata;
			this.info = {
				message: data.message,
				reason: data.reason,
				source: Object.entries(data.source).map(([_, value]) => `${value}`),
				firstSeen: data.firstTimestamp,
				lastSeen: data.lastTimestamp,
				count: data.count,
				type: data.type,
			};
			this.involvedObject = [this.getResource(data.involvedObject)];
			this.involvedObject[0]["fieldPath"] = data.involvedObject["fieldPath"]
		});
	},
	methods: {},
	beforeDestroy(){
		this.$nuxt.$off("view-data-read-completed");
	},
}
</script>
