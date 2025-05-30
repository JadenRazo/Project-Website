// Performance configuration and optimization utilities

export interface PerformanceConfig {
  lazyLoadingThreshold: number;
  imageLoadingDelay: number;
  prefetchDelay: number;
  minimumLoadingTime: number;
  enableServiceWorker: boolean;
  enableBundleSplitting: boolean;
  enablePreloading: boolean;
}

export const defaultPerformanceConfig: PerformanceConfig = {
  lazyLoadingThreshold: 0.1,
  imageLoadingDelay: 200,
  prefetchDelay: 150,
  minimumLoadingTime: 300,
  enableServiceWorker: false, // Disabled by default in development
  enableBundleSplitting: true,
  enablePreloading: true
};

// Environment-specific configurations
export const getPerformanceConfig = (): PerformanceConfig => {
  const isProduction = process.env.NODE_ENV === 'production';
  const isDevelopment = process.env.NODE_ENV === 'development';
  
  if (isProduction) {
    return {
      ...defaultPerformanceConfig,
      enableServiceWorker: true,
      minimumLoadingTime: 200,
      prefetchDelay: 100
    };
  }
  
  if (isDevelopment) {
    return {
      ...defaultPerformanceConfig,
      minimumLoadingTime: 100,
      prefetchDelay: 300
    };
  }
  
  return defaultPerformanceConfig;
};

// Bundle splitting configuration for Webpack
export const bundleSplittingConfig = {
  chunks: 'all' as const,
  cacheGroups: {
    // Vendor chunk for node_modules
    vendor: {
      test: /[\\/]node_modules[\\/]/,
      name: 'vendors',
      chunks: 'all' as const,
      priority: 20
    },
    
    // Common chunk for shared components
    common: {
      name: 'common',
      minChunks: 2,
      chunks: 'all' as const,
      priority: 10,
      reuseExistingChunk: true
    },
    
    // React and React-DOM in separate chunk
    react: {
      test: /[\\/]node_modules[\\/](react|react-dom)[\\/]/,
      name: 'react',
      chunks: 'all' as const,
      priority: 30
    },
    
    // Animation libraries
    animations: {
      test: /[\\/]node_modules[\\/](framer-motion|@react-spring|gsap)[\\/]/,
      name: 'animations',
      chunks: 'all' as const,
      priority: 25
    },
    
    // UI and styling libraries
    ui: {
      test: /[\\/]node_modules[\\/](styled-components|@emotion)[\\/]/,
      name: 'ui',
      chunks: 'all' as const,
      priority: 25
    },
    
    // Three.js and related 3D libraries
    threejs: {
      test: /[\\/]node_modules[\\/](three|@react-three)[\\/]/,
      name: 'threejs',
      chunks: 'all' as const,
      priority: 25
    }
  }
};

// Resource hints for preloading critical resources
export const criticalResourceHints = [
  { rel: 'preconnect', href: 'https://fonts.googleapis.com' },
  { rel: 'preconnect', href: 'https://fonts.gstatic.com', crossOrigin: 'anonymous' },
  { rel: 'dns-prefetch', href: '//api.github.com' }
];

// Service Worker configuration
export const serviceWorkerConfig = {
  swSrc: 'src/serviceWorker.ts',
  dontCacheBustURLsMatching: /\.\w{8}\./,
  exclude: [/\.map$/, /manifest$/, /\.htaccess$/],
  navigateFallback: '/index.html',
  navigateFallbackBlacklist: [/^\/_/, /\/[^/?]+\.[^/]+$/],
  runtimeCaching: [
    {
      urlPattern: /^https:\/\/fonts\.googleapis\.com\//,
      handler: 'StaleWhileRevalidate',
      options: {
        cacheName: 'google-fonts-stylesheets'
      }
    },
    {
      urlPattern: /^https:\/\/fonts\.gstatic\.com\//,
      handler: 'CacheFirst',
      options: {
        cacheName: 'google-fonts-webfonts',
        expiration: {
          maxAgeSeconds: 60 * 60 * 24 * 365 // 1 year
        }
      }
    },
    {
      urlPattern: /\.(?:png|jpg|jpeg|svg|gif|webp)$/,
      handler: 'CacheFirst',
      options: {
        cacheName: 'images',
        expiration: {
          maxEntries: 100,
          maxAgeSeconds: 60 * 60 * 24 * 30 // 30 days
        }
      }
    }
  ]
};

// Performance monitoring utilities
export const performanceObserver = {
  // Observe Largest Contentful Paint
  observeLCP: (callback: (entry: PerformanceEntry) => void) => {
    if ('PerformanceObserver' in window) {
      try {
        const observer = new PerformanceObserver((list) => {
          const entries = list.getEntries();
          const lastEntry = entries[entries.length - 1];
          callback(lastEntry);
        });
        observer.observe({ entryTypes: ['largest-contentful-paint'] });
        return observer;
      } catch (e) {
        console.warn('LCP observation not supported:', e);
      }
    }
    return null;
  },

  // Observe First Input Delay
  observeFID: (callback: (entry: PerformanceEntry) => void) => {
    if ('PerformanceObserver' in window) {
      try {
        const observer = new PerformanceObserver((list) => {
          const entries = list.getEntries();
          entries.forEach(callback);
        });
        observer.observe({ entryTypes: ['first-input'] });
        return observer;
      } catch (e) {
        console.warn('FID observation not supported:', e);
      }
    }
    return null;
  },

  // Observe Cumulative Layout Shift
  observeCLS: (callback: (entry: PerformanceEntry) => void) => {
    if ('PerformanceObserver' in window) {
      try {
        const observer = new PerformanceObserver((list) => {
          const entries = list.getEntries();
          entries.forEach(callback);
        });
        observer.observe({ entryTypes: ['layout-shift'] });
        return observer;
      } catch (e) {
        console.warn('CLS observation not supported:', e);
      }
    }
    return null;
  }
};

export default {
  config: getPerformanceConfig(),
  bundleSplitting: bundleSplittingConfig,
  resourceHints: criticalResourceHints,
  serviceWorker: serviceWorkerConfig,
  observer: performanceObserver
};