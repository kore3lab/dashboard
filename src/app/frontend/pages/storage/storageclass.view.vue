<template>
<div>
	<!-- 1. metadata -->
	<c-metadata v-model="value" dtCols="3" ddCols="9">
		<dt v-if="metadata.finalizers" class="col-sm-3">Finalizers</dt>
		<dd v-if="metadata.finalizers" class="col-sm-9">
			<ul class="list-unstyled mb-0">
				<li v-for="(value, name) in metadata.finalizers" v-bind:key="name">{{ value }}</li>
			</ul>
		</dd>
		<dt class="col-sm-3">Provisioner</dt><dd class="col-sm-9">{{ info.provisioner }}</dd>
		<dt class="col-sm-3">Volume Binding Mode</dt><dd class="col-sm-9">{{ info.volumeBindingMode }}</dd>
		<dt class="col-sm-3">Reclaim Policy</dt><dd class="col-sm-9">{{ info.reclaimPolicy }}</dd>
		<dt v-if="info.mountOption" class="col-sm-3">Mount Options</dt><dd v-if="info.mountOption" class="col-sm-9">{{ info.mountOption }} </dd>
	</c-metadata>
	<!-- 2. parameters -->
	<div v-if="parameters" class="row">
		<div class="col-md-12">
			<div class="card card-secondary card-outline">
				<div class="card-header p-2"><h3 class="card-title">Parameters</h3></div>
				<div class="card-body p-2">
					<dl v-for="(value, name) in parameters" v-bind:key="name" class="row mb-0">
						<dt class="col-sm-3">{{ name }}</dt>
						<dd class="col-sm-9">{{ value }}</dd>
					</dl>
				</div>
			</div>
		</div>
	</div>
	<!-- 3. persistent volumnes -->
	<div v-if="pvList.length>0" class="row">
		<div class="col-md-12">
			<div class="card card-secondary card-outline">
				<div class="card-header p-2"><h3 class="card-title">Persistent Volumes</h3></div>
				<div class="card-body p-2 overflow-auto">
					<b-table small :items="pvList" :fields="fields">
						<template v-slot:cell(name)="data">
							<a href="#" @click="$emit('navigate', getViewLink('', 'persistentvolumes','', data.item.name))">{{ data.item.name }}</a>
						</template>
						<template v-slot:cell(status)="data">
							<span v-bind:class="{'text-success': data.value=='Bound', 'text-warning': data.value=='Released'} ">{{ data.value }}</span>
						</template>
					</b-table>
				</div>
			</div>
		</div>
	</div>
	<!-- 4. events -->
	<c-events v-model="value" class="row"></c-events>

</div>
</template>
<script>
import VueMetadataView	from "@/components/view/metadataView.vue";
import VueEventsView	from "@/components/view/eventsView.vue";

export default {
	props:["value"],
	components: {
		"c-metadata": { extends: VueMetadataView },
		"c-events": { extends: VueEventsView }
	},
	data() {
		return {
			metadata: {},
			info: [],
			parameters: [],
			pvList: [],
			fields: [
				{ key: "name", label: "Name", sortable: true },
				{ key: "capacity", label: "Capacity", sortable: true, thClass:"text-right", tdClass:"text-right" },
				{ key: "status", label: "Status", sortable: true  },
			],
		}
	},
	watch: {
		value(d) { this.onSync(d) }
	},
	methods: {
		onSync(data) {
			if(!data) return
			this.metadata = data.metadata;
			this.info = {
				provisioner: data.provisioner,
				volumeBindingMode: data.volumeBindingMode,
				reclaimPolicy: data.reclaimPolicy,
				mountOption: data.mountOptions? data.mountOptions.join(",") : '',
			};
			this.parameters = data.parameters;

			// pv-list
			this.$axios.get(this.getApiUrl('','persistentvolumes'))
				.then(resp => {
					let list = [];
					resp.data.items.forEach(el => {
						if(el.spec.storageClassName && el.spec.storageClassName === this.metadata.name) {
							list.push({
								name: el.metadata.name,
								capacity: el.spec.capacity.storage,
								status: el.status.phase
							})
						}
					})
					this.pvList = list
			})
		}
	}
}
</script>
