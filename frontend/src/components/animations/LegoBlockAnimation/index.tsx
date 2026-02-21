import React, { useEffect, useState } from 'react';
import styled from 'styled-components';
import { motion } from 'framer-motion';
import { useTheme } from '../../../hooks/useTheme';
import { useLegoAnimation } from './useLegoAnimation';
import { LegoCanvas } from './LegoCanvas';
import type { LegoBlockAnimationProps } from './types';

const AnimationContainer = styled.section`
  position: relative;
  width: 100%;
  min-height: 100vh;
  overflow: hidden;
  background: ${({ theme }) => theme.colors.background};
  display: flex;
  align-items: center;
  justify-content: center;
`;

const CanvasOverlay = styled(motion.div)`
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 10;
  pointer-events: none;
  background: ${({ theme }) => theme.colors.background};
`;

const ContentWrapper = styled.div`
  position: relative;
  width: 100%;
  max-width: 900px;
  margin: 0 auto;
  padding: 4rem 1.5rem;
  z-index: 1;

  @media (min-width: 768px) {
    padding: 5rem 2rem;
  }
`;

const BioContent = styled(motion.div)`
  position: relative;
`;

const SectionHeading = styled(motion.h2)`
  font-size: 2.5rem;
  font-weight: 700;
  margin-bottom: 2rem;
  position: relative;
  display: inline-block;
  color: ${({ theme }) => theme.colors.primary};

  &::after {
    content: '';
    position: absolute;
    left: 0;
    bottom: -10px;
    height: 3px;
    width: 60%;
    background: ${({ theme }) => theme.colors.accent};
  }

  @media (max-width: 768px) {
    font-size: 2rem;
  }
`;

const BioText = styled(motion.p)`
  font-size: 1.25rem;
  line-height: 1.8;
  color: ${({ theme }) => theme.colors.text};
  max-width: 700px;

  @media (max-width: 768px) {
    font-size: 1.1rem;
    line-height: 1.7;
  }
`;

const bioVariants = {
  hidden: {
    opacity: 0,
    y: 30,
  },
  visible: {
    opacity: 1,
    y: 0,
    transition: {
      duration: 0.8,
      ease: 'easeOut',
      staggerChildren: 0.2,
      delayChildren: 0.1,
    },
  },
};

const itemVariants = {
  hidden: {
    opacity: 0,
    y: 20,
  },
  visible: {
    opacity: 1,
    y: 0,
    transition: {
      duration: 0.6,
      ease: 'easeOut',
    },
  },
};

export const LegoBlockAnimation: React.FC<LegoBlockAnimationProps> = ({
  onAnimationComplete,
}) => {
  const { theme } = useTheme();
  const {
    config,
    skipAnimation,
    isVisible,
    containerRef,
    phase,
    setPhase,
    showBio,
    setShowBio,
    dimensions,
  } = useLegoAnimation();

  const [animationReady, setAnimationReady] = useState(false);
  const [canvasOpacity, setCanvasOpacity] = useState(1);

  useEffect(() => {
    if (dimensions.width > 0 && dimensions.height > 0 && !skipAnimation) {
      const timer = setTimeout(() => {
        setAnimationReady(true);
      }, 100);
      return () => clearTimeout(timer);
    }
  }, [dimensions, skipAnimation]);

  const handleDissolveStart = () => {
    setShowBio(true);
  };

  const handlePhaseChange = (newPhase: typeof phase) => {
    setPhase(newPhase);
    if (newPhase === 'dissolve') {
      setCanvasOpacity(0.8);
    }
    if (newPhase === 'revealed') {
      setCanvasOpacity(0);
      if (onAnimationComplete) {
        onAnimationComplete();
      }
    }
  };

  const shouldRunAnimation = !skipAnimation &&
    animationReady &&
    dimensions.width > 0 &&
    dimensions.height > 0;

  const bioIsVisible = showBio || skipAnimation || phase === 'revealed' || phase === 'dissolve';

  return (
    <AnimationContainer ref={containerRef} id="about">
      <ContentWrapper>
        <BioContent
          variants={bioVariants}
          initial={skipAnimation ? "visible" : "hidden"}
          animate={bioIsVisible ? "visible" : "hidden"}
        >
          <SectionHeading variants={itemVariants}>
            About Me
          </SectionHeading>

          <BioText variants={itemVariants}>
            Building solutions block by block. I build robust, scalable
            applications with React, TypeScript, Go, and Python. From Go
            microservices to auto detailing platforms to AI-powered Discord bots,
            I build production systems that bridge functionality with clean,
            user-friendly design.
          </BioText>
        </BioContent>
      </ContentWrapper>

      {shouldRunAnimation && (
        <CanvasOverlay
          initial={{ opacity: 1 }}
          animate={{ opacity: canvasOpacity }}
          transition={{ duration: 0.5 }}
        >
          <LegoCanvas
            config={config}
            theme={theme}
            isVisible={isVisible}
            phase={phase}
            onPhaseChange={handlePhaseChange}
            onDissolveStart={handleDissolveStart}
            width={dimensions.width}
            height={dimensions.height}
          />
        </CanvasOverlay>
      )}
    </AnimationContainer>
  );
};

export default LegoBlockAnimation;
