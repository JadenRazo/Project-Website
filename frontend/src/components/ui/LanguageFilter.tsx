import React, { memo, useMemo } from 'react';
import styled from 'styled-components';
import { motion } from 'framer-motion';
import useDeviceCapabilities from '../../hooks/useDeviceCapabilities';
import usePerformanceOptimizations from '../../hooks/usePerformanceOptimizations';

interface LanguageFilterProps {
  languages: readonly string[];
  selectedLanguage: string;
  onSelectLanguage: (language: string) => void;
}

interface FilterButtonProps {
  isActive: boolean;
  isPowerfulDevice: boolean;
}

// Optimize styled components with efficient CSS transitions
const FilterContainer = styled(motion.div)`
  display: flex;
  justify-content: center;
  flex-wrap: wrap;
  gap: 12px;
  padding: 1.5rem 0;
  max-width: 100%;
  
  @media (max-width: 768px) {
    gap: 8px;
    padding: 1rem 0;
  }
  
  @media (max-width: 480px) {
    gap: 6px;
    justify-content: flex-start;
    overflow-x: auto;
    padding: 0.75rem 0;
    -webkit-overflow-scrolling: touch;
    scroll-snap-type: x mandatory;
    
    &::-webkit-scrollbar {
      display: none;
    }
  }
`;

// Filter button with hardware-accelerated CSS transitions for better performance
const FilterButton = styled.button<FilterButtonProps>`
  padding: 10px 20px;
  border: none;
  border-radius: 8px;
  background-color: ${({ isActive, theme }) => 
    isActive ? theme.colors.primary : `${theme.colors.primary}15`};
  color: ${({ isActive, theme }) => 
    isActive ? theme.colors.backgroundAlt : theme.colors.primary};
  cursor: pointer;
  font-size: clamp(0.875rem, 1vw, 1rem);
  font-family: inherit;
  font-weight: ${({ isActive }) => isActive ? '600' : '400'};
  position: relative;
  overflow: hidden;
  transform: translateZ(0); /* Hardware acceleration */
  
  /* Use efficient CSS transitions instead of animation libraries for simple effects */
  transition: 
    background-color 0.15s ease-out,
    color 0.15s ease-out,
    transform 0.15s ease-out,
    box-shadow 0.15s ease-out;
    
  &:hover {
    background-color: ${({ isActive, theme }) => 
      isActive ? theme.colors.primary : `${theme.colors.primary}25`};
    transform: ${({ isPowerfulDevice }) => 
      isPowerfulDevice ? 'translateY(-2px)' : 'none'};
    box-shadow: ${({ isPowerfulDevice }) => 
      isPowerfulDevice ? '0 4px 8px rgba(0, 0, 0, 0.1)' : 'none'};
  }
  
  &:active {
    transform: scale(0.97);
    transition: transform 0.1s ease-out;
  }
  
  /* Active button highlight effect */
  &::after {
    content: '';
    position: absolute;
    bottom: 0;
    left: 50%;
    width: ${({ isActive }) => isActive ? '40%' : '0'};
    height: 2px;
    background-color: ${({ theme }) => theme.colors.background};
    transform: translateX(-50%);
    transition: width 0.15s ease-out;
    opacity: ${({ isActive }) => isActive ? 1 : 0};
  }
  
  @media (max-width: 768px) {
    padding: 8px 16px;
    font-size: 0.875rem;
  }
  
  @media (max-width: 480px) {
    padding: 6px 14px;
    font-size: 0.8125rem;
    flex-shrink: 0;
    scroll-snap-align: start;
  }
`;

// Optimized ripple effect component that uses efficient DOM operations
const Ripple = styled(motion.span)`
  position: absolute;
  border-radius: 50%;
  background-color: rgba(255, 255, 255, 0.3);
  transform: scale(0);
  pointer-events: none;
`;

const LanguageFilter: React.FC<LanguageFilterProps> = ({ 
  languages, 
  selectedLanguage, 
  onSelectLanguage 
}) => {
  const { deviceCapabilities } = usePerformanceOptimizations();
  const isPowerfulDevice = deviceCapabilities.deviceCategory !== 'mobile' && 
                        !deviceCapabilities.isLowPoweredDevice;
                        
  // Memoize filter containers variant to prevent re-creation
  const containerVariant = useMemo(() => ({
    initial: { opacity: 0, y: -10 },
    animate: { 
      opacity: 1, 
      y: 0,
      transition: { staggerChildren: 0.05 }
    },
    exit: { opacity: 0 }
  }), []);
  
  // Memoize button variants to prevent re-creation
  const buttonVariant = useMemo(() => ({
    initial: { opacity: 0, y: -5 },
    animate: { 
      opacity: 1, 
      y: 0,
      transition: { 
        duration: 0.2,
        ease: "easeOut" 
      } 
    }
  }), []);
  
  // Avoid creating a new array on every render
  const languageOptions = useMemo(() => languages, [languages]);
  
  return (
    <FilterContainer
      initial="initial"
      animate="animate"
      exit="exit"
      variants={containerVariant}
    >
      {languageOptions.map(language => (
        <motion.div key={language} variants={buttonVariant}>
          <FilterButton
            isActive={selectedLanguage === language}
            isPowerfulDevice={isPowerfulDevice}
            onClick={() => onSelectLanguage(language)}
            role="button"
            aria-pressed={selectedLanguage === language}
          >
            {language}
          </FilterButton>
        </motion.div>
      ))}
    </FilterContainer>
  );
};

export default memo(LanguageFilter);
