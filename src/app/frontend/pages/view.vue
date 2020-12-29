<template>
<!-- content-wrapper -->
<div class="content-wrapper">

	<!-- Content Header (Page header) -->
	<div class="content-header">
		<div class="container-fluid">
			<c-navigator :group="group"></c-navigator>
			<div class="row mb-2">
				<div class="col-sm-10">
					<h1 class="m-0 text-dark"><span class="badge badge-info mr-2">{{ badge }}</span>{{ crd }}<small class="text-muted ml-2">/ {{ name }}<small><i class="fas fa-trash ml-2"></i></small></small></h1>
				</div>
				<div class="col-sm-2 text-right">
					<b-button variant="secondary" size="sm" @click="$router.go(-1)">Cancel</b-button>
				</div>
			</div>
		</div>
	</div>

	<!-- Main content -->
	<section class="content">
	<div class="container-fluid">
		<b-tabs content-class="col-md-12" card>
			<!-- summary tab -->
			<b-tab title="Summary" active>
				<div class="row">
					<div class="col-md-12">
						<div class="card">
							<div class="card-header"><h3 class="card-title">Meta data</h3></div>
							<div class="card-body">
								<dl class="row">
									<dt class="col-sm-2">Annotations<i class="fas fa-edit ml-2 text-secondary"></i></dt>
									<dd class="col-sm-10 text-truncate">
										<ul class="list-unstyled mb-0">
											<li v-for="(value, name) in raw.metadata.annotations" v-bind:key="name"><span class="badge badge-secondary font-weight-light text-sm mb-1">{{ name }}:{{ value }}</span></li>
										</ul>
									</dd>
									<dt class="col-sm-2">Labels<i class="fas fa-edit ml-2 text-secondary"></i></dt>
									<dd class="col-sm-10">
										<ul class="list-unstyled mb-0">
											<li v-for="(value, name) in raw.metadata.labels" v-bind:key="name"><span class="badge badge-secondary font-weight-light text-sm mb-1">{{ name }}:{{ value }}</span></li>
										</ul>
									</dd>
									<dt class="col-sm-2">Create at</dt><dd class="col-sm-10">{{ raw.metadata.creationTimestamp }}</dd>
									<dt class="col-sm-2">UID</dt><dd class="col-sm-10">{{ raw.metadata.uid }}</dd>
								</dl>
							</div>
						</div>
					</div>
				</div>
				<div class="row" v-if="crd != 'Storage Class'">
					<div class="col-md-12">
						<div class="card">
							<div class="card-header">
								<h3 class="card-title" v-if="crd=='Config Map' || crd=='Secret'">Data</h3>
								<h3 class="card-title" v-else-if="crd=='Service Account'">Secrets</h3>
								<h3 class="card-title" v-else-if="crd=='Role' || crd=='Cluster Role'">Rules</h3>
								<h3 class="card-title" v-else-if="crd=='Role Binding' || crd=='Cluster Role Binding'">Subjects,RoleRef</h3>
								<h3 class="card-title" v-else>Specification</h3>
							</div>
							<div class="card-body" ref="wrapSpec"></div>
						</div>
					</div>
				</div>
			</b-tab>
			<!-- yaml tab -->
			<b-tab title="Yaml">
				<div class="row mb-2">
					<div class="col-sm-12 text-right">
						<b-button variant="primary" size="sm" @click="onPatch">Apply</b-button>
						<b-button variant="secondary" size="sm" @click="onReset">Reset</b-button>
					</div>
				</div>
				<div class="row">
					<div class="col-md-12">
						<div class="card">
							<div class="card-header"></div>
							<c-aceeditor class="card-body" v-model="raw" v-on:error="onError" style="min-height: calc(100vh - 210px - 60px)"></c-aceeditor>
						</div>
					</div>
				</div>
			</b-tab>
		</b-tabs>
	</div>
	</section>
	<!-- /.content -->
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
			name : this.$route.query.name ? this.$route.query.name: "httpbin",
			origin: { metadata: {}, spec: {} },
			raw: { metadata: {}, spec: {} }
		}
	},
	layout: "default",
	created() {
		this.$nuxt.$on("navbar-context-selected", (ctx) => this.query() );
		if(this.currentContext()) this.$nuxt.$emit("navbar-context-selected");
	},
	beforeDestroy(){
		this.$nuxt.$off('navbar-context-selected')
	},
	methods: {
		// 조회
		query() {
			let params = this.$route.query;
			axios.get(`${this.dashboardUrl()}/api/v1/_raw/${params.url}?context=${this.currentContext()}`)
				.then( resp => {
					this.origin = Object.assign({}, resp.data);
					this.raw = resp.data;
					let spec = null; 

					if(this.$data.crd == "Config Map" || this.$data.crd == "Secret") spec = resp.data.data;
					else if(this.$data.crd == "Storage Class") spec = null; // 무시
					else if(this.$data.crd == "Role" || this.$data.crd == "Cluster Role") spec = resp.data.rules;
					else if(this.$data.crd == "Role Binding" || this.$data.crd == "Cluster Role Binding") spec = { subjects: resp.data.subjects, roleRef: resp.data.roleRef} ;
					else if(this.$data.crd == "Service Account") spec = resp.data.secrets;
					else spec = resp.data.spec;

					if (spec) this.$jsonTree.create(spec, this.$refs["wrapSpec"]);

				}).catch( error => {
					this.toast(error.message, "danger");
				});
		},
		// apply
		onPatch() {
			axios.put(`${this.backendUrl()}/api/_raw/${this.currentContext()}`, this.raw)
				.then( resp => {
					this.origin = Object.assign({}, resp.data);
					this.raw = resp.data;
				}).catch( error => {
					this.toast(error.message, "danger");
				});
		},
		// reset
		onError(error) {
			// this.mesbox(error.message);
			this.toast(error.message, "danger");
		},
		// reset
		onReset() {
			this.raw = Object.assign({}, this.origin);
		}
	}
}

</script>
