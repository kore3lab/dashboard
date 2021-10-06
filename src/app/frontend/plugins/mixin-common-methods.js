import Vue from "vue";
const CONSTRANTS = {
	"ITEMS_PER_PAGE" : [{text:"10",value:10}, {text:"20",value:20}, {text:"30",value:30}, {text:"50",value:50}, {text:"100",value:100}],
	"CHART_BORDER_COLOR" : {
		requests: "#999",
	},
	"CHART_BG_COLOR" : {
		cpu: "rgba(119,149,233,0.5)",
		memory:	"rgba(179,145,208,0.5)",
		requests: "rgba(108, 117, 125,0.2)",
		background: "#fff",
	}
}

Vue.filter("formatNumber", (value) => {
	try {
		return Number(value).toLocaleString();
	} catch {
		return "NaN";
	}
});

Vue.mixin({
	methods: {
		toast(msg, variant) {
			if (!variant) variant = "info";
			this.$bvToast.toast(msg, { title: variant, noCloseButton: false, variant: variant, autoHideDelay: 4000});
		},
		mesbox(msg) {
			this.$bvModal.msgBoxOk(msg, { title: "", variant: "info", buttonSize: "sm", footerClass: "p-1"});
		},
		confirm(msg, callback) {
			this.$bvModal
				.msgBoxConfirm(msg, { title: "", variant: "info", buttonSize: "sm", footerClass: "p-1"})
				.then(callback)
				.catch((_) => {});
		},
		msghttp(error) {
			if (error.response && error.response.data && error.response.data.message ) {
				this.toast(error.response.data.message, "warning");
			} else {
				this.toast(error.message, "danger");
			}
		},
		var(name) {
			return CONSTRANTS[name];
		},
		getElapsedTime(timestamp) {
			let elapsedTime = new Date() - Date.parse(timestamp);

			let second = Math.floor((elapsedTime % (1000 * 60)) / 1000);
			let minute = Math.floor((elapsedTime % (1000 * 60 * 60)) / (1000 * 60));
			let hour = Math.floor((elapsedTime % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
			let days = Math.floor(elapsedTime / (1000 * 60 * 60 * 24));
			let str = "";

			if (days > 0) {
				str += `${days}d`;
				if (days >= 10) return str;
			}
			if (hour > 0) {
				str += `${hour}h`;
				if (days < 10 && days > 0) return str;
			}
			if (minute > 0) {
				if (days > 0 || hour > 0) return str;
				str += `${minute}m`;
			}

			if (second > 0) {
				if (hour > 0 || minute > 9) return str;
				str += `${second}s`;
			}
			return str;
		},
		getTimestampString(timestamp) {
			let dt = Date.parse(timestamp);
			let seconds = Math.floor((new Date() - dt) / 1000);

			var interval = seconds / 31536000;
			if (interval > 1) return Math.floor(interval) + " years";

			interval = seconds / 2592000;
			if (interval > 1) return Math.floor(interval) + " months";

			interval = seconds / 86400;
			if (interval > 1) return Math.floor(interval) + " days";

			interval = seconds / 3600;
			if (interval > 1) return Math.floor(interval) + " hours";
			interval = seconds / 60;
			if (interval > 1) return Math.floor(interval) + " minutes";

			return Math.floor(seconds) + " seconds";
		},
		formatNumber(value, decimal) {
			try {
				return Number(decimal ? Number.parseFloat(value.toFixed(decimal)): value).toLocaleString();
			} catch {
				return "NaN";
			}
		},
		toPodStatus(deletionTimestamp, status) {
			// 삭제
			if (deletionTimestamp) {
				return { "value": "Terminating", "style": "text-secondary"}
			}

			// Pending
			if (!status.containerStatuses) {
				if(status.phase === "Failed") {
					return { "value": status.phase, "style": "text-danger"}
				} else {
					return { "value": status.phase, "style": "text-warning"}
				}
			}

			// [if]: Running, [else]: (CrashRoofBack / Completed / ContainerCreating)
			if(status.containerStatuses.filter(el => el.ready).length === status.containerStatuses.length) {
				const state = Object.keys(status.containerStatuses.find(el => el.ready).state)[0]
				return { "value": state.charAt(0).toUpperCase() + state.slice(1), "style": "text-success" }
			}
			else {
				const state = status.containerStatuses.find(el => !el.ready).state
				let style = "text-secondary"
				if ( state[Object.keys(state)].reason === "Completed") style = "text-success"
				if ( state[Object.keys(state)].reason === "Error") style = "text-danger"
				return { "value": state[Object.keys(state)].reason, "style": style }
			}
		},
		stringifyLabels(label) {
			if(!label) return [];
			try {
				return Object.entries(label).map(([name, value]) => `${name}=${value}`);
			} catch (e) {
				console.log(e);
			}
			return [];
		},
		// get a resource-info(selfLink)( metdata.ownerReferences, ...  )
		getResource(reference) {
			let resource = {};
			if(reference) {
				let version = reference.apiVersion?reference.apiVersion.split('/'):[];
				resource = {
					kind: reference.kind,
					group: (version.length>1)? version[0]: "",
					resource:`${reference.kind.toLowerCase()}${reference.kind.endsWith('s')?'es':'s'}`
				};
				if (reference["metadata"]) {
					resource["name"] = reference.metadata.name;
					resource["namespace"] = reference.metadata["namespace"]? reference.metadata["namespace"]: "";
				} else {
					resource["name"] = reference.name;
					resource["namespace"] = reference["namespace"]? reference["namespace"]: "";
				}
				return resource;
			} else {
				return {}
			}
		},
		getApiUrl(group, rscName, namespace, name, query, log) {
			if(!namespace) namespace = '';
			if(!log) log = false;
			let resource = this.resources()[group][rscName];
			if (resource) {
				let url
				if(namespace) {
					url = `/raw/clusters/${this.currentContext()}/${group ? "apis" : "api"}/${resource.groupVersion}/${resource.namespaced ? "namespaces/" + namespace + "/" : ""}${resource.name}`;
				}else {
					url = `/raw/clusters/${this.currentContext()}/${group ? "apis" : "api"}/${resource.groupVersion}/${resource.name}`;
				}
				if(log){
					url = `${url}/${name}/log${query ? '?' + query : ''}`;
				}else {
					url = name ? `${url}/${name}${query ? '?' + query : ''}` : url + (query ? '?' + query : '');
				}
				return url;
			} else {
				return "#";
			}
		},
		// Get currentContext's namespaces
		namespaces(_) {
			if (_) this.$store.commit("setNamespaces", _);
			else return this.$store.getters["getNamespaces"];
		},
		// Get contexts
		contexts(_) {
			if (_) this.$store.commit("setContexts", _);
			else return this.$store.getters["getContexts"];
		},
		// Get currentContext's resources
		resources(_) {
			if (_) this.$store.commit("setResources", _);
			else return this.$store.getters["getResources"];
		},
		// Get a currentContext
		currentContext(_) {
			if (_) this.$store.commit("setCurrentContext", _);
			else return this.$store.getters["getCurrentContext"];
		},
		selectNamespace(_) {
			if(_ === "") this.$store.commit("setSelectNamespace", _);
			if (_) this.$store.commit("setSelectNamespace", _);
			else return this.$store.getters["getSelectNamespace"];
		},
		statusbar(_) {
			if(_ === "") this.$store.commit("setStatusbar", _);
			if (_) this.$store.commit("setStatusbar", _);
			else return this.$store.getters["getStatusbar"];
		},
		contentPadding(_) {
			if(_) this.$store.commit("setContentPadding", _);
			else return this.$store.getters["getContentPadding"];
		}
	},
});