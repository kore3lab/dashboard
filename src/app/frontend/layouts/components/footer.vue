<template>
<footer>
		<div class="float-left pr-2 text-sm">{{ statusbar().kubernetesVersion }}</div>
		<div class="float-left pl-2 pr-2 mr-1 text-sm">{{ statusbar().platform }}</div>
		<div class="float-left  text-sm">
	</div>
	<div class= "margin-right">
		<a href="https://github.com/kore3lab/dashboard" target="_blank"><b-icon class="float-right" icon="github" font-scale="1.25"></b-icon></a>
		<div class="float-right text-sm mr-2">
		<b>Version</b> {{$config.version}}
		</div>
	</div>
</footer>
</template>
