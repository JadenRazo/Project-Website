import { useState, useEffect } from 'react';

/**
 * Interface describing various device capability properties
 */
export interface DeviceCapabilities {
  // Basic device detection
  isTouchDevice: boolean;
  hasPointer: boolean;
  hasFinePointer: boolean;
  hasCoarsePointer: boolean;
  
  // Motion preferences
  prefersReducedMotion: boolean;
  
  // Connection information
  connectionType: string;
  effectiveConnectionType: string; // 2g, 3g, 4g
  connectionSavingEnabled: boolean;
  
  // Device performance indicators
  isLowPoweredDevice: boolean;
  devicePixelRatio: number;
  
  // Screen properties
  orientation: 'portrait' | 'landscape';
  viewportWidth: number;
  viewportHeight: number;
  screenWidth: number;
  screenHeight: number;
  deviceCategory: 'mobile' | 'tablet' | 'desktop';
  
  // Additional accessibility preferences
  prefersColorScheme: 'dark' | 'light' | 'no-preference';
  prefersContrast: 'high' | 'low' | 'no-preference';
  
  // Browser features
  supportsWebP: boolean;
  supportsIntersectionObserver: boolean;
  supportsTouchEvents: boolean;
}

/**
 * Hook that detects and monitors various device capabilities
 * to enable optimized experiences across different devices
 */
export const useDeviceCapabilities = (): DeviceCapabilities => {
  // Default state with pessimistic assumptions
  const [capabilities, setCapabilities] = useState<DeviceCapabilities>({
    isTouchDevice: false,
    hasPointer: false,
    hasFinePointer: false,
    hasCoarsePointer: false,
    prefersReducedMotion: false,
    connectionType: 'unknown',
    effectiveConnectionType: 'unknown',
    connectionSavingEnabled: false,
    isLowPoweredDevice: false,
    devicePixelRatio: 1,
    orientation: 'portrait',
    viewportWidth: 0,
    viewportHeight: 0,
    screenWidth: 0,
    screenHeight: 0,
    deviceCategory: 'desktop',
    prefersColorScheme: 'no-preference',
    prefersContrast: 'no-preference',
    supportsWebP: false,
    supportsIntersectionObserver: false,
    supportsTouchEvents: false
  });
  
  // Flag to track if we've checked for WebP support
  const [webPChecked, setWebPChecked] = useState(false);
  
  useEffect(() => {
    // Skip detection during SSR
    if (typeof window === 'undefined') return;
    
    // Declare connection variable at the effect scope level so it's available for cleanup
    let connection: any = null;
    
    /**
     * Comprehensive feature detection function
     */
    const detectCapabilities = () => {
      // Touch detection
      const isTouchDevice = 'ontouchstart' in window || 
                           (window.navigator.maxTouchPoints || 0) > 0;
      
      // Pointer detection using media queries
      const hasPointer = window.matchMedia('(pointer: fine), (pointer: coarse)').matches;
      const hasFinePointer = window.matchMedia('(pointer: fine)').matches;
      const hasCoarsePointer = window.matchMedia('(pointer: coarse)').matches;
      
      // Reduced motion preference
      const prefersReducedMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches;
      
      // Color scheme preference
      const prefersColorScheme = window.matchMedia('(prefers-color-scheme: dark)').matches 
        ? 'dark' 
        : window.matchMedia('(prefers-color-scheme: light)').matches 
          ? 'light' 
          : 'no-preference';
      
      // Contrast preference
      const prefersContrast = window.matchMedia('(prefers-contrast: more)').matches 
        ? 'high' 
        : window.matchMedia('(prefers-contrast: less)').matches 
          ? 'low' 
          : 'no-preference';
      
      // Network connection information (when available)
      connection = 'connection' in navigator && 
                        (navigator as any).connection ? 
                        (navigator as any).connection : null;
                        
      const connectionType = connection ? connection.type || 'unknown' : 'unknown';
      const effectiveConnectionType = connection ? connection.effectiveType || 'unknown' : 'unknown';
      const connectionSavingEnabled = connection ? connection.saveData || false : false;
      
      // Device pixel ratio for high-DPI screens
      const devicePixelRatio = window.devicePixelRatio || 1;
      
      // Viewport and screen dimensions
      const viewportWidth = window.innerWidth;
      const viewportHeight = window.innerHeight;
      const screenWidth = window.screen.width;
      const screenHeight = window.screen.height;
      
      // Device category estimation
      let deviceCategory: 'mobile' | 'tablet' | 'desktop' = 'desktop';
      if (viewportWidth < 768 || (isTouchDevice && viewportWidth < 1024)) {
        deviceCategory = viewportWidth < 480 ? 'mobile' : 'tablet';
      }
      
      // Low power device estimation based on UA and connection
      const isLowPoweredDevice = 
        /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent) && 
        (effectiveConnectionType === '2g' || effectiveConnectionType === '3g' || 
         devicePixelRatio < 2 || viewportWidth < 768);
      
      // Orientation detection
      const orientation = viewportWidth > viewportHeight ? 'landscape' : 'portrait';
      
      // Feature detection for modern APIs
      const supportsIntersectionObserver = 'IntersectionObserver' in window;
      const supportsTouchEvents = 'ontouchstart' in window;
      
      // Apply fix for mobile viewport height (CSS custom property)
      document.documentElement.style.setProperty('--vh', `${viewportHeight * 0.01}px`);
      
      // Gather all detected capabilities
      setCapabilities({
        isTouchDevice,
        hasPointer,
        hasFinePointer,
        hasCoarsePointer,
        prefersReducedMotion,
        connectionType,
        effectiveConnectionType,
        connectionSavingEnabled,
        isLowPoweredDevice,
        devicePixelRatio,
        orientation,
        viewportWidth,
        viewportHeight,
        screenWidth,
        screenHeight,
        deviceCategory,
        prefersColorScheme: prefersColorScheme as any,
        prefersContrast: prefersContrast as any,
        supportsWebP: webPChecked ? capabilities.supportsWebP : false, // Will be updated by the WebP check
        supportsIntersectionObserver,
        supportsTouchEvents
      });
    };
    
    // Run detection immediately
    detectCapabilities();
    
    // Check for WebP support
    if (!webPChecked) {
      const webP = document.createElement('img');
      webP.onload = () => {
        setCapabilities(prev => ({ ...prev, supportsWebP: true }));
        setWebPChecked(true);
      };
      webP.onerror = () => {
        setCapabilities(prev => ({ ...prev, supportsWebP: false }));
        setWebPChecked(true);
      };
      webP.src = 'data:image/webp;base64,UklGRhoAAABXRUJQVlA4TA0AAAAvAAAAEAcQERGIiP4HAA==';
    }
    
    // Set up event listeners for window resize and orientation change
    const handleResize = () => {
      // Throttle the updates for better performance
      if (resizeTimeout) clearTimeout(resizeTimeout);
      resizeTimeout = setTimeout(detectCapabilities, 100);
    };
    
    let resizeTimeout: NodeJS.Timeout | null = null;
    window.addEventListener('resize', handleResize);
    
    // Also listen for orientation changes explicitly
    window.addEventListener('orientationchange', () => {
      // Small delay to ensure values are correct after orientation change
      setTimeout(detectCapabilities, 100);
    });
    
    // Listen for changes in user preferences
    const reducedMotionQuery = window.matchMedia('(prefers-reduced-motion: reduce)');
    const colorSchemeQuery = window.matchMedia('(prefers-color-scheme: dark)');
    const contrastQuery = window.matchMedia('(prefers-contrast: more)');
    
    // Update on preference changes
    const handlePreferenceChange = detectCapabilities;
    
    // Add event listeners with compatibility check for older browsers
    if (reducedMotionQuery.addEventListener) {
      reducedMotionQuery.addEventListener('change', handlePreferenceChange);
      colorSchemeQuery.addEventListener('change', handlePreferenceChange);
      contrastQuery.addEventListener('change', handlePreferenceChange);
    } else if ('addListener' in reducedMotionQuery) {
      // Older browsers support
      (reducedMotionQuery as any).addListener(handlePreferenceChange);
      (colorSchemeQuery as any).addListener(handlePreferenceChange);
      (contrastQuery as any).addListener(handlePreferenceChange);
    }
    
    // If available, listen for network connection changes
    if (connection && connection.addEventListener) {
      connection.addEventListener('change', detectCapabilities);
    }
    
    // Clean up event listeners
    return () => {
      if (resizeTimeout) clearTimeout(resizeTimeout);
      window.removeEventListener('resize', handleResize);
      window.removeEventListener('orientationchange', handleResize);
      
      if (reducedMotionQuery.removeEventListener) {
        reducedMotionQuery.removeEventListener('change', handlePreferenceChange);
        colorSchemeQuery.removeEventListener('change', handlePreferenceChange);
        contrastQuery.removeEventListener('change', handlePreferenceChange);
      } else if ('removeListener' in reducedMotionQuery) {
        // Older browsers support
        (reducedMotionQuery as any).removeListener(handlePreferenceChange);
        (colorSchemeQuery as any).removeListener(handlePreferenceChange);
        (contrastQuery as any).removeListener(handlePreferenceChange);
      }
      
      if (connection && connection.removeEventListener) {
        connection.removeEventListener('change', detectCapabilities);
      }
    };
  }, [capabilities.supportsWebP, webPChecked]);
  
  return capabilities;
};

export default useDeviceCapabilities;
