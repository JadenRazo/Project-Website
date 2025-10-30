import styled, { keyframes, css } from 'styled-components';

// Intersection observer hook for lazy loading
import { useEffect, useRef, useState } from 'react';

// Smooth fade in animation
export const fadeIn = keyframes`
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
`;

// Scale in animation
export const scaleIn = keyframes`
  from {
    opacity: 0;
    transform: scale(0.95);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
`;

// Slide in from right
export const slideInRight = keyframes`
  from {
    opacity: 0;
    transform: translateX(20px);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
`;

// Pulse animation for loading states
export const pulse = keyframes`
  0% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
  100% {
    opacity: 1;
  }
`;

// Rotate animation for refresh button
export const rotate = keyframes`
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
`;

// Shimmer effect for loading placeholders
export const shimmer = keyframes`
  0% {
    background-position: -1000px 0;
  }
  100% {
    background-position: 1000px 0;
  }
`;

// Smooth expand/collapse
export const expand = keyframes`
  from {
    max-height: 0;
    opacity: 0;
  }
  to {
    max-height: 1000px;
    opacity: 1;
  }
`;

// Animation utilities
export const animationMixin = css<{ 
  delay?: number; 
  duration?: number;
  timingFunction?: string;
}>`
  animation-delay: ${props => props.delay || 0}ms;
  animation-duration: ${props => props.duration || 300}ms;
  animation-timing-function: ${props => props.timingFunction || 'ease-out'};
  animation-fill-mode: both;
`;

// Animated container with stagger effect
export const AnimatedContainer = styled.div<{ 
  staggerDelay?: number;
  animationName?: any;
}>`
  animation-delay: 0ms;
  animation-duration: 300ms;
  animation-timing-function: ease-out;
  animation-fill-mode: both;
  animation-name: ${props => props.animationName || fadeIn};
  
  > * {
    animation-duration: 300ms;
    animation-timing-function: ease-out;
    animation-fill-mode: both;
    animation-name: ${props => props.animationName || fadeIn};
    
    ${props => props.staggerDelay && css`
      &:nth-child(1) { animation-delay: ${props.staggerDelay * 0}ms; }
      &:nth-child(2) { animation-delay: ${props.staggerDelay * 1}ms; }
      &:nth-child(3) { animation-delay: ${props.staggerDelay * 2}ms; }
      &:nth-child(4) { animation-delay: ${props.staggerDelay * 3}ms; }
      &:nth-child(5) { animation-delay: ${props.staggerDelay * 4}ms; }
      &:nth-child(6) { animation-delay: ${props.staggerDelay * 5}ms; }
      &:nth-child(7) { animation-delay: ${props.staggerDelay * 6}ms; }
      &:nth-child(8) { animation-delay: ${props.staggerDelay * 7}ms; }
      &:nth-child(9) { animation-delay: ${props.staggerDelay * 8}ms; }
      &:nth-child(10) { animation-delay: ${props.staggerDelay * 9}ms; }
    `}
  }
`;

// Loading skeleton
export const SkeletonLoader = styled.div<{ 
  width?: string; 
  height?: string;
  borderRadius?: string;
}>`
  width: ${props => props.width || '100%'};
  height: ${props => props.height || '20px'};
  border-radius: ${props => props.borderRadius || '4px'};
  background: linear-gradient(
    90deg,
    ${props => props.theme.colors.border} 0%,
    ${props => props.theme.colors.background} 50%,
    ${props => props.theme.colors.border} 100%
  );
  background-size: 1000px 100%;
  animation: ${shimmer} 2s ease-in-out infinite;
`;

// Smooth transition wrapper
export const SmoothTransition = styled.div<{ 
  isVisible: boolean;
  duration?: number;
}>`
  transition: all ${props => props.duration || 300}ms ease-out;
  opacity: ${props => props.isVisible ? 1 : 0};
  transform: ${props => props.isVisible ? 'translateY(0)' : 'translateY(10px)'};
  pointer-events: ${props => props.isVisible ? 'auto' : 'none'};
`;

// Animated metric card
export const AnimatedCard = styled.div<{ index?: number }>`
  animation-name: ${scaleIn};
  animation-delay: ${props => (props.index || 0) * 50}ms;
  animation-duration: 400ms;
  animation-timing-function: ease-out;
  animation-fill-mode: both;
  transition: all 200ms ease-out;
  
  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 8px 16px rgba(0, 0, 0, 0.1);
  }
`;

// Animated button with ripple effect
export const AnimatedButton = styled.button<{ isLoading?: boolean }>`
  position: relative;
  overflow: hidden;
  transition: all 200ms ease-out;
  
  &::after {
    content: '';
    position: absolute;
    top: 50%;
    left: 50%;
    width: 0;
    height: 0;
    border-radius: 50%;
    background: rgba(255, 255, 255, 0.5);
    transform: translate(-50%, -50%);
    transition: width 600ms, height 600ms;
  }
  
  &:active::after {
    width: 300px;
    height: 300px;
  }
  
  ${props => props.isLoading && css`
    pointer-events: none;
    opacity: 0.7;
    
    &::before {
      content: '';
      position: absolute;
      top: 50%;
      left: 50%;
      width: 16px;
      height: 16px;
      margin: -8px 0 0 -8px;
      border: 2px solid currentColor;
      border-right-color: transparent;
      border-radius: 50%;
      animation: ${rotate} 1s linear infinite;
    }
  `}
`;

// Performance-optimized list container
export const OptimizedList = styled.div`
  will-change: transform;
  transform: translateZ(0);
  backface-visibility: hidden;
  -webkit-font-smoothing: antialiased;
  
  > * {
    will-change: opacity, transform;
  }
`;

export const useIntersectionObserver = (options?: IntersectionObserverInit) => {
  const ref = useRef<HTMLElement>(null);
  const [isIntersecting, setIsIntersecting] = useState(false);
  const [hasIntersected, setHasIntersected] = useState(false);

  useEffect(() => {
    const observer = new IntersectionObserver(([entry]) => {
      setIsIntersecting(entry.isIntersecting);
      if (entry.isIntersecting) {
        setHasIntersected(true);
      }
    }, options);

    if (ref.current) {
      observer.observe(ref.current);
    }

    return () => {
      if (ref.current) {
        observer.unobserve(ref.current);
      }
    };
  }, [options]);

  return { ref, isIntersecting, hasIntersected };
};

// Animation hook for reduced motion preference
export const useReducedMotion = () => {
  const [reducedMotion, setReducedMotion] = useState(false);

  useEffect(() => {
    const mediaQuery = window.matchMedia('(prefers-reduced-motion: reduce)');
    setReducedMotion(mediaQuery.matches);

    const handleChange = (e: MediaQueryListEvent) => {
      setReducedMotion(e.matches);
    };

    mediaQuery.addEventListener('change', handleChange);
    return () => mediaQuery.removeEventListener('change', handleChange);
  }, []);

  return reducedMotion;
};