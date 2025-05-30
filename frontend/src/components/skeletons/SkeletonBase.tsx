import React from 'react';
import styled, { keyframes } from 'styled-components';

const shimmer = keyframes`
  0% {
    background-position: -468px 0;
  }
  100% {
    background-position: 468px 0;
  }
`;

const SkeletonElement = styled.div<{
  $width?: string;
  $height?: string;
  $borderRadius?: string;
  $margin?: string;
}>`
  background: linear-gradient(
    90deg,
    ${({ theme }) => theme.colors.surface || '#2a2a2a'} 25%,
    ${({ theme }) => theme.colors.border || '#3a3a3a'} 50%,
    ${({ theme }) => theme.colors.surface || '#2a2a2a'} 75%
  );
  background-size: 468px 104px;
  animation: ${shimmer} 1.2s ease-in-out infinite;
  border-radius: ${({ $borderRadius }) => $borderRadius || '4px'};
  width: ${({ $width }) => $width || '100%'};
  height: ${({ $height }) => $height || '20px'};
  margin: ${({ $margin }) => $margin || '0'};
  position: relative;
  overflow: hidden;
`;

interface SkeletonBaseProps {
  width?: string;
  height?: string;
  borderRadius?: string;
  margin?: string;
  className?: string;
}

const SkeletonBase: React.FC<SkeletonBaseProps> = ({
  width,
  height,
  borderRadius,
  margin,
  className
}) => {
  return (
    <SkeletonElement
      $width={width}
      $height={height}
      $borderRadius={borderRadius}
      $margin={margin}
      className={className}
    />
  );
};

export default SkeletonBase;