import { scrollToTop, scrollToElement, scrollToPosition, navigateToSection } from './scrollConfig';

interface ScrollTestResult {
  feature: string;
  status: 'pass' | 'fail';
  details?: string;
}

export const scrollTestUtils = {
  // Test basic scroll functionality
  testBasicScroll: async (): Promise<ScrollTestResult[]> => {
    const results: ScrollTestResult[] = [];
    
    // Test scroll to top
    try {
      const initialPosition = window.pageYOffset;
      scrollToTop({ behavior: 'auto' });
      await new Promise(resolve => setTimeout(resolve, 100));
      const newPosition = window.pageYOffset;
      
      results.push({
        feature: 'Scroll to Top (instant)',
        status: newPosition === 0 ? 'pass' : 'fail',
        details: `Initial: ${initialPosition}px, Final: ${newPosition}px`
      });
    } catch (error) {
      results.push({
        feature: 'Scroll to Top (instant)',
        status: 'fail',
        details: error instanceof Error ? error.message : 'Unknown error'
      });
    }
    
    // Test smooth scroll
    try {
      window.scrollTo(0, 500); // Set initial position
      await new Promise(resolve => setTimeout(resolve, 100));
      
      scrollToTop({ behavior: 'smooth' });
      await new Promise(resolve => setTimeout(resolve, 600)); // Wait for animation
      const finalPosition = window.pageYOffset;
      
      results.push({
        feature: 'Scroll to Top (smooth)',
        status: finalPosition < 10 ? 'pass' : 'fail',
        details: `Final position: ${finalPosition}px`
      });
    } catch (error) {
      results.push({
        feature: 'Scroll to Top (smooth)',
        status: 'fail',
        details: error instanceof Error ? error.message : 'Unknown error'
      });
    }
    
    return results;
  },
  
  // Test scroll to element
  testScrollToElement: async (): Promise<ScrollTestResult[]> => {
    const results: ScrollTestResult[] = [];
    
    const testElements = ['hero', 'about', 'skills', 'projects'];
    
    for (const elementId of testElements) {
      try {
        const element = document.getElementById(elementId);
        if (element) {
          scrollToElement(element, { behavior: 'auto' });
          await new Promise(resolve => setTimeout(resolve, 100));
          
          const rect = element.getBoundingClientRect();
          const isVisible = rect.top >= 0 && rect.top <= window.innerHeight;
          
          results.push({
            feature: `Scroll to #${elementId}`,
            status: isVisible ? 'pass' : 'fail',
            details: `Element top: ${rect.top}px`
          });
        } else {
          results.push({
            feature: `Scroll to #${elementId}`,
            status: 'fail',
            details: 'Element not found'
          });
        }
      } catch (error) {
        results.push({
          feature: `Scroll to #${elementId}`,
          status: 'fail',
          details: error instanceof Error ? error.message : 'Unknown error'
        });
      }
    }
    
    return results;
  },
  
  // Test browser compatibility
  testBrowserCompatibility: (): ScrollTestResult[] => {
    const results: ScrollTestResult[] = [];
    
    // Check ScrollBehavior support
    results.push({
      feature: 'ScrollBehavior API',
      status: 'scrollBehavior' in document.documentElement.style ? 'pass' : 'fail',
      details: 'Native smooth scroll support'
    });
    
    // Check IntersectionObserver support
    results.push({
      feature: 'IntersectionObserver API',
      status: 'IntersectionObserver' in window ? 'pass' : 'fail',
      details: 'Required for lazy loading and scroll reveal'
    });
    
    // Check requestAnimationFrame support
    results.push({
      feature: 'requestAnimationFrame',
      status: 'requestAnimationFrame' in window ? 'pass' : 'fail',
      details: 'Required for smooth animations'
    });
    
    return results;
  },
  
  // Run all tests
  runAllTests: async (): Promise<void> => {
    console.log('🧪 Running Scroll Tests...\n');
    
    // Browser compatibility
    console.log('📋 Browser Compatibility:');
    const compatResults = scrollTestUtils.testBrowserCompatibility();
    compatResults.forEach(result => {
      const icon = result.status === 'pass' ? '✅' : '❌';
      console.log(`${icon} ${result.feature}: ${result.details || ''}`);
    });
    
    // Basic scroll tests
    console.log('\n📜 Basic Scroll Tests:');
    const basicResults = await scrollTestUtils.testBasicScroll();
    basicResults.forEach(result => {
      const icon = result.status === 'pass' ? '✅' : '❌';
      console.log(`${icon} ${result.feature}: ${result.details || ''}`);
    });
    
    // Element scroll tests
    console.log('\n🎯 Element Scroll Tests:');
    const elementResults = await scrollTestUtils.testScrollToElement();
    elementResults.forEach(result => {
      const icon = result.status === 'pass' ? '✅' : '❌';
      console.log(`${icon} ${result.feature}: ${result.details || ''}`);
    });
    
    console.log('\n✨ Scroll tests completed!');
  }
};

// Attach to window in development
if (process.env.NODE_ENV === 'development' && typeof window !== 'undefined') {
  (window as any).scrollTests = scrollTestUtils;
  console.log('💡 Scroll test utilities available. Run: window.scrollTests.runAllTests()');
}