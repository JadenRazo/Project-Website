import { Cloud, Gauge, GitPullRequest, ShieldCheck } from 'lucide-react'

const principles = [
  {
    icon: Cloud,
    title: 'Infrastructure as a reviewed system',
    description:
      'Terraform, environment boundaries, least privilege, and explicit apply controls make infrastructure changes explainable before they are executable.',
  },
  {
    icon: GitPullRequest,
    title: 'Safe delivery paths',
    description:
      'Credential-free pull-request checks, immutable action references, OIDC, artifact provenance, and rollback ownership reduce the cost of a bad change.',
  },
  {
    icon: Gauge,
    title: 'Reliability with a denominator',
    description:
      'SLO math, bounded metrics, failure exercises, alerts, and runbooks turn “reliable” from a description into something a team can evaluate.',
  },
  {
    icon: ShieldCheck,
    title: 'Cost and security at design time',
    description:
      'Pricing drift, dependency risk, secrets, blast radius, and operational cost are design inputs, not cleanup work after launch.',
  },
]

export default function Services() {
  return (
    <section id="services" aria-labelledby="services-title" className="relative w-full border-b border-border py-16 sm:py-20 lg:py-28">
      <div className="portfolio-container">
        <div className="grid gap-6 lg:grid-cols-[0.7fr_1.3fr] lg:items-end">
          <div>
            <p className="mb-4 font-mono text-xs uppercase tracking-[0.2em] text-primary">
              Operating principles
            </p>
            <h2 id="services-title" className="font-display text-4xl font-bold tracking-[-0.04em] text-text-primary sm:text-5xl lg:text-6xl">
              How I approach the work.
            </h2>
          </div>
          <p className="max-w-2xl text-base leading-7 text-text-secondary lg:justify-self-end lg:text-lg lg:leading-8">
            The goal is not more tooling. It is a delivery system that makes
            risk visible, keeps writes intentional, and leaves the next operator
            with evidence they can trust.
          </p>
        </div>

        <div className="mt-12 grid border-l border-t border-border sm:grid-cols-2 lg:mt-16">
          {principles.map((principle, index) => (
            <article key={principle.title} className="border-b border-r border-border p-6 sm:p-8 lg:p-10">
              <div className="flex items-center justify-between">
                <principle.icon className="h-5 w-5 text-primary" aria-hidden="true" />
                <span className="font-mono text-xs text-text-muted">
                  {String(index + 1).padStart(2, '0')}
                </span>
              </div>
              <h3 className="mt-8 text-xl font-semibold leading-7 text-text-primary sm:text-2xl">
                {principle.title}
              </h3>
              <p className="mt-3 text-[15px] leading-7 text-text-secondary sm:text-base">
                {principle.description}
              </p>
            </article>
          ))}
        </div>
      </div>
    </section>
  )
}
