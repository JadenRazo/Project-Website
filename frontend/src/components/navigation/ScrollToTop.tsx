/**
 * ScrollToTop component
 * 
 * This component handles scrolling to the top of the page in multiple scenarios:
 * 1. When the route changes
 * 2. When the page refreshes
 * 3. When an internal link is clicked
 * 
 * It uses multiple scroll methods for maximum browser compatibility
 */
import { useEffect, useRef } from 'react';
import { useLocation, useNavigationType, useNavigate } from 'react-router-dom';

const ScrollToTop = () => {
  const { pathname, hash } = useLocation();
  const navigationType = useNavigationType();
  const navigate = useNavigate();
  const prevPathRef = useRef(pathname);
  
  useEffect(() => {
    // Create a forceful scroll function that uses multiple methods
    const forceScrollToTop = () => {
      // Method 1: Basic window.scrollTo
      window.scrollTo(0, 0);
      
      // Method 2: Document elements with options
      const scrollOptions = { top: 0, left: 0, behavior: 'auto' as ScrollBehavior };
      document.documentElement.scrollTo(scrollOptions);
      document.body.scrollTo(scrollOptions);
      
      // Method 3: Direct property setting
      document.documentElement.scrollTop = 0;
      document.body.scrollTop = 0;
      
      console.log(`Scrolled to top on route: ${pathname}`);
    };
    
    // Handle hash links differently
    if (hash) {
      // Special handling for #projects - redirect to projects page
      if (hash === '#projects' && pathname === '/') {
        navigate('/projects', { replace: true });
        return;
      }
      
      // Small timeout to ensure DOM is ready for other hashes
      const timeoutId = setTimeout(() => {
        try {
          const element = document.getElementById(hash.substring(1));
          if (element) {
            element.scrollIntoView({ behavior: 'smooth' });
            console.log(`Scrolled to hash: ${hash}`);
          } else {
            // If hash target doesn't exist, scroll to top instead
            forceScrollToTop();
          }
        } catch (error) {
          console.error('Error scrolling to hash:', error);
          forceScrollToTop();
        }
      }, 100);
      
      // Cleanup timeout on unmount
      return () => clearTimeout(timeoutId);
    }
    
    // If no hash and path changed, scroll to top
    if (prevPathRef.current !== pathname) {
      forceScrollToTop();
      prevPathRef.current = pathname;
    }
  }, [pathname, hash, navigationType, navigate]);
  
  // Add a global click handler for all internal navigation links
  useEffect(() => {
    const handleLinkClick = (e: MouseEvent) => {
      const target = e.target as HTMLElement;
      const link = target.closest('a');
      
      // Only process internal links (same origin)
      if (link && link.origin === window.location.origin && !link.getAttribute('href')?.includes('#')) {
        // Pre-scroll before navigation happens
        window.scrollTo(0, 0);
        console.log('Pre-scrolled on link click');
      }
    };
    
    // Add and remove click listener
    document.addEventListener('click', handleLinkClick);
    return () => {
      document.removeEventListener('click', handleLinkClick);
    };
  }, []);
  
  return null;
};

export default ScrollToTop; 