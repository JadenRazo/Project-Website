import React, { Suspense, lazy } from 'react';
import type { LazyExoticComponent, ComponentType } from 'react';
import SmartSkeleton from '../../components/skeletons/SmartSkeleton';
import { useComponentLoading } from '../../hooks/useLoadingState';

// Define proper types for lazy-loaded components
type LazyComponent = LazyExoticComponent<ComponentType<any>>;

// Lazy load components with proper typing
const Hero: LazyComponent = lazy(() => import('../../components/sections/Hero'));
const Projects: LazyComponent = lazy(() => import('../../components/sections/Projects'));
const About: LazyComponent = lazy(() => import('../../components/sections/About'));
const SkillsSection: LazyComponent = lazy(() => import('../../components/sections/SkillsSection'));

// Main Home component
const Home = () => {
  const { setLoading, isLoading } = useComponentLoading('Home');

  return (
    <div className="home-container">
      <Suspense fallback={<SmartSkeleton type="hero" />}>
        <section id="hero" className="section">
          <Hero />
        </section>
      </Suspense>
      
      <Suspense fallback={<SmartSkeleton type="skills" />}>
        <section id="skills" className="section content-section">
          <SkillsSection />
        </section>
      </Suspense>
      
      <Suspense fallback={<SmartSkeleton type="projects" />}>
        <section id="projects" className="section content-section">
          <Projects />
        </section>
      </Suspense>
      
      <Suspense fallback={<SmartSkeleton type="about" />}>
        <section id="about" className="section content-section">
          <About />
        </section>
      </Suspense>
    </div>
  );
};

export default Home; 