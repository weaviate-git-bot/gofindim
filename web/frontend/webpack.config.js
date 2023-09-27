const webpack = require("webpack"); // only add this if you don't have yet
const path = require("path");
const HtmlWebpackPlugin = require("html-webpack-plugin");
require("dotenv").config({ path: "./.env" });
module.exports = {
  entry: { app: path.join(__dirname, "src", "index.tsx") },
  module: {
    rules: [
      {
        test: /\.tsx?$/,
        loader: require.resolve("babel-loader"),
        exclude: /node_modules/,
      },
    ],
  },
  resolve: {
    extensions: [".tsx", ".ts", ".js", ".json"],
    modules: [path.resolve(__dirname, "src"), "node_modules"],
  },
  output: {
    filename: "bundle.js",
    path: path.resolve(__dirname, "dist"),
    publicPath: '', // The publicPath is set to an empty string

  },
  target: "web",
  devServer: {
    static: path.join(__dirname, "public"),
    compress: true,
    hot: true,
    open: false,
    port: 4000,
    historyApiFallback: true,
  },
  plugins: [
    new HtmlWebpackPlugin({
      template: path.join(__dirname, "public", "index.html"),
    }),
    new webpack.DefinePlugin({
      "process.env": JSON.stringify(process.env),
    }),
    new webpack.ProvidePlugin({
      "config": path.join(__dirname, "src", "config.ts"),
    }),
  ],
};
