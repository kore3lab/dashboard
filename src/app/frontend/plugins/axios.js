export default function ({ $axios, $config, app, store, route, redirect }) {

	if ($config.nodeEnv === "development") {
		$axios.defaults.baseURL = `${location.protocol}//${location.hostname}:${$config.backendPort}`;
	}

	$axios.onError(err => {
		if(route.fullPath === "/login") {
			return Promise.reject(err);
		} else {
			if (err.response.status === 401) {
				app.$auth.logout();
				redirect("/login")
			} else {
				return Promise.reject(err);
			}
		}

	})
}