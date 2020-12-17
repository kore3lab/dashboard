export default {
  ssr: false,
  srcDir: "src/app/frontend/",
  head: {
    title: process.env.npm_package_name || "",
    meta: [
      { charset: "utf-8" },
      { name: "viewport", content: "width=device-width, initial-scale=1" },
      {
        hid: "description",
        name: "description",
        content: process.env.npm_package_description || "",
      },
    ],
    link: [
      { rel: "icon", type: "image/x-icon", href: "/favicon.ico" },
      {
        rel: "stylesheet",
        href:
          "https://fonts.googleapis.com/css?family=Source+Sans+Pro:300,400,400i,700",
      },
    ],
    bodyAttrs: {
      class: "hold-transition sidebar-mini layout-fixed",
    },
  },
  css: ["~/assets/css/app.css"],
  plugins: [
    { src: "~/plugins/json-tree", mode: "client" },
    { src: "~/plugins/mixin-common-methods", mode: "client" },
  ],
  buildModules: ["@nuxt/typescript-build"],
  modules: ["bootstrap-vue/nuxt", "@nuxtjs/axios"],
  build: {
    extend(config, ctx) {},
    babel: {
      compact: true, // for build ERROR "bootstrap-vue icons.js as it exceeds the max of 500KB."
    },
  },
  // using "this.$config.backendPort",...
  publicRuntimeConfig: {
    backendPort: process.env.BACKEND_PORT || "3001",
    dashboardPort: process.env.DASHBOARD_PORT || "9090",
    kialiPort: process.env.KIALI_PORT || "20001",
  },
};
