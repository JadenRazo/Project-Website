// src/index.tsx
import React from 'react';
import { createRoot } from 'react-dom/client';
import App from './App';
import { ZIndexProvider } from './hooks/useZIndex';
import { StyleSheetManager } from 'styled-components';
import './index.css';

/**
 * Prop filtering function for styled-components
 * Prevents props like isActive, isReducedMotion, etc. from being passed to DOM elements
 */
const shouldForwardProp = (prop: string): boolean => {
  // List of props that should NOT be forwarded to DOM
  const filteredProps = [
    'isActive',
    'isReducedMotion',
    'isPowerfulDevice',
    'inView',
    'isVisible',
    'isMobile',
    'isTablet',
    'isDesktop'
  ];
  
  return !filteredProps.includes(prop);
};

/**
 * Main application entry point with integrated optimization systems
 * Wraps the entire application with necessary providers
 */
/** gang gang gang */
const AppWithProviders: React.FC = () => (
  <React.StrictMode>
    <ZIndexProvider>
      <StyleSheetManager shouldForwardProp={shouldForwardProp}>
        <App />
      </StyleSheetManager>
    </ZIndexProvider>
  </React.StrictMode>
);

// Initialize app with error boundary
const rootElement = document.getElementById('root');

if (!rootElement) {
  throw new Error('Root element not found. Failed to mount React application.');
}

const root = createRoot(rootElement);
root.render(<AppWithProviders />);
