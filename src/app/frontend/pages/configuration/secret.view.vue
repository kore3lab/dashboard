<template>
	<div class="card-body p-2">
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
							<dt v-if="metadata.ownerReferences" class="col-sm-2 text-truncate">Controlled By</dt>
							<dd v-if="metadata.ownerReferences" class="col-sm-10">{{ metadata.ownerReferences[0].kind }} <a href="#" @click="$emit('navigate', getViewLink(controller.g, controller.k, metadata.namespace, metadata.ownerReferences[0].name))">{{ metadata.ownerReferences[0].name }}</a></dd>
							<dt class="col-sm-2">Type</dt><dd class="col-sm-10">{{ origin.type }}</dd>
						</dl>
					</div>
				</div>
			</div>
		</div>

		<div v-show="secretData" class="row">
			<div class="col-md-12">
				<div class="card card-secondary card-outline m-0">
					<div class="card-header p-2"><h3 class="card-title text-md">Data</h3></div>
					<div class="card-body p-2">
						<dl v-for="(val, idx) in secretData" v-bind:key="idx" class="row mb-0 card-body p-2">
							<dt class="col-sm-12">
								<span class="card-title mb-2">{{ val.key }}</span>
								<button type="button" class="btn btn-tool" @click="onShow(idx)"><i v-show="isShow[idx]" class="fas fa-eye-slash bg-gray-light"></i><i v-show="!isShow[idx]" class="fas fa-eye bg-gray-light"></i></button>
							</dt>
							<dd class="col-sm-12" v-model="isShow[idx]">
								<b-form-textarea v-show="isShow[idx]" max-rows="100" v-model="val.val" class="card-body p-2 border text-sm"></b-form-textarea>
								<b-form-textarea v-show="!isShow[idx]" max-rows="100" v-model="val.decoval" class="card-body p-2 border text-sm"></b-form-textarea>
							</dd>
						</dl>
						{{ te }}
						<b-button variant="primary" size="sm" class="mb-1 ml-2" @click="onSave()">Save</b-button>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>
<script>
import axios			from "axios"
import VueAceEditor from "@/components/aceeditor";

export default {
	components: {
		"c-aceeditor": { extends: VueAceEditor },
	},
	data() {
		return {
			te: '',
			dataValue: [],
			metadata: {},
			info: [],
			origin: [],
			controller: [],
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
			this.onSync(data)
		});
		this.$nuxt.$emit("onCreated",'')
	},
	methods: {
		onSync(data) {
			this.controller = this.getController(data.metadata.ownerReferences);
			this.secretData = this.getData(data.data);
		},
		getData(data) {
			if(!data) return false
			let list = [];
			let i=0;
			let key = Object.keys(data)
			key.forEach(el => {
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
			axios.put(`${this.backendUrl()}/raw/clusters/${this.currentContext()}`, this.origin)
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
