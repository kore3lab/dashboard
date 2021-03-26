<template>
	<div class="card-body p-2">
		<div class="row mb-0">
			<div class="col-md-12">
				<div class="card card-secondary card-outline">
					<div class="card-body p-2">
						<b-tabs content-class="mt-3" >
							<b-tab title="CPU" active title-link-class="border-top-0 border-right-0  border-left-0">
								<div v-if="isCpu" class="chart">
									<c-linechart id="cpu" :chart-data="chart.data.cpu" :options="chart.options.cpu" class="mw-100" style="height: 14em;"></c-linechart>
								</div>
								<div v-if="!isCpu" class="text-center"><p> Metrics not available at the moment</p></div>
							</b-tab>
							<b-tab title="Memory"  title-link-class="border-top-0 border-right-0  border-left-0">
								<div v-if="isMemory" class="chart">
									<c-linechart id="memory" :chart-data="chart.data.memory" :options="chart.options.memory" class="mw-100" style="height: 14em;"></c-linechart>
								</div>
								<div v-if="!isMemory" class="text-center"><p> Metrics not available at the moment</p></div>
							</b-tab>
						</b-tabs>
					</div>
				</div>
			</div>
		</div>

		<div class="row">
			<div class="col-md-12">
				<div class="card card-secondary card-outline">
					<div class="card-body p-2">
						<dl class="row mb-0">
							<dt class="col-sm-2 text-truncate">Create at</dt><dd class="col-sm-10">{{ this.getTimestampString(metadata.creationTimestamp)}} ago ({{ metadata.creationTimestamp }})</dd>
							<dt class="col-sm-2">Name</dt><dd class="col-sm-10">{{ metadata.name }}</dd>
							<dt class="col-sm-2">Namespace</dt><dd class="col-sm-10">{{ metadata.namespace }}</dd>
							<dt class="col-sm-2 text-truncate">Annotations</dt>
							<dd class="col-sm-10 text-truncate">
								<ul class="list-unstyled mb-0">
									<li v-for="(value, name) in metadata.annotations" v-bind:key="name"><span class="badge badge-secondary font-weight-light text-sm mb-1">{{ name }}:{{ value }}</span></li>
								</ul>
							</dd>
							<dt class="col-sm-2 text-truncate">Labels</dt>
							<dd class="col-sm-10">
								<ul class="list-unstyled mb-0">
									<li v-for="(value, name) in metadata.labels" v-bind:key="name"><span class="badge badge-secondary font-weight-light text-sm mb-1">{{ name }}:{{ value }}</span></li>
								</ul>
							</dd>
							<dt class="col-sm-2 text-truncate">UID</dt><dd class="col-sm-10">{{ metadata.uid }}</dd>
							<dt v-if="metadata.ownerReferences" class="col-sm-2 text-truncate">Controlled By</dt>
							<dd v-if="metadata.ownerReferences" class="col-sm-10">{{ metadata.ownerReferences[0].kind }} <a href="#" @click="$emit('navigate', getViewLink(controller.g, controller.k, metadata.namespace, metadata.ownerReferences[0].name))">{{ metadata.ownerReferences[0].name }}</a></dd>
							<dt class="col-sm-2 text-truncate">Status</dt><dd class="col-sm-10" v-bind:class="status.style">{{ status.value }}</dd>
							<dt class="col-sm-2">Node</dt>
							<dd v-if="raw.spec.nodeName" class="col-sm-10"><a href="#" @click="$emit('navigate', getViewLink('', 'nodes', '', raw.spec.nodeName? raw.spec.nodeName : '' ))">{{ raw.spec.nodeName}}</a></dd>
							<dd v-if="!raw.spec.nodeName" class="col-sm-10">-</dd>
							<dt class="col-sm-2 text-truncate">Pod IP</dt>
							<dd class="col-sm-10">
								<ul class="list-unstyled mb-0">
									<li v-for="(d, idx) in podInfo.podIP" v-bind:key="idx" class="mb-1">{{ d }}</li>
								</ul>
							</dd>
							<dt class="col-sm-2 text-truncate">Priority Class</dt><dd class="col-sm-10">{{ podInfo.priorityClass}}</dd>
							<dt class="col-sm-2 text-truncate">QoS Class</dt><dd class="col-sm-10">{{ podInfo.qosClass }}</dd>
							<dt v-if="podInfo.conditions" class="col-sm-2 text-truncate">Conditions</dt><dd v-if="podInfo.conditions" class="col-sm-10"><span v-for="(d, idx) in podInfo.conditions" v-bind:key="idx" class="badge badge-secondary font-weight-light text-sm mr-1">{{ d }}</span></dd>
							<dt v-if="podInfo.isNodeSelector" class="col-sm-2 text-truncate">Node Selector</dt><dd v-if="podInfo.isNodeSelector" class="col-sm-10"><span v-for="(d, idx) in podInfo.nodeSelector" v-bind:key="idx" class="badge badge-secondary font-weight-light text-sm mr-1">{{ d.name }}: {{ d.value }}</span></dd>
							<dt class="col-sm-2 text-truncate">Tolerations</dt>
							<dd class="col-sm-10">{{ podInfo.tolerations? podInfo.tolerations.length: "-" }}<a class="float-right" v-b-toggle.tol href="#tol-table" @click.prevent @click="onTol">{{onTols ? 'Hide' : 'Show'}}</a></dd>
							<b-collapse class="col-sm-12" id="tol-table"><b-table striped hover small :items="podInfo.tolerations"></b-table></b-collapse>

							<dt v-show="podInfo.isAffinity" class="col-sm-2 text-truncate">Affinities</dt>
							<dd v-show="podInfo.isAffinity" class="col-sm-10">{{ podInfo.affinities? Object.keys(podInfo.affinities).length: "-" }}<a class="float-right" v-b-toggle.affi href="#affi-json" @click.prevent @click="onAffi">{{onAffis ? 'Hide' : 'Show'}}</a>
								<b-collapse id="affi-json"><c-jsontree id="txtSpec" v-model="podInfo.affinities" class="card-body p-2 border"></c-jsontree></b-collapse>
							</dd>
							<dt v-if="podInfo.secret" class="col-sm-2 text-truncate">Secrets</dt>
							<dd v-if="podInfo.secret" class="col-sm-10" >
								<ul class="list-unstyled">
									<li v-for="(d, idx) in podInfo.secret" v-bind:key="idx" >
										<a href="#" @click="$emit('navigate', getViewLink('', 'secrets', metadata.namespace, d))">{{ d }}</a>
									</li>
								</ul>
							</dd>
						</dl>
					</div>
				</div>
			</div>
		</div>

		<div v-show="isInit" class="row">
			<div class="col-md-12">
				<div class="card card-secondary card-outline m-0">
					<div class="card-header p-2"><h3 class="card-title text-md">Init Containers</h3></div>
					<div class="card-body p-2">
						<dl v-for="(val, idx) in initContainers" v-bind:key="idx" class="row mb-0 card-body p-2 border-bottom">
							<dt class="col-sm-12"><span class="badge font-weight-light text-sm ml-1" v-bind:class="val.status.badge">{{" "}}</span><span class="card-title mb-2">{{ val.name }}</span></dt>
							<dt v-if="val.status.value" class="col-sm-2 text-truncate">Status</dt><dd v-if="val.status.value" class="col-sm-10" v-bind:class="val.status.style">{{ val.status.value }}{{ (val.status.ready)? `, ${val.status.ready}` : '' }} {{ (val.status.reason.reason) ? `- ${val.status.reason.reason} (exit code: ${val.status.reason.exitCode})` :''}}</dd>
							<dt v-if="val.lastState" class="col-sm-2 text-truncate">Last Status</dt>
							<dd v-if="val.lastState" class="col-sm-10">
								<ul class="list-unstyled mb-0">
									<li v-for="(ls, idx) in val.lastState" v-bind:key="idx">{{ idx }} : {{ ls }}</li>
								</ul>
							</dd>
							<dt class="col-sm-2 text-truncate">Image</dt><dd class="col-sm-10">{{ val.image }}</dd>
							<dt v-if="val.ports" class="col-sm-2 text-truncate">Ports</dt>
							<dd v-if="val.ports" class="col-sm-10">
								<ul class="list-unstyled mb-0">
									<li v-for="(p, idx) in val.ports" v-bind:key="idx">{{p.name? p.name+':' : ""}}{{ p.port }}/{{ p.protocol }}</li>
								</ul>
							</dd>
							<dt class="col-sm-2 text-truncate">Environment</dt>
							<dd class="col-sm-10">
								<ul class="list-unstyled mb-0">
									<li v-for="(e, idx) in val.env" v-bind:key="idx"><span class="font-weight-bold">{{ e.name }}</span>: {{ e.value }} {{ e.v }}</li>
									<li v-if="!val.env">-</li>
								</ul>
							</dd>
							<dt class="col-sm-2 text-truncate">Mounts</dt>
							<dd class="col-sm-10">
								<ul v-for="(m, idx) in val.mounts" v-bind:key="idx" class="list-unstyled mb-0">
									<li><span class="badge badge-secondary font-weight-light text-sm mb-1">{{ m.path }}</span></li>
									<li>from {{ m.name }}({{m.ro}})</li>
								</ul>
							</dd>
							<dt v-if="val.command" class="col-sm-2 text-truncate">Command</dt><dd v-if="val.command" class="col-sm-10 text-sm">{{ val.command }}</dd>
						</dl>
					</div>
				</div>
			</div>
		</div>

		<div class="row">
			<div class="col-md-12">
				<div class="card card-secondary card-outline m-0">
					<div class="card-header p-2"><h3 class="card-title text-md">Containers</h3></div>
					<div class="card-body p-2">
						<dl v-for="(val, idx) in containers" v-bind:key="idx" class="row mb-0 card-body p-2 border-bottom">
							<dt class="col-sm-12">

								<span class="card-title mb-2"><b-badge :variant="val.status.badge" class="mt-1 mb-1 mr-1">&nbsp;</b-badge>{{ val.name }}</span></dt>
							<dt v-if="val.status.value" class="col-sm-2 text-truncate">Status</dt><dd v-if="val.status.value" class="col-sm-10" v-bind:class="val.status.style">{{ val.status.value }}{{ (val.status.ready)? `, ${val.status.ready}` : '' }} {{ (val.status.reason.reason) ? `- ${val.status.reason.reason} (exit code: ${val.status.reason.exitCode})` :''}}</dd>
							<dt v-if="val.lastState" class="col-sm-2 text-truncate">Last Status</dt>
							<dd v-if="val.lastState" class="col-sm-10">
								<ul class="list-unstyled mb-0">
									<li v-for="(ls, idx) in val.lastState" v-bind:key="idx">{{ idx }} : {{ ls }}</li>
								</ul>
							</dd>
							<dt class="col-sm-2 text-truncate">Image</dt><dd class="col-sm-10">{{ val.image }}</dd>
							<dt v-if="val.ports" class="col-sm-2 text-truncate">Ports</dt>
							<dd v-if="val.ports" class="col-sm-10">
								<ul class="list-unstyled mb-0">
									<li v-for="(p, idx) in val.ports" v-bind:key="idx">{{p.name? p.name+':' : ""}}{{ p.port }}/{{ p.protocol }}</li>
								</ul>
							</dd>
							<dt class="col-sm-2 text-truncate">Environment</dt>
							<dd class="col-sm-10">
								<ul class="list-unstyled mb-0">
									<li v-for="(e, idx) in val.env" v-bind:key="idx"><span class="font-weight-bold">{{ e.name }}</span>: {{ e.value }} {{ e.v }}</li>
									<li v-if="!val.env">-</li>
								</ul>
							</dd>
							<dt class="col-sm-2 text-truncate">Mounts</dt>
							<dd class="col-sm-10">
								<ul v-for="(m, idx) in val.mounts" v-bind:key="idx" class="list-unstyled mb-0">
									<li ><span class="badge badge-secondary font-weight-light text-sm mb-1">{{ m.path }}</span></li>
									<li>from {{ m.name }}({{m.ro}})</li>
								</ul>
							</dd>
							<dt v-if="val.command" class="col-sm-2 text-truncate">Command</dt><dd v-if="val.command" class="col-sm-10 text-sm">{{ val.command }}</dd>
							<dt v-if="val.liveness" class="col-sm-2 text-truncate">Liveness</dt><dd v-if="val.liveness" class="col-sm-10"><span v-for="(d, idx) in val.liveness" v-bind:key="idx" class="badge badge-secondary font-weight-light text-sm mb-1 mr-1">{{ d }}</span></dd>
							<dt v-if="val.readiness" class="col-sm-2 text-truncate">Readiness</dt><dd v-if="val.readiness" class="col-sm-10"><span v-for="(d, idx) in val.readiness" v-bind:key="idx" class="badge badge-secondary font-weight-light text-sm mb-1 mr-1">{{ d }}</span></dd>
							<dt v-if="val.startup" class="col-sm-2 text-truncate">Startup</dt><dd v-if="val.startup" class="col-sm-10"><span v-for="(d, idx) in val.startup" v-bind:key="idx" class="badge badge-secondary font-weight-light text-sm mb-1 mr-1">{{ d }}</span></dd>
							<dt v-if="val.args" class="col-sm-2 text-truncate">Arguments</dt><dd v-if="val.args" class="col-sm-10 text-sm"><span v-for="(d, idx) in val.args" v-bind:key="idx">{{ d }} </span></dd>
						</dl>
					</div>
				</div>
			</div>
		</div>

		<div class="row">
			<div class="col-md-12">
				<div class="card card-secondary card-outline m-0">
					<div class="card-header p-2"><h3 class="card-title text-md">volumes</h3></div>
					<div class="card-body p-2">
						<dl v-for="(val, idx) in volumes" v-bind:key="idx" class="row mb-0 card-body p-2 border-bottom">
							<dt class="col-sm-12"><p class="mb-1"><i class="fas fa-hdd mr-1 "></i> {{ val.name }}</p></dt>
							<dt class="col-sm-2 text-truncate">Type</dt><dd class="col-sm-10">{{ val.type }}</dd>
							<dt v-if="val.subName !== ''" class="col-sm-2 text-truncate">{{ val.subName }}</dt><dd v-if="val.subName !== ''" class="col-sm-10">{{ val.subValue }}</dd>
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
import axios			from "axios"
import VueChartJs from "vue-chartjs"
import VueJsonTree from "@/components/jsontree";

export default {
	components: {
		"c-jsontree": { extends: VueJsonTree },
		"c-linechart": {
			extends: VueChartJs.Line,
			props: ["options"],
			mixins: [VueChartJs.mixins.reactiveProp],
			mounted () {
				if(this.chartData) {
					this.renderChart(this.chartData, this.options)
				}
			}
		}
	},
	data() {
		return {
			raw: { metadata: {}, spec: {} },
			event: [],
			metadata: {},
			volumes: [],
			initContainers: [],
			containers: [],
			status: [],
			metrics: [],
			controller: [],
			podInfo: [],
			isCpu: false,
			isMemory: false,
			isInit: false,
			onTols: false,
			onAffis: false,
			chart: {
				options: {
					cpu: {
						maintainAspectRatio : false, responsive : true, legend: { display: false },
						scales: {
							xAxes: [{ gridLines : {display : false}}],
							yAxes: [{ gridLines : {display : false},  ticks: { beginAtZero: true, suggestedMax: 0} }]
						}
					},
					memory: {
						maintainAspectRatio : false, responsive : true, legend: { display: false },
						scales: {
							xAxes: [{ gridLines : {display : false}}],
							yAxes: [{ gridLines : {display : false},  ticks: { beginAtZero: true, suggestedMax: 0} }]
						}
					}
				},
				data: { cpu: {}, memory: {}}
			},
		}
	},
	mounted() {
		this.$nuxt.$on("onReadCompleted", (data) => {
			if(!data) return
			this.metadata = data.metadata;
			this.onCpu(data.spec)
			this.onMemory(data.spec)
			this.onSync(data)
		});
	},
	methods: {
		onSync(data) {
			this.raw = data;
			this.event = this.getEvents(data.metadata.uid);
			this.volumes = this.getVolumes(data.spec.volumes) || {};
			this.containers = this.getContainers(data) || {};
			this.status = this.toStatus(data.metadata.deletionTimestamp, data.status);
			this.controller = this.getController(data.metadata);
			this.podInfo = this.getPodInfo(data);
			this.initContainers = this.getInitContainers(data);
		},
		onCpu(spec) {
			axios.get(`${this.backendUrl()}/api/clusters/${this.currentContext()}/namespaces/${this.metadata.namespace}/pods/${this.metadata.name}/metrics/cpu`)
					.then(resp => {
						if (resp.data.items) {
							let data = resp.data.items[0]
							let labels =[], da= []; let top = 0;
							data.metricPoints.forEach(d => {
								if (d.value>top) top = d.value;
								let dt = new Date(d.timestamp);
								labels.push(`${dt.getHours()}:${dt.getMinutes()}`);
								da.push(d.value/1000);
							});
							let topData = [];
							if (spec.containers) {
								spec.containers.forEach(el => {
									if(el.resources) {
										if(el.resources.requests && el.resources.requests.cpu) topData.push(el.resources.requests.cpu)
									}
								})
							}
							let sum = 0;
							for(let i=0;i<topData.length;i++) {
								if(topData[i].includes('m')) {
									sum += Number(topData[i].slice(0,-1))
								} else {
									sum += topData[i]*1000
								}
							}
							this.isCpu = !!data;
							if (sum) top = sum
							else top = top*1.2
							if (top === 0) top = 1;
							this.$data.chart.options.cpu.scales.yAxes[0].ticks.suggestedMax = (top/1000);
							this.$data.chart.data.cpu = {
								labels: labels,
								datasets: [
									{ backgroundColor : "rgba(119,149,233,0.9)",data: da}
								]
							};
						} else {
							this.isCpu = false;
						}
					})
		},
		onMemory(spec) {
			axios.get(`${this.backendUrl()}/api/clusters/${this.currentContext()}/namespaces/${this.metadata.namespace}/pods/${this.metadata.name}/metrics/memory`)
					.then(resp => {
						if (resp.data.items){
							let data = resp.data.items[0]
							let labels =[], da= []; let top = 0;
							data.metricPoints.forEach(d => {
								if (d.value>top) top = d.value;
								let dt = new Date(d.timestamp);
								labels.push(`${dt.getHours()}:${dt.getMinutes()}`);
								da.push(Math.round(d.value/(1024)));
							});
							let topData = [];
							if (spec.containers) {
								spec.containers.forEach(el => {
									if(el.resources) {
										if(el.resources.requests && el.resources.requests.memory) topData.push(el.resources.requests.memory)
									}
								})
							}
							let sum = 0;
							for(let i=0;i<topData.length;i++) {
								if (topData[i].includes('Gi')){
									sum += Number(topData[i].slice(0,-2))*1024*1024*1024
								} else if(topData[i].includes('Mi')) {
									sum += Number(topData[i].slice(0,-2))*1024*1024
								} else if(topData[i].includes('Ki')){
									sum += Number(topData[i].slice(0,-2))*1024
								} else {
									sum += topData[i]*1024
								}
							}
							this.isMemory = !!data;
							if (sum) top = sum
							else top = top*1.2
							if( top === 0) top = 1;
							this.$data.chart.options.memory.scales.yAxes[0].ticks.suggestedMax = (top/1024);
							this.$data.chart.data.memory = {
								labels: labels,
								datasets: [
									{ backgroundColor : "rgba(179,145,208,1)",data: da}
								]
							};
						} else {
							this.isMemory = false;
						}
					})
		},
		isError(e) {
			this.errorcheck = true;
			this.raw = { metadata: {}, spec: {} };
			this.title = ""
			this.errorMessage = e.response? e.response.data.message: e;
		},
		getPodInfo(d) {
			let podIP = [];
			let conditions = [];
			let tolerations = [];
			let affinity = [];
			let secret = [];
			let nodeSelector = [];
			let isAffinity = false;
			let isNodeSelector = false;
			if(d.status.podIPs) {
				d.status.podIPs.forEach(el => {
					podIP.push(el.ip)
				})
			} else podIP =['-']
			if(d.status.conditions) {
				d.status.conditions.forEach(el =>{
					conditions.push(el.type)
				})
			}
			if(d.spec.tolerations) {
				d.spec.tolerations.forEach(el => {
					tolerations.push({
						key: el.key || '',
						operator: el.operator || '',
						effect: el.effect || '',
						seconds: el.tolerationSeconds || '',
					})
				})
			}
			if(d.spec.affinity) {
				affinity = d.spec.affinity;
				isAffinity = true;
			}
			if(d.spec.volumes) {
				d.spec.volumes.forEach(el => {
					if(el.secret) {
						secret.push(el.secret.secretName)
					}
				})
				if(secret.length === 0) secret = false
			}
			if(d.spec.nodeSelector) {
				let key = Object.keys(d.spec.nodeSelector)
				for(let i=0;i<key.length;i++) {
					nodeSelector.push({
						name: key[i],
						value: d.spec.nodeSelector[key[i]]
					})
				}
				isNodeSelector = true
			}
			return {
				podIP: podIP,
				priorityClass: d.spec.priorityClassName? d.spec.priorityClassName: '-',
				qosClass: d.status.qosClass? d.status.qosClass: '-',
				conditions: conditions,
				nodeSelector: nodeSelector,
				tolerations: tolerations,
				affinities: affinity,
				secret: secret,
				isAffinity: isAffinity,
				isNodeSelector: isNodeSelector,
			}
		},
		getInitContainers(d) {
			let statusCons = []
			let specCons = []
			let statusCon = d.status.initContainerStatuses
			let specCon = d.spec.initContainers
			this.isInit = !!d.spec.initContainers;
			if(statusCon) {
				statusCon.forEach(el => {
					statusCons.push({
						name: el.name,
						status: this.checkStatus(Object.keys(el.state),el),
						lastState: this.getLast(el.lastState),
						image: el.image,
					})
				})
			}
			if(specCon) {
				specCon.forEach(el => {
					specCons.push({
						name: el.name,
						args: el.args,
						image: el.image,
						env: this.getEnv(el.env),
						ports: this.getPorts(el.ports),
						mounts: this.getMounts(el.volumeMounts),
						command: this.getCommand(el.command),
						status: {value:'',style:''},
					})
				})
			}
			if(specCon) {
				for (let i = 0; i < specCon.length; i++) {
					Object.assign(specCons[i], statusCons[i])
				}
				return specCons
			} else if(statusCon) {
				return statusCon
			} else return false
		},
		getEvents(uid) {
			let events = [];
			axios.get(this.getApiUrl('events.k8s.io','events',''))
					.then( resp => {
						for(let i=0; i<resp.data.items.length; i++) {
							if(resp.data.items[i].regarding.uid === uid) {
								events.unshift({
									name: resp.data.items[i].note || "-",
									source: resp.data.items[i].deprecatedSource.host || resp.data.items[i].deprecatedSource.component || "undefined",
									count: resp.data.items[i].deprecatedCount || "-",
									subObject: resp.data.items[i].regarding.fieldPath || "-",
									lastSeen: resp.data.items[i].deprecatedLastTimestamp || "-",
									type: resp.data.items[i].type === "Warning"? "text-danger" : "text-secondary",
								})
							}
						}
					})
			return events
		},
		getVolumes(vol) {
			let vols = []
			if(vol) {
				vol.forEach(d => {
					if (d.persistentVolumeClaim) {
						vols.push({
							name: d.name,
							type: 'persistentVolumeClaim',
							subName: 'claimName',
							subValue: d.persistentVolumeClaim.claimName
						})
					} else if (d.configMap) {
						vols.push({
							name: d.name,
							type: 'configMap',
							subName: 'name',
							subValue: d.configMap.name
						})
					} else {
						vols.push({
							name: d.name,
							type: (Object.keys(d)[1] === 'name'? Object.keys(d)[0] : Object.keys(d)[1]),
							subName: '',
							subValue: '',
						})
					}
				})
				return vols
			}
			return false

		},
		getContainers(d) {
			let statusCons = []
			let specCons = []
			let statusCon = d.status.containerStatuses
			let specCon = d.spec.containers
			if(statusCon) {
				statusCon.forEach(el => {
					statusCons.push({
						name: el.name,
						status: this.checkStatus(Object.keys(el.state),el),
						lastState: this.getLast(el.lastState),
						image: el.image,
					})
				})
			}
			if(specCon) {
				specCon.forEach(el => {
					specCons.push({
						name: el.name,
						args: el.args,
						image: el.image,
						env: this.getEnv(el.env),
						ports: this.getPorts(el.ports),
						mounts: this.getMounts(el.volumeMounts),
						command: this.getCommand(el.command),
						status: {value:'',style:''},
						liveness: this.getProbe(el.livenessProbe),
						readiness: this.getProbe(el.readinessProbe),
						startup: this.getProbe(el.startupProbe),
					})
				})
			}
			if(specCon) {
				for (let i = 0; i < specCon.length; i++) {
					Object.assign(specCons[i], statusCons[i])
				}
				return specCons
			} else if(statusCon) {
				return statusCon
			} else return false
		},
		getEnv(env) {
			let list = []
			if(env) {
				env.forEach(el => {
					if(el.value){
						list.push({name: el.name, value: el.value})
					}else if(el.valueFrom){
						let val = el.valueFrom[Object.keys(el.valueFrom)[0]]
						val = Object.values(val)
						let v = `(${val[1]}.${val[0]})`
						list.push({name: el.name, value: Object.keys(el.valueFrom)[0], v: v })
					}
				})
				return list
			}
			return false
		},
		getPorts(ports) {
			let po = [];
			if(ports) {
				ports.forEach(el => {
					po.push({
						port: el.containerPort,
						name: el.name,
						protocol: el.protocol
					})
				})
				return po
			}
			return false
		},
		getMounts(m) {
			let list =[];
			if(m) {
				m.forEach(el => {
					list.push({
						path: el.mountPath,
						name: el.name,
						ro: (el.readOnly ? 'ro' : 'rw')
					})
				})
				return list
			}
			return false
		},
		getCommand(c) {
			let list = ""
			if(c) {
				c.forEach(el => {
					list += el + " "
				})
				return list
			}
			return false
		},
		toStatus(deletionTimestamp, status) {
			// 삭제
			if (deletionTimestamp) {
				return {
					"value": "Terminating",
					"style": "text-secondary",
				}
			}

			// Pending
			if (!status.containerStatuses) {
				if(status.phase === "Failed") {
					return {
						"value": status.phase,
						"style": "text-danger",
					}
				} else {
					return {
						"value": status.phase,
						"style": "text-warning",
					}
				}
			}

			// [if]: Running, [else]: (CrashRoofBack / Completed / ContainerCreating)
			if(status.containerStatuses.filter(el => el.ready).length === status.containerStatuses.length) {
				const state = Object.keys(status.containerStatuses.find(el => el.ready).state)[0]
				return {
					"value": state.charAt(0).toUpperCase() + state.slice(1),
					"style": "text-success",
				}
			}
			else {
				const state = status.containerStatuses.find(el => !el.ready).state
				let style = "text-secondary"
				if ( state[Object.keys(state)].reason === "Completed") style = "text-success"
				if ( state[Object.keys(state)].reason === "Error") style = "text-danger"
				return {
					"value": state[Object.keys(state)].reason,
					"style": style,
				}
			}
		},
		checkStatus(status,el) {
			status = status[0]
			let reason = this.checkReason(el.state)
			let rd;
			if(el.ready) {
				rd = 'ready';
			} else rd = ''
			if(status === "failed") {
				return {
					"value": status,
					"style": "text-danger",
					"ready": rd,
					"reason": reason,
					'badge': "danger",
				}
			} else if(status === "pending" || status === 'waiting') {
				return {
					"value": status,
					"style": "text-warning",
					"ready": rd,
					"reason": reason,
					'badge': "warning",
				}
			} else if(status === "running" || status === "completed" || status ==="ready") {
				return {
					"value": status,
					"style": "text-success",
					"ready": rd,
					"reason": reason,
					'badge': "success",
				}
			}else {
				return {
					"value": status,
					"style": "text-secondary",
					"ready": rd,
					"reason": reason,
					'badge': "secondary",
				}
			}
		},
		checkReason(state) {
			let key = Object.keys(state)
			return {
				reason: state[key].reason,
				exitCode: state[key].exitCode,
			}
		},
		getLast(s) {
			if(Object.keys(s).length !== 0) {
				return {
					Status: Object.keys(s)[0],
					Reason: s[Object.keys(s)].reason,
					ErrorCode: s[Object.keys(s)].exitCode,
					StartedAt: s[Object.keys(s)].startedAt,
					FinishedAt: s[Object.keys(s)].finishedAt,
				}
			}
			return false
		},
		getController(meta) {
			if (meta.ownerReferences) {
				let or = meta.ownerReferences[0]
				let k = (or.kind).toLowerCase() + 's'
				let g = (or.apiVersion).split('/')
				if (g.length === 2) {
					return {
						g: g[0],
						k: k
					}
				} else {
					return {
						g: '',
						k: k
					}
				}
			}
			return false
		},
		getProbe(probeData) {
			if (!probeData) {
				return false;
			}
			const {
				httpGet, exec, tcpSocket, initialDelaySeconds, timeoutSeconds,
				periodSeconds, successThreshold, failureThreshold
			} = probeData;
			let probe = [];
			if (httpGet) {
				const { path, port, host, scheme } = httpGet;

				probe.push(
						"http-get",
						`${scheme.toLowerCase()}://${host || ""}:${port || ""}${path || ""}`,
				);
			}
			if (exec && exec.command) {
				probe.push(`exec [${exec.command.join(" ")}]`);
			}

			if (tcpSocket && tcpSocket.port) {
				probe.push(`tcp-socket :${tcpSocket.port}`);
			}

			probe.push(
					`delay=${initialDelaySeconds || "0"}s`,
					`timeout=${timeoutSeconds || "0"}s`,
					`period=${periodSeconds || "0"}s`,
					`#success=${successThreshold || "0"}`,
					`#failure=${failureThreshold || "0"}`,
			);
			return probe;
		},
		onTol() {
			this.onTols = !this.onTols
		},
		onAffi() {
			this.onAffis = !this.onAffis
		},
	},
	beforeDestroy(){
		this.$nuxt.$off("onReadCompleted");
	},
}
</script>
