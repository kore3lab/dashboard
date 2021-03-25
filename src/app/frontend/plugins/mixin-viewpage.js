import Vue from "vue";

const CUSTOM_VIEWS	= {
	pods: "workload/pod.view",
}
Vue.mixin({
	methods: {
		getViewLink(group, resource, namespace, name) {
			let r = this.resources()[group][resource];
			return {
				title: r.kind,
				name: name, 
				src: CUSTOM_VIEWS[resource], 
				url: this.getApiUrl(group,resource, namespace, name)
			}
		}
	},
});
