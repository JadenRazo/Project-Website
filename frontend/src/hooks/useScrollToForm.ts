import { useRef, useCallback } from 'react';
import { useScrollTo } from './useScrollTo';
import { getScrollOffset } from '../utils/scrollConfig';

interface UseScrollToFormReturn {
  formRef: React.RefObject<HTMLDivElement>;
  scrollToForm: () => void;
  scrollToFormWithDelay: (delay?: number) => void;
}

export const useScrollToForm = (): UseScrollToFormReturn => {
  const formRef = useRef<HTMLDivElement>(null);
  const { scrollToElement } = useScrollTo();

  const scrollToForm = useCallback(() => {
    if (formRef.current) {
      // Add extra offset for better UX
      const extraOffset = 20;
      scrollToElement(formRef.current, { 
        offset: getScrollOffset() + extraOffset,
        behavior: 'smooth' 
      });
    }
  }, [scrollToElement]);

  const scrollToFormWithDelay = useCallback((delay: number = 100) => {
    setTimeout(() => {
      scrollToForm();
    }, delay);
  }, [scrollToForm]);

  return {
    formRef,
    scrollToForm,
    scrollToFormWithDelay,
  };
};