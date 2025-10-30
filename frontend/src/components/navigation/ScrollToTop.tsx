import { useEffect, useRef } from 'react';
import { useLocation, useNavigationType } from 'react-router-dom';

const ScrollToTop = () => {
  const { pathname, hash, state } = useLocation();
  const navigationType = useNavigationType();
  const prevPathRef = useRef(pathname);

  useEffect(() => {
    const locationState = state as { preventScroll?: boolean; fromFooter?: boolean };

    if (locationState?.preventScroll) {
      return;
    }

    if (hash) {
      const elementId = hash.substring(1);
      const element = document.getElementById(elementId);
      if (element) {
        element.scrollIntoView({ behavior: 'smooth' });
      }
    } else {
      if (prevPathRef.current !== pathname || locationState?.fromFooter) {
        const pageTop = document.getElementById('page-top');
        if (pageTop) {
          pageTop.scrollIntoView({ behavior: 'smooth' });
        }
      }
    }

    prevPathRef.current = pathname;
  }, [pathname, hash, state, navigationType]);

  return null;
};

export default ScrollToTop;