<template>
<div>
	<!-- 1. charts -->
	<c-charts class="row" v-model="selectUrl"></c-charts>
	<!-- 2. metadata -->
	<c-metadata v-model="metadata" dtCols="3" ddCols="9">
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
		<dt v-if="info.taints.length>0" class="col-sm-3">Taints</dt>
		<dd v-if="info.taints.length>0" class="col-sm-9">
			<span v-for="(val, idx) in info.taints" v-bind:key="idx" v-bind:class="val.style">{{ val.key }}: {{ val.effect }} ({{ val.value }})</span>
		</dd>
		<dt class="col-sm-3">Conditions</dt>
		<dd class="col-sm-9">
			<span v-for="(value, idx) in info.conditions" v-bind:key="idx" v-bind:class="{'badge-success':value.type=='Ready', 'badge-danger':value.type=='Failed','badge-warning':value.type=='Pending'}" class="badge font-weight-light text-sm mb-1 mr-1">{{ value.type }}</span>
		</dd>
	</c-metadata>
	<!-- 3. pods -->
	<c-podlist class="row" v-model="selectUrl" namespace="true" @navigate="$emit('navigate',arguments[0])"></c-podlist>
	<!-- 4. events -->
	<c-events class="row" v-model="metadata.uid"></c-events>

</div>
</template>
<script>
import VueMetadataView	from "@/components/view/metadataView.vue";
import VueJsonTree		from "@/components/jsontree";
import VueEventsView	from "@/components/view/eventsView.vue";
import VueChartsView	from "@/components/view/metricsChartsView.vue";
import VuePodListView	from "@/components/view/podListView.vue";

export default {
	components: {
		"c-metadata": { extends: VueMetadataView },
		"c-jsontree": { extends: VueJsonTree },
		"c-events": { extends: VueEventsView },
		"c-charts": { extends: VueChartsView },
		"c-podlist": { extends: VuePodListView }
	},
	data() {
		return {
			metadata: {},
			info: {
				taints: []
			},
			selectUrl: ""
		}
	},
	mounted() {
		this.$nuxt.$on("onReadCompleted", (data) => {
			if(!data) return
			this.metadata = data.metadata;
			this.selectUrl = `nodes/${data.metadata.name}`;

			let regexp = /\B(?=(\d{3})+(?!\d))/g;
			this.info = {
				capacity: `CPU: ${data.status.capacity.cpu}, Memory: ${(this.tranMemory(data.status.capacity.memory)/(1024**2)).toFixed(2).replace(regexp, ',')+'Mi'}, Pods: ${data.status.capacity.pods}`,
				allocatable: `CPU: ${data.status.allocatable.cpu}, Memory: ${(this.tranMemory(data.status.allocatable.memory)/(1024**2)).toFixed(2).replace(regexp, ',')+'Mi'}, Pods: ${data.status.allocatable.pods}`,
				addresses: data.status.addresses,
				os: `${data.status.nodeInfo.operatingSystem} (${data.status.nodeInfo.architecture})`,
				osImage: data.status.nodeInfo.osImage,
				kernelVersion: data.status.nodeInfo.kernelVersion,
				containerRuntime: data.status.nodeInfo.containerRuntimeVersion,
				kubeletVersion: data.status.nodeInfo.kubeletVersion,
				conditions: data.status.conditions? data.status.conditions.filter(el=> {return el.status=='True'}): [],
				taints: data.spec.taints?data.spec.taints:[],
			};
		});
		this.$nuxt.$emit("onCreated",'')
	},
	methods: {
		tranMemory(memory) {
			let mem;
			if (memory.includes('Gi')){
				mem = Number(memory.slice(0,-2))*1024**3
			} else if(memory.includes('Mi')) {
				mem = Number(memory.slice(0,-2))*1024**2
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
