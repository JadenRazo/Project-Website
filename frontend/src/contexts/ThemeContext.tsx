import React, { createContext, useContext, useEffect, useState, useMemo } from 'react';
import { ThemeProvider as StyledThemeProvider } from 'styled-components';
import { themes } from '../styles/themes';
import type { Theme, ThemeMode } from '../styles/theme.types';

interface ThemeContextValue {
  theme: Theme;
  themeMode: ThemeMode;
  toggleTheme: () => void;
  setThemeMode: (mode: ThemeMode) => void;
}

const ThemeContext = createContext<ThemeContextValue | undefined>(undefined);

interface ThemeProviderProps {
  children: React.ReactNode;
  defaultTheme?: ThemeMode;
}

const getPreferredTheme = (): ThemeMode => {
  if (typeof window !== 'undefined') {
    const savedTheme = localStorage.getItem('theme') as ThemeMode;
    
    if (savedTheme === 'dark' || savedTheme === 'light') {
      return savedTheme;
    }
    
    const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
    return prefersDark ? 'dark' : 'light';
  }
  
  return 'dark';
};

export const ThemeProvider: React.FC<ThemeProviderProps> = ({ 
  children, 
  defaultTheme = 'dark' 
}) => {
  const [themeMode, setThemeMode] = useState<ThemeMode>(getPreferredTheme);
  
  const theme = useMemo(() => themes[themeMode], [themeMode]);
  
  const toggleTheme = () => {
    setThemeMode(prev => prev === 'light' ? 'dark' : 'light');
  };
  
  const handleSetThemeMode = (mode: ThemeMode) => {
    localStorage.setItem('theme', mode);
    setThemeMode(mode);
  };
  
  useEffect(() => {
    const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');
    
    const handleChange = (e: MediaQueryListEvent) => {
      const savedTheme = localStorage.getItem('theme') as ThemeMode;
      
      if (!savedTheme) {
        setThemeMode(e.matches ? 'dark' : 'light');
      }
    };
    
    mediaQuery.addEventListener('change', handleChange);
    
    document.documentElement.dataset.theme = themeMode;
    
    return () => {
      mediaQuery.removeEventListener('change', handleChange);
    };
  }, []);
  
  useEffect(() => {
    document.body.dataset.theme = themeMode;
    
    const metaThemeColor = document.querySelector('meta[name="theme-color"]');
    if (metaThemeColor) {
      metaThemeColor.setAttribute(
        'content', 
        themeMode === 'dark' ? themes.dark.colors.background : themes.light.colors.background
      );
    }
  }, [themeMode]);
  
  const contextValue = useMemo<ThemeContextValue>(() => ({
    theme,
    themeMode,
    toggleTheme,
    setThemeMode: handleSetThemeMode,
  }), [theme, themeMode]);
  
  return (
    <ThemeContext.Provider value={contextValue}>
      <StyledThemeProvider theme={theme}>
        {children}
      </StyledThemeProvider>
    </ThemeContext.Provider>
  );
};

export const useTheme = (): ThemeContextValue => {
  const context = useContext(ThemeContext);
  
  if (!context) {
    throw new Error('useTheme must be used within a ThemeProvider');
  }
  
  return context;
};

export default ThemeContext;
