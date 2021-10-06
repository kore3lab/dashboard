<template>
  <div class="content-wrapper" v-bind:style="{paddingBottom: contentPadding()}">
		<section class="content-header">
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
									<template v-slot:cell(message)="data">
										<span v-bind:class="{'text-danger':data.item.type=='Warning'}">{{ data.value }}</span>
									</template>
									<template v-slot:cell(involvedObj)="data">
										<span>{{ data.value.kind }} : <a href="#" @click="viewModel=getViewLink(data.value.group, data.value.resource, data.item.namespace, data.value.name); isShowSidebar=true;">{{ data.value.name }}</a></span>
									</template>
									<template v-slot:cell(source)="data">
										<span v-for="(val, idx) in data.value" v-bind:key="idx" >{{ val }} </span>
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
			filterOn: ["involvedObj"],
			fields: [],
			fieldsAll: [
				{ key: "type", label: "Type", sortable: true },
				{ key: "message", label: "Message", sortable: true },
				{ key: "namespace", label: "Namespace", sortable: true },
				{ key: "involvedObj", label: "Involved Object", formatter: this.getResource },
				{ key: "source", label: "Source", formatter: this.formatSource },
				{ key: "count", label: "Count", sortable: true },
				{ key: "lastSeen", label: "Last Seen", sortable: true, formatter: this.getElapsedTime },
				{ key: "creationTimestamp", label: "Age", sortable: true, formatter: this.getElapsedTime },
			],
			isBusy: false,
			items: [],
			itemsPerPage: this.$storage.global.get("itemsPerPage",10),
			currentPage: 1,
			totalItems: 0,
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
			if (this.isShowSidebar) this.viewModel = this.getViewLink('', 'events', items[0].namespace, items[0].name)
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
								source: el.source,
								count: el.count,
								lastSeen: el.lastTimestamp,
								creationTimestamp: el.metadata.creationTimestamp
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
		formatSource(source, key, item) {
			return Object.entries(source).map(([_, value]) => `${value}`);
		}
	},
	beforeDestroy(){
		this.$nuxt.$off('navbar-context-selected')
	}
}
</script>
