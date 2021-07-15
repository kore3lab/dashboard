<template>
<div>
	<!-- 1. charts -->
	<c-charts class="row" v-model="chartsUrl"></c-charts>
	<!-- 2. metadata -->
	<div class="row">
		<div class="col-md-12">
			<div class="card card-secondary card-outline">
				<div class="card-body p-2">
					<dl class="row mb-0">
						<dt class="col-sm-3">Create at</dt><dd class="col-sm-9">{{ this.getTimestampString(metadata.creationTimestamp)}} ago ({{ metadata.creationTimestamp }})</dd>
						<dt class="col-sm-3">Name</dt><dd class="col-sm-9">{{ metadata.name }}</dd>
						<dt class="col-sm-3">Capacity</dt><dd class="col-sm-9">{{ info.capacity }}</dd>
						<dt class="col-sm-3">Allocatable</dt><dd class="col-sm-9">{{ info.allocatable }}</dd>
						<dt class="col-sm-3">Addresses</dt>
						<dd class="col-sm-9">
							<ul class="list-unstyled mb-0">
								<li v-for="(val, idx) in info.addresses" v-bind:key="idx">{{ val.type }}: {{ val.address }}</li>
							</ul>
						</dd>
						<dt class="col-sm-3">OS</dt><dd class="col-sm-9">{{ info.os }}</dd>
						<dt class="col-sm-3">OS Image</dt><dd class="col-sm-9">{{ info.osImage }}</dd>
						<dt class="col-sm-3">Kernel version</dt><dd class="col-sm-9">{{ info.kernelVersion }}</dd>
						<dt class="col-sm-3">Container runtime</dt><dd class="col-sm-9">{{ info.containerRuntime }}</dd>
						<dt class="col-sm-3">Kubelet version</dt><dd class="col-sm-9">{{ info.kubeletVersion }}</dd>
						<dt class="col-sm-3">UID</dt><dd class="col-sm-9">{{ metadata.uid }}</dd>
						<dt class="col-sm-3">Annotations</dt>
						<dd class="col-sm-9">
							<ul class="list-unstyled mb-0">
								<li v-for="(value, name) in metadata.annotations" v-bind:key="name">{{ name }}=<span class="font-weight-light">{{ value }}</span></li>
							</ul>
						</dd>
						<dt class="col-sm-3">Labels</dt>
						<dd class="col-sm-9">
							<span v-for="(value, name) in metadata.labels" v-bind:key="name" class="label">{{ name }}={{ value }}</span>
						</dd>
						<dt v-if="info.taints" class="col-sm-3">Taints</dt>
						<dd v-if="info.taints" class="col-sm-9">
							<span v-for="(val, idx) in info.taints" v-bind:key="idx" v-bind:class="val.style">{{ val.key }}: {{ val.effect }} ({{ val.value }})</span>
						</dd>
						<dt class="col-sm-3">Conditions</dt>
						<dd class="col-sm-9">
							<span v-for="(val, idx) in info.conditions" v-bind:key="idx" v-bind:class="val.style" class="badge font-weight-light text-sm mb-1 mr-1">{{ val.type }}</span>
						</dd>
					</dl>
				</div>
			</div>
		</div>
	</div>
	<!-- 3. pods -->
	<div class="row">
		<div class="col-md-12">
			<div class="card card-secondary card-outline">
				<div class="card-header p-2"><h3 class="card-title">Pods</h3></div>
				<div class="card-body p-2 overflow-auto">
					<b-table striped hover small :items="childPod" :fields="fields">
						<template v-slot:cell(name)="data">
							<a href="#" @click="$emit('navigate', getViewLink('', 'pods', data.item.namespace, data.item.name))">{{ data.item.name }}</a>
						</template>
						<template v-slot:cell(status)="data">
							<span v-bind:class="data.item.status.style">{{ data.item.status.value }}</span>
						</template>
						<template v-slot:cell(nowCpu)="data">
							<span v-if="data.value[data.item.idx]">{{ data.value[data.item.idx].val ? (data.value[data.item.idx].val*100/maxCpu).toFixed(2)+'%' : '' }}</span>
						</template>
						<template v-slot:cell(nowMemory)="data">
							<span v-if="data.item.nowMemory[data.item.idx]">{{ data.value[data.item.idx].val ? (data.value[data.item.idx].val*100/maxMemory).toFixed(2)+'%' : ''}}</span>
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
import VueJsonTree		from "@/components/jsontree";
import VueEventsView	from "@/components/view/eventsView.vue";
import VueChartsView	from "@/components/view/metricsChartsView.vue";

export default {
	components: {
		"c-jsontree": { extends: VueJsonTree },
		"c-events": { extends: VueEventsView },
		"c-charts": { extends: VueChartsView }
	},
	data() {
		return {
			metadata: {},
			info: [],
			chartsUrl: "",
			childPod: [],
			nowCpu: [],
			nowMemory: [],
			maxCpu: 0,
			maxMemory: 0,
			fields: [
				{ key: "name", label: "Name", sortable: true },
				{ key: "namespace", label: "Namespace", sortable: true  },
				{ key: "ready", label: "Ready", sortable: true  },
				{ key: "nowCpu", label: "CPU" },
				{ key: "nowMemory", label: "Memory" },
				{ key: "status", label: "Status" },
			],
		}
	},
	mounted() {
		this.$nuxt.$on("onReadCompleted", (data) => {
			if(!data) return
			this.metadata = data.metadata;
			this.onSync(data)
			this.chartsUrl = `nodes/${data.metadata.name}/metrics`;
		});
		this.$nuxt.$emit("onCreated",'')
	},
	methods: {
		onSync(data) {
			this.totalCpu = 0; this.totalMemory = 0; this.cpus = {};
			this.maxCpu = data.status.capacity.cpu;
			this.maxMemory = this.tranMemory(data.status.capacity.memory)/(1024);
			this.info = this.getInfo(data);
			this.childPod = this.getChildPod();
		},
		getInfo(data) {
			let regexp = /\B(?=(\d{3})+(?!\d))/g;
			let capacity = `CPU: ${data.status.capacity.cpu}, Memory: ${(this.tranMemory(data.status.capacity.memory)/(1024*1024)).toFixed(2).replace(regexp, ',')+'Mi'}, Pods: ${data.status.capacity.pods}`
			let allocatable = `CPU: ${data.status.allocatable.cpu}, Memory: ${(this.tranMemory(data.status.allocatable.memory)/(1024*1024)).toFixed(2).replace(regexp, ',')+'Mi'}, Pods: ${data.status.allocatable.pods}`
			let addresses = this.getAddress(data.status.addresses);
			let conditions = this.getConditions(data.status.conditions);
			let taints = this.getTaints(data.spec.taints);

			return {
				capacity: capacity,
				allocatable: allocatable,
				addresses: addresses,
				os: `${data.status.nodeInfo.operatingSystem} (${data.status.nodeInfo.architecture})`,
				osImage: data.status.nodeInfo.osImage,
				kernelVersion: data.status.nodeInfo.kernelVersion,
				containerRuntime: data.status.nodeInfo.containerRuntimeVersion,
				kubeletVersion: data.status.nodeInfo.kubeletVersion,
				conditions: conditions,
				taints: taints,
			}
		},
		getTaints(t) {
			let list = [];
			if(t) {
				t.forEach(el => {
					list.push({
						effect: el.effect,
						key: el.key,
						value: el.value
					})
				})
				return list
			}
		},
		getConditions(con) {
			let list = [];
			let style = ''
			con.forEach(el => {
				if(el.status === 'True'){
					if(el.type === 'Ready') {
						style = 'badge-success'
					} else if(el.type === 'Pending') {
						style = 'badge-warning'
					} else if(el.type === 'Failed') {
						style = 'badge-danger'
					} else {
						style = 'badge-secondary'
					}
					list.push({
						type: el.type,
						style: style,
					})
				}
			})
			return list
		},
		getAddress(address) {
			let list =[];
			if(address) {
				address.forEach(el => {
					list.push({
						type: el.type,
						address: el.address,
					})
				})
				return list
			}
		},
		tranMemory(memory) {
			let mem;
			if (memory.includes('Gi')){
				mem = Number(memory.slice(0,-2))*1024*1024*1024
			} else if(memory.includes('Mi')) {
				mem = Number(memory.slice(0,-2))*1024*1024
			} else if(memory.includes('Ki')){
				mem = Number(memory.slice(0,-2))*1024
			} else {
				mem = memory*1024
			}
			return mem
		},
		getChildPod() {
			let childPod = [];
			this.nowCpu = [];
			this.nowMemory = [];
			this.$axios.get(this.getApiUrl('','pods','','','fieldSelector=spec.nodeName=' + this.metadata.name))
					.then( resp => {
						let idx = 0;
						resp.data.items.forEach(el => {
							childPod.push({
								name: el.metadata.name,
								namespace: el.metadata.namespace,
								ready: this.toReady(el.status,el.spec),
								status: this.toStatus(el.metadata.deletionTimestamp, el.status),
								nowCpu: 0,
								nowMemory: 0,
								//nowCpu: this.getPodCpu(el,idx),
								//nowMemory: this.getPodMemory(el,idx),
								idx: idx,
							})
							idx++;
						})
					})
			return childPod
		},
		getPodCpu(el,idx) {
			this.$axios.get(`/api/clusters/${this.currentContext()}/namespaces/${el.metadata.namespace}/pods/${el.metadata.name}/metrics/cpu`)
					.then(resp =>{
						if(resp.data.items) {
							let data = resp.data.items[0];
							let da = [];
							data.metricPoints.forEach(d => {
								da.push(d.value/1000);
							})
							this.nowCpu.push({
								val: ((da[da.length-1]).toFixed(3)),
								idx:idx,
							})
						} else {
							this.nowCpu.push({
								val: '',
								idx: idx,
							})
						}
						this.nowCpu = this.sorted(this.nowCpu)
					})
			return this.nowCpu
		},
		getPodMemory(el,idx) {
			this.$axios.get(`/api/clusters/${this.currentContext()}/namespaces/${el.metadata.namespace}/pods/${el.metadata.name}/metrics/memory`)
					.then(resp =>{
						if(resp.data.items) {
							let data = resp.data.items[0];
							let da = [];
							data.metricPoints.forEach(d => {
								da.push(d.value/1024);
							})
							this.nowMemory.push({
								val: ((da[da.length-1])),
								idx:idx,
							})
						} else {
							this.nowMemory.push({
								val: '',
								idx: idx,
							})
						}
						this.nowMemory = this.sorted(this.nowMemory)
					})
			return this.nowMemory
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
	},
	beforeDestroy(){
		this.$nuxt.$off("onReadCompleted");
	},
}
</script>
