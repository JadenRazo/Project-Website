import React, { memo, useMemo, useState, useCallback, useRef, useEffect } from 'react';
import styled, { keyframes } from 'styled-components';
import { motion, AnimatePresence, HTMLMotionProps } from 'framer-motion';
import { useNavigate } from 'react-router-dom';
import { useTheme } from '../../hooks/useTheme';
import useDeviceCapabilities from '../../hooks/useDeviceCapabilities';
import usePerformanceOptimizations from '../../hooks/usePerformanceOptimizations';

interface BioItemProps {
  $active: boolean;
  $enablePulse: boolean;
}

interface Skill {
  id: string;
  name: string;
  description: string;
  icon?: string;
  projects?: string[];
}

// Define media breakpoints for responsive design
const breakpoints = {
  mobileSm: '320px',
  mobileMd: '375px',
  mobileLg: '425px',
  tablet: '768px',
  laptop: '1024px',
  laptopLg: '1440px',
  desktop: '1920px'
};

// Media query helper functions
const media = {
  mobileSm: `@media (max-width: ${breakpoints.mobileSm})`,
  mobileMd: `@media (max-width: ${breakpoints.mobileMd})`,
  mobileLg: `@media (max-width: ${breakpoints.mobileLg})`,
  tablet: `@media (max-width: ${breakpoints.tablet})`,
  laptop: `@media (max-width: ${breakpoints.laptop})`,
  laptopLg: `@media (max-width: ${breakpoints.laptopLg})`,
  desktop: `@media (min-width: ${breakpoints.laptopLg})`,
  touch: `@media (max-width: ${breakpoints.tablet})`,
  mouse: `@media (min-width: ${breakpoints.tablet})`
};

// Dynamic animation variants based on performance settings

// Create animation objects based on performance
const createAnimationObjects = (performanceSettings: any) => {
  const { transitionSpeed, enableHoverEffects } = performanceSettings;
  
  return {
    hover: enableHoverEffects ? {
      scale: 1.01,
      transition: { duration: transitionSpeed * 0.5, ease: "easeOut" }
    } : {},
    
    tap: {
      scale: 0.98,
      transition: { duration: transitionSpeed * 0.25 }
    }
  };
};

// Keyframes for animations - conditionally applied based on performance
const createKeyframes = (performanceSettings: any) => {
  const { enableComplexAnimations, reduceMotion } = performanceSettings;
  
  // Basic pulse animation with reduced intensity if needed
  const pulseIntensity = reduceMotion ? '5px' : '10px';
  const pulse = keyframes`
    0% { box-shadow: 0 0 0 0 rgba(var(--primary-rgb), 0.7); }
    70% { box-shadow: 0 0 0 ${pulseIntensity} rgba(var(--primary-rgb), 0); }
    100% { box-shadow: 0 0 0 0 rgba(var(--primary-rgb), 0); }
  `;
  
  // Optional complex animations
  const shimmerAnimation = enableComplexAnimations ? keyframes`
    0% { background-position: 200% center; }
    100% { background-position: -200% center; }
  ` : keyframes`
    0%, 100% { background-position: 0% center; }
  `;
  
  const floatingDistance = reduceMotion ? '3px' : '10px';
  const floatingAnimation = enableComplexAnimations ? keyframes`
    0% { transform: translateY(0); }
    50% { transform: translateY(-${floatingDistance}); }
    100% { transform: translateY(0); }
  ` : keyframes`
    0%, 100% { transform: translateY(0); }
  `;
  
  return {
    pulse,
    shimmerAnimation,
    floatingAnimation
  };
};

// Define default system fonts to replace theme.fonts references
const systemFonts = {
  sans: "'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif",
  primary: "'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif",
  mono: "'SF Mono', 'Fira Code', 'Fira Mono', 'Roboto Mono', monospace"
};

interface StyledMotionProps extends HTMLMotionProps<"div"> {
  $isHovered?: boolean;
  $enableAnimations?: boolean;
}

// Enhanced styled components with optimized spacing and smooth scroll
const HeroContainer = styled.div`
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 60px 40px 40px;
  position: relative;
  overflow: hidden;
  margin-top: 60px;
  text-align: center;
  background: ${({ theme }) => theme.colors.background};

  @media (max-width: ${({ theme }) => theme.breakpoints.tablet}) {
    padding: 40px 32px 32px;
    min-height: calc(100vh - 60px);
  }
  
  @media (max-width: ${({ theme }) => theme.breakpoints.mobile}) {
    padding: 32px 20px 24px;
  }
`;

const ContentWrapper = styled(motion.div)<StyledMotionProps>`
  max-width: 1200px;
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: 32px;
  z-index: 1;
  align-items: center;
  text-align: center;

  @media (max-width: ${({ theme }) => theme.breakpoints.tablet}) {
    gap: 28px;
    align-items: center;
    justify-content: center;
  }
  
  @media (max-width: ${({ theme }) => theme.breakpoints.mobile}) {
    gap: 24px;
  }
`;


const Name = styled(motion.h2)<StyledMotionProps>`
  font-size: 4rem;
  color: ${({ theme }) => theme.colors.primary};
  margin: 0;
  line-height: 1.2;
  transform: ${({ $isHovered }) => $isHovered ? 'scale(1.05)' : 'scale(1)'};
  transition: transform 0.3s ease;

  @media (max-width: ${({ theme }) => theme.breakpoints.tablet}) {
    font-size: 3rem;
  }

  @media (max-width: ${({ theme }) => theme.breakpoints.mobile}) {
    font-size: 2.5rem;
  }
`;



// Professional dropdown container with optimized spacing
const SkillDetailContainer = styled(motion.div)`
  width: 100%;
  max-width: 700px;
  background: linear-gradient(135deg, 
    ${({ theme }) => theme.colors.surface}f5 0%, 
    ${({ theme }) => theme.colors.surface}e8 100%);
  border-radius: 20px;
  padding: 28px;
  margin: 20px auto 0;
  color: ${({ theme }) => theme.colors.text};
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  position: relative;
  overflow: hidden;
  
  /* Animated gradient border */
  &::before {
    content: '';
    position: absolute;
    inset: 0;
    padding: 2px;
    background: linear-gradient(
      45deg,
      ${({ theme }) => theme.colors.primary}80,
      transparent 30%,
      transparent 70%,
      ${({ theme }) => theme.colors.primary}80
    );
    border-radius: inherit;
    mask: linear-gradient(#fff 0 0) content-box, linear-gradient(#fff 0 0);
    mask-composite: xor;
    animation: borderGlow 3s ease-in-out infinite;
  }
  
  /* Floating particles effect */
  &::after {
    content: '';
    position: absolute;
    top: 20%;
    left: 10%;
    width: 6px;
    height: 6px;
    background: ${({ theme }) => theme.colors.primary}60;
    border-radius: 50%;
    animation: float 6s ease-in-out infinite;
    box-shadow: 
      40px 20px 0 -2px ${({ theme }) => theme.colors.primary}40,
      80px -10px 0 -3px ${({ theme }) => theme.colors.primary}30,
      120px 30px 0 -1px ${({ theme }) => theme.colors.primary}50;
  }
  
  @keyframes borderGlow {
    0%, 100% { background-position: 0% 50%; }
    50% { background-position: 100% 50%; }
  }
  
  @keyframes float {
    0%, 100% { transform: translateY(0px) rotate(0deg); }
    33% { transform: translateY(-20px) rotate(120deg); }
    66% { transform: translateY(10px) rotate(240deg); }
  }
  
  ${media.touch} {
    padding: 22px;
    margin: 16px auto 0;
    border-radius: 16px;
  }
  
  ${media.mobileSm} {
    padding: 18px;
    margin: 12px auto 0;
    border-radius: 14px;
  }
`;

// Enhanced skill title with sophisticated animations
const SkillTitle = styled(motion.h3)`
  font-size: clamp(20px, 4vw, 28px);
  font-weight: 700;
  background: linear-gradient(
    135deg,
    ${({ theme }) => theme.colors.primary} 0%,
    ${({ theme }) => theme.colors.primary}cc 50%,
    ${({ theme }) => theme.colors.primary}80 100%
  );
  background-clip: text;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  margin: 0 0 20px 0;
  display: flex;
  align-items: center;
  gap: 12px;
  position: relative;
  
  /* Animated underline */
  &::after {
    content: '';
    position: absolute;
    bottom: -8px;
    left: 0;
    height: 3px;
    width: 60px;
    background: linear-gradient(
      90deg,
      ${({ theme }) => theme.colors.primary},
      transparent
    );
    border-radius: 2px;
    animation: expandLine 0.8s ease-out 0.2s both;
  }
  
  @keyframes expandLine {
    from { width: 0; opacity: 0; }
    to { width: 60px; opacity: 1; }
  }
  
  svg {
    width: 28px;
    height: 28px;
    filter: drop-shadow(0 2px 4px rgba(0,0,0,0.1));
    animation: iconBounce 0.6s ease-out 0.4s both;
  }
  
  @keyframes iconBounce {
    0% { transform: scale(0) rotate(-180deg); }
    50% { transform: scale(1.1) rotate(-10deg); }
    100% { transform: scale(1) rotate(0deg); }
  }
  
  ${media.mobileLg} {
    margin: 0 0 16px 0;
    font-size: clamp(18px, 3.5vw, 24px);
    gap: 10px;
    
    &::after {
      width: 50px;
    }
    
    svg {
      width: 24px;
      height: 24px;
    }
  }
`;

const SkillDescription = styled(motion.p)`
  font-size: 16px;
  line-height: 1.7;
  margin-bottom: 20px;
  color: ${({ theme }) => theme.colors.text}dd;
  text-align: left;
  position: relative;
  
  ${media.mobileLg} {
    font-size: 15px;
    line-height: 1.6;
    margin-bottom: 16px;
  }
  
  ${media.mobileSm} {
    font-size: 14px;
    margin-bottom: 14px;
  }
`;

// Enhanced project list with staggered animations
const ProjectList = styled(motion.ul)`
  list-style-type: none;
  padding: 0;
  margin: 20px 0 0 0;
  
  li {
    position: relative;
    padding: 8px 0 8px 28px;
    margin-bottom: 10px;
    line-height: 1.5;
    font-size: 15px;
    color: ${({ theme }) => theme.colors.text}cc;
    border-radius: 8px;
    transition: all 0.3s ease;
    
    &::before {
      content: '';
      position: absolute;
      left: 8px;
      top: 50%;
      transform: translateY(-50%);
      width: 8px;
      height: 8px;
      background: ${({ theme }) => theme.colors.primary};
      border-radius: 50%;
      box-shadow: 0 0 0 3px ${({ theme }) => theme.colors.primary}20;
      transition: all 0.3s ease;
    }
    
    &:hover {
      color: ${({ theme }) => theme.colors.text};
      background: rgba(255, 255, 255, 0.05);
      transform: translateX(4px);
      
      &::before {
        background: ${({ theme }) => theme.colors.primary};
        box-shadow: 0 0 0 6px ${({ theme }) => theme.colors.primary}30;
        transform: translateY(-50%) scale(1.2);
      }
    }
  }
  
  ${media.mobileLg} {
    margin: 16px 0 0 0;
    
    li {
      font-size: 14px;
      padding: 6px 0 6px 24px;
      margin-bottom: 8px;
      
      &::before {
        width: 6px;
        height: 6px;
        left: 6px;
      }
    }
  }
  
  ${media.mobileSm} {
    li {
      font-size: 13px;
      padding: 5px 0 5px 20px;
      margin-bottom: 6px;
    }
  }
`;

// Enhanced CTA button with animated gradient border on hover
const CTAButton = styled(motion.button)`
  display: inline-block;
  background-color: ${({ theme }) => theme.colors.primary};
  border: 1px solid ${({ theme }) => theme.colors.primary};
  border-radius: 8px;
  color: #fff;
  font-family: ${systemFonts.mono};
  font-size: clamp(14px, 2vw, 16px);
  padding: 1rem 1.5rem;
  text-decoration: none;
  position: relative;
  overflow: hidden;
  pointer-events: auto;
  transition: all 0.3s ease;
  max-width: 200px;
  text-align: center;
  margin: 0 auto;
  font-weight: 500;
  cursor: pointer;
  
  &:after {
    content: '';
    position: absolute;
    top: 0;
    left: -100%;
    width: 100%;
    height: 100%;
    background: ${({ theme }) => `linear-gradient(
      90deg,
      transparent,
      ${theme.colors.primary}20,
      transparent
    )`};
    transition: left 0.7s ease;
    z-index: -1;
  }
  
  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
    background-color: ${({ theme }) => theme.colors.backgroundAlt};
    color: ${({ theme }) => theme.colors.primary};
    
    &:after {
      left: 200%;
    }
  }
  
  @media (max-width: ${({ theme }) => theme.breakpoints.tablet}) {
    width: 100%;
    max-width: 250px;
    padding: 1rem 1.5rem;
    margin: 0 auto;
  }
  
  @media (max-width: ${({ theme }) => theme.breakpoints.mobile}) {
    padding: 0.8rem 1.2rem;
    font-size: 13px;
  }
`;

// Enhanced skill data with icon references
const SKILLS: Record<string, Skill> = {
  'ui': {
    id: 'ui',
    name: 'UI Designer',
    description: 'Creating intuitive and visually appealing user interfaces that prioritize user experience and accessibility. Proficient in design principles, color theory, and responsive layouts.',
    icon: 'design',
    projects: ['Portfolio Website Redesign', 'Discord Bot styling', 'Dashboard Interface for Analytics Platform']
  },
  'api': {
    id: 'api',
    name: 'API Coding',
    description: 'Building robust and secure APIs that connect front-end applications to back-end services. Experience with RESTful design principles, authentication, and data handling.',
    icon: 'code',
    projects: ['Developing Microservices Architecture', 'Discord Python Bot API', 'ChatGPT, Claude, Gemini API Integration']
  },
  'db': {
    id: 'db',
    name: 'Database Management',
    description: 'Designing efficient database structures and managing data storage solutions. Skilled in SQL and NoSQL databases, query optimization, and data security practices.',
    icon: 'database',
    projects: ['Diverse SQL Querys', 'Asyncpg in Python', 'Turning Data into Insights']
  }
};

// SVG icons for skills
const SkillIcon = memo(({ type }: { type: string }) => {
  switch (type) {
    case 'design':
      return (
        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
          <path d="M12 20h9"></path>
          <path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"></path>
        </svg>
      );
    case 'code':
      return (
        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
          <polyline points="16 18 22 12 16 6"></polyline>
          <polyline points="8 6 2 12 8 18"></polyline>
        </svg>
      );
    case 'database':
      return (
        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
          <ellipse cx="12" cy="5" rx="9" ry="3"></ellipse>
          <path d="M21 12c0 1.66-4 3-9 3s-9-1.34-9-3"></path>
          <path d="M3 5v14c0 1.66 4 3 9 3s9-1.34 9-3V5"></path>
        </svg>
      );
    default:
      return null;
  }
});
SkillIcon.displayName = 'SkillIcon';

// Interface for the AnimatedBioItem component props
interface AnimatedBioItemProps extends BioItemProps {
  text: string;
  onClick: () => void;
  $isActive: boolean;
  performanceSettings: any;
  animationVariants: any;
  animationObjects: any;
}

// Update the AnimatedBioItem component
const AnimatedBioItem = memo(({ 
  text, 
  onClick, 
  $isActive, 
  performanceSettings,
  animationVariants,
  animationObjects
}: AnimatedBioItemProps) => {
  
  return (
    <motion.button
      variants={animationVariants.item}
      whileHover={animationObjects.hover}
      whileTap={animationObjects.tap}
      onClick={onClick}
      className={$isActive ? 'active' : ''}
      style={{
        fontSize: 'clamp(16px, 3vw, 20px)',
        color: $isActive ? '#fff' : 'var(--colors-text)',
        padding: '8px 16px',
        borderRadius: '8px',
        backgroundColor: $isActive ? 'var(--colors-primary)' : 'rgba(255, 255, 255, 0.15)',
        border: '1px solid var(--colors-primary)',
        position: 'relative',
        cursor: 'pointer',
        fontFamily: 'inherit',
        boxShadow: '0 2px 10px rgba(0, 0, 0, 0.1)',
        overflow: 'hidden',
        transition: 'all 0.3s ease',
        fontWeight: 500
      }}
    >
      {text}
    </motion.button>
  );
});
AnimatedBioItem.displayName = 'AnimatedBioItem';

// Professional SkillDetail component with sophisticated animations
const SkillDetail = memo(({ 
  skill, 
  performanceSettings, 
  animationVariants,
  visibleProjects,
  projectRefs 
}: { 
  skill: Skill; 
  performanceSettings: any;
  animationVariants: any;
  visibleProjects: Set<number>;
  projectRefs: React.MutableRefObject<(HTMLLIElement | null)[]>;
}) => (
  <SkillDetailContainer
    initial={{
      opacity: 0,
      y: 30,
      scale: 0.95,
      rotateX: -15
    }}
    animate={{
      opacity: 1,
      y: 0,
      scale: 1,
      rotateX: 0
    }}
    exit={{
      opacity: 0,
      y: -20,
      scale: 0.98,
      rotateX: 10
    }}
    transition={{
      duration: 0.6,
      ease: [0.25, 0.46, 0.45, 0.94],
      opacity: { duration: 0.4 },
      scale: { duration: 0.5, delay: 0.1 }
    }}
    layout
  >
    <SkillTitle
      initial={{ opacity: 0, x: -30 }}
      animate={{ opacity: 1, x: 0 }}
      transition={{ delay: 0.2, duration: 0.5, ease: "easeOut" }}
    >
      {skill.icon && <SkillIcon type={skill.icon} />}
      {skill.name}
    </SkillTitle>
    
    <SkillDescription
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ delay: 0.4, duration: 0.5, ease: "easeOut" }}
    >
      {skill.description}
    </SkillDescription>
    
    {skill.projects && skill.projects.length > 0 && (
      <motion.div
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ delay: 0.6, duration: 0.4 }}
      >
        <motion.div 
          style={{ 
            fontSize: '16px', 
            fontWeight: '600', 
            color: 'var(--colors-primary)', 
            marginBottom: '12px',
            display: 'flex',
            alignItems: 'center',
            gap: '8px'
          }}
          initial={{ opacity: 0, x: -20 }}
          animate={{ opacity: 1, x: 0 }}
          transition={{ delay: 0.7, duration: 0.4 }}
        >
          <motion.span
            initial={{ scale: 0 }}
            animate={{ scale: 1 }}
            transition={{ delay: 0.8, duration: 0.3, type: "spring" }}
          >
            ✨
          </motion.span>
          Related Projects:
        </motion.div>
        
        <ProjectList
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.8, duration: 0.3 }}
        >
          {skill.projects.map((project, index) => (
            <motion.li 
              key={index}
              ref={(el) => {
                projectRefs.current[index] = el;
              }}
              data-index={index}
              initial={{ 
                opacity: 0, 
                x: -50, 
                y: 20,
                scale: 0.9,
                filter: 'blur(6px)',
                rotateX: -15
              }}
              animate={visibleProjects.has(index) ? {
                opacity: 1,
                x: 0,
                y: 0,
                scale: 1,
                filter: 'blur(0px)',
                rotateX: 0
              } : {
                opacity: 0,
                x: -50,
                y: 20,
                scale: 0.9,
                filter: 'blur(6px)',
                rotateX: -15
              }}
              transition={{
                duration: 0.7,
                ease: [0.165, 0.84, 0.44, 1], // Enhanced easing for smoother animation
                delay: index * 0.1, // Reduced stagger for faster sequence
                type: "tween",
                // Individual property transitions for more control
                opacity: { duration: 0.5, delay: index * 0.1 },
                x: { duration: 0.7, delay: index * 0.1 },
                y: { duration: 0.6, delay: index * 0.1 + 0.05 },
                scale: { duration: 0.6, delay: index * 0.1 + 0.1 },
                filter: { duration: 0.4, delay: index * 0.1 + 0.15 },
                rotateX: { duration: 0.5, delay: index * 0.1 + 0.05 }
              }}
              whileHover={{
                scale: 1.03,
                x: 12,
                y: -2,
                transition: { 
                  duration: 0.2,
                  ease: "easeOut"
                }
              }}
              style={{
                // Add subtle transform origin for better hover effects
                transformOrigin: "left center"
              }}
            >
              {project}
            </motion.li>
          ))}
        </ProjectList>
      </motion.div>
    )}
  </SkillDetailContainer>
));
SkillDetail.displayName = 'SkillDetail';

// 2+1 button layout with optimized spacing
const Bio = styled(motion.div)<StyledMotionProps>`
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin: 0;
  align-items: center;
  width: 100%;
  max-width: 700px;
  
  @media (max-width: ${({ theme }) => theme.breakpoints.tablet}) {
    gap: 14px;
    max-width: 600px;
  }
  
  @media (max-width: ${({ theme }) => theme.breakpoints.mobile}) {
    gap: 12px;
    max-width: 100%;
    padding: 0 16px;
  }
`;

// Container for the top two buttons
const TopButtonRow = styled.div`
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
  width: 100%;
  
  @media (max-width: ${({ theme }) => theme.breakpoints.tablet}) {
    gap: 14px;
  }
  
  @media (max-width: ${({ theme }) => theme.breakpoints.mobile}) {
    grid-template-columns: 1fr;
    gap: 12px;
  }
`;

// Container for the bottom single button
const BottomButtonRow = styled.div`
  display: flex;
  justify-content: center;
  width: 100%;
`;

const BioItem = styled(motion.button)<BioItemProps>`
  /* Responsive sizing for different positions */
  width: 100%;
  min-height: 64px;
  padding: 18px 20px;
  
  /* Typography - responsive but consistent */
  font-size: clamp(13px, 2vw, 15px);
  font-weight: 500;
  font-family: inherit;
  line-height: 1.3;
  text-align: center;
  
  /* Perfect text centering and handling */
  display: flex;
  align-items: center;
  justify-content: center;
  white-space: normal;
  word-break: break-word;
  hyphens: auto;
  
  /* Visual design */
  color: ${({ theme, $active }) => $active ? '#fff' : theme.colors.text};
  background-color: ${({ theme, $active }) => $active 
    ? theme.colors.primary 
    : 'rgba(255, 255, 255, 0.08)'};
  border: 2px solid ${({ theme, $active }) => $active 
    ? theme.colors.primary 
    : 'rgba(255, 255, 255, 0.12)'};
  border-radius: 14px;
  backdrop-filter: blur(10px);
  
  /* Sleek interactions */
  cursor: pointer;
  pointer-events: auto;
  position: relative;
  overflow: hidden;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  
  /* Subtle gradient overlay */
  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: -100%;
    width: 100%;
    height: 100%;
    background: linear-gradient(
      90deg,
      transparent,
      rgba(255, 255, 255, 0.1),
      transparent
    );
    transition: left 0.6s ease;
    z-index: 0;
  }
  
  /* Text stays above overlay */
  & > * {
    position: relative;
    z-index: 1;
  }

  &:hover {
    background-color: ${({ theme, $active }) => $active 
      ? theme.colors.primary 
      : 'rgba(255, 255, 255, 0.12)'};
    border-color: ${({ theme }) => theme.colors.primary};
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.2);
    transform: translateY(-2px);
    
    &::before {
      left: 100%;
    }
  }

  &:active {
    transform: translateY(0) scale(0.98);
    transition: all 0.15s ease;
  }

  /* Special styling for bottom button */
  &.bottom-button {
    max-width: 340px;
  }

  /* Responsive adjustments */
  @media (max-width: ${({ theme }) => theme.breakpoints.tablet}) {
    min-height: 60px;
    padding: 16px 18px;
    border-radius: 12px;
    font-size: clamp(12px, 2.2vw, 14px);
    
    &.bottom-button {
      max-width: 320px;
    }
  }
  
  @media (max-width: ${({ theme }) => theme.breakpoints.mobile}) {
    min-height: 56px;
    padding: 14px 16px;
    border-radius: 10px;
    font-size: clamp(12px, 3vw, 14px);
    
    &.bottom-button {
      max-width: 100%;
    }
  }
`;

const TypewriterContainer = styled.div`
  position: relative;
  display: block;
  width: 100%;
  text-align: center;
  word-break: normal;
  white-space: normal;
  letter-spacing: 0.05em;

  @media (max-width: ${({ theme }) => theme.breakpoints.mobile}) {
    span.word {
      font-size: 2.2rem;
    }
    
    span:nth-child(2) {
      &::after {
        content: '';
        display: block;
        width: 100%;
        height: 0;
      }
    }
  }
  
  @media (max-width: 375px) {
    span.word {
      font-size: 1.8rem;
    }
  }
`;

const WordWrapper = styled.span`
  display: inline-block;
  white-space: nowrap;
  margin-right: 12px;
  margin-bottom: 5px;
`;

const TypewriterCharacter = styled(motion.span)<{ $isSlash?: boolean }>`
  display: inline-block;
  color: ${props => props.$isSlash ? props.theme.colors.primary : 'inherit'};
  position: relative;
  letter-spacing: 0.03em;
`;

const TypewriterCursor = styled(motion.div)`
  position: absolute;
  right: -4px;
  top: 0;
  width: 3px;
  height: 100%;
  background-color: currentColor;
  border-radius: 1px;
`;

// Global smooth scroll enhancement
const GlobalSmoothScrollStyles = styled.div`
  /* Enhanced smooth scrolling for the entire page */
  html {
    scroll-behavior: smooth;
    scroll-padding-top: 80px;
  }
  
  /* Force smooth scrolling for all scroll operations */
  * {
    scroll-behavior: smooth;
  }
  
  /* Custom scroll timing for enhanced smoothness */
  html, body {
    scroll-behavior: smooth;
    scrollbar-width: thin;
  }
  
  /* Webkit browsers smooth scroll enhancement */
  ::-webkit-scrollbar {
    width: 8px;
  }
  
  ::-webkit-scrollbar-track {
    background: rgba(0,0,0,0.1);
  }
  
  ::-webkit-scrollbar-thumb {
    background: rgba(0,0,0,0.3);
    border-radius: 4px;
  }
  
  ::-webkit-scrollbar-thumb:hover {
    background: rgba(0,0,0,0.5);
  }
`;

// Main Hero component refactored to use the new hooks
export const Hero: React.FC = () => {
  // Get theme from context
  const { theme } = useTheme();
  
  // Initialize navigation
  const navigate = useNavigate();
  
  // Initialize component state
  const [isNameHovered, setIsNameHovered] = useState(false);
  const [activeSkill, setActiveSkill] = useState<string | null>(null);
  const [visibleProjects, setVisibleProjects] = useState<Set<number>>(new Set());
  const [initialScrollPosition, setInitialScrollPosition] = useState<number>(0);
  
  // Create refs for touch interactions and section visibility
  const bioRef = useRef<HTMLDivElement>(null);
  const skillDetailRef = useRef<HTMLDivElement>(null);
  const projectRefs = useRef<(HTMLLIElement | null)[]>([]);
  
  
  // Use our new hooks for device optimization
  const deviceCapabilities = useDeviceCapabilities();
  const { performanceSettings } = usePerformanceOptimizations();
  
  // Enhanced intersection observer for seamless project animations
  useEffect(() => {
    if (!activeSkill) {
      setVisibleProjects(new Set());
      return;
    }

    // Capture refs at the beginning of the effect to avoid stale closure
    const currentRefs = projectRefs.current;

    const observer = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting) {
            const index = parseInt(entry.target.getAttribute('data-index') || '0');
            setVisibleProjects(prev => new Set([...prev, index]));
          } else {
            // Optional: Remove from visible set when out of view for re-animation
            const index = parseInt(entry.target.getAttribute('data-index') || '0');
            setVisibleProjects(prev => {
              const newSet = new Set(prev);
              newSet.delete(index);
              return newSet;
            });
          }
        });
      },
      {
        threshold: 0.1, // Trigger when 10% of the project is visible (earlier)
        rootMargin: '0px 0px -20px 0px' // Start animation earlier
      }
    );

    // Small delay to ensure refs are set
    const setupObserver = () => {
      currentRefs.forEach((ref) => {
        if (ref) observer.observe(ref);
      });
    };

    // Setup observer after dropdown animation starts
    const timeoutId = setTimeout(setupObserver, 100);
    
    return () => {
      clearTimeout(timeoutId);
      currentRefs.forEach((ref) => {
        if (ref) observer.unobserve(ref);
      });
    };
  }, [activeSkill]); // Re-run when activeSkill changes

  // Auto-trigger animation for projects if they're already in view
  useEffect(() => {
    if (activeSkill) {
      // Trigger animations for projects in view after a delay
      setTimeout(() => {
        projectRefs.current.forEach((ref, index) => {
          if (ref) {
            const rect = ref.getBoundingClientRect();
            const isInView = rect.top < window.innerHeight && rect.bottom > 0;
            if (isInView) {
              setVisibleProjects(prev => new Set([...prev, index]));
            }
          }
        });
      }, 300); // Wait for dropdown to be visible
    }
  }, [activeSkill]);
  
  
  // Memoize animation settings based on performance
  const animationVariants = useMemo(
    () => ({
      container: {
        hidden: { opacity: 0 },
        visible: {
          opacity: 1,
          transition: {
            staggerChildren: 0.1,
            delayChildren: 0.2,
          },
        },
      },
      item: {
        hidden: { opacity: 0, y: 20 },
        visible: {
          opacity: 1,
          y: 0,
          transition: {
            type: "spring",
            stiffness: 100,
            damping: 10,
          },
        },
      },
    }),
    []
  );
  
  const animationObjects = useMemo(
    () => createAnimationObjects(performanceSettings),
    [performanceSettings]
  );
  
  // Memoize animation keyframes and add them to the theme
  useEffect(() => {
    const keyframesObj = createKeyframes(performanceSettings);
    
    // Add keyframes to the theme for use in styled components
    if (theme) {
      Object.entries(keyframesObj).forEach(([key, value]) => {
        (theme as any)[key] = value;
      });
    }
  }, [performanceSettings, theme]);
  
  // Determine if we're on a touch device for interaction optimizations
  const isTouchDevice = deviceCapabilities.isTouchDevice;
  
  
  
  
  // Clean bio items with skill mapping
  const bioItems = useMemo(() => [
    { label: 'UI Designer', skillId: 'ui' },
    { label: 'API Development & Integration', skillId: 'api' },
    { label: 'Database Management', skillId: 'db' }
  ], []);

  const nameText = "Hi there, I'm Jaden Razo/";
  const nameCharacters = nameText.split('');
  const [typedCount, setTypedCount] = useState(0);
  const [showCursor, setShowCursor] = useState(true);
  const [isTypingComplete, setIsTypingComplete] = useState(false);

  // Update typing effect with improved timing and consistent speed
  useEffect(() => {
    if (typedCount < nameCharacters.length) {
      const charDelay = Math.random() * 10 + 40; // Approximately 60 WPM (40-50ms per character)
      const timeout = setTimeout(() => {
        setTypedCount(prev => prev + 1);
      }, charDelay); 
      return () => clearTimeout(timeout);
    } else if (!isTypingComplete) {
      setIsTypingComplete(true);
      setTimeout(() => {
        setShowCursor(false); // Hide cursor after typing is complete with slight delay
      }, 800);
    }
  }, [typedCount, nameCharacters.length, isTypingComplete]);

  // Remove cursor blink effect since we're using overtype cursor
  useEffect(() => {
    if (isTypingComplete) {
      setShowCursor(false);
    }
  }, [isTypingComplete]);

  // Handlers for name hover - only enable effects on non-touch devices
  const handleNameMouseEnter = useCallback(() => {
    if (!isTouchDevice) {
      setIsNameHovered(true);
    }
  }, [isTouchDevice]);

  const handleNameMouseLeave = useCallback(() => {
    setIsNameHovered(false);
  }, []);

  // FIXED: Simple and clear scroll logic
  const handleSkillClick = useCallback((skillId: string) => {
    const isClosing = activeSkill === skillId;
    const isSwitching = activeSkill && activeSkill !== skillId; 
    const isFirstOpen = !activeSkill && initialScrollPosition === 0;
    
    
    if (isClosing) {
      // CLOSING: Scroll back up to buttons area and RESET state
      setActiveSkill(null);
      setVisibleProjects(new Set());
      
      // RESTORE WORKING SCROLL METHODS FOR UPWARD SCROLL
      if (initialScrollPosition !== undefined && initialScrollPosition !== 0) {
        const scrollBackPosition = initialScrollPosition; // Save before reset
        
        try {
          // Ensure smooth scroll CSS is applied
          document.documentElement.style.scrollBehavior = 'smooth';
          document.body.style.scrollBehavior = 'smooth';
          
          // Method 1: Standard scrollTo (WORKING VERSION)
          window.scrollTo({
            top: scrollBackPosition,
            left: 0,
            behavior: 'smooth'
          });
          
          // Method 2: Backup with document.documentElement (WORKING VERSION)
          setTimeout(() => {
            document.documentElement.scrollTo({
              top: scrollBackPosition,
              left: 0,
              behavior: 'smooth'
            });
          }, 100);
          
          // Method 3: Force with document.body as fallback (WORKING VERSION)
          setTimeout(() => {
            if (Math.abs(window.pageYOffset - scrollBackPosition) > 30) {
              document.body.scrollTop = scrollBackPosition;
              document.documentElement.scrollTop = scrollBackPosition;
            }
          }, 200);
          
        } catch (error) {
          console.error('Upward scroll failed:', error);
          // Emergency fallback
          document.body.scrollTop = scrollBackPosition;
          document.documentElement.scrollTop = scrollBackPosition;
        }
        
        // CRITICAL: Reset scroll position after closing so next button is treated as fresh open
        setTimeout(() => {
          setInitialScrollPosition(0);
        }, 500); // Wait for scroll to complete
      }
      return;
    }
    
    if (isSwitching) {
      // SWITCHING: Just change content, NO SCROLLING
      setActiveSkill(skillId);
      setVisibleProjects(new Set());
      return;
    }
    
    if (isFirstOpen) {
      // FIRST OPEN: Save position and scroll down
      const currentScrollY = window.pageYOffset || document.documentElement.scrollTop;
      setInitialScrollPosition(currentScrollY);
      
      setActiveSkill(skillId);
      setVisibleProjects(new Set());
      
      // RESTORE WORKING SCROLL METHODS
      setTimeout(() => {
        const scrollToDropdown = () => {
          const bioElement = bioRef.current;
          if (!bioElement) return;
          
          const bioRect = bioElement.getBoundingClientRect();
          const currentScrollY = window.pageYOffset || document.documentElement.scrollTop;
          const bioBottom = bioRect.bottom + currentScrollY;
          const dropdownHeight = 200;
          const targetY = bioBottom + dropdownHeight - (window.innerHeight * 0.6);
          
          
          // RESTORE THE PROVEN WORKING METHODS
          try {
            // Ensure smooth scroll CSS is applied
            document.documentElement.style.scrollBehavior = 'smooth';
            document.body.style.scrollBehavior = 'smooth';
            
            // Method 1: Standard scrollTo (WORKING VERSION)
            window.scrollTo({
              top: targetY,
              left: 0,
              behavior: 'smooth'
            });
            
            // Method 2: Backup with document.documentElement (WORKING VERSION)
            setTimeout(() => {
              document.documentElement.scrollTo({
                top: targetY,
                left: 0,
                behavior: 'smooth'
              });
            }, 100);
            
            // Method 3: Force with document.body as fallback (WORKING VERSION)
            setTimeout(() => {
              if (Math.abs(window.pageYOffset - targetY) > 30) {
                document.body.scrollTop = targetY;
                document.documentElement.scrollTop = targetY;
              }
            }, 200);
            
          } catch (error) {
            console.error('Scroll failed:', error);
            // Emergency fallback
            document.body.scrollTop = targetY;
            document.documentElement.scrollTop = targetY;
          }
        };
        
        scrollToDropdown();
      }, 100);
      
      return;
    }
    
  }, [activeSkill, initialScrollPosition]);

  // Handler for portfolio navigation
  const handlePortfolioClick = useCallback(() => {
    navigate('/projects');
  }, [navigate]);

  return (
    <>
      <GlobalSmoothScrollStyles />
      <HeroContainer>
        <ContentWrapper
          variants={animationVariants.container}
          initial="hidden"
          animate="visible"
        >
        <Name
          variants={animationVariants.item}
          $isHovered={isNameHovered}
          $enableAnimations={performanceSettings.enableComplexAnimations}
          onMouseEnter={handleNameMouseEnter}
          onMouseLeave={handleNameMouseLeave}
        >
          <TypewriterContainer>
            {/* Process typed text word by word with explicit space handling */}
            {nameText.substring(0, typedCount).split(' ').map((word, wordIndex) => {
              // Skip empty words
              if (word === '') return null;
              
              return (
                <WordWrapper 
                  key={`word-${wordIndex}`}
                  className="word"
                >
                  {word.split('').map((char, charIndex) => (
                    <TypewriterCharacter
                      key={`${wordIndex}-${charIndex}`}
                      $isSlash={char === '/'}
                      initial={{ opacity: 0, y: -10 }}
                      animate={{ opacity: 1, y: 0 }}
                      transition={{ duration: 0.1 }}
                    >
                      {char}
                    </TypewriterCharacter>
                  ))}
                </WordWrapper>
              );
            })}
            {showCursor && !isTypingComplete && (
              <TypewriterCursor
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                transition={{ duration: 0.1 }}
              />
            )}
          </TypewriterContainer>
        </Name>

        <Bio 
          ref={bioRef}
          variants={animationVariants.container}
        >
          {/* Top row: First two buttons side by side */}
          <TopButtonRow>
            {bioItems.slice(0, 2).map((item) => {
              const isActive = activeSkill === item.skillId;
              return (
                <BioItem
                  key={item.skillId}
                  onClick={() => handleSkillClick(item.skillId)}
                  $active={isActive}
                  $enablePulse={!performanceSettings?.reduceMotion}
                  animate={isActive ? {
                    scale: 1.02,
                    backgroundColor: theme.colors.primary,
                    color: '#fff',
                    boxShadow: '0 12px 40px rgba(0,0,0,0.3)',
                    borderColor: theme.colors.primary,
                  } : {
                    scale: 1,
                    backgroundColor: 'rgba(255, 255, 255, 0.08)',
                    color: theme.colors.text,
                    boxShadow: '0 4px 16px rgba(0,0,0,0.1)',
                    borderColor: 'rgba(255, 255, 255, 0.12)',
                  }}
                  transition={{ 
                    duration: 0.4, 
                    ease: [0.4, 0, 0.2, 1],
                    type: "tween"
                  }}
                >
                  {item.label}
                </BioItem>
              );
            })}
          </TopButtonRow>

          {/* Bottom row: Third button centered */}
          <BottomButtonRow>
            {bioItems.slice(2).map((item) => {
              const isActive = activeSkill === item.skillId;
              return (
                <BioItem
                  key={item.skillId}
                  className="bottom-button"
                  onClick={() => handleSkillClick(item.skillId)}
                  $active={isActive}
                  $enablePulse={!performanceSettings?.reduceMotion}
                  animate={isActive ? {
                    scale: 1.02,
                    backgroundColor: theme.colors.primary,
                    color: '#fff',
                    boxShadow: '0 12px 40px rgba(0,0,0,0.3)',
                    borderColor: theme.colors.primary,
                  } : {
                    scale: 1,
                    backgroundColor: 'rgba(255, 255, 255, 0.08)',
                    color: theme.colors.text,
                    boxShadow: '0 4px 16px rgba(0,0,0,0.1)',
                    borderColor: 'rgba(255, 255, 255, 0.12)',
                  }}
                  transition={{ 
                    duration: 0.4, 
                    ease: [0.4, 0, 0.2, 1],
                    type: "tween"
                  }}
                >
                  {item.label}
                </BioItem>
              );
            })}
          </BottomButtonRow>
        </Bio>

        {/* Animated skill details */}
        <AnimatePresence mode="wait">
          {activeSkill && SKILLS[activeSkill] && (
            <div ref={skillDetailRef}>
              <SkillDetail 
                skill={SKILLS[activeSkill]} 
                performanceSettings={performanceSettings}
                animationVariants={animationVariants}
                visibleProjects={visibleProjects}
                projectRefs={projectRefs}
              />
            </div>
          )}
        </AnimatePresence>

        <CTAButton
          onClick={handlePortfolioClick}
          variants={animationVariants.item}
          whileHover={!isTouchDevice ? animationObjects.hover : undefined}
          whileTap={!isTouchDevice ? animationObjects.tap : undefined}
        >
          View My Portfolio
        </CTAButton>
        </ContentWrapper>
      </HeroContainer>
    </>
  );
};

export default Hero;
