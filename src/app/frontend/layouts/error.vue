<template>
<div class="content-wrapper">
	<section class="content-header">&nbsp;</section>
	<section class="content">
		<div class="error-page">
			<h2 class="headline" v-bind:class="{'text-danger': error.statusCode!=404,'text-warning': error.statusCode==404 }">{{error.statusCode}}</h2>
			<div class="error-content">
				<h3 v-if="error.statusCode == 404"><i class="fas fa-exclamation-triangle text-warning" ></i> Page not found</h3>
				<h3 v-else><i class="fas fa-exclamation-triangle text-danger" ></i> An error occurred.</h3>
				<p v-if="error.statusCode == 404">{{message}}<br><nuxt-link :to="redirect">return to the page</nuxt-link>.</p>
				<p v-else>{{message}}<br><nuxt-link :to="redirect">return to the page</nuxt-link>.</p>
				<a v-if="error.statusCode != 404" v-b-toggle href="#message" @click.prevent>show detail</a>
				<b-collapse v-if="error.statusCode" id="message" class="text-secondary">{{message}}</b-collapse>
			</div>
		</div>
	</section>
</div>
</template>
<script>
export default {
	props: ["error"],
	computed: {
		redirect: {
			get() {
				return this.error.redirect ? this.error.redirect :"/";
			}
		},
		message: {
			get() {
				if(this.error.statusCode==404) {
					return this.error.message? this.error.message: "We could not find the page you were looking for. you may";
				} else  {
					return this.error.message? this.error.message: "Sorry, we had some technical problems during your last operation. you may";
				}
			}
		}
	},
	layout: "default"
}
</script>