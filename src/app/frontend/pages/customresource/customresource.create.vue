<template>
	<!-- content-wrapper -->
	<div class="content-wrapper">

		<!-- Content Header (Page header) -->
		<div class="content-header">
			<div class="container-fluid">
				<c-navigator :group="'Custom Resource / '+ crdQuery.group"></c-navigator>
				<div class="row mb-2">
					<div class="col-sm-10">
						<h1 class="m-0 text-dark"><span class="badge badge-info mr-2 text-capitalize">{{badge}}</span>{{ title }}</h1>
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
			badge: (this.$route.query.crd && this.$route.query.crd.length >0) ? this.$route.query.crd.substring(0,1): "C",
			raw: {}
		}
	},
	computed: {
		crdQuery: {
			get () {
				return {
					group: this.$route.query.group,		//ex. networking.istio.io
					crd: this.$route.query.crd,			//ex. virtualservices
					name: this.$route.query.name,		//ex. VirtualService
					version: this.$route.query.version	//ex. v1beta1
				}
			},
		},
	},
	layout: "default",
	created() {
		this.title = `Create ${["A", "E", "I", "O", "U", "a", "e", "i", "o", "u"].includes(this.badge)?"an":"a"} ${this.crdQuery.name}`;
		this.$nuxt.$on("context-selected", () => {
			// crd spec 읽어서 template 동적 생성
			this.$axios.get(`${this.getApiUrl("apiextensions.k8s.io","customresourcedefinitions")}/${this.crdQuery.crd}.${this.crdQuery.group}`)
				.then(resp => {
					// properties 재귀호출 함수
					let fn_el_loop =function(props, d) {
						if (!props) return
						for(const nm in props) {
							let el;
							if (props[nm]["type"] == "object") {
								el = {}
								fn_el_loop(props[nm]["properties"], el)
							} else if (props[nm]["type"] == "integer") {
								el = 0
							} else if (props[nm]["type"] == "number") {
								el = 0
							} else if (props[nm]["type"] == "boolean") {
								el = true
							} else if (props[nm]["type"] == "array") {
								el = [];
								if(props[nm].items["properties"]) fn_el_loop(props[nm].items["properties"], el);
							} else if (props[nm]["type"] == "string") {
								el = ""
							}
							if (el && Array.isArray(d)) d.push(el); else  d[nm] = el;
						}
					}
					// template 생성
					let d = {
						kind: resp.data.spec.names.kind,
						apiVersion: `${this.crdQuery.group}/${this.crdQuery.version}`,
						metadata: {
							name: "default"
						}
					}
					if(resp.data.spec.scope == "Namespaced") d.metadata.namespace = "default"

					// structure (https://kubernetes.io/docs/reference/using-api/deprecation-guide/#customresourcedefinition-v122)
					let properties;
					if( resp.data.apiVersion == "apiextensions.k8s.io/v1") {
						const version = resp.data.spec.versions.find(el => el.name= this.crdQuery.version);
						if(version && version.schema && version.schema.openAPIV3Schema ) properties = version.schema.openAPIV3Schema.properties;
					} else {
						if(resp.data.spec.validation && resp.data.spec.openAPIV3Schema)	properties = resp.data.spec.validation.openAPIV3Schema.properties;
					}

					if(properties) fn_el_loop(properties, d);
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
