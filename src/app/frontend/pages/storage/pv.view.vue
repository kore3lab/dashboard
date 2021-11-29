<template>
<div>
	<!-- 1. metadata -->
	<c-metadata v-model="value" dtCols="3" ddCols="9">
		<dt v-if="metadata.finalizers" class="col-sm-3">Finalizers</dt>
		<dd v-if="metadata.finalizers" class="col-sm-9">
			<ul class="list-unstyled mb-0">
				<li v-for="(value, name) in metadata.finalizers" v-bind:key="name">{{ value }}</li>
			</ul>
		</dd>
		<dt class="col-sm-3">Capacity</dt><dd class="col-sm-9">{{ info.capacity }}</dd>
		<dt v-if="info.mountOption" class="col-sm-3">Mount Options</dt><dd v-if="info.mountOption" class="col-sm-9">{{ info.mountOption }} </dd>
		<dt class="col-sm-3">Access Modes</dt><dd class="col-sm-9"><span v-for="(val, idx) in info.accessModes" v-bind:key="idx">{{ val }} </span></dd>
		<dt class="col-sm-3">Reclaim Policy</dt><dd class="col-sm-9">{{ info.reclaimPolicy }}</dd>
		<dt class="col-sm-3">Storage Class Name</dt><dd class="col-sm-9">{{ info.storageClassName }}</dd>
		<dt class="col-sm-3">Status</dt>
		<dd class="col-sm-9">
			<span v-bind:class="{'badge-success':info.status=='Available'||info.status=='Bound','badge-warning':info.status=='Released'}" class="badge font-weight-light text-sm mb-1">{{ info.status }}</span>
		</dd>
	</c-metadata>
	<!-- 2. host path -->
	<div v-if="hostPath" class="row">
		<div class="col-md-12">
			<div class="card card-secondary card-outline">
				<div class="card-header p-2"><h3 class="card-title">Host Path</h3></div>
				<div class="card-body p-2">
					<dl class="row mb-0">
						<dt class="col-sm-3">Type</dt><dd class="col-sm-9">{{ hostPath.type ? hostPath.type : '-' }}</dd>
						<dt class="col-sm-3">Path</dt><dd class="col-sm-9">{{ hostPath.path }}</dd>
					</dl>
				</div>
			</div>
		</div>
	</div>
	<!-- 3. nfs -->
	<div v-if="nfs" class="row">
		<div class="col-md-12">
			<div class="card card-secondary card-outline">
				<div class="card-header p-2"><h3 class="card-title">Network File System</h3></div>
				<div class="card-body p-2">
					<dl class="row mb-0">
						<dt class="col-sm-3">Server</dt><dd class="col-sm-9">{{ nfs["server"] }}</dd>
						<dt class="col-sm-3">Path</dt><dd class="col-sm-9">{{ nfs["path"] }}</dd>
					</dl>
				</div>
			</div>
		</div>
	</div>
	<!-- 4. flex volume -->
	<div v-if="flexVolume" class="row">
		<div class="col-md-12">
			<div class="card card-secondary card-outline">
				<div class="card-header p-2"><h3 class="card-title">FlexVolume</h3></div>
				<div class="card-body p-2">
					<dl class="row mb-0">
						<dt class="col-sm-3">Driver</dt><dd class="col-sm-9">{{ flexVolume.driver }}</dd>
					</dl>
					<dl class="row mb-0">
						<dt class="col-sm-3" v-if="flexVolume.fsType">Type</dt>
						<dd class="col-sm-9" v-if="flexVolume.fsType">{{ flexVolume.fsType }}</dd>
						<dt class="col-sm-3" v-if="flexVolume.options">Options</dt>
						<dd class="col-sm-9" v-if="flexVolume.options">{{ flexVolume.options }}</dd>
					</dl>
				</div>
			</div>
		</div>
	</div>
	<!-- 5. claim -->
	<div v-if="claim" class="row">
		<div class="col-md-12">
			<div class="card card-secondary card-outline">
				<div class="card-header p-2"><h3 class="card-title">Claim</h3></div>
				<div class="card-body p-2">
					<dl class="row mb-0">
						<dt class="col-sm-3">Type</dt><dd class="col-sm-9">{{ claim.kind }}</dd>
						<dt class="col-sm-3">Name</dt><dd class="col-sm-9"><a href="#" @click="$emit('navigate', getViewLink(claim.group, claim.resource, claim.namespace, claim.name ))">{{ claim.name }}</a></dd>
						<dt class="col-sm-3">Namespace</dt><dd class="col-sm-9">{{ claim.namespace }}</dd>
					</dl>
				</div>
			</div>
		</div>
	</div>
	<!-- 6. events -->
	<c-events v-model="value" class="row"></c-events>

</div>
</template>
<script>
import VueMetadataView	from "@/components/view/metadataView.vue";
import VueEventsView	from "@/components/view/eventsView.vue";

export default {
	props:["value"],
	components: {
		"c-metadata": { extends: VueMetadataView },
		"c-events": { extends: VueEventsView }
	},
	data() {
		return {
			metadata: {},
			info: [],
			nfs: {},
			claim: {},
			flexVolume: {},
			hostPath: {},
		}
	},
	watch: {
		value(d) { this.onSync(d) }
	},
	methods: {
		onSync(data) {
			if(!data) return
			this.metadata = data.metadata;
			this.info = {
				capacity: data.spec.capacity.storage,
				mountOption: data.spec.mountOptions? data.spec.mountOptions.join(",") : '',
				reclaimPolicy: data.spec.persistentVolumeReclaimPolicy,
				accessModes: data.spec.accessModes || '-',
				storageClassName: data.spec.storageClassName || '-',
				status: data.status? data.status.phase : '-'
			};
			this.hostPath = data.spec.hostPath;
			this.nfs = data.spec.nfs;
			this.claim = this.getResource(data.spec.claimRef);
			this.flexVolume = data.spec.flexVolume;
		}
	}
}
</script>
