import { useRef, useEffect, useCallback } from 'react';
import { throttle } from 'lodash';

interface ScrollHandlerOptions {
  threshold?: number;
  throttleMs?: number;
  enableParallax?: boolean;
}

interface ScrollState {
  scrollY: number;
  scrollDirection: 'up' | 'down' | null;
  scrollProgress: number; // 0-1
  isScrolling: boolean;
}

type ScrollCallback = (state: ScrollState) => void;

/**
 * Optimized scroll handler that replaces the complex ScrollTransformBackground
 * Provides efficient scroll event handling with minimal performance impact
 */
export const useOptimizedScrollHandler = (
  callback: ScrollCallback,
  options: ScrollHandlerOptions = {}
) => {
  const {
    threshold = 5,
    throttleMs = 16, // ~60fps
    enableParallax = true
  } = options;
  
  const lastScrollY = useRef(0);
  const scrollDirection = useRef<'up' | 'down' | null>(null);
  const isScrolling = useRef(false);
  const scrollTimeoutRef = useRef<number | null>(null);
  const ticking = useRef(false);
  
  // Use requestAnimationFrame for smooth updates
  const updateScroll = useCallback(() => {
    const currentScrollY = window.pageYOffset || document.documentElement.scrollTop;
    const scrollDelta = currentScrollY - lastScrollY.current;
    
    // Only update if scroll change is significant
    if (Math.abs(scrollDelta) < threshold) {
      ticking.current = false;
      return;
    }
    
    // Determine scroll direction
    const newDirection = scrollDelta > 0 ? 'down' : 'up';
    if (newDirection !== scrollDirection.current) {
      scrollDirection.current = newDirection;
    }
    
    // Calculate scroll progress
    const documentHeight = Math.max(
      document.body.scrollHeight,
      document.documentElement.scrollHeight
    );
    const windowHeight = window.innerHeight;
    const maxScroll = documentHeight - windowHeight;
    const scrollProgress = maxScroll > 0 ? Math.min(currentScrollY / maxScroll, 1) : 0;
    
    // Mark as scrolling
    isScrolling.current = true;
    
    // Create scroll state
    const scrollState: ScrollState = {
      scrollY: currentScrollY,
      scrollDirection: scrollDirection.current,
      scrollProgress,
      isScrolling: isScrolling.current
    };
    
    // Call the callback
    callback(scrollState);
    
    // Update last scroll position
    lastScrollY.current = currentScrollY;
    
    // Reset scrolling flag after a delay
    if (scrollTimeoutRef.current) {
      clearTimeout(scrollTimeoutRef.current);
    }
    
    scrollTimeoutRef.current = window.setTimeout(() => {
      isScrolling.current = false;
      callback({
        ...scrollState,
        isScrolling: false
      });
    }, 150);
    
    ticking.current = false;
  }, [callback, threshold]);
  
  // Throttled scroll handler
  const handleScroll = useCallback(() => {
    if (!ticking.current) {
      requestAnimationFrame(updateScroll);
      ticking.current = true;
    }
  }, [updateScroll]);
  
  // Throttled version for additional performance
  const throttledScrollHandler = useCallback(
    throttle(handleScroll, throttleMs, { leading: true, trailing: true }),
    [handleScroll, throttleMs]
  );
  
  useEffect(() => {
    if (!enableParallax) return;
    
    // Use passive listeners for better performance
    window.addEventListener('scroll', throttledScrollHandler, { passive: true });
    
    // Initial call
    updateScroll();
    
    return () => {
      window.removeEventListener('scroll', throttledScrollHandler);
      if (scrollTimeoutRef.current) {
        clearTimeout(scrollTimeoutRef.current);
      }
      throttledScrollHandler.cancel();
    };
  }, [throttledScrollHandler, enableParallax, updateScroll]);
  
  // Return current scroll state
  return {
    scrollY: lastScrollY.current,
    scrollDirection: scrollDirection.current,
    isScrolling: isScrolling.current
  };
};

/**
 * Simple parallax effect hook for individual elements
 */
export const useParallaxEffect = (
  elementRef: React.RefObject<HTMLElement>,
  speed: number = 0.5
) => {
  const handleScroll = useCallback((state: ScrollState) => {
    if (!elementRef.current) return;
    
    const element = elementRef.current;
    const rect = element.getBoundingClientRect();
    const elementTop = rect.top + state.scrollY;
    const elementHeight = rect.height;
    const windowHeight = window.innerHeight;
    
    // Check if element is in viewport
    const isInViewport = rect.top < windowHeight && rect.bottom > 0;
    
    if (isInViewport) {
      // Calculate parallax offset
      const elementProgress = (state.scrollY - elementTop + windowHeight) / (windowHeight + elementHeight);
      const parallaxOffset = elementProgress * speed * 100;
      
      // Apply transform
      element.style.transform = `translateY(${parallaxOffset}px)`;
    }
  }, [elementRef, speed]);
  
  useOptimizedScrollHandler(handleScroll, {
    threshold: 2,
    throttleMs: 16
  });
};


/**
 * Intersection Observer based visibility detection
 * More efficient than scroll-based visibility checks
 */
export const useInViewport = (
  elementRef: React.RefObject<HTMLElement>,
  options: IntersectionObserverInit = {}
) => {
  const isInViewRef = useRef(false);
  const observerRef = useRef<IntersectionObserver | null>(null);
  
  useEffect(() => {
    if (!elementRef.current) return;
    
    const element = elementRef.current;
    
    observerRef.current = new IntersectionObserver(
      ([entry]) => {
        isInViewRef.current = entry.isIntersecting;
      },
      {
        threshold: 0.1,
        rootMargin: '50px',
        ...options
      }
    );
    
    observerRef.current.observe(element);
    
    return () => {
      if (observerRef.current) {
        observerRef.current.disconnect();
      }
    };
  }, [elementRef, options]);
  
  return isInViewRef.current;
};

/**
 * Optimized virtual scrolling for large lists
 */
export const useVirtualScrolling = (
  itemHeight: number,
  totalItems: number,
  containerHeight: number
) => {
  const scrollTop = useRef(0);
  const startIndex = useRef(0);
  const endIndex = useRef(0);
  const visibleItems = useRef(0);
  
  const handleScroll = useCallback((state: ScrollState) => {
    scrollTop.current = state.scrollY;
    
    // Calculate visible range
    const start = Math.floor(scrollTop.current / itemHeight);
    const visible = Math.ceil(containerHeight / itemHeight);
    const end = Math.min(start + visible + 1, totalItems); // +1 for buffer
    
    startIndex.current = Math.max(0, start - 1); // -1 for buffer
    endIndex.current = end;
    visibleItems.current = end - startIndex.current;
  }, [itemHeight, totalItems, containerHeight]);
  
  useOptimizedScrollHandler(handleScroll);
  
  return {
    startIndex: startIndex.current,
    endIndex: endIndex.current,
    visibleItems: visibleItems.current,
    offsetY: startIndex.current * itemHeight
  };
};

/**
 * Performance-optimized scroll reveal animation
 */
export const useScrollReveal = (
  elements: React.RefObject<HTMLElement>[],
  options: {
    threshold?: number;
    distance?: number;
    duration?: number;
    delay?: number;
  } = {}
) => {
  const {
    threshold = 0.1,
    distance = 50,
    duration = 600,
    delay = 0
  } = options;
  
  const revealedElements = useRef(new Set<HTMLElement>());
  
  useEffect(() => {
    const observer = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting && !revealedElements.current.has(entry.target as HTMLElement)) {
            const element = entry.target as HTMLElement;
            
            // Apply reveal animation
            element.style.transition = `transform ${duration}ms ease-out, opacity ${duration}ms ease-out`;
            element.style.transform = 'translateY(0)';
            element.style.opacity = '1';
            
            revealedElements.current.add(element);
          }
        });
      },
      { threshold }
    );
    
    // Set initial state and observe elements
    elements.forEach((ref) => {
      if (ref.current) {
        const element = ref.current;
        element.style.transform = `translateY(${distance}px)`;
        element.style.opacity = '0';
        element.style.transition = 'none';
        
        setTimeout(() => {
          observer.observe(element);
        }, delay);
      }
    });
    
    return () => {
      observer.disconnect();
    };
  }, [elements, threshold, distance, duration, delay]);
};