window._env_ = {
  REACT_APP_API_URL: '',  // Empty string will use relative URLs (same domain)
  REACT_APP_WS_URL: `wss://${window.location.host}`,  // Use current domain for WebSocket
  REACT_APP_ENVIRONMENT: 'production'
};