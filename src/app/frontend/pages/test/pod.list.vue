<template>
<!-- content-wrapper -->
<div class="content-wrapper">

	<div class="content-header">
		<div class="container-fluid">
			<c-navigator group="Workload"></c-navigator>
			<div class="row mb-2">
				<div class="col-sm-2"><h1 class="m-0 text-dark">Pods</h1></div>
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
								<nuxt-link :to="{ path:'/view', query:{ context: currentContext(), group: 'Test', crd: 'Pod', name: data.item.name, url: `api/v1/namespaces/${data.item.namespace}/pods/${data.item.name}`, preurl: $router.currentRoute.fullPath}}">{{ data.value }}</nuxt-link>
							</template>
							<template v-slot:cell(labels)="data">
								<ul class="list-unstyled mb-0">
									<li v-for="(value, name) in data.item.labels" v-bind:key="name"><span class="badge badge-secondary font-weight-light text-sm mb-1">{{ name }}:{{ value }}</span></li>
								</ul>
							</template>
						</b-table>
					</div>
					<b-pagination v-model="currentPage" :per-page="$config.itemsPerPage" :total-rows="totalItems" size="sm" align="center"></b-pagination>
				</div>
			</div>
		</div><!-- //GRID-->
	</div>
	</section>

</div>
</template>
<script>
import axios		from "axios"
import VueNavigator from "@/components/navigator"

export default {
	components: {
		"c-navigator": { extends: VueNavigator }
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
				{ text: "Completed", value: "Completed" },
				{ text: "Failed", value: "Failed" },
				{ text: "Unknowen", value: "Unknowen" }
			],
			keyword: "",
			filterOn: ["name"],
			fields: [
				{ key: "name", label: "이름", sortable: true },
				{ key: "namespace", label: "네임스페이스", sortable: true  },
				{ key: "ready", label: "Ready", sortable: true  },
				{ key: "status", label: "상태", sortable: true  },
				{ key: "restartCount", label: "재시작", sortable: true  },
				{ key: "creationTimestamp", label: "생성시간", sortable: true },
				{ key: "nodeName", label: "노드", sortable: true  },
				{ key: "usageCpu", label: "CPU 사용량", sortable: true  },
				{ key: "usageMemory", label: "MEMORY 사용량", sortable: true  },
			],
			isBusy: false,
			metricsItems: [],
			origin: [],
			items: [],
			containerStatuses: [],
			currentPage: 1,
			totalItems: 0
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
			let selectedStatus = this.selectedStatus;
			this.items = this.origin.filter(el => {
				return (selectedStatus.length == 0) || selectedStatus.includes(el.status);
			});
			this.totalItems = this.items.length;
			this.currentPage = 1
		},
		// 조회
		query_All() {
			this.isBusy = true;
			this.getMetrics();
			axios.get(`${this.backendUrl()}/raw/clusters/${this.currentContext()}/api/v1/namespaces/${this.$data.selectedNamespace}/pods`)
				.then((resp) => {
					this.items = [];
					resp.data.items.forEach(el => {
						this.items.push({
							name: el.metadata.name,
							namespace: el.metadata.namespace,
							ready: this.toReady(el.status, el.spec.containers.length),
							status: this.toStatus(el.metadata.deletionTimestamp, el.status),
							restartCount: el.status.containerStatuses ? el.status.containerStatuses.map(el => el.restartCount).reduce((accumulator, currentValue) => accumulator + currentValue) : 0,
							creationTimestamp: this.$root.getElapsedTime(el.metadata.creationTimestamp),
							nodeName: el.spec.nodeName ? el.spec.nodeName : "<none>",
							usageCpu: this.toUsageHandler('cpu', el.metadata.name),
							usageMemory: this.toUsageHandler('memory', el.metadata.name),
						});
					});
					this.origin = this.items;
					this.onFiltered(this.items);
				})
				.catch(e => { this.msghttp(e);})
				.finally(()=> { this.isBusy = false;});
		},
		/**
		 * 메트릭 값을 리소스별로 단위 계산 후 반환 한다.
		 * 
		 * @param {string} resource 구분자 cpu/memory
		 * @param {string} podName 구분자 pod 이름으로 구분 한다.
		 * @return {string} 리소스 합산 값의 단위를 추가 해서 반환 한다.
		 */
		toUsageHandler(resource, podName) {
			if (!this.metricsItems.find(x => x.metadata.name === podName)) return "<none>"
			else {
				const cpuSize = ["n", "m"]
				const memorySize = ["Ki", "Mi", "Gi", "Ti", "Pi", "Ei"]
				const decimals = 2

				const calculator = this.metricsItems.find(x => x.metadata.name === podName)
					.containers.map(x => x.usage[resource])
						.map(el => {
							let calculate = 0

							if (resource == "cpu") {
								cpuSize.forEach((size, index) => {
									if (el.indexOf(size) > -1) {
										calculate= Number(el.split(size)[0]) * (index > 0 ? Math.pow(1000, index + 1) : 1)
									}
								}) 
							}
							if (resource == "memory") {
								memorySize.forEach((size, index) => {
									if (el.indexOf(size) > -1) {
										calculate= Number(el.split(size)[0]) * Math.pow(1024, index + 1)
									}
								})
							}
							return calculate
						})
						.reduce((accumulator, currentValue) => {
							return accumulator + currentValue
						})

				return this.$root.getFormatMetrics(resource, calculator, decimals)
			}
		},
		/**
		 * 컨테이너의 ready 수를 반환 한다.
		 * 
		 * @param {object} status 파드의 status 값.
		 * @param {number} containersLength 파드의 컨테이너 수.
		 * @return {string} 컨테이서 총 개수 대비 ready 수 를 반환 한다.
		 */
		toReady(status, containersLength) {
			let ready = ""
			if ( status.containerStatuses ) {
				ready = `${status.containerStatuses.filter(el => el.ready).length}/${containersLength}`
			} else {
				ready = `0/${containersLength}`
			}
			return ready
		},
		/**
		 * 파드의 상태를 반환 한다.
		 * 
		 * @param {date} deletionTimestamp 파드의 삭제 시간 값.
		 * @param {object} statusItems 파드의 status 값.
		 * @return {string} 파드의 상태 값을 반환 한다.
		 */
		toStatus(deletionTimestamp, statusItems) {
			// 삭제
			if (deletionTimestamp) {
				return "Terminating"
			}

			// Pending
			if (!statusItems.containerStatuses) {
				return statusItems.phase
			}

			// [if]: Running, [else]: (CrashRoofBack / Completed / ContainerCreating)
			if(statusItems.containerStatuses.filter(el => el.ready).length === statusItems.containerStatuses.length) {
				const state = Object.keys(statusItems.containerStatuses.find(el => el.ready).state)[0]
				return state.charAt(0).toUpperCase() + state.slice(1)
			}
			else {
				const state = statusItems.containerStatuses.find(el => !el.ready).state
				return state[Object.keys(state)].reason
			}
		},
		/**
		 * 리소스 메트릭 값을 반환 한다.
		 * 
		 * @async
		 * @function getMetrics
		 * @returns {Promise<object>} 리소스(cpu/memory) 메트릭 값을 반환 한다.
		 */
		async getMetrics() {
			this.metricsItems = [];
			let resp = await axios.get(`${this.backendUrl()}/raw/clusters/${this.currentContext()}/apis/metrics.k8s.io/v1beta1/namespaces/${this.$data.selectedNamespace}/pods`)
			resp.data.items.forEach(el => {
				this.metricsItems.push(el)
			});
		},
		onFiltered(filteredItems) {
			let status = { running:0, pending:0, failed:0, terminating:0, crashLoopBackOff:0, crashLoopBackOff:0, completed:0, failed:0, unknowen:0 }

			filteredItems.forEach(el=> {
				if(el.status == "Running") status.running++;
				if(el.status == "Pending") status.pending++;
				if(el.status == "Terminating") status.terminating++;
				if(el.status == "CrashLoopBackOff") status.crashLoopBackOff++;
				if(el.status == "Completed") status.completed++;
				if(el.status == "Failed") status.failed++;
				if(el.status == "Unknowen") status.unknowen++;
			});

			this.optionsStatus[0].text = status.running >0 ? `Running (${status.running})`: "Running";
			this.optionsStatus[1].text = status.pending >0 ? `Pending (${status.pending})`: "Pending";
			this.optionsStatus[2].text = status.terminating >0 ? `Terminating (${status.terminating})`: "Terminating";
			this.optionsStatus[3].text = status.crashLoopBackOff >0 ? `CrashLoopBackOff (${status.crashLoopBackOff})`: "CrashLoopBackOff";
			this.optionsStatus[4].text = status.completed >0 ? `Completed (${status.completed})`: "Completed";
			this.optionsStatus[5].text = status.failed >0 ? `Failed (${status.failed})`: "Failed";
			this.optionsStatus[6].text = status.unknowen >0 ? `Unknowen (${status.unknowen})`: "Unknowen";

			this.totalItems = filteredItems.length;
			this.currentPage = 1
		},
	},
	beforeDestroy(){
		this.$nuxt.$off('navbar-context-selected')
	}
}
</script>
<style>label {font-weight: 500;}</style>
