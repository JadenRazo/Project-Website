import React from 'react';
import { createRoot } from 'react-dom/client';
import App from './App';
import { ZIndexProvider } from './hooks/useZIndex';
import { StyleSheetManager } from 'styled-components';
import './index.css';

const shouldForwardProp = (prop: string): boolean => {
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

const AppWithProviders: React.FC = () => (
  <React.StrictMode>
    <ZIndexProvider>
      <StyleSheetManager shouldForwardProp={shouldForwardProp}>
        <App />
      </StyleSheetManager>
    </ZIndexProvider>
  </React.StrictMode>
);

const rootElement = document.getElementById('root');

if (!rootElement) {
  throw new Error('Root element not found. Failed to mount React application.');
}

const root = createRoot(rootElement);
root.render(<AppWithProviders />);

const hideInitialLoader = () => {
  const loader = document.getElementById('initial-loader');
  if (loader) {
    loader.classList.add('hidden');
    setTimeout(() => {
      loader.remove();
    }, 500);
  }
};

if (document.readyState === 'complete') {
  hideInitialLoader();
} else {
  window.addEventListener('load', hideInitialLoader);
}
