<template>
<div id="wrapEditor"></div>
</template>
<script>
import "@/components/json-tree/jsonTree.css";

import jsonTree			from "@/components/json-tree/jsonTree";

export default {
	props:["value"],
	data () {
		return {
			localValue: this.value
		}
	},
	mounted(){
	},
	watch: {
		value(newVal) {
			try {
				if(this.localValue != newVal) {
					jsonTree.create(newVal, this.$el).expand();
					this.localValue = newVal;
				}
			} catch (ex) {
				console.error(ex.message);
			}
		}
	}
}
</script>
