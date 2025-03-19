import React, { useState, useEffect } from 'react';
import { Link, useLocation } from 'react-router-dom';
import styled from 'styled-components';
import { motion, AnimatePresence } from 'framer-motion';
import type { ThemeMode } from '../../styles/theme.types';

interface NavigationBarProps {
  themeMode: ThemeMode;
  toggleTheme: () => void;
}

const NavContainer = styled.nav`
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: ${({ theme }) => theme.zIndex.header};
  background: ${({ theme }) => theme.colors.surface}CC;
  backdrop-filter: blur(8px);
  border-bottom: 1px solid ${({ theme }) => theme.colors.border};
  transition: all ${({ theme }) => theme.transitions.normal};
`;

const NavContent = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  max-width: 1200px;
  margin: 0 auto;
  padding: ${({ theme }) => `${theme.spacing.sm} ${theme.spacing.md}`};
`;

const Logo = styled(Link)`
  font-size: 1.5rem;
  font-weight: bold;
  color: ${({ theme }) => theme.colors.primary};
  text-decoration: none;
  display: flex;
  align-items: center;
  gap: ${({ theme }) => theme.spacing.xs};

  svg {
    width: 24px;
    height: 24px;
  }

  &:hover {
    color: ${({ theme }) => theme.colors.primaryHover};
  }
`;

const NavLinks = styled.div`
  display: flex;
  gap: ${({ theme }) => theme.spacing.md};
  align-items: center;

  @media (max-width: ${({ theme }) => theme.breakpoints.tablet}) {
    display: none;
  }
`;

const NavLink = styled(Link)<{ $isActive?: boolean }>`
  color: ${({ theme, $isActive }) => 
    $isActive ? theme.colors.primary : theme.colors.text};
  text-decoration: none;
  padding: ${({ theme }) => `${theme.spacing.xs} ${theme.spacing.sm}`};
  border-radius: ${({ theme }) => theme.borderRadius.small};
  transition: all ${({ theme }) => theme.transitions.fast};
  font-weight: 500;

  &:hover {
    color: ${({ theme }) => theme.colors.primary};
    background: ${({ theme }) => theme.colors.primaryLight};
  }
`;

const StatusIndicator = styled(motion.div)<{ $status: 'online' | 'offline' }>`
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: ${({ theme, $status }) => 
    $status === 'online' ? theme.colors.success : theme.colors.error};
  margin-left: ${({ theme }) => theme.spacing.xxs};
`;

const ServiceStatus = styled.div`
  display: flex;
  align-items: center;
  gap: ${({ theme }) => theme.spacing.xs};
  font-size: 0.875rem;
  color: ${({ theme }) => theme.colors.textSecondary};
  padding: ${({ theme }) => `${theme.spacing.xs} ${theme.spacing.sm}`};
  border-radius: ${({ theme }) => theme.borderRadius.small};
  background: ${({ theme }) => theme.colors.background};
`;

const NavigationBar: React.FC<NavigationBarProps> = ({ themeMode, toggleTheme }) => {
  const location = useLocation();
  const [servicesStatus, setServicesStatus] = useState({
    devpanel: false,
    urlshortener: false,
    messaging: false
  });

  // Check backend services status
  useEffect(() => {
    const checkServices = async () => {
      try {
        const services = {
          devpanel: 'http://localhost:8080/devpanel/health',
          urlshortener: 'http://localhost:8081/health',
          messaging: 'http://localhost:8082/health'
        };

        const results = await Promise.all(
          Object.entries(services).map(async ([service, url]) => {
            try {
              const response = await fetch(url);
              return [service, response.ok];
            } catch {
              return [service, false];
            }
          })
        );

        setServicesStatus(
          Object.fromEntries(results) as typeof servicesStatus
        );
      } catch (error) {
        console.error('Failed to check services status:', error);
      }
    };

    // Check initially and every 30 seconds
    checkServices();
    const interval = setInterval(checkServices, 30000);

    return () => clearInterval(interval);
  }, []);

  const isActive = (path: string) => location.pathname === path;

  return (
    <NavContainer>
      <NavContent>
        <Logo to="/">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
            <path d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5" />
          </svg>
          Jaden Razo
        </Logo>

        <NavLinks>
          <ServiceStatus>
            DevPanel
            <StatusIndicator
              $status={servicesStatus.devpanel ? 'online' : 'offline'}
              initial={{ scale: 0.8 }}
              animate={{ scale: 1 }}
              transition={{ repeat: Infinity, duration: 2 }}
            />
          </ServiceStatus>

          <ServiceStatus>
            URL Shortener
            <StatusIndicator
              $status={servicesStatus.urlshortener ? 'online' : 'offline'}
              initial={{ scale: 0.8 }}
              animate={{ scale: 1 }}
              transition={{ repeat: Infinity, duration: 2 }}
            />
          </ServiceStatus>

          <ServiceStatus>
            Messaging
            <StatusIndicator
              $status={servicesStatus.messaging ? 'online' : 'offline'}
              initial={{ scale: 0.8 }}
              animate={{ scale: 1 }}
              transition={{ repeat: Infinity, duration: 2 }}
            />
          </ServiceStatus>

          <NavLink to="/about" $isActive={isActive('/about')}>
            About
          </NavLink>

          <NavLink to="/contact" $isActive={isActive('/contact')}>
            Contact
          </NavLink>
        </NavLinks>
      </NavContent>
    </NavContainer>
  );
};

export default NavigationBar; 