<template>
	<div class="card-body p-2">
		<div class="row">
			<div class="col-md-12">
				<div class="card card-secondary card-outline">
					<div class="card-body p-2">
						<dl class="row mb-0">
							<dt class="col-sm-3 text-truncate">Create at</dt><dd class="col-sm-9">{{ this.getTimestampString(metadata.creationTimestamp)}} ago ({{ metadata.creationTimestamp }})</dd>
							<dt class="col-sm-3">Name</dt><dd class="col-sm-9">{{ metadata.name }}</dd>
							<dt class="col-sm-3">Annotations</dt>
							<dd class="col-sm-9 text-truncate">
								<ul class="list-unstyled mb-0">
									<li v-for="(value, name) in metadata.annotations" v-bind:key="name"><span class="badge badge-secondary font-weight-light text-sm mb-1">{{ name }}:{{ value }}</span></li>
								</ul>
							</dd>
							<dt class="col-sm-3">Labels</dt>
							<dd class="col-sm-9 text-truncate">
								<ul class="list-unstyled mb-0">
									<li v-for="(value, name) in metadata.labels" v-bind:key="name"><span class="badge badge-secondary font-weight-light text-sm mb-1">{{ name }}:{{ value }}</span></li>
								</ul>
							</dd>
							<dt class="col-sm-3">UID</dt><dd class="col-sm-9">{{ metadata.uid }}</dd>
							<dt v-if="metadata.finalizers" class="col-sm-3">Finalizers</dt>
							<dd v-if="metadata.finalizers" class="col-sm-9 text-truncate">
								<ul class="list-unstyled mb-0">
									<li v-for="(value, name) in metadata.finalizers" v-bind:key="name"><span class="badge badge-secondary font-weight-light text-sm mb-1">{{ value }}</span></li>
								</ul>
							</dd>
							<dt class="col-sm-3">Capacity</dt><dd class="col-sm-9">{{ info.capacity }}</dd>
							<dt v-if="info.mountOption" class="col-sm-3">Mount Options</dt><dd v-if="info.mountOption" class="col-sm-9">{{ info.mountOption }} </dd>
							<dt class="col-sm-3">Access Modes</dt><dd class="col-sm-9"><span v-for="(val, idx) in info.accessModes" v-bind:key="idx">{{ val }} </span></dd>
							<dt class="col-sm-3">Reclaim Policy</dt><dd class="col-sm-9">{{ info.reclaimPolicy }}</dd>
							<dt class="col-sm-3">Storage Class Name</dt><dd class="col-sm-9">{{ info.storageClassName }}</dd>
							<dt class="col-sm-3">Status</dt><dd class="col-sm-9"><span class="badge badge-secondary font-weight-light text-sm mb-1">{{ info.status }}</span></dd>
						</dl>
					</div>
				</div>
			</div>
		</div>

		<div v-if="nfs" class="row">
			<div class="col-md-12">
				<div class="card card-secondary card-outline">
					<div class="card-header p-2"><h3 class="card-title text-md">Network File System</h3></div>
					<div class="card-body p-2">
						<dl v-for="(val, idx) in nfs" v-bind:key="idx" class="row mb-0">
							<dt class="col-sm-2 text-truncate">{{ val.key }}</dt><dd class="col-sm-10">{{ val.val }}</dd>
						</dl>
					</div>
				</div>
			</div>
		</div>

		<div v-if="flexVolume" class="row">
			<div class="col-md-12">
				<div class="card card-secondary card-outline">
					<div class="card-header p-2"><h3 class="card-title text-md">FlexVolume</h3></div>
					<div class="card-body p-2">
						<dl class="row mb-0">
							<dt class="col-sm-2 text-truncate">Driver</dt><dd class="col-sm-2">{{ flexVolume.driver }}</dd>
						</dl>
						<dl v-for="(val, idx) in flexVolume.list" v-bind:key="idx" class="row mb-0">
							<dt class="col-sm-2">{{ val.key }}</dt><dd class="col-sm-10">{{ val.val }}</dd>
						</dl>
					</div>
				</div>
			</div>
		</div>

		<div v-if="claim" class="row">
			<div class="col-md-12">
				<div class="card card-secondary card-outline">
					<div class="card-header p-2"><h3 class="card-title text-md">Claim</h3></div>
					<div class="card-body p-2">
						<dl class="row mb-0">
							<dt class="col-sm-2 text-truncate">Type</dt><dd class="col-sm-10">{{ claim.type }}</dd>
							<dt class="col-sm-2 text-truncate">Name</dt><dd class="col-sm-10"><a href="#" @click="$emit('navigate', getViewLink(claim.group.g, claim.group.k, claim.namespace, claim.name ))">{{ claim.name }}</a></dd>
							<dt class="col-sm-2 text-truncate">Namespace</dt><dd class="col-sm-10">{{ claim.namespace }}</dd>
						</dl>
					</div>
				</div>
			</div>
		</div>

		<div class="row">
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
			info: [],
			nfs: [],
			claim: [],
			flexVolume: [],
		}
	},
	mounted() {
		this.$nuxt.$on("onReadCompleted", (data) => {
			if(!data) return
			this.metadata = data.metadata;
			this.onSync(data)
		});
		this.$nuxt.$emit("onCreated",'')
	},
	methods: {
		onSync(data) {
			this.event = this.getEvents(data.metadata.uid,'fieldSelector=involvedObject.name='+data.metadata.name);
			this.info = this.getInfo(data);
			this.nfs = this.getNFS(data.spec.nfs);
			this.claim = this.getClaim(data.spec.claimRef);
			this.flexVolume = this.getFV(data.spec.flexVolume);
		},
		getInfo(data) {

			return {
				capacity: data.spec.capacity.storage,
				mountOption: data.spec.mountOptions? data.spec.mountOptions.join(",") : '',
				reclaimPolicy: data.spec.persistentVolumeReclaimPolicy,
				accessModes: data.spec.accessModes || '-',
				storageClassName: data.spec.storageClassName || '-',
				status: data.status? data.status.phase : '-'
			}
		},
		getNFS(nfs) {
			let list = [];
			if(!nfs) return
			Object.entries(nfs).map(([name, value]) => {
				list.push({
					key: name,
					val: value
				})
			})
			return list
		},
		getClaim(ref) {
			if(!ref) return

			return {
				type: ref.kind,
				name: ref.name,
				namespace: ref.namespace,
				group: this.getController(ref)
			}
		},
		getFV(fv) {
			let list = []
			if(!fv) return
			Object.entries(fv).map(([name, value]) => {
				list.push({
					key: name,
					val: value
				})
			})
			return {
				list: list,
				driver: fv.driver
			}
		},
	},
	beforeDestroy(){
		this.$nuxt.$off("onReadCompleted");
	},
}
</script>
