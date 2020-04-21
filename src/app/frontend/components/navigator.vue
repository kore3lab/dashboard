<template>
	<div class="row mb-2">
		<div class="col-sm-12">
			<ol class="breadcrumb">
				<li class="breadcrumb-item"><nuxt-link to="/">Home</nuxt-link></li>
				<li class="breadcrumb-item active">{{ group }}</li>
			</ol>
		</div>
	</div>
</template>
<script>
export default {
	props:["group"]
}
</script>
