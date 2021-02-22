<template>
<div id="wrapEditor"></div>
</template>
<script>
import Ace from "ace-builds/src-noconflict/ace"
import "ace-builds/webpack-resolver"
import "ace-builds/src-noconflict/mode-yaml"
import {dump as toYaml, load as toJSON} from "js-yaml";

export default {
	props:["value"],
	data () {
		return {
			localValue: this.value
		}
	},
	mounted(){
		//  editor
		this.editor = Ace.edit(this.$el.id);
		this.editor.setOptions({
			showPrintMargin: false,
			highlightActiveLine: true,
			tabSize: 4,
			wrap: true,
			fontSize: 12,
			// fontFamily: `'Roboto Mono Regular', monospace`,
		});
		this.editor.setReadOnly(false);
		this.editor.session.setMode("ace/mode/yaml");
		this.editor.session.setUseWorker(false);
		this.editor.on('blur', () => {
			try {
				this.localValue = toJSON(this.editor.getValue())
				this.$emit("input", this.localValue);
			} catch (ex) {
				console.error(ex.message);
				this.$emit("error", ex);
			}
		});
	},
	watch: {
		value(newVal) {
			try {
				this.editor.setValue(toYaml(newVal), -1);
				this.value = newVal;
			} catch (ex) {
				console.error(ex.message);
			}
		}
	}
}
</script>
