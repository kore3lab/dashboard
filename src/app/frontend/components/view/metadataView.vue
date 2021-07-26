<template>
<div class="row">
	<div v-bind:class="allClass">
		<div class="card card-secondary card-outline">
			<div class="card-body p-2">
				<dl class="row mb-0">
					<dt v-bind:class="dtClass">Create at</dt><dd v-bind:class="ddClass">{{ this.getTimestampString(value.creationTimestamp)}} ago ({{ value.creationTimestamp }})</dd>
					<dt v-bind:class="dtClass">Name</dt><dd v-bind:class="ddClass">{{ value.name }}</dd>
					<dt v-bind:class="dtClass">Namespace</dt><dd v-bind:class="ddClass">{{ value.namespace }}</dd>
					<dt v-bind:class="dtClass">Annotations</dt>
					<dd v-bind:class="ddClass">
						<a v-if="value.annotations" href="#"><b-icon :icon="isEllipseAnnotations?'arrows-expand':'arrows-collapse'" class="float-right mt-2" @click="isEllipseAnnotations=!isEllipseAnnotations"></b-icon></a>
						<ul class="list-unstyled mb-0">
							<li v-for="(v, k) in value.annotations" v-bind:key="k" v-bind:class="{'text-truncate':isEllipseAnnotations}">{{ k }}=<span class="font-weight-light">{{ v }}</span></li>
						</ul>
					</dd>
					<dt v-bind:class="dtClass">Labels</dt>
					<dd v-bind:class="ddClass">
						<span v-for="(value, name) in value.labels" v-bind:key="name" class="label">{{ name }}={{ value }}</span>
					</dd>
					<dt v-bind:class="dtClass">UID</dt><dd  v-bind:class="ddClass">{{ value.uid }}</dd>
					<slot></slot>
				</dl>
			</div>
		</div>
	</div>
</div>
</template>
<script>


export default {
	props:["value","dtCols","ddCols","size" ],
	data () {
		return {
			allClass: `col-${this.size?this.size:'sm'}-12`,
			dtClass: `col-${this.size?this.size:'sm'}-${this.dtCols?this.dtCols:'2'}`,
			ddClass: `col-${this.size?this.size:'sm'}-${this.ddCols?this.ddCols:'10'}`,
			isEllipseAnnotations: true
		}
	}
}
</script>
