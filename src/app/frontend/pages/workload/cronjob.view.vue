<template>
<div>
	<!-- 1.metdata -->
	<c-metadata dtCols="2" ddCols="10" @navigate="$emit('navigate', arguments[0])">
		<dt class="col-sm-2">Schedule</dt><dd class="col-sm-10">{{ info.schedule }}</dd>
		<dt class="col-sm-2">Active</dt><dd class="col-sm-10">{{ info.active }}</dd>
		<dt class="col-sm-2">Suspend</dt><dd class="col-sm-10">{{ info.suspend }}</dd>
		<dt class="col-sm-2">Last schedule</dt><dd class="col-sm-10">{{ info.lastSchedule }}</dd>
	</c-metadata>
	<!-- 2.jobs -->
	<div v-if="jobs.length>0" class="row">
		<div class="col-md-12">
			<div class="card card-secondary card-outline">
				<div class="card-header p-2"><h3 class="card-title">Jobs</h3></div>
				<div class="card-body group">
					<ul>
						<li v-for="(val, idx) in jobs" v-bind:key="idx">
							<p class="title"><a href="#" @click="$emit('navigate', getViewLink('batch', 'jobs', val.namespace, val.name))">{{ val.name }}</a></p>
							<dl class="row">
								<dt class="col-sm-2">Condition</dt>
								<dd class="col-sm-10">
									<span v-for="(v,i) in val.condition" v-bind:key="i" v-bind:class="{'text-success':v.type=='Complete', 'text-danger':v.type=='Failed'}">{{ v.type }}</span>
								</dd>
								<dt class="col-sm-2">Selector</dt>
								<dd class="col-sm-10">
									<span v-for="(v,i) in val.selector" v-bind:key="i" class="border-box background">{{ v }}</span>
								</dd>
							</dl>
						</li>
					</ul>
				</div>
			</div>
		</div>
	</div>
	<!-- 3.events -->
	<c-events class="row"></c-events>
</div>
</template>
<script>
import VueMetadataView	from "@/components/view/metadataView.vue";
import VueEventsView 	from "@/components/view/eventsView.vue";

export default {
	components: {
		"c-metadata": { extends: VueMetadataView },
		"c-events": { extends: VueEventsView }
	},
	data() {
		return {
			info: [],
			jobs: []
		}
	},
	mounted() {
		this.$nuxt.$on("view-data-read-completed", (data) => {
			if(!data) return
			this.info = {
				schedule: data.spec.schedule,
				suspend: data.spec.suspend,
				active: data.status.active?data.status.active.length:"0",
				lastSchedule: this.getElapsedTime(data.status.lastScheduleTime)
			};

			//job-list
			this.$axios.get(this.getApiUrl("batch","jobs",data.metadata.namespace))
				.then(resp => {
					let list = [];
					resp.data.items.forEach(el => {
						if (el.metadata.ownerReferences && el.metadata.ownerReferences[0].uid === data.metadata.uid) {
							list.push({
								name: el.metadata.name,
								namespace: el.metadata.namespace,
								condition: el.status.conditions ? el.status.conditions.filter(e=>{return e.status === 'True'}): [],
								selector: this.stringifyLabels(el.spec.selector? el.spec.selector.matchLabels : ''),
							});
						}
					})
					this.jobs = list;
				}).catch(e => {
					this.msghttp(e);
					this.jobs = [];
				});

		});
	},
	methods: {},
	beforeDestroy(){
		this.$nuxt.$off("view-data-read-completed");
	}
}
</script>
