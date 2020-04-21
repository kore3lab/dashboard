const config = {
  backendPort:
    process.env.NODE_ENV !== "production"
      ? process.env.BACKEND_PORT || "3001"
      : "30081",
  dashboardPort:
    process.env.NODE_ENV !== "production"
      ? process.env.DASHBOARD_PORT || "9090"
      : "30090",
  kialiRootUrl:
    process.env.NODE_ENV !== "production"
      ? process.env.KIALI_ROOT_URL || "http://localhost:20001"
      : "http://101.55.69.105:32080",
};
export default {
  mode: "spa",
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
  loading: {
    color: "#0826c0",
    height: "6px",
    failedColor: "#d43c6d",
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
  },
  env: {
    backendPort: config.backendPort,
    dashboardPort: config.dashboardPort,
    kialiRootUrl: config.kialiRootUrl,
  },
};
