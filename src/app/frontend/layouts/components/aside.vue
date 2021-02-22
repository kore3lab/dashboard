<template>
<aside class="main-sidebar sidebar-dark-primary elevation-4">

	<c-context />

	<nuxt-link to="/" class="brand-link">
		<i class="fab fa-mixcloud mr-0 text-info font-weight-bold"> K</i><span class="brand-text font-weight-bolder">ore<sup><i class="fas fa-sm fa-cubes ml-1 text-warning"></i></sup></span>
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
						<!-- <li class="nav-item small"><nuxt-link to="/cluster/servicemesh"  class="nav-link"><i class="nav-icon text-sm mr-0">S</i><p>erviceMesh</p></nuxt-link></li> -->
						<li class="nav-item small"><nuxt-link to="/cluster/namespace.list"  class="nav-link"><i class="nav-icon text-sm mr-0">N</i><p>amespaces</p></nuxt-link></li>
					</ul>
				</li>
				<li class="nav-item  menu-open">
					<a href="pages/widgets.html" class="nav-link"><i class="nav-icon fas fa-server"></i><p>Workload</p></a>
					<ul class="nav nav-treeview">
						<!-- <li class="nav-item small"><nuxt-link to="/list"  class="nav-link"><i class="nav-icon text-sm mr-0">A</i><p>pplocations</p></nuxt-link></li> -->
						<!-- <li class="nav-item small"><nuxt-link to="/list"  class="nav-link"><i class="nav-icon text-sm mr-0">W</i><p>orkloads</p></nuxt-link></li> -->
						<li class="nav-item small"><nuxt-link to="/workload/pod.list"  class="nav-link"><i class="nav-icon text-sm mr-0">P</i><p>Pods</p></nuxt-link></li>
						<li class="nav-item small"><nuxt-link to="/workload/deployment.list"  class="nav-link"><i class="nav-icon text-sm mr-0">D</i><p>eployments</p></nuxt-link></li>
						<li class="nav-item small"><nuxt-link to="/workload/statefulset.list"  class="nav-link"><i class="nav-icon text-sm mr-0">S</i><p>tateful Sets</p></nuxt-link></li>
						<li class="nav-item small"><nuxt-link to="/workload/cornjob.list"  class="nav-link"><i class="nav-icon text-sm mr-0">C</i><p>ron Jobs</p></nuxt-link></li>
						<li class="nav-item small"><nuxt-link to="/workload/job.list"  class="nav-link"><i class="nav-icon text-sm mr-0">J</i><p>obs</p></nuxt-link></li>
						<li class="nav-item small"><nuxt-link to="/workload/daemonset.list"  class="nav-link"><i class="nav-icon text-sm mr-0">D</i><p>aemon Sets</p></nuxt-link></li>
						<li class="nav-item small"><nuxt-link to="/workload/replicaset.list"  class="nav-link"><i class="nav-icon text-sm mr-0">R</i><p>eplica Sets</p></nuxt-link></li>
						<li class="nav-item small"><nuxt-link to="/workload/replicationcontroller.list"  class="nav-link"><i class="nav-icon text-sm mr-0">R</i><p>eplication Controllers</p></nuxt-link></li>
						<!-- <li class="nav-item small"><nuxt-link to="/workload/hpa.list"  class="nav-link"><i class="nav-icon text-sm mr-0">H</i><p>orizontal Pod Autoscalers</p></nuxt-link></li> -->
					</ul>
				</li>
				<li class="nav-item  menu-open">
					<a href="pages/widgets.html" class="nav-link"><i class="nav-icon fas fa-network-wired"></i><p>Networking</p></a>
					<ul class="nav nav-treeview">
						<li class="nav-item small"><nuxt-link to="/networking/service.list"  class="nav-link"><i class="nav-icon text-sm mr-0">S</i><p>ervices</p></nuxt-link></li>
						<li class="nav-item small"><nuxt-link to="/networking/ingress.list"  class="nav-link"><i class="nav-icon text-sm mr-0">I</i><p>ngres</p></nuxt-link></li>
						<!-- <li class="nav-item small"><nuxt-link to="/networking/networkpolicy.list"  class="nav-link"><i class="nav-icon text-sm mr-0">N</i><p>etwork Plicies</p></nuxt-link></li> -->
						<!-- <li class="nav-item small"><nuxt-link to="/workload/pod.list"  class="nav-link"><i class="nav-icon text-sm mr-0">N</i><p>etwork Policies</p></nuxt-link></li> -->
						<!-- <li class="nav-item small"><nuxt-link to="/workload/pod.list"  class="nav-link"><i class="nav-icon text-sm mr-0">G</i><p>ateways</p></nuxt-link></li> -->
						<!-- <li class="nav-item small"><nuxt-link to="/workload/pod.list"  class="nav-link"><i class="nav-icon text-sm mr-0">V</i><p>irtual Services</p></nuxt-link></li> -->
						<!-- <li class="nav-item small"><nuxt-link to="/workload/pod.list"  class="nav-link"><i class="nav-icon text-sm mr-0">D</i><p>estination Rules</p></nuxt-link></li> -->
					</ul>
				</li>
				<li class="nav-item menu-open">
					<a href="pages/widgets.html" class="nav-link"><i class="nav-icon fas fa-hdd"></i><p>Storage</p></a>
					<ul class="nav nav-treeview">
						<li class="nav-item small"><nuxt-link to="/storage/pvc.list"  class="nav-link"><i class="nav-icon text-sm mr-0">P</i><p>ersistent Valumn Claims</p></nuxt-link></li>
						<li class="nav-item small"><nuxt-link to="/storage/storageclass.list"  class="nav-link"><i class="nav-icon text-sm mr-0">S</i><p>troage Classes</p></nuxt-link></li>
					</ul>
				</li>
				<li class="nav-item menu-open">
					<a href="pages/widgets.html" class="nav-link"><i class="nav-icon fas fa-users-cog"></i><p>Adminstrator</p></a>
					<ul class="nav nav-treeview">
						<li class="nav-item small"><nuxt-link to="/administrator/serviceaccount.list"  class="nav-link"><i class="nav-icon text-sm mr-0">S</i><p>ervice Accounts</p></nuxt-link></li>
						<li class="nav-item small"><nuxt-link to="/administrator/clusterrole.list"  class="nav-link"><i class="nav-icon text-sm mr-0">C</i><p>luster Roles</p></nuxt-link></li>
						<li class="nav-item small"><nuxt-link to="/administrator/clusterrolebinding.list"  class="nav-link"><i class="nav-icon text-sm mr-0">C</i><p>luster Role Bindings</p></nuxt-link></li>
						<li class="nav-item small"><nuxt-link to="/administrator/role.list"  class="nav-link"><i class="nav-icon text-sm mr-0">R</i><p>oles</p></nuxt-link></li>
						<li class="nav-item small"><nuxt-link to="/administrator/rolebinding.list"  class="nav-link"><i class="nav-icon text-sm mr-0">R</i><p>ole Bindings</p></nuxt-link></li>
						<li class="nav-item small"><nuxt-link to="/administrator/configmap.list"  class="nav-link"><i class="nav-icon text-sm mr-0">C</i><p>onfig Maps</p></nuxt-link></li>
						<li class="nav-item small"><nuxt-link to="/administrator/secret.list"  class="nav-link"><i class="nav-icon text-sm mr-0">S</i><p>ecrets</p></nuxt-link></li>
						<!-- <li class="nav-item small"><nuxt-link to="/workload/pod.list"  class="nav-link"><i class="nav-icon text-sm mr-0">R</i><p>esource Quotas</p></nuxt-link></li> -->
						<!-- <li class="nav-item small"><nuxt-link to="/workload/pod.list"  class="nav-link"><i class="nav-icon text-sm mr-0">L</i><p>imite Ranges</p></nuxt-link></li> -->
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