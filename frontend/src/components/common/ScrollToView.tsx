import React, { useEffect, useRef } from 'react';
import { useScrollTo } from '../../hooks/useScrollTo';

interface ScrollToViewProps {
  children: React.ReactNode;
  when: boolean;
  offset?: number;
  delay?: number;
  behavior?: ScrollBehavior;
  className?: string;
  as?: keyof JSX.IntrinsicElements;
}

export const ScrollToView: React.FC<ScrollToViewProps> = ({
  children,
  when,
  offset = 80,
  delay = 100,
  behavior = 'smooth',
  className,
  as: Component = 'div'
}) => {
  const elementRef = useRef<HTMLElement>(null);
  const { scrollToElement } = useScrollTo();

  useEffect(() => {
    if (when && elementRef.current) {
      const timeoutId = setTimeout(() => {
        scrollToElement(elementRef.current, {
          behavior,
          offset
        });
      }, delay);

      return () => clearTimeout(timeoutId);
    }
  }, [when, scrollToElement, offset, delay, behavior]);

  return React.createElement(
    Component as any,
    { ref: elementRef, className },
    children
  );
};