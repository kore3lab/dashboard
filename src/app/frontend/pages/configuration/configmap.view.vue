<template>
<div>
	<!-- 1. metadata -->
	<c-metadata dtCols="2" ddCols="10" @navigate="$emit('navigate', arguments[0])"></c-metadata>
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
	<c-events class="row"></c-events>
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
			configData: [],
			fields: [
				{ key: "name", label: "Name" },
				{ key: "endpoints", label: "Endpoints" },
			],
		}
	},
	mounted() {
		this.$nuxt.$on("view-data-read-completed", (data) => {
			if(!data) return
			this.origin = data;
			this.metadata = data.metadata;
			this.configData = this.getData(data.data);
		});
	},
	methods: {
		getData(data) {
			if(!data) return false
			let list = [];
			Object.keys(data).forEach(el => {
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
		this.$nuxt.$off("view-data-read-completed");
	},
}
</script>
