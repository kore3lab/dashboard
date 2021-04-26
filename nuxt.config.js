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
		middleware: ["authenticated"]
	},
	plugins: [
		{ src: "~/plugins/mixin-common-methods", mode: "client" },
		{ src: "~/plugins/mixin-viewpage", mode: "client" },
		{ src: "~/plugins/admin-lte", mode: "client" },
		{ src: '~/plugins/axios', mode: "client" }
	],
	buildModules: ["@nuxt/typescript-build"],
	modules: ["bootstrap-vue/nuxt", "@nuxtjs/axios", "cookie-universal-nuxt"],
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
	publicRuntimeConfig: {
		nodeEnv: process.env.NODE_ENV,						// production or development
		backendPort: process.env.BACKEND_PORT || "3001",	// only development mode
		itemsPerPage: process.env.ITEMS_PER_PAGE || "10"
	},
	build: {
		extend(config, ctx) {},
		babel: {
			compact: true, // for build ERROR "bootstrap-vue icons.js as it exceeds the max of 500KB."
		},
	},
};
