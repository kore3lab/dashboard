export const state = () => ({
	currentContext: "",		// currentContext
	contexts: [], 			// context list
	namespaces: [],			// currentContext namespace 리스트
	labelSelector: [],		// labelSelector
	selectNamespace: "",	// select Namespace
	resources: [],			// currentContext resource 리스트
	statusbar: {
		kubernetesVersion: "",
		platform: "",
		message: ""
	}
});
export const mutations = {
	setCurrentContext(state, ctx) {
		state.currentContext = ctx;
	},
	setContexts(state, ctxs) {
		state.contexts = ctxs;
	},
	setNamespaces(state, namespaces) {
		state.namespaces = namespaces;
	},
	setLabelSelector(state, labelSelector) {
		state.labelSelector = labelSelector;
	},
	setSelectNamespace(state, selectNamespace) {
		state.selectNamespace = selectNamespace;
	},
	setResources(state, resources) {
		state.resources = resources;
	},
	setStatusbar(state, statusbar) {
		state.statusbar = statusbar;
	}
};
export const getters = {
	getCurrentContext(state) {
		return state.currentContext;
	},
	getContexts(state) {
		return state.contexts;
	},
	getNamespaces(state) {
		return state.namespaces;
	},
	getLabelSelector(state) {
		return state.labelSelector;
	},
	getSelectNamespace(state) {
		return state.selectNamespace;
	},
	getResources(state) {
		return state.resources;
	},
	getStatusbar(state) {
		return state.statusbar;
	}
};
