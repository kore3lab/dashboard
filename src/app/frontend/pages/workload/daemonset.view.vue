<template>
<div>
	<!-- 1. chart -->
	<c-charts class="row" v-model="selectUrl"></c-charts>
	<!-- 2. metadata -->
	<c-metadata v-model="metadata" dtCols="2" ddCols="10">
		<dt v-if="info.selector" class="col-sm-2">Selector</dt>
		<dd v-if="info.selector" class="col-sm-10">
			<span v-for="(value, key) in info.selector" v-bind:key="key" class="border-box background">{{key}}={{value}}</span>
		</dd>
		<dt v-if="info.nodeSelector" class="col-sm-2">Node Selector</dt>
		<dd v-if="info.nodeSelector" class="col-sm-10">
			<span v-for="(value, key) in info.nodeSelector" v-bind:key="key" class="border-box background">{{key}}={{value}}</span>
		</dd>
		<dt class="col-sm-2">Images</dt>
		<dd class="col-sm-10">
			<ul class="list-unstyled">
				<li v-for="(val, idx) in info.image" v-bind:key="idx">{{ val }}</li>
			</ul>
		</dd>
		<dt class="col-sm-2">Strategy Type</dt><dd class="col-sm-10">{{ info.strategyType }}</dd>

		<dt v-if="info.isToleration" class="col-sm-2">Tolerations</dt>
		<dd v-if="info.isToleration" class="col-sm-10">{{ info.tolerations? info.tolerations.length: "-" }}
			<a href="#" class="float-right " @click="isTolerations=!isTolerations">{{isTolerations?'Hide':'Show'}}</a>
			<b-collapse v-model="isTolerations">
				<b-table-lite small :items="info.tolerations"></b-table-lite>
			</b-collapse>
		</dd>
		<dt v-show="info.isAffinity" class="col-sm-2">Affinities</dt>
		<dd v-show="info.isAffinity" class="col-sm-10">{{ info.affinities? Object.keys(info.affinities).length: "-" }}
			<a href="#" class="float-right " @click="isAffinities=!isAffinities">{{isAffinities?'Hide':'Show'}}</a>
			<b-collapse v-model="isAffinities">
				<c-jsontree v-model="info.affinities" class="card-body p-2 border"></c-jsontree>
			</b-collapse>
		</dd>
	</c-metadata>
	<!-- 3. pods -->
	<c-podlist class="row" v-model="selectUrl" @navigate="$emit('navigate',arguments[0])"></c-podlist>
	<!-- 4. events -->
	<c-events class="row" v-model="metadata.uid"></c-events>
</div>
</template>
<script>
import VueMetadataView	from "@/components/view/metadataView.vue";
import VueJsonTree		from "@/components/jsontree";
import VueEventsView	from "@/components/view/eventsView.vue";;
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
			selectUrl: "",
			info: [],
			isTolerations: false,
			isAffinities: false
		}
	},
	mounted() {
		this.$nuxt.$on("onReadCompleted", (data) => {
			if(!data) return
			this.metadata = data.metadata;
			this.selectUrl = `namespaces/${data.metadata.namespace}/daemonsets/${data.metadata.name}`;
			this.info = this.getInfo(data);
		});
		this.$nuxt.$emit("onCreated",'')
	},
	methods: {
		getInfo(data) {
			let tolerations = [];
			let affinity = [];
			let image = [];
			let bToleration = false;
			let bAffinity = false;

			if(data.spec.template.spec.containers) {
				data.spec.template.spec.containers.forEach(el =>{
					image.push(el.image)
				})
			}
			if(data.spec.template.spec.tolerations) {
				data.spec.template.spec.tolerations.forEach(el =>{
					tolerations.push({
						key: el.key || '',
						operator: el.operator || '',
						effect: el.effect || '',
						seconds: el.tolerationSeconds || '',
					})
					bToleration = true;
				})
			}
			if(data.spec.template.spec.affinity && Object.keys(data.spec.template.spec.affinity).length !== 0) {
				affinity = data.spec.template.spec.affinity;
				bAffinity = true;
			}
			return {
				selector: data.spec.selector.matchLabels || '',
				nodeSelector: data.spec.template.spec.nodeSelector || '',
				strategyType: data.spec.updateStrategy.type,
				image: image,
				tolerations: tolerations,
				affinities: Object.assign({},affinity),
				isAffinity: bAffinity,
				isToleration: bToleration,
			}
		}
	},
	beforeDestroy(){
		this.$nuxt.$off("onReadCompleted");
	},
}
</script>
<style>
    .ellipsis{
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
    } 
</style>