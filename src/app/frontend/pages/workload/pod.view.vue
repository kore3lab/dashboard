<template>
<div>
	<!-- 1. chart -->
	<c-charts v-model="value" class="row"></c-charts>
	<!-- 2. metadata -->
	<c-metadata v-model="value" dtCols="2" ddCols="10" @navigate="$emit('navigate', arguments[0])">
		<dt class="col-sm-2">Status</dt>
		<dd class="col-sm-10" v-bind:class="status.style">{{ status.value }}</dd>
		<dt class="col-sm-2">Node</dt>
		<dd class="col-sm-10">
			<a v-if="nodeName" href="#" @click="$emit('navigate', getViewLink('', 'nodes', '', nodeName))">{{ nodeName }}</a>
			<span v-if="!nodeName">-</span>
		</dd>
		<dt class="col-sm-2">Pod IP</dt>
		<dd class="col-sm-10">
			<ul class="list-unstyled mb-0">
				<li v-if="typeof info.podIP == 'string'" class="mb-1">{{ info.podIP }}</li>
				<li else v-for="(v, k) in info.podIP" v-bind:key="k" class="mb-1">{{ v.ip }}</li>
			</ul>
		</dd>
		<dt class="col-sm-2">Priority Class</dt><dd class="col-sm-10">{{ info.priorityClass}}</dd>
		<dt class="col-sm-2">QoS Class</dt><dd class="col-sm-10">{{ info.qosClass }}</dd>
		<dt v-if="info.conditions" class="col-sm-2">Conditions</dt>
		<dd v-if="info.conditions" class="col-sm-10">
			<span v-for="(v, k) in info.conditions" v-bind:key="k" class="border-box">{{ v.type }}</span>
		</dd>
		<dt v-if="info.secret" class="col-sm-2">Secrets</dt>
		<dd v-if="info.secret" class="col-sm-10" >
			<ul class="list-unstyled mb-0">
				<li v-for="(d, idx) in info.secret" v-bind:key="idx" >
					<a href="#" @click="$emit('navigate', getViewLink('', 'secrets', metadata.namespace, d.secret.secretName))">{{ d.secret.secretName }}</a>
				</li>
			</ul>
		</dd>
	</c-metadata>

	<!-- 2-2. graph -->
	<div class="row">
		<div class="col-md-12">
			<div class="card card-secondary card-outline">
				<div class="card-header p-2"><h3 class="card-title">Relation</h3></div>
				<div class="card-body mw-100" id="wrapRelationGraph"></div>
			</div>
		</div>
	</div>

	<!-- 3. init container -->
	<div class="row" v-show="initContainers.length > 0">
		<div class="col-md-12">
			<div class="card card-secondary card-outline">
				<div class="card-header p-2"><h3 class="card-title">Init Containers</h3></div>
				<div class="card-body group">
					<ul>
						<li v-for="(val, idx) in initContainers" v-bind:key="idx">
							<div class="title"><span v-bind:class=" {'badge-success': (val.status.value ==='running' || val.status.value ==='complete' || val.status.value ==='ready'), 'badge-danger':val.status.value ==='failed', 'badge-secondary':(val.status.value ==='unknown' || val.status.value ==='terminated'),'badge-warning':(val.status.value ==='pending' || val.status.value ==='waiting')}" class="badge mr-1">&nbsp;</span>{{ val.name }}</div>
							<dl class="row">
								<dt v-if="val.status.value" class="col-sm-2">Status</dt><dd v-if="val.status.value" class="col-sm-10" v-bind:class="{'text-success': (val.status.value ==='running' || val.status.value ==='complete' || val.status.value ==='ready'), 'text-danger':val.status.value ==='failed', 'text-warning': (val.status.value ==='pending' || val.status.value ==='waiting'),'text-secondary' :(val.status.value ==='unknown' || val.status.value ==='terminated')}" > {{ val.status.value }}{{ (val.status.ready)? `, ${val.status.ready}` : '' }} {{ (val.status.reason.reason) ? `- ${val.status.reason.reason} (exit code: ${val.status.reason.exitCode})` :''}}</dd>
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
										<li class="font-weight-bold">{{ m.mountPath }}</li>
										<li>from {{ m.name }}({{m.readOnly ? "ro" : "rw"}})</li>
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
							<div class="title">
                				<span v-bind:class=" {'badge-success': (val.status.value ==='running' || val.status.value ==='complete' || val.status.value ==='ready'), 'badge-danger':val.status.value==='failed', 'badge-secondary':(val.status.value ==='unknown' || val.status.value==='terminated'),'badge-warning':(val.status.value ==='pending' || val.status.value ==='waiting')}" class="badge mr-1">&nbsp;</span><span class="pr-3">{{ val.name }}</span>
								<span v-if="val.status.value === 'running'"/>
								<nuxt-link  v-if="val.status.value === 'running'" :to="{path: '/terminal', query: {termtype: 'container',pod: metadata.name, namespace: metadata.namespace, cluster: currentContext(),container:val.name}}" target="_blank">
									<button id="terminal" class="btn pl-0 pr-0 text-sm" v-b-tooltip.hover title="Shell"><b-icon icon="terminal-fill"></b-icon></button>
								</nuxt-link>
								<button class="btn pl-0 text-sm" @click="onClickShowLogs(val.name)" v-b-tooltip.hover title="Logs"><b-icon icon="card-text"></b-icon></button>
							</div>
							<dl class="row">
								<dt v-if="val.status.value" class="col-sm-2">Status</dt><dd v-if="val.status.value" class="col-sm-10" v-bind:class="{'text-success': (val.status.value ==='running' || val.status.value ==='complete' || val.status.value==='ready'), 'text-danger':val.status.value ==='failed', 'text-warning': (val.status.value ==='pending' || val.status.value ==='waiting'),'text-secondary' :(val.status.value ==='unknown' || val.status.value ==='terminated')}" > {{ val.status.value }}{{ (val.status.ready)? `, ${val.status.ready}` : '' }} {{ (val.status.reason.reason) ? `- ${val.status.reason.reason} (exit code: ${val.status.reason.exitCode})` :''}}</dd>
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
										<li v-for="(p, idx) in val.ports" v-bind:key="idx">{{p.name? p.name+':' : ""}}{{ p.containerPort }}/{{ p.protocol }}</li>
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
										<li class="font-weight-bold">{{ m.mountPath }}</li>
										<li>from {{ m.name }}({{m.readOnly ? "ro" : "rw"}})</li>
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
	<c-events v-model="value" class="row"></c-events>

</div>
</template>
<style scoped>
#wrapRelationGraph {min-height: 18em;}
#wrapRelationGraph >>> g.outline.vertical g.tree g.node path.background {fill:#f3f4f5}
</style>
<script>
import VueMetadataView	from "@/components/view/metadataView.vue";
import VueJsonTree		from "@/components/jsontree";
import VueEventsView	from "@/components/view/eventsView.vue";
import VueChartsView	from "@/components/view/metricsChartsView.vue";
import HierarchyGraph	from "@/components/graph/graph.hierarchy";

export default {
	props:["value"],
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
			spec: {},
			volumes: [],
			initContainers: [],
			containers: [],
			status: [],
			metrics: [],
			info: {}
		}
	},
	watch: {
		value(d) { this.onSync(d) }
	},
	methods: {
		onSync(data) {
			if (!data) return
			this.metadata = data.metadata;
			this.nodeName = data.spec.nodeName;
			this.containers = this.getContainers(data) || {};
			this.status = this.toPodStatus(data.metadata.deletionTimestamp, data.status);
			this.info = {
				podIP: data.status.podIPs || data.status.podIP || "-",
				priorityClass: data.spec.priorityClassName || '-',
				qosClass: data.status.qosClass || '-',
				conditions: data.status.conditions || [],
				secret: data.spec.volumes? data.spec.volumes.filter(el=>{return el.secret}): []
			};
			this.initContainers = this.getContainers(data, true);

			//graph
			let g = new HierarchyGraph("#wrapRelationGraph", {
				global: {
					toolbar: { visible:false }
				},
				extends: {
					hierarchy: {
						type: "vertical",
						group: {
							title: { display: "none" },
							box: {
								border: { width: 0 },
								background: { fill: "none" }
							}
						}
					}
				}
			})
			.on("nodeclick", (e,d)=> {
				if (d.data.namespace && d.data.name) {
					// (core) pod, (apps) daemonset, replicaset, deployment
					const v = d.data.apiVersion.split("/");
					const model = this.getViewLink(v.length == 1 ? "": v[0], `${d.data.kind.toLowerCase()}${d.data.kind.endsWith("s") ? "es": "s"}`, d.data.namespace, d.data.name);
					this.$emit("navigate", model)
				}
			})
			this.$axios.get(`/api/clusters/${this.currentContext()}/graph/pod/namespaces/${data.metadata.namespace}/pods/${data.metadata.name}`)
				.then( resp => {
					g.data(resp.data).render();
				})
				.catch(e => { this.msghttp(e);})

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
		},
		onClickShowLogs(name) {
			let containerList = []
			this.containers.forEach(item =>{
				containerList.push(item.name);
			})
			this.$nuxt.$emit("open-terminal", this.metadata.name, "logs", { metadata: this.metadata, container:name, containers:containerList });
			this.$parent.$parent.$emit('close')
		},
		getContainers(d, type) {
			let specCons = []
			let statusCon = type ?  d.status.initContainerStatuses : d.status.containerStatuses;
			let specCon = type ? d.spec.initContainers || d.status.initContainerStatuses || [] : d.spec.containers || d.status.containerStatuses || []
			
			if(specCon) {
				for (let i = 0; i < specCon.length; i++) {
					Object.assign(specCon[i], statusCon[i])
				}

				specCon.forEach((specCon, i) => {
					specCons.push({
						name : specCon.name,
						args: specCon.args,
						image: specCon.image,
						env: this.getEnv(specCon.env),
						ports: specCon.ports,
						mounts: specCon.volumeMounts,
						command: this.getCommand(specCon.command),
						status: this.checkStatus(Object.keys(specCon.state),specCon) || {value:'',style:''},
						lastState: this.getLast(specCon.lastState)
					})
					if(!type){
						specCons[i].liveness = this.getProbe(specCon.livenessProbe),
						specCons[i].readiness = this.getProbe(specCon.readinessProbe),
						specCons[i].startup = this.getProbe(specCon.startupProbe)
					}
				})
			}
			return specCons;
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
			return{
				value : status[0],
				ready : el.ready ? 'ready' : '',
				reason : {
					reason: status[Object.keys(status)].reason,
					exitCode: status[Object.keys(status)].exitCode,
				}
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
	}
}
</script>