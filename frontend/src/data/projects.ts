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
    name: 'WeenieSMP',
    description: 'A full-stack Minecraft server ecosystem featuring a Vue 3 website with Tebex e-commerce integration, Go microservices delivering real-time player statistics and leaderboards, and MariaDB storing custom plugin data including playtime tracking, land claims, and bounty systems. Built with Redis caching, multi-domain Nginx routing, and fully containerized Docker deployment.',
    repo_url: 'https://github.com/JadenRazo/Project-Website/tree/main/weeniesmp',
    live_url: 'https://weeniesmp.net',
    tags: ['Vue 3', 'TypeScript', 'Go', 'MariaDB', 'Redis', 'Tebex API', 'Docker'],
    techCategories: {
      frontend: ['Vue 3', 'TypeScript', 'Pinia', 'Tailwind CSS'],
      backend: ['Go', 'Microservices'],
      database: ['MariaDB', 'Redis'],
      infrastructure: ['Docker', 'Nginx', 'TLS 1.3'],
      apis: ['Tebex API'],
    },
    status: 'active',
    mediaUrl: '/videos/weeniesmp_gambling_demo_optimized.mp4',
    mediaType: 'video',
    badges: ['live', 'demo'],
  },
  {
    id: '5',
    name: 'WeenieSMP CI/CD Pipeline',
    description: 'A multi-stage GitLab CI/CD pipeline for the WeenieSMP Minecraft server ecosystem covering validation, build, test, security scanning, and automated deployment with Discord notifications.',
    repo_url: 'https://github.com/JadenRazo/Project-Website/tree/main/weeniesmp',
    live_url: '',
    tags: ['GitLab CI', 'Docker', 'Terraform', 'SAST', 'Trivy', 'Semgrep'],
    techCategories: {
      infrastructure: ['GitLab CI', 'Docker', 'Terraform', 'Nginx'],
      other: ['SAST', 'Trivy', 'Semgrep', 'Secret Detection', 'Discord Webhooks'],
    },
    status: 'active',
    mediaUrl: 'cicd-pipeline',
    mediaType: 'component',
    badges: ['live'],
  },
  {
    id: '6',
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
    id: '7',
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

export default mockProjects;