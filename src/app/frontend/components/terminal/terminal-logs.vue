<template>
<div>
	<div class="header-wrapper">
		<ul>
			<li><label>Namespace</label><span>{{ value.metadata.namespace }}</span></li>
			<li><label>Pod</label><span>{{ value.metadata.name }}</span></li>
			<li><label>Containers</label><b-form-select size="sm" :options="value.containers" :value="value.container" @input="onChangeContainer"></b-form-select></li>
			<li>
				<div class="input-group input-group-sm">
					<b-form-input id="txtKeyword" class="form-control float-right" placeholder="Search"></b-form-input>
					<div class="input-group-append"><button type="submit" class="btn btn-default"><i class="fas fa-search"></i></button></div>
				</div>
			</li>
		</ul>
	</div>
	<div class="body-wrapper" @scroll="handleScroll">
		<p v-for="d in logs" :key="d" class="logs" v-html="d"></p>
	</div>
</div>
</template>
<style>
div.terminal-content .body-wrapper p.logs {margin-bottom:0; font-size: .825rem;}
</style>
<script>
import AnsiConverter from "ansi-to-html"
const ansi = new AnsiConverter();

const TAIL_LINES = 300;
let controller = {};
export default {
	props: {
		value: {
			type: Object,
			default: function () {
				return { 
					metadata: { namespace: "", name: ""} ,
					container: "",
					containers: [],
				}
			}
		},
	},
	data() {
		return {
			logs: [],
			baseURL: ""
		}
	},
	computed: {
		key: function() {
			return `${this.value.metadata.namespace}-${this.value.metadata.name}`;
		}
	},
	mounted() {
		this.hostURL = (this.$config.nodeEnv === "development") ? `${location.protocol}//${location.hostname}:${ this.$config.backendPort}`: "";
		this.baseURL = this.getApiUrl("", "pods", this.value.metadata.namespace, this.value.metadata.name) + "/log";
		this.watchLogs();
	},
	methods: {
		async watchLogs(){

			let c = new AbortController();
			controller[this.key] = c;
			const wrapEl = this.$el.querySelector("div.body-wrapper");

			try {
				this.logs = [];
				const response = await fetch(`${this.hostURL}${this.baseURL}?follow=1&timestamps=0&container=${this.value.container}&tailLines=${TAIL_LINES}`, { signal: c.signal });
				const reader = response.body.getReader();
				while (true) {
					const { value, done } = await reader.read();
					if (done) break;

					const isScrollMove = (wrapEl.scrollTop >= wrapEl.scrollHeight-wrapEl.clientHeight);		// whitch scrolling to bottom (if scroll-position is bottom)

					const buffer =  ansi.toHtml(new TextDecoder().decode(value, {stream: true}));
					const data = buffer.split("\n").filter(d => d.length > 0);
					this.logs = this.logs.length == 0 ? data: this.logs.concat(data);

					if(isScrollMove) this.$nextTick(()=>{ wrapEl.scrollTop = wrapEl.scrollHeight; }); // scroll to bottom
				}
			} catch (e) {
				console.log(`fetch connection closed (namespace=${this.value.metadata.namespace}, pods=${this.value.metadata.name}, container=${this.value.container}, ${e.message})`);
			}
			
		},
		onChangeContainer(value) {
			controller[this.key].abort();
			this.value.container = value;
			this.watchLogs();
		},
		handleScroll(evt){
			if(evt.target.scrollTop === 0){
				const h = evt.target.scrollHeight
				this.$axios.get(`${this.baseURL}?follow=0&container=${this.value.container}&tailLines=${this.logs.length + TAIL_LINES}`)
					.then((resp) => {
						this.logs = ansi.toHtml(resp.data).split("\n");
						this.$nextTick(()=>{ evt.target.scrollTop = evt.target.scrollHeight-h;  });	// fix scroll position
					})
					.catch(e => { this.msghttp(e);});
			}
		}
	},
	beforeDestroy(){
		controller[this.key].abort();
	},
}
</script>