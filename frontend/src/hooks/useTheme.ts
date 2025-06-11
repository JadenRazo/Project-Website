import { useThemeStore } from '../stores';

export const useTheme = () => {
  const {
    theme,
    themeMode,
    toggleTheme,
    setThemeMode,
  } = useThemeStore();

  return {
    theme,
    themeMode,
    toggleTheme,
    setThemeMode,
  };
};