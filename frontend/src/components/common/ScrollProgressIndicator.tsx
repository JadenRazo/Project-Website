import React, { useEffect, useRef, useCallback } from 'react';
import styled from 'styled-components';
import { useLenis } from '../../providers/LenisProvider';

const ProgressContainer = styled.div<{ $navHeight: number }>`
  position: fixed;
  top: ${({ $navHeight }) => $navHeight}px;
  left: 0;
  width: 100%;
  height: 3px;
  background: ${({ theme }) => theme.colors.surface}20;
  z-index: 1001;
  pointer-events: none;

  @media (max-width: ${({ theme }) => theme.breakpoints.mobile}) {
    height: 4px;
  }
`;

const ProgressBar = styled.div`
  height: 100%;
  background: ${({ theme }) => theme.colors.primary};
  width: 100%;
  will-change: transform;
  transform-origin: left;
  transform: scaleX(0);
  box-shadow: 0 0 10px ${({ theme }) => theme.colors.primary}50;
`;

const ScrollProgressIndicator: React.FC = () => {
  const barRef = useRef<HTMLDivElement>(null);
  const navHeightRef = useRef(0);
  const containerRef = useRef<HTMLDivElement>(null);
  const { lenis } = useLenis();

  const measureNav = useCallback(() => {
    const nav = document.querySelector('nav');
    if (nav) {
      const h = nav.getBoundingClientRect().height;
      navHeightRef.current = h;
      if (containerRef.current) {
        containerRef.current.style.top = `${h}px`;
      }
    }
  }, []);

  useEffect(() => {
    measureNav();
    const timer = setTimeout(measureNav, 200);
    window.addEventListener('resize', measureNav);
    return () => {
      clearTimeout(timer);
      window.removeEventListener('resize', measureNav);
    };
  }, [measureNav]);

  useEffect(() => {
    if (lenis) {
      const onScroll = ({ progress }: { progress: number }) => {
        if (barRef.current) {
          barRef.current.style.transform = `scaleX(${progress})`;
        }
      };

      lenis.on('scroll', onScroll);

      if (barRef.current) {
        barRef.current.style.transform = `scaleX(${lenis.progress || 0})`;
      }

      return () => {
        lenis.off('scroll', onScroll);
      };
    }

    const onNativeScroll = () => {
      if (!barRef.current) return;
      const scrollTop = window.scrollY;
      const docHeight = document.documentElement.scrollHeight - window.innerHeight;
      const progress = docHeight > 0 ? Math.min(scrollTop / docHeight, 1) : 0;
      barRef.current.style.transform = `scaleX(${progress})`;
    };

    onNativeScroll();
    window.addEventListener('scroll', onNativeScroll, { passive: true });
    return () => {
      window.removeEventListener('scroll', onNativeScroll);
    };
  }, [lenis]);

  return (
    <ProgressContainer ref={containerRef} $navHeight={navHeightRef.current}>
      <ProgressBar ref={barRef} />
    </ProgressContainer>
  );
};

export default ScrollProgressIndicator;
