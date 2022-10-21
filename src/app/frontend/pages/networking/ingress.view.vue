<template>
<div>
	<!-- 1. metadata -->
	<c-metadata v-model="value" dtCols="2" ddCols="10">
		<dt class="col-sm-2">Ports</dt><dd class="col-sm-10">{{ info.ports }}</dd>
		<dt v-if="info.tls" class="col-sm-2">TLS</dt><dd v-if="info.tls" class="col-sm-10"><span v-for="(val,idx) in info.tls" v-bind:key="idx">{{ val }} </span></dd>
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
					<b-table-lite small :items="lbIp" :fields="lbFields" class="subset"></b-table-lite>
				</div>
			</div>
		</div>
	</div>
	<!-- 4. events -->
	<c-events v-model="value" class="row"></c-events>

</div>
</template>
<script>
import VueMetadataView	from "@/components/view/metadataView.vue";
import VueEventsView	from "@/components/view/eventsView.vue";

export default {
	props:["value"],
	components: {
		"c-metadata": { extends: VueMetadataView },
		"c-events": { extends: VueEventsView }
	},
	data() {
		return {
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
	watch: {
		value(d) { this.onSync(d) }
	},
	methods: {
		onSync(data) {
			if(!data) return
			this.info = {
				ports : this.getPorts(data),
				tls: data.spec.tls ? [...new Set((data.spec.tls).map(tls => tls.secretName))] ?? '' : ''
			}
			this.getRules(data.spec.rules);
			this.lbIp = [{
				ip : data.status.loadBalancer.ingress[0].ip ?? "-",
				hostname : data.status.loadBalancer.ingress[0].hostname ?? "-"
			}]
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
		getRules(r) {
			let list = [];
			let host=[];
			if(r) {
				r.map((rule, index) => {
					host = rule.host
					list = [];
					rule.http.paths.map((p,_) => {
						const serviceName = "service" in p.backend && p.backend.service ? p.backend.service.name : p.backend.serviceName;
						const servicePort = "service" in p.backend && p.backend.service ? p.backend.service.port.number ?? p.backend.service.port.name : p.backend.servicePort;
						list.push({
							path: p.path || "-",
							backends: `${serviceName}:${servicePort}`
						})
					})
					this.rules[index] = {host: host, value: list}
				})
			}
		},
	}
}
</script>