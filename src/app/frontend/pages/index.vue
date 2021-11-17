<template>
<div class="content-wrapper">
	<section class="content pt-3">
		<div class="container-fluid">
			<div class="row">
				<div class="col-lg-3 col-6">
					<div class="small-box bg-warning">
						<div class="inner"><h3>{{ Object.keys(nodes).length }}</h3><p>Nodes</p></div>
						<div class="icon"><i class="fas fa-server"></i></div>
					</div>
				</div>
				<div class="col-lg-3 col-6">
					<div class="small-box bg-success">
						<div class="inner"><h3>{{ Math.round(summary.percent["cpu"],0) }}<small>%</small></h3><p>CPU</p></div>
						<div class="icon"><i class="fas fa-microchip"></i></div>
					</div>
				</div>
				<div class="col-lg-3 col-6">
					<div class="small-box bg-info">
						<div class="inner"><h3>{{ Math.round(summary.percent["memory"],0) }}<small>%</small></h3><p>Memory</p></div>
						<div class="icon"><i class="fas fa-memory"></i></div>
					</div>
				</div>
				<div class="col-lg-3 col-6">
					<div class="small-box bg-secondary">
						<div class="inner"><h3>{{ Math.round(summary.percent["storage"],0) }}<small>%</small></h3><p>Storage</p></div>
						<div class="icon"><i class="fas fa-hdd"></i></div>
					</div>
				</div>
			</div>
			<div class="row">
				<div class="col-lg-7">
					<!-- Nodes -->
					<div class="card">
						<div class="card-header">
							<div class="d-flex justify-content-between">
								<h3 class="card-title">Nodes</h3>
								<div class="card-tools">
									<nuxt-link to="/cluster/node.list" class="btn-md text-info"><i class="fas fa-arrow-circle-right"></i></nuxt-link>
								</div>
							</div>
						</div>
						<div class="card-body p-0">
							<b-table-simple v-show="!isBusy"  small responsive borderless >
								<colgroup><col width="40%"><col><col><col><col><col><col><col><col><col></colgroup>
								<b-thead><b-tr><b-th>Name</b-th><b-th>CPU</b-th><b-th>Memory</b-th><b-th>Storage</b-th><b-th>Pods</b-th></b-tr></b-thead>
								<b-tbody>
									<b-tr v-for="(nd, key) in nodes" :key="key">
										<b-td>
											<span v-bind:class="{ 'bg-success': (nd.status=='Ready'), 'bg-warning': (nd.status!='Ready')}" class="badge font-weight-light">{{ nd.status }}</span>
											<span class="pl-1 text-lg">{{ key }}</span>
											<p class="text-muted text-sm font-weight-light ml-1">{{ nd.role }}</p>
										</b-td>
										<b-td>
											<p class="text-lg">{{ nd.metrics.percent.cpu }}<small>%</small></p>
											<p class="text-muted text-sm font-weight-light">{{ nd.metrics.usage.cpu/1000 | formatNumber }}/{{ nd.metrics.allocatable.cpu/1000  | formatNumber }} Core</p>
										</b-td>
										<b-td >
											<p class="text-lg">{{ nd.metrics.percent.memory }}<small>%</small></p>
											<p class="text-muted text-sm font-weight-light">{{ (nd.metrics.usage.memory/(1024**3)).toFixed(2) | formatNumber }}/{{ Math.round(nd.metrics.allocatable.memory/(1024**3),2)  | formatNumber }} GiB</p>
										</b-td>
										<b-td>
											<p class="text-lg">{{ nd.metrics.percent.storage }}<small>%</small></p>
											<p class="text-muted text-sm font-weight-light">{{ (nd.metrics.usage.storage/(1024**3)).toFixed(2) | formatNumber }}/{{ Math.round(nd.metrics.allocatable.storage/(1024**3),2) | formatNumber }} GiB</p>
										</b-td>
										<b-td>
											<p class="text-lg">{{ nd.metrics.percent.pods }}<small>%</small></p>
											<p class="text-muted text-sm font-weight-light">{{ nd.metrics.usage.pods | formatNumber }}/{{ nd.metrics.allocatable.pods | formatNumber }} ea</p>
										</b-td> 
									</b-tr>
								</b-tbody>
							</b-table-simple>
								<div v-show="isBusy" class="text-center text-success" style="line-height: 10rem;">
									<b-spinner type="grow" variant="success" class="align-middle mr-2"></b-spinner>
									<span class="text-lg align-middle">Loading...</span>
								</div>
						</div>
					</div><!-- / Nodes -->
				</div>
				<div class="col-lg-5">
					<!-- Nodes Map -->
					<div class="card">
						<div class="card-header">
							<div class="d-flex justify-content-between">
								<h3 class="card-title">Nodes Map</h3>
							</div>
						</div>
						<div class="card-body">
							<ul class="hexGrid">
								<li class="hex" v-for="(nd, key) in nodes" :key="key">
									<div class="hexIn">
										<a v-bind:class="{ 'bg-success': (nd.status=='Ready'), 'bg-warning': (nd.status!='Ready'), 'hexLink':true }" href="#"><h1>{{ key }}</h1><p class="text-truncate">{{ nd.address }}</p></a>
									</div>
								</li>
							</ul>
						</div>
					</div><!-- /Nodes Map -->
				</div>
			</div>
			<div class="row">
				<div class="col-sm">
					<div class="info-box">
						<div class="info-box-content">
							<span class="info-box-text">Daemon Sets</span>
							<span class="info-box-number">{{ workloads.daemonset.ready | formatNumber }} / {{ workloads.daemonset.available | formatNumber }}</span>
						</div>
						<nuxt-link to="/workload/daemonset.list" class="btn btn-md text-info"><i class="fas fa-arrow-circle-right"></i></nuxt-link>
					</div>
				</div>
				<div class="col-sm">
					<div class="info-box">
						<div class="info-box-content">
							<span class="info-box-text">Deployments</span>
							<span class="info-box-number">{{ workloads.deployment.ready | formatNumber }} / {{ workloads.deployment.available | formatNumber }}</span>
						</div>
						<nuxt-link to="/workload/deployment.list" class="btn btn-md text-info"><i class="fas fa-arrow-circle-right"></i></nuxt-link>
					</div>
				</div>
				<div class="col-sm">
					<div class="info-box">
						<div class="info-box-content">
							<span class="info-box-text">Replica Sets</span>
							<span class="info-box-number">{{ workloads.replicaset.ready | formatNumber }} / {{ workloads.replicaset.available | formatNumber }}</span>
						</div>
						<nuxt-link to="/workload/replicaset.list" class="btn btn-md text-info"><i class="fas fa-arrow-circle-right"></i></nuxt-link>
					</div>
				</div>
				<div class="col-sm">
					<div class="info-box">
						<div class="info-box-content">
							<span class="info-box-text">Stateful Sets</span>
							<span class="info-box-number">{{ workloads.statefulset.ready | formatNumber }} / {{ workloads.statefulset.available | formatNumber }}</span>
						</div>
						<nuxt-link to="/workload/statefulset.list" class="btn btn-md text-info"><i class="fas fa-arrow-circle-right"></i></nuxt-link>
					</div>
				</div>
				<div class="col-sm">
					<div class="info-box">
						<div class="info-box-content">
							<span class="info-box-text">Pods</span>
							<span class="info-box-number">{{ workloads.pod.ready | formatNumber }} / {{ workloads.pod.available | formatNumber }}</span>
						</div>
						<nuxt-link to="/workload/pod.list" class="btn btn-md text-info"><i class="fas fa-arrow-circle-right"></i></nuxt-link>
					</div>
				</div>
			</div>
			<!-- charts -->
			<div class="row">
				<div class="col-md-6">
					<div class="card">
						<div class="card-header border-0">
							<h3 class="card-title">CPU Usages</h3>
							<div class="text-right">
								<span v-if="allocatable.cpu>0" class="border-box background">allocatable : {{ allocatable.cpu  | formatNumber }}</span>
								<span class="border-box">Cores</span>
							</div>
						</div>
						<div class="card-body">
							<div class="chart">
								<c-linechart id="cpuChart" :chart-data="chart.data.cpu" :options="chart.options.cpu" class="mw-100 h-chart"></c-linechart>
							</div>
						</div>
					</div>
				</div>
				<div class="col-md-6">
					<div class="card">
						<div class="card-header border-0">
							<h3 class="card-title">Memory Usages</h3>
							<div class="text-right">
								<span v-if="allocatable.memory>0" class="border-box background">allocatable : {{ allocatable.memory.toFixed(2) | formatNumber }}</span>
								<span class="border-box">GiB</span>
							</div>
						</div>
						<div class="card-body">
							<div class="chart">
								<c-linechart id="memoryChart" :chart-data="chart.data.memory"  :options="chart.options.memory" class="mw-100 h-chart"></c-linechart>
							</div>
						</div>
					</div>
				</div>
			</div><!-- /charts -->
		</div>
	</section>
</div>
</template>
<style scoped>
.card {height: calc(100% - 1rem);}
th {font-weight: normal;font-size:.8em}
td > p {margin-bottom: 0;}
</style>
<script>
import "@/assets/css/hexagons.css"
import VueChartJs		from "vue-chartjs"

export default {
	data() {
		return {
			summary: { allocatable: {}, usage: {}, percent: {} },
			nodes: {},
			workloads: { daemonset: {}, deployment: {}, replicaset: {}, statefulset: {}, pod: {} },
			chart: {
				options: {
					cpu: {
						tooltips: {
							callbacks: {
								label: function(data) {
									return ` ${data.yLabel} core`
								}
							}
						},
						maintainAspectRatio : false, responsive : true, legend: { display: false },
						scales: {
							xAxes: [{ gridLines : {display : false}}],
							yAxes: [{
								gridLines : {display : false},
								ticks: {
									beginAtZero: true,
									suggestedMax: 0,
									callback: (value) => {
										return this.formatNumber(value);
									}
								}
							}]
						}
					},
					memory: {
						tooltips: {
							callbacks: {
								label: function(data) {
									return ` ${data.yLabel} Gi`
								}
							}
						},
						maintainAspectRatio : false, responsive : true, legend: { display: false },
						scales: {
							xAxes: [{ gridLines : {display : false}}],
							yAxes: [{
								gridLines : {display : false},  
								ticks: { 
									beginAtZero: true,
									suggestedMax: 0,
									callback: (value) => {
										return this.formatNumber(value);
									}
								}
							}]
						}
					}
				},
				data: { cpu: {}, memory: {} }
			},
			allocatable: {cpu: 0, memory: 0},
			timer: 0,
			isBusy: true
		}
	},
	layout: "default",
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
	created() {
		this.$nuxt.$on("context-selected", () => {
			let ctx = this.currentContext();
			if(!ctx) return;
			// workloads - available/ready
			this.$axios.get(`/api/clusters/${ctx}/dashboard`)
				.then((resp) => {
					this.workloads = resp.data;
				})
				.catch(e => { this.msghttp(e);});

			// nodes & summary
			this.$axios.get(`/api/clusters/${ctx}/nodes`)
				.then((resp) => {
					this.nodes = resp.data.nodes;
					this.summary = resp.data.summary;
				})
				.catch(e => { 
					this.msghttp(e);
				})
				.finally(()=> {
					this.isBusy = false;
				});

			// metrics
			this.$axios.get(`/api/clusters/${ctx}/metrics`)
				.then((resp) => {
					if (resp.data.metrics ) {
						let labels = [], cpus = [], memories = [];
						resp.data.metrics.forEach(d => {
							let dt = new Date(d.timestamp);
							labels.push(`${dt.getHours()}:${dt.getMinutes()}m`);
							cpus.push((d.cpu/1000).toFixed(3));
							memories.push((d.memory/1024**3).toFixed(2));

						});
						this.allocatable = {
							cpu: resp.data.allocatable ? resp.data.allocatable.cpu/1000: 0,
							memory: resp.data.allocatable ? resp.data.allocatable.memory/1024**3: 0,
						};
						this.chart.options.cpu.scales.yAxes[0].ticks.suggestedMax = this.allocatable.cpu || 1;
						this.chart.options.memory.scales.yAxes[0].ticks.suggestedMax = this.allocatable.memory || 1; 

						this.chart.data.cpu = {
							labels: labels,
							datasets: [
								{ backgroundColor : this.var("CHART_BG_COLOR").cpu, data: cpus }
							]
						};
						this.chart.data.memory = {
							labels: labels,
							datasets: [
								{ backgroundColor : this.var("CHART_BG_COLOR").memory, data: memories }
							]
						};
					}

				})
				.catch(e => { this.msghttp(e);})
		})

		this.$nuxt.$emit("context-selected",);
		this.timer = setInterval(function(){
			this.$nuxt.$emit("context-selected");
		}.bind(this), 30*1000);

	},
	beforeDestroy(){
		this.$nuxt.$off("context-selected")
		clearInterval(this.timer)
	}
}

</script>