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
							<dt class="col-sm-2">Annotations</dt>
							<dd class="col-sm-10 text-truncate">
								<ul class="list-unstyled mb-0">
									<li v-for="(value, name) in metadata.annotations" v-bind:key="name"><span class="badge badge-secondary font-weight-light text-sm mb-1">{{ name }}:{{ value }}</span></li>
								</ul>
							</dd>
							<dt class="col-sm-2">Labels</dt>
							<dd class="col-sm-10 text-truncate">
								<ul class="list-unstyled mb-0">
									<li v-for="(value, name) in metadata.labels" v-bind:key="name"><span class="badge badge-secondary font-weight-light text-sm mb-1">{{ name }}:{{ value }}</span></li>
								</ul>
							</dd>
							<dt class="col-sm-2">UID</dt><dd class="col-sm-10">{{ metadata.uid }}</dd>
							<dt class="col-sm-2">Replicas</dt><dd class="col-sm-10">{{ info.replicas }}</dd>
							<dt v-if="info.selector" class="col-sm-2">Selector</dt><dd v-if="info.selector" class="col-sm-10"><span v-for="(value, key) in info.selector" v-bind:key="key" class="badge badge-secondary font-weight-light text-sm mb-1 mr-1">{{key}}={{value}}</span></dd>
							<dt v-if="info.nodeSelector" class="col-sm-2">Node Selector</dt><dd v-if="info.nodeSelector" class="col-sm-10"><span v-for="(value, key) in info.nodeSelector" v-bind:key="key" class="badge badge-secondary font-weight-light text-sm mb-1 mr-1">{{key}}={{value}}</span></dd>
							<dt class="col-sm-2">Strategy Type</dt><dd class="col-sm-10">{{ info.strategyType }}</dd>
							<dt class="col-sm-2">Conditions</dt><dd class="col-sm-10"><span v-for="(d, idx) in info.conditions" v-bind:key="idx" v-bind:class="d.style" class="badge font-weight-light text-sm mb-1 mr-1">{{ d.type }}</span></dd>

							<dt v-if="info.isToleration" class="col-sm-2 text-truncate">Tolerations</dt>
							<dd v-if="info.isToleration" class="col-sm-10">{{ info.tolerations? info.tolerations.length: "-" }}<a class="float-right" v-b-toggle.tol href="#tol-table" @click.prevent @click="onTol">{{onTols ? 'Hide' : 'Show'}}</a></dd>
							<b-collapse class="col-sm-12" id="tol-table"><b-table striped hover small :items="info.tolerations"></b-table></b-collapse>

							<dt v-show="info.isAffinity" class="col-sm-2 text-truncate">Affinities</dt>
							<dd v-show="info.isAffinity" class="col-sm-10">{{ info.affinities? Object.keys(info.affinities).length: "-" }}<a class="float-right" v-b-toggle.affi href="#affi-json" @click.prevent @click="onAffi">{{onAffis ? 'Hide' : 'Show'}}</a>
								<b-collapse id="affi-json"><c-jsontree id="txtSpec" v-model="info.affinities" class="card-body p-2 border"></c-jsontree></b-collapse>
							</dd>
							<dt class="col-sm-2">CPU</dt><dd class="col-sm-10"><span v-if="totalCpu !== 0" class="badge badge-secondary font-weight-light text-sm mb-1">Usage : {{ totalCpu }}</span><span v-if="totalCpu === 0">-</span></dd>
							<dt class="col-sm-2">Memory</dt><dd class="col-sm-10"><span v-if="totalMemory !== 0" class="badge badge-secondary font-weight-light text-sm mb-1">Usage : {{ totalMemory }} Mi</span><span v-if="totalMemory === 0">-</span></dd>
						</dl>
					</div>
				</div>
			</div>
		</div>
		<div class="row">
			<div class="col-md-12">
				<div class="card card-secondary card-outline">
					<div class="card-header p-2"><h3 class="card-title text-md ">Pods</h3></div>
					<div class="card-body p-2">
						<b-table striped hover small :items="childPod" :fields="fields">
							<template v-slot:cell(name)="data">
								<a href="#" @click="$emit('navigate', getViewLink('', 'pods', data.item.namespace, data.item.name))">{{ data.item.name }}</a>
							</template>
							<template v-slot:cell(status)="data">
								<span v-bind:class="data.item.status.style">{{ data.item.status.value }}</span>
							</template>
							<template v-slot:cell(nowCpu)="data">
								<span v-if="data.item.nowCpu[data.item.idx]">{{ data.item.nowCpu[data.item.idx].val ? data.item.nowCpu[data.item.idx].val : '' }}</span>
							</template>
							<template v-slot:cell(nowMemory)="data">
								<span v-if="data.item.nowMemory[data.item.idx]">{{ data.item.nowMemory[data.item.idx].val ? data.item.nowMemory[data.item.idx].val+'Mi' : ''}}</span>
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
			sumCpu: [],
			sumMemory: [],
			nowCpu: [],
			topCpu: [],
			topMemory: [],
			nowMemory: {},
			totalCpu: 0,
			totalMemory: 0,

			isCpu: false,
			isMemory: false,
			onTols: false,
			onAffis: false,
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
			this.onSync(data)
		});
		this.$nuxt.$emit("onCreated",'')
	},
	methods: {
		onSync(data) {
			this.totalCpu = 0; this.totalMemory = 0; this.cpus = {};
			this.info = this.getInfo(data);
			this.event = this.getEvents(data.metadata.uid);
			this.childPod = this.getChildPod(data.spec.template.metadata.labels);
		},
		getInfo(data) {
			let replicas = data.spec.replicas +' desired, ' + (data.status.updatedReplicas || 0) + ' updated, ' + (data.status.replicas || 0) +' total, ' + (data.status.availableReplicas || 0) + ' available, ' + (data.status.unavailableReplicas || 0) + ' unavailable'
			let conditions = [];
			let tolerations = [];
			let affinity = [];
			let isToleration = false;
			let isAffinity = false;
			if (data.status.conditions) {
				data.status.conditions.forEach(el => {
					conditions.push({
						type: el.type,
						style: this.checkStyle(el.type),
					})
				})
				conditions.sort(function(a,b) {
					return a.type < b.type ? -1 : a.type > b.type ? 1 : 0;
				})
			}
			if(data.spec.template.spec.tolerations) {
				data.spec.template.spec.tolerations.forEach(el =>{
					tolerations.push({
						key: el.key || '',
						operator: el.operator || '',
						effect: el.effect || '',
						seconds: el.tolerationSeconds || '',
					})
					isToleration = true;
				})
			}
			if(data.spec.template.spec.affinity && Object.keys(data.spec.template.spec.affinity).length !== 0) {
				affinity = data.spec.template.spec.affinity;
				isAffinity = true;
			}
			return {
				replicas: replicas,
				selector: data.spec.selector.matchLabels || '',
				nodeSelector: data.spec.template.spec.nodeSelector || '',
				strategyType: data.spec.strategy.type,
				conditions: conditions,
				tolerations: tolerations,
				affinities: affinity,
				isAffinity: isAffinity,
				isToleration: isToleration,
			}
		},
		getChildPod(label) {
			label = this.stringifyLabels(label)
			let childPod = [];
			this.nowMemory = [];
			this.nowCpu = [];
			this.sumCpu = [];
			this.topCpu = [];
			this.topMemory = [];
			this.sumMemory = [];
			axios.get(this.getApiUrl('','pods',this.metadata.namespace))
					.then( resp => {
						let childLabel
						let count = 0, idx = 0;
						resp.data.items.forEach(el =>{
							childLabel= this.stringifyLabels(el.metadata.labels)
							label.forEach(e => {
								if(childLabel.indexOf(e) === -1) {
									count++;
								}
							})
							if(count === 0) {
								childPod.push({
									name: el.metadata.name,
									namespace: el.metadata.namespace,
									ready: this.toReady(el.status,el.spec),
									nowMemory: this.onMemory(el,idx),
									nowCpu: this.onCpu(el,idx),
									status: this.toStatus(el.metadata.deletionTimestamp, el.status),
									idx: idx,
								})
								idx++;
							}
							count = 0;
						})
					})
			return childPod
		},

		onCpu(el,idx) {
			if(el.spec.containers) {
				el.spec.containers.forEach(el => {
					if(el.resources && el.resources.requests) {
						if(el.resources.requests.cpu) {
							this.topCpu.push(el.resources.requests.cpu)
						}
					}
				})
			}
			axios.get(`${this.backendUrl()}/api/clusters/${this.currentContext()}/namespaces/${el.metadata.namespace}/pods/${el.metadata.name}/metrics/cpu`)
					.then(resp => {
						if (resp.data.items) {
							let data = resp.data.items[0];
							let labels=[], da= [];
							data.metricPoints.forEach(d => {
								let dt = new Date(d.timestamp);
								labels.push(`${dt.getHours()}:${dt.getMinutes()}`);
								da.push(d.value/1000);
								this.sumCpu.push({
									val: d.value,
									time: dt,
								})
							});
							this.nowCpu.push({
								val:((da[da.length-1])).toFixed(3),
								idx:idx,
							})
							this.getCpuGraph(this.sumCpu,labels)
							this.cpuSum(this.sumCpu)
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
		getCpuGraph(value,labels) {
			let top = 0; let da= [];
			let map = {};
			for (let i=0;i<value.length;i++) {
				if(value[i].time in map) {
					map[value[i].time] += value[i].val;
					if(top<map[value[i].time]) top = map[value[i].time]
				} else {
					map[value[i].time] = value[i].val;
				}
			}
			let keys = Object.keys(map)
			for (let i=0;i<keys.length;i++) {
				da[i] = map[keys[i]]/1000
			}
			let sum = 0;
			for(let i=0;i<this.topCpu.length;i++) {
				if(this.topCpu[i].includes('m')) {
					sum += Number(this.topCpu[i].slice(0,-1))
				} else {
					sum += this.topCpu[i]*1000
				}
			}
			this.isCpu = !!value
			if (sum) top = sum
			else top = top*1.2
			if (top === 0) top = 1
			this.$data.chart.options.cpu.scales.yAxes[0].ticks.suggestedMax = top/1000;
			this.$data.chart.data.cpu = {
				labels: labels,
				datasets: [
					{ backgroundColor : "rgba(119,149,233,0.9)",data: da}
				]
			}
		},
		onMemory(el,idx) {
			if(el.spec.containers) {
				el.spec.containers.forEach(el => {
					if(el.resources && el.resources.requests) {
						if(el.resources.requests.memory) {
							this.topMemory.push(el.resources.requests.memory)
						}
					}
				})
			}
			axios.get(`${this.backendUrl()}/api/clusters/${this.currentContext()}/namespaces/${el.metadata.namespace}/pods/${el.metadata.name}/metrics/memory`)
					.then(resp => {
						if (resp.data.items) {
							let data = resp.data.items[0];
							let labels=[], da= [];
							data.metricPoints.forEach(d => {
								let dt = new Date(d.timestamp);
								labels.push(`${dt.getHours()}:${dt.getMinutes()}`);
								da.push(Math.round(d.value/1024));
								this.sumMemory.push({
									val: d.value,
									time: dt,
								})
							});
							this.nowMemory.push({
								val:((da[da.length-1])/1024).toFixed(2),
								idx:idx,
							})
							this.getMemoryGraph(this.sumMemory,labels)
							this.memorySum(this.sumMemory)
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
		getMemoryGraph(value, labels) {
			let top = 0; let da= [];
			let map = {};
			if(value) {
				for (let i = 0; i < value.length; i++) {
					if (value[i].time in map) {
						map[value[i].time] += value[i].val;
						if (top < map[value[i].time]) top = map[value[i].time] * 2
					} else {
						map[value[i].time] = value[i].val;
					}
				}
				let keys = Object.keys(map)
				for (let i = 0; i < keys.length; i++) {
					da[i] = map[keys[i]] / 1024
				}
				let sum = 0;
				for(let i=0;i<this.topMemory.length;i++) {
					if (this.topMemory[i].includes('Gi')){
						sum += Number(this.topMemory[i].slice(0,-2))*1024*1024*1024
					} else if(this.topMemory[i].includes('Mi')) {
						sum += Number(this.topMemory[i].slice(0,-2))*1024*1024
					} else if(this.topMemory[i].includes('Ki')){
						sum += Number(this.topMemory[i].slice(0,-2))*1024
					} else {
						sum += this.topMemory[i]*1024
					}
				}
				this.isMemory = !!value
				if (sum) top = sum
				else top = top*1.2
				if (top === 0) top = 1
				this.$data.chart.options.memory.scales.yAxes[0].ticks.suggestedMax = top / 1024;
				this.$data.chart.data.memory = {
					labels: labels,
					datasets: [
						{backgroundColor: "rgba(179,145,208,1)", data: da}
					]
				}
			}
		},
		cpuSum(cpu) {
			if(cpu && cpu[cpu.length -1].val) {
				this.totalCpu += (cpu[cpu.length - 1].val / 1000)
				this.totalCpu = Number(this.totalCpu.toFixed(3))
			}
		},
		memorySum(memory) {
			if(memory && memory[memory.length -1].val) {
				this.totalMemory += (memory[memory.length - 1].val / (1024 * 1024))
				this.totalMemory = Number(this.totalMemory.toFixed(2))
			}
		},
		checkStyle(t) {
			if(t === 'Progressing') return 'badge-primary'
			if(t === 'Available') return 'badge-success'
			else return 'badge-danger'
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
