import { local as safeLocalStorage } from '../utils/safeStorage';
import { useState, useEffect } from 'react';
import type { ThemeMode } from '../styles/theme.types';
import { themes } from '../styles/themes';

export const useThemeToggle = () => {
  const getInitialTheme = (): ThemeMode => {
    const savedTheme = safeLocalStorage.getItem('theme') as ThemeMode;
    if (savedTheme && (savedTheme === 'light' || savedTheme === 'dark')) {
      return savedTheme;
    }
    
    if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
      return 'dark';
    }
    
    return 'dark'; // Default to dark theme
  };

  const [themeMode, setThemeMode] = useState<ThemeMode>(getInitialTheme);

  useEffect(() => {
    safeLocalStorage.setItem('theme', themeMode);
    document.documentElement.setAttribute('data-theme', themeMode);
    
    // Update meta theme-color
    const metaThemeColor = document.querySelector('meta[name="theme-color"]');
    if (metaThemeColor) {
      metaThemeColor.setAttribute(
        'content',
        themes[themeMode].colors.background
      );
    }
  }, [themeMode]);

  useEffect(() => {
    const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');
    const handleChange = (e: MediaQueryListEvent) => {
      if (!safeLocalStorage.getItem('theme')) {
        setThemeMode(e.matches ? 'dark' : 'light');
      }
    };

    mediaQuery.addEventListener('change', handleChange);
    return () => mediaQuery.removeEventListener('change', handleChange);
  }, []);

  const toggleTheme = () => {
    setThemeMode(prevTheme => prevTheme === 'light' ? 'dark' : 'light');
  };

  return {
    theme: themes[themeMode],
    themeMode,
    toggleTheme
  };
};

export default useThemeToggle; 