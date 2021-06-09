<template>
	<div id="aside-contexts" class="sidebar-contexts d-flex flex-column sidebar-dark-primary border-right border-secondary">
		<div v-for="(ctx, index) in contexts()" :key="ctx" :value="ctx">
			<b-overlay :show="showOverlay===ctx" rounded="sm">
				<b-overlay :show="showErrorOverlay===ctx" rounded="sm" variant="white">
				<template #overlay><b-icon icon="exclamation-circle-fill" font-scale="2" variant="danger"></b-icon></template>
				<b-button v-bind:id="'btn_aside_cluster_' + ctx" @click="onContextSelected(ctx)" v-bind:class="{active: ctx===currentContext()}" :value="ctx" class="w-100 text-uppercase">{{ ctx.substring(0,1) }}</b-button>
				</b-overlay>
			</b-overlay>
			<p class="text-center text-white text-truncate">{{ctx}}</p>
			<b-popover v-bind:target="'btn_aside_cluster_' + ctx" :title="ctx" triggers="hover" boundary="window" boundary-padding="0">
				<ul class="list-unstyled m-0 p-0">
					<li v-if="ctx!==currentContext()"><b-link href="#" @click="onContextDelete(ctx, index)"><b-icon icon="x-circle" class="mr-1 text-danger"></b-icon>Remove</b-link></li>
					<li v-if="index!==0"><b-link href="#" @click="onMoveTop(index)"><b-icon icon="chevron-double-up" class="mr-1 "></b-icon>Top</b-link></li>
					<li v-if="index!==0"><b-link href="#" @click="onMoveUp(index)"><b-icon icon="caret-up-fill" class="mr-1 "></b-icon>Up</b-link></li>
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
export default {
	data() {
		return {
			showOverlay: "",
			showErrorOverlay: "",
			list: [],
		}
	},
	async fetch() {
		if(!this.currentContext()) {

			// context 선택
			let ctx = this.$route.query.context ? this.$route.query.context: "";
			if (!ctx && localStorage.getItem("currentContext")!=null) ctx = localStorage.getItem("currentContext");
			if(ctx) {
				if( !this.contexts().find(el => el===ctx)) ctx = "";
			}

			// context list
			let equals = (a, b) => {
				if (!a) return false;
				if (!b) return false;
				if (a.length !== b.length) return false;
				for (let i=0; i < a.length; ++i) {
					let exist = b.find(v => v===a[i])
					if(!exist) return false;
				}
				return true;
			};

			this.$axios.get(`/api/contexts`)
				.then((resp)=>{
					if(resp.data.contexts) {
						let local;
						try { local = JSON.parse(localStorage.getItem("contexts")); } catch (e) {}
						if (equals(local,resp.data.contexts)) {
							this.contexts(local);
						} else {
							this.contexts(resp.data.contexts);
							localStorage.setItem("contexts",resp.data.contexts);
						}
					}
					if(!ctx) ctx = resp.data.currentContext;
					this.onContextSelected(ctx)
				}).catch(error=> {
					this.toast(error.message, "danger");
				})
		}
	},
	mounted() {
		this.$nuxt.$on("navbar-set-context-selected", (ctx) => this.onContextSelected(ctx) );
	},
	methods: {
		// context select
		onContextLoad(ctx) {

			// context list
			let equals = (a, b) => {
				if (!a) return false;
				if (!b) return false;
				if (a.length !== b.length) return false;
				for (let i=0; i < a.length; ++i) {
					let exist = b.find(v => v===a[i])
					if(!exist) return false;
				}
				return true;
			};

			this.$axios.get(`/api/contexts`)
				.then((resp)=>{
					if(resp.data.contexts) {
						let local;
						try { local = JSON.parse(localStorage.getItem("contexts")); } catch (e) {}
						if (equals(local,resp.data.contexts)) {
							this.contexts(local);
						} else {
							this.contexts(resp.data.contexts);
							localStorage.setItem("contexts",resp.data.contexts);
						}
					}
					if(!ctx) ctx = resp.data.currentContext;
					this.onContextSelected(ctx)
				}).catch(error=> {
					this.toast(error.message, "danger");
				})
		},
		// context select
		onContextSelected(ctx) {

			if (!ctx) return

			this.showOverlay = ctx;
			this.showErrorOverlay = "";
			this.$axios.get(`/api/contexts/${ctx}`)
				.then((resp)=>{
					if (resp.data["error"]) {
						this.showErrorOverlay = ctx
						this.toast(resp.data["error"], "warning")
					} else {
						this.currentContext(ctx ? ctx : resp.data.currentContext.name);
						let nsList = [{ value: "", text: "All Namespaces" }];
						if (resp.data.currentContext.namespaces) {
							resp.data.currentContext.namespaces.forEach(el => {
								nsList.push({ value: el, text: el });
							});
						}
						this.namespaces(nsList);
						this.resources(resp.data.currentContext.resources);
						this.statusbar({message: "", kubernetesVersion: resp.data.currentContext.kubernetesVersion, platform: resp.data.currentContext.platform})
						localStorage.setItem("currentContext", this.currentContext());
					}
			}).catch(error=> {
				this.toast(error.message, "danger");
			}).finally(() => { this.getCRD() });

		},
		getCRD() {
			let crList = []
			this.$axios.get(this.getApiUrl("apiextensions.k8s.io","customresourcedefinitions"))
			.then((resp) => {
				resp.data.items.forEach(el => {
					if(crList.find(e => e === el.spec.group)) {

					}else {
						crList.push(el.spec.group)
					}
				})
			}).catch(error=> {
				this.toast(error.message, "danger");
			}).finally(() => {
				this.drawList(crList)
				this.showOverlay = "";
				this.$nuxt.$emit("navbar-context-selected");
			})
		},
		drawList(crList) {
			this.list = []
			crList.forEach(el => {
				this.list.push({[el]:this.resources()[el]})
			})
			this.$nuxt.$emit("crList_up",this.list)
		},
		onContextDelete(ctx, index) {
			this.confirm(`Delete a selected cluster "${ctx}" , Are you sure?`, yes => {
				if(!yes) return;
				this.$axios.delete(`/api/contexts/${ctx}`)
					.then( resp => {
						let contexts = []
						this.contexts().forEach( d => {
							if (d!==ctx) contexts.push(d)
						});
						this.contexts(contexts);
						this.toast("Delete a selected cluster...OK", "success");
					}).catch(e => { this.msghttp(e);});

			})
		},
		onMoveTop(index) {
			if (index === 0) return;
			var list = JSON.parse(JSON.stringify(this.contexts()))	// deep-copy
			list.splice(0, 0, list.splice(index,1)[0]);
			this.contexts(list);
			localStorage.setItem("contexts", JSON.stringify(list));
		},
		onMoveUp(index) {
			if (index === 0) return;
			var list = JSON.parse(JSON.stringify(this.contexts()))	// deep copy
			list.splice(index-1, 0, list.splice(index, 1)[0]);
			this.contexts(list);
			localStorage.setItem("contexts", JSON.stringify(list));
		},
		onMoveDown(index) {
			if (index === this.contexts().length - 1) return;
			var list = JSON.parse(JSON.stringify(this.contexts()))	// deep copy
			list.splice(index + 1, 0, list.splice(index, 1)[0]);
			this.contexts(list);
			localStorage.setItem("contexts", JSON.stringify(list));
		}
	}
}
</script>