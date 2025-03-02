// src/App.tsx
import React, { useState, useEffect, useCallback, memo, useMemo } from 'react';
import { ThemeProvider, useTheme } from './contexts/ThemeContext';
import { motion, AnimatePresence, Variants } from 'framer-motion';
import { GlobalStyles } from './styles/GlobalStyles';
import { NavigationBar } from './components/layout/NavigationBar';
import { Hero } from './components/sections/Hero';
import { LoadingScreen } from './components/animations/LoadingScreen';
import { NetworkBackground } from './components/animations/NetworkBackground';
import { ScrollIndicator } from './components/animations/ScrollIndicator';
import { ProjectCard } from './components/ui/ProjectCard';
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

interface ProjectsSectionProps {
  languages: readonly string[];
  selectedLanguage: string;
  setSelectedLanguage: (lang: string) => void;
  filteredProjects: readonly Project[];
}

interface LanguageFilterProps {
  language: string;
  isSelected: boolean;
  onSelect: (language: string) => void;
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

const projectCardVariants: Variants = {
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
  }
};

// === Constants ===
const ANIMATION_VARIANTS = Object.freeze({
  pageTransition: fadeVariants,
  projectsSection: slideUpVariants,
  projectCard: projectCardVariants
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

// === Components ===
const LanguageFilter: React.FC<LanguageFilterProps> = memo(({
  language,
  isSelected,
  onSelect
}) => (
  <motion.button
    key={language}
    onClick={() => onSelect(language)}
    className={`filter-button ${isSelected ? 'active' : ''}`}
    style={{
      padding: '10px 20px',
      border: 'none',
      borderRadius: 'var(--border-radius)',
      backgroundColor: isSelected 
        ? 'var(--primary)' 
        : 'var(--primary-light)',
      color: isSelected 
        ? 'var(--surface-light)' 
        : 'var(--primary)',
      cursor: 'pointer',
      transition: 'var(--transition)',
      fontSize: '1rem',
      fontFamily: 'inherit'
    }}
    whileHover={{ scale: 1.05 }}
    whileTap={{ scale: 0.95 }}
  >
    {language}
  </motion.button>
));
LanguageFilter.displayName = 'LanguageFilter';

const ProjectsSection: React.FC<ProjectsSectionProps> = memo(({
  languages,
  selectedLanguage,
  setSelectedLanguage,
  filteredProjects
}) => (
  <motion.section
    className="projects-section container"
    {...ANIMATION_VARIANTS.projectsSection}
  >
    <motion.div
      className="language-filters"
      style={{
        display: 'flex',
        justifyContent: 'center',
        gap: '20px',
        flexWrap: 'wrap',
        padding: '20px',
        marginTop: '50px'
      }}
    >
      {languages.map((language) => (
        <LanguageFilter
          key={language}
          language={language}
          isSelected={selectedLanguage === language}
          onSelect={setSelectedLanguage}
        />
      ))}
    </motion.div>

    <motion.div
      className="projects-grid"
      style={{
        display: 'grid',
        gridTemplateColumns: 'repeat(auto-fit, minmax(300px, 1fr))',
        gap: '2rem',
        padding: '2rem',
        maxWidth: '1200px',
        margin: '0 auto'
      }}
      {...ANIMATION_VARIANTS.projectsSection}
    >
      {filteredProjects.map((project, index) => (
        <motion.div
          key={project.id}
          {...ANIMATION_VARIANTS.projectCard}
          transition={{ delay: index * 0.1 }}
        >
          <ProjectCard {...project} />
        </motion.div>
      ))}
    </motion.div>
  </motion.section>
));
ProjectsSection.displayName = 'ProjectsSection';

/**
 * Main content component that handles the application state and UI
 */
const MainContent: React.FC = () => {
  // Get theme context
  const { theme, themeMode, toggleTheme } = useTheme();
  
  // State
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [selectedLanguage, setSelectedLanguage] = useState<string>('All');
  const [mounted, setMounted] = useState<boolean>(false);

  // Effects
  useEffect(() => {
    setMounted(true);
    const timer = setTimeout(() => setIsLoading(false), 2000);
    return () => clearTimeout(timer);
  }, []);

  // Callbacks
  const handleLanguageSelect = useCallback((language: string) => {
    setSelectedLanguage(language);
  }, []);

  // Memoized values
  const filteredProjects = useMemo(() => {
    if (selectedLanguage === 'All') return PROJECTS_DATA;
    return PROJECTS_DATA.filter(project => project.language === selectedLanguage);
  }, [selectedLanguage]);

  // Early return if not mounted
  if (!mounted) return null;

  // App content with all components
  const appContent = (
    <>
      <NetworkBackground />
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
        <ProjectsSection
          languages={LANGUAGES}
          selectedLanguage={selectedLanguage}
          setSelectedLanguage={handleLanguageSelect}
          filteredProjects={filteredProjects}
        />
        <ScrollIndicator />
        <BurgerMenu />
      </motion.div>
    </>
  );

  return (
    <Layout>
      <GlobalStyles theme={theme} />
      {/* Using the new LoadingScreen component with proper props */}
      <LoadingScreen 
        isLoading={isLoading}
        template="profile"  // You can change this to 'card', 'article', or 'custom' based on your preference
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
            {appContent}
          </motion.div>
        </AnimatePresence>
      </LoadingScreen>
    </Layout>
  );
};

/**
 * Error Boundary component to catch and display errors gracefully
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
    console.error('Error:', error);
    console.error('Error Info:', errorInfo);
  }

  render(): React.ReactNode {
    if (this.state.hasError) {
      return (
        <div style={{ 
          padding: '20px', 
          textAlign: 'center',
          color: 'red',
          backgroundColor: '#ffebee',
          borderRadius: '4px',
          margin: '20px' 
        }}>
          <h2>Something went wrong</h2>
          <p>The application encountered an error. Please try refreshing the page.</p>
        </div>
      );
    }

    return this.props.children;
  }
}

/**
 * Root App component that sets up the theme provider and error boundary
 */
const App: React.FC = () => {
  // Using the ErrorBoundary to catch any errors in the app
  return (
    <ErrorBoundary>
      <ThemeProvider defaultTheme="dark">
        <MainContent />
      </ThemeProvider>
    </ErrorBoundary>
  );
};

export default App;
