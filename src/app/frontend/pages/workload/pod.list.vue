<template>
	<div class="content-wrapper">
		<div class="content-header">
			<div class="container-fluid">
				<c-navigator group="Workload"></c-navigator>
				<div class="row mb-2">
					<!-- title & search -->
					<div class="col-sm"><h1 class="m-0 text-dark"><span class="badge badge-info mr-2">P</span>Pods</h1></div>
					<div class="col-sm-2"><b-form-select v-model="selectedNamespace" :options="namespaces()" size="sm" @input="query_All(); selectNamespace(selectedNamespace);"></b-form-select></div>
					<div class="col-sm-2 float-left">
						<div class="input-group input-group-sm" >
							<b-form-input id="txtKeyword" v-model="keyword" class="form-control float-right" placeholder="Search"></b-form-input>
							<div class="input-group-append"><button type="submit" class="btn btn-default" @click="query_All"><i class="fas fa-search"></i></button></div>
						</div>
					</div>
					<!-- button -->
					<div class="col-sm-1 text-right">
						<b-button variant="primary" size="sm" @click="$router.push(`/create?context=${currentContext()}&group=Workload&crd=Pod`)">Create</b-button>
					</div>
				</div>
			</div>
		</div>

		<section class="content">
			<div class="container-fluid">
				<!-- search & filter -->
				<div class="row mb-2">
					<div class="col-10">
						<b-form-group class="mb-0 font-weight-light overflow-auto">
							<button type="submit" class="btn btn-default btn-sm float-left mr-2" @click="selectedClear">All</button>
							<b-form-checkbox-group v-model="selectedStatus" :options="optionsStatus" button-variant="light" font="light" switches size="sm" @input="onChangeStatus" class="float-left"></b-form-checkbox-group>
						</b-form-group>
					</div>
					<div class="col-2 text-right "><span class="text-sm align-middle">Total : {{ totalItems }}</span></div>
				</div>
				<!-- GRID-->
				<div class="row">
					<div class="col-12">
						<div class="card">
							<div class="card-body table-responsive p-0">
								<b-table id="list" hover selectable show-empty select-mode="single" @sort-changed="onSortChanged()" @row-selected="onRowSelected" ref="selectableTable" :items="items" :fields="fields" :filter="keyword" :filter-included-fields="filterOn" @filtered="onFiltered" :current-page="currentPage" :per-page="$config.itemsPerPage" :busy="isBusy" class="text-sm">
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
									<template v-slot:cell(status)="data">
										<div class="list-unstyled mb-0" v-if="data.item.status.value">
											<span v-bind:class="data.item.status.style">{{ data.item.status.value }}</span>
										</div>
									</template>
									<template v-slot:cell(containers)="data">
										<div class="list-unstyled mb-0 float-left" v-if="data.item.containers.containerStatuses.length > 0">
											<span v-for="(containerStatuses, idx) in data.item.containers.containerStatuses" v-bind:key="idx" v-bind:class="containerStatuses.style" class="badge font-weight-light text-sm ml-1"> {{" "}}</span>
										</div>
										<div class="list-unstyled mb-0 ml-0 float-left" v-if="data.item.containers.initContainerStatuses.length > 0">
											<span v-for="(initContainerStatuses, idx) in data.item.containers.initContainerStatuses" v-bind:key="idx" v-bind:class="initContainerStatuses.style" class="badge font-weight-light text-sm ml-1">{{" "}}</span>
										</div>
									</template>
									<template v-slot:cell(controller)="data">
										<a href="#" @click="viewModel=getViewLink(data.value.group,data.value.rs,data.item.namespace, data.value.name); isShowSidebar=true;">{{ data.value.kind }}</a>
									</template>
									<template v-slot:cell(node)="data">
										<a href="#" @click="viewModel=getViewLink(data.value.group,data.value.rs, '',  data.value.name); isShowSidebar=true;">{{ data.value.name }}</a>
									</template>
									<template v-slot:cell(creationTimestamp)="data">
										{{ data.value.str }}
									</template>
								</b-table>
							</div>
								<b-pagination v-model="currentPage" :per-page="$config.itemsPerPage" :total-rows="totalItems" size="sm" align="center"></b-pagination>
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
			selectedNamespace: "",
			selectedStatus: [],
			allStatus: ["Running", "Pending", "Terminating", "CrashLoopBackOff", "ImagePullBackOff", "Completed", "ContainerCreating", "Failed", "etc"],
			optionsStatus: [
				{ text: "Running", value: "Running" },
				{ text: "Pending", value: "Pending" },
				{ text: "Terminating", value: "Terminating" },
				{ text: "CrashLoopBackOff", value: "CrashLoopBackOff" },
				{ text: "ImagePullBackOff", value: "ImagePullBackOff" },
				{ text: "Completed", value: "Completed" },
				{ text: "ContainerCreating", value: "ContainerCreating" },
				{ text: "Failed", value: "Failed" },
				{ text: "etc", value: "etc" }
			],
			keyword: "",
			filterOn: ["name"],
			fields: [
				{ key: "name", label: "Name", sortable: true },
				{ key: "namespace", label: "Namespace", sortable: true  },
				{ key: "ready", label: "Ready", sortable: true  },
				{ key: "containers", label: "Containers", sortable: true  },
				{ key: "restartCount", label: "Restart", sortable: true  },
				{ key: "controller", label: "Controlled By", sortable: true  },
				{ key: "node", label: "Node", sortable: true  },
				{ key: "qos", label: "QoS", sortable: true  },
				{ key: "creationTimestamp", label: "Age", sortable: true  },
				{ key: "status", label: "Status", sortable: true  },
			],
			isBusy: false,
			metricsItems: [],
			origin: [],
			items: [],
			currentitems:[],
			selectIndex: 0,
			containerStatuses: [],
			currentPage: 1,
			totalItems: 0,
			isShowSidebar: false,
			viewModel:{},
		}
	},
	layout: "default",
	created() {
		this.$nuxt.$on("navbar-context-selected", (ctx) => {
			this.selectedNamespace = this.selectNamespace()
			this.selectedClear()
		});
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
					this.viewModel = this.getViewLink('', 'pods', items[0].namespace, items[0].name)
					if(this.currentitems.length ===0) this.currentitems = Object.assign({},this.viewModel)
					this.isShowSidebar = true
				} else {
					if(this.currentitems.title !== this.viewModel.title) {
						if(this.currentitems.length ===0) this.isShowSidebar = false
						else {
							this.viewModel = Object.assign({},this.currentitems)
							this.currentitems = []
							this.isShowSidebar = true
							this.$refs.selectableTable.selectRow(this.selectIndex)
						}
					} else {
						this.isShowSidebar = false
						this.$refs.selectableTable.clearSelected()
					}
				}
			} else {
				this.currentitems = []
				this.isShowSidebar = false
				this.$refs.selectableTable.clearSelected()
			}
		},
		//  status 필터링
		onChangeStatus() {
			let selectedStatus = this.selectedStatus;
			this.items = this.origin.filter(el => {
				if(selectedStatus.includes("etc")) {
					return (selectedStatus.length === 0) || selectedStatus.includes(el.status.value) || !(this.allStatus.includes(el.status.value));
				}
				return (selectedStatus.length === 0) || selectedStatus.includes(el.status.value);
			});
			this.totalItems = this.items.length;
			this.currentPage = 1
		},
		// 조회
		query_All() {
			this.isBusy = true;
			this.loadMetrics();
			this.$axios.get(this.getApiUrl("","pods",this.selectedNamespace))
					.then((resp) => {
						this.items = [];
						resp.data.items.forEach(el => {
							this.items.push({
								name: el.metadata.name,
								namespace: el.metadata.namespace,
								ready: this.toReady(el.status, el.spec),
								containers: this.toContainers(el.status),
								restartCount: el.status.containerStatuses ? el.status.containerStatuses.map(x => x.restartCount).reduce((accumulator, currentValue) => accumulator + currentValue, 0) : 0,
								controller: this.getController(el),
								status: this.toStatus(el.metadata.deletionTimestamp, el.status),
								creationTimestamp: this.getElapsedTime(el.metadata.creationTimestamp),
								node: this.getNode(el),
								qos: el.status.qosClass,
							});
						});
						this.origin = this.items;
						this.onFiltered(this.items);
						this.onChangeStatus()
					})
					.catch(e => { this.msghttp(e);})
					.finally(()=> { this.isBusy = false;});
		},
		selectedClear() {
			this.selectedStatus = [];
			this.query_All()
		},
		getNode(el) {
			return {
				"name" : el.spec.nodeName ? el.spec.nodeName : "",
				"group" : "",
				"rs" : "nodes"
			}
		},
		getController(el) {
			let version
			let group
			let gr = ""
			if ( el.metadata.ownerReferences ) {
				let or = el.metadata.ownerReferences[0]
				version = or.apiVersion.split('/')
				if (version.length>1) {
					group = version[0]
				}
				if (or.kind === "Node") {
					gr = "Cluster"
				} else {
					gr = "Workload"
				}
				let rs;
				let len = or.kind.length
				if(or.kind[len-1] === 's') rs = (or.kind).toLowerCase() + 'es'
				else rs = (or.kind).toLowerCase() + 's'
				return {
					"name" : or.name,
					"group" : group ||"",
					"kind" : or.kind,
					"rs" : rs,
					"spaceKind" : this.onKind(or.kind),
					"gr" : gr,
				}
			} else {
				return ""
			}
		},
		onKind(kind) {
			if (kind === "StatefulSet") {
				return "Stateful Set"
			} else if (kind === "CronJob") {
				return "Cron Job"
			} else if (kind === "DaemonSet") {
				return "Daemon Set"
			} else if (kind === "ReplicaSet") {
				return "Replica Set"
			} else {
				return kind
			}
		},
		// check pod's containers
		toContainers(status) {
			let initContainerStatuses = []
			let containerStatuses = []
			let style = ""
			if ( status.initContainerStatuses ) {
				initContainerStatuses = status.initContainerStatuses.map(x => {
					if ( !x.ready ) style = "badge-warning"
					else style = "badge-secondary"
					return Object.assign(x, {"style": style})
				})

			}
			if ( status.containerStatuses ) {
				containerStatuses = status.containerStatuses.map(x => {
					if ( !x.ready ) {
						style = "badge-warning"
						if ( x.state[Object.keys(x.state)].reason === "Completed") style = "badge-secondary"
					}
					else style = "badge-success"
					return Object.assign(x, {"style": style})
				})
			}
			return {
				"initContainerStatuses": initContainerStatuses,
				"containerStatuses": containerStatuses,
			}
		},
		// load pod's Metrics
		async loadMetrics() {
			this.metricsItems = [];
			let resp = await this.$axios.get(this.getApiUrl("metrics.k8s.io","pods",this.selectedNamespace))
			if (!resp) return
			this.metricsItems = resp.data.items
		},
		onFiltered(filteredItems) {
			let status = { running:0, pending:0, failed:0, terminating:0, crashLoopBackOff:0, imagePullBackOff:0, completed:0, containerCreating:0, etc:0 }

			filteredItems.forEach(el=> {
				if(el.status.value === "Running") status.running++;
				else if(el.status.value === "Pending") status.pending++;
				else if(el.status.value === "Terminating") status.terminating++;
				else if(el.status.value === "CrashLoopBackOff") status.crashLoopBackOff++;
				else if(el.status.value === "ImagePullBackOff") status.imagePullBackOff++;
				else if(el.status.value === "Completed") status.completed++;
				else if(el.status.value === "Failed") status.failed++;
				else if(el.status.value === "ContainerCreating") status.containerCreating++;
				else status.etc++;
			});

			this.optionsStatus[0].text = status.running >0 ? `Running (${status.running})`: "Running";
			this.optionsStatus[1].text = status.pending >0 ? `Pending (${status.pending})`: "Pending";
			this.optionsStatus[2].text = status.terminating >0 ? `Terminating (${status.terminating})`: "Terminating";
			this.optionsStatus[3].text = status.crashLoopBackOff >0 ? `CrashLoopBackOff (${status.crashLoopBackOff})`: "CrashLoopBackOff";
			this.optionsStatus[4].text = status.imagePullBackOff >0 ? `ImagePullBackOff (${status.imagePullBackOff})`: "ImagePullBackOff";
			this.optionsStatus[5].text = status.completed >0 ? `Completed (${status.completed})`: "Completed";
			this.optionsStatus[6].text = status.containerCreating >0 ? `ContainerCreating (${status.containerCreating})`: "ContainerCreating";
			this.optionsStatus[7].text = status.failed >0 ? `Failed (${status.failed})`: "Failed";
			this.optionsStatus[8].text = status.etc >0 ? `etc (${status.etc})`: "etc";

			this.totalItems = filteredItems.length;
			this.currentPage = 1
		},
	},
	beforeDestroy(){
		this.$nuxt.$off('navbar-context-selected')
	}
}
</script>
<style scoped></style>
