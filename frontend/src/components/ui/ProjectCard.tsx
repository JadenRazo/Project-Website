import React, { useEffect, useRef, useState } from 'react';
import styled from 'styled-components';
import { motion, useMotionValue, useTransform, useAnimation } from 'framer-motion';
import usePerformanceOptimizations from '../../hooks/usePerformanceOptimizations';

// Performance-optimized card with conditional effects
const Card = styled(motion.div)<{ isPowerfulDevice?: boolean }>`
  width: 300px;
  height: 400px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 20px;
  padding: 20px;
  position: relative;
  overflow: hidden;
  backdrop-filter: blur(10px);
  cursor: pointer;
  user-select: none;
  transform-style: preserve-3d; /* For 3D effect */
  transform: translateZ(0); /* Hardware acceleration */
  will-change: transform; /* Hint for browser optimization */
  
  @media (max-width: 768px) {
    width: 280px;
    height: 380px;
  }
  
  @media (max-width: 480px) {
    width: 100%;
    max-width: 320px;
    height: 360px;
  }
`;

const ProjectImage = styled.img`
  width: 100%;
  height: 200px;
  object-fit: cover;
  border-radius: 10px;
  transform: translateZ(20px); /* Subtle depth effect */
  
  @media (max-width: 480px) {
    height: 180px;
  }
`;

const ProjectContent = styled.div`
  margin-top: 20px;
  color: white;
  transform: translateZ(30px); /* More pronounced depth effect for text */
  
  h3 {
    font-size: 1.25rem;
    margin-bottom: 10px;
    font-weight: 600;
    
    @media (max-width: 480px) {
      font-size: 1.1rem;
    }
  }
  
  p {
    font-size: 0.9rem;
    line-height: 1.6;
    opacity: 0.8;
    
    @media (max-width: 480px) {
      font-size: 0.85rem;
      line-height: 1.5;
    }
  }
`;

// Optimized gradient overlay that appears on hover
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

// Actual link element for better accessibility 
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
  const isPowerfulDevice = !useSimplifiedEffects && performanceSettings?.performanceTier === 'high';
  
  // For tracking if image is loaded
  const [imageLoaded, setImageLoaded] = useState(false);
  const cardRef = useRef<HTMLDivElement>(null);
  
  // Define transform ranges based on device capability
  const rotateRange = isPowerfulDevice ? 10 : 5;
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
    if (!isPowerfulDevice) return; // Skip effect on less powerful devices
    
    const rect = event.currentTarget.getBoundingClientRect();
    const centerX = rect.left + rect.width / 2;
    const centerY = rect.top + rect.height / 2;
    
    x.set(event.clientX - centerX);
    y.set(event.clientY - centerY);
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
  };

  // Handle touch events for mobile
  const handleTouchMove = (event: React.TouchEvent<HTMLDivElement>) => {
    if (!isPowerfulDevice) return; // Skip effect on less powerful devices
    
    const rect = event.currentTarget.getBoundingClientRect();
    const centerX = rect.left + rect.width / 2;
    const centerY = rect.top + rect.height / 2;
    
    x.set(event.touches[0].clientX - centerX);
    y.set(event.touches[0].clientY - centerY);
  };

  // Handle image load
  const handleImageLoad = () => {
    setImageLoaded(true);
  };

  return (
    <Card
      ref={cardRef}
      whileHover={isPowerfulDevice ? { scale: 1.05 } : undefined}
      style={{ rotateX, rotateY, perspective: 1000 }}
      onMouseMove={handleMouseMove}
      onMouseLeave={handleMouseLeave}
      onTouchMove={handleTouchMove}
      onTouchEnd={handleMouseLeave}
      animate={controls}
      initial={{ scale: 0.98, opacity: 0.8 }}
      isPowerfulDevice={isPowerfulDevice}
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
      <ProjectContent>
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