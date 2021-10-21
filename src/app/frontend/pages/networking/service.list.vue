<template>
  <div class="content-wrapper">
		<section class="content-header">
			<div class="container-fluid">
				<c-navigator group="Networking"></c-navigator>
				<div class="row mb-2">
					<!-- title & search -->
					<div class="col-sm"><h1 class="m-0 text-dark"><span class="badge badge-info mr-2">S</span>Services</h1></div>
					<div class="col-sm-2"><b-form-select v-model="selectedNamespace" :options="namespaces()" size="sm" @input="query_All(); selectNamespace(selectedNamespace);"></b-form-select></div>
					<div class="col-sm-2 float-left">
						<div class="input-group input-group-sm" >
							<b-form-input id="txtKeyword" v-model="keyword" class="form-control float-right" placeholder="Search"></b-form-input>
							<div class="input-group-append"><button type="submit" class="btn btn-default" @click="query_All"><i class="fas fa-search"></i></button></div>
						</div>
					</div>
					<!-- button -->
					<div class="col-sm-1 text-right">
						<b-button variant="primary" size="sm" @click="$router.push(`/create?context=${currentContext()}&group=Networking&crd=Service`)">Create</b-button>
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
								<b-table hover selectable show-empty select-mode="single" @row-selected="onRowSelected" @sort-changed="currentPage=1" ref="grdSheet1" :items="items" :fields="fields" :filter="keyword" :filter-included-fields="filterOn" @filtered="onFiltered" :current-page="currentPage" :per-page="itemsPerPage" :busy="isBusy" class="text-sm">
									<template #table-busy>
										<div class="text-center text-success lh-vh-50">
											<b-spinner type="grow" variant="success" class="align-middle mr-2"></b-spinner>
											<span class="text-lg align-middle">Loading...</span>
										</div>
									</template>
									<template v-slot:cell(ports)="data">
										<ul class="list-unstyled mb-0">
											<li v-for="value in data.value" v-bind:key="value">{{ value }}</li>
										</ul>
									</template>
									<template v-slot:cell(externalIPs)="data">
										<ul class="list-unstyled mb-0">
											<li v-for="(value, idx) in data.value" v-bind:key="idx">{{ value }}</li>
										</ul>
									</template>
									<template v-slot:cell(selector)="data">
										<span v-for="(value, name) in data.item.selector" v-bind:key="name" class="border-box background">{{ name }}={{ value }}</span>
									</template>
									<template v-slot:cell(status)="data">
										<span v-bind:key="data.value" v-bind:class="{'text-success': data.value=='Active', 'text-warning': data.value=='Pending'}">{{ data.value }}</span>
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
			selectedNamespace: "",
			keyword: "",
			filterOn: ["name"],
			fields: [],
			fieldsAll: [
				{ key: "name", label: "Name", sortable: true },
				{ key: "namespace", label: "Namespace", sortable: true },
				{ key: "type", label: "Type", sortable: true },
				{ key: "clusterIP", label: "Cluster IP", sortable: true },
				{ key: "ports", label: "Ports", sortable: true, formatter: this.formatPorts },
				{ key: "externalIPs", label: "External IP", sortable: true, formatter: this.formatExternalIPs },
				{ key: "selector", label: "Selector", sortable: true },
				{ key: "creationTimestamp", label: "Age", sortable: true, formatter: this.getElapsedTime },
				{ key: "status", label: "Status", sortable: true, formatter: this.formatStatus }
			],
			isBusy: false,
			items: [],
			status: true,
			itemsPerPage: this.$storage.global.get("itemsPerPage",10),
			currentPage: 1,
			totalItems: 0,
			activeColor: 'green',
			pendingColor: 'orange',
			isShowSidebar: false,
			viewModel:{},
		}
	},
	watch: {
		itemsPerPage(n) {
			this.$storage.global.set("itemsPerPage",n)
		}
	},
	layout: "default",
	created() {
		this.$nuxt.$on("navbar-context-selected", (ctx) => {
			this.selectedNamespace = this.selectNamespace()
			this.query_All()
		});
		if(this.currentContext()) this.$nuxt.$emit("navbar-context-selected");
	},
	methods: {
		onRowSelected(items) {
			this.isShowSidebar = (items && items.length > 0)
			if (this.isShowSidebar) this.viewModel = this.getViewLink('', 'services', items[0].namespace, items[0].name)
		},
		// 조회
		query_All() {
			this.isBusy = true;
			this.$axios.get(this.getApiUrl("","services",this.selectedNamespace))
					.then((resp) => {
						this.items = [];
						resp.data.items.forEach(el => {
							this.items.push({
								name: el.metadata.name,
								namespace: el.metadata.namespace,
								type: el.spec.type,
								clusterIP: el.spec.clusterIP,
								ports: el.spec.ports,
								externalIPs: el.spec.externalIPs || el.spec.externalName,
								selector: el.spec.selector,
								creationTimestamp: el.metadata.creationTimestamp,
								status: el.status
							});
						});
						this.onFiltered(this.items);
					})
					.catch(e => { this.msghttp(e);})
					.finally(()=> { this.isBusy = false;});
		},
		onFiltered(filteredItems) {
			this.totalItems = filteredItems.length;
			this.currentPage = 1
		},
		formatStatus(status, key, item) {
			if (item.type == "LoadBalancer") {
				if(status.loadBalancer && status.loadBalancer.ingress) {
					return "Active";
				} else {
					return "Pending";
				}
			} else {
				return "Active";
			}
		},
		formatExternalIPs(externalIPs, key, item) {
			if (item.type == "LoadBalancer") {
				return item.status.loadBalancer.ingress? [item.status.loadBalancer] : ["-"];
			} else if(item.type == "ExternalName") {
				return [externalIPs];
			} else {
				return (externalIPs? externalIPs : ["-"])
			}
		},
		formatPorts(ports, key, item) {
			let list = [];
			if (ports) {
				for(let i =0; i < ports.length; i++) {
					if (item.type == "NodePort" || item.type == "LoadBalancer") {
						list.push(`${ports[i].port}:${ports[i].nodePort}/${ports[i].protocol}`)
					} else if(ports[i].targetPort === ports[i].port) {
						list.push(`${ports[i].port}/${ports[i].protocol}`)
					} else {
						list.push(`${ports[i].port}:${ports[i].targetPort}/${ports[i].protocol}`)
					}
				}
			}
			return list;
		}
	},
	beforeDestroy(){
		this.$nuxt.$off('navbar-context-selected')
	}
}
</script>
