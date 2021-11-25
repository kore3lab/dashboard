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
		<dt class="col-sm-3">Access Modes</dt><dd class="col-sm-9"><span v-for="(val, idx) in info.accessModes" v-bind:key="idx">{{ val }} </span></dd>
		<dt class="col-sm-3">Storage Class Name</dt><dd class="col-sm-9">{{ info.storageClassName }}</dd>
		<dt class="col-sm-3">Storage</dt><dd class="col-sm-9">{{ info.storage }}</dd>
		<dt v-if="info.volumeName" class="col-sm-3">Binding Volume</dt><dd v-if="info.volumeName" class="col-sm-9"><a href="#" @click="$emit('navigate', getViewLink('', 'persistentvolumes','', info.volumeName ))">{{ info.volumeName }}</a></dd>
		<dt class="col-sm-3">Pods</dt>
		<dd class="col-sm-9">
			<span v-for="(val, idx) in pods" v-bind:key="idx"><a href="#" @click="$emit('navigate', getViewLink('', 'pods', val.namespace, val.name ))">{{ val.name }}</a></span>
			<span v-if="pods.length==0">-</span>
		</dd>
		<dt class="col-sm-3">Status</dt><dd class="col-sm-9">{{ info.status }}</dd>
	</c-metadata>
	<!-- 2. selector -->
	<div class="row">
		<div class="col-md-12">
			<div class="card card-secondary card-outline">
				<div class="card-header p-2"><h3 class="card-title">Selector</h3></div>
				<div class="card-body p-2">
					<dl class="row mb-0">
						<dt class="col-sm-3">Match Labels</dt>
						<dd class="col-sm-9">
							<span v-for="(val, idx) in selector.label" v-bind:key="idx" class="border-box background">{{ val }} </span>
							<span v-if="!selector.label.length==0"> - </span>
						</dd>
						<dt class="col sm-3">Match Expressions</dt>
						<dd class="col-sm-9">
							<ul v-for="(val, idx) in selector.expression" v-bind:key="idx" class="row list-unstyled mb-0">
								<li class="col-sm-3 mb-1 text-bold">Key</li>
								<li class="col-sm-9">{{ val.key }}</li>
								<li class="col-sm-3 mb-1 text-bold">Operator</li>
								<li class="col-sm-9">{{ val.operator }}</li>
								<li class="col-sm-3 mb-1 text-bold">Values</li>
								<li class="col-sm-9">{{ val.values }}</li>
							</ul>
							<span v-if="selector.expression.length==0"> - </span>
						</dd>
					</dl>
				</div>
			</div>
		</div>
	</div>
	<!-- 3. events -->
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
			info: {},
			pods: [],
			selector: { label:[], expression:[] }
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
				accessModes: data.spec.accessModes || '-',
				storageClassName: data.spec.storageClassName || '-',
				storage: data.spec.resources && data.spec.resources.requests ? data.spec.resources.requests.storage: "-",
				status: data.status? data.status.phase : '-',
				volumeName: data.spec.volumeName,
			};
			this.selector = data.spec.selector?{label: this.stringifyLabels(data.spec.selector.matchLabels), expression: data.spec.selector.matchExpressions}: {label: [], expression: []};

			// pod-list
			this.$axios.get(this.getApiUrl('','pods',this.metadata.namespace))
				.then( resp => {
					let podList = [];
					resp.data.items.forEach(el => {
						el.spec.volumes.forEach(e => {
							if(e.persistentVolumeClaim) {
								if(e.persistentVolumeClaim.claimName === this.metadata.name) {
									podList.push(el.metadata)
								}
							}
						})
					});
					this.pods = podList;
				})
				.catch(e => { 
					this.pods = [];
					this.msghttp(e);
				});
		}
	}
}
</script>
