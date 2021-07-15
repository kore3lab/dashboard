<template>
<div>
	<!-- 1. graph -->
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
								<li v-for="(value, name) in metadata.annotations" v-bind:key="name">{{ name }}=<span class="font-weight-light">{{ value }}</span></li>
							</ul>
						</dd>
						<dt class="col-sm-2">Labels</dt>
						<dd class="col-sm-10">
							<span v-for="(value, name) in metadata.labels" v-bind:key="name" class="label">{{ name }}={{ value }}</span>
						</dd>
						<dt class="col-sm-2">UID</dt><dd class="col-sm-10">{{ metadata.uid }}</dd>
						<dt v-if="metadata.ownerReferences" class="col-sm-2">Controlled By</dt>
						<dd v-if="metadata.ownerReferences" class="col-sm-10">{{ metadata.ownerReferences[0].kind }} <a href="#" @click="$emit('navigate', getViewLink(controller.g, controller.k, metadata.namespace, metadata.ownerReferences[0].name))">{{ metadata.ownerReferences[0].name }}</a></dd>
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
						<dt v-if="info.isToleration" class="col-sm-2">Tolerations</dt>
						<dd v-if="info.isToleration" class="col-sm-10">{{ info.tolerations? info.tolerations.length: "-" }}<a class="float-right" v-b-toggle.tol href="#tol-table" @click.prevent @click="onTol">{{onTols ? 'Hide' : 'Show'}}</a></dd>
						<b-collapse class="col-sm-12" id="tol-table"><b-table striped hover small :items="info.tolerations"></b-table></b-collapse>

						<dt v-show="info.isAffinity" class="col-sm-2">Affinities</dt>
						<dd v-show="info.isAffinity" class="col-sm-10">{{ info.affinities? Object.keys(info.affinities).length: "-" }}<a class="float-right" v-b-toggle.affi href="#affi-json" @click.prevent @click="onAffi">{{onAffis ? 'Hide' : 'Show'}}</a>
							<b-collapse id="affi-json"><c-jsontree id="txtSpec" v-model="info.affinities" class="card-body p-2 border"></c-jsontree></b-collapse>
						</dd>
						<dt class="col-sm-2">Pod Status</dt><dd class="col-sm-10"><span v-for="(val,idx) in cs" v-bind:key="idx" v-bind:class="val.style">{{ val.status }} : {{ val.count }}  </span><span v-if="!isStatus">-</span></dd>
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
							<span v-if="data.item.nowMemory[data.item.idx]">{{ data.item.nowMemory[data.item.idx].val ? data.item.nowMemory[data.item.idx].val+'Mi' : ''}}</span>
						</template>
					</b-table>
				</div>
			</div>
		</div>
	</div>

	<!-- 3. events -->
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
			controller: [],
			temp: [],
			cs: [],
			sumCpu: [],
			sumMemory: [],
			nowCpu: [],
			topCpu: [],
			topMemory: [],
			nowMemory: {},
			totalCpu: 0,
			totalMemory: 0,
			cpuRequests: 0,
			cpuLimits: 0,
			memoryRequests: 0,
			memoryLimits: 0,
			isCpu: false,
			isMemory: false,
			isStatus: false,
			onTols: false,
			onAffis: false,
			isPods: false,
			fields: [
				{ key: "name", label: "Name", sortable: true },
				{ key: "namespace", label: "Namespace", sortable: true  },
				{ key: "ready", label: "Ready", sortable: true  },
				{ key: "nowCpu", label: "CPU" },
				{ key: "nowMemory", label: "Memory" },
				{ key: "status", label: "Status" },
			]
		}
	},
	mounted() {
		this.$nuxt.$on("onReadCompleted", (data) => {
			if(!data) return
			this.metadata = data.metadata;
			this.chartsUrl = `namespaces/${data.metadata.namespace}/statefulsets/${data.metadata.name}/metrics`;
			this.onSync(data)
		});
		this.$nuxt.$emit("onCreated",'')
	},
	methods: {
		onSync(data) {
			this.totalCpu = 0; this.totalMemory = 0; this.cpus = {};
			this.controller = this.getController(data.metadata.ownerReferences)
			this.info = this.getInfo(data);
			this.childPod = this.getChildPod(data.status.currentRevision);
		},
		getInfo(data) {
			let tolerations = [];
			let affinity = [];
			let image = [];
			let isToleration = false;
			let isAffinity = false;

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
					isToleration = true;
				})
			}
			if(data.spec.template.spec.affinity && Object.keys(data.spec.template.spec.affinity).length !== 0) {
				affinity = data.spec.template.spec.affinity;
				isAffinity = true;
			}
			return {
				selector: data.spec.selector.matchLabels || '',
				nodeSelector: data.spec.template.spec.nodeSelector || '',
				image: image,
				tolerations: tolerations,
				affinities: affinity,
				isAffinity: isAffinity,
				isToleration: isToleration,
			}
		},
		getChildPod(currentRevision) {
			let childPod = [];
			this.nowMemory = [];
			this.nowCpu = [];
			this.sumCpu = [];
			this.sumMemory = [];
			this.temp = [];
			this.cs = [];
			this.topCpu = [];
			this.topMemory = [];
			this.cpuRequests = 0
			this.cpuLimits = 0
			this.memoryRequests = 0
			this.memoryLimits = 0
			this.isStatus = false;
			this.isPods = false;
			this.isCpu = false;
			this.isMemory = false;
			this.$axios.get(this.getApiUrl('','pods',this.metadata.namespace,'','labelSelector=controller-revision-hash=' + currentRevision))
					.then( resp => {
						let idx = 0;
						resp.data.items.forEach(el =>{
							this.isPods = true;
							childPod.push({
								name: el.metadata.name,
								namespace: el.metadata.namespace,
								ready: this.toReady(el.status,el.spec),
								nowMemory: 0,
								nowCpu: 0,
								//nowMemory: this.onMemory(el,idx),
								//nowCpu: this.onCpu(el,idx),
								status: this.toStatus(el.metadata.deletionTimestamp, el.status),
								countStatus: this.countStatus(el.status),
								idx: idx,
							})
							idx++;
						})
					})
			return childPod
		},
		countStatus(status) {
			this.isStatus = true;
			let style;
			this.cs = [];
			if(status.phase in this.temp) {
				this.temp[status.phase] += 1;
			} else {
				this.temp[status.phase] = 1;
			}
			let key = Object.keys(this.temp)
			key.forEach(el =>{
				if(el === 'Running') {
					style = 'text-success'
				} else if (el === 'Pending') {
					style = 'text-warning'
				} else if (el === 'Failed') {
					style = 'text-danger'
				} else style = 'text-secondary'
				this.cs.push({
					status: el,
					count: this.temp[el],
					style: style,
				})
			})
			return this.cs
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
