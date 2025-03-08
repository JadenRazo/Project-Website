import React, { createContext, useContext, useMemo, ReactNode } from 'react';
import { Z_INDEX, getRelatedZIndex } from '../constants/zIndex';

interface ZIndexContextType {
  zIndex: typeof Z_INDEX;
  getRelated: (baseIndex: number, offset: number) => number;
}

const ZIndexContext = createContext<ZIndexContextType | undefined>(undefined);

interface ZIndexProviderProps {
  children: ReactNode;
  customValues?: Partial<typeof Z_INDEX>;
}

/**
 * Provider component for z-index management across the application
 * Allows for customization of z-index values during testing or for specific scenarios
 */
export const ZIndexProvider: React.FC<ZIndexProviderProps> = ({ 
  children, 
  customValues = {} 
}) => {
  const value = useMemo(() => {
    // Merge default values with any custom overrides
    const mergedZIndex = { ...Z_INDEX, ...customValues };
    
    return {
      zIndex: mergedZIndex,
      getRelated: getRelatedZIndex
    };
  }, [customValues]);

  return (
    <ZIndexContext.Provider value={value}>
      {children}
    </ZIndexContext.Provider>
  );
};

/**
 * Hook to access z-index values throughout the application
 * Provides consistent layering and helps maintain proper stacking contexts
 */
export const useZIndex = (): ZIndexContextType => {
  const context = useContext(ZIndexContext);
  
  if (!context) {
    throw new Error('useZIndex must be used within a ZIndexProvider');
  }
  
  return context;
};

export default ZIndexProvider; 