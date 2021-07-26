<template>
<div>
	<!-- 1. chart -->
	<c-charts class="row" v-model="chartsUrl"></c-charts>
	<!-- 2. metadata -->
	<c-metadata v-model="metadata" dtCols="2" ddCols="10">
		<dt class="col-sm-2">Controlled By</dt>
		<dd v-if="metadata.ownerReferences" class="col-sm-10">{{ controller.kind }} <a href="#" @click="$emit('navigate', getViewLink(controller.group, controller.resource, metadata.namespace, controller.name))">{{ controller.name }}</a></dd>
		<dt class="col-sm-2">Status</dt><dd class="col-sm-10" v-bind:class="status.style">{{ status.value }}</dd>
		<dt class="col-sm-2">Node</dt>
		<dd v-if="nodeName" class="col-sm-10"><a href="#" @click="$emit('navigate', getViewLink('', 'nodes', '', nodeName))">{{ nodeName}}</a></dd>
		<dd v-if="!nodeName" class="col-sm-10">-</dd>
		<dt class="col-sm-2">Pod IP</dt>
		<dd class="col-sm-10">
			<ul class="list-unstyled mb-0">
				<li v-for="(v, k) in info.podIP" v-bind:key="k" class="mb-1">{{ v.ip }}</li>
			</ul>
		</dd>
		<dt class="col-sm-2">Priority Class</dt><dd class="col-sm-10">{{ info.priorityClass}}</dd>
		<dt class="col-sm-2">QoS Class</dt><dd class="col-sm-10">{{ info.qosClass }}</dd>
		<dt v-if="info.conditions" class="col-sm-2">Conditions</dt>
		<dd v-if="info.conditions" class="col-sm-10">
			<span v-for="(v, k) in info.conditions" v-bind:key="k" class="border-box">{{ v.type }}</span>
		</dd>
		<dt v-if="info.nodeSelector" class="col-sm-2">Node Selector</dt>
		<dd v-if="info.nodeSelector" class="col-sm-10">
			<span v-for="(v, k) in info.nodeSelector" v-bind:key="k" class="border-box  background">{{ k }}: {{ v }}</span>
		</dd>
		<dt class="col-sm-2">Tolerations</dt>
		<dd class="col-sm-10">{{ info.tolerations? info.tolerations.length: "-" }}
			<a href="#" class="float-right " @click="isTolerations=!isTolerations">{{isTolerations?'Hide':'Show'}}</a>
		</dd>
		<b-collapse class="col-sm-12" v-model="isTolerations">
			<b-table-lite :items="info.tolerations" class="subset"></b-table-lite>
		</b-collapse>

		<dt v-show="Object.keys(affinities).length>0" class="col-sm-2">Affinities</dt>
		<dd v-show="Object.keys(affinities).length>0" class="col-sm-10">{{ Object.keys(affinities).length }}
			<a href="#" class="float-right " @click="isAffinities=!isAffinities">{{isAffinities?'Hide':'Show'}}</a>
			<b-collapse v-model="isAffinities">
				<c-jsontree v-model="affinities" class="card-body p-2 border"></c-jsontree>
			</b-collapse>
		</dd>
		<dt v-if="info.secret" class="col-sm-2">Secrets</dt>
		<dd v-if="info.secret" class="col-sm-10" >
			<ul class="list-unstyled">
				<li v-for="(d, idx) in info.secret" v-bind:key="idx" >
					<a href="#" @click="$emit('navigate', getViewLink('', 'secrets', metadata.namespace, d.secret.secretName))">{{ d.secret.secretName }}</a>
				</li>
			</ul>
		</dd>
	</c-metadata>

	<!-- 3. init container -->
	<div class="row" v-show="initContainers.length > 0">
		<div class="col-md-12">
			<div class="card card-secondary card-outline">
				<div class="card-header p-2"><h3 class="card-title">Init Containers</h3></div>
				<div class="card-body group">
					<ul>
						<li v-for="(val, idx) in initContainers" v-bind:key="idx">
							<div class="title"><b-badge :variant="val.status.badge" class="mr-1">&nbsp;</b-badge><span>{{ val.name }}</span></div>
							<dl class="row">
								<dt v-if="val.status.value" class="col-sm-2">Status</dt><dd v-if="val.status.value" class="col-sm-10" v-bind:class="val.status.style">{{ val.status.value }}{{ (val.status.ready)? `, ${val.status.ready}` : '' }} {{ (val.status.reason.reason) ? `- ${val.status.reason.reason} (exit code: ${val.status.reason.exitCode})` :''}}</dd>
								<dt v-if="val.lastState" class="col-sm-2">Last Status</dt>
								<dd v-if="val.lastState" class="col-sm-10">
									<ul class="list-unstyled mb-0">
										<li v-for="(ls, idx) in val.lastState" v-bind:key="idx">{{ idx }} : {{ ls }}</li>
									</ul>
								</dd>
								<dt class="col-sm-2">Image</dt>
								<dd class="col-sm-10"><span class="border border-secondary rounded pl-2 pr-2" id="copyTextInit">{{ val.image }}</span><button type="button" class="btn p-0 pl-2" @click="copy('init')"><i class="fas fa-copy"></i></button></dd>
								<dt v-if="val.ports" class="col-sm-2">Ports</dt>
								<dd v-if="val.ports" class="col-sm-10">
									<ul class="list-unstyled mb-0">
										<li v-for="(p, idx) in val.ports" v-bind:key="idx">{{p.name? p.name+':' : ""}}{{ p.port }}/{{ p.protocol }}</li>
									</ul>
								</dd>
								<dt class="col-sm-2">Environment</dt>
								<dd class="col-sm-10">
									<ul class="list-unstyled mb-0">
										<li v-for="(e, idx) in val.env" v-bind:key="idx"><span class="font-weight-bold">{{ e.name }}</span>: {{ e.value }} {{ e.v }}</li>
										<li v-if="!val.env">-</li>
									</ul>
								</dd>
								<dt class="col-sm-2">Mounts</dt>
								<dd class="col-sm-10">
									<ul v-for="(m, idx) in val.mounts" v-bind:key="idx" class="list-unstyled mb-0">
										<li style="font-weight-bold">{{ m.path }}</li>
										<li>from {{ m.name }}({{m.ro}})</li>
									</ul>
								</dd>
								<dt v-if="val.command" class="col-sm-2">Command</dt><dd v-if="val.command" class="col-sm-10">{{ val.command }}</dd>
							</dl>
						</li>
					</ul>
				</div>
			</div>
		</div>
	</div>

	<!-- 4. containers -->
	<div class="row">
		<div class="col-md-12">
			<div class="card card-secondary card-outline">
				<div class="card-header p-2"><h3 class="card-title">Containers</h3></div>
				<div class="card-body group">
					<ul>
						<li v-for="(val, idx) in containers" v-bind:key="idx">
							<div class="title"><b-badge :variant="val.status.badge" class="mr-1">&nbsp;</b-badge>{{ val.name }}</div>
							<dl class="row">
								<dt v-if="val.status.value" class="col-sm-2">Status</dt><dd v-if="val.status.value" class="col-sm-10" v-bind:class="val.status.style">{{ val.status.value }}{{ (val.status.ready)? `, ${val.status.ready}` : '' }} {{ (val.status.reason.reason) ? `- ${val.status.reason.reason} (exit code: ${val.status.reason.exitCode})` :''}}</dd>
								<dt v-if="val.lastState" class="col-sm-2">Last Status</dt>
								<dd v-if="val.lastState" class="col-sm-10">
									<ul class="list-unstyled mb-0">
										<li v-for="(ls, idx) in val.lastState" v-bind:key="idx">{{ idx }} : {{ ls }}</li>
									</ul>
								</dd>
								<dt class="col-sm-2">Image</dt><dd class="col-sm-10"><span id="copyTextCon">{{ val.image }}</span><button type="button" class="btn p-0 pl-2" @click="copy('con')"><i class="fas fa-copy"></i></button></dd>
								<dt v-if="val.ports" class="col-sm-2">Ports</dt>
								<dd v-if="val.ports" class="col-sm-10">
									<ul class="list-unstyled mb-0">
										<li v-for="(p, idx) in val.ports" v-bind:key="idx">{{p.name? p.name+':' : ""}}{{ p.port }}/{{ p.protocol }}</li>
									</ul>
								</dd>
								<dt class="col-sm-2">Environment</dt>
								<dd class="col-sm-10">
									<ul class="list-unstyled mb-0">
										<li v-for="(e, idx) in val.env" v-bind:key="idx"><span class="font-weight-bold">{{ e.name }}</span>: {{ e.value }} {{ e.v }}</li>
										<li v-if="!val.env">-</li>
									</ul>
								</dd>
								<dt class="col-sm-2">Mounts</dt>
								<dd class="col-sm-10">
									<ul v-for="(m, idx) in val.mounts" v-bind:key="idx" class="list-unstyled mb-0">
										<li class="font-weight-bold">{{ m.path }}</li>
										<li>from {{ m.name }}({{m.ro}})</li>
									</ul>
								</dd>
								<dt v-if="val.command" class="col-sm-2">Command</dt><dd v-if="val.command" class="col-sm-10">{{ val.command }}</dd>
								<dt v-if="val.liveness" class="col-sm-2">Liveness</dt><dd v-if="val.liveness" class="col-sm-10">
									<span v-for="(d, idx) in val.liveness" v-bind:key="idx" class="border-box">{{ d }}</span>
								</dd>
								<dt v-if="val.readiness" class="col-sm-2">Readiness</dt>
								<dd v-if="val.readiness" class="col-sm-10">
									<span v-for="(d, idx) in val.readiness" v-bind:key="idx" class="border-box">{{ d }}</span>
								</dd>
								<dt v-if="val.startup" class="col-sm-2">Startup</dt>
								<dd v-if="val.startup" class="col-sm-10">
									<span v-for="(d, idx) in val.startup" v-bind:key="idx" class="border-box">{{ d }}</span>
								</dd>
								<dt v-if="val.args" class="col-sm-2">Arguments</dt>
								<dd v-if="val.args" class="col-sm-10">
									<span v-for="(d, idx) in val.args" v-bind:key="idx">{{ d }} </span>
								</dd>
							</dl>
						</li>
					</ul>
				</div>
			</div>
		</div>
	</div>

	<!-- 5. volumnes -->
	<div class="row" v-show="volumes.length > 0">
		<div class="col-md-12">
			<div class="card card-secondary card-outline">
				<div class="card-header p-2"><h3 class="card-title">volumes</h3></div>
				<div class="card-body group">
					<ul>
						<li v-for="(val, idx) in volumes" v-bind:key="idx">
							<div class="title"><i class="fas fa-hdd mr-1 "></i> {{ val.name }}</div>
							<dl class="row">
								<dt class="col-sm-2">Type</dt>
								<dd class="col-sm-10">{{ val.type }}</dd>
								<dt v-if="val.subName" class="col-sm-2">{{ val.subName }}</dt>
								<dd v-if="val.subName" class="col-sm-10"><a href="#" @click="$emit('navigate', getViewLink('', val.type.toLowerCase()+'s', metadata.namespace, val.subValue))">{{ val.subValue }}</a></dd>
							</dl>
						</li>
					</ul>
				</div>
			</div>
		</div>
	</div>

	<!-- 6. events -->
	<c-events class="row" v-model="metadata.uid"></c-events>

</div>
</template>
<script>
import VueMetadataView	from "@/components/view/metadataView.vue";
import VueJsonTree		from "@/components/jsontree";
import VueEventsView	from "@/components/view/eventsView.vue";
import VueChartsView	from "@/components/view/metricsChartsView.vue";


export default {
	components: {
		"c-metadata": { extends: VueMetadataView },
		"c-jsontree": { extends: VueJsonTree },
		"c-events": { extends: VueEventsView },
		"c-charts": { extends: VueChartsView }
	},
	data() {
		return {
			nodeName: "",
			metadata: {},
			chartsUrl: "",
			volumes: [],
			initContainers: [],
			containers: [],
			status: [],
			metrics: [],
			controller: {},
			info: [],
			affinities: {},
			isTolerations: false,
			isAffinities: false
		}
	},
	mounted() {
		this.$nuxt.$on("onReadCompleted", (data) => {
			if (!data) return
			this.metadata = data.metadata;
			this.chartsUrl = `namespaces/${data.metadata.namespace}/pods/${data.metadata.name}`;
			this.nodeName = data.spec.nodeName;
			this.containers = this.getContainers(data) || {};
			this.status = this.toPodStatus(data.metadata.deletionTimestamp, data.status);
			this.controller = data.metadata.ownerReferences?this.getResource(data.metadata.ownerReferences[0]):{};
			this.info = {
				podIP: data.status.podIPs || ["-"],
				priorityClass: data.spec.priorityClassName || '-',
				qosClass: data.status.qosClass || '-',
				conditions: data.status.conditions || [],
				nodeSelector: data.spec.nodeSelector,
				tolerations: data.spec.tolerations || [],
				secret: data.spec.volumes? data.spec.volumes.filter(el=>{return el.secret}): []
			};
			this.affinities = data.spec.affinity || {},
			this.initContainers = this.getInitContainers(data);

			// volumns
			this.volumes = [];
			if(data.spec.volumes) {
				data.spec.volumes.forEach(d => {
					if (d.persistentVolumeClaim) {
						this.volumes.push({name: d.name, type: "persistentVolumeClaim", subName: "claimName", subValue: d.persistentVolumeClaim.claimName});
					} else if (d.configMap) {
						this.volumes.push({name: d.name, type: "configMap", subName: "name", subValue: d.configMap.name});
					} else {
						this.volumes.push({name: d.name, type: (Object.keys(d)[1] === "name"? Object.keys(d)[0] : Object.keys(d)[1]), subName: "",subValue: ""});
					}
				})
			}

		});
		this.$nuxt.$emit("onCreated","")
	},
	methods: {
		getInitContainers(d) {
			let statusCons = []
			let specCons = []
			let statusCon = d.status.initContainerStatuses
			let specCon = d.spec.initContainers
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
				this.$nuxt.$emit('Containers',specCons)
				return specCons
			} else if(statusCon) {
				return statusCon
			} else return false
		},
		copy(type) {
			let copyText
			if(type === 'con') {
				copyText = document.getElementById("copyTextCon").textContent;
			} else {
				copyText = document.getElementById("copyTextInit").textContent;
			}
			const textArea = document.createElement('textarea');
			textArea.textContent = 'docker pull '+ copyText;
			document.body.append(textArea);
			textArea.select();
			document.execCommand("copy")
			textArea.remove()
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
					ExitCode: s[Object.keys(s)].exitCode,
					StartedAt: s[Object.keys(s)].startedAt,
					FinishedAt: s[Object.keys(s)].finishedAt,
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
		}
	},
	beforeDestroy(){
		this.$nuxt.$off("onReadCompleted");
	},
}
</script>