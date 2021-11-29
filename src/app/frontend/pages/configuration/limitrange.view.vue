<template>
<div>
	<!-- 1. metadata -->
	<c-metadata v-model="value" dtCols="2" ddCols="10"></c-metadata>
	<!-- 2. Limits -->
	<div class="row">
		<div class="col-md-12">
			<div class="card card-secondary card-outline">
				<div class="card-header p-2"><h3 class="card-title">Limits</h3></div>
				<div class="card-body group">
					<ul>
						<li v-for="(key, idx) in Object.keys(limits)" v-bind:key="idx">
							<p class="title">{{key}} Limit</p>
							<div class="ml-3">
								<b-table-lite thead-class="d-none" :items="limits[key]" class="subset">
									<template v-slot:cell(value)="d">
										<span v-for="(v,k) in d.item.value" v-bind:key="k" class="border-box">{{k}}={{v}}</span>
									</template>
								</b-table-lite>
							</div>
						</li>
					</ul>
				</div>
			</div>
		</div>
	</div>

</div>
</template>
<script>
import VueMetadataView	from "@/components/view/metadataView.vue";

export default {
	props:["value"],
	components: {
		"c-metadata": { extends: VueMetadataView }
	},
	data() {
		return {
			limits: {}
		}
	},
	watch: {
		value(d) { this.onSync(d) }
	},
	methods: {
		onSync(data) {
			if(!data) return
			data.spec.limits.forEach(el=> {
				let d = []
				for(let key in el) {
					if(key != "type") d.push({key:key, value:el[key]});
				}
				this.limits[el["type"]] = d;
			});
		}
	}
}
</script>
