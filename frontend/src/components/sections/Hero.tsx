import React, { memo, useMemo, useState, useCallback, useRef, useEffect, forwardRef } from 'react';
import styled, { css, keyframes } from 'styled-components';
import { motion, Variants, HTMLMotionProps, AnimatePresence, useAnimation } from 'framer-motion';
import { useTheme } from '../../contexts/ThemeContext';
import useDeviceCapabilities from '../../hooks/useDeviceCapabilities';
import useTouchInteractions from '../../hooks/useTouchInteractions';
import usePerformanceOptimizations from '../../hooks/usePerformanceOptimizations';
import { useInView } from 'react-intersection-observer';
import { CreativeShaderBackground } from '../animations/CreativeShaderBackground';

// Types
interface AnimatedElementProps extends HTMLMotionProps<"div"> {
  isInteractive?: boolean;
}

interface BioItemProps extends AnimatedElementProps {
  text: string;
  onClick: () => void;
  isActive: boolean;
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
const createAnimationVariants = (performanceSettings: any) => {
  const { transitionSpeed, staggerDelay, reduceMotion } = performanceSettings;
  
  return {
    container: {
      hidden: { opacity: 0 },
      visible: {
        opacity: 1,
        transition: {
          staggerChildren: reduceMotion ? staggerDelay / 2 : staggerDelay,
          when: "beforeChildren",
        },
      },
    },
    item: {
      hidden: { opacity: 0, y: reduceMotion ? 5 : 15 },
      visible: {
        opacity: 1,
        y: 0,
        transition: {
          type: "tween",
          duration: transitionSpeed,
          ease: [0.25, 0.1, 0.25, 1.0],
        },
      },
    },
    skillDetail: {
      initial: { 
        opacity: 0, 
        height: 0,
        y: reduceMotion ? -5 : -20
      },
      animate: { 
        opacity: 1, 
        height: 'auto',
        y: 0,
        transition: {
          opacity: { duration: transitionSpeed * 0.75 },
          height: { duration: transitionSpeed },
          y: { duration: transitionSpeed * 0.75, ease: "easeOut" }
        }
      },
      exit: { 
        opacity: 0, 
        height: 0,
        y: reduceMotion ? -5 : -10,
        transition: {
          opacity: { duration: transitionSpeed * 0.5 },
          height: { duration: transitionSpeed * 0.75 },
          y: { duration: transitionSpeed * 0.5 }
        }
      }
    },
    parallax: {
      initial: { y: 0 },
      animate: (custom: number) => ({
        y: custom,
        transition: {
          type: reduceMotion ? "tween" : "spring",
          stiffness: 10,
          damping: 25,
          mass: 1
        }
      })
    }
  };
};

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

// Enhanced styled components with responsive and performance optimizations
const HeroContainer = styled.div`
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: ${({ theme }) => theme.spacing.xl};
  position: relative;
  overflow: hidden;
  margin-top: 60px;
  text-align: center;
  background: ${({ theme }) => theme.colors.background};

  @media (max-width: ${({ theme }) => theme.breakpoints.tablet}) {
    padding: ${({ theme }) => theme.spacing.lg};
    min-height: calc(100vh - 60px);
    display: flex;
    align-items: center;
    justify-content: center;
  }
`;

const ContentWrapper = styled(motion.div)<StyledMotionProps>`
  max-width: 1200px;
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: ${({ theme }) => theme.spacing.xl};
  z-index: 1;
  align-items: center;
  text-align: center;

  @media (max-width: ${({ theme }) => theme.breakpoints.tablet}) {
    align-items: center;
    justify-content: center;
  }
`;

const Greeting = styled(motion.h1)`
  font-size: 2rem;
  color: ${({ theme }) => theme.colors.text};
  margin: 0;
  line-height: 1.2;
  transition: all 0.3s ease;
  cursor: default;

  &:hover {
    font-size: 2.2rem;
    color: ${({ theme }) => theme.colors.primary};
  }

  @media (max-width: ${({ theme }) => theme.breakpoints.tablet}) {
    font-size: 1.75rem;
    &:hover {
      font-size: 1.95rem;
    }
  }

  @media (max-width: ${({ theme }) => theme.breakpoints.mobile}) {
    font-size: 1.5rem;
    &:hover {
      font-size: 1.7rem;
    }
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

const Title = styled(motion.h3)`
  font-size: 2rem;
  color: ${({ theme }) => theme.colors.text};
  margin: 0;
  line-height: 1.2;

  @media (max-width: ${({ theme }) => theme.breakpoints.tablet}) {
    font-size: 1.75rem;
  }

  @media (max-width: ${({ theme }) => theme.breakpoints.mobile}) {
    font-size: 1.5rem;
  }
`;

const Description = styled(motion.p)`
  font-size: 1.25rem;
  color: ${({ theme }) => theme.colors.text};
  max-width: 600px;
  line-height: 1.6;
  margin: 0;

  @media (max-width: ${({ theme }) => theme.breakpoints.tablet}) {
    font-size: 1.1rem;
  }

  @media (max-width: ${({ theme }) => theme.breakpoints.mobile}) {
    font-size: 1rem;
  }
`;

// Enhanced skill detail container with more visual feedback
const SkillDetailContainer = styled(motion.div)`
  width: 100%;
  background: ${({ theme }) => theme.colors.surface};
  border-radius: 12px;
  padding: 20px;
  margin-top: 10px;
  margin-bottom: 20px;
  color: ${({ theme }) => theme.colors.text};
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
  border: 1px solid ${({ theme }) => theme.colors.primary}20;
  overflow: hidden;
  position: relative;
  
  &:before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    width: 4px;
    height: 100%;
    background: ${({ theme }) => theme.colors.primary};
    border-radius: 4px 0 0 4px;
  }
  
  ${media.touch} {
    padding: 15px 15px 15px 20px;
    text-align: left;
  }
  
  ${media.mobileSm} {
    padding: 12px 12px 12px 16px;
    margin-top: 5px;
    margin-bottom: 15px;
  }
`;

// Enhanced skill title with icon animation
const SkillTitle = styled.h3`
  font-size: clamp(18px, 4vw, 24px);
  color: ${({ theme }) => theme.colors.primary};
  margin: 0 0 15px 0;
  display: flex;
  align-items: center;
  gap: 10px;
  
  svg {
    width: 24px;
    height: 24px;
  }
  
  ${media.mobileLg} {
    margin: 0 0 10px 0;
    font-size: clamp(16px, 3.5vw, 20px);
    
    svg {
      width: 20px;
      height: 20px;
    }
  }
`;

const SkillDescription = styled.p`
  font-size: 16px;
  line-height: 1.6;
  margin-bottom: 15px;
  
  ${media.mobileLg} {
    font-size: 14px;
    line-height: 1.5;
    margin-bottom: 10px;
  }
`;

// Enhanced project list with improved mobile styling
const ProjectList = styled.ul`
  list-style-type: none;
  padding: 0;
  margin: 15px 0 0 0;
  
  li {
    position: relative;
    padding-left: 20px;
    margin-bottom: 8px;
    line-height: 1.4;
    transition: transform 0.2s ease;
    
    &:before {
      content: '→';
      color: ${({ theme }) => theme.colors.primary};
      position: absolute;
      left: 0;
      transition: transform 0.2s ease;
    }
  }
  
  ${media.mobileLg} {
    margin: 10px 0 0 0;
    
    li {
      font-size: 13px;
      padding-left: 15px;
      margin-bottom: 6px;
    }
  }
`;

// Enhanced CTA button with animated gradient border on hover
const CTAButton = styled(motion.a)`
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
    projects: ['Portfolio Website Redesign', 'E-commerce Mobile App UI', 'Dashboard Interface for Analytics Platform']
  },
  'api': {
    id: 'api',
    name: 'API Coding',
    description: 'Building robust and secure APIs that connect front-end applications to back-end services. Experience with RESTful design principles, authentication, and data handling.',
    icon: 'code',
    projects: ['Weather Data API Integration', 'Payment Gateway API Implementation', 'Social Media Platform API Development']
  },
  'db': {
    id: 'db',
    name: 'Database Management',
    description: 'Designing efficient database structures and managing data storage solutions. Skilled in SQL and NoSQL databases, query optimization, and data security practices.',
    icon: 'database',
    projects: ['Customer Information System', 'Inventory Management Database', 'Analytics Data Warehouse']
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
  performanceSettings: any;
  animationVariants: any;
  animationObjects: any;
}

// Update the AnimatedBioItem component
const AnimatedBioItem = memo(({ 
  text, 
  onClick, 
  isActive, 
  performanceSettings,
  animationVariants,
  animationObjects
}: AnimatedBioItemProps) => {
  const skillId = text === 'UI Designer' ? 'ui' : text === 'API Coding' ? 'api' : 'db';
  
  return (
    <motion.button
      variants={animationVariants.item}
      whileHover={animationObjects.hover}
      whileTap={animationObjects.tap}
      onClick={onClick}
      className={isActive ? 'active' : ''}
      style={{
        fontSize: 'clamp(16px, 3vw, 20px)',
        color: isActive ? '#fff' : 'var(--colors-text)',
        padding: '8px 16px',
        borderRadius: '8px',
        backgroundColor: isActive ? 'var(--colors-primary)' : 'rgba(255, 255, 255, 0.15)',
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

// Optimized SkillDetail component with performance considerations
const SkillDetail = memo(({ 
  skill, 
  performanceSettings, 
  animationVariants 
}: { 
  skill: Skill; 
  performanceSettings: any;
  animationVariants: any;
}) => (
  <SkillDetailContainer
    variants={animationVariants.skillDetail}
    initial="initial"
    animate="animate"
    exit="exit"
    layout
  >
    <SkillTitle>
      {skill.icon && <SkillIcon type={skill.icon} />}
      {skill.name}
    </SkillTitle>
    <SkillDescription>{skill.description}</SkillDescription>
    {skill.projects && skill.projects.length > 0 && (
      <>
        <div style={{ fontSize: '14px', fontWeight: '600', color: 'var(--colors-primary)' }}>Related Projects:</div>
        <ProjectList>
          {skill.projects.map((project, index) => (
            <motion.li 
              key={index}
              initial={{ opacity: 0, x: -5 }}
              animate={{ opacity: 1, x: 0 }}
              transition={{
                delay: index * 0.1,
                duration: 0.3
              }}
            >
              {project}
            </motion.li>
          ))}
        </ProjectList>
      </>
    )}
  </SkillDetailContainer>
));
SkillDetail.displayName = 'SkillDetail';

// Add these styled components before the AnimatedBioItem component
const Bio = styled(motion.div)<StyledMotionProps>`
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  margin: 20px 0 30px;
  justify-content: center;
  
  @media (max-width: ${({ theme }) => theme.breakpoints.tablet}) {
    justify-content: center;
    gap: 8px;
    width: 100%;
  }
`;

const BioItem = styled(motion.button)<{ active: boolean; enablePulse: boolean }>`
  font-size: clamp(16px, 3vw, 20px);
  color: ${({ theme, active }) => active ? theme.colors.backgroundAlt : theme.colors.text};
  padding: 8px 16px;
  border-radius: 8px;
  background-color: ${({ theme, active }) => active 
    ? theme.colors.primary 
    : theme.colors.primaryLight};
  position: relative;
  cursor: pointer;
  pointer-events: auto;
  border: none;
  font-family: inherit;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  overflow: hidden;
  transition: all 0.3s ease;
  width: auto;
  min-width: 140px;

  &:hover {
    background-color: ${({ theme, active }) => active 
      ? theme.colors.primary 
      : `${theme.colors.primary}30`};
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    transform: translateY(-2px);
  }

  @media (max-width: ${({ theme }) => theme.breakpoints.tablet}) {
    font-size: clamp(14px, 2.5vw, 18px);
    padding: 6px 14px;
    width: ${props => props.active ? '100%' : 'auto'};
    max-width: ${props => props.active ? '100%' : '160px'};
    margin: 0 auto;
  }
`;

const TypewriterContainer = styled.div`
  position: relative;
  display: flex;
  flex-wrap: wrap;
  gap: 2px;
  justify-content: center;
  width: 100%;

  @media (max-width: ${({ theme }) => theme.breakpoints.mobile}) {
    span:nth-last-child(-n+5) {
      margin-left: auto;
      margin-right: auto;
    }
    
    // Add line break before "razo/"
    span:nth-last-child(5) {
      &::before {
        content: '';
        display: block;
        width: 100%;
        height: 0;
      }
    }
  }
`;

const TypewriterCharacter = styled(motion.span)<{ $isSlash?: boolean }>`
  display: inline-block;
  color: ${props => props.$isSlash ? props.theme.colors.primary : 'inherit'};
  position: relative;
  white-space: pre;
  
  @media (max-width: ${({ theme }) => theme.breakpoints.mobile}) {
    font-size: inherit;
    line-height: 1.2;
  }
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

// Main Hero component refactored to use the new hooks
export const Hero: React.FC = () => {
  // Get theme from context
  const { theme } = useTheme();
  
  // Initialize component state
  const [isNameHovered, setIsNameHovered] = useState(false);
  const [activeSkill, setActiveSkill] = useState<string | null>(null);
  
  // Create refs for touch interactions and section visibility
  const bioRef = useRef<HTMLDivElement>(null);
  const sectionRef = useRef<HTMLElement>(null);
  const heroRef = useRef<HTMLElement>(null);
  
  // Initialize animation controller for parallax effects
  const parallaxControls = useAnimation();
  
  // Use our new hooks for device optimization
  const deviceCapabilities = useDeviceCapabilities();
  const { performanceSettings } = usePerformanceOptimizations();
  const touchInteractions = useTouchInteractions(heroRef);
  
  // Create observation for element visibility
  const [contentRef, inView] = useInView({
    threshold: 0.1,
    triggerOnce: false
  });
  
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
  
  // Window size tracking for responsive adjustments
  const [windowSize, setWindowSize] = useState({ 
    width: deviceCapabilities.viewportWidth || 0, 
    height: deviceCapabilities.viewportHeight || 0 
  });
  
  // Update window size when device capabilities change
  useEffect(() => {
    setWindowSize({
      width: deviceCapabilities.viewportWidth,
      height: deviceCapabilities.viewportHeight
    });
  }, [deviceCapabilities.viewportWidth, deviceCapabilities.viewportHeight]);
  
  // Parallax effect on mouse move
  const handleMouseMove = useCallback((e: React.MouseEvent) => {
    // Skip parallax on touch devices or if disabled
    if (isTouchDevice || !performanceSettings.enableParallax) return;
    
    const { clientX, clientY } = e;
    const moveX = (clientX - windowSize.width / 2) / 50;
    const moveY = (clientY - windowSize.height / 2) / 50;
    
    parallaxControls.start((i) => ({
      x: moveX * (i * 0.5),
      y: moveY * (i * 0.5),
    }));
  }, [isTouchDevice, windowSize, parallaxControls, performanceSettings.enableParallax]);
  
  // Reset parallax on mouse leave
  const handleMouseLeave = useCallback(() => {
    parallaxControls.start({ x: 0, y: 0 });
  }, [parallaxControls]);
  
  // Memoize bio items
  const bioItems = useMemo(() => [
    'UI Designer',
    'API Coding',
    'Database Management'
  ], []);

  // Split name into characters for typing animation
  const nameCharacters = "Hi there, I'm Jaden Razo/".split('');
  const [typedCount, setTypedCount] = useState(0);
  const [showCursor, setShowCursor] = useState(true);
  const [isTypingComplete, setIsTypingComplete] = useState(false);

  // Update typing effect with improved timing
  useEffect(() => {
    if (typedCount < nameCharacters.length) {
      const timeout = setTimeout(() => {
        setTypedCount(prev => prev + 1);
      }, 80); // Slightly faster typing speed for better UX
      return () => clearTimeout(timeout);
    } else if (!isTypingComplete) {
      setIsTypingComplete(true);
      setShowCursor(false); // Hide cursor immediately after typing is complete
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

  // Handler for skill button click with improved mobile experience
  const handleSkillClick = useCallback((skill: string) => {
    const skillId = skill === 'UI Designer' ? 'ui' : skill === 'API Coding' ? 'api' : 'db';
    
    // Toggle skill display
    setActiveSkill(prev => prev === skillId ? null : skillId);
    
    // Scroll to bio section if skill is activated, with special handling for mobile
    if (activeSkill !== skillId && bioRef.current) {
      const scrollDelay = isTouchDevice ? 300 : 100; // Longer delay on touch for better UX
      
      setTimeout(() => {
        const yOffset = isTouchDevice ? -20 : 0; // Add offset on mobile to improve visibility
        if (bioRef.current) {
          const y = bioRef.current.getBoundingClientRect().top + window.pageYOffset + yOffset;
          window.scrollTo({
            top: y,
            behavior: performanceSettings.reduceMotion ? 'auto' : 'smooth'
          });
        }
      }, scrollDelay);
    }
  }, [activeSkill, isTouchDevice, performanceSettings.reduceMotion]);

  return (
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
            {nameCharacters.slice(0, typedCount).map((char, index) => (
              <TypewriterCharacter
                key={index}
                $isSlash={char === '/'}
                initial={{ opacity: 0, y: -10 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.1 }}
              >
                {char}
              </TypewriterCharacter>
            ))}
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
          {bioItems.map((item) => (
            <AnimatedBioItem 
              key={`bio-item-${item}`} 
              text={item} 
              onClick={() => handleSkillClick(item)} 
              isActive={activeSkill === (item === 'UI Designer' ? 'ui' : item === 'API Coding' ? 'api' : 'db')}
              performanceSettings={performanceSettings}
              animationVariants={animationVariants}
              animationObjects={animationObjects}
            />
          ))}
        </Bio>

        {/* Animated skill details */}
        <AnimatePresence mode="wait">
          {activeSkill && SKILLS[activeSkill] && (
            <SkillDetail 
              skill={SKILLS[activeSkill]} 
              performanceSettings={performanceSettings}
              animationVariants={animationVariants}
            />
          )}
        </AnimatePresence>

        <CTAButton
          href="#projects"
          variants={animationVariants.item}
          whileHover={!isTouchDevice ? animationObjects.hover : undefined}
          whileTap={!isTouchDevice ? animationObjects.tap : undefined}
        >
          Check out my work
        </CTAButton>
      </ContentWrapper>
    </HeroContainer>
  );
};

export default Hero;
