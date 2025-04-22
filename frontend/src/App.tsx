// /Project-Website/frontend/src/App.tsx
import React, { useEffect, useRef } from 'react';
import { BrowserRouter as Router, Routes, Route, useLocation, useNavigationType } from 'react-router-dom';
import { ThemeProvider, useTheme } from './contexts/ThemeContext';
import { Hero } from './components/sections/Hero';
import { About } from './components/sections/About';
import { GlobalStyles } from './styles/GlobalStyles';
import { Projects } from './components/sections/Projects';
import { SkillsSection } from './components/sections/SkillsSection';
import NavigationBar from './components/NavigationBar/NavigationBar';
import Contact from './pages/Contact/Contact';
import Footer from './components/Footer/Footer';
import ScrollToTop from './components/navigation/ScrollToTop';
import styled from 'styled-components';
import { usePerformanceOptimizations } from './hooks/usePerformanceOptimizations';
import { Layout } from './components/layout/Layout';
import DevPanel from './pages/devpanel/DevPanel';
import Messaging from './pages/messaging/Messaging';
import UrlShortener from './pages/urlshortener/UrlShortener';
import NotFound from './pages/NotFound/NotFound';
import Home from './pages/Home/Home';
import AboutPage from './pages/About/About';

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
    <Hero />
    <SkillsSection />
    <Projects />
    <About />
  </>
);

// Main application content with routing
function AppContent() {
  // Get theme context and performance settings
  const { theme, themeMode, toggleTheme } = useTheme();
  const { performanceSettings } = usePerformanceOptimizations();
  
  return (
    <Layout>
      <AppContainer>
        <GlobalStyles theme={theme} />
        <NavigationBar 
          themeMode={themeMode} 
          toggleTheme={toggleTheme} 
        />
        
        <div className="content">
          <Routes>
            {/* Main routes */}
            <Route path="/" element={<HomePage />} />
            <Route path="/about" element={<AboutPage />} />
            <Route path="/contact" element={<Contact />} />
            <Route path="/projects" element={<Projects />} />
            <Route path="/skills" element={<SkillsSection />} />
            <Route path="/home" element={<Home />} />

            {/* Application routes */}
            <Route path="/devpanel" element={<DevPanel />} />
            <Route path="/urlshortener" element={<UrlShortener />} />
            <Route path="/messaging" element={<Messaging />} />
            
            {/* 404 route */}
            <Route path="*" element={<NotFound />} />
          </Routes>
        </div>
        
        <Footer />
      </AppContainer>
    </Layout>
  );
}

// Root App component with provider
function App() {
  return (
    <ThemeProvider>
      <Router>
        <ScrollToTop />
        <AppContent />
      </Router>
    </ThemeProvider>
  );
}

export default App;
