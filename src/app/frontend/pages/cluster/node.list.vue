<template>
<!-- content-wrapper -->
<div class="content-wrapper">

	<div class="content-header">
		<div class="container-fluid">
			<c-navigator group="Cluster"></c-navigator>
			<div class="row mb-2">
				<div class="col-sm-2"><h1 class="m-0 text-dark">Nodes</h1></div>
				<div class="col-sm-2 float-left">
					<div class="input-group input-group-sm" >
						<input type="text" name="table_search" class="form-control float-right" placeholder="Search">
						<div class="input-group-append">
							<button type="submit" class="btn btn-default"><i class="fas fa-search"></i></button>
						</div>
					</div>
				</div>
				<div class="col-sm-8 float-left"></div>
			</div>
		</div>
	</div>

	<section class="content">
	<div class="container-fluid">
		<div class="row">
			<div class="col-12">
				<!-- GRID-->
				<div class="card">
					<div class="card-body table-responsive p-0">
						<b-table id="list" hover :items="items" :fields="fields" class="text-sm">
							<template v-slot:cell(name)="data">
								<nuxt-link :to="{ path:'/view', query:{ context: currentContext(), group: 'Cluster', crd: 'Node', name: data.item.name, url: `node/name/${data.item.name}`}}">{{ data.value }}</nuxt-link>
							</template>
						</b-table>
					</div>
				</div>
				<!-- //GRID-->
			</div>
		</div>
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
			fields: [
				{ key: "name", label: "이름", sortable: true },
				{ key: "ready", label: "준비", sortable: true  },
				{ key: "cpuRequests", label: "CPU 요청", sortable: true  },
				{ key: "cpuLimits", label: "CPU 상한", sortable: true  },
				{ key: "memoryRequests", label: "메모리 요청", sortable: true  },
				{ key: "memoryLimits", label: "메모리 상한", sortable: true  },
				{ key: "creationTimestamp", label: "생성시간" }
			],
			items: [],
			origin: []
		}
	},
	layout: "default",
	created() {
		this.$nuxt.$on("navbar-context-selected", (ctx) => this.query_All() );
		if(this.currentContext()) this.$nuxt.$emit("navbar-context-selected");
	},
	methods: {
		query_All() {
			this.$data.selected = [];
			axios.get(`${this.dashboardUrl()}/api/v1/node?sortBy=d,creationTimestamp&context=${this.currentContext()}`)
				.then((resp) => {
					this.$data.origin = [];
					resp.data.nodes.forEach(el => {
						this.$data.origin.push({
							name: el.objectMeta.name,
							ready: el.ready,
							cpuRequests: this.toCpuWord(el.allocatedResources.cpuRequests, el.allocatedResources.cpuRequestsFraction),
							cpuLimits: this.toCpuWord(el.allocatedResources.cpuLimits, el.allocatedResources.cpuLimitsFraction) ,
							memoryRequests: this.toMemoryWord( el.allocatedResources.memoryRequests, el.allocatedResources.memoryRequestsFraction),
							memoryLimits: this.toMemoryWord( el.allocatedResources.memoryLimits, el.allocatedResources.memoryLimitsFraction),
							creationTimestamp: this.$root.getTimestampString(el.objectMeta.creationTimestamp)
						});
					});
					this.$data.items = this.$data.origin;
				})
				.catch((error) => {
					this.$root.toast(error.message, "danger");
				});
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

		}
	},
	beforeDestroy(){
		this.$nuxt.$off('navbar-context-selected')
	}
}
</script>
<style>label {font-weight: 500;}</style>
