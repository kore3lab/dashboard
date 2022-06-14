<template>
	<div class="content-wrapper">
		<section class="content-header">
			<div class="container-fluid">
				<c-navigator group="Workload"></c-navigator>
				<div class="row mb-2">
					<div class="col-sm"><h1 class="m-0 text-dark"><span class="badge badge-info mr-2">P</span>Pods <nuxt-link :to="{path:'/create', query: {group:'Workload', crd:'Pod'}}"><b-icon-plus-circle variant="secondary" font-scale="0.7"></b-icon-plus-circle></nuxt-link></h1></div>
				</div>
			</div>
		</section>
		<section class="content">
			<div class="container-fluid">
				<!-- search & total count & items per page  -->
				<c-search-form @input="query_All" @keyword="(k)=>{keyword=k}"/>
				<div class="d-flex">
					<div class="p-2">
						<b-form-group class="mb-0 font-weight-light overflow-auto">
							<button type="submit" class="btn btn-default btn-sm float-left mr-2" @click="onSelectedStatus([])">All</button>
							<b-form-checkbox-group v-model="selectedStatus" :options="optionsStatus" button-variant="light" font="light" switches size="sm" @input="onSelectedStatus" class="float-left"></b-form-checkbox-group>
						</b-form-group>
					</div>
					<div class="ml-auto p-2">
						<b-form inline>
							<c-colums-selector name="grdSheet1" v-model="fields" :fields="fieldsAll" ></c-colums-selector>
							<i class="text-secondary ml-2 mr-2">|</i>
							<b-form-select size="sm" :options="this.var('ITEMS_PER_PAGE')" v-model="itemsPerPage"></b-form-select>
							<span class="text-sm align-middle ml-2">Total : {{ totalItems }}</span>
						</b-form>
					</div>
				</div>
				<!-- GRID-->
				<div class="row">
					<div class="col-12">
						<div class="card">
							<div class="card-body table-responsive p-0">
								<b-table hover selectable show-empty select-mode="single" @sort-changed="currentPage=1" @row-selected="onRowSelected" ref="grdSheet1" :items="items" :fields="fields" :filter="keyword" :filter-included-fields="filterOn" @filtered="onFiltered" :current-page="currentPage" :per-page="itemsPerPage" :busy="isBusy" class="text-sm">
									<template #table-busy>
										<div class="text-center text-success lh-vh-50">
											<b-spinner type="grow" variant="success" class="align-middle mr-2"></b-spinner>
											<span class="text-lg align-middle">Loading...</span>
										</div>
									</template>
									<template v-slot:cell(status)="data">
										<span v-bind:class="data.value.style">{{ data.value.value }}</span>
									</template>
									<template v-slot:cell(containers)="data">
										<span v-for="(value, idx) in data.value" v-bind:key="idx" v-bind:class="`badge-${value}`" class="badge font-weight-light text-sm ml-1"> {{" "}}</span>
									</template>
									<template v-slot:cell(controller)="data">
										<a href="#" @click="viewModel=getViewLink(data.value.group, data.value.resource, data.item.namespace, data.value.name); isShowSidebar=true;">{{ data.value.kind }}</a>
									</template>
									<template v-slot:cell(node)="data">
										<a href="#" @click="viewModel=getViewLink('','nodes', '',  data.value); isShowSidebar=true;">{{ data.value }}</a>
									</template>
									<template v-slot:cell(menu)="data">
										<b-dropdown dropdown no-caret no-prefetch variant="link" size="xs" class="m-n2 text-center">
											<template #button-content>
												<b-button class="btn btn-tool" variant="link"><i class="fas fa-ellipsis-v" aria-hidden="true"></i></b-button>
											</template>
											<b-dropdown-item v-for="(item,index) in data.value" :key="index" class="columns-context-menu"  @click="if(item.containers.length==1) openTerminal(item.containers[0])">
												<b-icon v-bind:icon="item.icon" class="mr-2"></b-icon><span>{{ item.title }}<b-icon-caret-right-fill :scale="0.6" v-if="item.containers.length>1"></b-icon-caret-right-fill></span>
												<div v-if="item.containers.length>1" class="sub">
													<b-dropdown-item dropleft v-for="(d,i) in item.containers" :key="i" @click="openTerminal(d)">
														<span :class="{'badge-success': d.ready, 'badge-warning': !d.ready}" class="badge ">&nbsp;</span><span class="pl-1 text-xs" >{{d.title}}</span>
													</b-dropdown-item>
												</div>
											</b-dropdown-item>
										</b-dropdown>
									</template>
								</b-table>
							</div>
							<b-pagination v-model="currentPage" :per-page="itemsPerPage" :total-rows="totalItems" size="sm" align="center"></b-pagination>
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
<style scoped>
.columns-context-menu .sub{ display:none; position: absolute; white-space: nowrap; min-width: 8rem; right: 100%; left: auto; transform : translate(0, -60%); background-color: #fff; border: 1px solid rgba(0,0,0,.15); border-radius: .25rem;}
.columns-context-menu:hover .sub { display: block;}
</style>

<script>
import VueNavigator			from "@/components/navigator"
import VueColumsSelector	from "@/components/list/columnsSelector"
import VueSearchForm		from "@/components/list/searchForm"
import VueView				from "@/pages/view";

export default {
	components: {
		"c-navigator": { extends: VueNavigator },
		"c-colums-selector": { extends: VueColumsSelector},
		"c-search-form": { extends: VueSearchForm},
		"c-view": { extends: VueView }
	},

	data() {
		return {
			dd: {},
			selectedStatus: [],
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
			fields: [],
			fieldsAll: [
				{ key: "name", label: "Name", sortable: true },
				{ key: "namespace", label: "Namespace", sortable: true  },
				{ key: "ready", label: "Ready", sortable: true  },
				{ key: "containers", label: "Containers", sortable: true, formatter: this.formatContainers },
				{ key: "restartCount", label: "Restart", sortable: true  },
				{ key: "controller", label: "Controlled By", sortable: true,  formatter: this.getResource },
				{ key: "node", label: "Node", sortable: true },
				{ key: "qos", label: "QoS", sortable: true  },
				{ key: "creationTimestamp", label: "Age", sortable: true, formatter: this.getElapsedTime },
				{ key: "status", label: "Status", sortable: true},
				{ key: "menu", label: "", sortable: true, formatter: this.formatMenu }
			],
			isBusy: true,
			origin: [],
			items: [],
			currentItems:[],
			selectIndex: 0,
			containerStatuses: [],
			itemsPerPage: this.$storage.global.get("itemsPerPage",10),
			currentPage: 1,
			totalItems: 0,
			isShowSidebar: false,
			viewModel:{}
		}
	},
	watch: {
		itemsPerPage(n) {
			this.$storage.global.set("itemsPerPage",n)	// save to localstorage
		}
	},
	layout: "default",
	created() {
		this.$nuxt.$on("context-selected", () => {this.selectedStatus = []});	// chanag context -> set empty the selected status
	},
	methods: {
		onRowSelected(items) {
			this.isShowSidebar = (items && items.length > 0)
			if (this.isShowSidebar) this.viewModel = this.getViewLink("", "pods", items[0].namespace, items[0].name)
		},
		// 조회
		query_All(d) {
			this.isBusy = true;
			this.$axios.get(this.getApiUrl("","pods", (d && d.namespace) ? d.namespace: this.selectNamespace(), "", d && d.labelSelector? `labelSelector=${d.labelSelector}`: ""))
				.then((resp) => {
					this.items = [];
					resp.data.items.forEach(el => {
						this.items.push({
							name: el.metadata.name,
							namespace: el.metadata.namespace,
							ready: `${el.status.containerStatuses ? el.status.containerStatuses.filter(e => e.ready).length: 0}/${el.spec.containers?el.spec.containers.length: 0}`,
							containers: el.status,
							restartCount: el.status.containerStatuses ? el.status.containerStatuses.map(x => x.restartCount).reduce((accumulator, currentValue) => accumulator + currentValue, 0) : 0,
							controller: el.metadata.ownerReferences?el.metadata.ownerReferences[0]:null,
							status: this.toPodStatus(el.metadata.deletionTimestamp, el.status),
							creationTimestamp: el.metadata.creationTimestamp,
							deletionTimestamp: el.metadata.deletionTimestamp,
							node: el.spec.nodeName,
							qos: el.status.qosClass,
							containerStatuses: el.status.containerStatuses
						});
					});
					this.origin = this.items;
					this.onFiltered(this.items);
					this.selectedStatus = [];
				})
				.catch(e => { this.msghttp(e);})
				.finally(()=> { this.isBusy = false;});
		},
		onFiltered(items) {
			//calc. status count
			let opts = { }
			items.forEach(d=> {
				opts[d.status.value] =  (opts[d.status.value] ? opts[d.status.value]: 0) + 1;
			})
			this.optionsStatus = [];
			this.totalItems = 0;
			for (const [k, v] of Object.entries(opts)) {
				this.optionsStatus.push({ text: `${k} (${v})`, value: k })
				this.totalItems += v;
			}
			this.currentPage = 1
		},
		//  selected status 
		onSelectedStatus(selectedStatus) {
			if (selectedStatus.length ==0) {
				this.items = this.origin;
			} else {
				this.items = this.origin.filter(el => selectedStatus.includes(el.status.value));
			}
			this.selectedStatus = selectedStatus;
			this.totalItems = this.items.length;
			this.currentPage = 1
		},
		formatContainers(status) {
			let list = []
			if ( status.initContainerStatuses ) {
				status.initContainerStatuses.forEach(el=> {
					el.ready ? "secondary" : "warning"
				});
			}
			if ( status.containerStatuses ) {
				status.containerStatuses.forEach(el=> {
					if (el.ready) {
						list.push("success");
					} else {
						list.push((el.state[Object.keys(el.state)].reason === "Completed") ? "secondary":  "warning");
					}
				});
			}
			return list;
		},
		formatMenu(v, key, item) {
			if(!item.containerStatuses) return [];

			let terminals = [];
			let logs = [];
			let containers = [];

			item.containerStatuses.forEach(d => {
				containers.push(d.name);
			});
			item.containerStatuses.forEach(d => {
				logs.push ({
					title: d.name,
					ready: d.ready,
					type: "logs",
					query: {
						metadata: { name: item.name, namespace:item.namespace }, 
						container: d.name, 
						containers: containers
					}
				});
				terminals.push ({
					title: d.name,
					ready: d.ready,
					type: "terminal",
					query: {
						termtype: "container", 
						pod: item.name, 
						namespace: item.namespace, 
						cluster: this.currentContext(), 
						container: d.name
					}
				});
			});

			return [
				{ title:"Logs", icon:"card-text", name:item.name, containers: logs  },
				{ title:"Terminal", icon:"terminal-fill",  name:item.name, containers: terminals }
			];

		},
		openTerminal(d) {
			if(d.type === "logs") {
				this.$nuxt.$emit("open-terminal", d.query.metadata.name, "logs", d.query);
			} else {
				const route = this.$router.resolve( { path: "/terminal", query: d.query } );
				window.open(route.href);
			}
		},
	},
	beforeDestroy(){
		this.$nuxt.$off("context-selected");
	}
}
</script>