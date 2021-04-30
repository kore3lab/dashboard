import Vue from "vue";

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
		/**
		 * timestamp를 day,hour,minute,second로 구분 봔환함
		 *
		 * @param {date} timestamp 변환할 date 값
		 * @return {{str: string, elapsedTime: number}} timestamp의 day/hour/minute/second 값으로 변환하여 반환함
		 */
		getElapsedTime(timestamp) {
			const dt = Date.parse(timestamp);
			const elapsedTime = new Date() - dt;

			const second = Math.floor((elapsedTime % (1000 * 60)) / 1000);
			const minute = Math.floor((elapsedTime % (1000 * 60 * 60)) / (1000 * 60));
			const hour = Math.floor(
				(elapsedTime % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60)
			);
			// const day = Math.floor((elapsedTime % (1000*60*60*24*30)) / (1000*60*60*24))
			const days = Math.floor(elapsedTime / (1000 * 60 * 60 * 24));
			// const month = Math.floor((elapsedTime % (1000*60*60*24*30*12)) / (1000*60*60*24*30))
			// const year = Math.floor(elapsedTime  / (1000*60*60*24*30*12))

			let str = "";
			// if(year > 0) str += `${year}y`
			// if(month > 0) str += `${month}m`
			// if(day > 0) str += `${day}d`
			if (days > 0) {
				str += `${days}d`;
				if (days >= 10) return {elapsedTime,str};
			}
			if (hour > 0) {
				str += `${hour}h`;
				if (days < 10 && days > 0) return {elapsedTime,str};
			}
			if (minute > 0) {
				if (days > 0 || hour > 0) return {elapsedTime,str};
				str += `${minute}m`;
			}

			if (second > 0) {
				if (hour > 0 || minute > 9) return {elapsedTime,str};
				str += `${second}s`;
			}
			return {elapsedTime,str};
		},
		unitsToBytes(value) {
			const base = 1024;
			const suffixes = ["K", "M", "G", "T", "P", "E"];

			if (!suffixes.some(suffix => value.includes(suffix))) {
				return parseFloat(value);
			}

			const suffix = value.replace(/[0-9]|i|\./g, "");
			const index = suffixes.indexOf(suffix);

			return parseInt(
				(parseFloat(value) * Math.pow(base, index + 1)).toFixed(1)
			);
		},
		cpuUnitsToNumber(cpu) {
			const thousand = 1000;
			const million = thousand * thousand;
			const shortBillion = thousand * million;

			const cpuNum = parseInt(cpu);
			if (cpu.includes("k")) return cpuNum * thousand;
			if (cpu.includes("m")) return cpuNum / thousand;
			if (cpu.includes("u")) return cpuNum / million;
			if (cpu.includes("n")) return cpuNum / shortBillion;

			return parseFloat(cpu);
		},
		metricUnitsToNumber(value) {
			const base = 1000;
			const suffixes = ["k", "m", "g", "t", "p"];

			const suffix = value.toLowerCase().slice(-1);
			const index = suffixes.indexOf(suffix);

			return parseInt(
				(parseFloat(value) * Math.pow(base, index + 1)).toFixed(1)
			);
		},
		memoryUnitsToNumber(value) {
			const base = 1024;
			const suffixes = ["ki", "mi", "gi","ti","pi"]

			const suffix = value.toLowerCase().slice(-2);
			const index = suffixes.indexOf(suffix);

			return (parseFloat(value) * Math.pow(base, index + 1))
		},
		cpuRL(cpu) {
			if(!cpu) return 0
			return this.cpuUnitsToNumber(cpu)
		},
		memoryRL(memory) {
			if(!memory) return 0
			return this.memoryUnitsToNumber(memory)/(1024 * 1024)
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
		getEvents(uid) {
			let events = [];
			this.$axios.get(this.getApiUrl('events.k8s.io','events',''))
				.then( resp => {
					for(let i=0; i<resp.data.items.length; i++) {
						if(resp.data.items[i].regarding.uid === uid) {
							events.unshift({
								name: resp.data.items[i].note || "-",
								source: resp.data.items[i].deprecatedSource.host || resp.data.items[i].deprecatedSource.component || "undefined",
								count: resp.data.items[i].deprecatedCount || "-",
								subObject: resp.data.items[i].regarding.fieldPath || "-",
								lastSeen: resp.data.items[i].deprecatedLastTimestamp || "-",
								type: resp.data.items[i].type === "Warning"? "text-danger" : "text-secondary",
							})
						}
					}
				})
			return events
		},
		toStatus(deletionTimestamp, status) {
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
		toReady(status, spec) {
			let containersReady = 0
			let containersLength = 0
			if ( spec.containers ) containersLength = spec.containers.length
			if ( status.containerStatuses ) containersReady = status.containerStatuses.filter(el => el.ready).length
			return `${containersReady}/${containersLength}`
		},
		sorted(val) {
			if(val) {
				let temp = [];
				for (let i = 0; i < val.length; i++) {
					for (let j = 0; j < val.length; j++) {
						if (val[i].idx < val[j].idx) {
							temp = val[i]
							val[i] = val[j]
							val[j] = temp
						}
					}
				}
			}
			return val
		},
		stringifyLabels(label) {
			if(!label) return [];

			return Object.entries(label).map(([name, value]) => `${name}=${value}`);
		},
		getController(ref) {
			if (!ref) return
			let or = ref[0] ? ref[0] : ref
			let len = or.kind.length
			let k;
			if(or.kind[len-1] === 's') k = (or.kind).toLowerCase() + 'es'
			else k = (or.kind).toLowerCase() + 's'
			let g = (or.apiVersion).split('/')
			if (g.length === 2) {
				return { g: g[0], k: k }
			} else {
				return { g: '', k: k}
			}
		},
		getApiUrl(group, rscName, namespace, name) {
			if(!namespace) namespace = ''
			let resource = this.resources()[group][rscName];
			if (resource) {
				let url
				if(namespace) {
					url = `/raw/clusters/${this.currentContext()}/${group ? "apis" : "api"}/${resource.groupVersion}/${resource.namespaced ? "namespaces/" + namespace + "/" : ""}${resource.name}`;
				}else {
					url = `/raw/clusters/${this.currentContext()}/${group ? "apis" : "api"}/${resource.groupVersion}/${resource.name}`;
				}
				return name ? `${url}/${name}` : url;
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
		isNamespace(_) {
			if (_) this.$store.commit("setIsNamespace",_);
			return this.$store.getters["getIsNamespace"];
		}
	},
});