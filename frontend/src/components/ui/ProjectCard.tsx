import React, { useEffect, useRef, useState } from 'react';
import styled, { css } from 'styled-components';
import { motion, useMotionValue, useTransform, useAnimation } from 'framer-motion';
import usePerformanceOptimizations from '../../hooks/usePerformanceOptimizations';

// Primary card container with hover effects
const Card = styled(motion.div)<{ $isHovered: boolean; $isReducedMotion?: boolean }>`
  position: relative;
  max-width: 280px;
  width: 100%;
  border-radius: 16px;
  overflow: hidden;
  background: ${({ theme }) => theme.colors.surface};
  box-shadow: ${({ theme }) => theme.shadows.medium};
  transition: transform 0.3s ease, box-shadow 0.3s ease;
  margin: 0 auto;
  box-sizing: border-box;
  cursor: pointer;
  transform-style: flat; // Prevent 3D transforms from affecting layout
  
  ${props => !props.$isReducedMotion && props.$isHovered && css`
    transform: scale(1.02);
    box-shadow: ${({ theme }) => theme.shadows.large};
  `}
  
  &:active {
    transform: scale(0.98);
  }
  
  @media (max-width: 480px) {
    max-width: 100%;
  }
`;

// Content container with depth effect
const ProjectContent = styled.div<{ $isHovered: boolean; $isReducedMotion?: boolean }>`
  position: relative;
  z-index: 1;
  padding: 1.5rem;
  height: 100%;
  width: 100%;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  
  ${props => !props.$isReducedMotion && css`
    transform: ${props.$isHovered ? 'translateZ(20px)' : 'translateZ(0)'};
    transition: transform 0.3s ease;
  `}
`;

// Image container with proper constraints
const ImageContainer = styled.div`
  width: 100%;
  height: 180px;
  overflow: hidden;
  position: relative;
  box-sizing: border-box;
  
  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    transition: transform 0.5s ease;
  }
`;

// Tech badge with contained sizing
const TechBadge = styled.span`
  display: inline-block;
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-size: 0.75rem;
  margin-right: 0.5rem;
  margin-bottom: 0.5rem;
  background: ${({ theme }) => theme.colors.surface || 'rgba(255,255,255,0.1)'};
  color: ${({ theme }) => theme.colors.text || 'white'};
  white-space: nowrap;
  box-sizing: border-box;
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

// Link wrapper for the entire card
const ProjectLink = styled.a`
  position: absolute;
  inset: 0;
  z-index: 10;
  text-decoration: none;
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
  
  // Define transform ranges based on device capability
  const rotateRange = isReducedMotion ? 2 : 5;
  const rotateX = useTransform(y, [-100, 100], [rotateRange, -rotateRange]);
  const rotateY = useTransform(x, [-100, 100], [-rotateRange, rotateRange]);

  // Use IntersectionObserver to only animate when in view
  useEffect(() => {
    if (!cardRef.current) return;
    
    const observer = new IntersectionObserver((entries) => {
      entries.forEach(entry => {
        if (entry.isIntersecting) {
          // Reset transform when card comes into view
          controls.start({ scale: 1, opacity: 1 });
        } else {
          // Optionally reset when out of view
          controls.set({ scale: 0.98, opacity: 0.8 });
        }
      });
    }, { threshold: 0.1 });
    
    observer.observe(cardRef.current);
    
    return () => {
      if (cardRef.current) observer.unobserve(cardRef.current);
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
        <h3>{title}</h3>
        <p>{description}</p>
      </ProjectContent>
      <ProjectLink 
        href={link}
        target="_blank" 
        rel="noopener noreferrer"
        aria-label={`View project: ${title}`}
      />
    </Card>
  );
};

export default ProjectCard;