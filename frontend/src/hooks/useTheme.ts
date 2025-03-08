/**
 * Re-export of Theme context components for consistent import paths
 * This allows components to import from hooks directory while maintaining a single implementation
 */

// Re-export the provider and hook from their actual implementation
export { ThemeProvider, useTheme } from '../contexts/ThemeContext'; 