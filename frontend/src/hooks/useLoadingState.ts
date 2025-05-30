import { useState, useCallback, useRef, useEffect } from 'react';

interface LoadingState {
  isLoading: boolean;
  error: Error | null;
  data: any;
}

interface UseLoadingStateOptions {
  initialLoading?: boolean;
  minimumLoadingTime?: number;
  onSuccess?: (data: any) => void;
  onError?: (error: Error) => void;
}

export const useLoadingState = (options: UseLoadingStateOptions = {}) => {
  const {
    initialLoading = false,
    minimumLoadingTime = 300,
    onSuccess,
    onError
  } = options;

  const [state, setState] = useState<LoadingState>({
    isLoading: initialLoading,
    error: null,
    data: null
  });

  const loadingStartTime = useRef<number | null>(null);
  const isMountedRef = useRef(true);

  useEffect(() => {
    return () => {
      isMountedRef.current = false;
    };
  }, []);

  const setLoading = useCallback((loading: boolean) => {
    if (!isMountedRef.current) return;
    
    if (loading) {
      loadingStartTime.current = Date.now();
      setState(prev => ({ ...prev, isLoading: true, error: null }));
    } else {
      const elapsed = loadingStartTime.current ? Date.now() - loadingStartTime.current : 0;
      const remainingTime = Math.max(0, minimumLoadingTime - elapsed);
      
      setTimeout(() => {
        if (isMountedRef.current) {
          setState(prev => ({ ...prev, isLoading: false }));
        }
      }, remainingTime);
    }
  }, [minimumLoadingTime]);

  const setData = useCallback((data: any) => {
    if (!isMountedRef.current) return;
    
    setState(prev => ({ ...prev, data, error: null }));
    setLoading(false);
    onSuccess?.(data);
  }, [setLoading, onSuccess]);

  const setError = useCallback((error: Error) => {
    if (!isMountedRef.current) return;
    
    setState(prev => ({ ...prev, error, data: null }));
    setLoading(false);
    onError?.(error);
  }, [setLoading, onError]);

  const reset = useCallback(() => {
    if (!isMountedRef.current) return;
    
    setState({
      isLoading: false,
      error: null,
      data: null
    });
  }, []);

  const execute = useCallback(async <T>(asyncFn: () => Promise<T>): Promise<T | void> => {
    try {
      setLoading(true);
      const result = await asyncFn();
      setData(result);
      return result;
    } catch (error) {
      setError(error instanceof Error ? error : new Error('Unknown error'));
    }
  }, [setLoading, setData, setError]);

  return {
    ...state,
    setLoading,
    setData,
    setError,
    reset,
    execute
  };
};

// Hook for component-level loading states
export const useComponentLoading = (componentName?: string) => {
  const [loadingStates, setLoadingStates] = useState<Record<string, boolean>>({});

  const setLoading = useCallback((key: string, loading: boolean) => {
    setLoadingStates(prev => ({
      ...prev,
      [key]: loading
    }));
  }, []);

  const isLoading = useCallback((key?: string) => {
    if (key) {
      return loadingStates[key] || false;
    }
    return Object.values(loadingStates).some(loading => loading);
  }, [loadingStates]);

  const clearAll = useCallback(() => {
    setLoadingStates({});
  }, []);

  return {
    setLoading,
    isLoading,
    clearAll,
    loadingStates
  };
};

// Hook for async operations with loading state
export const useAsyncOperation = <T = any>(
  operation: () => Promise<T>,
  dependencies: any[] = []
) => {
  const [state, setState] = useState<{
    data: T | null;
    loading: boolean;
    error: Error | null;
  }>({
    data: null,
    loading: false,
    error: null
  });

  const execute = useCallback(async () => {
    setState(prev => ({ ...prev, loading: true, error: null }));
    
    try {
      const result = await operation();
      setState({ data: result, loading: false, error: null });
      return result;
    } catch (error) {
      const errorObj = error instanceof Error ? error : new Error('Unknown error');
      setState({ data: null, loading: false, error: errorObj });
      throw errorObj;
    }
  }, dependencies);

  return {
    ...state,
    execute,
    retry: execute
  };
};

export default useLoadingState;