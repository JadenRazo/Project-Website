import React from 'react';
import styled from 'styled-components';

interface BadgeProps {
  variant: 'featured' | 'category' | 'status' | 'expiry' | 'expired';
  children: React.ReactNode;
  className?: string;
  style?: React.CSSProperties;
}

const StyledBadge = styled.span<{ variant: BadgeProps['variant'] }>`
  display: inline-flex;
  align-items: center;
  padding: 0.4rem 0.8rem;
  font-size: 0.75rem;
  border-radius: 6px;
  font-weight: 600;
  white-space: nowrap;
  transition: all 0.2s ease;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  
  &:hover {
    transform: translateY(-1px);
    box-shadow: 0 2px 6px rgba(0, 0, 0, 0.15);
  }
  
  ${({ theme, variant }) => {
    switch (variant) {
      case 'featured':
        return `
          background: linear-gradient(135deg, ${theme.colors.primary}25, ${theme.colors.primary}15);
          color: ${theme.colors.primary};
          border: 1px solid ${theme.colors.primary}30;
          
          &::before {
            content: '★';
            margin-right: 0.3rem;
            font-size: 0.8rem;
          }
        `;
      case 'category':
        return `
          background: linear-gradient(135deg, ${theme.colors.secondary || theme.colors.accent}25, ${theme.colors.secondary || theme.colors.accent}15);
          color: ${theme.colors.secondary || theme.colors.accent};
          border: 1px solid ${theme.colors.secondary || theme.colors.accent}30;
          
          &::before {
            content: '📁';
            margin-right: 0.3rem;
            font-size: 0.7rem;
          }
        `;
      case 'status':
        return `
          background: linear-gradient(135deg, ${theme.colors.warning || '#F59E0B'}25, ${theme.colors.warning || '#F59E0B'}15);
          color: ${theme.colors.warning || '#F59E0B'};
          border: 1px solid ${theme.colors.warning || '#F59E0B'}30;
          
          &::before {
            content: '🔒';
            margin-right: 0.3rem;
            font-size: 0.7rem;
          }
        `;
      case 'expiry':
        return `
          background: linear-gradient(135deg, ${theme.colors.warning || '#F59E0B'}25, ${theme.colors.warning || '#F59E0B'}15);
          color: ${theme.colors.warning || '#F59E0B'};
          border: 1px solid ${theme.colors.warning || '#F59E0B'}30;
          
          &::before {
            content: '⏰';
            margin-right: 0.3rem;
            font-size: 0.7rem;
          }
        `;
      case 'expired':
        return `
          background: linear-gradient(135deg, ${theme.colors.error || '#DC2626'}25, ${theme.colors.error || '#DC2626'}15);
          color: ${theme.colors.error || '#DC2626'};
          border: 1px solid ${theme.colors.error || '#DC2626'}30;
          
          &::before {
            content: '⚠️';
            margin-right: 0.3rem;
            font-size: 0.7rem;
          }
        `;
      default:
        return `
          background: linear-gradient(135deg, ${theme.colors.primary}25, ${theme.colors.primary}15);
          color: ${theme.colors.primary};
          border: 1px solid ${theme.colors.primary}30;
        `;
    }
  }}
  
  @media (max-width: 768px) {
    padding: 0.3rem 0.6rem;
    font-size: 0.7rem;
  }
`;

export const Badge: React.FC<BadgeProps> = ({ variant, children, className, style }) => {
  return (
    <StyledBadge variant={variant} className={className} style={style}>
      {children}
    </StyledBadge>
  );
};

export default Badge;