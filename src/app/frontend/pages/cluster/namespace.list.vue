<template>
	<div class="content-wrapper">
		<section class="content-header">
			<div class="container-fluid">
				<c-navigator group="Cluster"></c-navigator>
				<div class="row mb-2">
					<!-- title & search -->
					<div class="col-sm"><h1 class="m-0 text-dark"><span class="badge badge-info mr-2">N</span>Namespaces</h1></div>
					<div class="col-sm-2 float-left">
						<div class="input-group input-group-sm" >
							<b-form-input id="txtKeyword" v-model="keyword" class="form-control float-right" placeholder="Search"></b-form-input>
							<div class="input-group-append"><button type="submit" class="btn btn-default" @click="query_All"><i class="fas fa-search"></i></button></div>
						</div>
					</div>
					<!-- button -->
					<div class="col-sm-1 text-right">
						<b-button variant="primary" size="sm" @click="$router.push(`/create?context=${currentContext()}&group=Cluster&crd=Namespace`)">Create</b-button>
					</div>
				</div>
			</div>
		</section>

		<section class="content">
			<div class="container-fluid">
				<!-- search & filter -->
				<div class="d-flex">
					<div class="p-2">
						<b-form-group class="mb-0 font-weight-light overflow-auto">
							<button type="submit" class="btn btn-default btn-sm float-left mr-2" @click="selectedClear">All</button>
							<b-form-checkbox-group v-model="selectedPhase" :options="optionsPhase" button-variant="light" font="light" switches size="sm" @input="onChangePhase" class="float-left"></b-form-checkbox-group>
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
								<b-table hover selectable show-empty select-mode="single" @row-selected="onRowSelected" @sort-changed="currentPage=1" ref="grdSheet1" :items="items" :fields="fields" :filter="keyword" :filter-included-fields="filterOn" @filtered="onFiltered" :current-page="currentPage" :per-page="itemsPerPage" :busy="isBusy" class="text-sm">
									<template #table-busy>
										<div class="text-center text-success lh-vh-50">
											<b-spinner type="grow" variant="success" class="align-middle mr-2"></b-spinner>
											<span class="text-lg align-middle">Loading...</span>
										</div>
									</template>
									<template v-slot:cell(labels)="data">
										<span v-for="(value, name) in data.item.labels" v-bind:key="name" class="label">{{ name }}={{ value }}</span>
									</template>
									<template v-slot:cell(phase)="data">
										<span v-bind:class="{ 'text-success': data.item.phase=='Active' }">{{ data.item.phase }}</span>
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
			selectedPhase: [],
			optionsPhase: [
				{ text: "Active", value: "Active" },
				{ text: "Terminating", value: "Terminating" }
			],
			keyword: "",
			filterOn: ["name"],
			fields: [],
			fieldsAll: [
				{ key: "name", label: "Name", sortable: true },
				{ key: "labels", label: "Labels", sortable: true },
				{ key: "creationTimestamp", label: "Age", sortable: true, formatter: this.getElapsedTime },
				{ key: "phase", label: "Status", sortable: true },
			],
			isBusy: false,
			origin: [],
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
			this.selectedClear()
		} );
		if(this.currentContext()) this.$nuxt.$emit("navbar-context-selected");
	},
	methods: {
		onChangePhase() {
			let selectedPhase = this.selectedPhase;
			this.items = this.origin.filter(el => {
				return (selectedPhase.length === 0) || selectedPhase.includes(el.phase.status);
			});
			this.totalItems = this.items.length;
			this.currentPage = 1
		},
		onRowSelected(items) {
			this.isShowSidebar = (items && items.length > 0)
			if (this.isShowSidebar) this.viewModel = this.getViewLink('', 'namespaces', items[0].namespace, items[0].name)
		},
		// 조회
		query_All() {
			this.isBusy = true;
			let c = false
			let sel = this.selectNamespace()
			this.$axios.get(this.getApiUrl("","namespaces"))
				.then((resp) => {
					this.items = [];
					let nsList = [{ value: "", text: "All Namespaces" }];
					resp.data.items.forEach(el => {
						nsList.push({ value: el.metadata.name, text: el.metadata.name });
						if(el.metadata.name === sel) c = true
						this.items.push({
							name: el.metadata.name,
							labels: el.metadata.labels,
							phase: el.status.phase,
							creationTimestamp: el.metadata.creationTimestamp
						});
					});
					if(!c) this.selectNamespace(nsList[0].value)
					this.namespaces(nsList);
					this.origin = this.items;
					this.onFiltered(this.items);
					this.onChangePhase()
				})
				.catch(e => { this.msghttp(e);})
				.finally(()=> { this.isBusy = false;});
		},
		selectedClear() {
			this.selectedPhase = [];
			this.query_All()
		},
		//  status 필터링
		onFiltered(filteredItems) {
			let status = { active:0, terminating:0 }

			filteredItems.forEach(el=> {
				if(el.phase.status === "Active") status.active++;
				if(el.phase.status === "Terminating") status.terminating++;
			});

			this.optionsPhase[0].text = status.active >0 ? `Active (${status.active})`: "Active";
			this.optionsPhase[1].text = status.terminating >0 ? `Terminating (${status.terminating})`: "Terminating";

			this.totalItems = filteredItems.length;
			this.currentPage = 1
		}
	},
	beforeDestroy(){
		this.$nuxt.$off('navbar-context-selected')
	}
}
</script>
