<template>
<div>
	<!-- 1. metadata -->
	<div class="row">
		<div class="col-md-12">
			<div class="card card-secondary card-outline">
				<div class="card-body p-2">
					<dl class="row mb-0">
						<dt class="col-sm-3">Create at</dt><dd class="col-sm-9">{{ this.getTimestampString(metadata.creationTimestamp)}} ago ({{ metadata.creationTimestamp }})</dd>
						<dt class="col-sm-3">Name</dt><dd class="col-sm-9">{{ metadata.name }}</dd>
						<dt class="col-sm-3">Annotations</dt>
						<dd class="col-sm-9">
							<ul class="list-unstyled mb-0">
								<li v-for="(value, name) in metadata.annotations" v-bind:key="name">{{ name }}=<span class="font-weight-light">{{ value }}</span></li>
							</ul>
						</dd>
						<dt class="col-sm-3">Labels</dt>
						<dd class="col-sm-9">
							<span v-for="(value, name) in metadata.labels" v-bind:key="name" class="label">{{ name }}={{ value }}</span>
						</dd>
						<dt class="col-sm-3">UID</dt><dd class="col-sm-9">{{ metadata.uid }}</dd>
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
					</dl>
				</div>
			</div>
		</div>
	</div>
	<!-- 2. parameters -->
	<div v-if="param" class="row">
		<div class="col-md-12">
			<div class="card card-secondary card-outline">
				<div class="card-header p-2"><h3 class="card-title">Parameters</h3></div>
				<div class="card-body p-2">
					<dl v-for="(val, idx) in param" v-bind:key="idx" class="row mb-0">
						<dt class="col-sm-3">{{ val.type? val.type[0].toUpperCase()+val.type.slice(1): '' }}</dt><dd class="col-sm-9">{{ val.val }}</dd>
					</dl>
				</div>
			</div>
		</div>
	</div>
	<!-- 3. persistent volumnes -->
	<div v-if="isPV" class="row">
		<div class="col-md-12">
			<div class="card card-secondary card-outline">
				<div class="card-header p-2"><h3 class="card-title">Persistent Volumes</h3></div>
				<div class="card-body p-2 overflow-auto">
					<b-table striped hover small :items="pvList" :fields="fields">
						<template v-slot:cell(name)="data">
							<a href="#" @click="$emit('navigate', getViewLink('', 'persistentvolumes','', data.item.name))">{{ data.item.name }}</a>
						</template>
						<template v-slot:cell(status)="data">
							<span v-bind:class="data.value==='Bound' ? 'text-success': 'text-warning'">{{ data.value }}</span>
						</template>
					</b-table>
				</div>
			</div>
		</div>
	</div>
	<!-- 4. events -->
	<c-events class="row" v-model="metadata.uid"></c-events>

</div>
</template>
<script>
import VueEventsView	from "@/components/view/eventsView.vue";

export default {
	components: {
		"c-events": { extends: VueEventsView }
	},
	data() {
		return {
			metadata: {},
			info: [],
			param: [],
			pvList: [],
			isPV: false,
			fields: [
				{ key: "name", label: "Name", sortable: true },
				{ key: "capacity", label: "Capacity", sortable: true  },
				{ key: "status", label: "Status", sortable: true  },
			],
		}
	},
	mounted() {
		this.$nuxt.$on("onReadCompleted", (data) => {
			if(!data) return
			this.metadata = data.metadata;
			this.onSync(data)
		});
		this.$nuxt.$emit("onCreated",'')
	},
	methods: {
		onSync(data) {
			this.info = this.getInfo(data);
			this.param = this.getParam(data.parameters);
			this.pvList = this.getPv();

		},
		getInfo(data) {
			return {
				provisioner: data.provisioner,
				volumeBindingMode: data.volumeBindingMode,
				reclaimPolicy: data.reclaimPolicy,
				mountOption: data.mountOptions? data.mountOptions.join(",") : '',
			}
		},
		getParam(p) {
			let list = [];
			if(!p) return

			Object.entries(p).map(([name, val]) => {
				list.push({
					type: name,
					val: val
				})
			})
			return list
		},
		getPv() {
			this.isPV = false;
			let list = [];
			this.$axios.get(this.getApiUrl('','persistentvolumes'))
			.then(resp => {
				resp.data.items.forEach(el => {
					if(el.spec.storageClassName && el.spec.storageClassName === this.metadata.name) {
						this.isPV = true;
						list.push({
							name: el.metadata.name,
							capacity: el.spec.capacity.storage,
							status: el.status.phase
						})
					}
				})
				this.pvList = list
				return list
			})
		},
	},
	beforeDestroy(){
		this.$nuxt.$off("onReadCompleted");
	},
}
</script>
