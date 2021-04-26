<template>
	<div class="card-body p-2">
		<div class="row">
			<div class="col-md-12">
				<div class="card card-secondary card-outline">
					<div class="card-body p-2">
						<dl class="row mb-0">
							<dt class="col-sm-3 text-truncate">Create at</dt><dd class="col-sm-9">{{ this.getTimestampString(metadata.creationTimestamp)}} ago ({{ metadata.creationTimestamp }})</dd>
							<dt class="col-sm-3">Name</dt><dd class="col-sm-9">{{ metadata.name }}</dd>
							<dt class="col-sm-3">Namespace</dt><dd class="col-sm-9">{{ metadata.namespace }}</dd>
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
							<dt v-if="containerLimit" class="col-sm-4 text-truncate">Container Limits</dt>
							<dd v-if="containerLimit" class="col-sm-8">
								<dl v-for="(val, idx) in Object.keys(containerLimit)" v-bind:key="idx" class="row mb-0">
									<dt class="col-sm-3">{{ val }}</dt><dd class="col-sm-9"><span v-for="(v,i) in containerLimit[val]" v-bind:key="i" class="badge badge-secondary font-weight-light text-sm mb-1 mr-1">{{ v }}</span></dd>
								</dl>
								<hr>
							</dd>
							<dt v-if="podLimit" class="col-sm-4 text-truncate">Pod Limits</dt>
							<dd v-if="podLimit" class="col-sm-8">
								<dl v-for="(val, idx) in Object.keys(podLimit)" v-bind:key="idx" class="row mb-0">
									<dt class="col-sm-3">{{ val }}</dt><dd class="col-sm-9"><span v-for="(v,i) in podLimit[val]" v-bind:key="i" class="badge badge-secondary font-weight-light text-sm mb-1 mr-1">{{ v }}</span></dd>
								</dl>
								<hr>
							</dd>
							<dt v-if="pvcLimit" class="col-sm-4 text-truncate">Persistent Volume Claim Limits</dt>
							<dd v-if="pvcLimit" class="col-sm-8">
								<dl v-for="(val, idx) in Object.keys(pvcLimit)" v-bind:key="idx" class="row mb-0">
									<dt class="col-sm-3">{{ val }}</dt><dd class="col-sm-9"><span v-for="(v,i) in pvcLimit[val]" v-bind:key="i" class="badge badge-secondary font-weight-light text-sm mb-1 mr-1">{{ v }}</span></dd>
								</dl>
							</dd>
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
			containerLimit: [],
			podLimit: [],
			pvcLimit: [],
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
			this.containerLimit = this.getContainerLimit(data.spec.limits)
			this.podLimit = this.getPodLimit(data.spec.limits)
			this.pvcLimit = this.getPvcLimit(data.spec.limits)
		},
		getContainerLimit(limits) {
			this.containerLimit = [];
			let list = [];
			let cpu = [];
			let memory = [];
			let ephemeral = [];
			if(!limits) return
			list = limits.filter(limit => limit.type === 'Container')
			if(list.length === 0) return
			list.forEach(el => {
				let keys = Object.keys(el)
				keys.forEach(e => {
					if(el[e].cpu) cpu.push(`${ e } : ${ el[e].cpu }`)
					if(el[e].memory) memory.push(`${ e } : ${ el[e].memory }`)
					if(el[e]['ephemeral-storage']) ephemeral.push(`${ e } : ${ el[e]['ephemeral-storage'] }`)
				})
			})
			this.containerLimit.cpu = cpu
			this.containerLimit.memory = memory
			this.containerLimit.ephemeral = ephemeral
			return this.containerLimit

		},
		getPodLimit(limits) {
			this.podLimit = [];
			let list = [];
			let cpu = [];
			let memory = [];
			let ephemeral = [];
			if(!limits) return
			list = limits.filter(limit => limit.type === 'Pod')
			if(list.length === 0) return
			list.forEach(el => {
				let keys = Object.keys(el)
				keys.forEach(e => {
					if(el[e].cpu) cpu.push(`${ e } : ${ el[e].cpu }`)
					if(el[e].memory) memory.push(`${ e } : ${ el[e].memory }`)
					if(el[e]['ephemeral-storage']) ephemeral.push(`${ e } : ${ el[e]['ephemeral-storage'] }`)
				})
			})
			this.podLimit.cpu = cpu
			this.podLimit.memory = memory
			this.podLimit.ephemeral = ephemeral
			return this.podLimit
		},
		getPvcLimit(limits) {
			this.pvcLimit = [];
			let list = [];
			let storage = [];
			if(!limits) return
			list = limits.filter(limit => limit.type === 'PersistentVolumeClaim')
			if(list.length === 0) return
			list.forEach(el => {
				let keys = Object.keys(el)
				keys.forEach(e => {
					if(el[e].storage){
						storage.push(`${ e } : ${ el[e].storage }`)
					}
				})
			})
			this.pvcLimit.storage = storage
			return this.pvcLimit
		},
	},
	beforeDestroy(){
		this.$nuxt.$off("onReadCompleted");
	},
}
</script>
