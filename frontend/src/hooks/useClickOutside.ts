// src/hooks/useClickOutside.ts
import { useEffect, RefObject } from 'react';

/**
 * Hook that detects clicks outside of a specified element
 * @param ref - Reference to the element to monitor
 * @param handler - Function to call when a click outside occurs
 */
export const useClickOutside = <T extends HTMLElement = HTMLElement>(
  ref: RefObject<T>,
  handler: (event: MouseEvent | TouchEvent) => void
): void => {
  useEffect(() => {
    const listener = (event: MouseEvent | TouchEvent) => {
      // Handle the case where ref might be null
      const el = ref.current;
      
      // Don't do anything if the ref doesn't exist or if clicking ref element or descendents
      if (!el || el.contains(event.target as Node)) {
        return;
      }
      
      handler(event);
    };

    document.addEventListener('mousedown', listener);
    document.addEventListener('touchstart', listener);

    return () => {
      document.removeEventListener('mousedown', listener);
      document.removeEventListener('touchstart', listener);
    };
  }, [ref, handler]);
};
