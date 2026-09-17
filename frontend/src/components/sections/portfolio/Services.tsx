import { useSiteContent } from '../../../lib/site-content'
import { Cloud, Gauge, GitPullRequest, ShieldCheck } from 'lucide-react'

const principleIcons = [Cloud, GitPullRequest, Gauge, ShieldCheck]

export default function Services() {
  const siteContent = useSiteContent()
  const copy = siteContent['services-home']
  const principles = [siteContent['principle-1'], siteContent['principle-2'], siteContent['principle-3'], siteContent['principle-4']].map((copy, index) => ({ ...copy, icon: principleIcons[index] }))
  return (
    <section id="services" aria-labelledby="services-title" className="relative w-full border-b border-border py-16 sm:py-20 lg:py-28">
      <div className="portfolio-container">
        <div className="grid gap-6 lg:grid-cols-[0.7fr_1.3fr] lg:items-end">
          <div>
            <p className="mb-4 font-mono text-xs uppercase tracking-[0.2em] text-primary">
              <span data-rh="services-home.eyebrow">{copy.eyebrow}</span>
            </p>
            <h2 id="services-title" className="font-display text-4xl font-bold tracking-[-0.04em] text-text-primary sm:text-5xl lg:text-6xl">
              <span data-rh="services-home.title">{copy.title}</span>
            </h2>
          </div>
          <p className="max-w-2xl text-base leading-7 text-text-secondary lg:justify-self-end lg:text-lg lg:leading-8">
            <span data-rh="services-home.intro">{copy.intro}</span>
          </p>
        </div>

        <div className="mt-12 grid border-l border-t border-border sm:grid-cols-2 lg:mt-16">
          {principles.map((principle, index) => (
            <article key={index} className="border-b border-r border-border p-6 sm:p-8 lg:p-10">
              <div className="flex items-center justify-between">
                <principle.icon className="h-5 w-5 text-primary" aria-hidden="true" />
                <span className="font-mono text-xs text-text-muted">
                  {String(index + 1).padStart(2, '0')}
                </span>
              </div>
              <h3 className="mt-8 text-xl font-semibold leading-7 text-text-primary sm:text-2xl">
                <span data-rh={`principle-${index + 1}.title`}>{principle.title}</span>
              </h3>
              <p className="mt-3 text-[15px] leading-7 text-text-secondary sm:text-base">
                <span data-rh={`principle-${index + 1}.description`}>{principle.description}</span>
              </p>
            </article>
          ))}
        </div>
      </div>
    </section>
  )
}
