import React from 'react';
import styled from 'styled-components';
import { SkeletonLoader, AnimatedContainer, pulse } from '../animations/AnimatedComponents';

const LoadingContainer = styled.div`
  padding: 2rem;
  max-width: 1200px;
  margin: 0 auto;
`;

const HeaderSkeleton = styled.div`
  margin-bottom: 2rem;
  padding-bottom: 1rem;
  border-bottom: 1px solid ${({ theme }) => theme.colors.border};
`;

const MetricsGridSkeleton = styled.div`
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
`;

const CardSkeleton = styled.div`
  background: ${({ theme }) => theme.colors.card};
  border-radius: 8px;
  padding: 1.5rem;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.05);
  border: 1px solid ${({ theme }) => theme.colors.border};
`;

const ServiceGridSkeleton = styled.div`
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
  gap: 1.5rem;
`;

const ChartSkeleton = styled.div`
  background: ${({ theme }) => theme.colors.card};
  border-radius: 8px;
  padding: 1.5rem;
  height: 400px;
  margin-bottom: 2rem;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
`;

const PulsingIcon = styled.div`
  width: 60px;
  height: 60px;
  border-radius: 50%;
  background: ${({ theme }) => theme.colors.primary}20;
  display: flex;
  align-items: center;
  justify-content: center;
  animation: ${pulse} 1.5s ease-in-out infinite;
  
  &::after {
    content: '📊';
    font-size: 24px;
  }
`;

const ButtonGroupSkeleton = styled.div`
  display: flex;
  gap: 0.5rem;
  margin-top: 1rem;
`;

const DevPanelLoadingState: React.FC = () => {
  return (
    <LoadingContainer>
      <AnimatedContainer staggerDelay={100}>
        <HeaderSkeleton>
          <SkeletonLoader width="200px" height="40px" borderRadius="4px" />
          <div style={{ marginTop: '0.5rem' }}>
            <SkeletonLoader width="300px" height="20px" borderRadius="4px" />
          </div>
        </HeaderSkeleton>

        {/* System Metrics Skeleton */}
        <MetricsGridSkeleton>
          {[1, 2, 3, 4].map((i) => (
            <CardSkeleton key={`metric-${i}`}>
              <SkeletonLoader width="120px" height="16px" borderRadius="4px" />
              <div style={{ margin: '1rem 0' }}>
                <SkeletonLoader width="80px" height="32px" borderRadius="4px" />
              </div>
              <SkeletonLoader width="150px" height="14px" borderRadius="4px" />
            </CardSkeleton>
          ))}
        </MetricsGridSkeleton>

        {/* Charts Skeleton */}
        <ChartSkeleton>
          <PulsingIcon />
        </ChartSkeleton>

        {/* Services Skeleton */}
        <ServiceGridSkeleton>
          {[1, 2, 3, 4, 5, 6].map((i) => (
            <CardSkeleton key={`service-${i}`}>
              <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '1rem' }}>
                <SkeletonLoader width="100px" height="20px" borderRadius="4px" />
                <SkeletonLoader width="60px" height="24px" borderRadius="12px" />
              </div>
              
              <div style={{ display: 'grid', gridTemplateColumns: 'repeat(2, 1fr)', gap: '0.5rem', marginBottom: '1rem' }}>
                {[1, 2, 3, 4].map((j) => (
                  <div key={j} style={{ textAlign: 'center' }}>
                    <SkeletonLoader width="60px" height="24px" borderRadius="4px" />
                    <div style={{ marginTop: '0.25rem' }}>
                      <SkeletonLoader width="50px" height="12px" borderRadius="4px" />
                    </div>
                  </div>
                ))}
              </div>
              
              <ButtonGroupSkeleton>
                {[1, 2, 3, 4, 5].map((j) => (
                  <SkeletonLoader key={j} width="60px" height="32px" borderRadius="4px" />
                ))}
              </ButtonGroupSkeleton>
            </CardSkeleton>
          ))}
        </ServiceGridSkeleton>
      </AnimatedContainer>
    </LoadingContainer>
  );
};

export default DevPanelLoadingState;