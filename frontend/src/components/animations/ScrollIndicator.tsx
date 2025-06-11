import React, { useState, useCallback, useEffect } from 'react';
import styled from 'styled-components';
import { motion, useScroll, useMotionValueEvent, AnimatePresence } from 'framer-motion';
import { useTheme } from '../../hooks/useTheme';

interface ScrollIndicatorProps {
  targetId?: string;
  offset?: number;
  showAboveFold?: boolean;
}

/**
 * Styled container for the scroll indicator
 * Uses either fixed or absolute positioning based on context
 */
const ScrollContainer = styled(motion.div)`
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.8rem;
  cursor: pointer;
  color: ${({ theme }) => theme.colors.text};
  width: max-content;
  position: relative;
  pointer-events: auto;
  will-change: transform, opacity;
  
  &::before {
    content: '';
    position: absolute;
    top: -15px;
    left: -15px;
    right: -15px;
    bottom: -15px;
    border-radius: 16px;
    z-index: -1;
    background: ${({ theme }) => `${theme.colors.background}40`};
    backdrop-filter: blur(5px);
    opacity: 0;
    transition: opacity 0.2s ease;
  }
  
  &:hover::before {
    opacity: 1;
  }
  
  @media (max-width: 768px) {
    transform: scale(0.9);
  }
  
  @media (max-width: 480px) {
    transform: scale(0.85);
  }
`;

const ScrollText = styled.span`
  font-size: 0.9rem;
  font-weight: 500;
  letter-spacing: 0.05em;
  text-transform: uppercase;
  color: ${props => props.theme.colors.text};
  
  @media (max-width: 480px) {
    font-size: 0.8rem;
  }
`;

const ScrollArrow = styled(motion.div)`
  color: ${props => props.theme.colors.primary};
  
  svg {
    width: 24px;
    height: 14px;
    
    @media (max-width: 480px) {
      width: 20px;
      height: 12px;
    }
  }
`;

const arrowVariants = {
  animate: {
    y: [0, 8, 0],
    transition: { 
      repeat: Infinity,
      duration: 1.2,
      ease: "easeInOut" 
    }
  }
};

const containerVariants = {
  hidden: { 
    opacity: 0, 
    y: 10,
    transition: {
      duration: 0.2,
      ease: [0.43, 0.13, 0.23, 0.96]
    }
  },
  visible: { 
    opacity: 1, 
    y: 0,
    transition: {
      duration: 0.3,
      ease: [0.16, 1, 0.3, 1]
    }
  }
};

export const ScrollIndicator: React.FC<ScrollIndicatorProps> = ({ 
  targetId = 'skills',
  offset = 80,
  showAboveFold = true
}) => {
  const [showIndicator, setShowIndicator] = useState(showAboveFold);
  const { scrollY } = useScroll();
  const [isTouchDevice, setIsTouchDevice] = useState(false);
  
  // Detect touch devices on mount
  useEffect(() => {
    setIsTouchDevice('ontouchstart' in window || navigator.maxTouchPoints > 0);
  }, []);
  
  // Improved threshold detection for more responsive hiding/showing
  useMotionValueEvent(scrollY, "change", (latest) => {
    if (latest > window.innerHeight * 0.25) {
      setShowIndicator(false);
    } else if (showAboveFold && latest < window.innerHeight * 0.05) {
      setShowIndicator(true);
    }
  });
  
  // Smooth scroll to target section
  const scrollToTarget = useCallback(() => {
    const targetElement = document.getElementById(targetId);
    
    if (targetElement) {
      // Get target position with offset
      const targetPosition = targetElement.getBoundingClientRect().top + window.pageYOffset - offset;
      
      if (!isTouchDevice) {
        // Standard smooth scroll for non-touch devices
        window.scrollTo({
          top: targetPosition,
          behavior: 'smooth'
        });
      } else {
        // Custom smooth scroll for touch devices
        const startPosition = window.pageYOffset;
        const distance = targetPosition - startPosition;
        const duration = 600;
        let start: number | null = null;
        
        // Improved easing function for smoother movement
        const easeOutCubic = (t: number): number => {
          return 1 - Math.pow(1 - t, 3);
        };
        
        const animateScroll = (timestamp: number) => {
          if (!start) start = timestamp;
          const elapsed = timestamp - start;
          const progress = Math.min(elapsed / duration, 1);
          const eased = easeOutCubic(progress);
          
          window.scrollTo(0, startPosition + distance * eased);
          
          if (progress < 1) {
            window.requestAnimationFrame(animateScroll);
          }
        };
        
        window.requestAnimationFrame(animateScroll);
      }
    }
  }, [targetId, offset, isTouchDevice]);
  
  return (
    <AnimatePresence mode="wait">
      {showIndicator && (
        <ScrollContainer
          key="scroll-indicator"
          initial="hidden"
          animate="visible"
          exit="hidden"
          variants={containerVariants}
          onClick={scrollToTarget}
          style={{ maxWidth: '90vw' }}
        >
          <ScrollText>Explore My Work</ScrollText>
          <ScrollArrow variants={arrowVariants} animate="animate">
            <svg width="24" height="14" viewBox="0 0 24 14" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M12 14L0 2.32804L2.32804 0L12 9.67196L21.672 0L24 2.32804L12 14Z" fill="currentColor"/>
            </svg>
          </ScrollArrow>
        </ScrollContainer>
      )}
    </AnimatePresence>
  );
};

export default ScrollIndicator;