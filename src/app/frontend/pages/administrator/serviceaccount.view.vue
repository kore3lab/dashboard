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
									<li v-for="(value, name) in metadata.annotations" v-bind:key="name"><span class="badge badge-secondary font-weight-light text-sm mb-1">{{ name }}:{{ value }}</span></li>
								</ul>
							</dd>
							<dt class="col-sm-3">Labels</dt>
							<dd class="col-sm-9 text-truncate">
								<ul class="list-unstyled mb-0">
									<li v-for="(value, name) in metadata.labels" v-bind:key="name"><span class="badge badge-secondary font-weight-light text-sm mb-1">{{ name }}:{{ value }}</span></li>
								</ul>
							</dd>
							<dt v-if="IPS" class="col-sm-3 text-truncate">ImagePullSecrets</dt>
							<dd v-if="IPS" class="col-sm-9"><span v-for="(val, idx) in IPS" v-bind:key="idx" class="mr-1"><a href="#" @click="$emit('navigate', getViewLink('','secrets',metadata.namespace,val))">{{ val }}</a></span></dd>
						</dl>
					</div>
				</div>
			</div>
		</div>

		<div class="row">
			<div class="col-md-12">
				<div class="card card-secondary card-outline">
					<div class="card-header p-2"><h3 class="card-title text-md">Mountable secrets</h3></div>
					<div class="card-body p-2">
						<dl v-for="(val, idx) in secrets" v-bind:key="idx" class="row mb-1">
							<dt class="col-sm-2">Name</dt><dd class="col-sm-10">{{ val.name }}</dd>
							<dt class="col-sm-2">Value</dt>
							<dd class="col-sm-10">
								<span v-show="!isShow[idx].show">{{ '••••••••••••••••' }} <i v-on:click="isShow[idx].show = true" class="fas fa-unlock"></i></span>
								<span v-show="isShow[idx].show">{{ val.value }}</span>
							</dd>
							<dt class="col-sm-2">Created at</dt><dd class="col-sm-10">{{ val.createdAt }}</dd>
							<dt class="col-sm-2">Type</dt><dd class="col-sm-10">{{ val.type }}</dd>
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
			IPS: [],
			secrets: [],
			isShow: [],
		}
	},
	mounted() {
		this.$nuxt.$on("onReadCompleted", (data) => {
			if(!data) return
			this.origin = data;
			this.metadata = data.metadata;
			this.onSync(data)
		});
		this.$nuxt.$emit("onCreated",'')
	},
	methods: {
		onSync(data) {
			this.event = this.getEvents(data.metadata.uid,'fieldSelector=involvedObject.name='+data.metadata.name);
			this.IPS = this.getIPS(data.imagePullSecrets)
			this.secrets = this.getSecrets(data.secrets)
		},
		getIPS(ips) {
			if(!ips) return

			let list = []
			ips.map((v) => {
				if(v.name) list.push(v.name)
			})
			return list
		},
		getSecrets(secrets) {
			this.isShow = [];
			if(!secrets) return
			let list = []
			let count = 0;
			secrets.map((v) => {
				this.$axios.get(this.getApiUrl("","secrets",this.metadata.namespace,v.name))
						.then((resp) => {
							this.isShow.push({show: false})
							list.push({
								name: v.name,
								value: resp.data.data.token,
								createdAt: new Date(resp.data.metadata.creationTimestamp),
								type: resp.data.type,
								index: count
							})
							count++
						})
			})
			return list
		},
	},
	beforeDestroy(){
		this.$nuxt.$off("onReadCompleted");
	},
}
</script>
