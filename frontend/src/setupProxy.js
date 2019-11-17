const proxy = require("http-proxy-middleware");

module.exports = app => {
  app.use(
    "/api",
    proxy({
      target: "http://api:8080",
      changeOrigin: true,
      pathRewrite: {
        "^/api/": "/"
      }
    })
  );
};
