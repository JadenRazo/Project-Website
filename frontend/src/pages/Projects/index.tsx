import React, { useState, useRef, useEffect, useCallback, useMemo } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import {
  FaGithub as FaGithubIcon,
  FaExternalLinkAlt as FaExternalLinkAltIcon,
  FaLaptopCode as FaLaptopCodeIcon,
  FaCode as FaCodeIcon,
} from 'react-icons/fa';
import styled from 'styled-components';
import { mockProjects } from '../../data/projects';

// --- Types ---
interface Project {
  readonly id: string;
  readonly title: string;
  readonly description: string;
  readonly technologies: readonly string[];
  readonly mediaUrl?: string;
  readonly mediaType?: 'image' | 'gif' | 'video';
  readonly githubUrl: string;
  readonly liveUrl: string;
}

// --- Interface for backend data ---
interface BackendProject {
  readonly id: string;
  readonly name: string;
  readonly description: string;
  readonly repo_url: string;
  readonly live_url?: string;
  readonly tags: readonly string[];
  readonly status: string;
}

// --- Component Props ---
interface ProjectCardProps {
  readonly project: Project;
  readonly expanded: boolean;
  readonly onClick: () => void;
  readonly tabIndex: number;
  readonly index: number;
}

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
  
  @media (max-width: 1200px) {
    max-width: 100%;
    padding: calc(1.75rem + ${HEADER_HEIGHT}px) 1.5rem 3.5rem 1.5rem;
  }
  
  @media (max-width: 900px) {
    padding: calc(1.5rem + ${HEADER_HEIGHT}px) 1rem 3rem 1rem;
  }
  
  @media (max-width: 600px) {
    padding: calc(1rem + ${HEADER_HEIGHT}px) 0.75rem 2rem 0.75rem;
  }
  
  @media (max-width: 480px) {
    padding: calc(0.75rem + ${HEADER_HEIGHT}px) 0.5rem 1.5rem 0.5rem;
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
  text-align: center;
  
  @media (max-width: 768px) {
    font-size: 2.5rem;
  }
  
  @media (max-width: 480px) {
    font-size: 2rem;
  }
  
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
    
    @media (max-width: 480px) {
      width: 60px;
      height: 3px;
    }
  }
`;

const PageDescription = styled.p`
  max-width: 650px;
  margin: 1.5rem auto 0;
  color: ${({ theme }) => theme?.colors?.textSecondary || '#aaa'};
  font-size: 1.1rem;
  line-height: 1.6;
  margin-bottom: 2rem;
  text-align: center;
  
  @media (max-width: 768px) {
    font-size: 1rem;
    max-width: 90%;
    margin-bottom: 1.5rem;
  }
  
  @media (max-width: 480px) {
    font-size: 0.9rem;
    line-height: 1.5;
    margin-bottom: 1rem;
  }
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
  max-width: 100%;
  text-align: center;
  flex-wrap: wrap;

  @media (max-width: 768px) {
    padding: 0.875rem 1.25rem;
    margin-bottom: 2.5rem;
    font-size: 0.95rem;
    border-radius: 10px;
    gap: 0.6rem;
  }
  
  @media (max-width: 480px) {
    padding: 0.75rem 1rem;
    margin-bottom: 2rem;
    font-size: 0.9rem;
    border-radius: 8px;
    gap: 0.5rem;
    flex-direction: column;
  }

  svg {
    color: ${({ theme }) => theme?.colors?.primary || '#007bff'};
    font-size: 1.5rem;
    
    @media (max-width: 768px) {
      font-size: 1.35rem;
    }
    
    @media (max-width: 480px) {
      font-size: 1.25rem;
    }
  }

  strong {
    color: ${({ theme }) => theme?.colors?.text || '#fff'};
    font-weight: 600;
    font-size: 1.1rem;
    
    @media (max-width: 768px) {
      font-size: 1.05rem;
    }
    
    @media (max-width: 480px) {
      font-size: 1rem;
    }
  }
`;

const ProjectsGrid = styled.div`
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
  gap: 2rem;
  width: 100%;
  box-sizing: border-box;
  overflow: visible;
  justify-items: center;
  align-items: start;
  
  @media (max-width: 1200px) {
    grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
    gap: 1.75rem;
  }
  
  @media (max-width: 900px) {
    grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
    gap: 1.5rem;
  }
  
  @media (max-width: 768px) {
    grid-template-columns: 1fr;
    gap: 1.25rem;
    max-width: 500px;
    margin: 0 auto;
  }
  
  @media (max-width: 480px) {
    gap: 1rem;
    max-width: 100%;
    padding: 0 0.5rem;
  }
`;

// --- Project Card ---
const Card = styled(motion.div)<{ $expanded: boolean }>`
  background: ${({ theme }) => theme?.colors?.surface || 'rgba(255,255,255,0.08)'};
  border-radius: 16px;
  box-shadow: 0 4px 16px rgba(0,0,0,0.12);
  border: 1px solid ${({ theme }) => theme?.colors?.border || 'rgba(255,255,255,0.1)'};
  overflow: hidden;
  transition: box-shadow 0.3s, border 0.3s;
  outline: none;
  position: relative;
  display: flex;
  flex-direction: column;
  box-sizing: border-box;
  width: 100%;
  height: auto;
  z-index: ${({ $expanded }) => $expanded ? '999' : '1'};
  contain: layout style;
  
  ${({ $expanded }) => $expanded && `
    overflow: visible;
    box-shadow: 0 20px 40px rgba(0,0,0,0.3);
    border-color: var(--primary-color, #007bff);
  `}
  
  
  @media (max-width: 768px) {
    border-radius: 12px;
  }
  
  @media (max-width: 480px) {
    border-radius: 8px;
  }
  
  &:hover, &:focus {
    box-shadow: 0 8px 32px rgba(0,0,0,0.18);
    border-color: ${({ theme }) => theme?.colors?.primary || '#007bff'};
    transform: ${({ $expanded }) => $expanded ? 'none' : 'translateY(-5px)'};
    
    @media (max-width: 768px) {
      transform: ${({ $expanded }) => $expanded ? 'none' : 'translateY(-3px)'};
    }
    
    @media (max-width: 480px) {
      transform: none;
    }
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
    
    @media (max-width: 480px) {
      height: 3px;
    }
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
  padding: 1.5rem;
  width: 100%;
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 1rem;
  min-height: 0;
  
  @media (max-width: 768px) {
    padding: 1.25rem;
    gap: 0.875rem;
  }
  
  @media (max-width: 480px) {
    padding: 1rem;
    gap: 0.75rem;
  }
`;

const CardTitle = styled.h2`
  font-size: 1.25rem;
  margin: 0;
  color: ${({ theme }) => theme?.colors?.text || '#fff'};
  font-family: ${({ theme }) => theme?.fonts?.primary || 'system-ui'};
  display: flex;
  align-items: center;
  gap: 0.5rem;
  line-height: 1.3;
  font-weight: 600;
  
  @media (max-width: 768px) {
    font-size: 1.15rem;
  }
  
  @media (max-width: 480px) {
    font-size: 1.1rem;
    gap: 0.4rem;
  }
  
  svg {
    color: ${({ theme }) => theme?.colors?.primary || '#007bff'};
    flex-shrink: 0;
    
    @media (max-width: 480px) {
      font-size: 1rem;
    }
  }
`;

const CardDescription = styled.p`
  font-size: 0.9rem;
  margin: 0;
  color: ${({ theme }) => theme?.colors?.textSecondary || '#aaa'};
  line-height: 1.5;
  flex: 1;
  
  @media (max-width: 768px) {
    font-size: 0.85rem;
    line-height: 1.4;
  }
  
  @media (max-width: 480px) {
    font-size: 0.8rem;
    line-height: 1.4;
  }
`;

const TechStack = styled.div`
  display: flex;
  flex-wrap: wrap;
  gap: 0.4rem;
  margin: 0;
  
  @media (max-width: 768px) {
    gap: 0.35rem;
  }
  
  @media (max-width: 480px) {
    gap: 0.3rem;
  }
`;

const TechTag = styled.span`
  background: ${({ theme }) => theme?.colors?.primaryLight || 'rgba(0,123,255,0.1)'};
  color: ${({ theme }) => theme?.colors?.primary || '#007bff'};
  padding: 0.2rem 0.6rem;
  border-radius: 12px;
  font-size: 0.75rem;
  font-family: ${({ theme }) => theme?.fonts?.mono || 'monospace'};
  transition: transform 0.2s, background 0.2s;
  white-space: nowrap;
  font-weight: 500;
  
  @media (max-width: 768px) {
    font-size: 0.7rem;
    padding: 0.15rem 0.5rem;
    border-radius: 10px;
  }
  
  @media (max-width: 480px) {
    font-size: 0.65rem;
    padding: 0.125rem 0.4rem;
    border-radius: 8px;
  }
  
  &:hover {
    transform: translateY(-1px);
    background: ${({ theme }) => theme?.colors?.primary || '#007bff'};
    color: white;
    
    @media (max-width: 480px) {
      transform: none;
    }
  }
`;

const CardLinks = styled.div`
  display: flex;
  gap: 0.75rem;
  margin-top: auto;
  
  @media (max-width: 768px) {
    gap: 0.6rem;
  }
  
  @media (max-width: 480px) {
    gap: 0.5rem;
    flex-direction: row;
  }
  
  @media (max-width: 360px) {
    flex-direction: column;
    gap: 0.4rem;
  }
`;

const CardLink = styled.a`
  display: inline-flex;
  align-items: center;
  gap: 0.4rem;
  color: ${({ theme }) => theme?.colors?.primary || '#007bff'};
  text-decoration: none;
  font-family: ${({ theme }) => theme?.fonts?.mono || 'monospace'};
  font-size: 0.8rem;
  padding: 0.6rem 0.8rem;
  border-radius: 6px;
  background: ${({ theme }) => `${theme?.colors?.primary}15` || 'rgba(0,123,255,0.1)'};
  transition: background 0.2s, transform 0.2s;
  white-space: nowrap;
  flex: 1;
  text-align: center;
  justify-content: center;
  min-height: 36px;
  font-weight: 500;
  
  @media (max-width: 768px) {
    font-size: 0.75rem;
    padding: 0.5rem 0.6rem;
    min-height: 34px;
    gap: 0.3rem;
  }
  
  @media (max-width: 480px) {
    font-size: 0.7rem;
    padding: 0.4rem 0.5rem;
    min-height: 32px;
    gap: 0.25rem;
  }
  
  @media (max-width: 360px) {
    font-size: 0.65rem;
    padding: 0.35rem 0.4rem;
    min-height: 30px;
    gap: 0.2rem;
  }
  
  &:hover {
    background: ${({ theme }) => `${theme?.colors?.primary}25` || 'rgba(0,123,255,0.15)'};
    transform: translateY(-1px);
    
    @media (max-width: 480px) {
      transform: none;
      background: ${({ theme }) => `${theme?.colors?.primary}20` || 'rgba(0,123,255,0.12)'};
    }
  }
  
  svg {
    font-size: 0.9rem;
    flex-shrink: 0;
    
    @media (max-width: 480px) {
      font-size: 0.8rem;
    }
    
    @media (max-width: 360px) {
      font-size: 0.75rem;
    }
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
  contain: layout style;
  
  /* Fixed height to prevent grid disruption */
  height: 300px;
  
  @media (max-width: 768px) {
    height: 250px;
  }
  
  @media (max-width: 480px) {
    height: 200px;
  }
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

const MobileCloseButton = styled.button`
  position: absolute;
  top: 0.75rem;
  right: 0.75rem;
  background: rgba(0, 0, 0, 0.7);
  border: none;
  border-radius: 50%;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  cursor: pointer;
  z-index: 10;
  transition: background 0.2s;
  
  @media (min-width: 769px) {
    display: none;
  }
  
  &:hover {
    background: rgba(0, 0, 0, 0.9);
  }
  
  svg {
    font-size: 0.875rem;
  }
`;

// --- Project Card Animations ---
const cardVariants = {
  hidden: { opacity: 0, y: 20 },
  visible: (i: number) => ({
    opacity: 1,
    y: 0,
    transition: {
      delay: i * 0.1,
      duration: window.innerWidth <= 768 ? 0.4 : 0.5,
      ease: "easeOut"
    }
  })
};

// Expandable Card component
const ProjectCard: React.FC<ProjectCardProps> = React.memo(({ 
  project, 
  expanded, 
  onClick, 
  tabIndex, 
  index 
}) => {
  // Ref for outside click
  const cardRef = useRef<HTMLDivElement>(null);
  const [isMobile, setIsMobile] = useState(false);

  // Check if mobile on mount and resize
  useEffect(() => {
    const checkMobile = () => {
      setIsMobile(window.innerWidth <= 768);
    };
    
    checkMobile();
    window.addEventListener('resize', checkMobile);
    return () => window.removeEventListener('resize', checkMobile);
  }, []);

  // Close on outside click (desktop only)
  useEffect(() => {
    if (!expanded || isMobile) return;
    
    const cardElement = cardRef.current;
    const handleClick = (e: MouseEvent) => {
      if (cardElement && !cardElement.contains(e.target as Node)) {
        onClick();
      }
    };
    
    document.addEventListener('mousedown', handleClick);
    return () => document.removeEventListener('mousedown', handleClick);
  }, [expanded, onClick, isMobile]);

  return (
    <Card
      ref={cardRef}
      $expanded={expanded}
      aria-expanded={expanded}
      initial="hidden"
      variants={cardVariants}
      whileHover={{ scale: expanded ? 1 : 1.02 }}
      animate={{
        opacity: 1,
        y: 0, 
        boxShadow: expanded ? '0 8px 32px rgba(0,0,0,0.25)' : '0 4px 16px rgba(0,0,0,0.12)',
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
            initial={{ opacity: 0, scaleY: 0 }}
            animate={{ 
              opacity: 1, 
              scaleY: 1
            }}
            exit={{ opacity: 0, scaleY: 0 }}
            transition={{ 
              type: 'spring', 
              stiffness: 300, 
              damping: 30,
              duration: 0.4
            }}
            style={{ transformOrigin: 'top' }}
            onClick={e => e.stopPropagation()} // Prevent bubbling to header
          >
            {isMobile && (
              <MobileCloseButton
                onClick={(e) => {
                  e.stopPropagation();
                  onClick();
                }}
                aria-label="Close media"
              >
                ×
              </MobileCloseButton>
            )}
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
                    controls={window.innerWidth <= 768} // Show controls on mobile/tablet
                    preload="metadata"
                    style={{ 
                      width: '100%', 
                      height: '100%', 
                      objectFit: 'contain', // Changed from 'cover' to 'contain' to show full video
                      borderRadius: 0,
                      cursor: window.innerWidth <= 768 ? 'auto' : 'pointer',
                      backgroundColor: 'rgba(0, 0, 0, 0.05)' // Subtle background for letterboxing
                    }}
                    aria-label={project.title + ' demo video'}
                    onError={(e) => console.error('Video load error:', e)}
                    onLoadStart={() => console.log('Video loading started for:', project.title)}
                    onCanPlay={() => console.log('Video can play:', project.title)}
                    onClick={(e) => {
                      // On mobile, let video controls handle interaction
                      if (window.innerWidth <= 768) {
                        e.stopPropagation();
                      }
                    }}
                  />
                ) : (
                  <img
                    src={project.mediaUrl}
                    alt={project.title + ' demo'}
                    style={{ 
                      width: '100%', 
                      height: '100%', 
                      objectFit: 'cover', 
                      borderRadius: 0,
                      cursor: 'pointer'
                    }}
                    loading="lazy"
                    onError={(e) => console.error('Image load error:', e)}
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
}, (prevProps: ProjectCardProps, nextProps: ProjectCardProps): boolean => {
  // Custom comparison function for React.memo
  return (
    prevProps.project.id === nextProps.project.id &&
    prevProps.expanded === nextProps.expanded &&
    prevProps.tabIndex === nextProps.tabIndex &&
    prevProps.index === nextProps.index &&
    prevProps.onClick === nextProps.onClick
  );
});

ProjectCard.displayName = 'ProjectCard';

// --- Main Projects Page ---
const Projects: React.FC = () => {
  // Only one card can be open at a time
  const [expandedId, setExpandedId] = useState<string | null>(null);
  const [totalLinesOfCode, setTotalLinesOfCode] = useState<number | null>(null);
  const [isLoadingLines, setIsLoadingLines] = useState<boolean>(true);
  const [errorLines, setErrorLines] = useState<string | null>(null);
  
  // Projects state
  const [projects, setProjects] = useState<Project[]>([]);
  const [projectsLoading, setProjectsLoading] = useState<boolean>(true);
  const [projectsError, setProjectsError] = useState<string | null>(null);


  // Handler for card click
  const handleCardClick = useCallback((id: string) => {
    if (expandedId === id) {
      setExpandedId(null);
    } else {
      setExpandedId(id);
    }
  }, [expandedId]);

  // Memoized click handlers for each project to prevent unnecessary re-renders
  const clickHandlers = useMemo(() => {
    const handlers = new Map<string, () => void>();
    projects.forEach(project => {
      handlers.set(project.id, () => handleCardClick(project.id));
    });
    return handlers;
  }, [projects, handleCardClick]);
  
  // Close expanded card when clicking outside of any card (desktop only)
  useEffect(() => {
    if (!expandedId || window.innerWidth <= 768) return;
    
    const handleOutsideClick = (e: MouseEvent) => {
      const projectsGrid = document.querySelector('[data-projects-grid]');
      if (projectsGrid && !projectsGrid.contains(e.target as Node)) {
        setExpandedId(null);
      }
    };
    
    document.addEventListener('mousedown', handleOutsideClick);
    return () => document.removeEventListener('mousedown', handleOutsideClick);
  }, [expandedId]);

  // Fetch projects from API
  useEffect(() => {
    const fetchProjects = async (): Promise<void> => {
      setProjectsLoading(true);
      setProjectsError(null);
      
      try {
        const apiUrl = (window as any)._env_?.REACT_APP_API_URL || process.env.REACT_APP_API_URL || 'http://localhost:8080';
        const response = await fetch(`${apiUrl}/api/v1/projects?status=active`);
        
        if (!response.ok) {
          throw new Error(`Failed to fetch projects: ${response.status} ${response.statusText}`);
        }
        
        const data: { projects?: BackendProject[] } = await response.json();
        const backendProjects: BackendProject[] = data.projects ?? [];
        
        if (backendProjects.length === 0) {
          // Use mock data when API returns empty results
          const transformedMockProjects: Project[] = mockProjects.map((proj): Project => ({
            id: proj.id,
            title: proj.name,
            description: proj.description,
            technologies: proj.tags ?? [],
            mediaUrl: proj.mediaUrl ?? '',
            mediaType: (proj.mediaType as Project['mediaType']) ?? 'image',
            githubUrl: proj.repo_url.includes('github.com/JadenRazo') 
              ? proj.repo_url.replace('https://github.com/JadenRazo', '') 
              : proj.repo_url,
            liveUrl: proj.live_url ?? proj.repo_url
          }));
          
          setProjects(transformedMockProjects);
        } else {
          // Transform backend data and merge with mock data for media
          const transformedProjects: Project[] = backendProjects.map((proj: BackendProject): Project => {
            const mockProject = mockProjects.find(mock => mock.name === proj.name);
            return {
              id: proj.id,
              title: proj.name,
              description: proj.description,
              technologies: proj.tags ?? [],
              mediaUrl: mockProject?.mediaUrl ?? '',
              mediaType: (mockProject?.mediaType as Project['mediaType']) ?? 'image',
              githubUrl: proj.repo_url.includes('github.com/JadenRazo') 
                ? proj.repo_url.replace('https://github.com/JadenRazo', '') 
                : proj.repo_url,
              liveUrl: proj.live_url ?? proj.repo_url
            };
          });
          
          setProjects(transformedProjects);
        }
      } catch (err) {
        console.error("Error fetching projects:", err);
        const errorMessage = err instanceof Error ? err.message : "Failed to load projects";
        setProjectsError(errorMessage);
        
        // Fallback to mock data
        const transformedMockProjects: Project[] = mockProjects.map((proj): Project => ({
          id: proj.id,
          title: proj.name,
          description: proj.description,
          technologies: proj.tags ?? [],
          mediaUrl: proj.mediaUrl ?? '',
          mediaType: (proj.mediaType as Project['mediaType']) ?? 'image',
          githubUrl: proj.repo_url.includes('github.com/JadenRazo') 
            ? proj.repo_url.replace('https://github.com/JadenRazo', '') 
            : proj.repo_url,
          liveUrl: proj.live_url ?? proj.repo_url
        }));
        
        setProjects(transformedMockProjects);
        setProjectsError(null); // Clear error since we have fallback data
      } finally {
        setProjectsLoading(false);
      }
    };

    fetchProjects();
  }, []);

  useEffect(() => {
    const fetchLinesOfCode = async (): Promise<void> => {
      setIsLoadingLines(true);
      setErrorLines(null);
      
      try {
        // Try API endpoint first, fallback to static JSON
        let response: Response;
        try {
          // Use the API URL from environment or default to localhost
          const apiUrl = (window as any)._env_?.REACT_APP_API_URL || process.env.REACT_APP_API_URL || 'http://localhost:8080';
          response = await fetch(`${apiUrl}/api/v1/code/stats`);
        } catch {
          // Fallback to static JSON file
          response = await fetch('/code_stats.json');
        }
        
        if (!response.ok) {
          throw new Error(`Failed to fetch code stats: ${response.status} ${response.statusText}`);
        }
        
        const data: { totalLines?: number; error?: string } = await response.json();
        
        if (data.error) {
          throw new Error(`Error in code stats: ${data.error}`);
        }

        // API now returns totalLines directly
        const totalLines = data.totalLines;
        if (typeof totalLines === 'number') {
          setTotalLinesOfCode(totalLines);
        } else {
          setTotalLinesOfCode(null);
          console.warn("Warning: Total lines of code was not a number or was missing.");
        }
      } catch (err) {
        console.error("Error fetching or processing lines of code:", err);
        const errorMessage = err instanceof Error ? err.message : "An unknown error occurred while loading code stats.";
        setErrorLines(errorMessage);
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
      
      {projectsLoading ? (
        <CodeStatsDisplayContainer
          initial={{ opacity: 0, y: -10 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5, delay: 0.2 }}
        >
          <FaCodeIcon />
          <span>Loading projects...</span>
        </CodeStatsDisplayContainer>
      ) : projectsError ? (
        <CodeStatsDisplayContainer
          initial={{ opacity: 0, y: -10 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5, delay: 0.2 }}
        >
          <FaCodeIcon />
          <span>Error loading projects: {projectsError}</span>
        </CodeStatsDisplayContainer>
      ) : (
        <ProjectsGrid data-projects-grid>
          {projects.map((project, idx) => (
            <ProjectCard
              key={project.id}
              project={project}
              expanded={expandedId === project.id}
              onClick={clickHandlers.get(project.id)!}
              tabIndex={0}
              index={idx}
            />
          ))}
        </ProjectsGrid>
      )}
    </ProjectsContainer>
  );
};

export default Projects; 