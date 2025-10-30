export interface ScrollOptions {
  behavior?: ScrollBehavior;
  offset?: number;
}

export const SCROLL_CONFIG = {
  defaultDuration: 800,
  defaultEasing: 'ease-in-out',
  headerOffset: 80,
  mobileHeaderOffset: 60,
  defaultBehavior: 'smooth' as ScrollBehavior,
  instantBehavior: 'auto' as ScrollBehavior,
} as const;

export const SCROLL_DELAYS = {
  MODAL_OPEN: 300,
  ERROR_DISPLAY: 200,
  FORM_FOCUS: 150,
  NOTIFICATION: 100,
  ANIMATION_BUFFER: 100,
} as const;

export const SCROLL_ANIMATION = {
  DURATION: 500,
  CUSTOM_EASING: (t: number) => t < 0.5 ? 2 * t * t : -1 + (4 - 2 * t) * t,
} as const;

// Cache the scroll offset and update on resize
let cachedScrollOffset: number | null = null;
let lastWindowWidth: number | null = null;

// Mobile detection
export const isMobileDevice = (): boolean => {
  return /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent);
};

// Get viewport height accounting for mobile browser chrome
export const getViewportHeight = (): number => {
  // Use visualViewport if available (better for mobile)
  if (typeof window !== 'undefined' && window.visualViewport) {
    return window.visualViewport.height;
  }
  return window.innerHeight;
};

export const getScrollOffset = (): number => {
  const currentWidth = window.innerWidth;
  
  // Return cached value if window width hasn't changed
  if (cachedScrollOffset !== null && lastWindowWidth === currentWidth) {
    return cachedScrollOffset;
  }
  
  // Calculate and cache new value
  const isMobile = currentWidth < 768;
  let offset = isMobile ? SCROLL_CONFIG.mobileHeaderOffset : SCROLL_CONFIG.headerOffset;
  
  // Additional offset for mobile devices with browser chrome
  if (isMobile && isMobileDevice()) {
    const viewportDiff = window.innerHeight - getViewportHeight();
    if (viewportDiff > 0) {
      offset += Math.min(viewportDiff * 0.5, 20); // Add up to 20px extra for mobile chrome
    }
  }
  
  cachedScrollOffset = offset;
  lastWindowWidth = currentWidth;
  
  return cachedScrollOffset;
};

// Debounce function for resize events
let resizeTimeout: NodeJS.Timeout | null = null;
const debounce = (func: Function, wait: number) => {
  return (...args: any[]) => {
    if (resizeTimeout) clearTimeout(resizeTimeout);
    resizeTimeout = setTimeout(() => func(...args), wait);
  };
};

// Clear cache on window resize with debouncing
if (typeof window !== 'undefined') {
  const handleResize = debounce(() => {
    cachedScrollOffset = null;
  }, 250);
  
  window.addEventListener('resize', handleResize);
  
  // Cleanup function for when module is unloaded (development only)
  if (process.env.NODE_ENV === 'development' && typeof module !== 'undefined' && (module as any).hot) {
    (module as any).hot.dispose(() => {
      window.removeEventListener('resize', handleResize);
      if (resizeTimeout) clearTimeout(resizeTimeout);
      if (animationFrameId !== null) {
        cancelAnimationFrame(animationFrameId);
      }
    });
  }
}

let animationFrameId: number | null = null;

export const scrollToElement = (
  element: Element | null,
  options: ScrollOptions = {}
): Promise<void> => {
  return new Promise((resolve, reject) => {
    try {
      // Validate element
      if (!element) {
        resolve();
        return;
      }
      
      // Validate element is in the DOM
      if (!document.body.contains(element)) {
        resolve();
        return;
      }

      const { 
        behavior = SCROLL_CONFIG.defaultBehavior, 
        offset = getScrollOffset() 
      } = options;

      // Cancel any ongoing animation
      if (animationFrameId !== null) {
        cancelAnimationFrame(animationFrameId);
        animationFrameId = null;
      }

      const elementPosition = element.getBoundingClientRect().top;
      const currentScroll = window.pageYOffset || document.documentElement.scrollTop;
      const targetPosition = Math.max(0, elementPosition + currentScroll - offset);

      // For mobile touch devices, add a small delay to ensure touch events are processed
      const scrollDelay = isMobileDevice() ? 50 : 0;
      
      setTimeout(() => {
        // Use native scroll if available
        if ('scrollBehavior' in document.documentElement.style) {
          window.scrollTo({
            top: targetPosition,
            behavior: behavior
          });
          
          // Resolve after animation completes
          if (behavior === 'smooth') {
            setTimeout(resolve, SCROLL_ANIMATION.DURATION);
          } else {
            resolve();
          }
        } else if (behavior === 'smooth') {
        // Custom smooth scroll for older browsers
        const startPosition = currentScroll;
        const distance = targetPosition - startPosition;
        const startTime = performance.now();

        const animateScroll = (currentTime: number) => {
          const elapsed = currentTime - startTime;
          const progress = Math.min(elapsed / SCROLL_ANIMATION.DURATION, 1);
          
          const position = startPosition + distance * SCROLL_ANIMATION.CUSTOM_EASING(progress);
          
          window.scrollTo(0, position);
          
          if (progress < 1) {
            animationFrameId = requestAnimationFrame(animateScroll);
          } else {
            animationFrameId = null;
            resolve();
          }
        };
        
          animationFrameId = requestAnimationFrame(animateScroll);
        } else {
          // Instant scroll
          window.scrollTo(0, targetPosition);
          resolve();
        }
      }, scrollDelay);
    } catch (error) {
      if (process.env.NODE_ENV === 'development') {
        console.error('scrollToElement: Failed to scroll', error);
      }
      reject(error);
    }
  });
};

export const scrollToTop = (options: ScrollOptions = {}): Promise<void> => {
  return new Promise((resolve, reject) => {
    const { behavior = SCROLL_CONFIG.defaultBehavior } = options;
    
    try {
      // Cancel any ongoing animation
      if (animationFrameId !== null) {
        cancelAnimationFrame(animationFrameId);
        animationFrameId = null;
      }
      
      // First try native window.scrollTo with ScrollBehavior support
      if ('scrollBehavior' in document.documentElement.style) {
        window.scrollTo({
          top: 0,
          behavior: behavior
        });
        
        if (behavior === 'smooth') {
          setTimeout(resolve, SCROLL_ANIMATION.DURATION);
        } else {
          resolve();
        }
        return;
      }
    
      // Fallback to custom smooth scroll animation for older browsers
      if (behavior === 'smooth') {
        const startPosition = window.pageYOffset || document.documentElement.scrollTop || document.body.scrollTop;
        const startTime = performance.now();
        
        const animateScroll = (currentTime: number) => {
          const elapsed = currentTime - startTime;
          const progress = Math.min(elapsed / SCROLL_ANIMATION.DURATION, 1);
          
          const position = startPosition * (1 - SCROLL_ANIMATION.CUSTOM_EASING(progress));
          
          // Try all possible scroll methods
          window.scrollTo(0, position);
          document.documentElement.scrollTop = position;
          document.body.scrollTop = position;
          
          if (progress < 1) {
            animationFrameId = requestAnimationFrame(animateScroll);
          } else {
            animationFrameId = null;
            resolve();
          }
        };
        
        animationFrameId = requestAnimationFrame(animateScroll);
      } else {
        // Instant scroll - try all methods
        window.scrollTo(0, 0);
        document.documentElement.scrollTop = 0;
        document.body.scrollTop = 0;
        resolve();
      }
    } catch (error) {
      if (process.env.NODE_ENV === 'development') {
        console.error('scrollToTop: Failed to scroll', error);
      }
      // Ultimate fallback
      try {
        window.scrollTo(0, 0);
        resolve();
      } catch (e) {
        reject(e);
      }
    }
  });
};

export const scrollToPosition = (
  position: number,
  options: ScrollOptions = {}
): Promise<void> => {
  return new Promise((resolve, reject) => {
    const { behavior = SCROLL_CONFIG.defaultBehavior } = options;
    const safePosition = Math.max(0, position);
    
    try {
      // Cancel any ongoing animation
      if (animationFrameId !== null) {
        cancelAnimationFrame(animationFrameId);
        animationFrameId = null;
      }
      
      // Use native scroll if available
      if ('scrollBehavior' in document.documentElement.style) {
        window.scrollTo({
          top: safePosition,
          behavior
        });
        
        if (behavior === 'smooth') {
          setTimeout(resolve, SCROLL_ANIMATION.DURATION);
        } else {
          resolve();
        }
      } else if (behavior === 'smooth') {
        // Custom smooth scroll for older browsers
        const startPosition = window.pageYOffset || document.documentElement.scrollTop;
        const distance = safePosition - startPosition;
        const startTime = performance.now();

        const animateScroll = (currentTime: number) => {
          const elapsed = currentTime - startTime;
          const progress = Math.min(elapsed / SCROLL_ANIMATION.DURATION, 1);
          
          const currentPosition = startPosition + distance * SCROLL_ANIMATION.CUSTOM_EASING(progress);
          
          window.scrollTo(0, currentPosition);
          
          if (progress < 1) {
            animationFrameId = requestAnimationFrame(animateScroll);
          } else {
            animationFrameId = null;
            resolve();
          }
        };
        
        animationFrameId = requestAnimationFrame(animateScroll);
      } else {
        // Instant scroll
        window.scrollTo(0, safePosition);
        resolve();
      }
    } catch (error) {
      if (process.env.NODE_ENV === 'development') {
        console.error('scrollToPosition: Failed to scroll', error);
      }
      try {
        window.scrollTo(0, safePosition);
        resolve();
      } catch (e) {
        reject(e);
      }
    }
  });
};

export const shouldUseInstantScroll = (): boolean => {
  const prefersReducedMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches;
  return prefersReducedMotion;
};

export const getScrollBehavior = (): ScrollBehavior => {
  return shouldUseInstantScroll() ? SCROLL_CONFIG.instantBehavior : SCROLL_CONFIG.defaultBehavior;
};

// Helper function to handle navigation to sections
export const navigateToSection = (sectionId: string, navigate?: (to: string) => void): void => {
  const isHomePage = window.location.pathname === '/';
  
  if (isHomePage) {
    // If already on home page, just scroll to the section
    const element = document.getElementById(sectionId);
    if (element) {
      scrollToElement(element, { behavior: 'smooth' });
    }
  } else if (navigate) {
    // If on another page, navigate to home with hash
    navigate(`/#${sectionId}`);
  } else {
    // Fallback to window.location
    window.location.href = `/#${sectionId}`;
  }
};