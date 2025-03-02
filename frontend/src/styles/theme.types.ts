// src/styles/theme.types.ts
export type ThemeMode = 'light' | 'dark';

export interface Theme {
  colors: {
    background: string;
    backgroundAlt: string;
    text: string;
    primary: string;
    primaryLight: string;
    primaryHover: string;
    secondary: string;    
    accent: string;
    surfaceLight: string;
    surfaceMedium: string;
    surfaceDark: string;
    error: string;
    success: string;
    warning: string;
  };
  fonts: {
    primary: string;
    mono: string;
    sans: string;
  };
  effects: {
    gradient: string;
    glassEffect: string;
    shadow: string;
  };
}

// For styled-components
export interface ThemeProps {
  theme: Theme;
}
