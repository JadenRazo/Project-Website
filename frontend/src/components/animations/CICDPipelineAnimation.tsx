import React, { useState, useEffect, useRef } from 'react';
import styled, { keyframes, css } from 'styled-components';

const STAGES = [
  { name: 'validate', desc: 'Checking lint & formatting' },
  { name: 'build', desc: 'Building Docker images' },
  { name: 'test', desc: 'Running unit tests' },
  { name: 'security', desc: 'Scanning with SAST & Trivy' },
  { name: 'deploy', desc: 'Deploying to production' },
  { name: 'notify', desc: 'Sending Discord notification' },
];

const ELAPSED = ['0:23', '1:08', '2:31', '3:15', '3:52', '3:58'];
const JOBS_PER_STAGE = [2, 3, 4, 3, 2, 1];
const TOTAL_JOBS = 15;

type Status = 'pending' | 'running' | 'passed';

const INTERVAL_MS = 900;
const RUN_MS = 550;
const START_DELAY = 600;
const HOLD_MS = 2500;
const CYCLE_MS = START_DELAY + STAGES.length * INTERVAL_MS + RUN_MS + HOLD_MS;

const pulse = keyframes`
  0%, 100% { box-shadow: 0 0 0 0 rgba(56,132,255,0.5); }
  50% { box-shadow: 0 0 0 5px rgba(56,132,255,0); }
`;

const checkPop = keyframes`
  0% { transform: scale(0); }
  60% { transform: scale(1.3); }
  100% { transform: scale(1); }
`;

const dotPulse = keyframes`
  0%, 100% { opacity: 1; }
  50% { opacity: 0.4; }
`;

const Container = styled.div`
  width: 100%;
  height: 100%;
  background: linear-gradient(180deg, #0d1117 0%, #161b22 100%);
  display: flex;
  flex-direction: column;
  padding: 14px 16px 12px;
  overflow: hidden;
  position: relative;
  user-select: none;
`;

const Header = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 4px;
`;

const TitleGroup = styled.div`
  display: flex;
  align-items: center;
  gap: 7px;
`;

const PipelineIconBox = styled.div`
  width: 15px;
  height: 15px;
  border-radius: 3px;
  background: linear-gradient(135deg, #3fb950, #238636);
  display: grid;
  place-items: center;
  flex-shrink: 0;

  &::after {
    content: '';
    width: 5px;
    height: 5px;
    border: 1.5px solid rgba(255,255,255,0.9);
    border-radius: 50%;
  }
`;

const PipelineLabel = styled.span`
  font-size: 11.5px;
  font-weight: 600;
  color: #e6edf3;
  letter-spacing: 0.2px;
`;

const StatusBadge = styled.span<{ $s: string }>`
  font-size: 9.5px;
  font-weight: 600;
  padding: 2px 7px;
  border-radius: 10px;
  letter-spacing: 0.3px;
  transition: all 0.3s ease;
  ${({ $s }) => {
    if ($s === 'passed') return css`background: rgba(63,185,80,0.15); color: #3fb950; border: 1px solid rgba(63,185,80,0.25);`;
    if ($s === 'running') return css`background: rgba(56,132,255,0.12); color: #58a6ff; border: 1px solid rgba(56,132,255,0.25);`;
    return css`background: rgba(139,148,158,0.1); color: #8b949e; border: 1px solid rgba(139,148,158,0.2);`;
  }}
`;

const MetaRow = styled.div`
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 4px;
  font-size: 10px;
  color: #484f58;
  overflow: hidden;
  white-space: nowrap;
`;

const BranchTag = styled.span`
  background: rgba(56,132,255,0.08);
  color: #58a6ff;
  padding: 1px 5px;
  border-radius: 3px;
  font-family: 'SFMono-Regular', Consolas, 'Liberation Mono', Menlo, monospace;
  font-size: 10px;
  flex-shrink: 0;
`;

const Sha = styled.span`
  font-family: 'SFMono-Regular', Consolas, 'Liberation Mono', Menlo, monospace;
  color: #484f58;
  flex-shrink: 0;
`;

const TrackWrapper = styled.div`
  flex: 1;
  display: flex;
  align-items: center;
  min-height: 0;
`;

const TrackInner = styled.div`
  width: 100%;
  position: relative;
  padding: 0 2px;
`;

const LineContainer = styled.div`
  position: absolute;
  top: 11px;
  left: 24px;
  right: 24px;
  height: 2px;
`;

const LineBase = styled.div`
  width: 100%;
  height: 100%;
  background: #21262d;
  border-radius: 1px;
`;

const LineFill = styled.div<{ $w: number }>`
  position: absolute;
  top: 0;
  left: 0;
  height: 100%;
  width: ${({ $w }) => $w}%;
  background: linear-gradient(90deg, #238636, #3fb950);
  border-radius: 1px;
  transition: width 0.5s cubic-bezier(0.25, 0.46, 0.45, 0.94);
`;

const NodesRow = styled.div`
  display: flex;
  justify-content: space-between;
  position: relative;
`;

const StageNode = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 5px;
  z-index: 1;
`;

const Circle = styled.div<{ $s: Status }>`
  width: 24px;
  height: 24px;
  border-radius: 50%;
  display: grid;
  place-items: center;
  flex-shrink: 0;
  transition: background 0.3s ease, border-color 0.3s ease;
  border: 2px solid ${({ $s }) =>
    $s === 'passed' ? '#3fb950' :
    $s === 'running' ? '#58a6ff' : '#30363d'
  };
  background: ${({ $s }) =>
    $s === 'passed' ? '#3fb950' :
    $s === 'running' ? 'rgba(56,132,255,0.15)' : '#0d1117'
  };
  ${({ $s }) => $s === 'running' && css`animation: ${pulse} 1.5s ease-in-out infinite;`}
`;

const CheckMark = styled.span`
  color: white;
  font-size: 11px;
  font-weight: 800;
  line-height: 1;
  animation: ${checkPop} 0.3s ease forwards;
`;

const RunningDot = styled.div`
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #58a6ff;
  animation: ${dotPulse} 0.8s ease-in-out infinite;
`;

const NodeName = styled.span<{ $s: Status }>`
  font-size: 8.5px;
  font-weight: 500;
  letter-spacing: 0.2px;
  color: ${({ $s }) =>
    $s === 'passed' ? '#3fb950' :
    $s === 'running' ? '#58a6ff' : '#484f58'
  };
  font-family: 'SFMono-Regular', Consolas, 'Liberation Mono', Menlo, monospace;
  transition: color 0.3s ease;

  @media (max-width: 400px) {
    font-size: 7.5px;
  }
`;

const InfoPanel = styled.div`
  display: flex;
  flex-direction: column;
  gap: 5px;
`;

const ActivityLine = styled.div<{ $highlight: boolean }>`
  font-size: 10px;
  color: ${({ $highlight }) => $highlight ? '#c9d1d9' : '#484f58'};
  font-family: 'SFMono-Regular', Consolas, 'Liberation Mono', Menlo, monospace;
  min-height: 14px;
  transition: color 0.3s ease;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
`;

const BarTrack = styled.div`
  width: 100%;
  height: 3px;
  background: #21262d;
  border-radius: 2px;
  overflow: hidden;
`;

const BarFill = styled.div<{ $w: number; $done: boolean }>`
  height: 100%;
  border-radius: 2px;
  transition: width 0.4s ease;
  width: ${({ $w }) => $w}%;
  background: ${({ $done }) => $done
    ? 'linear-gradient(90deg, #238636, #3fb950)'
    : 'linear-gradient(90deg, #1f6feb, #58a6ff)'
  };
`;

const StatsRow = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 9.5px;
  color: #484f58;
  font-family: 'SFMono-Regular', Consolas, 'Liberation Mono', Menlo, monospace;
`;

const CICDPipelineAnimation: React.FC = () => {
  const [statuses, setStatuses] = useState<Status[]>(Array(6).fill('pending'));
  const [elapsed, setElapsed] = useState('0:00');
  const timeoutsRef = useRef<ReturnType<typeof setTimeout>[]>([]);

  useEffect(() => {
    const clear = () => {
      timeoutsRef.current.forEach(t => clearTimeout(t));
      timeoutsRef.current = [];
    };

    const runCycle = () => {
      clear();
      setStatuses(Array(6).fill('pending'));
      setElapsed('0:00');

      STAGES.forEach((_, i) => {
        const t0 = START_DELAY + i * INTERVAL_MS;
        timeoutsRef.current.push(setTimeout(() => {
          setStatuses(prev => prev.map((s, j) => j === i ? 'running' : s));
        }, t0));
        timeoutsRef.current.push(setTimeout(() => {
          setStatuses(prev => prev.map((s, j) => j === i ? 'passed' : s));
          setElapsed(ELAPSED[i]);
        }, t0 + RUN_MS));
      });
    };

    runCycle();
    const id = setInterval(runCycle, CYCLE_MS);
    return () => { clearInterval(id); clear(); };
  }, []);

  const passedCount = statuses.filter(s => s === 'passed').length;
  const runningIdx = statuses.indexOf('running');
  const allPassed = passedCount === STAGES.length;
  const jobsDone = JOBS_PER_STAGE.slice(0, passedCount).reduce((a, b) => a + b, 0);
  const lineW = passedCount >= 2 ? ((passedCount - 1) / (STAGES.length - 1)) * 100 : 0;
  const barW = (jobsDone / TOTAL_JOBS) * 100;

  const activity = allPassed
    ? 'Pipeline completed successfully'
    : runningIdx >= 0
      ? `${STAGES[runningIdx].desc}...`
      : 'Waiting for runner...';

  const badgeLabel = allPassed ? 'passed' : runningIdx >= 0 ? 'running' : 'created';

  return (
    <Container>
      <Header>
        <TitleGroup>
          <PipelineIconBox />
          <PipelineLabel>Pipeline #1847</PipelineLabel>
        </TitleGroup>
        <StatusBadge $s={badgeLabel}>{badgeLabel}</StatusBadge>
      </Header>

      <MetaRow>
        <BranchTag>main</BranchTag>
        <Sha>a3f8c2d</Sha>
        <span>Update CI/CD deployment config</span>
      </MetaRow>

      <TrackWrapper>
        <TrackInner>
          <LineContainer>
            <LineBase />
            <LineFill $w={lineW} />
          </LineContainer>
          <NodesRow>
            {STAGES.map((stage, i) => (
              <StageNode key={stage.name}>
                <Circle $s={statuses[i]}>
                  {statuses[i] === 'passed' && <CheckMark>&#10003;</CheckMark>}
                  {statuses[i] === 'running' && <RunningDot />}
                </Circle>
                <NodeName $s={statuses[i]}>{stage.name}</NodeName>
              </StageNode>
            ))}
          </NodesRow>
        </TrackInner>
      </TrackWrapper>

      <InfoPanel>
        <ActivityLine $highlight={runningIdx >= 0 || allPassed}>
          {activity}
        </ActivityLine>
        <BarTrack>
          <BarFill $w={barW} $done={allPassed} />
        </BarTrack>
        <StatsRow>
          <span>{jobsDone}/{TOTAL_JOBS} jobs</span>
          <span>{elapsed}</span>
        </StatsRow>
      </InfoPanel>
    </Container>
  );
};

export default CICDPipelineAnimation;
