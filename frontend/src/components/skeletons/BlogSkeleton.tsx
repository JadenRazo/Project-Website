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
  max-width: 1200px;
  margin: 0 auto;
  padding: calc(2rem + 60px) 2rem 4rem;
`;

const HeaderSkeleton = styled.div`
  text-align: center;
  margin-bottom: 2rem;
`;

const TitleBar = styled(SkeletonBox)`
  height: 3rem;
  width: 160px;
  margin: 0 auto 1rem;
`;

const DescBar = styled(SkeletonBox)`
  height: 1rem;
  width: 400px;
  max-width: 80%;
  margin: 0 auto;
`;

const SearchSkeleton = styled(SkeletonBox)`
  height: 48px;
  max-width: 500px;
  margin: 2rem auto;
`;

const TagsSkeleton = styled.div`
  display: flex;
  gap: 0.5rem;
  justify-content: center;
  margin-bottom: 2rem;
`;

const TagBar = styled(SkeletonBox)`
  height: 30px;
  width: 70px;
  border-radius: 20px;
`;

const Grid = styled.div`
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(340px, 1fr));
  gap: 1.5rem;

  @media (max-width: 768px) {
    grid-template-columns: 1fr;
  }
`;

const CardSkeleton = styled.div`
  border-radius: 12px;
  border: 1px solid ${({ theme }) => theme?.colors?.border || 'rgba(255,255,255,0.1)'};
  overflow: hidden;
`;

const CardImage = styled(SkeletonBox)`
  height: 200px;
  border-radius: 0;
`;

const CardBody = styled.div`
  padding: 1.25rem;
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
`;

const CardTitle = styled(SkeletonBox)`
  height: 1.2rem;
  width: 70%;
`;

const CardLine = styled(SkeletonBox)`
  height: 0.9rem;
  width: 100%;
`;

const CardMeta = styled(SkeletonBox)`
  height: 0.8rem;
  width: 40%;
`;

const BlogSkeleton: React.FC = () => (
  <Container>
    <HeaderSkeleton>
      <TitleBar />
      <DescBar />
    </HeaderSkeleton>
    <SearchSkeleton />
    <TagsSkeleton>
      {[1, 2, 3, 4].map((i) => (
        <TagBar key={i} />
      ))}
    </TagsSkeleton>
    <Grid>
      {[1, 2, 3].map((i) => (
        <CardSkeleton key={i}>
          <CardImage />
          <CardBody>
            <CardTitle />
            <CardLine />
            <CardLine style={{ width: '85%' }} />
            <CardMeta />
          </CardBody>
        </CardSkeleton>
      ))}
    </Grid>
  </Container>
);

export default BlogSkeleton;
