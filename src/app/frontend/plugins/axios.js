export default function ({ $axios, $config, route, redirect }) {

	if ($config.nodeEnv == "development") {
		$axios.onRequest(config => {
			config.baseURL = `${location.protocol}//${location.hostname}:${$config.backendPort}`;
		})
	}

	$axios.onError(err => {

		if(route.fullPath == "/login") {
			return Promise.reject(err);
		} else {
			if (err.response.status == 401) {
				redirect("/login")
			} else {
				return Promise.reject(err);
			}
		}

	})
}