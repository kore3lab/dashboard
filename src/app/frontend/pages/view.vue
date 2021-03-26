<template>
	<section class="content border-primary border-left" v-bind:class="{ 'h-100' : errorcheck}">
		<div class="card card-primary m-0 layer" v-bind:class="{ 'h-100' : errorcheck}">
			<!-- card-header -->
			<div class="card-header pt-2 pb-2 sticky-top" style="position:sticky">
				<h3 class="card-title text-truncate">{{ title }} / {{ name }}</h3>
				<div class="card-tools">
					<span v-show="!errorcheck">
						<button type="button" class="btn btn-tool" @click="onSync()"><i class="fas fa-sync-alt"></i></button>
						<button type="button" class="btn btn-tool" v-show="isJSON && component"  @click="isJSON=false"><i class="fas fa-list-alt"></i></button>
						<button type="button" class="btn btn-tool" v-show="!isJSON && component" @click="isJSON=true"><i>JSON</i></button>
						<button type="button" class="btn btn-tool" @click="isYaml=true"><i class="fas fa-edit"></i></button>
						<button type="button" class="btn btn-tool" @click="deleteOverlay.visible = true"><i class="fas fa-trash"></i></button>
					</span>
					<button type="button" class="btn btn-tool" @click="$emit('close')"><i class="fas fa-times"></i></button>
				</div>
			</div>
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
						<div class="col-sm-12 text-right">
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
			isYaml: false,
			isJSON: false,
			deleteOverlay: {
				visible : false,
				processing : false,
				timer: null
			},
			localUrl: "",
			localSrc: "",
			errorcheck: false,
			errorMessage: "",
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
	created: function() {
		window.addEventListener('click',this.clickCheck)
	},
	watch: {
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
	methods: {
		navigate(loc) {
			this.value.name = loc.name;
			this.value.title = loc.title;
			this.src = loc.src;
			this.url = loc.url;
		},
		// 조회
		onSync() {
			axios.get(this.localUrl)
					.then( resp => {
						this.errorcheck = false;
						this.origin = Object.assign({}, resp.data);
						this.raw = resp.data;
						if(this.title === "ConfigMap" || this.title === "Secret") this.raw.spec = resp.data.data || {};
						else if(this.title === "StorageClass") this.raw.spec = null; // 무시
						else if(this.title === "Role" || this.title === "ClusterRole") this.raw.spec = resp.data.rules || {};
						else if(this.title === "RoleBinding" || this.title === "ClusterRoleBinding") this.raw.spec = { subjects: resp.data.subjects, roleRef: resp.data.roleRef} || {} ;
						else if(this.title === "ServiceAccount") this.raw.spec = resp.data.secrets || {};
						else if(this.title === "Endpoints") this.raw.spec = resp.data.subsets || {};
						else if(this.title === "HorizontalPodAutoscaler") this.raw.spec = {spec: resp.data.spec, status: resp.data.status}
						else if(this.title === "Node") this.raw.spec = resp.data.status.nodeInfo
						else this.raw.spec = resp.data.spec || {};
						this.$nuxt.$emit("onReadCompleted", this.origin);
					})
					.catch(e => {
						this.isError(e);
					});
		},
		isError(e) {
			this.errorcheck = true;
			this.raw = { metadata: {}, spec: {} };
			this.title = ""
			this.errorMessage = e.response.data.message ;
		},
		onApply() {
			axios.put(`${this.backendUrl()}/raw/clusters/${this.currentContext()}`, this.raw)
					.then( resp => {
						this.origin = Object.assign({}, resp.data);
						this.raw = resp.data;
					})
					.catch(e => { this.msghttp(e);});
		},
		onDelete() {
			this.deleteOverlay.processing = true;
			axios.delete(`${this.backendUrl()}/raw/clusters/${this.currentContext()}${this.raw.metadata.selfLink}`)
					.then( _ => {
						this.watch();
					})
					.catch(e => { this.msghttp(e);});
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
			this.mesbox(error.message);
		},
		onReset() {
			this.raw = Object.assign({}, this.origin);
		},
		clickCheck(el) {
			if($(el.target).closest(`.layer`).length === 0 && $(el.target).closest(`a`).length === 0 && $(el.target).closest(`.b-sidebar-body`).length === 0) {
				this.$emit('close')
			}
		},
	}
}

</script>
