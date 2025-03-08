// src/styles/theme.types.ts
/**
 * Type definitions for the application theme system
 * Provides strong typing for theme colors and properties
 */

export type ThemeMode = 'light' | 'dark';

export interface ThemeColors {
  // Base colors
  background: string;
  backgroundAlt: string;
  backgroundSecondary: string; 
  backgroundTertiary: string;
  
  // Text colors
  text: string; 
  textSecondary: string;
  textMuted: string;
  textPrimary: string;
  
  // Brand colors
  primary: string;
  primaryLight: string;
  primaryDark: string;
  primaryHover: string;
  secondary: string;
  secondaryLight: string;
  secondaryDark: string;
  accent: string;
  
  // UI colors
  surfaceLight: string;
  surfaceMedium: string;
  surfaceDark: string;
  
  // Feedback colors
  error: string;
  success: string;
  warning: string;
  info: string;
  
  // Special effect colors
  glass: string;
  shadow: string;
  border: string;
}

export interface Theme {
  id: ThemeMode;
  colors: ThemeColors;
  
  // Additional theme properties
  borderRadius: {
    small: string;
    medium: string;
    large: string;
    pill: string;
  };
  
  shadows: {
    small: string;
    medium: string;
    large: string;
    button: string;
    text: string;
  };
  
  transitions: {
    fast: string;
    normal: string;
    slow: string;
  };
  
  spacing: {
    xxs: string;
    xs: string;
    sm: string;
    md: string;
    lg: string;
    xl: string;
    xxl: string;
  };
}

// Type for theme context value
export interface ThemeContextValue {
  theme: Theme;
  themeMode: ThemeMode;
  toggleTheme: () => void;
  setThemeMode: (mode: ThemeMode) => void;
}

// For styled-components
export interface ThemeProps {
  theme: Theme;
}
