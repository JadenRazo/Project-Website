// src/styles/themes.ts
import type { Theme, ThemeMode } from './theme.types';

export const themes: Record<ThemeMode, Theme> = {
  light: {
    id: 'light',
    colors: {
      // Base colors
      background: '#ffffff',
      backgroundAlt: '#f8f9fa',
      backgroundSecondary: '#f0f2f5',
      backgroundTertiary: '#e9ecef',
      
      // Text colors
      text: '#212529',
      textSecondary: '#495057',
      textMuted: '#6c757d',
      textPrimary: '#212529',
      
      // Brand colors
      primary: '#0066ff',
      primaryLight: '#4d94ff',
      primaryDark: '#0047b3',
      primaryHover: '#0052cc',
      secondary: '#6c757d',
      secondaryLight: '#adb5bd',
      secondaryDark: '#495057',
      accent: '#ff9500',
      
      // UI colors
      surfaceLight: '#e9ecef',
      surfaceMedium: '#dee2e6',
      surfaceDark: '#ced4da',
      
      // Feedback colors
      error: '#dc3545',
      success: '#28a745',
      warning: '#ffc107',
      info: '#17a2b8',
      
      // Special effect colors
      glass: 'rgba(255, 255, 255, 0.8)',
      shadow: 'rgba(0, 0, 0, 0.1)',
      border: '#dee2e6'
    },
    borderRadius: {
      small: '4px',
      medium: '8px',
      large: '16px',
      pill: '9999px'
    },
    shadows: {
      small: '0 2px 4px rgba(0, 0, 0, 0.05)',
      medium: '0 4px 8px rgba(0, 0, 0, 0.08)',
      large: '0 8px 16px rgba(0, 0, 0, 0.1)',
      button: '0 4px 6px rgba(0, 0, 0, 0.1)',
      text: '0 1px 2px rgba(0, 0, 0, 0.1)'
    },
    transitions: {
      fast: '150ms ease',
      normal: '300ms ease',
      slow: '500ms ease'
    },
    spacing: {
      xxs: '0.25rem',
      xs: '0.5rem',
      sm: '1rem',
      md: '1.5rem',
      lg: '2rem',
      xl: '3rem',
      xxl: '4rem'
    }
  },
  dark: {
    id: 'dark',
    colors: {
      // Base colors
      background: '#121212',
      backgroundAlt: '#1e1e1e',
      backgroundSecondary: '#252525',
      backgroundTertiary: '#2c2c2c',
      
      // Text colors
      text: '#e0e0e0',
      textSecondary: '#b0b0b0',
      textMuted: '#808080',
      textPrimary: '#ffffff',
      
      // Brand colors
      primary: '#4d94ff',
      primaryLight: '#80b3ff',
      primaryDark: '#0052cc',
      primaryHover: '#1a75ff',
      secondary: '#6c757d',
      secondaryLight: '#adb5bd',
      secondaryDark: '#495057',
      accent: '#ff9500',
      
      // UI colors
      surfaceLight: '#2c2c2c',
      surfaceMedium: '#3c3c3c',
      surfaceDark: '#4c4c4c',
      
      // Feedback colors
      error: '#f55a4e',
      success: '#5cb85c',
      warning: '#f0ad4e',
      info: '#5bc0de',
      
      // Special effect colors
      glass: 'rgba(30, 30, 30, 0.8)',
      shadow: 'rgba(0, 0, 0, 0.2)',
      border: '#333333'
    },
    borderRadius: {
      small: '4px',
      medium: '8px',
      large: '16px',
      pill: '9999px'
    },
    shadows: {
      small: '0 2px 4px rgba(0, 0, 0, 0.2)',
      medium: '0 4px 8px rgba(0, 0, 0, 0.3)',
      large: '0 8px 16px rgba(0, 0, 0, 0.4)',
      button: '0 4px 6px rgba(0, 0, 0, 0.3)',
      text: '0 1px 2px rgba(0, 0, 0, 0.5)'
    },
    transitions: {
      fast: '150ms ease',
      normal: '300ms ease',
      slow: '500ms ease'
    },
    spacing: {
      xxs: '0.25rem',
      xs: '0.5rem',
      sm: '1rem',
      md: '1.5rem',
      lg: '2rem',
      xl: '3rem',
      xxl: '4rem'
    }
  }
};

/**
 * Determines if a value is a valid theme
 */
export const isTheme = (theme: unknown): theme is Theme => {
  return (
    typeof theme === 'object' &&
    theme !== null &&
    'colors' in theme &&
    typeof theme.colors === 'object'
  );
};

// Helper function to ensure type safety when accessing theme values
export const getThemeValue = <T extends keyof Theme>(
  theme: Theme,
  property: T
): Theme[T] => {
  if (!theme[property]) {
    throw new Error(`Theme property "${String(property)}" not found`);
  }
  return theme[property];
};

export const getColorValue = (theme: Theme, colorKey: keyof Theme['colors']): string => {
  const color = theme.colors[colorKey];
  if (!color) {
    throw new Error(`Color "${String(colorKey)}" not found in theme`);
  }
  return color;
};

export const getBorderRadiusValue = (
  theme: Theme, 
  radiusKey: keyof Theme['borderRadius']
): string => {
  const radius = theme.borderRadius[radiusKey];
  if (!radius) {
    throw new Error(`Border radius "${String(radiusKey)}" not found in theme`);
  }
  return radius;
};

export const getShadowValue = (
  theme: Theme, 
  shadowKey: keyof Theme['shadows']
): string => {
  const shadow = theme.shadows[shadowKey];
  if (!shadow) {
    throw new Error(`Shadow "${String(shadowKey)}" not found in theme`);
  }
  return shadow;
};

export const getTransitionValue = (
  theme: Theme, 
  transitionKey: keyof Theme['transitions']
): string => {
  const transition = theme.transitions[transitionKey];
  if (!transition) {
    throw new Error(`Transition "${String(transitionKey)}" not found in theme`);
  }
  return transition;
};

export const getSpacingValue = (
  theme: Theme, 
  spacingKey: keyof Theme['spacing']
): string => {
  const spacing = theme.spacing[spacingKey];
  if (!spacing) {
    throw new Error(`Spacing "${String(spacingKey)}" not found in theme`);
  }
  return spacing;
};

// Utility function to combine multiple theme values
export const combineThemeValues = <T>(values: T[]): T[] => {
  return values;
};
