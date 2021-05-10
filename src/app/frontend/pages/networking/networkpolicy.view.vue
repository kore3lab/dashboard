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
									<li v-for="(value, name) in metadata.annotations" v-bind:key="name"><span class="badge badge-secondary font-weight-light text-sm mb-1">{{ name }}:{{ value }}</span></li>
								</ul>
							</dd>
							<dt class="col-sm-2">Labels</dt>
							<dd class="col-sm-10 text-truncate">
								<ul class="list-unstyled mb-0">
									<li v-for="(value, name) in metadata.labels" v-bind:key="name"><span class="badge badge-secondary font-weight-light text-sm mb-1">{{ name }}:{{ value }}</span></li>
								</ul>
							</dd>
							<dt class="col-sm-2">UID</dt><dd class="col-sm-10">{{ metadata.uid }}</dd>
							<dt class="col-sm-2">Pod Selector</dt>
							<dd class="col-sm-10 text-truncate">
								<ul class="list-unstyled mb-0">
									<li v-for="(value, name) in info.label" v-bind:key="name"><span class="badge badge-secondary font-weight-light text-sm mb-1">{{ value }}</span></li>
									<li v-if="!isLabel">(empty) (Allowing the specific traffic to all pods in this namespace)</li>
								</ul>
							</dd>
						</dl>
					</div>
				</div>
			</div>
		</div>

		<div v-if="isIngress" class="row">
			<div class="col-md-12">
				<div class="card card-secondary card-outline">
					<div class="card-header p-2"><h3 class="card-title text-md">Ingress</h3></div>
					<div class="card-body p-2">
						<dl class="row mb-1">
							<dt class="col-sm-3">Ports</dt><dd class="col-sm-9"><span v-for="(val, idx) in ingress.port" v-bind:key="idx">{{ val }} </span><span v-if="!isIngPort">-</span></dd>
						</dl>
						<dl v-if="isIngFrom" class="row mb-1">
							<dt class="col-sm-12 mt-1"><h5 class="">From</h5></dt>
							<dt v-if="isIngIpBlock" class="col-sm-3">ipBlock</dt><dd v-if="isIngIpBlock" class="col-sm-9"><span v-for="(val, idx) in ingress.ipblock" v-bind:key="idx">{{ val }} </span></dd>
							<dt v-if="ingress.nslabel" class="col-sm-3">namespaceSelector</dt><dd v-if="ingress.nslabel" class="col-sm-9"><span v-for="(val, idx) in ingress.nslabel" v-bind:key="idx">{{ val }} </span></dd>
							<dt v-if="ingress.podsel" class="col-sm-3">podSelector</dt><dd v-if="ingress.podsel" class="col-sm-9"><span v-for="(val, idx) in ingress.podsel" v-bind:key="idx">{{ val }} </span></dd>
						</dl>
					</div>
				</div>
			</div>
		</div>

		<div v-if="isEgress" class="row">
			<div class="col-md-12">
				<div class="card card-secondary card-outline">
					<div class="card-header p-2"><h3 class="card-title text-md">Egress</h3></div>
					<div class="card-body p-2">
						<dl class="row mb-1">
							<dt class="col-sm-3">Ports</dt><dd class="col-sm-9"><span v-for="(val, idx) in egress.port" v-bind:key="idx">{{ val }} </span><span v-if="!isEgPort">-</span></dd>
						</dl>
						<dl v-if="isEgTo" class="row mb-1">
							<dt class="col-sm-12 mt-1"><h5 class="">To</h5></dt>
							<dt v-if="isEgIpBlock" class="col-sm-3">ipBlock</dt><dd v-if="isEgIpBlock" class="col-sm-9"><span v-for="(val, idx) in egress.ipblock" v-bind:key="idx">{{ val }} </span></dd>
							<dt v-if="egress.nslabel" class="col-sm-3">namespaceSelector</dt><dd v-if="egress.nslabel" class="col-sm-9"><span v-for="(val, idx) in egress.nslabel" v-bind:key="idx">{{ val }} </span></dd>
							<dt v-if="egress.podsel" class="col-sm-3">podSelector</dt><dd v-if="egress.podsel" class="col-sm-9"><span v-for="(val, idx) in egress.podsel" v-bind:key="idx">{{ val }} </span></dd>
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
			ingress: [],
			egress: [],
			isLabel: false,
			isIngress: false,
			isIngFrom: false,
			isIngIpBlock: false,
			isIngPort: false,
			isEgress: false,
			isEgTo: false,
			isEgIpBlock: false,
			isEgPort: false,
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
			this.ingress = this.getIngress(data.spec.ingress);
			this.egress = this.getEgress(data.spec.egress);
		},
		getInfo(data) {
			this.isLabel = false;
			let label = this.stringifyLabels(data.spec.podSelector? data.spec.podSelector.matchLabels : '')
			if(label.length !== 0 ){
				this.isLabel = true
			}
			return {
				label: label,
			}
		},
		getIngress(ing) {
			this.isIngress = false;
			this.isIngFrom = false;
			this.isIngIpBlock = false;
			this.isIngPort = false;
			let port = [];
			let ipblock = [];
			let nslabel;
			let podsel;
			if(ing) {
				this.isIngress = true;
				ing.forEach(el => {
					if(el.ports) {
						this.isIngPort = true;
						el.ports.forEach(e => {
							port.push(`${e.protocol || ""}:${e.port || ""}`)
						})
					}
					if(el.from) {
						this.isIngFrom = true;
						el.from.forEach(e => {
							let key = Object.keys(e)
							if( key[0] === 'ipBlock') {
								if (e.ipBlock.cidr) {
									this.isIngIpBlock = true;
									ipblock.push(`cidr:${e.ipBlock.cidr} ${e.ipBlock.except? `, except:${e.ipBlock.except[0]}` : ''}`)
								}
							} else if(key[0] === 'namespaceSelector') {
								nslabel = this.stringifyLabels(e.namespaceSelector.matchLabels)
							} else {
								podsel = this.stringifyLabels(e.podSelector.matchLabels)
							}
						})
					}
				})
			} else this.isIngress = false;
			return {
				port: port || '-',
				ipblock: ipblock,
				nslabel: nslabel,
				podsel: podsel,
			}
		},
		getEgress(eg) {
			this.isEgress = false;
			this.isEgTo = false;
			this.isEgIpBlock = false;
			this.isEgPort = false;
			let port = [];
			let ipblock = [];
			let nslabel;
			let podsel;
			if(eg) {
				this.isEgress = true;
				eg.forEach(el => {
					if(el.ports) {
						this.isEgPort = true;
						el.ports.forEach(e => {
							port.push(`${e.protocol || ""}:${e.port || ""}`)
						})
					}
					if(el.to) {
						this.isEgTo = true;
						el.to.forEach(e => {
							let key = Object.keys(e)
							if( key[0] === 'ipBlock') {
								if (e.ipBlock.cidr) {
									this.isEgIpBlock = true;
									ipblock.push(`cidr:${e.ipBlock.cidr} ${e.ipBlock.except? `, except:${e.ipBlock.except[0]}` : ''}`)
								}
							} else if(key[0] === 'namespaceSelector') {
								nslabel = this.stringifyLabels(e.namespaceSelector.matchLabels)
							} else {
								podsel = this.stringifyLabels(e.podSelector.matchLabels)
							}
						})
					}
				})
			} else this.isEgress = false;
			return {
				port: port || '-',
				ipblock: ipblock,
				nslabel: nslabel,
				podsel: podsel,
			}
		},
	},
	beforeDestroy(){
		this.$nuxt.$off("onReadCompleted");
	},
}
</script>
