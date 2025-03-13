import React, { useRef, useState, useEffect, useMemo } from 'react';
import styled from 'styled-components';
import { motion, useAnimation } from 'framer-motion';
import { useInView } from 'react-intersection-observer';
import { useTheme } from '../../contexts/ThemeContext';
import { useZIndex } from '../../hooks/useZIndex';
import SkillBar from '../ui/SkillBar';
import CreativeShaderBackground from '../animations/CreativeShaderBackground';
import usePerformanceOptimizations from '../../hooks/usePerformanceOptimizations';

interface Skill {
  name: string;
  percentage: number;
  category: string;
}

// Main container with modern design and relative positioning for animation overlay
const SkillsContainer = styled.section`
  position: relative;
  min-height: 100vh;
  width: 100%;
  padding: 6rem 2rem;
  overflow: hidden;
  background: ${({ theme }) => `${theme.colors.backgroundAlt}80`};
  isolation: isolate; // Create stacking context for z-index
  
  @media (max-width: 768px) {
    padding: 5rem 1.5rem;
  }
`;

// Content wrapper with improved animation transitions
const ContentWrapper = styled(motion.div)<{ $visible: boolean }>`
  position: relative;
  max-width: 1200px;
  margin: 0 auto;
  z-index: 2;
  opacity: ${props => props.$visible ? 1 : 0};
  transform: translateY(${props => props.$visible ? 0 : '30px'});
  transition: opacity 0.6s ease, transform 0.6s ease;
  will-change: opacity, transform;
`;

// Section heading with animated underline
const SectionHeading = styled(motion.h2)`
  font-size: 2.5rem;
  margin-bottom: 1rem;
  position: relative;
  display: inline-block;
  color: ${({ theme }) => theme.colors.text};
  
  &::after {
    content: '';
    position: absolute;
    bottom: -10px;
    left: 0;
    width: 60px;
    height: 3px;
    background: ${({ theme }) => theme.colors.primary};
    transform-origin: left center;
  }
  
  @media (max-width: 768px) {
    font-size: 2rem;
  }
`;

const SectionDescription = styled(motion.p)`
  font-size: 1.1rem;
  margin-bottom: 3rem;
  max-width: 600px;
  line-height: 1.6;
  color: ${({ theme }) => theme.colors.textSecondary || theme.colors.text};
  opacity: 0.9;
`;

// Modern grid for skill categories
const SkillsGrid = styled(motion.div)`
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 3rem;
  
  @media (max-width: 768px) {
    grid-template-columns: 1fr;
    gap: 2rem;
  }
`;

// Glass card design for category containers
const CategoryContainer = styled(motion.div)`
  background: ${({ theme }) => `${theme.colors.background}80`};
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  border-radius: 16px;
  padding: 2rem;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  border: 1px solid ${({ theme }) => `${theme.colors.primary}30`};
  will-change: transform;
  
  /* Refined hover effect */
  transition: transform 0.3s cubic-bezier(0.34, 1.56, 0.64, 1), 
              box-shadow 0.3s ease,
              background-color 0.3s ease;
              
  &:hover {
    transform: translateY(-8px);
    box-shadow: 0 12px 40px rgba(0, 0, 0, 0.15);
    background: ${({ theme }) => `${theme.colors.background}95`};
  }
`;

const CategoryTitle = styled(motion.h3)`
  font-size: 1.5rem;
  margin-bottom: 1.5rem;
  color: ${({ theme }) => theme.colors.primary};
  padding-bottom: 0.75rem;
  border-bottom: 1px solid ${({ theme }) => `${theme.colors.primary}30`};
  display: flex;
  align-items: center;
  
  &::before {
    content: '';
    display: inline-block;
    width: 8px;
    height: 8px;
    background: ${({ theme }) => theme.colors.primary};
    border-radius: 50%;
    margin-right: 12px;
  }
`;

// Animation variants with refined timing
const containerVariants = {
  hidden: { 
    opacity: 0,
    transition: {
      staggerChildren: 0.1,
      staggerDirection: -1,
      when: "afterChildren"
    }
  },
  visible: {
    opacity: 1,
    transition: {
      staggerChildren: 0.15,
      delayChildren: 0.1,
      when: "beforeChildren"
    }
  }
};

const itemVariants = {
  hidden: { 
    opacity: 0, 
    y: 30,
    transition: { 
      duration: 0.4, 
      ease: [0.43, 0.13, 0.23, 0.96]
    }
  },
  visible: {
    opacity: 1,
    y: 0,
    transition: { 
      duration: 0.6, 
      ease: [0.22, 1, 0.36, 1]
    }
  }
};

// Invisible marker divs for better scroll detection
const ScrollDetector = styled.div`
  position: absolute;
  width: 100%;
  height: 10%;
  left: 0;
  pointer-events: none;
  opacity: 0;
`;

const TopDetector = styled(ScrollDetector)`
  top: 10%;
`;

const MiddleDetector = styled(ScrollDetector)`
  top: 50%;
  transform: translateY(-50%);
`;

const BottomDetector = styled(ScrollDetector)`
  bottom: 10%;
`;

// List of skills with their percentages and categories
const skillsData: Skill[] = [
  // Frontend skills
  { name: 'React.js', percentage: 95, category: 'Frontend' },
  { name: 'TypeScript', percentage: 90, category: 'Frontend' },
  { name: 'Next.js', percentage: 85, category: 'Frontend' },
  { name: 'CSS/SCSS', percentage: 92, category: 'Frontend' },
  { name: 'Framer Motion', percentage: 80, category: 'Frontend' },
  
  // Backend skills
  { name: 'Node.js', percentage: 88, category: 'Backend' },
  { name: 'Express', percentage: 85, category: 'Backend' },
  { name: 'MongoDB', percentage: 82, category: 'Backend' },
  { name: 'PostgreSQL', percentage: 78, category: 'Backend' },
  { name: 'GraphQL', percentage: 75, category: 'Backend' },
  
  // Design & Tools
  { name: 'Figma', percentage: 88, category: 'Design & Tools' },
  { name: 'UX/UI Design', percentage: 85, category: 'Design & Tools' },
  { name: 'Git/GitHub', percentage: 92, category: 'Design & Tools' },
  { name: 'DevOps', percentage: 70, category: 'Design & Tools' },
  { name: 'Responsive Design', percentage: 95, category: 'Design & Tools' },
];

export const Skills: React.FC = () => {
  // State and hooks
  const { theme } = useTheme();
  const { zIndex } = useZIndex();
  const { performanceSettings } = usePerformanceOptimizations();
  
  // Animation controls
  const headingControls = useAnimation();
  const descriptionControls = useAnimation();
  const gridControls = useAnimation();
  
  // Refs for sections
  const sectionRef = useRef<HTMLDivElement>(null);
  
  // Improved visibility detection with multiple detection points
  const { ref: topRef, inView: isTopVisible } = useInView({ 
    threshold: 0.1,
    triggerOnce: false 
  });
  
  const { ref: middleRef, inView: isMiddleVisible } = useInView({ 
    threshold: 0.1,
    triggerOnce: false 
  });
  
  const { ref: bottomRef, inView: isBottomVisible } = useInView({ 
    threshold: 0.1,
    triggerOnce: false 
  });
  
  // More reliable section visibility with debouncing
  const [isSectionVisible, setIsSectionVisible] = useState(false);
  const [isAnimatingOut, setIsAnimatingOut] = useState(false);
  const [shouldAnimateIn, setShouldAnimateIn] = useState(false);
  const [isFirstRender, setIsFirstRender] = useState(true);
  const visibilityTimerRef = useRef<NodeJS.Timeout | null>(null);
  
  // Group skills by category
  const skillsByCategory = useMemo(() => {
    return skillsData.reduce<Record<string, Skill[]>>((acc, skill) => {
      if (!acc[skill.category]) {
        acc[skill.category] = [];
      }
      acc[skill.category].push(skill);
      return acc;
    }, {});
  }, []);
  
  // Handle initial load animation
  useEffect(() => {
    if (isFirstRender) {
      setIsFirstRender(false);
      return;
    }
  }, [isFirstRender]);
  
  // Enhanced visibility detection with debouncing to prevent flickering
  useEffect(() => {
    const isCurrentlyVisible = isTopVisible || isMiddleVisible || isBottomVisible;
    
    // Clear any existing visibility timer
    if (visibilityTimerRef.current) {
      clearTimeout(visibilityTimerRef.current);
      visibilityTimerRef.current = null;
    }
    
    // Handle visibility changes with debouncing
    if (isCurrentlyVisible && !isSectionVisible) {
      // Immediately set as visible when entering viewport
      setIsSectionVisible(true);
      setShouldAnimateIn(true);
      setIsAnimatingOut(false);
    } else if (!isCurrentlyVisible && isSectionVisible) {
      // Delayed exit to prevent flickering
      setIsAnimatingOut(true);
      visibilityTimerRef.current = setTimeout(() => {
        setIsSectionVisible(false);
        setShouldAnimateIn(false);
        setIsAnimatingOut(false);
      }, 300);
    }
    
    return () => {
      if (visibilityTimerRef.current) {
        clearTimeout(visibilityTimerRef.current);
      }
    };
  }, [isTopVisible, isMiddleVisible, isBottomVisible, isSectionVisible]);
  
  // Handle section visibility and animations with improved timing
  useEffect(() => {
    // Skip animations on first server render
    if (isFirstRender) return;
    
    if (isSectionVisible && !isAnimatingOut) {
      // Section is in view - Animate elements in
      const animationDelay = performanceSettings.reduceMotion ? 0 : 100;
      
      // Begin staggered animations
      setTimeout(() => {
        headingControls.start("visible");
        
        setTimeout(() => {
          descriptionControls.start("visible");
          
          setTimeout(() => {
            gridControls.start("visible");
          }, animationDelay);
        }, animationDelay);
      }, animationDelay);
    } else {
      // Section is not in view - Reset animations with appropriate delay
      const exitDelay = performanceSettings.reduceMotion ? 0 : 100;
      setTimeout(() => {
        headingControls.start("hidden");
        descriptionControls.start("hidden");
        gridControls.start("hidden");
      }, exitDelay);
    }
  }, [
    isFirstRender,
    isSectionVisible,
    isAnimatingOut,
    headingControls,
    descriptionControls,
    gridControls,
    performanceSettings.reduceMotion
  ]);
  
  // Generate random animation delays for staggered entrance
  const getRandomDelay = useMemo(() => {
    return (index: number) => ({
      visible: { 
        transition: { 
          delay: 0.1 + (index * 0.07),
          duration: 0.6, 
          ease: [0.22, 1, 0.36, 1] 
        }
      }
    });
  }, []);
  
  return (
    <SkillsContainer 
      ref={sectionRef} 
      id="skills"
    >
      {/* Multiple visibility detection points for better scroll detection */}
      <TopDetector ref={topRef} />
      <MiddleDetector ref={middleRef} />
      <BottomDetector ref={bottomRef} />
      
      {/* Creative background with conditional rendering for performance */}
      {isSectionVisible && (
        <CreativeShaderBackground 
          disableParallax={performanceSettings.reduceMotion}
          intensity={0.7}
          speed={0.5}
          colorIntensity={0.6}
        />
      )}
      
      <ContentWrapper 
        $visible={isSectionVisible}
      >
        <SectionHeading
          variants={itemVariants}
          initial="hidden"
          animate={headingControls}
        >
          My Skills
        </SectionHeading>
        
        <SectionDescription
          variants={itemVariants}
          initial="hidden"
          animate={descriptionControls}
        >
          I specialize in building responsive, high-performance applications 
          with modern web technologies. My experience spans both frontend and 
          backend development with a focus on clean, maintainable code.
        </SectionDescription>
        
        <SkillsGrid
          variants={containerVariants}
          initial="hidden"
          animate={gridControls}
        >
          {Object.entries(skillsByCategory).map(([category, skills], categoryIndex) => (
            <CategoryContainer 
              key={category} 
              variants={itemVariants}
              custom={categoryIndex}
              {...getRandomDelay(categoryIndex)}
            >
              <CategoryTitle>{category}</CategoryTitle>
              {skills.map((skill, skillIndex) => (
                <SkillBar
                  key={skill.name}
                  skill={skill.name}
                  percentage={skill.percentage}
                  shouldAnimate={shouldAnimateIn}
                  delay={0.15 + (skillIndex * 0.08)}
                />
              ))}
            </CategoryContainer>
          ))}
        </SkillsGrid>
      </ContentWrapper>
    </SkillsContainer>
  );
};

export default Skills; 