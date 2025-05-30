import React from 'react';
import styled from 'styled-components';
import SkeletonBase from './SkeletonBase';

const HeroSkeletonContainer = styled.div`
  min-height: 100vh;
  padding: 120px 2rem 4rem;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
  max-width: 1200px;
  margin: 0 auto;
  gap: 2rem;
  
  @media (max-width: ${({ theme }) => theme.breakpoints.tablet}) {
    padding: 100px 1rem 3rem;
    gap: 1.5rem;
  }
`;

const TitleWrapper = styled.div`
  display: flex;
  flex-direction: column;
  gap: 1rem;
  width: 100%;
  max-width: 600px;
`;

const SubtitleWrapper = styled.div`
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  width: 100%;
  max-width: 500px;
`;

const ButtonWrapper = styled.div`
  display: flex;
  gap: 1rem;
  justify-content: center;
  flex-wrap: wrap;
  margin-top: 2rem;
`;

const HeroSkeleton: React.FC = () => {
  return (
    <HeroSkeletonContainer>
      <TitleWrapper>
        <SkeletonBase height="60px" borderRadius="8px" />
        <SkeletonBase height="48px" width="80%" margin="0 auto" borderRadius="8px" />
      </TitleWrapper>
      
      <SubtitleWrapper>
        <SkeletonBase height="24px" width="90%" margin="0 auto" borderRadius="4px" />
        <SkeletonBase height="24px" width="75%" margin="0 auto" borderRadius="4px" />
        <SkeletonBase height="24px" width="85%" margin="0 auto" borderRadius="4px" />
      </SubtitleWrapper>
      
      <ButtonWrapper>
        <SkeletonBase width="150px" height="50px" borderRadius="25px" />
        <SkeletonBase width="140px" height="50px" borderRadius="25px" />
      </ButtonWrapper>
    </HeroSkeletonContainer>
  );
};

export default HeroSkeleton;