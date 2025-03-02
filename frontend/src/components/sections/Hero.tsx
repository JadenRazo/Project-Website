import React, { memo, useMemo, useState, useCallback, useRef } from 'react';
import styled, { css, keyframes } from 'styled-components';
import { motion, Variants, HTMLMotionProps, AnimatePresence } from 'framer-motion';
import { useTheme } from '../../contexts/ThemeContext';

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

// Performance-optimized animation variants
const ANIMATION_VARIANTS: Record<string, Variants> = {
  container: {
    hidden: { opacity: 0 },
    visible: {
      opacity: 1,
      transition: {
        staggerChildren: 0.08,
        when: "beforeChildren",
      },
    },
  },
  item: {
    hidden: { opacity: 0, y: 15 },
    visible: {
      opacity: 1,
      y: 0,
      transition: {
        type: "tween",
        duration: 0.4,
        ease: [0.25, 0.1, 0.25, 1.0],
      },
    },
  },
  skillDetail: {
    initial: { 
      opacity: 0, 
      height: 0,
      y: -20
    },
    animate: { 
      opacity: 1, 
      height: 'auto',
      y: 0,
      transition: {
        opacity: { duration: 0.3 },
        height: { duration: 0.4 },
        y: { duration: 0.3, ease: "easeOut" }
      }
    },
    exit: { 
      opacity: 0, 
      height: 0,
      y: -10,
      transition: {
        opacity: { duration: 0.2 },
        height: { duration: 0.3 },
        y: { duration: 0.2 }
      }
    }
  }
};

// Simple animation objects
const HOVER_ANIMATION = {
  scale: 1.01,
  transition: { duration: 0.2, ease: "easeOut" }
};

const TAP_ANIMATION = {
  scale: 0.98,
  transition: { duration: 0.1 }
};

// Keyframes for animations
const pulse = keyframes`
  0% { box-shadow: 0 0 0 0 rgba(var(--primary-rgb), 0.7); }
  70% { box-shadow: 0 0 0 10px rgba(var(--primary-rgb), 0); }
  100% { box-shadow: 0 0 0 0 rgba(var(--primary-rgb), 0); }
`;

const shimmerAnimation = keyframes`
  0% { background-position: 200% center; }
  100% { background-position: -200% center; }
`;

// Styled components with performance optimizations
const HeroSection = styled.section`
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  justify-content: center;
  padding: 0 clamp(1rem, 5vw, 10rem);
  max-width: 1600px;
  margin: 0 auto;
  user-select: none;
  position: relative;
  z-index: 1;
  background-color: ${({ theme }) => theme.colors.background};
  color: ${({ theme }) => theme.colors.text};
  overscroll-behavior: none;
`;

const ContentWrapper = styled(motion.div)<AnimatedElementProps>`
  max-width: 800px;
  pointer-events: ${props => props.isInteractive ? 'auto' : 'none'};
  will-change: transform, opacity;
`;

const Greeting = styled(motion.span)`
  color: ${({ theme }) => theme.colors.primary};
  font-family: ${({ theme }) => theme.fonts.mono};
  font-size: clamp(14px, 2vw, 16px);
  font-weight: 400;
  margin-bottom: 20px;
  display: block;
  pointer-events: auto;
`;

interface NameProps {
  isHovered: boolean;
}

// Optimized name container with enhanced animations
const NameContainer = styled(motion.h1)<NameProps>`
  font-size: clamp(40px, 8vw, 80px);
  font-weight: 600;
  color: ${({ theme }) => theme.colors.text};
  line-height: 1.1;
  margin: 0;
  cursor: pointer;
  transform-style: preserve-3d;
  perspective: 1000px;
  
  ${({ isHovered, theme }) => isHovered && css`
    background-image: linear-gradient(
      90deg, 
      ${theme.colors.text} 0%,
      ${theme.colors.primary} 30%,
      ${theme.colors.accent || theme.colors.secondary || theme.colors.primary} 50%,
      ${theme.colors.primary} 70%,
      ${theme.colors.text} 100%
    );
    background-size: 200% auto;
    background-clip: text;
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    animation: ${shimmerAnimation} 3s linear infinite;
  `}
  
  transition: all 0.3s ease-out;
`;

// New technique for animating name as blocks on initial load
const NameBlockContainer = styled(motion.div)`
  display: flex;
  flex-wrap: wrap;
  margin: 0;
`;

const NameBlock = styled(motion.div)`
  display: inline-block;
  margin-right: 0.3em;
  transform-origin: center left;
`;

const Bio = styled(motion.div)`
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  margin: 20px 0 30px;
`;

interface BioItemStyledProps {
  active: boolean;
}

const BioItem = styled(motion.button)<BioItemStyledProps>`
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
  
  ${({ active }) => active && css`
    animation: ${pulse} 1.5s infinite;
  `}
  
  &:hover {
    background-color: ${({ theme, active }) => active 
      ? theme.colors.primary 
      : `${theme.colors.primary}30`};
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    transform: translateY(-2px);
  }
  
  &:after {
    content: '';
    position: absolute;
    top: -50%;
    left: -50%;
    width: 200%;
    height: 200%;
    background: radial-gradient(circle, rgba(255,255,255,0.2) 0%, rgba(255,255,255,0) 60%);
    opacity: 0;
    transition: opacity 0.3s ease;
    pointer-events: none;
    transform: scale(0.5);
  }
  
  &:active:after {
    opacity: 1;
    transform: scale(1);
    transition: transform 0.3s ease, opacity 0s;
  }
`;

const SkillDetailContainer = styled(motion.div)`
  width: 100%;
  background: ${({ theme }) => theme.colors.backgroundAlt};
  border-radius: 12px;
  padding: 20px;
  margin-top: 10px;
  margin-bottom: 20px;
  color: ${({ theme }) => theme.colors.text};
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
  border: 1px solid ${({ theme }) => theme.colors.primary}20;
  overflow: hidden;
`;

const SkillTitle = styled.h3`
  font-size: clamp(18px, 4vw, 24px);
  color: ${({ theme }) => theme.colors.primary};
  margin: 0 0 15px 0;
  display: flex;
  align-items: center;
  gap: 10px;
`;

const SkillDescription = styled.p`
  font-size: 16px;
  line-height: 1.6;
  margin-bottom: 15px;
`;

const ProjectList = styled.ul`
  list-style-type: none;
  padding: 0;
  margin: 15px 0 0 0;
  
  li {
    position: relative;
    padding-left: 20px;
    margin-bottom: 8px;
    line-height: 1.4;
    
    &:before {
      content: '→';
      color: ${({ theme }) => theme.colors.primary};
      position: absolute;
      left: 0;
    }
  }
`;

const CTAButton = styled(motion.a)`
  display: inline-block;
  background-color: transparent;
  border: 1px solid ${({ theme }) => theme.colors.primary};
  border-radius: 8px;
  color: ${({ theme }) => theme.colors.primary};
  font-family: ${({ theme }) => theme.fonts.mono};
  font-size: clamp(14px, 2vw, 16px);
  padding: 1.25rem 1.75rem;
  text-decoration: none;
  position: relative;
  overflow: hidden;
  pointer-events: auto;
  transition: all 0.3s ease;
  
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
    
    &:after {
      left: 200%;
    }
  }
`;

// Skill data
const SKILLS: Record<string, Skill> = {
  'ui': {
    id: 'ui',
    name: 'UI Designer',
    description: 'Creating intuitive and visually appealing user interfaces that prioritize user experience and accessibility. Proficient in design principles, color theory, and responsive layouts.',
    projects: ['Portfolio Website Redesign', 'E-commerce Mobile App UI', 'Dashboard Interface for Analytics Platform']
  },
  'api': {
    id: 'api',
    name: 'API Coding',
    description: 'Building robust and secure APIs that connect front-end applications to back-end services. Experience with RESTful design principles, authentication, and data handling.',
    projects: ['Weather Data API Integration', 'Payment Gateway API Implementation', 'Social Media Platform API Development']
  },
  'db': {
    id: 'db',
    name: 'Database Management',
    description: 'Designing efficient database structures and managing data storage solutions. Skilled in SQL and NoSQL databases, query optimization, and data security practices.',
    projects: ['Customer Information System', 'Inventory Management Database', 'Analytics Data Warehouse']
  }
};

// Optimized components using memoization
const AnimatedBioItem = memo(({ text, onClick, isActive }: BioItemProps) => {
  const skillId = text === 'UI Designer' ? 'ui' : text === 'API Coding' ? 'api' : 'db';
  
  return (
    <BioItem
      variants={ANIMATION_VARIANTS.item}
      whileHover={HOVER_ANIMATION}
      whileTap={TAP_ANIMATION}
      onClick={onClick}
      active={isActive}
    >
      {text}
    </BioItem>
  );
});
AnimatedBioItem.displayName = 'AnimatedBioItem';

const SkillDetail = memo(({ skill }: { skill: Skill }) => (
  <SkillDetailContainer
    variants={ANIMATION_VARIANTS.skillDetail}
    initial="initial"
    animate="animate"
    exit="exit"
    layout
  >
    <SkillTitle>{skill.name}</SkillTitle>
    <SkillDescription>{skill.description}</SkillDescription>
    {skill.projects && skill.projects.length > 0 && (
      <>
        <div style={{ fontSize: '14px', color: 'var(--primary)' }}>Related Projects:</div>
        <ProjectList>
          {skill.projects.map((project, index) => (
            <li key={index}>{project}</li>
          ))}
        </ProjectList>
      </>
    )}
  </SkillDetailContainer>
));
SkillDetail.displayName = 'SkillDetail';

// Optimized and modernized Hero component
export const Hero: React.FC = () => {
  const { theme } = useTheme();
  const [isNameHovered, setIsNameHovered] = useState(false);
  const [activeSkill, setActiveSkill] = useState<string | null>(null);
  const bioRef = useRef<HTMLDivElement>(null);
  
  // Memoize bio items
  const bioItems = useMemo(() => [
    'UI Designer',
    'API Coding',
    'Database Management'
  ], []);

  // Split name into words instead of letters for better performance
  const nameWords = useMemo(() => {
    return "Jaden Scott Razo".split(' ').map((word, i) => ({
      word,
      key: `word-${i}`
    }));
  }, []);

  // Handlers for name hover
  const handleNameMouseEnter = useCallback(() => {
    setIsNameHovered(true);
  }, []);

  const handleNameMouseLeave = useCallback(() => {
    setIsNameHovered(false);
  }, []);

  // Handler for skill button click
  const handleSkillClick = useCallback((skill: string) => {
    const skillId = skill === 'UI Designer' ? 'ui' : skill === 'API Coding' ? 'api' : 'db';
    
    // Toggle skill display
    setActiveSkill(prev => prev === skillId ? null : skillId);
    
    // Scroll to bio section if skill is activated
    if (activeSkill !== skillId && bioRef.current) {
      setTimeout(() => {
        bioRef.current?.scrollIntoView({ 
          behavior: 'smooth',
          block: 'start'
        });
      }, 100);
    }
  }, [activeSkill]);

  return (
    <HeroSection>
      <ContentWrapper
        variants={ANIMATION_VARIANTS.container}
        initial="hidden"
        animate="visible"
        isInteractive={true}
      >
        <Greeting
          variants={ANIMATION_VARIANTS.item}
          whileHover={HOVER_ANIMATION}
          whileTap={TAP_ANIMATION}
        >
          Hi there, I'm
        </Greeting>
        
        {/* Optimized name animation using word blocks instead of letters */}
        <NameContainer 
          variants={ANIMATION_VARIANTS.item}
          isHovered={isNameHovered}
          onMouseEnter={handleNameMouseEnter}
          onMouseLeave={handleNameMouseLeave}
        >
          <NameBlockContainer>
            {nameWords.map(({ word, key }, index) => (
              <NameBlock
                key={key}
                initial={{ opacity: 0, y: 25, rotateX: 30 }}
                animate={{ 
                  opacity: 1, 
                  y: 0, 
                  rotateX: 0,
                  transition: {
                    type: 'spring',
                    stiffness: 100,
                    damping: 15,
                    mass: 1,
                    delay: 0.3 + (index * 0.1)
                  }
                }}
              >
                {word}
              </NameBlock>
            ))}
          </NameBlockContainer>
        </NameContainer>

        <Bio 
          ref={bioRef}
          variants={ANIMATION_VARIANTS.container}
        >
          {bioItems.map((item) => (
            <AnimatedBioItem 
              key={`bio-item-${item}`} 
              text={item} 
              onClick={() => handleSkillClick(item)} 
              isActive={activeSkill === (item === 'UI Designer' ? 'ui' : item === 'API Coding' ? 'api' : 'db')}
            />
          ))}
        </Bio>

        {/* Animated skill details */}
        <AnimatePresence>
          {activeSkill && SKILLS[activeSkill] && (
            <SkillDetail skill={SKILLS[activeSkill]} />
          )}
        </AnimatePresence>

        <CTAButton
          href="#work"
          variants={ANIMATION_VARIANTS.item}
          whileHover={HOVER_ANIMATION}
          whileTap={TAP_ANIMATION}
        >
          Check out my work
        </CTAButton>
      </ContentWrapper>
    </HeroSection>
  );
};
