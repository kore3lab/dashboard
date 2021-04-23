<template>
	<div class="card-body p-2">
		<div class="row">
			<div class="col-md-12">
				<div class="card card-secondary card-outline">
					<div class="card-body p-2">
						<dl class="row mb-0">
							<dt class="col-sm-2 text-truncate">Create at</dt><dd class="col-sm-10">{{ this.getTimestampString(metadata.creationTimestamp)}} ago ({{ metadata.creationTimestamp }})</dd>
							<dt class="col-sm-2">Name</dt><dd class="col-sm-10">{{ metadata.name }}</dd>
							<dt class="col-sm-2">Namespace</dt><dd class="col-sm-10">{{ metadata.namespace }}</dd>
							<dt class="col-sm-2">Annotations</dt>
							<dd class="col-sm-10 text-truncate">
								<ul class="list-unstyled mb-0">
									<li v-for="(value, name) in metadata.annotations" v-bind:key="name"><span class="badge badge-secondary font-weight-light text-sm mb-1">{{ name }}:{{ value }}</span></li>
								</ul>
							</dd>
							<dt class="col-sm-2">Labels</dt>
							<dd class="col-sm-10 text-truncate">
								<ul class="list-unstyled mb-0">
									<li v-for="(value, name) in metadata.labels" v-bind:key="name"><span class="badge badge-secondary font-weight-light text-sm mb-1">{{ name }}:{{ value }}</span></li>
								</ul>
							</dd>
							<dt v-if="metadata.ownerReferences" class="col-sm-2 text-truncate">Controlled By</dt>
							<dd v-if="metadata.ownerReferences" class="col-sm-10">{{ metadata.ownerReferences[0].kind }} <a href="#" @click="$emit('navigate', getViewLink(controller.g, controller.k, metadata.namespace, metadata.ownerReferences[0].name))">{{ metadata.ownerReferences[0].name }}</a></dd>
						</dl>
					</div>
				</div>
			</div>
		</div>

		<div v-show="configData" class="row">
			<div class="col-md-12">
				<div class="card card-secondary card-outline m-0">
					<div class="card-header p-2"><h3 class="card-title text-md">Data</h3></div>
					<div class="card-body p-2">
						<dl v-for="(val, idx) in configData" v-bind:key="idx" class="row mb-0 card-body p-2">
							<dt class="col-sm-12"><span class="card-title mb-2">{{ val.key }}</span></dt>
							<dd class="col-sm-12"><b-form-textarea id="txtSpec" max-rows="10" v-model="val.val" class="card-body p-2 border text-sm"></b-form-textarea></dd>
						</dl>
						<b-button variant="primary" size="sm" class="mb-1 ml-2" @click="onSave()">Save</b-button>
					</div>
				</div>
			</div>
		</div>


		<div class="row">
			<div class="col-md-12">
				<div class="card card-secondary card-outline m-0">
					<div class="card-header p-2"><h3 class="card-title text-md">Events</h3></div>
					<div class="card-body p-2">
						<dl v-for="(val, idx) in event" v-bind:key="idx" class="row mb-0 card-body p-2 border-bottom">
							<dt class="col-sm-12"><p v-bind:class="val.type" class="mb-1">{{ val.name }}</p></dt>
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
</template>
<script>
import VueAceEditor from "@/components/aceeditor";

export default {
	components: {
		"c-aceeditor": { extends: VueAceEditor },
	},
	data() {
		return {
			metadata: {},
			event: [],
			info: [],
			origin: [],
			controller: [],
			configData: [],
			fields: [
				{ key: "name", label: "Name" },
				{ key: "endpoints", label: "Endpoints" },
			],
		}
	},
	mounted() {
		this.$nuxt.$on("onReadCompleted", (data) => {
			if(!data) return
			this.origin = data;
			this.metadata = data.metadata;
			this.onSync(data)
		});
		this.$nuxt.$emit("onCreated",'')
	},
	methods: {
		onSync(data) {
			this.event = this.getEvents(data.metadata.uid);
			this.controller = this.getController(data.metadata.ownerReferences);
			this.configData = this.getData(data.data);
		},
		getData(data) {
			if(!data) return false
			let list = [];
			let key = Object.keys(data)
			key.forEach(el => {
				list.push({
					key: el,
					val: data[el]
				})
			})
			return list
		},
		onSave() {
			let list = {};
			this.configData.forEach(el => {
				list[el.key] = el.val
			})
			this.origin.data = list
			this.$axios.put(`/raw/clusters/${this.currentContext()}`, this.origin)
					.then( resp => {
						this.origin = Object.assign({}, resp.data);
						this.toast(`ConfigMap ${ this.metadata.name } successfully updated.`, "success");
					})
					.catch(e => { this.msghttp(e);});
		},
	},
	beforeDestroy(){
		this.$nuxt.$off("onReadCompleted");
	},
}
</script>
