import React, { useState } from 'react';
import styled from 'styled-components';
import { motion, useScroll, useMotionValueEvent } from 'framer-motion';
import { useTheme } from '../../contexts/ThemeContext';

interface ScrollIndicatorProps {
  targetId?: string;
  offset?: number;
  showAboveFold?: boolean;
}

/**
 * Styled container for the scroll indicator
 * Uses either fixed or absolute positioning based on context
 */
const ScrollContainer = styled(motion.div)<{ $positionMode: 'absolute' | 'fixed' }>`
  position: ${props => props.$positionMode};
  bottom: ${props => props.$positionMode === 'absolute' ? '40px' : '15%'};
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  cursor: pointer;
  z-index: 15; /* Adjusted to be above Selected Work but below hero section */
  user-select: none;
  transition: opacity 0.3s ease, transform 0.3s ease;
  pointer-events: auto; /* Ensure it's clickable */
  filter: drop-shadow(0 0 10px rgba(0, 0, 0, 0.2)); /* Add subtle shadow to help visibility */
  
  /* Ensure visibility regardless of background content */
  &::before {
    content: '';
    position: absolute;
    top: -20px;
    left: -20px;
    right: -20px;
    bottom: -20px;
    z-index: -1;
    backdrop-filter: blur(2px);
    opacity: 0;
    pointer-events: none;
    transition: opacity 0.3s ease;
  }
  
  &:hover::before {
    opacity: 0.3;
  }
  
  @media (max-width: 768px) {
    bottom: ${props => props.$positionMode === 'absolute' ? '30px' : '12%'};
  }
  
  @media (max-width: 480px) {
    bottom: ${props => props.$positionMode === 'absolute' ? '20px' : '10%'};
  }
  
  &::after {
    content: '';
    position: absolute;
    bottom: -10px;
    left: 50%;
    transform: translateX(-50%);
    width: 40px;
    height: 2px;
    background: linear-gradient(90deg, ${props => props.theme.colors.primary}, ${props => props.theme.colors.secondary});
    border-radius: 2px;
    opacity: 0;
    transition: opacity 0.3s ease;
  }
  
  &:hover::after {
    opacity: 1;
  }
`;

const ScrollText = styled.span`
  font-size: 0.9rem;
  font-weight: 500;
  letter-spacing: 0.05em;
  text-transform: uppercase;
  color: ${props => props.theme.colors.text};
`;

const ScrollArrow = styled(motion.div)`
  color: ${props => props.theme.colors.primary};
`;

export const ScrollIndicator: React.FC<ScrollIndicatorProps> = ({ 
  targetId = 'skills',
  offset = 80,
  showAboveFold = true
}) => {
  const [showIndicator, setShowIndicator] = useState(showAboveFold);
  const { scrollY } = useScroll();
  const { theme } = useTheme();
  
  // Hide indicator once user scrolls beyond a threshold
  useMotionValueEvent(scrollY, "change", (latest) => {
    if (latest > window.innerHeight * 0.3) {
      setShowIndicator(false);
    } else if (showAboveFold) {
      setShowIndicator(true);
    }
  });
  
  // Smooth scroll to target section
  const scrollToTarget = () => {
    const targetElement = document.getElementById(targetId);
    
    if (targetElement) {
      // Calculate position, accounting for any fixed headers
      const targetPosition = targetElement.getBoundingClientRect().top + window.pageYOffset - offset;
      
      // Use smooth scrolling
      window.scrollTo({
        top: targetPosition,
        behavior: 'smooth'
      });
    }
  };
  
  return (
    <ScrollContainer
      $positionMode="absolute"
      initial={{ opacity: 0, y: 20 }}
      animate={{ 
        opacity: showIndicator ? 1 : 0, 
        y: showIndicator ? 0 : 20 
      }}
      transition={{ duration: 0.5 }}
      onClick={scrollToTarget}
    >
      <ScrollText>Explore My Work</ScrollText>
      <ScrollArrow 
        animate={{ y: [0, 8, 0] }}
        transition={{ 
          repeat: Infinity,
          duration: 1.5,
          ease: "easeInOut" 
        }}
      >
        <svg width="24" height="14" viewBox="0 0 24 14" fill="none" xmlns="http://www.w3.org/2000/svg">
          <path d="M12 14L0 2.32804L2.32804 0L12 9.67196L21.672 0L24 2.32804L12 14Z" fill="currentColor"/>
        </svg>
      </ScrollArrow>
    </ScrollContainer>
  );
};

export default ScrollIndicator;