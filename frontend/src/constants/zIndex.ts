/**
 * Z-index constants for consistent layering across the application
 * 
 * Layer groups (from back to front):
 * - Background (-10 to -1): Background elements that should always be behind content
 * - Base (0): Default z-index for standard content
 * - Low (1-10): Elements slightly above standard content
 * - Mid (11-100): Interactive elements, dropdowns, tooltips
 * - High (101-1000): Modals, overlays, notifications
 * - Top (1001+): Critical UI elements that must always be accessible
 */

export const Z_INDEX = {
  // Background layer elements
  BACKGROUND_BASE: -10,
  PARTICLE_BACKGROUND: -5,
  NETWORK_BACKGROUND: -3,
  
  // Base content layer
  CONTENT: 0,
  
  // Low priority interactive elements
  SCROLL_INDICATOR: 5,
  NAVIGATION_BACKDROP: 10,
  
  // Mid-level elements
  DROPDOWN: 50,
  TOOLTIP: 60,
  
  // High-level elements
  MODAL_BACKDROP: 100,
  MODAL: 110,
  
  // Top level elements
  LOADING_SCREEN: 1000,
  ERROR_BOUNDARY: 1001,
  TOAST_NOTIFICATION: 1010
};

/**
 * Helper function to get related z-index values
 * Used for creating stacking context relationships
 */
export const getRelatedZIndex = (baseIndex: number, offset: number): number => {
  return baseIndex + offset;
};

export default Z_INDEX; 