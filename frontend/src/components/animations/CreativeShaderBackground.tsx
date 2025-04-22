import React, { useRef, useEffect, useState, useMemo } from 'react';
import * as THREE from 'three';
import styled from 'styled-components';
import { useTheme } from '../../contexts/ThemeContext';
import usePerformanceOptimizations from '../../hooks/usePerformanceOptimizations';

const AnimationContainer = styled.div<{ visible: boolean }>`
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 0;
  pointer-events: none;
  overflow: hidden;
  transition: opacity 0.8s ease;
  opacity: ${({ visible }) => (visible ? 1 : 0)};
  will-change: transform;
  transform: translateZ(0);
  max-width: 100%;
  box-sizing: border-box;
  margin: 0;
  padding: 0;
`;

interface CreativeShaderBackgroundProps {
  disableParallax?: boolean;
  intensity?: number;
  speed?: number;
  colorIntensity?: number;
  pattern?: 'waves' | 'flow' | 'particles' | 'ribbons';
}

const fragmentShader = `
  uniform float time;
  uniform vec2 resolution;
  uniform vec3 baseColor;
  uniform vec3 accentColor;
  uniform float intensity;
  uniform float speed;
  uniform float colorIntensity;
  uniform int pattern;
  uniform vec2 mousePosition;
  
  // Constants
  const float PI = 3.14159265359;
  
  // Simplex noise functions
  vec3 mod289(vec3 x) { return x - floor(x * (1.0 / 289.0)) * 289.0; }
  vec2 mod289(vec2 x) { return x - floor(x * (1.0 / 289.0)) * 289.0; }
  vec3 permute(vec3 x) { return mod289(((x*34.0)+1.0)*x); }
  
  float snoise(vec2 v) {
    const vec4 C = vec4(0.211324865405187,  // (3.0-sqrt(3.0))/6.0
                        0.366025403784439,  // 0.5*(sqrt(3.0)-1.0)
                        -0.577350269189626,  // -1.0 + 2.0 * C.x
                        0.024390243902439); // 1.0 / 41.0
    vec2 i  = floor(v + dot(v, C.yy) );
    vec2 x0 = v -   i + dot(i, C.xx);
    vec2 i1;
    i1 = (x0.x > x0.y) ? vec2(1.0, 0.0) : vec2(0.0, 1.0);
    vec4 x12 = x0.xyxy + C.xxzz;
    x12.xy -= i1;
    i = mod289(i);
    vec3 p = permute( permute( i.y + vec3(0.0, i1.y, 1.0 ))
      + i.x + vec3(0.0, i1.x, 1.0 ));
    vec3 m = max(0.5 - vec3(dot(x0,x0), dot(x12.xy,x12.xy), dot(x12.zw,x12.zw)), 0.0);
    m = m*m ;
    m = m*m ;
    vec3 x = 2.0 * fract(p * C.www) - 1.0;
    vec3 h = abs(x) - 0.5;
    vec3 ox = floor(x + 0.5);
    vec3 a0 = x - ox;
    m *= 1.79284291400159 - 0.85373472095314 * ( a0*a0 + h*h );
    vec3 g;
    g.x  = a0.x  * x0.x  + h.x  * x0.y;
    g.yz = a0.yz * x12.xz + h.yz * x12.yw;
    return 130.0 * dot(m, g);
  }

  // Enhanced FBM with domain warping for more organic patterns
  float fbm(vec2 p, int octaves) {
    float value = 0.0;
    float amplitude = 0.5;
    float frequency = 1.0;
    
    // Domain warping parameters
    float warp = 0.35;
    
    for (int i = 0; i < octaves; i++) {
      // Apply domain warping for more organic flow
      vec2 q = p;
      q.x += 0.31 * sin(p.y * warp + time * speed * 0.15);
      q.y += 0.31 * sin(p.x * warp - time * speed * 0.15);
      
      value += amplitude * (snoise(q * frequency) * 0.5 + 0.5);
      
      // Evolve parameters for next octave
      amplitude *= 0.5;
      frequency *= 2.1;
      p = p * 1.15 + vec2(0.3, -0.4);
      warp *= 1.3;
    }
    
    return value;
  }
  
  // Smooth cell noise for additional texture
  float cellNoise(vec2 p) {
    vec2 ip = floor(p);
    vec2 fp = fract(p);
    
    float d = 1.0e10;
    for (int i = -1; i <= 1; i++) {
      for (int j = -1; j <= 1; j++) {
        vec2 offset = vec2(float(i), float(j));
        vec2 r = offset + vec2(
          snoise(ip + offset),
          snoise(ip + offset + vec2(31.1891, 17.3031))
        ) * 0.5 - fp;
        
        float d2 = dot(r, r);
        d = min(d, d2);
      }
    }
    
    return smoothstep(0.0, 1.0, d);
  }
  
  // Smooth sine wave function
  float smoothSin(float x, float phase, float amplitude) {
    return sin(x + phase) * amplitude;
  }
  
  // Flowing wave pattern with interference
  float wavePattern(vec2 uv, float time) {
    // Multi-layered waves with interference
    float wave1 = sin((uv.x * 10.0) + time * 1.5) * 0.5;
    float wave2 = sin((uv.y * 8.0) + time * 0.8) * 0.5;
    float wave3 = sin((uv.x * 6.0 + uv.y * 6.0) + time * 1.2) * 0.3;
    float wave4 = sin((uv.x * 12.0 - uv.y * 3.0) + time * 1.7) * 0.2;
    
    return (wave1 + wave2 + wave3 + wave4) * 0.4 + 0.5;
  }
  
  // Ribbon flow pattern
  float ribbonPattern(vec2 uv, float time) {
    float distToCenter = length(uv * 2.0 - 1.0);
    
    // Multiple moving ribbons with smooth transitions
    float a = atan(uv.y - 0.5, uv.x - 0.5) * 0.4;
    float r = length(uv - 0.5) * 2.0;
    
    float ribbon1 = smoothstep(0.1, 0.2, abs(sin(a * 8.0 + time + r * 3.0) * 0.8));
    float ribbon2 = smoothstep(0.05, 0.15, abs(sin(a * 5.0 - time * 0.7 + r * 2.0) * 0.8));
    float ribbon3 = smoothstep(0.02, 0.05, abs(sin(a * 12.0 + time * 1.3 + r * 4.0) * 0.5));
    
    return mix(ribbon1, ribbon2, 0.5) * (1.0 - ribbon3) * (1.0 - distToCenter * 0.7);
  }
  
  // Flowing particles system
  float particlePattern(vec2 uv, float time) {
    // Create a grid for particle system
    vec2 grid = fract(uv * 10.0) - 0.5;
    
    // Offset the grid over time with different speeds for each column/row
    float xOffset = snoise(vec2(uv.x * 2.0, time * 0.3)) * 0.3;
    float yOffset = snoise(vec2(uv.y * 2.0, time * 0.2)) * 0.3;
    
    // Calculate distance from particle centers
    vec2 particlePos = grid + vec2(xOffset, yOffset);
    float dist = length(particlePos);
    
    // Size oscillation
    float size = 0.15 + sin(time * speed + uv.x * 10.0 + uv.y * 8.0) * 0.05;
    
    // Particle glow
    float particle = smoothstep(size, size * 0.8, dist);
    
    // Add connections between particles
    float connections = smoothstep(0.4, 0.39, dist) * 0.15;
    
    return particle + connections;
  }
  
  // Dynamic flow pattern
  float flowPattern(vec2 uv, float time) {
    // Animate UV coordinates
    vec2 flow = uv + vec2(
      sin(uv.y * 4.0 + time * 0.5) * 0.1,
      cos(uv.x * 4.0 + time * 0.5) * 0.1
    );
    
    // Multiple noise layers flowing in different directions
    float noise1 = fbm(flow * 3.0 + vec2(time * 0.1, 0.0), 3);
    float noise2 = fbm(flow * 5.0 - vec2(time * 0.15, time * 0.05), 2);
    float noise3 = fbm(flow * 8.0 + vec2(0.0, time * 0.2), 1);
    
    // Combine noise layers with dynamic blending
    float blend = sin(time * 0.1) * 0.5 + 0.5;
    float finalNoise = mix(
      mix(noise1, noise2, 0.5),
      noise3,
      blend
    );
    
    // Create flowing stripes
    float stripes = smoothstep(0.45, 0.55, sin(finalNoise * 10.0 + time * 0.5) * 0.5 + 0.5);
    
    return mix(finalNoise, stripes, 0.3);
  }
  
  // Main shader function
  void main() {
    vec2 uv = gl_FragCoord.xy / resolution.xy;
    vec2 centeredUV = 2.0 * uv - 1.0;
    centeredUV.x *= resolution.x / resolution.y;
    
    float distFromCenter = length(centeredUV) * 0.6;
    
    // Mouse influence for interactive effects
    vec2 mouseUV = mousePosition / resolution.xy;
    mouseUV = mouseUV * 2.0 - 1.0;
    mouseUV.x *= resolution.x / resolution.y;
    float mouseDistance = length(centeredUV - mouseUV) * 2.0;
    float mouseInfluence = smoothstep(1.0, 0.0, mouseDistance);
    
    // Animate UV coordinates with flowing motion
    vec2 animatedUV = uv + vec2(
      sin(uv.y * 3.0 + time * speed * 0.2) * 0.03,
      cos(uv.x * 3.0 + time * speed * 0.3) * 0.03
    );
    
    // Apply mouse distortion if mouse is active
    if (mousePosition.x > 0.0) {
      animatedUV += (centeredUV - mouseUV) * 0.03 * mouseInfluence;
    }
    
    // Pattern selection and blending based on uniform
    float patternOutput = 0.0;
    
    if (pattern == 0) { // Waves
      patternOutput = wavePattern(animatedUV, time * speed);
    } else if (pattern == 1) { // Flow
      patternOutput = flowPattern(animatedUV, time * speed);
    } else if (pattern == 2) { // Particles
      patternOutput = particlePattern(animatedUV, time * speed);
    } else if (pattern == 3) { // Ribbons
      patternOutput = ribbonPattern(animatedUV, time * speed);
    } else {
      // Fallback to waves
      patternOutput = wavePattern(animatedUV, time * speed);
    }
    
    // Apply cell noise for additional detail in all patterns
    float cells = cellNoise(animatedUV * 8.0 + time * speed * 0.1);
    patternOutput = mix(patternOutput, cells, 0.1);
    
    // Apply intensity control with enhanced contrast
    float pattern = smoothstep(0.2, 0.8, patternOutput) * intensity;
    
    // Edge highlighting with improved definition
    float edge = abs(patternOutput - 0.5) * 2.0;
    edge = pow(edge, 2.0 + sin(time * speed * 0.5) * 0.5);
    
    // Dynamic vignette
    float vignetteStrength = 0.6 + 0.2 * sin(time * speed * 0.1);
    float vignette = 1.0 - smoothstep(0.4, 1.6 + 0.2 * sin(time * speed * 0.3), distFromCenter);
    vignette = pow(vignette, vignetteStrength);
    
    // Enhanced color blending with more dynamic gradients
    float colorBlend = smoothstep(0.3, 0.7, patternOutput) * colorIntensity;
    colorBlend += sin(patternOutput * 8.0 + time * speed) * 0.1;
    
    vec3 color = mix(
      baseColor,
      accentColor,
      colorBlend
    );
    
    // Apply subtle color variations
    color += vec3(
      sin(time * speed * 0.2) * 0.03,
      sin(time * speed * 0.3) * 0.03,
      sin(time * speed * 0.4) * 0.03
    ) * patternOutput;
    
    // Apply edge highlighting with color shift
    vec3 edgeColor = mix(accentColor, vec3(1.0), 0.5);
    color = mix(color, edgeColor, edge * 0.3 * intensity);
    
    // Apply vignette and dynamic alpha
    float alpha = clamp(pattern * vignette * 0.9, 0.0, 0.92);
    
    // Pulse effect
    alpha *= 0.85 + 0.15 * sin(time * speed * 0.2);
    
    // Mouse interaction highlight
    if (mousePosition.x > 0.0) {
      color = mix(color, accentColor, mouseInfluence * 0.3);
      alpha = mix(alpha, min(alpha + 0.2, 1.0), mouseInfluence * 0.5);
    }
    
    gl_FragColor = vec4(color, alpha);
  }
`;

const vertexShader = `
  varying vec2 vUv;
  
  void main() {
    vUv = uv;
    gl_Position = vec4(position, 1.0);
  }
`;

// Renderer cache to reduce instantiation overhead
let sharedRenderer: THREE.WebGLRenderer | null = null;

export const CreativeShaderBackground: React.FC<CreativeShaderBackgroundProps> = ({
  disableParallax = false,
  intensity = 1.0,
  speed = 1.0,
  colorIntensity = 0.7,
  pattern = 'waves'
}) => {
  const { theme } = useTheme();
  const { performanceSettings } = usePerformanceOptimizations();
  
  const containerRef = useRef<HTMLDivElement>(null);
  const rendererRef = useRef<THREE.WebGLRenderer | null>(null);
  const sceneRef = useRef<THREE.Scene | null>(null);
  const cameraRef = useRef<THREE.OrthographicCamera | null>(null);
  const materialRef = useRef<THREE.ShaderMaterial | null>(null);
  const geometryRef = useRef<THREE.PlaneGeometry | null>(null);
  const frameRef = useRef<number | null>(null);
  const mouseRef = useRef<{ x: number, y: number }>({ x: 0, y: 0 });
  const [visible, setVisible] = useState(false);
  const [dimensions, setDimensions] = useState({ width: 0, height: 0 });
  
  // Parse pattern type to shader-compatible format
  const patternValue = useMemo(() => {
    switch (pattern) {
      case 'waves': return 0;
      case 'flow': return 1;
      case 'particles': return 2;
      case 'ribbons': return 3;
      default: return 0;
    }
  }, [pattern]);
  
  // Convert theme colors to THREE.js vectors with enhanced palette
  const colors = useMemo(() => {
    const baseColor = new THREE.Color(theme.colors.background || '#121212');
    const accentColor = new THREE.Color(theme.colors.primary || '#6c63ff');
    
    // Adjust colors for better visual quality
    baseColor.multiplyScalar(0.8); // Slightly darken base color
    
    return {
      base: baseColor,
      accent: accentColor,
    };
  }, [theme.colors.background, theme.colors.primary]);
  
  // Set up initial scene
  useEffect(() => {
    if (!containerRef.current) return;
    
    // Skip animation for reduced motion settings
    if (performanceSettings.reduceMotion) {
      setVisible(true);
      return;
    }
    
    // Get container dimensions
    const updateDimensions = () => {
      if (!containerRef.current) return;
      const rect = containerRef.current.getBoundingClientRect();
      setDimensions({
        width: rect.width,
        height: rect.height
      });
    };
    
    updateDimensions();
    
    // Set up resize observer for responsive canvas
    const resizeObserver = new ResizeObserver(() => {
      updateDimensions();
      if (rendererRef.current && cameraRef.current && containerRef.current) {
        const { width, height } = containerRef.current.getBoundingClientRect();
        rendererRef.current.setSize(width, height);
        
        if (materialRef.current) {
          materialRef.current.uniforms.resolution.value.set(width, height);
        }
        
        if (width > 0 && height > 0) {
          // Only update when dimensions are valid
          rendererRef.current.setPixelRatio(window.devicePixelRatio);
        }
      }
    });
    
    if (containerRef.current) {
      resizeObserver.observe(containerRef.current);
    }
    
    // Create scene
    const scene = new THREE.Scene();
    sceneRef.current = scene;
    
    // Create orthographic camera
    const camera = new THREE.OrthographicCamera(-1, 1, 1, -1, 0.1, 10);
    camera.position.z = 1;
    cameraRef.current = camera;
    
    // Reuse renderer if possible for better performance
    if (!sharedRenderer) {
      sharedRenderer = new THREE.WebGLRenderer({
        antialias: performanceSettings.performanceTier !== 'low',
        alpha: true,
        powerPreference: 'high-performance',
        precision: performanceSettings.performanceTier === 'high' ? 'highp' : 'mediump'
      });
    }
    
    // Configure renderer
    const renderer = sharedRenderer;
    renderer.setClearColor(0x000000, 0);
    rendererRef.current = renderer;
    
    if (containerRef.current.firstChild) {
      containerRef.current.removeChild(containerRef.current.firstChild);
    }
    
    containerRef.current.appendChild(renderer.domElement);
    
    // Adjust for device pixel ratio
    const pixelRatio = Math.min(window.devicePixelRatio, 
      performanceSettings.performanceTier === 'high' ? 2 : 
      performanceSettings.performanceTier === 'medium' ? 1.5 : 1);
    renderer.setPixelRatio(pixelRatio);
    
    // Set up shader material with optimized settings
    const material = new THREE.ShaderMaterial({
      uniforms: {
        time: { value: 0 },
        resolution: { value: new THREE.Vector2(dimensions.width, dimensions.height) },
        baseColor: { value: colors.base },
        accentColor: { value: colors.accent },
        intensity: { value: intensity },
        speed: { value: speed },
        colorIntensity: { value: colorIntensity },
        pattern: { value: patternValue },
        mousePosition: { value: new THREE.Vector2(0, 0) }
      },
      vertexShader,
      fragmentShader,
      transparent: true,
      blending: THREE.NormalBlending,
      depthWrite: false,
      depthTest: false
    });
    
    materialRef.current = material;
    
    // Create geometry with appropriate detail level
    const geometry = new THREE.PlaneGeometry(2, 2);
    geometryRef.current = geometry;
    
    // Create and add mesh to scene
    const mesh = new THREE.Mesh(geometry, material);
    scene.add(mesh);
    
    // Fade in the effect
    setTimeout(() => setVisible(true), 300);
    
    // Update renderer size
    if (containerRef.current) {
      const { width, height } = containerRef.current.getBoundingClientRect();
      renderer.setSize(width, height);
    }
    
    // Set up pattern change animation
    let lastPatternValue = patternValue;
    const patternChangeTime = Date.now();
    
    // Animation transition helper function
    const animatePatternTransition = () => {
      const now = Date.now();
      const elapsed = (now - patternChangeTime) / 1000; // seconds
      const transitionDuration = 1.0; // seconds
      
      if (elapsed < transitionDuration && material.uniforms.pattern.value !== patternValue) {
        // Transition in progress - could implement cross-fade here if needed
        material.uniforms.pattern.value = patternValue;
      }
      
      if (lastPatternValue !== patternValue) {
        lastPatternValue = patternValue;
      }
    };
    
    return () => {
      // Clean up on unmount
      if (frameRef.current) {
        cancelAnimationFrame(frameRef.current);
        frameRef.current = null;
      }
      
      resizeObserver.disconnect();
      
      // Remove renderer from container
      if (containerRef.current?.contains(renderer.domElement)) {
        containerRef.current.removeChild(renderer.domElement);
      }
      
      // Properly clean up Three.js objects
      if (geometryRef.current) {
        geometryRef.current.dispose();
        geometryRef.current = null;
      }
      
      if (materialRef.current) {
        materialRef.current.dispose();
        materialRef.current = null;
      }
      
      // Note: We don't dispose the renderer here as it's shared
      // We also don't dispose the scene as it doesn't have a dispose method
      
      // Remove all children from the scene
      if (sceneRef.current) {
        while (sceneRef.current.children.length > 0) {
          const object = sceneRef.current.children[0];
          sceneRef.current.remove(object);
        }
        sceneRef.current = null;
      }
    };
  }, [performanceSettings.reduceMotion, performanceSettings.performanceTier, colors, dimensions.width, dimensions.height, intensity, speed, colorIntensity, patternValue]);
  
  // Update shader uniforms when theme or props change
  useEffect(() => {
    if (materialRef.current) {
      materialRef.current.uniforms.baseColor.value = colors.base;
      materialRef.current.uniforms.accentColor.value = colors.accent;
      materialRef.current.uniforms.intensity.value = intensity;
      materialRef.current.uniforms.speed.value = speed;
      materialRef.current.uniforms.colorIntensity.value = colorIntensity;
      materialRef.current.uniforms.pattern.value = patternValue;
    }
  }, [colors, intensity, speed, colorIntensity, patternValue]);
  
  // Animation loop with performance optimizations
  useEffect(() => {
    if (!rendererRef.current || !sceneRef.current || !cameraRef.current || !materialRef.current) return;
    if (performanceSettings.reduceMotion) return;
    
    let lastTime = 0;
    const targetFPS = performanceSettings.performanceTier === 'high' ? 60 : 
                     performanceSettings.performanceTier === 'medium' ? 45 : 30;
    const frameInterval = 1000 / targetFPS;
    
    // Pattern rotation timing
    const patternCycleDuration = 30000; // 30 seconds per pattern
    const startTime = Date.now();
    
    const animate = (currentTime: number) => {
      if (!materialRef.current || !rendererRef.current || !sceneRef.current || !cameraRef.current) return;
      frameRef.current = requestAnimationFrame(animate);
      
      // Skip frames based on performance settings
      const delta = currentTime - lastTime;
      if (delta < frameInterval) return;
      
      // Adjust time for consistent speed regardless of framerate
      lastTime = currentTime - (delta % frameInterval);
      
      // Update time uniform (convert to seconds and apply speed)
      materialRef.current.uniforms.time.value = currentTime * 0.001;
      
      // Apply mouse influence if parallax is enabled
      if (!disableParallax && mouseRef.current) {
        materialRef.current.uniforms.mousePosition.value.set(
          mouseRef.current.x * dimensions.width,
          (1 - mouseRef.current.y) * dimensions.height
        );
      } else {
        materialRef.current.uniforms.mousePosition.value.set(-1, -1); // Invalid position indicates no mouse
      }
      
      // Render the scene
      rendererRef.current.render(sceneRef.current, cameraRef.current);
    };
    
    frameRef.current = requestAnimationFrame(animate);
    
    return () => {
      if (frameRef.current) {
        cancelAnimationFrame(frameRef.current);
        frameRef.current = null;
      }
    };
  }, [
    performanceSettings.reduceMotion, 
    performanceSettings.performanceTier, 
    dimensions, 
    disableParallax
  ]);
  
  // Mouse movement parallax effect
  useEffect(() => {
    if (disableParallax || performanceSettings.reduceMotion) return;
    
    const handleMouseMove = (e: MouseEvent) => {
      if (!containerRef.current) return;
      
      // Calculate normalized mouse position
      const rect = containerRef.current.getBoundingClientRect();
      const mouseX = (e.clientX - rect.left) / rect.width;
      const mouseY = (e.clientY - rect.top) / rect.height;
      
      // Update mouse ref
      mouseRef.current = { x: mouseX, y: mouseY };
    };
    
    // Use passive listener for better performance
    window.addEventListener('mousemove', handleMouseMove, { passive: true });
    
    // Reset mouse position when cursor leaves the window
    const handleMouseLeave = () => {
      mouseRef.current = { x: -1, y: -1 };
    };
    
    document.addEventListener('mouseleave', handleMouseLeave);
    
    return () => {
      window.removeEventListener('mousemove', handleMouseMove);
      document.removeEventListener('mouseleave', handleMouseLeave);
    };
  }, [disableParallax, performanceSettings.reduceMotion]);
  
  // Touch event handling for mobile devices
  useEffect(() => {
    if (disableParallax || performanceSettings.reduceMotion) return;
    
    const handleTouch = (e: TouchEvent) => {
      if (e.touches.length > 0 && containerRef.current) {
        const touch = e.touches[0];
        const rect = containerRef.current.getBoundingClientRect();
        
        // Calculate normalized touch position
        const touchX = (touch.clientX - rect.left) / rect.width;
        const touchY = (touch.clientY - rect.top) / rect.height;
        
        // Update mouse ref with touch position
        mouseRef.current = { x: touchX, y: touchY };
        
        // Auto-reset touch influence after a short delay
        setTimeout(() => {
          mouseRef.current = { x: -1, y: -1 };
        }, 2000);
      }
    };
    
    const resetTouch = () => {
      mouseRef.current = { x: -1, y: -1 };
    };
    
    // Add touch event listeners
    if (containerRef.current) {
      containerRef.current.addEventListener('touchstart', handleTouch, { passive: true });
      containerRef.current.addEventListener('touchmove', handleTouch, { passive: true });
      containerRef.current.addEventListener('touchend', resetTouch, { passive: true });
    }
    
    return () => {
      if (containerRef.current) {
        containerRef.current.removeEventListener('touchstart', handleTouch);
        containerRef.current.removeEventListener('touchmove', handleTouch);
        containerRef.current.removeEventListener('touchend', resetTouch);
      }
    };
  }, [disableParallax, performanceSettings.reduceMotion]);
  
  // Clean up shared renderer when component unmounts
  useEffect(() => {
    return () => {
      // Clean up shared renderer on final unmount
      // This would typically be done at the application level,
      // but we're doing it here for completeness
      if (typeof window !== 'undefined' && sharedRenderer && document.querySelectorAll('.creative-shader-background').length <= 1) {
        sharedRenderer.dispose();
        sharedRenderer = null;
      }
    };
  }, []);
  
  return (
    <AnimationContainer ref={containerRef} visible={visible} className="creative-shader-background" />
  );
};

export default CreativeShaderBackground; 