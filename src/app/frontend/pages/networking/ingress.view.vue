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
						<dt class="col-sm-2">Ports</dt><dd class="col-sm-10">{{ info.ports }}</dd>
						<dt v-if="info.tls" class="col-sm-2">TLS</dt><dd v-if="info.tls" class="col-sm-10"><span v-for="(val,idx) in info.tls" v-bind:key="idx">{{ val }} </span></dd>
						<dt v-if="info.service" class="col-sm-2">Service</dt><dd v-if="info.service" class="col-sm-10">{{ info.type }}</dd>
					</dl>
				</div>
			</div>
		</div>
	</div>
	<!-- 2. rules -->
	<div class="row">
		<div class="col-md-12">
			<div class="card card-secondary card-outline">
				<div class="card-header p-2"><h3 class="card-title">Rules</h3></div>
				<div class="card-body group">
					<ul>
						<li v-for="(val, idx) in rules" v-bind:key="idx">
							<p class="title">Host: {{ val.host }}</p>
							<b-table striped hover small :items="val.value" :fields="ruleFields"></b-table>
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
					<b-table striped hover small :items="lbIp" :fields="lbFields"></b-table>
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
