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
						<dt class="col-sm-3">Reference</dt>
						<dd class="col-sm-9"><a href="#" @click="$emit('navigate', getViewLink(ref.group.g, ref.group.k, metadata.namespace, ref.name))">{{ ref.kind }}/{{ ref.name }}</a></dd>
						<dt class="col-sm-3">Min Pods</dt><dd class="col-sm-9">{{ info.minPods }}</dd>
						<dt class="col-sm-3">Max Pods</dt><dd class="col-sm-9">{{ info.maxPods }}</dd>
						<dt class="col-sm-3">Replicas</dt><dd class="col-sm-9">{{ info.replicas }}</dd>
					</dl>
				</div>
			</div>
		</div>
	</div>
	<!-- 2. evnets -->
	<c-events class="row" v-model="metadata.uid"></c-events>

</div>
</template>
<script>
import VueEventsView	from "@/components/view/eventsView.vue";

export default {
	components: {
		"c-events": { extends: VueEventsView }
	},
	data() {
		return {
			metadata: {},
			info: [],
			ref: []
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
			this.ref = this.getRef(data.spec.scaleTargetRef)
		},
		getInfo(data) {
			return {
				minPods: data.spec.minReplicas || 0,
				maxPods: data.spec.maxReplicas || 0,
				replicas: data.status.currentReplicas
			}
		},
		getRef(ref) {
			if(!ref) return

			return {
				group: this.getController(ref),
				name: ref.name,
				kind: ref.kind
			}
		},
	},
	beforeDestroy(){
		this.$nuxt.$off("onReadCompleted");
	},
}
</script>
