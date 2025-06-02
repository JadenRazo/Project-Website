import React, { memo, useMemo } from 'react';
import styled from 'styled-components';
import { motion } from 'framer-motion';
import usePerformanceOptimizations from '../../hooks/usePerformanceOptimizations';

interface LanguageFilterProps {
  languages: readonly string[];
  selectedLanguage: string;
  onSelectLanguage: (language: string) => void;
}

interface LanguageButtonProps {
  $isActive: boolean;
  $isPowerfulDevice: boolean;
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
const LanguageButton = styled.button<LanguageButtonProps>`
  background-color: ${({ $isActive, theme }) =>
    $isActive ? theme.colors.primary : `${theme.colors.primary}15`};
  color: ${({ $isActive, theme }) =>
    $isActive ? theme.colors.backgroundAlt : theme.colors.primary};
  border: none;
  padding: ${({ theme }) => `${theme.spacing.xs} ${theme.spacing.md}`};
  border-radius: ${({ theme }) => theme.borderRadius.pill};
  font-weight: ${({ $isActive }) => $isActive ? '600' : '400'};
  font-size: 0.9rem;
  cursor: pointer;
  transition: all ${({ theme }) => theme.transitions.normal};
  position: relative;
  overflow: hidden;

  &:hover {
    background-color: ${({ $isActive, theme }) =>
      $isActive ? theme.colors.primary : `${theme.colors.primary}25`};
    transform: ${({ $isPowerfulDevice }) =>
      $isPowerfulDevice ? 'translateY(-2px)' : 'none'};
    box-shadow: ${({ $isPowerfulDevice }) =>
      $isPowerfulDevice ? '0 4px 8px rgba(0, 0, 0, 0.1)' : 'none'};
  }
`;

const LanguageName = styled.span<{ $isActive: boolean }>`
  position: relative;
  z-index: 2;
`;

const ButtonBg = styled.div<{ $isActive: boolean }>`
  position: absolute;
  bottom: 0;
  left: 0;
  height: 2px;
  width: ${({ $isActive }) => $isActive ? '40%' : '0'};
  background: ${({ theme }) => theme.colors.backgroundAlt};
  transition: all ${({ theme }) => theme.transitions.normal};
  opacity: ${({ $isActive }) => $isActive ? 1 : 0};
  z-index: 1;
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
          <LanguageButton
            $isActive={selectedLanguage === language}
            $isPowerfulDevice={isPowerfulDevice}
            onClick={() => onSelectLanguage(language)}
            role="button"
            aria-pressed={selectedLanguage === language}
          >
            <ButtonBg $isActive={selectedLanguage === language} />
            <LanguageName $isActive={selectedLanguage === language}>
              {language}
            </LanguageName>
          </LanguageButton>
        </motion.div>
      ))}
    </FilterContainer>
  );
};

export default memo(LanguageFilter);
