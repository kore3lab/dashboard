<template>
<div>
	<!-- 1. metadata -->
	<c-metadata v-model="metadata" dtCols="3" ddCols="9">
		<dt class="col-sm-3">Group</dt><dd class="col-sm-9">{{ info.group }}</dd>
		<dt class="col-sm-3">Version</dt><dd class="col-sm-9">{{ info.version }}</dd>
		<dt class="col-sm-3">Stored versions</dt><dd class="col-sm-9">{{ info.storedVersions }}</dd>
		<dt class="col-sm-3">Scope</dt><dd class="col-sm-9">{{ info.scope }}</dd>
		<dt class="col-sm-3">Resource</dt><dd class="col-sm-9"><span v-for="(val, idx) in info.resource" v-bind:key="idx"><nuxt-link :to="{path: '/customresource/customresource.list', query: {crd: `${val.plural}.${info.group}`, version: info.version}}" class="mr-2">{{ val.plural }}</nuxt-link></span></dd>
		<dt class="col-sm-3">Conversion</dt><dd class="col-sm-9"><c-jsontree id="txtConversion" v-model="info.conversion" class="card-body p-2 border"></c-jsontree></dd>
		<dt class="col-sm-3">Conditions</dt>
		<dd class="col-sm-9">
			<span v-for="(val, idx) in info.conditions" v-bind:key="idx" v-bind:class="{'badge-success': (val.type=='Established' || val.type=='NamesAccepted'), 'badge-danger':val.type=='NonStructuralSchema', 'badge-secondary':val.type=='Terminating','badge-warning':val.type=='KubernetesAPIApprovalPolicyConformant'}" class="badge font-weight-light text-sm mr-1">{{ val.type }}</span>
		</dd>
	</c-metadata>
	<!-- 2. names -->
	<div class="row">
		<div class="col-md-12">
			<div class="card card-secondary card-outline">
				<div class="card-header p-2"><h3 class="card-title">Names</h3></div>
				<div class="card-body p-2">
					<b-table-lite small :items="info.resource" :fields="fields" class="subset"></b-table-lite>
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
					<b-table-lite small :items="printerColumns" :fields="columnsFields" class="subset"></b-table-lite>
				</div>
			</div>
		</div>
	</div>
	<!-- 4. validation -->
	<div v-show="Object.keys(validation).length > 0" class="row">
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
import VueMetadataView	from "@/components/view/metadataView.vue";
import VueJsonTree 		from "@/components/jsontree";

export default {
	components: {
		"c-metadata": { extends: VueMetadataView },
		"c-jsontree": {extends: VueJsonTree},
	},
	data() {
		return {
			metadata: {},
			info: [],
			printerColumns: [],
			validation: {},
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
			this.info = {
				group: data.spec.group,
				version: data.spec.versions[0]?.name ?? data.spec.version,
				storedVersions: data.status.storedVersions.join(),
				scope: data.spec.scope,
				resource: [data.spec.names],
				conversion: data.spec.conversion,
				conditions: data.status?.conditions ?? []
			}
			this.printerColumns = this.getPrinterColumns(data.spec)
			this.validation = data.spec.validation || data.spec.versions?.[0]?.schema || {}
		});
		this.$nuxt.$emit("onCreated",'')
	},
	methods: {
		getPrinterColumns(spec) {
			const columns = spec.versions.find(a => spec.versions[0]?.name ?? spec.version === a.name)?.additionalPrinterColumns ??
					spec.additionalPrinterColumns?.map(({ JSONPath, ...rest}) => ({ ...rest, jsonPath: JSONPath}))
					.filter(column => column.name !== 'Age') ?? [];
			return columns;
		}
	},
	beforeDestroy(){
		this.$nuxt.$off("onReadCompleted");
	},
}
</script>
