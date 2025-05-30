import React from 'react';
import styled from 'styled-components';
import SkeletonBase from './SkeletonBase';

const GenericSkeletonContainer = styled.div`
  min-height: calc(100vh - 200px);
  padding: calc(4rem + 60px) 2rem 4rem;
  background: ${({ theme }) => theme.colors.background};
  
  @media (max-width: ${({ theme }) => theme.breakpoints.tablet}) {
    padding: calc(3rem + 60px) 1rem 3rem;
  }
`;

const ContentWrapper = styled.div`
  max-width: 1200px;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  gap: 3rem;
`;

const Header = styled.div`
  text-align: center;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1.5rem;
`;

const MainContent = styled.div`
  display: flex;
  flex-direction: column;
  gap: 2rem;
`;

const ContentSection = styled.div`
  background: ${({ theme }) => theme.colors.surface};
  padding: 2rem;
  border-radius: 12px;
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
`;

const GenericPageSkeleton: React.FC = () => {
  return (
    <GenericSkeletonContainer>
      <ContentWrapper>
        <Header>
          <SkeletonBase width="300px" height="48px" borderRadius="8px" />
          <div style={{ maxWidth: '600px', width: '100%' }}>
            <SkeletonBase height="20px" margin="0 0 0.5rem 0" borderRadius="4px" />
            <SkeletonBase height="20px" width="85%" margin="0 auto 0.5rem" borderRadius="4px" />
            <SkeletonBase height="20px" width="70%" margin="0 auto" borderRadius="4px" />
          </div>
        </Header>
        
        <MainContent>
          <ContentSection>
            <SkeletonBase width="200px" height="28px" borderRadius="6px" />
            <div>
              <SkeletonBase height="18px" margin="0 0 0.75rem 0" borderRadius="4px" />
              <SkeletonBase height="18px" width="95%" margin="0 0 0.75rem 0" borderRadius="4px" />
              <SkeletonBase height="18px" width="88%" margin="0 0 0.75rem 0" borderRadius="4px" />
              <SkeletonBase height="18px" width="92%" margin="0 0 0.75rem 0" borderRadius="4px" />
              <SkeletonBase height="18px" width="85%" borderRadius="4px" />
            </div>
          </ContentSection>
          
          <ContentSection>
            <SkeletonBase width="180px" height="28px" borderRadius="6px" />
            <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(200px, 1fr))', gap: '1rem' }}>
              {Array.from({ length: 6 }).map((_, index) => (
                <SkeletonBase 
                  key={index}
                  height="60px" 
                  borderRadius="8px"
                />
              ))}
            </div>
          </ContentSection>
          
          <ContentSection>
            <SkeletonBase width="150px" height="28px" borderRadius="6px" />
            <div>
              <SkeletonBase height="18px" margin="0 0 0.75rem 0" borderRadius="4px" />
              <SkeletonBase height="18px" width="90%" margin="0 0 0.75rem 0" borderRadius="4px" />
              <SkeletonBase height="18px" width="87%" margin="0 0 1.5rem 0" borderRadius="4px" />
              
              <div style={{ display: 'flex', gap: '1rem', marginBottom: '1rem' }}>
                <SkeletonBase width="120px" height="40px" borderRadius="8px" />
                <SkeletonBase width="100px" height="40px" borderRadius="8px" />
              </div>
              
              <SkeletonBase height="18px" margin="0 0 0.75rem 0" borderRadius="4px" />
              <SkeletonBase height="18px" width="93%" borderRadius="4px" />
            </div>
          </ContentSection>
        </MainContent>
      </ContentWrapper>
    </GenericSkeletonContainer>
  );
};

export default GenericPageSkeleton;