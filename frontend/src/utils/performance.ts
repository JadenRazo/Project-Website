import { useRef, useEffect, useCallback } from 'react';

/**
 * Performance monitoring and optimization utilities
 * Provides tools for tracking memory usage and optimizing resource cleanup
 */

export interface PerformanceMetrics {
  jsHeapSizeLimit?: number;
  totalJSHeapSize?: number;
  usedJSHeapSize?: number;
  timestamp: number;
}

interface MemoryMonitorOptions {
  logToConsole?: boolean;
  logLevel?: 'debug' | 'warn' | 'error';
  thresholdPercent?: number;
  onThresholdExceeded?: (metrics: PerformanceMetrics) => void;
}

const DEFAULT_OPTIONS: MemoryMonitorOptions = {
  logToConsole: process.env.NODE_ENV === 'development',
  logLevel: 'warn',
  thresholdPercent: 80
};

/**
 * Captures current memory usage metrics if available
 * Falls back gracefully when memory API isn't available
 */
export const captureMemoryMetrics = (): PerformanceMetrics => {
  const metrics: PerformanceMetrics = {
    timestamp: Date.now()
  };

  if (performance && 'memory' in performance) {
    const memory = (performance as any).memory;
    metrics.jsHeapSizeLimit = memory.jsHeapSizeLimit;
    metrics.totalJSHeapSize = memory.totalJSHeapSize;
    metrics.usedJSHeapSize = memory.usedJSHeapSize;
  }

  return metrics;
};

/**
 * Formats bytes into a human-readable string
 */
export const formatBytes = (bytes: number | undefined): string => {
  if (bytes === undefined) return 'N/A';
  
  if (bytes === 0) return '0 Bytes';
  
  const k = 1024;
  const sizes = ['Bytes', 'KB', 'MB', 'GB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
};

/**
 * Monitors memory usage in the application
 * Provides threshold warnings and optimization suggestions
 */
export const monitorMemoryUsage = (
  componentName: string,
  options: MemoryMonitorOptions = {}
): PerformanceMetrics => {
  const mergedOptions = { ...DEFAULT_OPTIONS, ...options };
  const metrics = captureMemoryMetrics();
  
  if (mergedOptions.logToConsole && metrics.usedJSHeapSize && metrics.jsHeapSizeLimit) {
    const usagePercent = (metrics.usedJSHeapSize / metrics.jsHeapSizeLimit) * 100;
    const usageExceedsThreshold = usagePercent > (mergedOptions.thresholdPercent || 80);
    
    if (usageExceedsThreshold) {
      const logMethod = mergedOptions.logLevel === 'error' 
        ? console.error 
        : mergedOptions.logLevel === 'warn' 
          ? console.warn 
          : console.debug;
          
      logMethod(
        `%cMemory usage in ${componentName}: ${usagePercent.toFixed(1)}%`,
        'color: #ff6b6b; font-weight: bold;',
        `\nUsed: ${formatBytes(metrics.usedJSHeapSize)} / ${formatBytes(metrics.jsHeapSizeLimit)}`
      );
      
      if (mergedOptions.onThresholdExceeded) {
        mergedOptions.onThresholdExceeded(metrics);
      }
    }
  }
  
  return metrics;
};

/**
 * Hook for tracking memory leaks in components
 * Helps identify components that continue to consume memory after unmounting
 */
export const useMemoryTracker = (componentName: string, options: MemoryMonitorOptions = {}) => {
  const metricsRef = useRef<PerformanceMetrics[]>([]);
  
  useEffect(() => {
    // Capture initial memory state
    metricsRef.current.push(captureMemoryMetrics());
    
    // Setup interval for tracking
    const interval = setInterval(() => {
      metricsRef.current.push(captureMemoryMetrics());
      
      // Keep only last 10 measurements to avoid memory leak in the tracker itself
      if (metricsRef.current.length > 10) {
        metricsRef.current = metricsRef.current.slice(-10);
      }
      
      monitorMemoryUsage(componentName, options);
    }, 5000);
    
    return () => {
      clearInterval(interval);
      
      // Check memory difference at unmount
      setTimeout(() => {
        const finalMetrics = captureMemoryMetrics();
        const initialMetrics = metricsRef.current[0];
        
        if (finalMetrics.usedJSHeapSize && 
            initialMetrics.usedJSHeapSize && 
            finalMetrics.usedJSHeapSize > initialMetrics.usedJSHeapSize + 5 * 1024 * 1024) {
          console.warn(
            `%cPossible memory leak detected in ${componentName}`,
            'color: #ff6b6b; font-weight: bold;',
            `\nMemory increased by ${formatBytes(finalMetrics.usedJSHeapSize - initialMetrics.usedJSHeapSize)}`
          );
        }
      }, 1000);
    };
  }, [componentName]);
  
  return metricsRef.current;
};

/**
 * Creates a resource cleanup function with enhanced safety checks
 * Ensures resources are properly cleaned up even in exceptional conditions
 */
export const createCleanupManager = () => {
  const cleanupFunctions: Array<() => void> = [];
  
  const addCleanupTask = (cleanupFn: () => void) => {
    cleanupFunctions.push(cleanupFn);
    
    // Return function to remove this specific cleanup task
    return () => {
      const index = cleanupFunctions.indexOf(cleanupFn);
      if (index !== -1) {
        cleanupFunctions.splice(index, 1);
      }
    };
  };
  
  const performCleanup = () => {
    // Execute all cleanup functions in reverse order (LIFO)
    // This ensures dependent resources are cleaned up properly
    for (let i = cleanupFunctions.length - 1; i >= 0; i--) {
      try {
        cleanupFunctions[i]();
      } catch (error) {
        console.error('Error during resource cleanup:', error);
      }
    }
    
    // Clear the array
    cleanupFunctions.length = 0;
  };
  
  return {
    addCleanupTask,
    performCleanup
  };
};

/**
 * Hook that provides safe async operation management
 * Prevents state updates on unmounted components
 */
export const useSafeAsync = () => {
  const mountedRef = useRef(true);
  
  useEffect(() => {
    return () => {
      mountedRef.current = false;
    };
  }, []);
  
  const safeExecute = useCallback(<T extends any[]>(
    callback: (...args: T) => void
  ) => {
    return (...args: T) => {
      if (mountedRef.current) {
        callback(...args);
      }
    };
  }, []);
  
  return { safeExecute, isMounted: () => mountedRef.current };
};

/**
 * Hook to handle debounced operations safely
 * Helps reduce memory usage by limiting function call frequency
 */
export const useDebounce = <T extends (...args: any[]) => void>(
  callback: T,
  delay: number
): [T, () => void] => {
  const timerRef = useRef<number | null>(null);
  const callbackRef = useRef(callback);
  
  // Update the callback ref when the callback changes
  useEffect(() => {
    callbackRef.current = callback;
  }, [callback]);
  
  // Cleanup on unmount
  useEffect(() => {
    return () => {
      if (timerRef.current !== null) {
        clearTimeout(timerRef.current);
      }
    };
  }, []);
  
  const debouncedFn = useCallback((...args: Parameters<T>) => {
    if (timerRef.current !== null) {
      clearTimeout(timerRef.current);
    }
    
    timerRef.current = window.setTimeout(() => {
      callbackRef.current(...args);
      timerRef.current = null;
    }, delay) as unknown as number;
  }, [delay]) as T;
  
  const cancelDebounce = useCallback(() => {
    if (timerRef.current !== null) {
      clearTimeout(timerRef.current);
      timerRef.current = null;
    }
  }, []);
  
  return [debouncedFn, cancelDebounce];
};

/**
 * Utility function to safely clean up WebGL contexts
 * Prevents memory leaks associated with GPU resources
 */
export const cleanupWebGLContext = (
  canvas: HTMLCanvasElement | null | undefined, 
  contextType: 'webgl' | 'webgl2' = 'webgl'
) => {
  if (!canvas) return;
  
  try {
    const gl = canvas.getContext(contextType) as WebGLRenderingContext | WebGL2RenderingContext | null;
    if (gl) {
      const extension = gl.getExtension('WEBGL_lose_context');
      if (extension) {
        extension.loseContext();
      }
      
      // Force context loss
      const loseContextExt = 
        gl.getExtension('WEBGL_lose_context') || 
        gl.getExtension('webkit-WEBGL_lose_context');
        
      if (loseContextExt) {
        loseContextExt.loseContext();
      }
    }
  } catch (error) {
    console.error('Error cleaning up WebGL context:', error);
  }
};

/**
 * Creates an optimized version of a function that only runs when in viewport
 * Reduces unnecessary processing and memory usage for offscreen elements
 */
export const createViewportAwareFunction = <T extends (...args: any[]) => void>(
  element: React.RefObject<HTMLElement>,
  fn: T,
  options: IntersectionObserverInit = {}
): T => {
  let intersectionObserver: IntersectionObserver | null = null;
  let isInViewport = false;
  
  const wrappedFn = ((...args: Parameters<T>) => {
    if (isInViewport) {
      return fn(...args);
    }
    
    return undefined;
  }) as T;
  
  if (element.current) {
    intersectionObserver = new IntersectionObserver(
      (entries) => {
        isInViewport = entries[0]?.isIntersecting ?? false;
      },
      { threshold: 0.1, ...options }
    );
    
    intersectionObserver.observe(element.current);
  }
  
  // Add cleanup method to the wrapped function
  (wrappedFn as any).cleanup = () => {
    if (intersectionObserver) {
      intersectionObserver.disconnect();
      intersectionObserver = null;
    }
  };
  
  return wrappedFn;
}; 