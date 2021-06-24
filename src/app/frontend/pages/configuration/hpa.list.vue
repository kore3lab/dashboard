<template>
	<div class="content-wrapper">
		<div class="content-header">
			<div class="container-fluid">
				<c-navigator group="Configuration"></c-navigator>
				<div class="row mb-2">
					<!-- title & search -->
					<div class="col-sm"><h1 class="m-0 text-dark"><span class="badge badge-info mr-2">H</span>Horizontal Pod Autoscalers</h1></div>
					<div class="col-sm-2"><b-form-select v-model="selectedNamespace" :options="namespaces()" size="sm" @input="query_All(); selectNamespace(selectedNamespace);"></b-form-select></div>
					<div class="col-sm-2 float-left">
						<div class="input-group input-group-sm" >
							<b-form-input id="txtKeyword" v-model="keyword" class="form-control float-right" placeholder="Search"></b-form-input>
							<div class="input-group-append"><button type="submit" class="btn btn-default" @click="query_All"><i class="fas fa-search"></i></button></div>
						</div>
					</div>
					<!-- button -->
					<div class="col-sm-1 text-right">
						<b-button variant="primary" size="sm" @click="$router.push(`/create?context=${currentContext()}&group=Configuration&crd=HPA`)">Create</b-button>
					</div>
				</div>
			</div>
		</div>

		<section class="content">
			<div class="container-fluid">
				<!-- total count & items per page  -->
				<div class="d-flex flex-row-reverse">
					<div class="p-2">
						<b-form inline>
							<c-colums-selector name="grdSheet1" v-model="fields" :fields="fieldsAll" ></c-colums-selector>
							<i class="text-secondary ml-2 mr-2">|</i>
							<b-form-select size="sm" :options="this.var('ITEMS_PER_PAGE')" v-model="itemsPerPage"></b-form-select>
							<span class="text-sm align-middle ml-2">Total : {{ totalItems }}</span>
						</b-form>
					</div>
				</div>
				<!-- GRID-->
				<div class="row">
					<div class="col-12">
						<div class="card">
							<div class="card-body table-responsive p-0">
								<b-table hover selectable show-empty select-mode="single" @row-selected="onRowSelected" @sort-changed="currentPage=1" ref="grdSheet1" :items="items" :fields="fields" :filter="keyword" :filter-included-fields="filterOn" @filtered="onFiltered" :current-page="currentPage" :per-page="itemsPerPage" :busy="isBusy" class="text-sm">
									<template #table-busy>
										<div class="text-center text-success lh-vh-50">
											<b-spinner type="grow" variant="success" class="align-middle mr-2"></b-spinner>
											<span class="text-lg align-middle">Loading...</span>
										</div>
									</template>
									<template v-slot:cell(labels)="data">
										<ul class="list-unstyled mb-0">
											<li v-for="(value, name) in data.item.labels" v-bind:key="name"><span class="badge badge-secondary font-weight-light text-sm mb-1">{{ name }}={{ value }}</span></li>
										</ul>
									</template>
									<template v-slot:cell(target)="data">
										<a href="#" @click="viewModel=getViewLink(data.value.group,data.value.rs,data.item.namespace, data.value.name); isShowSidebar=true;">{{ data.value.kind }}/{{ data.value.name }}</a>
									</template>
									<template v-slot:cell(status)="data">
										<span v-for="(status, idx) in data.item.status" v-bind:key="idx" v-bind:class="status.style" class=" text-sm ml-1">{{ status.type }}</span>
									</template>
								</b-table>
							</div>
							<b-pagination v-model="currentPage" :per-page="itemsPerPage" :total-rows="totalItems" size="sm" align="center"></b-pagination>
						</div>
					</div>
				</div><!-- //GRID-->
			</div>
		</section>
		<b-sidebar v-model="isShowSidebar" width="50em" @hidden="$refs.grdSheet1.clearSelected()" right shadow no-header>
			<c-view v-model="viewModel" @delete="query_All()" @close="isShowSidebar=false"/>
		</b-sidebar>
	</div>
</template>
<script>
import VueNavigator			from "@/components/navigator"
import VueColumsSelector	from "@/components/columnsSelector"
import VueView				from "@/pages/view";

export default {
	components: {
		"c-navigator": { extends: VueNavigator },
		"c-colums-selector": { extends: VueColumsSelector},
		"c-view": { extends: VueView }
	},
	data() {
		return {
			selectedNamespace: "",
			keyword: "",
			filterOn: ["name"],
			fields: [],
			fieldsAll: [
				{ key: "name", label: "Name", sortable: true },
				{ key: "namespace", label: "Namespace", sortable: true  },
				{ key: "metrics", label: "Metrics", sortable: true  },
				{ key: "minpods", label: "Min Pods", sortable: true  },
				{ key: "maxpods", label: "Max Pods", sortable: true  },
				{ key: "replicas", label: "Replicas", sortable: true  },
				{ key: "target", label: "Target", sortable: true },
				{ key: "creationTimestamp", label: "Age", sortable: true, formatter: this.getElapsedTime },
			],
			isBusy: false,
			items: [],
			itemsPerPage: this.$storage.global.get("itemsPerPage",10),
			currentPage: 1,
			totalItems: 0,
			isShowSidebar: false,
			viewModel:{},
		}
	},
	watch: {
		itemsPerPage(n) {
			this.$storage.global.set("itemsPerPage",n)
		}
	},
	layout: "default",
	created() {
		this.$nuxt.$on("navbar-context-selected", (ctx) => {
			this.selectedNamespace = this.selectNamespace()
			this.query_All()
		});
		if(this.currentContext()) this.$nuxt.$emit("navbar-context-selected");
	},
	methods: {
		onRowSelected(items) {
			this.isShowSidebar = (items && items.length > 0)
			if (this.isShowSidebar) this.viewModel = this.getViewLink('autoscaling', 'horizontalpodautoscalers', items[0].namespace, items[0].name)
		},
		// 조회
		query_All() {
			this.isBusy = true;
			this.$axios.get(this.getApiUrl("autoscaling","horizontalpodautoscalers",this.selectedNamespace))
					.then((resp) => {
						this.items = [];
						resp.data.items.forEach(el => {
							this.items.push({
								name: el.metadata.name,
								namespace: el.metadata.namespace,
								metrics: this.getMetrics(el),
								minpods: el.spec.minReplicas,
								maxpods: el.spec.maxReplicas,
								replicas: el.status.currentReplicas,
								target: this.getTarget(el.spec.scaleTargetRef),
								creationTimestamp: el.metadata.creationTimestamp
							});
						});
						this.onFiltered(this.items);
					})
					.catch(e => { this.msghttp(e);})
					.finally(()=> { this.isBusy = false;});
		},
		onFiltered(filteredItems) {
			this.totalItems = filteredItems.length;
			this.currentPage = 1
		},
		getStatus(conditions) {
			let list = [];
			if(conditions) {
				for (let i = 0; i < conditions.length; i++) {
					if (conditions[i].status === "True") {
						if (conditions[i].type === "AbleToScale") {
							list.push({
								type: conditions[i].type,
								style: "text-success"
							})
						} else if (conditions[i].type === "ScalingActive") {
							list.push({
								type: conditions[i].type,
								style: "text-primary"
							})
						} else {
							list.push({
								type: conditions[i].type,
								style: "text-danger"
							})
						}
					}
				}
			}
			return list
		},


		getMetrics(el) {
			let current = ""
			let target = ""
			if (el.status.currentCPUUtilizationPercentage) {
				current = el.status.currentCPUUtilizationPercentage + "%"
			} else {
				current = "0"
			}
			if (el.spec.targetCPUUtilizationPercentage) {
				target = el.spec.targetCPUUtilizationPercentage + "%"
			} else {
				target = "0"
			}
			return current+' / '+target
			// v2beta2 에서 사용 가능.
			//   let specType,statusType,lowerType;
			//   let targetVal,currentVal;
			//   let currentValue = "unknown";
			//   let targetValue = "unknown";
			//   if (el.spec.metrics) {
			//     specType = el.spec.metrics[0].type
			//   }
			//   if (el.status.currentMetrics) {
			//     statusType = el.status.currentMetrics[0].type
			//   }
			//   if (specType) {
			//     lowerType = specType.toLowerCase()
			//     targetVal = el.spec.metrics[0][lowerType]
			//   }
			//   if (statusType) {
			//     lowerType = statusType.toLowerCase()
			//     currentVal = el.status.currentMetrics[0][lowerType]
			//   }
			//   if (targetVal) {
			//     targetValue = targetVal.target.averageUtilization || targetVal.target.averageValue || targetVal.target.value;
			//     if (targetVal.target.averageUtilization) {
			//       targetValue += "%";
			//     }
			//   }
			//   if (currentVal) {
			//     currentValue = currentVal.current.averageUtilization || currentVal.current.averageValue || currentVal.current.value;
			//     if (currentVal.current.averageUtilization) {
			//       currentValue += "%";
			//     }
			//   }
			//   if (el.spec.metrics && el.spec.metrics.length > 1) {
			//     targetValue += ` + ${el.spec.metrics.length-1} more...`
			//   }
			//   return `${currentValue} / ${targetValue}`
		},
		getTarget(ref) {
			let group = ""
			let version = ref.apiVersion.split('/')
			if (version.length>1) {
				group = version[0]
			}
			return {
				"name": ref.name,
				"group": group,
				"kind": ref.kind,
				"rs": ref.kind.toLowerCase()+'s'
			}
		},
	},
	beforeDestroy(){
		this.$nuxt.$off('navbar-context-selected')
	}
}
</script>
<style scoped>label {font-weight: 500;}</style>
