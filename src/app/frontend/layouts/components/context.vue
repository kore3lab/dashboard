<template>
<div id="aside-contexts" class="sidebar-contexts d-flex flex-column sidebar-dark-primary border-right border-secondary">
	<div v-for="(option, index) in contexts()" :key="option" :value="option">
		<b-overlay :show="showOverlay==option" rounded="sm">
		<b-button v-bind:id="'btn_aside_cluster_' + option" @click="onContextSelected(option, $event)" v-bind:class="{active: option==currentContext()}" :value="option" class="w-100 text-uppercase">{{ option.substring(0,1) }}</b-button>
		</b-overlay>
		<p class="text-center text-white text-truncate">{{option}}</p>
		<b-popover v-bind:target="'btn_aside_cluster_' + option" v-bind:title="option" triggers="hover" boundary="window" boundary-padding="0">
			<ul class="list-unstyled m-0">
				<li v-if="option!=currentContext()"><b-link href="#" @click="onContextDelete(option, index)"><b-icon icon="x-circle" class="mr-1 text-danger"></b-icon>Remove</b-link></li>
				<li v-if="index!=0"><b-link href="#" @click="onMoveTop(index)"><b-icon icon="chevron-double-up" class="mr-1 "></b-icon>Top</b-link></li>
				<li v-if="index!=0"><b-link href="#" @click="onMoveUp(index)"><b-icon icon="caret-up-fill" class="mr-1 "></b-icon>Up</b-link></li>
				<li v-if="index<(contexts().length-1)"><b-link href="#" @click="onMoveDown(index)"><b-icon icon="caret-down-fill" class="mr-1"></b-icon>Down</b-link></li>
			</ul>
		</b-popover>
	</div>
	<div>
		<b-button id="btn_aside_add_cluster" variant="primary" to="/kubeconfig" class="w-100">+</b-button>
		<b-tooltip target="btn_aside_add_cluster" placement="right" boundary="window">Add a cluster</b-tooltip>
	</div>
</div>
</template>
<script>
import axios from "axios"
export default {
	data() {
		return {
			showOverlay: ""
		}
	},
	async fetch() {
		// context 리스트 조회
		if(!this.currentContext()) {
			let equals = (a, b) => {
				if (!a) return false;
				if (!b) return false;
				if (a.length !== b.length) return false;
				for (var i=0; i < a.length; ++i) {
					let exist = b.find(v=> v==a[i])
					if(!exist) return false;
				}
				return true;
			};

			this.showOverlay = {"apps-05":false, "apps-06":false}

			let resp = await axios.get(`${this.backendUrl()}/api/clusters`);
			if(resp.data.contexts) {
				let local;
				try { local = JSON.parse(localStorage.getItem("contexts")); } catch (e) {};
				if (equals(local,resp.data.contexts)) {
					this.contexts(local);
				} else {
					this.contexts(resp.data.contexts);
					localStorage.setItem("contexts",resp.data.contexts);
				}
			}

			let ctx = this.$route.query.context;
			if (!ctx) ctx = localStorage.getItem("currentContext");
			if(ctx) {
				if( !this.contexts().find(el => el==ctx)) ctx = "";
			}
			this.currentContext(ctx ? ctx : resp.data.currentContext);
				
			await this.loadNamespaces(ctx);
			this.$nuxt.$emit("navbar-context-selected");
		}
	},
	mounted() {
	},
	methods: {
		/**
		 * 클러스터의 네임스페이스 목록을 가져 온다.
		 * 
		 * @async
		 * @function loadNamespaces
		 * @returns {Array} 네이스페이스 값을 [{valu, text}] 값으로 반환 한다.
		 */
		async loadNamespaces(ctx){
			let nsList = [{ value: "", text: "All Namespaces" }];
			if (ctx) {
				let resp = await axios.get(`${this.backendUrl()}/raw/clusters/${ctx}/api/v1/namespaces`)
				resp.data.items.forEach(el => {
					nsList.push({ value: el.metadata.name, text: el.metadata.name });
				});
			}
			this.namespaces(nsList);
		},
		onContextSelected(ctx, ev) {
			if (ctx == this.currentContext()) return;

			this.showOverlay = ctx;
			this.loadNamespaces(ctx).then(() => {
				this.currentContext(ctx);
				localStorage.setItem("currentContext", ctx);
				this.$nuxt.$emit("navbar-context-selected");
				document.querySelectorAll("#aside-contexts button").forEach(e => {
					if(e.value) {
						if(ctx == e.value) e.classList.add("active");
						else e.classList.remove("active");
					}
				});
			}).catch(error=> {
				this.toast(error.message, "danger");
			}).finally(() => {
				this.showOverlay = ""
			});

		},
		onContextDelete(ctx, index) {
			this.confirm(`Delete a selected cluster "${ctx}" , Are you sure?`, yes => {
				if(!yes) return;
				axios.delete(`${this.backendUrl()}/api/clusters/${ctx}`)
					.then( resp => {
						this.contexts().splice(index, 1);
						this.toast("Delete a selected cluster...OK", "success");
					}).catch(e => { this.msghttp(e);});

			})
		},
		onMoveTop(index) {
			if (index === 0) return;
			var list = JSON.parse(JSON.stringify(this.contexts()))	// deep clone
			list.splice(0, 0, list.splice(index)[0]);
			this.contexts(list);
			localStorage.setItem("contexts", JSON.stringify(list));
		},
		onMoveUp(index) {
			if (index === 0) return;
			var list = JSON.parse(JSON.stringify(this.contexts()))	// deep clone
			list.splice(index-1, 0, list.splice(index, 1)[0]);
			this.contexts(list);
			localStorage.setItem("contexts", JSON.stringify(list));
		},
		onMoveDown(index) {
			if (index === this.contexts().length - 1) return;
			var list = JSON.parse(JSON.stringify(this.contexts()))	// deep clone
			list.splice(index + 1, 0, list.splice(index, 1)[0]);
			this.contexts(list);
			localStorage.setItem("contexts", JSON.stringify(list));
		}
	}
}
</script>