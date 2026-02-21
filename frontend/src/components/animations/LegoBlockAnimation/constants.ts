import type { Theme } from '../../../styles/theme.types';
import type {
  PerformanceConfig,
  PhaseTimings,
  BlockShape,
  GradientColor,
  PerformanceTier
} from './types';

export const PHASE_TIMINGS: PhaseTimings = {
  build: 3500,
  settle: 800,
  dissolve: 2000,
};

export const GRID = {
  cellSize: 32,
  studRadius: 4,
  studSpacing: 16,
  blockBorderRadius: 4,
  padding: 8,
};

export const ANIMATION = {
  fallSpeed: 12,
  maxFallSpeed: 20,
  gravity: 0.8,
  bounceStrength: 0.3,
  bounceDuration: 150,
  spawnInterval: 80,
  initialDelay: 100,
};

export const DISSOLVE = {
  sweepSpeed: 0.004,
  particleDuration: 600,
  particleGravity: -0.1,
  particleDrag: 0.98,
};

export const SETTLE = {
  glowPulseDuration: 500,
  glowMaxIntensity: 0.3,
};

export const BLOCK_SHAPES: Record<BlockShape, { cols: number; rows: number }> = {
  '1x1': { cols: 1, rows: 1 },
  '1x2': { cols: 2, rows: 1 },
  '1x3': { cols: 3, rows: 1 },
  '1x4': { cols: 4, rows: 1 },
};

export const SHAPE_WEIGHTS: Record<BlockShape, number> = {
  '1x1': 0.15,
  '1x2': 0.35,
  '1x3': 0.30,
  '1x4': 0.20,
};

export const BLOCK_CONTENT = {
  technologies: [
    // Frontend Frameworks & Libraries
    'React', 'TypeScript', 'Astro', 'Next.js', 'Tailwind', 'Redux', 'Zustand',
    'Framer', 'Three.js', 'GSAP', 'Vite', 'Webpack', 'ESLint', 'Prettier',
    // Backend & Languages
    'Go', 'Node.js', 'Python', 'Express', 'Gin', 'FastAPI', 'GraphQL', 'REST',
    // Databases & Caching
    'PostgreSQL', 'Redis', 'MongoDB', 'SQLite', 'Prisma',
    // DevOps & Infrastructure
    'Docker', 'Nginx', 'Linux', 'Ubuntu', 'Debian', 'AWS', 'Vercel', 'Netlify',
    'GitHub', 'CI/CD', 'Prometheus', 'Grafana',
    // Tools & Utilities
    'Git', 'tmux', 'Vim', 'VSCode', 'Zsh', 'Bash', 'Fish',
    'npm', 'pnpm', 'Yarn', 'Make', 'systemd', 'PM2',
    // Networking & Security
    'DNS', 'SFTP', 'SSH', 'SSL/TLS', 'OAuth', 'JWT', 'CORS', 'HTTPS',
    // Web Technologies
    'HTML5', 'CSS3', 'SCSS', 'WebSocket', 'HTTP/2', 'gRPC',
    // Testing & Quality
    'Jest', 'Vitest', 'Cypress', 'Playwright',
    // Data & Formats
    'JSON', 'YAML', 'Markdown', 'SQL',
  ],
};

export const getBlockGradients = (theme: Theme): GradientColor[] => [
  {
    start: theme.colors.primary,
    end: theme.colors.primaryHover,
  },
  {
    start: theme.colors.accent,
    end: theme.colors.accentHover,
  },
  {
    start: `${theme.colors.primary}dd`,
    end: `${theme.colors.accent}dd`,
  },
  {
    start: theme.colors.secondary,
    end: theme.colors.secondaryHover,
  },
  {
    start: `${theme.colors.accent}cc`,
    end: `${theme.colors.primary}cc`,
  },
];

export const getPerformanceConfig = (tier: PerformanceTier): PerformanceConfig => {
  const configs: Record<PerformanceTier, PerformanceConfig> = {
    low: {
      blockCount: 18,
      particlesPerBlock: 2,
      enableGlow: false,
      enableGradients: false,
      enableStudDetail: false,
      frameTargetMs: 33,
      enableBlur: false,
      dissolveSimplified: true,
    },
    medium: {
      blockCount: 30,
      particlesPerBlock: 3,
      enableGlow: true,
      enableGradients: true,
      enableStudDetail: true,
      frameTargetMs: 22,
      enableBlur: false,
      dissolveSimplified: false,
    },
    high: {
      blockCount: 45,
      particlesPerBlock: 4,
      enableGlow: true,
      enableGradients: true,
      enableStudDetail: true,
      frameTargetMs: 16,
      enableBlur: true,
      dissolveSimplified: false,
    },
  };
  return configs[tier];
};
