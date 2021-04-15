<template>
	<section class="content border-primary border-left min-vh-100" v-bind:class="{ 'min-vh-100' : errorcheck}">
		<div class="card card-primary m-0 layer" v-bind:class="{ 'min-vh-100' : errorcheck}">
			<!-- card-header -->
			<div class="card-header pt-2 pb-2 sticky-top" style="position:sticky">
				<h3 class="card-title text-truncate">{{ title }} / {{ name }}</h3>
				<div class="card-tools">
					<span v-show="!errorcheck">
						<button type="button" class="btn btn-tool" @click="onSync()"><i class="fas fa-sync-alt"></i></button>
						<button type="button" class="btn btn-tool" v-show="isJSON && component"  @click="isJSON=false"><i class="fas fa-list-alt"></i></button>
						<button type="button" class="btn btn-tool" v-show="!isJSON && component" @click="isJSON=true"><i>JSON</i></button>
						<button type="button" class="btn btn-tool" @click="isYaml=true"><i class="fas fa-edit"></i></button>
						<button id="terminal" type="button" class="btn btn-tool" v-show="isTerminal"  @click="openTerminal()"><i class="fas fa-terminal"></i></button>
						<button type="button" class="btn btn-tool" @click="deleteOverlay.visible = true"><i class="fas fa-trash"></i></button>
					</span>
					<button type="button" class="btn btn-tool" @click="$emit('close')"><i class="fas fa-times"></i></button>
				</div>
			</div>
			<b-popover triggers="hover" target="terminal" placement="bottomleft" boundary="window" boundary-padding="0">
				<ul class="list-unstyled m-0">
					<li v-for="(val,idx) in containers" v-bind:key="idx" class="mb-1">
						<span v-if="val.status.value === 'running'" class="text-truncate"><button type="button" class="btn btn-tool" @click="openTerminal(val)"><b-badge :variant="val.status.badge" class="mt-1 mb-1 mr-1">&nbsp;</b-badge>{{ val.name }}</button></span>
					</li>
				</ul>
			</b-popover>
			<!-- error message-->
			<div class="card-body" v-show="errorcheck" style="padding-top: 50%">
				<div class="col-md-12 m-3 text-sm-center"><p>Resource loading has failed: <b>{{ errorMessage }}</b></p></div>
			</div>

			<b-overlay :show="deleteOverlay.visible" rounded="sm" no-center>
				<!--1. not Yaml(editor) -->
				<div v-show="!isYaml && !errorcheck" class="card-body p-2">
					<!-- 1.1 meta data -->
					<div class="row" v-show="isJSON">
						<div class="col-md-12">
							<div class="card card-secondary card-outline">
								<div class="card-body p-2">
									<dl class="row mb-0">
										<dt class="col-sm-2 text-truncate">Create at</dt><dd class="col-sm-10">{{ this.getTimestampString(raw.metadata.creationTimestamp)}} ago ({{ raw.metadata.creationTimestamp }})</dd>
										<dt class="col-sm-2">Name</dt><dd class="col-sm-10">{{ raw.metadata.name }}</dd>
										<dt v-if="raw.metadata.namespace" class="col-sm-2">Namespace</dt><dd v-if="raw.metadata.namespace" class="col-sm-10">{{ raw.metadata.namespace }}</dd>
										<dt class="col-sm-2 text-truncate">Annotations</dt>
										<dd class="col-sm-10 text-truncate">
											<ul class="list-unstyled mb-0">
												<li v-for="(value, name) in raw.metadata.annotations" v-bind:key="name"><span class="badge badge-secondary font-weight-light text-sm mb-1">{{ name }}:{{ value }}</span></li>
											</ul>
										</dd>
										<dt class="col-sm-2 text-truncate">Labels</dt>
										<dd class="col-sm-10">
											<ul class="list-unstyled mb-0">
												<li v-for="(value, name) in raw.metadata.labels" v-bind:key="name"><span class="badge badge-secondary font-weight-light text-sm mb-1">{{ name }}:{{ value }}</span></li>
											</ul>
										</dd>
										<dt class="col-sm-2 text-truncate">UID</dt><dd class="col-sm-10">{{ raw.metadata.uid }}</dd>
									</dl>
								</div>
							</div>
						</div>
					</div>
					<!-- 1.2  custom view -->
					<component :is="component" v-if="component" v-show="!isJSON" v-model="value" @navigate="navigate"/>
					<!-- 1.3 json-tree -->
					<div class="row" v-show="isJSON && title !== 'StorageClass'">
						<div class="col-md-12">
							<div class="card card-secondary card-outline m-0">
								<div class="card-header p-2">
									<h3 class="card-title text-md" v-if="title==='ConfigMap' || title==='Secret'">Data</h3>
									<h3 class="card-title text-md" v-else-if="title==='ServiceAccount'">Secrets</h3>
									<h3 class="card-title text-md" v-else-if="title==='Role' || title==='ClusterRole'">Rules</h3>
									<h3 class="card-title text-md" v-else-if="title==='RoleBinding' || title==='ClusterRoleBinding'">Subjects,RoleRef</h3>
									<h3 class="card-title text-md" v-else-if="title==='Endpoints'">Subsets</h3>
									<h3 class="card-title text-md" v-else-if="title==='HorizontalPodAutoscaler'">Spec,Status</h3>
									<h3 class="card-title text-md" v-else-if="title==='Node'">NodeInfo</h3>
									<h3 class="card-title text-md" v-else>Specification</h3>
								</div>
								<c-jsontree id="txtSpec" v-model="raw.spec" class="card-body p-2"></c-jsontree>
							</div>
						</div>
					</div>
				</div>
				<!-- 2. if Yaml(editor) then  -->
				<div v-show="isYaml && !errorcheck" class="card-body p-1">
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
					<div v-if="deleteOverlay.processing" class="text-center">
						<b-spinner small class="mr-2" label="please wait"></b-spinner><span>Watching DELETE status...</span>
					</div>
					<div v-else class="text-center">
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

import axios			from "axios"
import VueAceEditor 	from "@/components/aceeditor"
import VueJsonTree 		from "@/components/jsontree"
export default {

	props:["value"],

	components: {
		"c-aceeditor": { extends: VueAceEditor },
		"c-jsontree": { extends: VueJsonTree },
	},
	data() {
		return {
			component: null,
			src: "",
			url: "",
			origin: { metadata: {}, spec: {} },
			raw: { metadata: {}, spec: {} },
			containers: [],
			containerCount: 0,
			isYaml: false,
			isJSON: false,
			isTerminal: false,
			deleteOverlay: {
				visible : false,
				processing : false,
				timer: null
			},
			localUrl: "",
			localSrc: "",
			errorcheck: false,
			errorMessage: "",
			isCreated: false,
			delay: 0,
			disabled: false,
		}
	},
	computed: {
		loader() {
			if (!this.src) return null;
			return () => import(`@/pages/${this.src}`)
		},
		title() {
			return this.value.title;
		},
		name() {
			return this.value.name;
		}
	},
	watch: {
		isYaml() {
			return this.raw = Object.assign({}, this.raw)
		},
		value(newVal) {
			this.src = newVal.src;
			this.url = newVal.url;
			document.getElementsByClassName(`b-sidebar-body`)[0].scrollTop = 0
		},
		src(newVal) {
			if(newVal !== this.localSrc) {
				this.component = null;
				this.isJSON = true;
				if(this.loader) {
					this.loader()
							.then(() => {
								this.component = () => this.loader();
								this.isJSON = false;
							})
							.catch((ex) => {
								console.error(ex)
							})
				}
				this.localSrc = newVal;
			}
		},
		url(newVal) {
			if(!newVal) return;
			if(newVal !== this.localUrl) {
				this.localUrl =  newVal;
				this.onSync();
			} else {
				this.onSync();
			}
		},
	},
	mounted() {
		this.$emit('close');
		this.$nuxt.$on("onCreated",() => {
			this.isCreated = true;
			this.onSync()
		})
		this.$nuxt.$on('Containers', (data) => {
			this.containerCount = 0
			this.containers = data
			if(!data) return
			data.forEach(el => {
				if(el.status.value === 'running') this.containerCount++;
			})
			if(this.containerCount === 0) this.isTerminal = false;
		})
	},
	beforeUpdate() {
		let el = document.getElementsByTagName("body")
		el[0].style.removeProperty('height')
	},
	methods: {
		openTerminal(val) {
			let routeData
			if(val) {
				routeData = this.$router.resolve({path: '/terminal', query: {termtype: "container",pod: this.name, namespace: this.raw.metadata.namespace, cluster: this.currentContext(),container:val.name}});
			} else routeData = this.$router.resolve({path: '/terminal', query: {termtype: "pod",pod: this.name, namespace: this.raw.metadata.namespace, cluster: this.currentContext(),}});
			window.open("about:blank").location.href = routeData.href
		},
		navigate(loc) {
			this.value.name = loc.name;
			this.value.title = loc.title;
			this.src = loc.src;
			this.url = loc.url;
		},
		// 조회
		onSync() {
			this.isYaml = false;
			if (this.delay === 1) return;
			this.delay++;
			axios.get(this.localUrl)
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
						this.isTerminal = this.title === 'Pod';
						if (this.isCreated) {
							this.$nuxt.$emit("onReadCompleted", this.origin);
						}
						this.delay = 0;
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
			if(this.disabled) return this.disabled = false
			axios.put(`${this.backendUrl()}/raw/clusters/${this.currentContext()}`, this.raw)
					.then( resp => {
						this.origin = Object.assign({}, resp.data);
						this.raw = resp.data;
						this.toast('Patch Successful')
					})
					.catch(e => {this.msghttp(e);});
		},
		onDelete() {
			this.deleteOverlay.processing = true;
			axios.delete(`${this.backendUrl()}/raw/clusters/${this.currentContext()}${this.raw.metadata.selfLink}`)
					.then( _ => {
						this.watch();
					})
					.catch(e => {this.msghttp(e);});
		},
		watch() {
			axios.get(`${this.backendUrl()}/raw/clusters/${this.currentContext()}${this.raw.metadata.selfLink}`)
					.then( resp => {
						if (resp.status === 404) {
							this.deleteOverlay.visible = false;
							this.deleteOverlay.processing = false;
							this.$emit("delete");
							this.$emit("close");
						} else {
							setTimeout(() => { this.watch(); }, 2000);
						}
					}).catch( error => {
				if(error.response && error.response.status === 404) {
					this.deleteOverlay.visible = false;
					this.deleteOverlay.processing = false;
					this.$emit("delete");
					this.$emit("close");
				} else {
					this.msghttp(error);
				}
			});
		},
		onError(error) {
			this.disabled = true;
			this.msghttp(error)
		},
		onReset() {
			this.raw = Object.assign({}, this.origin);
			this.disabled = false
		},
	},
	beforeDestroy(){
		this.$nuxt.$off("onCreated",'')
		this.$nuxt.$off('Containers','')
	}
}

</script>
