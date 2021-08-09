<template>
<div>
	<!-- 1. metadata -->
	<c-metadata v-model="metadata" dtCols="2" ddCols="10">
		<dt class="col-sm-2">Quotas</dt>
		<dd class="col-sm-10">
			<ul v-for="(val, idx) in quotas" v-bind:key="idx" class="list-unstyled mb-0">
				<li><span>{{ val.name }}</span> <span class="float-right">{{ val.current }} / {{ val.max }} </span></li>
				<b-progress :max="val.max">
					<b-progress-bar :value="val.current" :label="`${((parseInt(val.current) / parseInt(val.max)) * 100).toFixed(0)}%`"></b-progress-bar>
				</b-progress>
			</ul>
		</dd>
	</c-metadata>
	<!-- 2. scope selector -->
	<div v-if="scopeSelector" class="row">
		<div class="col-md-12">
			<div class="card card-secondary card-outline m-0">
				<div class="card-header p-2"><h3 class="card-title">Scope Selector</h3></div>
				<div class="card-body p-2">
					<b-table-lite striped hover small :items="scopeSelector" :fields="fields"></b-table-lite>
				</div>
			</div>
		</div>
	</div>

</div>
</template>
<script>
import VueMetadataView	from "@/components/view/metadataView.vue";

export default {
	components: {
		"c-metadata": { extends: VueMetadataView }
	},
	data() {
		return {
			metadata: {},
			quotas: [],
			scopeSelector: [],
			fields: [
				{key: "operator", label:"Operator"},
				{key: "scopeName", label:"Scope name"},
				{key: "values", label:"Values"},
			]
		}
	},
	mounted() {
		this.$nuxt.$on("onReadCompleted", (data) => {
			if(!data) return
			this.metadata = data.metadata;
			this.quotas = this.getQuotas(data.status)
			this.scopeSelector = this.getScope(data.spec.scopeSelector)
		});
		this.$nuxt.$emit("onCreated",'')
	},
	methods: {
		getQuotas(status) {
			let list = [];
			let temp;
			const { hard = {}, used = {} } = status
			Object.entries(hard)
					.filter(([name]) => used[name])
					.map(([name, value]) => {
						let current = this.transformUnit(name, used[name]);
						let max = this.transformUnit(name, value);
						list.push({
							name: name,
							current: current,
							max: max
						})
					})
			return list
		},
		getScope(scope) {
			let list = []
			if(!scope) return
			if(scope.matchExpressions) {
				scope.matchExpressions.map((selector, index) => {
					const { operator, scopeName, values } = selector;
					list.push({
						operator: operator,
						scopeName: scopeName,
						values: values.join(",")
					})
				})
			}
			return list
		},
		transformUnit(name, value) {
			if(name.includes('memory') || name.includes('storage')) {
				return this.unitsToBytes(name, value);
			}
			if(name.includes('cpu')) {
				return this.cpuUnitsToNumber(value) + "m";
			}
			return value;
		},
		
		cpuUnitsToNumber(value) {
			const thousand = 1000;
			const million = thousand * thousand;
			const cpuNum = parseInt(value);
			
			if (value.includes("k")) return cpuNum * million;
			if (value.includes("m")) return cpuNum;
			if (value.includes("u")) return cpuNum / thousand;
			if (value.includes("n")) return cpuNum / million;

			return cpuNum * thousand;
		},
		unitsToBytes(name, value) {
			const base = 1024;
			const suffixes = ["K", "M", "G", "T", "P", "E"];
			const suffix = value.replace(/[0-9]|i|\./g, "");
			const index = suffixes.indexOf(suffix);
			const val = name.includes('memory') ? parseFloat(value) * Math.pow(base, index-1) + "Mi" : parseFloat(value) * Math.pow(base, index-2) + "Gi";
			
			return val;
		},
	},
	beforeDestroy(){
		this.$nuxt.$off("onReadCompleted");
	},
}
</script>
