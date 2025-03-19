// src/App.tsx
import React, { useEffect, useState, useRef, memo, useCallback, Suspense } from 'react';
import { motion, AnimatePresence, useScroll, useTransform } from 'framer-motion';
import { 
  BrowserRouter as Router, 
  Routes, 
  Route, 
  useLocation, 
  useNavigate 
} from 'react-router-dom';
import { ThemeProvider } from './contexts/ThemeContext';
import { Hero } from './components/sections/Hero';
import { About } from './components/sections/About';
import { ZIndexProvider, useZIndex } from './hooks/useZIndex';
import { GlobalStyles } from './styles/GlobalStyles';
import { Projects } from './components/sections/Projects';
import type { Theme, ThemeMode } from './styles/theme.types';
import { lightTheme, darkTheme } from './styles/themes';
import { LoadingScreen } from './components/animations/LoadingScreen';
import { ScrollTransformBackground } from './components/animations/ScrollTransformBackground';
import { ScrollIndicator } from './components/animations/ScrollIndicator';
import { BurgerMenu } from './components/navigation/BurgerMenu';
import { Layout } from './components/layout/Layout';
import { SkillsSection } from './components/sections/SkillsSection';
import { useTheme } from './contexts/ThemeContext';
import { useThemeToggle } from './hooks/useThemeToggle';
import NavigationBar from './components/NavigationBar/NavigationBar';
import Home from './pages/Home/Home';
import AboutPage from './pages/About/About';
import Contact from './pages/Contact/Contact';
import NotFound from './pages/NotFound/NotFound';
import Footer from './components/Footer/Footer';
import DevPanel from './pages/devpanel/DevPanel';
import UrlShortener from './pages/urlshortener/UrlShortener';
import Messaging from './pages/messaging/Messaging';

// Types
interface Project {
  id: string;
  title: string;
  description: string;
  image: string;
  link: string;
  language: string;
}

interface ScrollTransform {
  type: 'opacity' | 'translateY' | 'translateX' | 'scale' | 'rotate';
  target: string;
  from: string;
  to: string;
  unit?: string;
  easing?: string;
}

interface ScrollSection {
  id: string;
  startPercent: number;
  endPercent: number;
  transforms: ScrollTransform[];
}

interface ProjectsSectionProps {
  themeMode: ThemeMode;
}

interface AppContentProps {
  themeMode: ThemeMode;
  toggleTheme: () => void;
}

// Enhanced animation variants with proper spring physics
const ANIMATION_VARIANTS = {
  pageTransition: {
    initial: { opacity: 0, y: 20 },
    animate: { 
      opacity: 1, 
      y: 0, 
      transition: { 
        type: 'spring', 
        damping: 20, 
        stiffness: 100
      } 
    },
    exit: { 
      opacity: 0, 
      y: -20, 
      transition: { 
        type: 'spring', 
        damping: 25, 
        stiffness: 120
      } 
    }
  },
  staggerContainer: {
    animate: { transition: { staggerChildren: 0.1 } }
  }
};

// Sample projects data (you can replace with your actual data)
const PROJECTS_DATA: Project[] = [
  {
    id: '1',
    title: "Project 1",
    description: "This is the first project.",
    image: "https://via.placeholder.com/300",
    link: "#",
    language: "JavaScript"
  },
  // ... other projects
];

const LANGUAGES = Object.freeze(['All', 'JavaScript', 'Python', 'Java']);

// ProjectsSection component with isolated state
const ProjectsSection = memo<ProjectsSectionProps>(({ themeMode }) => {
  const [selectedLanguage, setSelectedLanguage] = useState<string>('All');
  
  const handleLanguageSelect = useCallback((language: string) => {
    setSelectedLanguage(language);
  }, []);
  
  return (
    <Projects 
      projects={PROJECTS_DATA}
      languages={LANGUAGES}
      selectedLanguage={selectedLanguage}
      onLanguageChange={handleLanguageSelect}
    />
  );
});

ProjectsSection.displayName = 'ProjectsSection';

// Section observer hook to track visible sections
const useSectionObserver = (sectionIds: string[]): string | null => {
  const [activeSection, setActiveSection] = useState<string | null>(null);
  const [userHasScrolled, setUserHasScrolled] = useState(false);
  const previousPosition = useRef(0);
  const scrollTimeoutRef = useRef<number | null>(null);
  
  useEffect(() => {
    // Track if user has manually scrolled
    const handleUserScroll = () => {
      setUserHasScrolled(true);
      
      // Reset the flag after some time to allow auto-scrolling again
      if (scrollTimeoutRef.current) {
        window.clearTimeout(scrollTimeoutRef.current);
      }
      
      scrollTimeoutRef.current = window.setTimeout(() => {
        setUserHasScrolled(false);
      }, 1000) as unknown as number;
    };
    
    window.addEventListener('wheel', handleUserScroll, { passive: true });
    window.addEventListener('touchmove', handleUserScroll, { passive: true });
    
    return () => {
      window.removeEventListener('wheel', handleUserScroll);
      window.removeEventListener('touchmove', handleUserScroll);
      if (scrollTimeoutRef.current) {
        window.clearTimeout(scrollTimeoutRef.current);
      }
    };
  }, []);
  
  useEffect(() => {
    if (sectionIds.length === 0) return;
    
    const handleIntersect = (entries: IntersectionObserverEntry[]) => {
      // If user is actively scrolling, just track position without forcing scroll
      if (userHasScrolled) {
        previousPosition.current = window.scrollY;
        return;
      }
      
      entries.forEach(entry => {
        if (entry.isIntersecting) {
          const sectionId = entry.target.id;
          
          if (sectionIds.includes(sectionId)) {
            setActiveSection(sectionId);
            
            // Only update URL without forcing scroll
            if (window.history && window.location.pathname) {
              window.history.replaceState(
                null,
                document.title,
                `${window.location.pathname}${sectionId ? `#${sectionId}` : ''}`
              );
            }
          }
        }
      });
    };
    
    const options = {
      rootMargin: '-10% 0px -10% 0px',
      threshold: 0.2
    };
    
    const observer = new IntersectionObserver(handleIntersect, options);
    
    sectionIds.forEach(id => {
      const element = document.getElementById(id);
      if (element) observer.observe(element);
    });
    
    return () => observer.disconnect();
  }, [sectionIds, userHasScrolled]);
  
  return activeSection;
};

// Scroll controls hook for manual navigation
const useScrollControls = () => {
  const [isScrollable, setIsScrollable] = useState(true);
  const [scrollDirection, setScrollDirection] = useState<'up' | 'down' | null>(null);
  const lastScrollPosition = useRef(0);
  const scrollLockTimeoutRef = useRef<number | null>(null);
  
  // Add back scrollToSection for compatibility with existing code
  const scrollToSection = useCallback((sectionId: string) => {
    const section = document.getElementById(sectionId);
    if (section) {
      // Get the element's position
      const rect = section.getBoundingClientRect();
      const scrollTop = window.scrollY || document.documentElement.scrollTop;
      const targetPosition = scrollTop + rect.top - 80; // Adjust for any header offset
      
      // Use native smooth scrolling instead of scrollIntoView
      window.scrollTo({
        top: targetPosition,
        behavior: 'smooth'
      });
    }
  }, []);
  
  useEffect(() => {
    const handleScroll = () => {
      const currentScrollPosition = window.scrollY;
      
      if (lastScrollPosition.current < currentScrollPosition) {
        setScrollDirection('down');
      } else if (lastScrollPosition.current > currentScrollPosition) {
        setScrollDirection('up');
      }
      
      lastScrollPosition.current = currentScrollPosition;
    };
    
    // Allow normal scrolling by default
    document.body.style.overflow = isScrollable ? 'auto' : 'hidden';
    
    window.addEventListener('scroll', handleScroll, { passive: true });
    
    return () => {
      window.removeEventListener('scroll', handleScroll);
      document.body.style.overflow = 'auto';
      if (scrollLockTimeoutRef.current) {
        window.clearTimeout(scrollLockTimeoutRef.current);
      }
    };
  }, [isScrollable]);
  
  const lockScroll = useCallback(() => {
    setIsScrollable(false);
    document.body.style.overflow = 'hidden';
  }, []);
  
  const unlockScroll = useCallback(() => {
    setIsScrollable(true);
    document.body.style.overflow = 'auto';
  }, []);
  
  const temporarilyLockScroll = useCallback((durationMs: number = 800) => {
    lockScroll();
    
    if (scrollLockTimeoutRef.current) {
      window.clearTimeout(scrollLockTimeoutRef.current);
    }
    
    scrollLockTimeoutRef.current = window.setTimeout(() => {
      unlockScroll();
      scrollLockTimeoutRef.current = null;
    }, durationMs) as unknown as number;
  }, [lockScroll, unlockScroll]);
  
  return {
    isScrollable,
    scrollDirection,
    lockScroll,
    unlockScroll,
    temporarilyLockScroll,
    scrollToSection
  };
};

// Implement a scroll restoration function to handle browser history navigation

const useScrollRestoration = () => {
  const scrollPositions = useRef<Record<string, number>>({});
  
  useEffect(() => {
    // Save scroll position before navigating away
    const handleBeforeUnload = () => {
      scrollPositions.current[window.location.pathname] = window.scrollY;
      sessionStorage.setItem('scrollPositions', JSON.stringify(scrollPositions.current));
    };
    
    // Restore scroll position when navigating back
    const handleLoad = () => {
      const savedPositions = sessionStorage.getItem('scrollPositions');
      if (savedPositions) {
        scrollPositions.current = JSON.parse(savedPositions);
        
        const savedPosition = scrollPositions.current[window.location.pathname];
        if (savedPosition) {
          // Allow the page to render properly before scrolling
          setTimeout(() => {
            window.scrollTo(0, savedPosition);
          }, 100);
        }
      }
    };
    
    window.addEventListener('beforeunload', handleBeforeUnload);
    window.addEventListener('load', handleLoad);
    
    // Handle browser back/forward buttons
    window.addEventListener('popstate', () => {
      const savedPosition = scrollPositions.current[window.location.pathname];
      if (savedPosition) {
        setTimeout(() => {
          window.scrollTo(0, savedPosition);
        }, 100);
      }
    });
    
    return () => {
      window.removeEventListener('beforeunload', handleBeforeUnload);
      window.removeEventListener('load', handleLoad);
    };
  }, []);
};

// Main App Content component with route handling
const AppContent: React.FC<AppContentProps> = ({ themeMode, toggleTheme }) => {
  const location = useLocation();
  const { scrollY } = useScroll();
  const sections = ['hero', 'skills', 'projects', 'about'];
  const defaultSectionId = 'hero';
  const activeSection = useSectionObserver(sections) || defaultSectionId;
  const { scrollToSection } = useScrollControls();
  
  // State to track viewport width
  const [isMobile, setIsMobile] = useState(false);
  
  // Set initial mobile state and add resize listener
  useEffect(() => {
    const checkMobile = () => {
      setIsMobile(window.innerWidth <= 768);
    };
    
    // Set initial value
    checkMobile();
    
    // Add resize listener
    window.addEventListener('resize', checkMobile);
    
    // Cleanup
    return () => {
      window.removeEventListener('resize', checkMobile);
    };
  }, []);
  
  // Sync URL with active section
  useEffect(() => {
    if (location.pathname === '/' && activeSection && activeSection !== 'hero') {
      window.history.replaceState(null, '', `#${activeSection}`);
    }
  }, [activeSection, location.pathname]);
  
  // Handle direct navigation from URL hash
  useEffect(() => {
    if (location.hash) {
      const sectionId = location.hash.substring(1);
      if (sections.includes(sectionId)) {
        setTimeout(() => scrollToSection(sectionId), 100);
      }
    }
  }, [location.hash, scrollToSection, sections]);
  
  // Parallax scroll effect for background
  const backgroundY = useTransform(scrollY, [0, 1000], [0, -200]);
  
  return (
    <motion.div 
      className="app-content content-wrapper"
      variants={ANIMATION_VARIANTS.pageTransition}
      initial="initial"
      animate="animate"
      exit="exit"
    >
      <motion.div 
        className="parallax-background"
        style={{ y: backgroundY }}
      />
      
      <AnimatePresence mode="wait">
        <Routes location={location} key={location.pathname}>
          <Route path="/" element={
            <>
              <section id="hero" className="section">
                <Hero />
                <div
                  style={{
                    position: 'relative',
                    width: '100%',
                    display: 'flex',
                    justifyContent: 'center',
                    alignItems: 'center',
                    marginTop: isMobile ? '10px' : '-90px',
                    left: 0,
                    right: 0,
                    zIndex: 40,
                    pointerEvents: 'none'
                  }}
                >
                  <ScrollIndicator 
                    targetId="skills" 
                    offset={isMobile ? 60 : 80}
                    showAboveFold={true} 
                  />
                </div>
              </section>
              
              <section id="skills" className="section content-section">
                <SkillsSection />
              </section>
              
              <section id="projects" className="section content-section">
                <ProjectsSection themeMode={themeMode} />
              </section>
              
              <section id="about" className="section content-section">
                <About />
              </section>
            </>
          } />
          <Route path="/projects" element={<ProjectsSection themeMode={themeMode} />} />
          <Route path="/about" element={<About />} />
          <Route path="/skills" element={<SkillsSection />} />
        </Routes>
      </AnimatePresence>
      
      <BurgerMenu 
        activeSection={activeSection} 
        toggleTheme={toggleTheme}
        themeMode={themeMode} 
        onNavigate={scrollToSection}
      />
    </motion.div>
  );
};

AppContent.displayName = 'AppContent';

// Progress indicator component
const ProgressIndicator = () => {
  const { scrollYProgress } = useScroll();
  
  return (
    <motion.div 
      className="progress-bar"
      style={{ 
        scaleX: scrollYProgress,
        position: 'fixed',
        top: 0,
        left: 0,
        right: 0,
        height: '3px',
        background: 'linear-gradient(90deg, #6c63ff, #ff6b6b)',
        transformOrigin: '0%',
        zIndex: 1000
      }} 
    />
  );
};

// Custom scroll transform sections
const getScrollSections = (): ScrollSection[] => [
  {
    id: 'hero-reveal',
    startPercent: 0,
    endPercent: 0.15,
    transforms: [
      {
        type: 'opacity',
        target: '#hero .hero-title',
        from: '0',
        to: '1',
        easing: 'easeOut'
      },
      {
        type: 'translateY',
        target: '#hero .hero-subtitle',
        from: '30',
        to: '0',
        unit: 'px',
        easing: 'easeOut'
      }
    ]
  },
  {
    id: 'skills-animation',
    startPercent: 0.1,
    endPercent: 0.3,
    transforms: [
      {
        type: 'opacity',
        target: '#skills .skill-item',
        from: '0',
        to: '1'
      },
      {
        type: 'translateX',
        target: '#skills .skill-item:nth-child(odd)',
        from: '-50',
        to: '0',
        unit: 'px'
      },
      {
        type: 'translateX',
        target: '#skills .skill-item:nth-child(even)',
        from: '50',
        to: '0',
        unit: 'px'
      }
    ]
  },
  {
    id: 'projects-reveal',
    startPercent: 0.25,
    endPercent: 0.5,
    transforms: [
      {
        type: 'opacity',
        target: '#projects .section-title',
        from: '0',
        to: '1'
      },
      {
        type: 'translateY',
        target: '#projects .project-card',
        from: '50',
        to: '0',
        unit: 'px'
      }
    ]
  }
];

// Asset preloader function
const preloadAssets = (assets: string[]): Promise<void[]> => {
  return Promise.all(
    assets.map(src => {
      return new Promise<void>((resolve, reject) => {
        if (src.match(/\.(jpe?g|png|gif|svg|webp)$/i)) {
          const img = new Image();
          img.src = src;
          img.onload = () => resolve();
          img.onerror = () => {
            console.warn(`Failed to preload image: ${src}`);
            resolve(); // Resolve anyway to not block loading
          };
        } else {
          resolve();
        }
      });
    })
  );
};

// Main content with loading screen
const MainContent = () => {
  const { theme, themeMode, toggleTheme } = useTheme();
  const { zIndex } = useZIndex();
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [mounted, setMounted] = useState<boolean>(false);
  const [showDebug, setShowDebug] = useState<boolean>(false);
  const [initialLoadComplete, setInitialLoadComplete] = useState<boolean>(false);
  const scrollSections = getScrollSections();

  // Debug panel keyboard shortcut
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.altKey && e.key === 'd') {
        setShowDebug(prev => !prev);
      }
    };
    
    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, []);

  // Initial loading sequence
  useEffect(() => {
    setMounted(true);
    
    // Critical assets to preload
    const criticalAssets: string[] = [
      // Add paths to your critical images here
    ];
    
    // Minimum loading time combined with asset preloading
    const minLoadingTime = new Promise<void>(resolve => {
      const timeoutId = setTimeout(() => resolve(), 1500);
      return () => clearTimeout(timeoutId);
    });
    
    Promise.all([
      preloadAssets(criticalAssets),
      minLoadingTime
    ])
      .then(() => {
        setIsLoading(false);
        const timeoutId = setTimeout(() => setInitialLoadComplete(true), 1000);
        return () => clearTimeout(timeoutId);
      })
      .catch(err => {
        console.error('Failed during loading sequence:', err);
        setIsLoading(false);
        const timeoutId = setTimeout(() => setInitialLoadComplete(true), 1000);
        return () => clearTimeout(timeoutId);
      });
  }, []);

  if (!mounted) return null;

  return (
    <Layout>
      <GlobalStyles theme={theme} />
      
      <div className="app-container" style={{ position: 'relative' }}>
        <ScrollTransformBackground 
          showDebug={showDebug} 
          customSections={scrollSections}
          enableFloatingOrbs={true}
        />
        
        {initialLoadComplete && <ProgressIndicator />}
        
        <LoadingScreen 
          isLoading={isLoading}
          template="profile"
          fullScreen={true}
          backgroundColor={theme.colors.background}
        >
          <Suspense fallback={<div>Loading component...</div>}>
            <AppContent 
              themeMode={themeMode}
              toggleTheme={toggleTheme}
            />
          </Suspense>
        </LoadingScreen>
      </div>
    </Layout>
  );
};

// Enhanced Error boundary with better UX
class ErrorBoundary extends React.Component<
  { children: React.ReactNode },
  { hasError: boolean; errorInfo: string }
> {
  constructor(props: { children: React.ReactNode }) {
    super(props);
    this.state = { hasError: false, errorInfo: '' };
  }

  static getDerivedStateFromError(error: Error) {
    return { hasError: true };
  }

  componentDidCatch(error: Error, errorInfo: React.ErrorInfo) {
    console.error("Error caught by boundary:", error, errorInfo);
    this.setState({ errorInfo: errorInfo.componentStack || 'Unknown component error' });
  }

  render() {
    if (this.state.hasError) {
      return (
        <div className="error-container" style={{
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
          justifyContent: 'center',
          height: '100vh',
          padding: '20px',
          background: '#121212',
          color: '#fff'
        }}>
          <h2>Something went wrong</h2>
          <p>An unexpected error occurred while rendering the application.</p>
          <button 
            onClick={() => window.location.reload()}
            style={{
              background: 'linear-gradient(135deg, #6c63ff, #ff6b6b)',
              border: 'none',
              padding: '10px 20px',
              borderRadius: '4px',
              color: 'white',
              fontSize: '16px',
              cursor: 'pointer',
              marginTop: '20px'
            }}
          >
            Refresh Page
          </button>
          {process.env.NODE_ENV === 'development' && (
            <details style={{ marginTop: '20px', maxWidth: '800px', overflow: 'auto' }}>
              <summary>Error Details</summary>
              <pre style={{ whiteSpace: 'pre-wrap', textAlign: 'left' }}>
                {this.state.errorInfo}
              </pre>
            </details>
          )}
        </div>
      );
    }

    return this.props.children;
  }
}

// Root App component with providers
const App = () => {
  const { theme, themeMode, toggleTheme } = useThemeToggle();

  return (
    <ErrorBoundary>
      <ThemeProvider defaultTheme="dark">
        <GlobalStyles theme={theme} />
        <Router>
          <NavigationBar themeMode={themeMode} toggleTheme={toggleTheme} />
          <Routes>
            <Route path="/" element={<Home />} />
            <Route path="/about" element={<AboutPage />} />
            <Route path="/portfolio" element={<Home />} />
            <Route path="/contact" element={<Contact />} />
            <Route path="/devpanel" element={<DevPanel />} />
            <Route path="/urlshortener" element={<UrlShortener />} />
            <Route path="/messaging" element={<Messaging />} />
            <Route path="*" element={<NotFound />} />
          </Routes>
          <Footer />
        </Router>
      </ThemeProvider>
    </ErrorBoundary>
  );
};

export default App;
