<template>
	<div>
		<div class="row">
			<div class="col-md-12">
				<div class="card card-secondary card-outline">
					<div class="card-body p-2">
						<dl class="row mb-0">
							<dt class="col-sm-2 text-truncate">Create at</dt><dd class="col-sm-10">{{ this.getTimestampString(metadata.creationTimestamp)}} ago ({{ metadata.creationTimestamp }})</dd>
							<dt class="col-sm-2">Name</dt><dd class="col-sm-10">{{ metadata.name }}</dd>
							<dt class="col-sm-2">Namespace</dt><dd class="col-sm-10">{{ metadata.namespace }}</dd>
							<dt class="col-sm-2">Annotations</dt>
							<dd class="col-sm-10 text-truncate">
								<ul class="list-unstyled mb-0">
									<li v-for="(value, name) in metadata.annotations" v-bind:key="name"><span class="badge badge-secondary font-weight-light text-sm mb-1">{{ name }}:{{ value }}</span></li>
								</ul>
							</dd>
							<dt class="col-sm-2">Labels</dt>
							<dd class="col-sm-10 text-truncate">
								<ul class="list-unstyled mb-0">
									<li v-for="(value, name) in metadata.labels" v-bind:key="name"><span class="badge badge-secondary font-weight-light text-sm mb-1">{{ name }}:{{ value }}</span></li>
								</ul>
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

		<div v-if="isJobs" class="row">
			<div class="col-md-12">
				<div class="card card-secondary card-outline m-0">
					<div class="card-header p-2"><h3 class="card-title text-md">Jobs</h3></div>
					<div class="card-body p-2">
						<dl v-for="(val, idx) in jobs" v-bind:key="idx" class="row mb-0 card-body p-2 border-bottom">
							<dt class="col-sm-12 mb-1"><a href="#" @click="$emit('navigate', getViewLink('batch', 'jobs', val.namespace, val.name))">{{ val.name }}</a></dt>
							<dt class="col-sm-2 text-truncate">Condition</dt><dd class="col-sm-10"><span v-for="(v,i) in val.condition" v-bind:key="i" v-bind:class="v.style" class="badge font-weight-light text-sm mr-1">{{ v.type }}</span></dd>
							<dt class="col-sm-2 text-truncate">Selector</dt><dd class="col-sm-10"><span v-for="(v,i) in val.selector" v-bind:key="i" class="badge badge-secondary font-weight-light text-sm mr-1">{{ v }}</span></dd>
						</dl>
					</div>
				</div>
			</div>
		</div>

		<div class="row">
			<div class="col-md-12">
				<div class="card card-secondary card-outline m-0">
					<div class="card-header p-2"><h3 class="card-title text-md">Events</h3></div>
					<div class="card-body p-2">
						<dl v-for="(val, idx) in event" v-bind:key="idx" class="row mb-0 card-body p-2 border-bottom">
							<dt class="col-sm-12"><p v-bind:class="val.type" class="mb-1">{{ val.name }}</p></dt>
							<dt class="col-sm-2 text-truncate">Source</dt><dd class="col-sm-10">{{ val.source }}</dd>
							<dt class="col-sm-2 text-truncate">Count</dt><dd class="col-sm-10">{{ val.count }}</dd>
							<dt class="col-sm-2 text-truncate">Sub-object</dt><dd class="col-sm-10">{{ val.subObject }}</dd>
							<dt class="col-sm-2 text-truncate">Last seen</dt><dd class="col-sm-10">{{ val.lastSeen }}</dd>
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
			event: [],
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
			this.event = this.getEvents(data.metadata.uid,'fieldSelector=involvedObject.name='+data.metadata.name);
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
						if (el.type === 'Complete') list.push({ type: 'Complete', style: 'badge-success' })
						else list.push({ type: 'Failed', style: 'badge-danger' })
					}
				})
				return list
			}
		},
		getScheduleTime(time) {
			let tran = this.getElapsedTime(time)
			if (tran.str === ""){
				return '-'
			}
			return tran.str
		},
	},
	beforeDestroy(){
		this.$nuxt.$off("onReadCompleted");
	},
}
</script>
