module.exports = {
  devServer: {
    socket: 'socket',
    port: 8081,
    proxy: {
      '/session': {
          target: 'http://localhost:8380',
          changeOrigin: true,
          pathRewrite: {
              '^/session': ''
          }
      },
      '/gateway': {
          target: 'http://localhost:8980',
          changeOrigin: true,
          pathRewrite: {
              '^/gateway': ''
          }
      }
    }
  }
}
