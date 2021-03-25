<template>
	<!-- content-wrapper -->
	<div class="content-wrapper">

		<!-- Content Header (Page header) -->
		<div class="content-header">
			<div class="container-flui">
				<c-navigator :group="group"></c-navigator>
				<div class="row mb-2">
					<div class="col-sm-10">
						<h1 class="m-0 text-dark"><span class="badge badge-info mr-2">{{ badge }}</span>Create {{ crd }}</h1>
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

import axios			from "axios"
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
			group: this.$route.query.group ?  this.$route.query.group: "Workload",
			crd : this.$route.query.crd ?  this.$route.query.crd: "pod",
			raw: { metadata: {}, spec: {} },
			template: null
		}
	},
	layout: "default",
	created() {
		if(this.currentContext()) this.$nuxt.$emit("navbar-context-selected");
	},
	mounted() {
		try {
			let filename = this.crd.toLowerCase().replaceAll(" ", "");
			this.template = require(`~/assets/template/${filename}.json`);
			this.raw = Object.assign({}, this.template);
		} catch (ex) {
			console.log(`can't find "${this.crd}" template on ~/assets/template`);
		}
	},
	beforeDestroy(){
		this.$nuxt.$off('navbar-context-selected')
	},
	methods: {
		onCreate() {
			axios.post(`${this.backendUrl()}/raw/clusters/${this.currentContext()}`, this.raw)
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
		},
	}
}
</script>
