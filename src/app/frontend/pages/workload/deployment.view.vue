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
						<dt class="col-sm-2">Create at</dt><dd class="col-sm-10">{{ this.getTimestampString(metadata.creationTimestamp)}} ago ({{ metadata.creationTimestamp }})</dd>
						<dt class="col-sm-2">Name</dt><dd class="col-sm-10">{{ metadata.name }}</dd>
						<dt class="col-sm-2">Namespace</dt><dd class="col-sm-10">{{ metadata.namespace }}</dd>
						<dt class="col-sm-2">Annotations</dt>
						<dd class="col-sm-10">
							<ul class="list-unstyled mb-0">
								<li v-for="(value, name) in metadata.annotations" v-bind:key="name" class="text-wrap">{{ name }}=<span class="font-weight-light">{{ value }}</span></li>
							</ul>
						</dd>
						<dt class="col-sm-2">Labels</dt>
						<dd class="col-sm-10">
							<span v-for="(value, name) in metadata.labels" v-bind:key="name" class="label">{{ name }}={{ value }}</span>
						</dd>
						<dt class="col-sm-2">UID</dt><dd class="col-sm-10">{{ metadata.uid }}</dd>
						<dt class="col-sm-2">Replicas</dt><dd class="col-sm-10">{{ info.replicas }}</dd>
						<dt v-if="info.selector" class="col-sm-2">Selector</dt><dd v-if="info.selector" class="col-sm-10">
							<span v-for="(value, key) in info.selector" v-bind:key="key" class="border-box background">{{key}}={{value}}</span>
						</dd>
						<dt v-if="info.nodeSelector" class="col-sm-2">Node Selector</dt>
						<dd v-if="info.nodeSelector" class="col-sm-10">
							<span v-for="(value, key) in info.nodeSelector" v-bind:key="key" class="border-box background">{{key}}={{value}}</span>
						</dd>
						<dt class="col-sm-2">Strategy Type</dt><dd class="col-sm-10">{{ info.strategyType }}</dd>
						<dt class="col-sm-2">Conditions</dt>
						<dd class="col-sm-10">
							<span v-for="(d, idx) in info.conditions" v-bind:key="idx" v-bind:class="d.style" class="badge font-weight-light text-sm mb-1 mr-1">{{ d.type }}</span>
						</dd>
						<dt v-if="info.isToleration" class="col-sm-2">Tolerations</dt>
						<dd v-if="info.isToleration" class="col-sm-10">{{ info.tolerations? info.tolerations.length: "-" }}<a class="float-right" v-b-toggle.tol href="#tol-table" @click.prevent @click="isTolerations=!isTolerations">{{isTolerations ? 'Hide' : 'Show'}}</a></dd>
						<b-collapse class="col-sm-12" id="tol-table"><b-table striped hover small :items="info.tolerations"></b-table></b-collapse>

						<dt v-show="info.isAffinity" class="col-sm-2">Affinities</dt>
						<dd v-show="info.isAffinity" class="col-sm-10">{{ info.affinities? Object.keys(info.affinities).length: "-" }}<a class="float-right" v-b-toggle.affi href="#affi-json" @click.prevent @click="isAffinities=!isAffinities">{{isAffinities ? 'Hide' : 'Show'}}</a>
							<b-collapse id="affi-json"><c-jsontree id="txtSpec" v-model="info.affinities" class="card-body p-2 border"></c-jsontree></b-collapse>
						</dd>
					</dl>
				</div>
			</div>
		</div>
	</div>

	<!-- 3. pods -->
	<c-podlist class="row" v-model="selectUrl" @navigate="$emit('navigate',arguments[0])"></c-podlist>
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
			this.selectUrl = `namespaces/${data.metadata.namespace}/deployments/${data.metadata.name}`;
			this.info = this.getInfo(data);
		});
		this.$nuxt.$emit("onCreated",'')
	},
	methods: {
		getInfo(data) {
			let replicas = data.spec.replicas +' desired, ' + (data.status.updatedReplicas || 0) + ' updated, ' + (data.status.replicas || 0) +' total, ' + (data.status.availableReplicas || 0) + ' available, ' + (data.status.unavailableReplicas || 0) + ' unavailable'
			let conditions = [];
			let tolerations = [];
			let affinity = [];
			let bToleration = false;
			let bAffinity = false;
			if (data.status.conditions) {
				data.status.conditions.forEach(el => {
					conditions.push({
						type: el.type,
						style: this.checkStyle(el.type),
					})
				})
				conditions.sort(function(a,b) {
					return a.type < b.type ? -1 : a.type > b.type ? 1 : 0;
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
				replicas: replicas,
				selector: data.spec.selector.matchLabels || '',
				nodeSelector: data.spec.template.spec.nodeSelector || '',
				strategyType: data.spec.strategy.type,
				conditions: conditions,
				tolerations: tolerations,
				affinities: Object.assign({},affinity),
				isAffinity: bAffinity,
				isToleration: bToleration,
			}
		},
		checkStyle(t) {
			if(t === 'Progressing') return 'badge-primary'
			if(t === 'Available') return 'badge-success'
			else return 'badge-danger'
		}
	},
	beforeDestroy(){
		this.$nuxt.$off("onReadCompleted");
	},
}
</script>
