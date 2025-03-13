// src/components/ui/SkillBar.tsx
import React, { useEffect, useRef, useMemo } from 'react';
import styled from 'styled-components';
import { motion, useAnimation } from 'framer-motion';
import { useInView } from 'react-intersection-observer';

interface SkillBarProps {
  skill: string;
  percentage: number;
  shouldAnimate?: boolean;
  delay?: number;
}

const SkillContainer = styled.div`
  width: 100%;
  margin: 1.25rem 0;
  will-change: opacity, transform;
`;

const SkillHeader = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
`;

const SkillName = styled.h4`
  margin: 0;
  color: ${({ theme }) => theme.colors.text};
  font-weight: 600;
  font-size: 0.95rem;
`;

const SkillPercentage = styled(motion.span)`
  font-weight: 400;
  font-size: 0.9rem;
  opacity: 0.8;
  transition: opacity 0.2s ease;
`;

const ProgressBarContainer = styled.div`
  width: 100%;
  height: 8px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 4px;
  overflow: hidden;
  position: relative;
  box-shadow: inset 0 1px 3px rgba(0, 0, 0, 0.1);
`;

const ProgressBar = styled(motion.div)`
  height: 100%;
  background: linear-gradient(90deg, 
    ${({ theme }) => theme.colors.primary} 0%, 
    ${({ theme }) => theme.colors.accent || theme.colors.primary} 100%
  );
  border-radius: 4px;
  transform-origin: left;
  will-change: transform;
  box-shadow: 0 0 8px ${({ theme }) => `${theme.colors.primary}50`};
`;

const ProgressShimmer = styled(motion.div)`
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(
    90deg,
    transparent 0%,
    rgba(255, 255, 255, 0.2) 50%,
    transparent 100%
  );
  transform: translateX(-100%);
  will-change: transform, opacity;
`;

export const SkillBar: React.FC<SkillBarProps> = ({ 
  skill, 
  percentage, 
  shouldAnimate = true,
  delay = 0 
}) => {
  // Optimized animation configuration with faster timings
  const animationConfig = useMemo(() => ({
    threshold: 0.2, // Reduced threshold for earlier triggering
    duration: {
      enter: 0.7, // Reduced from 1.2
      shimmer: 1.0, // Reduced from 1.5
      exit: 0.3 // Reduced from 0.4
    },
    easing: {
      enter: [0.16, 1, 0.3, 1], // Snappier easing for faster visual appearance
      exit: [0.43, 0.13, 0.23, 0.96]
    }
  }), []);

  // Refs for animation state
  const hasAnimatedIn = useRef(false);
  
  // Animation controls
  const progressControls = useAnimation();
  const shimmerControls = useAnimation();
  const percentageControls = useAnimation();
  
  // Track element visibility with optimized threshold
  const [ref, inView] = useInView({
    threshold: animationConfig.threshold,
    triggerOnce: false
  });

  // Responsive animation handling with performance optimizations
  useEffect(() => {
    const handleAnimation = async () => {
      if (inView && shouldAnimate) {
        // Start percentage counter animation immediately
        percentageControls.start({
          opacity: 1,
          y: 0,
          transition: {
            duration: animationConfig.duration.enter * 0.4, // Even faster for text
            delay: delay * 0.3, // Reduced delay multiplier
            ease: animationConfig.easing.enter
          }
        });
        
        // Animate progress bar fill with reduced delay
        await progressControls.start({
          scaleX: percentage / 100,
          transition: {
            duration: animationConfig.duration.enter,
            delay: delay * 0.5, // Reduced delay multiplier
            ease: animationConfig.easing.enter
          }
        });
        
        hasAnimatedIn.current = true;
        
        // Add shimmer effect immediately after fill completes
        shimmerControls.start({
          x: ['0%', '100%'],
          opacity: [0, 1, 0],
          transition: {
            duration: animationConfig.duration.shimmer,
            ease: "easeInOut",
            times: [0, 0.5, 1]
          }
        });
      } else if (hasAnimatedIn.current) {
        // Faster exit animations
        shimmerControls.stop();
        
        // Quick fade for percentage
        percentageControls.start({
          opacity: 0.3,
          y: 5,
          transition: {
            duration: animationConfig.duration.exit * 0.5,
            ease: animationConfig.easing.exit
          }
        });
        
        // Quick retraction for progress bar
        progressControls.start({
          scaleX: 0,
          transition: {
            duration: animationConfig.duration.exit,
            ease: animationConfig.easing.exit
          }
        });
      }
    };
    
    // Use requestAnimationFrame for smoother animation start
    if (typeof window !== 'undefined') {
      requestAnimationFrame(() => handleAnimation());
    } else {
      handleAnimation();
    }
    
    return () => {
      // Clean up animations on unmount
      progressControls.stop();
      shimmerControls.stop();
      percentageControls.stop();
    };
  }, [
    inView, 
    shouldAnimate, 
    percentage, 
    delay, 
    progressControls, 
    shimmerControls, 
    percentageControls,
    animationConfig
  ]);

  return (
    <SkillContainer ref={ref}>
      <SkillHeader>
        <SkillName>{skill}</SkillName>
        <SkillPercentage
          initial={{ opacity: 0, y: 10 }}
          animate={percentageControls}
        >
          {percentage}%
        </SkillPercentage>
      </SkillHeader>
      
      <ProgressBarContainer>
        <ProgressBar
          initial={{ scaleX: 0 }}
          animate={progressControls}
        >
          <ProgressShimmer
            initial={{ x: '-100%', opacity: 0 }}
            animate={shimmerControls}
          />
        </ProgressBar>
      </ProgressBarContainer>
    </SkillContainer>
  );
};

export default SkillBar;
