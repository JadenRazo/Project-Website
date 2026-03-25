import headshotImage from '../assets/images/headshot.webp';

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

    if (this.preloadedResources.size >= 100) {
      const oldest = this.preloadedResources.values().next().value;
      if (oldest) {
        this.preloadedResources.delete(oldest);
        this.preloadPromises.delete(oldest);
      }
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

      link.onload = () => {
        document.head.removeChild(link);
        resolve();
      };
      link.onerror = () => {
        document.head.removeChild(link);
        reject(new Error(`Failed to preload ${href}`));
      };

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
  headshot: headshotImage
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
    preloader.prefetchResource('/src/assets/data/code_stats.json');
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
  private cleanupIdleDetection: (() => void) | null = null;

  constructor() {
    this.trackUserBehavior();
    this.setupIdleDetection();
  }

  private trackUserBehavior(): void {
    const currentRoute = window.location.pathname;
    this.userBehavior[currentRoute] = (this.userBehavior[currentRoute] || 0) + 1;
  }

  private setupIdleDetection(): void {
    let idleTimer: ReturnType<typeof setTimeout>;

    const resetIdleTimer = () => {
      this.isIdle = false;
      clearTimeout(idleTimer);
      idleTimer = setTimeout(() => {
        this.isIdle = true;
        this.processPreloadQueue();
      }, 2000);
    };

    const events: [string, AddEventListenerOptions?][] = [
      ['mousedown'],
      ['scroll', { passive: true }],
      ['touchstart', { passive: true }],
    ];

    events.forEach(([event, opts]) => {
      document.addEventListener(event, resetIdleTimer, opts ?? true);
    });

    resetIdleTimer();

    this.cleanupIdleDetection = () => {
      clearTimeout(idleTimer);
      events.forEach(([event, opts]) => {
        document.removeEventListener(event, resetIdleTimer, opts ?? (true as any));
      });
    };
  }

  destroy(): void {
    this.cleanupIdleDetection?.();
    this.cleanupIdleDetection = null;
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