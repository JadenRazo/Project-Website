import React, { useRef, useEffect, useState, useCallback } from 'react';
import styled from 'styled-components';
import { useDebounce } from '../../utils/performance';

// Try to import react-window, but have fallback if not available
let FixedSizeList: any;
try {
  const ReactWindow = require('react-window');
  FixedSizeList = ReactWindow.FixedSizeList;
} catch (error) {
  // react-window not available, will use fallback implementation
  console.debug('react-window not available, using fallback virtualization');
}

interface VirtualizedListProps<T> {
  data: T[];
  renderItem: (item: T, index: number, isVisible: boolean) => React.ReactNode;
  height?: string | number;
  itemHeight: number;
  overscan?: number;
  className?: string;
  onEndReached?: () => void;
  onEndReachedThreshold?: number;
  onScroll?: (event: React.UIEvent<HTMLDivElement>) => void;
  scrollToIndex?: number;
  keyExtractor: (item: T, index: number) => string;
  emptyComponent?: React.ReactNode;
  loadingComponent?: React.ReactNode;
  isLoading?: boolean;
}

interface ViewportItemsInfo {
  startIndex: number;
  endIndex: number;
  visibleStartIndex: number;
  visibleEndIndex: number;
}

interface VirtualContainerProps {
  height: string | number;
}

// Styled components
const Container = styled.div<VirtualContainerProps>`
  height: ${props => typeof props.height === 'number' ? `${props.height}px` : props.height};
  overflow-y: auto;
  position: relative;
  will-change: transform; /* Optimize for GPU acceleration */
  -webkit-overflow-scrolling: touch; /* Better scrolling on iOS */
`;

const InnerContainer = styled.div`
  position: relative;
  width: 100%;
  height: 0;
`;

const ItemContainer = styled.div`
  position: absolute;
  left: 0;
  width: 100%;
  will-change: transform; /* Optimize for GPU acceleration */
`;

/**
 * High-performance virtualized list component
 * Renders only items visible in the viewport plus overscan buffer
 * Significantly reduces memory usage for large lists
 */
export function VirtualizedList<T>(props: VirtualizedListProps<T>): React.ReactElement {
  // If react-window is available, use optimized implementation
  if (FixedSizeList) {
    return <OptimizedVirtualizedList {...props} />;
  }
  
  // Fallback to basic implementation
  return <BasicVirtualizedList {...props} />;
}

// Rename original implementation to OptimizedVirtualizedList
function OptimizedVirtualizedList<T>({
  data,
  renderItem,
  height = '100%',
  itemHeight,
  overscan = 5,
  className,
  onEndReached,
  onEndReachedThreshold = 0.8,
  onScroll,
  scrollToIndex,
  keyExtractor,
  emptyComponent,
  loadingComponent,
  isLoading = false
}: VirtualizedListProps<T>): React.ReactElement {
  const containerRef = useRef<HTMLDivElement>(null);
  const [viewportItems, setViewportItems] = useState<ViewportItemsInfo>({
    startIndex: 0,
    endIndex: 0,
    visibleStartIndex: 0,
    visibleEndIndex: 0
  });
  const onEndReachedRef = useRef(false);
  
  const totalHeight = data.length * itemHeight;
  
  // Calculate visible items based on scroll position
  const calculateVisibleItems = useCallback(() => {
    if (!containerRef.current) return;
    
    const { scrollTop, clientHeight } = containerRef.current;
    
    // Calculate visible indices
    const visibleStartIndex = Math.floor(scrollTop / itemHeight);
    const visibleEndIndex = Math.min(
      Math.ceil((scrollTop + clientHeight) / itemHeight),
      data.length - 1
    );
    
    // Add overscan buffer
    const startIndex = Math.max(0, visibleStartIndex - overscan);
    const endIndex = Math.min(data.length - 1, visibleEndIndex + overscan);
    
    setViewportItems({
      startIndex,
      endIndex,
      visibleStartIndex,
      visibleEndIndex
    });
    
    // Check if we've reached the end
    if (
      onEndReached &&
      !onEndReachedRef.current &&
      scrollTop + clientHeight >= totalHeight * onEndReachedThreshold
    ) {
      onEndReachedRef.current = true;
      onEndReached();
    } else if (scrollTop + clientHeight < totalHeight * onEndReachedThreshold) {
      onEndReachedRef.current = false;
    }
  }, [data.length, itemHeight, overscan, onEndReached, totalHeight, onEndReachedThreshold]);
  
  // Debounce the scroll handler for better performance
  const [debouncedCalculateVisibleItems] = useDebounce(calculateVisibleItems, 16);
  
  // Handle scrolling
  const handleScroll = useCallback((event: React.UIEvent<HTMLDivElement>) => {
    debouncedCalculateVisibleItems();
    if (onScroll) onScroll(event);
  }, [debouncedCalculateVisibleItems, onScroll]);
  
  // Initialize visible items on mount and data changes
  useEffect(() => {
    calculateVisibleItems();
    // Reset end reached flag when data changes
    onEndReachedRef.current = false;
  }, [calculateVisibleItems, data]);
  
  // Handle scroll to index
  useEffect(() => {
    if (scrollToIndex !== undefined && containerRef.current) {
      containerRef.current.scrollTop = scrollToIndex * itemHeight;
      calculateVisibleItems();
    }
  }, [scrollToIndex, itemHeight, calculateVisibleItems]);
  
  // Optimization: Only render items when browser is ready for layout
  useEffect(() => {
    if (typeof requestAnimationFrame !== 'undefined') {
      const rafId = requestAnimationFrame(() => {
        calculateVisibleItems();
      });
      
      return () => cancelAnimationFrame(rafId);
    }
  }, [calculateVisibleItems]);
  
  // Render empty or loading state if applicable
  if (isLoading && loadingComponent) {
    return <>{loadingComponent}</>;
  }
  
  if (data.length === 0 && emptyComponent) {
    return <>{emptyComponent}</>;
  }
  
  // Generate visible items
  const items = [];
  for (let i = viewportItems.startIndex; i <= viewportItems.endIndex; i++) {
    if (i >= data.length) break;
    
    const item = data[i];
    const isVisible = i >= viewportItems.visibleStartIndex && i <= viewportItems.visibleEndIndex;
    
    items.push(
      <ItemContainer
        key={keyExtractor(item, i)}
        style={{
          height: `${itemHeight}px`,
          top: `${i * itemHeight}px`,
          transform: 'translate3d(0, 0, 0)' // Force GPU acceleration
        }}
      >
        {renderItem(item, i, isVisible)}
      </ItemContainer>
    );
  }
  
  return (
    <Container
      ref={containerRef}
      height={height}
      onScroll={handleScroll}
      className={className}
    >
      <InnerContainer style={{ height: `${totalHeight}px` }}>
        {items}
      </InnerContainer>
    </Container>
  );
}

// Add basic implementation without react-window
function BasicVirtualizedList<T>({
  data,
  renderItem,
  height = '100%',
  className,
  emptyComponent,
  loadingComponent,
  isLoading = false,
  keyExtractor
}: VirtualizedListProps<T>): React.ReactElement {
  if (isLoading && loadingComponent) {
    return <>{loadingComponent}</>;
  }
  
  if (data.length === 0 && emptyComponent) {
    return <>{emptyComponent}</>;
  }
  
  const containerStyle = {
    height: typeof height === 'number' ? `${height}px` : height,
    overflowY: 'auto' as const
  };
  
  return (
    <div className={className} style={containerStyle}>
      {data.map((item, index) => (
        <div key={keyExtractor(item, index)}>
          {renderItem(item, index, true)}
        </div>
      ))}
    </div>
  );
}

export default VirtualizedList; 