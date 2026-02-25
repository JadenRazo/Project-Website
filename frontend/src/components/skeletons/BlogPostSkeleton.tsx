import React from 'react';
import styled, { keyframes } from 'styled-components';

const shimmer = keyframes`
  0% { background-position: -200% 0; }
  100% { background-position: 200% 0; }
`;

const SkeletonBox = styled.div`
  background: linear-gradient(
    90deg,
    ${({ theme }) => theme?.colors?.surface || 'rgba(255,255,255,0.05)'} 25%,
    ${({ theme }) => theme?.colors?.border || 'rgba(255,255,255,0.08)'} 50%,
    ${({ theme }) => theme?.colors?.surface || 'rgba(255,255,255,0.05)'} 75%
  );
  background-size: 200% 100%;
  animation: ${shimmer} 1.5s ease-in-out infinite;
  border-radius: 8px;
`;

const Container = styled.div`
  max-width: 800px;
  margin: 0 auto;
  padding: calc(2rem + 60px) 2rem 4rem;
`;

const BackBar = styled(SkeletonBox)`
  height: 1rem;
  width: 120px;
  margin-bottom: 2rem;
`;

const TagRow = styled.div`
  display: flex;
  gap: 0.4rem;
  margin-bottom: 1rem;
`;

const TagBar = styled(SkeletonBox)`
  height: 24px;
  width: 60px;
  border-radius: 10px;
`;

const TitleBar = styled(SkeletonBox)`
  height: 2.5rem;
  width: 80%;
  margin-bottom: 1rem;
`;

const MetaRow = styled.div`
  display: flex;
  gap: 1.5rem;
  margin-bottom: 2.5rem;
`;

const MetaBar = styled(SkeletonBox)`
  height: 0.9rem;
  width: 100px;
`;

const ImageBar = styled(SkeletonBox)`
  width: 100%;
  aspect-ratio: 21 / 9;
  margin-bottom: 2.5rem;
  border-radius: 12px;
`;

const ContentLine = styled(SkeletonBox)<{ $width?: string }>`
  height: 1rem;
  width: ${({ $width }) => $width || '100%'};
  margin-bottom: 0.75rem;
`;

const BlogPostSkeleton: React.FC = () => (
  <Container>
    <BackBar />
    <TagRow>
      <TagBar />
      <TagBar style={{ width: '50px' }} />
    </TagRow>
    <TitleBar />
    <MetaRow>
      <MetaBar />
      <MetaBar style={{ width: '80px' }} />
      <MetaBar style={{ width: '70px' }} />
    </MetaRow>
    <ImageBar />
    {[100, 95, 100, 80, 0, 100, 90, 100, 70, 0, 100, 85, 60].map((w, i) =>
      w === 0 ? (
        <div key={i} style={{ height: '1.5rem' }} />
      ) : (
        <ContentLine key={i} $width={`${w}%`} />
      )
    )}
  </Container>
);

export default BlogPostSkeleton;
