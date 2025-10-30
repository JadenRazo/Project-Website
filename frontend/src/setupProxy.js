const { createProxyMiddleware } = require('http-proxy-middleware');

module.exports = function(app) {
  // DevPanel Auth API
  app.use(
    '/api/v1/auth/admin',
    createProxyMiddleware({
      target: 'http://localhost:8081',
      changeOrigin: true,
    })
  );

  // DevPanel API
  app.use(
    '/api/v1/devpanel',
    createProxyMiddleware({
      target: 'http://localhost:8081',
      changeOrigin: true,
    })
  );

  // Status API (Main backend)
  app.use(
    '/api/v1/status',
    createProxyMiddleware({
      target: 'http://localhost:8080',
      changeOrigin: true,
    })
  );

  // Messaging API
  app.use(
    '/api/messaging',
    createProxyMiddleware({
      target: 'http://localhost:8082',
      changeOrigin: true,
    })
  );

  // URL Shortener
  app.use(
    '/s',
    createProxyMiddleware({
      target: 'http://localhost:8083',
      changeOrigin: true,
    })
  );

  // Projects API (Main backend)
  app.use(
    '/api/v1/projects',
    createProxyMiddleware({
      target: 'http://localhost:8080',
      changeOrigin: true,
    })
  );

  // Code Stats API (Main backend)
  app.use(
    '/api/v1/code',
    createProxyMiddleware({
      target: 'http://localhost:8080',
      changeOrigin: true,
    })
  );

  // WebSocket for messaging
  app.use(
    '/ws',
    createProxyMiddleware({
      target: 'ws://localhost:8082',
      ws: true,
      changeOrigin: true,
    })
  );
};