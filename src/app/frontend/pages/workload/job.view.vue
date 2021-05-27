<template>
	<div>
		<div class="row">
			<div class="col-md-12">
				<div class="card card-secondary card-outline">
					<div class="card-body p-2">
						<dl class="row mb-0">
							<dt class="col-sm-2 text-truncate">Create at</dt><dd class="col-sm-10">{{ this.getTimestampString(metadata.creationTimestamp)}} ago ({{ metadata.creationTimestamp }})</dd>
							<dt class="col-sm-2">Name</dt><dd class="col-sm-10">{{ metadata.name }}</dd>
							<dt class="col-sm-2">Namespace</dt><dd class="col-sm-10">{{ metadata.namespace }}</dd>
							<dt class="col-sm-2">Annotations</dt>
							<dd class="col-sm-10 text-truncate">
								<ul class="list-unstyled mb-0">
									<li v-for="(value, name) in metadata.annotations" v-bind:key="name"><span class="badge badge-secondary font-weight-light text-sm mb-1">{{ name }}={{ value }}</span></li>
								</ul>
							</dd>
							<dt class="col-sm-2">Labels</dt>
							<dd class="col-sm-10 text-truncate">
								<ul class="list-unstyled mb-0">
									<li v-for="(value, name) in metadata.labels" v-bind:key="name"><span class="badge badge-secondary font-weight-light text-sm mb-1">{{ name }}={{ value }}</span></li>
								</ul>
							</dd>
							<dt class="col-sm-2">UID</dt><dd class="col-sm-10">{{ metadata.uid }}</dd>
							<dt v-if="metadata.ownerReferences" class="col-sm-2 text-truncate">Controlled By</dt>
							<dd v-if="metadata.ownerReferences" class="col-sm-10">{{ metadata.ownerReferences[0].kind }} <a href="#" @click="$emit('navigate', getViewLink(controller.g, controller.k, metadata.namespace, metadata.ownerReferences[0].name))">{{ metadata.ownerReferences[0].name }}</a></dd>
							<dt v-if="info.selector" class="col-sm-2">Selector</dt><dd v-if="info.selector" class="col-sm-10"><span v-for="(value, key) in info.selector" v-bind:key="key" class="badge badge-secondary font-weight-light text-sm mb-1 mr-1">{{ value }}</span></dd>
							<dt v-if="info.nodeSelector" class="col-sm-2">Node Selector</dt><dd v-if="info.nodeSelector" class="col-sm-10"><span v-for="(value, key) in info.nodeSelector" v-bind:key="key" class="badge badge-secondary font-weight-light text-sm mb-1 mr-1">{{key}}={{value}}</span></dd>
							<dt class="col-sm-2">Image</dt><dd class="col-sm-10">{{ info.image }}</dd>
							<dt class="col-sm-2">Conditions</dt><dd class="col-sm-10"><span v-for="(val, idx) in info.condition" v-bind:key="idx" v-bind:class="val.style" class="badge font-weight-light text-sm mb-1 mr-1"> {{ val.type }} </span></dd>
							<dt class="col-sm-2">Completions</dt><dd class="col-sm-10">{{ info.completions }}</dd>
							<dt class="col-sm-2">Parallelism</dt><dd class="col-sm-10">{{ info.parallelism }}</dd>
							<dt class="col-sm-2">Pod Status</dt><dd class="col-sm-10"><span v-for="(val,idx) in cs" v-bind:key="idx" v-bind:class="val.style">{{ val.status }} : {{ val.count }}  </span><span v-if="!isStatus">-</span></dd>
						</dl>
					</div>
				</div>
			</div>
		</div>

		<div v-show="isPods" class="row">
			<div class="col-md-12">
				<div class="card card-secondary card-outline">
					<div class="card-header p-2"><h3 class="card-title text-md">Pods</h3></div>
					<div class="card-body p-2 overflow-auto">
						<b-table striped hover small :items="childPod" :fields="fields" class="text-truncate">
							<template v-slot:cell(name)="data">
								<a href="#" @click="$emit('navigate', getViewLink('', 'pods', data.item.namespace, data.item.name))">{{ data.item.name }}</a>
							</template>
							<template v-slot:cell(status)="data">
								<span v-bind:class="data.item.status.style">{{ data.item.status.value }}</span>
							</template>
							<template v-slot:cell(nowCpu)="data">
								<span v-if="data.item.nowCpu[data.item.idx]">{{ data.item.nowCpu[data.item.idx].val ? data.item.nowCpu[data.item.idx].val : '-' }}</span>
							</template>
							<template v-slot:cell(nowMemory)="data">
								<span v-if="data.item.nowMemory[data.item.idx]">{{ data.item.nowMemory[data.item.idx].val ? data.item.nowMemory[data.item.idx].val+'Mi' : '-'}}</span>
							</template>
						</b-table>
					</div>
				</div>
			</div>
		</div>

		<div class="row">
			<div class="col-md-12">
				<div class="card card-secondary card-outline m-0">
					<div class="card-header p-2"><h3 class="card-title text-md">Events</h3></div>
					<div class="card-body p-2">
						<dl v-for="(val, idx) in event" v-bind:key="idx" class="row mb-0 card-body p-2 border-bottom">
							<dt class="col-sm-12"><p v-bind:class="val.type" class="mb-1">{{ val.name }}</p></dt>
							<dt class="col-sm-2 text-truncate">Source</dt><dd class="col-sm-10">{{ val.source }}</dd>
							<dt class="col-sm-2 text-truncate">Count</dt><dd class="col-sm-10">{{ val.count }}</dd>
							<dt class="col-sm-2 text-truncate">Sub-object</dt><dd class="col-sm-10">{{ val.subObject }}</dd>
							<dt class="col-sm-2 text-truncate">Last seen</dt><dd class="col-sm-10">{{ val.lastSeen }}</dd>
						</dl>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>
<script>

export default {
	data() {
		return {
			metadata: {},
			info: [],
			event: [],
			childPod: [],
			controller: [],
			temp: [],
			cs: [],
			nowCpu: [],
			nowMemory: [],
			isStatus: false,
			isPods: false,
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
			this.onSync(data)
		});
		this.$nuxt.$emit("onCreated",'')
	},
	methods: {
		onSync(data) {
			this.controller = this.getController(data.metadata.ownerReferences)
			this.info = this.getInfo(data);
			this.event = this.getEvents(data.metadata.uid,'fieldSelector=involvedObject.name='+data.metadata.name);
			this.childPod = this.getChildPod(data.metadata.uid);
		},
		getInfo(data) {

			return {
				selector: this.stringifyLabels(data.spec.selector? data.spec.selector.matchLabels : ''),
				image: data.spec.template.spec.containers[0].image,
				nodeSelector: data.spec.template.spec.nodeSelector || '',
				condition: this.getCondition(data.status.conditions),
				completions: data.spec.completions,
				parallelism: data.spec.parallelism,

			}
		},
		getCondition(con) {
			let list = []
			if(con) {
				con.forEach(el => {
					if (el.status === 'True') {
						if (el.type === 'Complete') list.push({ type: 'Complete', style: 'badge-success' })
						else list.push({ type: 'Failed', style: 'badge-danger' })
					}
				})
				return list
			}
		},
		getChildPod(uid) {
			let childPod = [];
			this.nowMemory = [];
			this.nowCpu = [];
			this.temp = [];
			this.cs = [];
			this.isStatus = false;
			this.isPods = false;
			this.$axios.get(this.getApiUrl('','pods',this.metadata.namespace,'','labelSelector=controller-uid=' + uid))
					.then( resp => {
						let idx = 0;
						resp.data.items.forEach(el =>{
							this.isPods = true;
							childPod.push({
								name: el.metadata.name,
								namespace: el.metadata.namespace,
								ready: this.toReady(el.status,el.spec),
								nowMemory: this.onMemory(el,idx),
								nowCpu: this.onCpu(el,idx),
								status: this.toStatus(el.metadata.deletionTimestamp, el.status),
								countStatus: this.countStatus(el.status),
								idx: idx,
							})
							idx++;
						})
					})
			return childPod
		},

		onCpu(el,idx) {
			this.$axios.get(`/api/clusters/${this.currentContext()}/namespaces/${el.metadata.namespace}/pods/${el.metadata.name}/metrics/cpu`)
					.then(resp => {
						if (resp.data.items) {
							let data = resp.data.items[0];
							let da= [];
							data.metricPoints.forEach(d => {
								da.push(d.value/1000);
							});
							this.nowCpu.push({
								val:((da[da.length-1])).toFixed(3),
								idx:idx,
							})
						} else {
							this.nowCpu.push({
								val: '',
								idx: idx,
							})
						}
						this.nowCpu = this.sorted(this.nowCpu)
					})
			return this.nowCpu
		},
		onMemory(el,idx) {
			this.$axios.get(`/api/clusters/${this.currentContext()}/namespaces/${el.metadata.namespace}/pods/${el.metadata.name}/metrics/memory`)
					.then(resp => {
						if (resp.data.items) {
							let data = resp.data.items[0];
							let da= [];
							data.metricPoints.forEach(d => {
								da.push(Math.round(d.value/1024));
							});
							this.nowMemory.push({
								val:((da[da.length-1])/1024).toFixed(2),
								idx:idx,
							})
						} else {
							this.nowMemory.push({
								val: '',
								idx: idx,
							})
						}
						this.nowMemory = this.sorted(this.nowMemory)
					})
			return this.nowMemory
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
				if(el === 'Running' || el === 'Succeeded') {
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
	},
	beforeDestroy(){
		this.$nuxt.$off("onReadCompleted");
	},
}
</script>
