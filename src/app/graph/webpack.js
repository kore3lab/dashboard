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
      mesh: "./src/app/graph/core/mesh/graph.ts",
      topology: "./src/app/graph/core/topology/graph.ts",
      rbac: "./src/app/graph/core/rbac/graph.ts",
    },
    resolve: {
      extensions: [".ts", ".js", ".scss"],
    },
    output: {
      path: path.resolve(__dirname, "../frontend/static/acorn-graph"),
      library: ["acorn", "graph"],
      libraryTarget: "umd",
      filename: "acorn.graph.[name].js",
      globalObject: "this",
      publicPath: "/graph",
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
        template: "src/app/graph/mesh.html",
        filename: "mesh.html",
      }),
      new HtmlWebpackPlugin({
        chunks: ["topology"],
        template: "src/app/graph/topology.html",
        filename: "topology.html",
      }),
      new HtmlWebpackPlugin({
        chunks: ["rbac"],
        template: "src/app/graph/rbac.html",
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
        "/api/*": "http://localhost:3001",
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
