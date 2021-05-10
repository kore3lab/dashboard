<template>
	<div class="card-body p-2">
		<b-overlay :show="deleteOverlay" rounded="sm" no-center>
			<div class="row">
				<div class="col-md-12">
					<div class="card card-secondary card-outline">
						<div class="card-body p-2">
							<dl class="row mb-0">
								<dt class="col-sm-2 text-truncate">Create at</dt><dd class="col-sm-10">{{ this.getTimestampString(metadata.creationTimestamp)}} ago ({{ metadata.creationTimestamp }})</dd>
								<dt class="col-sm-2">Name</dt><dd class="col-sm-10">{{ metadata.name }}</dd>
								<dt class="col-sm-2">Namespace</dt><dd class="col-sm-10">{{ metadata.namespace }}</dd>
								<dt v-if="metadata.ownerReferences" class="col-sm-2">Controlled By</dt><dd v-if="metadata.ownerReferences" class="col-sm-10">{{ controller.k }} <a href="#" @click="$emit('navigate', getViewLink(controller.g, controller.k, metadata.namespace, metadata.ownerReferences[0].name))">{{ metadata.ownerReferences[0].name }}</a></dd>
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
							</dl>
						</div>
					</div>
				</div>
			</div>

			<div class="row">
				<div class="col-md-12">
					<div class="card card-secondary card-outline">
						<div class="card-header p-2"><h3 class="card-title text-md">Reference</h3></div>
						<div class="card-body p-2">
							<b-table striped hover small :items="ref" :fields="fields"></b-table>
						</div>
					</div>
				</div>
			</div>

			<div class="row">
				<div class="col-md-12">
					<div class="card card-secondary card-outline">
						<div class="card-header p-2"><h3 class="card-title text-md">Bindings</h3></div>
						<div class="card-body p-2">
							<div v-if="bindings">
								<b-table hover small selectable select-mode="multi" ref="selectableTable" responsive="sm" @row-selected="onRowSelected" :items="bindings" :fields="bindFields">
									<template v-slot:cell(selected)="{ rowSelected }">
										<template v-if="rowSelected">
											<span aria-hidden="true">&check;</span>
											<span class="sr-only">Selected</span>
										</template>
										<template v-else>
											<span aria-hidden="true">&nbsp;</span>
											<span class="sr-only">Not Selected</span>
										</template>
									</template>
									<template v-slot:cell(namespace)="data">
										{{ data.value? data.value : '-'}}
									</template>
								</b-table>
								<b-button v-show="selected.length !== 0" variant="primary" size="sm" class="mb-1 ml-2" @click="deleteOverlay = true">Delete</b-button>
							</div>
						</div>
					</div>
				</div>
			</div>

			<div class="row">
				<div class="col-md-12">
					<div class="card card-secondary card-outline m-0">
						<div class="card-header p-2"><h3 class="card-title text-md">Events</h3></div>
						<div class="card-body p-2">
							<dl v-for="(val, idx) in event" v-bind:key="idx" class="row mb-0 card-body p-2 border-bottom">
								<dt class="col-sm-12"><p v-bind:class="val.type" class="mb-1">{{ val.name }}</p></dt>
								<dt class="col-sm-2 text-truncate">Source</dt><dd class="col-sm-10">{{ val.source }}</dd>
								<dt class="col-sm-2 text-truncate">Count</dt><dd class="col-sm-10">{{ val.count }}</dd>
								<dt class="col-sm-2 text-truncate">Sub-object</dt><dd class="col-sm-10">{{ val.subObject }}</dd>
								<dt class="col-sm-2 text-truncate">Last seen</dt><dd class="col-sm-10">{{ val.lastSeen }}</dd>
							</dl>
						</div>
					</div>
				</div>
			</div>
			<template #overlay>
				<div class="text-center">
					<p>Are you sure DELETE it?</p>
					<div class="text-center">
						<b-button variant="outline-danger" size="sm" class="mr-1" @click="deleteOverlay = false">Cancel</b-button>
						<b-button variant="success" size="sm" @click="deleteBinding">OK</b-button>
					</div>
				</div>
			</template>
		</b-overlay>
	</div>
</template>
<script>

export default {
	data() {
		return {
			metadata: {},
			event: [],
			controller: [],
			ref: [],
			bindings: [],
			fields: [
				{ key: "kind", label: "Kind" },
				{ key: "name", label: "Name" },
				{ key: "apiGroup", label: "Api Group" },
			],
			bindFields: [
				{ key: "selected", label: "Selected"},
				{ key: "name", label: "Binding" },
				{ key: "kind", label: "Type" },
				{ key: "namespace", label: "Namespace" },
			],
			selected: [],
			deleteOverlay: false,
		}
	},
	mounted() {
		this.$nuxt.$on("onReadCompleted", (data) => {
			if(!data) return
			this.origin = data;
			this.metadata = data.metadata;
			this.onSync(data)
		});
		this.$nuxt.$emit("onCreated",'')
	},
	methods: {
		onRowSelected(items) {
			this.selected = items
		},
		onSync(data) {
			this.controller = this.getController(data.metadata.ownerReferences);
			this.ref = this.getRef(data.roleRef)
			this.bindings = this.getBindings(data.subjects)
			this.event = this.getEvents(data.metadata.uid,'fieldSelector=involvedObject.name='+data.metadata.name);
		},
		getRef(roleRef) {
			if(!roleRef) return

			let list = []
			list.push({
				kind: roleRef.kind,
				name: roleRef.name,
				apiGroup: roleRef.apiGroup
			})
			return list
		},
		getBindings(subjects) {
			if(!subjects) return

			let list = []
			subjects.map((val,idx) => {
				list.push({
					apiGroup: val.apiGroup,
					name: val.name || '',
					kind: val.kind || '',
					namespace: val.namespace || '',
					idx: idx,
				})
			})
			return list
		},
		deleteBinding() {
			let list = this.origin.subjects
			let temp = this.origin
			this.selected.forEach(el => {
				list.forEach(e => {
					if(e.apiGroup? e.apiGroup : "" === el.apiGroup && e.kind ? e.kind : "" === el.kind && e.name ? e.name : "" === el.name && e.namespace? e.namespace : "" === el.namespace) {
						list.splice(el.idx,1)
					}
				})
			})
			this.origin.subjects = list
			this.$axios.put(`/raw/clusters/${this.currentContext()}`, this.origin)
					.then( _ => {
						this.onSync(this.origin)
					})
					.catch(e => {
						this.origin = temp
						this.msghttp(e);
					});
			this.deleteOverlay = false;

		},
	},
	beforeDestroy(){
		this.$nuxt.$off("onReadCompleted");
	},
}
</script>
