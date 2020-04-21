<template>
<!-- content-wrapper -->
<div class="content-wrapper">
	<div class="content-header">
		<div class="container-fluid">
			<div class="row mb-2">
				<div class="col-sm-12">
					<ol class="breadcrumb">
					<li class="breadcrumb-item"><nuxt-link to="/">Home</nuxt-link></li>
					<li class="breadcrumb-item active">{{ group }}</li>
					</ol>
				</div>
			</div>
			<div class="row mb-2">
				<div class="col-sm-10">
					<h1 class="m-0 text-dark">Create {{ crd }}</h1>
				</div>
				<div class="col-sm-2 text-right">
					<b-button variant="primary" size="sm">Apply</b-button>
					<b-button variant="secondary" size="sm" @click="$router.go(-1)">Cancel</b-button>
				</div>
			</div>
		</div>
	</div>
	<section class="content">
	<div class="container-fluid">
		<div class="row">
			<div class="col-md-12">
				<div class="card">
					<div class="card-header"><h3 class="card-title"></h3></div>
					<c-aceeditor class="card-body" style="min-height: calc(100vh - 210px - 60px)" id="wrapEditor"></c-aceeditor>
				</div>
			</div>
		</div>
	</div>
	</section>
</div>
<!-- /.content-wrapper -->

</template>
<script>
import VueAceEditor from "@/components/aceeditor"

export default {
	data() {
		return {
			badge: this.$route.query.crd ? this.$route.query.crd.substring(0,1): "P",
			group: this.$route.query.group ?  this.$route.query.group: "Workload", 
			crd : this.$route.query.crd ?  this.$route.query.crd: "pod",
			raw: { metadata: {}, spec: {} }
		}
	},
	layout: "default",
	components: {
		"c-aceeditor": {
			extends: VueAceEditor
		}
	}
}
</script>
