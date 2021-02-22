<template>
	<!-- Navbar -->
	<nav class="main-header navbar navbar-expand navbar-white navbar-light">
		<ul class="navbar-nav">
			<!-- 햄버거 메뉴 -->
			<li class="nav-item">
				<a class="nav-link" data-widget="pushmenu" href="#" role="button"><i class="fas fa-bars"></i></a>
			</li>
			<!-- Selected context -->
			<li class="nav-item pt-1">
				<span class="text-lg text-muted">{{this.currentContext()}}</span>
			</li>
		</ul>
		<!-- <ul class="navbar-nav ml-auto">
			<li class="nav-item">
				<a class="nav-link" data-widget="control-sidebar" data-slide="true" href="#" role="button">
					<i class="fas fa-th-large"></i>
				</a>
			</li>
		</ul> -->
	</nav>
</template>
