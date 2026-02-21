import React, { useEffect, useRef, useState } from 'react';
import styled from 'styled-components';
import { motion } from 'framer-motion';
import { FaCode } from 'react-icons/fa';
import {
  SiReact,
  SiVuedotjs,
  SiPython,
  SiNodedotjs,
  SiDocker,
  SiTypescript,
  SiGo,
  SiPostgresql,
  SiRedis,
  SiNginx,
  SiTailwindcss,
  SiAstro
} from 'react-icons/si';
import { HiServer, HiCog, HiCloud, HiViewGrid, HiDatabase } from 'react-icons/hi';
import { TechCategory } from '../../data/projects';
import CICDPipelineAnimation from '../animations/CICDPipelineAnimation';

const Card = styled(motion.div)`
  position: relative;
  width: 100%;
  max-width: 420px;
  border-radius: 16px;
  overflow: visible;
  background: ${({ theme }) => theme.colors.surface};
  box-shadow: ${({ theme }) => theme.shadows.large};
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  flex-direction: column;
  border: 1px solid ${({ theme }) => theme.colors.border || 'rgba(255,255,255,0.1)'};
  transform: translateZ(0);
  backface-visibility: hidden;

  &::before {
    content: '';
    position: absolute;
    inset: -2px;
    border-radius: 18px;
    padding: 2px;
    background: linear-gradient(135deg, var(--color-neon) 0%, transparent 50%, var(--color-neon) 100%);
    -webkit-mask: linear-gradient(#fff 0 0) content-box, linear-gradient(#fff 0 0);
    -webkit-mask-composite: xor;
    mask-composite: exclude;
    opacity: 0;
    transition: opacity 0.4s ease;
    pointer-events: none;
  }

  &:hover {
    transform: translateY(-12px) rotateX(2deg);
    box-shadow: 0 0 40px var(--color-neon-glow), 0 16px 48px rgba(0,0,0,0.3);
    border-color: var(--color-neon);

    &::before {
      opacity: 1;
    }
  }

  @media (max-width: 768px) {
    max-width: 100%;
    border-radius: 12px;

    &::before {
      border-radius: 14px;
    }

    &:hover {
      transform: translateY(-6px);
    }
  }
`;

const MediaContainer = styled.div`
  width: 100%;
  height: 220px;
  overflow: hidden;
  position: relative;
  background: ${({ theme }) => theme.colors.background};
  border-top-left-radius: 16px;
  border-top-right-radius: 16px;

  @media (max-width: 768px) {
    height: 200px;
    border-top-left-radius: 12px;
    border-top-right-radius: 12px;
  }
`;

const BadgeContainer = styled.div`
  position: absolute;
  top: 12px;
  right: 12px;
  display: flex;
  gap: 0.5rem;
  z-index: 3;
  flex-wrap: wrap;
  justify-content: flex-end;
`;

const ProjectBadge = styled.span<{ $type: string }>`
  padding: 0.35rem 0.65rem;
  border-radius: 20px;
  font-size: 0.65rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  white-space: nowrap;
  backdrop-filter: blur(10px);

  ${({ $type, theme }) => {
    const styles: Record<string, string> = {
      live: `
        background: linear-gradient(135deg, rgba(16, 185, 129, 0.95) 0%, rgba(5, 150, 105, 0.95) 100%);
        color: white;
        box-shadow: 0 4px 12px rgba(16, 185, 129, 0.4);
      `,
      demo: `
        background: rgba(59, 130, 246, 0.2);
        color: #3b82f6;
        border: 1.5px solid #3b82f6;
        box-shadow: 0 4px 12px rgba(59, 130, 246, 0.3);
      `,
      client: `
        background: linear-gradient(135deg, rgba(245, 158, 11, 0.95) 0%, rgba(217, 119, 6, 0.95) 100%);
        color: white;
        box-shadow: 0 4px 12px rgba(245, 158, 11, 0.4);
      `,
      internal: `
        background: rgba(100, 116, 139, 0.2);
        color: ${theme.colors.textSecondary};
        border: 1px solid ${theme.colors.border};
      `
    };
    return styles[$type] || styles.internal;
  }}
`;

const ProjectImage = styled.img`
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: all 0.5s cubic-bezier(0.4, 0, 0.2, 1);
  filter: grayscale(20%);

  ${Card}:hover & {
    transform: scale(1.08);
    filter: grayscale(0%);
  }
`;

const ProjectVideo = styled.video`
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: all 0.5s cubic-bezier(0.4, 0, 0.2, 1);
  filter: grayscale(20%);

  ${Card}:hover & {
    transform: scale(1.08);
    filter: grayscale(0%);
  }
`;

const ContentSection = styled.div`
  padding: 1.75rem;
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
  flex: 1;

  @media (max-width: 768px) {
    padding: 1.5rem;
    gap: 1rem;
  }
`;

const ProjectTitle = styled.h3`
  font-size: 1.4rem;
  font-weight: 700;
  margin: 0;
  color: ${({ theme }) => theme.colors.text};
  line-height: 1.3;

  @media (max-width: 768px) {
    font-size: 1.25rem;
  }
`;

const ProjectDescription = styled.p`
  font-size: 0.95rem;
  line-height: 1.6;
  color: ${({ theme }) => theme.colors.textSecondary};
  margin: 0;

  @media (max-width: 768px) {
    font-size: 0.9rem;
    line-height: 1.5;
  }
`;

const TechStackSection = styled.div`
  display: flex;
  flex-direction: column;
  gap: 1rem;
  margin-top: 0.5rem;
`;

const TechStackHeader = styled.div`
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.85rem;
  font-weight: 600;
  color: ${({ theme }) => theme.colors.primary};
  text-transform: uppercase;
  letter-spacing: 0.5px;

  svg {
    font-size: 1rem;
  }
`;

const TechCategoriesGrid = styled.div`
  display: grid;
  gap: 0.875rem;
`;

const TechCategoryRow = styled.div`
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
`;

const CategoryLabel = styled.div`
  display: flex;
  align-items: center;
  gap: 0.4rem;
  font-size: 0.75rem;
  font-weight: 600;
  color: ${({ theme }) => theme.colors.textSecondary};
  text-transform: uppercase;
  letter-spacing: 0.3px;

  svg {
    font-size: 0.85rem;
    color: ${({ theme }) => theme.colors.primary};
  }
`;

const TechTagsContainer = styled.div`
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
`;

const TechTag = styled.span<{ $hasIcon?: boolean }>`
  display: inline-flex;
  align-items: center;
  gap: 0.35rem;
  background: ${({ theme }) => `${theme.colors.primary}15`};
  color: ${({ theme }) => theme.colors.primary};
  padding: ${({ $hasIcon }) => $hasIcon ? '0.4rem 0.75rem' : '0.4rem 0.85rem'};
  border-radius: 8px;
  font-size: 0.8rem;
  font-weight: 500;
  transition: all 0.2s ease;
  border: 1px solid ${({ theme }) => `${theme.colors.primary}25`};

  svg {
    font-size: 0.95rem;
  }

  &:hover {
    background: ${({ theme }) => theme.colors.primary};
    color: white;
    transform: translateY(-2px);
    box-shadow: 0 4px 8px rgba(0,0,0,0.15);
  }

  @media (max-width: 768px) {
    font-size: 0.75rem;
    padding: ${({ $hasIcon }) => $hasIcon ? '0.35rem 0.65rem' : '0.35rem 0.75rem'};

    svg {
      font-size: 0.85rem;
    }
  }
`;

const TechCountBadge = styled.div`
  display: inline-flex;
  align-items: center;
  gap: 0.4rem;
  background: ${({ theme }) => `${theme.colors.primary}10`};
  color: ${({ theme }) => theme.colors.primary};
  padding: 0.4rem 0.75rem;
  border-radius: 20px;
  font-size: 0.75rem;
  font-weight: 600;
  border: 1px solid ${({ theme }) => `${theme.colors.primary}20`};

  svg {
    font-size: 0.9rem;
  }
`;

const ProjectActions = styled.div`
  display: flex;
  gap: 0.75rem;
  margin-top: auto;
  padding-top: 1rem;
  border-top: 1px solid ${({ theme }) => theme.colors.border || 'rgba(255,255,255,0.1)'};

  @media (max-width: 768px) {
    gap: 0.6rem;
  }
`;

const ActionButton = styled.a`
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  padding: 0.875rem 1.25rem;
  border-radius: 8px;
  background: ${({ theme }) => theme.colors.primary};
  color: #f8fafc !important;
  text-decoration: none !important;
  font-size: 0.9rem;
  font-weight: 700;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  flex: 1;
  text-align: center;
  position: relative;
  overflow: hidden;

  span {
    position: relative;
    z-index: 1;
    color: #f8fafc !important;
  }

  &:hover {
    background: ${({ theme }) => theme.colors.primaryHover || theme.colors.primary};
    color: #f8fafc !important;
    transform: translateY(-4px);
    box-shadow: 0 8px 24px ${({ theme }) => `${theme.colors.primary}60`}, 0 4px 12px rgba(0,0,0,0.2);

    span {
      color: #f8fafc !important;
    }
  }

  &:visited,
  &:active,
  &:link {
    color: #f8fafc !important;

    span {
      color: #f8fafc !important;
    }
  }

  @media (max-width: 768px) {
    padding: 0.75rem 1rem;
    font-size: 0.85rem;

    &:hover {
      transform: translateY(-2px);
    }
  }
`;

const SecondaryButton = styled(ActionButton)`
  background: ${({ theme }) => `${theme.colors.primary}15`};
  color: ${({ theme }) => theme.colors.primary};
  border: 2px solid ${({ theme }) => theme.colors.primary};

  span {
    color: ${({ theme }) => theme.colors.primary};
  }

  &:hover {
    background: ${({ theme }) => `${theme.colors.primary}25`};
    transform: translateY(-4px);
    box-shadow: 0 8px 24px ${({ theme }) => `${theme.colors.primary}40`}, 0 4px 12px rgba(0,0,0,0.2);

    span {
      color: ${({ theme }) => theme.colors.primary};
    }
  }

  &:visited,
  &:active {
    color: ${({ theme }) => theme.colors.primary};

    span {
      color: ${({ theme }) => theme.colors.primary};
    }
  }

  @media (max-width: 768px) {
    &:hover {
      transform: translateY(-2px);
    }
  }
`;

interface EnhancedProjectCardProps {
  title: string;
  description: string;
  image: string;
  repoLink: string;
  liveLink: string;
  techCategories?: TechCategory;
  mediaType?: 'image' | 'video' | 'component';
  badges?: Array<'live' | 'demo' | 'client' | 'internal'>;
}

const componentMap: Record<string, React.FC> = {
  'cicd-pipeline': CICDPipelineAnimation,
};

const getTechIcon = (tech: string): JSX.Element | null => {
  const lowerTech = tech.toLowerCase();

  if (lowerTech.includes('react')) return <SiReact />;
  if (lowerTech.includes('vue')) return <SiVuedotjs />;
  if (lowerTech.includes('typescript')) return <SiTypescript />;
  if (lowerTech.includes('python')) return <SiPython />;
  if (lowerTech.includes('node')) return <SiNodedotjs />;
  if (lowerTech.includes('go')) return <SiGo />;
  if (lowerTech.includes('docker')) return <SiDocker />;
  if (lowerTech.includes('postgres')) return <SiPostgresql />;
  if (lowerTech.includes('redis')) return <SiRedis />;
  if (lowerTech.includes('nginx')) return <SiNginx />;
  if (lowerTech.includes('tailwind')) return <SiTailwindcss />;
  if (lowerTech.includes('astro')) return <SiAstro />;

  return null;
};

const getCategoryIcon = (category: string): JSX.Element => {
  switch(category) {
    case 'frontend': return <FaCode />;
    case 'backend': return <HiServer />;
    case 'database': return <HiDatabase />;
    case 'infrastructure': return <HiCloud />;
    case 'apis': return <HiCog />;
    default: return <HiViewGrid />;
  }
};

const getCategoryLabel = (category: string) => {
  const labels: { [key: string]: string } = {
    frontend: 'Frontend',
    backend: 'Backend',
    database: 'Database',
    infrastructure: 'Infrastructure',
    apis: 'APIs & Integrations',
    other: 'Other'
  };
  return labels[category] || category;
};

export const EnhancedProjectCard: React.FC<EnhancedProjectCardProps> = ({
  title,
  description,
  image,
  repoLink,
  liveLink,
  techCategories,
  mediaType = 'image',
  badges
}) => {
  const [imageLoaded, setImageLoaded] = useState(false);
  const cardRef = useRef<HTMLDivElement>(null);

  const handleImageLoad = () => {
    setImageLoaded(true);
  };

  const totalTechCount = techCategories ?
    Object.values(techCategories).reduce((sum, arr) => sum + (arr?.length || 0), 0) : 0;

  return (
    <Card
      ref={cardRef}
      initial={{ opacity: 0, y: 30, filter: 'blur(10px)', rotateX: -5 }}
      whileInView={{ opacity: 1, y: 0, filter: 'blur(0px)', rotateX: 0 }}
      viewport={{ once: true, margin: "-50px" }}
      transition={{
        duration: 0.6,
        ease: [0.25, 0.46, 0.45, 0.94],
        filter: { duration: 0.4 }
      }}
    >
      <MediaContainer>
        {badges && badges.length > 0 && (
          <BadgeContainer>
            {badges.map(badge => (
              <ProjectBadge key={badge} $type={badge}>
                {badge === 'live' && 'Live'}
                {badge === 'demo' && 'Demo'}
                {badge === 'client' && 'Client Work'}
                {badge === 'internal' && 'Internal'}
              </ProjectBadge>
            ))}
          </BadgeContainer>
        )}
        {mediaType === 'component' && componentMap[image] ? (
          React.createElement(componentMap[image])
        ) : mediaType === 'video' ? (
          <ProjectVideo
            src={image}
            poster={image.replace('.mp4', '_poster.jpg').replace('_optimized.mp4', '_poster.jpg')}
            autoPlay
            loop
            muted
            playsInline
            preload="none"
            onLoadedData={handleImageLoad}
            style={{ opacity: imageLoaded ? 1 : 0 }}
            aria-label={`${title} demo video`}
          />
        ) : (
          <ProjectImage
            src={image}
            alt={title}
            loading="lazy"
            onLoad={handleImageLoad}
            style={{ opacity: imageLoaded ? 1 : 0 }}
            decoding="async"
          />
        )}
      </MediaContainer>

      <ContentSection>
        <div>
          <ProjectTitle>{title}</ProjectTitle>
          <ProjectDescription>{description}</ProjectDescription>
        </div>

        {techCategories && (
          <TechStackSection>
            <TechStackHeader>
              <HiViewGrid />
              <span>Tech Stack</span>
              {totalTechCount > 0 && (
                <TechCountBadge>
                  {totalTechCount} {totalTechCount === 1 ? 'Technology' : 'Technologies'}
                </TechCountBadge>
              )}
            </TechStackHeader>

            <TechCategoriesGrid>
              {Object.entries(techCategories).map(([category, technologies]) => {
                if (!technologies || technologies.length === 0) return null;

                return (
                  <TechCategoryRow key={category}>
                    <CategoryLabel>
                      {getCategoryIcon(category)}
                      {getCategoryLabel(category)}
                    </CategoryLabel>
                    <TechTagsContainer>
                      {technologies.map((tech: string) => {
                        const icon = getTechIcon(tech);
                        return (
                          <TechTag key={tech} $hasIcon={!!icon}>
                            {icon}
                            {tech}
                          </TechTag>
                        );
                      })}
                    </TechTagsContainer>
                  </TechCategoryRow>
                );
              })}
            </TechCategoriesGrid>
          </TechStackSection>
        )}

        <ProjectActions>
          <SecondaryButton
            href={repoLink}
            target="_blank"
            rel="noopener noreferrer"
            aria-label={`View ${title} repository`}
            style={{ color: '#0078ff' }}
          >
            <span style={{ color: '#0078ff' }}>View Code</span>
          </SecondaryButton>
          <ActionButton
            href={liveLink}
            target="_blank"
            rel="noopener noreferrer"
            aria-label={`View ${title} live demo`}
            style={{ color: '#ffffff' }}
          >
            <span style={{ color: '#ffffff' }}>Live Demo</span>
          </ActionButton>
        </ProjectActions>
      </ContentSection>
    </Card>
  );
};

export default EnhancedProjectCard;
