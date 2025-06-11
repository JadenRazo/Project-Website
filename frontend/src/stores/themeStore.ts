import { create } from 'zustand';
import { devtools, persist } from 'zustand/middleware';
import { themes } from '../styles/themes';
import type { Theme, ThemeMode, ThemeState } from './types';

interface ThemeActions {
  toggleTheme: () => void;
  setThemeMode: (mode: ThemeMode) => void;
  initializeTheme: () => void;
}

type ThemeStore = ThemeState & ThemeActions;

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

export const useThemeStore = create<ThemeStore>()(
  devtools(
    persist(
      (set, get) => ({
        theme: themes.dark,
        themeMode: 'dark',

        toggleTheme: () => {
          const currentMode = get().themeMode;
          const newMode = currentMode === 'light' ? 'dark' : 'light';
          const newTheme = themes[newMode];

          localStorage.setItem('theme', newMode);
          
          set({
            themeMode: newMode,
            theme: newTheme,
          }, false, 'toggleTheme');

          updateDocumentTheme(newMode, newTheme);
        },

        setThemeMode: (mode: ThemeMode) => {
          const newTheme = themes[mode];
          
          localStorage.setItem('theme', mode);
          
          set({
            themeMode: mode,
            theme: newTheme,
          }, false, 'setThemeMode');

          updateDocumentTheme(mode, newTheme);
        },

        initializeTheme: () => {
          const preferredMode = getPreferredTheme();
          const preferredTheme = themes[preferredMode];

          set({
            themeMode: preferredMode,
            theme: preferredTheme,
          }, false, 'initializeTheme');

          updateDocumentTheme(preferredMode, preferredTheme);
          setupMediaQueryListener();
        },
      }),
      {
        name: 'theme-storage',
        partialize: (state) => ({ 
          themeMode: state.themeMode 
        }),
      }
    ),
    { name: 'theme-store' }
  )
);

const updateDocumentTheme = (mode: ThemeMode, theme: Theme) => {
  if (typeof document !== 'undefined') {
    document.documentElement.dataset.theme = mode;
    document.body.dataset.theme = mode;
    
    const metaThemeColor = document.querySelector('meta[name="theme-color"]');
    if (metaThemeColor) {
      metaThemeColor.setAttribute('content', theme.colors.background);
    }
  }
};

const setupMediaQueryListener = () => {
  if (typeof window !== 'undefined') {
    const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');
    
    const handleChange = (e: MediaQueryListEvent) => {
      const savedTheme = localStorage.getItem('theme') as ThemeMode;
      
      if (!savedTheme) {
        const newMode = e.matches ? 'dark' : 'light';
        useThemeStore.getState().setThemeMode(newMode);
      }
    };
    
    mediaQuery.addEventListener('change', handleChange);
    
    return () => {
      mediaQuery.removeEventListener('change', handleChange);
    };
  }
};

export const initializeTheme = () => {
  useThemeStore.getState().initializeTheme();
};