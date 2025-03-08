import React, { useRef, useState, useEffect, useCallback, useMemo } from 'react';
import { motion, useScroll, useTransform, useSpring, useMotionValueEvent, useAnimation, MotionValue } from 'framer-motion';
import { useTheme } from '../../contexts/ThemeContext';
import { useZIndex } from '../../contexts/ZIndexContext';
import styled from 'styled-components';
import { useInView } from 'react-intersection-observer';
import { debounce } from 'lodash';

// Define types for scroll transforms
interface ScrollTransformSection {
  id: string;
  startPercent: number;  // When to start the effect (0-1)
  endPercent: number;    // When to end the effect (0-1)
  transforms: ScrollTransformEffect[];
}

interface ScrollTransformEffect {
  type: 'scale' | 'opacity' | 'translateX' | 'translateY' | 'rotate' | 'blur' | 'color';
  target: string;        // CSS selector or element ID
  from: number | string;
  to: number | string;
  unit?: string;         // px, %, deg, etc.
  easing?: string;       // "linear", "easeIn", "easeOut", "easeInOut"
}

interface DeviceCapabilities {
  isHighPerformance: boolean;
  isTouch: boolean;
  prefersReducedMotion: boolean;
  isMobile: boolean;
  browserSupportsBackdropFilter: boolean;
}

interface Point {
  x: number;
  y: number;
}

interface TransformElement {
  id: string;
  element: HTMLElement;
  initialState: any;
  effects: ScrollTransformEffect[];
  controls: any;
}

interface DebugInfo {
  fps: number;
  transformElements: number;
  activeTransforms: number;
  viewportWidth: number;
  viewportHeight: number;
  scrollPosition: number;
  isActive: boolean;
  gravityFactor: number;
}

interface Orb {
  id: string;
  x: number;
  y: number;
  size: number;
  speed: number;
  color: string;
  opacity: number;
  popped: boolean;
  rotationSpeed: number;
  currentRotation: number;
  type: 'circle' | 'square' | 'triangle';
}

// Controller class for handling scroll transformations
class ScrollTransformController {
  private readonly transformElements: Map<string, TransformElement> = new Map();
  private readonly scrollYProgress: MotionValue<number>;
  private deviceCapabilities: DeviceCapabilities;
  private sections: ScrollTransformSection[] = [];
  private isActive: boolean = true;
  private observer: IntersectionObserver | null = null;
  private container: HTMLElement | null = null;
  private lastScrollY: number = 0;
  private lastTime: number = 0;
  private fps: number = 60;
  private activeTransforms: number = 0;
  private gravityFactor: number = 1;
  
  constructor(scrollYProgress: MotionValue<number>, deviceCapabilities: DeviceCapabilities) {
    this.scrollYProgress = scrollYProgress;
    this.deviceCapabilities = deviceCapabilities;
    this.lastTime = performance.now();
  }

  public setContainer(container: HTMLElement): void {
    this.container = container;
    
    // Setup intersection observer to disable animations when not visible
    this.observer = new IntersectionObserver(
      (entries) => {
        entries.forEach(entry => {
          this.isActive = entry.isIntersecting;
        });
      },
      { threshold: 0.1 }
    );
    
    this.observer.observe(container);
  }

  public setSections(sections: ScrollTransformSection[]): void {
    // Filter out complex effects on mobile/low-performance devices
    if (this.deviceCapabilities.isTouch || this.deviceCapabilities.isMobile || !this.deviceCapabilities.isHighPerformance) {
      sections = this.optimizeSections(sections);
    }
    
    this.sections = sections;
    this.initializeElements();
  }

  private optimizeSections(sections: ScrollTransformSection[]): ScrollTransformSection[] {
    return sections.map(section => ({
      ...section,
      transforms: section.transforms.filter(transform => 
        // Keep only simple transforms on lower-end devices
        ['opacity', 'translateY'].includes(transform.type)
      )
    }));
  }

  private initializeElements(): void {
    if (!this.container) return;
    
    // Clear existing elements
    this.transformElements.clear();
    
    // Get all targets from all sections
    const allTargets = new Set<string>();
    this.sections.forEach(section => {
      section.transforms.forEach(transform => {
        allTargets.add(transform.target);
      });
    });
    
    // Initialize each element
    allTargets.forEach(target => {
      let elements: NodeListOf<Element>;
      
      try {
        // Try to find elements within document
        elements = document.querySelectorAll(target);
        
        if (elements.length === 0 && this.container) {
          // Try to find elements within container
          elements = this.container.querySelectorAll(target);
        }
      } catch (err) {
        console.warn(`Invalid selector: ${target}`);
        return;
      }
      
      elements.forEach((el, index) => {
        const element = el as HTMLElement;
        const id = `${target.replace(/[^\w-]/g, '')}-${index}`;
        
        // Store the initial state of the element for reference
        const initialState = {
          opacity: parseFloat(window.getComputedStyle(element).opacity) || 1,
          transform: window.getComputedStyle(element).transform || 'none',
          filter: window.getComputedStyle(element).filter || 'none',
          color: window.getComputedStyle(element).color || 'inherit'
        };
        
        // Get all effects targeting this element
        const effects = this.sections.flatMap(section => 
          section.transforms.filter(transform => transform.target === target)
        );
        
        // Create animation controls for this element
        const controls = useAnimation();
        
        this.transformElements.set(id, {
          id,
          element,
          initialState,
          effects,
          controls
        });
      });
    });
  }

  public updateOnScroll(scrollY: number, viewportHeight: number, gravityFactor: number): void {
    if (!this.isActive || !this.container) return;
    
    // Calculate FPS
    const now = performance.now();
    const delta = now - this.lastTime;
    if (delta > 0) {
      this.fps = 1000 / delta;
    }
    this.lastTime = now;
    
    // Store last scroll position
    this.lastScrollY = scrollY;
    this.gravityFactor = gravityFactor;
    
    // Calculate scroll percentage (clamped between 0 and 1)
    const scrollPercent = Math.max(0, Math.min(1, scrollY / (document.body.scrollHeight - viewportHeight)));
    
    // Reset active transforms counter
    this.activeTransforms = 0;
    
    // Apply transforms from each section
    this.sections.forEach(section => {
      // Check if we're within this section's scroll range
      if (scrollPercent >= section.startPercent && scrollPercent <= section.endPercent) {
        // Calculate progress within this section
        const sectionProgress = (scrollPercent - section.startPercent) / 
          (section.endPercent - section.startPercent);
        
        // Apply transforms for this section
        this.applyTransforms(section, sectionProgress);
      }
    });
  }

  private applyTransforms(section: ScrollTransformSection, progress: number): void {
    section.transforms.forEach(transform => {
      let elements: NodeListOf<Element>;
      
      try {
        // Try document-wide
        elements = document.querySelectorAll(transform.target);
        
        // If still no elements found, try the container
        if (elements.length === 0 && this.container) {
          elements = this.container.querySelectorAll(transform.target);
        }
      } catch (err) {
        console.warn(`Invalid selector: ${transform.target}`);
        return;
      }
      
      elements.forEach((el) => {
        const element = el as HTMLElement;
        this.activeTransforms++;
        
        // Apply the easing function if specified
        let easedProgress = progress;
        if (transform.easing) {
          switch (transform.easing) {
            case 'easeIn':
              easedProgress = progress * progress;
              break;
            case 'easeOut':
              easedProgress = 1 - Math.pow(1 - progress, 2);
              break;
            case 'easeInOut':
              easedProgress = progress < 0.5 
                ? 2 * progress * progress 
                : 1 - Math.pow(-2 * progress + 2, 2) / 2;
              break;
          }
        }
        
        // Apply the appropriate transform
        this.applySpecificTransform(element, transform, easedProgress);
      });
    });
  }
  
  private applySpecificTransform(element: HTMLElement, transform: ScrollTransformEffect, progress: number): void {
    switch (transform.type) {
      case 'opacity':
        const fromOpacity = parseFloat(transform.from as string);
        const toOpacity = parseFloat(transform.to as string);
        element.style.opacity = String(fromOpacity + (toOpacity - fromOpacity) * progress);
        break;
        
      case 'translateY':
        const fromY = parseFloat(transform.from as string);
        const toY = parseFloat(transform.to as string);
        const unit = transform.unit || 'px';
        const value = fromY + (toY - fromY) * progress;
        
        this.applyTransformProperty(element, `translateY(${value}${unit})`);
        break;
        
      case 'translateX':
        const fromX = parseFloat(transform.from as string);
        const toX = parseFloat(transform.to as string);
        const unitX = transform.unit || 'px';
        const valueX = fromX + (toX - fromX) * progress;
        
        this.applyTransformProperty(element, `translateX(${valueX}${unitX})`);
        break;
        
      case 'scale':
        const fromScale = parseFloat(transform.from as string);
        const toScale = parseFloat(transform.to as string);
        const valueScale = fromScale + (toScale - fromScale) * progress;
        
        this.applyTransformProperty(element, `scale(${valueScale})`);
        break;
        
      case 'rotate':
        const fromRotate = parseFloat(transform.from as string);
        const toRotate = parseFloat(transform.to as string);
        const unitRotate = transform.unit || 'deg';
        const valueRotate = fromRotate + (toRotate - fromRotate) * progress;
        
        this.applyTransformProperty(element, `rotate(${valueRotate}${unitRotate})`);
        break;
        
      case 'blur':
        const fromBlur = parseFloat(transform.from as string);
        const toBlur = parseFloat(transform.to as string);
        const unitBlur = transform.unit || 'px';
        const valueBlur = fromBlur + (toBlur - fromBlur) * progress;
        
        // Preserve other filter values
        const currentFilter = window.getComputedStyle(element).filter;
        const hasExistingFilter = currentFilter && currentFilter !== 'none';
        
        if (hasExistingFilter && !currentFilter.includes('blur')) {
          element.style.filter = `${currentFilter} blur(${valueBlur}${unitBlur})`;
        } else {
          element.style.filter = `blur(${valueBlur}${unitBlur})`;
        }
        break;
        
      case 'color':
        if (typeof transform.from === 'string' && typeof transform.to === 'string') {
          // Use RGBA interpolation for better color transitions
          const fromColor = this.parseColor(transform.from);
          const toColor = this.parseColor(transform.to);
          
          if (fromColor && toColor) {
            const r = Math.round(fromColor.r + (toColor.r - fromColor.r) * progress);
            const g = Math.round(fromColor.g + (toColor.g - fromColor.g) * progress);
            const b = Math.round(fromColor.b + (toColor.b - fromColor.b) * progress);
            const a = fromColor.a + (toColor.a - fromColor.a) * progress;
            
            element.style.color = `rgba(${r}, ${g}, ${b}, ${a})`;
          } else {
            // Fallback if color parsing fails
            element.style.color = progress > 0.5 ? transform.to : transform.from;
          }
        }
        break;
    }
  }
  
  private applyTransformProperty(element: HTMLElement, newTransform: string): void {
    const currentTransform = window.getComputedStyle(element).transform;
    const hasExistingTransform = currentTransform && currentTransform !== 'none';
    
    if (!hasExistingTransform) {
      element.style.transform = newTransform;
      return;
    }
    
    // Extract the transform type from newTransform
    const transformType = newTransform.split('(')[0];
    
    // Combine with existing transforms, replacing the same type if it exists
    const transforms = this.parseTransforms(currentTransform);
    let updated = false;
    
    for (let i = 0; i < transforms.length; i++) {
      if (transforms[i].startsWith(transformType)) {
        transforms[i] = newTransform;
        updated = true;
        break;
      }
    }
    
    if (!updated) {
      transforms.push(newTransform);
    }
    
    element.style.transform = transforms.join(' ');
  }
  
  private parseTransforms(transformString: string): string[] {
    if (transformString === 'none') return [];
    
    // Handle matrix and multiple transforms
    const transforms: string[] = [];
    const regex = /([\w]+)\s*\(([^)]*)\)/g;
    let match;
    
    while ((match = regex.exec(transformString)) !== null) {
      transforms.push(`${match[1]}(${match[2]})`);
    }
    
    return transforms;
  }
  
  private parseColor(color: string): { r: number; g: number; b: number; a: number } | null {
    // Handle hex colors
    if (color.startsWith('#')) {
      const hex = color.substring(1);
      const r = parseInt(hex.substring(0, 2), 16);
      const g = parseInt(hex.substring(2, 4), 16);
      const b = parseInt(hex.substring(4, 6), 16);
      const a = hex.length === 8 ? parseInt(hex.substring(6, 8), 16) / 255 : 1;
      
      return { r, g, b, a };
    }
    
    // Handle rgba
    if (color.startsWith('rgba')) {
      const values = color.match(/rgba\((\d+),\s*(\d+),\s*(\d+),\s*([\d.]+)\)/);
      if (values) {
        return {
          r: parseInt(values[1], 10),
          g: parseInt(values[2], 10),
          b: parseInt(values[3], 10),
          a: parseFloat(values[4])
        };
      }
    }
    
    // Handle rgb
    if (color.startsWith('rgb')) {
      const values = color.match(/rgb\((\d+),\s*(\d+),\s*(\d+)\)/);
      if (values) {
        return {
          r: parseInt(values[1], 10),
          g: parseInt(values[2], 10),
          b: parseInt(values[3], 10),
          a: 1
        };
      }
    }
    
    return null;
  }

  public resize(): void {
    if (this.isActive) {
      this.initializeElements();
    }
  }

  public cleanup(): void {
    if (this.observer && this.container) {
      this.observer.unobserve(this.container);
      this.observer.disconnect();
    }
    
    this.transformElements.clear();
  }
  
  public getDebugInfo(): DebugInfo {
    return {
      fps: Math.round(this.fps),
      transformElements: this.transformElements.size,
      activeTransforms: this.activeTransforms,
      viewportWidth: window.innerWidth,
      viewportHeight: window.innerHeight,
      scrollPosition: this.lastScrollY,
      isActive: this.isActive,
      gravityFactor: this.gravityFactor
    };
  }
}

// Default transform sections for professional website
const websiteTransformSections: ScrollTransformSection[] = [
  {
    id: 'header-reveal',
    startPercent: 0,
    endPercent: 0.15,
    transforms: [
      {
        type: 'opacity',
        target: 'header',
        from: '0',
        to: '1',
        easing: 'easeOut'
      },
      {
        type: 'translateY',
        target: 'header',
        from: '-20',
        to: '0',
        unit: 'px',
        easing: 'easeOut'
      }
    ]
  },
  {
    id: 'hero-parallax',
    startPercent: 0,
    endPercent: 0.3,
    transforms: [
      {
        type: 'translateY',
        target: '#hero .hero-title',
        from: '0',
        to: '-50',
        unit: 'px',
        easing: 'easeInOut'
      },
      {
        type: 'opacity',
        target: '#hero .hero-subtitle',
        from: '1',
        to: '0.6',
        easing: 'easeInOut'
      },
      {
        type: 'scale',
        target: '#hero .hero-image',
        from: '1',
        to: '0.95',
        easing: 'easeInOut'
      }
    ]
  },
  {
    id: 'skills-reveal',
    startPercent: 0.1,
    endPercent: 0.4,
    transforms: [
      {
        type: 'opacity',
        target: '#skills .section-title',
        from: '0',
        to: '1',
        easing: 'easeOut'
      },
      {
        type: 'translateY',
        target: '#skills .section-title',
        from: '30',
        to: '0',
        unit: 'px',
        easing: 'easeOut'
      },
      {
        type: 'opacity',
        target: '#skills .skill-item:nth-child(1)',
        from: '0',
        to: '1',
        easing: 'easeOut'
      },
      {
        type: 'translateX',
        target: '#skills .skill-item:nth-child(1)',
        from: '-50',
        to: '0',
        unit: 'px',
        easing: 'easeOut'
      },
      {
        type: 'opacity',
        target: '#skills .skill-item:nth-child(2)',
        from: '0',
        to: '1',
        easing: 'easeOut'
      },
      {
        type: 'translateY',
        target: '#skills .skill-item:nth-child(2)',
        from: '50',
        to: '0',
        unit: 'px',
        easing: 'easeOut'
      },
      {
        type: 'opacity',
        target: '#skills .skill-item:nth-child(3)',
        from: '0',
        to: '1',
        easing: 'easeOut'
      },
      {
        type: 'translateX',
        target: '#skills .skill-item:nth-child(3)',
        from: '50',
        to: '0',
        unit: 'px',
        easing: 'easeOut'
      }
    ]
  },
  {
    id: 'projects-reveal',
    startPercent: 0.3,
    endPercent: 0.6,
    transforms: [
      {
        type: 'opacity',
        target: '#projects .section-title',
        from: '0',
        to: '1',
        easing: 'easeOut'
      },
      {
        type: 'translateY',
        target: '#projects .section-title',
        from: '30',
        to: '0',
        unit: 'px',
        easing: 'easeOut'
      },
      {
        type: 'opacity',
        target: '#projects .project-card',
        from: '0',
        to: '1',
        easing: 'easeOut'
      },
      {
        type: 'translateY',
        target: '#projects .project-card',
        from: '50',
        to: '0',
        unit: 'px',
        easing: 'easeOut'
      }
    ]
  },
  {
    id: 'about-reveal',
    startPercent: 0.5,
    endPercent: 0.8,
    transforms: [
      {
        type: 'opacity',
        target: '#about .section-title',
        from: '0',
        to: '1',
        easing: 'easeOut'
      },
      {
        type: 'translateY',
        target: '#about .section-title',
        from: '30',
        to: '0',
        unit: 'px',
        easing: 'easeOut'
      },
      {
        type: 'opacity',
        target: '#about .about-content',
        from: '0',
        to: '1',
        easing: 'easeOut'
      },
      {
        type: 'translateY',
        target: '#about .about-content',
        from: '50',
        to: '0',
        unit: 'px',
        easing: 'easeOut'
      }
    ]
  }
];

// Hook to detect device capabilities
const useDeviceCapabilities = (): DeviceCapabilities => {
  return useMemo(() => {
    const isTouch = 'ontouchstart' in window || navigator.maxTouchPoints > 0;
    const prefersReducedMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches;
    const isMobile = /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent);
    
    // Check if browser supports backdrop-filter
    const testEl = document.createElement('div');
    testEl.style.backdropFilter = 'blur(1px)';
    const browserSupportsBackdropFilter = testEl.style.backdropFilter.length > 0;
    
    // Determine if high performance device
    const isHighPerformance = !isMobile && !prefersReducedMotion && window.devicePixelRatio >= 1;
    
    return {
      isHighPerformance,
      isTouch,
      prefersReducedMotion,
      isMobile,
      browserSupportsBackdropFilter
    };
  }, []);
};

// Styled components
const BackgroundContainer = styled(motion.div)`
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100vh;
  z-index: 0;
  overflow: hidden;
  pointer-events: none;
`;

const BackgroundOverlay = styled(motion.div)<{ $reducedMotion: boolean }>`
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: radial-gradient(circle at center, transparent 0%, ${props => props.theme.colors.background || '#0f0f17'} 80%);
  opacity: ${props => props.$reducedMotion ? 0.9 : 0.7};
  mix-blend-mode: normal;
`;

const OrbsContainer = styled.div`
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  overflow: hidden;
  pointer-events: none;
  z-index: 1;
`;

const OrbElement = styled(motion.div)`
  position: absolute;
  border-radius: 50%;
  cursor: pointer;
  pointer-events: all;
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  
  &::after {
    content: "";
    position: absolute;
    width: 40%;
    height: 40%;
    top: 15%;
    left: 15%;
    background: rgba(255, 255, 255, 0.1);
    border-radius: 50%;
  }
`;

const SquareOrb = styled(motion.div)`
  position: absolute;
  border-radius: 6px;
  cursor: pointer;
  pointer-events: all;
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  
  &::after {
    content: "";
    position: absolute;
    width: 50%;
    height: 50%;
    top: 25%;
    left: 25%;
    background: rgba(255, 255, 255, 0.1);
    border-radius: 2px;
  }
`;

const TriangleOrb = styled(motion.div)`
  position: absolute;
  cursor: pointer;
  pointer-events: all;
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  clip-path: polygon(50% 0%, 100% 100%, 0% 100%);
  
  &::after {
    content: "";
    position: absolute;
    width: 35%;
    height: 35%;
    top: 40%;
    left: 32.5%;
    background: rgba(255, 255, 255, 0.1);
    clip-path: polygon(50% 10%, 90% 90%, 10% 90%);
  }
`;

const DebugPanel = styled(motion.div)`
  position: fixed;
  bottom: 20px;
  right: 20px;
  background: rgba(0, 0, 0, 0.8);
  color: white;
  padding: 12px;
  border-radius: 8px;
  font-family: monospace;
  font-size: 12px;
  z-index: 9999;
  pointer-events: all;
  max-width: 300px;
  user-select: none;
`;

// Main component
export const ScrollTransformBackground: React.FC<{
  customSections?: ScrollTransformSection[];
  showDebug?: boolean;
  enableFloatingOrbs?: boolean;
}> = ({ customSections, showDebug = false, enableFloatingOrbs = true }) => {
  // Refs
  const containerRef = useRef<HTMLDivElement>(null);
  const controllerRef = useRef<ScrollTransformController | null>(null);
  const animationFrameRef = useRef<number | null>(null);
  const lastFrameTime = useRef<number>(0);
  
  // Get custom hooks
  const { theme } = useTheme();
  const { zIndex } = useZIndex();
  const { scrollY, scrollYProgress } = useScroll();
  const deviceCapabilities = useDeviceCapabilities();
  
  // State
  const [debugInfo, setDebugInfo] = useState<DebugInfo | null>(null);
  const [orbs, setOrbs] = useState<Orb[]>([]);
  const [gravityFactor, setGravityFactor] = useState<number>(1);
  const prevScrollY = useRef<number>(0);
  
  // Initialize controller on mount
  useEffect(() => {
    if (controllerRef.current === null) {
      controllerRef.current = new ScrollTransformController(
        scrollYProgress,
        deviceCapabilities
      );
    }
    
    return () => {
      if (controllerRef.current) {
        controllerRef.current.cleanup();
      }
      
      if (animationFrameRef.current) {
        cancelAnimationFrame(animationFrameRef.current);
      }
    };
  }, [scrollYProgress, deviceCapabilities]);
  
  // Set container when ref is available
  useEffect(() => {
    if (containerRef.current && controllerRef.current) {
      controllerRef.current.setContainer(containerRef.current);
    }
  }, [containerRef.current]);
  
  // Set sections (either custom or default website sections)
  useEffect(() => {
    if (controllerRef.current) {
      controllerRef.current.setSections(customSections || websiteTransformSections);
    }
  }, [customSections, deviceCapabilities]);
  
  // Generate orbs with different shapes
  useEffect(() => {
    if (!enableFloatingOrbs) return;
    
    const generateOrbs = () => {
      if (deviceCapabilities.isMobile && !deviceCapabilities.isHighPerformance) {
        setOrbs([]);
        return;
      }
      
      const orbCount = deviceCapabilities.isHighPerformance ? 20 : 10;
      
      // Use theme colors for a cohesive look
      const primaryColor = theme.colors?.primary || '#6c63ff';
      const secondaryColor = theme.colors?.secondary || '#ff6b6b';
      const tertiaryColor = theme.colors?.accent || '#4ecdc4';
      
      const colorPalette = [primaryColor, secondaryColor, tertiaryColor];
      const shapeTypes: Array<'circle' | 'square' | 'triangle'> = ['circle', 'square', 'triangle'];
      
      const newOrbs = Array.from({ length: orbCount }).map((_, i) => ({
        id: `orb-${i}`,
        x: Math.random() * 100,
        y: Math.random() * 100,
        size: 20 + Math.random() * 60,
        speed: 0.1 + Math.random() * 0.3,
        color: colorPalette[i % colorPalette.length],
        opacity: 0.1 + Math.random() * 0.3,
        popped: false,
        rotationSpeed: (Math.random() - 0.5) * 0.5,
        currentRotation: Math.random() * 360,
        type: shapeTypes[i % shapeTypes.length]
      }));
      
      setOrbs(newOrbs);
    };
    
    generateOrbs();
    
    const handleResize = debounce(generateOrbs, 500);
    window.addEventListener('resize', handleResize);
    
    return () => {
      window.removeEventListener('resize', handleResize);
      if (animationFrameRef.current) {
        cancelAnimationFrame(animationFrameRef.current);
      }
    };
  }, [deviceCapabilities, theme.colors, enableFloatingOrbs]);
  
  // Optimized orb animation loop
  useEffect(() => {
    if (!enableFloatingOrbs || orbs.length === 0) return;
    
    const animateOrbs = (timestamp: number) => {
      const targetFPS = deviceCapabilities.isHighPerformance ? 60 : 30;
      const frameInterval = 1000 / targetFPS;
      
      if (timestamp - lastFrameTime.current < frameInterval) {
        animationFrameRef.current = requestAnimationFrame(animateOrbs);
        return;
      }
      
      lastFrameTime.current = timestamp;
      
      setOrbs(currentOrbs => 
        currentOrbs.map(orb => {
          if (orb.popped) return orb;
          
          // Create smooth, natural-looking movement patterns
          let newX = orb.x + Math.sin(timestamp * 0.0001 * orb.speed * 5) * 0.05;
          
          // Apply gravity factor to vertical movement
          let newY = orb.y - (orb.speed / gravityFactor) * 0.2;
          
          // Loop orbs when they exit the view
          if (newY < -10) {
            return {
              ...orb,
              y: 110,
              x: Math.random() * 100
            };
          }
          
          // Rotate orbs slightly
          const newRotation = (orb.currentRotation + orb.rotationSpeed) % 360;
          
          return { 
            ...orb, 
            x: newX, 
            y: newY,
            currentRotation: newRotation
          };
        })
      );
      
      animationFrameRef.current = requestAnimationFrame(animateOrbs);
    };
    
    animationFrameRef.current = requestAnimationFrame(animateOrbs);
    
    return () => {
      if (animationFrameRef.current) {
        cancelAnimationFrame(animationFrameRef.current);
      }
    };
  }, [orbs, deviceCapabilities, gravityFactor, enableFloatingOrbs]);
  
  // Handle scroll events for transformations and gravity
  useMotionValueEvent(scrollY, "change", (latest) => {
    if (controllerRef.current) {
      controllerRef.current.updateOnScroll(latest, window.innerHeight, gravityFactor);
      
      if (showDebug) {
        setDebugInfo(controllerRef.current.getDebugInfo());
      }
    }
    
    if (enableFloatingOrbs) {
      // Calculate scroll direction and adjust gravity
      const direction = latest > prevScrollY.current ? 1 : -1;
      const scrollSpeed = Math.abs(latest - prevScrollY.current);
      const speedFactor = Math.min(1, scrollSpeed / 50);
      
      // Apply more dramatic gravity changes with faster scrolling
      const newGravityFactor = direction > 0 ? 
        Math.min(3, gravityFactor + 0.1 * speedFactor) : 
        Math.max(0.3, gravityFactor - 0.2 * speedFactor);
      
      setGravityFactor(newGravityFactor);
      prevScrollY.current = latest;
    }
  });
  
  // Handle window resizing
  useEffect(() => {
    const handleResize = debounce(() => {
      if (controllerRef.current) {
        controllerRef.current.resize();
      }
    }, 100);
    
    window.addEventListener('resize', handleResize);
    return () => window.removeEventListener('resize', handleResize);
  }, []);
  
  // Handle orb interaction
  const handleOrbClick = useCallback((id: string) => {
    if (!enableFloatingOrbs) return;
    
    // Pop the orb with animation
    setOrbs(currentOrbs => 
      currentOrbs.map(orb => 
        orb.id === id ? 
          { ...orb, popped: true, opacity: 0 } : 
          orb
      )
    );
    
    // Create new orb after a delay
    setTimeout(() => {
      setOrbs(currentOrbs => {
        const poppedIndex = currentOrbs.findIndex(orb => orb.id === id);
        if (poppedIndex === -1) return currentOrbs;
        
        const primaryColor = theme.colors?.primary || '#6c63ff';
        const secondaryColor = theme.colors?.secondary || '#ff6b6b';
        const tertiaryColor = theme.colors?.accent || '#4ecdc4';
        const colorPalette = [primaryColor, secondaryColor, tertiaryColor];
        const shapeTypes: Array<'circle' | 'square' | 'triangle'> = ['circle', 'square', 'triangle'];
        
        const updatedOrbs = [...currentOrbs];
        updatedOrbs[poppedIndex] = {
          ...updatedOrbs[poppedIndex],
          popped: false,
          opacity: 0.1 + Math.random() * 0.3,
          x: Math.random() * 100,
          y: 110, // Start from bottom
          size: 20 + Math.random() * 60,
          color: colorPalette[Math.floor(Math.random() * colorPalette.length)],
          rotationSpeed: (Math.random() - 0.5) * 0.5,
          currentRotation: Math.random() * 360,
          type: shapeTypes[Math.floor(Math.random() * shapeTypes.length)]
        };
        
        return updatedOrbs;
      });
    }, 2000);
  }, [theme.colors]);
  
  // Render the background with interactive orbs
  return (
    <BackgroundContainer 
      ref={containerRef}
      style={{ zIndex: zIndex.PARTICLE_BACKGROUND }}
      className="scroll-transform-background"
    >
      {enableFloatingOrbs && (
  <OrbsContainer>
    {orbs.map(orb => {
      // Render different orb shapes based on type
      if (orb.type === 'circle') {
        return (
          <OrbElement
            key={orb.id}
            style={{
              left: `${orb.x}%`,
              top: `${orb.y}%`,
              width: `${orb.size}px`,
              height: `${orb.size}px`,
              opacity: orb.opacity,
              backgroundColor: `${orb.color}10`,
              boxShadow: `0 0 20px ${orb.color}30, inset 0 0 15px ${orb.color}30`,
              border: `1px solid ${orb.color}40`,
              transform: orb.popped 
                ? 'scale(0) rotate(0deg)' 
                : `scale(1) rotate(${orb.currentRotation}deg)`,
              transition: orb.popped 
                ? 'transform 0.5s cubic-bezier(0.34, 1.56, 0.64, 1), opacity 0.3s ease-out' 
                : 'none'
            }}
            onClick={() => handleOrbClick(orb.id)}
            whileHover={{ 
              scale: 1.1, 
              opacity: orb.opacity * 1.5,
              boxShadow: `0 0 30px ${orb.color}50, inset 0 0 20px ${orb.color}50`
            }}
            whileTap={{ scale: 0.9 }}
          />
        );
      } else if (orb.type === 'square') {
        return (
          <SquareOrb
            key={orb.id}
            style={{
              left: `${orb.x}%`,
              top: `${orb.y}%`,
              width: `${orb.size * 0.8}px`,
              height: `${orb.size * 0.8}px`,
              opacity: orb.opacity,
              backgroundColor: `${orb.color}10`,
              boxShadow: `0 0 20px ${orb.color}30, inset 0 0 15px ${orb.color}30`,
              border: `1px solid ${orb.color}40`,
              transform: orb.popped 
                ? 'scale(0) rotate(0deg)' 
                : `scale(1) rotate(${orb.currentRotation}deg)`,
              transition: orb.popped 
                ? 'transform 0.5s cubic-bezier(0.34, 1.56, 0.64, 1), opacity 0.3s ease-out' 
                : 'none'
            }}
            onClick={() => handleOrbClick(orb.id)}
            whileHover={{ 
              scale: 1.1, 
              opacity: orb.opacity * 1.5,
              boxShadow: `0 0 30px ${orb.color}50, inset 0 0 20px ${orb.color}50`
            }}
            whileTap={{ scale: 0.9 }}
          />
        );
      } else {
        return (
          <TriangleOrb
            key={orb.id}
            style={{
              left: `${orb.x}%`,
              top: `${orb.y}%`,
              width: `${orb.size}px`,
              height: `${orb.size}px`,
              opacity: orb.opacity,
              backgroundColor: `${orb.color}10`,
              boxShadow: `0 0 20px ${orb.color}30, inset 0 0 15px ${orb.color}30`,
              border: `1px solid ${orb.color}40`,
              transform: orb.popped 
                ? 'scale(0) rotate(0deg)' 
                : `scale(1) rotate(${orb.currentRotation}deg)`,
              transition: orb.popped 
                ? 'transform 0.5s cubic-bezier(0.34, 1.56, 0.64, 1), opacity 0.3s ease-out' 
                : 'none'
            }}
            onClick={() => handleOrbClick(orb.id)}
            whileHover={{ 
              scale: 1.1, 
              opacity: orb.opacity * 1.5,
              boxShadow: `0 0 30px ${orb.color}50, inset 0 0 20px ${orb.color}50`
            }}
            whileTap={{ scale: 0.9 }}
          />
        );
      }
    })}
  </OrbsContainer>
)}
      
      <BackgroundOverlay 
        $reducedMotion={deviceCapabilities.prefersReducedMotion}
        initial={{ opacity: 0 }}
        animate={{ opacity: deviceCapabilities.prefersReducedMotion ? 0.9 : 0.7 }}
        transition={{ duration: 1 }}
      />
      
      {showDebug && debugInfo && (
        <DebugPanel initial={{ opacity: 0 }} animate={{ opacity: 1 }}>
          <div>FPS: {debugInfo.fps}</div>
          <div>Elements: {debugInfo.transformElements}</div>
          <div>Active: {debugInfo.activeTransforms}</div>
          <div>Scroll: {Math.round(debugInfo.scrollPosition)}px</div>
          <div>Viewport: {debugInfo.viewportWidth}x{debugInfo.viewportHeight}</div>
          <div>Gravity: {gravityFactor.toFixed(2)}</div>
          <div>Orbs: {orbs.length}</div>
        </DebugPanel>
      )}
    </BackgroundContainer>
  );
}; 