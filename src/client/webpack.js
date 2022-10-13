
const Version = "0.2.0";

module.exports = {
  mode: "production",
  entry: './src/index.ts',
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
