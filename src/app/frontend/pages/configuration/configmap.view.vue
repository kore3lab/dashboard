<template>
<div>
	<!-- 1. metadata -->
	<c-metadata v-model="metadata" dtCols="2" ddCols="10">
		<dt v-if="metadata.ownerReferences" class="col-sm-2">Controlled By</dt>
		<dd v-if="metadata.ownerReferences" class="col-sm-10">{{ controller.kind }} <a href="#" @click="$emit('navigate', getViewLink(controller.group, controller.resource, metadata.namespace, controller.name))">{{ controller.name }}</a></dd>
	</c-metadata>
	<!-- 2. data -->
	<div v-show="configData" class="row">
		<div class="col-md-12">
			<div class="card card-secondary card-outline m-0">
				<div class="card-header p-2"><h3 class="card-title">Data</h3></div>
				<div class="card-body group">
					<ul>
						<li v-for="(val, idx) in configData" v-bind:key="idx">
							<p class="title">{{ val.key }}</p>
							<b-form-textarea id="txtSpec" max-rows="10" v-model="val.val" class="card-body p-2 border text-sm"></b-form-textarea>
						</li>
					</ul>
					<b-button variant="primary" size="sm" class="m-1" @click="onSave()">Save</b-button>
				</div>
			</div>
		</div>
	</div>
	<!-- 4. events -->
	<c-events class="row" v-model="metadata.uid"></c-events>
</div>
</template>
<script>
import VueMetadataView	from "@/components/view/metadataView.vue";
import VueEventsView	from "@/components/view/eventsView.vue";

export default {
	components: {
		"c-metadata": { extends: VueMetadataView },
		"c-events": { extends: VueEventsView }
	},
	data() {
		return {
			metadata: {},
			info: [],
			origin: [],
			controller: {},
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
			this.controller = data.metadata.ownerReferences? this.getResource(data.metadata.ownerReferences[0]): {};
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
