<template>
<!-- content-wrapper -->
<div class="content-wrapper">

	<div class="content-header">
		<div class="container-fluid">
			<c-navigator group="Storage"></c-navigator>
			<div class="row mb-2">
				<div class="col-sm-2"><h1 class="m-0 text-dark">Storage Classes</h1></div>
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
						<b-dropdown-item class="dropdown-item"><nuxt-link :to="{ path:'/create', query:{ context: currentContext(), group: 'Storage', crd: 'Storage Class' } }">from Yaml</nuxt-link></b-dropdown-item>
					</b-dropdown>
				</div>
			</div>
		</div>
	</div>

	<section class="content">
	<div class="container-fluid">
		<div class="row">
			<div class="col-12">
				<!-- GRID-->
				<div class="card">
					<div class="card-body table-responsive p-0">
						<b-table id="list" hover :items="items" :fields="fields" class="text-sm">
							<template v-slot:cell(name)="data">
								<nuxt-link :to="{ path:'/view', query:{ context: currentContext(), group: 'Storage', crd: 'Storage Class', name: data.item.name, url: `storageclass/name/${data.item.name}`}}">{{ data.value }}</nuxt-link>
							</template>
							<template v-slot:cell(parameters)="data">
								<ul class="list-unstyled mb-0">
									<li v-for="(value, name) in data.item.parameters" v-bind:key="name"><span class="badge badge-secondary font-weight-light text-sm mb-1">{{ name }}:{{ value }}</span></li>
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
import axios	from "axios"
import VueNavigator from "@/components/navigator"
export default {
	components: {
		"c-navigator": { extends: VueNavigator }
	},
	data() {
		return {
			fields: [
				{ key: "name", label: "이름", sortable: true },
				{ key: "provisioner", label: "제공자", sortable: true  },
				{ key: "parameters", label: "파라미터", sortable: true  },
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
		query_All() {
			this.$data.selected = [];
			axios.get(`${this.dashboardUrl()}/api/v1/storageclass?sortBy=d,creationTimestamp&context=${this.currentContext()}`)
				.then((resp) => {
					this.$data.origin = [];
					resp.data.storageClasses.forEach(el => {
						this.$data.origin.push({
							name: el.objectMeta.name,
							provisioner: el.provisioner,
							parameters: el.parameters,
							creationTimestamp: this.$root.getTimestampString(el.objectMeta.creationTimestamp)
						});
					});
					this.$data.items = this.$data.origin;
				})
				.catch((error) => {
					this.$root.toast(error.message, "danger");
				});
		}
	},
	beforeDestroy(){
		this.$nuxt.$off('navbar-context-selected')
	}
}
</script>
<style>label {font-weight: 500;}</style>
