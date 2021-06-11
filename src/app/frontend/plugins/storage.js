// LocalStorage plugin
//---------------------------------
//"_global = {
//	"itmesPerPage" : "20"
//}
//"_local" = {
//	"/workload/pod.list" : {
//		"columns.grdSheet1" : "name|namespace|age|"
//	}
//}
export default ({ }, inject) => {

	function _storage() {

		// get, set
		this.set = (key, obj) => {
			localStorage.setItem(key, JSON.stringify(obj));
		}
		this.get = (key, defaultValue) => {
			let obj = JSON.parse(localStorage.getItem(key));
			return obj ? obj: defaultValue;
		}

		// local data, global data
		this._local = this.get("_local", {});
		this._global = this.get("_global", {});

		//local
		this.local = {
			get: (key, defaultValue) => {
				return this._local[location.pathname] ? (this._local[location.pathname][key] ? this._local[location.pathname][key]: defaultValue? defaultValue:""): (defaultValue? defaultValue:"");
			},
			set: (key, value) => {
				if(!this._local[location.pathname]) this._local[location.pathname] = {}
				this._local[location.pathname][key] = value;
				this.set("_local", this._local);
			},
			delete: (key) => {
				if (key && this._local[location.pathname] && this._local[location.pathname][key]) delete this._local[location.pathname][key];
				else return
				this.set("_local", this._local);
			}
		}
		// global
		this.global = {
			get: (key, defaultValue) => {
				return this._global[key] ? this._global[key]: (defaultValue? defaultValue:"");
			},
			set: (key, value) => {
				this._global[key] = value;
				this.set("_global", this._global);
			},
			delete: (key) => {
				if(this._global[key]) delete this._global[key];
				else return
				this.set("_global", this._global);
			}
		}

	}
	inject("storage", new _storage());
}

