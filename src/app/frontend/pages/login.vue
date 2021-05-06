<template>
	<div class="content-wrapper min-vh-100 overflow-hidden">
		<!-- token -->
		<div class="row" v-show="(auth.scheme=='token')" style="padding-top:10em">
			<div class="col-md-3"></div>
			<div class="col-md-6">
				<div class="login-logo"><b>Kore</b>board</div>
				<div class="card card-info">
					<div class="card-header"><h3 class="card-title">Sign in</h3></div>
					<div class="card-body">
						<p v-show="(auth.provider!='static-token')" class="card-text">Service account has a secret with valid token that can be used to login to dashboard. </p>
						<p v-show="(auth.provider=='static-token')" class="card-text">Enter your static token that can be used to login to dashboard. </p>
						<pre v-show="(auth.provider!='static-token')" class="text-sm border bg-white"><small id="copyText"><i class="fas fa-dollar-sign mr-1"></i>SECRET="$(kubectl get sa -n kore -l app=kore-board -o jsonpath='{.items[0].secrets[0].name}')"
<i class="fas fa-dollar-sign mr-1"></i>echo "$(kubectl get secret ${SECRET} -n kore -o jsonpath='{.data.token}' | base64 --decode)"</small><button type="button" class="btn p-0 pl-2" @click="copy()"><i class="fas fa-copy"></i></button></pre>
						<div class="input-group">
							<b-form-input v-model="form.token" autofocus size="sm" type="password" :state="isLoginState" placeholder="Enter token" trim></b-form-input>
    						<b-form-invalid-feedback id="token-feedback">{{message}}</b-form-invalid-feedback>
						</div>
					</div>
					<div class="card-footer">
						<b-button size="sm" variant="info" class="float-right" @click="doLogin">Sign in</b-button>
					</div>
				</div>
			</div>
			<div class="col-md-3"></div>
		</div>
		<!-- user -->
		<div class="row" v-show="(auth.scheme=='user')" style="padding-top:10em">
			<div class="col-md-4"></div>
			<div class="col-md-4">
				<div class="login-logo"><b>Kore</b>board</div>
				<div class="card card-info">
					<div class="card-header"><h3 class="card-title">Sign in</h3></div>
					<div class="card-body">
						<b-form-group label-cols="4" label-cols-lg="3" label-size="sm" label="Username" label-for="input-sm">
							<b-form-input v-model="form.username" type="text" size="sm" :state="isLoginState" placeholder="Username"></b-form-input>
						</b-form-group>
						<b-form-group label-cols="4" label-cols-lg="3" label-size="sm" label="Password" label-for="input-sm">
							<b-form-input v-model="form.password" type="password" size="sm" :state="isLoginState" placeholder="Password"></b-form-input>
						<b-form-invalid-feedback id="token-feedback">{{message}}</b-form-invalid-feedback>
						</b-form-group>
					</div>
					<div class="card-footer">
						<b-button size="sm" variant="info" class="float-right" @click="doLogin">Sign in</b-button>
					</div>
				</div>
			</div>
			<div class="col-md-4"></div>
		</div>
		<!-- none -->
		<div class="row" v-show="(auth.scheme=='')" style="padding-top:10em">
			<div class="col-md-4"></div>
			<div class="col-md-4">
				<div class="card card-info">
					<div class="card-header"><h3 class="card-title">Sign in</h3></div>
					<div class="card-body">
						<div class="text-center text-success">
							<b-spinner type="grow" variant="success" class="align-middle mr-2"></b-spinner>
							<span class="text-lg align-middle">Loading...</span>
						</div>
					</div>
				</div>
			</div>
			<div class="col-md-4"></div>
		</div>
	</div>
</template>
<script>
export default {
	components: {
	},
	computed: {
		isLoginState() {
			return this.message ? false: null;
		}
    },
	data() {
		return {
			message: "",
			auth: {
				strategy: "",
				scheme: "",
				provider: "",
			},
			form: {
				username: "",
				password: "",
				token: "",
			},
		}
	},
	layout: "blank",
	beforeCreate() {
		this.$axios.get(`/api/auth/login`, {})
			.then( d => {
				this.auth = d.data;
				if (this.auth.scheme == "") {
					setTimeout(this.doLogin, 1000)	//auto-singin
				}
			})
			.catch(error => { 
				if ( error.response && error.response.data && error.response.data.message) {
					this.message = error.response.data.message;
				}
			});
	},
	methods: {
		copy() {
			const copyText = document.getElementById("copyText").textContent;
			const textArea = document.createElement('textarea');
			textArea.textContent = copyText;
			document.body.append(textArea);
			textArea.select();
			document.execCommand("copy")
			textArea.remove()
		},
		doLogin() {
			this.$auth.loginWith(this.auth.strategy, { data: this.form })
				.catch(error => { 
					if ( error.response && error.response.data && error.response.data.message) {
						this.message = error.response.data.message;
					}
				});
		},
		copy() {
			const copyText = document.getElementById("copyText").textContent;
			const textArea = document.createElement('textarea');
			textArea.textContent = copyText;
			document.body.append(textArea);
			textArea.select();
			document.execCommand("copy")
			textArea.remove()
		}

	}
}
</script>
