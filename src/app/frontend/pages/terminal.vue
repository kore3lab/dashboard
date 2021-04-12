<template>
		<div id="terminal">
    </div>
</template>

<script>
import axios		from "axios";
import { Xterm } from "../static/terminal/js/xterm";
import { WebTTY, protocols } from "../static/terminal/js/webtty";
import { ConnectionFactory } from "../static/terminal/js/websocket";

export default {
	layout: 'termlayout',

    data() {
    return {
      elem:"",
      instterm:"",
      requrl:"",
      token: ""
    };
  },
  mounted() {
      if(!this.$route.query.termtype){
          this.toast("request parameter(termtype) error");
      }

      this.elem = document.getElementById("terminal");
      this.instterm = new Xterm(this.elem);

       window.addEventListener("resize", this.handleResize); 
      
      this.getToken();
   },
   beforeDestroy(){
      window.removeEventListener("resize", this.handleResize);
   },
  methods: {
        getToken() {
          switch(this.$route.query.termtype){
            case 'cluster':
              this.requrl = '/api/clusters/' +  this.$route.query.cluster + '/terminal';
            break;
            case 'pod':
              this.requrl = '/api/clusters/' +  this.$route.query.cluster + '/namespaces/' + this.$route.query.namespace + '/pods/' + this.$route.query.pod + '/terminal';
            break;
            case 'container':
              this.requrl = '/api/clusters/' +  this.$route.query.cluster + '/namespaces/' + this.$route.query.namespace + '/pods/' + this.$route.query.pod + '/containers/' + this.$route.query.container + '/terminal';
            break;
            default:
              return;
          }

          axios.get(this.backendUrl() + this.requrl)
          .then( resp => {
						this.token = resp.data.Token;
            this.connWs();
            this.instterm.focus();
					})
					.catch(e => { this.msghttp(e);});
      },
     connWs() {       
                //웹소캣 접속
                const httpsEnabled = `${location.protocol}` == "https:";
                const url = (httpsEnabled ? 'wss://' : 'ws://') + `${location.hostname}:${this.$config.backendPort}` + '/api/terminal/' + 'ws';
          
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
        }
      
  }      

};
</script>
<style>
 @import 'xterm/dist/xterm.css';
 @import '../static/terminal/css/xterm_customize.css';
 @import '../static/terminal/css/terminal.css';
</style>