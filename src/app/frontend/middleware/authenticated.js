export default function(context) {

	if(context.route.fullPath == "/login") return

	if( !context.app.$cookies.get("refresh-token") ) {
		context.redirect("/login")
		return
	}

	context.app.$axios.interceptors.response.use(
		resp=> {
			return resp
		}, err => {
			if(context.route.fullPath == "/login") {
				return Promise.reject(err);
			} else {
				if (err.response.status == 401) {
					context.redirect("/login")
				} else {
					return Promise.reject(err);
				}
			}
		})
}
 