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
									<li v-for="(value, name) in metadata.annotations" v-bind:key="name"><span class="badge badge-secondary font-weight-light text-sm mb-1">{{ name }}={{ value }}</span></li>
								</ul>
							</dd>
							<dt class="col-sm-3">Labels</dt>
							<dd class="col-sm-9 text-truncate">
								<ul class="list-unstyled mb-0">
									<li v-for="(value, name) in metadata.labels" v-bind:key="name"><span class="badge badge-secondary font-weight-light text-sm mb-1">{{ name }}={{ value }}</span></li>
								</ul>
							</dd>
							<dt class="col-sm-3">UID</dt><dd class="col-sm-9">{{ metadata.uid }}</dd>
							<dt v-if="metadata.finalizers" class="col-sm-3">Finalizers</dt>
							<dd v-if="metadata.finalizers" class="col-sm-9 text-truncate">
								<ul class="list-unstyled mb-0">
									<li v-for="(value, name) in metadata.finalizers" v-bind:key="name"><span class="badge badge-secondary font-weight-light text-sm mb-1">{{ value }}</span></li>
								</ul>
							</dd>
							<dt class="col-sm-3">Access Modes</dt><dd class="col-sm-9"><span v-for="(val, idx) in info.accessModes" v-bind:key="idx">{{ val }} </span></dd>
							<dt class="col-sm-3">Storage Class Name</dt><dd class="col-sm-9">{{ info.storageClassName }}</dd>
							<dt class="col-sm-3">Storage</dt><dd class="col-sm-9">{{ info.storage }}</dd>
							<dt v-if="info.volumeName" class="col-sm-3">Binding Volume</dt><dd v-if="info.volumeName" class="col-sm-9"><a href="#" @click="$emit('navigate', getViewLink('', 'persistentvolumes','', info.volumeName ))">{{ info.volumeName }}</a></dd>
							<dt class="col-sm-3">Pods</dt>
							<dd class="col-sm-9">
								<span v-for="(val, idx) in pods" v-bind:key="idx"><a href="#" @click="$emit('navigate', getViewLink('', 'pods', metadata.namespace, val ))">{{ val }}</a></span>
								<span v-if="!isPod">-</span>
							</dd>
							<dt class="col-sm-3">Status</dt><dd class="col-sm-9">{{ info.status }}</dd>
						</dl>
					</div>
				</div>
			</div>
		</div>

		<div class="row">
			<div class="col-md-12">
				<div class="card card-secondary card-outline">
					<div class="card-header p-2"><h3 class="card-title text-md">Selector</h3></div>
					<div class="card-body p-2">
						<dl class="row mb-0">
							<dt class="col-sm-3">Match Labels</dt><dd class="col-sm-9"><span v-for="(val, idx) in selector.label" v-bind:key="idx" class="badge badge-secondary font-weight-light text-sm mr-1">{{ val }} </span><span v-if="!isLabel"> - </span></dd>
							<dt class="col sm-3">Match Expressions</dt>
							<dd class="col-sm-9">
								<ul v-for="(val, idx) in selector.expression" v-bind:key="idx" class="row list-unstyled mb-0">
									<li class="col-sm-3 mb-1 text-bold">Key</li>
									<li class="col-sm-9">{{ val.key }}</li>
									<li class="col-sm-3 mb-1 text-bold">Operator</li>
									<li class="col-sm-9">{{ val.operator }}</li>
									<li class="col-sm-3 mb-1 text-bold">Values</li>
									<li class="col-sm-9">{{ val.values }}</li>
								</ul>
								<span v-if="!isEx"> - </span>
							</dd>
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
			pods: [],
			selector: { label:{}, expression:{} },
			isPod: false,
			isLabel: false,
			isEx: false,
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
			this.pods = this.getPods()
			this.selector = this.getSelector(data.spec.selector);
		},
		getInfo(data) {
			let storage;
			if(!data.spec.resources || !data.spec.resources.requests) storage = '-'
			else storage = data.spec.resources.requests.storage

			return {
				accessModes: data.spec.accessModes || '-',
				storageClassName: data.spec.storageClassName || '-',
				storage: storage,
				status: data.status? data.status.phase : '-',
				volumeName: data.spec.volumeName,
			}
		},
		getPods() {
			this.isPod = false;
			let podList = [];
			this.$axios.get(this.getApiUrl('','pods',this.metadata.namespace))
					.then( resp => {
						resp.data.items.forEach(el => {
							el.spec.volumes.forEach(e => {
								if(e.persistentVolumeClaim) {
									if(e.persistentVolumeClaim.claimName === this.metadata.name) {
										podList.push(el.metadata.name)
										this.isPod = true;
									}
								}
							})
						})
					})
			return podList
		},
		getSelector(sel) {
			this.isLabel = false;
			this.isEx = false;
			if(!sel) return { label: '', expression: ''}
			let label = [];
			let expression = [];
			let key =Object.keys(sel.matchLabels)
			key.forEach(el => {
				this.isLabel = true
				label.push(`${el}=${sel.matchLabels[el]}`)
			})
			sel.matchExpressions.forEach(el => {
				el.values.forEach(e => {
					this.isEx = true
					expression.push({
						key: el.key,
						operator: el.operator,
						values: e,
					})
				})
			})
			return {
				label: label,
				expression: expression,
			}
		},
	},
	beforeDestroy(){
		this.$nuxt.$off("onReadCompleted");
	},
}
</script>
