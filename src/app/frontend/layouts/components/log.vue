<template>
  <div class="main-log " v-if="tabs.length > 0">
    <section class="content margin-right">
      <div class="card card-secondary card-outline">
        <b-tabs card grow v-model="active_tab" ref="tabs">
          <!-- Render Tabs, supply a unique `key` to each tab -->
          <b-tab v-for="val in tabs" :key="val.index" class="text-truncate m-n2">
            <template #title>
              {{val.podName}}
              <button type="button" class="btn btn-tool" @click="closeTab(val.containerName, val.podName)" ><i class="fas fa-times"></i></button>
            </template>
            <div class="row mb-1" role="group" >
              <div class="col-xs-1 "> &nbsp; Namespace : </div>
              <div class="col-sm-2 text-truncate float-left"><h5><span class="badge badge-light">{{val.namespace}}</span></h5></div>
              <div class="col-xs-1  float-left"> &nbsp; Pod : </div>
              <div class="col-sm-4 text-truncate float-left max-width:10px;"><h5><span class="badge badge-light">{{val.podName}}</span></h5></div>
              <div class="col-xs-1  float-left"> &nbsp; Container : </div>
              <div class="col-sm-2 float-left"><b-form-select v-model="val.containerName" :options="val.containerList" size="sm" @input="getContainerName(val.containerName, val.index);"></b-form-select></div>
              <br/>
            </div>
            <div class="overflow-auto bg-black px-2 " style="height:15rem;" v-bind:id="val.podName + '_' + val.index" @scroll="handleScroll">
              <div class="d-flex align-items-center justify-content-center" style="height:15rem;" v-if="val.data === '' ">
                <span>
                  There are no logs available for container
                </span>
              </div>
              <div v-else class="row" v-for="d in val.data">
                <span class="text-white"> &nbsp; {{d.replace(/^\d+.*?\s/gm, "")}} </span>
              </div>
            </div>
          </b-tab>
        </b-tabs>
      </div>
    </section>
  </div>
</template>
<script>

export default {
  data() {
    return {
      tabs: [],
      active_tab: 0,
      tabCounter: 0,
      tailLines: 300,
      baseUrl: ''
    }
  },
  created() {
    this.$nuxt.$on('log-data-created', (metadata, containerName, containerList) => {
      if (metadata.name) {
        if (this.tabs.length < 10) {
          this.contentPadding('25rem');

          let flag = true;
          this.tabs.forEach((tabs, index) => {
            if (tabs.containerName === containerName && tabs.podName === metadata.name) {
              flag = false;
              this.active_tab = index;
            }
          })
          if (flag) {
            this.tabs.push({
              index: this.tabCounter,
              containerName: containerName,
              podName: metadata.name,
              namespace: metadata.namespace,
              response: {},
              data: '',
              firstTimestamp: '',
              containerList: containerList
            })
            this.tabCounter++;

            setTimeout(() => {
              this.active_tab = this.tabs[this.tabs.length - 1].index;
            }, 100)

            this.getLogData(metadata.namespace, metadata.name, containerName, '');
          }
        }
      }
    });
  },
  methods: {
    getContainerName(containerName, idx){
      this.tabs.forEach((item, index) => {
        if(item.index == idx){
          if (item.response.hasOwnProperty("read")) {
            item.response.read.cancel();
          }
          item.data = ''
          this.getLogData(item.namespace, item.podName, item.containerName, index);
        }
      })
    },
    getLogData(namespace, podName, containerName, index){
      this.baseUrl = `${this.$config.nodeEnv}` === "development" ? `${location.protocol}//${location.hostname}:${this.$config.backendPort}` : `${location.host}`;
      fetch(this.baseUrl + this.getApiUrl("", "pods", namespace, podName, "follow=1&container=" + containerName + "&tailLines=" + this.tailLines, true)).then(response => {
        const reader = response.body.getReader();
        const Vue = this;
        let buffer = '';
        let tabIndex = this.tabs.length-1;
        if(index !== ''){
          tabIndex = index;
        }
        this.tabs[tabIndex].response = {containerName: containerName, podName: podName, read: reader};
        return reader.read().then(function process(result) {
          if (result.done) {
            return;
          }
          buffer = new TextDecoder().decode(result.value, {stream: true});

          Vue.tabs.forEach(tabs => {
            if (tabs.containerName === containerName && tabs.podName === podName) {
              const logData = buffer.replace(/\[(?:\d{1}|\d{2})m/gi, '').split("\n");

              tabs.data = tabs.data.length === 0 ? logData.filter(item => item.length > 0) : tabs.data.concat(logData).filter(item => item.length > 0);

              setTimeout(() => {
                const logDiv = Vue.$refs.tabs.$refs.content;
                for (let i = 0; i < logDiv.children.length; i++) {
                  const logId = logDiv.children[i].children[1].id;
                  const divInfo = logId.split("_");
                  if(divInfo[1] == tabs.index){
                    logDiv.children[i].children[1].scrollTop = logDiv.children[i].children[1].scrollHeight;
                  }
                }
              });
            }
          })
          return reader.read().then(process);
        })
      })
    },
    handleScroll(evt){
      if(evt.target.scrollTop === 0){
        const logDivId = evt.target.id;
        const divInfo = logDivId.split("_");
        this.tabs.forEach((item, index) => {
          if(divInfo[1] == item.index && item.data !== ''){
            item.firstTimestamp = this.getTimestamps(item.data[0]);
            let tailLine = item.data.length + this.tailLines;
            this.$axios.get(this.getApiUrl("", "pods", item.namespace, item.podName, "follow=0&container=" + item.containerName + "&tailLines=" + tailLine , true))
                .then((resp) => {
                  this.items = [];
                  const prevData = resp.data.split(item.firstTimestamp);
                  const logData = prevData[0].replace(/\[(?:\d{1}|\d{2})m/gi, '').split("\n");
                  this.tabs[index].data = logData.concat(item.data).filter(item => item.length > 0);
                })
                .catch(e => { this.msghttp(e);});
          }
        });
      }
    },
    getTimestamps(logs) {
      return logs.match(/^\d+\S+/gm);
    },
    closeTab(containerName, podName) {
      this.tabs.forEach((item, index) => {
        if (item.containerName === containerName && item.podName === podName) {
          if (item.response.hasOwnProperty("read")) {
            item.response.read.cancel();
          }
          this.tabs.splice(index, 1);
        }
      })
      if(this.tabs.length === 0){
        this.contentPadding('2rem');
      }
    }
  }
}
</script>