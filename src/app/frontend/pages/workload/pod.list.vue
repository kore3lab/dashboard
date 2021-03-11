<template>
  <!-- content-wrapper -->
  <div class="content-wrapper">

    <div class="content-header">
      <div class="container-fluid">
        <c-navigator group="Workload"></c-navigator>
        <div class="row mb-2">
          <div class="col-sm-2"><h1 class="m-0 text-dark"><span class="badge badge-info mr-2">P</span>Pods</h1></div>
          <!-- 검색 (namespace) -->
          <div class="col-sm-2">
            <b-form-select v-model="selectedNamespace" :options="namespaces()" size="sm" @input="query_All"></b-form-select>
          </div><!--//END -->
          <!-- 검색 (검색어) -->
          <div class="col-sm-2 float-left">
            <div class="input-group input-group-sm" >
              <b-form-input id="txtKeyword" v-model="keyword" class="form-control float-right" placeholder="Search"></b-form-input>
              <div class="input-group-append">
                <button type="submit" class="btn btn-default" @click="query_All"><i class="fas fa-search"></i></button>
              </div>
            </div>
          </div><!--//END -->
          <!-- 버튼 -->
          <div class="col-sm-6 text-right">
            <b-button variant="primary" size="sm" @click="$router.push(`/create?context=${currentContext()}&group=Workload&crd=Pod`)">Create</b-button>
          </div><!--//END -->
        </div>
      </div>
    </div>

    <section class="content">
      <div class="container-fluid">
        <!-- 검색 (상태) -->
        <div class="row mb-2">
          <div class="col-11">
            <b-form-group class="mb-0 font-weight-light">
              <button type="submit" class="btn btn-default btn-sm" @click="query_All">All</button>
              <b-form-checkbox-group v-model="selectedStatus" :options="optionsStatus" button-variant="light"  font="light" buttons size="sm" @input="onChangeStatus"></b-form-checkbox-group>
            </b-form-group>
          </div>
          <div class="col-1 text-right "><span class="text-sm align-middle">Total : {{ totalItems }}</span></div>
        </div><!--//END -->
        <!-- GRID-->
        <div class="row">
          <div class="col-12">
            <div class="card">
              <div class="card-body table-responsive p-0">
                <b-table id="list" hover :items="items" :fields="fields" :filter="keyword" :filter-included-fields="filterOn" @filtered="onFiltered" :current-page="currentPage" :per-page="$config.itemsPerPage" :busy="isBusy" class="text-sm">
                  <template #table-busy>
                    <div class="text-center text-success" style="margin:150px 0">
                      <b-spinner type="grow" variant="success" class="align-middle mr-2"></b-spinner>
                      <span class="align-middle text-lg">Loading...</span>
                    </div>
                  </template>
                  <template v-slot:cell(name)="data">
                    <a href="#" @click="sidebar={visible:true, name:data.item.name, crd:'Pod', src:`${getApiUrl('','pods',data.item.namespace)}/${data.item.name}`}">{{ data.value }}</a>
                  </template>
                  <template v-slot:cell(status)="data">
                    <div class="list-unstyled mb-0" v-if="data.item.status.value">
                      <span v-bind:class="data.item.status.style">{{ data.item.status.value }}</span>
                    </div>
                  </template>
                  <template v-slot:cell(containers)="data">
                    <div class="list-unstyled mb-0 float-left" v-if="data.item.containers.containerStatuses.length > 0">
                      <span v-for="(containerStatuses, idx) in data.item.containers.containerStatuses" v-bind:key="idx" v-bind:class="containerStatuses.style" class="badge font-weight-light text-sm ml-1"> {{" "}}</span>
                    </div>
                    <div class="list-unstyled mb-0 ml-0 float-left" v-if="data.item.containers.initContainerStatuses.length > 0">
                      <span v-for="(initContainerStatuses, idx) in data.item.containers.initContainerStatuses" v-bind:key="idx" v-bind:class="initContainerStatuses.style" class="badge font-weight-light text-sm ml-1">{{" "}}</span>
                    </div>
                  </template>
                  <template v-slot:cell(controller)="data">
                    <a href="#" @click="sidebar={visible:true, name:data.value.name, crd:data.value.spaceKind, group:data.value.gr, src:`${getApiUrl(data.value.group,data.value.rs,data.item.namespace)}/${data.value.name}`}">{{ data.value.kind }}</a>
                  </template>
                  <template v-slot:cell(node)="data">
                    <a href="#" @click="sidebar={visible:true, name:data.value.name, crd:'Node', src:`${getApiUrl(data.value.group,data.value.rs)}/${data.value.name}`}">{{ data.value.name }}</a>
                  </template>
                </b-table>
              </div>
              <b-pagination v-model="currentPage" :per-page="$config.itemsPerPage" :total-rows="totalItems" size="sm" align="center"></b-pagination>
            </div>
          </div>
        </div><!-- //GRID-->
      </div>
    </section>
    <b-sidebar v-model="sidebar.visible" width="50em" right shadow no-header>
      <c-view :crd="sidebar.crd" :group="sidebar.crd" :name="sidebar.name" :url="sidebar.src" @delete="query_All()" @close="sidebar.visible=false"/>
    </b-sidebar>
  </div>
</template>
<script>
import axios		from "axios"
import VueNavigator from "@/components/navigator"
import VueView from "@/pages/view";
export default {
  components: {
    "c-navigator": { extends: VueNavigator },
    "c-view": { extends: VueView }
  },
  data() {
    return {
      selectedNamespace: "",
      selectedStatus: [],
      optionsStatus: [
        { text: "Running", value: "Running" },
        { text: "Pending", value: "Pending" },
        { text: "Terminating", value: "Terminating" },
        { text: "CrashLoopBackOff", value: "CrashLoopBackOff" },
        { text: "ImagePullBackOff", value: "ImagePullBackOff" },
        { text: "Completed", value: "Completed" },
        { text: "Failed", value: "Failed" },
        { text: "Unknown", value: "Unknown" }
      ],
      keyword: "",
      filterOn: ["name"],
      fields: [
        { key: "name", label: "Name", sortable: true },
        { key: "namespace", label: "Namespace", sortable: true  },
        { key: "ready", label: "Ready", sortable: true  },
        { key: "containers", label: "Containers", sortable: true  },
        { key: "restartCount", label: "Restart", sortable: true  },
        { key: "controller", label: "Controlled By", sortable: true  },
        { key: "node", label: "Node", sortable: true  },
        { key: "qos", label: "QoS", sortable: true },
        { key: "creationTimestamp", label: "Age", sortable: true },
        { key: "status", label: "Status", sortable: true  },
      ],
      isBusy: false,
      metricsItems: [],
      origin: [],
      items: [],
      containerStatuses: [],
      currentPage: 1,
      totalItems: 0,
      sidebar: {
        visible: false,
        name: "",
        src: "",
      },
    }
  },
  layout: "default",
  created() {
    this.$nuxt.$on("navbar-context-selected", (ctx) => this.query_All() );
    if(this.currentContext()) this.$nuxt.$emit("navbar-context-selected");
  },
  methods: {
    //  status 필터링
    onChangeStatus() {
      var selectedStatus = this.selectedStatus;
      this.items = this.origin.filter(el => {
        return (selectedStatus.length === 0) || selectedStatus.includes(el.status.value);
      });
      this.totalItems = this.items.length;
      this.currentPage = 1
    },
    // 조회
    query_All() {
      this.isBusy = true;
      this.loadMetrics();
      axios.get(this.getApiUrl("","pods",this.selectedNamespace))
          .then((resp) => {
            this.items = [];
            resp.data.items.forEach(el => {
              this.items.push({
                name: el.metadata.name,
                namespace: el.metadata.namespace,
                ready: this.toReady(el.status, el.spec),
                containers: this.toContainers(el.status),
                restartCount: el.status.containerStatuses ? el.status.containerStatuses.map(x => x.restartCount).reduce((accumulator, currentValue) => accumulator + currentValue, 0) : 0,
                controller: this.getController(el),
                status: this.toStatus(el.metadata.deletionTimestamp, el.status),
                creationTimestamp: this.$root.getElapsedTime(el.metadata.creationTimestamp),
                node: this.getNode(el),
                qos: el.status.qosClass,
              });
            });
            this.origin = this.items;
            this.onFiltered(this.items);
          })
          .catch(e => { this.msghttp(e);})
          .finally(()=> { this.isBusy = false;});
    },
    getNode(el) {
      return {
        "name" : el.spec.nodeName ? el.spec.nodeName : "",
        "group" : "",
        "rs" : "nodes"
      }
    },
    getController(el) {
      let version
      let group
      let gr = ""
      if ( el.metadata.ownerReferences ) {
        version = el.metadata.ownerReferences[0].apiVersion.split('/')
        if (version.length>1) {
          group = version[0]
        }
        if (el.metadata.ownerReferences[0].kind === "Node") {
          gr = "Cluster"
        } else {
          gr = "Workload"
        }
        return {
          "name" : el.metadata.ownerReferences[0].name,
          "group" : group ||"",
          "kind" : el.metadata.ownerReferences[0].kind,
          "rs" : (el.metadata.ownerReferences[0].kind).toLowerCase()+'s',
          "spaceKind" : this.onKind(el.metadata.ownerReferences[0].kind),
          "gr" : gr,
        }
      } else {
        return ""
      }
    },
    onKind(kind) {
      if (kind === "StatefulSet") {
        return "Stateful Set"
      } else if (kind === "CronJob") {
        return "Cron Job"
      } else if (kind === "DaemonSet") {
        return "Daemon Set"
      } else if (kind === "ReplicaSet") {
        return "Replica Set"
      } else if (kind === "ReplicationController") {
        return "Replication Controller"
      } else {
        return kind
      }
    },
    // check Ready pod
    toReady(status, spec) {
      let containersReady = 0
      let containersLength = 0
      if ( spec.containers ) containersLength = spec.containers.length
      if ( status.containerStatuses ) containersReady = status.containerStatuses.filter(el => el.ready).length
      return `${containersReady}/${containersLength}`
    },
    // check pod's containers
    toContainers(status) {
      let initContainerStatuses = []
      let containerStatuses = []
      let style = ""
      if ( status.initContainerStatuses ) {
        initContainerStatuses = status.initContainerStatuses.map(x => {
          if ( !x.ready ) style = "badge-warning"
          else style = "badge-secondary"
          return Object.assign(x, {"style": style})
        })

      }
      if ( status.containerStatuses ) {
        containerStatuses = status.containerStatuses.map(x => {
          if ( !x.ready ) {
            style = "badge-warning"
            if ( x.state[Object.keys(x.state)].reason === "Completed") style = "badge-secondary"
          }
          else style = "badge-success"
          return Object.assign(x, {"style": style})
        })
      }
      return {
        "initContainerStatuses": initContainerStatuses,
        "containerStatuses": containerStatuses,
      }
    },
    // pod status check
    toStatus(deletionTimestamp, status) {
      // 삭제
      if (deletionTimestamp) {
        return {
          "value": "Terminating",
          "style": "text-secondary",
        }
      }

      // Pending
      if (!status.containerStatuses) {
        if(status.phase === "Failed") {
          return {
            "value": status.phase,
            "style": "text-danger",
          }
        } else {
          return {
            "value": status.phase,
            "style": "text-warning",
          }
        }
      }

      // [if]: Running, [else]: (CrashRoofBack / Completed / ContainerCreating)
      if(status.containerStatuses.filter(el => el.ready).length === status.containerStatuses.length) {
        const state = Object.keys(status.containerStatuses.find(el => el.ready).state)[0]
        return {
          "value": state.charAt(0).toUpperCase() + state.slice(1),
          "style": "text-success",
        }
      }
      else {
        const state = status.containerStatuses.find(el => !el.ready).state
        let style = "text-secondary"
        if ( state[Object.keys(state)].reason === "Completed") style = "text-success"
        return {
          "value": state[Object.keys(state)].reason,
          "style": style,
        }
      }
    },
    // load pod's Metrics
    async loadMetrics() {
      this.metricsItems = [];
      let resp = await axios.get(this.getApiUrl("metrics.k8s.io","pods",this.selectedNamespace))
      if (!resp) return
      this.metricsItems = resp.data.items
    },
    onFiltered(filteredItems) {
      let status = { running:0, pending:0, failed:0, terminating:0, crashLoopBackOff:0, imagePullBackOff:0, completed:0, unknown:0 }

      filteredItems.forEach(el=> {
        if(el.status.value === "Running") status.running++;
        if(el.status.value === "Pending") status.pending++;
        if(el.status.value === "Terminating") status.terminating++;
        if(el.status.value === "CrashLoopBackOff") status.crashLoopBackOff++;
        if(el.status.value === "ImagePullBackOff") status.imagePullBackOff++;
        if(el.status.value === "Completed") status.completed++;
        if(el.status.value === "Failed") status.failed++;
        if(el.status.value === "Unknown") status.unknown++;
      });

      this.optionsStatus[0].text = status.running >0 ? `Running (${status.running})`: "Running";
      this.optionsStatus[1].text = status.pending >0 ? `Pending (${status.pending})`: "Pending";
      this.optionsStatus[2].text = status.terminating >0 ? `Terminating (${status.terminating})`: "Terminating";
      this.optionsStatus[3].text = status.crashLoopBackOff >0 ? `CrashLoopBackOff (${status.crashLoopBackOff})`: "CrashLoopBackOff";
      this.optionsStatus[4].text = status.imagePullBackOff >0 ? `ImagePullBackOff (${status.imagePullBackOff})`: "ImagePullBackOff";
      this.optionsStatus[5].text = status.completed >0 ? `Completed (${status.completed})`: "Completed";
      this.optionsStatus[6].text = status.failed >0 ? `Failed (${status.failed})`: "Failed";
      this.optionsStatus[7].text = status.unknown >0 ? `Unknown (${status.unknown})`: "Unknown";

      this.totalItems = filteredItems.length;
      this.currentPage = 1
    },
  },
  beforeDestroy(){
    this.$nuxt.$off('navbar-context-selected')
  }
}
</script>
<style scoped></style>
