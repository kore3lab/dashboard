<template>
<div>
	<!-- 1. metadata -->
	<c-metadata v-model="metadata" dtCols="2" ddCols="10">
		<dt v-if="metadata.ownerReferences" class="col-sm-2">Controlled By</dt>
		<dd v-if="metadata.ownerReferences" class="col-sm-10">{{ controller.kind }} <a href="#" @click="$emit('navigate', getViewLink(controller.group, controller.resource, metadata.namespace, controller.name))">{{ controller.name }}</a></dd>
		<dt v-if="info.selector" class="col-sm-2">Selector</dt>
		<dd v-if="info.selector" class="col-sm-10">
			<span v-for="(value, key) in info.selector" v-bind:key="key" class="border-box background">{{ value }}</span>
		</dd>
		<dt v-if="info.nodeSelector" class="col-sm-2">Node Selector</dt>
		<dd v-if="info.nodeSelector" class="col-sm-10">
			<span v-for="(value, key) in info.nodeSelector" v-bind:key="key" class="border-box background">{{key}}={{value}}</span>
		</dd>
		<dt class="col-sm-2">Image</dt><dd class="col-sm-10">{{ info.image }}</dd>
		<dt class="col-sm-2">Conditions</dt>
		<dd class="col-sm-10">
			<span v-for="(value, idx) in info.conditions" v-bind:key="idx" v-bind:class="{'badge-success':value.type=='Complete', 'badge-danger':value.type=='Failed'}" class="badge font-weight-light text-sm mb-1 mr-1"> {{ value.type }} </span>
		</dd>
		<dt class="col-sm-2">Completions</dt><dd class="col-sm-10">{{ info.completions }}</dd>
		<dt class="col-sm-2">Parallelism</dt><dd class="col-sm-10">{{ info.parallelism }}</dd>
	</c-metadata>
	<!-- 2. pods -->
	<c-podlist class="row" v-model="selectUrl" @navigate="$emit('navigate',arguments[0])"></c-podlist>
	<!-- 3. events -->
	<c-events class="row" v-model="metadata.uid"></c-events>
</div>
</template>
<script>
import VueMetadataView	from "@/components/view/metadataView.vue";
import VueEventsView	from "@/components/view/eventsView.vue";
import VuePodListView	from "@/components/view/podListView.vue";

export default {
	components: {
		"c-metadata": { extends: VueMetadataView },
		"c-events": { extends: VueEventsView },
		"c-podlist": { extends: VuePodListView }
	},
	data() {
		return {
			metadata: {},
			selectUrl: "",
			info: {},
			controller: {}
		}
	},
	mounted() {
		this.$nuxt.$on("onReadCompleted", (data) => {
			if(!data) return
			this.metadata = data.metadata;
			this.selectUrl = `namespaces/${data.metadata.namespace}/jobs/${data.metadata.name}`;
			this.info = {
				selector: this.stringifyLabels(data.spec.selector? data.spec.selector.matchLabels : ''),
				image: data.spec.template.spec.containers[0].image,
				nodeSelector: data.spec.template.spec.nodeSelector || '',
				conditions: data.status.conditions?data.status.conditions.filter(el=>{return el.status === 'True'}):[],
				completions: data.spec.completions,
				parallelism: data.spec.parallelism,
			};
			this.controller = data.metadata.ownerReferences?this.getResource(data.metadata.ownerReferences[0]):{};
		});
		this.$nuxt.$emit("onCreated",'')
	},
	methods: {
	},
	beforeDestroy(){
		this.$nuxt.$off("onReadCompleted");
	},
}
</script>
