<template>
<div class="row">
	<div class="col-md-12">
		<div class="card card-secondary card-outline">
			<div class="card-body p-2">
				<b-tabs content-class="mt-3" >
					<div v-if="isEmpty" class="text-center"><p> Metrics not available at the moment</p></div>
					<b-tab v-if="!isEmpty" title="CPU" active title-link-class="border-top-0 border-right-0  border-left-0">
						<div class="chart">
							<c-linechart :chart-data="data.cpu" :options="options.cpu" class="mw-100 h-chart"></c-linechart>
						</div>
					</b-tab>
					<b-tab v-if="!isEmpty" title="Memory"  title-link-class="border-top-0 border-right-0  border-left-0">
						<div class="chart">
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
import {CHART_BG_COLOR} from "static/constrants";

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
					maintainAspectRatio : false, responsive : true, legend: { display: false, position: "bottom" },
					scales: {
						xAxes: [{ gridLines : {display : false}}],
						yAxes: [{ gridLines : {display : false},  ticks: { beginAtZero: true, suggestedMax: 0} }]
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
					maintainAspectRatio : false, responsive : true, legend: { display: false, position: "bottom" },
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
		}
	},
	watch: {
		value(newVal) {
			this.$axios.get(`/api/clusters/${this.currentContext()}/${newVal}`)
				.then(resp => {
					this.isEmpty = true
					this.resources = { limits: resp.data.limits, requests: resp.data.requests};

					if (resp.data.metrics) {
						let labels =[], cpus= [], memories = [];
						resp.data.metrics.forEach(d => {
							let dt = new Date(d.timestamp);
							labels.push(`${dt.getHours()}:${dt.getMinutes()}`);
							cpus.push((d.cpu/1000).toFixed(3));
							memories.push(Math.round(d.memory/(1024*1024)));
						});
						this.data.cpu = {
							labels: labels,
							datasets: [
								{ backgroundColor : CHART_BG_COLOR.cpu, label:"Usage", data: cpus }
							]
						};
						this.data.memory = {
							labels: labels,
							datasets: [
								{ backgroundColor : CHART_BG_COLOR.memory, label:"Usage", data: memories }
							]
						};
						this.isEmpty = resp.data.metrics.length==0;

						// limits
						if(resp.data.limits) {
							this.options.cpu.scales.yAxes[0].ticks.suggestedMax = resp.data.limits.cpu/1000;
							this.options.memory.scales.yAxes[0].ticks.suggestedMax = resp.data.limits.memory/(1024*1024);
						}

						// requests
						let legend = {cpu: false, memory: false};
						if(resp.data.requests) {
							legend = {cpu: resp.data.requests.cpu>0, memory: resp.data.requests.memory>0};
							if(resp.data.requests.cpu) {
								let requestCpus = []; 
								for(let i=0; i<cpus.length; i++) requestCpus.push((resp.data.requests.cpu/1000).toFixed(3));
								this.data.cpu.datasets.push({
									backgroundColor: CHART_BG_COLOR.white,
									borderColor: CHART_BG_COLOR.requests,
									label:'Requests',
									pointRadius:0,
									borderWidth:1,
									borderDash: [2, 2],
									data: requestCpus
								});
							}
							if(resp.data.requests.memory) {
								let requestMemories = []; 
								for(let i=0; i<cpus.length; i++) requestMemories.push(resp.data.requests.memory/(1024*1024));
								this.data.memory.datasets.push({
									backgroundColor : CHART_BG_COLOR.white,
									borderColor: CHART_BG_COLOR.limits,
									label:"Limits",
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

		}
	}

}
</script>
