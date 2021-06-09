<template>
	<div class="content-wrapper">
		<div class="content-header">
			<div class="container-fluid">
				<c-navigator group="Storage"></c-navigator>
				<div class="row mb-2">
					<!-- title & search -->
					<div class="col-sm"><h1 class="m-0 text-dark"><span class="badge badge-info mr-2">S</span>Storage Classes</h1></div>
					<div class="col-sm-2 float-left">
						<div class="input-group input-group-sm" >
							<b-form-input id="txtKeyword" v-model="keyword" class="form-control float-right" placeholder="Search"></b-form-input>
							<div class="input-group-append"><button type="submit" class="btn btn-default" @click="query_All"><i class="fas fa-search"></i></button></div>
						</div>
					</div>
					<!-- button -->
					<div class="col-sm-1 text-right">
						<b-button variant="primary" size="sm" @click="$router.push(`/create?context=${currentContext()}&group=Storage&crd=Storage Class`)">Create</b-button>
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
									<template v-slot:cell(name)="data">
										{{ data.value }}
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
			keyword: "",
			filterOn: ["name"],
			fields: [
				{ key: "name", label: "Name", sortable: true },
				{ key: "provisioner", label: "Provisioner", sortable: true  },
				{ key: "reclaimPolicy", label: "Reclaim Policy", sortable: true  },
				{ key: "onDefault", label: "Default"},
				{ key: "creationTimestamp", label: "Age", sortable: true }
			],
			isBusy: false,
			items: [],
			currentItems:[],
			columnOpt: [],
			selected: [],
			selectIndex: 0,
			currentPage: 1,
			totalItems: 0,
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
			localStorage.setItem('columns_storageclass',this.selected)
		}
	},
	created() {
		this.columnOpt = Object.assign([],this.fields)
		this.$nuxt.$on("navbar-context-selected", (ctx) => {
			if(localStorage.getItem('columns_storageclass')) {
				this.selected = (localStorage.getItem('columns_storageclass')).split(',')
			} else {
				this.fields.forEach(el => {
					this.selected.push(el.label)
				})
			}
			this.query_All()
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
					this.viewModel = this.getViewLink('storage.k8s.io', 'storageclasses', items[0].namespace, items[0].name)
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
			this.$axios.get(this.getApiUrl("storage.k8s.io","storageclasses"))
				.then((resp) => {
					this.items = [];
					resp.data.items.forEach(el => {
						this.items.push({
							name: el.metadata.name,
							provisioner: el.provisioner,
							reclaimPolicy: el.reclaimPolicy,
							onDefault: this.isDefault(el),
							creationTimestamp: this.getElapsedTime(el.metadata.creationTimestamp)
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
		isDefault(el) {
			if (el.metadata.annotations && el.metadata.annotations["storageclass.kubernetes.io/is-default-class"] === "true")
			{
				return 'Yes'
			}
			return ""
		},
	},
	beforeDestroy(){
		this.$nuxt.$off('navbar-context-selected')
	}
}
</script>
<style scoped>label {font-weight: 500;}</style>
