<template>
<!-- content-wrapper -->
<div class="content-wrapper">

	<div class="content-header">
		<div class="container-fluid">
			<c-navigator group="Cluster"></c-navigator>
			<div class="row mb-2">
				<div class="col-sm-2"><h1 class="m-0 text-dark">Namespaces</h1></div>
				<!-- 검색 (검색어) -->
				<div class="col-sm-2 float-left">
					<div class="input-group input-group-sm" >
						<b-form-input id="txtKeyword" v-model="keyword" class="form-control float-right" placeholder="Search"></b-form-input>
						<div class="input-group-append">
							<button type="submit" class="btn btn-default" @click="query_All"><i class="fas fa-search"></i></button>
						</div>
					</div>
				</div><!--//END -->
				<!-- 버튼 -->
				<div class="col-sm-8 text-right dropdown">
					<b-button variant="primary" size="sm" @click="$router.push(`/create?context=${currentContext()}&group=Cluster&crd=Namespace`)">Create</b-button>
				</div><!--//END -->
			</div>
		</div>
	</div>

	<section class="content">
	<div class="container-fluid">
		<!-- 검색 (상태) -->
		<div class="row mb-2">
			<div class="col-11">
				<b-form-group class="mb-0 font-weight-light">
					<button type="submit" class="btn btn-default btn-sm" @click="query_All">All</button>
					<b-form-checkbox-group v-model="selectedPhase" :options="optionsPhase" button-variant="light"  font="light" buttons size="sm" @input="onChangePhase"></b-form-checkbox-group>
				</b-form-group>
			</div>
			<div class="col-1 text-right "><span class="text-sm align-middle">Total : {{ totalItems }}</span></div>
		</div><!--//END -->
		<!-- GRID-->
		<div class="row">
			<div class="col-12">
				<div class="card">
					<div class="card-body table-responsive p-0">
						<b-table id="list" hover :items="items" :fields="fields" :filter="keyword" :filter-included-fields="filterOn" @filtered="onFiltered" :current-page="currentPage" :per-page="$config.itemsPerPage" :busy="isBusy" class="text-sm">
							<template #table-busy>
								<div class="text-center text-success" style="margin:150px 0">
									<b-spinner type="grow" variant="success" class="align-middle mr-2"></b-spinner>
									<span class="align-middle text-lg">Loading...</span>
								</div>
							</template>
							<template v-slot:cell(name)="data">
								<nuxt-link :to="{ path:'/view', query:{ context: currentContext(), group: 'Cluster', crd: 'Namespace', name: data.item.name, url: `api/v1/namespaces/${data.item.name}`}}">{{ data.value }}</nuxt-link>
							</template>
							<template v-slot:cell(labels)="data">
								<ul class="list-unstyled mb-0">
									<li v-for="(value, name) in data.item.labels" v-bind:key="name"><span class="badge badge-secondary font-weight-light text-sm mb-1">{{ name }}:{{ value }}</span></li>
								</ul>
							</template>
						</b-table>
					</div>
					<b-pagination v-model="currentPage" :per-page="$config.itemsPerPage" :total-rows="totalItems" size="sm" align="center"></b-pagination>
				</div>
			</div>
		</div><!-- //GRID-->
	</div>
	</section>
</div>
</template>
<script>
import axios		from "axios"
import VueNavigator from "@/components/navigator"
export default {
	components: {
		"c-navigator": { extends: VueNavigator }
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
			fields: [
				{ key: "name", label: "이름", sortable: true },
				{ key: "phase", label: "STATUS", sortable: true },
				{ key: "creationTimestamp", label: "AGE" }
			],
			isBusy: false,
			items: [],
			currentPage: 1,
			totalItems: 0
		}
	},
	layout: "default",
	created() {
		this.$nuxt.$on("navbar-context-selected", (ctx) => this.query_All() );
		if(this.currentContext()) this.$nuxt.$emit("navbar-context-selected");
	},
	methods: {
		//  Phase 필터링
		onChangePhase() {
			let selectedPhase = this.selectedPhase;
			this.items = this.origin.filter(el => {
				return (this.selectedPhase.length == 0) || this.selectedPhase.includes(el.phase);
			});
			this.totalItems = this.items.length;
			this.currentPage = 1
		},
		// 조회
		query_All() {
			this.isBusy = true;
			axios.get(`${this.backendUrl()}/raw/clusters/${this.currentContext()}/api/v1/namespaces`)
				.then((resp) => {
					this.items = [];
					console.log("resp.data == ", resp.data)
					resp.data.items.forEach(el => {
						this.items.push({
							name: el.metadata.name,
							// labels: el.metadata.labels,
							phase: el.status.phase,
							creationTimestamp: this.$root.getTimestampString(el.metadata.creationTimestamp)
						});
					});
					this.origin = this.items;
					this.onFiltered(this.items);
				})
				.catch(e => { this.msghttp(e);})
				.finally(()=> { this.isBusy = false;});
		},
		//  status 필터링
		onFiltered(filteredItems) {
			let status = { active:0, terminating:0 }

			filteredItems.forEach(el=> {
				if(el.phase == "Active") status.active++;
				if(el.phase == "Terminating") status.terminating++;
			});

			this.optionsPhase[0].text = status.active >0 ? `Active (${status.active})`: "Active";
			this.optionsPhase[1].text = status.terminating >0 ? `Terminating (${status.terminating})`: "Terminating";

			this.totalItems = filteredItems.length;
			this.currentPage = 1
		},
		toEndpointList(p) {
			let list = [];
			for(let i =0; i < p.ports.length; i++) list.push(`${p.host}:${p.ports[i]["port"]} ${p.ports[i].protocol}`)
			return list;
		},
	},
	beforeDestroy(){
		this.$nuxt.$off('navbar-context-selected')
	}
}
</script>
<style>label {font-weight: 500;}</style>
