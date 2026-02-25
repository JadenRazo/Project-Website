// src/styles/GlobalStyles.ts
import { createGlobalStyle } from 'styled-components';
import type { Theme } from './theme.types';

// Define default system fonts for fallback
const systemFonts = {
  sans: "'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif",
  primary: "'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif",
  mono: "'SF Mono', 'Fira Code', 'Fira Mono', 'Roboto Mono', monospace"
};

export const GlobalStyles = createGlobalStyle<{ theme: Theme }>`
  :root {
    --primary-color: ${({ theme }) => theme.colors.primary};
    --primary-light: ${({ theme }) => theme.colors.primaryLight};
    --primary-hover: ${({ theme }) => theme.colors.primaryHover};
    --secondary: ${({ theme }) => theme.colors.secondary};
    --secondary-light: ${({ theme }) => theme.colors.secondaryLight};
    --secondary-hover: ${({ theme }) => theme.colors.secondaryHover};
    --accent: ${({ theme }) => theme.colors.accent};
    --accent-light: ${({ theme }) => theme.colors.accentLight};
    --accent-hover: ${({ theme }) => theme.colors.accentHover};
    --background: ${({ theme }) => theme.colors.background};
    --background-alt: ${({ theme }) => theme.colors.backgroundAlt};
    --background-hover: ${({ theme }) => theme.colors.backgroundHover};
    --surface: ${({ theme }) => theme.colors.surface};
    --surface-light: ${({ theme }) => theme.colors.surfaceLight};
    --surface-hover: ${({ theme }) => theme.colors.surfaceHover};
    --surface-active: ${({ theme }) => theme.colors.surfaceActive};
    --surface-disabled: ${({ theme }) => theme.colors.surfaceDisabled};
    --text: ${({ theme }) => theme.colors.text};
    --text-hover: ${({ theme }) => theme.colors.textHover};
    --text-secondary: ${({ theme }) => theme.colors.textSecondary};
    --text-inverse: ${({ theme }) => theme.colors.textInverse};
    --text-disabled: ${({ theme }) => theme.colors.textDisabled};
    --border: ${({ theme }) => theme.colors.border};
    --border-hover: ${({ theme }) => theme.colors.borderHover};
    --border-active: ${({ theme }) => theme.colors.borderActive};
    --border-disabled: ${({ theme }) => theme.colors.borderDisabled};
    --error: ${({ theme }) => theme.colors.error};
    --error-light: ${({ theme }) => theme.colors.errorLight};
    --error-hover: ${({ theme }) => theme.colors.errorHover};
    --success: ${({ theme }) => theme.colors.success};
    --success-light: ${({ theme }) => theme.colors.successLight};
    --success-hover: ${({ theme }) => theme.colors.successHover};
    --warning: ${({ theme }) => theme.colors.warning};
    --warning-light: ${({ theme }) => theme.colors.warningLight};
    --warning-hover: ${({ theme }) => theme.colors.warningHover};
    
    --primary-rgb: 0, 120, 255; // Default RGB values for primary color
    
    /* Z-index layers */
    --z-background: -5;
    --z-content-base: 10;
    --z-scroll-indicator: 15;
    --z-navigation: 20;
    --z-modal: 50;
    --z-toast: 100;
    --z-loading-screen: 1000;
  }

  * {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
  }

  html {
    font-size: 16px;
    scroll-padding-top: 80px;

    @media (max-width: ${({ theme }) => theme.breakpoints.mobile}) {
      scroll-padding-top: 60px;
    }
  }

  body {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Oxygen',
      'Ubuntu', 'Cantarell', 'Fira Sans', 'Droid Sans', 'Helvetica Neue',
      sans-serif;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
    background-color: ${({ theme }) => theme.colors.background};
    color: ${({ theme }) => theme.colors.text};
    line-height: 1.5;
    transition: background-color ${({ theme }) => theme.transitions.normal}, 
                color ${({ theme }) => theme.transitions.normal};
    overflow-x: hidden;
    /* Force a new stacking context at the body level */
    isolation: isolate;
  }

  #root {
    position: relative;
    min-height: 100vh;
    /* Create a stacking context at the root level */
    isolation: isolate;
    z-index: 0;
  }

  /* Animation background container */
  .animation-background {
    position: fixed;
    top: 0;
    left: 0;
    width: 100vw;
    height: 100vh;
    z-index: var(--z-background);
    pointer-events: none;
    will-change: transform;
  }

  /* Content layering */
  .content-wrapper {
    position: relative;
    z-index: var(--z-content-base);
    isolation: isolate;
  }

  /* Page wrapper for smooth transitions */
  .page-transition-wrapper {
    position: relative;
    width: 100%;
    overflow-x: hidden;
    z-index: var(--z-content-base);
  }

  /* Hero section should appear above scroll indicator */
  section.hero {
    position: relative;
    z-index: var(--z-navigation);
  }

  /* Projects section should appear below scroll indicator */
  section#projects {
    position: relative;
    z-index: calc(var(--z-scroll-indicator) - 1);
  }

  /* Navigation elements */
  nav, header {
    position: relative;
    z-index: var(--z-navigation);
  }

  /* Reset for potentially problematic elements */
  canvas {
    display: block;
  }

  h1, h2, h3, h4, h5, h6 {
    font-family: ${systemFonts.primary};
    font-weight: bold;
    line-height: 1.2;
    margin-bottom: 1rem;
  }

  p {
    margin-bottom: 1rem;
  }

  a {
    color: ${({ theme }) => theme.colors.primary};
    text-decoration: none;
    transition: color ${({ theme }) => theme.transitions.fast};

    &:hover {
      color: ${({ theme }) => theme.colors.primaryHover};
    }
  }

  img {
    max-width: 100%;
    height: auto;
  }

  button {
    cursor: pointer;
    font-family: inherit;
  }

  /* Fix for iOS fixed positioning */
  @supports (-webkit-touch-callout: none) {
    .fixed-background {
      background-attachment: scroll;
    }
  }

  /* Custom scrollbar */
  ::-webkit-scrollbar {
    width: 8px;
    height: 8px;
  }

  ::-webkit-scrollbar-track {
    background: ${({ theme }) => theme.colors.background};
  }

  ::-webkit-scrollbar-thumb {
    background: ${({ theme }) => theme.colors.primaryLight};
    border-radius: 4px;

    &:hover {
      background: ${({ theme }) => theme.colors.primary};
    }
  }

  /* Utility classes */
  .visually-hidden {
    position: absolute;
    width: 1px;
    height: 1px;
    padding: 0;
    margin: -1px;
    overflow: hidden;
    clip: rect(0, 0, 0, 0);
    white-space: nowrap;
    border: 0;
  }
`;

export default GlobalStyles;
