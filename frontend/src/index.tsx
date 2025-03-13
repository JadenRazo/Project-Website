// src/index.tsx
import React from 'react';
import { createRoot } from 'react-dom/client';
import App from './App';
import { ThemeProvider } from './contexts/ThemeContext';
import { ZIndexProvider } from './hooks/useZIndex';
import { MemoryManagerProvider } from './utils/MemoryManager';
import './index.css';

/**
 * Main application entry point with integrated optimization systems
 * Wraps the entire application with necessary providers
 */
const AppWithProviders: React.FC = () => (
  <React.StrictMode>
    <MemoryManagerProvider 
      monitoringInterval={30000}
      memoryThreshold={75}
    >
      <ThemeProvider defaultTheme="dark">
        <ZIndexProvider>
          <App />
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
