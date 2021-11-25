<template>
<div class="row">
	<div class="col-md-12">
		<div class="card card-secondary card-outline">
			<div class="card-body p-2">
				<b-tabs content-class="mt-3" >
					<div v-if="isEmpty" class="text-center"><p> Metrics not available at the moment</p></div>
					<b-tab v-if="!isEmpty" title="CPU" active title-link-class="border-top-0 border-right-0  border-left-0">
						<div class="chart">
							<div class="text-right">
								<span v-if="quotas.cpu.requests>0" class="border-box background">requests : {{ quotas.cpu.requests | formatNumber }}m</span>
								<span v-if="quotas.cpu.limits>0" class="border-box background">limits : {{ quotas.cpu.limits  | formatNumber }}m</span>
								<span v-if="quotas.cpu.allocatable>0" class="border-box background">allocatable : {{ quotas.cpu.allocatable/1000  | formatNumber }}</span>
								<span v-if="quotas.cpu.allocatable==0" class="border-box">milli-cores</span>
								<span v-if="quotas.cpu.allocatable>0" class="border-box">Cores</span>
							</div>
							<c-linechart :chart-data="data.cpu" :options="options.cpu" class="mw-100 h-chart"></c-linechart>
						</div>
					</b-tab>
					<b-tab v-if="!isEmpty" title="Memory"  title-link-class="border-top-0 border-right-0  border-left-0">
						<div class="chart">
							<div class="text-right">
								<span v-if="quotas.memory.requests>0" class="border-box background">requests : {{ quotas.memory.requests/(1024**2) | formatNumber }}Mi</span>
								<span v-if="quotas.memory.limits>0" class="border-box background">limits : {{ quotas.memory.limits/(1024**2) | formatNumber }}Mi</span>
								<span v-if="quotas.memory.allocatable>0" class="border-box background">allocatable : {{ quotas.memory.allocatable/(1024**3) | formatNumber }}Gi</span>
								<span v-if="quotas.memory.allocatable==0" class="border-box">MiB</span>
								<span v-if="quotas.memory.allocatable>0" class="border-box">GiB</span>
							</div>
							<c-linechart :chart-data="data.memory" :options="options.memory" class="mw-100 h-chart"></c-linechart>
						</div>
					</b-tab>
				</b-tabs>
			</div>
		</div>
	</div>
</div>
</template>
<script>
import VueChartJs		from "vue-chartjs";

export default {
	props:["value"],
	components: {
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
	data () {
		return {
			isEmpty: true,
			options: {
				cpu: {
					tooltips: {
						callbacks: {
							label: function(data) {
								return ` ${data.yLabel}m`
							}
						}
					},
					maintainAspectRatio : false, responsive : true, legend: { display: false, position: "bottom" },
					scales: {
						xAxes: [{ gridLines : {display : false}}],
						yAxes: [{ gridLines : {display : true, borderDash: [2, 2]},  ticks: { beginAtZero: true, suggestedMax: 0, callback: this.formatNumber} }]
					}
				},
				memory: {
					tooltips: {
						callbacks: {
							label: function(data) {
								return ` ${(data.yLabel).toFixed(2)}Mi`
							}
						}
					},
					maintainAspectRatio : false, responsive : true, legend: { display: false, position: "bottom" },
					scales: {
						xAxes: [{ gridLines : {display : false}}],
						yAxes: [{
								gridLines : {display : true, borderDash: [2, 2]},
								ticks: { beginAtZero: true, suggestedMax: 0, callback: this.formatNumber}
						}]
					}
				}
			},
			quotas: {
				cpu: { requests:0, limits:0 },
				memory: { requests:0, limits:0 }
			},
			data: { cpu: {}, memory: {}}
		}
	},
	watch: {
		value(val) {
			if(val) this.onSync(val);
		}
	},
	methods: {
		onSync(data) {
			if(!data) return;
			let selectUrl = data.kind=="Node"? `nodes/${data.metadata.name}`: `namespaces/${data.metadata.namespace}/${data.kind.toLowerCase()}s/${data.metadata.name}`;
			this.$axios.get(`/api/clusters/${this.currentContext()}/${selectUrl}/metrics`)
				.then(resp => {
					this.isEmpty = true;
					this.cpu = {limits: 0, requests: 0};
					this.memory = {limits: 0, requests: 0};
					this.quotas = {
						cpu: {
							requests: resp.data.requests?resp.data.requests.cpu:0,
							limits: resp.data.limits?resp.data.limits.cpu:0,
							allocatable: resp.data.allocatable?resp.data.allocatable.cpu:0	//case is 'node'
						},
						memory: {
							requests: resp.data.requests?resp.data.requests.memory:0,
							limits: resp.data.limits?resp.data.limits.memory:0,
							allocatable: resp.data.allocatable?resp.data.allocatable.memory:0	//case is 'node'
						}
					};
					if (resp.data.metrics) {
						let labels =[], cpus= [], memories = [];
						resp.data.metrics.forEach(d => {
							let dt = new Date(d.timestamp);
							labels.push(`${dt.getHours()}:${dt.getMinutes()}`);
							cpus.push(resp.data.allocatable?d.cpu/1000:d.cpu);
							memories.push(Math.round(d.memory/(resp.data.allocatable?1024**3:1024**2)));
						});
						this.data.cpu = {
							labels: labels,
							datasets: [
								{ backgroundColor : this.var("CHART_BG_COLOR").cpu, label:"Usage", data: cpus }
							]
						};
						this.data.memory = {
							labels: labels,
							datasets: [
								{ backgroundColor : this.var("CHART_BG_COLOR").memory, label:"Usage", data: memories }
							]
						};
						this.isEmpty = resp.data.metrics.length==0;

						// allocatable / limits
						if(resp.data.allocatable) {
							this.options.cpu.scales.yAxes[0].ticks.suggestedMax = resp.data.allocatable.cpu/1000 ||  50;
							this.options.memory.scales.yAxes[0].ticks.suggestedMax = resp.data.allocatable.memory/(1024**3);

							this.options.cpu.tooltips.callbacks.label = (data) => {
								return ` ${data.yLabel}`
							}
							this.options.memory.tooltips.callbacks.label = (data) => {
								return ` ${(data.yLabel).toFixed(2)} Gi`
							}

						} else if(resp.data.limits) {
							this.options.cpu.scales.yAxes[0].ticks.suggestedMax = resp.data.limits.cpu ||  50;
							this.options.memory.scales.yAxes[0].ticks.suggestedMax = resp.data.limits.memory/(1024**2);
						}

						// requests
						let legend = {cpu: false, memory: false};
						if(resp.data.requests) {
							legend = {cpu: resp.data.requests.cpu>0, memory: resp.data.requests.memory>0};
							if(resp.data.requests.cpu) {
								let requestCpus = []; 
								for(let i=0; i<cpus.length; i++) requestCpus.push(resp.data.requests.cpu);
								this.data.cpu.datasets.push({
									borderColor : "#999",
									backgroundColor: this.var("CHART_BG_COLOR").requests,
									label:'Requests',
									pointRadius:0,
									borderWidth:1,
									borderDash: [2, 2],
									data: requestCpus
								});
							}
							if(resp.data.requests.memory) {
								let requestMemories = []; 
								for(let i=0; i<cpus.length; i++) requestMemories.push(resp.data.requests.memory/(1024**2));
								this.data.memory.datasets.push({
									borderColor : this.var("CHART_BORDER_COLOR").requests,
									backgroundColor : this.var("CHART_BG_COLOR").requests,
									label:"Requests",
									pointRadius:0,
									borderWidth:1,
									borderDash: [2, 2],
									data: requestMemories
								});
							}
						}
						this.options.cpu.legend.display = legend.cpu;
						this.options.memory.legend.display = legend.memory;

					}
				})
				.catch(e => {
					this.msghttp(e);
				});
		}

	}
}
</script>
