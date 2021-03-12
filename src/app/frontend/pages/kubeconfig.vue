<template>
	<div class="content-wrapper">
		<div class="row"  style="padding-top:10%">
			<div class="col-md-3"></div>
			<div class="col-md-6 m-3">
				<div class="text-right"><a href="#" @click="$router.back(-1)"><b-icon icon="x" class="h2 text-right border" variant="secondary"></b-icon></a></div>
				<h2>Add Clusters from Kubeconfig</h2>
				<p class="text-left">Add clusters by clicking the Add Cluster button. You'll need to obtain a working kubeconfig for the cluster you want to add. You can either browse it from the file system or paste it as a text from the clipboard. </p>
				<b-tabs content-class="col-md-12" card>
					<b-tab title="Select kubeconfig file" active>
						<b-form-group label="Select a file:" label-cols-sm="2" label-size="sm">
							<b-form-file v-model="yamlFile" size="sm" placeholder="Choose a file" drop-placeholder="Drop file here..." @input="onFileSelected()"></b-form-file>
						</b-form-group>
					</b-tab>
					<b-tab title="Paste as text">
						<b-form-textarea v-model="yamlText" rows="5" max-rows="6" @change="onTextChange()"></b-form-textarea>
					</b-tab>
				</b-tabs>
				<div class="row">
					<div class="col-10">
						<b-form-select size="sm" v-model="selected" :options="contextList"></b-form-select>
					</div>
					<div class="col-2">
						<b-button variant="primary" size="sm" @click="onAddContext()">Add a cluster</b-button>
					</div>
				</div>
			</div>
			<div class="col-md-3"></div>
		</div>
	</div>
</template>
<script>
import axios			from "axios"
import {load as toJSON} from "js-yaml";
export default {
	components: {
	},
	data() {
		return {
			kubeconfig: "aaaa",
			yamlFile: null,
			yamlText: "",
			selected: null,
			contextList: [
				{ value: null, text: "Please select a context" }
			]
		}
	},
	layout: "default",
	created() {
	},
	mounted() {
	},
	beforeDestroy(){
	},
	methods: {
		onFileSelected() {
			let configReader = this.doReadConfig;
			let reader = new FileReader();
			reader.onload = function(ev) {
				configReader(ev.target.result);
			};
			reader.readAsText(this.yamlFile);
		},
		onTextChange() {
			this.doReadConfig(this.yamlText);
		},
		doReadConfig(json) {
			let conf;
			try {
				conf = toJSON(json);
				if(conf && typeof conf === "object") {
					this.contextList.splice(1, this.contextList.length-2);
					conf["contexts"].forEach(el => {
						this.contextList.push( {text: el.name, value: el} );
					});
					this.kubeconfig = conf;
				} else {
					this.toast("The selected file is invalid.","danger");
				}
			} catch (e) { this.toast("The selected file is invalid.","danger") }
		},
		onAddContext() {
			try {
				if (!this.selected) {
					this.mesbox("Not selected a context");
					return
				}

				let cluster = this.kubeconfig.clusters.find(el=> el.name == this.selected.context["cluster"]);
				let user = this.kubeconfig.users.find(el=> el.name == this.selected.context["user"]);

				axios.post(`${this.backendUrl()}/api/clusters/${this.selected.name}`,
						{
							"cluster": Object.assign({}, cluster.cluster),
							"user": Object.assign({}, user.user)
						}
				).then( resp => {
					this.contexts(resp.data.contexts);
					this.toast("Add a cluster.. OK", "success");
				}).catch(e => { this.msghttp(e);});
			} catch (e) { this.toast("Add a cluster failed","danger") }

		}
	}
}
</script>
