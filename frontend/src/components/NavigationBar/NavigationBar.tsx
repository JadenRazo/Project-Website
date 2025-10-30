import React, { useState, useEffect } from 'react';
import { Link, useLocation } from 'react-router-dom';
import styled from 'styled-components';
import { motion, AnimatePresence } from 'framer-motion';
import type { ThemeMode } from '../../styles/theme.types';
import { useScrollTo } from '../../hooks/useScrollTo';

interface NavigationBarProps {
  themeMode: ThemeMode;
  toggleTheme: () => void;
}

interface ServiceStatus {
  name: string;
  status: 'operational' | 'degraded' | 'down';
  latency_ms?: number;
  uptime_percentage?: number;
  error?: string;
  description: string;
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

const NavLinks = styled.div<{ $isOpen: boolean }>`
  display: flex;
  gap: ${({ theme }) => theme.spacing.md};
  align-items: center;

  @media (max-width: ${({ theme }) => theme.breakpoints.tablet}) {
    display: ${({ $isOpen }) => ($isOpen ? 'flex' : 'none')};
    position: fixed;
    top: 60px;
    left: 0;
    right: 0;
    background: ${({ theme }) => theme.colors.surface};
    padding: ${({ theme }) => theme.spacing.md};
    flex-direction: column;
    align-items: flex-start;
    border-bottom: 1px solid ${({ theme }) => theme.colors.border};
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
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
  width: 100%;
  text-align: left;

  &:hover {
    color: ${({ theme }) => theme.colors.primary};
    background: ${({ theme }) => theme.colors.primaryLight};
  }

  @media (max-width: ${({ theme }) => theme.breakpoints.tablet}) {
    padding: ${({ theme }) => `${theme.spacing.sm} ${theme.spacing.md}`};
  }
`;

const ServiceStatusButton = styled(motion.button)`
  display: flex;
  align-items: center;
  gap: ${({ theme }) => theme.spacing.xs};
  background: transparent;
  border: none;
  padding: ${({ theme }) => `${theme.spacing.xs} ${theme.spacing.sm}`};
  color: ${({ theme }) => theme.colors.text};
  font-size: 0.875rem;
  cursor: pointer;
  transition: all ${({ theme }) => theme.transitions.fast};
  position: relative;
  border-radius: ${({ theme }) => theme.borderRadius.small};

  &:hover {
    background: ${({ theme }) => theme.colors.primaryLight};
  }

  svg {
    width: 16px;
    height: 16px;
    transition: transform ${({ theme }) => theme.transitions.fast};
  }

  &[aria-expanded="true"] svg {
    transform: rotate(180deg);
  }
`;

const ServiceStatusDropdown = styled(motion.div)`
  position: absolute;
  top: 100%;
  right: 0;
  margin-top: ${({ theme }) => theme.spacing.xs};
  background: ${({ theme }) => theme.colors.surface};
  border: 1px solid ${({ theme }) => theme.colors.border};
  border-radius: ${({ theme }) => theme.borderRadius.medium};
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  min-width: 280px;
  overflow: hidden;
  z-index: 1000;

  @media (max-width: ${({ theme }) => theme.breakpoints.tablet}) {
    position: static;
    margin-top: ${({ theme }) => theme.spacing.sm};
    width: 100%;
    box-shadow: none;
  }
`;

const ServiceItem = styled(motion.div)`
  padding: ${({ theme }) => theme.spacing.sm};
  border-bottom: 1px solid ${({ theme }) => theme.colors.border};
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: ${({ theme }) => theme.spacing.sm};

  &:last-child {
    border-bottom: none;
  }
`;

const ServiceInfo = styled.div`
  flex: 1;
`;

const ServiceName = styled.div`
  font-weight: 500;
  color: ${({ theme }) => theme.colors.text};
`;

const ServiceDescription = styled.div`
  font-size: 0.75rem;
  color: ${({ theme }) => theme.colors.textSecondary};
  margin-top: 2px;
`;

const StatusIndicator = styled(motion.div)<{ $status: 'operational' | 'degraded' | 'down' }>`
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: ${({ theme, $status }) => {
    switch ($status) {
      case 'operational': return '#10b981'; // Green
      case 'degraded': return '#f59e0b';    // Yellow
      case 'down': return '#ef4444';        // Red
      default: return theme.colors.textSecondary;
    }
  }};
  position: relative;

  &::after {
    content: '';
    position: absolute;
    inset: -4px;
    border-radius: 50%;
    background: ${({ theme, $status }) => {
      switch ($status) {
        case 'operational': return '#10b981'; // Green
        case 'degraded': return '#f59e0b';    // Yellow
        case 'down': return '#ef4444';        // Red
        default: return theme.colors.textSecondary;
      }
    }};
    opacity: 0.2;
    animation: ${({ $status }) => $status === 'operational' ? 'pulse 2s infinite' : 'none'};
  }

  @keyframes pulse {
    0% {
      transform: scale(1);
      opacity: 0.2;
    }
    50% {
      transform: scale(1.5);
      opacity: 0;
    }
    100% {
      transform: scale(1);
      opacity: 0;
    }
  }
`;

const HamburgerButton = styled.button`
  display: none;
  background: none;
  border: none;
  cursor: pointer;
  padding: ${({ theme }) => theme.spacing.xs};
  color: ${({ theme }) => theme.colors.text};
  z-index: 1000;

  @media (max-width: ${({ theme }) => theme.breakpoints.tablet}) {
    display: block;
  }
`;


const NavigationBar: React.FC<NavigationBarProps> = ({ themeMode, toggleTheme }) => {
  const location = useLocation();
  const [isMenuOpen, setIsMenuOpen] = useState(false);
  const [isServicesOpen, setIsServicesOpen] = useState(false);
  const [servicesStatus, setServicesStatus] = useState<ServiceStatus[]>([
    {
      name: 'API',
      status: 'operational',
      description: 'Main API service'
    },
    {
      name: 'Database',
      status: 'operational',
      description: 'PostgreSQL database'
    },
    {
      name: 'LOC Counter',
      status: 'operational',
      description: 'Code statistics service'
    }
  ]);


  useEffect(() => {
    setIsMenuOpen(false);
  }, [location]);

  const checkServices = React.useCallback(async () => {
    try {
      const apiUrl = (window as any)._env_?.REACT_APP_API_URL || process.env.REACT_APP_API_URL || '';
      const endpoint = apiUrl ? `${apiUrl}/api/v1/status/` : '/api/v1/status/';
      const response = await fetch(endpoint, {
        method: 'GET',
        headers: { 'Cache-Control': 'no-cache' }
      });
      
      if (response.ok) {
        const data = await response.json();
        const mappedServices = data.services.map((service: any) => ({
          name: getServiceDisplayName(service.name),
          status: service.status,
          latency_ms: service.latency_ms,
          uptime_percentage: service.uptime_percentage,
          error: service.error,
          description: getServiceDescription(service.name)
        }));
        setServicesStatus(mappedServices);
      } else {
        setServicesStatus(current => 
          current.map(service => ({ ...service, status: 'down' as const }))
        );
      }
    } catch (error) {
      console.error("Error checking services:", error);
      setServicesStatus(current => 
        current.map(service => ({ ...service, status: 'down' as const }))
      );
    }
  }, []);

  const getServiceDisplayName = (name: string): string => {
    switch (name.toLowerCase()) {
      case 'api': return 'API';
      case 'database': return 'Database';
      case 'code_stats': return 'LOC Counter';
      default: return name;
    }
  };

  const getServiceDescription = (name: string): string => {
    switch (name.toLowerCase()) {
      case 'api': return 'Main API service';
      case 'database': return 'PostgreSQL database';
      case 'code_stats': return 'Code statistics service';
      default: return `${name} service`;
    }
  };

  const getOverallStatus = (): 'operational' | 'degraded' | 'down' => {
    if (servicesStatus.length === 0) return 'down';
    
    const operationalCount = servicesStatus.filter(s => s.status === 'operational').length;
    const degradedCount = servicesStatus.filter(s => s.status === 'degraded').length;
    
    if (operationalCount === servicesStatus.length) return 'operational';
    if (operationalCount > 0 || degradedCount > 0) return 'degraded';
    return 'down';
  };
  
  useEffect(() => {
    checkServices();
    
    const interval = setInterval(checkServices, 30000);
    
    return () => clearInterval(interval);
  }, [checkServices]);

  const { scrollToTop } = useScrollTo();

  const handleLinkClick = () => {
    scrollToTop();
    setIsMenuOpen(false);
    setIsServicesOpen(false);
  };


  const isActive = (path: string) => location.pathname === path;

  return (
    <NavContainer>
      <NavContent>
        <Logo to="/" onClick={handleLinkClick}>
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
            <path d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5" />
          </svg>
          Jaden Razo
        </Logo>

        <HamburgerButton onClick={() => {
          setIsMenuOpen(!isMenuOpen);
          setIsServicesOpen(false);
        }}>
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
            {isMenuOpen ? (
              <path d="M18 6L6 18M6 6l12 12" />
            ) : (
              <path d="M3 12h18M3 6h18M3 18h18" />
            )}
          </svg>
        </HamburgerButton>

        <NavLinks $isOpen={isMenuOpen}>
          <div style={{ position: 'relative' }}>
            <ServiceStatusButton
              onClick={() => setIsServicesOpen(!isServicesOpen)}
              aria-expanded={isServicesOpen}
              whileHover={{ scale: 1.02 }}
              whileTap={{ scale: 0.98 }}
            >
              <StatusIndicator
                $status={getOverallStatus()}
              />
              Services Status
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <path d="M6 9l6 6 6-6" />
              </svg>
            </ServiceStatusButton>

            <AnimatePresence>
              {isServicesOpen && (
                <ServiceStatusDropdown
                  initial={{ opacity: 0, y: -10 }}
                  animate={{ opacity: 1, y: 0 }}
                  exit={{ opacity: 0, y: -10 }}
                  transition={{ duration: 0.2 }}
                >
                  {servicesStatus.map((service, index) => (
                    <ServiceItem
                      key={service.name}
                      initial={{ opacity: 0, x: -20 }}
                      animate={{ opacity: 1, x: 0 }}
                      exit={{ opacity: 0, x: -20 }}
                      transition={{ delay: index * 0.05 }}
                    >
                      <ServiceInfo>
                        <ServiceName>{service.name}</ServiceName>
                        <ServiceDescription>{service.description}</ServiceDescription>
                      </ServiceInfo>
                      <StatusIndicator
                        $status={service.status}
                        initial={{ scale: 0.8 }}
                        animate={{ scale: 1 }}
                        transition={{ repeat: Infinity, duration: 2 }}
                      />
                    </ServiceItem>
                  ))}
                  <ServiceItem
                    as={Link}
                    to="/status"
                    style={{ 
                      borderTop: '1px solid',
                      borderTopColor: 'inherit',
                      cursor: 'pointer',
                      textDecoration: 'none',
                      display: 'block'
                    }}
                    onClick={() => setIsServicesOpen(false)}
                  >
                    <div style={{ 
                      color: '#0078ff',
                      fontSize: '0.875rem',
                      fontWeight: '500',
                      width: '100%'
                    }}>
                      View Detailed Status Page →
                    </div>
                  </ServiceItem>
                </ServiceStatusDropdown>
              )}
            </AnimatePresence>
          </div>

          <NavLink to="/about" $isActive={isActive('/about')} onClick={handleLinkClick}>
            About
          </NavLink>

          <NavLink to="/contact" $isActive={isActive('/contact')} onClick={handleLinkClick}>
            Contact
          </NavLink>

          <NavLink to="/portfolio" $isActive={isActive('/portfolio')} onClick={handleLinkClick}>
            Portfolio
          </NavLink>

        </NavLinks>
      </NavContent>
    </NavContainer>
  );
};

export default NavigationBar; 