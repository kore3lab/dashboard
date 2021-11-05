<template>
<div class="content-wrapper">
	<section class="content-header">&nbsp;</section>
	<section class="content">
		<div class="error-page">
			<h2 class="headline" v-bind:class="{'text-danger': error.statusCode!=404,'text-warning': error.statusCode==404 }">{{error.statusCode}}</h2>
			<div class="error-content">
				<h3 v-if="error.statusCode == 404"><i class="fas fa-exclamation-triangle text-warning" ></i> Page not found</h3>
				<h3 v-else><i class="fas fa-exclamation-triangle text-danger" ></i> An error occurred.</h3>
				<p v-if="error.statusCode == 404">We could not find the page you were looking for. you may <nuxt-link to="/">return to dashboard</nuxt-link>.</p>
				<p v-else>Sorry, we had some technical problems during your last operation. you may <nuxt-link to="/">return to dashboard</nuxt-link>.</p>
				<a v-if="error.statusCode != 404" v-b-toggle href="#message" @click.prevent>show detail</a>
				<b-collapse v-if="error.statusCode" id="message" class="text-secondary">{{error.message}}</b-collapse>
			</div>
		</div>
	</section>
</div>
</template>
<script>
export default {
	props: ["error"],
	layout: "default"
}
</script>