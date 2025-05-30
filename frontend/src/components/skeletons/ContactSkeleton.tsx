import React from 'react';
import styled from 'styled-components';
import SkeletonBase from './SkeletonBase';

const ContactSkeletonContainer = styled.div`
  min-height: calc(100vh - 200px);
  padding: calc(4rem + 60px) 2rem 4rem;
  background: ${({ theme }) => theme.colors.background};
  
  @media (max-width: ${({ theme }) => theme.breakpoints.tablet}) {
    padding: calc(3rem + 60px) 1rem 3rem;
  }
`;

const ContentWrapper = styled.div`
  max-width: 800px;
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

const ContactForm = styled.div`
  background: ${({ theme }) => theme.colors.surface};
  padding: 2.5rem;
  border-radius: 12px;
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
`;

const FormGroup = styled.div`
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
`;

const FormRow = styled.div`
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
  
  @media (max-width: ${({ theme }) => theme.breakpoints.mobile}) {
    grid-template-columns: 1fr;
  }
`;

const ContactInfo = styled.div`
  background: ${({ theme }) => theme.colors.surface};
  padding: 2rem;
  border-radius: 12px;
  display: flex;
  flex-direction: column;
  gap: 2rem;
`;

const ContactItem = styled.div`
  display: flex;
  align-items: center;
  gap: 1rem;
`;

const SocialLinks = styled.div`
  display: flex;
  justify-content: center;
  gap: 1rem;
  flex-wrap: wrap;
`;

const ContactSkeleton: React.FC = () => {
  return (
    <ContactSkeletonContainer>
      <ContentWrapper>
        <Header>
          <SkeletonBase width="200px" height="40px" borderRadius="8px" />
          <div style={{ maxWidth: '500px', width: '100%' }}>
            <SkeletonBase height="18px" margin="0 0 0.5rem 0" borderRadius="4px" />
            <SkeletonBase height="18px" width="85%" margin="0 auto" borderRadius="4px" />
          </div>
        </Header>
        
        <ContactForm>
          <SkeletonBase width="150px" height="28px" borderRadius="6px" />
          
          <FormRow>
            <FormGroup>
              <SkeletonBase width="60px" height="16px" borderRadius="4px" />
              <SkeletonBase height="44px" borderRadius="8px" />
            </FormGroup>
            <FormGroup>
              <SkeletonBase width="50px" height="16px" borderRadius="4px" />
              <SkeletonBase height="44px" borderRadius="8px" />
            </FormGroup>
          </FormRow>
          
          <FormGroup>
            <SkeletonBase width="70px" height="16px" borderRadius="4px" />
            <SkeletonBase height="44px" borderRadius="8px" />
          </FormGroup>
          
          <FormGroup>
            <SkeletonBase width="80px" height="16px" borderRadius="4px" />
            <SkeletonBase height="120px" borderRadius="8px" />
          </FormGroup>
          
          <SkeletonBase width="120px" height="48px" borderRadius="24px" margin="1rem 0 0 0" />
        </ContactForm>
        
        <ContactInfo>
          <SkeletonBase width="140px" height="24px" borderRadius="6px" />
          
          <div style={{ display: 'flex', flexDirection: 'column', gap: '1.5rem' }}>
            <ContactItem>
              <SkeletonBase width="24px" height="24px" borderRadius="4px" />
              <div style={{ flex: 1 }}>
                <SkeletonBase width="180px" height="18px" borderRadius="4px" />
              </div>
            </ContactItem>
            
            <ContactItem>
              <SkeletonBase width="24px" height="24px" borderRadius="4px" />
              <div style={{ flex: 1 }}>
                <SkeletonBase width="200px" height="18px" borderRadius="4px" />
              </div>
            </ContactItem>
            
            <ContactItem>
              <SkeletonBase width="24px" height="24px" borderRadius="4px" />
              <div style={{ flex: 1 }}>
                <SkeletonBase width="160px" height="18px" borderRadius="4px" />
              </div>
            </ContactItem>
          </div>
          
          <div>
            <SkeletonBase width="120px" height="20px" margin="0 0 1rem 0" borderRadius="4px" />
            <SocialLinks>
              {Array.from({ length: 5 }).map((_, index) => (
                <SkeletonBase 
                  key={index}
                  width="48px" 
                  height="48px" 
                  borderRadius="50%" 
                />
              ))}
            </SocialLinks>
          </div>
        </ContactInfo>
      </ContentWrapper>
    </ContactSkeletonContainer>
  );
};

export default ContactSkeleton;