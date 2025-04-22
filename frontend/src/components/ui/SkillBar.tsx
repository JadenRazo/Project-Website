// src/components/ui/SkillBar.tsx
import React, { useRef, useEffect } from 'react';
import styled from 'styled-components';
import { motion, useAnimationControls } from 'framer-motion';

interface SkillBarProps {
  skill: string;
  percentage: number;
  shouldAnimate?: boolean;
  delay?: number;
}

// More compact container with reduced margins
const SkillBarContainer = styled.div`
  margin-bottom: 1rem;
  width: 100%;
`;

// Skill name and percentage display
const SkillInfo = styled.div`
  display: flex;
  justify-content: space-between;
  margin-bottom: 0.4rem;
  font-size: 0.9rem;
`;

const SkillName = styled.span`
  color: ${({ theme }) => theme.colors.text};
  font-weight: 500;
`;

const SkillPercentage = styled.span`
  color: ${({ theme }) => theme.colors.primary};
  font-weight: 600;
`;

// Bar background with reduced height
const BarBackground = styled.div`
  background-color: ${({ theme }) => theme.colors.backgroundAlt || 'rgba(255, 255, 255, 0.1)'};
  height: 5px;
  border-radius: 3px;
  overflow: hidden;
`;

// Fill bar with gradient options
const BarFill = styled(motion.div)<{ percentage: number }>`
  height: 100%;
  width: ${props => props.percentage}%;
  background: ${({ theme }) => theme.colors.primary};
  background: linear-gradient(90deg, 
    ${({ theme }) => theme.colors.primary} 0%, 
    ${({ theme }) => theme.colors.accent || theme.colors.primary} 100%);
  border-radius: 3px;
  transform-origin: left center;
`;

const SkillBar: React.FC<SkillBarProps> = ({ 
  skill, 
  percentage, 
  shouldAnimate = true,
  delay = 0 
}) => {
  const controls = useAnimationControls();
  const hasAnimated = useRef(false);
  
  useEffect(() => {
    if (shouldAnimate) {
      // Only animate if we haven't animated before or if it's explicitly requested
      if (!hasAnimated.current) {
        controls.set({ scaleX: 0 });
        controls.start({
          scaleX: 1,
          transition: {
            duration: 0.8,
            delay: delay,
            ease: [0.22, 1, 0.36, 1]
          }
        });
        hasAnimated.current = true;
      }
    }
  }, [controls, shouldAnimate, delay]);
  
  return (
    <SkillBarContainer>
      <SkillInfo>
        <SkillName>{skill}</SkillName>
        <SkillPercentage>{percentage}%</SkillPercentage>
      </SkillInfo>
      <BarBackground>
        <BarFill
          percentage={percentage}
          initial={{ scaleX: 0 }}
          animate={controls}
        />
      </BarBackground>
    </SkillBarContainer>
  );
};

export default SkillBar;
