import { useMemo, useRef, useState, useEffect, useCallback } from 'react';
import { useInView } from 'framer-motion';
import { usePerformanceOptimizations } from '../../../hooks/usePerformanceOptimizations';
import type { PerformanceConfig, AnimationPhase } from './types';
import { getPerformanceConfig } from './constants';

interface UseLegoAnimationReturn {
  config: PerformanceConfig;
  skipAnimation: boolean;
  isVisible: boolean;
  containerRef: React.RefObject<HTMLDivElement>;
  phase: AnimationPhase;
  setPhase: (phase: AnimationPhase) => void;
  showBio: boolean;
  setShowBio: (show: boolean) => void;
  dimensions: { width: number; height: number };
}

export const useLegoAnimation = (): UseLegoAnimationReturn => {
  const { performanceSettings } = usePerformanceOptimizations();
  const containerRef = useRef<HTMLDivElement>(null);
  const isInView = useInView(containerRef, { amount: 0.1, once: false });
  const [phase, setPhase] = useState<AnimationPhase>('build');
  const [showBio, setShowBio] = useState(false);
  const [dimensions, setDimensions] = useState({ width: 0, height: 0 });
  const [isInitialized, setIsInitialized] = useState(false);

  const skipAnimation = useMemo(() => {
    return performanceSettings.reduceMotion;
  }, [performanceSettings.reduceMotion]);

  const config = useMemo(() => {
    return getPerformanceConfig(performanceSettings.performanceTier);
  }, [performanceSettings.performanceTier]);

  const updateDimensions = useCallback(() => {
    let newWidth = window.innerWidth;
    let newHeight = window.innerHeight;

    if (containerRef.current) {
      const rect = containerRef.current.getBoundingClientRect();
      if (rect.width > 0) newWidth = rect.width;
      if (rect.height > 0) newHeight = rect.height;
    }

    newHeight = Math.max(newHeight, window.innerHeight * 0.8);

    if (newWidth > 0 && newHeight > 0) {
      setDimensions({
        width: newWidth,
        height: newHeight,
      });
      if (!isInitialized) {
        setIsInitialized(true);
      }
    }
  }, [isInitialized]);

  useEffect(() => {
    updateDimensions();

    const timeoutId = setTimeout(updateDimensions, 100);
    const timeoutId2 = setTimeout(updateDimensions, 500);

    const resizeObserver = new ResizeObserver(() => {
      updateDimensions();
    });

    if (containerRef.current) {
      resizeObserver.observe(containerRef.current);
    }

    window.addEventListener('resize', updateDimensions);

    return () => {
      clearTimeout(timeoutId);
      clearTimeout(timeoutId2);
      resizeObserver.disconnect();
      window.removeEventListener('resize', updateDimensions);
    };
  }, [updateDimensions]);

  useEffect(() => {
    if (skipAnimation) {
      setShowBio(true);
      setPhase('revealed');
    }
  }, [skipAnimation]);

  const effectiveIsVisible = isInView || !isInitialized;

  return {
    config,
    skipAnimation,
    isVisible: effectiveIsVisible,
    containerRef: containerRef as React.RefObject<HTMLDivElement>,
    phase,
    setPhase,
    showBio,
    setShowBio,
    dimensions,
  };
};

export const useAnimationLoop = (
  callback: (timestamp: number, deltaTime: number) => void,
  isActive: boolean,
  targetFrameMs: number
): void => {
  const animationRef = useRef<number>();
  const lastTimeRef = useRef<number>(0);
  const callbackRef = useRef(callback);

  callbackRef.current = callback;

  useEffect(() => {
    if (!isActive) {
      if (animationRef.current) {
        cancelAnimationFrame(animationRef.current);
      }
      return;
    }

    const animate = (timestamp: number) => {
      const deltaTime = lastTimeRef.current ? timestamp - lastTimeRef.current : targetFrameMs;

      if (deltaTime >= targetFrameMs * 0.8) {
        callbackRef.current(timestamp, Math.min(deltaTime, targetFrameMs * 2));
        lastTimeRef.current = timestamp;
      }

      animationRef.current = requestAnimationFrame(animate);
    };

    animationRef.current = requestAnimationFrame(animate);

    return () => {
      if (animationRef.current) {
        cancelAnimationFrame(animationRef.current);
      }
    };
  }, [isActive, targetFrameMs]);
};

export default useLegoAnimation;
