export const state = () => ({
	currentContext: "",		// currentContext
	contexts: [], 			// context list
	namespaces: [],			// currentContext namespace 리스트
	selectNamespace: "",	// select Namespace
	resources: [],			// currentContext resource 리스트
	statusbar: {
		kubernetesVersion: "",
		platform: "",
		message: ""
	},
	contentPadding:'2rem'
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
	setSelectNamespace(state, selectNamespace) {
		state.selectNamespace = selectNamespace;
	},
	setResources(state, resources) {
		state.resources = resources;
	},
	setStatusbar(state, statusbar) {
		state.statusbar = statusbar;
	},
	setContentPadding(state, height){
		state.contentPadding = height;
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
	getSelectNamespace(state) {
		return state.selectNamespace;
	},
	getResources(state) {
		return state.resources;
	},
	getStatusbar(state) {
		return state.statusbar;
	},
	getContentPadding(state){
		return state.contentPadding;
	}
};
