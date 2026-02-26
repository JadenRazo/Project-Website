import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import styled from 'styled-components';
import { motion } from 'framer-motion';
import { useScrollTo } from '../../hooks/useScrollTo';
import { useTheme } from '../../hooks/useTheme';

const FooterContainer = styled.footer`
  background: ${({ theme }) => theme.colors.surface};
  border-top: 1px solid ${({ theme }) => theme.colors.border};
  padding: ${({ theme }) => `${theme.spacing.xl} 0`};
  margin-top: auto;
  color: ${({ theme }) => theme.colors.text};
`;

const FooterContent = styled.div`
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 ${({ theme }) => theme.spacing.md};
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: ${({ theme }) => theme.spacing.xl};

  @media (max-width: ${({ theme }) => theme.breakpoints.tablet}) {
    grid-template-columns: 1fr;
    text-align: center;
    gap: ${({ theme }) => theme.spacing.lg};
    padding: 0 ${({ theme }) => theme.spacing.lg};
  }
`;

const FooterSection = styled.div`
  display: flex;
  flex-direction: column;
  gap: ${({ theme }) => theme.spacing.md};
  align-items: flex-start;
  
  p {
    color: ${({ theme }) => theme.colors.text};
    opacity: 0.8;
    margin: 0;
  }

  @media (max-width: ${({ theme }) => theme.breakpoints.tablet}) {
    align-items: center;
    gap: ${({ theme }) => theme.spacing.sm};
  }
`;

const FooterTitle = styled.h3`
  color: ${({ theme }) => theme.colors.primary};
  font-size: 1.2rem;
  margin-bottom: ${({ theme }) => theme.spacing.sm};
`;

const FooterLink = styled(Link)<{ $isNavigating?: boolean }>`
  color: ${({ theme }) => theme.colors.text};
  text-decoration: none;
  transition: all ${({ theme }) => theme.transitions.fast};
  font-size: 0.9rem;
  position: relative;
  opacity: ${({ $isNavigating }) => $isNavigating ? 0.7 : 1};

  &:hover {
    color: ${({ theme }) => theme.colors.primary};
  }

  &::after {
    content: '';
    position: absolute;
    bottom: -2px;
    left: 0;
    width: ${({ $isNavigating }) => $isNavigating ? '100%' : '0'};
    height: 2px;
    background: ${({ theme }) => theme.colors.primary};
    transition: width ${({ theme }) => theme.transitions.fast};
  }
`;

const FooterExternalLink = styled.a`
  color: ${({ theme }) => theme.colors.text};
  text-decoration: none;
  transition: color ${({ theme }) => theme.transitions.fast};
  font-size: 0.9rem;
  display: inline-flex;
  align-items: center;
  gap: ${({ theme }) => theme.spacing.xs};
  width: fit-content;

  &:hover {
    color: ${({ theme }) => theme.colors.primary};
  }

  svg {
    width: 16px;
    height: 16px;
  }

  @media (max-width: ${({ theme }) => theme.breakpoints.tablet}) {
    justify-content: center;
  }
`;

const CopyrightRow = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  gap: ${({ theme }) => theme.spacing.md};
  padding: ${({ theme }) => theme.spacing.md};
  color: ${({ theme }) => theme.colors.text};
  opacity: 0.7;
  font-size: 0.9rem;
  border-top: 1px solid ${({ theme }) => theme.colors.border};
  margin-top: ${({ theme }) => theme.spacing.xl};

  @media (max-width: ${({ theme }) => theme.breakpoints.mobile}) {
    flex-direction: column;
    gap: ${({ theme }) => theme.spacing.sm};
  }
`;

const ThemeToggleButton = styled.button`
  display: inline-flex;
  align-items: center;
  gap: 6px;
  background: none;
  border: 1px solid ${({ theme }) => theme.colors.border};
  border-radius: ${({ theme }) => theme.borderRadius.pill};
  padding: 6px 14px;
  color: ${({ theme }) => theme.colors.text};
  font-size: 0.85rem;
  cursor: pointer;
  transition: all ${({ theme }) => theme.transitions.fast};
  opacity: 1;

  &:hover {
    border-color: ${({ theme }) => theme.colors.primary};
    color: ${({ theme }) => theme.colors.primary};
  }

  svg {
    width: 14px;
    height: 14px;
  }
`;

const SocialLinks = styled.div`
  display: flex;
  gap: ${({ theme }) => theme.spacing.md};
  flex-wrap: wrap;
  
  @media (max-width: ${({ theme }) => theme.breakpoints.tablet}) {
    justify-content: center;
    width: 100%;
  }
`;

const TechStack = styled.div`
  display: flex;
  flex-wrap: wrap;
  gap: ${({ theme }) => theme.spacing.xs};

  @media (max-width: ${({ theme }) => theme.breakpoints.tablet}) {
    justify-content: center;
  }
`;

const TechBadge = styled(motion.span)`
  background: ${({ theme }) => theme.colors.primaryLight || `${theme.colors.primary}20`};
  color: ${({ theme }) => theme.colors.primary};
  padding: ${({ theme }) => `${theme.spacing.xxs} ${theme.spacing.xs}`};
  border-radius: ${({ theme }) => theme.borderRadius.small};
  font-size: 0.8rem;
  font-weight: 500;
`;

const Footer: React.FC = () => {
  const currentYear = new Date().getFullYear();
  const navigate = useNavigate();
  const [navigatingTo, setNavigatingTo] = useState<string | null>(null);
  const { themeMode, toggleTheme } = useTheme();
  
  const techStack = [
    'TypeScript',
    'React',
    'Go',
    'Python',
    'PowerShell',
    'Bash',
    'PostgreSQL',
    'Docker',
    'AWS',
    'Terraform',
    'Active Directory',
    'Linux',
    'Nginx',
    'Prometheus',
    'Grafana',
    'Git',
  ];

  const { scrollToTop } = useScrollTo();

  // Handler for footer link clicks
  const handleLinkClick = (e: React.MouseEvent<HTMLAnchorElement>) => {
    e.preventDefault();
    const href = e.currentTarget.getAttribute('href');
    const currentPath = window.location.pathname;
    
    // If staying on the same page, just scroll to top
    if (href === currentPath) {
      scrollToTop({ behavior: 'smooth' });
    } else if (href) {
      // Set navigating state for visual feedback
      setNavigatingTo(href);
      
      // Small delay for visual feedback before navigation
      setTimeout(() => {
        // Navigate with state to indicate footer navigation
        navigate(href, { state: { fromFooter: true } });
        
        // Clear navigation state after a delay
        setTimeout(() => setNavigatingTo(null), 500);
      }, 100);
    }
  };

  return (
    <FooterContainer>
      <FooterContent>
        <FooterSection>
          <FooterTitle>About</FooterTitle>
          <p>Systems Administrator &amp; Full Stack Developer. Building scalable infrastructure and web applications with modern technologies.</p>
          <TechStack>
            {techStack.map((tech, index) => (
              <TechBadge
                key={tech}
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: index * 0.1 }}
              >
                {tech}
              </TechBadge>
            ))}
          </TechStack>
        </FooterSection>

        <FooterSection>
          <FooterTitle>Projects</FooterTitle>
          <FooterLink to="/devpanel" onClick={handleLinkClick} $isNavigating={navigatingTo === '/devpanel'}>Developer Panel</FooterLink>
          <FooterLink to="/urlshortener" onClick={handleLinkClick} $isNavigating={navigatingTo === '/urlshortener'}>URL Shortener</FooterLink>
          <FooterLink to="/messaging" onClick={handleLinkClick} $isNavigating={navigatingTo === '/messaging'}>Real-time Messaging</FooterLink>
          <FooterLink to="/projects" onClick={handleLinkClick} $isNavigating={navigatingTo === '/projects'}>View All Projects</FooterLink>
          <FooterLink to="/blog" onClick={handleLinkClick} $isNavigating={navigatingTo === '/blog'}>Blog</FooterLink>
        </FooterSection>

        <FooterSection>
          <FooterTitle>Connect</FooterTitle>
          <SocialLinks>
            <FooterExternalLink
              href="https://github.com/JadenRazo"
              target="_blank"
              rel="noopener noreferrer"
            >
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <path d="M9 19c-5 1.5-5-2.5-7-3m14 6v-3.87a3.37 3.37 0 0 0-.94-2.61c3.14-.35 6.44-1.54 6.44-7A5.44 5.44 0 0 0 20 4.77 5.07 5.07 0 0 0 19.91 1S18.73.65 16 2.48a13.38 13.38 0 0 0-7 0C6.27.65 5.09 1 5.09 1A5.07 5.07 0 0 0 5 4.77a5.44 5.44 0 0 0-1.5 3.78c0 5.42 3.3 6.61 6.44 7A3.37 3.37 0 0 0 9 18.13V22" />
              </svg>
              GitHub
            </FooterExternalLink>
            <FooterExternalLink
              href="https://jadenrazo.dev/s/linkedin"
              target="_blank"
              rel="noopener noreferrer"
            >
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <path d="M16 8a6 6 0 0 1 6 6v7h-4v-7a2 2 0 0 0-2-2 2 2 0 0 0-2 2v7h-4v-7a6 6 0 0 1 6-6z" />
                <rect x="2" y="9" width="4" height="12" />
                <circle cx="4" cy="4" r="2" />
              </svg>
              LinkedIn
            </FooterExternalLink>
          </SocialLinks>
          <FooterExternalLink href="mailto:contact@jadenrazo.dev">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z" />
              <polyline points="22,6 12,13 2,6" />
            </svg>
            contact@jadenrazo.dev
          </FooterExternalLink>
        </FooterSection>
      </FooterContent>
      
      <CopyrightRow>
        <span>© {currentYear} Jaden Razo. All rights reserved.</span>
        <ThemeToggleButton onClick={toggleTheme} aria-label={`Switch to ${themeMode === 'dark' ? 'light' : 'dark'} mode`}>
          {themeMode === 'dark' ? (
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
              <circle cx="12" cy="12" r="5" />
              <line x1="12" y1="1" x2="12" y2="3" />
              <line x1="12" y1="21" x2="12" y2="23" />
              <line x1="4.22" y1="4.22" x2="5.64" y2="5.64" />
              <line x1="18.36" y1="18.36" x2="19.78" y2="19.78" />
              <line x1="1" y1="12" x2="3" y2="12" />
              <line x1="21" y1="12" x2="23" y2="12" />
              <line x1="4.22" y1="19.78" x2="5.64" y2="18.36" />
              <line x1="18.36" y1="5.64" x2="19.78" y2="4.22" />
            </svg>
          ) : (
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
              <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z" />
            </svg>
          )}
          {themeMode === 'dark' ? 'Light' : 'Dark'}
        </ThemeToggleButton>
      </CopyrightRow>
    </FooterContainer>
  );
};

export default Footer; 