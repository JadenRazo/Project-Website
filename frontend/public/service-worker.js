const CACHE_NAME = 'portfolio-v2';
const urlsToCache = [
  '/',
  '/static/css/main.css',
  '/static/js/bundle.js',
  '/favicon.ico',
  '/manifest.json'
];

self.addEventListener('install', (event) => {
  event.waitUntil(
    caches.open(CACHE_NAME)
      .then((cache) => cache.addAll(urlsToCache))
  );
});

self.addEventListener('fetch', (event) => {
  // Skip caching for video files and env-config.js to prevent loading issues
  const url = new URL(event.request.url);
  const isVideo = url.pathname.endsWith('.mp4') || 
                  url.pathname.endsWith('.webm') || 
                  url.pathname.endsWith('.ogg') || 
                  event.request.headers.get('accept')?.includes('video');
  
  const isEnvConfig = url.pathname.endsWith('env-config.js');
  
  if (isVideo || isEnvConfig) {
    // For video files and env-config.js, always fetch from network
    event.respondWith(fetch(event.request));
    return;
  }

  event.respondWith(
    caches.match(event.request)
      .then((response) => {
        // Cache hit - return response
        if (response) {
          return response;
        }

        return fetch(event.request).then((response) => {
          // Check if valid response
          if (!response || response.status !== 200 || response.type !== 'basic') {
            return response;
          }

          // Clone the response
          const responseToCache = response.clone();

          caches.open(CACHE_NAME)
            .then((cache) => {
              cache.put(event.request, responseToCache);
            });

          return response;
        });
      })
  );
});