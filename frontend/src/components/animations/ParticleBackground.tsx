// src/components/animations/ParticleBackground.tsx
import React, { useEffect, useRef } from 'react';
import styled from 'styled-components';

// Container for particles that never changes size
const ParticleContainer = styled.div`
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  width: 100%;
  height: 100%;
  z-index: 1;
  pointer-events: none;
`;

// Use a singleton pattern to ensure particlesJS is only initialized once
let particlesInitialized = false;

export const ParticleBackground = React.memo(() => {
  const containerRef = useRef<HTMLDivElement>(null);
  const mountedRef = useRef<boolean>(true);

  useEffect(() => {
    // Skip if already initialized or window/particlesJS not available
    if (particlesInitialized || typeof window === 'undefined') return;
    
    const particlesJS = (window as any).particlesJS;
    if (!particlesJS) {
      console.warn('particlesJS not loaded');
      return;
    }

    // Set flag immediately to prevent double initialization
    particlesInitialized = true;
    
    // Initialize with optimal settings for performance
    particlesJS('particles-js', {
      particles: {
        number: {
          value: 60,
          density: {
            enable: true,
            value_area: 800
          }
        },
        color: {
          value: '#ffffff'
        },
        shape: {
          type: 'circle',
          stroke: {
            width: 0,
            color: '#000000'
          }
        },
        opacity: {
          value: 0.3,
          random: true,
          anim: {
            enable: true,
            speed: 0.5,
            opacity_min: 0.1,
            sync: false
          }
        },
        size: {
          value: 3,
          random: true,
          anim: {
            enable: true,
            speed: 2,
            size_min: 0.1,
            sync: false
          }
        },
        move: {
          enable: true,
          speed: 1.5, // Reduced for smoother performance
          direction: 'none',
          random: false,
          straight: false,
          out_mode: 'out',
          bounce: false
        }
      },
      interactivity: {
        detect_on: 'canvas',
        events: {
          onhover: {
            enable: true,
            mode: 'grab' // Changed from 'repulse' for better performance
          },
          onclick: {
            enable: true,
            mode: 'push'
          },
          resize: true
        }
      },
      retina_detect: true
    });
    
    // Cleanup function
    return () => {
      mountedRef.current = false;
      // particlesJS doesn't provide a destroy method, so we would need
      // to manually clean up if a destroy method becomes available
    };
  }, []); // Empty dependency array - run only once

  return <ParticleContainer ref={containerRef} id="particles-js" />;
}, () => true); // Always return true from memo comparison to prevent rerenders

ParticleBackground.displayName = 'ParticleBackground';
