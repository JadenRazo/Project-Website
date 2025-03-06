import React, { useEffect, useState } from 'react';
import styled from 'styled-components';
import { motion, useAnimation } from 'framer-motion';

interface ScrollIndicatorProps {
  targetId?: string;
  offset?: number;
  showAboveFold?: boolean;
}

const ScrollContainer = styled(motion.div)<{ positionMode: 'absolute' | 'fixed' }>`
  position: ${props => props.positionMode};
  bottom: ${props => props.positionMode === 'absolute' ? '40px' : '15%'};
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  cursor: pointer;
  z-index: 1000;
  user-select: none;
  transition: opacity 0.3s ease;
  
  @media (max-width: 768px) {
    bottom: ${props => props.positionMode === 'absolute' ? '30px' : '12%'};
  }
  
  @media (max-width: 480px) {
    bottom: ${props => props.positionMode === 'absolute' ? '20px' : '10%'};
  }
`;

const ScrollText = styled(motion.span)`
  font-size: 14px;
  color: var(--primary);
  opacity: 0.8;
  font-weight: 500;
  letter-spacing: 0.5px;
  pointer-events: none;
  white-space: nowrap;
`;

const ArrowContainer = styled(motion.div)`
  width: 28px;
  height: 44px;
  border: 2px solid var(--primary);
  border-radius: 14px;
  display: flex;
  justify-content: center;
  padding-top: 12px;
  box-sizing: border-box;
  
  @media (max-width: 480px) {
    width: 24px;
    height: 38px;
    padding-top: 10px;
  }
`;

const ArrowDot = styled(motion.div)`
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background-color: var(--primary);
  
  @media (max-width: 480px) {
    width: 5px;
    height: 5px;
  }
`;

export const ScrollIndicator: React.FC<ScrollIndicatorProps> = ({ 
  targetId = 'projects',
  offset = 80,
  showAboveFold = true
}) => {
  const [isVisible, setIsVisible] = useState(true);
  const [positionMode, setPositionMode] = useState<'absolute' | 'fixed'>(
    showAboveFold ? 'absolute' : 'fixed'
  );
  const controls = useAnimation();
  
  useEffect(() => {
    // Animate on mount
    controls.start({
      y: [0, 6, 0],
      transition: {
        duration: 1.5,
        repeat: Infinity,
        repeatType: 'loop',
        ease: 'easeInOut'
      }
    });
    
    // Setup scroll listener for visibility
    const handleScroll = () => {
      // Hide when scrolled down
      const scrollPosition = window.scrollY;
      const viewportHeight = window.innerHeight;
      
      if (scrollPosition > viewportHeight * 0.3) {
        setIsVisible(false);
      } else {
        setIsVisible(true);
      }
    };
    
    window.addEventListener('scroll', handleScroll);
    return () => window.removeEventListener('scroll', handleScroll);
  }, [controls]);
  
  const scrollToContent = () => {
    const targetElement = document.getElementById(targetId);
    if (targetElement) {
      const targetPosition = targetElement.getBoundingClientRect().top;
      const offsetPosition = targetPosition + window.scrollY - offset;
      
      window.scrollTo({
        top: offsetPosition,
        behavior: 'smooth'
      });
    }
  };
  
  return (
    <ScrollContainer 
      positionMode={positionMode}
      onClick={scrollToContent}
      animate={{ opacity: isVisible ? 1 : 0 }}
      initial={{ opacity: 0 }}
      whileHover={{ scale: 1.05 }}
      whileTap={{ scale: 0.95 }}
      aria-label="Scroll to projects"
      role="button"
    >
      <ScrollText>Scroll Down</ScrollText>
      <ArrowContainer>
        <ArrowDot animate={controls} />
      </ArrowContainer>
    </ScrollContainer>
  );
};

export default ScrollIndicator;