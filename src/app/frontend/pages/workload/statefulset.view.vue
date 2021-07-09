<template>
<div>
	<!-- 1. graph -->
	<div class="row">
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

	<!-- 2. metadata -->
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
						<dt v-if="metadata.ownerReferences" class="col-sm-2">Controlled By</dt>
						<dd v-if="metadata.ownerReferences" class="col-sm-10">{{ metadata.ownerReferences[0].kind }} <a href="#" @click="$emit('navigate', getViewLink(controller.g, controller.k, metadata.namespace, metadata.ownerReferences[0].name))">{{ metadata.ownerReferences[0].name }}</a></dd>
						<dt v-if="info.selector" class="col-sm-2">Selector</dt>
						<dd v-if="info.selector" class="col-sm-10">
							<span v-for="(value, key) in info.selector" v-bind:key="key" class="border-box background">{{key}}={{value}}</span>
						</dd>
						<dt v-if="info.nodeSelector" class="col-sm-2">Node Selector</dt>
						<dd v-if="info.nodeSelector" class="col-sm-10">
							<span v-for="(value, key) in info.nodeSelector" v-bind:key="key" class="border-box background">{{key}}={{value}}</span>
						</dd>
						<dt class="col-sm-2">Images</dt>
						<dd class="col-sm-10">
							<ul class="list-unstyled">
								<li v-for="(val, idx) in info.image" v-bind:key="idx">{{ val }}</li>
							</ul>
						</dd>
						<dt v-if="info.isToleration" class="col-sm-2">Tolerations</dt>
						<dd v-if="info.isToleration" class="col-sm-10">{{ info.tolerations? info.tolerations.length: "-" }}<a class="float-right" v-b-toggle.tol href="#tol-table" @click.prevent @click="onTol">{{onTols ? 'Hide' : 'Show'}}</a></dd>
						<b-collapse class="col-sm-12" id="tol-table"><b-table striped hover small :items="info.tolerations"></b-table></b-collapse>

						<dt v-show="info.isAffinity" class="col-sm-2">Affinities</dt>
						<dd v-show="info.isAffinity" class="col-sm-10">{{ info.affinities? Object.keys(info.affinities).length: "-" }}<a class="float-right" v-b-toggle.affi href="#affi-json" @click.prevent @click="onAffi">{{onAffis ? 'Hide' : 'Show'}}</a>
							<b-collapse id="affi-json"><c-jsontree id="txtSpec" v-model="info.affinities" class="card-body p-2 border"></c-jsontree></b-collapse>
						</dd>
						<dt class="col-sm-2">Pod Status</dt><dd class="col-sm-10"><span v-for="(val,idx) in cs" v-bind:key="idx" v-bind:class="val.style">{{ val.status }} : {{ val.count }}  </span><span v-if="!isStatus">-</span></dd>
						<dt class="col-sm-2">CPU</dt>
						<dd class="col-sm-10">
							<span v-if="totalCpu !== 0" class="border-box">Usage : {{ totalCpu }}</span>
							<span v-if="cpuRequests !== 0" class="border-box">Requests : {{ cpuRequests.toFixed(2) }}</span>
							<span v-if="cpuLimits !== 0" class="border-box">Limits : {{ cpuLimits.toFixed(2) }}</span>
							<span v-if="totalCpu === 0 && cpuRequests === 0 && cpuLimits === 0">-</span>
						</dd>
						<dt class="col-sm-2">Memory</dt>
						<dd class="col-sm-10">
							<span v-if="totalMemory !== 0" class="border-box">Usage : {{ totalMemory | formatNumber }} Mi</span>
							<span v-if="memoryRequests !== 0" class="border-box">Requests : {{ memoryRequests.toFixed(1) | formatNumber }} Mi</span>
							<span v-if="memoryLimits !== 0" class="border-box">Limits : {{ memoryLimits.toFixed(1) | formatNumber }} Mi</span>
							<span v-if="totalMemory === 0 && memoryRequests === 0 && memoryLimits === 0">-</span>
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

	<!-- 3. events -->
	<div class="row" v-show="event.length>0">
		<div class="col-md-12">
			<div class="card card-secondary card-outline">
				<div class="card-header p-2"><h3 class="card-title">Events</h3></div>
				<div class="card-body p-2">
					<dl v-for="(val, idx) in event" v-bind:key="idx" class="row mb-0 card-body p-2 border-bottom">
						<dt class="col-sm-12"><p v-bind:class="val.type" class="mb-1">{{ val.name }}</p></dt>
						<dt class="col-sm-2">Source</dt><dd class="col-sm-10">{{ val.source }}</dd>
						<dt class="col-sm-2">Count</dt><dd class="col-sm-10">{{ val.count }}</dd>
						<dt class="col-sm-2">Sub-object</dt><dd class="col-sm-10">{{ val.subObject }}</dd>
						<dt class="col-sm-2">Last seen</dt><dd class="col-sm-10">{{ val.lastSeen }}</dd>
					</dl>
				</div>
			</div>
		</div>
	</div>
</div>
</template>
<script>
import VueJsonTree from "@/components/jsontree";
import VueChartJs from "vue-chartjs";
import {CHART_BG_COLOR} from "static/constrants";

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
			controller: [],
			temp: [],
			cs: [],
			sumCpu: [],
			sumMemory: [],
			nowCpu: [],
			topCpu: [],
			topMemory: [],
			nowMemory: {},
			totalCpu: 0,
			totalMemory: 0,
			cpuRequests: 0,
			cpuLimits: 0,
			memoryRequests: 0,
			memoryLimits: 0,
			isCpu: false,
			isMemory: false,
			isStatus: false,
			onTols: false,
			onAffis: false,
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
						maintainAspectRatio : false, responsive : true, legend: { display: true, position: 'bottom' },
						scales: {
							xAxes: [{ gridLines : {display : false}}],
							yAxes: [{ gridLines : {display : false},  ticks: { beginAtZero: true, suggestedMax: 0, callback: function(value) {return value.toFixed(3)}} }]
						}
					},
					memory: {
						tooltips: {
							callbacks: {
								label: function(data) {
									return (data.yLabel).toFixed(2) + "Mi"
								}
							}
						},
						maintainAspectRatio : false, responsive : true, legend: { display: true, position: 'bottom' },
						scales: {
							xAxes: [{ gridLines : {display : false}}],
							yAxes: [{ gridLines : {display : false},  ticks: { beginAtZero: true, suggestedMax: 0, callback: function(value) {
										if(value === 0) return value
										let regexp = /\B(?=(\d{3})+(?!\d))/g;
										return value.toString().replace(regexp, ',')+'Mi';}
								}}]
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
			this.controller = this.getController(data.metadata.ownerReferences)
			this.info = this.getInfo(data);
			this.event = this.getEvents(data.metadata.uid);
			this.childPod = this.getChildPod(data.status.currentRevision);
		},
		getInfo(data) {
			let tolerations = [];
			let affinity = [];
			let image = [];
			let isToleration = false;
			let isAffinity = false;

			if(data.spec.template.spec.containers) {
				data.spec.template.spec.containers.forEach(el =>{
					image.push(el.image)
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
				selector: data.spec.selector.matchLabels || '',
				nodeSelector: data.spec.template.spec.nodeSelector || '',
				image: image,
				tolerations: tolerations,
				affinities: affinity,
				isAffinity: isAffinity,
				isToleration: isToleration,
			}
		},
		getChildPod(currentRevision) {
			let childPod = [];
			this.nowMemory = [];
			this.nowCpu = [];
			this.sumCpu = [];
			this.sumMemory = [];
			this.temp = [];
			this.cs = [];
			this.topCpu = [];
			this.topMemory = [];
			this.cpuRequests = 0
			this.cpuLimits = 0
			this.memoryRequests = 0
			this.memoryLimits = 0
			this.isStatus = false;
			this.isPods = false;
			this.isCpu = false;
			this.isMemory = false;
			this.$axios.get(this.getApiUrl('','pods',this.metadata.namespace,'','labelSelector=controller-revision-hash=' + currentRevision))
					.then( resp => {
						let idx = 0;
						resp.data.items.forEach(el =>{
							this.isPods = true;
							childPod.push({
								name: el.metadata.name,
								namespace: el.metadata.namespace,
								ready: this.toReady(el.status,el.spec),
								nowMemory: this.onMemory(el,idx),
								nowCpu: this.onCpu(el,idx),
								status: this.toStatus(el.metadata.deletionTimestamp, el.status),
								countStatus: this.countStatus(el.status),
								idx: idx,
							})
							idx++;
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
							this.cpuRequests += this.cpuRL(el.resources.requests.cpu)
						}
					}
					if(el.resources.limits && el.resources.limits.cpu) {
						this.cpuLimits += this.cpuRL(el.resources.limits.cpu)
					}
				})
			}
			this.$axios.get(`/api/clusters/${this.currentContext()}/namespaces/${el.metadata.namespace}/pods/${el.metadata.name}/metrics/cpu`)
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
			let re=[], li=[];
			if(value)
			{
				for (let i=0;i<value.length;i++) {
					if(value[i].time in map) {
						map[value[i].time] += value[i].val;
					} else {
						map[value[i].time] = value[i].val;
					}
					if(top<map[value[i].time]) top = map[value[i].time]
				}
				let keys = Object.keys(map)
				for (let i=0;i<keys.length;i++) {
					da[i] = map[keys[i]]/1000
					if(this.cpuRequests) re[i] = this.cpuRequests
					if(this.cpuLimits) li[i] = this.cpuLimits
				}
				if(this.cpuLimits > 0) {
					top = this.cpuLimits
				} else if ( this.cpuRequests > 0) {
					top = this.cpuRequests
				} else {
					top = top*1.2 / 1000
				}
				if(top === 0) top = 1
				this.isCpu = !!value
				this.$data.chart.options.cpu.scales.yAxes[0].ticks.suggestedMax = top;
				this.$data.chart.data.cpu = {
					labels: labels,
					datasets: [
						{ backgroundColor : CHART_BG_COLOR.cpu,data: da,label:'Usage'},
					]
				}
				if(this.cpuRequests) this.$data.chart.data.cpu.datasets.push({ backgroundColor: CHART_BG_COLOR.white,data: re, borderColor: CHART_BG_COLOR.requests,label:'Requests',pointRadius:0,borderWidth:1})
				if(this.cpuLimits) this.$data.chart.data.cpu.datasets.push({ backgroundColor : CHART_BG_COLOR.white,data: li, borderColor: CHART_BG_COLOR.limits,label:"Limits",pointRadius:0,borderWidth:1})
			}
		},
		onMemory(el,idx) {
			if(el.spec.containers) {
				el.spec.containers.forEach(el => {
					if(el.resources && el.resources.requests) {
						if(el.resources.requests.memory) {
							this.topMemory.push(el.resources.requests.memory)
							this.memoryRequests += this.memoryRL(el.resources.requests.memory)
						}
					}
					if(el.resources.limits && el.resources.limits.memory) {
						this.memoryLimits += this.memoryRL(el.resources.limits.memory)
					}
				})
			}
			this.$axios.get(`/api/clusters/${this.currentContext()}/namespaces/${el.metadata.namespace}/pods/${el.metadata.name}/metrics/memory`)
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
			let re=[], li=[];
			if(value) {
				for (let i = 0; i < value.length; i++) {
					if (value[i].time in map) {
						map[value[i].time] += value[i].val;
					} else {
						map[value[i].time] = value[i].val;
					}
					if (top < map[value[i].time]) top = map[value[i].time] / 1024 / 1024
				}
				let keys = Object.keys(map)
				for (let i = 0; i < keys.length; i++) {
					da[i] = map[keys[i]] / 1024 / 1024
					if(this.memoryRequests) re[i] = this.memoryRequests
					if(this.memoryLimits) li[i] = this.memoryLimits
				}
				if(this.memoryLimits > 0) {
					top = this.memoryLimits
				} else if ( this.memoryRequests > 0) {
					top = this.memoryRequests
				} else {
					top = top*1.2
				}
				if (top === 0) top = 1024
				this.isMemory = !!value
				this.$data.chart.options.memory.scales.yAxes[0].ticks.suggestedMax = top;
				this.$data.chart.data.memory = {
					labels: labels,
					datasets: [
						{backgroundColor: CHART_BG_COLOR.memory, data: da, label:'Usage'},
					]
				}
				if(this.memoryRequests) this.$data.chart.data.memory.datasets.push({ backgroundColor: CHART_BG_COLOR.white,data: re, borderColor: CHART_BG_COLOR.requests,label:'Requests',pointRadius:0,borderWidth:1})
				if(this.memoryLimits) this.$data.chart.data.memory.datasets.push({ backgroundColor : CHART_BG_COLOR.white,data: li, borderColor: CHART_BG_COLOR.limits,label:"Limits",pointRadius:0,borderWidth:1})
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
		countStatus(status) {
			this.isStatus = true;
			let style;
			this.cs = [];
			if(status.phase in this.temp) {
				this.temp[status.phase] += 1;
			} else {
				this.temp[status.phase] = 1;
			}
			let key = Object.keys(this.temp)
			key.forEach(el =>{
				if(el === 'Running') {
					style = 'text-success'
				} else if (el === 'Pending') {
					style = 'text-warning'
				} else if (el === 'Failed') {
					style = 'text-danger'
				} else style = 'text-secondary'
				this.cs.push({
					status: el,
					count: this.temp[el],
					style: style,
				})
			})
			return this.cs
		},
		onTol() {
			this.onTols = !this.onTols
		},
		onAffi() {
			this.info.affinities = Object.assign({},this.info.affinities)
			this.onAffis = !this.onAffis
		},
	},
	beforeDestroy(){
		this.$nuxt.$off("onReadCompleted");
	},
}
</script>
