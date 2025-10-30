import React, { useEffect, useState, useRef } from 'react';
import { useLocation } from 'react-router-dom';
import styled from 'styled-components';
import { motion, AnimatePresence } from 'framer-motion';


interface PageTransitionProps {
  children: React.ReactNode;
  onContentReady?: () => void;
}

const TransitionContainer = styled(motion.div)`
  min-height: 100vh;
  width: 100%;
`;

const pageVariants = {
  initial: {
    opacity: 0,
    y: 20,
  },
  in: {
    opacity: 1,
    y: 0,
  },
  out: {
    opacity: 0,
    y: -20,
  },
};

const pageTransition = {
  type: 'tween',
  ease: 'anticipate',
  duration: 0.4,
};

export const PageTransition: React.FC<PageTransitionProps> = ({ children, onContentReady }) => {
  const location = useLocation();
  const [isContentReady, setIsContentReady] = useState(false);
  const contentRef = useRef<HTMLDivElement>(null);
  const observerRef = useRef<MutationObserver | null>(null);
  const timeoutRef = useRef<NodeJS.Timeout | null>(null);

  useEffect(() => {
    setIsContentReady(false);

    const checkContentReady = () => {
      if (!contentRef.current) return false;

      const contentHeight = contentRef.current.scrollHeight;
      const hasSubstantialContent = contentHeight > window.innerHeight * 0.5;
      const hasImages = contentRef.current.querySelectorAll('img');
      let imagesLoaded = true;

      hasImages.forEach(img => {
        if (!img.complete || img.naturalHeight === 0) {
          imagesLoaded = false;
        }
      });

      return hasSubstantialContent && imagesLoaded;
    };

    const handleContentReady = () => {
      if (checkContentReady() && !isContentReady) {
        setIsContentReady(true);
        
        
        if (onContentReady) {
          onContentReady();
        }
      }
    };

    // Set up MutationObserver to watch for content changes
    if (contentRef.current) {
      observerRef.current = new MutationObserver(() => {
        handleContentReady();
      });

      observerRef.current.observe(contentRef.current, {
        childList: true,
        subtree: true,
        attributes: false,
      });
    }

    // Also check periodically and on image load events
    const images = contentRef.current?.querySelectorAll('img') || [];
    images.forEach(img => {
      if (!img.complete) {
        img.addEventListener('load', handleContentReady);
        img.addEventListener('error', handleContentReady);
      }
    });

    // Initial check
    timeoutRef.current = setTimeout(() => {
      handleContentReady();
    }, 100);

    // Fallback timeout
    const fallbackTimeout = setTimeout(() => {
      if (!isContentReady) {
        setIsContentReady(true);
        

        
        if (onContentReady) {
          onContentReady();
        }
      }
    }, 1000);

    return () => {
      if (observerRef.current) {
        observerRef.current.disconnect();
      }
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
      clearTimeout(fallbackTimeout);
      
      images.forEach(img => {
        img.removeEventListener('load', handleContentReady);
        img.removeEventListener('error', handleContentReady);
      });
    };
  }, [location.pathname, onContentReady]);

  return (
    <AnimatePresence mode="wait">
      <TransitionContainer
        key={location.pathname}
        ref={contentRef}
        initial="initial"
        animate="in"
        exit="out"
        variants={pageVariants}
        transition={pageTransition}
        onAnimationComplete={(definition) => {
          if (definition === 'in' && !isContentReady) {
            // Additional check when animation completes
            const checkTimer = setTimeout(() => {
              if (!isContentReady) {
                setIsContentReady(true);
                
                // Emit content ready event

                
                if (onContentReady) {
                  onContentReady();
                }
              }
            }, 100);
            return () => clearTimeout(checkTimer);
          }
        }}
      >
        {children}
      </TransitionContainer>
    </AnimatePresence>
  );
};

export default PageTransition;