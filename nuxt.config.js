export default {
	ssr: false,
	srcDir: "src/app/frontend/",
	head: {
		title: process.env.npm_package_name || "",
		meta: [
			{ charset: "utf-8" },
			{ name: "viewport", content: "width=device-width, initial-scale=1" },
			{ name: "description", content: process.env.npm_package_description || "", hid: "description"},
		],
		link: [
			{ rel: "icon", type: "image/svg+xml", href: "/favicon.svg" },
			{ rel: "stylesheet", href: "https://fonts.googleapis.com/css?family=Source+Sans+Pro:300,400,400i,700",},
		],
		bodyAttrs: {
			class: "hold-transition sidebar-mini layout-fixed",
		},
	},
	css: [
		"@fortawesome/fontawesome-free/css/all.css",
		"~/plugins/adminLte/css/adminlte.min.css",
		"~/assets/css/app.css",
	],
	router: {
		middleware: ["auth"]
	},
	plugins: [
		{ src: "~/plugins/mixin-common-methods", mode: "client" },
		{ src: "~/plugins/mixin-viewpage", mode: "client" },
		{ src: "~/plugins/admin-lte", mode: "client" },
		{ src: '~/plugins/axios', mode: "client" }
	],
	buildModules: ["@nuxt/typescript-build"],
	modules: ["bootstrap-vue/nuxt", "@nuxtjs/axios", "cookie-universal-nuxt","@nuxtjs/auth-next"],
	bootstrapVue: {
		bootstrapCSS: false, 
		icons: true,
	},
	axios: {
		baseURL: "/",
		credentials: true,
		headers: {
			common: { "Accept": "application/json"}
		}
	},
	auth: {
		redirect: {
			login: "/login",	//User will be redirected to this path if login is required.
			logout: "/login",	//User will be redirected to this path if after logout, current route is protected.
			home: "/"			//User will be redirected to this path after login. (rewriteRedirects will rewrite this path)
		},
		localStorage: {
			prefix: false		// false or "auth."
		},
		cookie: {
			prefix: "auth."
		},
		strategies: {
			cookie: {
				user: {
					property: "user"
				},
				endpoints: {
					login:	{ url: "/api/auth/login", method: "POST" },
					logout:	{ url: "/api/auth/logout", method: "GET" },
					user:	{ url: '/api/auth/user', method: "GET" },
				}
			},
			local: {
				scheme: "refresh",
				token: {
					property: 'token',
					maxAge: 60 * 15,
					type: false
				},
				refreshToken: {
					property: 'refreshToken',
					data: 'refreshToken',
					maxAge: 60 * 60 * 24 * 7
				},
				user: {
					property: 'user'
				},
				endpoints: {
					login:		{ url: "/api/auth/login", method: "POST" },
					logout:		{ url: "/api/auth/logout", method: "GET" },
					refresh:	{ url: '/api/auth/token/refresh', method: "POST" },
					user:		{ url: '/api/auth/user', method: "GET" },
				}
			}
		}
	},
	publicRuntimeConfig: {
		nodeEnv: process.env.NODE_ENV,						// production or development
		backendPort: process.env.BACKEND_PORT || "3001",	// only development mode
		terminalPort: process.env.TERMINAL_PORT || "3003",	// only development mode
		version: JSON.stringify(require('./package.json').version)
	},
	build: {
		extend(config, ctx) {},
		babel: {
			compact: true, // for build ERROR "bootstrap-vue icons.js as it exceeds the max of 500KB."
		},
	},
};
