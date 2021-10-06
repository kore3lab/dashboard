<template>
	<!-- Navbar -->
	<nav class="main-header navbar navbar-expand navbar-white navbar-light">
		<ul class="navbar-nav">
			<li class="nav-item"><a class="nav-link" data-widget="pushmenu" href="#" role="button"><i class="fas fa-bars"></i></a></li>
			<li class="nav-item pt-1"><span class="text-lg text-muted">{{this.currentContext()}}</span></li>
		</ul>
		<ul class="navbar-nav ml-auto margin-right">
			<li v-if="currentContext()" class="nav-item"><nuxt-link :to="{path: '/terminal', query: {termtype: 'cluster', cluster: this.currentContext()}}" target="_blank" class="nav-link"><b-icon icon="terminal-fill" font-scale="1.2"></b-icon></nuxt-link></li>
			<li class="nav-item"><a class="nav-link" data-widget="fullscreen" href="#" role="button"><i class="fas fa-expand-arrows-alt"></i></a></li>
			<!--<li class="nav-item"><a class="nav-link" data-widget="control-sidebar" data-slide="true" href="#" role="button"><i class="fas fa-th-large"></i></a></li>-->
			<li class="nav-item"><a class="nav-link" href="#" @click="logout"><b-icon icon="door-closed" font-scale="1.2"></b-icon></a></li>
		</ul>
	</nav>
</template>
<script>

export default {
	methods: {
		logout() {
			this.$auth.logout().catch(e => { this.msghttp(e);})
		},
	}
}
</script>

