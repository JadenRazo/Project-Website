import type { Theme, ThemeMode } from '../styles/theme.types';

export type { Theme, ThemeMode };

export interface User {
  id: string;
  email: string;
  isAdmin: boolean;
  token: string;
}

export interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  authModalOpen: boolean;
  authModalMode: 'login' | 'register';
}

export interface ThemeState {
  theme: Theme;
  themeMode: ThemeMode;
}

export interface PerformanceMetrics {
  timestamp: number;
  usedJSHeapSize?: number;
  totalJSHeapSize?: number;
  jsHeapSizeLimit?: number;
  [key: string]: any;
}

export interface ApplicationState {
  effectsEnabled: boolean;
  animationsEnabled: boolean;
  backgroundEffectsEnabled: boolean;
  highQualityImagesEnabled: boolean;
  virtualizationEnabled: boolean;
}

export interface PerformanceState {
  memoryUsage: PerformanceMetrics;
  isMemoryConstrained: boolean;
  applicationState: ApplicationState;
}