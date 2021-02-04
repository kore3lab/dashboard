<template>
<!-- content-wrapper -->
<div class="content-wrapper">

	<div class="content-header">
		<div class="container-fluid">
			<c-navigator group="Cluster"></c-navigator>
			<div class="row mb-2">
				<div class="col-sm-2"><h1 class="m-0 text-dark">Nodes</h1></div>
				<!-- 검색 (검색어) -->
				<div class="col-sm-2 float-left">
					<div class="input-group input-group-sm" >
						<b-form-input id="txtKeyword" v-model="keyword" class="form-control float-right" placeholder="Search"></b-form-input>
						<div class="input-group-append">
							<button type="submit" class="btn btn-default" @click="query_All"><i class="fas fa-search"></i></button>
						</div>
					</div>
				</div><!--//END -->
				<div class="col-sm-8 float-left"></div>
			</div>
		</div>
	</div>

	<section class="content">
	<div class="container-fluid">
		<!-- 검색 -->
		<div class="row mb-2">
			<div class="col-12 text-right "><span class="text-sm align-middle">Total : {{ totalItems }}</span></div>
		</div><!--//END -->
		<!-- GRID-->
		<div class="row">
			<div class="col-12">
				<div class="card">
					<div class="card-body table-responsive p-0">
						<b-table id="list" hover :items="items" :fields="fields" :filter="keyword" :filter-included-fields="filterOn" @filtered="onFiltered" :busy="isBusy" class="text-sm">
							<template #table-busy>
								<div class="text-center text-success" style="margin:150px 0">
									<b-spinner type="grow" variant="success" class="align-middle mr-2"></b-spinner>
									<span class="align-middle text-lg">Loading...</span>
								</div>
							</template>
							<template v-slot:cell(name)="data">
								<nuxt-link :to="{ path:'/view', query:{ context: currentContext(), group: 'Test', crd: 'Node', name: data.item.name, url: `api/v1/nodes/${data.item.name}`, preurl: $router.currentRoute.fullPath}}">{{ data.value }}</nuxt-link>
							</template>
						</b-table>
					</div>
				</div>
			</div>
		</div><!-- //GRID-->
	</div>
	</section>
</div>
</template>
<script>
import axios	from "axios"
import VueNavigator from "@/components/navigator"
export default {
	components: {
		"c-navigator": { extends: VueNavigator }
	},
	data() {
		return {
			keyword: "",
			filterOn: ["name"],
			fields: [
				{ key: "name", label: "이름", sortable: true },
				{ key: "ready", label: "상태", sortable: true  },
				{ key: "creationTimestamp", label: "생성시간" },
				{ key: "k8sVersion", label: "VERSION" },
				{ key: "interaalIp", label: "INTERNAL-IP", sortable: true  },
				{ key: "externalIp", label: "EXTERNAL-IP", sortable: true  },
				{ key: "usageCpu", label: "CPU 사용량", sortable: true  },
				{ key: "usageMemory", label: "MEMORY 사용량", sortable: true  },
			],
			isBusy: false,
			metricsItems: [],
			items: [],
			totalItems: 0
		}
	},
	layout: "default",
	created() {
		this.$nuxt.$on("navbar-context-selected", (ctx) => this.query_All() );
		if(this.currentContext()) this.$nuxt.$emit("navbar-context-selected");
	},
	methods: {
		// 조회
		query_All() {
			console.log("this.$router.currentRoute == ", this.$router.currentRoute.fullPath)
			let interaalIp = {}
			let exteraalIp = {}
			this.isBusy = true;
			axios.get(`${this.backendUrl()}/raw/clusters/${this.currentContext()}/api/v1/nodes`)
				.then((resp) => {
					this.items = [];
					resp.data.items.forEach(el => {
						console.log("adfasdfasdf == ", el)
						const addresses = el.status.addresses
						this.items.push({
							name: el.metadata.name,
							ready: this.toConditions(el.status.conditions),
							creationTimestamp: this.$root.getElapsedTime(el.metadata.creationTimestamp),
							k8sVersion: el.status.nodeInfo.kubeletVersion,
							interaalIp: addresses.find(x => x.type === "InternalIP") ? addresses.find(x => x.type === "InternalIP").address : "<none>",
							externalIp: addresses.find(x => x.type === "ExternalIP") ? addresses.find(x => x.type === "ExternalIP").address : "<none>",
							// cpuRequests: this.toCpuWord(el.allocatedResources.cpuRequests, el.allocatedResources.cpuRequestsFraction),
							// cpuLimits: this.toCpuWord(el.allocatedResources.cpuLimits, el.allocatedResources.cpuLimitsFraction) ,
							// memoryRequests: this.toMemoryWord( el.allocatedResources.memoryRequests, el.allocatedResources.memoryRequestsFraction),
							// memoryLimits: this.toMemoryWord( el.allocatedResources.memoryLimits, el.allocatedResources.memoryLimitsFraction),
						});
					});
					this.onFiltered(this.items);
				})
				.catch(e => { this.msghttp(e);})
				.finally(()=> { this.isBusy = false;});
		},
		async getMetrics() {
			this.metricsItems = [];
			let resp = await axios.get(`${this.backendUrl()}/raw/clusters/${this.currentContext()}/apis/metrics.k8s.io/v1beta1/nodes`)
			resp.data.items.forEach(el => {
				this.metricsItems.push(el)
			});
		},
		onFiltered(filteredItems) {
			this.totalItems = filteredItems.length;
		},
		toCpuWord(cpu, percent) {
			return cpu<1000 ? `${cpu.toFixed(2)}m (${percent.toFixed(2)}%)`: `${(cpu/1000).toFixed(2)} (${percent.toFixed(2)}%)`;
		},
		toMemoryWord(memory, percent) {
			let mi  = 1024*1024
			let gi  = mi*1024

			if(memory > gi) {
				return `${(memory/gi).toFixed(2)}Gi (${percent.toFixed(2)}%)`
			} else if(memory > mi) {
				return `${(memory/mi).toFixed(2)}Mi (${percent.toFixed(2)}%)`
			} else {
				return `${(memory/1024).toFixed(2)} (${percent.toFixed(2)}%)`
			}

		},
		toConditions(conditions){
			return conditions.filter(el => el.type === "Ready")[0].status === "True" ? "Ready" : "None"
		}
	},
	beforeDestroy(){
		this.$nuxt.$off('navbar-context-selected')
	}
}
</script>
<style>label {font-weight: 500;}</style>
