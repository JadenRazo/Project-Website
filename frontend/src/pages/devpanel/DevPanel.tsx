import React, { useState, useEffect, useCallback } from 'react';
import styled from 'styled-components';
import Collapsible from 'react-collapsible';
import ProjectManager from '../../components/devpanel/ProjectManager';
import ProjectPathsManager from '../../components/devpanel/ProjectPathsManager';
import CertificationsManager from '../../components/devpanel/CertificationsManager';
import PromptsManager from '../../components/devpanel/PromptsManager';
import VisitorAnalytics from '../../components/devpanel/VisitorAnalytics';
import AdminLogin from '../../components/devpanel/AdminLogin';
import DevPanelLoadingState from '../../components/devpanel/DevPanelLoadingState';
import { api } from '../../utils/apiConfig';
import { handleError } from '../../utils/errorHandler';
import { useAuthStore } from '../../stores/authStore';
import { 
  AnimatedContainer, 
  useReducedMotion 
} from '../../components/animations/AnimatedComponents';
import { useScrollTo } from '../../hooks/useScrollTo';

import ScrollToTopButton from '../../components/common/ScrollToTopButton';
import {
  LineChart,
  Line,
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer
} from 'recharts';

// Types

interface ServiceMetricsData {
  name: string;
  status: string;
  uptime: string;
  memoryUsage: number;
  cpuUsage: number;
  requestCount: number;
  errorCount: number;
  averageResponseTime: number;
  lastError: string;
  lastUpdate: string;
}

interface SystemMetrics {
  cpuUsage: number;
  memoryUsage: number;
  diskUsage: number;
  uptime: string;
  goVersion: string;
  numGoroutines: number;
  lastUpdate: string;
}

interface ServiceControlResponse {
  success: boolean;
  message: string;
}

// Styled Components
const PageContainer = styled.div`
  padding: 2rem;
  padding-top: calc(80px + 2rem);
  max-width: 1200px;
  margin: 0 auto;
  min-height: 100vh;
  position: relative;
  
  @media (max-width: ${({ theme }) => theme.breakpoints.tablet}) {
    padding: 1.5rem;
    padding-top: calc(70px + 1.5rem);
  }
  
  @media (max-width: ${({ theme }) => theme.breakpoints.mobile}) {
    padding: 1rem;
    padding-top: calc(60px + 1rem);
  }
`;

const Header = styled.header`
  margin-bottom: 2rem;
  border-bottom: 1px solid ${({ theme }) => theme.colors.border};
  padding-bottom: 1rem;
  
  @media (max-width: ${({ theme }) => theme.breakpoints.mobile}) {
    margin-bottom: 1.5rem;
  }
`;

const Title = styled.h1`
  font-size: 2.5rem;
  margin-bottom: 0.5rem;
  color: ${({ theme }) => theme.colors.primary};
  
  @media (max-width: ${({ theme }) => theme.breakpoints.mobile}) {
    font-size: 2rem;
  }
`;

const Subtitle = styled.p`
  font-size: 1.1rem;
  color: ${({ theme }) => theme.colors.textSecondary};
  
  @media (max-width: ${({ theme }) => theme.breakpoints.mobile}) {
    font-size: 1rem;
  }
`;

const MetricsGrid = styled.div`
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
  
  @media (max-width: ${({ theme }) => theme.breakpoints.tablet}) {
    grid-template-columns: repeat(2, 1fr);
    gap: 1rem;
  }
  
  @media (max-width: ${({ theme }) => theme.breakpoints.mobile}) {
    grid-template-columns: 1fr;
    gap: 0.75rem;
  }
`;

const MetricCard = styled.div`
  background: ${({ theme }) => theme.colors.card};
  border-radius: 8px;
  padding: 1.5rem;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.05);
  border: 1px solid ${({ theme }) => theme.colors.border};
  
  @media (max-width: ${({ theme }) => theme.breakpoints.mobile}) {
    padding: 1rem;
  }
`;

const MetricTitle = styled.h3`
  font-size: 1.1rem;
  margin-bottom: 1rem;
  color: ${({ theme }) => theme.colors.text};
`;

const SystemMetricValue = styled.div`
  font-size: 1.5rem;
  font-weight: 600;
  color: ${({ theme }) => theme.colors.primary};
  margin-bottom: 0.5rem;
  
  @media (max-width: ${({ theme }) => theme.breakpoints.mobile}) {
    font-size: 1.25rem;
  }
`;

const SystemMetricLabel = styled.div`
  font-size: 0.9rem;
  color: ${({ theme }) => theme.colors.textSecondary};
  margin-top: 0.5rem;
  
  @media (max-width: ${({ theme }) => theme.breakpoints.mobile}) {
    font-size: 0.8rem;
  }
`;

const ChartContainer = styled.div`
  background: ${({ theme }) => theme.colors.card};
  border-radius: 8px;
  padding: 1.5rem;
  margin-bottom: 2rem;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.05);
  border: 1px solid ${({ theme }) => theme.colors.border};
  height: 400px;
  overflow: hidden;
  
  @media (max-width: ${({ theme }) => theme.breakpoints.tablet}) {
    height: 350px;
    padding: 1.25rem;
  }
  
  @media (max-width: ${({ theme }) => theme.breakpoints.mobile}) {
    height: 280px;
    padding: 1rem;
    margin-bottom: 1.5rem;
  }
  
  @media (max-width: 480px) {
    height: 250px;
    padding: 0.75rem;
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


// Service-specific styled components
const ServiceGrid = styled.div`
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
  
  @media (max-width: 1024px) {
    grid-template-columns: repeat(2, 1fr);
    gap: 1.25rem;
  }
  
  @media (max-width: ${({ theme }) => theme.breakpoints.tablet}) {
    grid-template-columns: 1fr;
    gap: 1rem;
  }
  
  @media (max-width: ${({ theme }) => theme.breakpoints.mobile}) {
    gap: 0.75rem;
  }
`;

const ServiceCard = styled.div`
  background: ${({ theme }) => theme.colors.card};
  border-radius: 8px;
  padding: 1.5rem;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.05);
  border: 1px solid ${({ theme }) => theme.colors.border};
  display: flex;
  flex-direction: column;
  gap: 1rem;
  
  @media (max-width: ${({ theme }) => theme.breakpoints.mobile}) {
    padding: 1rem;
  }
`;

const ServiceHeader = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
`;

const ServiceTitle = styled.h3`
  font-size: 1.25rem;
  color: ${({ theme }) => theme.colors.text};
  margin: 0;
  
  @media (max-width: ${({ theme }) => theme.breakpoints.mobile}) {
    font-size: 1.1rem;
  }
`;

const ServiceStatus = styled.span<{ status: string }>`
  padding: 0.25rem 0.75rem;
  border-radius: 12px;
  font-size: 0.875rem;
  font-weight: 500;
  background: ${({ theme, status }) => 
    status === 'running' ? theme.colors.success + '20' :
    status === 'stopped' ? theme.colors.error + '20' :
    theme.colors.warning + '20'
  };
  color: ${({ theme, status }) => 
    status === 'running' ? theme.colors.success :
    status === 'stopped' ? theme.colors.error :
    theme.colors.warning
  };
`;

const ServiceMetrics = styled.div`
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(120px, 1fr));
  gap: 1rem;
  margin-bottom: 1rem;
  
  @media (max-width: ${({ theme }) => theme.breakpoints.mobile}) {
    grid-template-columns: repeat(2, 1fr);
  }
`;

const MetricItem = styled.div`
  text-align: center;
  padding: 0.75rem;
  background: ${({ theme }) => theme.colors.background};
  border-radius: 6px;
  
  @media (max-width: ${({ theme }) => theme.breakpoints.mobile}) {
    padding: 0.5rem;
  }
`;

const ServiceControls = styled.div`
  display: flex;
  gap: 0.5rem;
  margin-top: auto;
  
  @media (max-width: ${({ theme }) => theme.breakpoints.tablet}) {
    flex-wrap: wrap;
    gap: 0.4rem;
  }
  
  @media (max-width: ${({ theme }) => theme.breakpoints.mobile}) {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 0.5rem;
    
    button:nth-child(4),
    button:nth-child(5) {
      grid-column: span 1.5;
    }
  }
  
  @media (max-width: 480px) {
    grid-template-columns: repeat(2, 1fr);
    
    button:nth-child(4),
    button:nth-child(5) {
      grid-column: span 1;
    }
  }
`;

const ControlButton = styled.button<{ variant: 'start' | 'stop' | 'restart' | 'logs' | 'config' }>`
  flex: 1;
  padding: 0.5rem;
  border: none;
  border-radius: 4px;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  
  background: ${({ theme, variant }) => 
    variant === 'start' ? theme.colors.success :
    variant === 'stop' ? theme.colors.error :
    variant === 'restart' ? theme.colors.warning :
    variant === 'logs' ? theme.colors.accent :
    theme.colors.primary
  };
  color: white;
  
  &:hover {
    opacity: 0.9;
    transform: translateY(-1px);
  }
  
  &:disabled {
    opacity: 0.5;
    cursor: not-allowed;
    transform: none;
  }
  
  @media (max-width: ${({ theme }) => theme.breakpoints.mobile}) {
    padding: 0.4rem;
    font-size: 0.8rem;
  }
`;

// Rename MetricValue to ServiceMetricValue
const ServiceMetricValue = styled.div`
  font-size: 1.25rem;
  font-weight: 600;
  color: ${({ theme }) => theme.colors.primary};
  margin-bottom: 0.25rem;
  
  @media (max-width: ${({ theme }) => theme.breakpoints.mobile}) {
    font-size: 1.1rem;
  }
`;

// Rename MetricLabel to ServiceMetricLabel
const ServiceMetricLabel = styled.div`
  font-size: 0.75rem;
  color: ${({ theme }) => theme.colors.textSecondary};
  text-transform: uppercase;
  letter-spacing: 0.05em;
`;

// Styled wrapper for collapsible sections with scroll margin
const CollapsibleSection = styled.div`
  scroll-margin-top: 100px;
  
  @media (max-width: ${({ theme }) => theme.breakpoints.tablet}) {
    scroll-margin-top: 80px;
  }
  
  @media (max-width: ${({ theme }) => theme.breakpoints.mobile}) {
    scroll-margin-top: 70px;
  }
`;

// New styled components for collapsible sections
const SectionHeader = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1rem;
  background: ${({ theme }) => theme.colors.card};
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
  margin-bottom: 1rem;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
  border: 1px solid ${({ theme }) => theme.colors.border};
  
  &:hover {
    background: ${({ theme }) => theme.colors.background};
    transform: translateY(-1px);
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  }
  
  h2 {
    margin: 0;
    font-size: 1.5rem;
    color: ${({ theme }) => theme.colors.text};
  }
  
  @media (max-width: ${({ theme }) => theme.breakpoints.mobile}) {
    padding: 0.75rem;
    
    h2 {
      font-size: 1.25rem;
    }
  }
`;

const SectionIcon = styled.span<{ isOpen: boolean }>`
  transform: ${({ isOpen }) => isOpen ? 'rotate(180deg)' : 'rotate(0)'};
  transition: transform 0.3s ease;
  font-size: 1.5rem;
  color: ${({ theme }) => theme.colors.textSecondary};
  display: inline-block;
  width: 1.5rem;
  height: 1.5rem;
  display: flex;
  align-items: center;
  justify-content: center;
`;

const SectionContent = styled.div`
  padding: 1rem;
  background: ${({ theme }) => theme.colors.background};
  border-radius: 8px;
  margin-bottom: 1rem;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
  border: 1px solid ${({ theme }) => theme.colors.border};
  
  @media (max-width: ${({ theme }) => theme.breakpoints.mobile}) {
    padding: 0.75rem;
  }
`;

const RefreshButton = styled.button`
  position: fixed;
  bottom: 2rem;
  right: 2rem;
  padding: 1rem;
  border-radius: 50%;
  background: ${({ theme }) => theme.colors.primary};
  color: white;
  border: none;
  cursor: pointer;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 50;
  font-size: 1.25rem;
  width: 56px;
  height: 56px;
  
  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 6px 8px rgba(0, 0, 0, 0.15);
  }
  
  &:active {
    transform: translateY(0);
  }
  
  &:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
  
  @media (max-width: ${({ theme }) => theme.breakpoints.tablet}) {
    bottom: 1.5rem;
    right: 1.5rem;
    width: 48px;
    height: 48px;
    font-size: 1.1rem;
  }
  
  @media (max-width: ${({ theme }) => theme.breakpoints.mobile}) {
    bottom: 1rem;
    right: 1rem;
    width: 40px;
    height: 40px;
    padding: 0.5rem;
    font-size: 1rem;
  }
`;

const LoadingOverlay = styled.div`
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: ${({ theme }) => theme.zIndex.modal};
`;

const LoadingContent = styled.div`
  background: ${({ theme }) => theme.colors.card};
  padding: 2rem;
  border-radius: 8px;
  text-align: center;
  color: ${({ theme }) => theme.colors.text};
`;

const LoadingSpinner = styled.div`
  text-align: center;
  padding: 2rem;
  color: ${({ theme }) => theme.colors.textSecondary};
  font-size: 1.1rem;
`;

const GlobalStyle = `
  @keyframes spin {
    0% { transform: rotate(0deg); }
    100% { transform: rotate(360deg); }
  }
`;

// Inject global styles
if (typeof document !== 'undefined') {
  const style = document.createElement('style');
  style.textContent = GlobalStyle;
  document.head.appendChild(style);
}

const HeaderContent = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 1rem;
  
  @media (max-width: ${({ theme }) => theme.breakpoints.mobile}) {
    flex-direction: column;
    align-items: flex-start;
  }
`;

const HeaderActions = styled.div`
  display: flex;
  align-items: center;
  gap: 1rem;
  
  @media (max-width: ${({ theme }) => theme.breakpoints.mobile}) {
    width: 100%;
    justify-content: space-between;
  }
`;

const WelcomeText = styled.span`
  color: ${({ theme }) => theme.colors.textSecondary};
  font-size: 0.9rem;
  
  @media (max-width: 480px) {
    display: none;
  }
`;

// Main Component
const DevPanel: React.FC = () => {

  
  const reducedMotion = useReducedMotion();
  const { user: adminUser, isAuthenticated: isAdminAuthenticated, isLoading: checkingAuth } = useAuthStore();
  const { scrollToTop, scrollToElement, scrollToId } = useScrollTo();
  const [systemMetrics, setSystemMetrics] = useState<SystemMetrics | null>(null);
  const [serviceMetrics, setServiceMetrics] = useState<ServiceMetricsData[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [controlError, setControlError] = useState<string | null>(null);
  const [isRefreshing, setIsRefreshing] = useState<boolean>(false);
  const [lastUpdateTime, setLastUpdateTime] = useState<string>('');
  const [openSections, setOpenSections] = useState<{
    system: boolean;
    charts: boolean;
    services: boolean;
    projects: boolean;
    projectPaths: boolean;
    certifications: boolean;
    prompts: boolean;
    visitorAnalytics: boolean;
  }>({
    system: true,
    charts: true,
    services: true,
    projects: true,
    projectPaths: true,
    certifications: true,
    prompts: true,
    visitorAnalytics: true,
  });

  const toggleSection = (section: keyof typeof openSections) => {
    setOpenSections(prev => ({
      ...prev,
      [section]: !prev[section],
    }));
    
    // Scroll to the section when opening
    if (!openSections[section]) {
      // Small delay to allow the section to start opening
      setTimeout(() => {
        scrollToId(`section-${section}`, { 
          behavior: 'smooth',
          offset: 80
        });
      }, 150);
    }
  };

  const fetchData = useCallback(async (showRefresh = true) => {
    try {
      if (!systemMetrics && !serviceMetrics.length) {
        setLoading(true);
      }
      if (showRefresh) {
        setIsRefreshing(true);
      }
      setError(null);

      // Fetch system metrics and service metrics in parallel
      const [systemData, servicesData] = await Promise.all([
        api.get<SystemMetrics>('/api/v1/devpanel/system'),
        api.get<ServiceMetricsData[]>('/api/v1/devpanel/services'),
      ]);

      setSystemMetrics(systemData);
      setServiceMetrics(servicesData);
      setLastUpdateTime(new Date().toLocaleTimeString());
    } catch (err) {
      handleError(err, { context: 'DevPanel.fetchData' });
      setError('Failed to load data. Please try again later.');
    } finally {
      setLoading(false);
      if (showRefresh) {
        setIsRefreshing(false);
      }
    }
  }, [systemMetrics, serviceMetrics]);

  const controlService = async (serviceName: string, action: 'start' | 'stop' | 'restart') => {
    try {
      setControlError(null);
      
      const data = await api.post<ServiceControlResponse>(
        `/api/v1/devpanel/services/${serviceName}/${action}`
      );
      
      if (!data.success) {
        throw new Error(data.message);
      }

      // Refresh data after successful control action
      await fetchData(false);
    } catch (error: unknown) {
      handleError(error, { 
        context: 'DevPanel.controlService',
        service: serviceName,
        action 
      });
      const errorMessage = error instanceof Error ? error.message : 'Unknown error occurred';
      setControlError(`Failed to ${action} service: ${errorMessage}`);
    }
  };

  useEffect(() => {
    checkAdminAuth();
  }, []);

  useEffect(() => {
    if (isAdminAuthenticated && adminUser) {
      fetchData(false);
      const interval = setInterval(() => fetchData(false), 60000); // Update every 60 seconds instead of 30
      return () => clearInterval(interval);
    }
  }, [isAdminAuthenticated, adminUser, fetchData]);

  const checkAdminAuth = async () => {
    const token = localStorage.getItem('auth_token');
    if (token) {
      useAuthStore.getState().validateToken(token);
    } else {
      useAuthStore.getState().setLoading(false);
    }
  };

  const handleLoginSuccess = (userData: any) => {
    useAuthStore.getState().setUser(userData);
    useAuthStore.getState().setAuthenticated(true);
  };

  const handleLogout = () => {
    useAuthStore.getState().logout();
  };

  if (checkingAuth) {
    return (
      <PageContainer>
        <DevPanelLoadingState />
      </PageContainer>
    );
  }

  if (!isAdminAuthenticated) {
    return <AdminLogin onLoginSuccess={handleLoginSuccess} />;
  }

  // Prepare data for charts
  const serviceMemoryData = serviceMetrics.map(service => ({
    name: service.name,
    memory: service.memoryUsage,
    cpu: service.cpuUsage,
  }));

  const serviceResponseData = serviceMetrics.map(service => ({
    name: service.name,
    responseTime: service.averageResponseTime,
    requests: service.requestCount,
  }));


  return (
    <PageContainer>
      <Header>
        <HeaderContent>
          <div>
            <Title>Developer Panel</Title>
            <Subtitle>System monitoring and service management</Subtitle>
          </div>
          <HeaderActions>
            <WelcomeText>
              Welcome, {adminUser?.email}
            </WelcomeText>
            <ControlButton 
              variant="config" 
              onClick={handleLogout}
              style={{ padding: '0.5rem 1rem', fontSize: '0.9rem' }}
            >
              Logout
            </ControlButton>
          </HeaderActions>
        </HeaderContent>
      </Header>

      {error && <ErrorMessage>{error}</ErrorMessage>}
      {controlError && <ErrorMessage>{controlError}</ErrorMessage>}

      {loading ? (
        <LoadingSpinner>Loading...</LoadingSpinner>
      ) : (
        <>
          <CollapsibleSection>
            <Collapsible
              id="section-system"
              trigger={
                <SectionHeader onClick={() => toggleSection('system')}>
                  <h2>System Overview</h2>
                  <SectionIcon isOpen={openSections.system}>▼</SectionIcon>
                </SectionHeader>
              }
              open={openSections.system}
              onTriggerOpening={() => {}}
            onTriggerClosing={() => {}}
            transitionTime={200}
          >
            <SectionContent>
              <MetricsGrid>
                {systemMetrics && (
                  <>
                    <MetricCard>
                      <MetricTitle>System Memory Usage</MetricTitle>
                      <SystemMetricValue>{systemMetrics.memoryUsage?.toFixed(2) ?? '0.00'}%</SystemMetricValue>
                      <SystemMetricLabel>Total system memory utilization</SystemMetricLabel>
                    </MetricCard>
                    <MetricCard>
                      <MetricTitle>CPU Usage</MetricTitle>
                      <SystemMetricValue>{systemMetrics.cpuUsage?.toFixed(2) ?? '0.00'}%</SystemMetricValue>
                      <SystemMetricLabel>Total CPU utilization</SystemMetricLabel>
                    </MetricCard>
                    <MetricCard>
                      <MetricTitle>Disk Usage</MetricTitle>
                      <SystemMetricValue>{systemMetrics.diskUsage?.toFixed(2) ?? '0.00'}%</SystemMetricValue>
                      <SystemMetricLabel>Total disk space utilization</SystemMetricLabel>
                    </MetricCard>
                    <MetricCard>
                      <MetricTitle>Active Goroutines</MetricTitle>
                      <SystemMetricValue>{systemMetrics.numGoroutines}</SystemMetricValue>
                      <SystemMetricLabel>Currently running goroutines</SystemMetricLabel>
                    </MetricCard>
                  </>
                )}
              </MetricsGrid>
            </SectionContent>
          </Collapsible>
          </CollapsibleSection>

          <CollapsibleSection>
            <Collapsible
              id="section-charts"
            trigger={
              <SectionHeader onClick={() => toggleSection('charts')}>
                <h2>Service Analytics</h2>
                <SectionIcon isOpen={openSections.charts}>▼</SectionIcon>
              </SectionHeader>
            }
            open={openSections.charts}
            onTriggerOpening={() => {}}
            onTriggerClosing={() => {}}
            transitionTime={200}
          >
            <SectionContent>
              <ChartContainer>
                <MetricTitle>Service Resource Usage</MetricTitle>
                <ResponsiveContainer width="100%" height="100%">
                  <BarChart data={serviceMemoryData}>
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis 
                      dataKey="name" 
                      tick={{ fontSize: window.innerWidth < 768 ? 10 : 12 }}
                    />
                    <YAxis 
                      tick={{ fontSize: window.innerWidth < 768 ? 10 : 12 }}
                    />
                    <Tooltip 
                      contentStyle={{ 
                        fontSize: window.innerWidth < 768 ? '12px' : '14px',
                        backgroundColor: 'rgba(255, 255, 255, 0.9)',
                        border: '1px solid #ccc',
                        borderRadius: '4px'
                      }}
                    />
                    <Legend 
                      wrapperStyle={{ 
                        fontSize: window.innerWidth < 768 ? '10px' : '12px',
                        paddingTop: '10px'
                      }}
                    />
                    <Bar dataKey="memory" name="Memory Usage (%)" fill="#8884d8" />
                    <Bar dataKey="cpu" name="CPU Usage (%)" fill="#82ca9d" />
                  </BarChart>
                </ResponsiveContainer>
              </ChartContainer>

              <ChartContainer>
                <MetricTitle>Service Performance</MetricTitle>
                <ResponsiveContainer width="100%" height="100%">
                  <LineChart data={serviceResponseData}>
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis 
                      dataKey="name" 
                      tick={{ fontSize: window.innerWidth < 768 ? 10 : 12 }}
                    />
                    <YAxis 
                      yAxisId="left"
                      tick={{ fontSize: window.innerWidth < 768 ? 10 : 12 }}
                    />
                    <YAxis 
                      yAxisId="right" 
                      orientation="right"
                      tick={{ fontSize: window.innerWidth < 768 ? 10 : 12 }}
                    />
                    <Tooltip 
                      contentStyle={{ 
                        fontSize: window.innerWidth < 768 ? '12px' : '14px',
                        backgroundColor: 'rgba(255, 255, 255, 0.9)',
                        border: '1px solid #ccc',
                        borderRadius: '4px'
                      }}
                    />
                    <Legend 
                      wrapperStyle={{ 
                        fontSize: window.innerWidth < 768 ? '10px' : '12px',
                        paddingTop: '10px'
                      }}
                    />
                    <Line
                      yAxisId="left"
                      type="monotone"
                      dataKey="responseTime"
                      name="Response Time (ms)"
                      stroke="#8884d8"
                      activeDot={{ r: window.innerWidth < 768 ? 4 : 8 }}
                    />
                    <Line
                      yAxisId="right"
                      type="monotone"
                      dataKey="requests"
                      name="Total Requests"
                      stroke="#82ca9d"
                      activeDot={{ r: window.innerWidth < 768 ? 4 : 8 }}
                    />
                  </LineChart>
                </ResponsiveContainer>
              </ChartContainer>
            </SectionContent>
          </Collapsible>
          </CollapsibleSection>

          <CollapsibleSection>
            <Collapsible
              id="section-services"
            trigger={
              <SectionHeader onClick={() => toggleSection('services')}>
                <h2>Service Management</h2>
                <SectionIcon isOpen={openSections.services}>▼</SectionIcon>
              </SectionHeader>
            }
            open={openSections.services}
            onTriggerOpening={() => {}}
            onTriggerClosing={() => {}}
            transitionTime={200}
          >
            <SectionContent>
              <AnimatedContainer staggerDelay={reducedMotion ? 0 : 75}>
                <ServiceGrid>
                  {serviceMetrics.map((service, index) => (
                    <ServiceCard key={service.name}>
                    <ServiceHeader>
                      <ServiceTitle>{service.name}</ServiceTitle>
                      <ServiceStatus status={service.status.toLowerCase()}>
                        {service.status}
                      </ServiceStatus>
                    </ServiceHeader>
                    
                    <ServiceMetrics>
                      <MetricItem>
                        <ServiceMetricValue>{service.memoryUsage?.toFixed(2) ?? '0.00'}%</ServiceMetricValue>
                        <ServiceMetricLabel>Memory</ServiceMetricLabel>
                      </MetricItem>
                      <MetricItem>
                        <ServiceMetricValue>{service.cpuUsage?.toFixed(2) ?? '0.00'}%</ServiceMetricValue>
                        <ServiceMetricLabel>CPU</ServiceMetricLabel>
                      </MetricItem>
                      <MetricItem>
                        <ServiceMetricValue>{service.averageResponseTime?.toFixed(2) ?? '0.00'}ms</ServiceMetricValue>
                        <ServiceMetricLabel>Response Time</ServiceMetricLabel>
                      </MetricItem>
                      <MetricItem>
                        <ServiceMetricValue>{service.requestCount}</ServiceMetricValue>
                        <ServiceMetricLabel>Requests</ServiceMetricLabel>
                      </MetricItem>
                    </ServiceMetrics>
                    
                    <ServiceControls>
                      <ControlButton
                        variant="start"
                        onClick={() => controlService(service.name, 'start')}
                        disabled={service.status.toLowerCase() === 'running'}
                      >
                        Start
                      </ControlButton>
                      <ControlButton
                        variant="stop"
                        onClick={() => controlService(service.name, 'stop')}
                        disabled={service.status.toLowerCase() === 'stopped'}
                      >
                        Stop
                      </ControlButton>
                      <ControlButton
                        variant="restart"
                        onClick={() => controlService(service.name, 'restart')}
                        disabled={service.status.toLowerCase() === 'stopped'}
                      >
                        Restart
                      </ControlButton>
                      <ControlButton
                        variant="logs"
                        onClick={() => window.open(`/api/v1/devpanel/logs/${service.name}`, '_blank')}
                      >
                        Logs
                      </ControlButton>
                      <ControlButton
                        variant="config"
                        onClick={() => window.open(`/api/v1/devpanel/config/${service.name}`, '_blank')}
                      >
                        Config
                      </ControlButton>
                    </ServiceControls>
                  </ServiceCard>
                  ))}
                </ServiceGrid>
              </AnimatedContainer>
            </SectionContent>
          </Collapsible>
          </CollapsibleSection>

          <CollapsibleSection>
            <Collapsible
              id="section-visitorAnalytics"
            trigger={
              <SectionHeader onClick={() => toggleSection('visitorAnalytics')}>
                <h2>Visitor Analytics</h2>
                <SectionIcon isOpen={openSections.visitorAnalytics}>▼</SectionIcon>
              </SectionHeader>
            }
            open={openSections.visitorAnalytics}
            onTriggerOpening={() => {}}
            onTriggerClosing={() => {}}
            transitionTime={200}
          >
            <SectionContent>
              <VisitorAnalytics />
            </SectionContent>
          </Collapsible>
          </CollapsibleSection>

          <CollapsibleSection>
            <Collapsible
              id="section-projects"
            trigger={
              <SectionHeader onClick={() => toggleSection('projects')}>
                <h2>Project Management</h2>
                <SectionIcon isOpen={openSections.projects}>▼</SectionIcon>
              </SectionHeader>
            }
            open={openSections.projects}
            onTriggerOpening={() => {}}
            onTriggerClosing={() => {}}
            transitionTime={200}
          >
            <SectionContent>
              <ProjectManager />
            </SectionContent>
          </Collapsible>
          </CollapsibleSection>

          <CollapsibleSection>
            <Collapsible
              id="section-projectPaths"
            trigger={
              <SectionHeader onClick={() => toggleSection('projectPaths')}>
                <h2>Lines of Code Project Paths</h2>
                <SectionIcon isOpen={openSections.projectPaths}>▼</SectionIcon>
              </SectionHeader>
            }
            open={openSections.projectPaths}
            onTriggerOpening={() => {}}
            onTriggerClosing={() => {}}
            transitionTime={200}
          >
            <SectionContent>
              <ProjectPathsManager />
            </SectionContent>
          </Collapsible>
          </CollapsibleSection>

          <CollapsibleSection>
            <Collapsible
              id="section-certifications"
            trigger={
              <SectionHeader onClick={() => toggleSection('certifications')}>
                <h2>Certifications Management</h2>
                <SectionIcon isOpen={openSections.certifications}>▼</SectionIcon>
              </SectionHeader>
            }
            open={openSections.certifications}
            onTriggerOpening={() => {}}
            onTriggerClosing={() => {}}
            transitionTime={200}
          >
            <SectionContent>
              <CertificationsManager />
            </SectionContent>
          </Collapsible>
          </CollapsibleSection>

          <CollapsibleSection>
            <Collapsible
              id="section-prompts"
            trigger={
              <SectionHeader onClick={() => toggleSection('prompts')}>
                <h2>Prompts Management</h2>
                <SectionIcon isOpen={openSections.prompts}>▼</SectionIcon>
              </SectionHeader>
            }
            open={openSections.prompts}
            onTriggerOpening={() => {}}
            onTriggerClosing={() => {}}
            transitionTime={200}
          >
            <SectionContent>
              <PromptsManager />
            </SectionContent>
          </Collapsible>
          </CollapsibleSection>

          <RefreshButton
            onClick={() => fetchData(true)}
            disabled={isRefreshing}
            title={`Refresh data (Last updated: ${lastUpdateTime})`}
          >
            {isRefreshing ? '⟳' : '↻'}
          </RefreshButton>

          {isRefreshing && (
            <LoadingOverlay>
              <LoadingContent>
                <div style={{ display: 'flex', alignItems: 'center', gap: '1rem' }}>
                  <div style={{ 
                    width: '24px', 
                    height: '24px', 
                    border: '2px solid #ccc', 
                    borderTop: '2px solid #007bff', 
                    borderRadius: '50%', 
                    animation: 'spin 1s linear infinite' 
                  }}></div>
                  Refreshing data...
                </div>
              </LoadingContent>
            </LoadingOverlay>
          )}
        </>
      )}
      <ScrollToTopButton />
    </PageContainer>
  );
};

export default DevPanel; 