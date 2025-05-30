import React from 'react';
import styled from 'styled-components';
import SkeletonBase from './SkeletonBase';

const ProjectsSkeletonContainer = styled.div`
  max-width: 1200px;
  width: 100%;
  margin: 0 auto;
  padding: calc(2rem + 60px) 2rem 4rem 2rem;
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  background: ${({ theme }) => theme.colors.background};
  
  @media (max-width: 900px) {
    padding: calc(1.5rem + 60px) 1rem 3rem 1rem;
  }
  
  @media (max-width: 600px) {
    padding: calc(1rem + 60px) 0.75rem 2rem 0.75rem;
  }
`;

const PageHeader = styled.div`
  text-align: center;
  margin-bottom: 3rem;
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1.5rem;
`;

const StatsContainer = styled.div`
  background: ${({ theme }) => theme.colors.surface};
  padding: 1rem 1.5rem;
  border-radius: 12px;
  margin-bottom: 3rem;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.75rem;
`;

const ProjectsGrid = styled.div`
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(340px, 1fr));
  gap: 2rem;
  width: 100%;
  
  @media (max-width: 900px) {
    grid-template-columns: 1fr;
    gap: 1.5rem;
  }
`;

const ProjectCard = styled.div`
  background: ${({ theme }) => theme.colors.surface};
  border-radius: 16px;
  overflow: hidden;
  height: 260px;
  display: flex;
  flex-direction: column;
`;

const CardContent = styled.div`
  padding: 1.75rem 1.5rem;
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 1rem;
`;

const TechStack = styled.div`
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
`;

const ButtonGroup = styled.div`
  display: flex;
  gap: 1rem;
  margin-top: auto;
`;

const ProjectsSkeleton: React.FC = () => {
  return (
    <ProjectsSkeletonContainer>
      <PageHeader>
        <SkeletonBase width="300px" height="48px" borderRadius="8px" />
        <div style={{ maxWidth: '650px', width: '100%' }}>
          <SkeletonBase height="20px" margin="0 0 0.5rem 0" borderRadius="4px" />
          <SkeletonBase height="20px" width="85%" margin="0 auto 0.5rem" borderRadius="4px" />
          <SkeletonBase height="20px" width="70%" margin="0 auto" borderRadius="4px" />
        </div>
      </PageHeader>
      
      <StatsContainer>
        <SkeletonBase width="24px" height="24px" borderRadius="4px" />
        <SkeletonBase width="200px" height="20px" borderRadius="4px" />
        <SkeletonBase width="80px" height="24px" borderRadius="4px" />
      </StatsContainer>
      
      <ProjectsGrid>
        {Array.from({ length: 3 }).map((_, index) => (
          <ProjectCard key={index}>
            <CardContent>
              <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem' }}>
                <SkeletonBase width="24px" height="24px" borderRadius="4px" />
                <SkeletonBase width="180px" height="28px" borderRadius="6px" />
              </div>
              
              <div>
                <SkeletonBase height="18px" margin="0 0 0.5rem 0" borderRadius="4px" />
                <SkeletonBase height="18px" width="90%" margin="0 0 0.5rem 0" borderRadius="4px" />
                <SkeletonBase height="18px" width="75%" borderRadius="4px" />
              </div>
              
              <TechStack>
                {Array.from({ length: 4 }).map((_, techIndex) => (
                  <SkeletonBase 
                    key={techIndex}
                    width={`${Math.random() * 40 + 60}px`}
                    height="28px" 
                    borderRadius="14px"
                  />
                ))}
              </TechStack>
              
              <ButtonGroup>
                <SkeletonBase width="80px" height="36px" borderRadius="8px" />
                <SkeletonBase width="90px" height="36px" borderRadius="8px" />
              </ButtonGroup>
            </CardContent>
          </ProjectCard>
        ))}
      </ProjectsGrid>
    </ProjectsSkeletonContainer>
  );
};

export default ProjectsSkeleton;