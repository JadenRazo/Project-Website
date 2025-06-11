import { create } from 'zustand';
import { devtools } from 'zustand/middleware';
import { captureMemoryMetrics, monitorMemoryUsage } from '../utils/performance';
import type { PerformanceState, PerformanceMetrics, ApplicationState } from './types';

interface PerformanceActions {
  freeMemory: () => PerformanceMetrics;
  optimizePerformance: (level?: 'low' | 'medium' | 'high') => ApplicationState;
  updateMemoryMetrics: () => void;
  setMemoryConstrained: (constrained: boolean) => void;
  updateApplicationState: (updates: Partial<ApplicationState>) => void;
  startMonitoring: (interval?: number, threshold?: number) => void;
  stopMonitoring: () => void;
}

type PerformanceStore = PerformanceState & PerformanceActions;

const isPerformanceApiAvailable = typeof performance !== 'undefined' && 
  typeof performance.now === 'function';

const isMemoryApiAvailable = isPerformanceApiAvailable && 
  typeof (performance as any).memory !== 'undefined';

const isCachesApiAvailable = typeof caches !== 'undefined';

let monitoringInterval: NodeJS.Timeout | null = null;

export const usePerformanceStore = create<PerformanceStore>()(
  devtools(
    (set, get) => ({
      memoryUsage: { timestamp: Date.now() },
      isMemoryConstrained: false,
      applicationState: {
        effectsEnabled: true,
        animationsEnabled: true,
        backgroundEffectsEnabled: true,
        highQualityImagesEnabled: true,
        virtualizationEnabled: true,
      },

      updateMemoryMetrics: () => {
        if (!isMemoryApiAvailable) return;
        
        const metrics = captureMemoryMetrics();
        set({ memoryUsage: metrics }, false, 'updateMemoryMetrics');
      },

      setMemoryConstrained: (constrained: boolean) => {
        set({ isMemoryConstrained: constrained }, false, 'setMemoryConstrained');
      },

      updateApplicationState: (updates: Partial<ApplicationState>) => {
        set(
          (state) => ({
            applicationState: { ...state.applicationState, ...updates }
          }),
          false,
          'updateApplicationState'
        );
      },

      freeMemory: () => {
        const images = document.querySelectorAll('img[data-src]');
        images.forEach(img => {
          if (!img.classList.contains('critical-image')) {
            (img as HTMLImageElement).src = 'data:image/gif;base64,R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7';
          }
        });
        
        if (isCachesApiAvailable) {
          caches.keys().then(cacheNames => {
            cacheNames.forEach(cacheName => {
              if (cacheName.includes('non-critical')) {
                caches.delete(cacheName);
              }
            });
          });
        }
        
        if (typeof (window as any).gc === 'function') {
          try {
            (window as any).gc();
          } catch (e) {
            console.debug('Manual GC not available');
          }
        }
        
        if (process.env.NODE_ENV === 'development') {
          console.info('%cMemory cleanup performed', 'color: #4ecdc4;');
        }
        
        if (isMemoryApiAvailable) {
          const updatedMetrics = captureMemoryMetrics();
          set({ memoryUsage: updatedMetrics }, false, 'freeMemory');
          return updatedMetrics;
        }
        
        return { timestamp: Date.now() };
      },

      optimizePerformance: (level: 'low' | 'medium' | 'high' = 'medium') => {
        let newApplicationState: ApplicationState;
        
        switch (level) {
          case 'low':
            newApplicationState = {
              effectsEnabled: true,
              animationsEnabled: true,
              backgroundEffectsEnabled: false,
              highQualityImagesEnabled: true,
              virtualizationEnabled: true,
            };
            break;
            
          case 'medium':
            newApplicationState = {
              effectsEnabled: true,
              animationsEnabled: true,
              backgroundEffectsEnabled: false,
              highQualityImagesEnabled: false,
              virtualizationEnabled: true,
            };
            break;
            
          case 'high':
            newApplicationState = {
              effectsEnabled: false,
              animationsEnabled: false,
              backgroundEffectsEnabled: false,
              highQualityImagesEnabled: false,
              virtualizationEnabled: true,
            };
            break;
        }
        
        set({ applicationState: newApplicationState }, false, `optimizePerformance/${level}`);
        
        setTimeout(() => {
          get().freeMemory();
        }, 500);
        
        return newApplicationState;
      },

      startMonitoring: (interval: number = 30000, threshold: number = 80) => {
        if (!isMemoryApiAvailable) {
          if (process.env.NODE_ENV === 'development') {
            console.info('Memory API not available, memory monitoring disabled');
          }
          return;
        }

        if (monitoringInterval) {
          clearInterval(monitoringInterval);
        }

        const updateMetrics = () => {
          const metrics = captureMemoryMetrics();
          set({ memoryUsage: metrics }, false, 'monitoring/updateMetrics');
          
          if (metrics.usedJSHeapSize && metrics.jsHeapSizeLimit) {
            const usagePercent = (metrics.usedJSHeapSize / metrics.jsHeapSizeLimit) * 100;
            const newMemoryConstrained = usagePercent > threshold;
            
            const currentState = get();
            if (newMemoryConstrained !== currentState.isMemoryConstrained) {
              set({ isMemoryConstrained: newMemoryConstrained }, false, 'monitoring/memoryConstraints');
              
              if (newMemoryConstrained) {
                setTimeout(() => get().optimizePerformance('medium'), 0);
                
                if (process.env.NODE_ENV === 'development') {
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

        updateMetrics();
        monitoringInterval = setInterval(updateMetrics, interval);

        if (process.env.NODE_ENV === 'development') {
          monitorMemoryUsage('Application', {
            logToConsole: true,
            thresholdPercent: threshold,
            logLevel: 'warn'
          });
        }
      },

      stopMonitoring: () => {
        if (monitoringInterval) {
          clearInterval(monitoringInterval);
          monitoringInterval = null;
        }
      },
    }),
    { name: 'performance-store' }
  )
);

export const initializePerformanceMonitoring = (interval?: number, threshold?: number) => {
  usePerformanceStore.getState().startMonitoring(interval, threshold);
};