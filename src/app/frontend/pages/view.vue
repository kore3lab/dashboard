<template>
	<section class="content border-primary border-left" v-bind:class="{ 'h-100' : errorcheck}">
		<div class="card card-primary m-0 layer" v-bind:class="{ 'h-100' : errorcheck}">
			<!-- card-header -->
			<div class="card-header pt-2 pb-2 sticky-top" style="position:sticky">
				<h3 class="card-title text-truncate">{{ title }}</h3>
				<div class="card-tools">
					<span v-show="!errorcheck">
						<button type="button" class="btn btn-tool" @click="onSync()"><i class="fas fa-sync-alt"></i></button>
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
				<!-- summary -->
				<div v-show="!isYaml && !errorcheck" class="card-body p-2">
					<div class="row">
						<div class="col-md-12">
							<div class="card card-secondary card-outline">
								<div class="card-header p-2"><h3 class="card-title text-md">Meta data</h3></div>
								<div class="card-body p-2">
									<dl class="row mb-0">
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
										<dt class="col-sm-2 text-truncate">Create at</dt><dd class="col-sm-10">{{ this.getTimestampString(raw.metadata.creationTimestamp)}} ago ({{ raw.metadata.creationTimestamp }})</dd>
										<dt class="col-sm-2 text-truncate">UID</dt><dd class="col-sm-10">{{ raw.metadata.uid }}</dd>
									</dl>
								</div>
							</div>
						</div>
					</div>
					<div class="row" v-if="crd !== 'Storage Class'">
						<div class="col-md-12">
							<div class="card card-secondary card-outline m-0">
								<div class="card-header p-2">
									<h3 class="card-title text-md" v-if="crd==='Config Map' || crd==='Secret'">Data</h3>
									<h3 class="card-title text-md" v-else-if="crd==='Service Account'">Secrets</h3>
									<h3 class="card-title text-md" v-else-if="crd==='Role' || crd==='Cluster Role'">Rules</h3>
									<h3 class="card-title text-md" v-else-if="crd==='Role Binding' || crd==='Cluster Role Binding'">Subjects,RoleRef</h3>
									<h3 class="card-title text-md" v-else-if="crd==='Endpoint'">Subsets</h3>
									<h3 class="card-title text-md" v-else-if="crd==='HPA'">Spec,Status</h3>
									<h3 class="card-title text-md" v-else-if="crd==='Node'">NodeInfo</h3>
									<h3 class="card-title text-md" v-else>Specification</h3>
								</div>
								<c-jsontree id="txtSpec" v-model="raw.spec" class="card-body p-2"></c-jsontree>
							</div>
						</div>
					</div>
					<div class="row" v-if="showEvent">
						<div class="col-md-12">
							<div class="card card-secondary card-outline">
								<div class="card-header p-2"><h3 class="card-title text-md">Events</h3></div>
								<div class="card-body p-2">
									<dl v-for="(val,idx) in event" v-bind:key="idx" class="row mb-0 card-body p-2 border-top">
										<dt class="col-sm-12 mb-1 text-truncate">{{ val.name }}</dt>
										<dt class="col-sm-2 text-truncate">Source</dt><dd class="col-sm-10">{{ val.source }}</dd>
										<dt class="col-sm-2 text-truncate">Count</dt><dd class="col-sm-10">{{ val.count }}</dd>
										<dt class="col-sm-2 text-truncate">Sub-object</dt><dd class="col-sm-10">{{ val.subObject }}</dd>
										<dt class="col-sm-2 text-truncate">Last seen</dt><dd class="col-sm-10">{{ val.lastSeen }}</dd>
									</dl>
								</div>
							</div>
						</div>
					</div>
				</div>
				<!-- yaml tab -->
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
				<!-- delete overlay -->
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
	props:["crd", "badge", "group", "name", "url","visible"],
	components: {
		"c-aceeditor": { extends: VueAceEditor },
		"c-jsontree": { extends: VueJsonTree },
	},
	data() {
		return {
			origin: { metadata: {}, spec: {} },
			raw: { metadata: {}, spec: {} },
			event: [],
			showEvent: false,
			isYaml: false,
			deleteOverlay: {
				visible : false,
				processing : false,
				timer: null
			},
			localUrl: "",
			title: "",
			errorcheck: false,
			errorMessage: "",
		}
	},
	created: function() {
		window.addEventListener('click',this.clickCheck)
	},
	watch: {
		url(newVal) {
			if(newVal && (newVal !== this.localUrl)) {
				this.onSync();
				document.getElementsByClassName(`b-sidebar-body`)[0].scrollTop = 0
				this.localUrl = newVal;
			}
		}
	},
	methods: {
		// 조회
		onSync() {
			axios.get(this.url)
					.then( resp => {
						this.errorcheck = false;
						this.showEvent = false;
						this.origin = Object.assign({}, resp.data);
						this.raw = resp.data;
						this.event = this.getEvents(resp.data.metadata.uid)
						this.title = this.crd +' / '+this.name;
						if(this.crd === "Config Map" || this.crd === "Secret") this.raw.spec = resp.data.data || {};
						else if(this.crd === "Storage Class") this.raw.spec = null; // 무시
						else if(this.crd === "Role" || this.crd === "Cluster Role") this.raw.spec = resp.data.rules || {};
						else if(this.crd === "Role Binding" || this.crd === "Cluster Role Binding") this.raw.spec = { subjects: resp.data.subjects, roleRef: resp.data.roleRef} || {} ;
						else if(this.crd === "Service Account") this.raw.spec = resp.data.secrets || {};
						else if(this.crd === "Endpoint") this.raw.spec = resp.data.subsets || {};
						else if(this.crd === "HPA") this.raw.spec = {spec: resp.data.spec, status: resp.data.status}
						else if(this.crd === "Node") this.raw.spec = resp.data.status.nodeInfo
						else this.raw.spec = resp.data.spec || {};
					})
					.catch(e => {
						// this.msghttp(e);
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
			axios.put(this.backendUrl() + '/raw/clusters/' + this.currentContext(), this.raw)
					.then( resp => {
						this.origin = Object.assign({}, resp.data);
						this.raw = resp.data;
					})
					.catch(e => { this.msghttp(e);});
		},
		onDelete() {
			this.deleteOverlay.processing = true;
			axios.delete(`${this.backendUrl() + '/raw/clusters/' + this.currentContext()}${this.raw.metadata.selfLink}`)
					.then( _ => {
						this.watch();
					})
					.catch(e => { this.msghttp(e);});
		},
		watch() {
			axios.get(`${this.backendUrl() + '/raw/clusters/' + this.currentContext()}${this.raw.metadata.selfLink}`)
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
		getEvents(uid) {
			let events = [];
			axios.get(this.getApiUrl('events.k8s.io','events'))
					.then( resp => {
						for(let i=0; i<resp.data.items.length; i++) {
							if(resp.data.items[i].regarding.uid === uid) {
								this.showEvent = true;
								events.push({
									name: resp.data.items[i].note || "",
									source: resp.data.items[i].deprecatedSource.host || resp.data.items[i].deprecatedSource.component || "",
									count: resp.data.items[i].deprecatedCount || "0",
									subObject: resp.data.items[i].regarding.fieldPath || "",
									lastSeen: resp.data.items[i].deprecatedLastTimestamp || ""
								})
							}
						}
					})
			return events
		}
	}
}

</script>
