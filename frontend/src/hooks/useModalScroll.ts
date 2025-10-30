import { useEffect, useRef } from 'react';
import { useScrollTo } from './useScrollTo';

interface UseModalScrollOptions {
  offset?: number;
  delay?: number;
}

interface UseModalScrollReturn<T extends HTMLElement = HTMLDivElement> {
  modalRef: React.RefObject<T>;
  scrollToModal: () => void;
}

export function useModalScroll<T extends HTMLElement = HTMLDivElement>(
  isOpen: boolean,
  options: UseModalScrollOptions = {}
): UseModalScrollReturn<T> {
  const { offset = 100, delay = 100 } = options;
  const modalRef = useRef<T>(null);
  const { scrollToElement } = useScrollTo();

  const scrollToModal = () => {
    if (modalRef.current) {
      setTimeout(() => {
        scrollToElement(modalRef.current, {
          behavior: 'smooth',
          offset
        });
      }, delay);
    }
  };

  useEffect(() => {
    if (isOpen && modalRef.current) {
      scrollToModal();
    }
  }, [isOpen]); // eslint-disable-line react-hooks/exhaustive-deps

  return {
    modalRef,
    scrollToModal
  };
}

export default useModalScroll;