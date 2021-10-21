<template>
	<aside class="sidebar-dark-primary elevation-4">

		<c-context class="aside-contexts"/>

		<nuxt-link to="/" class="brand-link align-bottom">
			<img src="/favicon.svg" class="brand-image">
			<span class="brand-text text-info font-weight-bold align-bottom"> K</span><span class="align-bottom">ore</span><sup><i class="fas fa-sm fa-cubes ml-1 text-warning"></i></sup><sub class="text-sm ml-1">dashboard</sub>
		</nuxt-link>

		<!-- Sidebar -->
		<div class="aside-menus">
			<!-- Sidebar Menu -->
			<nav class="mt-2">
				<ul class="nav nav-pills nav-sidebar flex-column nav-child-indent nav-compact" data-widget="treeview" role="menu" data-accordion="false">
					<li class="nav-item has-treeview menu-open">
						<a href="#" class="nav-link" v-on:click="isVisible.cluster=!isVisible.cluster"><i class="nav-icon fas fa-cloud"></i><p>Cluster<i v-bind:class="isVisible.cluster?'fa-angle-left':'fa-angle-up'" class="right fas"></i></p></a>
						<b-collapse v-model="isVisible.cluster">
							<ul class="nav d-block nav-treeview">
								<li class="nav-item small"><nuxt-link to="/cluster/node.list"  class="nav-link"><i class="nav-icon text-sm mr-0">N</i><p>odes</p></nuxt-link></li>
								<li class="nav-item small"><nuxt-link to="/cluster/topology"  class="nav-link"><i class="nav-icon text-sm mr-0">T</i><p>opology</p></nuxt-link></li>
								<li class="nav-item small"><nuxt-link to="/cluster/namespace.list"  class="nav-link"><i class="nav-icon text-sm mr-0">N</i><p>amespaces</p></nuxt-link></li>
								<li class="nav-item small"><nuxt-link to="/cluster/event.list"  class="nav-link"><i class="nav-icon text-sm mr-0">E</i><p>vents</p></nuxt-link></li>
							</ul>
						</b-collapse>
					</li>
					<li class="nav-item menu-open">
						<a href="#" class="nav-link" v-on:click="isVisible.workload=!isVisible.workload"><i class="nav-icon fas fa-server"></i><p>Workload<i v-bind:class="isVisible.workload?'fa-angle-left':'fa-angle-up'" class="right fas"></i></p></a>
						<b-collapse v-model="isVisible.workload">
							<ul class="nav d-block nav-treeview">
								<!--								<li class="nav-item small"><nuxt-link to="/workload/overview"  class="nav-link"><i class="nav-icon text-sm mr-0">O</i><p>verview</p></nuxt-link></li>-->
								<li class="nav-item small"><nuxt-link to="/workload/pod.list"  class="nav-link"><i class="nav-icon text-sm mr-0">P</i><p>ods</p></nuxt-link></li>
								<li class="nav-item small"><nuxt-link to="/workload/deployment.list"  class="nav-link"><i class="nav-icon text-sm mr-0">D</i><p>eployments</p></nuxt-link></li>
								<li class="nav-item small"><nuxt-link to="/workload/statefulset.list"  class="nav-link"><i class="nav-icon text-sm mr-0">S</i><p>tateful Sets</p></nuxt-link></li>
								<li class="nav-item small"><nuxt-link to="/workload/cornjob.list"  class="nav-link"><i class="nav-icon text-sm mr-0">C</i><p>ron Jobs</p></nuxt-link></li>
								<li class="nav-item small"><nuxt-link to="/workload/job.list"  class="nav-link"><i class="nav-icon text-sm mr-0">J</i><p>obs</p></nuxt-link></li>
								<li class="nav-item small"><nuxt-link to="/workload/daemonset.list"  class="nav-link"><i class="nav-icon text-sm mr-0">D</i><p>aemon Sets</p></nuxt-link></li>
								<li class="nav-item small"><nuxt-link to="/workload/replicaset.list"  class="nav-link"><i class="nav-icon text-sm mr-0">R</i><p>eplica Sets</p></nuxt-link></li>
							</ul>
						</b-collapse>
					</li>
					<li class="nav-item menu-open">
						<a href="#" class="nav-link" v-on:click="isVisible.networking=!isVisible.networking"><i class="nav-icon fas fa-network-wired"></i><p>Networking<i v-bind:class="isVisible.networking?'fa-angle-left':'fa-angle-up'" class="right fas"></i></p></a>
						<b-collapse v-model="isVisible.networking">
							<ul class="nav d-block nav-treeview">
								<li class="nav-item small"><nuxt-link to="/networking/service.list"  class="nav-link"><i class="nav-icon text-sm mr-0">S</i><p>ervices</p></nuxt-link></li>
								<li class="nav-item small"><nuxt-link to="/networking/ingress.list"  class="nav-link"><i class="nav-icon text-sm mr-0">I</i><p>ngresses</p></nuxt-link></li>
								<li class="nav-item small"><nuxt-link to="/networking/endpoint.list"  class="nav-link"><i class="nav-icon text-sm mr-0">E</i><p>ndpoints</p></nuxt-link></li>
								<li class="nav-item small"><nuxt-link to="/networking/networkpolicy.list"  class="nav-link"><i class="nav-icon text-sm mr-0">N</i><p>etwork Policies</p></nuxt-link></li>
							</ul>
						</b-collapse>
					</li>
					<li class="nav-item menu-open">
						<a href="#" class="nav-link" v-on:click="isVisible.storage=!isVisible.storage"><i class="nav-icon fas fa-hdd"></i><p>Storage<i v-bind:class="isVisible.storage?'fa-angle-left':'fa-angle-up'" class="right fas"></i></p></a>
						<b-collapse v-model="isVisible.storage">
							<ul class="nav d-block nav-treeview">
								<li class="nav-item small"><nuxt-link to="/storage/pvc.list"  class="nav-link"><i class="nav-icon text-sm mr-0">P</i><p>ersistent Volume Claims</p></nuxt-link></li>
								<li class="nav-item small"><nuxt-link to="/storage/pv.list"  class="nav-link"><i class="nav-icon text-sm mr-0">P</i><p>ersistent Volumes</p></nuxt-link></li>
								<li class="nav-item small"><nuxt-link to="/storage/storageclass.list"  class="nav-link"><i class="nav-icon text-sm mr-0">S</i><p>troage Classes</p></nuxt-link></li>
							</ul>
						</b-collapse>
					</li>
					<li class="nav-item menu-open">
						<a href="#" class="nav-link" v-on:click="isVisible.configuration=!isVisible.configuration"><i class="nav-icon fas fa-cog"></i><p>Configuration<i v-bind:class="isVisible.configuration?'fa-angle-left':'fa-angle-up'" class="right fas"></i></p></a>
						<b-collapse v-model="isVisible.configuration">
							<ul class="nav d-block nav-treeview">
								<li class="nav-item small"><nuxt-link to="/configuration/configmap.list"  class="nav-link"><i class="nav-icon text-sm mr-0">C</i><p>onfig Maps</p></nuxt-link></li>
								<li class="nav-item small"><nuxt-link to="/configuration/secret.list"  class="nav-link"><i class="nav-icon text-sm mr-0">S</i><p>ecrets</p></nuxt-link></li>
								<li class="nav-item small"><nuxt-link to="/configuration/resourcequota.list"  class="nav-link"><i class="nav-icon text-sm mr-0">R</i><p>esource Quotas</p></nuxt-link></li>
								<li class="nav-item small"><nuxt-link to="/configuration/limitrange.list"  class="nav-link"><i class="nav-icon text-sm mr-0">L</i><p>imit Ranges</p></nuxt-link></li>
								<li class="nav-item small"><nuxt-link to="/configuration/hpa.list"  class="nav-link"><i class="nav-icon text-sm mr-0">H</i><p>orizontal Pod Autoscalers</p></nuxt-link></li>
								<li class="nav-item small"><nuxt-link to="/configuration/poddisruptionbudget.list"  class="nav-link"><i class="nav-icon text-sm mr-0">P</i><p>od Disruption Budgets</p></nuxt-link></li>
							</ul>
						</b-collapse>
					</li>
					<li class="nav-item menu-open">
						<a href="#" class="nav-link" v-on:click="isVisible.administrator=!isVisible.administrator"><i class="nav-icon fas fa-users-cog"></i><p>Administrator<i v-bind:class="isVisible.administrator?'fa-angle-left':'fa-angle-up'" class="right fas"></i></p></a>
						<b-collapse v-model="isVisible.administrator">
							<ul class="nav d-block nav-treeview">
								<li class="nav-item small"><nuxt-link to="/administrator/serviceaccount.list"  class="nav-link"><i class="nav-icon text-sm mr-0">S</i><p>ervice Accounts</p></nuxt-link></li>
								<li class="nav-item small"><nuxt-link to="/administrator/clusterrole.list"  class="nav-link"><i class="nav-icon text-sm mr-0">C</i><p>luster Roles</p></nuxt-link></li>
								<li class="nav-item small"><nuxt-link to="/administrator/clusterrolebinding.list"  class="nav-link"><i class="nav-icon text-sm mr-0">C</i><p>luster Role Bindings</p></nuxt-link></li>
								<li class="nav-item small"><nuxt-link to="/administrator/role.list"  class="nav-link"><i class="nav-icon text-sm mr-0">R</i><p>oles</p></nuxt-link></li>
								<li class="nav-item small"><nuxt-link to="/administrator/rolebinding.list"  class="nav-link"><i class="nav-icon text-sm mr-0">R</i><p>ole Bindings</p></nuxt-link></li>
							</ul>
						</b-collapse>
					</li>
					<!--custom resource -->
					<li class="nav-item menu-open mw-100">
						<a href="#" class="nav-link" v-on:click="isVisible.crd=!isVisible.crd"><i class="nav-icon fas fa-puzzle-piece"></i><p>Custom Resource<i v-bind:class="isVisible.crd?'fa-angle-left':'fa-angle-up'" class="right fas"></i></p></a>
						<b-collapse v-model="isVisible.crd">
							<!--custom resource : group -->
							<ul class="nav d-block nav-treeview">
								<li class="nav-item small"><nuxt-link to="/customresource/customresourcedefinitions.list"  class="nav-link"><i class="nav-icon text-sm mr-0">D</i><p>efinitions</p></nuxt-link></li>
								<ul class="nav nav-pills nav-sidebar flex-column nav-child-indent nav-compact" data-widget="treeview" role="menu" data-accordion="false">
									<li v-for="(group, nm) in crdList" :key="nm" class="nav-item small mw-100">
										<a href="#" class="nav-link" v-on:click="isVisible.crdGroup[nm] = !isVisible.crdGroup[nm]"><p class="text-truncate mw-100 pr-2">{{ nm }}<i v-bind:class="1==1 ? 'fa-angle-down':'fa-angle-left'" class="right fas"></i></p></a>
										<!--custom resource : group > crd -->
										<b-collapse v-model="isVisible.crdGroup[nm]">
											<ul class="nav d-block nav-treeview">
												<li v-for="(crd, crdnm) in group"  :key="crdnm" class="nav-item" >
													<nuxt-link :to="toCRDLink(crd)"  class="nav-link pl-2 pt-0 pb-0" replace>
														<p class="text-truncate mw-100"><b-icon icon="dot"></b-icon>{{ crd.name }}</p>
													</nuxt-link>
												</li>
											</ul>
										</b-collapse>
									</li>
								</ul>
							</ul>
						</b-collapse>
					</li><!--// custom resource -->

				</ul>
			</nav><!-- /.sidebar-menu -->
		</div>
	</aside>
</template>
<script>
import Context	from './context.vue'
export default {
	components: {
		"c-context": Context,
	},
	data() {
		return {
			isVisible: {
				cluster: true,
				workload: true,
				networking: true,
				storage: true,
				configuration: true,
				administrator: true,
				crd : true,
				crdGroup: {}
			},
			crdList: {}
		}
	},
	created() {
		this.$nuxt.$on("aside-context-selected", (_) => {
			this.crdList = {}
			for(let r in this.resources()) {
				if( !r.endsWith("k8s.io") && r.indexOf(".") > 0 ) {
					this.crdList[r] = this.resources()[r];
				}
			};
			this.isVisible.crdGroup = {};
		});
	},
	methods: {
		toCRDLink(crd) {
			let a = crd.groupVersion.split("/");
			//return {path: `/customresource/customresource.list?crd=${crd.name}.${a[0]}&version=${a[1]}`}
			let query = (a.length > 1) ? {crd: `${crd.name}.${a[0]}`, version: a[1]} : {crd: `${crd.name}.${a[0]}`, version: ""}
			return {path: '/customresource/customresource.list', query: query}
		}
	},
	beforeDestroy(){
		this.$nuxt.$off("aside-context-selected")
	}
}
</script>