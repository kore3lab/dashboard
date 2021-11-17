<template>
<div :id="id" class="cursor text-secondary">
	<b-icon icon="layout-three-columns"></b-icon>
	<b-popover triggers="hover" :target="id" placement="top">
		<b-form-group class="mb-0">
			<b-button size="sm" variant="light" class="w-100 text-left p-0" @click="selectedFields=fields;raiseModelChange();"><i :id="id"  class="fas fa-eraser ml-0 mr-2"></i>Select All</b-button>
			<b-form-checkbox v-for="d in fields" v-model="selectedFields" :key="d.key" :value="d" @change="raiseModelChange">{{ d.label }}</b-form-checkbox>
		</b-form-group>
	</b-popover>
</div>
</template>
<script>

export default {
	props:["value", "name", "fields"],
	model: {
		prop: "value",
		event: "change"
	},
	data () {
		return {
			id: `grid_colums_selector_${Math.floor(Math.random() * 99)}`,
			selectedFields: [],
		}
	},
	mounted(){
		//get selected fields from localstorage
		let v = this.$storage.local.get(`columns.${this.name}`);
		if(!v) this.selectedFields = this.fields;
		else {
			this.selectedFields = this.fields.filter(d => {
				return v.includes(`${d.key}|`);
			})
		}
		
		this.raiseModelChange();	//v-model  emit event
	},
	methods: {
		raiseModelChange() {
			// selected fileds 
			let r = [];
			this.fields.forEach(el => {
				if(this.selectedFields.includes(el)) r.push(el)
			});

			// save to localstorage / clean-up localstorage 
			if(r.length ==0 || this.selectedFields == this.fields ) this.$storage.local.delete(`columns.${this.name}`);
			else {
				let v = "";
				this.selectedFields.forEach(d => {
					v += `${d.key}|`;
				})
				this.$storage.local.set(`columns.${this.name}`, v);
			}

			this.$emit("change", r)
		}
	}

}
</script>
