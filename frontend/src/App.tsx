// src/App.tsx
import React, { useState, useEffect, useCallback, memo, useMemo, useRef } from 'react';
import { ThemeProvider, useTheme } from './contexts/ThemeContext';
import { motion, AnimatePresence, Variants } from 'framer-motion';
import { GlobalStyles } from './styles/GlobalStyles';
import { NavigationBar } from './components/layout/NavigationBar';
import { Hero } from './components/sections/Hero';
import { Projects } from './components/sections/Projects';
import { LoadingScreen } from './components/animations/LoadingScreen';
import { NetworkBackground } from './components/animations/NetworkBackground';
import { ScrollIndicator } from './components/animations/ScrollIndicator';
import { BurgerMenu } from './components/navigation/BurgerMenu';
import { Layout } from './components/layout/Layout';

// === Types ===
interface Project {
  id: string;
  title: string;
  description: string;
  image: string;
  link: string;
  language: string;
}

// === Animation Variants ===
const fadeVariants: Variants = {
  initial: {
    opacity: 0
  },
  animate: {
    opacity: 1,
    transition: {
      duration: 0.5,
      ease: 'easeOut'
    }
  },
  exit: {
    opacity: 0,
    transition: {
      duration: 0.3,
      ease: 'easeIn'
    }
  }
};

const slideUpVariants: Variants = {
  initial: {
    opacity: 0,
    y: 20
  },
  animate: {
    opacity: 1,
    y: 0,
    transition: {
      duration: 0.5,
      ease: [0.645, 0.045, 0.355, 1]
    }
  },
  exit: {
    opacity: 0,
    y: -20,
    transition: {
      duration: 0.3,
      ease: [0.645, 0.045, 0.355, 1]
    }
  }
};

// === Constants ===
const ANIMATION_VARIANTS = Object.freeze({
  pageTransition: fadeVariants,
  projectsSection: slideUpVariants
});

const PROJECTS_DATA: readonly Project[] = Object.freeze([
  {
    id: '1',
    title: "Project 1",
    description: "This is the first project.",
    image: "https://via.placeholder.com/300",
    link: "#",
    language: "JavaScript"
  },
  {
    id: '2',
    title: "Project 2",
    description: "This is the second project.",
    image: "https://via.placeholder.com/300",
    link: "#",
    language: "Python"
  },
  {
    id: '3',
    title: "Project 3",
    description: "This is the third project.",
    image: "https://via.placeholder.com/300",
    link: "#",
    language: "Java"
  },
  {
    id: '4',
    title: "Project 4",
    description: "This is the fourth project.",
    image: "https://via.placeholder.com/300",
    link: "#",
    language: "JavaScript"
  }
]);

const LANGUAGES: readonly string[] = Object.freeze(['All', 'JavaScript', 'Python', 'Java']);

// Completely isolated background component that never rerenders
// This is a critical optimization for animation-heavy components
const BackgroundLayer = memo(() => {
  // Using useRef to ensure this component never needs to rerender
  const mountedRef = useRef(true);
  
  return (
    <div className="background-layer" style={{ position: 'fixed', top: 0, left: 0, right: 0, bottom: 0, zIndex: -1 }}>
      <NetworkBackground />
    </div>
  );
});

BackgroundLayer.displayName = 'BackgroundLayer';

// Create a completely isolated project filter state container
// This prevents filter state changes from affecting other components
const ProjectsSection: React.FC<{
  themeMode: string;
}> = memo(({ themeMode }) => {
  // Filter state is contained entirely within this component
  const [selectedLanguage, setSelectedLanguage] = useState<string>('All');
  
  // Filter handler doesn't affect parent components
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

// Main UI content excluding background animations
const AppContent: React.FC<{
  themeMode: string;
  toggleTheme: () => void;
}> = memo(({ themeMode, toggleTheme }) => {
  return (
    <motion.div 
      className="app-content"
      variants={ANIMATION_VARIANTS.pageTransition}
      initial="initial"
      animate="animate"
      exit="exit"
    >
      <NavigationBar 
        isDarkMode={themeMode === 'dark'}
        toggleTheme={toggleTheme}
      />
      <Hero />
      <ScrollIndicator targetId="projects" showAboveFold={true} offset={80} />
      <ProjectsSection themeMode={themeMode} />
      <BurgerMenu />
    </motion.div>
  );
});

AppContent.displayName = 'AppContent';

/**
 * Main content component that handles the application state and UI
 */
const MainContent: React.FC = () => {
  // Get theme context
  const { theme, themeMode, toggleTheme } = useTheme();
  
  // State
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [mounted, setMounted] = useState<boolean>(false);

  // Effects
  useEffect(() => {
    setMounted(true);
    const timer = setTimeout(() => setIsLoading(false), 2000);
    return () => clearTimeout(timer);
  }, []);

  // No hooks after this point, so conditional returns are safe now
  if (!mounted) return null;

  return (
    <Layout>
      <GlobalStyles theme={theme} />
      {/* BackgroundLayer rendered once and never rerenders */}
      <BackgroundLayer />
      
      <LoadingScreen 
        isLoading={isLoading}
        template="profile"
        fullScreen={true}
      >
        <AnimatePresence mode="wait">
          <motion.div
            key="content"
            initial="initial"
            animate="animate"
            exit="exit"
            variants={ANIMATION_VARIANTS.pageTransition}
          >
            <AppContent 
              themeMode={themeMode}
              toggleTheme={toggleTheme}
            />
          </motion.div>
        </AnimatePresence>
      </LoadingScreen>
    </Layout>
  );
};

/**
 * Root App component that provides theme context
 */
const App: React.FC = () => {
  return (
    <ThemeProvider>
      <ErrorBoundary>
        <MainContent />
      </ErrorBoundary>
    </ThemeProvider>
  );
};

/**
 * Error boundary to catch rendering errors
 */
class ErrorBoundary extends React.Component<
  { children: React.ReactNode },
  { hasError: boolean }
> {
  constructor(props: { children: React.ReactNode }) {
    super(props);
    this.state = { hasError: false };
  }

  static getDerivedStateFromError(): { hasError: boolean } {
    return { hasError: true };
  }

  componentDidCatch(error: Error, errorInfo: React.ErrorInfo): void {
    console.error('Error caught by ErrorBoundary:', error, errorInfo);
  }

  render(): React.ReactNode {
    if (this.state.hasError) {
      return (
        <div className="error-boundary">
          <h1>Something went wrong.</h1>
          <p>Please refresh the page or try again later.</p>
        </div>
      );
    }

    return this.props.children;
  }
}

export default App;
