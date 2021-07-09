<template>
<div>
	<!-- 1. metadata -->
	<div class="row">
		<div class="col-md-12">
			<div class="card card-secondary card-outline">
				<div class="card-body p-2">
					<dl class="row mb-0">
						<dt class="col-sm-3">Create at</dt><dd class="col-sm-9">{{ this.getTimestampString(metadata.creationTimestamp)}} ago ({{ metadata.creationTimestamp }})</dd>
						<dt class="col-sm-3">Name</dt><dd class="col-sm-9">{{ metadata.name }}</dd>
						<dt class="col-sm-3">Namespace</dt><dd class="col-sm-9">{{ metadata.namespace }}</dd>
						<dt class="col-sm-3">Annotations</dt>
						<dd class="col-sm-9">
							<ul class="list-unstyled mb-0">
								<li v-for="(value, name) in metadata.annotations" v-bind:key="name">{{ name }}=<span class="font-weight-light">{{ value }}</span></li>
							</ul>
						</dd>
						<dt class="col-sm-3">Labels</dt>
						<dd class="col-sm-9">
							<span v-for="(value, name) in metadata.labels" v-bind:key="name" class="label">{{ name }}={{ value }}</span>
						</dd>
						<dt v-if="selector" class="col-sm-3">Selector</dt>
						<dd v-if="selector" class="col-sm-9">
							<span v-for="(value, name) in selector" v-bind:key="name" class="border-box background">{{ value }}</span>
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
			this.info =  {
				minAvailable: data.spec.minAvailable || "N/A",
				maxUnavailable: data.spec.maxUnavailable || "N/A",
				currentHealthy: data.status.currentHealthy,
				desiredHealthy: data.status.desiredHealthy,
			}
			this.selector = this.stringifyLabels(data.spec.selector.matchLabels)
		});
		this.$nuxt.$emit("onCreated",'')
	},
	beforeDestroy(){
		this.$nuxt.$off("onReadCompleted");
	},
}
</script>
