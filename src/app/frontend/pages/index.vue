<template>
	<!-- content-wrapper -->
	<div class="content-wrapper">

		<!-- Content Header (Page header) -->
		<div class="content-header">
			<div class="container-fluid">
			</div>
		</div>

		<!-- Main content -->
		<section class="content">
			<div class="container-fluid">
				<!-- @@@@ -->
				<!-- row -->
				<div class="row">
					<div class="col-lg-3 col-6">
						<div class="small-box bg-warning">
							<div class="inner"><h3>{{ summary.nodes.usage }}</h3><p>Nodes</p></div>
							<div class="icon"><i class="fas fa-server"></i></div>
						</div>
					</div>
					<div class="col-lg-3 col-6">
						<div class="small-box bg-success">
							<div class="inner"><h3>{{ Math.round(summary.cpu["percent"],0) }}<small>%</small></h3><p>CPU</p></div>
							<div class="icon"><i class="fas fa-microchip"></i></div>
						</div>
					</div>
					<div class="col-lg-3 col-6">
						<div class="small-box bg-info">
							<div class="inner"><h3>{{ Math.round(summary.memory["percent"],0) }}<small>%</small></h3><p>Memory</p></div>
							<div class="icon"><i class="fas fa-memory"></i></div>
						</div>
					</div>
					<div class="col-lg-3 col-6">
						<div class="small-box bg-secondary">
							<div class="inner"><h3>{{ Math.round(summary.storage["percent"],0) }}<small>%</small></h3><p>Storage</p></div>
							<div class="icon"><i class="fas fa-hdd"></i></div>
						</div>
					</div>
				</div><!-- /.row -->
				<!-- row -->
				<div class="row">
					<div class="col-lg-7">
						<div class="card">
							<div class="card-header">
								<div class="d-flex justify-content-between">
									<h3 class="card-title">Nodes</h3>
								</div>
							</div>
							<div class="card-body">
								<div class="d-flex flex-row justify-content-between align-items-center"  v-for="(nd, key) in nodes" :key="key">
									<p class="text-lg mr-auto">{{ key }}</p>
									<p class="d-flex text-left p-2">
										<span class="badge badge-success font-weight-light ml-1">{{ nd.status }}</span>
										<span class="badge badge-secondary font-weight-light ml-1">{{ nd.roles }}</span>
									</p>
									<p class="d-flex flex-column text-center p-2">
										<span class="text-lg">{{ nd.usage.cpu.percent }}<small>%</small></span>
										<span class="text-muted text-sm font-weight-light">{{ Number(nd.usage.cpu.usage).toLocaleString() }}/{{ Number(nd.usage.cpu.allocatable).toLocaleString() }} m</span>
									</p>
									<p class="d-flex flex-column text-center p-2">
										<span class="text-lg">{{ nd.usage.memory.percent }}<small>%</small></span>
										<span class="text-muted text-sm font-weight-light">{{ Number(Math.round(nd.usage.memory.usage/(1024*1024),2)).toLocaleString() }}/{{ Number(Math.round(nd.usage.memory.allocatable/(1024*1024),2)).toLocaleString() }} Mib</span>
									</p>
									<p class="d-flex flex-column text-center p-2">
										<span class="text-lg">{{ nd.usage.storage.percent }}<small>%</small></span>
										<span class="text-muted text-sm font-weight-light">{{ Number(Math.round(nd.usage.storage.usage/(1024*1024),2)).toLocaleString() }}/{{ Number(Math.round(nd.usage.storage.allocatable/(1024*1024),2)).toLocaleString() }} Gib</span>
									</p>
									<p class="d-flex flex-column text-center  p-2">
										<span class="text-lg">{{ nd.usage.pod.percent }}<small>%</small></span>
										<span class="text-muted text-sm font-weight-light">{{ nd.usage.pod.usage }}/{{ nd.usage.pod.allocatable }} ea</span>
									</p>
								</div>
							</div>
						</div>
					</div>
					<div class="col-lg-5">
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
											<a v-bind:class="{ 'bg-success': (nd.status=='Ready'), 'bg-warning': (nd.status=='NotReady'), 'hexLink':true }" href="#"><h1>{{ key }}</h1><p class="text-truncate">{{ nd.address }}</p></a>
										</div>
									</li>
								</ul>
							</div>
						</div>
					</div>
				</div><!-- /.row -->
				<!-- row -->
				<div class="row">
					<div class="col-sm">
						<div class="info-box">
							<div class="info-box-content">
								<span class="info-box-text">Daemon Sets</span>
								<span class="info-box-number">{{ workloads.daemonset.ready }} / {{ workloads.daemonset.available }}</span>
							</div>
							<nuxt-link to="/workload/daemonset.list" class="small-box-footer"><i class="fas fa-arrow-circle-right"></i></nuxt-link>
						</div>
					</div>
					<div class="col-sm">
						<div class="info-box">
							<div class="info-box-content">
								<span class="info-box-text">Deployments</span>
								<span class="info-box-number">{{ workloads.deployment.ready }} / {{ workloads.deployment.available }}</span>
							</div>
							<nuxt-link to="/workload/deployment.list" class="small-box-footer"><i class="fas fa-arrow-circle-right"></i></nuxt-link>
						</div>
					</div>
					<div class="col-sm">
						<div class="info-box">
							<div class="info-box-content">
								<span class="info-box-text">Replica Sets</span>
								<span class="info-box-number">{{ workloads.replicaset.ready }} / {{ workloads.replicaset.available }}</span>
							</div>
							<nuxt-link to="/workload/replicaset.list" class="small-box-footer"><i class="fas fa-arrow-circle-right"></i></nuxt-link>
						</div>
					</div>
					<div class="col-sm">
						<div class="info-box">
							<div class="info-box-content">
								<span class="info-box-text">Stateful Sets</span>
								<span class="info-box-number">{{ workloads.statefulset.ready }} / {{ workloads.statefulset.available }}</span>
							</div>
							<nuxt-link to="/workload/statefulset.list" class="small-box-footer"><i class="fas fa-arrow-circle-right"></i></nuxt-link>
						</div>
					</div>
					<div class="col-sm">
						<div class="info-box">
							<div class="info-box-content">
								<span class="info-box-text">Pods</span>
								<span class="info-box-number">{{ workloads.pod.ready }} / {{ workloads.pod.available }}</span>
							</div>
							<nuxt-link to="/workload/pod.list" class="small-box-footer"><i class="fas fa-arrow-circle-right"></i></nuxt-link>
						</div>
					</div>
				</div><!-- /.row -->
				<!-- row -->
				<div class="row">
					<div class="col-md-6">
						<div class="card">
							<div class="card-header border-0">
								<h3 class="card-title">CPU Usages</h3>
							</div>
							<div class="card-body">
								<div class="chart">
									<c-linechart id="cpuChart" :chart-data="chart.data.cpu" :options="chart.options.cpu" style="min-height: 250px; height: 250px; max-height: 250px; max-width: 100%;"></c-linechart>
								</div>
							</div>
						</div>
					</div>
					<div class="col-md-6">
						<div class="card">
							<div class="card-header border-0">
								<h3 class="card-title">Memory Usages</h3>
							</div>
							<div class="card-body">
								<div class="chart">
									<c-linechart id="memoryChart" :chart-data="chart.data.memory"  :options="chart.options.memory" style="min-height: 250px; height: 250px; max-height: 250px; max-width: 100%;"></c-linechart>
								</div>
							</div>
						</div>
					</div>
				</div><!-- /.row -->
				<!-- //@@@@ -->
			</div>
		</section>
	</div>
	<!-- /.content-wrapper -->
</template>
<script>
import "@/assets/css/hexagons.css"
import VueChartJs	from "vue-chartjs"
import axios		from "axios"

export default {
	data() {
		return {
			summary: { nodes: {}, cpu: {}, memory: {}, storage: {} },
			nodes: {},
			workloads: { daemonset: {}, deployment: {}, replicaset: {}, statefulset: {}, pod: {} },
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
				data: { cpu: {}, memory: {} }
			},
			timer: 0
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
				}
			}
		}
	},
	created() {
		this.$nuxt.$on("navbar-context-selected", () => {
			let ctx = this.currentContext();
			if(!ctx) return;
			axios.get(`${this.backendUrl()}/api/clusters/${ctx}/dashboard`)
					.then((resp) => {
						this.$data.summary = resp.data.summary;
						this.$data.nodes = resp.data.nodes;
						this.$data.workloads = resp.data.workloads;

						if (resp.data.metrics && resp.data.metrics.cpu && resp.data.metrics.cpu.dataPoints ) {
							let labels = [], data = [];
							resp.data.metrics.cpu.dataPoints.forEach(d => {
								let dt = new Date(d.x*1000);
								labels.push(`${dt.getHours()}:${dt.getMinutes()}m`);
								data.push(d.y);
							});
							this.$data.chart.options.cpu.scales.yAxes[0].ticks.suggestedMax = resp.data.summary.cpu.allocatable;
							this.$data.chart.data.cpu = {
								labels: labels,
								datasets: [
									{ backgroundColor : "rgba(60,141,188,0.9)", data: data }
								]
							};
						}
						if (resp.data.metrics && resp.data.metrics.memory && resp.data.metrics.memory.dataPoints ) {
							let labels = [], data = [];
							resp.data.metrics.memory.dataPoints.forEach(d => {
								let dt = new Date(d.x*1000);
								labels.push(`${dt.getHours()}:${dt.getMinutes()}m`);
								data.push(Math.round(d.y/(1024*1024)));
							});
							this.$data.chart.options.memory.scales.yAxes[0].ticks.suggestedMax = resp.data.summary.memory.allocatable/(1024*1024);
							this.$data.chart.data.memory = {
								labels: labels,
								datasets: [
									{ backgroundColor : "rgba(210, 214, 222, 1)", data: data }
								]
							};
						}

					})
					.catch(e => { this.msghttp(e);})

		})

		this.$nuxt.$emit("navbar-context-selected",);
		this.timer = setInterval(function(){
			this.$nuxt.$emit("navbar-context-selected");
		}.bind(this), 30*1000);

	},
	beforeDestroy(){
		this.$nuxt.$off("navbar-context-selected")
		clearInterval(this.timer)
	}
}

</script>
