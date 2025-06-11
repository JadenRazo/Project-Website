import { usePerformanceStore } from '../stores';

export const useMemoryManager = () => {
  const {
    memoryUsage,
    isMemoryConstrained,
    applicationState,
    freeMemory,
    optimizePerformance,
    updateApplicationState,
  } = usePerformanceStore();

  return {
    memoryUsage,
    isMemoryConstrained,
    applicationState,
    freeMemory,
    optimizePerformance,
    updateApplicationState,
  };
};