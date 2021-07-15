<template>
<div>
	<!-- 1. charts -->
	<c-charts class="row" v-model="chartsUrl"></c-charts>
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
						<dd v-if="info.isToleration" class="col-sm-10">{{ info.tolerations? info.tolerations.length: "-" }}<a class="float-right" v-b-toggle.tol href="#tol-table" @click.prevent @click="onTol">{{onTols ? 'Hide' : 'Show'}}</a></dd>
						<b-collapse class="col-sm-12" id="tol-table"><b-table striped hover small :items="info.tolerations"></b-table></b-collapse>

						<dt v-show="info.isAffinity" class="col-sm-2">Affinities</dt>
						<dd v-show="info.isAffinity" class="col-sm-10">{{ info.affinities? Object.keys(info.affinities).length: "-" }}<a class="float-right" v-b-toggle.affi href="#affi-json" @click.prevent @click="onAffi">{{onAffis ? 'Hide' : 'Show'}}</a>
							<b-collapse id="affi-json"><c-jsontree id="txtSpec" v-model="info.affinities" class="card-body p-2 border"></c-jsontree></b-collapse>
						</dd>
					</dl>
				</div>
			</div>
		</div>
	</div>

	<!-- 3. pods -->
	<div class="row">
		<div class="col-md-12">
			<div class="card card-secondary card-outline">
				<div class="card-header p-2"><h3 class="card-title">Pods</h3></div>
				<div class="card-body p-2">
					<b-table striped hover small :items="childPod" :fields="fields">
						<template v-slot:cell(name)="data">
							<a href="#" @click="$emit('navigate', getViewLink('', 'pods', data.item.namespace, data.item.name))">{{ data.item.name }}</a>
						</template>
						<template v-slot:cell(status)="data">
							<span v-bind:class="data.item.status.style">{{ data.item.status.value }}</span>
						</template>
						<template v-slot:cell(nowCpu)="data">
							<span v-if="data.item.nowCpu[data.item.idx]">{{ data.item.nowCpu[data.item.idx].val ? data.item.nowCpu[data.item.idx].val : '' }}</span>
						</template>
						<template v-slot:cell(nowMemory)="data">
							<span v-if="data.item.nowMemory[data.item.idx]">{{ data.item.nowMemory[data.item.idx].val ? data.item.nowMemory[data.item.idx].val+'Mi' : '' }}</span>
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
import VueJsonTree		from "@/components/jsontree";
import VueEventsView	from "@/components/view/eventsView.vue";
import VueChartsView	from "@/components/view/metricsChartsView.vue";

export default {
	components: {
		"c-jsontree": { extends: VueJsonTree },
		"c-events": { extends: VueEventsView },
		"c-charts": { extends: VueChartsView }
	},
	data() {
		return {
			metadata: {},
			chartsUrl: "",
			info: [],
			childPod: [],
			onTols: false,
			onAffis: false,
			fields: [
				{ key: "name", label: "Name", sortable: true },
				{ key: "namespace", label: "Namespace", sortable: true  },
				{ key: "ready", label: "Ready", sortable: true  },
				{ key: "nowCpu", label: "CPU" },
				{ key: "nowMemory", label: "Memory" },
				{ key: "status", label: "Status" },
			],
		}
	},
	mounted() {
		this.$nuxt.$on("onReadCompleted", (data) => {
			if(!data) return
			this.metadata = data.metadata;
			this.chartsUrl = `namespaces/${data.metadata.namespace}/deployments/${data.metadata.name}/metrics`;
			this.onSync(data)
		});
		this.$nuxt.$emit("onCreated",'')
	},
	methods: {
		onSync(data) {
			this.info = this.getInfo(data);
			this.childPod = this.getChildPod(data.spec.template.metadata.labels);
		},
		getInfo(data) {
			let replicas = data.spec.replicas +' desired, ' + (data.status.updatedReplicas || 0) + ' updated, ' + (data.status.replicas || 0) +' total, ' + (data.status.availableReplicas || 0) + ' available, ' + (data.status.unavailableReplicas || 0) + ' unavailable'
			let conditions = [];
			let tolerations = [];
			let affinity = [];
			let isToleration = false;
			let isAffinity = false;
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
					isToleration = true;
				})
			}
			if(data.spec.template.spec.affinity && Object.keys(data.spec.template.spec.affinity).length !== 0) {
				affinity = data.spec.template.spec.affinity;
				isAffinity = true;
			}
			return {
				replicas: replicas,
				selector: data.spec.selector.matchLabels || '',
				nodeSelector: data.spec.template.spec.nodeSelector || '',
				strategyType: data.spec.strategy.type,
				conditions: conditions,
				tolerations: tolerations,
				affinities: affinity,
				isAffinity: isAffinity,
				isToleration: isToleration,
			}
		},
		getChildPod(label) {
			label = this.stringifyLabels(label)
			let childPod = [];
			this.$axios.get(this.getApiUrl('','pods',this.metadata.namespace,'','labelSelector=' + label))
					.then( resp => {
						let idx = 0;
						resp.data.items.forEach(el =>{
							if(el.metadata.ownerReferences) {
								childPod.push({
									name: el.metadata.name,
									namespace: el.metadata.namespace,
									ready: this.toReady(el.status,el.spec),
									nowMemory: 0,
									nowCpu: 0,
									//nowMemory: this.onMemory(el,idx),
									//nowCpu: this.onCpu(el,idx),
									status: this.toStatus(el.metadata.deletionTimestamp, el.status),
									idx: idx,
								})
								idx++;
							}
						})
					})
			return childPod
		},
		checkStyle(t) {
			if(t === 'Progressing') return 'badge-primary'
			if(t === 'Available') return 'badge-success'
			else return 'badge-danger'
		},
		onTol() {
			this.onTols = !this.onTols
		},
		onAffi() {
			this.info.affinities = Object.assign({},this.info.affinities)
			this.onAffis = !this.onAffis
		},
	},
	beforeDestroy(){
		this.$nuxt.$off("onReadCompleted");
	},
}
</script>
