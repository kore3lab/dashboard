<template>
	<footer class="main-footer" style="">
		<div class="float-left pr-2 text-sm">{{ statusbar().kubernetesVersion }}</div>
		<div class="float-left pl-2 pr-2 mr-1 text-sm">{{ statusbar().platform }}</div>
		<div class="float-left  text-sm"></div>
		<a href="https://github.com/acornsoftlab/dashboard" target="_blank"><b-icon class="float-right" icon="github" font-scale="1.25"></b-icon></a>
		<div class="float-right text-sm mr-2">
			<b>Version</b> {{$config.version}}
		</div>
	</footer>
</template>
