import React, { useState, useEffect } from 'react';
import styled from 'styled-components';
import { motion, AnimatePresence } from 'framer-motion';
import { useScrollTo } from '../../hooks/useScrollTo';
import { useOptimizedScrollHandler } from '../../hooks/useOptimizedScrollHandler';

const ButtonContainer = styled(motion.button)`
  position: fixed;
  bottom: 2rem;
  right: 2rem;
  width: 48px;
  height: 48px;
  border-radius: 50%;
  background: ${({ theme }) => theme.colors.primary};
  color: ${({ theme }) => theme.colors.surface};
  border: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  transition: all ${({ theme }) => theme.transitions.fast};
  z-index: 100;
  
  &:hover {
    background: ${({ theme }) => theme.colors.primaryHover};
    transform: translateY(-2px);
    box-shadow: 0 6px 16px rgba(0, 0, 0, 0.2);
  }
  
  &:active {
    transform: translateY(0);
  }
  
  svg {
    width: 24px;
    height: 24px;
    stroke-width: 2.5;
  }
  
  @media (max-width: ${({ theme }) => theme.breakpoints.mobile}) {
    bottom: 1.5rem;
    right: 1.5rem;
    width: 44px;
    height: 44px;
    
    svg {
      width: 20px;
      height: 20px;
    }
  }
`;

const ScrollToTopButton: React.FC = () => {
  const [isVisible, setIsVisible] = useState(false);
  const { scrollToTop } = useScrollTo();
  
  // Show button when scrolled down 300px
  const handleScroll = () => {
    const scrolled = window.scrollY > 300;
    setIsVisible(scrolled);
  };
  
  // Use optimized scroll handler for better performance
  useOptimizedScrollHandler(handleScroll);
  
  // Check initial scroll position
  useEffect(() => {
    handleScroll();
  }, []);
  
  const handleClick = () => {
    if (process.env.NODE_ENV === 'development') {
      console.log('ScrollToTopButton clicked!');
      console.log('Current scroll position:', window.pageYOffset);
    }
    scrollToTop({ behavior: 'smooth' });
  };
  
  return (
    <AnimatePresence>
      {isVisible && (
        <ButtonContainer
          onClick={handleClick}
          initial={{ opacity: 0, scale: 0.8 }}
          animate={{ opacity: 1, scale: 1 }}
          exit={{ opacity: 0, scale: 0.8 }}
          transition={{ duration: 0.2 }}
          whileHover={{ scale: 1.1 }}
          whileTap={{ scale: 0.95 }}
          aria-label="Scroll to top"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              d="M5 10l7-7m0 0l7 7m-7-7v18"
            />
          </svg>
        </ButtonContainer>
      )}
    </AnimatePresence>
  );
};

export default ScrollToTopButton;