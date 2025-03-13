import React, { useState, useRef, useEffect } from 'react';
import styled from 'styled-components';
import { useSafeAsync } from '../../utils/performance';

interface OptimizedImageProps extends React.ImgHTMLAttributes<HTMLImageElement> {
  src: string;
  lowResSrc?: string;
  alt: string;
  aspectRatio?: number;
  lazyLoad?: boolean;
  blurPlaceholder?: boolean;
  priority?: boolean;
  onLoad?: () => void;
  onError?: () => void;
  backgroundColor?: string;
}

interface ImageContainerProps {
  aspectRatio?: number;
  $loaded: boolean;
  $blurPlaceholder: boolean;
  $backgroundColor?: string;
}

const ImageContainer = styled.div<ImageContainerProps>`
  position: relative;
  overflow: hidden;
  width: 100%;
  height: ${props => props.aspectRatio ? '0' : 'auto'};
  padding-bottom: ${props => props.aspectRatio ? `${(1 / props.aspectRatio) * 100}%` : '0'};
  background-color: ${props => props.$backgroundColor || 'transparent'};
  
  ${props => props.$blurPlaceholder && !props.$loaded ? `
    &::before {
      content: '';
      position: absolute;
      inset: 0;
      backdrop-filter: blur(8px);
      -webkit-backdrop-filter: blur(8px);
    }
  ` : ''}
`;

const StyledImage = styled.img<{ $visible: boolean; $loaded: boolean; $lowResLoaded: boolean }>`
  position: ${props => props.$loaded ? 'relative' : 'absolute'};
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
  opacity: ${props => (props.$visible && (props.$loaded || props.$lowResLoaded)) ? 1 : 0};
  transition: opacity 0.3s ease-in-out;
  will-change: opacity;
`;

const LowResImage = styled(StyledImage)`
  filter: blur(10px);
  transform: scale(1.1);
  z-index: 1;
  opacity: ${props => props.$visible && props.$lowResLoaded && !props.$loaded ? 0.7 : 0};
`;

const MainImage = styled(StyledImage)`
  z-index: 2;
`;

/**
 * Memory-optimized image component with progressive loading
 * Features lazy loading, aspect ratio preservation, and blur-up technique
 */
export const OptimizedImage: React.FC<OptimizedImageProps> = ({
  src,
  lowResSrc,
  alt,
  aspectRatio,
  lazyLoad = true,
  blurPlaceholder = true,
  priority = false,
  onLoad,
  onError,
  backgroundColor,
  ...props
}) => {
  const [visible, setVisible] = useState(!lazyLoad || priority);
  const [mainImageLoaded, setMainImageLoaded] = useState(false);
  const [lowResLoaded, setLowResLoaded] = useState(false);
  const imgRef = useRef<HTMLImageElement>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const { safeExecute } = useSafeAsync();
  
  // Handle intersection observer for lazy loading
  useEffect(() => {
    if (!lazyLoad || priority) return;
    
    const observer = new IntersectionObserver(
      entries => {
        const [entry] = entries;
        if (entry.isIntersecting) {
          safeExecute(() => setVisible(true));
          observer.disconnect();
        }
      },
      {
        rootMargin: '200px 0px', // Start loading when image is 200px from viewport
        threshold: 0.01
      }
    );
    
    if (containerRef.current) {
      observer.observe(containerRef.current);
    }
    
    return () => {
      observer.disconnect();
    };
  }, [lazyLoad, priority, safeExecute]);
  
  // Handle image load state
  const handleMainImageLoad = safeExecute(() => {
    setMainImageLoaded(true);
    if (onLoad) onLoad();
    
    // Free memory of low-res image once main image is loaded
    if (lowResLoaded && imgRef.current) {
      const lowResImg = imgRef.current.previousSibling as HTMLImageElement;
      if (lowResImg) {
        lowResImg.src = 'data:image/gif;base64,R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7';
      }
    }
  });
  
  const handleLowResImageLoad = safeExecute(() => {
    setLowResLoaded(true);
  });
  
  const handleError = safeExecute(() => {
    if (onError) onError();
  });
  
  return (
    <ImageContainer 
      ref={containerRef}
      aspectRatio={aspectRatio}
      $loaded={mainImageLoaded}
      $blurPlaceholder={blurPlaceholder}
      $backgroundColor={backgroundColor}
    >
      {lowResSrc && (
        <LowResImage
          src={visible ? lowResSrc : ''}
          alt=""
          $visible={visible}
          $loaded={mainImageLoaded}
          $lowResLoaded={lowResLoaded}
          onLoad={handleLowResImageLoad}
          onError={handleError}
          aria-hidden="true"
          loading="lazy"
        />
      )}
      
      <MainImage
        ref={imgRef}
        src={visible ? src : ''}
        alt={alt}
        $visible={visible}
        $loaded={mainImageLoaded}
        $lowResLoaded={lowResLoaded}
        onLoad={handleMainImageLoad}
        onError={handleError}
        loading={priority ? 'eager' : 'lazy'}
        decoding={priority ? 'sync' : 'async'}
        {...props}
      />
    </ImageContainer>
  );
};

export default OptimizedImage; 