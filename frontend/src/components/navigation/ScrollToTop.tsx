import { useEffect, useRef } from 'react';
import { useLocation, useNavigationType } from 'react-router-dom';
import { useLenis } from '../../providers/LenisProvider';

const ScrollToTop = () => {
  const { pathname, hash, state } = useLocation();
  const navigationType = useNavigationType();
  const prevPathRef = useRef(pathname);
  const { lenis } = useLenis();

  useEffect(() => {
    const locationState = state as { preventScroll?: boolean; fromFooter?: boolean };

    if (locationState?.preventScroll) {
      return;
    }

    if (hash) {
      const elementId = hash.substring(1);
      const element = document.getElementById(elementId);
      if (element) {
        if (lenis) {
          lenis.scrollTo(element, { offset: 0 });
        } else {
          element.scrollIntoView({ behavior: 'smooth' });
        }
      }
    } else {
      if (prevPathRef.current !== pathname || locationState?.fromFooter) {
        if (lenis) {
          lenis.scrollTo(0, { immediate: true });
        } else {
          window.scrollTo(0, 0);
        }
      }
    }

    prevPathRef.current = pathname;
  }, [pathname, hash, state, navigationType, lenis]);

  return null;
};

export default ScrollToTop;
