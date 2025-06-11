import React, { useRef, useEffect, useState, useMemo } from 'react';
import * as THREE from 'three';
import styled from 'styled-components';
import { useTheme } from '../../hooks/useTheme';
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
`;

interface OptimizedShaderBackgroundProps {
  intensity?: number;
  speed?: number;
  colorIntensity?: number;
  pattern?: 'waves' | 'flow' | 'simple';
}

// Simplified and optimized fragment shader (reduced from 600+ lines to ~100 lines)
const optimizedFragmentShader = `
  uniform float time;
  uniform vec2 resolution;
  uniform vec3 baseColor;
  uniform vec3 accentColor;
  uniform float intensity;
  uniform float speed;
  uniform float colorIntensity;
  uniform int pattern;
  
  // Simplified noise function (much faster than previous implementation)
  float hash(vec2 p) {
    return fract(sin(dot(p, vec2(12.9898, 78.233))) * 43758.5453);
  }
  
  float noise(vec2 p) {
    vec2 i = floor(p);
    vec2 f = fract(p);
    f = f * f * (3.0 - 2.0 * f); // Smooth interpolation
    
    float a = hash(i);
    float b = hash(i + vec2(1.0, 0.0));
    float c = hash(i + vec2(0.0, 1.0));
    float d = hash(i + vec2(1.0, 1.0));
    
    return mix(mix(a, b, f.x), mix(c, d, f.x), f.y);
  }
  
  // Simple wave pattern
  float wavePattern(vec2 uv, float time) {
    float wave1 = sin(uv.x * 8.0 + time * speed) * 0.5;
    float wave2 = sin(uv.y * 6.0 + time * speed * 0.8) * 0.3;
    return (wave1 + wave2) * 0.5 + 0.5;
  }
  
  // Simplified flow pattern
  float flowPattern(vec2 uv, float time) {
    vec2 flow = uv + vec2(
      sin(uv.y * 3.0 + time * speed * 0.5) * 0.1,
      cos(uv.x * 3.0 + time * speed * 0.5) * 0.1
    );
    
    return noise(flow * 4.0 + time * speed * 0.1);
  }
  
  // Very simple pattern for low-end devices
  float simplePattern(vec2 uv, float time) {
    return sin(uv.x * 4.0 + time * speed) * sin(uv.y * 4.0 + time * speed * 0.7) * 0.5 + 0.5;
  }
  
  void main() {
    vec2 uv = gl_FragCoord.xy / resolution.xy;
    
    float patternOutput = 0.0;
    
    if (pattern == 0) { // Waves
      patternOutput = wavePattern(uv, time);
    } else if (pattern == 1) { // Flow
      patternOutput = flowPattern(uv, time);
    } else { // Simple
      patternOutput = simplePattern(uv, time);
    }
    
    // Apply intensity
    patternOutput *= intensity;
    
    // Simple color blending
    vec3 color = mix(baseColor, accentColor, patternOutput * colorIntensity);
    
    // Simple vignette
    float dist = length(uv - 0.5);
    float vignette = 1.0 - smoothstep(0.3, 0.8, dist);
    
    // Final alpha with vignette
    float alpha = patternOutput * vignette * 0.6;
    
    gl_FragColor = vec4(color, alpha);
  }
`;

const vertexShader = `
  void main() {
    gl_Position = vec4(position, 1.0);
  }
`;

// Renderer pool to reuse WebGL contexts
class RendererPool {
  private static instance: RendererPool;
  private renderers: THREE.WebGLRenderer[] = [];
  private inUse: Set<THREE.WebGLRenderer> = new Set();
  
  static getInstance(): RendererPool {
    if (!RendererPool.instance) {
      RendererPool.instance = new RendererPool();
    }
    return RendererPool.instance;
  }
  
  getRenderer(performanceTier: string): THREE.WebGLRenderer {
    // Find available renderer
    const available = this.renderers.find(renderer => !this.inUse.has(renderer));
    
    if (available) {
      this.inUse.add(available);
      return available;
    }
    
    // Create new renderer if none available
    const renderer = new THREE.WebGLRenderer({
      antialias: performanceTier !== 'low',
      alpha: true,
      powerPreference: performanceTier === 'high' ? 'high-performance' : 'default',
      precision: performanceTier === 'high' ? 'highp' : 'mediump'
    });
    
    this.renderers.push(renderer);
    this.inUse.add(renderer);
    
    return renderer;
  }
  
  releaseRenderer(renderer: THREE.WebGLRenderer): void {
    this.inUse.delete(renderer);
  }
  
  cleanup(): void {
    this.renderers.forEach(renderer => {
      renderer.dispose();
    });
    this.renderers = [];
    this.inUse.clear();
  }
}

export const OptimizedShaderBackground: React.FC<OptimizedShaderBackgroundProps> = ({
  intensity = 1.0,
  speed = 1.0,
  colorIntensity = 0.7,
  pattern = 'simple'
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
  const [visible, setVisible] = useState(false);
  const [dimensions, setDimensions] = useState({ width: 0, height: 0 });
  
  // Get renderer pool instance
  const rendererPool = useMemo(() => RendererPool.getInstance(), []);
  
  // Parse pattern type to shader-compatible format
  const patternValue = useMemo(() => {
    switch (pattern) {
      case 'waves': return 0;
      case 'flow': return 1;
      case 'simple': return 2;
      default: return 2; // Default to simple
    }
  }, [pattern]);
  
  // Convert theme colors to THREE.js vectors
  const colors = useMemo(() => {
    const baseColor = new THREE.Color(theme.colors.background || '#121212');
    const accentColor = new THREE.Color(theme.colors.primary || '#6c63ff');
    
    baseColor.multiplyScalar(0.8);
    
    return { base: baseColor, accent: accentColor };
  }, [theme.colors.background, theme.colors.primary]);
  
  // Set up scene with optimized settings
  useEffect(() => {
    if (!containerRef.current) return;
    
    // Skip animation for reduced motion
    if (performanceSettings.reduceMotion) {
      setVisible(true);
      return;
    }
    
    // Force simple pattern on low performance devices
    const effectivePattern = performanceSettings.performanceTier === 'low' ? 2 : patternValue;
    
    const updateDimensions = () => {
      if (!containerRef.current) return;
      const rect = containerRef.current.getBoundingClientRect();
      setDimensions({ width: rect.width, height: rect.height });
    };
    
    updateDimensions();
    
    // Set up resize observer
    const resizeObserver = new ResizeObserver(updateDimensions);
    if (containerRef.current) {
      resizeObserver.observe(containerRef.current);
    }
    
    // Create scene
    const scene = new THREE.Scene();
    sceneRef.current = scene;
    
    // Create camera
    const camera = new THREE.OrthographicCamera(-1, 1, 1, -1, 0.1, 10);
    camera.position.z = 1;
    cameraRef.current = camera;
    
    // Get renderer from pool
    const renderer = rendererPool.getRenderer(performanceSettings.performanceTier);
    renderer.setClearColor(0x000000, 0);
    rendererRef.current = renderer;
    
    if (containerRef.current.firstChild) {
      containerRef.current.removeChild(containerRef.current.firstChild);
    }
    containerRef.current.appendChild(renderer.domElement);
    
    // Optimize pixel ratio based on performance tier
    const pixelRatio = Math.min(
      window.devicePixelRatio,
      performanceSettings.performanceTier === 'high' ? 2 : 1
    );
    renderer.setPixelRatio(pixelRatio);
    
    // Create optimized shader material
    const material = new THREE.ShaderMaterial({
      uniforms: {
        time: { value: 0 },
        resolution: { value: new THREE.Vector2(dimensions.width, dimensions.height) },
        baseColor: { value: colors.base },
        accentColor: { value: colors.accent },
        intensity: { value: intensity },
        speed: { value: speed },
        colorIntensity: { value: colorIntensity },
        pattern: { value: effectivePattern }
      },
      vertexShader,
      fragmentShader: optimizedFragmentShader,
      transparent: true,
      blending: THREE.NormalBlending,
      depthWrite: false,
      depthTest: false
    });
    
    materialRef.current = material;
    
    // Create geometry
    const geometry = new THREE.PlaneGeometry(2, 2);
    geometryRef.current = geometry;
    
    // Create mesh
    const mesh = new THREE.Mesh(geometry, material);
    scene.add(mesh);
    
    // Fade in
    setTimeout(() => setVisible(true), 100);
    
    // Update size
    if (containerRef.current) {
      const { width, height } = containerRef.current.getBoundingClientRect();
      renderer.setSize(width, height);
    }
    
    return () => {
      // Cleanup
      if (frameRef.current) {
        cancelAnimationFrame(frameRef.current);
      }
      
      resizeObserver.disconnect();
      
      // Remove from container
      if (containerRef.current?.contains(renderer.domElement)) {
        containerRef.current.removeChild(renderer.domElement);
      }
      
      // Dispose resources
      if (geometryRef.current) {
        geometryRef.current.dispose();
      }
      if (materialRef.current) {
        materialRef.current.dispose();
      }
      
      // Return renderer to pool
      rendererPool.releaseRenderer(renderer);
      
      // Clean up scene
      if (sceneRef.current) {
        while (sceneRef.current.children.length > 0) {
          sceneRef.current.remove(sceneRef.current.children[0]);
        }
      }
    };
  }, [performanceSettings, colors, dimensions.width, dimensions.height, intensity, speed, colorIntensity, patternValue, rendererPool]);
  
  // Update uniforms when props change
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
  
  // Optimized animation loop with frame skipping
  useEffect(() => {
    if (!rendererRef.current || !sceneRef.current || !cameraRef.current || !materialRef.current) return;
    if (performanceSettings.reduceMotion) return;
    
    let lastTime = 0;
    const targetFPS = performanceSettings.performanceTier === 'high' ? 60 : 
                     performanceSettings.performanceTier === 'medium' ? 45 : 30;
    const frameInterval = 1000 / targetFPS;
    
    const animate = (currentTime: number) => {
      if (!materialRef.current || !rendererRef.current || !sceneRef.current || !cameraRef.current) return;
      
      frameRef.current = requestAnimationFrame(animate);
      
      // Frame skipping for performance
      const delta = currentTime - lastTime;
      if (delta < frameInterval) return;
      
      lastTime = currentTime - (delta % frameInterval);
      
      // Update time uniform
      materialRef.current.uniforms.time.value = currentTime * 0.001;
      
      // Update resolution if changed
      if (materialRef.current.uniforms.resolution.value.x !== dimensions.width ||
          materialRef.current.uniforms.resolution.value.y !== dimensions.height) {
        materialRef.current.uniforms.resolution.value.set(dimensions.width, dimensions.height);
        rendererRef.current.setSize(dimensions.width, dimensions.height);
      }
      
      // Render
      rendererRef.current.render(sceneRef.current, cameraRef.current);
    };
    
    frameRef.current = requestAnimationFrame(animate);
    
    return () => {
      if (frameRef.current) {
        cancelAnimationFrame(frameRef.current);
      }
    };
  }, [performanceSettings, dimensions]);
  
  return (
    <AnimationContainer 
      ref={containerRef} 
      visible={visible} 
      className="optimized-shader-background" 
    />
  );
};

export default OptimizedShaderBackground;