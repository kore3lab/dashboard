<template>
<div>
	<div class="col-md-12" v-show="events.length>0">
		<div class="card card-secondary card-outline">
			<div class="card-header p-2"><h3 class="card-title">Events</h3></div>
			<div class="card-body group">
				<ul>
					<li  v-for="(d, idx) in events" v-bind:key="idx">
						<p v-bind:class="{'text-warning': d.type!='Normal'}" class="mb-1 title">{{ d.message }}</p>
						<dl class="row">
							<dt class="col-sm-2">Source</dt><dd class="col-sm-10">{{ d.source }}</dd>
							<dt class="col-sm-2">Count</dt><dd class="col-sm-10">{{ d.count }}</dd>
							<dt class="col-sm-2">Sub-object</dt><dd class="col-sm-10">{{ d.subObject }}</dd>
							<dt class="col-sm-2">Last seen</dt><dd class="col-sm-10">{{ d.lastTimestamp || "-" }}</dd>
						</dl>
					</li>
				</ul>
			</div>
		</div>
	</div>
</div>
</template>
<script>

export default {
	props:["value"],
	data () {
		return {
			events: []
		}
	},
	watch: {
		value(d) { this.onSync(d) }
	},
	methods: {
		onSync(data) {
			if(!data) return
			this.events = [];
			this.$axios.get(this.getApiUrl('','events','', '',`fieldSelector=involvedObject.uid=${data.metadata.uid}`))
				.then(resp => {
					this.events = resp.data.items;
				}).catch(ex => {
					console.error(ex.message);
				});
		}
	}
}
</script>
