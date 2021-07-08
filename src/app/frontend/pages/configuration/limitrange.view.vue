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
					</dl>
				</div>
			</div>
		</div>
	</div>

	<!-- 2. Container Limits -->
	<div class="row">
		<div class="col-md-12">
			<div class="card card-secondary card-outline">
				<div class="card-header p-2"><h3 class="card-title">Limits</h3></div>
				<div class="card-body group">
					<ul>
						<li v-for="(key, idx) in Object.keys(limits)" v-bind:key="idx">
							<p class="title">{{key}}</p>
							<b-table-lite thead-class="d-none" small :items="limits[key]" borderless>
								<template v-slot:cell(value)="d">
									<span v-for="(v,k) in d.item.value" v-bind:key="k" class="border-box">{{k}}={{v}}</span>
								</template>
							</b-table-lite>
						</li>
					</ul>
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
			limits: {}
		}
	},
	mounted() {
		this.$nuxt.$on("onReadCompleted", (data) => {
			if(!data) return
			this.metadata = data.metadata;
			this.limits = [];
			data.spec.limits.forEach(el=> {
				let d = []
				for(let key in el) {
					if(key != "type") d.push({key:key, value:el[key]});
				}
				this.limits[el["type"]] = d;
			});
		});
		this.$nuxt.$emit("onCreated",'')
	},
	beforeDestroy(){
		this.$nuxt.$off("onReadCompleted");
	},
}
</script>
