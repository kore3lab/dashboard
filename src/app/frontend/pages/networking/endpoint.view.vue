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
							<dt class="col-sm-2">UID</dt><dd class="col-sm-10">{{ metadata.uid }}</dd>
						</dl>
					</div>
				</div>
			</div>
		</div>

		<div class="row">
			<div class="col-md-12">
				<div class="card card-secondary card-outline">
					<div class="card-header p-2"><h3 class="card-title text-md">Subsets</h3></div>
					<div v-if="isSubset" class="card-body p-2">
						<div class="col-sm-12 m-1"><h3 class="text-lg">Addresses</h3></div>
						<div class="overflow-auto">
							<b-table striped hover small :items="address" :fields="adressesFields" class="text-truncate">
								<template v-slot:cell(target)="data">
									<a href="#" @click="$emit('navigate', getViewLink('', data.item.targetKind, data.item.targetNamespace, data.item.targetName))">{{ data.item.targetName }}</a>
								</template>
							</b-table>
						</div>
						<div class="col-sm-12 m-1"><h3 class="text-lg">Ports</h3></div>
						<div class="overflow-auto">
							<b-table striped hover small :items="ports" :fields="portsFields" class="text-truncate"></b-table>
						</div>
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
			address: [],
			ports: [],
			isSubset: false,
			adressesFields: [
				{ key: "ip", label: "IP" },
				{ key: "hostname", label: "Hostname" },
				{ key: "target", label: "Target" },
			],
			portsFields: [
				{ key: "port", label: "Port" },
				{ key: "name", label: "Name" },
				{ key: "protocol", label: "Protocol" },
			],
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
			this.getSubsets(data.subsets);
		},
		getSubsets(sub) {
			if(sub) {
				this.isSubset = true
				let adr = sub[0].addresses
				let port = sub[0].ports
				let adrlist = []
				let portlist = []
				if(adr){
					adr.forEach(el => {
						adrlist.push({
							ip: el.ip || '',
							hostname: el.hostname || '',
							targetName: el.targetRef ? el.targetRef.name : '',
							targetNamespace: el.targetRef ? el.targetRef.namespace : '',
							targetKind: this.getKind(el.targetRef) || '',
						})
					})
				}

				if(port) {
					port.forEach(el => {
						portlist.push({
							port: el.port,
							name: el.name,
							protocol: el.protocol
						})

					})
				}
				this.address = adrlist
				this.ports = portlist
			} else this.isSubset = false
		},
		getKind(ref){
			if(!ref) return
			let k = ref.kind
			let kind
			let len = k.length
			if(k[len-1] === 's') kind = k.toLowerCase() + 'es'
			else kind = k.toLowerCase() + 's'
			return kind
		}
	},
	beforeDestroy(){
		this.$nuxt.$off("onReadCompleted");
	},
}
</script>
