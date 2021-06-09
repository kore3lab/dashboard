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
								<b-table id="list" hover selectable show-empty select-mode="single" @row-selected="onRowSelected" @sort-changed="onSortChanged()" ref="selectableTable" :items="items" :fields="fields" :filter="keyword" :filter-included-fields="filterOn" @filtered="onFiltered" :busy="isBusy" class="text-sm">
									<template #table-busy>
										<div class="text-center text-success lh-vh-50">
											<b-spinner type="grow" variant="success" class="align-middle mr-2"></b-spinner>
											<span class="text-lg align-middle">Loading...</span>
										</div>
									</template>
									<template #empty="scope">
										<h4 class="text-center">does not exist.</h4>
									</template>
									<template v-slot:cell(name)="data">
										{{ data.value }}
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
									<template v-slot:cell(creationTimestamp)="data">
										{{ data.value.str }}
									</template>
									<template #head(button)>
										<div class="text-right">
											<a id="colOpt" class="nav-link" href="#"><i class="fas fa-ellipsis-v"></i></a>
										</div>
										<b-popover triggers="focus" ref="popover" target="colOpt" placement="bottomleft">
											<b-form-group>
												<b-form-checkbox v-for="option in columnOpt" v-model="selected" :key="option.key" :value="option.label" name="flavour-3a">
													{{ option.label }}
												</b-form-checkbox>
											</b-form-group>
										</b-popover>
									</template>
								</b-table>
							</div>
						</div>
					</div>
				</div><!-- //GRID-->
			</div>
		</section>
		<b-sidebar v-model="isShowSidebar" width="50em" right shadow no-header>
			<c-view v-model="viewModel" @delete="query_All()" @close="onRowSelected"/>
		</b-sidebar>
	</div>
</template>
<script>
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
				{ key: "name", label: "Name", sortable: true },
				{ key: "usageCpu", label: "CPU", sortable: true  },
				{ key: "usageMemory", label: "Memory", sortable: true  },
				{ key: "usageDisk", label: "Disk", sortable: true },
				{ key: "taints", label: "Taints", sortable: true },
				{ key: "roles", label: "Roles", sortable: true },
				{ key: "k8sVersion", label: "Version", sortable: true },
				{ key: "creationTimestamp", label: "Age", sortable: true },
				{ key: "ready", label: "Status", sortable: true },
			],
			isBusy: false,
			items: [],
			currentItems:[],
			columnOpt: [],
			selected: [],
			selectIndex: 0,
			totalItems: 0,
			metrics: [],
			isShowSidebar: false,
			viewModel:{},
		}
	},
	layout: "default",
	watch: {
		selected() {
			this.fields = []
			this.columnOpt.forEach(el => {
				this.selected.forEach(e => {
					if(el.label === e) {
						this.fields.push(el)
					}
				})
			})
			this.fields.push({ key: "button", label: "button", thClass: "wt10"})
			localStorage.setItem('columns_node',this.selected)
		}
	},
	created() {
		this.columnOpt = Object.assign([],this.fields)
		this.$nuxt.$on("navbar-context-selected", (ctx) => {
			if(localStorage.getItem('columns_node')) {
				this.selected = (localStorage.getItem('columns_node')).split(',')
			} else {
				this.fields.forEach(el => {
					this.selected.push(el.label)
				})
			}
			this.onUsage()
		} );
		if(this.currentContext()) this.$nuxt.$emit("navbar-context-selected");
	},
	methods: {
		onSortChanged() {
			this.currentPage = 1
		},
		onRowSelected(items) {
			if(items) {
				if(items.length) {
					for(let i=0;i<this.$config.itemsPerPage;i++) {
						if (this.$refs.selectableTable.isRowSelected(i)) this.selectIndex = i
					}
					this.viewModel = this.getViewLink('', 'nodes', items[0].namespace, items[0].name)
					if(this.currentItems.length ===0) this.currentItems = Object.assign({},this.viewModel)
					this.isShowSidebar = true
				} else {
					if(this.currentItems.title !== this.viewModel.title) {
						if(this.currentItems.length ===0) this.isShowSidebar = false
						else {
							this.viewModel = Object.assign({},this.currentItems)
							this.currentItems = []
							this.isShowSidebar = true
							this.$refs.selectableTable.selectRow(this.selectIndex)
						}
					} else {
						this.isShowSidebar = false
						this.$refs.selectableTable.clearSelected()
					}
				}
			} else {
				this.currentItems = []
				this.isShowSidebar = false
				this.$refs.selectableTable.clearSelected()
			}
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
<style scoped>label {font-weight: 500;}</style>
