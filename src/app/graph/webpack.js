const path = require("path");
const HtmlWebpackPlugin = require("html-webpack-plugin");
// const CleanWebpackPlugin = require("clean-webpack-plugin");

/**
 * 참고
 * 		output: https://webpack.js.org/configuration/output/
 */
module.exports = [
  {
    entry: {
      mesh: path.join(__dirname, "/core/mesh/graph.ts"),
      topology: path.join(__dirname, "/core/topology/graph.ts"),
      rbac: path.join(__dirname, "/core/rbac/graph.ts"),
    },
    resolve: {
      extensions: [".ts", ".js", ".scss"],
    },
    output: {
      // path: path.resolve(__dirname, "../frontend/static/acorn-graph"),
      // publicPath: "/graph",
      path: path.resolve(__dirname, "dist"),
      library: ["acorn", "graph"],
      libraryTarget: "umd",
      filename: "acorn.graph.[name].js",
      globalObject: "this",
    },
    module: {
      rules: [
        { test: /\.tsx?$/, use: "ts-loader", exclude: /node_modules/ },
        { test: /\.css$/, use: ["style-loader", "css-loader"] },
      ],
    },
    plugins: [
      // new CleanWebpackPlugin(),
      new HtmlWebpackPlugin({
        chunks: ["mesh"],
        template: "./examples/mesh.html",
        filename: "mesh.html",
      }),
      new HtmlWebpackPlugin({
        chunks: ["topology"],
        template: "./examples/topology.html",
        filename: "topology.html",
      }),
      new HtmlWebpackPlugin({
        chunks: ["rbac"],
        template: "./examples/rbac.html",
        filename: "rbac.html",
      }),
    ],
    devtool: "source-map",
    devServer: {
      historyApiFallback: true,
      compress: true,
      host: "0.0.0.0",
      port: 3002,
      proxy: {
        "/api/clusters/apps-05/*": "http://localhost:3002/examples/data",
      },
    },
    node: {
      fs: "empty",
      net: "empty",
      tls: "empty",
      child_process: "empty",
    },
  },
];
