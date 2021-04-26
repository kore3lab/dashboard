<template>
	<div class="card-body p-2">
		<div class="row">
			<div class="col-md-12">
				<div class="card card-secondary card-outline">
					<div class="card-body p-2">
						<dl class="row mb-0">
							<dt class="col-sm-2 text-truncate">Create at</dt><dd class="col-sm-10">{{ this.getTimestampString(metadata.creationTimestamp)}} ago ({{ metadata.creationTimestamp }})</dd>
							<dt class="col-sm-2">Name</dt><dd class="col-sm-10">{{ metadata.name }}</dd>
							<dt class="col-sm-2">Namespace</dt><dd class="col-sm-10">{{ metadata.namespace }}</dd>
							<dt class="col-sm-2">Annotations</dt>
							<dd class="col-sm-10 text-truncate">
								<ul class="list-unstyled mb-0">
									<li v-for="(value, name) in metadata.annotations" v-bind:key="name"><span class="badge badge-secondary font-weight-light text-sm mb-1">{{ name }}:{{ value }}</span></li>
								</ul>
							</dd>
							<dt class="col-sm-2">Labels</dt>
							<dd class="col-sm-10 text-truncate">
								<ul class="list-unstyled mb-0">
									<li v-for="(value, name) in metadata.labels" v-bind:key="name"><span class="badge badge-secondary font-weight-light text-sm mb-1">{{ name }}:{{ value }}</span></li>
								</ul>
							</dd>
							<dt class="col-sm-2">Quotas</dt><dd class="col-sm-10">
							<ul v-for="(val, idx) in quotas" v-bind:key="idx" class="list-unstyled mb-0">
								<li><span>{{ val.name }}</span> <span class="float-right">{{ val.current }} / {{ val.temp }}</span></li>
								<b-progress :value="val.current" :max="val.max" show-value class="mb-3"></b-progress>
							</ul>
						</dd>
						</dl>
					</div>
				</div>
			</div>
		</div>

		<div v-if="scopeSelector" class="row">
			<div class="col-md-12">
				<div class="card card-secondary card-outline m-0">
					<div class="card-header p-2"><h3 class="card-title text-md">Scope Selector</h3></div>
					<div class="card-body p-2">
						<b-table striped hover small :items="scopeSelector" :fields="fields"></b-table>
					</div>
				</div>
			</div>
		</div>

	</div>
</template>
<script>
import VueAceEditor from "@/components/aceeditor";

export default {
	components: {
		"c-aceeditor": { extends: VueAceEditor },
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
	},
	beforeDestroy(){
		this.$nuxt.$off("onReadCompleted");
	},
}
</script>
