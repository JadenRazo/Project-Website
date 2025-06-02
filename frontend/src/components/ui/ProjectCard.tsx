import React, { useEffect, useRef, useState } from 'react';
import styled, { css } from 'styled-components';
import { motion, useMotionValue, useAnimation } from 'framer-motion';
import usePerformanceOptimizations from '../../hooks/usePerformanceOptimizations';

// Primary card container with hover effects
const Card = styled(motion.div)<{ $isHovered: boolean; $isReducedMotion?: boolean }>`
  position: relative;
  max-width: 320px;
  width: 100%;
  border-radius: 16px;
  overflow: hidden;
  background: ${({ theme }) => theme.colors.surface};
  box-shadow: ${({ theme }) => theme.shadows.medium};
  transition: transform 0.3s ease, box-shadow 0.3s ease;
  margin: 0 auto;
  box-sizing: border-box;
  cursor: pointer;
  transform-style: flat;
  display: flex;
  flex-direction: column;
  min-height: 420px;
  
  ${props => !props.$isReducedMotion && props.$isHovered && css`
    transform: scale(1.02);
    box-shadow: ${({ theme }) => theme.shadows.large};
  `}
  
  &:active {
    transform: scale(0.98);
  }
  
  @media (max-width: 768px) {
    max-width: 100%;
    min-height: 380px;
    border-radius: 12px;
  }
  
  @media (max-width: 480px) {
    max-width: 100%;
    min-height: 360px;
    border-radius: 8px;
  }
`;

// Content container with depth effect
const ProjectContent = styled.div<{ $isHovered: boolean; $isReducedMotion?: boolean }>`
  position: relative;
  z-index: 1;
  padding: 1.5rem;
  flex: 1;
  width: 100%;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  
  ${props => !props.$isReducedMotion && css`
    transform: ${props.$isHovered ? 'translateZ(20px)' : 'translateZ(0)'};
    transition: transform 0.3s ease;
  `}
  
  @media (max-width: 768px) {
    padding: 1.25rem;
  }
  
  @media (max-width: 480px) {
    padding: 1rem;
  }
`;



// Image component with proper styling
const ProjectImage = styled.img`
  width: 100%;
  height: 180px;
  object-fit: cover;
  transition: transform 0.5s ease;
`;

// Gradient overlay for hover effects
const GradientOverlay = styled(motion.div)`
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, rgba(255,255,255,0.1) 0%, rgba(255,255,255,0) 100%);
  opacity: 0;
  pointer-events: none;
  z-index: 1;
`;

// Project info section
const ProjectInfo = styled.div`
  display: flex;
  flex-direction: column;
  gap: 1rem;
  flex: 1;
`;

// Project actions section
const ProjectActions = styled.div`
  display: flex;
  gap: 0.75rem;
  margin-top: auto;
  padding-top: 1rem;
  
  @media (max-width: 480px) {
    gap: 0.5rem;
    flex-direction: row;
  }
  
  @media (max-width: 360px) {
    flex-direction: column;
    gap: 0.5rem;
  }
`;

// Action button
const ActionButton = styled.a`
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  padding: 0.75rem 1rem;
  border-radius: 8px;
  background: ${({ theme }) => theme.colors.primary}15;
  color: ${({ theme }) => theme.colors.primary};
  text-decoration: none;
  font-size: 0.875rem;
  font-weight: 500;
  transition: all 0.2s ease;
  flex: 1;
  text-align: center;
  min-height: 44px;
  
  &:hover {
    background: ${({ theme }) => theme.colors.primary}25;
    transform: translateY(-2px);
  }
  
  @media (max-width: 768px) {
    padding: 0.6rem 0.8rem;
    font-size: 0.8rem;
    min-height: 40px;
  }
  
  @media (max-width: 480px) {
    padding: 0.5rem 0.6rem;
    font-size: 0.75rem;
    min-height: 36px;
    gap: 0.3rem;
    
    &:hover {
      transform: none;
      background: ${({ theme }) => theme.colors.primary}20;
    }
  }
`;

// Project title
const ProjectTitle = styled.h3`
  font-size: 1.25rem;
  font-weight: 600;
  margin-bottom: 0.5rem;
  color: ${({ theme }) => theme.colors.text};
  line-height: 1.3;
  
  @media (max-width: 768px) {
    font-size: 1.15rem;
  }
  
  @media (max-width: 480px) {
    font-size: 1.1rem;
  }
`;

// Project description
const ProjectDescription = styled.p`
  font-size: 0.9rem;
  line-height: 1.5;
  color: ${({ theme }) => theme.colors.textSecondary};
  margin-bottom: 0.5rem;
  
  @media (max-width: 480px) {
    font-size: 0.85rem;
    line-height: 1.4;
  }
`;

interface ProjectCardProps {
  id?: string;
  title: string;
  description: string;
  image: string;
  link: string;
  language?: string;
  useSimplifiedEffects?: boolean; 
  supportsBackdropFilter?: boolean;
}

export const ProjectCard: React.FC<ProjectCardProps> = ({ 
  title, 
  description, 
  image, 
  link,
  useSimplifiedEffects = false,
  supportsBackdropFilter = true
}) => {
  // Motion values for 3D effect
  const x = useMotionValue(0);
  const y = useMotionValue(0);
  const controls = useAnimation();
  
  // Get performance settings to conditionally enable effects
  const { performanceSettings } = usePerformanceOptimizations();
  const isReducedMotion = useSimplifiedEffects || performanceSettings?.reduceMotion;
  
  // For tracking if image is loaded and hover state
  const [imageLoaded, setImageLoaded] = useState(false);
  const [isHovered, setIsHovered] = useState(false);
  const cardRef = useRef<HTMLDivElement>(null);
  

  // Use IntersectionObserver to only animate when in view
  useEffect(() => {
    if (!cardRef.current) return;
    
    const cardElement = cardRef.current;
    const observer = new IntersectionObserver((entries) => {
      entries.forEach(entry => {
        if (entry.isIntersecting) {
          // Reset transform when card comes into view
          controls.start({ scale: 1, opacity: 1 });
        } else {
          // Optionally reset when out of view
          controls.start({ scale: 0.98, opacity: 0.8 });
        }
      });
    }, { threshold: 0.1 });
    
    observer.observe(cardElement);
    
    return () => {
      if (cardElement) observer.unobserve(cardElement);
    };
  }, [controls]);

  // Handle mouse movement for 3D effect
  const handleMouseMove = (event: React.MouseEvent<HTMLDivElement>) => {
    if (isReducedMotion) return; // Skip effect on reduced motion
    
    const rect = event.currentTarget.getBoundingClientRect();
    const centerX = rect.left + rect.width / 2;
    const centerY = rect.top + rect.height / 2;
    
    x.set(event.clientX - centerX);
    y.set(event.clientY - centerY);
    setIsHovered(true);
  };

  // Reset card position on mouse leave
  const handleMouseLeave = () => {
    controls.start({ 
      rotateX: 0, 
      rotateY: 0, 
      transition: { duration: 0.5 } 
    });
    
    // Reset motion values
    x.set(0);
    y.set(0);
    setIsHovered(false);
  };

  // Handle touch events for mobile
  const handleTouchMove = (event: React.TouchEvent<HTMLDivElement>) => {
    if (isReducedMotion) return; // Skip effect on reduced motion
    
    const rect = event.currentTarget.getBoundingClientRect();
    const centerX = rect.left + rect.width / 2;
    const centerY = rect.top + rect.height / 2;
    
    x.set(event.touches[0].clientX - centerX);
    y.set(event.touches[0].clientY - centerY);
    setIsHovered(true);
  };

  // Handle image load
  const handleImageLoad = () => {
    setImageLoaded(true);
  };

  return (
    <Card
      ref={cardRef}
      $isHovered={isHovered}
      $isReducedMotion={isReducedMotion}
      whileHover={!isReducedMotion ? { scale: 1.03 } : undefined}
      onMouseMove={handleMouseMove}
      onMouseLeave={handleMouseLeave}
      onTouchMove={handleTouchMove}
      onTouchEnd={handleMouseLeave}
      animate={controls}
      initial={{ scale: 0.98, opacity: 0.8 }}
    >
      <GradientOverlay 
        initial={{ opacity: 0 }}
        whileHover={{ opacity: supportsBackdropFilter ? 1 : 0.5 }} 
      />
      <ProjectImage 
        src={image} 
        alt={title} 
        loading="lazy"
        onLoad={handleImageLoad}
        style={{ opacity: imageLoaded ? 1 : 0, transition: 'opacity 0.3s ease' }}
      />
      <ProjectContent $isHovered={isHovered} $isReducedMotion={isReducedMotion}>
        <ProjectInfo>
          <ProjectTitle>{title}</ProjectTitle>
          <ProjectDescription>{description}</ProjectDescription>
        </ProjectInfo>
        <ProjectActions>
          <ActionButton
            href={link}
            target="_blank" 
            rel="noopener noreferrer"
            aria-label={`View project: ${title}`}
          >
            View Project
          </ActionButton>
        </ProjectActions>
      </ProjectContent>
    </Card>
  );
};

export default ProjectCard;