const TerserPlugin = require("terser-webpack-plugin");
const Version = "0.3.1";

module.exports = {
  mode: "production",
  entry: './src/',
    resolve: {
        extensions: [".ts", ".tsx", ".js", ".jsx"]
    },
  output: {
    filename: `main(${Version}).js`,
    path: __dirname + '/../../release/public/js/',
    library: 'FormLibrary',
    libraryTarget:'umd',
    umdNamedDefine: true
  },
  optimization: {
    minimizer: [
      new TerserPlugin({
        terserOptions: {
          keep_fnames: true,
        },
      }),
    ],
  },
    module: {
        rules: [
            {
                test: /\.ts(x?)$/,
                exclude: /node_modules/,
                use: [
                    {
                        loader: "ts-loader"
                    }
                ]
            },
            {
            test: /\.ts(x?)$/,
            enforce: 'pre',
            exclude: /(node_modules|bower_components|\.spec\.js)/,
            use: [
              {
                loader: 'webpack-strip-block'
              }
            ]
          },
        ],
    }
};
