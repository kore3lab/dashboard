<template>
<div id="wrapEditor"></div>
</template>
<script>
import Ace from "ace-builds/src-noconflict/ace"
import "ace-builds/webpack-resolver"
import "ace-builds/src-noconflict/mode-yaml"
import {dump as toYaml, load as fromYaml} from "js-yaml";

export default {
	props:["value"],
	mounted(){
		//  editor
		this.editor = Ace.edit(this.$el.id);
		this.editor.setOptions({
			showPrintMargin: false,
			highlightActiveLine: true,
			tabSize: 4,
			wrap: true,
			fontSize: 14,
			fontFamily: `'Roboto Mono Regular', monospace`,
		});
		this.editor.setReadOnly(false);
		this.editor.session.setMode("ace/mode/yaml");
		this.editor.session.setUseWorker(false);
	},
	watch: {
		value(newVal) {
			this.value = newVal;
			this.editor.setValue(toYaml(this.value), -1);
		}
	}
}
</script>
