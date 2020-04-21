export const state = () => ({
  currentContext: "", // 선택된 cluster context
  contexts: [], // context 리스트
  namespaces: [], // 현재 context 의 namespace 리스트
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
};
