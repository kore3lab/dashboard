import Vue from "vue";

const CUSTOM_VIEWS	= {
	//Workload
	cronjobs: "workload/cronjob.view",
	daemonsets: "workload/daemonset.view",
	deployments: "workload/deployment.view",
	jobs: "workload/job.view",
	pods: "workload/pod.view",
	replicasets: "workload/replicaset.view",
	replicationcontrollers: "workload/replicationcontroller.view",
	statefulsets: "workload/statefulset.view",

	//Storage
	persistentvolumes: "storage/pv.view",
	persistentvolumeclaims: "storage/pvc.view",
	storageclasses: "storage/storageclass.view",

	//Network
	endpoints: "networking/endpoint.view",
	ingresses: "networking/ingress.view",
	networkpolicies: "networking/networkpolicy.view",
	services: "networking/service.view",

	//Cluster
	nodes: "cluster/node.view",
	namespaces: "cluster/namespace.view",

	//Configuration
	configmaps: "configuration/configmap.view",
	customresourcedefinitions: "configuration/customresourcedefinitions.view",
	horizontalpodautoscalers: "configuration/hpa.view",
	poddisruptionbudgets: "configuration/poddisruptionbudget.view",
	secrets: "configuration/secret.view",
	limitranges: "configuration/limitrange.view",
	resourcequotas: "configuration/resourcequota.view",

	//Administrator
	serviceaccounts: "administrator/serviceaccount.view",
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
