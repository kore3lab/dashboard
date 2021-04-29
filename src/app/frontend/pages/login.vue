<template>
	<div class="content-wrapper min-vh-100 overflow-hidden">
		<div class="row" style="padding-top:10em">
			<div class="col-md-3"></div>
			<div class="col-md-6">
				<h2>Login</h2>
				<p class="text-left">Kore-board service account has a secret with valid token that can be used to login to dashboard. </p>
				<pre class="text-sm border bg-white"><small>$ SECRET="$(kubectl get sa -n kore -l app=kore-board -o jsonpath='{.items[0].secrets[0].name}')"
$ echo "$(kubectl get secret ${SECRET} -n kore -o jsonpath='{.data.token}' | base64 --decode)"</small></pre>
				<div class="row">
					<div class="col-11">
						<b-form-input v-model="token" autofocus size="sm" type="password" placeholder="Enter token"></b-form-input>
					</div>
					<div class="col-1">
						<b-button variant="primary" size="sm" @click="doLogin">Login</b-button>
					</div>
				</div>
			</div>
			<div class="col-md-3"></div>
		</div>
		<div  class="row mt-2">>
			<div class="col-md-3"></div>
			<div class="col-md-6"><p v-show="message!==''" class="text-danger">{{message}}</p></div>
			<div class="col-md-3"></div>
		</div>
	</div>
</template>
<script>
export default {
	components: {
	},
	data() {
		return {
			message: "",
			token: "",
		}
	},
	layout: "blank",
	methods: {
		doLogin() {
			this.message = this.token ? "": "token is empty";
			if(this.message) return
			this.$axios.post(`/api/token`, {secret:this.token})
				.then( d => {
					this.$router.push("/");
				})
				.catch(error => { 
					if ( error.response && error.response.data && error.response.data.message) {
						this.message = error.response.data.message;
					}
				});
		}
	}
}
</script>
