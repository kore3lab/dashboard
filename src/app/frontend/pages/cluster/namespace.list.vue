<template>
<!-- content-wrapper -->
<div class="content-wrapper">

	<div class="content-header">
		<div class="container-fluid">
			<c-navigator group="Cluster"></c-navigator>
			<div class="row mb-2">
				<div class="col-sm-2"><h1 class="m-0 text-dark">Namespaces</h1></div>
				<div class="col-sm-2 float-left">
					<div class="input-group input-group-sm" >
						<input type="text" name="table_search" class="form-control float-right" placeholder="Search">
						<div class="input-group-append">
							<button type="submit" class="btn btn-default"><i class="fas fa-search"></i></button>
						</div>
					</div>
				</div>
				<div class="col-sm-8 text-right dropdown">
					<b-dropdown text="Create" variant="primary" size="sm">
						<b-dropdown-item class="dropdown-item"><nuxt-link :to="{ path:'/create', query:{ context: currentContext(), group: 'Cluster', crd: 'Namespace' } }">from Yaml</nuxt-link></b-dropdown-item>
					</b-dropdown>
				</div>
			</div>
		</div>
	</div>

	<section class="content">
	<div class="container-fluid">
		<div class="row mb-2">
			<div class="col-12">
				<!-- 조회조건 -->
				<b-form-group class="mb-0 font-weight-light">
					<button type="submit" class="btn btn-default btn-sm" @click="query_All();">All</button>
					<b-form-checkbox-group v-model="selected" :options="conditions" button-variant="light"  font="light" buttons size="sm" @input="query"></b-form-checkbox-group>
				</b-form-group>
				<!-- //조회조건 -->
			</div>
		</div>
		<div class="row">
			<div class="col-12">
				<!-- GRID-->
				<div class="card">
					<div class="card-body table-responsive p-0">
						<b-table id="list" hover :items="items" :fields="fields" class="text-sm">
							<template v-slot:cell(name)="data">
								<nuxt-link :to="{ path:'/view', query:{ context: currentContext(), group: 'Cluster', crd: 'Namespace', name: data.item.name, url: `namespace/name/${data.item.name}`}}">{{ data.value }}</nuxt-link>
							</template>
							<template v-slot:cell(labels)="data">
								<ul class="list-unstyled mb-0">
									<li v-for="(value, name) in data.item.labels" v-bind:key="name"><span class="badge badge-secondary font-weight-light text-sm mb-1">{{ name }}:{{ value }}</span></li>
								</ul>
							</template>
						</b-table>
					</div>
				</div>
				<!-- //GRID-->
			</div>
		</div>
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
			selected: [],
			conditions: [
				{ text: "Active", value: "Active" },
				{ text: "Terminating", value: "Terminating" }
			],
			fields: [
				{ key: "name", label: "이름", sortable: true },
				{ key: "labels", label: "레이블", sortable: true },
				{ key: "phase", label: "단계", sortable: true },
				{ key: "creationTimestamp", label: "생성시간" }
			],
			items: [],
			origin: []
		}
	},
	layout: "default",
	created() {
		this.$nuxt.$on("navbar-context-selected", (ctx) => this.query_All() );
		if(this.currentContext()) this.$nuxt.$emit("navbar-context-selected");
	},
	methods: {
		// 상태 조건에 따른 조회 (namespace, status)
		query() {
			let selected = this.$data.selected;

			if ( selected.length == 0) {
				this.$data.items = this.$data.origin;
			} else {
				this.$data.items = this.$data.origin.filter(el => {
					return selected.includes(el.phase);
				});
			}
		},
		query_All() {
			this.$data.selected = [];

			axios.get(`${this.dashboardUrl()}/api/v1/namespace?sortBy=d,creationTimestamp&context=${this.currentContext()}`)
				.then((resp) => {
					let data = []
					resp.data.namespaces.forEach(el => {
						data.push({
							name: el.objectMeta.name,
							labels: el.objectMeta.labels,
							phase: el.phase,
							creationTimestamp: this.$root.getTimestampString(el.objectMeta.creationTimestamp)
						});
					});
					this.$data.origin = data;
					this.$data.items = data;
				})
				.catch((error) => {
					this.$root.toast(error.message, "danger");
				});
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
