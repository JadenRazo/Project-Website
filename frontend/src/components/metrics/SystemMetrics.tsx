import React, { useState, useEffect, useCallback, useMemo } from 'react';
import styled, { keyframes, css } from 'styled-components';
import { Line } from 'react-chartjs-2';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler,
} from 'chart.js';

// Register Chart.js components
ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler
);

interface LatencyMetric {
  timestamp: string;
  latency: number;
}

interface LatencyStats {
  average: number;
  min: number;
  max: number;
  count: number;
  period: string;
}

interface SystemMetricsData {
  period: string;
  has_sufficient_data: boolean;
  metrics: LatencyMetric[];
  stats: LatencyStats;
  last_updated: string;
  message?: string;
}

type TimePeriod = 'day' | 'week' | 'month';

// Animations
const pulse = keyframes`
  0% { opacity: 1; }
  50% { opacity: 0.7; }
  100% { opacity: 1; }
`;

const slideIn = keyframes`
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
`;

const shimmer = keyframes`
  0% {
    background-position: -468px 0;
  }
  100% {
    background-position: 468px 0;
  }
`;

// Styled Components
const MetricsContainer = styled.div`
  background: linear-gradient(135deg, ${props => props.theme.colors.surface} 0%, ${props => props.theme.colors.background} 100%);
  border-radius: 16px;
  padding: 24px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  border: 1px solid ${props => props.theme.colors.border};
  animation: ${slideIn} 0.6s ease-out;
  transition: all 0.3s ease;

  &:hover {
    box-shadow: 0 12px 48px rgba(0, 0, 0, 0.15);
    transform: translateY(-2px);
  }
`;

const Header = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  flex-wrap: wrap;
  gap: 16px;

  @media (max-width: 768px) {
    flex-direction: column;
    align-items: stretch;
  }
`;

const MetricsTitle = styled.h2`
  color: ${props => props.theme.colors.text};
  font-size: 24px;
  font-weight: 700;
  margin: 0;
  display: flex;
  align-items: center;
  gap: 12px;

  &::before {
    content: '📊';
    font-size: 28px;
  }
`;

const PeriodButtons = styled.div`
  display: flex;
  gap: 8px;
  background: ${props => props.theme.colors.background};
  padding: 4px;
  border-radius: 12px;
  border: 1px solid ${props => props.theme.colors.border};

  @media (max-width: 768px) {
    justify-content: center;
  }
`;

const PeriodButton = styled.button<{ $active: boolean }>`
  padding: 10px 20px;
  border: none;
  border-radius: 8px;
  background: ${props => props.$active 
    ? `linear-gradient(135deg, ${props.theme.colors.primary} 0%, ${props.theme.colors.accent} 100%)` 
    : 'transparent'};
  color: ${props => props.$active ? 'white' : props.theme.colors.text};
  font-weight: ${props => props.$active ? '600' : '500'};
  cursor: pointer;
  transition: all 0.2s ease;
  font-size: 14px;
  position: relative;
  overflow: hidden;

  &:hover {
    background: ${props => props.$active 
      ? `linear-gradient(135deg, ${props.theme.colors.primary} 0%, ${props.theme.colors.accent} 100%)` 
      : props.theme.colors.surface};
    transform: translateY(-1px);
  }

  &:active {
    transform: translateY(0);
  }

  &::after {
    content: '';
    position: absolute;
    top: 50%;
    left: 50%;
    width: 0;
    height: 0;
    border-radius: 50%;
    background: rgba(255, 255, 255, 0.3);
    transition: width 0.3s ease, height 0.3s ease;
    transform: translate(-50%, -50%);
  }

  &:active::after {
    width: 120px;
    height: 120px;
  }
`;

const ChartContainer = styled.div<{ $hasData: boolean }>`
  height: 400px;
  margin-bottom: 24px;
  border-radius: 12px;
  overflow: hidden;
  background: ${props => props.theme.colors.background};
  border: 1px solid ${props => props.theme.colors.border};
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;

  ${props => !props.$hasData && css`
    background: linear-gradient(90deg, 
      ${props.theme.colors.surface} 25%, 
      rgba(255, 255, 255, 0.1) 50%, 
      ${props.theme.colors.surface} 75%
    );
    background-size: 200% 100%;
    animation: ${shimmer} 2s infinite linear;
  `}

  canvas {
    max-height: 100%;
  }

  @media (max-width: 768px) {
    height: 300px;
  }
`;

const NoDataMessage = styled.div`
  text-align: center;
  padding: 40px 20px;
  color: ${props => props.theme.colors.textSecondary};
  background: rgba(255, 255, 255, 0.05);
  border-radius: 12px;
  border: 2px dashed ${props => props.theme.colors.border};
  animation: ${pulse} 2s infinite ease-in-out;

  h3 {
    color: ${props => props.theme.colors.text};
    margin-bottom: 16px;
    font-size: 20px;
    font-weight: 600;
  }

  p {
    line-height: 1.6;
    font-size: 16px;
    max-width: 500px;
    margin: 0 auto;
  }

  .icon {
    font-size: 48px;
    margin-bottom: 16px;
    display: block;
  }
`;

const StatsGrid = styled.div`
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(120px, 1fr));
  gap: 16px;
  margin-top: 16px;

  @media (max-width: 768px) {
    grid-template-columns: repeat(2, 1fr);
  }
`;

const StatCard = styled.div`
  background: linear-gradient(135deg, ${props => props.theme.colors.surface} 0%, rgba(255, 255, 255, 0.05) 100%);
  padding: 16px;
  border-radius: 12px;
  border: 1px solid ${props => props.theme.colors.border};
  text-align: center;
  transition: all 0.2s ease;

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
  }

  .label {
    color: ${props => props.theme.colors.textSecondary};
    font-size: 12px;
    font-weight: 500;
    text-transform: uppercase;
    letter-spacing: 0.5px;
    margin-bottom: 8px;
  }

  .value {
    color: ${props => props.theme.colors.text};
    font-size: 20px;
    font-weight: 700;
  }

  .unit {
    color: ${props => props.theme.colors.textSecondary};
    font-size: 14px;
    font-weight: 500;
    margin-left: 4px;
  }
`;

const LastUpdated = styled.div`
  text-align: center;
  color: ${props => props.theme.colors.textSecondary};
  font-size: 12px;
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid ${props => props.theme.colors.border};

  .time {
    font-weight: 600;
    color: ${props => props.theme.colors.text};
  }
`;

const HoverInfo = styled.div<{ $x: number; $y: number; $visible: boolean }>`
  position: fixed;
  left: ${props => Math.max(10, Math.min(props.$x, window.innerWidth - 200))}px;
  top: ${props => Math.max(10, props.$y - 80)}px;
  background: ${props => props.theme.colors.surface};
  color: ${props => props.theme.colors.text};
  padding: 12px 16px;
  border-radius: 8px;
  border: 1px solid ${props => props.theme.colors.border};
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.2);
  pointer-events: none;
  z-index: 10000;
  opacity: ${props => props.$visible ? 1 : 0};
  transition: opacity 0.2s ease;
  font-size: 14px;
  font-weight: 500;
  white-space: nowrap;
  min-width: 180px;

  .date {
    color: ${props => props.theme.colors.textSecondary};
    font-size: 12px;
    margin-bottom: 4px;
  }

  .latency {
    color: ${props => props.theme.colors.primary};
    font-weight: 700;
  }
`;

const SystemMetrics: React.FC = () => {
  const [selectedPeriod, setSelectedPeriod] = useState<TimePeriod>('day');
  const [metricsData, setMetricsData] = useState<SystemMetricsData | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [hoverInfo, setHoverInfo] = useState<{
    x: number;
    y: number;
    visible: boolean;
    date: string;
    latency: number;
  }>({
    x: 0,
    y: 0,
    visible: false,
    date: '',
    latency: 0,
  });

  const fetchMetrics = useCallback(async (period: TimePeriod) => {
    try {
      setLoading(true);
      setError(null);
      
      const response = await fetch(`/api/v1/status/metrics/${period}`);
      if (!response.ok) {
        throw new Error(`Failed to fetch metrics: ${response.statusText}`);
      }
      
      const data: SystemMetricsData = await response.json();
      
      // Debug logging for development
      if (process.env.NODE_ENV === 'development') {
        console.log(`[SystemMetrics] Fetched ${period} metrics:`, {
          period: data.period,
          hasData: data.has_sufficient_data,
          metricCount: data.metrics?.length || 0,
          sampleTimestamps: data.metrics?.slice(0, 3).map(m => m.timestamp) || [],
        });
      }
      
      setMetricsData(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch metrics');
      console.error('Error fetching metrics:', err);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchMetrics(selectedPeriod);
    
    // Auto-refresh every minute
    const interval = setInterval(() => {
      fetchMetrics(selectedPeriod);
    }, 60000);
    
    return () => clearInterval(interval);
  }, [selectedPeriod, fetchMetrics]);

  const sortedMetrics = useMemo(() => {
    if (!metricsData?.metrics?.length) {
      return [];
    }
    // Sort metrics by timestamp to ensure chronological order
    const sorted = [...metricsData.metrics].sort((a, b) => 
      new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime()
    );
    
    // Debug logging for development
    if (process.env.NODE_ENV === 'development') {
      console.log(`[SystemMetrics] Sorted metrics for ${selectedPeriod}:`, {
        originalCount: metricsData.metrics.length,
        sortedCount: sorted.length,
        firstTimestamp: sorted[0]?.timestamp,
        lastTimestamp: sorted[sorted.length - 1]?.timestamp,
        sampleLabels: sorted.slice(0, 5).map(m => {
          const date = new Date(m.timestamp);
          return selectedPeriod === 'day' 
            ? date.toLocaleTimeString('en-US', { hour: 'numeric', minute: '2-digit', hour12: true })
            : date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
        }),
      });
    }
    
    return sorted;
  }, [metricsData?.metrics, selectedPeriod]);

  const chartData = useMemo(() => {
    if (!metricsData?.has_sufficient_data || !sortedMetrics.length) {
      return null;
    }

    const labels = sortedMetrics.map(metric => {
      const date = new Date(metric.timestamp);
      switch (selectedPeriod) {
        case 'day':
          return date.toLocaleTimeString('en-US', { 
            hour: 'numeric', 
            minute: '2-digit',
            hour12: true 
          });
        case 'week':
          return date.toLocaleDateString('en-US', { 
            month: 'short', 
            day: 'numeric'
          }) + ' ' + date.toLocaleTimeString('en-US', { 
            hour: 'numeric',
            hour12: true 
          });
        case 'month':
          return date.toLocaleDateString('en-US', { 
            month: 'short', 
            day: 'numeric' 
          });
        default:
          return date.toLocaleTimeString('en-US', { 
            hour: 'numeric', 
            minute: '2-digit',
            hour12: true 
          });
      }
    });

    const data = sortedMetrics.map(metric => metric.latency);

    return {
      labels,
      datasets: [
        {
          label: 'Latency (ms)',
          data,
          borderColor: '#3b82f6',
          backgroundColor: 'rgba(59, 130, 246, 0.1)',
          borderWidth: 2,
          fill: true,
          tension: 0.4,
          pointBackgroundColor: '#3b82f6',
          pointBorderColor: '#ffffff',
          pointBorderWidth: 2,
          pointRadius: 4,
          pointHoverRadius: 6,
          pointHoverBackgroundColor: '#1d4ed8',
          pointHoverBorderColor: '#ffffff',
          pointHoverBorderWidth: 3,
        },
      ],
    };
  }, [sortedMetrics, selectedPeriod, metricsData?.has_sufficient_data]);

  const chartOptions = useMemo(() => ({
    responsive: true,
    maintainAspectRatio: false,
    interaction: {
      intersect: false,
      mode: 'index' as const,
    },
    plugins: {
      legend: {
        display: false,
      },
      tooltip: {
        enabled: false,
        external: (context: any) => {
          const { chart, tooltip } = context;
          
          if (tooltip.opacity === 0) {
            setHoverInfo(prev => ({ ...prev, visible: false }));
            return;
          }

          const position = chart.canvas.getBoundingClientRect();
          const dataPoint = tooltip.dataPoints?.[0];
          
          if (dataPoint && sortedMetrics[dataPoint.dataIndex]) {
            const metric = sortedMetrics[dataPoint.dataIndex];
            const date = new Date(metric.timestamp);
            
            setHoverInfo({
              x: position.left + window.scrollX + tooltip.caretX,
              y: position.top + window.scrollY + tooltip.caretY,
              visible: true,
              date: date.toLocaleString('en-US', {
                weekday: 'short',
                month: 'short',
                day: 'numeric',
                hour: 'numeric',
                minute: '2-digit',
                hour12: true,
              }),
              latency: metric.latency,
            });
          }
        },
      },
    },
    scales: {
      x: {
        grid: {
          color: 'rgba(255, 255, 255, 0.1)',
          drawBorder: false,
        },
        ticks: {
          color: '#9ca3af',
          maxTicksLimit: selectedPeriod === 'day' ? 24 : selectedPeriod === 'week' ? 14 : 12,
          maxRotation: 45,
          minRotation: 0,
        },
      },
      y: {
        grid: {
          color: 'rgba(255, 255, 255, 0.1)',
          drawBorder: false,
        },
        ticks: {
          color: '#9ca3af',
          callback: (value: any) => `${value}ms`,
        },
        beginAtZero: true,
      },
    },
    elements: {
      point: {
        hoverBorderWidth: 3,
      },
    },
  }), [selectedPeriod, sortedMetrics]);

  const handlePeriodChange = (period: TimePeriod) => {
    setSelectedPeriod(period);
  };

  const formatLatency = (latency: number) => {
    return latency < 1000 ? `${latency.toFixed(1)}ms` : `${(latency / 1000).toFixed(2)}s`;
  };

  const getPeriodDisplayName = (period: TimePeriod) => {
    switch (period) {
      case 'day': return '24 Hours';
      case 'week': return '7 Days';
      case 'month': return '30 Days';
      default: return period;
    }
  };

  if (error) {
    return (
      <MetricsContainer>
        <NoDataMessage>
          <span className="icon">⚠️</span>
          <h3>Unable to Load Metrics</h3>
          <p>We're experiencing issues loading the system metrics. Please try again later.</p>
        </NoDataMessage>
      </MetricsContainer>
    );
  }

  return (
    <MetricsContainer>
      <Header>
        <MetricsTitle>System Metrics</MetricsTitle>
        <PeriodButtons>
          {(['day', 'week', 'month'] as TimePeriod[]).map((period) => (
            <PeriodButton
              key={period}
              $active={selectedPeriod === period}
              onClick={() => handlePeriodChange(period)}
            >
              {getPeriodDisplayName(period)}
            </PeriodButton>
          ))}
        </PeriodButtons>
      </Header>

      <ChartContainer $hasData={!loading && metricsData?.has_sufficient_data === true}>
        {loading ? (
          <NoDataMessage>
            <span className="icon">📊</span>
            <h3>Loading Metrics...</h3>
            <p>Fetching the latest system performance data.</p>
          </NoDataMessage>
        ) : !metricsData?.has_sufficient_data ? (
          <NoDataMessage>
            <span className="icon">⏳</span>
            <h3>Insufficient Data</h3>
            <p>{metricsData?.message || `We're still collecting data for the ${getPeriodDisplayName(selectedPeriod).toLowerCase()} view. Please check back later.`}</p>
          </NoDataMessage>
        ) : chartData ? (
          <>
            <Line data={chartData} options={chartOptions} />
            <HoverInfo
              $x={hoverInfo.x}
              $y={hoverInfo.y}
              $visible={hoverInfo.visible}
            >
              <div className="date">{hoverInfo.date}</div>
              <div className="latency">{formatLatency(hoverInfo.latency)}</div>
            </HoverInfo>
          </>
        ) : (
          <NoDataMessage>
            <span className="icon">📈</span>
            <h3>No Data Available</h3>
            <p>No metrics data is currently available for the selected time period.</p>
          </NoDataMessage>
        )}
      </ChartContainer>

      {metricsData?.has_sufficient_data && metricsData.stats && (
        <StatsGrid>
          <StatCard>
            <div className="label">Average</div>
            <div className="value">
              {formatLatency(metricsData.stats.average)}
            </div>
          </StatCard>
          <StatCard>
            <div className="label">Minimum</div>
            <div className="value">
              {formatLatency(metricsData.stats.min)}
            </div>
          </StatCard>
          <StatCard>
            <div className="label">Maximum</div>
            <div className="value">
              {formatLatency(metricsData.stats.max)}
            </div>
          </StatCard>
          <StatCard>
            <div className="label">Data Points</div>
            <div className="value">
              {metricsData.stats.count.toLocaleString()}
            </div>
          </StatCard>
        </StatsGrid>
      )}

      {metricsData?.last_updated && (
        <LastUpdated>
          Last updated: <span className="time">{new Date(metricsData.last_updated).toLocaleString()}</span>
        </LastUpdated>
      )}
    </MetricsContainer>
  );
};

export default SystemMetrics;