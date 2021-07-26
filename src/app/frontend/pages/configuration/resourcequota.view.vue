<template>
<div>
	<!-- 1. metadata -->
	<c-metadata v-model="metadata" dtCols="2" ddCols="10">
		<dt class="col-sm-2">Quotas</dt>
		<dd class="col-sm-10">
			<ul v-for="(val, idx) in quotas" v-bind:key="idx" class="list-unstyled mb-0">
				<li><span>{{ val.name }}</span> <span class="float-right">{{ val.current }} / {{ val.temp | formatNumber }}</span></li>
				<b-progress :value="val.current" :max="val.max" show-value class="mb-3"></b-progress>
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
			this.onSync(data)
		});
		this.$nuxt.$emit("onCreated",'')
	},
	methods: {
		onSync(data) {
			this.quotas = this.getQuotas(data.status)
			this.scopeSelector = this.getScope(data.spec.scopeSelector)
		},
		getQuotas(status) {
			let list = [];
			let temp;
			const { hard = {}, used = {} } = status
			Object.entries(hard)
					.filter(([name]) => used[name])
					.map(([name, value]) => {
						let current = this.transformUnit(name, used[name]);
						let max = this.transformUnit(name, value);
						temp = max
						let usage = max === 0 ? 100 : Math.ceil(current / max * 100);
						if(name === 'cpu') {
							temp = Number(max)
							if(temp >= 1000) {
								temp = max / 1000
								temp+='k'
							}
						}
						if(name === 'memory' || name ==='storage') {
							temp = Number(max)
							if(temp >= 1024) {
								temp = temp / 1024
								if( temp >= 1024) {
									temp = temp / 1024
									if( temp >= 1024) {
										temp = temp / 1024
										if ( temp >= 1024) {
											max = temp / 1024
											if (temp >= 1024) {
												temp = temp / 1024
											} else temp += 'Ti'
										} else temp += 'Gi'
									} else temp += 'Mi'
								} else temp += 'Ki'
							}
						}
						list.push({
							name: name,
							current: current,
							max: max,
							usage: usage,
							temp: temp,
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
				return this.unitsToBytes(value);
			}
			if(name.includes('cpu')) {
				return this.cpuUnitsToNumber(value);
			}
			return this.metricUnitsToNumber(value);
		},
		metricUnitsToNumber(value) {
			const base = 1000;
			const suffixes = ["k", "m", "g", "t", "p"];

			const suffix = value.toLowerCase().slice(-1);
			const index = suffixes.indexOf(suffix);

			return parseInt(
				(parseFloat(value) * Math.pow(base, index + 1)).toFixed(1)
			);
		},
		cpuUnitsToNumber(cpu) {
			const thousand = 1000;
			const million = thousand * thousand;
			const shortBillion = thousand * million;

			const cpuNum = parseInt(cpu);
			if (cpu.includes("k")) return cpuNum * thousand;
			if (cpu.includes("m")) return cpuNum / thousand;
			if (cpu.includes("u")) return cpuNum / million;
			if (cpu.includes("n")) return cpuNum / shortBillion;

			return parseFloat(cpu);
		},
		unitsToBytes(value) {
			const base = 1024;
			const suffixes = ["K", "M", "G", "T", "P", "E"];

			if (!suffixes.some(suffix => value.includes(suffix))) {
				return parseFloat(value);
			}

			const suffix = value.replace(/[0-9]|i|\./g, "");
			const index = suffixes.indexOf(suffix);

			return parseInt(
				(parseFloat(value) * Math.pow(base, index + 1)).toFixed(1)
			);
		},
	},
	beforeDestroy(){
		this.$nuxt.$off("onReadCompleted");
	},
}
</script>
