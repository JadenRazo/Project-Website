import React, { useEffect, useRef, useState } from 'react';
import styled from 'styled-components';

// Animation container with guaranteed visibility across the entire viewport
const AnimationContainer = styled.div`
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  pointer-events: none;
  z-index: -5; /* Low enough to be behind all content */
  overflow: hidden;
  will-change: transform; /* Optimize for hardware acceleration */
  transform: translateZ(0); /* Force GPU rendering */
  backface-visibility: hidden; /* Improve performance */
`;

const Canvas = styled.canvas`
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  display: block;
`;

// Animation configuration
const ANIMATION_CONFIG = {
  PARTICLE_COUNT: 80,
  PARTICLE_SIZE_MIN: 1.5,
  PARTICLE_SIZE_MAX: 4,
  PARTICLE_SPEED: 0.4,
  CONNECTION_DISTANCE: 150,
  BACKGROUND_COLOR: 'transparent', // Ensure transparent background
  PARTICLE_COLOR: '#6c63ff', // Purple particle color
  CONNECTION_COLOR: '#6c63ff80', // Semi-transparent connection color
  FRAME_SKIP_RATE: 1, // 1 = render every frame, 2 = every second frame
};

interface Particle {
  x: number;
  y: number;
  size: number;
  speedX: number;
  speedY: number;
  opacity: number;
}

export const AnimationBackground: React.FC = () => {
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const animationFrameRef = useRef<number | null>(null);
  const particlesRef = useRef<Particle[]>([]);
  const [dimensions, setDimensions] = useState({ width: 0, height: 0 });
  const frameCountRef = useRef(0);

  // Initialize animation on mount
  useEffect(() => {
    if (!canvasRef.current || !containerRef.current) return;

    // Ensure container is visible at all times
    if (containerRef.current) {
      containerRef.current.style.visibility = 'visible';
      containerRef.current.style.display = 'block';
      
      // Debug - log position to help with troubleshooting
      if (process.env.NODE_ENV === 'development') {
        const rect = containerRef.current.getBoundingClientRect();
        console.log('Animation container position:', {
          top: rect.top,
          left: rect.left,
          width: rect.width,
          height: rect.height,
          zIndex: getComputedStyle(containerRef.current).zIndex
        });
      }
    }

    // Measure container
    const updateDimensions = () => {
      if (!containerRef.current) return;
      
      const { width, height } = containerRef.current.getBoundingClientRect();
      setDimensions({ width, height });
    };

    // Initialize dimensions
    updateDimensions();

    // Set up resize observer
    const resizeObserver = new ResizeObserver(() => {
      updateDimensions();
    });

    resizeObserver.observe(containerRef.current);

    // Clean up on unmount
    return () => {
      if (animationFrameRef.current !== null) {
        cancelAnimationFrame(animationFrameRef.current);
      }
      resizeObserver.disconnect();
    };
  }, []);

  // Update and draw particles
  const updateAndDrawParticles = (ctx: CanvasRenderingContext2D) => {
    particlesRef.current.forEach(particle => {
      // Update position
      particle.x += particle.speedX;
      particle.y += particle.speedY;
      
      // Wrap around edges
      if (particle.x > dimensions.width + particle.size) {
        particle.x = -particle.size;
      } else if (particle.x < -particle.size) {
        particle.x = dimensions.width + particle.size;
      }
      
      if (particle.y > dimensions.height + particle.size) {
        particle.y = -particle.size;
      } else if (particle.y < -particle.size) {
        particle.y = dimensions.height + particle.size;
      }
      
      // Draw particle
      ctx.beginPath();
      ctx.arc(particle.x, particle.y, particle.size, 0, Math.PI * 2);
      
      // Use particle opacity to make it more interesting
      const opacity = Math.floor(particle.opacity * 255).toString(16).padStart(2, '0');
      ctx.fillStyle = `${ANIMATION_CONFIG.PARTICLE_COLOR}${opacity}`;
      ctx.fill();
    });
  };

  // Draw connections between particles
  const drawConnections = (ctx: CanvasRenderingContext2D) => {
    ctx.strokeStyle = ANIMATION_CONFIG.CONNECTION_COLOR;
    ctx.lineWidth = 0.8;
    
    for (let i = 0; i < particlesRef.current.length; i++) {
      const p1 = particlesRef.current[i];
      
      for (let j = i + 1; j < particlesRef.current.length; j++) {
        const p2 = particlesRef.current[j];
        
        // Calculate distance between particles
        const dx = p1.x - p2.x;
        const dy = p1.y - p2.y;
        const distance = Math.sqrt(dx * dx + dy * dy);
        
        // Only draw connections if particles are close enough
        if (distance < ANIMATION_CONFIG.CONNECTION_DISTANCE) {
          // Make connections fade with distance
          const opacity = 1 - distance / ANIMATION_CONFIG.CONNECTION_DISTANCE;
          ctx.globalAlpha = opacity * 0.5; // Make connections semi-transparent
          
          ctx.beginPath();
          ctx.moveTo(p1.x, p1.y);
          ctx.lineTo(p2.x, p2.y);
          ctx.stroke();
        }
      }
    }
    
    // Reset global alpha
    ctx.globalAlpha = 1;
  };

  // Initialize particles when dimensions change
  useEffect(() => {
    if (dimensions.width === 0 || dimensions.height === 0) return;

    // Initialize canvas with correct dimensions
    const canvas = canvasRef.current;
    if (!canvas) return;

    // Get device pixel ratio for high DPI displays
    const dpr = window.devicePixelRatio || 1;
    canvas.width = dimensions.width * dpr;
    canvas.height = dimensions.height * dpr;
    
    // Set CSS dimensions
    canvas.style.width = `${dimensions.width}px`;
    canvas.style.height = `${dimensions.height}px`;

    // Initialize particles
    particlesRef.current = Array.from({ length: ANIMATION_CONFIG.PARTICLE_COUNT }, () => {
      const size = ANIMATION_CONFIG.PARTICLE_SIZE_MIN + 
        Math.random() * (ANIMATION_CONFIG.PARTICLE_SIZE_MAX - ANIMATION_CONFIG.PARTICLE_SIZE_MIN);
      
      // Distribute particles across the entire canvas
      return {
        x: Math.random() * dimensions.width,
        y: Math.random() * dimensions.height,
        size,
        speedX: (Math.random() - 0.5) * ANIMATION_CONFIG.PARTICLE_SPEED,
        speedY: (Math.random() - 0.5) * ANIMATION_CONFIG.PARTICLE_SPEED,
        opacity: 0.3 + Math.random() * 0.5, // Random opacity between 0.3 and 0.8
      };
    });

    // Start animation
    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    // Scale for high DPI displays
    ctx.scale(dpr, dpr);

    // Start animation loop
    const animate = () => {
      frameCountRef.current += 1;
      
      // Skip frames if needed for performance
      if (frameCountRef.current % ANIMATION_CONFIG.FRAME_SKIP_RATE !== 0) {
        animationFrameRef.current = requestAnimationFrame(animate);
        return;
      }
      
      // Clear canvas with transparent background
      ctx.clearRect(0, 0, dimensions.width, dimensions.height);
      
      // Update and draw connections first
      drawConnections(ctx);
      
      // Update and draw particles
      updateAndDrawParticles(ctx);
      
      animationFrameRef.current = requestAnimationFrame(animate);
    };
    
    animate();
    
    // Clean up current animation before starting a new one
    return () => {
      if (animationFrameRef.current !== null) {
        cancelAnimationFrame(animationFrameRef.current);
      }
    };
  }, [dimensions, updateAndDrawParticles, drawConnections]);

  return (
    <AnimationContainer ref={containerRef} className="animation-background">
      <Canvas ref={canvasRef} />
    </AnimationContainer>
  );
};

export default AnimationBackground; 