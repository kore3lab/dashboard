<template>
	<div class="card-body p-2">
		<div class="row">
			<div class="col-md-12">
				<div class="card card-secondary card-outline">
					<div class="card-body p-2">
						<dl class="row mb-0">
							<dt class="col-sm-3 text-truncate">Create at</dt><dd class="col-sm-9">{{ this.getTimestampString(metadata.creationTimestamp)}} ago ({{ metadata.creationTimestamp }})</dd>
							<dt class="col-sm-3">Name</dt><dd class="col-sm-9">{{ metadata.name }}</dd>
							<dt class="col-sm-3">Namespace</dt><dd class="col-sm-9">{{ metadata.namespace }}</dd>
							<dt class="col-sm-3">Annotations</dt>
							<dd class="col-sm-9 text-truncate">
								<ul class="list-unstyled mb-0">
									<li v-for="(value, name) in metadata.annotations" v-bind:key="name"><span class="badge badge-secondary font-weight-light text-sm mb-1">{{ name }}={{ value }}</span></li>
								</ul>
							</dd>
							<dt class="col-sm-3">Labels</dt>
							<dd class="col-sm-9 text-truncate">
								<ul class="list-unstyled mb-0">
									<li v-for="(value, name) in metadata.labels" v-bind:key="name"><span class="badge badge-secondary font-weight-light text-sm mb-1">{{ name }}={{ value }}</span></li>
								</ul>
							</dd>
							<dt v-if="selector" class="col-sm-3">Selector</dt>
							<dd v-if="selector" class="col-sm-9 text-truncate">
								<span v-for="(value, name) in selector" v-bind:key="name" class="badge badge-secondary font-weight-light text-sm mb-1 mr-1">{{ value }}</span>
							</dd>
							<dt class="col-sm-3">Min Available</dt><dd class="col-sm-9">{{ info.minAvailable }}</dd>
							<dt class="col-sm-3">max Unavailable</dt><dd class="col-sm-9">{{ info.maxUnavailable }}</dd>
							<dt class="col-sm-3">Current Healthy</dt><dd class="col-sm-9">{{ info.currentHealthy }}</dd>
							<dt class="col-sm-3">Desired Healthy</dt><dd class="col-sm-9">{{ info.desiredHealthy }}</dd>
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
			selector: [],
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
			this.info = this.getInfo(data)
			this.selector = this.getSelector(data.spec.selector)
		},
		getInfo(data) {
			return {
				minAvailable: data.spec.minAvailable || "N/A",
				maxUnavailable: data.spec.maxUnavailable || "N/A",
				currentHealthy: data.status.currentHealthy,
				desiredHealthy: data.status.desiredHealthy,
			}
		},
		getSelector(selector) {
			if(!selector) return
			let label = [];
			label = this.stringifyLabels(selector.matchLabels)
			return label
		},
	},
	beforeDestroy(){
		this.$nuxt.$off("onReadCompleted");
	},
}
</script>
