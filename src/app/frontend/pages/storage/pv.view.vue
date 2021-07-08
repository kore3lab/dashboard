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
						<dt class="col-sm-3">Annotations</dt>
						<dd class="col-sm-9">
							<ul class="list-unstyled mb-0">
								<li v-for="(value, name) in metadata.annotations" v-bind:key="name">{{ name }}=<span class="font-weight-light">{{ value }}</span></li>
							</ul>
						</dd>
						<dt class="col-sm-3">Labels</dt>
						<dd class="col-sm-9">
							<span v-for="(value, name) in metadata.labels" v-bind:key="name" class="label">{{ name }}={{ value }}</span>
						</dd>
						<dt class="col-sm-3">UID</dt><dd class="col-sm-9">{{ metadata.uid }}</dd>
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
							<span class="badge badge-secondary font-weight-light text-sm mb-1">{{ info.status }}</span>
						</dd>
					</dl>
				</div>
			</div>
		</div>
	</div>
	<!-- 2. host path -->
	<div v-if="hostPath" class="row">
		<div class="col-md-12">
			<div class="card card-secondary card-outline">
				<div class="card-header p-2"><h3 class="card-title">Host Path</h3></div>
				<div class="card-body p-2">
					<dl class="row mb-0">
						<dt class="col-sm-2">Type</dt><dd class="col-sm-10">{{ hostPath.type ? hostPath.type : '-' }}</dd>
						<dt class="col-sm-2">Path</dt><dd class="col-sm-10">{{ hostPath.path }}</dd>
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
					<dl v-for="(val, idx) in nfs" v-bind:key="idx" class="row mb-0">
						<dt class="col-sm-2">{{ val.key }}</dt><dd class="col-sm-10">{{ val.val }}</dd>
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
						<dt class="col-sm-2">Driver</dt><dd class="col-sm-2">{{ flexVolume.driver }}</dd>
					</dl>
					<dl v-for="(val, idx) in flexVolume.list" v-bind:key="idx" class="row mb-0">
						<dt class="col-sm-2">{{ val.key }}</dt><dd class="col-sm-10">{{ val.val }}</dd>
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
						<dt class="col-sm-2">Type</dt><dd class="col-sm-10">{{ claim.type }}</dd>
						<dt class="col-sm-2">Name</dt><dd class="col-sm-10"><a href="#" @click="$emit('navigate', getViewLink(claim.group.g, claim.group.k, claim.namespace, claim.name ))">{{ claim.name }}</a></dd>
						<dt class="col-sm-2">Namespace</dt><dd class="col-sm-10">{{ claim.namespace }}</dd>
					</dl>
				</div>
			</div>
		</div>
	</div>
	<!-- 6. events -->
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
			info: [],
			nfs: [],
			claim: [],
			flexVolume: [],
			hostPath: [],
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
			this.info = this.getInfo(data);
			this.hostPath = this.getHP(data.spec.hostPath)
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
		getHP(hp) {
			if(!hp) return

			return {
				type: hp.type,
				path: hp.path,
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
