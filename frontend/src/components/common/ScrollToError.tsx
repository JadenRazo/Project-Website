import React, { useEffect, useRef } from 'react';
import { useScrollTo } from '../../hooks/useScrollTo';

interface ScrollToErrorProps {
  error: string | null;
  children: React.ReactNode;
  offset?: number;
}

export const ScrollToError: React.FC<ScrollToErrorProps> = ({ 
  error, 
  children, 
  offset = 80 
}) => {
  const errorRef = useRef<HTMLDivElement>(null);
  const { scrollToElement } = useScrollTo();

  useEffect(() => {
    if (error && errorRef.current) {
      scrollToElement(errorRef.current, { 
        behavior: 'smooth',
        offset
      });
    }
  }, [error, scrollToElement, offset]);

  if (!error) return null;

  return <div ref={errorRef}>{children}</div>;
};

export default ScrollToError;