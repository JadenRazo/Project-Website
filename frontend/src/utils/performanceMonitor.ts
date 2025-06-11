// Performance Monitoring Utility
import React from 'react';

export interface PerformanceMetrics {
  timestamp: number;
  // Memory metrics
  memoryUsage: {
    used: number;
    total: number;
    limit: number;
    percentage: number;
  };
  // Rendering metrics
  fps: number;
  frameDrops: number;
  // Timing metrics
  domContentLoaded: number;
  firstContentfulPaint: number;
  largestContentfulPaint: number;
  firstInputDelay: number;
  cumulativeLayoutShift: number;
  // Custom metrics
  componentRenderTime: Map<string, number>;
  apiResponseTimes: Map<string, number>;
  bundleSize: number;
  // Device metrics
  deviceInfo: {
    isLowMemory: boolean;
    isMobile: boolean;
    pixelRatio: number;
    hardwareConcurrency: number;
  };
}

class PerformanceMonitor {
  private static instance: PerformanceMonitor;
  private isMonitoring: boolean = false;
  private frameCount: number = 0;
  private lastFrameTime: number = 0;
  private frameDrops: number = 0;
  private renderTimes: Map<string, number> = new Map();
  private apiTimes: Map<string, number> = new Map();
  private observer: PerformanceObserver | null = null;
  private intervalId: number | null = null;

  static getInstance(): PerformanceMonitor {
    if (!PerformanceMonitor.instance) {
      PerformanceMonitor.instance = new PerformanceMonitor();
    }
    return PerformanceMonitor.instance;
  }

  startMonitoring(): void {
    if (this.isMonitoring || typeof window === 'undefined') return;

    this.isMonitoring = true;
    this.setupPerformanceObserver();
    this.startFrameMonitoring();
    this.startMemoryMonitoring();
  }

  stopMonitoring(): void {
    if (!this.isMonitoring) return;

    this.isMonitoring = false;
    
    if (this.observer) {
      this.observer.disconnect();
      this.observer = null;
    }

    if (this.intervalId) {
      clearInterval(this.intervalId);
      this.intervalId = null;
    }
  }

  private setupPerformanceObserver(): void {
    if (!window.PerformanceObserver) return;

    this.observer = new PerformanceObserver((list) => {
      const entries = list.getEntries();
      
      entries.forEach((entry) => {
        if (entry.entryType === 'navigation') {
          console.log(`Navigation timing: ${entry.name} - ${entry.duration}ms`);
        } else if (entry.entryType === 'paint') {
          console.log(`Paint timing: ${entry.name} - ${entry.startTime}ms`);
        } else if (entry.entryType === 'largest-contentful-paint') {
          console.log(`LCP: ${entry.startTime}ms`);
        } else if (entry.entryType === 'first-input') {
          console.log(`FID: ${(entry as any).processingStart - entry.startTime}ms`);
        } else if (entry.entryType === 'layout-shift') {
          console.log(`CLS: ${(entry as any).value}`);
        }
      });
    });

    try {
      this.observer.observe({ entryTypes: ['navigation', 'paint', 'largest-contentful-paint', 'first-input', 'layout-shift'] });
    } catch (e) {
      console.warn('Some performance entry types not supported', e);
    }
  }

  private startFrameMonitoring(): void {
    const monitorFrame = (timestamp: number) => {
      if (!this.isMonitoring) return;

      this.frameCount++;
      
      if (this.lastFrameTime) {
        const deltaTime = timestamp - this.lastFrameTime;
        if (deltaTime > 16.67) { // Frame took longer than 60fps (16.67ms)
          this.frameDrops++;
        }
      }

      this.lastFrameTime = timestamp;
      requestAnimationFrame(monitorFrame);
    };

    requestAnimationFrame(monitorFrame);
  }

  private startMemoryMonitoring(): void {
    this.intervalId = window.setInterval(() => {
      this.logCurrentMetrics();
    }, 5000); // Log every 5 seconds
  }

  private getMemoryMetrics() {
    if ('memory' in performance) {
      const memory = (performance as any).memory;
      return {
        used: memory.usedJSHeapSize,
        total: memory.totalJSHeapSize,
        limit: memory.jsHeapSizeLimit,
        percentage: (memory.usedJSHeapSize / memory.jsHeapSizeLimit) * 100
      };
    }
    return { used: 0, total: 0, limit: 0, percentage: 0 };
  }

  private getFPS(): number {
    const fps = this.frameCount;
    this.frameCount = 0; // Reset for next interval
    return fps / 5; // 5-second intervals
  }

  private getDeviceInfo() {
    return {
      isLowMemory: navigator.deviceMemory ? navigator.deviceMemory <= 4 : false,
      isMobile: /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent),
      pixelRatio: window.devicePixelRatio || 1,
      hardwareConcurrency: navigator.hardwareConcurrency || 1
    };
  }

  recordComponentRender(componentName: string, duration: number): void {
    this.renderTimes.set(componentName, duration);
  }

  recordAPICall(endpoint: string, duration: number): void {
    this.apiTimes.set(endpoint, duration);
  }

  getCurrentMetrics(): PerformanceMetrics {
    const memoryMetrics = this.getMemoryMetrics();
    const deviceInfo = this.getDeviceInfo();

    return {
      timestamp: Date.now(),
      memoryUsage: memoryMetrics,
      fps: this.getFPS(),
      frameDrops: this.frameDrops,
      domContentLoaded: this.getPerformanceTiming('domContentLoadedEventEnd'),
      firstContentfulPaint: this.getPerformanceTiming('first-contentful-paint'),
      largestContentfulPaint: this.getPerformanceTiming('largest-contentful-paint'),
      firstInputDelay: this.getPerformanceTiming('first-input-delay'),
      cumulativeLayoutShift: this.getPerformanceTiming('cumulative-layout-shift'),
      componentRenderTime: new Map(this.renderTimes),
      apiResponseTimes: new Map(this.apiTimes),
      bundleSize: this.estimateBundleSize(),
      deviceInfo
    };
  }

  private getPerformanceTiming(metricName: string): number {
    try {
      const entries = performance.getEntriesByName(metricName);
      return entries.length > 0 ? entries[0].startTime : 0;
    } catch (e) {
      return 0;
    }
  }

  private estimateBundleSize(): number {
    // Rough estimation based on loaded scripts
    let totalSize = 0;
    const scripts = document.querySelectorAll('script[src]');
    
    scripts.forEach((script) => {
      const src = (script as HTMLScriptElement).src;
      if (src.includes('static/js/')) {
        // Estimate based on typical React bundle sizes
        totalSize += 250 * 1024; // 250KB average per chunk
      }
    });

    return totalSize;
  }

  private logCurrentMetrics(): void {
    if (!this.isMonitoring) return;

    const metrics = this.getCurrentMetrics();
    
    console.group('🔍 Performance Metrics');
    console.log('Memory Usage:', `${(metrics.memoryUsage.used / 1024 / 1024).toFixed(2)}MB (${metrics.memoryUsage.percentage.toFixed(1)}%)`);
    console.log('FPS:', metrics.fps);
    console.log('Frame Drops:', this.frameDrops);
    console.log('Device Info:', metrics.deviceInfo);
    
    if (metrics.componentRenderTime.size > 0) {
      console.log('Component Render Times:');
      metrics.componentRenderTime.forEach((time, component) => {
        console.log(`  ${component}: ${time.toFixed(2)}ms`);
      });
    }

    if (metrics.apiResponseTimes.size > 0) {
      console.log('API Response Times:');
      metrics.apiResponseTimes.forEach((time, endpoint) => {
        console.log(`  ${endpoint}: ${time.toFixed(2)}ms`);
      });
    }

    // Performance warnings
    if (metrics.memoryUsage.percentage > 80) {
      console.warn('⚠️ High memory usage detected!');
    }

    if (metrics.fps < 30) {
      console.warn('⚠️ Low FPS detected!');
    }

    if (this.frameDrops > 10) {
      console.warn('⚠️ High frame drop count!');
      this.frameDrops = 0; // Reset after warning
    }

    console.groupEnd();
  }

  // React component wrapper for measuring render performance
  measureComponentRender<T extends React.ElementType>(
    Component: T,
    componentName: string
  ) {
    type RefElement = React.ElementRef<T>;
    type Props = React.ComponentProps<T>;
    
    const WrappedComponent = React.forwardRef<RefElement, Props>((props, ref) => {
      const startTime = performance.now();
      
      React.useEffect(() => {
        const endTime = performance.now();
        this.recordComponentRender(componentName, endTime - startTime);
      });

      return React.createElement(Component as any, { ...props, ref });
    });

    WrappedComponent.displayName = `Measured(${componentName || 'Component'})`;
    return WrappedComponent;
  }

  // API call wrapper for measuring response times
  measureAPICall<T>(apiCall: () => Promise<T>, endpoint: string): Promise<T> {
    const startTime = performance.now();
    
    return apiCall().finally(() => {
      const endTime = performance.now();
      this.recordAPICall(endpoint, endTime - startTime);
    });
  }

  // Export metrics for analysis
  exportMetrics(): string {
    const metrics = this.getCurrentMetrics();
    return JSON.stringify(metrics, (key, value) => {
      if (value instanceof Map) {
        return Object.fromEntries(value);
      }
      return value;
    }, 2);
  }

  // Generate performance report
  generateReport(): {
    overall: 'excellent' | 'good' | 'fair' | 'poor';
    recommendations: string[];
    metrics: PerformanceMetrics;
  } {
    const metrics = this.getCurrentMetrics();
    const recommendations: string[] = [];
    let score = 100;

    // Memory score
    if (metrics.memoryUsage.percentage > 90) {
      score -= 30;
      recommendations.push('Critical memory usage - implement memory optimization');
    } else if (metrics.memoryUsage.percentage > 70) {
      score -= 15;
      recommendations.push('High memory usage - consider lazy loading and component optimization');
    }

    // FPS score
    if (metrics.fps < 30) {
      score -= 25;
      recommendations.push('Low FPS - optimize animations and reduce render complexity');
    } else if (metrics.fps < 45) {
      score -= 10;
      recommendations.push('Moderate FPS issues - consider reducing animation intensity');
    }

    // Frame drops
    if (this.frameDrops > 20) {
      score -= 20;
      recommendations.push('High frame drop count - optimize render loop and reduce computations');
    }

    // LCP score
    if (metrics.largestContentfulPaint > 2500) {
      score -= 15;
      recommendations.push('Slow LCP - optimize image loading and critical resource delivery');
    }

    // Bundle size
    if (metrics.bundleSize > 1024 * 1024) { // 1MB
      score -= 10;
      recommendations.push('Large bundle size - implement code splitting and tree shaking');
    }

    let overall: 'excellent' | 'good' | 'fair' | 'poor';
    if (score >= 85) overall = 'excellent';
    else if (score >= 70) overall = 'good';
    else if (score >= 50) overall = 'fair';
    else overall = 'poor';

    return { overall, recommendations, metrics };
  }
}

// Singleton instance
export const performanceMonitor = PerformanceMonitor.getInstance();

// React hook for using performance monitoring
export function usePerformanceMonitoring(enabled: boolean = process.env.NODE_ENV === 'development') {
  React.useEffect(() => {
    if (enabled) {
      performanceMonitor.startMonitoring();
      return () => performanceMonitor.stopMonitoring();
    }
  }, [enabled]);

  return {
    monitor: performanceMonitor,
    getCurrentMetrics: () => performanceMonitor.getCurrentMetrics(),
    generateReport: () => performanceMonitor.generateReport(),
    exportMetrics: () => performanceMonitor.exportMetrics()
  };
}

// Decorator for measuring component performance
export function withPerformanceMonitoring<T extends React.ElementType>(
  Component: T,
  componentName?: string
) {
  const name = componentName || 
    (typeof Component === 'function' ? Component.displayName || Component.name : null) || 
    'UnknownComponent';
  return performanceMonitor.measureComponentRender(Component, name);
}

// API wrapper for measuring fetch performance
export async function monitoredFetch(url: string, options?: RequestInit): Promise<Response> {
  return performanceMonitor.measureAPICall(
    () => fetch(url, options),
    `${options?.method || 'GET'} ${url}`
  );
}

export default performanceMonitor;