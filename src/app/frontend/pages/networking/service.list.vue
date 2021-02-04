<template>
<!-- content-wrapper -->
<div class="content-wrapper">

	<div class="content-header">
		<div class="container-fluid">
			<c-navigator group="Networking"></c-navigator>
			<div class="row mb-2">
				<div class="col-sm-2"><h1 class="m-0 text-dark">Services</h1></div>
				<!-- 검색 (namespace) -->
				<div class="col-sm-2">
					<b-form-select v-model="selectedNamespace" :options="namespaces()" size="sm" @input="query_All"></b-form-select>
				</div><!--//END -->
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
				<div class="col-sm-6 text-right dropdown">
					<b-button variant="primary" size="sm" @click="$router.push(`/create?context=${currentContext()}&group=Networking&crd=Service`)">Create</b-button>
				</div><!--//END -->
			</div>
		</div>
	</div>

	<section class="content">
	<div class="container-fluid">
		<!-- 검색 -->
		<div class="row mb-2">
			<div class="col-12 text-right "><span class="text-sm align-middle">Total : {{ totalItems }}</span></div>
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
								<nuxt-link :to="{ path:'/view', query:{ context: currentContext(), group: 'Networking', crd: 'Service', name: data.item.name, url: `api/v1/namespaces/${data.item.namespace}/services/${data.item.name}`}}">{{ data.value }}</nuxt-link>
							</template>
							<template v-slot:cell(selector)="data">
								<ul class="list-unstyled mb-0">
									<li v-for="(value, name) in data.item.selector" v-bind:key="name"><span class="badge badge-secondary font-weight-light text-sm">{{ name }}:{{ value }}</span></li>
								</ul>
							</template>
							<template v-slot:cell(internalEndpoints)="data">
								<ul class="list-unstyled mb-0">
									<li v-for="value in data.item.internalEndpoints" v-bind:key="value" class="mr-1 font-weight-light text-md">{{ value }}</li>
								</ul>
							</template>
							<template v-slot:cell(externalEndpoints)="data">
								<ul class="list-unstyled mb-0">
									<li v-for="value in data.item.externalEndpoints" v-bind:key="value" class="mr-1 font-weight-light text-md">{{ value }}</li>
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
			selectedNamespace: "",
			keyword: "",
			filterOn: ["name"],
			fields: [
				{ key: "name", label: "이름", sortable: true },
				{ key: "namespace", label: "네임스페이스", sortable: true },
				{ key: "selector", label: "셀렉터", sortable: true  },
				{ key: "cpuRequests", label: "클러스터IP", sortable: true  },
				{ key: "internalEndpoints", label: "내부 엔드포인트", sortable: true  },
				{ key: "externalEndpoints", label: "외부 엔드포인트", sortable: true  },
				{ key: "creationTimestamp", label: "생성시간" }
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
		// 조회
		query_All() {
			this.isBusy = true;
			axios.get(`${this.backendUrl()}/raw/clusters/${this.currentContext()}/api/v1/namespaces/${this.$data.selectedNamespace}/services`)
				.then((resp) => {
					this.items = [];
					resp.data.items.forEach(el => {
						this.items.push({
							name: el.metadata.name,
							namespace: el.metadata.namespace,
							// selector: el.selector,
							// internalEndpoints: this.toEndpointList(el.internalEndpoint),
							// externalEndpoints: el.externalEndpoints,
							// clusterIP: el.clusterIP,
							creationTimestamp: this.$root.getTimestampString(el.metadata.creationTimestamp)
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
