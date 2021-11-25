<template>
  <div class="content-wrapper">
		<section class="content-header">
			<div class="container-fluid">
				<c-navigator group="Custom Resource"></c-navigator>
				<div class="row mb-2">
					<div class="col-sm"><h1 class="m-0 text-dark"><span class="badge badge-info mr-2">C</span>Custom Resource Definitions</h1></div>
				</div>
			</div>
		</section>

		<section class="content">
			<div class="container-fluid">
				<!-- search & total count & items per page  -->
				<div class="row pb-2">
					<div class="col-sm-2"><b-form-select v-model="selectedGroup" :options="groupList" size="sm" @input="onChangeGroup"></b-form-select></div>
					<div class="col-sm-2 float-left">
						<div class="input-group input-group-sm" >
							<b-form-input id="txtKeyword" v-model="keyword" class="form-control float-right" placeholder="Search"></b-form-input>
							<div class="input-group-append"><button type="submit" class="btn btn-default" @click="query_All"><i class="fas fa-search"></i></button></div>
						</div>
					</div>
					<div class="col-sm-8">
						<b-form inline class="float-right">
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
									<template v-slot:cell(name)="data">
										<nuxt-link :to="{path: '/customresource/customresource.list', query: {group:data.item.group, crd: data.item.plural, name: data.item.name.name, version: data.item.version }}" class="mr-2">{{ data.value.name }}</nuxt-link>
									</template>
									<template v-slot:cell(labels)="data">
										<ul class="list-unstyled mb-0">
											<li v-for="(value, name) in data.item.labels" v-bind:key="name"><span class="badge badge-secondary font-weight-light text-sm mb-1">{{ name }}={{ value }}</span></li>
										</ul>
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
import VueColumsSelector	from "@/components/list/columnsSelector"
import VueView				from "@/pages/view";

export default {
	components: {
		"c-navigator": { extends: VueNavigator },
		"c-colums-selector": { extends: VueColumsSelector},
		"c-view": { extends: VueView }
	},
	data() {
		return {
			selectedGroup: "",
			keyword: "",
			filterOn: ["name"],
			fields: [],
			fieldsAll: [
				{ key: "name", label: "Name", sortable: true },
				{ key: "group", label: "Group", sortable: true  },
				{ key: "version", label: "Version", sortable: true },
				{ key: "scope", label: "Scope", sortable: true },
				{ key: "creationTimestamp", label: "Age", sortable: true, formatter: this.getElapsedTime },
			],
			isBusy: false,
			origin: [],
			items: [],
			itemsPerPage: this.$storage.global.get("itemsPerPage",10),
			currentPage: 1,
			totalItems: 0,
			groupList: [{value: "", text: "All Groups"}],
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
		this.$nuxt.$on("context-selected", (_) => {
			this.query_All()
		} );
		if(this.currentContext()) this.query_All();
	},
	methods: {
		onRowSelected(items) {
			this.isShowSidebar = (items && items.length > 0)
			if (this.isShowSidebar) this.viewModel = this.getViewLink('apiextensions.k8s.io', 'customresourcedefinitions', items[0].namespace, items[0].name.origin)
		},
		onChangeGroup() {
			let selectedGroup = this.selectedGroup;
			this.items = this.origin.filter(el => {
				return (selectedGroup.length === 0) || selectedGroup.includes(el.group);
			});
			this.totalItems = this.items.length;
			this.currentPage = 1

		},
		// 조회
		query_All() {
			this.isBusy = true;
			this.groupList = [{value: "", text: "All Groups"}]
			this.$axios.get(this.getApiUrl("apiextensions.k8s.io","customresourcedefinitions"))
				.then((resp) => {
					this.items = [];
					resp.data.items.forEach(el => {
						if(!this.groupList.find(d=> {return d.value == el.spec.group})) this.groupList.push({value: el.spec.group, text: el.spec.group});
						this.items.push({
							name: this.getName(el.spec.names.kind,el.metadata.name),
							group: el.spec.group,
							version: this.getVersion(el.spec.versions),
							scope: el.spec.scope,
							plural: el.spec.names.plural,
							creationTimestamp: el.metadata.creationTimestamp,
							printerColumns: this.conv(this.getPrinterColumns(el.spec))
						});
					});
					this.origin = this.items;
					if(!this.groupList.find(d=> {return d.value == this.selectedGroup})) this.selectedGroup = "";
					if(this.selectedGroup !="") this.onChangeGroup();
					this.onFiltered(this.items);
				})
				.catch(e => { this.msghttp(e);})
				.finally(()=> { this.isBusy = false;});
		},
		onFiltered(filteredItems) {
			this.totalItems = filteredItems.length;
			this.currentPage = 1
		},
		getName(name, origin) {
			return {
				name : name,
				origin : origin
			}
		},
		getVersion(v) {
			for(let i=0;i<v.length;i++) {
				if(v[i].storage === true){
					return v[i].name
				}
			}
		},
		getPrinterColumns(spec) {
			const columns = spec.versions.find(a => this.getColumnsVersion(spec) === a.name)?.additionalPrinterColumns ??
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
		getColumnsVersion(spec) {
			return spec.versions[0]?.name ?? spec.version
		},
	},
	beforeDestroy(){
		this.$nuxt.$off("context-selected");
	}
}
</script>
