import { useEffect, useCallback } from 'react';
import { useLocation } from 'react-router-dom';
import { preloader, smartPreloader, preloadRouteAssets, preloadCriticalAssets } from '../utils/preloader';

interface UsePreloaderOptions {
  enableSmartPreloading?: boolean;
  enableRoutePreloading?: boolean;
  enableHoverPreloading?: boolean;
  preloadDelay?: number;
}

export const usePreloader = (options: UsePreloaderOptions = {}) => {
  const {
    enableSmartPreloading = true,
    enableRoutePreloading = true,
    enableHoverPreloading = true,
    preloadDelay = 200
  } = options;

  const location = useLocation();

  // Preload critical assets on app initialization
  useEffect(() => {
    preloadCriticalAssets();
  }, []);

  // Route change preloading
  useEffect(() => {
    if (!enableRoutePreloading) return;

    const currentRoute = location.pathname;
    const preloadFn = preloadRouteAssets[currentRoute as keyof typeof preloadRouteAssets];
    
    if (preloadFn) {
      // Small delay to ensure the current page is fully loaded first
      setTimeout(preloadFn, 100);
    }

    // Queue smart preloading for predicted routes
    if (enableSmartPreloading) {
      const predictedRoutes = smartPreloader.getPredictedRoutes();
      predictedRoutes.forEach(route => {
        if (route !== currentRoute) {
          smartPreloader.queuePreload(route);
        }
      });
    }
  }, [location.pathname, enableRoutePreloading, enableSmartPreloading]);

  // Hover-based link preloading
  useEffect(() => {
    if (!enableHoverPreloading) return;

    let hoverTimeout: NodeJS.Timeout;

    const handleLinkHover = (event: MouseEvent) => {
      const target = event.target as HTMLElement;
      const link = target.closest('a[href]') as HTMLAnchorElement;
      
      if (!link || !link.href) return;
      
      const url = new URL(link.href);
      
      // Only preload internal links
      if (url.origin !== window.location.origin) return;
      
      hoverTimeout = setTimeout(() => {
        const route = url.pathname;
        const preloadFn = preloadRouteAssets[route as keyof typeof preloadRouteAssets];
        
        if (preloadFn) {
          try {
            preloadFn();
          } catch (error) {
            console.warn(`Failed to preload on hover for route ${route}:`, error);
          }
        }
      }, preloadDelay);
    };

    const handleLinkLeave = () => {
      if (hoverTimeout) {
        clearTimeout(hoverTimeout);
      }
    };

    document.addEventListener('mouseover', handleLinkHover);
    document.addEventListener('mouseout', handleLinkLeave);

    return () => {
      document.removeEventListener('mouseover', handleLinkHover);
      document.removeEventListener('mouseout', handleLinkLeave);
      if (hoverTimeout) {
        clearTimeout(hoverTimeout);
      }
    };
  }, [enableHoverPreloading, preloadDelay]);

  // Manual preloading functions
  const preloadRoute = useCallback((route: string) => {
    const preloadFn = preloadRouteAssets[route as keyof typeof preloadRouteAssets];
    if (preloadFn) {
      try {
        preloadFn();
      } catch (error) {
        console.warn(`Failed to manually preload route ${route}:`, error);
      }
    }
  }, []);

  const preloadImage = useCallback((src: string) => {
    return preloader.preload({ href: src, as: 'image' });
  }, []);

  const preloadImages = useCallback((srcs: string[]) => {
    return preloader.preloadImages(srcs);
  }, []);

  const prefetchResource = useCallback((href: string) => {
    preloader.prefetchResource(href);
  }, []);

  return {
    preloadRoute,
    preloadImage,
    preloadImages,
    prefetchResource,
    smartPreloader
  };
};

export default usePreloader;