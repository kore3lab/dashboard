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
						<a href="#" class="nav-link"><i class="nav-icon fas fa-cloud"></i><p>Cluster</p></a>
						<ul class="nav nav-treeview">
							<li class="nav-item small"><nuxt-link to="/cluster/node.list"  class="nav-link"><i class="nav-icon text-sm mr-0">N</i><p>odes</p></nuxt-link></li>
							<li class="nav-item small"><nuxt-link to="/cluster/topology"  class="nav-link"><i class="nav-icon text-sm mr-0">T</i><p>opology</p></nuxt-link></li>
							<li class="nav-item small"><nuxt-link to="/cluster/namespace.list"  class="nav-link"><i class="nav-icon text-sm mr-0">N</i><p>amespaces</p></nuxt-link></li>
						</ul>
					</li>
					<li class="nav-item menu-open">
						<a href="pages/widgets.html" class="nav-link"><i class="nav-icon fas fa-server"></i><p>Workload</p></a>
						<ul class="nav nav-treeview">
							<li class="nav-item small"><nuxt-link to="/workload/pod.list"  class="nav-link"><i class="nav-icon text-sm mr-0">P</i><p>ods</p></nuxt-link></li>
							<li class="nav-item small"><nuxt-link to="/workload/deployment.list"  class="nav-link"><i class="nav-icon text-sm mr-0">D</i><p>eployments</p></nuxt-link></li>
							<li class="nav-item small"><nuxt-link to="/workload/statefulset.list"  class="nav-link"><i class="nav-icon text-sm mr-0">S</i><p>tateful Sets</p></nuxt-link></li>
							<li class="nav-item small"><nuxt-link to="/workload/cornjob.list"  class="nav-link"><i class="nav-icon text-sm mr-0">C</i><p>ron Jobs</p></nuxt-link></li>
							<li class="nav-item small"><nuxt-link to="/workload/job.list"  class="nav-link"><i class="nav-icon text-sm mr-0">J</i><p>obs</p></nuxt-link></li>
							<li class="nav-item small"><nuxt-link to="/workload/daemonset.list"  class="nav-link"><i class="nav-icon text-sm mr-0">D</i><p>aemon Sets</p></nuxt-link></li>
							<li class="nav-item small"><nuxt-link to="/workload/replicaset.list"  class="nav-link"><i class="nav-icon text-sm mr-0">R</i><p>eplica Sets</p></nuxt-link></li>
						</ul>
					</li>
					<li class="nav-item menu-open">
						<a href="pages/widgets.html" class="nav-link"><i class="nav-icon fas fa-network-wired"></i><p>Networking</p></a>
						<ul class="nav nav-treeview">
							<li class="nav-item small"><nuxt-link to="/networking/service.list"  class="nav-link"><i class="nav-icon text-sm mr-0">S</i><p>ervices</p></nuxt-link></li>
							<li class="nav-item small"><nuxt-link to="/networking/ingress.list"  class="nav-link"><i class="nav-icon text-sm mr-0">I</i><p>ngresses</p></nuxt-link></li>
							<li class="nav-item small"><nuxt-link to="/networking/endpoint.list"  class="nav-link"><i class="nav-icon text-sm mr-0">E</i><p>ndpoints</p></nuxt-link></li>
							<li class="nav-item small"><nuxt-link to="/networking/networkpolicy.list"  class="nav-link"><i class="nav-icon text-sm mr-0">N</i><p>etwork Policies</p></nuxt-link></li>
						</ul>
					</li>
					<li class="nav-item menu-open">
						<a href="pages/widgets.html" class="nav-link"><i class="nav-icon fas fa-hdd"></i><p>Storage</p></a>
						<ul class="nav nav-treeview">
							<li class="nav-item small"><nuxt-link to="/storage/pvc.list"  class="nav-link"><i class="nav-icon text-sm mr-0">P</i><p>ersistent Volume Claims</p></nuxt-link></li>
							<li class="nav-item small"><nuxt-link to="/storage/pv.list"  class="nav-link"><i class="nav-icon text-sm mr-0">P</i><p>ersistent Volumes</p></nuxt-link></li>
							<li class="nav-item small"><nuxt-link to="/storage/storageclass.list"  class="nav-link"><i class="nav-icon text-sm mr-0">S</i><p>troage Classes</p></nuxt-link></li>
						</ul>
					</li>
					<li class="nav-item menu-open">
						<a href="pages/widgets.html" class="nav-link"><i class="nav-icon fas fa-cog"></i><p>Configuration</p></a>
						<ul class="nav nav-treeview">
							<li class="nav-item small"><nuxt-link to="/configuration/configmap.list"  class="nav-link"><i class="nav-icon text-sm mr-0">C</i><p>onfig Maps</p></nuxt-link></li>
							<li class="nav-item small"><nuxt-link to="/configuration/secret.list"  class="nav-link"><i class="nav-icon text-sm mr-0">S</i><p>ecrets</p></nuxt-link></li>
							<li class="nav-item small"><nuxt-link to="/configuration/resourcequota.list"  class="nav-link"><i class="nav-icon text-sm mr-0">R</i><p>esource Quotas</p></nuxt-link></li>
							<li class="nav-item small"><nuxt-link to="/configuration/limitrange.list"  class="nav-link"><i class="nav-icon text-sm mr-0">L</i><p>imit Ranges</p></nuxt-link></li>
							<li class="nav-item small"><nuxt-link to="/configuration/hpa.list"  class="nav-link"><i class="nav-icon text-sm mr-0">H</i><p>orizontal Pod Autoscalers</p></nuxt-link></li>
							<li class="nav-item small"><nuxt-link to="/configuration/poddisruptionbudget.list"  class="nav-link"><i class="nav-icon text-sm mr-0">P</i><p>od Disruption Budgets</p></nuxt-link></li>
							<li class="nav-item small"><nuxt-link to="/configuration/customresourcedefinitions.list"  class="nav-link"><i class="nav-icon text-sm mr-0">C</i><p>ustom Resource Definitions</p></nuxt-link></li>
						</ul>
					</li>
					<li class="nav-item menu-open">
						<a href="pages/widgets.html" class="nav-link"><i class="nav-icon fas fa-users-cog"></i><p>Administrator</p></a>
						<ul class="nav nav-treeview">
							<li class="nav-item small"><nuxt-link to="/administrator/serviceaccount.list"  class="nav-link"><i class="nav-icon text-sm mr-0">S</i><p>ervice Accounts</p></nuxt-link></li>
							<li class="nav-item small"><nuxt-link to="/administrator/clusterrole.list"  class="nav-link"><i class="nav-icon text-sm mr-0">C</i><p>luster Roles</p></nuxt-link></li>
							<li class="nav-item small"><nuxt-link to="/administrator/clusterrolebinding.list"  class="nav-link"><i class="nav-icon text-sm mr-0">C</i><p>luster Role Bindings</p></nuxt-link></li>
							<li class="nav-item small"><nuxt-link to="/administrator/role.list"  class="nav-link"><i class="nav-icon text-sm mr-0">R</i><p>oles</p></nuxt-link></li>
							<li class="nav-item small"><nuxt-link to="/administrator/rolebinding.list"  class="nav-link"><i class="nav-icon text-sm mr-0">R</i><p>ole Bindings</p></nuxt-link></li>
						</ul>
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
}
</script>