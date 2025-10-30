import React, { useState, useEffect } from 'react';
import styled from 'styled-components';
import { useOptimizedScrollHandler } from '../../hooks/useOptimizedScrollHandler';

interface ProgressBarProps {
  $progress: number;
}

const ProgressContainer = styled.div`
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 3px;
  background: ${({ theme }) => theme.colors.surface}20;
  z-index: 1000;
  transition: opacity ${({ theme }) => theme.transitions.fast};
`;

const ProgressBar = styled.div<ProgressBarProps>`
  height: 100%;
  background: ${({ theme }) => theme.colors.primary};
  width: ${({ $progress }) => $progress}%;
  transition: width 0.1s ease-out;
  box-shadow: 0 0 10px ${({ theme }) => theme.colors.primary}50;
`;

interface ScrollProgressIndicatorProps {
  showOnlyWhenScrolling?: boolean;
  hideThreshold?: number;
}

const ScrollProgressIndicator: React.FC<ScrollProgressIndicatorProps> = ({
  showOnlyWhenScrolling = false,
  hideThreshold = 100
}) => {
  const [progress, setProgress] = useState(0);
  const [isVisible, setIsVisible] = useState(!showOnlyWhenScrolling);

  const handleScroll = (state: { scrollProgress: number; isScrolling: boolean; scrollY: number }) => {
    // Update progress
    setProgress(state.scrollProgress * 100);

    // Handle visibility
    if (showOnlyWhenScrolling) {
      setIsVisible(state.isScrolling && state.scrollY > hideThreshold);
    } else {
      setIsVisible(state.scrollY > hideThreshold);
    }
  };

  useOptimizedScrollHandler(handleScroll, {
    throttleMs: 16 // ~60fps
  });

  // Set initial state
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
    setIsVisible(!showOnlyWhenScrolling && initialScroll > hideThreshold);
  }, [showOnlyWhenScrolling, hideThreshold]);

  if (!isVisible) return null;

  return (
    <ProgressContainer>
      <ProgressBar $progress={progress} />
    </ProgressContainer>
  );
};

export default ScrollProgressIndicator;