<template>
	<div class="content-wrapper">
		<div class="content-header">
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
		</div>

		<section class="content">
			<div class="container-fluid">
				<!-- count & filter -->
				<div class="row mb-2">
					<div class="col-11">
						<b-form-group class="mb-0 font-weight-light">
							<button type="submit" class="btn btn-default btn-sm" @click="onChangePhase('All')">All</button>
							<b-form-checkbox-group v-model="selectedPhase" :options="optionsPhase" button-variant="light"  font="light" buttons size="sm" @input="onChangePhase"></b-form-checkbox-group>
						</b-form-group>
					</div>
					<div class="col-1 text-right "><span class="text-sm align-middle">Total : {{ totalItems }}</span></div>
				</div>
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
										<a href="#" @click="sidebar={visible:true, name:data.item.name, src:`${getApiUrl('','namespaces')}/${data.item.name}`}">{{ data.value }}</a>
									</template>
									<template v-slot:cell(labels)="data">
										<ul class="list-unstyled mb-0">
											<li v-for="(value, name) in data.item.labels" v-bind:key="name"><span class="badge badge-secondary font-weight-light text-sm mb-1">{{ name }}:{{ value }}</span></li>
										</ul>
									</template>
									<template v-slot:cell(phase)="data">
										<div v-bind:class="data.item.phase.style">{{ data.item.phase.status }}</div>
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
		<b-sidebar v-model="sidebar.visible" width="50em" right shadow no-header>
			<c-view crd="Namespace" group="Cluster" :name="sidebar.name" :url="sidebar.src" @delete="query_All()" @close="sidebar.visible=false"/>
		</b-sidebar>
	</div>
</template>
<script>
import axios		from "axios"
import VueNavigator from "@/components/navigator"
import VueView from "@/pages/view";
export default {
	components: {
		"c-navigator": { extends: VueNavigator },
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
			fields: [
				{ key: "name", label: "Name", sortable: true },
				{ key: "labels", label: "Labels", sortable: true },
				{ key: "creationTimestamp", label: "Age", sortable: true },
				{ key: "phase", label: "Status", sortable: true },
			],
			isBusy: false,
			origin: [],
			items: [],
			currentPage: 1,
			totalItems: 0,
			sidebar: {
				visible: false,
				name: "",
				src: "",
			},
		}
	},
	layout: "default",
	created() {
		this.$nuxt.$on("navbar-context-selected", (ctx) => this.query_All() );
		if(this.currentContext()) this.$nuxt.$emit("navbar-context-selected");
	},
	methods: {
		//  Phase 필터링
		onChangePhase(a) {
			if(a === "All") this.selectedPhase = [];
			let selectedPhase = this.selectedPhase;
			this.items = this.origin.filter(el => {
				return (selectedPhase.length === 0) || selectedPhase.includes(el.phase.status);
			});
			this.totalItems = this.items.length;
			this.currentPage = 1
		},
		// 조회
		query_All() {
			this.isBusy = true;
			axios.get(this.getApiUrl("","namespaces"))
					.then((resp) => {
						this.items = [];
						resp.data.items.forEach(el => {
							this.items.push({
								name: el.metadata.name,
								labels: el.metadata.labels,
								phase: this.onPhase(el.status.phase),
								creationTimestamp: this.getElapsedTime(el.metadata.creationTimestamp)
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
				if(el.phase.status === "Active") status.active++;
				if(el.phase.status === "Terminating") status.terminating++;
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
		// status 확인
		onPhase(phase) {
			if (phase === "Active") {
				return {
					"status" : "Active",
					"style" : "text-success"
				}
			} else {
				return {
					"status" : "Terminating",
					"style" : "text-secondary"
				}
			}
		}
	},
	beforeDestroy(){
		this.$nuxt.$off('navbar-context-selected')
	}
}
</script>
<style scoped>label {font-weight: 500;}</style>
