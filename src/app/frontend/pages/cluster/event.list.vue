<template>
	<div class="content-wrapper">
		<div class="content-header">
			<div class="container-fluid">
				<c-navigator group="Administrator"></c-navigator>
				<div class="row mb-2">
					<!-- title & search -->
					<div class="col-sm"><h1 class="m-0 text-dark"><span class="badge badge-info mr-2">E</span>Events</h1></div>
					<div class="col-sm-2"><b-form-select v-model="selectedNamespace" :options="namespaces()" size="sm" @input="query_All(); selectNamespace(selectedNamespace);"></b-form-select></div>
					<div class="col-sm-2 float-left">
						<div class="input-group input-group-sm" >
							<b-form-input id="txtKeyword" v-model="keyword" class="form-control float-right" placeholder="Search"></b-form-input>
							<div class="input-group-append"><button type="submit" class="btn btn-default" @click="query_All"><i class="fas fa-search"></i></button></div>
						</div>
					</div>
					<!-- button -->
					<div class="col-sm-1 text-right">
						<b-button variant="primary" size="sm" @click="$router.push(`/create?context=${currentContext()}&group=Configuration&crd=ResourceQuota`)">Create</b-button>
					</div>
				</div>
			</div>
		</div>

		<section class="content">
			<div class="container-fluid">
				<!-- count -->
				<div class="row mb-2">
					<div class="col-12 text-right "><span class="text-sm align-middle"><i class="fas fa-question-circle" id="limitMark"></i> Total : {{ totalItems }}</span></div>
					<b-tooltip target="limitMark" placement="left" boundary="window" boundary-padding="0">limited to 1000</b-tooltip>
				</div>
				<!-- GRID-->
				<div class="row">
					<div class="col-12">
						<div class="card">
							<div class="card-body table-responsive p-0">
								<b-table id="list" hover selectable show-empty select-mode="single" @row-selected="onRowSelected" @sort-changed="onSortChanged()" ref="selectableTable" :items="items" :fields="fields" :filter="keyword" :filter-included-fields="filterOn" @filtered="onFiltered" :current-page="currentPage" :per-page="$config.itemsPerPage" :busy="isBusy" class="text-sm">
									<template #table-busy>
										<div class="text-center text-success lh-vh-50">
											<b-spinner type="grow" variant="success" class="align-middle mr-2"></b-spinner>
											<span class="text-lg align-middle">Loading...</span>
										</div>
									</template>
									<template #empty="scope">
										<h4 class="text-center">does not exist.</h4>
									</template>
									<template v-slot:cell(message)="data">
										<span v-bind:class="data.item.style">{{ data.value }}</span>
									</template>
									<template v-slot:cell(involvedObj)="data">
										<a href="#" @click="viewModel=getViewLink(data.item.controllers.g, data.item.controllers.k, data.value.namespace, data.value.name); isShowSidebar=true;">{{ data.value.kind }}: {{ data.value.name }}</a>
									</template>
									<template v-slot:cell(source)="data">
										<span v-for="(val, idx) in data.value" v-bind:key="idx" >{{ val }} </span>
									</template>
									<template v-slot:cell(lastSeen)="data">
										{{ data.value.str }}
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
			keyword: "",
			filterOn: ["involvedObj"],
			fields: [
				{ key: "type", label: "Type", sortable: true },
				{ key: "message", label: "Message", sortable: true },
				{ key: "namespace", label: "Namespace", sortable: true },
				{ key: "involvedObj", label: "Involved Object" },
				{ key: "source", label: "Source" },
				{ key: "count", label: "Count", sortable: true },
				{ key: "lastSeen", label: "Last Seen", sortable: true },
				{ key: "creationTimestamp", label: "Age", sortable: true },
			],
			isBusy: false,
			items: [],
			currentItems:[],
			selectIndex: 0,
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
			this.query_All()
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
					this.viewModel = this.getViewLink('', 'events', items[0].namespace, items[0].name)
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
			this.isBusy = true;
			this.$axios.get(this.getApiUrl("","events",this.selectedNamespace,'','limit=1000'))
					.then((resp) => {
						this.items = [];
						resp.data.items.forEach(el => {
							this.items.push({
								type: el.type,
								message: el.message,
								name: el.metadata.name,
								namespace: el.metadata.namespace,
								involvedObj: el.involvedObject,
								controllers: this.getController(el.involvedObject),
								source: Object.entries(el.source).map(([_, value]) => `${value}`),
								count: el.count,
								lastSeen: this.getElapsedTime(el.lastTimestamp),
								creationTimestamp: this.getElapsedTime(el.metadata.creationTimestamp),
								style: el.type === 'Warning' ? 'text-danger' : '',
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
	},
	beforeDestroy(){
		this.$nuxt.$off('navbar-context-selected')
	}
}
</script>
<style scoped>label {font-weight: 500;}</style>
