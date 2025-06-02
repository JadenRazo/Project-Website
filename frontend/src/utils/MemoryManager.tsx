import React, { createContext, useContext, useEffect, useRef, useState, useCallback } from 'react';
import { monitorMemoryUsage, captureMemoryMetrics, PerformanceMetrics } from './performance';

// Check if performance API is available
const isPerformanceApiAvailable = typeof performance !== 'undefined' && 
  typeof performance.now === 'function';

// Check if memory API is available
const isMemoryApiAvailable = isPerformanceApiAvailable && 
  typeof (performance as any).memory !== 'undefined';

// Check if caches API is available
const isCachesApiAvailable = typeof caches !== 'undefined';

interface MemoryManagerContextType {
  memoryUsage: PerformanceMetrics;
  freeMemory: () => void;
  optimizePerformance: (level?: 'low' | 'medium' | 'high') => void;
  isMemoryConstrained: boolean;
  applicationState: {
    effectsEnabled: boolean;
    animationsEnabled: boolean;
    backgroundEffectsEnabled: boolean;
    highQualityImagesEnabled: boolean;
    virtualizationEnabled: boolean;
  };
}

// Default context value with feature detection
const defaultContextValue: MemoryManagerContextType = {
  memoryUsage: { timestamp: Date.now() },
  freeMemory: () => {},
  optimizePerformance: () => {},
  isMemoryConstrained: false,
  applicationState: {
    effectsEnabled: true,
    animationsEnabled: true,
    backgroundEffectsEnabled: true,
    highQualityImagesEnabled: true,
    virtualizationEnabled: true
  }
};

// Create context
const MemoryManagerContext = createContext<MemoryManagerContextType>(defaultContextValue);

interface MemoryManagerProviderProps {
  children: React.ReactNode;
  monitoringInterval?: number;
  memoryThreshold?: number;
  enableLogging?: boolean;
}

/**
 * Provider component that manages application memory usage
 * Implements intelligent resource management and performance optimization
 */
export const MemoryManagerProvider: React.FC<MemoryManagerProviderProps> = ({
  children,
  monitoringInterval = 30000,
  memoryThreshold = 80,
  enableLogging = process.env.NODE_ENV === 'development'
}) => {
  // Skip monitoring if memory API is not available
  const shouldMonitorMemory = isMemoryApiAvailable;
  
  const [memoryUsage, setMemoryUsage] = useState<PerformanceMetrics>({ timestamp: Date.now() });
  const [isMemoryConstrained, setIsMemoryConstrained] = useState(false);
  const metricsHistory = useRef<PerformanceMetrics[]>([]);
  
  // Application performance states
  const [applicationState, setApplicationState] = useState({
    effectsEnabled: true,
    animationsEnabled: true,
    backgroundEffectsEnabled: true,
    highQualityImagesEnabled: true,
    virtualizationEnabled: true
  });

  // Free up memory (with API checks)
  const freeMemory = useCallback(() => {
    // Clear image cache (works in all browsers)
    const images = document.querySelectorAll('img[data-src]');
    images.forEach(img => {
      if (!img.classList.contains('critical-image')) {
        (img as HTMLImageElement).src = 'data:image/gif;base64,R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7';
      }
    });
    
    // Clear any application caches (only if available)
    if (isCachesApiAvailable) {
      caches.keys().then(cacheNames => {
        cacheNames.forEach(cacheName => {
          if (cacheName.includes('non-critical')) {
            caches.delete(cacheName);
          }
        });
      });
    }
    
    // Try to force a garbage collection if browser allows it
    if (typeof (window as any).gc === 'function') {
      try {
        (window as any).gc();
      } catch (e) {
        console.debug('Manual GC not available');
      }
    }
    
    if (enableLogging) {
      console.info('%cMemory cleanup performed', 'color: #4ecdc4;');
    }
    
    // Update metrics after cleanup (if API available)
    if (shouldMonitorMemory) {
      const updatedMetrics = captureMemoryMetrics();
      setMemoryUsage(updatedMetrics);
      return updatedMetrics;
    }
    
    return { timestamp: Date.now() };
  }, [enableLogging, shouldMonitorMemory]);

  // Update memory metrics at regular intervals (only if API is available)
  // Apply performance optimizations based on memory constraints
  const optimizePerformance = useCallback((level: 'low' | 'medium' | 'high' = 'medium') => {
    let returnState: any = null;
    
    setApplicationState(prevState => {
      const newState = { ...prevState };
      
      switch (level) {
        case 'low':
          // Minimal optimizations
          newState.backgroundEffectsEnabled = false;
          break;
          
        case 'medium':
          // Moderate optimizations
          newState.backgroundEffectsEnabled = false;
          newState.highQualityImagesEnabled = false;
          newState.virtualizationEnabled = true;
          break;
          
        case 'high':
          // Aggressive optimizations
          newState.backgroundEffectsEnabled = false;
          newState.highQualityImagesEnabled = false;
          newState.animationsEnabled = false;
          newState.effectsEnabled = false;
          newState.virtualizationEnabled = true;
          break;
      }
      
      returnState = newState;
      return newState;
    });
    
    // Force cleanup after optimization (outside setState to avoid loops)
    setTimeout(freeMemory, 500);
    
    return returnState;
  }, [freeMemory]);

  useEffect(() => {
    if (!shouldMonitorMemory) {
      if (enableLogging) {
        console.info('Memory API not available, memory monitoring disabled');
      }
      return;
    }
    
    const updateMemoryMetrics = () => {
      const metrics = captureMemoryMetrics();
      setMemoryUsage(metrics);
      
      // Keep history for trend analysis
      metricsHistory.current.push(metrics);
      if (metricsHistory.current.length > 10) {
        metricsHistory.current = metricsHistory.current.slice(-10);
      }
      
      // Check memory constraints
      if (metrics.usedJSHeapSize && metrics.jsHeapSizeLimit) {
        const usagePercent = (metrics.usedJSHeapSize / metrics.jsHeapSizeLimit) * 100;
        const newMemoryConstrained = usagePercent > memoryThreshold;
        
        if (newMemoryConstrained !== isMemoryConstrained) {
          setIsMemoryConstrained(newMemoryConstrained);
          
          // Auto-optimize if memory becomes constrained
          if (newMemoryConstrained) {
            setTimeout(() => optimizePerformance('medium'), 0);
            
            if (enableLogging) {
              console.warn(
                `%cMemory usage exceeded threshold (${usagePercent.toFixed(1)}%)`,
                'color: #ff6b6b; font-weight: bold;',
                '\nAuto-optimizing performance...'
              );
            }
          }
        }
      }
    };
    
    // Initial update
    updateMemoryMetrics();
    
    // Regular updates
    const intervalId = setInterval(updateMemoryMetrics, monitoringInterval);
    
    // Setup peak memory usage warnings
    if (enableLogging) {
      monitorMemoryUsage('Application', {
        logToConsole: true,
        thresholdPercent: memoryThreshold,
        logLevel: 'warn'
      });
    }
    
    return () => clearInterval(intervalId);
  }, [monitoringInterval, memoryThreshold, enableLogging, shouldMonitorMemory, isMemoryConstrained, optimizePerformance]);
  
  const contextValue: MemoryManagerContextType = {
    memoryUsage,
    freeMemory,
    optimizePerformance,
    isMemoryConstrained,
    applicationState
  };
  
  return (
    <MemoryManagerContext.Provider value={contextValue}>
      {children}
    </MemoryManagerContext.Provider>
  );
};

/**
 * Hook to access memory management features
 * Use in components to implement adaptive performance based on memory constraints
 */
export const useMemoryManager = (): MemoryManagerContextType => {
  const context = useContext(MemoryManagerContext);
  
  if (!context) {
    throw new Error('useMemoryManager must be used within a MemoryManagerProvider');
  }
  
  return context;
};

/**
 * HOC to add memory management capabilities to any component
 * Automatically optimizes rendering based on memory constraints
 */
export function withMemoryOptimization<P extends object>(
  Component: React.ComponentType<P>,
  options: {
    optimizationLevel?: 'low' | 'medium' | 'high';
    disableWhenConstrained?: boolean;
  } = {}
): React.FC<P> {
  const { 
    disableWhenConstrained = false
  } = options;
  
  const WithMemoryOptimization: React.FC<P> = (props) => {
    const { isMemoryConstrained, applicationState } = useMemoryManager();
    
    // For expensive components, optionally don't render them at all when memory is constrained
    if (disableWhenConstrained && isMemoryConstrained) {
      return null;
    }
    
    // Pass memory optimization props to the wrapped component
    return (
      <Component
        {...props}
        memoryOptimized={isMemoryConstrained}
        effectsEnabled={applicationState.effectsEnabled}
        animationsEnabled={applicationState.animationsEnabled}
        highQualityEnabled={applicationState.highQualityImagesEnabled}
      />
    );
  };
  
  WithMemoryOptimization.displayName = `WithMemoryOptimization(${
    Component.displayName || Component.name || 'Component'
  })`;
  
  return WithMemoryOptimization;
};

export default MemoryManagerProvider; 