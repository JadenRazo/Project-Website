import { useCallback, useRef, useEffect } from 'react';
import { useScrollTo } from './useScrollTo';

interface UseInlineFormScrollOptions {
  autoScroll?: boolean;
  scrollOffset?: number;
  scrollDelay?: number;
}

interface UseInlineFormScrollReturn<T extends HTMLElement = HTMLDivElement> {
  formRef: React.RefObject<T>;
  scrollToForm: () => void;
  triggerScroll: () => void;
}

export function useInlineFormScroll<T extends HTMLElement = HTMLDivElement>(
  isVisible: boolean,
  options: UseInlineFormScrollOptions = {}
): UseInlineFormScrollReturn<T> {
  const {
    autoScroll = true,
    scrollOffset = 80,
    scrollDelay = 150
  } = options;

  const formRef = useRef<T>(null);
  const { scrollToElement } = useScrollTo();

  const scrollToForm = useCallback(() => {
    if (formRef.current) {
      scrollToElement(formRef.current, {
        behavior: 'smooth',
        offset: scrollOffset
      });
    }
  }, [scrollToElement, scrollOffset]);

  const triggerScroll = useCallback(() => {
    setTimeout(scrollToForm, scrollDelay);
  }, [scrollToForm, scrollDelay]);

  useEffect(() => {
    if (isVisible && autoScroll) {
      triggerScroll();
    }
  }, [isVisible, autoScroll, triggerScroll]);

  return {
    formRef,
    scrollToForm,
    triggerScroll
  };
}