<template>
	<section class="content-view border-primary border-left min-vh-100" v-bind:class="{ 'min-vh-100' : errorcheck}">
		<div class="card card-primary m-0 layer" v-bind:class="{ 'min-vh-100' : errorcheck}">
			<!-- header -->
			<div class="card-header pt-2 pb-2 position-sticky fixed-top">
				<h3 class="card-title">{{ title }} / {{ (name && name.length >60) ? name.substring(0,60)+'...' : name }}</h3>
				<div class="card-tools">
					<span v-show="!errorcheck">
 						<button v-if="navigateStack.length>0" type="button" class="btn btn-tool" @click="onNavigateBack"><i class="fa fa-arrow-left"></i></button>
 						<button type="button" class="btn btn-tool" @click="onSync()" v-b-tooltip.hover title="Reload"><i class="fas fa-sync-alt"></i></button>
						<button type="button" class="btn btn-tool" v-show="isJSON && component"  @click="isJSON=false"><i><del>JSON</del></i></button>
						<button type="button" class="btn btn-tool" v-show="!isJSON && component" @click="isJSON=true;isYaml=false"><i>JSON</i></button>
						<button type="button" class="btn btn-tool" @click="isYaml=true" v-b-tooltip.hover title="Edit"><i class="fas fa-edit"></i></button>
						<button type="button" class="btn btn-tool" @click="deleteOverlay.visible = true" v-b-tooltip.hover title="Delete"><i class="fas fa-trash"></i></button>
					</span>
					<button type="button" class="btn btn-tool" @click="$emit('close')"><i class="fas fa-times"></i></button>
				</div>
			</div>
			<!-- error message-->
			<div class="col-md-12 mt-5 lh-vh-50" v-show="errorcheck"><p class="align-middle text-sm-center">Resource loading has failed: <b>{{ errorMessage }}</b></p></div>

			<!-- body -->
			<b-overlay :show="deleteOverlay.visible" no-center>
				<!--1. not Yaml(editor) -->
				<div v-show="!isYaml && !errorcheck" class="card-body p-2">
					<!-- 1.1 meta data -->
					<c-metadata v-show="isJSON" v-model="raw" dtCols="2" ddCols="10"></c-metadata>
					<!-- 1.2  view -->
					<component :is="component" v-if="component" v-show="!isJSON" v-model="raw" @navigate="onNavigate"/>
					<!-- 1.3 json-tree -->
					<div class="row" v-show="isJSON && title !== 'StorageClass'">
						<div class="col-md-12">
							<div class="card card-secondary card-outline m-0">
								<div class="card-header p-2">
									<h3 class="card-title" v-if="title==='ConfigMap' || title==='Secret'">Data</h3>
									<h3 class="card-title" v-else-if="title==='ServiceAccount'">Secrets</h3>
									<h3 class="card-title" v-else-if="title==='Role' || title==='ClusterRole'">Rules</h3>
									<h3 class="card-title" v-else-if="title==='RoleBinding' || title==='ClusterRoleBinding'">Subjects,RoleRef</h3>
									<h3 class="card-title" v-else-if="title==='Endpoints'">Subsets</h3>
									<h3 class="card-title" v-else-if="title==='HorizontalPodAutoscaler'">Spec,Status</h3>
									<h3 class="card-title" v-else-if="title==='Node'">NodeInfo</h3>
									<h3 class="card-title" v-else>Specification</h3>
								</div>
								<c-jsontree id="txtSpec" v-model="raw.spec" class="card-body p-2"></c-jsontree>
							</div>
						</div>
					</div>
				</div>
				<!-- 2. if Yaml(editor) then  -->
				<div v-show="isYaml && !errorcheck" class="card-body p-1 overflow-hidden">
					<div class="row">
						<div class="col-sm-12 text-right mb-1">
							<b-button variant="primary" size="sm" @click="onApply">Apply</b-button>
							<b-button variant="secondary" size="sm" @click="onReset">Reset</b-button>
							<b-button variant="secondary" size="sm" @click="isYaml=false">Close</b-button>
						</div>
					</div>
					<div class="row">
						<div class="col-md-12">
							<c-aceeditor id="txtYaml" v-model="raw" v-on:error="onError" style="min-height: calc(100vh - 85px)"></c-aceeditor>
						</div>
					</div>
				</div>
				<!-- 3. delete overlay -->
				<template #overlay>
					<div v-if="deleteOverlay.processing" class="text-center floating">
						<b-spinner small class="mr-2" label="please wait"></b-spinner><span>Watching DELETE status...</span>
					</div>
					<div v-else class="text-center floating">
						<p>Are you sure DELETE it?</p>
						<div class="text-center">
							<b-button variant="outline-danger" size="sm" class="mr-1" @click="deleteOverlay.visible = false">Cancel</b-button>
							<b-button variant="success" size="sm" @click="onDelete">OK</b-button>
						</div>
					</div>
				</template>
			</b-overlay>
		</div>
	</section>
</template>
<script>
import VueMetadataView	from "@/components/view/metadataView.vue";
import VueAceEditor 		from "@/components/aceeditor"
import VueJsonTree 			from "@/components/jsontree"

export default {
	props:["value"],
	components: {
		"c-metadata": { extends: VueMetadataView },
		"c-aceeditor": { extends: VueAceEditor },
		"c-jsontree": { extends: VueJsonTree },
	},
	data() {
		return {
			component: null,
			title: "",
			name: "",
			src: "",
			url: "",
			origin: { metadata: {}, spec: {} },
			raw: { metadata: {}, spec: {} },
			deleteLink: '',
			isYaml: false,
			isJSON: false,
			deleteOverlay: {
				visible : false,
				processing : false,
				timer: null
			},
			navigateStack: [],
			errorcheck: false,
			errorMessage: "",
			errorM: [],
			delay: 0,
			disabled: false,
			selfLink: "",
		}
	},
	computed: {
		loader() {
			if (!this.src) return null;
			return () => import(`@/pages/${this.src}`)
		}
	},
	watch: {
		isYaml() {
			return this.raw = Object.assign({}, this.origin)
		},
		raw() {
			if(!this.raw) {
				this.raw = { metadata: {}, spec: {} }
			}
		},
		value(newVal) {
			this.navigateStack = [];
			this.navigate(newVal);
		},
		url(newVal) {
			if(!newVal) return;
			this.onSync();
		},
	},
	mounted() {
		this.$emit("close");
	},
	methods: {
		// 보기 페이지 이동
		navigate(loc) {
			if (loc.src && this.src != loc.src) {
				this.src = loc.src;
				if(this.loader) {
					this.loader()
						.then(() => {
							this.component = () => this.loader();
							this.isJSON = false;
							if(this.url != loc.url) this.url = loc.url;
						})
						.catch((ex) => {
							console.error(ex)
						})
				}
			} else if(this.url != loc.url) {
				this.isJSON = (!loc.src);	//if parameter 'src' is empty  ->  display 'json'
				this.url = loc.url;
			}

			if(this.name != loc.name) this.name = loc.name;
			if(this.title != loc.title) this.title = loc.title;
		},
		// component 페이지 이동 이벤트 처리
		onNavigate(loc) {
			this.navigateStack.push( { src : this.src, url : this.url, title: this.title, name : this.name } );
			this.navigate(loc);
		},
		// 뒤로 버튼
		onNavigateBack() {
			const pos = (this.navigateStack.length-1);
			if(pos >= 0) {
				this.navigate( this.navigateStack[pos] );
				this.navigateStack = this.navigateStack.splice(0, pos);
			}
		},
		// 조회
		onSync() {
			this.deleteOverlay.visible = false;
			this.isYaml = false;
			if (this.delay === 1) return;
			this.delay++;
			this.$axios.get(this.url)
				.then(resp => {
					this.errorcheck = false;
					this.origin = Object.assign({}, resp.data);
					this.raw = resp.data;
					if(this.title === "ConfigMap" || this.title === "Secret") this.raw.spec = resp.data.data || {};
					else if(this.title === "StorageClass") this.raw.spec = {}; // 무시
					else if(this.title === "Role" || this.title === "ClusterRole") this.raw.spec = resp.data.rules || {};
					else if(this.title === "RoleBinding" || this.title === "ClusterRoleBinding") this.raw.spec = { subjects: resp.data.subjects, roleRef: resp.data.roleRef} || {} ;
					else if(this.title === "ServiceAccount") this.raw.spec = resp.data.secrets || {};
					else if(this.title === "Endpoints") this.raw.spec = resp.data.subsets || {};
					else if(this.title === "HorizontalPodAutoscaler") this.raw.spec = {spec: resp.data.spec, status: resp.data.status}
					else if(this.title === "Node") this.raw.spec = resp.data.status.nodeInfo
					else this.raw.spec = resp.data.spec || {};
					this.delay = 0;
					// calcuate selfLink (deprecated metadata.selfLink)
					let c = this.getResource(this.raw)
					this.selfLink = this.getApiUrl(c.group, c.resource, this.raw.metadata.namespace, c.name);
					if(this.$el && this.$el.parentElement) this.$el.parentElement.scrollTop = 0;
				})
				.catch(e => {
					this.isError(e);
				});
		},
		isError(e) {
			this.errorcheck = true;
			this.raw = { metadata: {}, spec: {} };
			this.title = ""
			this.delay = 0
			this.errorMessage = e.response.data.message ;
		},
		onApply() {
			if(this.disabled) {
				this.msghttp(this.errorM)
				return this.disabled = false
			}
			this.$axios.put(`/raw/clusters/${this.currentContext()}`, this.raw)
				.then( resp => {
					if(!resp.data) {
						return this.toast('Invalid modification.','warning')
					}
					this.origin = Object.assign({}, resp.data);
					this.raw = resp.data;
					this.toast('Patch Successful')
				})
				.catch(e => {this.msghttp(e);});
		},
		onDelete() {
			this.deleteLink = this.selfLink
			this.deleteOverlay.processing = true;
			this.$axios.delete(this.selfLink)
				.then( _ => {
					this.watch();
				})
				.catch(e => {
					this.deleteOverlay.visible = false; 
					this.deleteOverlay.processing = false;
					this.msghttp(e);
				});
		},
		watch() {
			if(this.selfLink !== this.deleteLink) {
				this.deleteOverlay.visible = false
				this.deleteOverlay.processing = false
			} else {
				this.deleteOverlay.visible = true
				this.deleteOverlay.processing = true
			}
			this.$axios.get(this.deleteLink)
				.then( resp => {
					if (resp.status === 404) {
						if(this.selfLink === this.deleteLink) {
							this.deleteOverlay.visible = false;
							this.deleteOverlay.processing = false;
							this.$emit("delete");
							this.$emit("close");
						} else this.toast('delete Successful')
					} else {
						setTimeout(() => { this.watch(); }, 2000);
					}
				}).catch( error => {
				if(error.response && error.response.status === 404) {
					if(this.selfLink === this.deleteLink) {
						this.deleteOverlay.visible = false;
						this.deleteOverlay.processing = false;
						this.$emit("delete");
						this.$emit("close");
					} else this.toast('delete Successful')
				} else {
					this.msghttp(error);
				}
			});
		},
		onError(error) {
			this.disabled = true;
			this.errorM = error
			// this.msghttp(error)
		},
		onReset() {
			this.raw = Object.assign({}, this.origin);
			this.disabled = false
		},
	}
}

</script>
<style scoped>
.floating { position: fixed; margin: 0 auto; left: 0; right: 0; top: 5%;}
</style>
