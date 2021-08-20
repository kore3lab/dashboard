<template>
<div>
	<!-- 1. metadata -->
	<c-metadata dtCols="3" ddCols="9" @navigate="$emit('navigate', arguments[0])">
		<dt v-if="imagePullSecrets.length>0" class="col-sm-3">ImagePullSecrets</dt>
		<dd v-if="imagePullSecrets.length>0" class="col-sm-9"><span v-for="(value, idx) in imagePullSecrets" v-bind:key="idx" class="mr-1"><a href="#" @click="$emit('navigate', getViewLink('','secrets',metadata.namespace,value.name))">{{ value.name }}</a></span></dd>
	</c-metadata>
	<!-- 2. mountable secrets -->
	<div class="row">
		<div class="col-md-12">
			<div class="card card-secondary card-outline">
				<div class="card-header p-2"><h3 class="card-title">Mountable secrets</h3></div>
				<div class="card-body group">
					<ul>
						<li v-for="(val, idx) in secrets" v-bind:key="idx">
							<dl class="row">
								<dt class="col-sm-2">Name</dt><dd class="col-sm-10">{{ val.name }}</dd>
								<dt class="col-sm-2">Value</dt>
								<dd class="col-sm-10">
									<span v-show="!isShow[idx].show">{{ '••••••••••••••••' }} <i v-on:click="isShow[idx].show = true" class="fas fa-unlock"></i></span>
									<span v-show="isShow[idx].show">{{ val.value }}</span>
								</dd>
								<dt class="col-sm-2">Created at</dt><dd class="col-sm-10">{{ val.createdAt }}</dd>
								<dt class="col-sm-2">Type</dt><dd class="col-sm-10">{{ val.type }}</dd>
							</dl>
						</li>
					</ul>
				</div>
			</div>
		</div>
	</div>
	<!--3. events -->
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
			metadata: {},
			imagePullSecrets: [],
			secrets: [],
			isShow: {},
		}
	},
	mounted() {
		this.$nuxt.$on("view-data-read-completed", (data) => {
			if(!data) return
			this.metadata = data.metadata;
			this.imagePullSecrets = data.imagePullSecrets? data.imagePullSecrets: [];
			this.secrets = this.secrets ?this.getSecrets(data.secrets):[];
		});
	},
	methods: {
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
		this.$nuxt.$off("view-data-read-completed");
	},
}
</script>
