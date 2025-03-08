import { useState, useEffect, useRef, useCallback } from 'react';
import { useDeviceCapabilities, DeviceCapabilities } from './useDeviceCapabilities';
import { usePerformanceOptimizations, PerformanceSettings } from './usePerformanceOptimizations';
import { debounce } from 'lodash';

export interface AnimationConfig {
  // Network animation settings
  networkParticleCount: number;
  networkConnectionDistance: number;
  networkParticleSize: [number, number]; // [min, max]
  networkParticleSpeed: number;
  networkGlowIntensity: number;
  
  // Particle animation settings
  particleCount: number;
  particleSize: [number, number]; // [min, max]
  particleSpeed: [number, number]; // [min, max]
  particleOpacity: [number, number]; // [min, max]
  particleConnectionDistance: number;
  
  // Shared settings
  frameSkipRate: number;
  useHardwareAcceleration: boolean;
  useHighQualityRendering: boolean;
  debugMode: boolean;
}

export interface AnimationState {
  isActive: boolean;
  isPaused: boolean;
  isVisible: boolean;
  fps: number;
  particleCount: number;
  connectionCount: number;
}

export const useAnimationController = () => {
  // Get device capabilities and performance settings using existing hooks
  // Use try-catch to handle potential errors or missing hooks with defaults
  let deviceCapabilities: DeviceCapabilities;
  let performanceSettings: PerformanceSettings;
  
  try {
    const deviceData = usePerformanceOptimizations();
    deviceCapabilities = deviceData.deviceCapabilities;
    performanceSettings = deviceData.performanceSettings;
  } catch (e) {
    // Default values if hooks are not available
    deviceCapabilities = {
      isTouchDevice: 'ontouchstart' in window,
      hasPointer: true,
      hasFinePointer: true,
      hasCoarsePointer: false,
      prefersReducedMotion: false,
      connectionType: 'unknown',
      effectiveConnectionType: '4g',
      connectionSavingEnabled: false,
      isLowPoweredDevice: false,
      devicePixelRatio: window.devicePixelRatio || 1,
      orientation: window.innerWidth > window.innerHeight ? 'landscape' : 'portrait',
      viewportWidth: window.innerWidth,
      viewportHeight: window.innerHeight,
      screenWidth: window.screen.width,
      screenHeight: window.screen.height,
      deviceCategory: window.innerWidth < 768 ? 'mobile' : window.innerWidth < 1024 ? 'tablet' : 'desktop',
      prefersColorScheme: window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light',
      prefersContrast: 'no-preference',
      supportsWebP: true,
      supportsIntersectionObserver: 'IntersectionObserver' in window,
      supportsTouchEvents: 'ontouchstart' in window
    };
    
    performanceSettings = {
      enableParallax: true,
      enableComplexAnimations: true,
      enableBlurEffects: true,
      enableShimmerEffects: true,
      transitionSpeed: 0.3,
      staggerDelay: 0.1,
      enableHoverEffects: true,
      reduceMotion: false,
      useSimplifiedLayout: false,
      useLazyLoading: true,
      imageQualityLevel: 'high',
      imagesToPreload: 'essential',
      useWebpImages: true,
      batchAnimations: true,
      skipAnimations: false,
      useHardwareAcceleration: true,
      aggressiveGarbageCollection: false,
      throttleNonVisibleAnimations: true,
      useDeferredRendering: true,
      useIntersectionObserver: true,
      performanceTier: 'high'
    };
  }
  
  // Animation configuration based on device capabilities and performance settings
  const animationConfig = useRef<AnimationConfig>(generateAnimationConfig(deviceCapabilities, performanceSettings));
  
  // Animation state
  const [animationState, setAnimationState] = useState<AnimationState>({
    isActive: false,
    isPaused: false,
    isVisible: true,
    fps: 0,
    particleCount: 0,
    connectionCount: 0
  });
  
  // Update animation state
  const updateAnimationState = useCallback((partialState: Partial<AnimationState>) => {
    setAnimationState(prevState => ({
      ...prevState,
      ...partialState
    }));
  }, []);
  
  // Update animation config based on device changes or performance settings
  useEffect(() => {
    animationConfig.current = generateAnimationConfig(deviceCapabilities, performanceSettings);
    
    // Debug logging in development
    if (process.env.NODE_ENV === 'development') {
      console.log('Animation config updated:', animationConfig.current);
    }
  }, [deviceCapabilities, performanceSettings]);
  
  // Pause animations when page is not visible
  useEffect(() => {
    const handleVisibilityChange = () => {
      const isVisible = document.visibilityState === 'visible';
      updateAnimationState({ isPaused: !isVisible, isVisible });
    };
    
    document.addEventListener('visibilitychange', handleVisibilityChange);
    
    return () => {
      document.removeEventListener('visibilitychange', handleVisibilityChange);
    };
  }, [updateAnimationState]);
  
  // Handle visibility intersection for animation elements
  const createIntersectionObserver = useCallback((element: HTMLElement, onVisible: (visible: boolean) => void) => {
    // Check if IntersectionObserver is supported
    if (!window.IntersectionObserver || !performanceSettings.useIntersectionObserver) {
      onVisible(true); // Default to visible if not using IntersectionObserver
      return null;
    }
    
    const observer = new IntersectionObserver(
      entries => {
        const isVisible = entries[0]?.isIntersecting ?? false;
        onVisible(isVisible);
      },
      { threshold: 0.1 }
    );
    
    observer.observe(element);
    return observer;
  }, [performanceSettings.useIntersectionObserver]);
  
  // Monitor FPS with throttling to avoid excessive updates
  const monitorFps = useCallback(debounce((fps: number) => {
    updateAnimationState({ fps });
    
    // Auto-adjust settings if FPS is too low
    if (fps < 30 && animationConfig.current.frameSkipRate < 3) {
      animationConfig.current.frameSkipRate += 1;
      
      if (process.env.NODE_ENV === 'development') {
        console.log('FPS too low, increasing frame skip rate:', animationConfig.current.frameSkipRate);
      }
    }
  }, 1000), [updateAnimationState]);
  
  return {
    animationConfig: animationConfig.current,
    animationState,
    updateAnimationState,
    createIntersectionObserver,
    monitorFps,
    deviceCapabilities,
    performanceSettings
  };
};

// Generate optimal animation configuration based on device capabilities and performance settings
function generateAnimationConfig(
  deviceCapabilities: DeviceCapabilities,
  performanceSettings: PerformanceSettings
): AnimationConfig {
  // Determine baseline device tier
  const { performanceTier } = performanceSettings;
  const { isTouchDevice, devicePixelRatio, isLowPoweredDevice, connectionType } = deviceCapabilities;
  
  // Base configuration
  const baseConfig: AnimationConfig = {
    networkParticleCount: 60,
    networkConnectionDistance: 150,
    networkParticleSize: [1.5, 3],
    networkParticleSpeed: 0.3,
    networkGlowIntensity: 0.8,
    
    particleCount: 60,
    particleSize: [1.5, 4],
    particleSpeed: [0.2, 0.6],
    particleOpacity: [0.2, 0.7],
    particleConnectionDistance: 150,
    
    frameSkipRate: 1,
    useHardwareAcceleration: true,
    useHighQualityRendering: true,
    debugMode: false
  };
  
  // Apply performance tier adjustments
  if (performanceTier === 'low' || isLowPoweredDevice) {
    return {
      ...baseConfig,
      networkParticleCount: 30,
      networkConnectionDistance: 120,
      particleCount: 30,
      frameSkipRate: 3,
      useHighQualityRendering: false
    };
  }
  
  if (performanceTier === 'medium') {
    return {
      ...baseConfig,
      networkParticleCount: 50,
      particleCount: 50,
      frameSkipRate: 2
    };
  }
  
  // Touch device optimizations
  if (isTouchDevice) {
    baseConfig.networkParticleSpeed *= 0.7; // Slower particles look better on touch
    baseConfig.networkGlowIntensity *= 1.2; // More glow for touch devices
  }
  
  // Network connection optimizations
  if (connectionType === '2g' || connectionType === 'slow-2g') {
    baseConfig.networkParticleCount = 20;
    baseConfig.particleCount = 20;
    baseConfig.useHighQualityRendering = false;
    baseConfig.frameSkipRate = 3;
  }
  
  // Screen density optimizations
  if (devicePixelRatio > 2) {
    baseConfig.networkParticleSize = [2, 4]; // Larger particles for high-DPI screens
    baseConfig.particleSize = [2, 5];
  }
  
  // Reduced motion preference
  if (performanceSettings.reduceMotion) {
    baseConfig.networkParticleSpeed *= 0.5;
    baseConfig.particleSpeed = [baseConfig.particleSpeed[0] * 0.5, baseConfig.particleSpeed[1] * 0.5];
    baseConfig.frameSkipRate = Math.max(2, baseConfig.frameSkipRate);
  }
  
  // Debug mode in development
  if (process.env.NODE_ENV === 'development' && 
      typeof window !== 'undefined' && 
      window.location.search.includes('debug=true')) {
    baseConfig.debugMode = true;
  }
  
  return baseConfig;
} 