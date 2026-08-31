// Mock project data that matches our backend structure
export interface TechCategory {
  frontend?: string[];
  backend?: string[];
  database?: string[];
  infrastructure?: string[];
  apis?: string[];
  other?: string[];
}

export interface ProjectData {
  id: string;
  name: string;
  description: string;
  repo_url: string;
  live_url: string;
  tags: string[];
  techCategories?: TechCategory;
  status: string;
  mediaUrl?: string;
  mediaType?: 'image' | 'video' | 'component';
  badges?: Array<'live' | 'demo' | 'client' | 'internal'>;
}

export const mockProjects: ProjectData[] = [
  {
    id: '1',
    name: 'Portfolio Website',
    description: 'A modern, responsive portfolio website built with React, TypeScript, and styled-components featuring real-time messaging, URL shortener, and developer panel.',
    repo_url: 'https://github.com/JadenRazo/Project-Website',
    live_url: 'https://jadenrazo.dev',
    tags: ['React', 'TypeScript', 'Go', 'PostgreSQL', 'WebSocket', 'REST API'],
    techCategories: {
      frontend: ['React', 'TypeScript', 'Styled Components'],
      backend: ['Go', 'REST API', 'WebSocket'],
      database: ['PostgreSQL', 'Redis'],
      infrastructure: ['Docker', 'Nginx', 'Prometheus'],
    },
    status: 'active',
    mediaUrl: '/images/projects/portfolio-workspace_optimized.jpg',
    mediaType: 'image',
    badges: ['live'],
  },
  {
    id: '2',
    name: 'Showers Auto Detail',
    description: 'A mobile-first auto detailing booking platform with instant quotes, online booking, Square payment integration, before/after gallery, and admin dashboard with 2FA authentication.',
    repo_url: 'https://github.com/JadenRazo/showersautodetail',
    live_url: 'https://showersautodetail.com',
    tags: ['Astro', 'React', 'Node.js', 'PostgreSQL', 'Tailwind CSS', 'Square API', 'Docker'],
    techCategories: {
      frontend: ['Astro', 'React', 'Tailwind CSS'],
      backend: ['Node.js'],
      database: ['PostgreSQL'],
      infrastructure: ['Docker'],
      apis: ['Square API'],
    },
    status: 'active',
    mediaUrl: '/images/projects/showers-auto-detail.jpg',
    mediaType: 'image',
    badges: ['live', 'client'],
  },
  {
    id: '3',
    name: 'Educational Quiz Discord Bot',
    description: 'An advanced Discord bot that leverages LLMs to create educational quizzes with multi-guild support, achievement system, and real-time leaderboards.',
    repo_url: 'https://github.com/JadenRazo/Quiz-Bot',
    live_url: '',
    tags: ['Python', 'Discord.py', 'PostgreSQL', 'OpenAI API', 'Anthropic Claude', 'Google Gemini'],
    techCategories: {
      backend: ['Python', 'Discord.py'],
      database: ['PostgreSQL'],
      apis: ['OpenAI API', 'Anthropic Claude', 'Google Gemini'],
    },
    status: 'active',
    mediaUrl: '/videos/web_ready_quizbot_example_video_optimized.mp4',
    mediaType: 'video',
    badges: ['demo'],
  },
  {
    id: '4',
    name: 'URL Shortener Service',
    description: 'A high-performance URL shortening service with analytics, custom short codes, and comprehensive statistics tracking.',
    repo_url: 'https://github.com/JadenRazo/Project-Website/tree/main/backend/internal/urlshortener',
    live_url: 'https://jadenrazo.dev/s/',
    tags: ['Go', 'PostgreSQL', 'Analytics', 'REST API', 'Microservice'],
    techCategories: {
      backend: ['Go', 'REST API'],
      database: ['PostgreSQL', 'Redis'],
      infrastructure: ['Docker', 'Microservices'],
    },
    status: 'active',
    mediaUrl: '/images/projects/url-shortener.svg',
    mediaType: 'image',
    badges: ['internal'],
  },
  {
    id: '5',
    name: 'Code Statistics Tracker',
    description: 'Automated system for tracking lines of code across projects with scheduled updates and API integration.',
    repo_url: 'https://github.com/JadenRazo/Project-Website/tree/main/scripts',
    live_url: 'https://jadenrazo.dev/api/v1/code/stats',
    tags: ['Go', 'Automation', 'CLI', 'Statistics', 'CRON'],
    techCategories: {
      backend: ['Go', 'CLI'],
      infrastructure: ['CRON', 'Automation'],
    },
    status: 'active',
    mediaUrl: '/images/projects/code-stats.svg',
    mediaType: 'image',
    badges: ['internal'],
  },
];

export interface EvidenceProject extends ProjectData {
  evidence: string;
  discipline: string;
}

export const evidenceProjects: EvidenceProject[] = [
  {
    id: 'cloudcostmcp',
    name: 'CloudCostMCP',
    discipline: 'Cost engineering',
    description:
      'Multi-IaC cost analysis for AWS, Azure, and GCP, with provider-backed pricing verification and incident-driven regression controls.',
    evidence:
      'Pricing-drift defect fixed; 41 live catalog checks and 1,564 deterministic tests exercise the corrected behavior.',
    repo_url: 'https://github.com/JadenRazo/CloudCostMCP',
    live_url: 'https://www.npmjs.com/package/@jadenrazo/cloudcost-mcp',
    tags: ['TypeScript', 'AWS', 'Azure', 'GCP', 'MCP'],
    status: 'active',
  },
  {
    id: 'sre-reference-app',
    name: 'SRE Reference App',
    discipline: 'Reliability engineering',
    description:
      'An ECS Fargate reference system that keeps SLO math, burn-rate alarms, failure injection, runbooks, and postmortem material beside the Terraform.',
    evidence:
      'A controlled task-failure exercise records 78-second recovery, the observed alarm behavior, and the limits of the result.',
    repo_url: 'https://github.com/JadenRazo/sre-reference-app',
    live_url: '',
    tags: ['AWS', 'Terraform', 'ECS', 'Prometheus', 'SLOs'],
    status: 'active',
  },
  {
    id: 'llm-lint',
    name: 'llm-lint',
    discipline: 'Developer tooling',
    description:
      'A Go policy scanner for generated repository artifacts and boundary violations, distributed as native binaries and an npm package.',
    evidence:
      'CI publishes SARIF, runs CodeQL and release verification, generates SBOMs, and signs release artifacts.',
    repo_url: 'https://github.com/JadenRazo/llm-lint',
    live_url: 'https://www.npmjs.com/package/@jadenrazo/llm-lint',
    tags: ['Go', 'SARIF', 'CodeQL', 'SBOM', 'Cosign'],
    status: 'active',
  },
  {
    id: 'aws-supply-chain-security',
    name: 'AWS Supply Chain Security',
    discipline: 'Secure delivery',
    description:
      'A container supply-chain reference that separates credential-free review checks from explicitly dispatched AWS write paths.',
    evidence:
      'Immutable action pins, OIDC role chaining, SBOM generation, image scanning, and keyless signing are visible in the workflows.',
    repo_url: 'https://github.com/JadenRazo/aws-supply-chain-security',
    live_url: '',
    tags: ['AWS', 'GitHub Actions', 'OIDC', 'Terraform', 'Cosign'],
    status: 'active',
  },
  {
    id: 'tts-raizhost',
    name: 'TTS RaizHost',
    discipline: 'Kubernetes reliability',
    description:
      'A self-hosted reader that routes synthesis between a preferred GPU service and an in-cluster CPU fallback, with bounded metrics and operator runbooks.',
    evidence:
      'Probe and circuit transitions are tested across TypeScript and Python contracts; the repository explicitly distinguishes those checks from an unexecuted live failure drill.',
    repo_url: 'https://github.com/JadenRazo/tts-raizhost',
    live_url: '',
    tags: ['Kubernetes', 'Next.js', 'Python', 'Prometheus', 'OpenTelemetry'],
    status: 'active',
  },
  {
    id: 'raizhost-architecture',
    name: 'RaizHost Architecture',
    discipline: 'Technical decision record',
    description:
      'A public architecture record for the AWS platform behind RaizHost, including service choices, cost constraints, and operational boundaries.',
    evidence:
      'The record documents roughly ten sites and three applications at about $50 per month while clearly separating public evidence from private source.',
    repo_url: 'https://github.com/JadenRazo/raizhost-architecture',
    live_url: 'https://raizhost.com',
    tags: ['AWS', 'CloudFront', 'Lambda', 'Graviton', 'Cost'],
    status: 'active',
  },
]

export default mockProjects;
