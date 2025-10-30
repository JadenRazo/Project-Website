import React, { useRef, useEffect, useState } from 'react';
import styled from 'styled-components';
import { motion, Variants } from 'framer-motion';

interface ScrollRevealProps {
  children: React.ReactNode;
  threshold?: number;
  rootMargin?: string;
  delay?: number;
  duration?: number;
  direction?: 'up' | 'down' | 'left' | 'right' | 'fade';
  once?: boolean;
  className?: string;
}

const RevealContainer = styled(motion.div)`
  width: 100%;
`;

const createVariants = (direction: string, duration: number): Variants => {
  const distance = 50;
  
  const hidden = {
    opacity: 0,
    ...(direction === 'up' && { y: distance }),
    ...(direction === 'down' && { y: -distance }),
    ...(direction === 'left' && { x: distance }),
    ...(direction === 'right' && { x: -distance }),
  };
  
  const visible = {
    opacity: 1,
    x: 0,
    y: 0,
    transition: {
      duration: duration / 1000,
      ease: 'easeOut',
    },
  };
  
  return { hidden, visible };
};

const ScrollReveal: React.FC<ScrollRevealProps> = ({
  children,
  threshold = 0.1,
  rootMargin = '0px',
  delay = 0,
  duration = 600,
  direction = 'up',
  once = true,
  className,
}) => {
  const ref = useRef<HTMLDivElement>(null);
  const [isVisible, setIsVisible] = useState(false);
  const [hasAnimated, setHasAnimated] = useState(false);
  
  useEffect(() => {
    const element = ref.current;
    if (!element) return;
    
    // Skip if already animated and once is true
    if (once && hasAnimated) return;
    
    const observer = new IntersectionObserver(
      ([entry]) => {
        if (entry.isIntersecting) {
          setIsVisible(true);
          if (once) {
            setHasAnimated(true);
            observer.unobserve(element);
          }
        } else if (!once) {
          setIsVisible(false);
        }
      },
      {
        threshold,
        rootMargin,
      }
    );
    
    observer.observe(element);
    
    return () => {
      observer.disconnect();
    };
  }, [threshold, rootMargin, once, hasAnimated]);
  
  const variants = createVariants(direction, duration);
  
  return (
    <RevealContainer
      ref={ref}
      initial="hidden"
      animate={isVisible ? 'visible' : 'hidden'}
      variants={variants}
      className={className}
      style={{ transitionDelay: `${delay}ms` }}
    >
      {children}
    </RevealContainer>
  );
};

export default ScrollReveal;