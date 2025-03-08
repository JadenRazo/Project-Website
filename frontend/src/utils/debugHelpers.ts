/**
 * Debug helper utilities for troubleshooting layout and animations
 */

// Visualize z-index stacking in the browser console
export const visualizeZIndexStacking = (): void => {
  if (process.env.NODE_ENV !== 'development') {
    console.warn('Z-index visualization only available in development mode');
    return;
  }
  
  // Get all elements with non-auto z-index
  const elementsWithZIndex: Array<{ el: Element, zIndex: number }> = [];
  const allElements = document.querySelectorAll('*');
  
  allElements.forEach(el => {
    const style = window.getComputedStyle(el);
    const zIndex = style.zIndex;
    
    if (zIndex !== 'auto') {
      const zIndexNum = parseInt(zIndex, 10);
      if (!isNaN(zIndexNum)) {
        elementsWithZIndex.push({ el, zIndex: zIndexNum });
      }
    }
  });
  
  // Sort by z-index value
  elementsWithZIndex.sort((a, b) => a.zIndex - b.zIndex);
  
  // Group by z-index
  const groups: Record<number, Element[]> = {};
  elementsWithZIndex.forEach(({ el, zIndex }) => {
    if (!groups[zIndex]) {
      groups[zIndex] = [];
    }
    groups[zIndex].push(el);
  });
  
  // Log nicely formatted stacking order
  console.log('%c Z-Index Stacking Visualization ', 'background: #333; color: #bada55; font-size: 16px;');
  console.log('-------------------------------------');
  
  Object.entries(groups).sort((a, b) => parseInt(a[0]) - parseInt(b[0])).forEach(([zIndex, elements]) => {
    console.groupCollapsed(`z-index: ${zIndex} (${elements.length} elements)`);
    elements.forEach(el => {
      const classAttr = el.getAttribute('class') || '';
      const classes = classAttr.split(' ').filter(Boolean).join('.');
      const idAttr = el.getAttribute('id');
      const tagName = el.tagName.toLowerCase();
      
      const selector = [
        tagName,
        idAttr ? `#${idAttr}` : '',
        classes ? `.${classes}` : ''
      ].filter(Boolean).join('');
      
      console.log(selector, el);
    });
    console.groupEnd();
  });
};

// Temporarily highlight elements to see their positions and z-index
export const highlightElementsWithZIndex = (targetZIndex?: number): void => {
  if (process.env.NODE_ENV !== 'development') {
    console.warn('Element highlighting only available in development mode');
    return;
  }
  
  // Remove any existing highlights
  const existingHighlights = document.querySelectorAll('.z-index-highlight');
  existingHighlights.forEach(el => el.remove());
  
  // Get all elements with non-auto z-index
  const allElements = document.querySelectorAll('*');
  
  allElements.forEach(el => {
    const style = window.getComputedStyle(el);
    const zIndex = style.zIndex;
    
    if (zIndex !== 'auto') {
      const zIndexNum = parseInt(zIndex, 10);
      
      if (!isNaN(zIndexNum) && (targetZIndex === undefined || zIndexNum === targetZIndex)) {
        const rect = el.getBoundingClientRect();
        
        // Create highlight element
        const highlight = document.createElement('div');
        highlight.className = 'z-index-highlight';
        highlight.style.position = 'absolute';
        highlight.style.top = `${rect.top + window.scrollY}px`;
        highlight.style.left = `${rect.left + window.scrollX}px`;
        highlight.style.width = `${rect.width}px`;
        highlight.style.height = `${rect.height}px`;
        highlight.style.border = '2px solid red';
        highlight.style.backgroundColor = 'rgba(255, 0, 0, 0.1)';
        highlight.style.zIndex = '9999';
        highlight.style.pointerEvents = 'none';
        
        // Add label with z-index value
        const label = document.createElement('div');
        label.style.position = 'absolute';
        label.style.top = '0';
        label.style.left = '0';
        label.style.backgroundColor = '#333';
        label.style.color = 'white';
        label.style.padding = '2px 5px';
        label.style.fontSize = '10px';
        label.style.borderRadius = '0 0 3px 0';
        label.textContent = `z: ${zIndexNum}`;
        
        highlight.appendChild(label);
        document.body.appendChild(highlight);
        
        // Remove after 5 seconds
        setTimeout(() => highlight.remove(), 5000);
      }
    }
  });
  
  console.log('Elements highlighted for 5 seconds');
};

// Check for WebGL support and capabilities
export const checkWebGLSupport = (): void => {
  const canvas = document.createElement('canvas');
  let gl: WebGLRenderingContext | null = null;
  
  try {
    // Cast to WebGLRenderingContext to fix TypeScript error
    gl = (canvas.getContext('webgl') || canvas.getContext('experimental-webgl')) as WebGLRenderingContext | null;
  } catch (e) {
    console.error('Error creating WebGL context:', e);
  }
  
  if (!gl) {
    console.warn('WebGL not supported. Animations may be limited.');
    return;
  }
  
  const debugInfo = gl.getExtension('WEBGL_debug_renderer_info');
  
  console.log('%c WebGL Support Information ', 'background: #333; color: #bada55; font-size: 16px;');
  console.log('-------------------------------------');
  
  if (debugInfo) {
    const vendor = gl.getParameter(debugInfo.UNMASKED_VENDOR_WEBGL);
    const renderer = gl.getParameter(debugInfo.UNMASKED_RENDERER_WEBGL);
    console.log('GPU Vendor:', vendor);
    console.log('GPU Renderer:', renderer);
  } else {
    console.log('GPU info not available');
  }
  
  console.log('WebGL Version:', gl.getParameter(gl.VERSION));
  console.log('GLSL Version:', gl.getParameter(gl.SHADING_LANGUAGE_VERSION));
  console.log('Max Texture Size:', gl.getParameter(gl.MAX_TEXTURE_SIZE));
  console.log('Max Viewport Dimensions:', gl.getParameter(gl.MAX_VIEWPORT_DIMS));
  
  const extensions = gl.getSupportedExtensions();
  console.groupCollapsed('Supported Extensions:');
  extensions?.forEach(ext => console.log(ext));
  console.groupEnd();
  
  // Check for potential issues
  const maxTextureSize = gl.getParameter(gl.MAX_TEXTURE_SIZE);
  if (maxTextureSize < 4096) {
    console.warn('Limited texture size may cause rendering issues with large images.');
  }
};

// Check canvas performance
export const benchmarkCanvasPerformance = (duration: number = 5000): void => {
  const canvas = document.createElement('canvas');
  canvas.width = 1000;
  canvas.height = 1000;
  const ctx = canvas.getContext('2d');
  
  if (!ctx) {
    console.error('Canvas 2D context not available.');
    return;
  }
  
  const results = {
    totalFrames: 0,
    avgFps: 0,
    minFps: Infinity,
    maxFps: 0
  };
  
  const frameTimes: number[] = [];
  let lastTime = performance.now();
  let running = true;
  
  console.log(`Starting canvas benchmark for ${duration}ms...`);
  
  const drawFrame = () => {
    if (!running) return;
    
    const currentTime = performance.now();
    const deltaTime = currentTime - lastTime;
    lastTime = currentTime;
    
    const fps = 1000 / deltaTime;
    frameTimes.push(fps);
    
    results.totalFrames++;
    results.minFps = Math.min(results.minFps, fps);
    results.maxFps = Math.max(results.maxFps, fps);
    
    // Draw random shapes
    for (let i = 0; i < 100; i++) {
      ctx.beginPath();
      ctx.fillStyle = `rgba(${Math.random() * 255}, ${Math.random() * 255}, ${Math.random() * 255}, 0.5)`;
      ctx.arc(Math.random() * canvas.width, Math.random() * canvas.height, Math.random() * 20 + 5, 0, Math.PI * 2);
      ctx.fill();
    }
    
    requestAnimationFrame(drawFrame);
  };
  
  drawFrame();
  
  setTimeout(() => {
    running = false;
    
    // Calculate average FPS
    results.avgFps = frameTimes.reduce((sum, fps) => sum + fps, 0) / frameTimes.length;
    
    console.log('%c Canvas Performance Results ', 'background: #333; color: #bada55; font-size: 16px;');
    console.log('-------------------------------------');
    console.log(`Total Frames: ${results.totalFrames}`);
    console.log(`Average FPS: ${results.avgFps.toFixed(2)}`);
    console.log(`Min FPS: ${results.minFps.toFixed(2)}`);
    console.log(`Max FPS: ${results.maxFps.toFixed(2)}`);
    
    // Provide performance assessment
    if (results.avgFps > 55) {
      console.log('%c Excellent canvas performance. Complex animations should work well.', 'color: green');
    } else if (results.avgFps > 30) {
      console.log('%c Good canvas performance. Most animations should work properly.', 'color: orange');
    } else {
      console.log('%c Poor canvas performance. Consider reducing animation complexity.', 'color: red');
    }
  }, duration);
};

export default {
  visualizeZIndexStacking,
  highlightElementsWithZIndex,
  checkWebGLSupport,
  benchmarkCanvasPerformance
}; 