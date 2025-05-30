type ResourceType = 'script' | 'style' | 'image' | 'font' | 'audio' | 'video' | 'document' | 'fetch';

interface PreloadResource {
  href: string;
  as: ResourceType;
  crossOrigin?: 'anonymous' | 'use-credentials';
  type?: string;
  media?: string;
}

class ResourcePreloader {
  private preloadedResources = new Set<string>();
  private preloadPromises = new Map<string, Promise<void>>();

  preload(resource: PreloadResource): Promise<void> {
    const { href, as, crossOrigin, type, media } = resource;
    
    if (this.preloadedResources.has(href)) {
      return this.preloadPromises.get(href) || Promise.resolve();
    }
    
    const promise = new Promise<void>((resolve, reject) => {
      const link = document.createElement('link');
      link.rel = 'preload';
      link.href = href;
      link.as = as;
      
      if (crossOrigin) {
        link.crossOrigin = crossOrigin;
      }
      
      if (type) {
        link.type = type;
      }
      
      if (media) {
        link.media = media;
      }
      
      link.onload = () => resolve();
      link.onerror = () => reject(new Error(`Failed to preload ${href}`));
      
      document.head.appendChild(link);
    });
    
    this.preloadedResources.add(href);
    this.preloadPromises.set(href, promise);
    
    return promise;
  }

  preloadImages(urls: string[]): Promise<void[]> {
    return Promise.all(
      urls.map(url => this.preload({ href: url, as: 'image' }))
    );
  }

  preloadScripts(urls: string[]): Promise<void[]> {
    return Promise.all(
      urls.map(url => this.preload({ href: url, as: 'script' }))
    );
  }

  preloadFonts(fonts: Array<{ url: string; type?: string }>): Promise<void[]> {
    return Promise.all(
      fonts.map(font => 
        this.preload({ 
          href: font.url, 
          as: 'font', 
          type: font.type || 'font/woff2',
          crossOrigin: 'anonymous'
        })
      )
    );
  }

  prefetchResource(href: string): void {
    if (this.preloadedResources.has(href)) return;
    
    const link = document.createElement('link');
    link.rel = 'prefetch';
    link.href = href;
    document.head.appendChild(link);
    
    this.preloadedResources.add(href);
  }

  dns(hostnames: string[]): void {
    hostnames.forEach(hostname => {
      const link = document.createElement('link');
      link.rel = 'dns-prefetch';
      link.href = `//${hostname}`;
      document.head.appendChild(link);
    });
  }

  preconnect(urls: string[]): void {
    urls.forEach(url => {
      const link = document.createElement('link');
      link.rel = 'preconnect';
      link.href = url;
      link.crossOrigin = 'anonymous';
      document.head.appendChild(link);
    });
  }
}

export const preloader = new ResourcePreloader();

// Import images that need preloading
const preloadableImages = {
  headshot: require('../assets/images/headshot.jpg')
};

// Route-based preloading
export const preloadRouteAssets = {
  '/about': () => {
    // Preload the headshot image using the webpack-bundled URL
    if (preloadableImages.headshot) {
      preloader.preloadImages([preloadableImages.headshot]);
    }
  },
  
  '/projects': () => {
    preloader.prefetchResource('/code_stats.json');
    // Preload any project media files
  },
  
  '/contact': () => {
    // Preload contact form related assets
  },
  
  '/devpanel': () => {
    // Preload devpanel specific assets
  },
  
  '/messaging': () => {
    // Preload messaging related assets
  },
  
  '/urlshortener': () => {
    // Preload URL shortener assets
  }
};

// Critical resource preloading
export const preloadCriticalAssets = (): void => {
  // Preload critical fonts
  preloader.preloadFonts([
    { url: '/fonts/primary-font.woff2', type: 'font/woff2' },
    { url: '/fonts/mono-font.woff2', type: 'font/woff2' }
  ]).catch(console.warn);

  // DNS prefetch for external domains
  preloader.dns([
    'fonts.googleapis.com',
    'api.github.com'
  ]);

  // Preconnect to important external resources
  preloader.preconnect([
    'https://fonts.gstatic.com'
  ]);
};

// Smart preloading based on user behavior
export class SmartPreloader {
  private userBehavior: { [route: string]: number } = {};
  private preloadQueue: string[] = [];
  private isIdle = false;

  constructor() {
    this.trackUserBehavior();
    this.setupIdleDetection();
  }

  private trackUserBehavior(): void {
    // Track route visits
    const currentRoute = window.location.pathname;
    this.userBehavior[currentRoute] = (this.userBehavior[currentRoute] || 0) + 1;
  }

  private setupIdleDetection(): void {
    let idleTimer: NodeJS.Timeout;
    
    const resetIdleTimer = () => {
      this.isIdle = false;
      clearTimeout(idleTimer);
      idleTimer = setTimeout(() => {
        this.isIdle = true;
        this.processPreloadQueue();
      }, 2000); // 2 seconds of inactivity
    };

    ['mousedown', 'mousemove', 'keypress', 'scroll', 'touchstart'].forEach(event => {
      document.addEventListener(event, resetIdleTimer, true);
    });

    resetIdleTimer();
  }

  queuePreload(route: string): void {
    if (!this.preloadQueue.includes(route)) {
      this.preloadQueue.push(route);
    }
    
    if (this.isIdle) {
      this.processPreloadQueue();
    }
  }

  private processPreloadQueue(): void {
    if (this.preloadQueue.length === 0) return;
    
    const route = this.preloadQueue.shift();
    if (route && preloadRouteAssets[route as keyof typeof preloadRouteAssets]) {
      try {
        preloadRouteAssets[route as keyof typeof preloadRouteAssets]();
      } catch (error) {
        console.warn(`Failed to preload assets for route ${route}:`, error);
      }
    }
    
    // Process next item in queue after a delay
    if (this.preloadQueue.length > 0) {
      setTimeout(() => this.processPreloadQueue(), 100);
    }
  }

  getPredictedRoutes(): string[] {
    // Return routes sorted by user behavior
    return Object.entries(this.userBehavior)
      .sort(([, a], [, b]) => b - a)
      .map(([route]) => route)
      .slice(0, 3); // Top 3 most visited routes
  }
}

export const smartPreloader = new SmartPreloader();