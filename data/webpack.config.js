const path = require("path");

module.exports = {
	devServer: {
		contentBase: path.join(__dirname, 'dist'),
		port: 3000
	},
	entry: "./js/index.js",
	output: {
		filename: "bundle.js",
		path: path.resolve(__dirname, "dist")
	},
	module: {
		rules: [
			{
				test: /\.(js|jsx)$/,
				exclude: /node_modules/,
				use: {
					loader: "babel-loader"
				}

			}

		]

	},
	// TODO: Add Prod Rules.
	mode: "development"
};
