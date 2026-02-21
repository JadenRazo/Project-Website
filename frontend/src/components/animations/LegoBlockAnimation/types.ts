import type { Theme } from '../../../styles/theme.types';

export type BlockShape = '1x1' | '1x2' | '1x3' | '1x4';

export type BlockContentType = 'icon' | 'code' | 'terminal' | 'symbol' | 'empty';

export type AnimationPhase = 'build' | 'settle' | 'dissolve' | 'revealed';

export type PerformanceTier = 'low' | 'medium' | 'high';

export interface BlockContent {
  type: BlockContentType;
  value?: string;
}

export interface GradientColor {
  start: string;
  end: string;
  blur?: boolean;
}

export interface StudPosition {
  x: number;
  y: number;
}

export interface Particle {
  x: number;
  y: number;
  vx: number;
  vy: number;
  life: number;
  maxLife: number;
  size: number;
  color: string;
  rotation: number;
  rotationSpeed: number;
}

export interface LegoBlock {
  id: string;
  shape: BlockShape;
  x: number;
  y: number;
  targetX: number;
  targetY: number;
  width: number;
  height: number;
  rotation: number;
  color: GradientColor;
  studs: StudPosition[];
  content: BlockContent;
  opacity: number;
  velocity: number;
  settled: boolean;
  dissolving: boolean;
  particles: Particle[];
  spawnDelay: number;
  spawned: boolean;
  glowIntensity: number;
  scaleX: number;
  scaleY: number;
  bouncePhase: number;
  settleStartTime: number;
}

export interface PerformanceConfig {
  blockCount: number;
  particlesPerBlock: number;
  enableGlow: boolean;
  enableGradients: boolean;
  enableStudDetail: boolean;
  frameTargetMs: number;
  enableBlur: boolean;
  dissolveSimplified: boolean;
}

export interface PhaseTimings {
  build: number;
  settle: number;
  dissolve: number;
}

export interface BlockGeneratorConfig {
  viewportWidth: number;
  viewportHeight: number;
  blockCount: number;
  gridCellSize: number;
}

export interface LegoCanvasProps {
  config: PerformanceConfig;
  theme: Theme;
  isVisible: boolean;
  phase: AnimationPhase;
  onPhaseChange: (phase: AnimationPhase) => void;
  onDissolveStart: () => void;
  width: number;
  height: number;
}

export interface LegoBlockAnimationProps {
  onAnimationComplete?: () => void;
}

export interface AnimationState {
  phase: AnimationPhase;
  startTime: number;
  phaseStartTime: number;
  blocks: LegoBlock[];
  showBio: boolean;
}

export interface GridCell {
  col: number;
  row: number;
  occupied: boolean;
}
