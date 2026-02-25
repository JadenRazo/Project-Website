// Optimized version of App.tsx with performance improvements
import React, { Suspense, lazy, useEffect, useMemo } from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { ThemeProvider as StyledThemeProvider } from 'styled-components';
import { GlobalStyles } from './styles/GlobalStyles';
import NavigationBar from './components/NavigationBar/NavigationBar';
import Footer from './components/Footer/Footer';
import ScrollToTop from './components/navigation/ScrollToTop';
import styled from 'styled-components';
import { Layout } from './components/layout/Layout';
import SmartSkeleton from './components/skeletons/SmartSkeleton';
import { usePreloader } from './hooks/usePreloader';
import { devCacheManager } from './utils/devCacheManager';
import { StoreInitializer } from './components/StoreInitializer';
import { useTheme } from './hooks/useTheme';
import { useOptimizedScrollHandler } from './hooks/useOptimizedScrollHandler';

// Optimized imports - removed unused animation libraries
import { OptimizedShaderBackground } from './components/animations/OptimizedShaderBackground';

// Lazy load pages with better chunk splitting
const Hero = lazy(() => 
  import('./components/sections/Hero').then(module => ({ default: module.Hero }))
);
const About = lazy(() => 
  import('./components/sections/About').then(module => ({ default: module.About }))
);
const SkillsSection = lazy(() => 
  import('./components/sections/SkillsSection').then(module => ({ default: module.SkillsSection }))
);
const Contact = lazy(() => import('./pages/Contact/Contact'));
const DevPanel = lazy(() => import('./pages/devpanel/DevPanel'));
const Messaging = lazy(() => import('./pages/messaging/Messaging'));
const UrlShortener = lazy(() => import('./pages/urlshortener/UrlShortener'));
const NotFound = lazy(() => import('./pages/NotFound/NotFound'));
const AboutPage = lazy(() => import('./pages/About/About'));
const ProjectsPage = lazy(() => import('./pages/Projects'));
const Status = lazy(() => import('./pages/Status/Status'));

const OptimizedAppContainer = styled.div`
  max-width: 100vw;
  width: 100%;
  overflow-x: hidden;
  contain: layout style paint; // Better containment
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  background: ${({ theme }) => theme.colors.background};
  color: ${({ theme }) => theme.colors.text};
  --page-max-width: 1200px;
  --content-max-width: 1000px;
  
  // Hardware acceleration
  transform: translateZ(0);
  will-change: transform;
`;

// Optimized Background Container with reduced complexity
const BackgroundContainer = styled.div`
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100vh;
  z-index: -1;
  pointer-events: none;
  overflow: hidden;
`;

// Optimized Home page component with better memory management
const OptimizedHomePage = React.memo(() => {
  return (
    <>
      <Suspense fallback={<SmartSkeleton type="hero" />}>
        <Hero />
      </Suspense>
      <Suspense fallback={<SmartSkeleton type="skills" />}>
        <SkillsSection />
      </Suspense>
      <Suspense fallback={<SmartSkeleton type="about" />}>
        <About />
      </Suspense>
    </>
  );
});

OptimizedHomePage.displayName = 'OptimizedHomePage';

// Main application content with optimized performance
function OptimizedAppContent() {
  const { theme, themeMode, toggleTheme } = useTheme();
  
  // Optimized preloader configuration
  usePreloader({
    enableSmartPreloading: true,
    enableRoutePreloading: true,
    enableHoverPreloading: false, // Disabled for better performance
    preloadDelay: 300
  });
  
  // Use optimized scroll handler instead of complex scroll transforms
  useOptimizedScrollHandler(
    (scrollState) => {
      // Simple scroll-based effects only when needed
      if (scrollState.isScrolling) {
        // Minimal scroll effects here
        document.body.style.setProperty('--scroll-progress', scrollState.scrollProgress.toString());
      }
    },
    {
      throttleMs: 16, // 60fps
      enableParallax: false // Disabled for performance
    }
  );
  
  // Optimized dev cache management
  useEffect(() => {
    if (process.env.NODE_ENV === 'development') {
      const hasInitialized = sessionStorage.getItem('dev-cache-initialized');
      if (!hasInitialized) {
        // Async cache clearing to not block rendering
        setTimeout(() => {
          devCacheManager.setupDevTools();
          devCacheManager.clearAllCaches();
          sessionStorage.setItem('dev-cache-initialized', 'true');
        }, 100);
      }
    }
  }, []);
  
  // Memoize background component to prevent unnecessary re-renders
  const backgroundComponent = useMemo(() => (
    <BackgroundContainer>
      <OptimizedShaderBackground
        intensity={0.8}
        speed={0.5}
        colorIntensity={0.6}
        pattern="simple" // Use simple pattern for better performance
      />
    </BackgroundContainer>
  ), []);
  
  return (
    <StyledThemeProvider theme={theme}>
      <Layout>
        <OptimizedAppContainer>
          <GlobalStyles theme={theme} />
          
          {/* Optimized background */}
          {backgroundComponent}
          
          <NavigationBar 
            themeMode={themeMode} 
            toggleTheme={toggleTheme} 
          />
          
          <main className="content" role="main">
            <Suspense fallback={<SmartSkeleton />}>
              <Routes>
                {/* Main routes */}
                <Route path="/" element={<OptimizedHomePage />} />
                <Route path="/about" element={<AboutPage />} />
                <Route path="/contact" element={<Contact />} />
                <Route path="/projects" element={<ProjectsPage />} />
                <Route path="/portfolio" element={<ProjectsPage />} />
                <Route path="/skills" element={<SkillsSection />} />
                {/* Application routes */}
                <Route path="/devpanel" element={<DevPanel />} />
                <Route path="/urlshortener" element={<UrlShortener />} />
                <Route path="/messaging" element={<Messaging />} />
                <Route path="/status" element={<Status />} />
                
                {/* 404 route */}
                <Route path="*" element={<NotFound />} />
              </Routes>
            </Suspense>
          </main>
          
          <Footer />
        </OptimizedAppContainer>
      </Layout>
    </StyledThemeProvider>
  );
}

// Root App component with optimized providers
function OptimizedApp() {
  return (
    <StoreInitializer>
      <Router>
        <ScrollToTop />
        <OptimizedAppContent />
      </Router>
    </StoreInitializer>
  );
}

export default OptimizedApp;