import { useCallback, useEffect } from 'react';
import { 
  scrollToElement, 
  scrollToTop, 
  scrollToPosition,
  ScrollOptions,
  getScrollBehavior 
} from '../utils/scrollConfig';

interface UseScrollToReturn {
  scrollToElement: (element: Element | null, options?: ScrollOptions) => Promise<void>;
  scrollToTop: (options?: ScrollOptions) => Promise<void>;
  scrollToPosition: (position: number, options?: ScrollOptions) => Promise<void>;
  scrollToId: (elementId: string, options?: ScrollOptions) => Promise<void>;
}

export const useScrollTo = (): UseScrollToReturn => {
  const scrollToElementCallback = useCallback(
    (element: Element | null, options?: ScrollOptions) => {
      const behavior = options?.behavior || getScrollBehavior();
      return scrollToElement(element, { ...options, behavior });
    },
    []
  );

  const scrollToTopCallback = useCallback((options?: ScrollOptions) => {
    const behavior = options?.behavior || getScrollBehavior();
    return scrollToTop({ ...options, behavior });
  }, []);

  const scrollToPositionCallback = useCallback(
    (position: number, options?: ScrollOptions) => {
      const behavior = options?.behavior || getScrollBehavior();
      return scrollToPosition(position, { ...options, behavior });
    },
    []
  );

  const scrollToIdCallback = useCallback(
    (elementId: string, options?: ScrollOptions) => {
      const element = document.getElementById(elementId);
      if (element) {
        return scrollToElementCallback(element, options);
      }
      return Promise.resolve();
    },
    [scrollToElementCallback]
  );

  return {
    scrollToElement: scrollToElementCallback,
    scrollToTop: scrollToTopCallback,
    scrollToPosition: scrollToPositionCallback,
    scrollToId: scrollToIdCallback,
  };
};