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

type ObserverCallback = (isIntersecting: boolean) => void;

interface SharedObserverEntry {
  callback: ObserverCallback;
  once: boolean;
  hasTriggered: boolean;
}

const observerCallbacks = new Map<Element, SharedObserverEntry>();

let sharedObserver: IntersectionObserver | null = null;
let currentThreshold = 0.1;
let currentRootMargin = '-50px 0px';

function getSharedObserver(threshold: number, rootMargin: string): IntersectionObserver {
  if (
    sharedObserver &&
    threshold === currentThreshold &&
    rootMargin === currentRootMargin
  ) {
    return sharedObserver;
  }

  if (sharedObserver) {
    sharedObserver.disconnect();
  }

  currentThreshold = threshold;
  currentRootMargin = rootMargin;

  sharedObserver = new IntersectionObserver(
    (entries) => {
      entries.forEach((entry) => {
        const record = observerCallbacks.get(entry.target);
        if (!record) return;

        if (entry.isIntersecting) {
          record.callback(true);
          if (record.once) {
            record.hasTriggered = true;
            sharedObserver?.unobserve(entry.target);
            observerCallbacks.delete(entry.target);
          }
        } else if (!record.once) {
          record.callback(false);
        }
      });
    },
    { threshold, rootMargin }
  );

  return sharedObserver;
}

function registerElement(
  element: Element,
  callback: ObserverCallback,
  once: boolean,
  threshold: number,
  rootMargin: string
): void {
  const observer = getSharedObserver(threshold, rootMargin);
  observerCallbacks.set(element, { callback, once, hasTriggered: false });
  observer.observe(element);
}

function unregisterElement(element: Element): void {
  sharedObserver?.unobserve(element);
  observerCallbacks.delete(element);
}

const ScrollReveal: React.FC<ScrollRevealProps> = ({
  children,
  threshold = 0.1,
  rootMargin = '-50px 0px',
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

    if (once && hasAnimated) return;

    registerElement(
      element,
      (intersecting) => {
        setIsVisible(intersecting);
        if (intersecting && once) {
          setHasAnimated(true);
        }
      },
      once,
      threshold,
      rootMargin
    );

    return () => {
      unregisterElement(element);
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
