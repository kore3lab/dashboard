<template>
<div>
	<div class="header-wrapper">
		<ul>
			<li><label>Namespace</label><span>{{ value.metadata.namespace }}</span></li>
			<li><label>Pod</label><span>{{ value.metadata.name }}</span></li>
			<li><label>Containers</label><b-form-select size="sm" :options="value.containers" :value="value.container" @input="onChangeContainer"></b-form-select></li>
			<li>
				<div class="input-group input-group-sm">
					<b-form-input id="txtKeyword" v-model="searchVal" @keydown="onSearchDown" class="form-control float-right" placeholder="Search"></b-form-input>
					<div class="input-group-append"><button @click="onSearchDown" class="btn btn-default"><i class="fas fa-search"></i></button></div>
				</div>
			</li>
		</ul>
	</div>
	<div v-if="logs.length > 0" class="body-wrapper" @scroll="handleScroll">
		<p v-for="d in logs" :key="d.index" class="logs" v-html="d.replace(/^\d+.*?\s/gm, '')"></p>
	</div>
	<div v-else class="body-wrapper d-flex align-items-center justify-content-center">
		There are no logs available for container
	</div>
</div>
</template>
<style scroped>
div.terminal-content .body-wrapper p.logs {margin-bottom:0; font-size: .825rem;}
mark{background-color: #3593f8;}
</style>
<script>
import AnsiConverter from "ansi-to-html"
const ansi = new AnsiConverter();

const TAIL_LINES = 300;
export default {
	props: {
		value: {
			type: Object,
			default: function () {
				return {
					metadata: { namespace: "", name: ""} ,
					container: "",
					containers: []
				}
			}
		},
	},
	data() {
		return {
			logs: [],
			searchVal: "",
			matchVal: "",
			lastTimestamp:"",
			offsetFlag: false,
			logInterval : ''
		}
	},
	mounted() {
		this.baseURL = this.getApiUrl("", "pods", this.value.metadata.namespace, this.value.metadata.name) + "/log";
		this.watchLogs();
	},
	methods: {
		intervalLogData(callback){
			const tic = () => {
				this.logInterval = setTimeout(() => {
					callback();
					tic();
				}, 10000)
			}
			tic();
		},
		setLogsData(respData){
			const wrapEl = this.$el.querySelector("div.body-wrapper");
			const isScrollMove = (wrapEl.scrollTop >= wrapEl.scrollHeight - wrapEl.clientHeight);
			const data = ansi.toHtml(respData).split("\n").filter(d => d.length > 0);
			let list = [];
			if(this.lastTimestamp){
				for(let i=data.length-1; i >= 0; i--) {
					if(data[i].indexOf(this.lastTimestamp) !== -1) break;
					list.push(data[i]);
				}
				list.reverse();
			}else{
				list = data;
			}
			this.logs = [...new Set(this.logs.concat(list))];
			this.lastTimestamp = this.logs[this.logs.length - 1]?.match(/^\d+\S+/gm);
			if (isScrollMove) this.$nextTick(() => { wrapEl.scrollTop = wrapEl.scrollHeight; this.offsetFlag = false; this.wordMark()});	// fix scroll position
		},
		loadMore(){
			this.$axios.get(`${this.baseURL}?follow=0&timestamps=1&container=${this.value.container}&tailLines=${this.logs.length + TAIL_LINES}&sinceTime=${this.lastTimestamp}`)
				.then((resp) => { this.setLogsData(resp.data); })
				.catch(e => { this.msghttp(e);});
		},
		watchLogs(){
			this.logs = [];
			this.$axios.get(`${this.baseURL}?follow=0&timestamps=1&container=${this.value.container}&tailLines=${TAIL_LINES}`)
				.then((resp) => {
					this.setLogsData(resp.data);
					this.intervalLogData(this.loadMore);
				})
				.catch(e => { this.msghttp(e);});
		},
		onSearchDown(e){
			if(e.key === 'Escape' || e.key === 'Backspace'){
				let list = [];
				if(e.key === 'Escape') this.searchVal = '';
				this.logs.forEach( item => { list.push(item.replace(/<(\/mark|mark)([^>]*)>/gi,"")) })
				this.logs = list;
				this.matchVal = "";
			}
			else if(e.key === 'Enter' || e.type === "click") {
				if(e.isComposing) return;
				this.offsetFlag = true;
				this.wordMark(true);
			}
		},
		wordMark(search){
			if(search){
				this.matchVal = this.searchVal.trim();
			}
			if(this.matchVal !== ""){
				let list = [];
				this.logs.forEach( item => { 
					list.push(item.replace(/<(\/mark|mark)([^>]*)>/gi,"")) })
					list.forEach((item, index) => {
						const startIndex = item.toUpperCase().indexOf(this.matchVal.toUpperCase());
						if (startIndex !== -1) {
							let re1 = new RegExp(this.matchVal, "gi");
							const word = item.slice(startIndex, startIndex+this.matchVal.length)
							list[index] = item.replace(re1, `<mark>${word}</mark>`);
							if(this.offsetFlag){
								this.$el.querySelectorAll("p.logs")[index].scrollIntoView();
								this.offsetFlag = false;
							}
						}
					})
					this.logs = list;
			} else {
				this.logs = this.logs.map( d=> { return d.replace(/<(\/mark|mark)([^>]*)>/gi,"") });
			}
		},
		onChangeContainer(value) {
			clearTimeout(this.logInterval);
			this.value.container = value;
			this.watchLogs();
		},
		handleScroll(evt){
			if(evt.target.scrollTop === 0){
				const h = evt.target.scrollHeight
				this.$axios.get(`${this.baseURL}?follow=0&timestamps=1&container=${this.value.container}&tailLines=${this.logs.length + TAIL_LINES}`)
					.then((resp) => {
						this.logs = ansi.toHtml(resp.data).split("\n").filter(d => d.length > 0);
						this.lastTimestamp = this.logs[this.logs.length - 1]?.match(/^\d+\S+/gm);
						this.$nextTick(()=>{ evt.target.scrollTop = evt.target.scrollHeight-h; this.offsetFlag = false; this.wordMark()});	// fix scroll position
					})
					.catch(e => { this.msghttp(e);});
			}
		}
	},
	beforeDestroy(){
		clearTimeout(this.logInterval);
	},
}
</script>