<template>
<div>
	<!-- 1. metadta -->
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
						<dt class="col-sm-2">Selector</dt>
						<dd class="col-sm-10">
							<span v-for="(val, idx) in info.selector" v-bind:key="idx" class="border-box background">{{ val }}</span>
						</dd>
						<dt v-if="metadata.ownerReferences" class="col-sm-2">Controlled By</dt>
						<dd v-if="metadata.ownerReferences" class="col-sm-10">{{ metadata.ownerReferences[0].kind }} <a href="#" @click="$emit('navigate', getViewLink(controller.g, controller.k, metadata.namespace, metadata.ownerReferences[0].name))">{{ metadata.ownerReferences[0].name }}</a></dd>
						<dt class="col-sm-2">Type</dt><dd class="col-sm-10">{{ info.type }}</dd>
						<dt class="col-sm-2">Session Affinity</dt><dd class="col-sm-10">{{ info.sessionAffinity }}</dd>
					</dl>
				</div>
			</div>
		</div>
	</div>
	<!-- 2. connection -->
	<div class="row">
		<div class="col-md-12">
			<div class="card card-secondary card-outline">
				<div class="card-header p-2"><h3 class="card-title">Connection</h3></div>
				<div class="card-body p-2">
					<dl class="row mb-0">
						<dt class="col-sm-2">Cluster IP</dt><dd class="col-sm-10">{{ connection.clusterIP }}</dd>
						<dt class="col-sm-2">Ports</dt><dd class="col-sm-10">
						<ul class="list-unstyled">
							<li v-for="(val,idx) in connection.ports" v-bind:key="idx">{{ val }}</li>
						</ul>
					</dd>
					</dl>
				</div>
			</div>
		</div>
	</div>
	<!-- 3. endpoint -->
	<div class="row">
		<div class="col-md-12">
			<div class="card card-secondary card-outline">
				<div class="card-header p-2"><h3 class="card-title">Endpoint</h3></div>
				<div v-show="isEndpoint" class="card-body p-2 overflow-auto">
					<b-table striped hover small :items="endpoints" :fields="fields">
						<template v-slot:cell(name)="data">
							<a href="#" @click="$emit('navigate', getViewLink('', 'endpoints', data.item.namespace, data.item.name))">{{ data.value }}</a>
						</template>
						<template v-slot:cell(endpoints)="data">
							<span v-for="(val, idx) in data.item.endpoints" v-bind:key="idx">{{ val }} </span>
						</template>
					</b-table>
				</div>
			</div>
		</div>
	</div>
	<!-- 4. events -->
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
			connection: [],
			endpoints: [],
			controller: [],
			isEndpoint: false,
			fields: [
				{ key: "name", label: "Name" },
				{ key: "endpoints", label: "Endpoints" },
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
			this.controller = this.getController(data.metadata.ownerReferences);
			this.info = this.getInfo(data);
			this.connection = this.getConnection(data);
			this.endpoints = this.getEndpoints(data);
		},
		getInfo(data) {
			let selector = [];
			if(data.spec.selector) {
				selector = this.stringifyLabels(data.spec.selector)
			}
			return {
				selector: selector,
				type: data.spec.type || "",
				sessionAffinity: data.spec.sessionAffinity || 'None',
			}
		},
		getConnection(data) {
			return {
				clusterIP: data.spec.clusterIP || "-",
				ports: this.toEndpointList(data.spec.ports,data.spec.type) || "",
			}
		},
		getEndpoints(data) {
			let list =[];
			this.$axios.get(`${this.getApiUrl('', 'endpoints', data.metadata.namespace)}/${data.metadata.name}`)
				.then(resp => {
					this.isEndpoint = true;
					list.push({
						name: resp.data.metadata.name,
						namespace: resp.data.metadata.namespace,
						endpoints: this.onEndpoints(resp.data)
					})
					this.endpoints = list
					return list
				})
				.catch(_ => this.isEndpoint = false)
			return list
		},
		toEndpointList(p,type) {
			let list = [];
			if (p === undefined) return;
			for(let i =0; i < p.length; i++) {
				if (type === 'NodePort' || type === 'LoadBalancer') {
					list.push(`${p[i].port}:${p[i].nodePort}/${p[i].protocol}`)
				}else if(p[i].targetPort === p[i].port){
					list.push(`${p[i].port}/${p[i].protocol}`)
				}
				else{
					list.push(`${p[i].port}:${p[i].targetPort}/${p[i].protocol}`)
				}
			}
			return list;
		},
		onEndpoints(el){
			let list = [];
			if (el.subsets !== undefined) {
				if (el.subsets[0].notReadyAddresses) {
					return "-"
				}
				for (let i =0;i<el.subsets[0].addresses.length;i++){
					list.push(`${el.subsets[0].addresses[i].ip}`)
					if(i !== el.subsets[0].addresses.length -1) {
						list[i] += ','
					}
				}
				return list
			}
			return "-"
		},
	},
	beforeDestroy(){
		this.$nuxt.$off("onReadCompleted");
	},
}
</script>
