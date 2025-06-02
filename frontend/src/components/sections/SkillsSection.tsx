import React, { useState, useEffect, useMemo } from 'react';
import styled from 'styled-components';
import { motion, useAnimation, AnimatePresence, Variants } from 'framer-motion';
import { useInView } from 'react-intersection-observer';

// SVG Icons for skill categories
const DesignIcon = () => (
  <svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
    <path d="M12 22C17.5228 22 22 17.5228 22 12C22 6.47715 17.5228 2 12 2C6.47715 2 2 6.47715 2 12C2 17.5228 6.47715 22 12 22Z" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
    <path d="M12 18C15.3137 18 18 15.3137 18 12C18 8.68629 15.3137 6 12 6C8.68629 6 6 8.68629 6 12C6 15.3137 8.68629 18 12 18Z" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
    <path d="M12 14C13.1046 14 14 13.1046 14 12C14 10.8954 13.1046 10 12 10C10.8954 10 10 10.8954 10 12C10 13.1046 10.8954 14 12 14Z" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
  </svg>
);

const CodeIcon = () => (
  <svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
    <path d="M16 18L22 12L16 6" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
    <path d="M8 6L2 12L8 18" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
  </svg>
);

const ServerIcon = () => (
  <svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
    <path d="M2 9H22" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
    <path d="M2 15H22" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
    <path d="M5 3H19C20.1046 3 21 3.89543 21 5V19C21 20.1046 20.1046 21 19 21H5C3.89543 21 3 20.1046 3 19V5C3 3.89543 3.89543 3 5 3Z" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
    <path d="M6 6H6.01" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
    <path d="M6 12H6.01" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
    <path d="M6 18H6.01" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
  </svg>
);

interface Skill {
  id: string;
  name: string;
  icon: React.ReactNode;
  level: number;
  description: string;
  projects: string[];
  color: string;
}

const SkillsSectionContainer = styled.section`
  position: relative;
  width: 100%;
  padding: 6rem 2rem;
  min-height: 100vh;
  overflow-x: hidden;
  background: ${({ theme }) => theme.colors.background};
  color: ${({ theme }) => theme.colors.text};
  
  @media (min-width: 768px) {
    padding: 8rem 4rem;
  }
`;

const ContentWrapper = styled(motion.div)`
  max-width: 1400px;
  margin: 0 auto;
  will-change: opacity, transform;
`;

const SectionHeader = styled(motion.div)`
  margin-bottom: 3.5rem;
  position: relative;
`;

const SectionTitle = styled(motion.h2)`
  font-size: 2.5rem;
  font-weight: 700;
  margin-bottom: 1.5rem;
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
`;

const SectionDescription = styled(motion.p)`
  font-size: 1.1rem;
  line-height: 1.6;
  color: ${({ theme }) => theme.colors.text};
  max-width: 800px;
  opacity: 0.9;
`;

const SkillsGrid = styled(motion.div)`
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 2rem;
  
  @media (min-width: 992px) {
    gap: 3rem;
  }
`;

const SkillCard = styled(motion.div)`
  background: ${({ theme }) => theme.colors.backgroundAlt || theme.colors.surface};
  border-radius: 12px;
  padding: 2rem;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  height: 100%;
  display: flex;
  flex-direction: column;
  transform: translateZ(0);
  backface-visibility: hidden;
  perspective: 1000px;
  will-change: transform, box-shadow;
  
  &:hover {
    transform: translateY(-8px);
    box-shadow: 0 8px 30px rgba(0, 0, 0, 0.12);
  }
`;

const SkillHeader = styled.div`
  display: flex;
  align-items: center;
  margin-bottom: 1.5rem;
`;

const IconWrapper = styled.div<{ color: string }>`
  display: flex;
  align-items: center;
  justify-content: center;
  width: 48px;
  height: 48px;
  border-radius: 12px;
  background: ${props => props.color}20;
  color: ${props => props.color};
  margin-right: 1rem;
`;

const SkillName = styled.h3`
  font-size: 1.25rem;
  font-weight: 600;
  margin: 0;
  color: ${({ theme }) => theme.colors.text};
`;

const SkillLevel = styled.div`
  margin-bottom: 1.5rem;
`;

const ProgressBar = styled.div`
  height: 6px;
  width: 100%;
  background: ${({ theme }) => theme.colors.surfaceLight || 'rgba(255, 255, 255, 0.1)'};
  border-radius: 3px;
  overflow: hidden;
  margin-top: 0.5rem;
`;

const Progress = styled(motion.div)<{ level: number; color: string }>`
  height: 100%;
  width: ${props => props.level}%;
  background: ${props => props.color};
  border-radius: 3px;
  transform-origin: left center;
`;

const SkillDescription = styled.p`
  font-size: 1rem;
  line-height: 1.6;
  color: ${({ theme }) => theme.colors.text};
  opacity: 0.7;
  margin-bottom: 1.5rem;
  flex-grow: 1;
`;

const ProjectsList = styled.div`
  margin-top: auto;
`;

const ProjectsTitle = styled.h4`
  font-size: 0.875rem;
  font-weight: 600;
  margin-bottom: 0.75rem;
  color: ${({ theme }) => theme.colors.text};
`;

const ProjectTags = styled.div`
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
`;

const ProjectTag = styled.span<{ color: string }>`
  font-size: 0.75rem;
  padding: 0.25rem 0.75rem;
  border-radius: 20px;
  background: ${props => props.color}15;
  color: ${props => props.color};
  transition: background-color 0.2s ease, transform 0.15s ease;
  transform: translateZ(0);
  cursor: pointer;
  
  &:hover {
    background: ${props => props.color}30;
    transform: translateY(-2px);
  }
`;

export const SkillsSection: React.FC = () => {
  
  // Optimized animation configuration for faster animations
  const animationConfig = useMemo(() => ({
    staggerChildren: 0.05, // Reduced from 0.1 for faster sequence
    childrenDelay: 0.08, // Reduced from 0.15 for faster initial delay
    duration: {
      enter: 0.5, // Reduced from 0.7 for faster entrance
      exit: 0.3 // Reduced from 0.5 for faster exit
    },
    threshold: 0.1, // Slightly reduced for earlier triggering
    rootMargin: "-5% 0px -5% 0px" // Adjusted for earlier triggering
  }), []);
  
  // InView hook with optimized threshold
  const [sectionRef, inView] = useInView({
    threshold: animationConfig.threshold,
    rootMargin: animationConfig.rootMargin,
    triggerOnce: false
  });
  
  // Animation controls
  const controls = useAnimation();
  
  // Animation state
  const [hasAnimated, setHasAnimated] = useState(false);
  
  // Optimized animation variants with faster transitions
  const variants: Record<string, Variants> = useMemo(() => ({
    section: {
      hidden: { opacity: 0 },
      visible: { 
        opacity: 1,
        transition: { 
          when: "beforeChildren",
          staggerChildren: animationConfig.staggerChildren,
          duration: animationConfig.duration.enter * 0.8, // Even faster for container
          ease: [0.16, 1, 0.3, 1] // Snappier easing
        }
      },
      exit: {
        opacity: 0,
        transition: {
          when: "afterChildren",
          staggerChildren: animationConfig.staggerChildren / 2,
          staggerDirection: -1,
          duration: animationConfig.duration.exit,
          ease: [0.43, 0.13, 0.23, 0.96]
        }
      }
    },
    title: {
      hidden: { 
        opacity: 0, 
        y: 15 // Reduced from 20 for less movement
      },
      visible: { 
        opacity: 1, 
        y: 0,
        transition: { 
          duration: animationConfig.duration.enter * 0.8, // Faster title animation
          ease: [0.16, 1, 0.3, 1] // Snappier easing
        }
      },
      exit: {
        opacity: 0,
        y: 5, // Reduced from 10 for less movement
        transition: {
          duration: animationConfig.duration.exit,
          ease: [0.43, 0.13, 0.23, 0.96]
        }
      }
    },
    description: {
      hidden: { 
        opacity: 0, 
        y: 15 // Reduced from 20 for less movement
      },
      visible: { 
        opacity: 1, 
        y: 0,
        transition: { 
          duration: animationConfig.duration.enter * 0.8, // Faster description animation
          delay: animationConfig.childrenDelay / 2,
          ease: [0.16, 1, 0.3, 1] // Snappier easing
        }
      },
      exit: {
        opacity: 0,
        y: 5, // Reduced from 10 for less movement
        transition: {
          duration: animationConfig.duration.exit,
          ease: [0.43, 0.13, 0.23, 0.96]
        }
      }
    },
    card: {
      hidden: { 
        opacity: 0, 
        y: 20, // Reduced from 30 for less movement
        scale: 0.98 // Closer to final scale for faster visual appearance
      },
      visible: (index: number) => ({ 
        opacity: 1, 
        y: 0,
        scale: 1,
        transition: { 
          duration: animationConfig.duration.enter * 0.8, // Faster card animation
          delay: animationConfig.childrenDelay + (index * 0.05), // Reduced from 0.08 for faster staggering
          ease: [0.16, 1, 0.3, 1] // Snappier easing
        }
      }),
      exit: {
        opacity: 0,
        y: 10, // Reduced from 20 for less movement
        scale: 0.98,
        transition: {
          duration: animationConfig.duration.exit,
          ease: [0.43, 0.13, 0.23, 0.96]
        }
      }
    },
    progress: {
      hidden: { 
        scaleX: 0,
        originX: 0
      },
      visible: (index: number) => ({ 
        scaleX: 1,
        transition: { 
          duration: 0.8, // Reduced from 1.2 for faster progress animation
          delay: animationConfig.childrenDelay + (index * 0.05) + 0.1, // Reduced delays
          ease: [0.16, 1, 0.3, 1] // Snappier easing
        }
      }),
      exit: {
        scaleX: 0,
        transition: {
          duration: animationConfig.duration.exit / 2,
          ease: [0.43, 0.13, 0.23, 0.96]
        }
      }
    }
  }), [animationConfig]);
  
  // Animation synchronization based on visibility with faster response
  useEffect(() => {
    if (inView) {
      // Start animations immediately when in view
      controls.start("visible").catch(() => {});
      setHasAnimated(true);
    } else if (hasAnimated) {
      controls.start("exit").catch(() => {});
    }
  }, [controls, inView, hasAnimated]);
  
  
  // Skills data
  const skills: Skill[] = [
    {
      id: 'frontend',
      name: 'Frontend Development',
      icon: <CodeIcon />,
      level: 95,
      description: 'Creating responsive and interactive web applications with modern frameworks like React and Next.js.',
      projects: ['Personal Portfolio', 'E-commerce Platform', 'Dashboard UI'],
      color: '#3498db'
    },
    {
      id: 'backend',
      name: 'Backend Development',
      icon: <ServerIcon />,
      level: 85,
      description: 'Building scalable APIs and server-side applications with Node.js and database integration.',
      projects: ['REST API', 'Authentication System', 'Data Processing Service'],
      color: '#2ecc71'
    },
    {
      id: 'design',
      name: 'UI/UX Design',
      icon: <DesignIcon />,
      level: 80,
      description: 'Designing intuitive user interfaces and experiences with a focus on accessibility and usability.',
      projects: ['Mobile App Design', 'Design System', 'Website Redesign'],
      color: '#e74c3c'
    },
    {
      id: 'architecture',
      name: 'System Architecture',
      icon: <ServerIcon />,
      level: 75,
      description: 'Designing and implementing scalable system architectures and deployment processes.',
      projects: ['Microservices Setup', 'Cloud Migration', 'Performance Optimization'],
      color: '#9b59b6'
    },
    {
      id: 'animation',
      name: 'Motion & Animation',
      icon: <DesignIcon />,
      level: 90,
      description: 'Creating smooth, engaging animations and transitions to enhance user experience.',
      projects: ['Interactive Tutorial', 'Animated Landing Page', 'WebGL Experiments'],
      color: '#f39c12'
    },
    {
      id: 'performance',
      name: 'Performance Optimization',
      icon: <CodeIcon />,
      level: 85,
      description: 'Optimizing web applications for speed, accessibility, and user experience.',
      projects: ['Core Web Vitals', 'Lighthouse Audit', 'Bundle Optimization'],
      color: '#1abc9c'
    }
  ];

  return (
    <SkillsSectionContainer ref={sectionRef} id="skills">
      <AnimatePresence mode="wait">
        <ContentWrapper
          key="skills-content"
          initial="hidden"
          animate={controls}
          variants={variants.section}
          exit="exit"
        >
          <SectionHeader>
            <SectionTitle variants={variants.title}>
              My Skills
            </SectionTitle>
            
            <SectionDescription variants={variants.description}>
              I specialize in building responsive, high-performance applications 
              with modern web technologies. My experience spans both frontend and 
              backend development with a focus on clean, maintainable code.
            </SectionDescription>
          </SectionHeader>
          
          <SkillsGrid>
            {skills.map((skill, index) => (
              <SkillCard
                key={skill.id}
                custom={index}
                variants={variants.card}
                whileHover={{ y: -8, boxShadow: "0 10px 30px rgba(0, 0, 0, 0.15)" }}
                whileTap={{ scale: 0.98 }}
              >
                <SkillHeader>
                  <IconWrapper color={skill.color}>
                    {skill.icon}
                  </IconWrapper>
                  <SkillName>{skill.name}</SkillName>
                </SkillHeader>
                
                <SkillLevel>
                  <ProgressBar>
                    <Progress 
                      level={skill.level} 
                      color={skill.color}
                      variants={variants.progress}
                      custom={index}
                    />
                  </ProgressBar>
                </SkillLevel>
                
                <SkillDescription>
                  {skill.description}
                </SkillDescription>
                
                <ProjectsList>
                  <ProjectsTitle>Related Projects</ProjectsTitle>
                  <ProjectTags>
                    {skill.projects.map((project, projectIndex) => (
                      <ProjectTag 
                        key={`${skill.id}-project-${projectIndex}`}
                        color={skill.color}
                      >
                        {project}
                      </ProjectTag>
                    ))}
                  </ProjectTags>
                </ProjectsList>
              </SkillCard>
            ))}
          </SkillsGrid>
        </ContentWrapper>
      </AnimatePresence>
    </SkillsSectionContainer>
  );
};

export default SkillsSection;
