<template>
	<aside class="main-sidebar sidebar-dark-primary elevation-4">

		<c-context />

		<nuxt-link to="/" class="brand-link align-bottom">
			<img src="/favicon.svg" class="brand-image">
			<span class="brand-text text-info font-weight-bold align-bottom"> K</span><span class="align-bottom">ore</span><sup><i class="fas fa-sm fa-cubes ml-1 text-warning"></i></sup><sub class="text-sm ml-1">dashboard</sub>
		</nuxt-link>

		<!-- Sidebar -->
		<div class="sidebar">
			<!-- Sidebar Menu -->
			<nav class="mt-2">
				<ul class="nav nav-pills nav-sidebar flex-column nav-child-indent nav-compact" data-widget="treeview" role="menu" data-accordion="false">
					<li class="nav-item has-treeview menu-open">
						<a href="#" class="nav-link" v-on:click="isVisible.cluster=!isVisible.cluster"><i class="nav-icon fas fa-cloud"></i><p>Cluster<i v-bind:class="isVisible.cluster?'fa-angle-left':'fa-angle-up'" class="right fas"></i></p></a>
						<b-collapse v-model="isVisible.cluster">
							<ul class="nav d-block nav-treeview">
								<li class="nav-item small" v-on:click="loads()"><nuxt-link to="/cluster/node.list"  class="nav-link"><i class="nav-icon text-sm mr-0">N</i><p>odes</p></nuxt-link></li>
								<li class="nav-item small" v-on:click="loads()"><nuxt-link to="/cluster/topology"  class="nav-link"><i class="nav-icon text-sm mr-0">T</i><p>opology</p></nuxt-link></li>
								<li class="nav-item small" v-on:click="loads()"><nuxt-link to="/cluster/namespace.list"  class="nav-link"><i class="nav-icon text-sm mr-0">N</i><p>amespaces</p></nuxt-link></li>
								<li class="nav-item small" v-on:click="loads()"><nuxt-link to="/cluster/event.list"  class="nav-link"><i class="nav-icon text-sm mr-0">E</i><p>vents</p></nuxt-link></li>
							</ul>
						</b-collapse>
					</li>
					<li class="nav-item menu-open">
						<a href="#" class="nav-link" v-on:click="isVisible.workload=!isVisible.workload"><i class="nav-icon fas fa-server"></i><p>Workload<i v-bind:class="isVisible.workload?'fa-angle-left':'fa-angle-up'" class="right fas"></i></p></a>
						<b-collapse v-model="isVisible.workload">
							<ul class="nav d-block nav-treeview">
								<!--								<li class="nav-item small"><nuxt-link to="/workload/overview"  class="nav-link"><i class="nav-icon text-sm mr-0">O</i><p>verview</p></nuxt-link></li>-->
								<li class="nav-item small" v-on:click="loads()"><nuxt-link to="/workload/pod.list"  class="nav-link"><i class="nav-icon text-sm mr-0">P</i><p>ods</p></nuxt-link></li>
								<li class="nav-item small" v-on:click="loads()"><nuxt-link to="/workload/deployment.list"  class="nav-link"><i class="nav-icon text-sm mr-0">D</i><p>eployments</p></nuxt-link></li>
								<li class="nav-item small" v-on:click="loads()"><nuxt-link to="/workload/statefulset.list"  class="nav-link"><i class="nav-icon text-sm mr-0">S</i><p>tateful Sets</p></nuxt-link></li>
								<li class="nav-item small" v-on:click="loads()"><nuxt-link to="/workload/cornjob.list"  class="nav-link"><i class="nav-icon text-sm mr-0">C</i><p>ron Jobs</p></nuxt-link></li>
								<li class="nav-item small" v-on:click="loads()"><nuxt-link to="/workload/job.list"  class="nav-link"><i class="nav-icon text-sm mr-0">J</i><p>obs</p></nuxt-link></li>
								<li class="nav-item small" v-on:click="loads()"><nuxt-link to="/workload/daemonset.list"  class="nav-link"><i class="nav-icon text-sm mr-0">D</i><p>aemon Sets</p></nuxt-link></li>
								<li class="nav-item small" v-on:click="loads()"><nuxt-link to="/workload/replicaset.list"  class="nav-link"><i class="nav-icon text-sm mr-0">R</i><p>eplica Sets</p></nuxt-link></li>
							</ul>
						</b-collapse>
					</li>
					<li class="nav-item menu-open">
						<a href="#" class="nav-link" v-on:click="isVisible.networking=!isVisible.networking"><i class="nav-icon fas fa-network-wired"></i><p>Networking<i v-bind:class="isVisible.networking?'fa-angle-left':'fa-angle-up'" class="right fas"></i></p></a>
						<b-collapse v-model="isVisible.networking">
							<ul class="nav d-block nav-treeview">
								<li class="nav-item small" v-on:click="loads()"><nuxt-link to="/networking/service.list"  class="nav-link"><i class="nav-icon text-sm mr-0">S</i><p>ervices</p></nuxt-link></li>
								<li class="nav-item small" v-on:click="loads()"><nuxt-link to="/networking/ingress.list"  class="nav-link"><i class="nav-icon text-sm mr-0">I</i><p>ngresses</p></nuxt-link></li>
								<li class="nav-item small" v-on:click="loads()"><nuxt-link to="/networking/endpoint.list"  class="nav-link"><i class="nav-icon text-sm mr-0">E</i><p>ndpoints</p></nuxt-link></li>
								<li class="nav-item small" v-on:click="loads()"><nuxt-link to="/networking/networkpolicy.list"  class="nav-link"><i class="nav-icon text-sm mr-0">N</i><p>etwork Policies</p></nuxt-link></li>
							</ul>
						</b-collapse>
					</li>
					<li class="nav-item menu-open">
						<a href="#" class="nav-link" v-on:click="isVisible.storage=!isVisible.storage"><i class="nav-icon fas fa-hdd"></i><p>Storage<i v-bind:class="isVisible.storage?'fa-angle-left':'fa-angle-up'" class="right fas"></i></p></a>
						<b-collapse v-model="isVisible.storage">
							<ul class="nav d-block nav-treeview">
								<li class="nav-item small" v-on:click="loads()"><nuxt-link to="/storage/pvc.list"  class="nav-link"><i class="nav-icon text-sm mr-0">P</i><p>ersistent Volume Claims</p></nuxt-link></li>
								<li class="nav-item small" v-on:click="loads()"><nuxt-link to="/storage/pv.list"  class="nav-link"><i class="nav-icon text-sm mr-0">P</i><p>ersistent Volumes</p></nuxt-link></li>
								<li class="nav-item small" v-on:click="loads()"><nuxt-link to="/storage/storageclass.list"  class="nav-link"><i class="nav-icon text-sm mr-0">S</i><p>troage Classes</p></nuxt-link></li>
							</ul>
						</b-collapse>
					</li>
					<li class="nav-item menu-open">
						<a href="#" class="nav-link" v-on:click="isVisible.configuration=!isVisible.configuration"><i class="nav-icon fas fa-cog"></i><p>Configuration<i v-bind:class="isVisible.configuration?'fa-angle-left':'fa-angle-up'" class="right fas"></i></p></a>
						<b-collapse v-model="isVisible.configuration">
							<ul class="nav d-block nav-treeview">
								<li class="nav-item small" v-on:click="loads()"><nuxt-link to="/configuration/configmap.list"  class="nav-link"><i class="nav-icon text-sm mr-0">C</i><p>onfig Maps</p></nuxt-link></li>
								<li class="nav-item small" v-on:click="loads()"><nuxt-link to="/configuration/secret.list"  class="nav-link"><i class="nav-icon text-sm mr-0">S</i><p>ecrets</p></nuxt-link></li>
								<li class="nav-item small" v-on:click="loads()"><nuxt-link to="/configuration/resourcequota.list"  class="nav-link"><i class="nav-icon text-sm mr-0">R</i><p>esource Quotas</p></nuxt-link></li>
								<li class="nav-item small" v-on:click="loads()"><nuxt-link to="/configuration/limitrange.list"  class="nav-link"><i class="nav-icon text-sm mr-0">L</i><p>imit Ranges</p></nuxt-link></li>
								<li class="nav-item small" v-on:click="loads()"><nuxt-link to="/configuration/hpa.list"  class="nav-link"><i class="nav-icon text-sm mr-0">H</i><p>orizontal Pod Autoscalers</p></nuxt-link></li>
								<li class="nav-item small" v-on:click="loads()"><nuxt-link to="/configuration/poddisruptionbudget.list"  class="nav-link"><i class="nav-icon text-sm mr-0">P</i><p>od Disruption Budgets</p></nuxt-link></li>
							</ul>
						</b-collapse>
					</li>
					<li class="nav-item menu-open">
						<a href="#" class="nav-link" v-on:click="isVisible.administrator=!isVisible.administrator"><i class="nav-icon fas fa-users-cog"></i><p>Administrator<i v-bind:class="isVisible.administrator?'fa-angle-left':'fa-angle-up'" class="right fas"></i></p></a>
						<b-collapse v-model="isVisible.administrator">
							<ul class="nav d-block nav-treeview">
								<li class="nav-item small" v-on:click="loads()"><nuxt-link to="/administrator/serviceaccount.list"  class="nav-link"><i class="nav-icon text-sm mr-0">S</i><p>ervice Accounts</p></nuxt-link></li>
								<li class="nav-item small" v-on:click="loads()"><nuxt-link to="/administrator/clusterrole.list"  class="nav-link"><i class="nav-icon text-sm mr-0">C</i><p>luster Roles</p></nuxt-link></li>
								<li class="nav-item small" v-on:click="loads()"><nuxt-link to="/administrator/clusterrolebinding.list"  class="nav-link"><i class="nav-icon text-sm mr-0">C</i><p>luster Role Bindings</p></nuxt-link></li>
								<li class="nav-item small" v-on:click="loads()"><nuxt-link to="/administrator/role.list"  class="nav-link"><i class="nav-icon text-sm mr-0">R</i><p>oles</p></nuxt-link></li>
								<li class="nav-item small" v-on:click="loads()"><nuxt-link to="/administrator/rolebinding.list"  class="nav-link"><i class="nav-icon text-sm mr-0">R</i><p>ole Bindings</p></nuxt-link></li>
							</ul>
						</b-collapse>
					</li>
					<li class="nav-item menu-open mw-100">
						<a href="#" class="nav-link" v-on:click="isVisible.cr=!isVisible.cr"><i class="nav-icon fas fa-puzzle-piece"></i><p>Custom Resource<i v-bind:class="isVisible.cr?'fa-angle-left':'fa-angle-up'" class="right fas"></i></p></a>
						<b-collapse v-model="isVisible.cr">
							<ul class="nav d-block nav-treeview">
								<li class="nav-item small" v-on:click="loads()"><nuxt-link to="/customresource/customresourcedefinitions.list"  class="nav-link"><i class="nav-icon text-sm mr-0">D</i><p>efinitions</p></nuxt-link></li>
								<ul class="nav nav-pills nav-sidebar flex-column nav-child-indent nav-compact" data-widget="treeview" role="menu" data-accordion="false">
									<li v-for="(val, idx) in crList" v-bind:key="idx" class="nav-item small mw-100">
										<a href="#" class="nav-link" v-on:click="checkStatus(val)"><p class="text-truncate mw-100 pr-2">{{ Object.keys(val)[0] }} {{Object.keys(val)[0]}}<i v-bind:class="getStatus(val) ? 'fa-angle-down':'fa-angle-left'" class="right fas"></i></p></a>
										<b-collapse v-model="crVisible[Object.keys(val)[0]]">
											<ul v-for="(v, i) in val" v-bind:key="i" class="nav d-block nav-treeview">
												<li v-for="(a, b) in v" class="nav-item" v-bind:key="a.name" v-on:click="sideCRD()"><nuxt-link :to="{path: '/customresource/customresource.list', query: { isSide: true, gV: a.groupVersion, k: a.kind, n: a.name, ns: a.namespaced}}"  class="nav-link"><i class="nav-icon text-sm mr-0"></i><p>{{ (b && b.length > 25) ? b.substring(0,25)+'...' : b }}</p></nuxt-link></li>
											</ul>
										</b-collapse>
									</li>
								</ul>

							</ul>
						</b-collapse>
					</li>
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
				cr : true,
			},
			isWorkload : true,
			crList: [],
			crVisible: {},
		}
	},
	created() {
		this.$nuxt.$on("crList_up",(list) => {
			this.crList = list
			this.crList.forEach(el => {
				this.crVisible[Object.keys(el)[0]] = false
			})
		})
		this.$nuxt.$on('crCol_up',(name) => {
			this.crVisible[name] = true
			this.crVisible = Object.assign({},this.crVisible)
		})
	},
	methods: {
		checkStatus(data) {
			let key = Object.keys(data)[0]
			this.$set(this.crVisible,key,!this.crVisible[key])
			this.crVisible = Object.assign({},this.crVisible)
		},
		getStatus(val) {
			return this.crVisible[Object.keys(val)[0]]
		},
		sideCRD() {
			this.$nuxt.$emit('sideCRD_click','')
		},
		loads() {
			if(this.currentContext()) this.$nuxt.$emit('navbar-context-selected','')
		},
	},
	beforeDestroy() {
		this.$nuxt.$off("crList_up")
		this.$nuxt.$off("crCol_up")
	}
}
</script>