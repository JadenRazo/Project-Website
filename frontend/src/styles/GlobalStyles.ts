// src/styles/GlobalStyles.ts
import { createGlobalStyle } from 'styled-components';

// Define default system fonts for fallback
const systemFonts = {
  sans: "'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif",
  primary: "'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif",
  mono: "'SF Mono', 'Fira Code', 'Fira Mono', 'Roboto Mono', monospace"
};

export const GlobalStyles = createGlobalStyle`
  :root {
    --primary: ${({ theme }) => theme.colors.primary};
    --background: ${({ theme }) => theme.colors.background};
    --text: ${({ theme }) => theme.colors.text};
    --secondary: ${({ theme }) => theme.colors.secondary};
    --accent: ${({ theme }) => theme.colors.accent};
    
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
    box-sizing: border-box;
    margin: 0;
    padding: 0;
  }

  html, body {
    overflow-x: hidden;
    width: 100%;
    height: 100%;
    position: relative;
    scroll-behavior: smooth;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
  }

  body {
    background-color: ${({ theme }) => theme.colors.background};
    color: ${({ theme }) => theme.colors.text};
    font-family: ${systemFonts.sans};
    line-height: 1.6;
    overflow-y: auto;
    margin: 0;
    padding: 0;
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
    transition: color 0.3s ease;

    &:hover {
      color: ${({ theme }) => theme.colors.secondary};
    }
  }

  img {
    max-width: 100%;
    height: auto;
  }

  button {
    cursor: pointer;
    font-family: ${systemFonts.sans};
  }

  /* Fix for iOS fixed positioning */
  @supports (-webkit-touch-callout: none) {
    .fixed-background {
      background-attachment: scroll;
    }
  }

  /* Fix for scrollbar inconsistencies */
  @media screen and (min-width: 768px) {
    html {
      scrollbar-width: thin;
      scrollbar-color: ${({ theme }) => theme.colors.secondary} ${({ theme }) => theme.colors.background};
    }

    ::-webkit-scrollbar {
      width: 8px;
    }

    ::-webkit-scrollbar-track {
      background: ${({ theme }) => theme.colors.background};
    }

    ::-webkit-scrollbar-thumb {
      background-color: ${({ theme }) => theme.colors.secondary};
      border-radius: 4px;
    }
  }
`;

export default GlobalStyles;
