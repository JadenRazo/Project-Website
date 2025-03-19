// src/styles/themes.ts
import type { Theme } from './theme.types';

const baseTheme = {
  breakpoints: {
    mobile: '320px',
    tablet: '768px',
    desktop: '1024px',
  },
  transitions: {
    fast: '0.2s',
    normal: '0.3s',
    slow: '0.5s',
  },
  shadows: {
    small: '0 2px 4px rgba(0,0,0,0.1)',
    medium: '0 4px 8px rgba(0,0,0,0.1)',
    large: '0 8px 16px rgba(0,0,0,0.1)',
  },
  zIndex: {
    modal: 1000,
    overlay: 900,
    dropdown: 800,
    header: 700,
  },
  borderRadius: {
    small: '4px',
    medium: '8px',
    large: '16px',
    pill: '9999px',
  },
  spacing: {
    xxs: '0.25rem',
    xs: '0.5rem',
    sm: '1rem',
    md: '1.5rem',
    lg: '2rem',
    xl: '3rem',
    xxl: '4rem',
  },
};

export const lightTheme: Theme = {
  ...baseTheme,
  colors: {
    primary: '#0078ff',
    primaryLight: '#e6f3ff',
    primaryHover: '#0056b3',
    secondary: '#6c757d',
    secondaryLight: '#f8f9fa',
    secondaryHover: '#5a6268',
    accent: '#4ecdc4',
    accentLight: '#a7ebe7',
    accentHover: '#3dbeb5',
    background: '#f8f9fa',
    backgroundAlt: '#ffffff',
    backgroundHover: '#e9ecef',
    surface: '#ffffff',
    surfaceLight: '#f8f9fa',
    surfaceHover: '#e9ecef',
    surfaceActive: '#dee2e6',
    surfaceDisabled: '#f8f9fa',
    text: '#212529',
    textHover: '#000000',
    textSecondary: '#6c757d',
    textInverse: '#ffffff',
    textDisabled: '#adb5bd',
    border: '#dee2e6',
    borderHover: '#adb5bd',
    borderActive: '#6c757d',
    borderDisabled: '#dee2e6',
    error: '#dc3545',
    errorLight: '#f8d7da',
    errorHover: '#c82333',
    success: '#28a745',
    successLight: '#d4edda',
    successHover: '#218838',
    warning: '#ffc107',
    warningLight: '#fff3cd',
    warningHover: '#e0a800',
  },
};

export const darkTheme: Theme = {
  ...baseTheme,
  colors: {
    primary: '#0078ff',
    primaryLight: '#1a1f24',
    primaryHover: '#339dff',
    secondary: '#6c757d',
    secondaryLight: '#2a2d3a',
    secondaryHover: '#868e96',
    accent: '#4ecdc4',
    accentLight: '#2a4a47',
    accentHover: '#5fe0d7',
    background: '#121212',
    backgroundAlt: '#1e1e1e',
    backgroundHover: '#2a2a2a',
    surface: '#1e1e1e',
    surfaceLight: '#2a2a2a',
    surfaceHover: '#343a40',
    surfaceActive: '#495057',
    surfaceDisabled: '#2a2a2a',
    text: '#e9ecef',
    textHover: '#ffffff',
    textSecondary: '#adb5bd',
    textInverse: '#212529',
    textDisabled: '#6c757d',
    border: '#343a40',
    borderHover: '#495057',
    borderActive: '#6c757d',
    borderDisabled: '#343a40',
    error: '#dc3545',
    errorLight: '#481a1f',
    errorHover: '#c82333',
    success: '#28a745',
    successLight: '#1a3b28',
    successHover: '#218838',
    warning: '#ffc107',
    warningLight: '#4d3b04',
    warningHover: '#e0a800',
  },
};

export const themes = {
  light: lightTheme,
  dark: darkTheme,
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
