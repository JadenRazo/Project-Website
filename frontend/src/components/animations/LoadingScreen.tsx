// src/components/animations/LoadingScreen.tsx
import React, { FC, ReactNode } from 'react';
import styled, { keyframes, css } from 'styled-components';
import { motion, AnimatePresence } from 'framer-motion';

// Types
interface SkeletonProps {
  width?: string;
  height?: string;
  borderRadius?: string;
  margin?: string;
  animated?: boolean;
}

interface SkeletonTextProps extends SkeletonProps {
  lines?: number;
  lineHeight?: string;
}

interface SkeletonScreenProps {
  isLoading: boolean;
  children?: ReactNode;
  template?: 'profile' | 'card' | 'article' | 'custom';
  customContent?: ReactNode;
  backgroundColor?: string;
  fullScreen?: boolean;
}

// Animations
const shimmer = keyframes`
  0% {
    background-position: -1000px 0;
  }
  100% {
    background-position: 1000px 0;
  }
`;

const fadeInOut = {
  initial: { opacity: 0 },
  animate: { opacity: 1, transition: { duration: 0.3 } },
  exit: { opacity: 0, transition: { duration: 0.2 } }
};

// Styled Components
const LoadingContainer = styled(motion.div)<{ fullScreen?: boolean }>`
  ${props => props.fullScreen ? `
    position: fixed;
    top: 0;
    left: 0;
    z-index: 9999;
  ` : `
    position: relative;
  `}
  width: 100%;
  height: ${props => props.fullScreen ? '100vh' : '100%'};
  background: ${props => props.theme.colors.background || '#0f0f0f'};
  display: flex;
  flex-direction: column;
  justify-content: flex-start;
  align-items: center;
  padding: ${props => props.fullScreen ? '20vh 5%' : '20px'};
  overflow: hidden;
`;

const ContentContainer = styled.div`
  width: 100%;
  max-width: 1200px;
  display: flex;
  flex-direction: column;
  align-items: center;
`;

const SkeletonElement = styled.div<SkeletonProps>`
  width: ${props => props.width || '100%'};
  height: ${props => props.height || '20px'};
  border-radius: ${props => props.borderRadius || '4px'};
  margin: ${props => props.margin || '0 0 10px 0'};
  background-color: ${props => props.theme.colors.primaryLight || '#2a2a2a'};
  opacity: 0.7;
  position: relative;
  overflow: hidden;
  
  ${props => props.animated && css`
    background: linear-gradient(to right, 
      ${props.theme.colors.primaryLight || '#2a2a2a'} 8%, 
      ${props.theme.colors.primary || '#444'} 18%, 
      ${props.theme.colors.primaryLight || '#2a2a2a'} 33%);
    background-size: 2000px 100%;
    animation: ${shimmer} 1.5s linear infinite;
  `}
`;

const SkeletonRow = styled.div`
  display: flex;
  width: 100%;
  justify-content: space-between;
  gap: 20px;
  margin-bottom: 15px;
`;

const SkeletonContainer = styled.div`
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
`;

// Skeleton Components
const SkeletonText: FC<SkeletonTextProps> = ({ 
  lines = 3, 
  lineHeight = '20px', 
  animated = true,
  ...rest 
}) => (
  <SkeletonContainer>
    {Array(lines).fill(0).map((_, i) => (
      <SkeletonElement 
        key={`text-line-${i}`}
        height={lineHeight}
        width={i === lines - 1 && lines > 1 ? `${Math.floor(Math.random() * 40) + 30}%` : '100%'}
        margin="0 0 10px 0"
        animated={animated}
        {...rest}
      />
    ))}
  </SkeletonContainer>
);

const SkeletonCircle: FC<SkeletonProps> = (props) => (
  <SkeletonElement 
    borderRadius="50%" 
    animated={true}
    {...props} 
  />
);

const SkeletonImage: FC<SkeletonProps> = (props) => (
  <SkeletonElement 
    height="200px" 
    animated={true}
    {...props} 
  />
);

const SkeletonButton: FC<SkeletonProps> = (props) => (
  <SkeletonElement
    width="120px"
    height="40px"
    borderRadius="4px"
    animated={true}
    {...props}
  />
);

// Template Generators
const generateProfileTemplate = () => (
  <ContentContainer>
    <SkeletonRow>
      <SkeletonCircle width="100px" height="100px" margin="0 20px 0 0" />
      <SkeletonContainer>
        <SkeletonElement width="70%" height="30px" margin="10px 0 20px 0" />
        <SkeletonText lines={2} lineHeight="16px" />
      </SkeletonContainer>
    </SkeletonRow>
    <SkeletonRow>
      <SkeletonText lines={4} />
    </SkeletonRow>
    <SkeletonRow>
      <SkeletonButton width="120px" />
      <SkeletonButton width="120px" />
    </SkeletonRow>
  </ContentContainer>
);

const generateCardTemplate = () => (
  <ContentContainer>
    <SkeletonRow>
      <SkeletonContainer>
        <SkeletonImage width="100%" height="250px" margin="0 0 15px 0" />
        <SkeletonElement width="60%" height="28px" margin="0 0 15px 0" />
        <SkeletonText lines={3} lineHeight="18px" />
        <SkeletonElement width="100px" height="35px" margin="15px 0 0 0" />
      </SkeletonContainer>
      <SkeletonContainer>
        <SkeletonImage width="100%" height="250px" margin="0 0 15px 0" />
        <SkeletonElement width="60%" height="28px" margin="0 0 15px 0" />
        <SkeletonText lines={3} lineHeight="18px" />
        <SkeletonElement width="100px" height="35px" margin="15px 0 0 0" />
      </SkeletonContainer>
    </SkeletonRow>
  </ContentContainer>
);

const generateArticleTemplate = () => (
  <ContentContainer>
    <SkeletonElement width="80%" height="40px" margin="0 0 20px 0" />
    <SkeletonElement width="50%" height="25px" margin="0 0 30px 0" />
    <SkeletonImage width="100%" height="300px" margin="0 0 30px 0" />
    <SkeletonText lines={7} lineHeight="22px" margin="0 0 15px 0" />
    
    <SkeletonElement width="100%" height="1px" margin="30px 0" />
    
    <SkeletonText lines={5} lineHeight="22px" />
    
    <SkeletonRow>
      <SkeletonCircle width="50px" height="50px" />
      <SkeletonContainer>
        <SkeletonElement width="150px" height="20px" margin="5px 0" />
        <SkeletonElement width="100px" height="16px" />
      </SkeletonContainer>
    </SkeletonRow>
  </ContentContainer>
);

// Main LoadingScreen Component
export const LoadingScreen: FC<SkeletonScreenProps> = ({
  isLoading,
  children,
  template = 'article',
  customContent,
  backgroundColor,
  fullScreen = true
}) => {
  // Select the appropriate template
  const getTemplateContent = () => {
    switch (template) {
      case 'profile':
        return generateProfileTemplate();
      case 'card':
        return generateCardTemplate();
      case 'article':
        return generateArticleTemplate();
      case 'custom':
        return customContent;
      default:
        return generateArticleTemplate();
    }
  };

  return (
    <AnimatePresence mode="wait">
      {isLoading ? (
        <LoadingContainer 
          key="loading"
          fullScreen={fullScreen}
          style={{ background: backgroundColor }}
          {...fadeInOut}
        >
          {getTemplateContent()}
        </LoadingContainer>
      ) : (
        children && (
          <motion.div
            key="content"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
          >
            {children}
          </motion.div>
        )
      )}
    </AnimatePresence>
  );
};

// Exporting individual skeleton components for reuse
export const Skeleton = {
  Text: SkeletonText,
  Circle: SkeletonCircle,
  Image: SkeletonImage,
  Button: SkeletonButton,
  Element: SkeletonElement,
  Row: SkeletonRow,
  Container: SkeletonContainer
};
