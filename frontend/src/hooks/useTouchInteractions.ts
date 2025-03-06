import { useState, useEffect, RefObject } from 'react';

/**
 * Interface describing touch interaction state
 */
export interface TouchState {
  // Basic touch state
  isTouching: boolean;
  
  // Touch coordinates
  touchStartX: number;
  touchStartY: number;
  touchCurrentX: number;
  touchCurrentY: number;
  touchEndX: number;
  touchEndY: number;
  
  // Touch metrics
  touchDuration: number;
  touchStartTime: number;
  isLongPress: boolean;
  
  // Swipe information
  swipeDirection: 'left' | 'right' | 'up' | 'down' | null;
  swipeDistance: number;
  swipeVelocity: number;
  
  // Multi-touch information
  touchCount: number;
  pinchDistance: number | null;
  pinchScale: number;
  pinchStartDistance: number | null;
  
  // Tap tracking
  tapCount: number;
  lastTapTime: number;
  doubleTapDetected: boolean;
}

/**
 * Configuration options for touch interactions
 */
export interface TouchInteractionOptions {
  longPressThreshold: number;  // Time in ms to trigger a long press
  doubleTapThreshold: number;  // Max time in ms between taps for double-tap
  swipeThreshold: number;      // Min distance in px to register a swipe
  swipeVelocityThreshold: number; // Min velocity for a "fast" swipe
  preventScrollOnSwipeX: boolean; // Prevent scrolling when swiping horizontally 
  preventScrollOnSwipeY: boolean; // Prevent scrolling when swiping vertically
}

/**
 * Default options for touch interactions
 */
const DEFAULT_OPTIONS: TouchInteractionOptions = {
  longPressThreshold: 500,
  doubleTapThreshold: 300,
  swipeThreshold: 30,
  swipeVelocityThreshold: 0.3,
  preventScrollOnSwipeX: false,
  preventScrollOnSwipeY: false
};

/**
 * Hook for enhanced touch interactions on any element
 * 
 * @param elementRef Reference to the element to track touch interactions on
 * @param options Configuration options for customizing touch behavior
 */
export const useTouchInteractions = (
  elementRef: RefObject<HTMLElement>, 
  options: Partial<TouchInteractionOptions> = {}
): TouchState => {
  // Merge default options with provided options
  const config = { ...DEFAULT_OPTIONS, ...options };
  
  // Initialize touch state
  const [touchState, setTouchState] = useState<TouchState>({
    isTouching: false,
    touchStartX: 0,
    touchStartY: 0,
    touchCurrentX: 0,
    touchCurrentY: 0,
    touchEndX: 0,
    touchEndY: 0,
    touchDuration: 0,
    touchStartTime: 0,
    isLongPress: false,
    swipeDirection: null,
    swipeDistance: 0,
    swipeVelocity: 0,
    touchCount: 0,
    pinchDistance: null,
    pinchScale: 1,
    pinchStartDistance: null,
    tapCount: 0,
    lastTapTime: 0,
    doubleTapDetected: false
  });
  
  useEffect(() => {
    const element = elementRef.current;
    if (!element || typeof window === 'undefined') return;
    
    let longPressTimer: NodeJS.Timeout | null = null;
    let touchUpdateInterval: NodeJS.Timeout | null = null;
    
    /**
     * Calculate distance between two touch points
     */
    const getTouchDistance = (touches: TouchList): number => {
      if (touches.length < 2) return 0;
      
      const dx = touches[0].clientX - touches[1].clientX;
      const dy = touches[0].clientY - touches[1].clientY;
      return Math.sqrt(dx * dx + dy * dy);
    };
    
    /**
     * Handle the start of a touch event
     */
    const handleTouchStart = (e: TouchEvent) => {
      // Clean up any existing timers
      if (longPressTimer) clearTimeout(longPressTimer);
      if (touchUpdateInterval) clearInterval(touchUpdateInterval);
      
      const touches = e.touches;
      const firstTouch = touches[0];
      const now = Date.now();
      const isMultiTouch = touches.length > 1;
      
      // Calculate multi-touch information if needed
      const pinchStartDistance = isMultiTouch ? getTouchDistance(touches) : null;
      
      // Calculate tap counts for double-tap detection
      const tapInfo = calculateTapInfo(now);
      
      // Set up new touch state
      setTouchState(prev => ({
        ...prev,
        isTouching: true,
        touchStartX: firstTouch.clientX,
        touchStartY: firstTouch.clientY,
        touchCurrentX: firstTouch.clientX,
        touchCurrentY: firstTouch.clientY,
        touchEndX: 0,
        touchEndY: 0,
        touchStartTime: now,
        touchDuration: 0,
        isLongPress: false,
        swipeDirection: null,
        swipeDistance: 0,
        swipeVelocity: 0,
        touchCount: touches.length,
        pinchDistance: pinchStartDistance,
        pinchStartDistance,
        pinchScale: 1,
        ...tapInfo
      }));
      
      // Set up a long press timer
      longPressTimer = setTimeout(() => {
        setTouchState(prev => ({
          ...prev,
          isLongPress: true
        }));
      }, config.longPressThreshold);
      
      // Update touch duration periodically
      touchUpdateInterval = setInterval(() => {
        setTouchState(prev => ({
          ...prev,
          touchDuration: Date.now() - prev.touchStartTime
        }));
      }, 100);
    };
    
    /**
     * Handle touch movement
     */
    const handleTouchMove = (e: TouchEvent) => {
      if (!touchState.isTouching) return;
      
      const touches = e.touches;
      const firstTouch = touches[0];
      const isMultiTouch = touches.length > 1;
      
      // Calculate swipe metrics
      const deltaX = firstTouch.clientX - touchState.touchStartX;
      const deltaY = firstTouch.clientY - touchState.touchStartY;
      const distance = Math.sqrt(deltaX * deltaX + deltaY * deltaY);
      const timeDelta = Date.now() - touchState.touchStartTime || 1; // Avoid division by zero
      const velocity = distance / timeDelta;
      
      // Calculate pinch distance if multi-touch
      const pinchDistance = isMultiTouch ? getTouchDistance(touches) : null;
      const pinchScale = pinchDistance && touchState.pinchStartDistance 
        ? pinchDistance / touchState.pinchStartDistance 
        : 1;
      
      // Determine swipe direction if distance threshold is met
      let swipeDirection = touchState.swipeDirection;
      if (distance > config.swipeThreshold) {
        const absX = Math.abs(deltaX);
        const absY = Math.abs(deltaY);
        
        if (absX > absY) {
          swipeDirection = deltaX > 0 ? 'right' : 'left';
          
          // Prevent default to disable scrolling if configured
          if (config.preventScrollOnSwipeX) {
            e.preventDefault();
          }
        } else {
          swipeDirection = deltaY > 0 ? 'down' : 'up';
          
          // Prevent default to disable scrolling if configured
          if (config.preventScrollOnSwipeY) {
            e.preventDefault();
          }
        }
      }
      
      // If we've moved more than the threshold, cancel the long press
      if (distance > 10 && longPressTimer) {
        clearTimeout(longPressTimer);
        longPressTimer = null;
      }
      
      // Update touch state
      setTouchState(prev => ({
        ...prev,
        touchCurrentX: firstTouch.clientX,
        touchCurrentY: firstTouch.clientY,
        touchCount: touches.length,
        swipeDirection,
        swipeDistance: distance,
        swipeVelocity: velocity,
        pinchDistance,
        pinchScale
      }));
    };
    
    /**
     * Handle touch end event
     */
    const handleTouchEnd = (e: TouchEvent) => {
      // Clean up timers
      if (longPressTimer) {
        clearTimeout(longPressTimer);
        longPressTimer = null;
      }
      
      if (touchUpdateInterval) {
        clearInterval(touchUpdateInterval);
        touchUpdateInterval = null;
      }
      
      // Calculate final touch metrics
      const endTime = Date.now();
      const touchDuration = endTime - touchState.touchStartTime;
      
      // Capture remaining touches if any
      const finalTouch = e.changedTouches[0];
      
      // Update final touch state
      setTouchState(prev => ({
        ...prev,
        isTouching: false,
        touchEndX: finalTouch.clientX,
        touchEndY: finalTouch.clientY,
        touchDuration,
        touchCount: e.touches.length
      }));
    };
    
    /**
     * Calculate tap-related information
     */
    const calculateTapInfo = (currentTime: number) => {
      // Check if this could be considered part of a double tap
      const timeSinceLastTap = currentTime - touchState.lastTapTime;
      const isDoubleTapCandidate = timeSinceLastTap < config.doubleTapThreshold;
      
      // Increment tap count if this could be a double tap, otherwise reset to 1
      const newTapCount = isDoubleTapCandidate ? touchState.tapCount + 1 : 1;
      
      // Detect double tap
      const doubleTapDetected = newTapCount >= 2;
      
      return {
        tapCount: newTapCount,
        lastTapTime: currentTime,
        doubleTapDetected
      };
    };
    
    // Add event listeners with passive option for better scrolling performance
    element.addEventListener('touchstart', handleTouchStart, { passive: !config.preventScrollOnSwipeY });
    element.addEventListener('touchmove', handleTouchMove, { passive: !config.preventScrollOnSwipeY });
    element.addEventListener('touchend', handleTouchEnd);
    element.addEventListener('touchcancel', handleTouchEnd);
    
    // Clean up
    return () => {
      if (longPressTimer) clearTimeout(longPressTimer);
      if (touchUpdateInterval) clearInterval(touchUpdateInterval);
      
      element.removeEventListener('touchstart', handleTouchStart);
      element.removeEventListener('touchmove', handleTouchMove);
      element.removeEventListener('touchend', handleTouchEnd);
      element.removeEventListener('touchcancel', handleTouchEnd);
    };
  }, [
    elementRef, 
    config.longPressThreshold, 
    config.swipeThreshold, 
    config.doubleTapThreshold,
    config.preventScrollOnSwipeX,
    config.preventScrollOnSwipeY,
    touchState.touchStartX,
    touchState.touchStartY,
    touchState.touchStartTime,
    touchState.pinchStartDistance,
    touchState.isTouching,
    touchState.lastTapTime,
    touchState.tapCount,
    touchState.swipeDirection
  ]);
  
  return touchState;
};

export default useTouchInteractions;
