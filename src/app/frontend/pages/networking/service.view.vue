<template>
<div>
	<!-- 1. metadta -->
	<c-metadata v-model="metadata" dtCols="2" ddCols="10">
		<dt class="col-sm-2">Selector</dt>
		<dd class="col-sm-10">
			<span v-for="(val, idx) in info.selector" v-bind:key="idx" class="border-box background">{{ val }}</span>
		</dd>
		<dt v-if="metadata.ownerReferences" class="col-sm-2">Controlled By</dt>
		<dd v-if="metadata.ownerReferences" class="col-sm-10">{{ controller.kind }} <a href="#" @click="$emit('navigate', getViewLink(controller.group, controller.resource, metadata.namespace, controller.name))">{{ controller.name }}</a></dd>
		<dt class="col-sm-2">Type</dt><dd class="col-sm-10">{{ info.type }}</dd>
		<dt class="col-sm-2">Session Affinity</dt><dd class="col-sm-10">{{ info.sessionAffinity }}</dd>
	</c-metadata>
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
import VueMetadataView	from "@/components/view/metadataView.vue";
import VueEventsView	from "@/components/view/eventsView.vue";

export default {
	components: {
		"c-metadata": { extends: VueMetadataView },
		"c-events": { extends: VueEventsView }
	},
	data() {
		return {
			metadata: {},
			info: [],
			connection: [],
			endpoints: [],
			controller: {},
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
			this.controller = data.metadata.ownerReferences? this.getResource(data.metadata.ownerReferences[0]):{};
			this.info = {
				selector: data.spec.selector? this.stringifyLabels(data.spec.selector):[] ,
				type: data.spec.type || "",
				sessionAffinity: data.spec.sessionAffinity || 'None',
			};
			this.connection = {
				clusterIP: data.spec.clusterIP || "-",
				ports: this.toEndpointList(data.spec.ports,data.spec.type) || "",
			};
			this.endpoints = this.getEndpoints(data);
		});
		this.$nuxt.$emit("onCreated",'')
	},
	methods: {
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
