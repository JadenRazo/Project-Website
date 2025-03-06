// hooks/usePerformanceOptimizations.ts
import { useMemo } from 'react';
import useDeviceCapabilities, { DeviceCapabilities } from './useDeviceCapabilities';

/**
 * Interface for performance optimization settings
 */
export interface PerformanceSettings {
  // Visual effects settings
  enableParallax: boolean;        // Enable parallax effects
  enableComplexAnimations: boolean; // Enable visually complex animations
  enableBlurEffects: boolean;     // Enable backdrop blur and similar effects
  enableShimmerEffects: boolean;  // Enable shimmer/gradient animations
  
  // Animation timing settings
  transitionSpeed: number;        // Base transition speed in seconds
  staggerDelay: number;           // Delay between staggered animations
  
  // Interactive behavior settings
  enableHoverEffects: boolean;    // Enable hover animations/effects
  reduceMotion: boolean;          // Respect reduced motion preferences
  
  // Layout optimizations
  useSimplifiedLayout: boolean;   // Use simpler layouts on low-powered devices
  useLazyLoading: boolean;        // Lazy load content as needed
  
  // Image optimizations
  imageQualityLevel: 'low' | 'medium' | 'high';  // Quality level for images
  imagesToPreload: 'all' | 'essential' | 'none'; // Which images to preload
  useWebpImages: boolean;         // Use WebP format when supported
  
  // Animation optimizations
  batchAnimations: boolean;       // Group animations for better performance
  skipAnimations: boolean;        // Skip non-essential animations
  useHardwareAcceleration: boolean; // Use hardware acceleration for animations
  
  // Resource management
  aggressiveGarbageCollection: boolean; // More frequent cleanup of resources
  throttleNonVisibleAnimations: boolean; // Reduce animation rates when not visible
  
  // Rendering optimizations
  useDeferredRendering: boolean;  // Defer non-critical rendering
  useIntersectionObserver: boolean; // Use IntersectionObserver for optimized rendering
  
  // Performance tier classification
  performanceTier: 'low' | 'medium' | 'high'; // Overall device performance category
}

/**
 * Hook that provides performance optimization settings based on device capabilities
 */
export const usePerformanceOptimizations = (): {
  deviceCapabilities: DeviceCapabilities;
  performanceSettings: PerformanceSettings;
} => {
  // Get current device capabilities
  const deviceCapabilities = useDeviceCapabilities();
  
  // Calculate performance settings based on device capabilities
  const performanceSettings = useMemo<PerformanceSettings>(() => {
    const {
      isLowPoweredDevice,
      connectionType,
      effectiveConnectionType,
      connectionSavingEnabled,
      prefersReducedMotion,
      isTouchDevice,
      hasPointer,
      hasFinePointer,
      devicePixelRatio,
      viewportWidth,
      viewportHeight,
      deviceCategory,
      supportsWebP,
      supportsIntersectionObserver
    } = deviceCapabilities;
    
    // Determine performance tier based on multiple factors
    let performanceTier: 'low' | 'medium' | 'high' = 'high';
    
    // Check for low performance conditions
    if (
      isLowPoweredDevice ||
      effectiveConnectionType === '2g' ||
      connectionSavingEnabled ||
      (viewportWidth * viewportHeight < 480000) // Small screen (e.g., 800×600)
    ) {
      performanceTier = 'low';
    } 
    // Check for medium performance conditions
    else if (
      effectiveConnectionType === '3g' ||
      deviceCategory === 'mobile' ||
      (deviceCategory === 'tablet' && devicePixelRatio < 2)
    ) {
      performanceTier = 'medium';
    }
    
    // Calculate transition speed based on performance tier
    const baseTransitionSpeed = prefersReducedMotion ? 0.2 : 0.4;
    const transitionSpeed = 
      performanceTier === 'low' ? baseTransitionSpeed * 0.5 :
      performanceTier === 'medium' ? baseTransitionSpeed * 0.75 :
      baseTransitionSpeed;
    
    // Calculate stagger delay for sequential animations
    const baseStaggerDelay = prefersReducedMotion ? 0.04 : 0.08;
    const staggerDelay = 
      performanceTier === 'low' ? baseStaggerDelay * 0.5 :
      performanceTier === 'medium' ? baseStaggerDelay * 0.75 :
      baseStaggerDelay;
    
    // Determine appropriate image quality level
    const imageQualityLevel =
      performanceTier === 'low' ? 'low' :
      performanceTier === 'medium' ? 'medium' :
      'high';
    
    // Image preloading strategy based on connection and performance
    const imagesToPreload =
      performanceTier === 'low' || connectionSavingEnabled ? 'none' :
      performanceTier === 'medium' ? 'essential' :
      'all';
    
    return {
      // Visual effects settings
      enableParallax: performanceTier === 'high' && !prefersReducedMotion && !isTouchDevice,
      enableComplexAnimations: performanceTier !== 'low' && !prefersReducedMotion,
      enableBlurEffects: performanceTier === 'high' && devicePixelRatio < 3,
      enableShimmerEffects: performanceTier !== 'low' && !prefersReducedMotion,
      
      // Animation timing settings
      transitionSpeed,
      staggerDelay,
      
      // Interactive behavior settings
      enableHoverEffects: hasPointer && hasFinePointer,
      reduceMotion: prefersReducedMotion,
      
      // Layout optimizations
      useSimplifiedLayout: performanceTier === 'low' || deviceCategory === 'mobile',
      useLazyLoading: true, // Always use lazy loading as a best practice
      
      // Image optimizations
      imageQualityLevel,
      imagesToPreload,
      useWebpImages: supportsWebP,
      
      // Animation optimizations
      batchAnimations: performanceTier !== 'high',
      skipAnimations: performanceTier === 'low' || prefersReducedMotion,
      useHardwareAcceleration: performanceTier !== 'low',
      
      // Resource management
      aggressiveGarbageCollection: performanceTier === 'low',
      throttleNonVisibleAnimations: performanceTier !== 'high',
      
      // Rendering optimizations
      useDeferredRendering: performanceTier !== 'high',
      useIntersectionObserver: supportsIntersectionObserver,
      
      // Overall performance classification
      performanceTier
    };
  }, [deviceCapabilities]);
  
  return { deviceCapabilities, performanceSettings };
};

export default usePerformanceOptimizations;