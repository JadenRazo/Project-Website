// /Project-Website/frontend/src/App.tsx
import React, { Suspense, lazy, useEffect } from 'react';
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

// Lazy load pages and sections for better performance
const Hero = lazy(() => import('./components/sections/Hero').then(module => ({ default: module.Hero })));
const About = lazy(() => import('./components/sections/About').then(module => ({ default: module.About })));
const SkillsSection = lazy(() => import('./components/sections/SkillsSection').then(module => ({ default: module.SkillsSection })));
const Contact = lazy(() => import('./pages/Contact/Contact'));
const DevPanel = lazy(() => import('./pages/devpanel/DevPanel'));
const Messaging = lazy(() => import('./pages/messaging/Messaging'));
const UrlShortener = lazy(() => import('./pages/urlshortener/UrlShortener'));
const NotFound = lazy(() => import('./pages/NotFound/NotFound'));
const Home = lazy(() => import('./pages/Home/Home'));
const AboutPage = lazy(() => import('./pages/About/About'));
const ProjectsPage = lazy(() => import('./pages/Projects'));
const Status = lazy(() => import('./pages/Status/Status'));

const AppContainer = styled.div`
  max-width: 100vw;
  width: 100%;
  overflow-x: hidden;
  contain: layout style;
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  background: ${({ theme }) => theme.colors.background};
  color: ${({ theme }) => theme.colors.text};
  --page-max-width: 1200px;
  --content-max-width: 1000px;
`;

// Home page component with properly ordered sections
const HomePage = () => (
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

// Main application content with routing
function AppContent() {
  // Get theme context
  const { theme, themeMode, toggleTheme } = useTheme();
  
  // Initialize preloader for optimized loading
  usePreloader({
    enableSmartPreloading: true,
    enableRoutePreloading: true,
    enableHoverPreloading: true,
    preloadDelay: 150
  });
  
  // Initialize dev cache manager in development
  useEffect(() => {
    if (process.env.NODE_ENV === 'development') {
      devCacheManager.setupDevTools();
      // Only clear caches on first load, not on every hot reload
      const hasInitialized = sessionStorage.getItem('dev-cache-initialized');
      if (!hasInitialized) {
        devCacheManager.clearAllCaches();
        sessionStorage.setItem('dev-cache-initialized', 'true');
      }
    }
  }, []);
  
  return (
    <StyledThemeProvider theme={theme}>
      <Layout>
        <AppContainer>
          <GlobalStyles theme={theme} />
          <NavigationBar 
            themeMode={themeMode} 
            toggleTheme={toggleTheme} 
          />
          
          <div className="content">
            <Suspense fallback={<SmartSkeleton />}>
              <Routes>
                {/* Main routes */}
                <Route path="/" element={<HomePage />} />
                <Route path="/about" element={<AboutPage />} />
                <Route path="/contact" element={<Contact />} />
                <Route path="/projects" element={<ProjectsPage />} />
                <Route path="/portfolio" element={<ProjectsPage />} />
                <Route path="/skills" element={<SkillsSection />} />
                <Route path="/home" element={<Home />} />

                {/* Application routes */}
                <Route path="/devpanel" element={<DevPanel />} />
                <Route path="/urlshortener" element={<UrlShortener />} />
                <Route path="/messaging" element={<Messaging />} />
                <Route path="/status" element={<Status />} />
                
                {/* 404 route */}
                <Route path="*" element={<NotFound />} />
              </Routes>
            </Suspense>
          </div>
          
          <Footer />
        </AppContainer>
      </Layout>
    </StyledThemeProvider>
  );
}

// Root App component with provider
function App() {
  return (
    <StoreInitializer>
      <Router>
        <ScrollToTop />
        <AppContent />
      </Router>
    </StoreInitializer>
  );
}

export default App;
