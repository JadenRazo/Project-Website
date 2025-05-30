import React from 'react';
import styled, { keyframes } from 'styled-components';
import { motion } from 'framer-motion';

interface LoadingButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  loading?: boolean;
  variant?: 'primary' | 'secondary' | 'outline';
  size?: 'small' | 'medium' | 'large';
  loadingText?: string;
  children: React.ReactNode;
}

const spin = keyframes`
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
`;

const StyledButton = styled(motion.button)<{
  $variant: LoadingButtonProps['variant'];
  $size: LoadingButtonProps['size'];
  $loading: boolean;
}>`
  position: relative;
  border: none;
  border-radius: 8px;
  font-weight: 600;
  cursor: ${({ $loading }) => $loading ? 'not-allowed' : 'pointer'};
  transition: all 0.2s ease;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  opacity: ${({ $loading }) => $loading ? 0.7 : 1};
  
  ${({ $size }) => {
    switch ($size) {
      case 'small':
        return `
          padding: 0.5rem 1rem;
          font-size: 0.875rem;
          min-height: 32px;
        `;
      case 'large':
        return `
          padding: 1rem 2rem;
          font-size: 1.125rem;
          min-height: 48px;
        `;
      default:
        return `
          padding: 0.75rem 1.5rem;
          font-size: 1rem;
          min-height: 40px;
        `;
    }
  }}
  
  ${({ $variant, theme }) => {
    switch ($variant) {
      case 'secondary':
        return `
          background: ${theme.colors.surface || '#f8f9fa'};
          color: ${theme.colors.text || '#333'};
          border: 1px solid ${theme.colors.border || '#dee2e6'};
          
          &:hover:not(:disabled) {
            background: ${theme.colors.border || '#e9ecef'};
          }
        `;
      case 'outline':
        return `
          background: transparent;
          color: ${theme.colors.primary || '#007bff'};
          border: 1px solid ${theme.colors.primary || '#007bff'};
          
          &:hover:not(:disabled) {
            background: ${theme.colors.primary || '#007bff'};
            color: white;
          }
        `;
      default:
        return `
          background: ${theme.colors.primary || '#007bff'};
          color: white;
          
          &:hover:not(:disabled) {
            background: ${theme.colors.primaryHover || '#0056b3'};
          }
        `;
    }
  }}
  
  &:disabled {
    cursor: not-allowed;
    opacity: 0.6;
  }
  
  &:focus {
    outline: none;
    box-shadow: 0 0 0 3px ${({ theme }) => theme.colors.primary || '#007bff'}40;
  }
`;

const LoadingSpinner = styled.div<{ $size: LoadingButtonProps['size'] }>`
  border: 2px solid transparent;
  border-top-color: currentColor;
  border-radius: 50%;
  animation: ${spin} 1s linear infinite;
  
  ${({ $size }) => {
    switch ($size) {
      case 'small':
        return 'width: 14px; height: 14px;';
      case 'large':
        return 'width: 20px; height: 20px;';
      default:
        return 'width: 16px; height: 16px;';
    }
  }}
`;

const ButtonContent = styled.span<{ $loading: boolean }>`
  display: flex;
  align-items: center;
  gap: 0.5rem;
  opacity: ${({ $loading }) => $loading ? 0 : 1};
  transition: opacity 0.2s ease;
`;

const LoadingContent = styled.div`
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  display: flex;
  align-items: center;
  gap: 0.5rem;
`;

const LoadingButton: React.FC<LoadingButtonProps> = ({
  loading = false,
  variant = 'primary',
  size = 'medium',
  loadingText,
  children,
  disabled,
  onClick,
  ...props
}) => {
  const handleClick = (event: React.MouseEvent<HTMLButtonElement>) => {
    if (loading || disabled) {
      event.preventDefault();
      return;
    }
    onClick?.(event);
  };

  // Exclude HTML events that conflict with Framer Motion
  const {
    onDrag,
    onDragEnd,
    onDragStart,
    onAnimationStart,
    onAnimationEnd,
    onPointerDown,
    onPointerMove,
    onPointerUp,
    onPointerCancel,
    onPointerEnter,
    onPointerLeave,
    onPointerOut,
    onPointerOver,
    ...buttonProps
  } = props;

  return (
    <StyledButton
      $variant={variant}
      $size={size}
      $loading={loading}
      disabled={disabled || loading}
      onClick={handleClick}
      whileHover={{ scale: loading ? 1 : 1.02 }}
      whileTap={{ scale: loading ? 1 : 0.98 }}
      {...buttonProps}
    >
      <ButtonContent $loading={loading}>
        {children}
      </ButtonContent>
      
      {loading && (
        <LoadingContent>
          <LoadingSpinner $size={size} />
          {loadingText && <span>{loadingText}</span>}
        </LoadingContent>
      )}
    </StyledButton>
  );
};

export default LoadingButton;