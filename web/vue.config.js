module.exports = {
  devServer: {
    proxy: {
      '^/api': {
        target: 'http://localhost:8080',
        ws: true,
        changeOrigin: true,
        pathRewrite: {
          '^/api': '/' // remove base path
        }
      },
    }
  }
};
