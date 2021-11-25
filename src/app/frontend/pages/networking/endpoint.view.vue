<template>
<div>
	<!-- 1. metadata -->
	<c-metadata v-model="value" dtCols="2" ddCols="10"></c-metadata>
	<!-- 2. subsets -->
	<div v-if="subsets.length>0" class="row">
		<div class="col-md-12">
			<div class="card card-secondary card-outline">
				<div class="card-header p-2"><h3 class="card-title">Subsets</h3></div>
				<div class="card-body group">
					<ul>
						<li v-for="(subset, i) in subsets" v-bind:key="i">
							<p class="title">Addresses</p>
							<div  class="ml-3">
								<b-table-lite :items="subset.addresses" :fields="adressesFields" class="subset">
									<template v-slot:cell(targetRef)="data">
										<a href="#" @click="$emit('navigate', getViewLink('', data.item.resource.resource, data.item.resource.namespace, data.item.resource.name))">{{ data.item.resource.name }}</a>
									</template>
								</b-table-lite>
							</div>
							<p class="title">Ports</p>
							<div  class="ml-3">
								<b-table-lite :items="subset.ports" :fields="portsFields"  class="subset"></b-table-lite>
							</div>
						</li>
					</ul>
				</div>
			</div>
		</div>
	</div>
	<!-- 3. evenets -->
	<c-events v-model="value" class="row"></c-events>
</div>
</template>
<script>
import VueMetadataView	from "@/components/view/metadataView.vue";
import VueEventsView	from "@/components/view/eventsView.vue";

export default {
	props:["value"],
	components: {
		"c-metadata": { extends: VueMetadataView },
		"c-events": { extends: VueEventsView }
	},
	data() {
		return {
			address: [],
			ports: [],
			subsets: [],
			adressesFields: [
				{ key: "ip", label: "IP"},
				{ key: "hostname", label: "Hostname"},
				{ key: "targetRef", label: "Target"},
			],
			portsFields: [
				{ key: "port", label: "Port"},
				{ key: "name", label: "Name"},
				{ key: "protocol", label: "Protocol"},
			],
		}
	},
	watch: {
		value(d) { this.onSync(d) }
	},
	methods: {
		onSync(data) {
			if(!data) return
			this.subsets = data.subsets?data.subsets:[];
			//populate - subset resource(for view link)
			this.subsets.forEach(subset => {
				if(subset.addresses)  {
					subset.addresses.forEach(addr => {
						addr.resource = this.getResource(addr.targetRef);
					});
				}
			});
		}
	}
}
</script>
