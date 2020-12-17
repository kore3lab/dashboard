import Vue from "vue";

Vue.mixin({
  methods: {
    toast(msg, variant) {
      if (!variant) variant = "info";
      this.$bvToast.toast(msg, {
        title: "Kore3",
        variant: variant,
        autoHideDelay: 4000,
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
    namespaces(_) {
      if (_) this.$store.commit("setNamespaces", _);
      else return this.$store.getters["getNamespaces"];
    },
    contexts(_) {
      if (_) this.$store.commit("setContexts", _);
      else return this.$store.getters["getContexts"];
    },
    currentContext(_) {
      if (_) this.$store.commit("setCurrentContext", _);
      else return this.$store.getters["getCurrentContext"];
    },
  },
});
