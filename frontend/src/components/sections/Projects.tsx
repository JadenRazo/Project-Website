// src/components/sections/Projects.tsx
import React, { useState, useMemo, useCallback, useRef, useEffect } from 'react';
import styled from 'styled-components';
import { motion, AnimatePresence, useAnimation } from 'framer-motion';
import { useInView } from 'react-intersection-observer';
import useDeviceCapabilities from '../../hooks/useDeviceCapabilities';
import usePerformanceOptimizations from '../../hooks/usePerformanceOptimizations';
import useTouchInteractions from '../../hooks/useTouchInteractions';
import { ProjectCard } from '../ui/ProjectCard';
import LanguageFilter from '../ui/LanguageFilter';

// Types
interface Project {
  id: string;
  title: string;
  description: string;
  image: string;
  link: string;
  language: string;
}

interface ProjectsProps {
  projects?: readonly Project[];
  title?: string;
  subtitle?: string;
  languages?: readonly string[];
  selectedLanguage?: string;
  onLanguageChange?: (language: string) => void;
}

// Responsive container with dynamic spacing based on viewport
const ProjectsSection = styled.section<{ isReducedMotion?: boolean }>`
  padding: clamp(2rem, 4vw, 4rem) clamp(0.5rem, 2vw, 1.5rem);
  max-width: var(--page-max-width, 100vw);
  width: 100%;
  margin: 0 auto;
  overflow-x: hidden;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  background: ${({ theme }) => theme.colors.background};
  color: ${({ theme }) => theme.colors.text};
  
  @media (max-width: 768px) {
    padding: clamp(1.5rem, 3vw, 2.5rem) clamp(0.5rem, 2vw, 1rem);
  }
  
  @media (max-width: 480px) {
    padding: 1.5rem 0.75rem;
  }
`;

// Content wrapper with intersection-based animations
const ContentWrapper = styled(motion.div)<{ inView: boolean }>`
  max-width: var(--content-max-width, 900px);
  width: 100%;
  margin: 0 auto;
  opacity: ${props => props.inView ? 1 : 0};
  transform: translateY(${props => props.inView ? 0 : '20px'});
  transition: opacity 0.6s ease, transform 0.6s ease;
  will-change: opacity, transform;
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  box-sizing: border-box;
`;

// Responsive heading with dynamic sizing
const Heading = styled(motion.h2)`
  font-size: clamp(2rem, 4vw, 3rem);
  margin-bottom: 0.5rem;
  position: relative;
  display: inline-block;
  text-rendering: optimizeLegibility;
  text-align: center;
  color: ${({ theme }) => theme.colors.primary};
  
  &::after {
    content: '';
    position: absolute;
    bottom: -5px;
    left: 50%;
    transform: translateX(-50%);
    width: 60px;
    height: 3px;
    background-color: ${({ theme }) => theme.colors?.primary || 'var(--primary)'};
  }
  
  @media (max-width: 768px) {
    font-size: clamp(1.75rem, 7vw, 2.5rem);
  }
`;

const Subtitle = styled(motion.p)`
  font-size: clamp(1rem, 1.5vw, 1.25rem);
  margin-bottom: 2rem;
  max-width: 600px;
  opacity: 0.8;
  text-align: center;
  color: ${({ theme }) => theme.colors.text};
  
  @media (max-width: 768px) {
    font-size: 1rem;
    margin-bottom: 1.5rem;
    padding: 0 1rem;
  }
`;

// Project grid with responsive layout
const ProjectsGrid = styled(motion.div)`
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: clamp(1.5rem, 3vw, 2.5rem);
  margin-top: 2rem;
  width: 100%;
  justify-items: center;
  box-sizing: border-box;
  
  @media (max-width: 768px) {
    grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
    gap: 1.5rem;
    padding: 0;
  }
  
  @media (max-width: 480px) {
    grid-template-columns: 1fr;
    gap: 1.25rem;
  }
`;

// Empty state message
const EmptyState = styled(motion.div)`
  text-align: center;
  padding: 3rem 1rem;
  margin-top: 2rem;
  background-color: ${({ theme }) => `${theme.colors.surface}80`};
  border-radius: 8px;
  backdrop-filter: blur(5px);
  -webkit-backdrop-filter: blur(5px);
  width: 100%;
  max-width: 600px;
  margin-left: auto;
  margin-right: auto;
  
  h3 {
    font-size: 1.5rem;
    margin-bottom: 1rem;
    color: ${({ theme }) => theme.colors?.primary || 'var(--primary)'};
  }
  
  p {
    opacity: 0.8;
    max-width: 500px;
    margin: 0 auto;
    color: ${({ theme }) => theme.colors.text};
  }
`;

// Project card wrapper for consistent sizing and animation
const ProjectCardWrapper = styled(motion.div)`
  width: 100%;
  max-width: 280px;
  display: flex;
  justify-content: center;
  transform-style: flat; // Prevent 3D transforms from affecting layout
  box-sizing: border-box;
  
  @media (max-width: 480px) {
    max-width: 100%;
    padding: 0 0.5rem;
  }
`;

export const Projects: React.FC<ProjectsProps> = ({
  projects = [],
  title = 'Selected Work',
  subtitle = 'A showcase of my recent projects and experiments.',
  languages = ['All'],
  selectedLanguage = 'All',
  onLanguageChange
}) => {
  // State for filtering
  const [currentLanguage, setCurrentLanguage] = useState(selectedLanguage);
  const [isInitialRender, setIsInitialRender] = useState(true);
  
  // Refs and hooks
  const containerRef = useRef<HTMLElement>(null);
  const [contentRef, inView] = useInView({
    threshold: 0.1,
    triggerOnce: false
  });
  
  // Get device capabilities and performance settings for optimization
  const deviceCapabilities = useDeviceCapabilities();
  const { performanceSettings } = usePerformanceOptimizations();
  const touchInteractions = useTouchInteractions(containerRef);
  
  // Animation controls for grid items
  const controls = useAnimation();
  
  // After first render, update initialRender state
  useEffect(() => {
    if (isInitialRender) {
      requestAnimationFrame(() => {
        setIsInitialRender(false);
      });
    }
  }, [isInitialRender]);
  
  // Handle scroll into view
  useEffect(() => {
    if (inView) {
      controls.start('visible');
    }
  }, [controls, inView]);
  
  // Effect to sync external selectedLanguage prop if provided
  useEffect(() => {
    if (selectedLanguage !== currentLanguage) {
      setCurrentLanguage(selectedLanguage);
    }
  }, [selectedLanguage]);
  
  // Memoize filter handler to prevent unnecessary re-renders
  const handleLanguageSelect = useCallback((language: string) => {
    setCurrentLanguage(language);
    if (onLanguageChange) {
      onLanguageChange(language);
    }
  }, [onLanguageChange]);
  
  // Memoize filtered projects for performance
  const filteredProjects = useMemo(() => {
    if (currentLanguage === 'All') return projects;
    return projects.filter(project => project.language === currentLanguage);
  }, [currentLanguage, projects]);
  
  // Create animation variants based on performance settings and device capabilities
  const containerVariants = useMemo(() => ({
    hidden: { opacity: 0 },
    visible: {
      opacity: 1,
      transition: {
        staggerChildren: performanceSettings?.reduceMotion ? 0.03 : 0.08,
        delayChildren: 0.05,
        when: "beforeChildren"
      }
    }
  }), [performanceSettings?.reduceMotion]);
  
  const itemVariants = useMemo(() => ({
    hidden: { 
      opacity: 0, 
      y: performanceSettings?.reduceMotion ? 10 : 20,
      scale: 0.98
    },
    visible: { 
      opacity: 1, 
      y: 0,
      scale: 1,
      transition: {
        type: performanceSettings?.reduceMotion ? 'tween' : 'spring',
        duration: performanceSettings?.transitionSpeed || 0.4,
        bounce: performanceSettings?.reduceMotion ? 0 : 0.2
      }
    }
  }), [performanceSettings]);
  
  // Simplified animation for low-powered devices
  const shouldUseSimplifiedAnimation = useMemo(() => 
    deviceCapabilities.isLowPoweredDevice || 
    deviceCapabilities.prefersReducedMotion ||
    (performanceSettings?.reduceMotion ?? false),
  [deviceCapabilities, performanceSettings]);
  
  // Detect if browser supports backdrop-filter for visual effects
  const supportsBackdropFilter = useMemo(() => {
    if (typeof window !== 'undefined') {
      // Check if CSS.supports is available (modern browsers)
      return typeof CSS !== 'undefined' && CSS.supports 
        ? CSS.supports('(backdrop-filter: blur(10px))') || CSS.supports('(-webkit-backdrop-filter: blur(10px))')
        : false;
    }
    return false;
  }, []);
  
  return (
    <ProjectsSection 
      ref={containerRef}
      isReducedMotion={shouldUseSimplifiedAnimation} 
      id="projects"
      role="region"
      aria-label="Projects Section"
    >
      <ContentWrapper 
        ref={contentRef}
        inView={inView}
      >
        <Heading
          initial={isInitialRender ? false : { opacity: 0, x: -20 }}
          animate={inView ? { opacity: 1, x: 0 } : { opacity: 0, x: -20 }}
          transition={{ duration: 0.5 }}
        >
          {title}
        </Heading>
        
        <Subtitle
          initial={isInitialRender ? false : { opacity: 0 }}
          animate={inView ? { opacity: 1 } : { opacity: 0 }}
          transition={{ duration: 0.5, delay: 0.1 }}
        >
          {subtitle}
        </Subtitle>
        
        <LanguageFilter
          languages={languages}
          selectedLanguage={currentLanguage}
          onSelectLanguage={handleLanguageSelect}
        />
        
        <AnimatePresence mode="wait" initial={false}>
          <motion.div
            key={currentLanguage}
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            transition={{ duration: 0.3 }}
          >
            {filteredProjects.length > 0 ? (
              <ProjectsGrid
                variants={containerVariants}
                initial="hidden"
                animate={controls}
              >
                {filteredProjects.map((project, index) => (
                  <ProjectCardWrapper
                    key={project.id || `project-${index}`}
                    variants={itemVariants}
                    custom={index}
                    layoutId={`project-${project.id || index}`}
                  >
                    <ProjectCard 
                      {...project} 
                      useSimplifiedEffects={shouldUseSimplifiedAnimation}
                      supportsBackdropFilter={supportsBackdropFilter}
                    />
                  </ProjectCardWrapper>
                ))}
              </ProjectsGrid>
            ) : (
              <EmptyState
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.4 }}
              >
                <h3>No projects found</h3>
                <p>There are no projects matching the selected filter. Try selecting a different language.</p>
              </EmptyState>
            )}
          </motion.div>
        </AnimatePresence>
      </ContentWrapper>
    </ProjectsSection>
  );
};

export default Projects;
