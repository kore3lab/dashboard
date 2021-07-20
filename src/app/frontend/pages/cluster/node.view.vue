<template>
<div>
	<!-- 1. charts -->
	<c-charts class="row" v-model="selectUrl"></c-charts>
	<!-- 2. metadata -->
	<div class="row">
		<div class="col-md-12">
			<div class="card card-secondary card-outline">
				<div class="card-body p-2">
					<dl class="row mb-0">
						<dt class="col-sm-3">Create at</dt><dd class="col-sm-9">{{ this.getTimestampString(metadata.creationTimestamp)}} ago ({{ metadata.creationTimestamp }})</dd>
						<dt class="col-sm-3">Name</dt><dd class="col-sm-9">{{ metadata.name }}</dd>
						<dt class="col-sm-3">Capacity</dt><dd class="col-sm-9">{{ info.capacity }}</dd>
						<dt class="col-sm-3">Allocatable</dt><dd class="col-sm-9">{{ info.allocatable }}</dd>
						<dt class="col-sm-3">Addresses</dt>
						<dd class="col-sm-9">
							<ul class="list-unstyled mb-0">
								<li v-for="(val, idx) in info.addresses" v-bind:key="idx">{{ val.type }}: {{ val.address }}</li>
							</ul>
						</dd>
						<dt class="col-sm-3">OS</dt><dd class="col-sm-9">{{ info.os }}</dd>
						<dt class="col-sm-3">OS Image</dt><dd class="col-sm-9">{{ info.osImage }}</dd>
						<dt class="col-sm-3">Kernel version</dt><dd class="col-sm-9">{{ info.kernelVersion }}</dd>
						<dt class="col-sm-3">Container runtime</dt><dd class="col-sm-9">{{ info.containerRuntime }}</dd>
						<dt class="col-sm-3">Kubelet version</dt><dd class="col-sm-9">{{ info.kubeletVersion }}</dd>
						<dt class="col-sm-3">UID</dt><dd class="col-sm-9">{{ metadata.uid }}</dd>
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
						<dt v-if="info.taints" class="col-sm-3">Taints</dt>
						<dd v-if="info.taints" class="col-sm-9">
							<span v-for="(val, idx) in info.taints" v-bind:key="idx" v-bind:class="val.style">{{ val.key }}: {{ val.effect }} ({{ val.value }})</span>
						</dd>
						<dt class="col-sm-3">Conditions</dt>
						<dd class="col-sm-9">
							<span v-for="(val, idx) in info.conditions" v-bind:key="idx" v-bind:class="val.style" class="badge font-weight-light text-sm mb-1 mr-1">{{ val.type }}</span>
						</dd>
					</dl>
				</div>
			</div>
		</div>
	</div>
	<!-- 3. pods -->
	<c-podlist class="row" v-model="selectUrl" namespace="true" @navigate="$emit('navigate',arguments[0])"></c-podlist>
	<!-- 4. events -->
	<c-events class="row" v-model="metadata.uid"></c-events>

</div>
</template>
<script>
import VueJsonTree		from "@/components/jsontree";
import VueEventsView	from "@/components/view/eventsView.vue";
import VueChartsView	from "@/components/view/metricsChartsView.vue";
import VuePodListView	from "@/components/view/podListView.vue";

export default {
	components: {
		"c-jsontree": { extends: VueJsonTree },
		"c-events": { extends: VueEventsView },
		"c-charts": { extends: VueChartsView },
		"c-podlist": { extends: VuePodListView }
	},
	data() {
		return {
			metadata: {},
			info: [],
			selectUrl: ""
		}
	},
	mounted() {
		this.$nuxt.$on("onReadCompleted", (data) => {
			if(!data) return
			this.metadata = data.metadata;
			this.info = this.getInfo(data);
			this.selectUrl = `nodes/${data.metadata.name}`;
		});
		this.$nuxt.$emit("onCreated",'')
	},
	methods: {
		getInfo(data) {
			let regexp = /\B(?=(\d{3})+(?!\d))/g;
			let capacity = `CPU: ${data.status.capacity.cpu}, Memory: ${(this.tranMemory(data.status.capacity.memory)/(1024*1024)).toFixed(2).replace(regexp, ',')+'Mi'}, Pods: ${data.status.capacity.pods}`
			let allocatable = `CPU: ${data.status.allocatable.cpu}, Memory: ${(this.tranMemory(data.status.allocatable.memory)/(1024*1024)).toFixed(2).replace(regexp, ',')+'Mi'}, Pods: ${data.status.allocatable.pods}`
			let addresses = this.getAddress(data.status.addresses);
			let conditions = this.getConditions(data.status.conditions);
			let taints = this.getTaints(data.spec.taints);

			return {
				capacity: capacity,
				allocatable: allocatable,
				addresses: addresses,
				os: `${data.status.nodeInfo.operatingSystem} (${data.status.nodeInfo.architecture})`,
				osImage: data.status.nodeInfo.osImage,
				kernelVersion: data.status.nodeInfo.kernelVersion,
				containerRuntime: data.status.nodeInfo.containerRuntimeVersion,
				kubeletVersion: data.status.nodeInfo.kubeletVersion,
				conditions: conditions,
				taints: taints,
			}
		},
		getTaints(t) {
			let list = [];
			if(t) {
				t.forEach(el => {
					list.push({
						effect: el.effect,
						key: el.key,
						value: el.value
					})
				})
				return list
			}
		},
		getConditions(con) {
			let list = [];
			let style = ''
			con.forEach(el => {
				if(el.status === 'True'){
					if(el.type === 'Ready') {
						style = 'badge-success'
					} else if(el.type === 'Pending') {
						style = 'badge-warning'
					} else if(el.type === 'Failed') {
						style = 'badge-danger'
					} else {
						style = 'badge-secondary'
					}
					list.push({
						type: el.type,
						style: style,
					})
				}
			})
			return list
		},
		getAddress(address) {
			let list =[];
			if(address) {
				address.forEach(el => {
					list.push({
						type: el.type,
						address: el.address,
					})
				})
				return list
			}
		},
		tranMemory(memory) {
			let mem;
			if (memory.includes('Gi')){
				mem = Number(memory.slice(0,-2))*1024*1024*1024
			} else if(memory.includes('Mi')) {
				mem = Number(memory.slice(0,-2))*1024*1024
			} else if(memory.includes('Ki')){
				mem = Number(memory.slice(0,-2))*1024
			} else {
				mem = memory*1024
			}
			return mem
		}
	},
	beforeDestroy(){
		this.$nuxt.$off("onReadCompleted");
	},
}
</script>
