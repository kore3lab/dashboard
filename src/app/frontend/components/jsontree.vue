<template>
<div></div>
</template>
<script>
import "@/components/json-tree/jsonTree.css";

import jsonTree			from "@/components/json-tree/jsonTree";

export default {
	props:["value"],
	data () {
		return {
			localValue: this.value,
			object: null
		}
	},
	mounted(){
	},
	watch: {
		value(newVal) {
			try {
				if(this.localValue != newVal) {
					if(this.object) {
						this.object.loadData(newVal);
					} else {
						this.object = jsonTree.create(newVal, this.$el);
					}
					this.localValue = newVal;
					this.object.expand();
				}
			} catch (ex) {
				console.error(ex.message);
			}
		}
	}
}
</script>
