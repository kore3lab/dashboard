<template>
	<div id="terminal">
	</div>
</template>

<script>
import { Xterm } from "../static/terminal/js/xterm";
import { WebTTY, protocols } from "../static/terminal/js/webtty";
import { ConnectionFactory } from "../static/terminal/js/websocket";
import {debounce} from 'lodash';
const axios = require('axios');

export default {
	layout: 'termlayout',

	data() {
		return {
			elem:"",
			instterm:"",
			requrl:"",
			token: "",
			debouncedFit_: null,
			axiosInst: "",
			reqHost: ""
		};
	},
	mounted() {
		if(!this.$route.query.termtype){
			this.toast("request parameter(termtype) error");
		}
		//개발환경에서 terminalPort 설정위한 별도 axios설정
		if (`${this.$config.nodeEnv}` === "development") {
			this.reqHost = `${location.hostname}:${this.$config.terminalPort}`;
			this.axiosInst = axios.create({
				baseURL: `${location.protocol}//` + this.reqHost
			});
			
		} else {
			this.reqHost = `${location.host}`;
			this.axiosInst = this.$axios;
		}


		this.elem = document.getElementById("terminal");
		this.instterm = new Xterm(this.elem);

		this.debouncedFit_ = debounce(() => {
			this.handleResize();
		}, 100);
		this.debouncedFit_();
		window.addEventListener('resize', () => this.debouncedFit_());

		//window.addEventListener("resize", this.handleResize);

		this.getToken();
		this.handleResize()
	},
	beforeDestroy(){
		window.removeEventListener("resize", this.debouncedFit_());
	},
	methods: {
		getToken() {
		   	switch(this.$route.query.termtype){
				case 'cluster':
					this.requrl = '/api/terminal/clusters/' +  this.$route.query.cluster + '/termtype/' + this.$route.query.termtype;
					break;
				case 'pod':
					this.requrl = '/api/terminal/clusters/' +  this.$route.query.cluster + '/namespaces/' + this.$route.query.namespace + '/pods/' + this.$route.query.pod + '/termtype/' + this.$route.query.termtype;
					break;
				case 'container':
					this.requrl = '/api/terminal/clusters/' +  this.$route.query.cluster + '/namespaces/' + this.$route.query.namespace + '/pods/' + this.$route.query.pod + '/containers/' + this.$route.query.container + '/termtype/' + this.$route.query.termtype;
					break;
				default:
					return;
			}

			this.axiosInst.get(this.requrl)
					.then( resp => {
						this.token = resp.data.Token;
						this.connWs();
						this.instterm.focus();
					})
					.catch(e => { this.msghttp(e);});
		},
		connWs() {
			//웹소캣 접속
			const httpsEnabled = `${location.protocol}` === "https:";
			const reqProtocol = httpsEnabled ? 'wss://' : 'ws://';
			const reqTailURL = '/api/terminal/ws';

			const url = reqProtocol + this.reqHost + reqTailURL;

            let authToken = this.token;
			let factory = new ConnectionFactory(url, protocols);
			let wt = new WebTTY(this.instterm, factory, "", authToken);
			let closer = wt.open();

			window.addEventListener("unload", () => {
				closer();
				term.close();
			});
		},
		handleResize(event) {
			this.elem.style.width = window.innerWidth + 'px';
			this.elem.style.height = window.innerHeight + 'px';
			this.instterm.resizeListener();
		},
		checkUndefined(source) {
			return (typeof source !== 'undefined')? source : '';
		},

	}

};
</script>
<style>
@import 'xterm/css/xterm.css';
@import '../static/terminal/css/xterm_customize.css';
</style>