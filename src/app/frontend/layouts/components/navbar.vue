<template>
	<!-- Navbar -->
	<nav class="main-header navbar navbar-expand navbar-white navbar-light">
		<ul class="navbar-nav">
			<li class="nav-item"><a class="nav-link" data-widget="pushmenu" href="#" role="button"><i class="fas fa-bars"></i></a></li>
			<li class="nav-item pt-1"><span class="text-lg text-muted">{{this.currentContext()}}</span></li>
		</ul>
		<ul class="navbar-nav ml-auto">
			<li class="nav-item">
				<nuxt-link :to="{path: '/terminal', query: {termtype: 'cluster', cluster: option}}" target="_blank" class="mr-2"><b-icon icon="terminal-fill" font-scale="1.7"></b-icon></nuxt-link>
			</li>
		</ul>
	</nav>
</template>
