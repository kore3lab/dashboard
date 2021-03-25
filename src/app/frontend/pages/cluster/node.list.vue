<template>
	<div class="content-wrapper">
		<div class="content-header">
			<div class="container-fluid">
				<c-navigator group="Cluster"></c-navigator>
				<div class="row mb-2">
					<!-- title & search -->
					<div class="col-sm"><h1 class="m-0 text-dark"><span class="badge badge-info mr-2">N</span>Nodes</h1></div>
					<div class="col-sm-2 float-left">
						<div class="input-group input-group-sm" >
							<b-form-input id="txtKeyword" v-model="keyword" class="form-control float-right" placeholder="Search"></b-form-input>
							<div class="input-group-append"><button type="submit" class="btn btn-default" @click="query_All"><i class="fas fa-search"></i></button></div>
						</div>
					</div>
				</div>
			</div>
		</div>

		<section class="content">
			<div class="container-fluid">
				<!-- count -->
				<div class="row mb-2">
					<div class="col-12 text-right "><span class="text-sm align-middle">Total : {{ totalItems }}</span></div>
				</div>
				<!-- GRID-->
				<div class="row">
					<div class="col-12">
						<div class="card">
							<div class="card-body table-responsive p-0">
								<b-table id="list" hover :items="items" :fields="fields" :filter="keyword" :filter-included-fields="filterOn" @filtered="onFiltered" :busy="isBusy" fixed class="text-sm">
									<template #table-busy>
										<div class="text-center text-success" style="margin:150px 0">
											<b-spinner type="grow" variant="success" class="align-middle mr-2"></b-spinner>
											<span class="align-middle text-lg">Loading...</span>
										</div>
									</template>
									<template v-slot:cell(name)="data">
										<a href="#" @click="viewModel=getViewLink('','nodes',data.item.namespace, data.item.name); isShowSidebar=true;">{{ data.value }}</a>
									</template>
									<template v-slot:cell(ready)="data">
										<span v-for="(value, idx) in data.item.ready" v-bind:key="idx" v-bind:class="value.style" >{{ value.value }}  </span>
									</template>
									<template v-slot:cell(usageCpu)="data">
										<b-progress :value="data.item.usageCpu" :max="100" variant="info" show-value class="mb-3"></b-progress>
									</template>
									<template v-slot:cell(usageMemory)="data">
										<b-progress :value="data.item.usageMemory" :max="100" variant="info" show-value class="mb-3"></b-progress>
									</template>
									<template v-slot:cell(usageDisk)="data">
										<b-progress :value="data.item.usageDisk" :max="100" variant="info" show-value class="mb-3"></b-progress>
									</template>
									<template v-slot:cell(creationTimestamp)="data">
										{{ data.value.str }}
									</template>
								</b-table>
							</div>
						</div>
					</div>
				</div><!-- //GRID-->
			</div>
		</section>
		<b-sidebar v-model="isShowSidebar" width="50em" right shadow no-header>
			<c-view v-model="viewModel" @delete="query_All()" @close="isShowSidebar=false"/>
		</b-sidebar>
	</div>
</template>
<script>
import axios		from "axios"
import VueNavigator from "@/components/navigator"
import VueView from "@/pages/view";

export default {
	components: {
		"c-navigator": { extends: VueNavigator },
		"c-view": { extends: VueView }
	},
	data() {
		return {
			keyword: "",
			filterOn: ["name"],
			fields: [
				{ key: "name", label: "Name", sortable: true, class:"text-truncate" },
				{ key: "usageCpu", label: "CPU", sortable: true, class:"text-truncate"  },
				{ key: "usageMemory", label: "Memory", sortable: true, class:"text-truncate"  },
				{ key: "usageDisk", label: "Disk", sortable: true, class:"text-truncate" },
				{ key: "taints", label: "Taints", sortable: true, class:"text-truncate" },
				{ key: "roles", label: "Roles", sortable: true, class:"text-truncate" },
				{ key: "k8sVersion", label: "Version", sortable: true, class:"text-truncate" },
				{ key: "creationTimestamp", label: "Age", sortable: true, class:"text-truncate" },
				{ key: "ready", label: "Status", sortable: true, class:"text-truncate"  },
			],
			isBusy: false,
			items: [],
			totalItems: 0,
			metrics: [],
			isShowSidebar: false,
			viewModel:{},
		}
	},
	layout: "default",
	created() {
		this.$nuxt.$on("navbar-context-selected", (ctx) =>this.onUsage() );
		if(this.currentContext()) this.$nuxt.$emit("navbar-context-selected");
	},
	methods: {
		// 조회
		query_All() {
			this.isBusy = true;
			axios.get(this.getApiUrl("","nodes"))
					.then((resp) => {
						this.items = [];
						resp.data.items.forEach(el => {
							const addresses = el.status.addresses
							this.items.push({
								name: el.metadata.name,
								ready: this.getConditions(el),
								creationTimestamp: this.getElapsedTime(el.metadata.creationTimestamp),
								k8sVersion: el.status.nodeInfo.kubeletVersion,
								taints: this.getTaints(el.spec),
								roles: this.getRoles(el.metadata.labels),
								usageCpu: this.getCpu(el.metadata.name),
								usageMemory: this.getMemory(el.metadata.name),
								usageDisk: this.getDisk(el.metadata.name),
							});
						});
						this.onFiltered(this.items);
					})
					.catch(e => { this.msghttp(e);})
					.finally(()=> { this.isBusy = false;});
		},
		onFiltered(filteredItems) {
			this.totalItems = filteredItems.length;
		},

		// node condition 체크
		getConditions(el) {
			let condition = [];
			if (el.spec.unschedulable) {
				condition.push({
							"value": "SchedulingDisabled",
							"style": "text-warning"
						}
				)
			}
			condition.push(el.status.conditions.filter(con => con.type === "Ready")[0].status === "True" ? { "value" : "Ready", "style" : "text-success" } : { "value" : "NotReady", "style" : "text-secondary" })
			return condition
		},
		getTaints(spec) {
			if (spec.taints) {
				return spec.taints.length
			} else return 0
		},
		getRoles(labels) {
			let roleLabels = Object.keys(labels).filter(key =>
					key.includes("node-role.kubernetes.io")
			).map(key => key.match(/([^/]+$)/)[0]);

			if (labels["kubernetes.io/role"] !== undefined) {
				roleLabels.push(labels["kubernetes.io/role"]);
			}

			return roleLabels.join(", ");
		},
		// node cpu,memory,disk 사용량 먼저 읽은 후 전체리스트 조회
		onUsage() {
			axios.get(`${this.backendUrl()}/api/clusters/${this.currentContext()}/dashboard`)
					.then((resp) => {
						this.metrics = resp.data.nodes
					}).finally(()=> { this.query_All()} )
		},
		getCpu(name) {
			return this.metrics[name].usage.cpu.percent

		},
		getMemory(name) {
			return this.metrics[name].usage.memory.percent
		},
		getDisk(name) {
			return this.metrics[name].usage.storage.percent
		}
	},
	beforeDestroy(){
		this.$nuxt.$off('navbar-context-selected')
	}
}
</script>
<style scoped>label {font-weight: 500;}</style>
