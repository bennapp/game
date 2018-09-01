const HtmlWebpackPlugin = require("html-webpack-plugin");

const path = require('path');

module.exports = {
  entry: "./src/client/index.js",
  output: {
    filename: "bundle.js",
    path: path.resolve(__dirname, 'dist')
  },
  node: {
    fs: "empty"
  },
  externals: {
    uws: "uws"
  },
  mode: "development",
  watch: true,
  plugins: [
    new HtmlWebpackPlugin({
      template: "./public/index.html",
    })
  ]
};
