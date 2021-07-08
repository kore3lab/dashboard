<template>
<div>
	<!-- 1.metdata -->
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
						<dt class="col-sm-2">Schedule</dt><dd class="col-sm-10">{{ info.schedule }}</dd>
						<dt class="col-sm-2">Active</dt><dd class="col-sm-10">{{ info.active }}</dd>
						<dt class="col-sm-2">Suspend</dt><dd class="col-sm-10">{{ info.suspend }}</dd>
						<dt class="col-sm-2">Last schedule</dt><dd class="col-sm-10">{{ info.lastSchedule }}</dd>
					</dl>
				</div>
			</div>
		</div>
	</div>

	<!-- 2.jobs -->
	<div v-if="isJobs" class="row">
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
									<span v-for="(v,i) in val.condition" v-bind:key="i">{{ v.type }}</span>
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
	<c-events class="row" v-model="metadata.uid"></c-events>

</div>
</template>
<script>
import VueEventsView from "@/components/view/eventsView.vue";

export default {
	components: {
		"c-events": { extends: VueEventsView }
	},
	data() {
		return {
			metadata: {},
			info: [],
			jobs: [],
			isJobs: false,
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
			this.isJobs= false;
			this.info = this.getInfo(data);
			this.jobs = this.getJobs(data.metadata.uid)
		},
		getJobs(uid) {
			let list = [];
			this.$axios.get(this.getApiUrl("batch","jobs",this.selectedNamespace))
					.then(resp => {
						resp.data.items.forEach(el => {
							if (el.metadata.ownerReferences && el.metadata.ownerReferences[0].uid === uid) {
								this.isJobs = true;
								list.push({
									name: el.metadata.name,
									namespace: el.metadata.namespace,
									condition: this.getCondition(el.status.conditions),
									selector: this.stringifyLabels(el.spec.selector? el.spec.selector.matchLabels : ''),
								});
							}
						})
						this.jobs = list
						return list
					})
			return list
		},
		getInfo(data) {
			return {
				schedule: data.spec.schedule,
				suspend: data.spec.suspend,
				active: data.status.active?data.status.active.length:"0",
				lastSchedule: this.getScheduleTime(data.status.lastScheduleTime),
			}
		},
		getCondition(con) {
			let list = []
			if(con) {
				con.forEach(el => {
					if (el.status === 'True') {
						if (el.type === 'Complete') list.push({ type: 'Complete' })
						else list.push({ type: 'Failed' })
					}
				})
				return list
			}
		},
		getScheduleTime(d) {
			let tran = this.getElapsedTime(d)
			return tran == "" ? "-": tran;
		},
	},
	beforeDestroy(){
		this.$nuxt.$off("onReadCompleted");
	},
}
</script>
