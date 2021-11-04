<template>
<div class="h-100">
	<b-tabs ref="tabs" v-model="current_tab" @changed="current_tab = (Object.keys(tabs).length -1)" @activate-tab="onActivateTab" class="h-100" content-class="terminal-content-class" nav-class="terminal-nav-class" nav-wrapper-class="terminal-nav-wrapper-class" active-tab-class="terminal-active-tab-class" active-nav-item-class="terminal-active-nav-item-class text-bold">
		<b-tab v-for="(d, k) in  tabs" :key="k" title-item-class="terminal-title-item-class" title-link-class="terminal-title-link-class text-sm">
			<template #title>
				<i class="fas fa-file-alt mr-1"></i>{{ d.title }}<i class="fas cursor fa-times ml-1 text-secondary" @click="onTerminalCloseClick(k)"></i>
			</template>
			<c-terminal-logs v-if="d.type=='logs'" v-model="tabs[k]" class="terminal-content"/>
			<div v-else-if="d.type=='shell'"></div>
			<div v-else-if="d.type=='cluster'"></div>
		</b-tab>
	</b-tabs>
</div>
</template>
<style>
/* terminal) b-tabs classes */
div.terminal-nav-wrapper-class { padding-right: 2.5rem; border-bottom: 1px solid #ccc; background-color:var(--bg-color)}
	ul.terminal-nav-class { height:2rem; border-bottom: 0; flex-wrap: nowrap; white-space: nowrap; overflow: auto; -ms-overflow-style: none; /* IE and Edge */ scrollbar-width: none; /* Firefox */}
	ul.terminal-nav-class::-webkit-scrollbar { display: none; /* Chrome, Safari, Opera*/ }
		li.terminal-title-item-class { background-color: transparent; }
			a.terminal-active-nav-item-class { border:0!important; border-bottom: .3rem solid #3593f8!important; font-weight: 500; }
			a.terminal-title-link-class { line-height: 1.7rem; padding: 0 .5em 0 .5em; color:#999; }
			a.terminal-title-link-class:not(.active):hover { border-color: #fff!important; }
div.terminal-content-class { height: 100%; margin-top: -2rem; padding-top:2rem }	/* content */
	div.terminal-active-tab-class { height: 100%; }
/* terminal) custom  */
div.terminal-content { height:100%;margin-top:-1.8rem; padding-top:1.8rem; }
div.terminal-content .header-wrapper { height:2rem;  }
div.terminal-content .body-wrapper { height:100%; padding: .2rem .5rem; color:#fff; background-color: #000; overflow: auto; scrollbar-color: auto;}
div.terminal-content .header-wrapper ul {display:flex; line-height:1.9rem; font-size: .875rem; margin:0; padding: 0; list-style: none; }
div.terminal-content .header-wrapper ul li { padding-left: 1rem; }
div.terminal-content .header-wrapper label { height:1.9rem; margin:0; padding-right:.5rem; font-weight: 700!important;}
div.terminal-content .header-wrapper ul li select { height: 1.6rem; width: 10rem; margin-left: .5rem; padding-top:0; padding-bottom:0;margin-top: -.1rem}
div.terminal-content .header-wrapper ul li .input-group-sm { height:1.85rem; margin:0; line-height: 1.8rem; }
div.terminal-content .header-wrapper ul li .input-group-sm input { height: 1.7rem; margin-top:.1rem; }
div.terminal-content .header-wrapper ul li .input-group-sm button { height: 1.7rem; margin-top:.1rem; padding: .2rem .3rem; }
</style>

<script>
import VueTerminalLogsView	from "~/components/terminal/terminal-logs.vue";

export default {
	components: {
		"c-terminal-logs": { extends: VueTerminalLogsView },
	},
	data() {
		return {
			tabs : {},
			current_tab: 0,
		}
	},
	created() {
		this.$nuxt.$on("open-terminal", (title, type, param) => {
			const key = `${title}-${type}`;
			if(!this.tabs[key]) {
				// logs				: { title: "", type: "logs", metadata: {namespace:"", name:""}, container:"", containers: ["",""] }
				// shell(pod)		: { title: "", type: "shell", metadata: {namespace:"", name:""}, container:"", containers: ["",""] }
				// shell(cluster)	: { title: "", type: "cluster", cluster: "apps-06" }
				if (type == "logs") this.$set(this.tabs, key, { title: title, type: "logs", metadata: param.metadata, container: param.container, containers: param.containers } );
				if (type == "shell") this.$set(this.tabs, key, { title: title, type: "shell", metadata: param.metadata, container: param.container, containers: param.containers } );
				if (type == "cluster") this.$set(this.tabs, key, { title: title, type: "cluster", cluster: param } );
				this.$emit("opened", key, Object.keys(this.tabs).length);
			} else{
        if (type === "logs") this.$set(this.tabs, key, { title: title, type: "logs", metadata: param.metadata, container: param.container, containers: param.containers } );
        this.current_tab = Object.keys(this.tabs).findIndex((v) => v === key);
      }
		});
	},
	methods: {
		onActivateTab(idx, prevIdx, ev) {
			const ul = ev.vueTarget.$el.querySelector("ul");
			const el = ul.querySelector(`li:nth-child(${idx+1})`);

			let x1 = el.offsetLeft - ul.offsetLeft - ul.scrollLeft;
			let x2 = el.offsetLeft + el.offsetWidth - ul.offsetLeft;
			let x3 = ul.scrollLeft + ul.offsetWidth;

			if (x1 < 0) {
				ul.scrollLeft += (x1-100);		// hidden left (move to right)
			}  else if( x2 > x3 ) {
				ul.scrollLeft += x2 - x3 + 100;	// hidden right (move to left)
			}
		},
		onTerminalCloseClick(k) {
			this.$delete(this.tabs, k);
			this.$emit("closed", k, Object.keys(this.tabs).length);
		},
	},
	beforeDestroy(){
		this.$nuxt.$off("open-terminal");
	}
}
</script>