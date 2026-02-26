import React, { useState, useEffect, useRef } from 'react';
import styled from 'styled-components';
import { useOptimizedScrollHandler } from '../../hooks/useOptimizedScrollHandler';

interface ProgressBarProps {
  $progress: number;
}

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

const ProgressBar = styled.div<ProgressBarProps>`
  height: 100%;
  background: ${({ theme }) => theme.colors.primary};
  width: ${({ $progress }) => $progress}%;
  transition: width 0.1s ease-out;
  box-shadow: 0 0 10px ${({ theme }) => theme.colors.primary}50;
`;

const ScrollProgressIndicator: React.FC = () => {
  const [progress, setProgress] = useState(0);
  const [navHeight, setNavHeight] = useState(0);
  const measuredRef = useRef(false);

  useEffect(() => {
    const measureNav = () => {
      const nav = document.querySelector('nav');
      if (nav) {
        setNavHeight(nav.getBoundingClientRect().height);
        measuredRef.current = true;
      }
    };

    measureNav();

    if (!measuredRef.current) {
      const timer = setTimeout(measureNav, 200);
      return () => clearTimeout(timer);
    }

    window.addEventListener('resize', measureNav);
    return () => window.removeEventListener('resize', measureNav);
  }, []);

  const handleScroll = (state: { scrollProgress: number; isScrolling: boolean; scrollY: number }) => {
    setProgress(state.scrollProgress * 100);
  };

  useOptimizedScrollHandler(handleScroll, {
    throttleMs: 16
  });

  useEffect(() => {
    const initialScroll = window.pageYOffset || document.documentElement.scrollTop;
    const documentHeight = Math.max(
      document.body.scrollHeight,
      document.documentElement.scrollHeight
    );
    const windowHeight = window.innerHeight;
    const maxScroll = documentHeight - windowHeight;
    const initialProgress = maxScroll > 0 ? Math.min(initialScroll / maxScroll, 1) : 0;

    setProgress(initialProgress * 100);
  }, []);

  return (
    <ProgressContainer $navHeight={navHeight}>
      <ProgressBar $progress={progress} />
    </ProgressContainer>
  );
};

export default ScrollProgressIndicator;
