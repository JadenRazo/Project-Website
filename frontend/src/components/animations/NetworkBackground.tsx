// src/components/animations/NetworkBackground.tsx
import React, { useEffect, useRef, useCallback, useState } from 'react';
import styled, { css, keyframes } from 'styled-components';
import { useTheme } from '../../contexts/ThemeContext';

// Enhanced configuration with better visibility
const ANIMATION_CONFIG = {
  MAX_PARTICLES: 120,          // Increased for better visibility
  MIN_PARTICLES: 50,           // Higher minimum for guaranteed visibility
  DENSITY_FACTOR: 15000,       // Adjusted for better particle density
  CONNECTION_DISTANCE: 150,    // Increased connection distance
  HOVER_DETECTION_RADIUS: 60,
  GLOW_DECAY_RATE: 0.02,
  BASE_OPACITY: 0.6,           // Higher base opacity
  MAX_OPACITY: 0.85,           // Higher max opacity
  GLOW_INTENSITY: 0.8,         // Increased glow intensity
  PARTICLE_MIN_SIZE: 1.5,      // Larger minimum particle size
  PARTICLE_MAX_SIZE: 3,        // Larger maximum particle size
  PARTICLE_OPACITY: 0.8,       // Higher particle opacity
  PARTICLE_SPEED: 0.3,
  LINE_WIDTH: 1.2,             // Thicker lines
  DEBUG_MODE: false,           // Set to true to see debugging info
};

// Enhanced animations
const pulseGlow = keyframes`
  0% { filter: blur(10px) brightness(1); }
  50% { filter: blur(15px) brightness(1.3); }
  100% { filter: blur(10px) brightness(1); }
`;

// Fixed and enhanced styled components with emphasis on full viewport coverage
const Container = styled.div`
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 2;               // Behind content but visible (critical fix)
  overflow: hidden;
  pointer-events: none;      // Allow clicks to pass through to content
  // Debug outline
  ${ANIMATION_CONFIG.DEBUG_MODE && css`
    border: 2px solid red;
  `}
`;

const Canvas = styled.canvas`
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  display: block;
  opacity: 1;                // Ensure full opacity
  // Debug outline
  ${ANIMATION_CONFIG.DEBUG_MODE && css`
    border: 1px solid yellow;
  `}
`;

const GlowEffect = styled.div<{ active: boolean; color: string }>`
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
  transition: opacity 0.3s ease;
  opacity: ${({ active }) => (active ? 0.2 : 0)};  // Increased opacity
  ${({ active, color }) =>
    active &&
    css`
      background-image: radial-gradient(
        circle at var(--mouse-x, 50%) var(--mouse-y, 50%), 
        ${color}60 0%, 
        ${color}30 30%, 
        transparent 70%
      );
      animation: ${pulseGlow} 3s ease-in-out infinite;
    `}
  // Debug outline
  ${ANIMATION_CONFIG.DEBUG_MODE && css`
    border: 1px solid green;
  `}
`;

// More transparent background overlay to ensure content is visible
const BackgroundOverlay = styled.div`
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: ${({ theme }) => `${theme.colors.background}60`};  // Reduced opacity further
  pointer-events: none;
  // Debug outline
  ${ANIMATION_CONFIG.DEBUG_MODE && css`
    border: 1px solid blue;
  `}
`;

// Debug panel for troubleshooting
const DebugPanel = styled.div`
  position: absolute;
  top: 10px;
  left: 10px;
  background: rgba(0, 0, 0, 0.7);
  color: #fff;
  padding: 10px;
  border-radius: 5px;
  font-family: monospace;
  font-size: 12px;
  z-index: 100;
  max-width: 300px;
  pointer-events: none;
`;

// Types
interface Point {
  x: number;
  y: number;
}

interface DebugInfo {
  fps: number;
  particles: number;
  connections: number;
  canvasWidth: number;
  canvasHeight: number;
  isActive: boolean;
  glowIntensity: number;
}

class Particle {
  private x: number;
  private y: number;
  private size: number;
  private speedX: number;
  private speedY: number;
  private readonly color: string;
  private readonly opacity: number;
  private readonly canvas: HTMLCanvasElement;

  constructor(canvas: HTMLCanvasElement, color: string) {
    this.canvas = canvas;
    this.color = color;
    this.opacity = ANIMATION_CONFIG.PARTICLE_OPACITY;
    
    // Distribute particles more evenly across the entire canvas
    this.x = Math.random() * canvas.width;
    this.y = Math.random() * canvas.height;
    
    // Random size within config limits
    this.size = Math.random() * 
      (ANIMATION_CONFIG.PARTICLE_MAX_SIZE - ANIMATION_CONFIG.PARTICLE_MIN_SIZE) + 
      ANIMATION_CONFIG.PARTICLE_MIN_SIZE;
    
    // Random speed with direction
    this.speedX = (Math.random() - 0.5) * 2 * ANIMATION_CONFIG.PARTICLE_SPEED;
    this.speedY = (Math.random() - 0.5) * 2 * ANIMATION_CONFIG.PARTICLE_SPEED;
  }

  update(): void {
    // Update position
    this.x += this.speedX;
    this.y += this.speedY;

    // Wrap around screen with margin to ensure smooth transitions
    const margin = 50;
    
    if (this.x > this.canvas.width + margin) {
      this.x = -margin;
    } else if (this.x < -margin) {
      this.x = this.canvas.width + margin;
    }
    
    if (this.y > this.canvas.height + margin) {
      this.y = -margin;
    } else if (this.y < -margin) {
      this.y = this.canvas.height + margin;
    }
  }

  draw(ctx: CanvasRenderingContext2D): void {
    ctx.beginPath();
    ctx.arc(this.x, this.y, this.size, 0, Math.PI * 2);
    ctx.fillStyle = `${this.color}${Math.floor(this.opacity * 255).toString(16).padStart(2, '0')}`;
    ctx.fill();
  }

  getPosition(): Point {
    return { x: this.x, y: this.y };
  }
}

class NetworkController {
  private particles: Particle[] = [];
  private readonly canvas: HTMLCanvasElement;
  private readonly ctx: CanvasRenderingContext2D;
  private readonly color: string;
  private mousePos: Point = { x: 0, y: 0 };
  private isActive: boolean = false;
  private glowIntensity: number = 0;
  private lastTime: number = 0;
  private fps: number = 0;
  private connectionCount: number = 0;

  constructor(canvas: HTMLCanvasElement, color: string) {
    this.canvas = canvas;
    const ctx = canvas.getContext('2d', { alpha: true });
    if (!ctx) throw new Error('Canvas context not available');
    this.ctx = ctx;
    this.color = color;
    this.lastTime = performance.now();
    this.initializeParticles();
  }

  private initializeParticles(): void {
    // Clear existing particles
    this.particles = [];
    
    // Calculate appropriate number of particles based on screen size
    const area = this.canvas.width * this.canvas.height;
    const count = Math.min(
      Math.max(
        Math.floor(area / ANIMATION_CONFIG.DENSITY_FACTOR), 
        ANIMATION_CONFIG.MIN_PARTICLES
      ), 
      ANIMATION_CONFIG.MAX_PARTICLES
    );
    
    // Create particles
    for (let i = 0; i < count; i++) {
      this.particles.push(new Particle(this.canvas, this.color));
    }
  }

  public updateMousePosition(x: number, y: number): void {
    this.mousePos = { x, y };
    this.isActive = true;
  }

  public mouseleave(): void {
    this.isActive = false;
  }

  private calculateGlowIntensity(): void {
    if (!this.isActive) {
      // Decay glow when mouse is inactive
      this.glowIntensity = Math.max(0, this.glowIntensity - ANIMATION_CONFIG.GLOW_DECAY_RATE);
      return;
    }

    // Check if mouse is near any connection
    let isNearConnection = false;
    
    for (let i = 0; i < this.particles.length; i++) {
      for (let j = i + 1; j < this.particles.length; j++) {
        const p1 = this.particles[i].getPosition();
        const p2 = this.particles[j].getPosition();
        const distance = this.getDistance(p1, p2);
        
        if (distance < ANIMATION_CONFIG.CONNECTION_DISTANCE) {
          // Calculate midpoint of connection
          const midpoint = {
            x: (p1.x + p2.x) / 2,
            y: (p1.y + p2.y) / 2
          };
          
          // Check if mouse is near this connection
          const mouseDistance = this.getDistance(this.mousePos, midpoint);
          if (mouseDistance < ANIMATION_CONFIG.HOVER_DETECTION_RADIUS) {
            isNearConnection = true;
            break;
          }
        }
      }
      if (isNearConnection) break;
    }
    
    // Update glow intensity
    if (isNearConnection) {
      this.glowIntensity = Math.min(
        this.glowIntensity + 0.05, 
        ANIMATION_CONFIG.GLOW_INTENSITY
      );
    } else {
      this.glowIntensity = Math.max(
        0, 
        this.glowIntensity - ANIMATION_CONFIG.GLOW_DECAY_RATE
      );
    }
  }

  private getDistance(p1: Point, p2: Point): number {
    const dx = p1.x - p2.x;
    const dy = p1.y - p2.y;
    return Math.sqrt(dx * dx + dy * dy);
  }

  public draw(): void {
    const currentTime = performance.now();
    const deltaTime = currentTime - this.lastTime;
    this.lastTime = currentTime;
    
    // Update FPS calculation
    this.fps = Math.round(1000 / (deltaTime || 1));
    
    // Clear canvas with transparent background
    this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);
    
    // Reset connection counter
    this.connectionCount = 0;
    
    // Update and draw particles
    this.particles.forEach(particle => {
      particle.update();
      particle.draw(this.ctx);
    });
    
    // Calculate glow
    this.calculateGlowIntensity();
    
    // Draw connections
    this.drawConnections();
  }

  private drawConnections(): void {
    // Set up context for better performance
    this.ctx.lineCap = 'round';
    
    for (let i = 0; i < this.particles.length; i++) {
      const p1 = this.particles[i].getPosition();
      
      for (let j = i + 1; j < this.particles.length; j++) {
        const p2 = this.particles[j].getPosition();
        const distance = this.getDistance(p1, p2);
        
        // Only draw connections within range
        if (distance < ANIMATION_CONFIG.CONNECTION_DISTANCE) {
          this.connectionCount++;
          
          // Calculate connection opacity based on distance
          const opacity = 
            (1 - distance / ANIMATION_CONFIG.CONNECTION_DISTANCE) * 
            ANIMATION_CONFIG.BASE_OPACITY;
          
          // Get midpoint for mouse proximity check
          const midpoint = {
            x: (p1.x + p2.x) / 2,
            y: (p1.y + p2.y) / 2
          };
          
          // Check if mouse is near this connection
          const mouseDistance = this.getDistance(this.mousePos, midpoint);
          const isHovered = mouseDistance < ANIMATION_CONFIG.HOVER_DETECTION_RADIUS && this.isActive;
          
          // Draw glow if hovered or global glow is active
          if (isHovered || this.glowIntensity > 0.1) {
            this.ctx.beginPath();
            this.ctx.strokeStyle = `${this.color}${Math.floor((isHovered ? 0.8 : this.glowIntensity * 0.5) * 255).toString(16).padStart(2, '0')}`;
            this.ctx.lineWidth = isHovered ? 2.5 : 1.8;
            this.ctx.shadowBlur = isHovered ? 18 : 12 * this.glowIntensity;
            this.ctx.shadowColor = this.color;
            this.ctx.moveTo(p1.x, p1.y);
            this.ctx.lineTo(p2.x, p2.y);
            this.ctx.stroke();
            
            // Reset shadow for better performance
            this.ctx.shadowBlur = 0;
          }
          
          // Draw regular connection
          this.ctx.beginPath();
          this.ctx.strokeStyle = `${this.color}${Math.floor((opacity + 0.1) * 255).toString(16).padStart(2, '0')}`;
          this.ctx.lineWidth = ANIMATION_CONFIG.LINE_WIDTH;
          this.ctx.moveTo(p1.x, p1.y);
          this.ctx.lineTo(p2.x, p2.y);
          this.ctx.stroke();
        }
      }
    }
  }
  
  public resize(width: number, height: number): void {
    // Update canvas size
    this.canvas.width = width;
    this.canvas.height = height;
    
    // Reinitialize particles
    this.initializeParticles();
  }
  
  public getGlowIntensity(): number {
    return this.glowIntensity;
  }
  
  public getDebugInfo(): DebugInfo {
    return {
      fps: this.fps,
      particles: this.particles.length,
      connections: this.connectionCount,
      canvasWidth: this.canvas.width,
      canvasHeight: this.canvas.height,
      isActive: this.isActive,
      glowIntensity: this.glowIntensity
    };
  }
}

export const NetworkBackground: React.FC = () => {
  // Setup refs for canvas and animation controller
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const networkRef = useRef<NetworkController | null>(null);
  const animationFrameRef = useRef<number>(0);
  const containerRef = useRef<HTMLDivElement>(null);
  
  // State for glow effect tracking
  const [isGlowActive, setIsGlowActive] = useState<boolean>(false);
  const [mousePos, setMousePos] = useState<{ x: number; y: number }>({ x: 50, y: 50 });
  const [debugInfo, setDebugInfo] = useState<DebugInfo | null>(null);
  
  // Get theme colors
  const { theme } = useTheme();
  const primaryColor = theme.colors.primary;
  
  // Mouse event handlers
  const handleMouseMove = useCallback((e: MouseEvent) => {
    const x = e.clientX;
    const y = e.clientY;
    
    // Update mouse position for glow effect CSS variables
    setMousePos({ x, y });
    
    // Update controller if available
    if (networkRef.current) {
      networkRef.current.updateMousePosition(x, y);
    }
  }, []);
  
  const handleMouseLeave = useCallback(() => {
    if (networkRef.current) {
      networkRef.current.mouseleave();
    }
  }, []);

  // Touch event handlers for mobile devices
  const handleTouchStart = useCallback((e: TouchEvent) => {
    if (e.touches.length > 0) {
      const touch = e.touches[0];
      const x = touch.clientX;
      const y = touch.clientY;
      
      setMousePos({ x, y });
      
      if (networkRef.current) {
        networkRef.current.updateMousePosition(x, y);
      }
    }
  }, []);
  
  const handleTouchMove = useCallback((e: TouchEvent) => {
    if (e.touches.length > 0) {
      const touch = e.touches[0];
      const x = touch.clientX;
      const y = touch.clientY;
      
      setMousePos({ x, y });
      
      if (networkRef.current) {
        networkRef.current.updateMousePosition(x, y);
      }
    }
  }, []);
  
  const handleTouchEnd = useCallback(() => {
    if (networkRef.current) {
      networkRef.current.mouseleave();
    }
  }, []);
  
  // Window resize handler with debounce
  const handleResize = useCallback(() => {
    if (canvasRef.current && networkRef.current && containerRef.current) {
      // Get window dimensions directly to ensure full viewport coverage
      const width = window.innerWidth;
      const height = window.innerHeight;
      
      // Ensure we're using the correct DPR for crisp rendering
      const dpr = window.devicePixelRatio || 1;
      
      // Update canvas dimensions for full viewport
      canvasRef.current.width = width * dpr;
      canvasRef.current.height = height * dpr;
      
      // Scale context to match DPR
      const ctx = canvasRef.current.getContext('2d');
      if (ctx) {
        ctx.scale(dpr, dpr);
      }
      
      // Update container style to explicitly match viewport
      containerRef.current.style.width = `${width}px`;
      containerRef.current.style.height = `${height}px`;
      
      // Resize the network
      networkRef.current.resize(width, height);
      
      console.log(`NetworkBackground: Resized to ${width}x${height}`);
    }
  }, []);
  
  // Animation loop
  const animate = useCallback(() => {
    if (networkRef.current) {
      networkRef.current.draw();
      
      // Check if glow should be active
      const intensity = networkRef.current.getGlowIntensity();
      setIsGlowActive(intensity > 0.1);
      
      // Update debug info
      if (ANIMATION_CONFIG.DEBUG_MODE) {
        setDebugInfo(networkRef.current.getDebugInfo());
      }
    }
    
    // Continue animation loop
    animationFrameRef.current = requestAnimationFrame(animate);
  }, []);
  
  // Initialize and cleanup
  useEffect(() => {
    console.log('NetworkBackground: Initializing...');
    
    // Get canvas element
    const canvas = canvasRef.current;
    const container = containerRef.current;
    
    if (!canvas || !container) {
      console.error('NetworkBackground: Canvas or container ref is null');
      return;
    }
    
    try {
      // Get window dimensions directly to ensure full viewport coverage
      const width = window.innerWidth;
      const height = window.innerHeight;
      
      console.log(`NetworkBackground: Window dimensions - ${width}x${height}`);
      
      // Set container dimensions explicitly to match viewport
      container.style.width = `${width}px`;
      container.style.height = `${height}px`;
      
      // Ensure we're using the correct DPR for crisp rendering
      const dpr = window.devicePixelRatio || 1;
      
      // Set canvas dimensions
      canvas.width = width * dpr;
      canvas.height = height * dpr;
      
      // Scale context to match DPR
      const ctx = canvas.getContext('2d');
      if (ctx) {
        ctx.scale(dpr, dpr);
      }
      
      // Create network controller
      networkRef.current = new NetworkController(canvas, primaryColor);
      console.log('NetworkBackground: Controller initialized');
      
      // Force an initial resize to ensure everything is sized correctly
      handleResize();
      
      // Add event listeners
      window.addEventListener('mousemove', handleMouseMove, { passive: true });
      document.body.addEventListener('mouseleave', handleMouseLeave, { passive: true });
      window.addEventListener('touchstart', handleTouchStart as unknown as EventListener, { passive: true });
      window.addEventListener('touchmove', handleTouchMove as unknown as EventListener, { passive: true });
      window.addEventListener('touchend', handleTouchEnd, { passive: true });
      
      // Add resize observer for more reliable size updates
      const resizeObserver = new ResizeObserver(() => {
        handleResize();
      });
      
      resizeObserver.observe(document.body);
      
      // Also add window resize event with debounce
      let resizeTimeout: NodeJS.Timeout;
      const debouncedResize = () => {
        clearTimeout(resizeTimeout);
        resizeTimeout = setTimeout(handleResize, 200);
      };
      
      window.addEventListener('resize', debouncedResize);
      
      // Move to the end of the event queue to ensure DOM is fully rendered
      setTimeout(() => {
        handleResize();
        
        // Start animation loop
        animate();
        console.log('NetworkBackground: Animation started');
      }, 0);
      
      // Cleanup function
      return () => {
        console.log('NetworkBackground: Cleaning up...');
        cancelAnimationFrame(animationFrameRef.current);
        window.removeEventListener('mousemove', handleMouseMove);
        document.body.removeEventListener('mouseleave', handleMouseLeave);
        window.removeEventListener('touchstart', handleTouchStart as unknown as EventListener);
        window.removeEventListener('touchmove', handleTouchMove as unknown as EventListener);
        window.removeEventListener('touchend', handleTouchEnd);
        window.removeEventListener('resize', debouncedResize);
        resizeObserver.disconnect();
        clearTimeout(resizeTimeout);
      };
    } catch (error) {
      console.error('NetworkBackground: Failed to initialize:', error);
      return () => {}; // Empty cleanup if initialization failed
    }
  }, [primaryColor, handleMouseMove, handleMouseLeave, handleTouchStart, handleTouchMove, handleTouchEnd, handleResize, animate]);
  
  // Set CSS variables for glow effect positioning
  useEffect(() => {
    document.documentElement.style.setProperty('--mouse-x', `${mousePos.x}px`);
    document.documentElement.style.setProperty('--mouse-y', `${mousePos.y}px`);
  }, [mousePos]);

  return (
    <Container ref={containerRef}>
      <BackgroundOverlay theme={theme} />
      <Canvas ref={canvasRef} />
      <GlowEffect 
        active={isGlowActive} 
        color={primaryColor}
      />
      
      {ANIMATION_CONFIG.DEBUG_MODE && debugInfo && (
        <DebugPanel>
          <div>FPS: {debugInfo.fps}</div>
          <div>Particles: {debugInfo.particles}</div>
          <div>Connections: {debugInfo.connections}</div>
          <div>Canvas: {debugInfo.canvasWidth}x{debugInfo.canvasHeight}</div>
          <div>Active: {debugInfo.isActive ? 'Yes' : 'No'}</div>
          <div>Glow: {debugInfo.glowIntensity.toFixed(2)}</div>
        </DebugPanel>
      )}
    </Container>
  );
};

export default NetworkBackground;
