<template>
<!-- content-wrapper -->
<div class="content-wrapper">

	<div class="content-header">
		<div class="container-fluid">
			<c-navigator group="Workload"></c-navigator>
			<div class="row mb-2">
				<div class="col-sm-2"><h1 class="m-0 text-dark">Deployments</h1></div>
				<div class="col-sm-2">
					<b-form-select v-model="selectedNamespace" :options="namespaces()" size="sm" @input="query"></b-form-select>
				</div>
				<div class="col-sm-2 float-left">
					<div class="input-group input-group-sm" >
						<input type="text" name="table_search" class="form-control float-right" placeholder="Search">
						<div class="input-group-append">
							<button type="submit" class="btn btn-default"><i class="fas fa-search"></i></button>
						</div>
					</div>
				</div>
				<div class="col-sm-6 text-right dropdown">
					<b-dropdown text="Create" variant="primary" size="sm">
						<b-dropdown-item class="dropdown-item"><nuxt-link :to="{ path:'/create', query:{ context: currentContext(), group: 'Workload', crd: 'Depoyment' } }">from Yaml</nuxt-link></b-dropdown-item>
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
								<nuxt-link :to="{ path:'/view', query:{ context: currentContext(), group: 'Workload', crd: 'Deployment', name: data.item.name, url: `deployment/namespace/${data.item.namespace}/name/${data.item.name}`}}">{{ data.value }}</nuxt-link>
							</template>
							<template v-slot:cell(images)="data">
								<ul class="list-unstyled mb-0">
									<li v-for="image in data.item.images" v-bind:key="image">{{ image }}</li>
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
			selectedNamespace: " ",
			fields: [
				{ key: "name", label: "이름", sortable: true },
				{ key: "namespace", label: "네임스페이스", sortable: true  },
				{ key: "pods", label: "파드", sortable: true  },
				{ key: "creationTimestamp", label: "생성시간" },
				{ key: "images", label: "이미지"  },
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
			let ns = this.$data.selectedNamespace;

			if (ns) {
				this.$data.items = this.$data.origin.filter(el => {
					return (ns == " " || el.namespace == ns);
				});
			} else {
				this.$data.items = this.$data.origin;
			}
		},
		query_All() {
			this.$data.selected = [];

			axios.get(`${this.dashboardUrl()}/api/v1/deployment/${this.$data.selectedNamespace}?sortBy=d,creationTimestamp&context=${this.currentContext()}`)
				.then((resp) => {
					let data = []
					resp.data.deployments.forEach(el => {
						data.push({
							name: el.objectMeta.name,
							namespace: el.objectMeta.namespace,
							images: el.containerImages,
							pods: `${el.pods.current}/${el.pods.desired}`,
							creationTimestamp: this.$root.getTimestampString(el.objectMeta.creationTimestamp)
						});
					});
					this.$data.origin = data;
					this.$data.items = data;
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
