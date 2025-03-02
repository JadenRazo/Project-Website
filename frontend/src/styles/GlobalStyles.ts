// src/styles/GlobalStyles.ts
import { createGlobalStyle, DefaultTheme } from 'styled-components';
import { Theme } from './theme.types';

export const GlobalStyles = createGlobalStyle<{ theme: Theme }>`
  :root {
    /* Colors */
    --background: ${({ theme }) => theme.colors.background};
    --background-alt: ${({ theme }) => theme.colors.backgroundAlt};
    --text: ${({ theme }) => theme.colors.text};
    --primary: ${({ theme }) => theme.colors.primary};
    --primary-light: ${({ theme }) => theme.colors.primaryLight};
    --primary-hover: ${({ theme }) => theme.colors.primaryHover};
    --secondary: ${({ theme }) => theme.colors.secondary};
    --accent: ${({ theme }) => theme.colors.accent};
    --surface-light: ${({ theme }) => theme.colors.surfaceLight};
    --surface-medium: ${({ theme }) => theme.colors.surfaceMedium};
    --surface-dark: ${({ theme }) => theme.colors.surfaceDark};
    --error: ${({ theme }) => theme.colors.error};
    --success: ${({ theme }) => theme.colors.success};
    --warning: ${({ theme }) => theme.colors.warning};

    /* Effects */
    --gradient: ${({ theme }) => theme.effects.gradient};
    --glass-effect: ${({ theme }) => theme.effects.glassEffect};
    --shadow: ${({ theme }) => theme.effects.shadow};

    /* Layout */
    --border-radius: 8px;
    --nav-height: 100px;
    --transition: all 0.25s cubic-bezier(0.645, 0.045, 0.355, 1);

    /* Typography */
    --font-mono: ${({ theme }) => theme.fonts.mono};
    --font-primary: ${({ theme }) => theme.fonts.primary};
    --font-sans: ${({ theme }) => theme.fonts.sans};
  }

  * {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
  }

  html {
    scroll-behavior: smooth;
    font-size: 16px;
    -webkit-text-size-adjust: 100%;
  }

  body {
    background-color: var(--background);
    color: var(--text);
    font-family: var(--font-sans);
    line-height: 1.5;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
  }

  .app {
    min-height: 100vh;
    display: flex;
    flex-direction: column;
  }

  .app-content {
    position: relative;
    z-index: 1;
    flex: 1;
    width: 100%;
  }

  .container {
    position: relative;
    z-index: 2;
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 2rem;
    width: 100%;
  }

  /* Form Elements */
  button, 
  input, 
  textarea {
    font-family: var(--font-primary);
  }

  /* Accessibility */
  @media (prefers-reduced-motion: reduce) {
    * {
      animation-duration: 0.01ms !important;
      animation-iteration-count: 1 !important;
      transition-duration: 0.01ms !important;
      scroll-behavior: auto !important;
    }
  }

  /* Selection */
  ::selection {
    background-color: var(--primary-light);
    color: var(--text);
  }

  /* Scrollbar */
  ::-webkit-scrollbar {
    width: 8px;
  }

  ::-webkit-scrollbar-track {
    background: var(--background-alt);
  }

  ::-webkit-scrollbar-thumb {
    background: var(--primary-light);
    border-radius: 4px;
  }

  ::-webkit-scrollbar-thumb:hover {
    background: var(--primary);
  }
`;
