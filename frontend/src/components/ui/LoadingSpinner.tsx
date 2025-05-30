import React from 'react';
import styled, { keyframes } from 'styled-components';

interface LoadingSpinnerProps {
  size?: 'small' | 'medium' | 'large' | number;
  color?: string;
  variant?: 'spin' | 'pulse' | 'dots' | 'bars';
  className?: string;
}

const spin = keyframes`
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
`;

const pulse = keyframes`
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
`;

const dotBounce = keyframes`
  0%, 80%, 100% { transform: scale(0); }
  40% { transform: scale(1); }
`;

const barStretch = keyframes`
  0%, 40%, 100% { transform: scaleY(0.4); }
  20% { transform: scaleY(1.0); }
`;

const getSize = (size: LoadingSpinnerProps['size']) => {
  if (typeof size === 'number') return `${size}px`;
  
  switch (size) {
    case 'small': return '16px';
    case 'large': return '32px';
    default: return '24px';
  }
};

const SpinnerContainer = styled.div<{
  $size: LoadingSpinnerProps['size'];
  $color?: string;
}>`
  display: inline-flex;
  align-items: center;
  justify-content: center;
  color: ${({ $color, theme }) => $color || theme.colors.primary || '#007bff'};
`;

const SpinSpinner = styled.div<{
  $size: LoadingSpinnerProps['size'];
}>`
  width: ${({ $size }) => getSize($size)};
  height: ${({ $size }) => getSize($size)};
  border: 2px solid transparent;
  border-top-color: currentColor;
  border-radius: 50%;
  animation: ${spin} 1s linear infinite;
`;

const PulseSpinner = styled.div<{
  $size: LoadingSpinnerProps['size'];
}>`
  width: ${({ $size }) => getSize($size)};
  height: ${({ $size }) => getSize($size)};
  background-color: currentColor;
  border-radius: 50%;
  animation: ${pulse} 1.5s ease-in-out infinite;
`;

const DotsContainer = styled.div<{
  $size: LoadingSpinnerProps['size'];
}>`
  display: flex;
  gap: ${({ $size }) => {
    const sizeValue = getSize($size);
    const num = parseInt(sizeValue.replace('px', ''));
    return `${num * 0.2}px`;
  }};
`;

const Dot = styled.div<{
  $size: LoadingSpinnerProps['size'];
  $delay: number;
}>`
  width: ${({ $size }) => {
    const sizeValue = getSize($size);
    const num = parseInt(sizeValue.replace('px', ''));
    return `${num * 0.3}px`;
  }};
  height: ${({ $size }) => {
    const sizeValue = getSize($size);
    const num = parseInt(sizeValue.replace('px', ''));
    return `${num * 0.3}px`;
  }};
  background-color: currentColor;
  border-radius: 50%;
  animation: ${dotBounce} 1.4s ease-in-out infinite both;
  animation-delay: ${({ $delay }) => $delay}s;
`;

const BarsContainer = styled.div<{
  $size: LoadingSpinnerProps['size'];
}>`
  display: flex;
  gap: ${({ $size }) => {
    const sizeValue = getSize($size);
    const num = parseInt(sizeValue.replace('px', ''));
    return `${num * 0.1}px`;
  }};
  align-items: center;
`;

const Bar = styled.div<{
  $size: LoadingSpinnerProps['size'];
  $delay: number;
}>`
  width: ${({ $size }) => {
    const sizeValue = getSize($size);
    const num = parseInt(sizeValue.replace('px', ''));
    return `${num * 0.2}px`;
  }};
  height: ${({ $size }) => getSize($size)};
  background-color: currentColor;
  animation: ${barStretch} 1.2s infinite ease-in-out;
  animation-delay: ${({ $delay }) => $delay}s;
`;

const LoadingSpinner: React.FC<LoadingSpinnerProps> = ({
  size = 'medium',
  color,
  variant = 'spin',
  className
}) => {
  const renderSpinner = () => {
    switch (variant) {
      case 'pulse':
        return <PulseSpinner $size={size} />;
      
      case 'dots':
        return (
          <DotsContainer $size={size}>
            <Dot $size={size} $delay={-0.32} />
            <Dot $size={size} $delay={-0.16} />
            <Dot $size={size} $delay={0} />
          </DotsContainer>
        );
      
      case 'bars':
        return (
          <BarsContainer $size={size}>
            <Bar $size={size} $delay={-0.4} />
            <Bar $size={size} $delay={-0.3} />
            <Bar $size={size} $delay={-0.2} />
            <Bar $size={size} $delay={-0.1} />
            <Bar $size={size} $delay={0} />
          </BarsContainer>
        );
      
      default:
        return <SpinSpinner $size={size} />;
    }
  };

  return (
    <SpinnerContainer $size={size} $color={color} className={className}>
      {renderSpinner()}
    </SpinnerContainer>
  );
};

export default LoadingSpinner;