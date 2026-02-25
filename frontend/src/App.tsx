import React, { Suspense, lazy, useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route, useLocation } from 'react-router-dom';
import { HelmetProvider } from 'react-helmet-async';
import { ThemeProvider as StyledThemeProvider } from 'styled-components';
import { GlobalStyles } from './styles/GlobalStyles';
import NavigationBar from './components/NavigationBar/NavigationBar';
import Footer from './components/Footer/Footer';
import ScrollToTop from './components/navigation/ScrollToTop';
import ScrollToTopButton from './components/common/ScrollToTopButton';
import ScrollProgressIndicator from './components/common/ScrollProgressIndicator';
import styled from 'styled-components';
import PageTop from './components/navigation/PageTop';
import { Layout } from './components/layout/Layout';
import SmartSkeleton from './components/skeletons/SmartSkeleton';
import { usePreloader } from './hooks/usePreloader';
import { devCacheManager } from './utils/devCacheManager';
import { StoreInitializer } from './components/StoreInitializer';
import { useTheme } from './hooks/useTheme';
import ErrorNotification from './components/notifications/ErrorNotification';
import AuthModal from './components/auth/AuthModal';
import { useAuthStore } from './stores/authStore';
import PageTransition from './components/navigation/PageTransition';
import { useVisitorTracking } from './hooks/useVisitorTracking';
import PortfolioLayout from './components/layout/PortfolioLayout';
import LenisProvider from './providers/LenisProvider';
import ErrorBoundary from './components/common/ErrorBoundary';

const PortfolioHome = lazy(() => import('./pages/PortfolioHome'));
const Contact = lazy(() => import('./pages/Contact/Contact'));
const DevPanel = lazy(() => import('./pages/devpanel/DevPanel'));
const Messaging = lazy(() => import('./pages/messaging/Messaging'));
const UrlShortener = lazy(() => import('./pages/urlshortener/UrlShortener'));
const NotFound = lazy(() => import('./pages/NotFound/NotFound'));
const AboutPage = lazy(() => import('./pages/About/About'));
const ProjectsPage = lazy(() => import('./pages/Projects'));
const Status = lazy(() => import('./pages/Status/Status'));
const BlogPage = lazy(() => import('./pages/Blog/Blog'));
const BlogPostPage = lazy(() => import('./pages/Blog/BlogPost'));

const AppContainer = styled.div`
  max-width: 100vw;
  width: 100%;
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  background: ${({ theme }) => theme.colors.background};
  color: ${({ theme }) => theme.colors.text};
  --page-max-width: 1200px;
  --content-max-width: 1000px;
`;

const SkipLink = styled.a`
  position: absolute;
  top: 1rem;
  left: 1rem;
  z-index: 10000;
  background: ${({ theme }) => theme.colors.primary};
  color: white;
  padding: 1rem;
  border-radius: ${({ theme }) => theme.borderRadius.small};
  text-decoration: none;
  font-weight: 500;
  transform: translateY(-200%);
  transition: transform 0.3s ease;

  &:focus {
    transform: translateY(0);
  }
`;

function PortfolioPage() {
  const { theme } = useTheme();

  return (
    <StyledThemeProvider theme={theme}>
      <GlobalStyles theme={theme} />
      <LenisProvider>
        <PortfolioLayout>
          <Suspense fallback={<SmartSkeleton type="hero" />}>
            <PortfolioHome />
          </Suspense>
        </PortfolioLayout>
      </LenisProvider>
    </StyledThemeProvider>
  );
}

function StandardLayout({ children }: { children: React.ReactNode }) {
  const { theme, themeMode, toggleTheme } = useTheme();
  const { authModalOpen, authModalMode, setAuthModalOpen } = useAuthStore();

  return (
    <StyledThemeProvider theme={theme}>
      <Layout>
        <PageTop />
        <SkipLink href="#main-content">Skip to content</SkipLink>
        <AppContainer>
          <GlobalStyles theme={theme} />
          <ErrorNotification />
          <NavigationBar
            themeMode={themeMode}
            toggleTheme={toggleTheme}
          />
          <ScrollProgressIndicator hideThreshold={50} />
          <div id="main-content" className="content">
            <PageTransition>
              {children}
            </PageTransition>
          </div>
          <Footer />
          <ScrollToTopButton />
          <AuthModal
            isOpen={authModalOpen}
            onClose={() => setAuthModalOpen(false)}
            initialMode={authModalMode}
          />
        </AppContainer>
      </Layout>
    </StyledThemeProvider>
  );
}

function AppContent() {
  const location = useLocation();
  const isPortfolioHome = location.pathname === '/';

  useVisitorTracking();

  usePreloader({
    enableSmartPreloading: true,
    enableRoutePreloading: true,
    enableHoverPreloading: true,
    preloadDelay: 150
  });

  useEffect(() => {
    if (process.env.NODE_ENV === 'development') {
      devCacheManager.setupDevTools();
      const hasInitialized = sessionStorage.getItem('dev-cache-initialized');
      if (!hasInitialized) {
        devCacheManager.clearAllCaches();
        sessionStorage.setItem('dev-cache-initialized', 'true');
      }
    }
  }, []);

  if (isPortfolioHome) {
    return <PortfolioPage />;
  }

  return (
    <StandardLayout>
      <ErrorBoundary>
        <Suspense fallback={<SmartSkeleton />}>
          <Routes>
            <Route path="/about" element={<AboutPage />} />
            <Route path="/contact" element={<Contact />} />
            <Route path="/projects" element={<ProjectsPage />} />
            <Route path="/portfolio" element={<ProjectsPage />} />
            <Route path="/devpanel" element={<DevPanel />} />
            <Route path="/urlshortener" element={<UrlShortener />} />
            <Route path="/messaging" element={<Messaging />} />
            <Route path="/status" element={<Status />} />
            <Route path="/blog" element={<BlogPage />} />
            <Route path="/blog/:slug" element={<BlogPostPage />} />

            <Route path="*" element={<NotFound />} />
          </Routes>
        </Suspense>
      </ErrorBoundary>
    </StandardLayout>
  );
}

function App() {
  return (
    <HelmetProvider>
      <StoreInitializer>
        <Router>
          <ScrollToTop />
          <AppContent />
        </Router>
      </StoreInitializer>
    </HelmetProvider>
  );
}

export default App;
