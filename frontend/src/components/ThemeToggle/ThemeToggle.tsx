import React from 'react';
import styled from 'styled-components';
import { motion } from 'framer-motion';
import { useTheme } from '../../hooks/useTheme';

const ToggleContainer = styled.div`
  display: flex;
  align-items: center;
  gap: ${({ theme }) => theme.spacing.sm};
`;

const ToggleLabel = styled.span`
  font-size: 0.9rem;
  color: ${({ theme }) => theme.colors.textSecondary};
  user-select: none;
`;

const ToggleButton = styled(motion.button)`
  position: relative;
  width: 60px;
  height: 32px;
  background: ${({ theme }) => theme.colors.surface};
  border: 2px solid ${({ theme }) => theme.colors.border};
  border-radius: ${({ theme }) => theme.borderRadius.pill};
  cursor: pointer;
  transition: all ${({ theme }) => theme.transitions.fast};
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 4px;
  overflow: hidden;
  
  &:hover {
    border-color: ${({ theme }) => theme.colors.primary};
    background: ${({ theme }) => theme.colors.surfaceHover};
    transform: scale(1.02);
  }

  &:focus {
    outline: none;
    box-shadow: 0 0 0 2px ${({ theme }) => theme.colors.primary}40;
  }
  
  &:active {
    transform: scale(0.98);
  }
`;

const ToggleSlider = styled(motion.div)<{ $isDark: boolean }>`
  position: absolute;
  width: 24px;
  height: 24px;
  background: ${({ theme, $isDark }) => $isDark ? theme.colors.primary : theme.colors.warning};
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: ${({ theme }) => theme.colors.textInverse};
  font-size: 12px;
  box-shadow: ${({ theme }) => theme.shadows.small};
  transition: all ${({ theme }) => theme.transitions.normal};
`;

const IconWrapper = styled(motion.div)`
  width: 14px;
  height: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
`;

const SunIcon = () => (
  <IconWrapper
    initial={{ rotate: 0, scale: 0.8 }}
    animate={{ rotate: 360, scale: 1 }}
    transition={{
      rotate: { duration: 0.6, ease: "easeOut" },
      scale: { duration: 0.2, ease: "easeOut" }
    }}
  >
    <svg viewBox="0 0 24 24" fill="currentColor">
      <circle cx="12" cy="12" r="4"/>
      <path d="M12 2v2m0 16v2M4.93 4.93l1.41 1.41m11.31 11.32l1.41 1.41M2 12h2m16 0h2M6.34 17.66l-1.41 1.41M19.07 4.93l-1.41 1.41"/>
    </svg>
  </IconWrapper>
);

const MoonIcon = () => (
  <IconWrapper
    initial={{ rotate: 0, scale: 0.8 }}
    animate={{ rotate: -30, scale: 1 }}
    transition={{
      rotate: { duration: 0.4, ease: "easeOut" },
      scale: { duration: 0.2, ease: "easeOut" }
    }}
  >
    <svg viewBox="0 0 24 24" fill="currentColor">
      <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/>
    </svg>
  </IconWrapper>
);

interface ThemeToggleProps {
  showLabel?: boolean;
  size?: 'small' | 'medium' | 'large';
}

const ThemeToggle: React.FC<ThemeToggleProps> = ({ 
  showLabel = true,
  size = 'medium'
}) => {
  const { themeMode, toggleTheme } = useTheme();
  const isDark = themeMode === 'dark';

  const sliderVariants = {
    light: { x: 0 },
    dark: { x: 28 }
  };

  const buttonVariants = {
    tap: { scale: 0.98 },
    hover: { scale: 1.01 }
  };

  return (
    <ToggleContainer>
      {showLabel && (
        <ToggleLabel>
          {isDark ? 'Dark' : 'Light'} Theme
        </ToggleLabel>
      )}
      <ToggleButton
        onClick={toggleTheme}
        variants={buttonVariants}
        whileTap="tap"
        whileHover="hover"
        aria-label={`Switch to ${isDark ? 'light' : 'dark'} theme`}
        role="switch"
        aria-checked={isDark}
      >
        <ToggleSlider
          $isDark={isDark}
          variants={sliderVariants}
          animate={isDark ? 'dark' : 'light'}
          transition={{
            type: "spring",
            stiffness: 400,
            damping: 25,
            mass: 0.8
          }}
        >
          {isDark ? <MoonIcon /> : <SunIcon />}
        </ToggleSlider>
      </ToggleButton>
    </ToggleContainer>
  );
};

export default ThemeToggle;