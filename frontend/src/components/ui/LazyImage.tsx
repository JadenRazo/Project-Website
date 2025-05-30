import React, { useState, useRef, useCallback, useEffect } from 'react';
import styled, { keyframes } from 'styled-components';
import { useInView } from 'react-intersection-observer';

interface LazyImageProps extends React.ImgHTMLAttributes<HTMLImageElement> {
  src: string;
  alt: string;
  lowQualitySrc?: string;
  blurDataURL?: string;
  aspectRatio?: string;
  priority?: boolean;
  onLoadComplete?: () => void;
  fallbackSrc?: string;
  placeholderColor?: string;
  showLoader?: boolean;
}

const shimmerAnimation = keyframes`
  0% {
    background-position: -468px 0;
  }
  100% {
    background-position: 468px 0;
  }
`;

const ImageContainer = styled.div<{
  $aspectRatio?: string;
  $placeholderColor?: string;
}>`
  position: relative;
  width: 100%;
  height: 0;
  padding-bottom: ${({ $aspectRatio }) => $aspectRatio || '56.25%'};
  background-color: ${({ $placeholderColor, theme }) => 
    $placeholderColor || theme.colors.surface || '#f0f0f0'};
  overflow: hidden;
  border-radius: inherit;
`;

const PlaceholderDiv = styled.div<{
  $showShimmer: boolean;
  $blurDataURL?: string;
}>`
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-image: ${({ $blurDataURL }) => 
    $blurDataURL ? `url(${$blurDataURL})` : 'none'};
  background-size: cover;
  background-position: center;
  filter: ${({ $blurDataURL }) => $blurDataURL ? 'blur(20px)' : 'none'};
  transform: ${({ $blurDataURL }) => $blurDataURL ? 'scale(1.1)' : 'none'};
  
  ${({ $showShimmer, theme }) => $showShimmer && `
    background: linear-gradient(
      90deg,
      ${theme.colors.surface || '#e0e0e0'} 25%,
      ${theme.colors.border || '#f0f0f0'} 50%,
      ${theme.colors.surface || '#e0e0e0'} 75%
    );
    background-size: 468px 104px;
    animation: ${shimmerAnimation} 1.2s ease-in-out infinite;
  `}
`;

const Image = styled.img<{
  $loaded: boolean;
  $error: boolean;
}>`
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
  opacity: ${({ $loaded, $error }) => ($loaded && !$error) ? 1 : 0};
  transition: opacity 0.3s ease-in-out;
  will-change: opacity;
`;

const LowQualityImage = styled(Image)`
  filter: blur(5px);
  transform: scale(1.05);
  z-index: 1;
`;

const HighQualityImage = styled(Image)`
  z-index: 2;
`;

const ErrorMessage = styled.div`
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  color: ${({ theme }) => theme.colors.textSecondary || '#666'};
  font-size: 0.875rem;
  text-align: center;
  z-index: 3;
`;

const LoadingSpinner = styled.div`
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 20px;
  height: 20px;
  border: 2px solid ${({ theme }) => theme.colors.border || '#e0e0e0'};
  border-top-color: ${({ theme }) => theme.colors.primary || '#007bff'};
  border-radius: 50%;
  animation: spin 1s linear infinite;
  z-index: 3;
  
  @keyframes spin {
    to {
      transform: translate(-50%, -50%) rotate(360deg);
    }
  }
`;

const LazyImage: React.FC<LazyImageProps> = ({
  src,
  alt,
  lowQualitySrc,
  blurDataURL,
  aspectRatio = '56.25%',
  priority = false,
  onLoadComplete,
  fallbackSrc = 'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMjQiIGhlaWdodD0iMjQiIHZpZXdCb3g9IjAgMCAyNCAyNCIgZmlsbD0ibm9uZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPHJlY3Qgd2lkdGg9IjI0IiBoZWlnaHQ9IjI0IiBmaWxsPSIjRjNGNEY2Ii8+CjxwYXRoIGQ9Ik0xOSAzSDVDMy44OTU0MyAzIDMgMy44OTU0MyAzIDVWMTlDMyAyMC4xMDQ2IDMuODk1NDMgMjEgNSAyMUgxOUMyMC4xMDQ2IDIxIDIxIDIwLjEwNDYgMjEgMTlWNUMyMSAzLjg5NTQzIDIwLjEwNDYgMyAxOSAzWiIgc3Ryb2tlPSIjOUI5RkE2IiBzdHJva2Utd2lkdGg9IjIiIGZpbGw9Im5vbmUiLz4KPHBhdGggZD0iTTguNSA5QzkuMzI4NDMgOSAxMCA4LjMyODQzIDEwIDcuNUMxMCA2LjY3MTU3IDkuMzI4NDMgNiA4LjUgNkM3LjY3MTU3IDYgNyA2LjY3MTU3IDcgNy41QzcgOC4zMjg0MyA3LjY3MTU3IDkgOC41IDlaIiBmaWxsPSIjOUI5RkE2Ii8+CjxwYXRoIGQ9Ik0yMSAxNUwxNiAxMEw1IDIxIiBzdHJva2U9IiM5QjlGQTYiIHN0cm9rZS13aWR0aD0iMiIgZmlsbD0ibm9uZSIvPgo8L3N2Zz4K',
  placeholderColor,
  showLoader = true,
  className,
  style,
  ...props
}) => {
  const [lowQualityLoaded, setLowQualityLoaded] = useState(false);
  const [highQualityLoaded, setHighQualityLoaded] = useState(false);
  const [error, setError] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [useFallback, setUseFallback] = useState(false);
  
  const lowQualityRef = useRef<HTMLImageElement>(null);
  const highQualityRef = useRef<HTMLImageElement>(null);
  
  const { ref: intersectionRef, inView } = useInView({
    triggerOnce: true,
    threshold: 0.1,
    rootMargin: '200px 0px',
    skip: priority
  });
  
  const shouldLoad = priority || inView;
  
  const handleLowQualityLoad = useCallback(() => {
    setLowQualityLoaded(true);
    setIsLoading(false);
  }, []);
  
  const handleHighQualityLoad = useCallback(() => {
    setHighQualityLoaded(true);
    setIsLoading(false);
    onLoadComplete?.();
  }, [onLoadComplete]);
  
  const handleError = useCallback(() => {
    setError(true);
    setIsLoading(false);
    if (!useFallback && fallbackSrc) {
      setUseFallback(true);
      setError(false);
    }
  }, [useFallback, fallbackSrc]);
  
  const currentSrc = useFallback ? fallbackSrc : src;
  
  useEffect(() => {
    if (shouldLoad && !highQualityLoaded && !error) {
      setIsLoading(true);
    }
  }, [shouldLoad, highQualityLoaded, error]);
  
  // Preload images when in view
  useEffect(() => {
    if (!shouldLoad || typeof window === 'undefined') return;
    
    // Preload low quality image first
    if (lowQualitySrc && !lowQualityLoaded) {
      const img = document.createElement('img');
      img.onload = handleLowQualityLoad;
      img.onerror = handleError;
      img.src = lowQualitySrc;
    }
    
    // Preload high quality image
    if (!highQualityLoaded) {
      const img = document.createElement('img');
      img.onload = handleHighQualityLoad;
      img.onerror = handleError;
      img.src = currentSrc;
    }
  }, [shouldLoad, lowQualitySrc, currentSrc, lowQualityLoaded, highQualityLoaded, handleLowQualityLoad, handleHighQualityLoad, handleError]);
  
  return (
    <ImageContainer
      ref={intersectionRef}
      $aspectRatio={aspectRatio}
      $placeholderColor={placeholderColor}
      className={className}
      style={style}
    >
      {/* Placeholder/Blur background */}
      <PlaceholderDiv
        $showShimmer={isLoading && showLoader && !blurDataURL}
        $blurDataURL={blurDataURL}
      />
      
      {/* Loading spinner */}
      {isLoading && showLoader && !blurDataURL && (
        <LoadingSpinner />
      )}
      
      {/* Low quality image */}
      {lowQualitySrc && shouldLoad && (
        <LowQualityImage
          ref={lowQualityRef}
          src={lowQualitySrc}
          alt=""
          $loaded={lowQualityLoaded}
          $error={false}
          loading="lazy"
          decoding="async"
          aria-hidden="true"
        />
      )}
      
      {/* High quality image */}
      {shouldLoad && (
        <HighQualityImage
          ref={highQualityRef}
          src={currentSrc}
          alt={alt}
          $loaded={highQualityLoaded}
          $error={error}
          loading={priority ? 'eager' : 'lazy'}
          decoding={priority ? 'sync' : 'async'}
          {...props}
        />
      )}
      
      {/* Error message */}
      {error && !useFallback && (
        <ErrorMessage>
          Failed to load image
        </ErrorMessage>
      )}
    </ImageContainer>
  );
};

export default LazyImage;