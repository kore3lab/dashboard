<template>
  <div class="content-wrapper">
		<section class="content-header">
			<div class="container-fluid">
				<c-navigator group="Administrator"></c-navigator>
				<div class="row mb-2">
					<div class="col-sm"><h1 class="m-0 text-dark"><span class="badge badge-info mr-2">E</span>Events</h1></div>
				</div>
			</div>
		</section>
		<section class="content">
			<div class="container-fluid">
				<!-- search & total count & items per page  -->
				<div class="row pb-2">
					<div class="col-sm-10"><c-search-form no-label-selector @input="query_All" @keyword="(k)=>{keyword=k}"/></div>
					<div class="col-sm-2">
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
			keyword: "",
			filterOn: ["involvedObj","message"],
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
	methods: {
		onRowSelected(items) {
			this.isShowSidebar = (items && items.length > 0)
			if (this.isShowSidebar) this.viewModel = this.getViewLink('', 'events', items[0].namespace, items[0].name)
		},
		// 조회
		query_All(d) {
			this.isBusy = true;
			this.$axios.get(this.getApiUrl("","events", (d && d.namespace) ? d.namespace: this.selectNamespace(), "", d && d.labelSelector? `labelSelector=${d.labelSelector}&limit=1000`: "limit=1000"))
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
	}
}
</script>
