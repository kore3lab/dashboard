export const state = () => ({
  currentContext: "", // currentContext
  contexts: [], // context list
  namespaces: [], // currentContext namespace 리스트
  selectNamespace: "", // select Namespace
  isNamespace: 'no',
  resources: [], // currentContext resource 리스트
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
  setResources(state, resources) {
    state.resources = resources;
  },
  setSelectNamespace(state, selectNamespace) {
    state.selectNamespace = selectNamespace;
  },
  setIsNamespace(state, isNamespace) {
    state.isNamespace = isNamespace
  },
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
  getResources(state) {
    return state.resources;
  },
  getSelectNamespace(state) {
    return state.selectNamespace;
  },
  getIsNamespace(state) {
    return state.isNamespace
  },
};
