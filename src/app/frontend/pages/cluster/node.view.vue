<template>
	<div class="card-body p-2">
		<div class="row mb-0">
			<div class="col-md-12">
				<div class="card card-secondary card-outline">
					<div class="card-body p-2">
						<b-tabs content-class="mt-3" >
							<b-tab title="CPU" active title-link-class="border-top-0 border-right-0  border-left-0">
								<div v-if="isCpu" class="chart">
									<c-linechart id="cpu" :chart-data="chart.data.cpu" :options="chart.options.cpu" class="mw-100 h-chart"></c-linechart>
								</div>
								<div v-if="!isCpu" class="text-center"><p> Metrics not available at the moment</p></div>
							</b-tab>
							<b-tab title="Memory"  title-link-class="border-top-0 border-right-0  border-left-0">
								<div v-if="isMemory" class="chart">
									<c-linechart id="memory" :chart-data="chart.data.memory" :options="chart.options.memory" class="mw-100 h-chart"></c-linechart>
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
							<dt class="col-sm-3 text-truncate">Create at</dt><dd class="col-sm-9">{{ this.getTimestampString(metadata.creationTimestamp)}} ago ({{ metadata.creationTimestamp }})</dd>
							<dt class="col-sm-3 text-truncate">Name</dt><dd class="col-sm-9">{{ metadata.name }}</dd>
							<dt class="col-sm-3 text-truncate">Capacity</dt><dd class="col-sm-9">{{ info.capacity }}</dd>
							<dt class="col-sm-3 text-truncate">Allocatable</dt><dd class="col-sm-9">{{ info.allocatable }}</dd>
							<dt class="col-sm-3 text-truncate">Addresses</dt>
							<dd class="col-sm-9">
								<ul class="list-unstyled mb-0">
									<li v-for="(val, idx) in info.addresses" v-bind:key="idx">{{ val.type }}: {{ val.address }}</li>
								</ul>
							</dd>
							<dt class="col-sm-3 text-truncate">OS</dt><dd class="col-sm-9">{{ info.os }}</dd>
							<dt class="col-sm-3 text-truncate">OS Image</dt><dd class="col-sm-9">{{ info.osImage }}</dd>
							<dt class="col-sm-3 text-truncate">Kernel version</dt><dd class="col-sm-9">{{ info.kernelVersion }}</dd>
							<dt class="col-sm-3 text-truncate">Container runtime</dt><dd class="col-sm-9">{{ info.containerRuntime }}</dd>
							<dt class="col-sm-3 text-truncate">Kubelet version</dt><dd class="col-sm-9">{{ info.kubeletVersion }}</dd>
							<dt class="col-sm-3 text-truncate">UID</dt><dd class="col-sm-9">{{ metadata.uid }}</dd>
							<dt class="col-sm-3 text-truncate">Annotations</dt>
							<dd class="col-sm-9 text-truncate">
								<ul class="list-unstyled mb-0">
									<li v-for="(value, name) in metadata.annotations" v-bind:key="name"><span class="badge badge-secondary font-weight-light text-sm mb-1">{{ name }}:{{ value }}</span></li>
								</ul>
							</dd>
							<dt class="col-sm-3">Labels</dt>
							<dd class="col-sm-9 text-truncate">
								<ul class="list-unstyled mb-0">
									<li v-for="(value, name) in metadata.labels" v-bind:key="name"><span class="badge badge-secondary font-weight-light text-sm mb-1">{{ name }}:{{ value }}</span></li>
								</ul>
							</dd>
							<dt v-if="info.taints" class="col-sm-3 text-truncate">Taints</dt><dd v-if="info.taints" class="col-sm-9"><span v-for="(val, idx) in info.taints" v-bind:key="idx" v-bind:class="val.style" class="badge badge-secondary font-weight-light text-sm mb-1 mr-1">{{ val.key }}: {{ val.effect }} ({{ val.value }})</span></dd>
							<dt class="col-sm-3 text-truncate">Conditions</dt><dd class="col-sm-9"><span v-for="(val, idx) in info.conditions" v-bind:key="idx" v-bind:class="val.style" class="badge font-weight-light text-sm mb-1 mr-1">{{ val.type }}</span></dd>
						</dl>
					</div>
				</div>
			</div>
		</div>
		<div class="row">
			<div class="col-md-12">
				<div class="card card-secondary card-outline">
					<div class="card-header p-2"><h3 class="card-title text-md ">Pods</h3></div>
					<div class="card-body p-2 overflow-auto">
						<b-table striped hover small :items="childPod" :fields="fields" class="text-truncate">
							<template v-slot:cell(name)="data">
								<a href="#" @click="$emit('navigate', getViewLink('', 'pods', data.item.namespace, data.item.name))">{{ data.item.name }}</a>
							</template>
							<template v-slot:cell(status)="data">
								<span v-bind:class="data.item.status.style">{{ data.item.status.value }}</span>
							</template>
							<template v-slot:cell(nowCpu)="data">
								<span v-if="data.item.nowCpu[data.item.idx]">{{ data.item.nowCpu[data.item.idx].val ? (data.item.nowCpu[data.item.idx].val*100/maxCpu).toFixed(2)+'%' : '' }}</span>
							</template>
							<template v-slot:cell(nowMemory)="data">
								<span v-if="data.item.nowMemory[data.item.idx]">{{ data.item.nowMemory[data.item.idx].val ? (data.item.nowMemory[data.item.idx].val*100/maxMemory).toFixed(2)+'%' : ''}}</span>
							</template>
						</b-table>
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
import VueJsonTree from "@/components/jsontree";
import VueChartJs from "vue-chartjs";

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
					this.update()
				}
			},
			watch: {
				chartData: function () {
					this.update();
				},
				options: function() {
					this.update();
				},
			},
			methods: {
				update: function() {
					this.renderChart(this.chartData, this.options);
				},
			},
		}
	},
	data() {
		return {
			metadata: {},
			info: [],
			event: [],
			childPod: [],
			nowCpu: [],
			nowMemory: [],
			maxCpu: 0,
			maxMemory: 0,
			isCpu: false,
			isMemory: false,
			isStatus: false,
			isPods: false,
			fields: [
				{ key: "name", label: "Name", sortable: true },
				{ key: "namespace", label: "Namespace", sortable: true  },
				{ key: "ready", label: "Ready", sortable: true  },
				{ key: "nowCpu", label: "CPU" },
				{ key: "nowMemory", label: "Memory" },
				{ key: "status", label: "Status" },
			],
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
							yAxes: [{ gridLines : {display : false},  ticks: { beginAtZero: true, suggestedMax: 0, callback: function(value) {return value + 'Mi'}} }]
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
			this.onSync(data)
			this.onCpu(data.status.allocatable.cpu)
			this.onMemory(data.status.allocatable.memory)
		});
		this.$nuxt.$emit("onCreated",'')
	},
	methods: {
		onSync(data) {
			this.totalCpu = 0; this.totalMemory = 0; this.cpus = {};
			this.maxCpu = data.status.capacity.cpu;
			this.maxMemory = this.tranMemory(data.status.capacity.memory)/(1024);
			this.info = this.getInfo(data);
			this.event = this.getEvents(data.metadata.uid);
			this.childPod = this.getChildPod();
		},
		getInfo(data) {
			let capacity = `CPU: ${data.status.capacity.cpu}, Memory: ${(this.tranMemory(data.status.capacity.memory)/(1024*1024)).toFixed(2)}Mi, Pods: ${data.status.capacity.pods}`
			let allocatable = `CPU: ${data.status.allocatable.cpu}, Memory: ${(this.tranMemory(data.status.allocatable.memory)/(1024*1024)).toFixed(2)}Mi, Pods: ${data.status.allocatable.pods}`
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
		onCpu(top) {
			axios.get(`${this.backendUrl()}/api/clusters/${this.currentContext()}/nodes/${this.metadata.name}/metrics/cpu`)
					.then(resp => {
						if (resp.data.items) {
							let data = resp.data.items[0]
							let labels =[], da= [];
							data.metricPoints.forEach(d => {
								let dt = new Date(d.timestamp);
								labels.push(`${dt.getHours()}:${dt.getMinutes()}`);
								da.push(d.value/1000);
							});
							this.isCpu = !!data;
							this.$data.chart.options.cpu.scales.yAxes[0].ticks.suggestedMax = top;
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
		onMemory(top) {
			axios.get(`${this.backendUrl()}/api/clusters/${this.currentContext()}/nodes/${this.metadata.name}/metrics/memory`)
					.then(resp => {
						if (resp.data.items){
							let data = resp.data.items[0]
							let labels =[], da= [];
							data.metricPoints.forEach(d => {
								let dt = new Date(d.timestamp);
								labels.push(`${dt.getHours()}:${dt.getMinutes()}`);
								da.push(Math.round(d.value/(1024*1024)));
							});
							top = this.tranMemory(top)
							this.isMemory = !!data;
							this.$data.chart.options.memory.scales.yAxes[0].ticks.suggestedMax = (top/(1024*1024));
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
			axios.get(this.getApiUrl('','pods'))
					.then( resp => {
						let idx = 0;
						resp.data.items.forEach(el => {
							if(el.spec.nodeName === this.metadata.name) {
								childPod.push({
									name: el.metadata.name,
									namespace: el.metadata.namespace,
									ready: this.toReady(el.status,el.spec),
									status: this.toStatus(el.metadata.deletionTimestamp, el.status),
									nowCpu: this.getPodCpu(el,idx),
									nowMemory: this.getPodMemory(el,idx),
									idx: idx,
								})
								idx++;
							}
						})
					})
			return childPod
		},
		getPodCpu(el,idx) {
			axios.get(`${this.backendUrl()}/api/clusters/${this.currentContext()}/namespaces/${el.metadata.namespace}/pods/${el.metadata.name}/metrics/cpu`)
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
			axios.get(`${this.backendUrl()}/api/clusters/${this.currentContext()}/namespaces/${el.metadata.namespace}/pods/${el.metadata.name}/metrics/memory`)
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
