import React from 'react';
import styled from 'styled-components';
import SkeletonBase from './SkeletonBase';

const SkillsSkeletonContainer = styled.div`
  padding: 4rem 2rem;
  background: ${({ theme }) => theme.colors.background};
  min-height: 60vh;
  
  @media (max-width: ${({ theme }) => theme.breakpoints.tablet}) {
    padding: 3rem 1rem;
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
  gap: 1rem;
`;

const SkillsGrid = styled.div`
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 2rem;
  
  @media (max-width: ${({ theme }) => theme.breakpoints.tablet}) {
    grid-template-columns: 1fr;
    gap: 1.5rem;
  }
`;

const SkillCategory = styled.div`
  background: ${({ theme }) => theme.colors.surface};
  padding: 2rem;
  border-radius: 12px;
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
`;

const SkillItems = styled.div`
  display: flex;
  flex-direction: column;
  gap: 1rem;
`;

const SkillItem = styled.div`
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
`;

const SkillBar = styled.div`
  width: 100%;
  height: 8px;
  background: ${({ theme }) => theme.colors.border};
  border-radius: 4px;
  overflow: hidden;
  position: relative;
`;

const SkillProgress = styled(SkeletonBase)<{ width: string }>`
  position: absolute;
  top: 0;
  left: 0;
  height: 100%;
  width: ${({ width }) => width};
  border-radius: 0;
`;

const SkillsSkeleton: React.FC = () => {
  const skillCategories = [
    { name: 'Frontend', skills: 5 },
    { name: 'Backend', skills: 4 },
    { name: 'Tools & Platforms', skills: 6 },
    { name: 'Languages', skills: 4 }
  ];

  return (
    <SkillsSkeletonContainer>
      <ContentWrapper>
        <Header>
          <SkeletonBase width="200px" height="40px" borderRadius="8px" />
          <div style={{ maxWidth: '600px', width: '100%' }}>
            <SkeletonBase height="18px" margin="0 0 0.5rem 0" borderRadius="4px" />
            <SkeletonBase height="18px" width="80%" margin="0 auto" borderRadius="4px" />
          </div>
        </Header>
        
        <SkillsGrid>
          {skillCategories.map((category, categoryIndex) => (
            <SkillCategory key={categoryIndex}>
              <SkeletonBase width="120px" height="24px" borderRadius="6px" />
              
              <SkillItems>
                {Array.from({ length: category.skills }).map((_, skillIndex) => (
                  <SkillItem key={skillIndex}>
                    <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                      <SkeletonBase width="100px" height="16px" borderRadius="4px" />
                      <SkeletonBase width="30px" height="14px" borderRadius="4px" />
                    </div>
                    <SkillBar>
                      <SkillProgress 
                        width={`${Math.random() * 40 + 60}%`}
                        height="8px"
                      />
                    </SkillBar>
                  </SkillItem>
                ))}
              </SkillItems>
            </SkillCategory>
          ))}
        </SkillsGrid>
      </ContentWrapper>
    </SkillsSkeletonContainer>
  );
};

export default SkillsSkeleton;