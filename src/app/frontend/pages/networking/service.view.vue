<template>
<div>
	<!-- 1. metadta -->
	<c-metadata dtCols="2" ddCols="10" @navigate="$emit('navigate', arguments[0])">
		<dt class="col-sm-2">Selector</dt>
		<dd class="col-sm-10">
			<span v-for="(val, idx) in info.selector" v-bind:key="idx" class="border-box background">{{ val }}</span>
		</dd>
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
						<ul class="list-unstyled mb-0">
							<li v-for="(val,idx) in connection.ports" v-bind:key="idx">{{ val }}</li>
						</ul>
					</dd>
					</dl>
				</div>
			</div>
		</div>
	</div>
	<!-- 3. endpoint -->
	<div v-if="endpoints.length>0" class="row">
		<div class="col-md-12">
			<div class="card card-secondary card-outline">
				<div class="card-header p-2"><h3 class="card-title">Endpoint</h3></div>
				<div class="card-body p-2 overflow-auto">
					<b-table-lite small :items="endpoints" :fields="fields" class="subset">
						<template v-slot:cell(name)="data">
							<a href="#" @click="$emit('navigate', getViewLink('', 'endpoints', data.item.namespace, data.item.name))">{{ data.value }}</a>
						</template>
						<template v-slot:cell(endpoints)="data">
							<span v-if = "data.item.endpoints.length > 0" >
								<span v-for="(val, idx) in data.item.endpoints" v-bind:key="idx"> {{val.ip}} </span>
							</span>
							<span v-else> - </span>
						</template>
					</b-table-lite>
				</div>
			</div>
		</div>
	</div>
	<!-- 4. events -->
	<c-events class="row"></c-events>

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
			info: [],
			connection: [],
			endpoints: [],
			fields: [
				{ key: "name", label: "Name" },
				{ key: "endpoints", label: "Endpoints" },
			],
		}
	},
	mounted() {
		this.$nuxt.$on("view-data-read-completed", (data) => {
			if(!data) return
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
	},
	methods: {
		getEndpoints(data) {
			let list =[];
			this.$axios.get(`${this.getApiUrl('', 'endpoints', data.metadata.namespace)}/${data.metadata.name}`)
				.then(resp => {
					list.push({
						name: resp.data.metadata.name,
						namespace: resp.data.metadata.namespace,
						endpoints: resp.data.subsets[0].addresses || []
					})
					return list
				})
			return list
		},
		toEndpointList(p,type) {
			let list = [];
			if (p === undefined) return;
			for(let i =0; i < p.length; i++) {
				(type === 'NodePort' || type === 'LoadBalancer') ? list.push(`${p[i].port}:${p[i].nodePort}/${p[i].protocol}`)
				: list.push(`${p[i].port}${p[i].port === p[i].targetPort ? "" : `:${p[i].targetPort}`}/${p[i].protocol}`)
			}
			return list;
		},
	},
	beforeDestroy(){
		this.$nuxt.$off("view-data-read-completed");
	},
}
</script>
