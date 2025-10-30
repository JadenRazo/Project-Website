import React, { useState, useEffect } from 'react';
import styled, { keyframes } from 'styled-components';
import { motion, useAnimation } from 'framer-motion';
import { useInView } from 'react-intersection-observer';
import SystemMetrics from '../../components/metrics/SystemMetrics';

// Types
interface ServiceStatus {
  name: string;
  status: 'operational' | 'degraded' | 'down';
  latency_ms: number;
  last_checked: string;
  error?: string;
  uptime_percentage: number;
  details?: Record<string, any>;
}

interface SystemStatus {
  status: 'operational' | 'partial_outage' | 'major_outage';
  services: ServiceStatus[];
  last_updated: string;
  incidents: Incident[];
}

interface Incident {
  id: number;
  service: string;
  title: string;
  description: string;
  status: 'investigating' | 'identified' | 'monitoring' | 'resolved';
  severity: 'minor' | 'major' | 'critical';
  started_at: string;
  resolved_at?: string;
  created_at: string;
  updated_at: string;
}

// Styled Components matching your site's design system
const StatusContainer = styled.div`
  min-height: calc(100vh - 200px);
  padding: ${({ theme }) => theme.spacing.xxl} ${({ theme }) => theme.spacing.xl};
  background: ${({ theme }) => theme.colors.background};
  margin-top: 60px;
  overflow-x: hidden;

  @media (max-width: ${({ theme }) => theme.breakpoints.tablet}) {
    padding: ${({ theme }) => theme.spacing.xl} ${({ theme }) => theme.spacing.md};
  }
`;

const ContentWrapper = styled.div`
  max-width: 1200px;
  margin: 0 auto;
`;

const PageHeader = styled(motion.div)`
  text-align: center;
  margin-bottom: ${({ theme }) => theme.spacing.xxl};
  
  h1 {
    font-size: 3rem;
    font-weight: 700;
    margin-bottom: ${({ theme }) => theme.spacing.lg};
    background: linear-gradient(135deg, ${({ theme }) => theme.colors.primary}, ${({ theme }) => theme.colors.accent});
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
    
    @media (max-width: ${({ theme }) => theme.breakpoints.tablet}) {
      font-size: 2.5rem;
    }
  }
  
  p {
    font-size: 1.2rem;
    color: ${({ theme }) => theme.colors.textSecondary};
    max-width: 600px;
    margin: 0 auto;
    line-height: 1.6;
  }
`;

const OverallStatusCard = styled(motion.div)<{ status: string }>`
  background: ${({ theme, status }) => {
    switch (status) {
      case 'operational': return `linear-gradient(135deg, ${theme.colors.success}, ${theme.colors.successLight})`;
      case 'partial_outage': return `linear-gradient(135deg, ${theme.colors.warning}, ${theme.colors.warningLight})`;
      case 'major_outage': return `linear-gradient(135deg, ${theme.colors.error}, ${theme.colors.errorLight})`;
      default: return `linear-gradient(135deg, ${theme.colors.surface}, ${theme.colors.surfaceLight})`;
    }
  }};
  color: white;
  padding: ${({ theme }) => theme.spacing.xl};
  border-radius: ${({ theme }) => theme.borderRadius.large};
  text-align: center;
  margin-bottom: ${({ theme }) => theme.spacing.xxl};
  box-shadow: ${({ theme }) => theme.shadows.large};
  border: 1px solid ${({ theme }) => theme.colors.border};
  
  h2 {
    margin: 0;
    font-size: 1.8rem;
    font-weight: 600;
    margin-bottom: ${({ theme }) => theme.spacing.xs};
  }
  
  p {
    margin: 0;
    font-size: 1rem;
    opacity: 0.9;
  }
`;

const ServicesGrid = styled.div`
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
  gap: ${({ theme }) => theme.spacing.xl};
  margin-bottom: ${({ theme }) => theme.spacing.xxl};
  
  @media (max-width: ${({ theme }) => theme.breakpoints.tablet}) {
    grid-template-columns: 1fr;
    gap: ${({ theme }) => theme.spacing.lg};
  }
`;

const ServiceCard = styled(motion.div)<{ status: string }>`
  background: ${({ theme }) => theme.colors.surface};
  border: 1px solid ${({ theme }) => theme.colors.border};
  border-radius: ${({ theme }) => theme.borderRadius.medium};
  padding: ${({ theme }) => theme.spacing.xl};
  transition: all ${({ theme }) => theme.transitions.normal};
  position: relative;
  overflow: hidden;
  
  &:hover {
    transform: translateY(-4px);
    box-shadow: ${({ theme }) => theme.shadows.large};
    border-color: ${({ theme }) => theme.colors.primary};
  }
  
  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 4px;
    background: ${({ status, theme }) => {
      switch (status) {
        case 'operational': return theme.colors.success;
        case 'degraded': return theme.colors.warning;
        case 'down': return theme.colors.error;
        default: return theme.colors.border;
      }
    }};
  }
`;

const ServiceHeader = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: ${({ theme }) => theme.spacing.lg};
  
  h3 {
    margin: 0;
    font-size: 1.3rem;
    font-weight: 600;
    color: ${({ theme }) => theme.colors.text};
    text-transform: capitalize;
  }
`;

const StatusBadge = styled.span<{ status: string }>`
  background: ${({ status, theme }) => {
    switch (status) {
      case 'operational': return theme.colors.success;
      case 'degraded': return theme.colors.warning;
      case 'down': return theme.colors.error;
      default: return theme.colors.secondary;
    }
  }};
  color: white;
  padding: ${({ theme }) => `${theme.spacing.xxs} ${theme.spacing.xs}`};
  border-radius: ${({ theme }) => theme.borderRadius.pill};
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  display: inline-block;
`;

const ServiceMetrics = styled.div`
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: ${({ theme }) => theme.spacing.md};
  margin-top: ${({ theme }) => theme.spacing.lg};
  
  .metric {
    text-align: center;
    padding: ${({ theme }) => theme.spacing.sm};
    background: ${({ theme }) => theme.colors.backgroundAlt};
    border-radius: ${({ theme }) => theme.borderRadius.small};
    
    .label {
      font-size: 0.8rem;
      color: ${({ theme }) => theme.colors.textSecondary};
      margin-bottom: ${({ theme }) => theme.spacing.xs};
      text-transform: uppercase;
      letter-spacing: 0.5px;
      font-weight: 500;
    }
    
    .value {
      font-weight: 700;
      font-size: 1.1rem;
      color: ${({ theme }) => theme.colors.text};
    }
  }
`;

const ErrorDisplay = styled.div`
  background: ${({ theme }) => theme.colors.errorLight};
  color: ${({ theme }) => theme.colors.error};
  padding: ${({ theme }) => theme.spacing.sm};
  border-radius: ${({ theme }) => theme.borderRadius.small};
  font-size: 0.9rem;
  margin-top: ${({ theme }) => theme.spacing.sm};
  border: 1px solid ${({ theme }) => theme.colors.error};
`;

const IncidentsSection = styled(motion.div)`
  margin-top: ${({ theme }) => theme.spacing.xxl};
  
  h2 {
    font-size: 2rem;
    font-weight: 600;
    margin-bottom: ${({ theme }) => theme.spacing.xl};
    color: ${({ theme }) => theme.colors.text};
    display: flex;
    align-items: center;
    gap: ${({ theme }) => theme.spacing.sm};
    
    &::before {
      content: '⚠️';
      font-size: 1.5rem;
    }
  }
`;

const IncidentCard = styled(motion.div)<{ severity: string }>`
  background: ${({ theme }) => theme.colors.surface};
  border: 1px solid ${({ theme }) => theme.colors.border};
  border-left: 4px solid ${({ severity, theme }) => {
    switch (severity) {
      case 'critical': return theme.colors.error;
      case 'major': return theme.colors.warning;
      case 'minor': return theme.colors.accent;
      default: return theme.colors.border;
    }
  }};
  border-radius: ${({ theme }) => theme.borderRadius.medium};
  padding: ${({ theme }) => theme.spacing.xl};
  margin-bottom: ${({ theme }) => theme.spacing.lg};
  transition: all ${({ theme }) => theme.transitions.normal};
  
  &:hover {
    box-shadow: ${({ theme }) => theme.shadows.medium};
  }
  
  .incident-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: ${({ theme }) => theme.spacing.sm};
    
    h3 {
      margin: 0;
      font-size: 1.2rem;
      font-weight: 600;
      color: ${({ theme }) => theme.colors.text};
    }
  }
  
  .incident-meta {
    display: flex;
    gap: ${({ theme }) => theme.spacing.lg};
    margin-bottom: ${({ theme }) => theme.spacing.lg};
    font-size: 0.9rem;
    color: ${({ theme }) => theme.colors.textSecondary};
    
    @media (max-width: ${({ theme }) => theme.breakpoints.tablet}) {
      flex-direction: column;
      gap: ${({ theme }) => theme.spacing.xs};
    }
  }
  
  .incident-description {
    line-height: 1.7;
    color: ${({ theme }) => theme.colors.text};
  }
`;

const spin = keyframes`
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
`;

const LoadingContainer = styled.div`
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  height: 400px;
  
  .spinner {
    width: 50px;
    height: 50px;
    border: 4px solid ${({ theme }) => theme.colors.border};
    border-top: 4px solid ${({ theme }) => theme.colors.primary};
    border-radius: 50%;
    animation: ${spin} 1s linear infinite;
    margin-bottom: ${({ theme }) => theme.spacing.lg};
  }
  
  p {
    color: ${({ theme }) => theme.colors.textSecondary};
    font-size: 1.1rem;
  }
`;

const LastUpdated = styled.div`
  text-align: center;
  margin-top: ${({ theme }) => theme.spacing.xxl};
  padding: ${({ theme }) => theme.spacing.lg};
  background: ${({ theme }) => theme.colors.backgroundAlt};
  border-radius: ${({ theme }) => theme.borderRadius.medium};
  color: ${({ theme }) => theme.colors.textSecondary};
  font-size: 0.9rem;
  
  strong {
    color: ${({ theme }) => theme.colors.text};
  }
`;

const EmptyState = styled.div`
  text-align: center;
  padding: ${({ theme }) => theme.spacing.xxl};
  color: ${({ theme }) => theme.colors.textSecondary};
  
  h3 {
    margin-bottom: ${({ theme }) => theme.spacing.sm};
    color: ${({ theme }) => theme.colors.text};
  }
`;

// Helper functions
const getServiceDisplayName = (name: string): string => {
  switch (name.toLowerCase()) {
    case 'api': return 'API';
    case 'database': return 'Database';
    case 'code_stats': return 'LOC Counter';
    default: return name;
  }
};

const getStatusText = (status: string) => {
  switch (status) {
    case 'operational': return 'All Systems Operational';
    case 'partial_outage': return 'Partial System Outage';
    case 'major_outage': return 'Major System Outage';
    default: return 'System Status Unknown';
  }
};

const getStatusIcon = (status: string) => {
  switch (status) {
    case 'operational': return '✅';
    case 'partial_outage': return '⚠️';
    case 'major_outage': return '🚨';
    default: return '❓';
  }
};

const formatUptime = (uptime: number) => {
  if (uptime == null) return '0.00%';
  return `${uptime.toFixed(2)}%`;
};

const formatLatency = (latency: number) => {
  return `${latency}ms`;
};

const formatTime = (timestamp: string) => {
  return new Date(timestamp).toLocaleString();
};

const Status: React.FC = () => {
  const [systemStatus, setSystemStatus] = useState<SystemStatus | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Animation controls for in-view animations
  const controls = useAnimation();
  const [ref, inView] = useInView({
    threshold: 0.1,
    triggerOnce: true,
  });

  useEffect(() => {
    if (inView) {
      controls.start('visible');
    }
  }, [controls, inView]);

  const fetchStatus = async () => {
    try {
      const apiUrl = (window as any)._env_?.REACT_APP_API_URL || process.env.REACT_APP_API_URL || '';
      const endpoint = apiUrl ? `${apiUrl}/api/v1/status/` : '/api/v1/status/';
      const response = await fetch(endpoint);
      if (!response.ok) {
        throw new Error(`Failed to fetch status: ${response.status}`);
      }
      const data = await response.json();
      setSystemStatus(data);
      setError(null);
    } catch (err) {
      console.error('Error fetching status:', err);
      setError(err instanceof Error ? err.message : 'Unknown error');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchStatus();
    
    // Refresh every 30 seconds
    const interval = setInterval(fetchStatus, 30000);
    return () => clearInterval(interval);
  }, []);

  // Animation variants
  const containerVariants = {
    hidden: { opacity: 0 },
    visible: {
      opacity: 1,
      transition: {
        staggerChildren: 0.1,
        delayChildren: 0.2,
      },
    },
  };

  const itemVariants = {
    hidden: { opacity: 0, y: 20 },
    visible: {
      opacity: 1,
      y: 0,
      transition: {
        duration: 0.6,
        ease: "easeOut",
      },
    },
  };

  if (loading) {
    return (
      <StatusContainer>
        <ContentWrapper>
          <PageHeader
            initial={{ opacity: 0, y: -30 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.6 }}
          >
            <h1>System Status</h1>
            <p>Real-time monitoring of all services</p>
          </PageHeader>
          <LoadingContainer>
            <div className="spinner" />
            <p>Loading system status...</p>
          </LoadingContainer>
        </ContentWrapper>
      </StatusContainer>
    );
  }

  if (error) {
    return (
      <StatusContainer>
        <ContentWrapper>
          <PageHeader
            initial={{ opacity: 0, y: -30 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.6 }}
          >
            <h1>System Status</h1>
            <p>Unable to load status information</p>
          </PageHeader>
          <OverallStatusCard 
            status="major_outage"
            initial={{ opacity: 0, scale: 0.9 }}
            animate={{ opacity: 1, scale: 1 }}
            transition={{ duration: 0.5, delay: 0.2 }}
          >
            <h2>🚨 Status Unavailable</h2>
            <p>{error}</p>
          </OverallStatusCard>
        </ContentWrapper>
      </StatusContainer>
    );
  }

  if (!systemStatus) {
    return null;
  }

  return (
    <StatusContainer>
      <ContentWrapper>
        <PageHeader
          initial={{ opacity: 0, y: -30 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.6 }}
        >
          <h1>System Status</h1>
          <p>Real-time monitoring of all backend services and infrastructure</p>
        </PageHeader>

        <motion.div
          ref={ref}
          variants={containerVariants}
          initial="hidden"
          animate={controls}
        >
          <motion.div variants={itemVariants}>
            <OverallStatusCard status={systemStatus.status}>
              <h2>
                {getStatusIcon(systemStatus.status)} {getStatusText(systemStatus.status)}
              </h2>
              <p>All services are being monitored continuously</p>
            </OverallStatusCard>
          </motion.div>

          <ServicesGrid>
            {systemStatus.services.map((service, index) => (
              <motion.div key={service.name} variants={itemVariants}>
                <ServiceCard status={service.status}>
                  <ServiceHeader>
                    <h3>{getServiceDisplayName(service.name)}</h3>
                    <StatusBadge status={service.status}>
                      {service.status}
                    </StatusBadge>
                  </ServiceHeader>
                  
                  {service.error && (
                    <ErrorDisplay>
                      <strong>Error:</strong> {service.error}
                    </ErrorDisplay>
                  )}
                  
                  <ServiceMetrics>
                    <div className="metric">
                      <div className="label">Uptime</div>
                      <div className="value">{formatUptime(service.uptime_percentage)}</div>
                    </div>
                    <div className="metric">
                      <div className="label">Latency</div>
                      <div className="value">{formatLatency(service.latency_ms)}</div>
                    </div>
                    <div className="metric">
                      <div className="label">Last Check</div>
                      <div className="value">
                        {new Date(service.last_checked).toLocaleTimeString()}
                      </div>
                    </div>
                  </ServiceMetrics>
                </ServiceCard>
              </motion.div>
            ))}
          </ServicesGrid>

          {systemStatus.incidents && systemStatus.incidents.length > 0 && (
            <motion.div variants={itemVariants}>
              <IncidentsSection>
                <h2>Active Incidents</h2>
                {systemStatus.incidents.map((incident, index) => (
                  <IncidentCard
                    key={incident.id}
                    severity={incident.severity}
                    initial={{ opacity: 0, x: -20 }}
                    animate={{ opacity: 1, x: 0 }}
                    transition={{ duration: 0.5, delay: index * 0.1 }}
                  >
                    <div className="incident-header">
                      <h3>{incident.title}</h3>
                      <StatusBadge status={incident.status}>
                        {incident.status}
                      </StatusBadge>
                    </div>
                    <div className="incident-meta">
                      <span><strong>Service:</strong> {incident.service}</span>
                      <span><strong>Severity:</strong> {incident.severity}</span>
                      <span><strong>Started:</strong> {formatTime(incident.started_at)}</span>
                    </div>
                    <div className="incident-description">
                      {incident.description}
                    </div>
                  </IncidentCard>
                ))}
              </IncidentsSection>
            </motion.div>
          )}

          {(!systemStatus.incidents || systemStatus.incidents.length === 0) && (
            <motion.div variants={itemVariants}>
              <EmptyState>
                <h3>🎉 No Active Incidents</h3>
                <p>All systems are operating normally with no reported issues.</p>
              </EmptyState>
            </motion.div>
          )}

          <motion.div variants={itemVariants}>
            <SystemMetrics />
          </motion.div>

          <motion.div variants={itemVariants}>
            <LastUpdated>
              <strong>Last updated:</strong> {formatTime(systemStatus.last_updated)}
              <br />
              <small>Status checks run automatically every 30 seconds</small>
            </LastUpdated>
          </motion.div>
        </motion.div>
      </ContentWrapper>
    </StatusContainer>
  );
};

export default Status;