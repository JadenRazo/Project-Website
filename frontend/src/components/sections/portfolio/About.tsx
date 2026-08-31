import { ExternalLink } from 'lucide-react'

const signals = [
  {
    label: 'Incident analysis',
    title: 'A pricing defect became a regression suite.',
    description:
      'The CloudCostMCP incident record shows the bad selection logic, the detection gap, live provider verification, and the controls added after the fix.',
    href: 'https://github.com/JadenRazo/CloudCostMCP/blob/main/docs/incidents/2026-08-pricing-drift.md',
  },
  {
    label: 'Controlled failure',
    title: 'Recovery is measured, not implied.',
    description:
      'The SRE reference app records a 78-second ECS recovery exercise beside its SLO math, runbook, GameDay, postmortem, and explicit limitations.',
    href: 'https://github.com/JadenRazo/sre-reference-app',
  },
  {
    label: 'Release integrity',
    title: 'Artifacts carry their own evidence.',
    description:
      'llm-lint uses CodeQL, SARIF, release verification, SBOMs, native package checks, and signed release artifacts instead of asking users to trust a build script.',
    href: 'https://github.com/JadenRazo/llm-lint',
  },
]

export default function About() {
  return (
    <section id="about" aria-labelledby="about-title" className="relative w-full border-b border-border py-16 sm:py-20 lg:py-28">
      <div className="portfolio-container grid gap-12 lg:grid-cols-[0.72fr_1.28fr] lg:gap-20">
        <div>
          <p className="mb-4 font-mono text-xs uppercase tracking-[0.2em] text-primary">
            Engineering signal
          </p>
          <h2 id="about-title" className="font-display text-4xl font-bold leading-tight tracking-[-0.04em] text-text-primary sm:text-5xl lg:text-6xl">
            Evidence over adjectives.
          </h2>
          <p className="mt-6 max-w-xl text-base leading-7 text-text-secondary sm:text-lg sm:leading-8">
            I run RaizHost and build public cloud, reliability, and developer-tooling
            projects around one rule: architecture claims should resolve to code,
            tests, measurements, or an honest limitation.
          </p>
          <p className="mt-4 max-w-xl text-[15px] leading-7 text-text-muted sm:text-base">
            My strongest work sits where AWS infrastructure, secure delivery,
            observability, and operational ownership meet. I&apos;m CompTIA A+ and
            Network+ certified and currently focused on cloud, DevOps, platform,
            and SRE roles.
          </p>
        </div>

        <div className="border-t border-border">
          {signals.map((signal, index) => (
            <a
              key={signal.label}
              href={signal.href}
              target="_blank"
              rel="noopener noreferrer"
              className="group grid gap-3 border-b border-border py-6 sm:grid-cols-[3rem_1fr_auto] sm:gap-5 sm:py-7"
            >
              <span className="font-mono text-xs text-text-muted">
                {String(index + 1).padStart(2, '0')}
              </span>
              <span>
                <span className="font-mono text-[11px] uppercase tracking-[0.16em] text-primary">
                  {signal.label}
                </span>
                <span className="mt-2 block text-lg font-semibold leading-7 text-text-primary group-hover:text-primary sm:text-xl">
                  {signal.title}
                </span>
                <span className="mt-2 block text-[15px] leading-7 text-text-secondary sm:text-base">
                  {signal.description}
                </span>
              </span>
              <ExternalLink className="hidden h-4 w-4 text-text-muted group-hover:text-primary sm:block" aria-hidden="true" />
            </a>
          ))}
        </div>
      </div>
    </section>
  )
}
