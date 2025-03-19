import React, { useState, useEffect } from 'react';
import styled from 'styled-components';

// Types
interface Project {
  id: string;
  name: string;
  description: string;
  url: string;
  status: string;
}

// Styled Components
const PageContainer = styled.div`
  padding: 2rem;
  max-width: 1200px;
  margin: 0 auto;
`;

const Header = styled.header`
  margin-bottom: 2rem;
  border-bottom: 1px solid ${({ theme }) => theme.colors.border};
  padding-bottom: 1rem;
`;

const Title = styled.h1`
  font-size: 2.5rem;
  margin-bottom: 0.5rem;
  color: ${({ theme }) => theme.colors.primary};
`;

const Subtitle = styled.p`
  font-size: 1.1rem;
  color: ${({ theme }) => theme.colors.textSecondary};
`;

const ProjectsGrid = styled.div`
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 1.5rem;
  margin-top: 2rem;
`;

const ProjectCard = styled.div`
  background: ${({ theme }) => theme.colors.card};
  border-radius: 8px;
  padding: 1.5rem;
  transition: transform 0.3s ease, box-shadow 0.3s ease;
  border: 1px solid ${({ theme }) => theme.colors.border};
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.05);
  
  &:hover {
    transform: translateY(-5px);
    box-shadow: 0 10px 20px rgba(0, 0, 0, 0.08);
  }
`;

const ProjectTitle = styled.h3`
  font-size: 1.25rem;
  margin-bottom: 0.5rem;
  color: ${({ theme }) => theme.colors.text};
`;

const ProjectDescription = styled.p`
  color: ${({ theme }) => theme.colors.textSecondary};
  margin-bottom: 1rem;
  line-height: 1.5;
`;

const ProjectStatus = styled.span<{ status: string }>`
  display: inline-block;
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-size: 0.75rem;
  font-weight: 500;
  background: ${({ theme, status }) => 
    status === 'active' ? theme.colors.success + '20' :
    status === 'development' ? theme.colors.warning + '20' :
    theme.colors.error + '20'
  };
  color: ${({ theme, status }) => 
    status === 'active' ? theme.colors.success :
    status === 'development' ? theme.colors.warning :
    theme.colors.error
  };
`;

const ProjectLink = styled.a`
  display: inline-block;
  margin-top: 1rem;
  color: ${({ theme }) => theme.colors.primary};
  text-decoration: none;
  
  &:hover {
    text-decoration: underline;
  }
`;

const Button = styled.button`
  background: ${({ theme }) => theme.colors.primary};
  color: white;
  border: none;
  border-radius: 4px;
  padding: 0.75rem 1.5rem;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  transition: background 0.3s ease;
  
  &:hover {
    background: ${({ theme }) => theme.colors.primaryHover};
  }
`;

const ErrorMessage = styled.div`
  padding: 1rem;
  margin: 1rem 0;
  background: ${({ theme }) => theme.colors.error}20;
  color: ${({ theme }) => theme.colors.error};
  border-radius: 4px;
  border-left: 4px solid ${({ theme }) => theme.colors.error};
`;

// Main Component
const DevPanel: React.FC = () => {
  const [projects, setProjects] = useState<Project[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchProjects = async () => {
      try {
        setLoading(true);
        // Replace with your actual API endpoint
        const response = await fetch('/devpanel/api/projects');
        
        // If the API is not ready, use mock data
        if (!response.ok) {
          // Simulate API response with mock data
          setTimeout(() => {
            setProjects([
              {
                id: '1',
                name: 'Portfolio Website',
                description: 'Personal portfolio website showcasing skills and projects',
                url: 'https://jadenrazo.dev',
                status: 'active'
              },
              {
                id: '2',
                name: 'URL Shortener Service',
                description: 'Custom URL shortening service with analytics',
                url: '/url-shortener',
                status: 'active'
              },
              {
                id: '3',
                name: 'Messaging Platform',
                description: 'Real-time messaging platform with WebSocket support',
                url: '/messaging',
                status: 'development'
              }
            ]);
            setLoading(false);
          }, 1000);
          return;
        }
        
        const data = await response.json();
        setProjects(data);
      } catch (err) {
        console.error('Error fetching projects:', err);
        setError('Failed to load projects. Please try again later.');
        
        // Fallback to mock data
        setProjects([
          {
            id: '1',
            name: 'Portfolio Website',
            description: 'Personal portfolio website showcasing skills and projects',
            url: 'https://jadenrazo.dev',
            status: 'active'
          },
          {
            id: '2',
            name: 'URL Shortener Service',
            description: 'Custom URL shortening service with analytics',
            url: '/url-shortener',
            status: 'active'
          },
          {
            id: '3',
            name: 'Messaging Platform',
            description: 'Real-time messaging platform with WebSocket support',
            url: '/messaging',
            status: 'development'
          }
        ]);
      } finally {
        setLoading(false);
      }
    };

    fetchProjects();
  }, []);

  return (
    <PageContainer>
      <Header>
        <Title>Developer Panel</Title>
        <Subtitle>Manage your projects and website settings</Subtitle>
      </Header>

      {error && <ErrorMessage>{error}</ErrorMessage>}

      {loading ? (
        <p>Loading projects...</p>
      ) : (
        <>
          <Button>Create New Project</Button>
          
          <ProjectsGrid>
            {projects.map(project => (
              <ProjectCard key={project.id}>
                <ProjectTitle>{project.name}</ProjectTitle>
                <ProjectStatus status={project.status}>
                  {project.status.charAt(0).toUpperCase() + project.status.slice(1)}
                </ProjectStatus>
                <ProjectDescription>{project.description}</ProjectDescription>
                <ProjectLink href={project.url} target={project.url.startsWith('http') ? "_blank" : undefined}>
                  View Project
                </ProjectLink>
              </ProjectCard>
            ))}
          </ProjectsGrid>
        </>
      )}
    </PageContainer>
  );
};

export default DevPanel; 