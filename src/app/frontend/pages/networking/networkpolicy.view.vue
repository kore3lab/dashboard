<template>
<div>
	<!-- 1. metadata -->
	<c-metadata dtCols="2" ddCols="10">
		<dt class="col-sm-2">Pod Selector</dt>
		<dd class="col-sm-10">
			<span v-for="(value, name) in podSelector" v-bind:key="name" class="border-box background">{{ value }}</span>
			<span v-if="podSelector.length==0">(empty) (Allowing the specific traffic to all pods in this namespace)</span>
		</dd>
	</c-metadata>
	<!-- 2. ingress -->
	<div v-if="ingress.length>0" class="row">
		<div class="col-md-12">
			<div class="card card-secondary card-outline">
				<div class="card-header p-2"><h3 class="card-title">Ingress</h3></div>
				<div class="card-body group">
					<ul>
						<li v-for="(ing, idx) in ingress" v-bind:key="idx">
							<dl class="row">
								<dt class="col-sm-2">Ports</dt>
								<dd class="col-sm-10">
									<span v-for="(value, idx) in ing.ports" v-bind:key="idx">{{ value["protocol"] }}:{{ value["port"] }} </span>
								</dd>
								<dt class="col-sm-2" v-if="ing.from">From</dt>
								<dd class="col-sm-10" v-if="ing.from">
									<b-table-lite thead-class="d-none" small :items="ing.from" borderless></b-table-lite>
								</dd>
							</dl>
						</li>
					</ul>
				</div>
			</div>
		</div>
	</div>
	<!-- 3. egress -->
	<div v-if="egress.length>0" class="row">
		<div class="col-md-12">
			<div class="card card-secondary card-outline">
				<div class="card-header p-2"><h3 class="card-title">Egress</h3></div>
				<div class="card-body group">
					<ul>
						<li v-for="(eg, idx) in egress" v-bind:key="idx">
							<dl class="row">
								<dt class="col-sm-2">Ports</dt>
								<dd class="col-sm-10">
									<span v-for="(port, idx) in eg.ports" v-bind:key="idx">{{ port }} </span>
								</dd>
								<dt class="col-sm-2" v-if="eg.to">To</dt>
								<dd class="col-sm-10" v-if="eg.to">
									<b-table-lite thead-class="d-none" borderless small :items="eg.to"></b-table-lite>
								</dd>
							</dl>
						</li>
					</ul>
				</div>
			</div>
		</div>
	</div>
	<!-- 4. events -->
	<c-events class="row"></c-events>

</div>
</template>
<script>
import VueMetadataView	from "@/components/view/metadataView.vue";
import VueEventsView	from "@/components/view/eventsView.vue";

export default {
	components: {
		"c-metadata": { extends: VueMetadataView },
		"c-events": { extends: VueEventsView }
	},
	data() {
		return {
			podSelector: [],
			ingress: [],
			egress: []
		}
	},
	mounted() {
		this.$nuxt.$on("view-data-read-completed", (data) => {
			if(!data) return
			this.podSelector = this.stringifyLabels(data.spec.podSelector? data.spec.podSelector.matchLabels : '');
			this.ingress = data.spec.ingress?data.spec.ingress:[];
			this.egress = data.spec.egress?data.spec.egress:[];
		});
	},
	methods: {},
	beforeDestroy(){
		this.$nuxt.$off("view-data-read-completed");
	},
}
</script>
