import { useEffect } from 'react';
import { initializeAuth, initializeTheme, initializePerformanceMonitoring } from '../stores';

interface StoreInitializerProps {
  children: React.ReactNode;
}

export const StoreInitializer: React.FC<StoreInitializerProps> = ({ children }) => {
  useEffect(() => {
    initializeAuth();
    initializeTheme();
    initializePerformanceMonitoring();
  }, []);

  return <>{children}</>;
};