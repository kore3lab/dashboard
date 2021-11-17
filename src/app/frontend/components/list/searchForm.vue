<template>
<div class="row">
	<div v-if="!noNamespace" class="col-sm-2"><b-form-select v-model="selectedNamespace" :options="namespaces()" size="sm" @input="dispatchInput"></b-form-select></div>
	<div class="col-sm-2 float-left">
		<div class="input-group input-group-sm" >
			<b-form-input v-model="keyword" class="form-control float-right" placeholder="Search" @input="$emit('keyword',keyword)"></b-form-input>
			<div class="input-group-append"><button type="submit" class="btn btn-default" @click="dispatchInput"><i class="fas fa-search"></i></button></div>
		</div>
	</div>
	<div v-if="!noLabelSelector" class="col-sm-8"><b-form-tags no-outer-focus class="form-labels-filter" size="sm" v-model="selectedLabels" @input="dispatchInput" tag-variant="info" separator=" ," placeholder="Enter new label selector"></b-form-tags></div>
</div>
</template>
<style scoped>
.form-labels-filter { border: 0; border-bottom: 1px solid #ced4da; background: transparent; padding-left: 0; }
.form-labels-filter >>> span { font-weight: 400; line-height: 1.1; font-size: .9rem;} 
.form-labels-filter >>> li { border-radius: .5rem;}
</style>
<script>
export default {
	props: {
		noLabelSelector: Boolean,
		noNamespace: Boolean
	},
	data () {
		return {
			keyword: ""
		}
	},
	computed: {
		selectedLabels: {
			get () {
				return this.labelSelector()
			},
			set (value) {
				this.labelSelector(value)
			}
		},
		selectedNamespace: {
			get () {
				return this.selectNamespace()
			},
			set (value) {
				this.selectNamespace(value)
			}
		}
	},
	created(){
		this.$nuxt.$on("context-selected", this.dispatchInput);
		if(this.currentContext()) this.dispatchInput();
	},
	methods: {
		dispatchInput() {
			this.$emit("input", {
				namespace: this.noNamespace? "": this.selectedNamespace,
				labelSelector: this.noLabelSelector ? "": this.labelSelector().length>0? this.labelSelector().reduce( (accumulator, d) => { return `${accumulator}&${d}`;}): "",
				keyword: this.keyword,
			});
		}
	},
	beforeDestroy(){
		this.$nuxt.$off("context-selected");
	}
}
</script>
