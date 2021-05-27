<template>
	<div class="content-wrapper">
		<div class="content-header">
			<div class="container-fluid">
				<c-navigator :group="'Custom Resource / '+group"></c-navigator>
				<div class="row mb-2">
					<!-- title & search -->
					<div class="col-sm"><h1 class="m-0 text-dark"><span class="badge badge-info mr-2">{{ kind.charAt(0).toUpperCase() }}</span>{{ kind }}</h1></div>
					<div v-if="$route.query.scope === 'Namespaced' || $route.query.ns === 'true'" class="col-sm-2"><b-form-select v-model="selectedNamespace" :options="namespaces()" size="sm" @input="query_All(); selectNamespace(selectedNamespace);"></b-form-select></div>
					<div class="col-sm-2 float-left">
						<div class="input-group input-group-sm" >
							<b-form-input id="txtKeyword" v-model="keyword" class="form-control float-right" placeholder="Search"></b-form-input>
							<div class="input-group-append"><button type="submit" class="btn btn-default" @click="query_All"><i class="fas fa-search"></i></button></div>
						</div>
					</div>
					<!-- button -->
					<div class="col-sm-1 text-right">
						<b-button variant="primary" size="sm" @click="$router.push(`/create?context=${currentContext()}&group=customresource&crd=${kind}&type=cr&url=${url}`)">Create</b-button>
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
										<div class="text-center text-success" style="margin:150px 0">
											<b-spinner type="grow" variant="success" class="align-middle mr-2"></b-spinner>
											<span class="align-middle text-lg">Loading...</span>
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
import jp from 'jsonpath'
export default {
	components: {
		"c-navigator": { extends: VueNavigator },
		"c-view": { extends: VueView }
	},
	data() {
		return {
			selectedNamespace: "",
			keyword: "",
			filterOn: ["name"],
			fields: [
				{ key: "name", label: "Name", sortable: true },
				{ key: "creationTimestamp", label: "Age", sortable: true },
			],
			isBusy: false,
			origin: [],
			items: [],
			currentItems:[],
			columnOpt: [],
			selected: [],
			selectIndex: 0,
			currentPage: 1,
			totalItems: 0,
			isShowSidebar: false,
			viewModel:{},
			group: "",
			plural: "",
			kind: "",
			columnsName: [],
			columnsPath: [],
			colList: [],
			url : "",
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
			localStorage.setItem('columns_cr_'+this.kind,this.selected)
		},
		group() {
			this.$nuxt.$emit('crCol_up',this.group)
		},
	},
	created() {
		this.$nuxt.$on('sideCRD_click',() =>{
			this.getInit()
		})
		this.$nuxt.$on("navbar-context-selected", (_) => {
			if(this.$route.query.isSide) {
				this.getInit()
			} else {
				if(this.$route.query.scope ==='Namespaced') this.fields.splice(1,0,{ key: "namespace", label: "Namespace", sortable: true })
				this.group = this.$route.query.group
				this.plural = this.$route.query.plural
				this.kind = this.$route.query.kind
				this.getColumns(this.$route.query.columnsName,this.$route.query.columnsPath)
				this.getCol()
			}
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
					this.viewModel = this.getViewLink(this.group, this.plural, items[0].namespace, items[0].name)
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
			if(!this.resources()[this.group]) return
			let resource = this.resources()[this.group][this.plural]
			this.url = resource.groupVersion+','+resource.name
			this.$axios.get(this.getApiUrl(this.group,this.plural,this.selectedNamespace))
				.then((resp) => {
					this.items = [];
					let idx = 0
					resp.data.items.forEach(el => {
						this.items[idx] = {
							name: el.metadata.name,
							namespace: el.metadata.namespace,
							creationTimestamp: this.getElapsedTime(el.metadata.creationTimestamp),
						}
						this.colList.forEach(e => {
							let val = jp.query(el,e[Object.keys(e)])
							val = val.join(', ')
							this.items[idx][Object.keys(e)] = val
						})
						idx ++
					});
					this.origin = this.items;
					this.onFiltered(this.items);
				})
				.catch(e => { this.msghttp(e);})
				.finally(()=> { this.isBusy = false;});
		},
		onFiltered(filteredItems) {
			this.totalItems = filteredItems.length;
			this.currentPage = 1
		},
		getColumns(name,path) {
			if(!name || !path) return

			let listName = name.split(',')
			let listPath = path.split(',')
			listName.forEach(el => {
				this.fields.splice(-1,0,{key: el, label: el, sortable: true})
				this.columnsName.push(el)
			})
			listPath.forEach(el => {
				this.columnsPath.push(el)
			})
			for(let i=0;i<listName.length;i++) {
				this.colList.push({
					[this.columnsName[i]] : this.columnsPath[i].substring(1)
				})
			}
		},
		getInit() {
			this.fields = [{ key: "name", label: "Name", sortable: true },{ key: "creationTimestamp", label: "Age", sortable: true }]
			if(this.$route.query.ns === 'true') this.fields.splice(1,0,{ key: "namespace", label: "Namespace", sortable: true })
			let col = []
			let g = (this.$route.query.gV).split('/')
			let n = this.$route.query.n
			let url = this.getApiUrl("apiextensions.k8s.io",'customresourcedefinitions')+'/'+n+'.'+g[0]
			this.$axios.get(url)
				.then((resp) => {
					this.group = g[0]
					this.plural = n
					this.kind = this.$route.query.k
					col = this.conv(this.getPrinterColumns(resp.data.spec))
					if(col) this.getColumns(col.name,col.path)
				}).finally(() => this.getCol())
		},
		getCol() {
			this.columnOpt = Object.assign([],this.fields)

			if(localStorage.getItem('columns_cr_'+this.kind)) {
				this.selected = (localStorage.getItem('columns_cr_'+this.kind)).split(',')
			} else {
				this.fields.forEach(el => {
					this.selected.push(el.label)
				})
			}
			this.selectedNamespace = this.selectNamespace()
			this.query_All()
		},
		getPrinterColumns(spec) {
			const columns = spec.versions.find(a => this.getVersion(spec) === a.name)?.additionalPrinterColumns ??
				spec.additionalPrinterColumns?.map(({ JSONPath, ...rest}) => ({ ...rest, jsonPath: JSONPath})) ?? [];
			return columns
				.filter(column => column.name !== 'Age')
				.filter(column => column.jsonPath ? true : !column.priority );
		},
		conv(col) {
			if(col.length === 0) return

			let name = []
			let path = []
			col.forEach(el => {
				name.push(el.name)
				path.push(el.jsonPath)
			})
			return {
				name: name.join(','),
				path: path.join(','),
			}
		},
		getVersion(spec) {
			return spec.versions[0]?.name ?? spec.version
		},
	},
	beforeDestroy(){
		this.$nuxt.$off('navbar-context-selected')
		this.$nuxt.$off('sideCRD_click')
	}
}
</script>
<style scoped>label {font-weight: 500;}</style>