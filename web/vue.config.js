const fs = require('fs')

module.exports = {
    publicPath: './',
    devServer: {
        port: 25702,
        host: '0.0.0.0',
        public: 'kb:25702',
        allowedHosts: ['kb'],
        disableHostCheck: true,
        compress: true,
        hot: true,
        historyApiFallback: true,
        proxy: {
          '/api/': {
            target: 'http://localhost:8006',
            changeOrigin: true,
            ws: true,
            pathRewrite: {
              '^/api/': '/api/'
            }
          },
        }
    },
    configureWebpack: config => {
      if (process.env.NODE_ENV === 'production') {
        // 为生产环境修改配置...
      } else {
        // 为开发环境修改配置...
      }
    },
    chainWebpack: config => {
      config.resolve.alias.set('vue-i18n', 'vue-i18n/dist/vue-i18n.cjs.js')
      config.module.rule("i18n").resourceQuery(/blockType=i18n/)
        .type('javascript/auto')
        .use("i18n")
        .loader("@intlify/vue-i18n-loader")
        .end()
        .use('yaml')
        .loader('yaml-loader')
        .end()
    },
    css: {
      extract: false,
      loaderOptions: {
        scss: {
          additionalData: `@import "~@/styles/theme-variables.scss";`
          // prependData: fs.readFileSync('src/styles/theme-variables.scss', 'utf-8')
        }
      }
    }
  }