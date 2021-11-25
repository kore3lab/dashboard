<template>
	<div class="content-wrapper overflow-hidden">
	<b-overlay :show="showOverlay" rounded="sm">
		<div class="row pt-5">
			<div class="col-md-3"></div>
			<div class="col-md-6 m-3">
				<div class="text-right"><a href="#" @click="$router.back(-1)"><b-icon icon="x" class="h2 text-right border" variant="secondary"></b-icon></a></div>
				<h2>Add Clusters from Kubeconfig</h2>
				<p class="text-left">Add clusters by clicking the Add Cluster button. You'll need to obtain a working kubeconfig for the cluster you want to add. You can either browse it from the file system or paste it as a text from the clipboard. </p>
				<b-tabs content-class="col-md-12" card  v-on:activate-tab="onActiveTab">
					<b-tab title="Select kubeconfig file" active>
						<b-form-group label="Select a file:" label-cols-sm="2" label-size="sm">
							<b-form-file v-model="yamlFile" size="sm" placeholder="Choose a file" drop-placeholder="Drop file here..." @input="onFileSelected"></b-form-file>
						</b-form-group>
					</b-tab>
					<b-tab title="Paste as text">
						<b-form-textarea v-model="yamlText" rows="5" max-rows="6"  @change="onTextChange"></b-form-textarea>
					</b-tab>
				</b-tabs>
				<div class="row">
					<div class="col-4">
						<b-form-select size="sm" v-model="selected" :options="contextList"></b-form-select>
					</div>
					<div class="col-4">
						<b-form-input size="sm" v-model="selected.name"></b-form-input>
					</div>
					<div class="col-2">
						<b-button variant="primary" size="sm" @click="addCluster">Add a cluster</b-button>
					</div>
				</div>
			</div>
			<div class="col-md-3"></div>
		</div>
		</b-overlay>
	</div>
</template>
<script>
import {load as toJSON} from "js-yaml";
export default {
	components: {
	},
	data() {
		return {
			showOverlay: false,
			kubeconfig: "",
			yamlFile: null,
			yamlText: "",
			clusterName: "",
			selected: {name:""},
			contextList: [
				{ value: {name:""}, text: "Please select a context" }
			]
		}
	},
	layout: "default",
	methods: {
		onActiveTab(idx) {
			if (idx===0) this.onFileSelected()
			else this.onTextChange()
		},
		onFileSelected() {
			if(!this.yamlFile) return

			let r = this.doReadConfig;
			let reader = new FileReader();
			reader.onload = function(ev) {
				r(ev.target.result);
			};
			reader.readAsText(this.yamlFile);
		},
		onTextChange() {
			if(!this.yamlText) return

			this.doReadConfig(this.yamlText);
		},
		doReadConfig(json) {
			let conf;
			try {
				this.contextList.splice(1, 2);
				this.selected = {name:""};

				conf = toJSON(json);
				if(conf && typeof conf === "object") {
					conf["contexts"].forEach(el => {
						this.contextList.push( {text: el.name, value: el} );
					});
					this.kubeconfig = conf;
				}
			} catch (e) { this.toast("The selected file is invalid.","danger") }
		},
		addCluster() {
			try {
				let clusterName = this.selected.name;
				if (!clusterName) {
					this.mesbox("Not selected a context");
					return
				}
				let cluster = this.kubeconfig.clusters.find(el=> el.name === this.selected.context["cluster"]);
				let user = this.kubeconfig.users.find(el=> el.name === this.selected.context["user"]);

				this.showOverlay = true;
				this.$axios.post(`/api/contexts/${clusterName}`,{ cluster: Object.assign({}, cluster.cluster), user: Object.assign({}, user.user)})
					.then( resp => {
						if(resp.data && resp.data.contexts) {
							this.contexts(resp.data.contexts);	// set contexts
							// 현재 context 가 목록에 없다면 새로 추가한 cluster 를  context로 선택
							const currentContext = this.currentContext();
							if( !resp.data.contexts.find(d=> {return d==currentContext }) ) {
								this.$nuxt.$emit("set-context-selected", clusterName);
							}
							this.toast("Add a cluster.. OK", "success");
						}
					}).catch(e => { this.msghttp(e) } )
					.finally( () => { this.showOverlay = false; } );
			} catch (ex) {
				this.toast(`Add a cluster failed ${ex}`,"danger") 
			}
		}
	}
}
</script>
