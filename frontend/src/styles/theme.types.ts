// src/styles/theme.types.ts
/**
 * Type definitions for the application theme system
 * Provides strong typing for theme colors and properties
 */

export type ThemeMode = 'light' | 'dark';

export interface Theme {
  colors: {
    primary: string;
    primaryLight: string;
    primaryHover: string;
    secondary: string;
    secondaryLight: string;
    secondaryHover: string;
    accent: string;
    accentLight: string;
    accentHover: string;
    background: string;
    backgroundAlt: string;
    backgroundHover: string;
    surface: string;
    surfaceLight: string;
    surfaceHover: string;
    surfaceActive: string;
    surfaceDisabled: string;
    text: string;
    textHover: string;
    textSecondary: string;
    textInverse: string;
    textDisabled: string;
    border: string;
    borderHover: string;
    borderActive: string;
    borderDisabled: string;
    error: string;
    errorLight: string;
    errorHover: string;
    success: string;
    successLight: string;
    successHover: string;
    warning: string;
    warningLight: string;
    warningHover: string;
  };
  shadows: {
    small: string;
    medium: string;
    large: string;
  };
  breakpoints: {
    mobile: string;
    tablet: string;
    desktop: string;
  };
  transitions: {
    fast: string;
    normal: string;
    slow: string;
  };
  zIndex: {
    modal: number;
    overlay: number;
    dropdown: number;
    header: number;
  };
  borderRadius: {
    small: string;
    medium: string;
    large: string;
    pill: string;
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
