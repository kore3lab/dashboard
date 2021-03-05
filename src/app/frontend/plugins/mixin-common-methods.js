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
    confirm(msg, callback) {
      this.$bvModal
        .msgBoxConfirm(msg, {
          title: "",
          variant: "info",
          buttonSize: "sm",
          footerClass: "p-1",
        })
        .then(callback)
        .catch((err) => {});
    },
    msghttp(error) {
      if (
        error.response &&
        error.response.data &&
        error.response.data.message
      ) {
        this.toast(error.response.data.message, "warning");
      } else {
        this.toast(error.message, "danger");
      }
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
    /**
     * timestamp를 day,hour,minute,second로 구분 봔환함
     *
     * @param {date} timestamp 변환할 date 값
     * @return {string} timestamp의 day/hour/minute/second 값으로 변환하여 반환함
     */
    getElapsedTime(timestamp) {
      const dt = Date.parse(timestamp);
      const elapsedTime = new Date() - dt;

      const second = Math.floor((elapsedTime % (1000 * 60)) / 1000);
      const minute = Math.floor((elapsedTime % (1000 * 60 * 60)) / (1000 * 60));
      const hour = Math.floor(
        (elapsedTime % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60)
      );
      // const day = Math.floor((elapsedTime % (1000*60*60*24*30)) / (1000*60*60*24))
      const days = Math.floor(elapsedTime / (1000 * 60 * 60 * 24));
      // const month = Math.floor((elapsedTime % (1000*60*60*24*30*12)) / (1000*60*60*24*30))
      // const year = Math.floor(elapsedTime  / (1000*60*60*24*30*12))

      let str = "";
      // if(year > 0) str += `${year}y`
      // if(month > 0) str += `${month}m`
      // if(day > 0) str += `${day}d`
      if (days > 0) {
        str += `${days}d`;
        if (days >= 10) return str;
      }
      if (hour > 0) {
        str += `${hour}h`;
        if (days < 10 && days > 0) return str;
      }
      if (minute > 0) {
        if (days > 0 || hour > 0) return str;
        str += `${minute}m`;
      }

      if (second > 0) {
        if (hour > 0 || minute > 9) return str;
        str += `${second}s`;
      }
      return str;
    },
    /**
     * 리소스 메트릭 수집 값을 Formatting 한다.
     * CPU / milliCPU 또는 CPU의 1/1000 단위로 처리 할 때 "표현"됩니다. 나노 코어 / 나노 CPU는 CPU의 1/1000000000 (10 억분의 1)입니다.
     * Memory 단위는 이진접두어를 사용 한다.
     *
     * @param {string} resource 구분자 cpu/memory
     * @param {number} metrics 리소스 사용량 합계 값
     * @param {number} decimals 소수점 아래 자릿수. 기본값은 2 이다
     * @return {string} 리소스 합산 값의 단위를 변환한다.
     */
    getFormatMetrics(resource, metrics, decimals = 2) {
      if (metrics === 0) return 0;

      const k = 1024;
      const dm = decimals < 0 ? 0 : decimals;
      const memorySize = ["", "Ki", "Mi", "Gi", "Ti", "Pi", "Ei"];
      // const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];

      const i = Math.floor(Math.log(metrics) / Math.log(k));

      if (resource == "cpu") {
        if (metrics / 1000000000 > 1)
          return `${parseFloat((metrics / 1000000000).toFixed(dm))}core`;
        if (metrics / 1000000 > 1)
          return `${parseFloat((metrics / 1000000).toFixed(dm))}m`;
        return `${parseFloat(metrics)}n`;
      }

      return `${parseFloat((metrics / Math.pow(k, i)).toFixed(dm))}${
        memorySize[i]
      }`;
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
    // Get currentContext's namespaces
    namespaces(_) {
      if (_) this.$store.commit("setNamespaces", _);
      else return this.$store.getters["getNamespaces"];
    },
    // Get contexts
    contexts(_) {
      if (_) this.$store.commit("setContexts", _);
      else return this.$store.getters["getContexts"];
    },
    // Get currentContext's resources
    resources(_) {
      if (_) this.$store.commit("setResources", _);
      else return this.$store.getters["getResources"];
    },
    // Get a currentContext
    currentContext(_) {
      if (_) this.$store.commit("setCurrentContext", _);
      else return this.$store.getters["getCurrentContext"];
    },
  },
});
