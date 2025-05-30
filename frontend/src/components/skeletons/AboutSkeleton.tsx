import React from 'react';
import styled from 'styled-components';
import SkeletonBase from './SkeletonBase';

const AboutSkeletonContainer = styled.div`
  min-height: calc(100vh - 200px);
  padding: 4rem 2rem;
  background: ${({ theme }) => theme.colors.background};
  margin-top: 60px;
  
  @media (max-width: ${({ theme }) => theme.breakpoints.tablet}) {
    padding: 3rem 1rem;
  }
`;

const ContentWrapper = styled.div`
  max-width: 1200px;
  margin: 0 auto;
  display: grid;
  grid-template-columns: 1fr;
  gap: 2rem;
  
  @media (min-width: ${({ theme }) => theme.breakpoints.tablet}) {
    grid-template-columns: 1fr 2fr;
    gap: 3rem;
  }
`;

const ProfileSection = styled.div`
  text-align: center;
  display: flex;
  flex-direction: column;
  align-items: center;
  
  @media (min-width: ${({ theme }) => theme.breakpoints.tablet}) {
    text-align: left;
    align-items: flex-start;
  }
`;

const MainContent = styled.div`
  display: flex;
  flex-direction: column;
  gap: 2rem;
`;

const Section = styled.div`
  background: ${({ theme }) => theme.colors.surface};
  padding: 2rem;
  border-radius: 12px;
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
`;

const SkillsGrid = styled.div`
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
  gap: 0.75rem;
`;

const ExperienceItem = styled.div`
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  margin-bottom: 2rem;
`;

const AboutSkeleton: React.FC = () => {
  return (
    <AboutSkeletonContainer>
      <ContentWrapper>
        <ProfileSection>
          <SkeletonBase 
            width="200px" 
            height="200px" 
            borderRadius="50%" 
            margin="0 0 1.5rem 0" 
          />
          <SkeletonBase width="180px" height="36px" margin="0 0 0.5rem 0" borderRadius="6px" />
          <SkeletonBase width="150px" height="24px" margin="0 0 2rem 0" borderRadius="4px" />
        </ProfileSection>

        <MainContent>
          {/* About Me Section */}
          <Section>
            <SkeletonBase width="120px" height="28px" borderRadius="6px" />
            <div>
              <SkeletonBase height="20px" margin="0 0 0.75rem 0" borderRadius="4px" />
              <SkeletonBase height="20px" width="95%" margin="0 0 0.75rem 0" borderRadius="4px" />
              <SkeletonBase height="20px" width="88%" borderRadius="4px" />
            </div>
          </Section>

          {/* Technical Skills Section */}
          <Section>
            <SkeletonBase width="140px" height="28px" borderRadius="6px" />
            <SkillsGrid>
              {Array.from({ length: 12 }).map((_, index) => (
                <SkeletonBase 
                  key={index}
                  height="40px" 
                  borderRadius="8px"
                />
              ))}
            </SkillsGrid>
          </Section>

          {/* Experience Section */}
          <Section>
            <SkeletonBase width="100px" height="28px" borderRadius="6px" />
            <ExperienceItem>
              <SkeletonBase width="200px" height="24px" borderRadius="4px" />
              <SkeletonBase width="180px" height="20px" borderRadius="4px" />
              <SkeletonBase width="120px" height="16px" borderRadius="4px" />
              <div style={{ marginTop: '1rem' }}>
                <SkeletonBase height="18px" margin="0 0 0.5rem 0" borderRadius="4px" />
                <SkeletonBase height="18px" width="92%" margin="0 0 0.5rem 0" borderRadius="4px" />
                <SkeletonBase height="18px" width="85%" margin="0 0 0.5rem 0" borderRadius="4px" />
                <SkeletonBase height="18px" width="78%" borderRadius="4px" />
              </div>
            </ExperienceItem>
          </Section>

          {/* Education Section */}
          <Section>
            <SkeletonBase width="80px" height="28px" borderRadius="6px" />
            <ExperienceItem>
              <SkeletonBase width="220px" height="24px" borderRadius="4px" />
              <SkeletonBase width="200px" height="20px" borderRadius="4px" />
              <SkeletonBase width="80px" height="16px" borderRadius="4px" />
              <div style={{ marginTop: '1rem' }}>
                <SkeletonBase height="18px" margin="0 0 0.5rem 0" borderRadius="4px" />
                <SkeletonBase height="18px" width="88%" borderRadius="4px" />
              </div>
            </ExperienceItem>
            <ExperienceItem>
              <SkeletonBase width="160px" height="24px" borderRadius="4px" />
              <SkeletonBase width="180px" height="20px" borderRadius="4px" />
              <SkeletonBase width="100px" height="16px" borderRadius="4px" />
              <div style={{ marginTop: '1rem' }}>
                <SkeletonBase height="18px" margin="0 0 0.5rem 0" borderRadius="4px" />
                <SkeletonBase height="18px" width="75%" borderRadius="4px" />
              </div>
            </ExperienceItem>
          </Section>
        </MainContent>
      </ContentWrapper>
    </AboutSkeletonContainer>
  );
};

export default AboutSkeleton;