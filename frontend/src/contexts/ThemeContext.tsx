import React, { createContext, useState, useContext, useCallback, useMemo } from 'react';
import { ThemeProvider as StyledThemeProvider } from 'styled-components';
import { themes } from '../styles/themes';
import { Theme, ThemeMode } from '../styles/theme.types';

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

export const ThemeProvider: React.FC<ThemeProviderProps> = ({
  children,
  defaultTheme = 'dark'
}) => {
  const [themeMode, setThemeMode] = useState<ThemeMode>(() => {
    // Try to get saved theme from localStorage
    const savedTheme = localStorage.getItem('theme') as ThemeMode;
    return savedTheme || defaultTheme;
  });

  const toggleTheme = useCallback(() => {
    setThemeMode((current) => {
      const newTheme = current === 'light' ? 'dark' : 'light';
      localStorage.setItem('theme', newTheme);
      return newTheme;
    });
  }, []);

  const handleSetThemeMode = useCallback((mode: ThemeMode) => {
    localStorage.setItem('theme', mode);
    setThemeMode(mode);
  }, []);

  const currentTheme = themes[themeMode];

  const contextValue = useMemo(
    () => ({
      theme: currentTheme,
      themeMode,
      toggleTheme,
      setThemeMode: handleSetThemeMode,
    }),
    [currentTheme, themeMode, toggleTheme, handleSetThemeMode]
  );

  return (
    <ThemeContext.Provider value={contextValue}>
      <StyledThemeProvider theme={currentTheme}>
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
