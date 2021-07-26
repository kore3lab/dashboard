<template>
	<div class="content-wrapper">
		<section class="content-header">
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
		</section>

		<section class="content">
			<div class="container-fluid">
				<!-- total count & items per page  -->
				<div class="d-flex flex-row-reverse">
					<div class="p-2">
						<b-form inline>
							<c-colums-selector name="grdSheet1" v-model="fields" :fields="fieldsAll" ></c-colums-selector>
							<i class="text-secondary ml-2 mr-2">|</i>
							<span class="text-sm align-middle">Total : {{ totalItems }}</span>
						</b-form>
					</div>
				</div>
				<!-- GRID-->
				<div class="row">
					<div class="col-12">
						<div class="card">
							<div class="card-body table-responsive p-0">
								<b-table hover selectable show-empty select-mode="single" @row-selected="onRowSelected" @sort-changed="currentPage=1" ref="grdSheet1" :items="items" :fields="fields" :filter="keyword" :filter-included-fields="filterOn" @filtered="onFiltered" :busy="isBusy" class="text-sm">
									<template #table-busy>
										<div class="text-center text-success lh-vh-50">
											<b-spinner type="grow" variant="success" class="align-middle mr-2"></b-spinner>
											<span class="text-lg align-middle">Loading...</span>
										</div>
									</template>
									<template v-slot:cell(ready)="data">
										<span v-for="(value, idx) in data.item.ready" v-bind:key="idx" v-bind:class="value.style" class="mr-1" >{{ value.value }}</span>
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
								</b-table>
							</div>
						</div>
					</div>
				</div><!-- //GRID-->
			</div>
		</section>
		<b-sidebar v-model="isShowSidebar" width="50em" @hidden="$refs.grdSheet1.clearSelected()" right shadow no-header>
			<c-view v-model="viewModel" @delete="query_All()" @close="isShowSidebar=false"/>
		</b-sidebar>
	</div>
</template>
<script>
import VueNavigator			from "@/components/navigator"
import VueColumsSelector	from "@/components/columnsSelector"
import VueView				from "@/pages/view";

export default {
	components: {
		"c-navigator": { extends: VueNavigator },
		"c-colums-selector": { extends: VueColumsSelector},
		"c-view": { extends: VueView }
	},
	data() {
		return {
			keyword: "",
			filterOn: ["name"],
			fields: [],
			fieldsAll: [
				{ key: "name", label: "Name", sortable: true },
				{ key: "usageCpu", label: "CPU", sortable: true  },
				{ key: "usageMemory", label: "Memory", sortable: true  },
				{ key: "usageDisk", label: "Disk", sortable: true },
				{ key: "taints", label: "Taints", sortable: true },
				{ key: "roles", label: "Roles", sortable: true },
				{ key: "k8sVersion", label: "Version", sortable: true },
				{ key: "creationTimestamp", label: "Age", sortable: true, formatter: this.getElapsedTime },
				{ key: "ready", label: "Status", sortable: true },
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
		this.$nuxt.$on("navbar-context-selected", (ctx) => {
			this.onUsage()
		} );
		if(this.currentContext()) this.$nuxt.$emit("navbar-context-selected");
	},
	methods: {
		onRowSelected(items) {
			this.isShowSidebar = (items && items.length > 0)
			if (this.isShowSidebar) this.viewModel = this.getViewLink('', 'nodes', items[0].namespace, items[0].name)
		},
		// 조회
		query_All() {
			this.$axios.get(this.getApiUrl("","nodes"))
					.then((resp) => {
						this.items = [];
						resp.data.items.forEach(el => {
							this.items.push({
								name: el.metadata.name,
								ready: this.getConditions(el),
								creationTimestamp: el.metadata.creationTimestamp,
								k8sVersion: el.status.nodeInfo.kubeletVersion,
								taints: el.spec.taints ? el.spec.taints.length: 0,
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
		getRoles(labels) {
			let roleLabels = Object.keys(labels).filter(key =>
					key.includes("node-role.kubernetes.io")
			).map(key => key.match(/([^/]+$)/)[0]);

			if (labels["kubernetes.io/role"] !== undefined) {
				roleLabels.push(labels["kubernetes.io/role"]);
			}

			return roleLabels.join(", ");
		},
		// node cpu,memory 사용량 먼저 읽은 후 전체리스트 조회
		onUsage() {
			this.isBusy = true;
			this.$axios.get(`/api/clusters/${this.currentContext()}/dashboard`)
					.then((resp) => {
						this.metrics = resp.data.nodes
					}).finally(()=> { this.query_All()} )
		},
		getCpu(name) {
			if(!this.metrics[name]) return

			return this.metrics[name].usage.cpu.percent

		},
		getMemory(name) {
			if(!this.metrics[name]) return

			return this.metrics[name].usage.memory.percent
		},
		getDisk(name) {
			if(!this.metrics[name]) return

			return this.metrics[name].usage.storage.percent
		}
	},
	beforeDestroy(){
		this.$nuxt.$off('navbar-context-selected')
	}
}
</script>
