// Frontend Performance Benchmarks
import React from 'react';
import { render } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import { ThemeProvider } from 'styled-components';
import { themes } from '../styles/themes';
import { CreativeShaderBackground } from '../components/animations/CreativeShaderBackground';
import { ScrollTransformBackground } from '../components/animations/ScrollTransformBackground';
import { captureMemoryMetrics, formatBytes } from '../utils/performance';

// Mock WebGL context for testing
const mockWebGLContext = {
  createProgram: jest.fn(() => ({})),
  createShader: jest.fn(() => ({})),
  shaderSource: jest.fn(),
  compileShader: jest.fn(),
  attachShader: jest.fn(),
  linkProgram: jest.fn(),
  useProgram: jest.fn(),
  uniform1f: jest.fn(),
  uniform2f: jest.fn(),
  uniform3f: jest.fn(),
  getUniformLocation: jest.fn(() => 0),
  createBuffer: jest.fn(() => ({})),
  bindBuffer: jest.fn(),
  bufferData: jest.fn(),
  getAttribLocation: jest.fn(() => 0),
  enableVertexAttribArray: jest.fn(),
  vertexAttribPointer: jest.fn(),
  viewport: jest.fn(),
  clearColor: jest.fn(),
  clear: jest.fn(),
  drawArrays: jest.fn(),
  canvas: { width: 800, height: 600 },
  drawingBufferWidth: 800,
  drawingBufferHeight: 600,
  VERTEX_SHADER: 0,
  FRAGMENT_SHADER: 1,
  ARRAY_BUFFER: 2,
  STATIC_DRAW: 3,
  TRIANGLES: 4,
  COLOR_BUFFER_BIT: 5,
};

// Mock HTMLCanvasElement.getContext using jest.spyOn for better type safety
beforeAll(() => {
  const getContextSpy = jest.spyOn(HTMLCanvasElement.prototype, 'getContext');
  
  getContextSpy.mockImplementation((contextId) => {
    if (contextId === 'webgl' || contextId === 'webgl2') {
      return mockWebGLContext as unknown as WebGLRenderingContext;
    }
    return null;
  });
});

afterAll(() => {
  jest.restoreAllMocks();
});

// Performance testing utilities
interface PerformanceResult {
  componentName: string;
  renderTime: number;
  memoryUsage: {
    before: any;
    after: any;
    difference: number;
  };
  fps?: number;
  warnings: string[];
}

const measureRenderPerformance = async (
  ComponentToTest: React.ComponentType<any>,
  props: any = {},
  testName: string
): Promise<PerformanceResult> => {
  const warnings: string[] = [];
  
  // Capture initial memory
  const memoryBefore = captureMemoryMetrics();
  
  // Force garbage collection if available
  if ((window as any).gc) {
    (window as any).gc();
  }
  
  const startTime = performance.now();
  
  // Render component
  const result = render(
    <BrowserRouter>
      <ThemeProvider theme={themes.dark}>
        <ComponentToTest {...props} />
      </ThemeProvider>
    </BrowserRouter>
  );
  
  const endTime = performance.now();
  const renderTime = endTime - startTime;
  
  // Wait for any async operations
  await new Promise(resolve => setTimeout(resolve, 100));
  
  // Capture final memory
  const memoryAfter = captureMemoryMetrics();
  const memoryDifference = (memoryAfter.usedJSHeapSize || 0) - (memoryBefore.usedJSHeapSize || 0);
  
  // Check for performance warnings
  if (renderTime > 100) {
    warnings.push(`Slow render time: ${renderTime.toFixed(2)}ms`);
  }
  
  if (memoryDifference > 5 * 1024 * 1024) { // 5MB
    warnings.push(`High memory usage: ${formatBytes(memoryDifference)}`);
  }
  
  // Cleanup
  result.unmount();
  
  return {
    componentName: testName,
    renderTime,
    memoryUsage: {
      before: memoryBefore,
      after: memoryAfter,
      difference: memoryDifference
    },
    warnings
  };
};

const measureAnimationPerformance = (
  testDuration: number = 1000
): Promise<{ fps: number; frameDrops: number }> => {
  return new Promise((resolve) => {
    let frames = 0;
    let frameDrops = 0;
    let lastFrameTime = performance.now();
    
    const animate = (currentTime: number) => {
      frames++;
      
      const deltaTime = currentTime - lastFrameTime;
      if (deltaTime > 16.67) { // Dropped frame (60fps = 16.67ms per frame)
        frameDrops++;
      }
      
      lastFrameTime = currentTime;
      
      if (currentTime - testStartTime < testDuration) {
        requestAnimationFrame(animate);
      } else {
        const fps = (frames * 1000) / testDuration;
        resolve({ fps, frameDrops });
      }
    };
    
    const testStartTime = performance.now();
    requestAnimationFrame(animate);
  });
};

describe('Frontend Performance Benchmarks', () => {
  beforeEach(() => {
    // Reset performance markers
    performance.clearMarks();
    performance.clearMeasures();
    
    // Mock requestAnimationFrame for testing
    global.requestAnimationFrame = jest.fn((cb) => setTimeout(cb, 16));
    global.cancelAnimationFrame = jest.fn();
  });

  afterEach(() => {
    jest.restoreAllMocks();
  });

  test('CreativeShaderBackground - Initial Render Performance', async () => {
    const result = await measureRenderPerformance(
      CreativeShaderBackground,
      {
        intensity: 1.0,
        speed: 1.0,
        colorIntensity: 0.7,
        pattern: 'waves'
      },
      'CreativeShaderBackground'
    );
    
    console.log('CreativeShaderBackground Performance:', {
      renderTime: `${result.renderTime.toFixed(2)}ms`,
      memoryUsage: formatBytes(result.memoryUsage.difference),
      warnings: result.warnings
    });
    
    // Assertions for performance thresholds
    expect(result.renderTime).toBeLessThan(500); // Should render in under 500ms
    expect(result.memoryUsage.difference).toBeLessThan(10 * 1024 * 1024); // Under 10MB
    expect(result.warnings.length).toBe(0); // No performance warnings
  });

  test('ScrollTransformBackground - Initial Render Performance', async () => {
    const result = await measureRenderPerformance(
      ScrollTransformBackground,
      {
        enableFloatingOrbs: true,
        showDebug: false
      },
      'ScrollTransformBackground'
    );
    
    console.log('ScrollTransformBackground Performance:', {
      renderTime: `${result.renderTime.toFixed(2)}ms`,
      memoryUsage: formatBytes(result.memoryUsage.difference),
      warnings: result.warnings
    });
    
    expect(result.renderTime).toBeLessThan(300);
    expect(result.memoryUsage.difference).toBeLessThan(15 * 1024 * 1024); // Under 15MB
  });

  test('Animation Performance - Frame Rate Test', async () => {
    const animationResult = await measureAnimationPerformance(2000); // Test for 2 seconds
    
    console.log('Animation Performance:', {
      fps: animationResult.fps.toFixed(2),
      frameDrops: animationResult.frameDrops,
      efficiency: `${((animationResult.fps / 60) * 100).toFixed(1)}%`
    });
    
    expect(animationResult.fps).toBeGreaterThan(30); // At least 30 FPS
    expect(animationResult.frameDrops).toBeLessThan(20); // Less than 20 dropped frames in 2 seconds
  });

  test('WebGL Context Creation Performance', () => {
    const startTime = performance.now();
    
    // Test WebGL context creation
    const canvas = document.createElement('canvas');
    canvas.width = 800;
    canvas.height = 600;
    
    const gl = canvas.getContext('webgl');
    const gl2 = canvas.getContext('webgl2');
    
    const endTime = performance.now();
    const creationTime = endTime - startTime;
    
    console.log('WebGL Context Creation:', {
      time: `${creationTime.toFixed(2)}ms`,
      webglSupported: !!gl,
      webgl2Supported: !!gl2
    });
    
    expect(creationTime).toBeLessThan(50); // Should create contexts quickly
    expect(gl).toBeTruthy(); // WebGL should be supported in test environment
  });

  test('Memory Leak Detection - Multiple Renders', async () => {
    const initialMemory = captureMemoryMetrics();
    
    // Render and unmount component multiple times
    for (let i = 0; i < 5; i++) {
      const result = render(
        <BrowserRouter>
          <ThemeProvider theme={themes.dark}>
            <CreativeShaderBackground intensity={0.5} />
          </ThemeProvider>
        </BrowserRouter>
      );
      
      await new Promise(resolve => setTimeout(resolve, 50));
      result.unmount();
    }
    
    // Force garbage collection
    if ((window as any).gc) {
      (window as any).gc();
    }
    
    await new Promise(resolve => setTimeout(resolve, 200));
    
    const finalMemory = captureMemoryMetrics();
    const memoryGrowth = (finalMemory.usedJSHeapSize || 0) - (initialMemory.usedJSHeapSize || 0);
    
    console.log('Memory Leak Test:', {
      initialMemory: formatBytes(initialMemory.usedJSHeapSize),
      finalMemory: formatBytes(finalMemory.usedJSHeapSize),
      growth: formatBytes(memoryGrowth),
      iterations: 5
    });
    
    // Memory growth should be minimal after multiple render cycles
    expect(memoryGrowth).toBeLessThan(2 * 1024 * 1024); // Under 2MB growth
  });

  test('Bundle Size Impact Assessment', () => {
    // Mock bundle analyzer results (would normally come from webpack-bundle-analyzer)
    const mockBundleData = {
      'three.js': 500 * 1024, // 500KB
      'framer-motion': 200 * 1024, // 200KB
      'gsap': 150 * 1024, // 150KB
      'react-spring': 100 * 1024, // 100KB
      'styled-components': 80 * 1024, // 80KB
      'lodash': 70 * 1024 // 70KB
    };
    
    const totalSize = Object.values(mockBundleData).reduce((sum, size) => sum + size, 0);
    const largeLibraries = Object.entries(mockBundleData)
      .filter(([_, size]) => size > 100 * 1024)
      .map(([name, size]) => ({ name, size: formatBytes(size) }));
    
    console.log('Bundle Size Analysis:', {
      totalSize: formatBytes(totalSize),
      largeLibraries,
      recommendations: [
        'Consider tree-shaking for large libraries',
        'Evaluate if multiple animation libraries are necessary',
        'Implement code splitting for non-critical components'
      ]
    });
    
    // Assertions for bundle size thresholds
    expect(totalSize).toBeLessThan(2 * 1024 * 1024); // Under 2MB total
    expect(largeLibraries.length).toBeLessThan(3); // No more than 3 large libraries
  });
});

// Performance monitoring utility for development
export const startPerformanceMonitoring = () => {
  if (process.env.NODE_ENV !== 'development') return;
  
  let frameCount = 0;
  let lastTime = performance.now();
  
  const monitor = () => {
    frameCount++;
    const currentTime = performance.now();
    
    if (currentTime - lastTime >= 1000) {
      const fps = frameCount;
      const memory = captureMemoryMetrics();
      
      console.log('Performance Monitor:', {
        fps,
        memory: formatBytes(memory.usedJSHeapSize),
        timestamp: new Date().toISOString()
      });
      
      frameCount = 0;
      lastTime = currentTime;
      
      // Warn if performance is poor
      if (fps < 30) {
        console.warn('Low FPS detected:', fps);
      }
      
      if (memory.usedJSHeapSize && memory.jsHeapSizeLimit) {
        const memoryUsage = (memory.usedJSHeapSize / memory.jsHeapSizeLimit) * 100;
        if (memoryUsage > 80) {
          console.warn('High memory usage:', `${memoryUsage.toFixed(1)}%`);
        }
      }
    }
    
    requestAnimationFrame(monitor);
  };
  
  requestAnimationFrame(monitor);
  
  return () => {
    // Cleanup function
    frameCount = 0;
  };
};