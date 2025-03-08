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
import { LoadingScreen } from './components/animations/LoadingScreen';
import { ScrollTransformBackground } from './components/animations/ScrollTransformBackground';
import { ScrollIndicator } from './components/animations/ScrollIndicator';
import { BurgerMenu } from './components/navigation/BurgerMenu';
import { Layout } from './components/layout/Layout';
import { SkillsSection } from './components/sections/SkillsSection';
import { useTheme } from './contexts/ThemeContext';

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
const useSectionObserver = (sectionIds: string[]) => {
  const [activeSection, setActiveSection] = useState<string>(sectionIds[0]);
  
  useEffect(() => {
    const observerOptions = {
      root: null,
      rootMargin: '-10% 0px -90% 0px',
      threshold: 0
    };
    
    const handleIntersect = (entries: IntersectionObserverEntry[]) => {
      entries.forEach((entry) => {
        if (entry.isIntersecting) {
          setActiveSection(entry.target.id);
        }
      });
    };
    
    const observer = new IntersectionObserver(handleIntersect, observerOptions);
    
    sectionIds.forEach((id) => {
      const element = document.getElementById(id);
      if (element) observer.observe(element);
    });
    
    return () => observer.disconnect();
  }, [sectionIds]);
  
  return activeSection;
};

// Scroll controls hook for manual navigation
const useScrollControls = () => {
  const scrollToSection = useCallback((sectionId: string) => {
    const section = document.getElementById(sectionId);
    if (section) {
      section.scrollIntoView({ behavior: 'smooth' });
    }
  }, []);
  
  return { scrollToSection };
};

// Main App Content component with route handling
const AppContent = memo<AppContentProps>(({ themeMode, toggleTheme }) => {
  const location = useLocation();
  const { scrollY } = useScroll();
  const sections = ['hero', 'skills', 'projects', 'about'];
  const activeSection = useSectionObserver(sections);
  const { scrollToSection } = useScrollControls();
  
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
                <ScrollIndicator 
                  targetId="skills" 
                  offset={80} 
                  showAboveFold={true} 
                />
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
});

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
  return (
    <Router>
      <ThemeProvider>
        <ZIndexProvider>
          <ErrorBoundary>
            <MainContent />
          </ErrorBoundary>
        </ZIndexProvider>
      </ThemeProvider>
    </Router>
  );
};

export default App;
