import React, { useState, useEffect } from 'react';
import styled from 'styled-components';
import {
  LineChart,
  Line,
  BarChart,
  Bar,
  PieChart,
  Pie,
  Cell,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
  Area,
  AreaChart
} from 'recharts';
import { api } from '../../utils/apiConfig';
import { handleError } from '../../utils/errorHandler';

// Types
interface VisitorAnalyticsData {
  today: {
    uniqueVisitors: number;
    totalPageViews: number;
  };
  last7Days: {
    uniqueVisitors: number;
    totalPageViews: number;
  };
  last30Days: {
    uniqueVisitors: number;
    totalPageViews: number;
  };
  allTime: {
    uniqueVisitors: number;
    totalPageViews: number;
  };
  realTimeCount: number;
}

interface RealtimeData {
  type: string;
  path: string;
  realtime: number;
}

interface MetricsSummary {
  uniqueVisitors: number;
  totalPageViews: number;
  avgSessionDuration: number;
  bounceRate: number;
  newVisitors: number;
  returningVisitors: number;
}

interface TrendData {
  todayVsYesterday: string;
  thisWeekVsLastWeek: string;
  thisMonthVsLastMonth: string;
}

interface TimelineData {
  timestamp: string;
  visitors: number;
  pageViews: number;
  avgSessionTime: number;
}

interface LocationData {
  countryCode: string;
  countryName: string;
  visitorCount: number;
  percentage: number;
}

interface DeviceBreakdown {
  period: string;
  devices: { [key: string]: number };
  browsers: { [key: string]: number };
  os: { [key: string]: number };
}

interface VisitorStats {
  today: MetricsSummary;
  last7Days: MetricsSummary;
  last30Days: MetricsSummary;
  allTime: MetricsSummary;
  trends: TrendData;
  realTimeCount: number;
}

// Styled Components
const Container = styled.div`
  padding: 1rem;
`;

const Header = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
  flex-wrap: wrap;
  gap: 1rem;
`;

const Title = styled.h2`
  font-size: 1.75rem;
  color: ${({ theme }) => theme.colors.text};
  margin: 0;
`;

const TimeRangeSelector = styled.div`
  display: flex;
  gap: 0.5rem;
  background: ${({ theme }) => theme.colors.background};
  padding: 0.25rem;
  border-radius: 8px;
  border: 1px solid ${({ theme }) => theme.colors.border};
`;

const TimeRangeButton = styled.button<{ active: boolean }>`
  padding: 0.5rem 1rem;
  border: none;
  background: ${({ active, theme }) => active ? theme.colors.primary : 'transparent'};
  color: ${({ active, theme }) => active ? 'white' : theme.colors.text};
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.875rem;
  font-weight: 500;
  transition: all 0.2s ease;

  &:hover {
    background: ${({ active, theme }) => active ? theme.colors.primary : theme.colors.card};
  }
`;

const MetricsGrid = styled.div`
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
`;

const MetricCard = styled.div<{ trend?: 'up' | 'down' | 'neutral' }>`
  background: ${({ theme }) => theme.colors.card};
  border-radius: 12px;
  padding: 1.5rem;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
  border: 1px solid ${({ theme }) => theme.colors.border};
  position: relative;
  overflow: hidden;

  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 4px;
    background: ${({ trend, theme }) => 
      trend === 'up' ? theme.colors.success :
      trend === 'down' ? theme.colors.error :
      theme.colors.primary
    };
  }
`;

const MetricLabel = styled.div`
  font-size: 0.875rem;
  color: ${({ theme }) => theme.colors.textSecondary};
  margin-bottom: 0.5rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
`;

const MetricValue = styled.div`
  font-size: 2rem;
  font-weight: 600;
  color: ${({ theme }) => theme.colors.text};
  margin-bottom: 0.5rem;
`;

const MetricTrend = styled.div<{ positive: boolean }>`
  font-size: 0.875rem;
  color: ${({ positive, theme }) => positive ? theme.colors.success : theme.colors.error};
  display: flex;
  align-items: center;
  gap: 0.25rem;
`;

const ChartContainer = styled.div`
  background: ${({ theme }) => theme.colors.card};
  border-radius: 12px;
  padding: 1.5rem;
  margin-bottom: 2rem;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
  border: 1px solid ${({ theme }) => theme.colors.border};
`;

const ChartTitle = styled.h3`
  font-size: 1.25rem;
  color: ${({ theme }) => theme.colors.text};
  margin: 0 0 1.5rem 0;
`;

const LiveIndicator = styled.div`
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  background: ${({ theme }) => theme.colors.success}20;
  color: ${({ theme }) => theme.colors.success};
  padding: 0.25rem 0.75rem;
  border-radius: 20px;
  font-size: 0.875rem;
  font-weight: 500;

  &::before {
    content: '';
    width: 8px;
    height: 8px;
    background: ${({ theme }) => theme.colors.success};
    border-radius: 50%;
    animation: pulse 2s infinite;
  }

  @keyframes pulse {
    0% {
      box-shadow: 0 0 0 0 ${({ theme }) => theme.colors.success}40;
    }
    70% {
      box-shadow: 0 0 0 10px ${({ theme }) => theme.colors.success}00;
    }
    100% {
      box-shadow: 0 0 0 0 ${({ theme }) => theme.colors.success}00;
    }
  }
`;

const LocationsGrid = styled.div`
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
  margin-top: 1rem;
`;

const LocationItem = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.75rem;
  background: ${({ theme }) => theme.colors.background};
  border-radius: 8px;
`;

const LoadingState = styled.div`
  text-align: center;
  padding: 4rem 2rem;
  color: ${({ theme }) => theme.colors.textSecondary};
`;

// Color palette for charts
const COLORS = ['#3B82F6', '#10B981', '#F59E0B', '#EF4444', '#8B5CF6', '#EC4899', '#14B8A6', '#F97316'];

// Component
const VisitorAnalytics: React.FC = () => {
  const [state, setState] = useState({
    timeRange: 'today' as 'today' | 'last7Days' | 'last30Days' | 'allTime',
    visitorStats: null as VisitorStats | null,
    timelineData: [] as TimelineData[],
    locationData: [] as LocationData[],
    deviceData: null as DeviceBreakdown | null,
    realtimeData: [] as RealtimeData[],
    loading: true,
    error: null as string | null,
  });

  useEffect(() => {
    const fetchVisitorData = async () => {
      try {
        setState(prevState => ({ ...prevState, loading: true }));
        const [stats, locations, devices] = await Promise.all([
          api.get<VisitorStats>('/api/v1/devpanel/visitors/stats'),
          api.get<{ locations: LocationData[] }>('/api/v1/devpanel/visitors/locations'),
          api.get<DeviceBreakdown>('/api/v1/devpanel/visitors/breakdown')
        ]);

        setState(prevState => ({
          ...prevState,
          visitorStats: stats,
          locationData: locations.locations,
          deviceData: devices,
          error: null,
        }));
      } catch (err) {
        handleError(err, { context: 'VisitorAnalytics.fetchVisitorData' });
        setState(prevState => ({ ...prevState, error: 'Failed to load visitor data' }));
      } finally {
        setState(prevState => ({ ...prevState, loading: false }));
      }
    };

    const fetchTimelineData = async () => {
      try {
        const period = state.timeRange === 'today' ? '1d' : 
                       state.timeRange === 'last7Days' ? '7d' : 
                       state.timeRange === 'last30Days' ? '30d' : '1y';
        
        const response = await api.get<{ data: TimelineData[] }>(
          `/api/v1/devpanel/visitors/timeline?period=${period}`
        );
        
        setState(prevState => ({ ...prevState, timelineData: response.data }));
      } catch (err) {
        handleError(err, { context: 'VisitorAnalytics.fetchTimelineData' });
      }
    };

    fetchVisitorData();
    fetchTimelineData();

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const ws = new WebSocket(`${protocol}//${window.location.host}/ws/analytics`);

    ws.onmessage = (event) => {
      const message = JSON.parse(event.data);
      if (message.type === 'pageview') {
        setState(prevState => ({ ...prevState, realtimeData: [...prevState.realtimeData, message] }));
      }
    };

    return () => {
      ws.close();
    };
  }, [state.timeRange]);

  const { timeRange, visitorStats, timelineData, locationData, deviceData, loading, error } = state;

  const formatDuration = (seconds: number): string => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins}m ${secs}s`;
  };

  const getTrendDirection = (trend: string): 'up' | 'down' | 'neutral' => {
    if (trend.startsWith('+')) return 'up';
    if (trend.startsWith('-')) return 'down';
    return 'neutral';
  };

  const getTrendValue = (trend: string): boolean => {
    return !trend.startsWith('-');
  };

  if (loading) {
    return <LoadingState>Loading visitor analytics...</LoadingState>;
  }

  if (error) {
    return <LoadingState>{error}</LoadingState>;
  }

  if (!visitorStats) {
    return <LoadingState>No data available</LoadingState>;
  }

  const currentMetrics = visitorStats[timeRange];
  const currentTrend = timeRange === 'today' ? visitorStats.trends.todayVsYesterday :
                       timeRange === 'last7Days' ? visitorStats.trends.thisWeekVsLastWeek :
                       visitorStats.trends.thisMonthVsLastMonth;

  // Prepare data for pie charts
  const deviceChartData = deviceData ? Object.entries(deviceData.devices).map(([name, value]) => ({
    name: name ? name.charAt(0).toUpperCase() + name.slice(1) : 'Unknown',
    value
  })) : [];

  const browserChartData = deviceData ? Object.entries(deviceData.browsers).map(([name, value]) => ({
    name: name ? name.charAt(0).toUpperCase() + name.slice(1) : 'Unknown',
    value
  })) : [];

  return (
    <Container>
      <Header>
        <Title>Visitor Analytics</Title>
        <TimeRangeSelector>
          <TimeRangeButton
            active={timeRange === 'today'}
            onClick={() => setState(s => ({...s, timeRange: 'today'}))}
          >
            Today
          </TimeRangeButton>
          <TimeRangeButton
            active={timeRange === 'last7Days'}
            onClick={() => setState(s => ({...s, timeRange: 'last7Days'}))}
          >
            Last 7 Days
          </TimeRangeButton>
          <TimeRangeButton
            active={timeRange === 'last30Days'}
            onClick={() => setState(s => ({...s, timeRange: 'last30Days'}))}
          >
            Last 30 Days
          </TimeRangeButton>
          <TimeRangeButton
            active={timeRange === 'allTime'}
            onClick={() => setState(s => ({...s, timeRange: 'allTime'}))}
          >
            All Time
          </TimeRangeButton>
        </TimeRangeSelector>
      </Header>

      <MetricsGrid>
        <MetricCard trend={getTrendDirection(currentTrend)}>
          <MetricLabel>Unique Visitors</MetricLabel>
          <MetricValue>{currentMetrics.uniqueVisitors.toLocaleString()}</MetricValue>
          <MetricTrend positive={getTrendValue(currentTrend)}>
            {currentTrend}
          </MetricTrend>
        </MetricCard>

        <MetricCard>
          <MetricLabel>Page Views</MetricLabel>
          <MetricValue>{currentMetrics.totalPageViews.toLocaleString()}</MetricValue>
          <MetricTrend positive={true}>
            {(currentMetrics.uniqueVisitors > 0 ? currentMetrics.totalPageViews / currentMetrics.uniqueVisitors : 0).toFixed(1)} per visitor
          </MetricTrend>
        </MetricCard>

        <MetricCard>
          <MetricLabel>Avg. Session Duration</MetricLabel>
          <MetricValue>{formatDuration(currentMetrics.avgSessionDuration)}</MetricValue>
        </MetricCard>

        <MetricCard>
          <MetricLabel>Bounce Rate</MetricLabel>
          <MetricValue>{currentMetrics.bounceRate.toFixed(1)}%</MetricValue>
          <MetricTrend positive={currentMetrics.bounceRate < 50}>
            {currentMetrics.bounceRate < 50 ? 'Good' : 'Needs improvement'}
          </MetricTrend>
        </MetricCard>

        <MetricCard>
          <MetricLabel>New Visitors</MetricLabel>
          <MetricValue>{currentMetrics.newVisitors.toLocaleString()}</MetricValue>
          <MetricTrend positive={true}>
            {(currentMetrics.uniqueVisitors > 0 ? (currentMetrics.newVisitors / currentMetrics.uniqueVisitors) * 100 : 0).toFixed(1)}%
          </MetricTrend>
        </MetricCard>

        <MetricCard>
          <MetricLabel>Real-Time Visitors</MetricLabel>
          <MetricValue>{visitorStats.realTimeCount}</MetricValue>
          <LiveIndicator>Live</LiveIndicator>
        </MetricCard>
      </MetricsGrid>

      <ChartContainer>
        <ChartTitle>Visitor Timeline</ChartTitle>
        <ResponsiveContainer width="100%" height={300}>
          <AreaChart data={timelineData}>
            <CartesianGrid strokeDasharray="3 3" stroke="#e0e0e0" />
            <XAxis 
              dataKey="timestamp" 
              tick={{ fontSize: 12 }}
              tickFormatter={(value) => {
                const date = new Date(value);
                return timeRange === 'today' ? 
                  date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }) :
                  date.toLocaleDateString([], { month: 'short', day: 'numeric' });
              }}
            />
            <YAxis tick={{ fontSize: 12 }} />
            <Tooltip 
              contentStyle={{ 
                backgroundColor: 'rgba(255, 255, 255, 0.95)',
                border: '1px solid #ccc',
                borderRadius: '4px'
              }}
            />
            <Legend />
            <Area
              type="monotone"
              dataKey="visitors"
              stroke="#3B82F6"
              fill="#3B82F6"
              fillOpacity={0.6}
              name="Unique Visitors"
            />
            <Area
              type="monotone"
              dataKey="pageViews"
              stroke="#10B981"
              fill="#10B981"
              fillOpacity={0.6}
              name="Page Views"
            />
          </AreaChart>
        </ResponsiveContainer>
      </ChartContainer>

      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(300px, 1fr))', gap: '1.5rem' }}>
        <ChartContainer>
          <ChartTitle>Top Countries</ChartTitle>
          <LocationsGrid>
            {locationData && locationData.slice(0, 8).map((location) => (
              <LocationItem key={location.countryCode}>
                <div>
                  <div style={{ fontWeight: 500 }}>{location.countryName}</div>
                  <div style={{ fontSize: '0.875rem', color: '#666' }}>
                    {location.visitorCount.toLocaleString()} visitors
                  </div>
                </div>
                <div style={{ fontWeight: 600, color: '#3B82F6' }}>
                  {location.percentage.toFixed(1)}%
                </div>
              </LocationItem>
            ))}
          </LocationsGrid>
        </ChartContainer>

        <ChartContainer>
          <ChartTitle>Device Distribution</ChartTitle>
          <ResponsiveContainer width="100%" height={250}>
            <PieChart>
              <Pie
                data={deviceChartData}
                cx="50%"
                cy="50%"
                outerRadius={80}
                fill="#8884d8"
                dataKey="value"
                label={({ name, percent }) => `${name} ${(percent * 100).toFixed(0)}%`}
              >
                {deviceChartData.map((entry, index) => (
                  <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                ))}
              </Pie>
              <Tooltip />
            </PieChart>
          </ResponsiveContainer>
        </ChartContainer>

        <ChartContainer>
          <ChartTitle>Browser Distribution</ChartTitle>
          <ResponsiveContainer width="100%" height={250}>
            <BarChart data={browserChartData}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="name" tick={{ fontSize: 12 }} />
              <YAxis tick={{ fontSize: 12 }} />
              <Tooltip />
              <Bar dataKey="value" fill="#8B5CF6" />
            </BarChart>
          </ResponsiveContainer>
        </ChartContainer>
      </div>
    </Container>
  );
};

export default VisitorAnalytics;
