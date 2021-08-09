<template>
<div>
	<!-- 1. metadata -->
	<c-metadata v-model="metadata" dtCols="2" ddCols="10" @navigate="$emit('navigate', arguments[0])">
		<dt class="col-sm-2">Type</dt><dd class="col-sm-10">{{ origin.type }}</dd>
	</c-metadata>
	<!-- 2. data -->
	<div v-show="secretData" class="row">
		<div class="col-md-12">
			<div class="card card-secondary card-outline m-0">
				<div class="card-header p-2"><h3 class="card-title">Data</h3></div>
				<div class="card-body group">
					<ul>
						<li v-for="(val, idx) in secretData" v-bind:key="idx">
							<p class="title">
								<span class="card-title mb-2">{{ val.key }}</span>
								<button type="button" class="btn btn-tool" @click="onShow(idx)"><i v-show="isShow[idx]" class="fas fa-eye-slash bg-gray-light"></i><i v-show="!isShow[idx]" class="fas fa-eye bg-gray-light"></i></button>
							</p>
							<div class="col-sm-12" :v-model="isShow[idx]">
								<b-form-textarea v-show="isShow[idx]" max-rows="100" v-model="val.val" class="card-body p-2 border text-sm"></b-form-textarea>
								<b-form-textarea v-show="!isShow[idx]" max-rows="100" v-model="val.decoval" class="card-body p-2 border text-sm"></b-form-textarea>
							</div>
						</li>
					</ul>
					{{ te }}
					<b-button variant="primary" size="sm" class="m-1" @click="onSave()">Save</b-button>
				</div>
			</div>
		</div>
	</div>
</div>
</template>
<script>
import VueMetadataView	from "@/components/view/metadataView.vue";

export default {
	components: {
		"c-metadata": { extends: VueMetadataView }
	},
	data() {
		return {
			te: '',
			dataValue: [],
			metadata: {},
			origin: [],
			secretData: [],
			isShow: [],
			fields: [
				{ key: "name", label: "Name" },
				{ key: "endpoints", label: "Endpoints" },
			],
		}
	},
	mounted() {
		this.$nuxt.$on("onReadCompleted", (data) => {
			if(!data) return
			this.origin = data;
			this.metadata = data.metadata;
			this.secretData = this.getData(data.data);
		});
		this.$nuxt.$emit("onCreated",'')
	},
	methods: {
		getData(data) {
			if(!data) return false
			let list = [];
			let i=0;
			Object.keys(data).forEach(el => {
				this.isShow[i] = true;
				this.dataValue.push(data[el])
				list.push({
					key: el,
					val: data[el],
					decoval: atob(data[el]),
				})
				i++
			})
			return list
		},
		onSave() {
			let list = {};
			let count = 0;
			this.secretData.forEach(el => {
				if(this.isShow[count]) {
					list[el.key] = el.val
				} else list[el.key] = btoa(el.decoval)
				count++
			})
			this.origin.data = list
			this.$axios.put(`/raw/clusters/${this.currentContext()}`, this.origin)
					.then( resp => {
						this.origin = Object.assign({}, resp.data);
						this.toast(`Secret ${ this.metadata.name } successfully updated.`, "success");
					})
					.catch(e => { this.msghttp(e);});
		},
		onShow(idx) {
			if(this.isShow[idx]) {
				try{
					this.secretData[idx].decoval = atob(this.secretData[idx].val)
				} catch(error) {
					console.log(error)
				}
			} else {
				try{
					this.secretData[idx].val = btoa(this.secretData[idx].decoval)
				} catch(error) {
					console.log(error)
				}
			}
			this.te = ''
			this.isShow[idx] = !this.isShow[idx]
			this.te = ' '
		},
	},
	beforeDestroy(){
		this.$nuxt.$off("onReadCompleted");
	},
}
</script>
