<template>
	<!-- content-wrapper -->
	<div class="content-wrapper">

		<!-- Content Header (Page header) -->
		<div class="content-header">
			<div class="container-fluid">
				<c-navigator :group="'Custom Resource / '+ crdQuery.group"></c-navigator>
				<div class="row mb-2">
					<div class="col-sm-10">
						<h1 class="m-0 text-dark"><span class="badge badge-info mr-2">C</span>Create {{ crdQuery.name }}</h1>
					</div>
					<div class="col-sm-2 text-right">
						<b-button variant="primary" size="sm"  @click="onCreate">Create</b-button>
						<b-button variant="secondary" size="sm" @click="$router.go(-1)">Back</b-button>
					</div>
				</div>
			</div>
		</div>
		<section class="content">
			<div class="container-fluid">
				<div class="row">
					<div class="col-md-12">
						<div class="card">
							<div class="card-header"></div>
							<c-aceeditor class="card-body" v-model="raw" v-on:error="onError" style="min-height: calc(100vh - 270px)"></c-aceeditor>
						</div>
					</div>
				</div>
			</div>
		</section>
	</div>
	<!-- /.content-wrapper -->

</template>
<script>

import VueAceEditor 	from "@/components/aceeditor"
import VueNavigator 	from "@/components/navigator"

export default {
	components: {
		"c-aceeditor": { extends: VueAceEditor },
		"c-navigator": { extends: VueNavigator }
	},
	data() {
		return {
			badge: this.$route.query.crd ? this.$route.query.crd.substring(0,1): "P",
			url: this.$route.query.url,
			raw: {},
			crdQuery: {
				group: this.$route.query.group,
				name: this.$route.query.name,
				crd: this.$route.query.crd,
				version: this.$route.query.version
			}
		}
	},
	layout: "default",
	created() {
		this.$nuxt.$on("context-selected", () => {
			// crd spec 읽어서 template 동적 생성
			this.$axios.get(`${this.getApiUrl("apiextensions.k8s.io","customresourcedefinitions")}/${this.crdQuery.crd}`)
				.then(resp => {
					// properties 재귀호출 함수
					let fn_el_loop =function(props, d) {
						if (!props) return
						for(let nm in props) {
							if (props[nm]["type"] == "object") {
								d[nm] = {};
								fn_el_loop(props[nm]["properties"], d[nm])
							} else if (props[nm]["type"] == "integer") {
								d[nm] = 0
							} else if (props[nm]["type"] == "boolean") {
								d[nm] = true
							} else {
								d[nm] = ""
							}
						}
					}

					// template 생성
					let d = {
						kind: resp.data.spec.names.kind,
						apiVersion: `${resp.data.spec.group}/${this.crdQuery.version}`,
						metadata: {
							name: "default"
						}
					}
					if(resp.data.spec.scope == "Namespaced") d.metadata.namespace = "default"
					resp.data.spec.versions.find(el => {
						if(el.name === this.crdQuery.version && el.schema.openAPIV3Schema) {
							if(el.schema.openAPIV3Schema.properties) {
								fn_el_loop(el.schema.openAPIV3Schema.properties, d)
							}
							
						}
					})
					if(!d.kind) d.kind = resp.data.spec.names.kind;
					if(!d.apiVersion) d.apiVersion = `${resp.data.spec.group}/${this.crdQuery.version}`;
					if(!d.metadata["name"]) d.metadata["name"] = "default";
					if(resp.data.spec.scope == "Namespaced") d.metadata.namespace = "default"

					this.raw = d;

				})
		});
		if(this.currentContext()) this.$nuxt.$emit("context-selected");
	},
	beforeDestroy(){
		this.$nuxt.$off("context-selected");
	},
	methods: {
		onCreate() {
			if(!this.raw.metadata.namespace) this.raw.metadata.namespace = 'default'
			this.$axios.post(`/raw/clusters/${this.currentContext()}`, this.raw)
				.then( resp => {
					this.origin = Object.assign({}, resp.data);
					this.raw = resp.data;
					this.toast("Apply OK", "info");
					this.$router.go(-1);
				})
				.catch(e => { this.msghttp(e);});
		},
		onError(error) {
			this.toast(error.message, "danger");
		}
	}
}
</script>
