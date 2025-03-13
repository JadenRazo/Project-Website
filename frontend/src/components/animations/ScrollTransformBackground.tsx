import React, { useRef, useState, useEffect, useCallback, useMemo } from 'react';
import { motion, useScroll, useTransform, useSpring, useMotionValueEvent, useAnimation, MotionValue } from 'framer-motion';
import { useTheme } from '../../contexts/ThemeContext';
import { useZIndex } from '../../contexts/ZIndexContext';
import styled from 'styled-components';
import { useInView } from 'react-intersection-observer';
import { debounce } from 'lodash';
import { useMemoryManager, withMemoryOptimization } from '../../utils/MemoryManager';
import { cleanupWebGLContext, useDebounce } from '../../utils/performance';

// Interface for explosion particles that would be used later in the component
interface ExplosionParticle {
  id: string;
  x: number;
  y: number;
  size: number;
  speedX: number;
  speedY: number;
  rotation: number;
  rotationSpeed: number;
  opacity: number;
  color: string;
  shape: 'circle' | 'square' | 'triangle';
  life: number;
  maxLife: number;
  gravity: number;
}

// Styled components for orbs with enhanced visuals and glow effects
const CircleOrb = styled(motion.div)`
  position: absolute;
  border-radius: 50%;
  transition: box-shadow 0.3s ease;
  cursor: pointer;
  filter: brightness(1.2);
  box-shadow: 0 0 15px rgba(255, 255, 255, 0.2), 
              inset 0 0 8px rgba(255, 255, 255, 0.2);
  will-change: transform, opacity, box-shadow;
  
  &:hover {
    filter: brightness(1.4);
    box-shadow: 0 0 25px rgba(255, 255, 255, 0.35), 
                inset 0 0 15px rgba(255, 255, 255, 0.35);
  }
  
  &.popped {
    opacity: 0 !important;
    transform: scale(0) !important;
    transition: transform 0.3s cubic-bezier(0.34, 1.56, 0.64, 1), opacity 0.3s ease-out;
  }
`;

const SquareOrb = styled(motion.div)`
  position: absolute;
  border-radius: 4px;
  transition: box-shadow 0.3s ease;
  cursor: pointer;
  filter: brightness(1.2);
  box-shadow: 0 0 15px rgba(255, 255, 255, 0.2), 
              inset 0 0 8px rgba(255, 255, 255, 0.2);
  will-change: transform, opacity, box-shadow;
  
  &:hover {
    filter: brightness(1.4);
    box-shadow: 0 0 25px rgba(255, 255, 255, 0.35), 
                inset 0 0 15px rgba(255, 255, 255, 0.35);
  }
  
  &.popped {
    opacity: 0 !important;
    transform: scale(0) !important;
    transition: transform 0.3s cubic-bezier(0.34, 1.56, 0.64, 1), opacity 0.3s ease-out;
  }
`;

const TriangleOrb = styled(motion.div)`
  position: absolute;
  clip-path: polygon(50% 0%, 0% 100%, 100% 100%);
  transition: filter 0.3s ease, box-shadow 0.3s ease;
  cursor: pointer;
  filter: brightness(1.2);
  box-shadow: 0 0 15px rgba(255, 255, 255, 0.2);
  will-change: transform, opacity, filter;
  
  &:hover {
    filter: brightness(1.4) drop-shadow(0 0 8px rgba(255, 255, 255, 0.35));
  }
  
  &.popped {
    opacity: 0 !important;
    transform: scale(0) !important;
    transition: transform 0.3s cubic-bezier(0.34, 1.56, 0.64, 1), opacity 0.3s ease-out;
  }
`;

// Add a styled component for explosion particles with better visual effects
const ExplosionParticleElement = styled.div<{ 
  $shape: 'circle' | 'square' | 'triangle',
  $life: number,
  $maxLife: number
}>`
  position: absolute;
  pointer-events: none;
  will-change: transform, opacity;
  z-index: 3;
  opacity: ${props => props.$life / props.$maxLife};
  filter: brightness(1.5);
  
  ${props => props.$shape === 'circle' && `
    border-radius: 50%;
  `}
  
  ${props => props.$shape === 'square' && `
    border-radius: 2px;
  `}
  
  ${props => props.$shape === 'triangle' && `
    clip-path: polygon(50% 0%, 0% 100%, 100% 100%);
  `}
`;

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
  exitBehavior?: 'reverse' | 'continue' | 'hold'; // How to animate when exiting viewport
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
  controls?: any;
  lastProgress?: number;
  isExiting?: boolean;
  exitBehavior?: string;
}

interface DebugInfo {
  fps: number;
  transformElements: number;
  activeTransforms: number;
  viewportWidth: number;
  viewportHeight: number;
  scrollPosition: number;
  scrollDirection: 'up' | 'down' | null;
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

// Add a throttle implementation at the top of the file
const throttle = <T extends (...args: any[]) => any>(
  func: T,
  limit: number
): ((...args: Parameters<T>) => void) => {
  let inThrottle = false;
  let lastFunc: ReturnType<typeof setTimeout> | null = null;
  let lastRan = 0;

  return function(this: any, ...args: Parameters<T>) {
    if (!inThrottle) {
      func.apply(this, args);
      lastRan = Date.now();
      inThrottle = true;
      
      setTimeout(() => {
        inThrottle = false;
      }, limit);
    } else {
      if (lastFunc) {
        clearTimeout(lastFunc);
      }
      
      lastFunc = setTimeout(() => {
        if (Date.now() - lastRan >= limit) {
          func.apply(this, args);
          lastRan = Date.now();
        }
      }, limit - (Date.now() - lastRan)) as unknown as ReturnType<typeof setTimeout>;
    }
  };
};

// Add a ControllerState interface to expose important state fields
interface ControllerState {
  scrollDirection: 'up' | 'down' | null;
  scrollPosition: number;
  isActive: boolean;
  gravityFactor: number;
  activeTransforms: number;
  hasUserScrolled: boolean;
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
  private shouldPreserveUserScroll: boolean = true;
  private isInitialLoad: boolean = true;
  private userHasScrolled: boolean = false;
  private disableScrollInterference: boolean = true;
  private scrollLockTimeout: number | null = null;
  private lastScrollPosition: number = 0;
  private scrollDirection: 'up' | 'down' | null = null;
  
  constructor(scrollYProgress: MotionValue<number>, deviceCapabilities: DeviceCapabilities) {
    this.scrollYProgress = scrollYProgress;
    this.deviceCapabilities = deviceCapabilities;
    this.lastTime = performance.now();
    this.lastScrollPosition = window.pageYOffset || document.documentElement.scrollTop;
    this.detectUserScrolling();
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
    
    // User has scrolled manually, don't interfere with their position
    if (this.userHasScrolled && this.disableScrollInterference) {
      // Determine scroll direction
      this.scrollDirection = scrollY > this.lastScrollPosition ? 'down' : 'up';
      
      // Apply transformations but don't modify scroll position
      this.applyTransformsWithoutScrolling();
      
      // Update last known position
      this.lastScrollPosition = scrollY;
      return;
    }
    
    // Calculate scroll direction
    this.scrollDirection = scrollY > this.lastScrollPosition ? 'down' : 'up';
    this.lastScrollPosition = scrollY;
    
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
    // Find all elements for this section's transforms
    for (const transform of section.transforms) {
      // Get all elements matching the selector
      const elements = this.container!.querySelectorAll(transform.target);
      
      if (elements.length === 0) continue;
      
      // Track if elements for this transform were previously active
      const transformElementId = `${section.id}-${transform.target}`;
      const wasActive = this.transformElements.has(transformElementId);
      
      // Process each matching element
      elements.forEach((element, index) => {
        // Unique ID for this specific element
        const elementId = `${transformElementId}-${index}`;
        
        // Assess if element is entering or exiting the view
        const isEntering = this.scrollDirection === 'down' && progress < 0.5;
        const isExiting = this.scrollDirection === 'up' && progress > 0.5;
        
        // If exiting, adjust progress based on exit behavior
        const exitBehavior = transform.exitBehavior || 'reverse';
        let effectiveProgress = progress;
        
        if (isExiting) {
          switch (exitBehavior) {
            case 'reverse':
              // Reverse the animation (0 becomes 1, 1 becomes 0)
              effectiveProgress = 1 - progress;
              break;
            case 'continue':
              // Let the animation continue past its range
              effectiveProgress = progress;
              break;
            case 'hold':
              // Freeze at the final state
              effectiveProgress = 1;
              break;
          }
        }
        
        // Apply the transform
        this.applySpecificTransform(element as HTMLElement, transform, effectiveProgress);
        
        // Add to active transforms map if not already there
        if (!this.transformElements.has(elementId)) {
          this.transformElements.set(elementId, {
            id: elementId,
            element: element as HTMLElement,
            initialState: this.captureInitialState(element as HTMLElement, transform),
            effects: [transform],
            lastProgress: effectiveProgress,
            isExiting: isExiting,
            exitBehavior
          });
        } else {
          // Update existing transform element
          const existing = this.transformElements.get(elementId)!;
          existing.lastProgress = effectiveProgress;
          existing.isExiting = isExiting;
        }
      });
    }
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
    
    // Clear any pending timeouts
    if (this.scrollLockTimeout) {
      window.clearTimeout(this.scrollLockTimeout);
    }
    
    // Remove scroll detection event listeners
    window.removeEventListener('wheel', this.detectUserScrolling);
    window.removeEventListener('touchmove', this.detectUserScrolling);
  }
  
  public getDebugInfo(): DebugInfo {
    return {
      fps: this.fps,
      transformElements: this.transformElements.size,
      activeTransforms: this.activeTransforms,
      viewportWidth: window.innerWidth,
      viewportHeight: window.innerHeight,
      scrollPosition: this.lastScrollPosition,
      scrollDirection: this.scrollDirection,
      isActive: this.isActive,
      gravityFactor: this.gravityFactor
    };
  }

  private detectUserScrolling(): void {
    if (typeof window === 'undefined') return;
    
    const handleUserScroll = () => {
      this.userHasScrolled = true;
      // If user has scrolled, we should preserve their position
      this.shouldPreserveUserScroll = true;
      
      // Set a flag to ignore programmatic scrolling temporarily
      if (this.isInitialLoad) {
        this.isInitialLoad = false;
      }
    };
    
    window.addEventListener('wheel', handleUserScroll, { passive: true });
    window.addEventListener('touchmove', handleUserScroll, { passive: true });
    
    // Also detect keyboard navigation
    window.addEventListener('keydown', (e) => {
      if (e.key === 'ArrowUp' || e.key === 'ArrowDown' || 
          e.key === 'PageUp' || e.key === 'PageDown' || 
          e.key === 'Home' || e.key === 'End' || e.key === ' ') {
        handleUserScroll();
      }
    });
  }

  private applyTransformsWithoutScrolling(): void {
    if (!this.container) return;
    
    const scrollY = window.pageYOffset || document.documentElement.scrollTop;
    const viewportHeight = window.innerHeight;
    const documentHeight = Math.max(
      document.body.scrollHeight, 
      document.documentElement.scrollHeight,
      document.body.offsetHeight,
      document.documentElement.offsetHeight
    );
    
    // Calculate overall scroll progress (0-1)
    const scrollProgress = Math.min(1, Math.max(0, scrollY / (documentHeight - viewportHeight)));
    
    let activeTransforms = 0;
    
    // Track which sections are active and which are exiting
    const activeSectionIds = new Set<string>();
    const exitingSectionIds = new Set<string>();
    
    // First pass: determine active and exiting sections
    for (const section of this.sections) {
      const sectionStartY = section.startPercent * (documentHeight - viewportHeight);
      const sectionEndY = section.endPercent * (documentHeight - viewportHeight);
      
      if (scrollY >= sectionStartY && scrollY <= sectionEndY) {
        // Section is active
        activeSectionIds.add(section.id);
        
        // Calculate progress within this section (0-1)
        const sectionProgress = 
          (scrollY - sectionStartY) / (sectionEndY - sectionStartY);
        
        this.applyTransforms(section, sectionProgress);
        activeTransforms += section.transforms.length;
      } else if (scrollY < sectionStartY && this.scrollDirection === 'up') {
        // We're scrolling up and about to enter this section
        exitingSectionIds.add(section.id);
      } else if (scrollY > sectionEndY && this.scrollDirection === 'down') {
        // We're scrolling down and just left this section
        exitingSectionIds.add(section.id);
      }
    }
    
    // Second pass: apply exit animations to sections that are no longer active
    for (const section of this.sections) {
      if (exitingSectionIds.has(section.id)) {
        // Calculate approximate exit progress
        const sectionStartY = section.startPercent * (documentHeight - viewportHeight);
        const sectionEndY = section.endPercent * (documentHeight - viewportHeight);
        const sectionMiddleY = (sectionStartY + sectionEndY) / 2;
        
        // Calculate how far we are from the section's edge
        const distanceFromSection = this.scrollDirection === 'up' 
          ? sectionStartY - scrollY 
          : scrollY - sectionEndY;
        
        // Normalize to 0-1 range for exit progress
        const exitDistance = Math.min(viewportHeight / 2, distanceFromSection);
        const exitProgress = Math.min(1, exitDistance / (viewportHeight / 2));
        
        // Apply transforms with exit behavior
        for (const transform of section.transforms) {
          const elements = this.container!.querySelectorAll(transform.target);
          
          elements.forEach((element) => {
            const exitBehavior = transform.exitBehavior || 'reverse';
            let effectiveProgress: number;
            
            switch (exitBehavior) {
              case 'reverse':
                // Reverse the animation
                effectiveProgress = 1 - exitProgress;
                break;
              case 'continue':
                // Continue in the same direction
                effectiveProgress = 1 + exitProgress;
                break;
              case 'hold':
              default:
                // Keep at final state
                effectiveProgress = 1;
                break;
            }
            
            this.applySpecificTransform(element as HTMLElement, transform, effectiveProgress);
            activeTransforms++;
          });
        }
      }
    }
    
    this.activeTransforms = activeTransforms;
  }
  
  public setScrollInterference(allow: boolean): void {
    this.disableScrollInterference = !allow;
  }
  
  public unlockScrollPreservation(delay: number = 1000): void {
    if (this.scrollLockTimeout) {
      window.clearTimeout(this.scrollLockTimeout);
    }
    
    this.scrollLockTimeout = window.setTimeout(() => {
      this.shouldPreserveUserScroll = false;
      this.scrollLockTimeout = null;
    }, delay) as unknown as number;
  }

  private captureInitialState(element: HTMLElement, transform: ScrollTransformEffect): any {
    const state: any = {};
    
    switch (transform.type) {
      case 'opacity':
        state.opacity = window.getComputedStyle(element).opacity;
        break;
      case 'translateX':
      case 'translateY':
        const transform = window.getComputedStyle(element).transform;
        state.transform = transform;
        break;
      case 'scale':
        const scale = window.getComputedStyle(element).transform;
        state.transform = scale;
        break;
      case 'rotate':
        const rotate = window.getComputedStyle(element).transform;
        state.transform = rotate;
        break;
      case 'blur':
        const filter = window.getComputedStyle(element).filter;
        state.filter = filter;
        break;
      case 'color':
        const color = window.getComputedStyle(element).color;
        state.color = color;
        break;
    }
    
    return state;
  }

  // Add a public getter for scroll direction
  public getScrollDirection(): 'up' | 'down' | null {
    return this.scrollDirection;
  }

  /**
   * Returns a snapshot of the current controller state
   * Provides key state information for components to use in animations and rendering
   */
  public getControllerState(): ControllerState {
    return {
      scrollDirection: this.scrollDirection,
      scrollPosition: this.lastScrollPosition,
      isActive: this.isActive,
      gravityFactor: this.gravityFactor,
      activeTransforms: this.activeTransforms,
      hasUserScrolled: this.userHasScrolled
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
        easing: 'easeOut',
        exitBehavior: 'reverse'
      },
      {
        type: 'translateY',
        target: 'header',
        from: '-20',
        to: '0',
        unit: 'px',
        easing: 'easeOut',
        exitBehavior: 'reverse'
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
        easing: 'easeInOut',
        exitBehavior: 'reverse'
      },
      {
        type: 'opacity',
        target: '#hero .hero-subtitle',
        from: '1',
        to: '0.6',
        easing: 'easeInOut',
        exitBehavior: 'reverse'
      },
      {
        type: 'scale',
        target: '#hero .hero-image',
        from: '1',
        to: '0.95',
        easing: 'easeInOut',
        exitBehavior: 'reverse'
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
        easing: 'easeOut',
        exitBehavior: 'reverse'
      },
      {
        type: 'translateY',
        target: '#skills .section-title',
        from: '30',
        to: '0',
        unit: 'px',
        easing: 'easeOut',
        exitBehavior: 'reverse'
      },
      {
        type: 'opacity',
        target: '#skills .skill-item:nth-child(1)',
        from: '0',
        to: '1',
        easing: 'easeOut',
        exitBehavior: 'reverse'
      },
      {
        type: 'translateX',
        target: '#skills .skill-item:nth-child(1)',
        from: '-50',
        to: '0',
        unit: 'px',
        easing: 'easeOut',
        exitBehavior: 'reverse'
      },
      {
        type: 'opacity',
        target: '#skills .skill-item:nth-child(2)',
        from: '0',
        to: '1',
        easing: 'easeOut',
        exitBehavior: 'reverse'
      },
      {
        type: 'translateY',
        target: '#skills .skill-item:nth-child(2)',
        from: '50',
        to: '0',
        unit: 'px',
        easing: 'easeOut',
        exitBehavior: 'reverse'
      },
      {
        type: 'opacity',
        target: '#skills .skill-item:nth-child(3)',
        from: '0',
        to: '1',
        easing: 'easeOut',
        exitBehavior: 'reverse'
      },
      {
        type: 'translateX',
        target: '#skills .skill-item:nth-child(3)',
        from: '50',
        to: '0',
        unit: 'px',
        easing: 'easeOut',
        exitBehavior: 'reverse'
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
        easing: 'easeOut',
        exitBehavior: 'reverse'
      },
      {
        type: 'translateY',
        target: '#projects .section-title',
        from: '30',
        to: '0',
        unit: 'px',
        easing: 'easeOut',
        exitBehavior: 'reverse'
      },
      {
        type: 'opacity',
        target: '#projects .project-card',
        from: '0',
        to: '1',
        easing: 'easeOut',
        exitBehavior: 'reverse'
      },
      {
        type: 'translateY',
        target: '#projects .project-card',
        from: '50',
        to: '0',
        unit: 'px',
        easing: 'easeOut',
        exitBehavior: 'reverse'
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
        easing: 'easeOut',
        exitBehavior: 'reverse'
      },
      {
        type: 'translateY',
        target: '#about .section-title',
        from: '30',
        to: '0',
        unit: 'px',
        easing: 'easeOut',
        exitBehavior: 'reverse'
      },
      {
        type: 'opacity',
        target: '#about .about-content',
        from: '0',
        to: '1',
        easing: 'easeOut',
        exitBehavior: 'reverse'
      },
      {
        type: 'translateY',
        target: '#about .about-content',
        from: '50',
        to: '0',
        unit: 'px',
        easing: 'easeOut',
        exitBehavior: 'reverse'
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
  pointer-events: none;
  overflow: hidden;
  z-index: 1;
  
  /* Add styles for entering and exiting orbs */
  .orb {
    transition-property: transform, opacity;
    will-change: transform, opacity;
    pointer-events: auto;
    
    &.entering {
      animation: fadeIn 0.6s ease-out forwards;
    }
    
    &.exiting {
      animation: fadeOut 0.6s ease-in forwards;
    }
    
    &.popped {
      animation: popEffect 0.5s cubic-bezier(0.34, 1.56, 0.64, 1) forwards;
    }
  }
  
  @keyframes fadeIn {
    from { opacity: 0; transform: translateY(20px) scale(0.8); }
    to { opacity: 1; transform: translateY(0) scale(1); }
  }
  
  @keyframes fadeOut {
    from { opacity: 1; transform: translateY(0) scale(1); }
    to { opacity: 0; transform: translateY(-20px) scale(0.8); }
  }
  
  @keyframes popEffect {
    0% { transform: scale(1) rotate(0deg); }
    50% { transform: scale(1.4) rotate(10deg); opacity: 0.7; }
    100% { transform: scale(0) rotate(20deg); opacity: 0; }
  }
  
  /* Specific orb type treatments */
  [data-type="circle"] {
    border-radius: 50%;
  }
  
  [data-type="square"] {
    border-radius: 4px;
    
    &.exiting {
      animation: squareExit 0.8s ease-in forwards;
    }
  }
  
  [data-type="triangle"] {
    clip-path: polygon(50% 0%, 0% 100%, 100% 100%);
    
    &.exiting {
      animation: triangleExit 0.8s ease-in forwards;
    }
  }
  
  @keyframes squareExit {
    0% { transform: rotate(0deg) scale(1); }
    100% { transform: rotate(90deg) scale(0.6); opacity: 0; }
  }
  
  @keyframes triangleExit {
    0% { transform: rotate(0deg) scale(1); }
    100% { transform: rotate(-60deg) scale(0.6); opacity: 0; }
  }
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
  preserveUserScroll?: boolean;
}> = ({ 
  customSections, 
  showDebug = false, 
  enableFloatingOrbs = true,
  preserveUserScroll = true
}) => {
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
  
  // Add memory optimization
  const { isMemoryConstrained, applicationState } = useMemoryManager();
  const actualEnableFloatingOrbs = enableFloatingOrbs && 
    applicationState.backgroundEffectsEnabled &&
    !isMemoryConstrained;
  
  // Add state to track user-initiated scrolling
  const [userScrolled, setUserScrolled] = useState(false);
  
  // Add reference to store scroll restoration position
  const scrollPositionRef = useRef<number>(0);
  
  // Add state for controller state snapshot
  const [controllerState, setControllerState] = useState<ControllerState>({
    scrollDirection: null,
    scrollPosition: 0,
    isActive: true,
    gravityFactor: 1,
    activeTransforms: 0,
    hasUserScrolled: false
  });
  
  // Add state for explosion particles
  const [explosionParticles, setExplosionParticles] = useState<ExplosionParticle[]>([]);
  const explosionParticlesRef = useRef<ExplosionParticle[]>([]);
  explosionParticlesRef.current = explosionParticles;
  
  // Add animation frame reference for explosion animation
  const explosionAnimationRef = useRef<number | null>(null);
  
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
    if (!actualEnableFloatingOrbs) return;
    
    const generateOrbs = () => {
      if (deviceCapabilities.isMobile && !deviceCapabilities.isHighPerformance) {
        setOrbs([]);
        return;
      }
      
      const orbCount = isMemoryConstrained 
        ? 5 
        : applicationState.effectsEnabled 
          ? 20 
          : 10;
      
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
  }, [deviceCapabilities, theme.colors, actualEnableFloatingOrbs]);
  
  // Optimized orb animation loop
  useEffect(() => {
    if (!actualEnableFloatingOrbs || orbs.length === 0) return;
    
    const animateOrbs = (timestamp: number) => {
      // Skip animation frames when memory is constrained
      if (isMemoryConstrained && timestamp % 2 !== 0) {
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
  }, [orbs, deviceCapabilities, gravityFactor, actualEnableFloatingOrbs]);
  
  // Handle scroll events for transformations and gravity
  const handleScrollEvent = useCallback(() => {
    if (controllerRef.current) {
      const latest = window.pageYOffset || document.documentElement.scrollTop;
      const viewportHeight = window.innerHeight;
      
      // Calculate scroll direction and speed
      const direction = latest > prevScrollY.current ? 1 : -1;
      const scrollDelta = Math.abs(latest - prevScrollY.current);
      const speed = Math.min(scrollDelta / viewportHeight, 1);
      
      // Track scroll direction for orb animations
      const scrollDirection = direction > 0 ? 'down' : 'up';
      
      // Calculate gravity factor with better bounds
      let newGravityFactor = 1 + (direction * speed * 2);
      newGravityFactor = Math.max(0.2, Math.min(newGravityFactor, 2.5));
      
      // Update scroll progress with controlled gravity
      controllerRef.current.updateOnScroll(latest, viewportHeight, newGravityFactor);
      
      // Update controller state snapshot
      setControllerState(controllerRef.current.getControllerState());
      
      // Update debug information if enabled
      if (showDebug) {
        setDebugInfo(controllerRef.current.getDebugInfo());
      }
      
      // Apply gravity changes to floating orbs if enabled
      if (actualEnableFloatingOrbs) {
        setGravityFactor(newGravityFactor);
        
        // Update orbs based on scroll direction
        updateOrbsOnScroll(scrollDirection, speed);
      }
      
      // Update previous value
      prevScrollY.current = latest;
    }
  }, [controllerRef, showDebug, actualEnableFloatingOrbs]);
  
  // Add a function to update orbs based on scroll direction
  const updateOrbsOnScroll = useCallback((direction: 'up' | 'down', speed: number) => {
    if (!actualEnableFloatingOrbs || orbs.length === 0) return;
    
    setOrbs(currentOrbs => {
      return currentOrbs.map(orb => {
        // Skip popped orbs
        if (orb.popped) return orb;
        
        // When scrolling down, increase orb opacity briefly for entering effect
        if (direction === 'down') {
          return {
            ...orb,
            opacity: Math.min(1, orb.opacity + speed * 0.1)
          };
        }
        
        // When scrolling up, add subtle upward acceleration and fade slightly
        if (direction === 'up') {
          return {
            ...orb,
            y: Math.max(0, orb.y - speed * 2),
            opacity: Math.max(0.3, orb.opacity - speed * 0.05)
          };
        }
        
        return orb;
      });
    });
  }, [actualEnableFloatingOrbs, orbs.length]);
  
  // Handle window resize to update transformations
  const handleResizeEvent = useCallback(() => {
    if (controllerRef.current) {
      controllerRef.current.resize();
    }
  }, [controllerRef]);
  
  // Apply throttling to the scroll handler for better performance
  const throttledScrollHandler = useMemo(() => 
    throttle(handleScrollEvent, 16), // ~60fps for smooth scrolling
  [handleScrollEvent]);
  
  // Set up scroll and resize event listeners
  useEffect(() => {
    // Add event listeners for scroll and resize
    window.addEventListener('scroll', throttledScrollHandler, { passive: true });
    window.addEventListener('resize', handleResizeEvent);
    
    // Initial setup
    handleScrollEvent();
    
    // Clean up event listeners on unmount
    return () => {
      window.removeEventListener('scroll', throttledScrollHandler);
      window.removeEventListener('resize', handleResizeEvent);
    };
  }, [throttledScrollHandler, handleResizeEvent]);
  
  // Enhanced orb click handler with explosion effect
  const handleOrbClick = useCallback((id: string, x: number, y: number, size: number, color: string, shape: 'circle' | 'square' | 'triangle') => {
    if (!actualEnableFloatingOrbs) return;
    
    // Mark the orb as popped with a scale animation
    setOrbs(currentOrbs => 
      currentOrbs.map(orb => 
        orb.id === id ? { ...orb, popped: true } : orb
      )
    );
    
    // Play a pop sound effect (if supported)
    if (typeof window !== 'undefined' && window.navigator.vibrate && deviceCapabilities.isTouch) {
      window.navigator.vibrate(50); // Haptic feedback on supported devices
    }
    
    // Create explosion particles with more dynamic behavior
    const particleCount = deviceCapabilities.isHighPerformance ? 30 : 15;
    const newParticles: ExplosionParticle[] = [];
    
    // Generate unique ID for this explosion
    const explosionId = `explosion-${Date.now()}`;
    
    // Create particles
    for (let i = 0; i < particleCount; i++) {
      // Particle physics properties
      const particleSize = size * (0.1 + Math.random() * 0.3);
      const particleAngle = Math.random() * Math.PI * 2; // Random direction
      const particleSpeed = 2 + Math.random() * 5; // Varied speed
      
      // Calculate velocity components
      const particleSpeedX = Math.cos(particleAngle) * particleSpeed;
      const particleSpeedY = Math.sin(particleAngle) * particleSpeed - (Math.random() * 2); // Add upward boost
      
      // Color variation
      const hueShift = Math.random() * 30 - 15;
      const particleColor = shiftHue(color, hueShift);
      
      // Create the particle
      newParticles.push({
        id: `${explosionId}-particle-${i}`,
        x, // Initial position at orb center
        y,
        size: particleSize,
        speedX: particleSpeedX,
        speedY: particleSpeedY,
        rotation: Math.random() * 360,
        rotationSpeed: (Math.random() - 0.5) * 15,
        opacity: 0.7 + Math.random() * 0.3,
        color: particleColor,
        shape,
        life: 80 + Math.random() * 40,
        maxLife: 120,
        gravity: 0.08 + Math.random() * 0.06
      });
    }
    
    // Add particles to state
    setExplosionParticles(current => [...current, ...newParticles]);
    
    // Start animation if not already running
    if (!explosionAnimationRef.current) {
      animateExplosionParticles();
    }
    
    // Create new orb after a delay to replace the popped one
    setTimeout(() => {
      setOrbs(currentOrbs => {
        const poppedIndex = currentOrbs.findIndex(orb => orb.id === id);
        if (poppedIndex === -1) return currentOrbs; // Orb not found
        
        const primaryColor = theme.colors?.primary || '#6c63ff';
        const secondaryColor = theme.colors?.secondary || '#ff6b6b';
        const tertiaryColor = theme.colors?.accent || '#4ecdc4';
        const colorPalette = [primaryColor, secondaryColor, tertiaryColor];
        const shapeTypes: Array<'circle' | 'square' | 'triangle'> = ['circle', 'square', 'triangle'];
        
        // Create a new orb with fresh properties to replace the popped one
        const updatedOrbs = [...currentOrbs];
        updatedOrbs[poppedIndex] = {
          ...updatedOrbs[poppedIndex],
          popped: false,
          y: 80 + Math.random() * 15, // Start near bottom of screen
          x: Math.random() * 100, // Random horizontal position
          speed: 0.2 + Math.random() * 0.3,
          opacity: 0.2 + Math.random() * 0.5,
          size: 20 + Math.random() * 60,
          color: colorPalette[Math.floor(Math.random() * colorPalette.length)],
          rotationSpeed: (Math.random() - 0.5) * 0.5,
          currentRotation: Math.random() * 360,
          type: shapeTypes[Math.floor(Math.random() * shapeTypes.length)]
        };
        
        return updatedOrbs;
      });
    }, 2000); // Wait 2 seconds before replacing
  }, [theme.colors, deviceCapabilities, actualEnableFloatingOrbs]);
  
  // Function to animate explosion particles with enhanced physics
  const animateExplosionParticles = useCallback(() => {
    if (explosionParticlesRef.current.length === 0) {
      explosionAnimationRef.current = null;
      return;
    }
    
    setExplosionParticles(particles => {
      return particles
        .map(particle => {
          // Apply physics: gravity, momentum, and drag
          const newSpeedY = particle.speedY + particle.gravity;
          
          // Add slight drag/air resistance
          const dragFactor = 0.98;
          const newSpeedX = particle.speedX * dragFactor;
          
          // Calculate life reduction based on size (smaller particles fade faster)
          const lifeReduction = 1 + (1 - particle.size / 30) * 0.5;
          
          return {
            ...particle,
            x: particle.x + newSpeedX,
            y: particle.y + newSpeedY,
            speedX: newSpeedX,
            speedY: newSpeedY,
            rotation: particle.rotation + particle.rotationSpeed,
            // Gradually slow down rotation
            rotationSpeed: particle.rotationSpeed * 0.98,
            // Reduce life more rapidly for smaller particles
            life: particle.life - lifeReduction,
            // Fade out as life decreases
            opacity: (particle.life / particle.maxLife) * particle.opacity
          };
        })
        .filter(particle => {
          // Remove particles that are out of bounds or expired
          return (
            particle.life > 0 &&
            particle.x > -100 &&
            particle.x < window.innerWidth + 100 &&
            particle.y > -100 &&
            particle.y < window.innerHeight + 100
          );
        });
    });
    
    // Continue animation loop
    explosionAnimationRef.current = requestAnimationFrame(animateExplosionParticles);
  }, []);
  
  // Cleanup animation frame
  useEffect(() => {
    return () => {
      if (explosionAnimationRef.current) {
        cancelAnimationFrame(explosionAnimationRef.current);
      }
    };
  }, []);
  
  // Helper function to shift hue of a color
  const shiftHue = (hexColor: string, shift: number): string => {
    // Convert hex to RGB
    const hex = hexColor.replace('#', '');
    const r = parseInt(hex.substring(0, 2), 16);
    const g = parseInt(hex.substring(2, 4), 16);
    const b = parseInt(hex.substring(4, 6), 16);
    
    // Convert RGB to HSL
    const [h, s, l] = rgbToHsl(r, g, b);
    
    // Shift hue and convert back to RGB
    const newHue = (h + shift) % 360;
    const [r2, g2, b2] = hslToRgb(newHue, s, l);
    
    // Convert back to hex
    return `#${Math.round(r2).toString(16).padStart(2, '0')}${Math.round(g2).toString(16).padStart(2, '0')}${Math.round(b2).toString(16).padStart(2, '0')}`;
  };
  
  // RGB to HSL conversion
  const rgbToHsl = (r: number, g: number, b: number): [number, number, number] => {
    r /= 255;
    g /= 255;
    b /= 255;
    
    const max = Math.max(r, g, b);
    const min = Math.min(r, g, b);
    let h = 0, s, l = (max + min) / 2;
    
    if (max === min) {
      h = s = 0; // achromatic
    } else {
      const d = max - min;
      s = l > 0.5 ? d / (2 - max - min) : d / (max + min);
      switch (max) {
        case r: h = (g - b) / d + (g < b ? 6 : 0); break;
        case g: h = (b - r) / d + 2; break;
        case b: h = (r - g) / d + 4; break;
      }
      h /= 6;
    }
    
    return [h * 360, s, l];
  };
  
  // HSL to RGB conversion
  const hslToRgb = (h: number, s: number, l: number): [number, number, number] => {
    h /= 360;
    let r, g, b;
    
    if (s === 0) {
      r = g = b = l; // achromatic
    } else {
      const hue2rgb = (p: number, q: number, t: number) => {
        if (t < 0) t += 1;
        if (t > 1) t -= 1;
        if (t < 1/6) return p + (q - p) * 6 * t;
        if (t < 1/2) return q;
        if (t < 2/3) return p + (q - p) * (2/3 - t) * 6;
        return p;
      };
      
      const q = l < 0.5 ? l * (1 + s) : l + s - l * s;
      const p = 2 * l - q;
      r = hue2rgb(p, q, h + 1/3);
      g = hue2rgb(p, q, h);
      b = hue2rgb(p, q, h - 1/3);
    }
    
    return [r * 255, g * 255, b * 255];
  };
  
  // Render the background with interactive orbs
  return (
    <BackgroundContainer 
      ref={containerRef}
      style={{ 
        zIndex: zIndex.PARTICLE_BACKGROUND,
        filter: isMemoryConstrained ? 'blur(1px)' : 'none',
        opacity: isMemoryConstrained ? 0.8 : 1
      }}
      className="scroll-transform-background"
    >
      {actualEnableFloatingOrbs && (
        <OrbsContainer>
          {orbs.map(orb => {
            // Get pixel position for explosion calculation
            const orbPositionX = (orb.x / 100) * (window.innerWidth);
            const orbPositionY = (orb.y / 100) * (window.innerHeight);
            
            // Calculate additional transform properties based on scroll direction
            const scrolling = prevScrollY.current > 0;
            const scrollingUp = controllerState.scrollDirection === 'up';
            const orbExiting = scrollingUp && orb.y < 30; // Orb near top of screen when scrolling up
            
            // Scale down orbs that are exiting when scrolling up
            const exitScale = orbExiting ? Math.max(0.7, 1 - (30 - orb.y) / 30) : 1;
            
            // Adjust orb rotation based on scroll direction
            const rotationAdjustment = scrollingUp ? orb.rotationSpeed * 2 : 0;
            
            // Apply subtle transitions based on orb's position
            const isLowOrb = orb.y > 70;
            const isHighOrb = orb.y < 30;
            const positionClass = isLowOrb ? 'entering' : isHighOrb ? 'exiting' : '';
            
            // Enhanced color for brighter appearance
            const enhancedColor = orb.color;
            
            if (orb.type === 'circle') {
              return (
                <CircleOrb
                  key={orb.id}
                  style={{
                    left: `${orb.x}%`,
                    top: `${orb.y}%`,
                    width: `${orb.size}px`,
                    height: `${orb.size}px`,
                    opacity: isMemoryConstrained ? orb.opacity * 0.7 : orb.opacity,
                    backgroundColor: enhancedColor + (isMemoryConstrained ? '30' : '15'),
                    transform: `
                      rotate(${orb.currentRotation + rotationAdjustment}deg)
                      scale(${exitScale})
                      ${isMemoryConstrained ? '' : 'translateZ(0)'}
                    `,
                    transition: orbExiting ? 'opacity 0.3s ease-out, transform 0.3s ease-out' : 'none'
                  }}
                  onClick={() => !isMemoryConstrained && handleOrbClick(
                    orb.id, 
                    orbPositionX, 
                    orbPositionY, 
                    orb.size, 
                    enhancedColor,
                    'circle'
                  )}
                  whileHover={{ 
                    scale: orbExiting ? exitScale : 1.2, 
                    opacity: 1,
                    boxShadow: `0 0 25px ${enhancedColor}50, inset 0 0 15px ${enhancedColor}50` 
                  }}
                  className={`orb ${orb.popped ? 'popped' : ''} ${positionClass}`}
                  data-type="circle"
                />
              );
            }
            
            // Similar updates for square and triangle orbs...
            if (orb.type === 'square') {
              return (
                <SquareOrb
                  key={orb.id}
                  style={{
                    left: `${orb.x}%`,
                    top: `${orb.y}%`,
                    width: `${orb.size * 0.8}px`,
                    height: `${orb.size * 0.8}px`,
                    opacity: isMemoryConstrained ? orb.opacity * 0.7 : orb.opacity,
                    backgroundColor: enhancedColor + (isMemoryConstrained ? '30' : '15'),
                    transform: `
                      rotate(${orb.currentRotation + rotationAdjustment}deg)
                      scale(${exitScale})
                      ${isMemoryConstrained ? '' : 'translateZ(0)'}
                    `,
                    transition: orbExiting ? 'opacity 0.3s ease-out, transform 0.3s ease-out' : 'none'
                  }}
                  onClick={() => !isMemoryConstrained && handleOrbClick(
                    orb.id, 
                    orbPositionX, 
                    orbPositionY, 
                    orb.size * 0.8, 
                    enhancedColor,
                    'square'
                  )}
                  whileHover={{ 
                    scale: orbExiting ? exitScale : 1.2, 
                    opacity: 1,
                    boxShadow: `0 0 25px ${enhancedColor}50, inset 0 0 15px ${enhancedColor}50` 
                  }}
                  className={`orb ${orb.popped ? 'popped' : ''} ${positionClass}`}
                  data-type="square"
                />
              );
            }
            
            if (orb.type === 'triangle') {
              return (
                <TriangleOrb
                  key={orb.id}
                  style={{
                    left: `${orb.x}%`,
                    top: `${orb.y}%`,
                    width: `${orb.size}px`,
                    height: `${orb.size}px`,
                    opacity: isMemoryConstrained ? orb.opacity * 0.7 : orb.opacity,
                    backgroundColor: enhancedColor + (isMemoryConstrained ? '30' : '15'),
                    transform: `
                      rotate(${orb.currentRotation + rotationAdjustment}deg)
                      scale(${exitScale})
                      ${isMemoryConstrained ? '' : 'translateZ(0)'}
                    `,
                    transition: orbExiting ? 'opacity 0.3s ease-out, transform 0.3s ease-out' : 'none'
                  }}
                  onClick={() => !isMemoryConstrained && handleOrbClick(
                    orb.id, 
                    orbPositionX, 
                    orbPositionY, 
                    orb.size, 
                    enhancedColor,
                    'triangle'
                  )}
                  whileHover={{ 
                    scale: orbExiting ? exitScale : 1.2, 
                    opacity: 1,
                    filter: `brightness(1.4) drop-shadow(0 0 8px ${enhancedColor}50)` 
                  }}
                  className={`orb ${orb.popped ? 'popped' : ''} ${positionClass}`}
                  data-type="triangle"
                />
              );
            }
          })}
        </OrbsContainer>
      )}
      
      {/* Explosion particles layer */}
      {explosionParticles.map(particle => (
        <ExplosionParticleElement
          key={particle.id}
          $shape={particle.shape}
          $life={particle.life}
          $maxLife={particle.maxLife}
          style={{
            left: `${particle.x}px`,
            top: `${particle.y}px`,
            width: `${particle.size}px`,
            height: `${particle.size}px`,
            backgroundColor: particle.color,
            opacity: particle.opacity * (particle.life / particle.maxLife),
            transform: `rotate(${particle.rotation}deg) scale(${0.5 + (particle.life / particle.maxLife) * 0.5})`,
            boxShadow: `0 0 ${particle.size * 0.8}px ${particle.color}99`,
            transition: 'box-shadow 0.1s ease'
          }}
        />
      ))}
      
      <BackgroundOverlay 
        $reducedMotion={deviceCapabilities.prefersReducedMotion}
        initial={{ opacity: 0 }}
        animate={{ opacity: deviceCapabilities.prefersReducedMotion ? 0.9 : 0.7 }}
        transition={{ duration: 1 }}
      />
      
      {showDebug && debugInfo && (
        <DebugPanel initial={{ opacity: 0 }} animate={{ opacity: 1 }}>
          <div>FPS: {Math.round(debugInfo.fps)}</div>
          <div>Elements: {debugInfo.transformElements}</div>
          <div>Active: {debugInfo.activeTransforms}</div>
          <div>Scroll: {Math.round(debugInfo.scrollPosition)}px</div>
          <div>Direction: {debugInfo.scrollDirection || 'none'}</div>
          <div>Viewport: {debugInfo.viewportWidth}x{debugInfo.viewportHeight}</div>
          <div>Gravity: {gravityFactor.toFixed(2)}</div>
          <div>Orbs: {orbs.length}</div>
          <div>Particles: {explosionParticles.length}</div>
        </DebugPanel>
      )}
    </BackgroundContainer>
  );
}; 