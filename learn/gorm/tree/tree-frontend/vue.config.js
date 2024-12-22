// vue.config.js
module.exports = {
  devServer: {
    proxy: {
      '/api': {
        target: 'http://localhost:8080', // 后端服务器地址
        changeOrigin: true,
        pathRewrite: {
          '^/api': '' // 可选：将 /api 前缀去掉
        }
      }
    }
  }
}
