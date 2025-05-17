import React, { useState, useRef, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import {
  FaGithub as FaGithubIcon,
  FaExternalLinkAlt as FaExternalLinkAltIcon,
  FaLaptopCode as FaLaptopCodeIcon,
  FaCode as FaCodeIcon,
} from 'react-icons/fa';
import styled from 'styled-components';

// --- Types ---
interface Project {
  id: string;
  title: string;
  description: string;
  technologies: string[];
  mediaUrl?: string; // GIF or video
  mediaType?: 'image' | 'gif' | 'video';
  githubUrl: string;
  liveUrl: string;
}

// --- Sample Data ---
const projects: Project[] = [
  {
    id: 'portfolio',
    title: 'Portfolio Website',
    description: 'A modern, responsive portfolio website built with React, TypeScript, and styled-components.',
    technologies: ['React', 'TypeScript', 'Styled Components', 'Framer Motion'],
    mediaUrl: '',
    mediaType: 'gif',
    githubUrl: '/Project-Website',
    liveUrl: 'https://jadenrazo.dev'
  },
  {
    id: 'devpanel',
    title: 'DevPanel',
    description: 'A development environment management system with real-time monitoring and control.',
    technologies: ['React', 'Node.js', 'WebSocket', 'TypeScript'],
    mediaUrl: '',
    mediaType: 'video',
    githubUrl: '/Project-Website/backend/internal/devpanel',
    liveUrl: 'https://jadenrazo.dev/devpanel'
  },
  {
    id: 'messaging',
    title: 'Messaging Platform',
    description: 'A real-time messaging platform with WebSocket integration and modern UI.',
    technologies: ['React', 'WebSocket', 'Node.js', 'TypeScript'],
    mediaUrl: '',
    mediaType: 'gif',
    githubUrl: '/Project-Website/backend/internal/messaging',
    liveUrl: 'https://jadenrazo.dev/messaging'
  }
];

// --- Constants ---
const HEADER_HEIGHT = 60; // Adjust if your header is taller

const ProjectsContainer = styled.div`
  max-width: var(--page-max-width, 1200px);
  width: 100%;
  margin: 0 auto;
  padding: calc(2rem + ${HEADER_HEIGHT}px) 2rem 4rem 2rem;
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  background: ${({ theme }) => theme?.colors?.background || '#111'};
  color: ${({ theme }) => theme?.colors?.text || '#fff'};
  box-sizing: border-box;
  @media (max-width: 900px) {
    padding: calc(1.5rem + ${HEADER_HEIGHT}px) 1rem 3rem 1rem;
  }
  @media (max-width: 600px) {
    padding: calc(1rem + ${HEADER_HEIGHT}px) 0.75rem 2rem 0.75rem;
  }
`;

const PageHeader = styled.div`
  text-align: center;
  margin-bottom: 3rem;
  width: 100%;
`;

const PageTitle = styled.h1`
  font-size: 3rem;
  margin-bottom: 0.5rem;
  color: ${({ theme }) => theme?.colors?.primary || '#007bff'};
  font-weight: 700;
  position: relative;
  display: inline-block;
  
  &::after {
    content: '';
    position: absolute;
    bottom: -10px;
    left: 50%;
    transform: translateX(-50%);
    width: 80px;
    height: 4px;
    background: ${({ theme }) => theme?.colors?.primary || '#007bff'};
    border-radius: 2px;
  }
`;

const PageDescription = styled.p`
  max-width: 650px;
  margin: 1.5rem auto 0;
  color: ${({ theme }) => theme?.colors?.textSecondary || '#aaa'};
  font-size: 1.1rem;
  line-height: 1.6;
  margin-bottom: 2rem;
`;

const CodeStatsDisplayContainer = styled(motion.div)`
  background: ${({ theme }) => theme?.colors?.surface || 'rgba(255,255,255,0.05)'};
  padding: 1rem 1.5rem;
  border-radius: 12px;
  border: 1px solid ${({ theme }) => theme?.colors?.border || 'rgba(255,255,255,0.1)'};
  margin-bottom: 3rem;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.75rem;
  color: ${({ theme }) => theme?.colors?.textSecondary || '#aaa'};
  font-size: 1rem;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);

  svg {
    color: ${({ theme }) => theme?.colors?.primary || '#007bff'};
    font-size: 1.5rem;
  }

  strong {
    color: ${({ theme }) => theme?.colors?.text || '#fff'};
    font-weight: 600;
    font-size: 1.1rem;
  }
`;

const ProjectsGrid = styled.div`
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(340px, 1fr));
  gap: 2rem;
  width: 100%;
  box-sizing: border-box;
  overflow: visible;
  @media (max-width: 900px) {
    grid-template-columns: 1fr;
    gap: 1.5rem;
  }
`;

// --- Project Card ---
const Card = styled(motion.div)<{ $expanded: boolean }>`
  background: ${({ theme }) => theme?.colors?.surface || 'rgba(255,255,255,0.08)'};
  border-radius: 16px;
  box-shadow: 0 4px 16px rgba(0,0,0,0.12);
  border: 1px solid ${({ theme }) => theme?.colors?.border || 'rgba(255,255,255,0.1)'};
  overflow: hidden;
  transition: box-shadow 0.3s, border 0.3s, transform 0.3s;
  outline: none;
  position: relative;
  height: ${({ $expanded }) => $expanded ? 'auto' : '260px'};
  min-height: ${({ $expanded }) => $expanded ? '440px' : '260px'};
  max-height: ${({ $expanded }) => $expanded ? '600px' : '260px'};
  display: flex;
  flex-direction: column;
  justify-content: flex-start;
  align-items: stretch;
  box-sizing: border-box;
  width: 100%;
  ${({ $expanded }) => $expanded ? 'z-index: 2;' : 'z-index: 1;'}
  
  &:hover, &:focus {
    box-shadow: 0 8px 32px rgba(0,0,0,0.18);
    border-color: ${({ theme }) => theme?.colors?.primary || '#007bff'};
    transform: ${({ $expanded }) => $expanded ? 'none' : 'translateY(-5px)'};
  }
  
  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 4px;
    background: linear-gradient(
      90deg, 
      ${({ theme }) => theme?.colors?.primary || '#007bff'}, 
      ${({ theme }) => theme?.colors?.secondary || '#6c757d'}
    );
    opacity: 0;
    transition: opacity 0.3s;
  }
  
  &:hover::before {
    opacity: 1;
  }
`;

const CardHeader = styled.div`
  width: 100%;
  cursor: pointer;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  overflow: hidden;
`;

const CardContent = styled.div`
  padding: 1.75rem 1.5rem;
  width: 100%;
`;

const CardTitle = styled.h2`
  font-size: 1.5rem;
  margin-bottom: 0.75rem;
  color: ${({ theme }) => theme?.colors?.text || '#fff'};
  font-family: ${({ theme }) => theme?.fonts?.primary || 'system-ui'};
  display: flex;
  align-items: center;
  gap: 0.5rem;
  
  svg {
    color: ${({ theme }) => theme?.colors?.primary || '#007bff'};
  }
`;

const CardDescription = styled.p`
  font-size: 1rem;
  margin-bottom: 1.5rem;
  color: ${({ theme }) => theme?.colors?.textSecondary || '#aaa'};
  line-height: 1.6;
`;

const TechStack = styled.div`
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
  margin-bottom: 1.5rem;
`;

const TechTag = styled.span`
  background: ${({ theme }) => theme?.colors?.primaryLight || 'rgba(0,123,255,0.1)'};
  color: ${({ theme }) => theme?.colors?.primary || '#007bff'};
  padding: 0.25rem 0.75rem;
  border-radius: 15px;
  font-size: 0.875rem;
  font-family: ${({ theme }) => theme?.fonts?.mono || 'monospace'};
  transition: transform 0.2s, background 0.2s;
  
  &:hover {
    transform: translateY(-2px);
    background: ${({ theme }) => theme?.colors?.primary || '#007bff'};
    color: white;
  }
`;

const CardLinks = styled.div`
  display: flex;
  gap: 1rem;
  margin-top: 1rem;
`;

const CardLink = styled.a`
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  color: ${({ theme }) => theme?.colors?.primary || '#007bff'};
  text-decoration: none;
  font-family: ${({ theme }) => theme?.fonts?.mono || 'monospace'};
  font-size: 0.875rem;
  padding: 0.5rem 0.75rem;
  border-radius: 8px;
  background: ${({ theme }) => `${theme?.colors?.primary}15` || 'rgba(0,123,255,0.1)'};
  transition: background 0.2s, transform 0.2s;
  
  &:hover {
    background: ${({ theme }) => `${theme?.colors?.primary}25` || 'rgba(0,123,255,0.15)'};
    transform: translateY(-2px);
  }
  
  svg {
    font-size: 1rem;
  }
`;

const GithubExternalLink = styled(CardLink)`
  svg {
    width: 16px;
    height: 16px;
  }
`;

const MediaWrapper = styled(motion.div)`
  width: 100%;
  position: relative;
  overflow: hidden;
  border-top: 1px solid ${({ theme }) => theme?.colors?.border || 'rgba(255,255,255,0.1)'};
  background: ${({ theme }) => theme?.colors?.surface || 'rgba(255,255,255,0.08)'};
`;

const Placeholder = styled.div`
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: ${({ theme }) => theme?.colors?.surface || 'rgba(255,255,255,0.08)'};
  color: ${({ theme }) => theme?.colors?.textSecondary || 'rgba(255,255,255,0.5)'};
`;

// --- Project Card Animations ---
const cardVariants = {
  hidden: { opacity: 0, y: 20 },
  visible: (i: number) => ({
    opacity: 1,
    y: 0,
    transition: {
      delay: i * 0.1,
      duration: 0.5,
      ease: "easeOut"
    }
  })
};

// Expandable Card component
const ProjectCard: React.FC<{
  project: Project;
  expanded: boolean;
  onClick: () => void;
  tabIndex: number;
  index: number;
}> = ({ project, expanded, onClick, tabIndex, index }) => {
  // Ref for outside click
  const cardRef = useRef<HTMLDivElement>(null);

  // Close on outside click
  useEffect(() => {
    if (!expanded) return;
    const handleClick = (e: MouseEvent) => {
      if (cardRef.current && !cardRef.current.contains(e.target as Node)) {
        onClick();
      }
    };
    document.addEventListener('mousedown', handleClick);
    return () => document.removeEventListener('mousedown', handleClick);
  }, [expanded, onClick]);

  return (
    <Card
      ref={cardRef}
      $expanded={expanded}
      tabIndex={tabIndex}
      aria-expanded={expanded}
      role="button"
      initial="hidden"
      variants={cardVariants}
      whileHover={{ scale: expanded ? 1 : 1.02 }}
      animate={{
        opacity: 1,
        y: 0, 
        boxShadow: expanded ? '0 8px 32px rgba(0,0,0,0.25)' : '0 4px 16px rgba(0,0,0,0.12)',
        height: expanded ? 'auto' : '260px',
        minHeight: expanded ? '440px' : '260px',
        maxHeight: expanded ? '600px' : '260px',
      }}
      transition={{ 
        type: 'spring', 
        stiffness: 200, 
        damping: 25,
        delay: index * 0.1,
        duration: 0.5,
        ease: "easeOut"
      }}
    >
      <CardHeader
        onClick={onClick}
        onKeyDown={e => (e.key === 'Enter' || e.key === ' ') && onClick()}
        tabIndex={0}
        aria-label={expanded ? `Collapse ${project.title}` : `Expand ${project.title}`}
      >
        <CardContent>
          <CardTitle>
            <FaLaptopCodeIcon />
            {project.title}
          </CardTitle>
          <CardDescription>{project.description}</CardDescription>
          <TechStack>
            {project.technologies.map((tech) => (
              <TechTag key={tech}>{tech}</TechTag>
            ))}
          </TechStack>
          <CardLinks>
            <GithubExternalLink
              href={`https://github.com/JadenRazo${project.githubUrl}`}
              target="_blank"
              rel="noopener noreferrer"
              tabIndex={-1}
              onClick={(e) => e.stopPropagation()} // Prevent card expansion when clicking link
            >
              <FaGithubIcon />
              GitHub
            </GithubExternalLink>
            <CardLink 
              href={project.liveUrl} 
              target="_blank" 
              rel="noopener noreferrer" 
              tabIndex={-1}
              onClick={(e) => e.stopPropagation()} // Prevent card expansion when clicking link
            >
              <FaExternalLinkAltIcon />
              Live Demo
            </CardLink>
          </CardLinks>
        </CardContent>
      </CardHeader>
      <AnimatePresence>
        {expanded && (
          <MediaWrapper
            key={"media-" + project.id}
            initial={{ opacity: 0, height: 0 }}
            animate={{ opacity: 1, height: 220 }}
            exit={{ opacity: 0, height: 0 }}
            transition={{ 
              type: 'spring', 
              stiffness: 200, 
              damping: 25
            }}
            onClick={e => e.stopPropagation()} // Prevent bubbling to header
          >
            <motion.div
              key={"media-content-" + project.id}
              style={{ width: '100%', height: '100%' }}
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              transition={{ duration: 0.35 }}
            >
              {project.mediaUrl ? (
                project.mediaType === 'video' ? (
                  <video
                    src={project.mediaUrl}
                    autoPlay
                    loop
                    muted
                    playsInline
                    style={{ width: '100%', height: '100%', objectFit: 'cover', borderRadius: 0 }}
                    aria-label={project.title + ' demo video'}
                    poster=""
                  />
                ) : (
                  <img
                    src={project.mediaUrl}
                    alt={project.title + ' demo'}
                    style={{ width: '100%', height: '100%', objectFit: 'cover', borderRadius: 0 }}
                    loading="lazy"
                  />
                )
              ) : (
                <Placeholder>No media available</Placeholder>
              )}
            </motion.div>
          </MediaWrapper>
        )}
      </AnimatePresence>
    </Card>
  );
};

// --- Main Projects Page ---
const Projects: React.FC = () => {
  // Only one card can be open at a time
  const [expandedId, setExpandedId] = useState<string | null>(null);
  const [totalLinesOfCode, setTotalLinesOfCode] = useState<number | null>(null);
  const [isLoadingLines, setIsLoadingLines] = useState<boolean>(true);
  const [errorLines, setErrorLines] = useState<string | null>(null);

  // Handler for card click
  const handleCardClick = (id: string) => {
    setExpandedId(prev => (prev === id ? null : id));
  };
  
  // Close expanded card when clicking outside of any card
  useEffect(() => {
    if (!expandedId) return;
    
    const handleOutsideClick = (e: MouseEvent) => {
      const projectsGrid = document.querySelector('[data-projects-grid]');
      if (projectsGrid && !projectsGrid.contains(e.target as Node)) {
        setExpandedId(null);
      }
    };
    
    document.addEventListener('mousedown', handleOutsideClick);
    return () => document.removeEventListener('mousedown', handleOutsideClick);
  }, [expandedId]);

  useEffect(() => {
    const fetchLinesOfCode = async () => {
      setIsLoadingLines(true);
      setErrorLines(null);
      try {
        // Fetch the static JSON file generated by the shell script
        const response = await fetch('/code_stats.json'); // Assumes code_stats.json is in the public folder
        if (!response.ok) {
          throw new Error(`Failed to fetch code_stats.json: ${response.status} ${response.statusText}`);
        }
        const data = await response.json();
        
        if (data.error) {
          throw new Error(`Error in code_stats.json: ${data.error}`);
        }

        if (typeof data.totalLines === 'number') {
          setTotalLinesOfCode(data.totalLines);
        } else {
          // This case might happen if totalLines is null but no explicit error was set
          // or if the format is unexpectedly different.
          setTotalLinesOfCode(null); // Explicitly set to null if not a number
          console.warn("Warning: totalLines in code_stats.json was not a number or was missing.");
          // Optionally, set an error message if this is unexpected
          // setErrorLines("Invalid or missing lines of code data from statistics file.");
        }
      } catch (err) {
        console.error("Error fetching or processing lines of code:", err);
        setErrorLines(err instanceof Error ? err.message : "An unknown error occurred while loading code stats.");
        setTotalLinesOfCode(null); // Clear any previous data
      } finally {
        setIsLoadingLines(false);
      }
    };

    fetchLinesOfCode();
    // The script is responsible for updates; frontend fetches once on load.
  }, []);

  return (
    <ProjectsContainer>
      <PageHeader>
        <PageTitle>My Projects</PageTitle>
        <PageDescription>
          Here's a showcase of my recent work. Each project represents different skills and technologies
          I've mastered. Click on any project to learn more about it.
        </PageDescription>
      </PageHeader>
      
      <CodeStatsDisplayContainer
        initial={{ opacity: 0, y: -10 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5, delay: 0.2 }}
      >
        <FaCodeIcon />
        {isLoadingLines && <span>Loading project stats...</span>}
        {errorLines && <span>Error loading stats: {errorLines}</span>}
        {!isLoadingLines && !errorLines && totalLinesOfCode !== null && (
          <>
            <span>Total Lines of Code Across Projects:</span>
            <strong>{totalLinesOfCode.toLocaleString()}</strong>
          </>
        )}
        {!isLoadingLines && !errorLines && totalLinesOfCode === null && (
            <span>Lines of code data not available.</span>
        )}
      </CodeStatsDisplayContainer>
      
      <ProjectsGrid data-projects-grid>
        {projects.map((project, idx) => (
          <ProjectCard
            key={project.id}
            project={project}
            expanded={expandedId === project.id}
            onClick={() => handleCardClick(project.id)}
            tabIndex={0}
            index={idx}
          />
        ))}
      </ProjectsGrid>
    </ProjectsContainer>
  );
};

export default Projects; 