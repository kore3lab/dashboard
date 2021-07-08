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
						<dt class="col-sm-2">UID</dt><dd class="col-sm-10">{{ metadata.uid }}</dd>
					</dl>
				</div>
			</div>
		</div>
	</div>
	<!-- 2. subsets -->
	<div class="row">
		<div class="col-md-12">
			<div class="card card-secondary card-outline">
				<div class="card-header p-2"><h3 class="card-title">Subsets</h3></div>
				<div class="card-body group">
					<ul>
						<li>
							<p class="title">Addresses</p>
							<b-table striped hover small :items="address" :fields="adressesFields">
								<template v-slot:cell(target)="data">
									<a href="#" @click="$emit('navigate', getViewLink('', data.item.targetKind, data.item.targetNamespace, data.item.targetName))">{{ data.item.targetName }}</a>
								</template>
							</b-table>
						</li>
						<li>
							<p class="title">Ports</p>
							<b-table striped hover small :items="ports" :fields="portsFields"></b-table>
						</li>
					</ul>
				</div>
			</div>
		</div>
	</div>
	<!-- 3. evenets -->
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
