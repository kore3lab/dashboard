<template>
<div v-if="items.length>0" class="row">
	<div class="col-md-12">
		<div class="card card-secondary card-outline">
			<div class="card-header p-2"><h3 class="card-title">Pods</h3></div>
			<div class="card-body p-2">
				<b-table striped hover small :items="items" :fields="fields" :busy="isBusy" >
					<template #table-busy>
						<div class="text-center text-success">
							<b-spinner size="md" type="grow" variant="success" class="align-middle mr-2"></b-spinner>
							<span class="text-md align-middle">Loading...</span>
						</div>
					</template>
					<template v-slot:cell(name)="data">
						<a href="#" @click="$emit('navigate', getViewLink('', 'pods', data.item.namespace, data.item.name))">{{ data.item.name }}</a>
					</template>
					<template v-slot:cell(cpu)="data">
						<span>{{ data.item.metrics.cpu | formatNumber }} m</span>
					</template>
					<template v-slot:cell(memory)="data">
						<span>{{ (data.item.metrics.memory/(1024*1024)).toFixed(2) | formatNumber}} Mi</span>
					</template>
					<template v-slot:cell(status)="data">
						<span v-bind:class="{'text-success':data.value=='Completed'||data.value=='Running'}">{{ data.value }}</span>
					</template>
				</b-table>
			</div>
		</div>
	</div>
</div>
</template>
<script>

export default {
	props:["value", "namespace"],
	data () {
		return {
			isBusy: false,
			items: [],
			fields: [
				{ key: "name", label: "Name", sortable: true },
				{ key: "namespace", label: "Namespace", sortable: true, thClass: this.namespace?"":"d-none", tdClass: this.namespace?"":"d-none" },
				{ key: "ready", label: "Ready" , thClass:"text-center", tdClass:"text-center" },
				{ key: "cpu", label: "CPU", sortable: true, thClass:"text-right", tdClass:"text-right" },
				{ key: "memory", label: "Memory", sortable: true, thClass:"text-right", tdClass:"text-right" },
				{ key: "status", label: "Status" },
			],
		}
	},
	watch: {
		value(newVal) {
			this.isBusy = true;
			this.$axios.get(`/api/clusters/${this.currentContext()}/${newVal}/pods`)
				.then(resp => {
					this.items = resp.data
				}).catch(e => {
					this.msghttp(e);
					this.items = [];
				}).finally(()=>{
					this.isBusy = false
				});
		}
	}
}
</script>
