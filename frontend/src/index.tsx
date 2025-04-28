// src/index.tsx
import React from 'react';
import { createRoot } from 'react-dom/client';
import App from './App';
import { ThemeProvider } from './contexts/ThemeContext';
import { ZIndexProvider } from './hooks/useZIndex';
import { MemoryManagerProvider } from './utils/MemoryManager';
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
    <MemoryManagerProvider 
      monitoringInterval={30000}
      memoryThreshold={75}
    >
      <ThemeProvider defaultTheme="dark">
        <ZIndexProvider>
          <StyleSheetManager shouldForwardProp={shouldForwardProp}>
            <App />
          </StyleSheetManager>
        </ZIndexProvider>
      </ThemeProvider>
    </MemoryManagerProvider>
  </React.StrictMode>
);

// Initialize app with error boundary
const rootElement = document.getElementById('root');

if (!rootElement) {
  throw new Error('Root element not found. Failed to mount React application.');
}

const root = createRoot(rootElement);
root.render(<AppWithProviders />);
