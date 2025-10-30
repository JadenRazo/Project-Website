import { useRef, useEffect, useState } from 'react';

interface UseLazyLoadOptions {
  threshold?: number;
  rootMargin?: string;
  placeholder?: string;
}

interface UseLazyLoadReturn {
  ref: React.RefObject<HTMLImageElement>;
  isLoaded: boolean;
  isInView: boolean;
}

export const useLazyLoad = (
  src: string,
  options: UseLazyLoadOptions = {}
): UseLazyLoadReturn => {
  const {
    threshold = 0.01,
    rootMargin = '50px',
    placeholder = 'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMSIgaGVpZ2h0PSIxIiB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciPjxyZWN0IHdpZHRoPSIxMDAlIiBoZWlnaHQ9IjEwMCUiIGZpbGw9IiNlZWUiLz48L3N2Zz4=',
  } = options;
  
  const ref = useRef<HTMLImageElement>(null);
  const [isInView, setIsInView] = useState(false);
  const [isLoaded, setIsLoaded] = useState(false);
  
  useEffect(() => {
    const element = ref.current;
    if (!element) return;
    
    // Set placeholder
    element.src = placeholder;
    
    const observer = new IntersectionObserver(
      ([entry]) => {
        if (entry.isIntersecting) {
          setIsInView(true);
          
          // Create a new image to preload
          const img = new Image();
          img.src = src;
          
          img.onload = () => {
            if (element.src !== src) {
              element.src = src;
              setIsLoaded(true);
            }
          };
          
          img.onerror = () => {
            console.error(`Failed to load image: ${src}`);
          };
          
          // Unobserve after loading starts
          observer.unobserve(element);
        }
      },
      {
        threshold,
        rootMargin,
      }
    );
    
    observer.observe(element);
    
    return () => {
      observer.disconnect();
    };
  }, [src, threshold, rootMargin, placeholder]);
  
  return { ref, isLoaded, isInView };
};

// Hook for lazy loading background images
export const useLazyBackgroundImage = (
  url: string,
  elementRef: React.RefObject<HTMLElement>,
  options: UseLazyLoadOptions = {}
): boolean => {
  const {
    threshold = 0.01,
    rootMargin = '50px',
  } = options;
  
  const [isLoaded, setIsLoaded] = useState(false);
  
  useEffect(() => {
    const element = elementRef.current;
    if (!element || !url) return;
    
    const observer = new IntersectionObserver(
      ([entry]) => {
        if (entry.isIntersecting) {
          // Preload the image
          const img = new Image();
          img.src = url;
          
          img.onload = () => {
            element.style.backgroundImage = `url(${url})`;
            setIsLoaded(true);
          };
          
          img.onerror = () => {
            console.error(`Failed to load background image: ${url}`);
          };
          
          observer.unobserve(element);
        }
      },
      {
        threshold,
        rootMargin,
      }
    );
    
    observer.observe(element);
    
    return () => {
      observer.disconnect();
    };
  }, [url, elementRef, threshold, rootMargin]);
  
  return isLoaded;
};