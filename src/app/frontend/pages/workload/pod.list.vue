<template>
	<div class="content-wrapper">
		<section class="content-header">
			<div class="container-fluid">
				<c-navigator group="Workload"></c-navigator>
				<div class="row mb-2">
					<div class="col-sm"><h1 class="m-0 text-dark"><span class="badge badge-info mr-2">P</span>Pods</h1></div>
					<div class="col-sm-1 text-right">
						<b-button variant="primary" size="sm" @click="$router.push(`/create?context=${currentContext()}&group=Workload&crd=Pod`)">Create</b-button>
					</div>
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
							<button type="submit" class="btn btn-default btn-sm float-left mr-2" @click="selectedClear">All</button>
							<b-form-checkbox-group v-model="selectedStatus" :options="optionsStatus" button-variant="light" font="light" switches size="sm" @input="onChangeStatus" class="float-left"></b-form-checkbox-group>
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
										<div class="text-center">
											<b-dropdown dropdown no-caret variant="link" size="xs" class="m-n2">
												<template #button-content>
													<b-button class="btn btn-tool" variant="link" @click="onClickInitMenu(data.value.status)">
														<i class="fas fa-ellipsis-v" aria-hidden="true"></i>
													</b-button>
												</template>
												<b-dropdown-item v-for="(item,index) in menu" :key="index" @mouseover.native="toggleSubMenu($event, index)" @mouseout.native="toggleSubMenu($event, index)" @click="onClickShowLogs(data.value, item.name, item.type)">
													<b-icon v-bind:icon="item.icon" class="mr-2"></b-icon><span>{{item.title}} <b-icon-caret-right-fill :scale="0.6" v-if="item.children"></b-icon-caret-right-fill></span>
													<div v-if="item.children" class="dropdown-sub-menu" v-show="item.showSubMenu">
														<b-dropdown-item dropleft v-for="(d,index) in item.children" :key="index" @click="onClickShowLogs(data.value, d.name, item.title)"> <span class="badge badge-success"> {{""}}</span><span class="pl-1 text-xs">{{d.title}}</span>
														</b-dropdown-item>
													</div>
												</b-dropdown-item>
											</b-dropdown>
										</div>
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
div .dropdown-sub-menu{ position: absolute; white-space: nowrap; min-width: 8rem; right: 100%; left: auto; transform : translate(0, -60%); background-color: #fff; border: 1px solid rgba(0,0,0,.15); border-radius: .25rem;}
.dropdown-menu {min-width:5rem !important;}
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
				{ key: "status", label: "Status", sortable: true, formatter: this.formatStatus },
				{ key: "menu", label: "" }
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
			viewModel:{},
			menu: [{title:"terminal", icon:"terminal-fill", type:'terminal', visible:true},{title:'log', icon:"card-text", type:'log', visible:true}]
		}
	},
	watch: {
		itemsPerPage(n) {
			this.$storage.global.set("itemsPerPage",n)	// save to localstorage
		}
	},
	layout: "default",
	methods: {
		onClickInitMenu(status){
			const menu = [{title:"terminal", icon:"terminal-fill", type:'terminal', visible:true},{title:'log', icon:"card-text", type:'log', visible:true}]
			if(status.containerStatuses && status.containerStatuses.length > 1){
				const list = [];
				status.containerStatuses.forEach(item =>{list.push({ title: item.name, name: item.name, visible: item.ready})})
				let menuList = [];
				this.menu.forEach(item => {
					menuList.push({ title: item.title, icon: item.icon, showSubMenu : false, visible: true, children : list})
				})
				this.menu = menuList;
			}else{
				menu.forEach((item, index) => {
					menu[index].name = status.containerStatuses[0].name;
					menu[index].visible = status.containerStatuses[0].ready;
				})
				this.menu = menu;
			}
		},
		toggleSubMenu (e, index){
			if(this.menu[index].children) this.menu[index].showSubMenu = e.type === 'mouseover';
		},
		onClickShowLogs(data, name, type) {
			if(name !== undefined){
				let containerList = []
				data.status.containerStatuses.forEach(item =>{
					containerList.push(item.name);
				})
				if(type === 'log') {this.$nuxt.$emit("open-terminal", data.metadata.name, "logs", { metadata: data.metadata, container:name, containers:containerList });}
				else {
					let route = this.$router.resolve({path: "/terminal", query: {termtype: 'container',pod: data.metadata.name, namespace: data.metadata.namespace, cluster: this.currentContext(), container:name}})
					window.open(route.href);
				}
			}
		},
		onRowSelected(items) {
			this.isShowSidebar = (items && items.length > 0)
			if (this.isShowSidebar) this.viewModel = this.getViewLink("", "pods", items[0].namespace, items[0].name)
		},
		//  status 필터링
		onChangeStatus() {
			let selectedStatus = this.selectedStatus;
			this.items = this.origin.filter(el => {
				if(selectedStatus.includes("etc")) {
					return (selectedStatus.length === 0) || selectedStatus.includes(el.status) || !(this.allStatus.includes(el.status));
				}
				return (selectedStatus.length === 0) || selectedStatus.includes(el.status);
			});
			this.totalItems = this.items.length;
			this.currentPage = 1
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
							status: el.status.phase,
							creationTimestamp: el.metadata.creationTimestamp,
							deletionTimestamp: el.metadata.deletionTimestamp,
							node: el.spec.nodeName,
							qos: el.status.qosClass,
							menu: el,
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
		onFilter(keyword) {
			this.keyword=keyword
		},
		onFiltered(filteredItems) {
			let status = { running:0, pending:0, failed:0, terminating:0, crashLoopBackOff:0, imagePullBackOff:0, completed:0, containerCreating:0, etc:0 }

			filteredItems.forEach(el=> {
				if(el.status === "Running") status.running++;
				else if(el.status === "Pending") status.pending++;
				else if(el.status === "Terminating") status.terminating++;
				else if(el.status === "CrashLoopBackOff") status.crashLoopBackOff++;
				else if(el.status === "ImagePullBackOff") status.imagePullBackOff++;
				else if(el.status === "Completed") status.completed++;
				else if(el.status === "Failed") status.failed++;
				else if(el.status === "ContainerCreating") status.containerCreating++;
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
		formatStatus(status, key, item)  {
			return this.toPodStatus(item.deletionTimestamp, item.containers);
		},
		formatContainers(status, key, item) {
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
		}
	}
}
</script>