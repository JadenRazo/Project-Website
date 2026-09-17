import content from '../../../raizhost/content.json'
import assets from './project-film-assets.json'

export interface FilmChapter {
  at: number
  title: string
  description: string
}
export interface ProjectFilm {
  id: string
  title: string
  category: string
  color: string
  headline: string
  summary: string
  kind: string
  duration: number
  tags: string[]
  chapters: FilmChapter[]
  sourceNote: string
  evidenceUrl: string
  evidenceLabel: string
  liveUrl?: string
  liveLabel?: string
  src: string
  mobileSrc: string
  webmSrc: string
  mobileWebmSrc: string
  poster: string
  mobilePoster: string
  captions: string
}

const stories = [
  {
    id: 'raizhost',
    color: '#9ac4ae',
    kind: 'Product & architecture tour',
    duration: 24,
    tags: ['AWS', 'CloudFront', 'S3', 'GitHub Actions'],
    evidenceUrl: 'https://github.com/JadenRazo/raizhost-architecture',
    evidenceLabel: 'Architecture',
    liveUrl: 'https://raizhost.com',
    liveLabel: 'Explore the site',
    sourceNote:
      'Public website capture, September 2026. Publishing is a source-backed architecture walkthrough; no customer account or production publish is shown.',
    chapters: [
      {
        at: 0,
        title: 'The experience',
        description:
          'A visitor explores the RaizHost website and its examples. The film shows the public site, captured in September 2026.',
      },
      {
        at: 8,
        title: 'The edit',
        description:
          'The owner editor changes a bounded content file while the site’s design stays in code. This portion illustrates the documented editing contract.',
      },
      {
        at: 16,
        title: 'The delivery',
        description:
          'Publishing creates a commit in the client’s repository. Its workflow builds the site, updates S3, and invalidates CloudFront. The public architecture record explains the boundaries.',
      },
    ],
  },
  {
    id: 'cloudcostmcp',
    color: '#a5baf4',
    kind: 'Recorded CLI output',
    duration: 24,
    tags: ['TypeScript', 'Terraform', 'MCP', 'AWS · Azure · GCP'],
    evidenceUrl: 'https://github.com/JadenRazo/CloudCostMCP',
    evidenceLabel: 'Source code',
    liveUrl: 'https://www.npmjs.com/package/@jadenrazo/cloudcost-mcp',
    liveLabel: 'Get the package',
    sourceNote:
      'Real CLI output from a small sample Terraform file. Estimates reflect the capture’s pricing source and assumptions; they are not a quote or a current cloud bill.',
    chapters: [
      {
        at: 0,
        title: 'Read the code',
        description:
          'CloudCost parses a sample Terraform configuration and identifies its compute resource and attributes.',
      },
      {
        at: 8,
        title: 'Estimate spend',
        description:
          'The estimate command returns a resource-level monthly cost with pricing metadata. The demo uses a sample workload, not a production account.',
      },
      {
        at: 16,
        title: 'See assumptions',
        description:
          'The report separates a live compute price from a synthetic data-transfer allowance and includes pricing metadata, confidence, and warnings.',
      },
    ],
  },
  {
    id: 'tickethacker',
    color: '#d8bca0',
    kind: 'Interface demo · sample data',
    duration: 24,
    tags: ['React', 'NestJS', 'PostgreSQL', 'Socket.IO'],
    evidenceUrl: 'https://github.com/JadenRazo/TicketHacker',
    evidenceLabel: 'Source code',
    sourceNote:
      'The actual dashboard runs locally with intercepted sample API responses. All names and conversations are fictional. This demonstrates the interface, not a live messaging integration.',
    chapters: [
      {
        at: 0,
        title: 'Triage the inbox',
        description:
          'The actual agent dashboard displays sample tickets from several channels, with status, priority, and assignment in one list.',
      },
      {
        at: 8,
        title: 'Focus the queue',
        description:
          'Filtering the actual inbox by high priority narrows the list to a sample conversation.',
      },
      {
        at: 16,
        title: 'Move work forward',
        description:
          'The agent changes the sample ticket’s status. The local capture uses fixture responses, so no customer messages or accounts are modified.',
      },
    ],
  },
  {
    id: 'llm-lint',
    color: '#c3b0ed',
    kind: 'Recorded CLI output',
    duration: 24,
    tags: ['Go', 'CLI', 'SARIF', 'GitHub Actions'],
    evidenceUrl: 'https://github.com/JadenRazo/llm-lint',
    evidenceLabel: 'Source code',
    liveUrl: 'https://www.npmjs.com/package/@jadenrazo/llm-lint',
    liveLabel: 'Get the package',
    sourceNote:
      'Real scan and fix-preview output from a disposable sample repository. Rules are configurable; required authorship, attribution, and disclosures must be retained.',
    chapters: [
      {
        at: 0,
        title: 'Scan the repo',
        description:
          'A disposable sample repository intentionally includes a local tool configuration file. The scanner identifies the policy violation.',
      },
      {
        at: 8,
        title: 'Explain the finding',
        description:
          'The result identifies the file, rule, and severity so an engineer can understand what needs attention.',
      },
      {
        at: 16,
        title: 'Preview the repair',
        description:
          'Fix-preview explains the proposed cleanup without modifying the sample repository. Teams decide which rules match their own publication policy.',
      },
    ],
  },
  {
    id: 'sre-reference-app',
    color: '#9dcce3',
    kind: 'Documented exercise walkthrough',
    duration: 24,
    tags: ['ECS Fargate', 'Terraform', 'CloudWatch', 'SLOs'],
    evidenceUrl:
      'https://github.com/JadenRazo/sre-reference-app/blob/main/docs/chaos-experiments.md',
    evidenceLabel: 'Read the exercise',
    sourceNote:
      'Animated explanation of the repository’s historical controlled task-stop exercise. The recorded 78-second recovery is one observation, not an uptime or future recovery guarantee. The timeline is compressed.',
    chapters: [
      {
        at: 0,
        title: 'Serve traffic',
        description:
          'An application load balancer routes requests to two ECS Fargate tasks. CloudWatch observes the service.',
      },
      {
        at: 8,
        title: 'Lose a task',
        description:
          'The documented exercise stopped one task with the ECS API. The surviving task absorbed traffic as the load balancer drained the stopped task.',
      },
      {
        at: 16,
        title: 'Measure recovery',
        description:
          'ECS replaced the task. The exercise recorded a 78-second recovery window and both burn-rate alarms remained OK. This is a compressed historical explanation, not a new live drill.',
      },
    ],
  },
  {
    id: 'sre-landing-zone',
    color: '#d4c495',
    kind: 'Architecture walkthrough',
    duration: 24,
    tags: ['AWS Organizations', 'Terraform', 'IAM', 'Disaster recovery'],
    evidenceUrl: 'https://github.com/JadenRazo/sre-landing-zone',
    evidenceLabel: 'Explore the design',
    sourceNote:
      'Source-backed architecture animation of the May 2026 lab. It explains the recorded design and does not assert that the AWS resources are currently running.',
    chapters: [
      {
        at: 0,
        title: 'Separate concerns',
        description:
          'The recorded organization uses five accounts: management, log archive, audit and security, development workloads, and production workloads.',
      },
      {
        at: 8,
        title: 'Prepare recovery',
        description:
          'The pilot-light design has standby infrastructure in a second region with ECS desired count at zero until an operator scales it up.',
      },
      {
        at: 16,
        title: 'Control idle cost',
        description:
          'A scheduled Lambda assumes a limited cross-account role and scales development ECS services to zero. The permission is conditioned on the Environment=dev resource tag.',
      },
    ],
  },
]

const descriptions = {
  'raizhost': content['film-raizhost'],
  'cloudcostmcp': content['film-cloudcostmcp'],
  'tickethacker': content['film-tickethacker'],
  'llm-lint': content['film-llm-lint'],
  'sre-reference-app': content['film-sre-reference-app'],
  'sre-landing-zone': content['film-sre-landing-zone'],
}

export const projectFilms: ProjectFilm[] = stories.map((story) => ({
  ...story,
  ...descriptions[story.id as keyof typeof descriptions],
  ...assets[story.id as keyof typeof assets],
}))
