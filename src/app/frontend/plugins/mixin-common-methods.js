import Vue from "vue";

Vue.mixin({
  methods: {
    toast(msg, variant) {
      if (!variant) variant = "info";
      this.$bvToast.toast(msg, {
        title: "",
        noCloseButton: true,
        variant: variant,
        autoHideDelay: 4000,
      });
    },
    mesbox(msg) {
      this.$bvModal.msgBoxOk(msg, {
        title: "",
        variant: "info",
        buttonSize: "sm",
        footerClass: "p-1",
      });
    },
    backendUrl() {
      return `${location.protocol}//${location.hostname}:${this.$config.backendPort}`;
    },
    dashboardUrl() {
      return `${location.protocol}//${location.hostname}:${this.$config.dashboardPort}`;
    },
    kialiRootUrl() {
      return `${location.protocol}//${location.hostname}:${this.$config.kialiPort}`;
    },
    getTimestampString(timestamp) {
      let dt = Date.parse(timestamp);
      let seconds = Math.floor((new Date() - dt) / 1000);

      var interval = seconds / 31536000;
      if (interval > 1) return Math.floor(interval) + " years";

      interval = seconds / 2592000;
      if (interval > 1) return Math.floor(interval) + " months";

      interval = seconds / 86400;
      if (interval > 1) return Math.floor(interval) + " days";

      interval = seconds / 3600;
      if (interval > 1) return Math.floor(interval) + " hours";
      interval = seconds / 60;
      if (interval > 1) return Math.floor(interval) + " minutes";

      return Math.floor(seconds) + " seconds";
    },
    // namespace 리스트 조회
    namespaces(_) {
      if (_) this.$store.commit("setNamespaces", _);
      else return this.$store.getters["getNamespaces"];
    },
    // context (cluster) 리스트 조회
    contexts(_) {
      if (_) this.$store.commit("setContexts", _);
      else return this.$store.getters["getContexts"];
    },
    // 현재 선택된 context (cluster) 조회
    currentContext(_) {
      if (_) this.$store.commit("setCurrentContext", _);
      else return this.$store.getters["getCurrentContext"];
    },
  },
});
