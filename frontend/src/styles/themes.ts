// src/styles/themes.ts
import { Theme } from './theme.types';

const baseTheme = {
  fonts: {
    primary: "'SF Pro Display', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif",
    mono: "'SF Mono', 'Fira Code', 'Fira Mono', monospace",
  },
} as const;

export const lightTheme: Theme = {
  colors: {
    background: '#FFFFFF',
    backgroundAlt: '#F8F9FA',
    text: '#2C3E50',
    primary: '#007AFF',
    primaryLight: 'rgba(0, 122, 255, 0.1)',
    primaryHover: '#0056b3',
    secondary: '#6c757d',
    accent: '#B24BF3',
    surfaceLight: '#FFFFFF',
    surfaceMedium: '#F8F9FA',
    surfaceDark: '#E9ECEF',
    error: '#dc3545',
    success: '#28a745',
    warning: '#ffc107',
  },
  fonts: {
    primary: '"Inter", -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif',
    mono: '"SF Mono", "Fira Code", "Fira Mono", monospace',
    sans: '"Inter", -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif',
  },
  effects: {
    gradient: 'linear-gradient(120deg, var(--primary), var(--accent))',
    glassEffect: 'backdrop-filter: blur(10px) saturate(180%)',
    shadow: '0 4px 6px rgba(0, 0, 0, 0.1)',
  },
};

export const darkTheme: Theme = {
  colors: {
    background: '#0F172A',
    backgroundAlt: '#1E293B',
    text: '#F1F5F9',
    primary: '#60A5FA',
    primaryLight: 'rgba(96, 165, 250, 0.1)',
    primaryHover: '#3B82F6',
    secondary: '#94A3B8',
    accent: '#C084FC',
    surfaceLight: '#1E293B',
    surfaceMedium: '#334155',
    surfaceDark: '#475569',
    error: '#EF4444',
    success: '#22C55E',
    warning: '#F59E0B',
  },
  fonts: {
    primary: '"Inter", -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif',
    mono: '"SF Mono", "Fira Code", "Fira Mono", monospace',
    sans: '"Inter", -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif',
  },
  effects: {
    gradient: 'linear-gradient(120deg, var(--primary), var(--accent))',
    glassEffect: 'backdrop-filter: blur(10px) saturate(180%)',
    shadow: '0 4px 6px rgba(0, 0, 0, 0.3)',
  },
};


export const themes = {
  light: lightTheme,
  dark: darkTheme,
} as const;

// Type guard for theme checking
export const isTheme = (theme: unknown): theme is Theme => {
  if (!theme || typeof theme !== 'object') return false;
  
  const requiredKeys = ['colors', 'fonts', 'effects'];
  return requiredKeys.every(key => key in theme);
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

export const getFontValue = (theme: Theme, fontKey: keyof Theme['fonts']): string => {
  const font = theme.fonts[fontKey];
  if (!font) {
    throw new Error(`Font "${String(fontKey)}" not found in theme`);
  }
  return font;
};

export const getEffectValue = (theme: Theme, effectKey: keyof Theme['effects']): string => {
  const effect = theme.effects[effectKey];
  if (!effect) {
    throw new Error(`Effect "${String(effectKey)}" not found in theme`);
  }
  return effect;
};
