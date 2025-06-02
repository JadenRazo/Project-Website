// src/components/animations/FloatingElement.tsx
import React, { useState, useCallback, useRef, useEffect } from 'react';
import styled from 'styled-components';
import { motion, AnimatePresence, useTransform, useMotionValue } from 'framer-motion';

// Enhanced container with 3D support
const FloatingContainer = styled(motion.div)`
  position: absolute;
  cursor: pointer;
  z-index: 5;
  transform: translateZ(0);
  backface-visibility: hidden;
  will-change: transform;
  transform-style: preserve-3d;
  transition: filter 0.3s ease;
  
  &:hover {
    filter: brightness(1.2) drop-shadow(0 0 5px rgba(255, 255, 255, 0.3));
  }
`;

// Optimized particle with minimal rendering overhead
const Particle = styled(motion.div)<{ 
  $color: string; 
  $size: number; 
  $glow: number;
}>`
  position: absolute;
  width: ${props => props.$size}px;
  height: ${props => props.$size}px;
  border-radius: 50%;
  background-color: ${props => props.$color};
  box-shadow: 0 0 ${props => props.$glow}px ${props => props.$glow / 2}px ${props => props.$color};
  pointer-events: none;
  will-change: transform, opacity;
`;

// Particle data structure
interface ParticleData {
  id: number;
  size: number;
  color: string;
  angle: number;
  distance: number;
  speed: number;
  delay: number;
  glow: number;
  initialX: number;
  initialY: number;
}

// Enhanced props interface with new animation options
interface FloatingElementProps {
  children: React.ReactNode;
  x?: number;
  y?: number;
  particleCount?: number;
  particleColors?: string[];
  particleLifetime?: number;
  shouldReset?: boolean;
  glowIntensity?: number;
  floatDuration?: number;
  floatRadius?: number;
  scale?: number;
  rotationSpeed?: number;
  rotationAxis?: 'x' | 'y' | 'both';
}

export const FloatingElement: React.FC<FloatingElementProps> = ({
  children,
  x = 0,
  y = 0,
  particleCount = 30,
  particleColors = ['#6c63ff', '#ff6b6b', '#48dbfb', '#1dd1a1'],
  particleLifetime = 2500,
  shouldReset = true,
  glowIntensity = 5,
  floatDuration = 20,
  floatRadius = 25,
  scale = 1,
  rotationSpeed = 1,
  rotationAxis = 'both'
}) => {
  // Component state
  const [isPopped, setIsPopped] = useState(false);
  const [particles, setParticles] = useState<ParticleData[]>([]);
  const [elementSize, setElementSize] = useState({ width: 0, height: 0 });
  
  // Motion values for smooth animations
  const rotateX = useMotionValue(0);
  const rotateY = useMotionValue(0);
  
  // Refs for optimization
  const elementRef = useRef<HTMLDivElement>(null);
  const isMounted = useRef(true);
  const animationFrame = useRef<number | null>(null);
  const resizeObserver = useRef<ResizeObserver | null>(null);
  
  // Cleanup on unmount
  useEffect(() => {
    return () => {
      isMounted.current = false;
      if (animationFrame.current) {
        cancelAnimationFrame(animationFrame.current);
      }
      if (resizeObserver.current) {
        resizeObserver.current.disconnect();
      }
    };
  }, []);
  
  // Update element size with ResizeObserver for performance
  useEffect(() => {
    if (!elementRef.current) return;
    
    const updateSize = () => {
      if (!elementRef.current || !isMounted.current) return;
      
      const rect = elementRef.current.getBoundingClientRect();
      if (rect.width > 0 && rect.height > 0) {
        setElementSize({ 
          width: rect.width, 
          height: rect.height 
        });
      }
    };
    
    // Initial size update
    updateSize();
    
    try {
      // Use ResizeObserver for better performance than resize event
      resizeObserver.current = new ResizeObserver(() => {
        if (isMounted.current) {
          animationFrame.current = requestAnimationFrame(updateSize);
        }
      });
      
      resizeObserver.current.observe(elementRef.current);
    } catch (error) {
      // Fallback to window resize for older browsers
      const handleResize = () => {
        if (isMounted.current) {
          animationFrame.current = requestAnimationFrame(updateSize);
        }
      };
      
      window.addEventListener('resize', handleResize, { passive: true });
      
      return () => {
        window.removeEventListener('resize', handleResize);
      };
    }
  }, []);
  
  // Helper function to create particles - defined inline to avoid dependency issues
  const createParticle = (x: number, y: number, angle: number): ParticleData => {
    // Generate a unique ID efficiently
    const id = Math.floor(Math.random() * 1000000);
    
    // Randomize properties within constrained ranges for performance
    const size = 2 + Math.random() * 4;
    const distance = 30 + Math.random() * 70;
    const speed = 0.5 + Math.random() * 1.5;
    const delay = Math.random() * 150;
    const colorIndex = Math.floor(Math.random() * particleColors.length);
    const glow = glowIntensity * (0.5 + Math.random() * 0.5);
    
    return {
      id,
      size,
      color: particleColors[colorIndex],
      angle,
      distance,
      speed,
      delay,
      glow,
      initialX: x,
      initialY: y
    };
  };
  
  // Generate outline particles when shape is clicked
  const generateOutlineParticles = useCallback(() => {
    if (!elementRef.current || !isMounted.current) return;
    
    const { width, height } = elementSize;
    if (width <= 0 || height <= 0) return;
    
    // Create array for particles
    const newParticles: ParticleData[] = [];
    
    // Calculate number of particles per side based on perimeter
    const perimeter = 2 * (width + height);
    const particlesPerPixel = particleCount / perimeter;
    
    // Distribute particles along the outline
    // Top edge
    const topCount = Math.max(2, Math.floor(width * particlesPerPixel));
    for (let i = 0; i < topCount; i++) {
      const x = (width * i) / (topCount - 1);
      const y = 0;
      const angle = Math.random() * Math.PI - Math.PI / 2; // Upward with variation
      
      newParticles.push(createParticle(x, y, angle));
    }
    
    // Right edge
    const rightCount = Math.max(2, Math.floor(height * particlesPerPixel));
    for (let i = 0; i < rightCount; i++) {
      const x = width;
      const y = (height * i) / (rightCount - 1);
      const angle = Math.random() * Math.PI + Math.PI / 2; // Rightward with variation
      
      newParticles.push(createParticle(x, y, angle));
    }
    
    // Bottom edge
    const bottomCount = Math.max(2, Math.floor(width * particlesPerPixel));
    for (let i = 0; i < bottomCount; i++) {
      const x = width - (width * i) / (bottomCount - 1);
      const y = height;
      const angle = Math.random() * Math.PI + Math.PI / 2; // Downward with variation
      
      newParticles.push(createParticle(x, y, angle));
    }
    
    // Left edge
    const leftCount = Math.max(2, Math.floor(height * particlesPerPixel));
    for (let i = 0; i < leftCount; i++) {
      const x = 0;
      const y = height - (height * i) / (leftCount - 1);
      const angle = Math.random() * Math.PI - Math.PI / 2; // Leftward with variation
      
      newParticles.push(createParticle(x, y, angle));
    }
    
    // Update state
    setParticles(newParticles);
    
    // Schedule cleanup if needed
    if (shouldReset) {
      const timerId = setTimeout(() => {
        if (isMounted.current) {
          setIsPopped(false);
          setParticles([]);
        }
      }, particleLifetime + 200);
      
      return () => clearTimeout(timerId);
    }
  }, [elementSize, particleCount, particleColors, particleLifetime, shouldReset, glowIntensity, createParticle]);
  
  // Handle element click with proper event handling
  const handlePop = useCallback((e: React.MouseEvent<HTMLDivElement>) => {
    e.preventDefault();
    e.stopPropagation();
    
    if (!isPopped && isMounted.current) {
      setIsPopped(true);
      generateOutlineParticles();
    }
  }, [isPopped, generateOutlineParticles]);
  
  // Add 3D tilt effect on mouse movement
  const handleMouseMove = useCallback((e: React.MouseEvent<HTMLDivElement>) => {
    const rect = e.currentTarget.getBoundingClientRect();
    const centerX = rect.left + rect.width / 2;
    const centerY = rect.top + rect.height / 2;
    
    const mouseXFromCenter = e.clientX - centerX;
    const mouseYFromCenter = e.clientY - centerY;
    
    // Calculate rotation based on mouse position relative to center
    const rotX = (mouseYFromCenter / rect.height) * 20; // Max 20 degrees
    const rotY = (mouseXFromCenter / rect.width) * 20; // Max 20 degrees
    
    rotateX.set(rotX);
    rotateY.set(rotY);
  }, [rotateX, rotateY]);
  
  // Reset rotation when mouse leaves
  const handleMouseLeave = useCallback(() => {
    rotateX.set(0);
    rotateY.set(0);
  }, [rotateX, rotateY]);
  
  // Generate enhanced floating animation
  const getFloatAnimation = useCallback(() => {
    // Create more complex float path with multiple points
    const numPoints = 6; // More points for smoother path
    const radius = floatRadius;
    const duration = floatDuration;
    
    // Create slightly elliptical path with random variations
    const xPath: number[] = [];
    const yPath: number[] = [];
    
    for (let i = 0; i <= numPoints; i++) {
      const anglePoint = (i / numPoints) * Math.PI * 2;
      const xOffset = Math.cos(anglePoint) * radius * (0.8 + Math.random() * 0.4);
      const yOffset = Math.sin(anglePoint) * radius * (0.8 + Math.random() * 0.4);
      
      xPath.push(x + xOffset);
      yPath.push(y + yOffset);
    }
    
    // Add starting point to end for seamless loop
    xPath.push(xPath[0]);
    yPath.push(yPath[0]);
    
    // Calculate rotation animation based on selected axis
    const rotationX = rotationAxis === 'x' || rotationAxis === 'both' 
      ? [0, 10, 0, -10, 0] 
      : [0];
      
    const rotationY = rotationAxis === 'y' || rotationAxis === 'both' 
      ? [0, -10, 0, 10, 0] 
      : [0];
    
    return {
      x: xPath,
      y: yPath,
      rotateX: rotationX,
      rotateY: rotationY,
      scale: [scale, scale * 1.05, scale, scale * 0.95, scale],
      transition: {
        x: {
          duration,
          repeat: Infinity,
          ease: "easeInOut",
          times: Array.from({ length: xPath.length }).map((_, i) => i / (xPath.length - 1))
        },
        y: {
          duration,
          repeat: Infinity,
          ease: "easeInOut",
          times: Array.from({ length: yPath.length }).map((_, i) => i / (yPath.length - 1))
        },
        rotateX: {
          duration: duration * 0.6,
          repeat: Infinity,
          ease: "easeInOut",
        },
        rotateY: {
          duration: duration * 0.8,
          repeat: Infinity,
          ease: "easeInOut",
        },
        scale: {
          duration: duration * 0.7,
          repeat: Infinity,
          ease: "easeInOut",
        }
      }
    };
  }, [x, y, floatDuration, floatRadius, scale, rotationAxis]);
  
  // Convert mouse motion to rotation transforms
  const mouseRotateX = useTransform(rotateX, [-20, 20], [-10, 10]);
  const mouseRotateY = useTransform(rotateY, [-20, 20], [10, -10]);
  
  // Memoize the animation values
  const floatAnimation = React.useMemo(() => getFloatAnimation(), [getFloatAnimation]);
  
  // Ensure we don't render on server
  if (typeof window === 'undefined') {
    return null;
  }
  
  return (
    <>
      {/* Main floating element */}
      <AnimatePresence mode="wait">
        {!isPopped && (
          <FloatingContainer
            ref={elementRef}
            onClick={handlePop}
            onMouseMove={handleMouseMove}
            onMouseLeave={handleMouseLeave}
            initial={{ opacity: 0, scale: 0.8 }}
            animate={{ 
              opacity: 1, 
              scale: floatAnimation.scale,
              x: floatAnimation.x,
              y: floatAnimation.y,
              rotateX: floatAnimation.rotateX,
              rotateY: floatAnimation.rotateY,
              transition: floatAnimation.transition
            }}
            style={{
              rotateX: mouseRotateX,
              rotateY: mouseRotateY,
              scale
            }}
            exit={{ 
              scale: [scale, scale * 1.2, 0], 
              opacity: [1, 1, 0], 
              transition: { duration: 0.3 }
            }}
            whileHover={{ scale: scale * 1.1 }}
            whileTap={{ scale: scale * 0.95 }}
          >
            {children}
          </FloatingContainer>
        )}
      </AnimatePresence>

      {/* Explosion particles */}
      <AnimatePresence>
        {isPopped && particles.length > 0 && particles.map(particle => (
          <Particle
            key={`particle-${particle.id}`}
            $color={particle.color}
            $size={particle.size}
            $glow={particle.glow}
            initial={{ 
              x: x + particle.initialX, 
              y: y + particle.initialY,
              opacity: 1,
              scale: 1
            }}
            animate={{ 
              x: x + particle.initialX + Math.cos(particle.angle) * particle.distance * particle.speed,
              y: y + particle.initialY + Math.sin(particle.angle) * particle.distance * particle.speed,
              opacity: 0,
              scale: 0,
              transition: {
                duration: particleLifetime / 1000,
                delay: particle.delay / 1000,
                ease: [0.32, 0.72, 0, 1] // Custom easing for more natural movement
              }
            }}
          />
        ))}
      </AnimatePresence>
    </>
  );
};
