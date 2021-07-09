<template>
<div>
	<!-- 1. metadata -->
	<div class="row">
		<div class="col-md-12">
			<div class="card card-secondary card-outline">
				<div class="card-body p-2">
					<dl class="row mb-0">
						<dt class="col-sm-3">Create at</dt><dd class="col-sm-9">{{ this.getTimestampString(metadata.creationTimestamp)}} ago ({{ metadata.creationTimestamp }})</dd>
						<dt class="col-sm-3">Name</dt><dd class="col-sm-9">{{ metadata.name }}</dd>
						<dt class="col-sm-3">Namespace</dt><dd class="col-sm-9">{{ metadata.namespace }}</dd>
						<dt class="col-sm-3">Annotations</dt>
						<dd class="col-sm-9">
							<ul class="list-unstyled mb-0">
								<li v-for="(value, name) in metadata.annotations" v-bind:key="name">{{ name }}=<span class="font-weight-light">{{ value }}</span></li>
							</ul>
						</dd>
						<dt class="col-sm-3">Labels</dt>
						<dd class="col-sm-9">
							<span v-for="(value, name) in metadata.labels" v-bind:key="name" class="label">{{ name }}={{ value }}</span>
						</dd>
						<dt class="col-sm-3">Group</dt><dd class="col-sm-9">{{ info.group }}</dd>
						<dt class="col-sm-3">Version</dt><dd class="col-sm-9">{{ info.version }}</dd>
						<dt class="col-sm-3">Stored versions</dt><dd class="col-sm-9">{{ info.storedVersions }}</dd>
						<dt class="col-sm-3">Scope</dt><dd class="col-sm-9">{{ info.scope }}</dd>
						<dt class="col-sm-3">Resource</dt><dd class="col-sm-9"><span v-for="(val, idx) in info.resource" v-bind:key="idx"><nuxt-link :to="{path: '/customresource/customresource.list', query: {group: info.group, plural:info.resource[0].plural, kind: info.resource[0].kind, columnsName: info.printerColumns? info.printerColumns.name : '', columnsPath: info.printerColumns? info.printerColumns.path : '', scope: info.scope }}" class="mr-2">{{ val.plural }}</nuxt-link></span></dd>
						<dt class="col-sm-3">Conversion</dt><dd class="col-sm-9"><c-jsontree id="txtConversion" v-model="info.conversion" class="card-body p-2 border"></c-jsontree></dd>
						<dt class="col-sm-3">Conditions</dt>
						<dd class="col-sm-9">
							<span v-for="(val, idx) in info.conditions" v-bind:key="idx" v-bind:class="val.style" class="badge font-weight-light text-sm mr-1">{{ val.type }}</span>
						</dd>
					</dl>
				</div>
			</div>
		</div>
	</div>
	<!-- 2. names -->
	<div class="row">
		<div class="col-md-12">
			<div class="card card-secondary card-outline">
				<div class="card-header p-2"><h3 class="card-title">Names</h3></div>
				<div class="card-body p-2">
					<b-table-lite small :items="info.resource" :fields="fields"></b-table-lite>
				</div>
			</div>
		</div>
	</div>
	<!-- 3. additional printer columns -->
	<div v-if="printerColumns.length !== 0" class="row">
		<div class="col-md-12">
			<div class="card card-secondary card-outline">
				<div class="card-header p-2"><h3 class="card-title">Additional Printer Columns</h3></div>
				<div class="card-body p-2">
					<b-table-lite small :items="printerColumns" :fields="columnsFields"></b-table-lite>
				</div>
			</div>
		</div>
	</div>
	<!-- 4. validation -->
	<div v-if="validation" class="row">
		<div class="col-md-12">
			<div class="card card-secondary card-outline">
				<div class="card-header p-2"><h3 class="card-title">Validation</h3></div>
				<div class="card-body p-2">
					<c-jsontree id="txtValidation" v-model="validation" class="card-body p-2 border"></c-jsontree>
				</div>
			</div>
		</div>
	</div>
</div>
</template>
<script>

import VueJsonTree from "@/components/jsontree";

export default {
	components: {
		"c-jsontree": {extends: VueJsonTree},
	},
	data() {
		return {
			metadata: {},
			info: [],
			printerColumns: [],
			validation: [],
			fields: [
				{ key: "plural", label: "plural" },
				{ key: "singular", label: "singular" },
				{ key: "kind", label: "kind" },
				{ key: "listKind", label: "listKind" },
			],
			columnsFields: [
				{ key: "name", label: "Name" },
				{ key: "type", label: "Type" },
				{ key: "jsonPath", label: "JSON Path" },
			],
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
			this.info = this.getInfo(data)
			this.printerColumns = this.getPrinterColumns(data.spec)
			this.validation = this.getValidation(data.spec)
		},
		getInfo(data) {
			let resource = [];
			resource.push({
				group: data.spec.group,
				plural: data.spec.names.plural,
				kind: data.spec.names.kind,
				singular: data.spec.names.singular,
				listKind: data.spec.names.listKind,
			})
			let conditions = [];
			if(data.status?.conditions) {
				data.status.conditions.map(condition => {
					conditions.push({
						type: condition.type,
						status: condition.status,
						style: this.styleCheck(condition.type)
					})
				})
			}

			return {
				group: data.spec.group,
				version: data.spec.versions[0]?.name ?? data.spec.version,
				storedVersions: data.status.storedVersions.join(","),
				scope: data.spec.scope,
				resource: resource,
				conversion: data.spec.conversion,
				conditions: conditions,
				printerColumns: this.conv(this.getPrinterColumns(data.spec))
			}
		},
		getPrinterColumns(spec) {
			const columns = spec.versions.find(a => this.getVersion(spec) === a.name)?.additionalPrinterColumns ??
					spec.additionalPrinterColumns?.map(({ JSONPath, ...rest}) => ({ ...rest, jsonPath: JSONPath})) ?? [];
			return columns
					.filter(column => column.name !== 'Age')
					.filter(column => column.jsonPath ? true : !column.priority );
		},
		getValidation(spec) {
			return spec.validation?? spec.versions?.[0]?.schema
		},
		getVersion(spec) {
			return spec.versions[0]?.name ?? spec.version
		},
		conv(col) {
			if(col.length === 0) return

			let name = []
			let path = []
			col.forEach(el => {
				name.push(el.name)
				path.push(el.jsonPath)
			})
			return {
				name: name.join(','),
				path: path.join(','),
			}
		},
		styleCheck(type) {
			if(type === 'Established') {
				return 'badge-success'
			}else if(type === 'NamesAccepted') {
				return 'badge-success'
			}else if(type === 'NonStructuralSchema') {
				return 'badge-danger'
			}else if(type === 'Terminating') {
				return 'badge-secondary'
			}else if(type === 'KubernetesAPIApprovalPolicyConformant') {
				return 'badge-warning'
			}
		},

	},
	beforeDestroy(){
		this.$nuxt.$off("onReadCompleted");
	},
}
</script>
