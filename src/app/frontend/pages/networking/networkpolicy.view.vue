<template>
<div>
	<!-- 1. metadata -->
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
						<dt class="col-sm-2">Pod Selector</dt>
						<dd class="col-sm-10">
							<span v-for="(value, name) in podSelector" v-bind:key="name" class="border-box background">{{ value }}</span>
							<span v-if="podSelector.length==0">(empty) (Allowing the specific traffic to all pods in this namespace)</span>
						</dd>
					</dl>
				</div>
			</div>
		</div>
	</div>

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
									<span v-for="(port, idx) in ing.ports" v-bind:key="idx">{{ port }} </span>
								</dd>
								<dt class="col-sm-2">From</dt>
								<dd class="col-sm-10">
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
								<dt class="col-sm-2">To</dt>
								<dd class="col-sm-10">
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
			podSelector: [],
			ingress: {},
			egress: {}
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
			this.podSelector = this.stringifyLabels(data.spec.podSelector? data.spec.podSelector.matchLabels : '')
			// ingress
			this.ingress = [];
			Object.assign([], data.spec.ingress).forEach(el=> {
				let from = []
				if(el["from"]) {
					el["from"].forEach(e=> {
						let key = Object.keys(e)[0];
						from.push({column:key, value: e[key]})
					});
				}
				
				this.ingress.push({from: from, ports: el.ports});
			});

			// egress
			this.egress = [];
			Object.assign([], data.spec.egress).forEach(el=> {
				let to = []
				if(el["to"]) {
					el["to"].forEach(e=> {
						let key = Object.keys(e)[0];
						to.push({column:key, value: e[key]})
					});
				}
				
				this.egress.push({to: to, ports: el.ports});
			});

		},
	},
	beforeDestroy(){
		this.$nuxt.$off("onReadCompleted");
	},
}
</script>
