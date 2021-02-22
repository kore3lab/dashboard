<template>
	<!-- Navbar -->
	<nav class="main-header navbar navbar-expand navbar-white navbar-light">
		<ul class="navbar-nav">
			<!-- 햄버거 메뉴 -->
			<li class="nav-item">
				<a class="nav-link" data-widget="pushmenu" href="#" role="button"><i class="fas fa-bars"></i></a>
			</li><!--// END -->
			<!-- 클러스터 선택 -->
			<li class="nav-item">
				<b-dropdown :text="ctx" variant="danger" size="sm" class="btn-group mt-1">
					<b-dropdown-item v-for="option in contexts()" :key="option" :value="option" @click="onContextSelected(option)">{{option}}</b-dropdown-item>
				</b-dropdown>
			</li><!--// END -->
		</ul>
		<ul class="navbar-nav ml-auto">
			<!-- <li class="nav-item">
				<a class="nav-link" data-widget="control-sidebar" data-slide="true" href="#" role="button">
					<i class="fas fa-th-large"></i>
				</a>
			</li> -->
		</ul>
	</nav>
</template>
<script>
import axios from "axios"
export default {
	data() {
		return {
			ctx: ""
		}
	},
	async fetch() {
		// context 리스트 조회
		let resp = await axios.get(`${this.backendUrl()}/api/clusters`);
		if(resp.data.contexts) this.contexts(resp.data.contexts);
		this.$data.ctx = this.$route.query.context ? this.$route.query.context: resp.data.currentContext;
		// namespace 로딩
		await this.loadNamespaces(this.$data.ctx);
		this.$nuxt.$emit("navbar-context-selected");
	},
	methods: {
		/**
		 * 클러스터의 네임스페이스 목록을 가져 온다.
		 * 
		 * @async
		 * @function loadNamespaces
		 * @returns {Array} 네이스페이스 값을 [{valu, text}] 값으로 반환 한다.
		 */
		async loadNamespaces(){
			
			// namespace 리스트 조회
			let nsList = [{ value: "", text: "All Namespaces" }];

			if (this.$data.ctx) {
				let resp = await axios.get(`${this.backendUrl()}/raw/clusters/${this.$data.ctx}/api/v1/namespaces`)
				resp.data.items.forEach(el => {
					nsList.push({ value: el.metadata.name, text: el.metadata.name });
				});
			}
			this.currentContext(this.$data.ctx);
			this.namespaces(nsList);

		},
		onContextSelected(ctx) {
			this.$data.ctx = ctx;
			this.loadNamespaces().then(() => {
				this.$nuxt.$emit("navbar-context-selected");
			}).catch(error=> {
				this.toast(error.message, "danger");
			});
			
		}
		
	}
}
</script>