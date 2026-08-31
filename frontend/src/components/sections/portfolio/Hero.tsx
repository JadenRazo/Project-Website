import { Activity, ArrowDownRight, ExternalLink, Github } from 'lucide-react'

import { useScroll } from '../../../providers/ScrollProvider'

const evidence = [
  {
    label: 'Cost correctness',
    detail: 'Live catalog drift regression',
    href: 'https://github.com/JadenRazo/CloudCostMCP/blob/main/docs/incidents/2026-08-pricing-drift.md',
  },
  {
    label: 'Failure recovery',
    detail: '78-second controlled ECS exercise',
    href: 'https://github.com/JadenRazo/sre-reference-app',
  },
  {
    label: 'Release trust',
    detail: 'CodeQL, SBOMs, and signed artifacts',
    href: 'https://github.com/JadenRazo/llm-lint',
  },
]

export default function Hero() {
  const { scrollTo } = useScroll()

  const inspectWork = (event: React.MouseEvent<HTMLAnchorElement>) => {
    event.preventDefault()
    scrollTo('#projects')
  }

  return (
    <section
      id="home"
      aria-labelledby="hero-title"
      className="relative flex min-h-[760px] w-full items-center border-b border-border pb-14 pt-28 sm:pt-32 lg:min-h-[820px] lg:pb-20"
    >
      <div className="portfolio-container grid items-center gap-12 lg:grid-cols-[minmax(0,1.2fr)_minmax(340px,0.8fr)] lg:gap-16">
        <div className="max-w-4xl">
          <div className="mb-7 inline-flex items-center gap-2 border border-primary/40 bg-primary/10 px-3 py-2 font-mono text-[11px] uppercase tracking-[0.16em] text-primary sm:text-xs">
            <span className="h-2 w-2 rounded-full bg-neon" aria-hidden="true" />
            Open to full-time cloud / DevOps roles
          </div>

          <p className="mb-4 font-mono text-xs uppercase tracking-[0.22em] text-text-secondary sm:text-sm">
            AWS · Terraform · Go · SRE
          </p>
          <h1
            id="hero-title"
            className="max-w-4xl font-display text-[clamp(2.85rem,8.5vw,6.75rem)] font-bold leading-[0.94] tracking-[-0.055em] text-text-primary"
          >
            I build reliable <span className="text-primary">cloud systems.</span>
          </h1>
          <p className="mt-7 max-w-2xl text-base leading-7 text-text-secondary sm:text-lg sm:leading-8 lg:text-xl">
            I&apos;m Jaden Razo, a cloud and DevOps engineer focused on secure
            delivery, measurable reliability, and cost-aware AWS infrastructure.
            I turn architecture intent into tests, runbooks, and evidence a team
            can inspect.
          </p>

          <div className="mt-8 flex flex-col gap-3 min-[430px]:flex-row sm:mt-10">
            <a href="#projects" onClick={inspectWork} className="btn-primary">
              <span>Inspect the work</span>
              <ArrowDownRight size={17} aria-hidden="true" />
            </a>
            <a
              href="https://github.com/JadenRazo"
              target="_blank"
              rel="noopener noreferrer"
              className="btn-secondary"
            >
              <Github size={17} aria-hidden="true" />
              <span>GitHub profile</span>
            </a>
          </div>

          <div className="mt-8 flex flex-wrap items-center gap-x-5 gap-y-2 text-xs text-text-muted sm:text-sm">
            <span className="font-mono uppercase tracking-[0.14em]">Operating stack</span>
            <span>AWS</span>
            <span>Terraform</span>
            <span>GitHub Actions</span>
            <span>Linux</span>
          </div>
        </div>

        <aside
          aria-label="Selected engineering evidence"
          className="border border-border bg-background-secondary/90 p-5 shadow-2xl shadow-black/20 sm:p-7"
        >
          <div className="flex items-center justify-between border-b border-border pb-4">
            <div className="flex items-center gap-2">
              <Activity className="h-4 w-4 text-neon" aria-hidden="true" />
              <span className="font-mono text-xs uppercase tracking-[0.18em] text-text-primary">
                Selected evidence
              </span>
            </div>
            <span className="font-mono text-[10px] uppercase tracking-[0.16em] text-text-muted">
              Public
            </span>
          </div>

          <div className="divide-y divide-border">
            {evidence.map((item, index) => (
              <a
                key={item.label}
                href={item.href}
                target="_blank"
                rel="noopener noreferrer"
                className="group grid grid-cols-[2rem_1fr_auto] items-start gap-3 py-5"
              >
                <span className="pt-0.5 font-mono text-xs text-text-muted">
                  {String(index + 1).padStart(2, '0')}
                </span>
                <span>
                  <span className="block text-sm font-semibold text-text-primary group-hover:text-primary">
                    {item.label}
                  </span>
                  <span className="mt-1 block text-sm leading-6 text-text-secondary">
                    {item.detail}
                  </span>
                </span>
                <ExternalLink className="mt-1 h-4 w-4 text-text-muted group-hover:text-primary" aria-hidden="true" />
              </a>
            ))}
          </div>

          <a
            href="https://status.raizhost.com"
            target="_blank"
            rel="noopener noreferrer"
            className="mt-2 flex min-h-11 items-center justify-between border-t border-border pt-5 text-sm text-text-secondary hover:text-primary"
          >
            <span>View service status</span>
            <ExternalLink className="h-4 w-4" aria-hidden="true" />
          </a>
        </aside>
      </div>
    </section>
  )
}
