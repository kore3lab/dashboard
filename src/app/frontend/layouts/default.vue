<template>
	<div class="wrapper toggled">
		<c-navbar class="main-header"/>
		<c-aside class="main-sidebar"/>
		<c-splitpanes class="main-body vh-100" :class="{'terminal-closed': !isTerminal }" horizontal>
			<c-pane class="body-pane" :size="100-paneSize">
				<nuxt/>
			</c-pane>
			<c-pane v-show="isTerminal" class="terminal-pane" :size="paneSize">
				<i @click="isMinimize = !isMinimize" v-bind:class="{'fa-angle-up':isMinimize, 'fa-angle-down':!isMinimize }" class="fas cursor text-secondary"></i>
				<c-terminal @opened="onTerminalOpened" @closed="onTerminalClosed" />
			</c-pane>
		</c-splitpanes>
		<c-footer class="main-footer"/>
	</div>
</template>
<script>
import Aside				from "./components/aside.vue"
import Navbar				from "./components/navbar.vue"
import Terminal				from "./components/terminal.vue"
import Footer				from "./components/footer.vue"
import {Splitpanes,Pane}	from "splitpanes";
import "splitpanes/dist/splitpanes.css";

const TERMINAL_HEIGHT=35
export default {
	data() {
		return {
			isTerminal: false,
			isMinimize: false,
			paneSize: 0,
			beforePaneSize: 0
		}
	},
	components: {
		"c-aside": Aside,
		"c-navbar": Navbar,
		"c-terminal" : Terminal,
		"c-footer": Footer,
		"c-splitpanes": Splitpanes,
		"c-pane": Pane
	},
	computed: {
	},
	watch: {
		isMinimize(newVal) {
			if(newVal) {
				let pEl = document.querySelector(".main-body.splitpanes");			// outer pane
				let hEl = document.querySelector(".main-body .splitpanes__pane");	// body pane
				let nEl = document.querySelector(".terminal-nav-wrapper-class");	// terminal pane's nav-header

				if (pEl && hEl && nEl) {
					this.beforePaneSize = (1 - hEl.clientHeight / pEl.clientHeight)*100;
					this.paneSize = nEl.clientHeight / pEl.clientHeight * 100;
				}
			} else {
				this.paneSize = this.beforePaneSize;
			}
		}
	},
	created() {
	},
	methods: {
		onTerminalOpened(key, len) {
			if(!this.isTerminal) {
				this.paneSize = TERMINAL_HEIGHT;
				this.isTerminal = true;
				this.isMinimize = false;
			}
		},
		onTerminalClosed(key, len) {
			if(len == 0) {
				this.paneSize = 0;
				this.isTerminal = false;
				if(this.isMinimize) {
					this.beforePaneSize = 0;
					this.isMinimize = false;
				}
				
			}
		}
	},
	beforeDestroy(){
	}
}
</script>
