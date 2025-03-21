import React, { Suspense, lazy } from 'react';
import type { LazyExoticComponent, ComponentType } from 'react';

// Define proper types for lazy-loaded components
type LazyComponent = LazyExoticComponent<ComponentType<any>>;

// Lazy load components with proper typing
const Hero: LazyComponent = lazy(() => import('../../components/sections/Hero'));
const Projects: LazyComponent = lazy(() => import('../../components/sections/Projects'));
const About: LazyComponent = lazy(() => import('../../components/sections/About'));
const SkillsSection: LazyComponent = lazy(() => import('../../components/sections/SkillsSection'));

// Loading fallback component
const LoadingFallback = () => (
  <div style={{ 
    height: '100vh', 
    display: 'flex', 
    justifyContent: 'center', 
    alignItems: 'center' 
  }}>
    Loading...
  </div>
);

// Main Home component
const Home = () => {
  return (
    <div className="home-container">
      <Suspense fallback={<LoadingFallback />}>
        <section id="hero" className="section">
          <Hero />
        </section>
      </Suspense>
      
      <Suspense fallback={<LoadingFallback />}>
        <section id="skills" className="section content-section">
          <SkillsSection />
        </section>
      </Suspense>
      
      <Suspense fallback={<LoadingFallback />}>
        <section id="projects" className="section content-section">
          <Projects />
        </section>
      </Suspense>
      
      <Suspense fallback={<LoadingFallback />}>
        <section id="about" className="section content-section">
          <About />
        </section>
      </Suspense>
    </div>
  );
};

export default Home; 