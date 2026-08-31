import { useEffect, useRef } from 'react';
import { useLocation, useNavigationType } from 'react-router-dom';
import { useScroll } from '../../providers/ScrollProvider';

const ScrollToTop = () => {
  const { pathname, hash, state } = useLocation();
  const navigationType = useNavigationType();
  const prevPathRef = useRef(pathname);
  const { scrollTo } = useScroll();

  useEffect(() => {
    const locationState = state as { preventScroll?: boolean; fromFooter?: boolean };

    if (locationState?.preventScroll) {
      return;
    }

    if (hash) {
      const elementId = hash.substring(1);
      const element = document.getElementById(elementId);
      if (element) {
        scrollTo(element);
      }
    } else {
      if (prevPathRef.current !== pathname || locationState?.fromFooter) {
        scrollTo(0, { immediate: true });
      }
    }

    prevPathRef.current = pathname;
  }, [pathname, hash, state, navigationType, scrollTo]);

  return null;
};

export default ScrollToTop;
