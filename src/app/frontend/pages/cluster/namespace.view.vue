<template>
	<div class="card-body p-2">
		<div class="row">
			<div class="col-md-12">
				<div class="card card-secondary card-outline">
					<div class="card-body p-2">
						<dl class="row mb-0">
							<dt class="col-sm-3 text-truncate">Create at</dt><dd class="col-sm-9">{{ this.getTimestampString(metadata.creationTimestamp)}} ago ({{ metadata.creationTimestamp }})</dd>
							<dt class="col-sm-3">Name</dt><dd class="col-sm-9">{{ metadata.name }}</dd>
							<dt class="col-sm-3">Annotations</dt>
							<dd class="col-sm-9 text-truncate">
								<ul class="list-unstyled mb-0">
									<li v-for="(value, name) in metadata.annotations" v-bind:key="name"><span class="badge badge-secondary font-weight-light text-sm mb-1">{{ name }}:{{ value }}</span></li>
								</ul>
							</dd>
							<dt class="col-sm-3">Labels</dt>
							<dd class="col-sm-9 text-truncate">
								<ul class="list-unstyled mb-0">
									<li v-for="(value, name) in metadata.labels" v-bind:key="name"><span class="badge badge-secondary font-weight-light text-sm mb-1">{{ name }}:{{ value }}</span></li>
								</ul>
							</dd>
							<dt v-if="metadata.finalizers" class="col-sm-3">Finalizers</dt><dd v-if="metadata.finalizers" class="col-sm-9"><span v-for="(val, idx) in metadata.finalizers" v-bind:key="idx" class="badge badge-secondary font-weight-light text-sm mb-1 mr-1">{{ val }}</span></dd>
							<dt class="col-sm-3">UID</dt><dd class="col-sm-9">{{ metadata.uid }}</dd>
							<dt v-if="metadata.ownerReferences" class="col-sm-3 text-truncate">Controlled By</dt>
							<dd v-if="metadata.ownerReferences" class="col-sm-9">{{ metadata.ownerReferences[0].kind }} <a href="#" @click="$emit('navigate', getViewLink(controller.g, controller.k, metadata.namespace, metadata.ownerReferences[0].name))">{{ metadata.ownerReferences[0].name }}</a></dd>
							<dt class="col-sm-3">Status</dt><dd class="col-sm-9" v-bind:class="status.style">{{ status.type }}</dd>
							<dt v-if="isQuota" class="col-sm-3">Resource Quotas</dt><dd v-if="isQuota" class="col-sm-9"><span v-for="(val, idx) in quotas" v-bind:key="idx"><a href="#" @click="$emit('navigate', getViewLink('', 'resourcequotas', metadata.name,val))">{{ val }} </a></span></dd>
							<dt v-if="isLimit" class="col-sm-3">Limit Ranges</dt><dd v-if="isLimit" class="col-sm-9"><span v-for="(val, idx) in limits" v-bind:key="idx"><a href="#" @click="$emit('navigate', getViewLink('', 'limitranges', metadata.name,val))">{{ val }} </a></span></dd>
						</dl>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>
<script>
import axios			from "axios"

export default {
	data() {
		return {
			metadata: {},
			controller: [],
			status: [],
			quotas: [],
			limits: [],
			isQuota: false,
			isLimit: false,
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
			this.status = this.getStatus(data.status.phase)
			this.controller = this.getController(data.metadata.ownerReferences)
			this.getQuotas(data.metadata.name)
			this.getLimit(data.metadata.name)
		},
		getStatus(phase) {
			let status
			if(phase) {
				status = phase
				if(phase === 'Active') status = {type: phase, style: 'text-success'}
				else status = {type: phase, style: 'text-secondary'}
				return status
			}
		},
		getQuotas(name) {
			this.isQuota = false;
			this.quotas = [];
			axios.get(this.getApiUrl("","resourcequotas",name))
			.then(resp => {
				resp.data.items.forEach(el =>{
					this.isQuota = true
					this.quotas.push(el.metadata.name)
				})
			})
		},
		getLimit(name) {
			this.isLimit = false;
			this.limits = [];
			axios.get(this.getApiUrl("","limitranges",name))
			.then(resp => {
				resp.data.items.forEach(el =>{
					this.isLimit = true
					this.limits.push(el.metadata.name)
				})
			})
		},
	},
	beforeDestroy(){
		this.$nuxt.$off("onReadCompleted");
	},
}
</script>
