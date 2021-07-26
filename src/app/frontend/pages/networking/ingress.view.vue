<template>
<div>
	<!-- 1. metadata -->
	<c-metadata v-model="metadata" dtCols="2" ddCols="10">
		<dt class="col-sm-2">Ports</dt><dd class="col-sm-10">{{ info.ports }}</dd>
		<dt v-if="info.tls" class="col-sm-2">TLS</dt><dd v-if="info.tls" class="col-sm-10"><span v-for="(val,idx) in info.tls" v-bind:key="idx">{{ val }} </span></dd>
		<dt v-if="info.service" class="col-sm-2">Service</dt><dd v-if="info.service" class="col-sm-10">{{ info.type }}</dd>
	</c-metadata>
	<!-- 2. rules -->
	<div class="row">
		<div class="col-md-12">
			<div class="card card-secondary card-outline">
				<div class="card-header p-2"><h3 class="card-title">Rules</h3></div>
				<div class="card-body group">
					<ul>
						<li v-for="(val, idx) in rules" v-bind:key="idx">
							<p class="title">Host: {{ val.host }}</p>
							<div class="ml-3">
								<b-table-lite :items="val.value" :fields="ruleFields" class="subset"></b-table-lite>
							</div>
						</li>
					</ul>
				</div>
			</div>
		</div>
	</div>
	<!-- 3. Load-Balancer Ingress Point -->
	<div class="row">
		<div class="col-md-12">
			<div class="card card-secondary card-outline">
				<div class="card-header p-2"><h3 class="card-title">Load-Balancer Ingress Points</h3></div>
				<div class="card-body group">
					<b-table-lite small :items="lbIp" :fields="lbFields"></b-table-lite>
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
			rules: [],
			lbIp: [],
			ruleFields: [
				{ key: "path", label: "Path" },
				{ key: "backends", label: "Backends" },
			],
			lbFields: [
				{ key: "hostname", label: "Hostname" },
				{ key: "ip", label: "IP" },
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
			this.rules = [];
			this.info = this.getInfo(data);
			this.getRules(data.spec.rules);
			this.lbIp = this.getLb(data.status.loadBalancer);
		},
		getInfo(data) {
			let ports = this.getPorts(data)
			let tls = this.getTls(data.spec)
			let service = this.getServiceNamePort(data.spec)
			return {
				ports: ports,
				tls: tls,
				service: service,

			}
		},
		getLb(lb) {
			let list = []
			if(lb.ingress) {
				lb.ingress.map((val, _) => {
					list.push({
						hostname: val.hostname? val.hostname: "-",
						ip: val.ip? val.ip: "-",
					})
				})
			}
			return list
		},
		getPorts(data) {
			const ports = [];
			const { spec: { tls, rules, backend, defaultBackend } } = data;
			const httpPort = 80;
			const tlsPort = 443;
			const servicePort = defaultBackend?.service.port.number ?? backend?.servicePort;

			if (rules && rules.length > 0) {
				if (rules.some(rule => rule.hasOwnProperty("http"))) {
					ports.push(httpPort);
				}
			} else if (servicePort !== undefined) {
				ports.push(Number(servicePort));
			}

			if (tls && tls.length > 0) {
				ports.push(tlsPort);
			}
			return ports.join(", ");
		},
		getTls(spec) {
			let list = [];
			if(spec.tls){
				spec.tls.map((tls, index) =>{
					list.push(tls.secretName)
				})
				return list
			}
		},
		getServiceNamePort(spec) {
			const serviceName = spec?.defaultBackend?.service.name ?? spec?.backend?.serviceName;
			const servicePort = spec?.defaultBackend?.service.port.number ?? spec?.defaultBackend?.service.port.name ?? spec?.backend?.servicePort;
			if(serviceName) return `${serviceName}:${servicePort}`
			else if(servicePort) return `${serviceName}:${servicePort}`
			else return
		},
		getRules(r) {
			let list = [];
			let host=[];
			if(r) {
				r.map((rule, index) => {
					host = rule.host
					list = [];
					rule.http.paths.map((p,_) => {
						const { serviceName, servicePort } = this.getBackendServiceNamePort(p.backend);
						const backend =`${serviceName}:${servicePort}`;
						list.push({
							path: p.path || "-",
							backends: backend
						})
					})
					this.rules[index] = {host: host, value: list}
				})
			}
		},
		getBackendServiceNamePort(b) {
			const serviceName = "service" in b ? b.service.name : b.serviceName
			const servicePort = "service" in b ? b.service.port.number ?? b.service.port.name : b.servicePort
			return { serviceName, servicePort };
		},
	},
	beforeDestroy(){
		this.$nuxt.$off("onReadCompleted");
	},
}
</script>