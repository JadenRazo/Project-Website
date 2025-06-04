import React, { useState, useEffect, useCallback } from 'react';
import styled from 'styled-components';
import Collapsible from 'react-collapsible';
import ProjectManager from '../../components/devpanel/ProjectManager';
import AdminLogin from '../../components/devpanel/AdminLogin';
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
  max-width: 1200px;
  margin: 0 auto;
  
  @media (max-width: ${({ theme }) => theme.breakpoints.mobile}) {
    padding: 1rem;
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
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
  
  @media (max-width: ${({ theme }) => theme.breakpoints.mobile}) {
    grid-template-columns: 1fr;
    gap: 1rem;
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
  
  @media (max-width: ${({ theme }) => theme.breakpoints.mobile}) {
    height: 300px;
    padding: 1rem;
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

const LoadingSpinner = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 200px;
`;

// Service-specific styled components
const ServiceGrid = styled.div`
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
  
  @media (max-width: ${({ theme }) => theme.breakpoints.tablet}) {
    grid-template-columns: 1fr;
    gap: 1rem;
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
  
  @media (max-width: ${({ theme }) => theme.breakpoints.mobile}) {
    flex-wrap: wrap;
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
  z-index: 1000;
  
  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 6px 8px rgba(0, 0, 0, 0.15);
  }
  
  &:active {
    transform: translateY(0);
  }
  
  @media (max-width: ${({ theme }) => theme.breakpoints.mobile}) {
    bottom: 1rem;
    right: 1rem;
    padding: 0.75rem;
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
  z-index: 1000;
`;

const LoadingContent = styled.div`
  background: ${({ theme }) => theme.colors.card};
  padding: 2rem;
  border-radius: 8px;
  text-align: center;
  color: ${({ theme }) => theme.colors.text};
`;

// Main Component
const DevPanel: React.FC = () => {
  const [adminUser, setAdminUser] = useState<any>(null);
  const [isAdminAuthenticated, setIsAdminAuthenticated] = useState(false);
  const [checkingAuth, setCheckingAuth] = useState(true);
  const [systemMetrics, setSystemMetrics] = useState<SystemMetrics | null>(null);
  const [serviceMetrics, setServiceMetrics] = useState<ServiceMetricsData[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [controlError, setControlError] = useState<string | null>(null);
  const [isRefreshing, setIsRefreshing] = useState<boolean>(false);
  const [openSections, setOpenSections] = useState<{
    system: boolean;
    charts: boolean;
    services: boolean;
    projects: boolean;
  }>({
    system: true,
    charts: true,
    services: true,
    projects: true,
  });

  const toggleSection = (section: keyof typeof openSections) => {
    setOpenSections(prev => ({
      ...prev,
      [section]: !prev[section],
    }));
  };

  const fetchData = useCallback(async () => {
    const token = localStorage.getItem('auth_token');
    if (!token) {
      setError('Authentication required');
      return;
    }
    
      try {
        setLoading(true);
      setIsRefreshing(true);
      const headers = {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      };

      // Fetch system metrics
      const systemResponse = await fetch('/api/v1/devpanel/system', { headers });
      if (!systemResponse.ok) throw new Error('Failed to fetch system metrics');
      const systemData = await systemResponse.json();
      setSystemMetrics(systemData);

      // Fetch service metrics
      const servicesResponse = await fetch('/api/v1/devpanel/services', { headers });
      if (!servicesResponse.ok) throw new Error('Failed to fetch service metrics');
      const servicesData = await servicesResponse.json();
      setServiceMetrics(servicesData);


    } catch (err) {
      console.error('Error fetching data:', err);
      setError('Failed to load data. Please try again later.');
    } finally {
      setLoading(false);
      setIsRefreshing(false);
    }
  }, []);

  const controlService = async (serviceName: string, action: 'start' | 'stop' | 'restart') => {
    const token = localStorage.getItem('auth_token');
    if (!token) {
      setControlError('Authentication required');
      return;
    }
        
    try {
      setControlError(null);
      const headers = {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      };

      const response = await fetch(`/api/v1/devpanel/services/${serviceName}/${action}`, {
        method: 'POST',
        headers,
      });

      if (!response.ok) throw new Error(`Failed to ${action} service`);
      
      const data: ServiceControlResponse = await response.json();
      if (!data.success) throw new Error(data.message);

      // Refresh data after successful control action
      await fetchData();
    } catch (error: unknown) {
      console.error(`Error controlling service:`, error);
      const errorMessage = error instanceof Error ? error.message : 'Unknown error occurred';
      setControlError(`Failed to ${action} service: ${errorMessage}`);
    }
  };

  useEffect(() => {
    checkAdminAuth();
  }, []);

  useEffect(() => {
    if (isAdminAuthenticated && adminUser) {
      fetchData();
      const interval = setInterval(fetchData, 30000);
      return () => clearInterval(interval);
    }
  }, [isAdminAuthenticated, adminUser, fetchData]);

  const checkAdminAuth = async () => {
    const token = localStorage.getItem('auth_token');
    if (!token) {
      setCheckingAuth(false);
      return;
    }

    try {
      const response = await fetch('/api/v1/auth/admin/validate', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
      });

      if (response.ok) {
        const userData = await response.json();
        setAdminUser(userData);
        setIsAdminAuthenticated(true);
      } else {
        localStorage.removeItem('auth_token');
      }
    } catch (error) {
      console.error('Auth check failed:', error);
      localStorage.removeItem('auth_token');
    } finally {
      setCheckingAuth(false);
    }
  };

  const handleLoginSuccess = (userData: any) => {
    setAdminUser(userData);
    setIsAdminAuthenticated(true);
  };

  const handleLogout = () => {
    localStorage.removeItem('auth_token');
    setAdminUser(null);
    setIsAdminAuthenticated(false);
  };

  if (checkingAuth) {
    return (
      <PageContainer>
        <LoadingSpinner>Checking authentication...</LoadingSpinner>
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
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <div>
            <Title>Developer Panel</Title>
            <Subtitle>System monitoring and service management</Subtitle>
          </div>
          <div style={{ display: 'flex', alignItems: 'center', gap: '1rem' }}>
            <span style={{ color: 'var(--color-text-secondary)', fontSize: '0.9rem' }}>
              Welcome, {adminUser?.email}
            </span>
            <ControlButton 
              variant="config" 
              onClick={handleLogout}
              style={{ padding: '0.5rem 1rem', fontSize: '0.9rem' }}
            >
              Logout
            </ControlButton>
          </div>
        </div>
      </Header>

      {error && <ErrorMessage>{error}</ErrorMessage>}
      {controlError && <ErrorMessage>{controlError}</ErrorMessage>}

      {loading ? (
        <LoadingSpinner>Loading...</LoadingSpinner>
      ) : (
        <>
          <Collapsible
            trigger={
              <SectionHeader>
                <h2>System Overview</h2>
                <SectionIcon isOpen={openSections.system}>▼</SectionIcon>
              </SectionHeader>
            }
            open={openSections.system}
            onOpen={() => toggleSection('system')}
            transitionTime={200}
          >
            <SectionContent>
              <MetricsGrid>
                {systemMetrics && (
                  <>
                    <MetricCard>
                      <MetricTitle>System Memory Usage</MetricTitle>
                      <SystemMetricValue>{systemMetrics.memoryUsage.toFixed(2)}%</SystemMetricValue>
                      <SystemMetricLabel>Total system memory utilization</SystemMetricLabel>
                    </MetricCard>
                    <MetricCard>
                      <MetricTitle>CPU Usage</MetricTitle>
                      <SystemMetricValue>{systemMetrics.cpuUsage.toFixed(2)}%</SystemMetricValue>
                      <SystemMetricLabel>Total CPU utilization</SystemMetricLabel>
                    </MetricCard>
                    <MetricCard>
                      <MetricTitle>Disk Usage</MetricTitle>
                      <SystemMetricValue>{systemMetrics.diskUsage.toFixed(2)}%</SystemMetricValue>
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

          <Collapsible
            trigger={
              <SectionHeader>
                <h2>Service Analytics</h2>
                <SectionIcon isOpen={openSections.charts}>▼</SectionIcon>
              </SectionHeader>
            }
            open={openSections.charts}
            onOpen={() => toggleSection('charts')}
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

          <Collapsible
            trigger={
              <SectionHeader>
                <h2>Service Management</h2>
                <SectionIcon isOpen={openSections.services}>▼</SectionIcon>
              </SectionHeader>
            }
            open={openSections.services}
            onOpen={() => toggleSection('services')}
            transitionTime={200}
          >
            <SectionContent>
              <ServiceGrid>
                {serviceMetrics.map(service => (
                  <ServiceCard key={service.name}>
                    <ServiceHeader>
                      <ServiceTitle>{service.name}</ServiceTitle>
                      <ServiceStatus status={service.status.toLowerCase()}>
                        {service.status}
                      </ServiceStatus>
                    </ServiceHeader>
                    
                    <ServiceMetrics>
                      <MetricItem>
                        <ServiceMetricValue>{service.memoryUsage.toFixed(2)}%</ServiceMetricValue>
                        <ServiceMetricLabel>Memory</ServiceMetricLabel>
                      </MetricItem>
                      <MetricItem>
                        <ServiceMetricValue>{service.cpuUsage.toFixed(2)}%</ServiceMetricValue>
                        <ServiceMetricLabel>CPU</ServiceMetricLabel>
                      </MetricItem>
                      <MetricItem>
                        <ServiceMetricValue>{service.averageResponseTime.toFixed(2)}ms</ServiceMetricValue>
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
            </SectionContent>
          </Collapsible>

          <Collapsible
            trigger={
              <SectionHeader>
                <h2>Project Management</h2>
                <SectionIcon isOpen={openSections.projects}>▼</SectionIcon>
              </SectionHeader>
            }
            open={openSections.projects}
            onOpen={() => toggleSection('projects')}
            transitionTime={200}
          >
            <SectionContent>
              <ProjectManager />
            </SectionContent>
          </Collapsible>

          <RefreshButton
            onClick={fetchData}
            disabled={isRefreshing}
            title="Refresh data"
          >
            {isRefreshing ? '⟳' : '↻'}
          </RefreshButton>

          {isRefreshing && (
            <LoadingOverlay>
              <LoadingContent>
                Refreshing data...
              </LoadingContent>
            </LoadingOverlay>
          )}
        </>
      )}
    </PageContainer>
  );
};

export default DevPanel; 