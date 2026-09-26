import { useSiteContent } from '../../../lib/site-content'
import { ExternalLink } from 'lucide-react'



export default function About() {
  const siteContent = useSiteContent()
  const copy = siteContent['about-home']
  const signals = siteContent['about-signals']
  return (
    <section id="about" aria-labelledby="about-title" className="relative w-full border-b border-border py-16 sm:py-20 lg:py-28">
      <div className="portfolio-container grid gap-12 lg:grid-cols-[0.72fr_1.28fr] lg:gap-20">
        <div>
          <p className="mb-4 font-mono text-xs uppercase tracking-[0.2em] text-primary">
            <span data-rh="about-home.eyebrow">{copy.eyebrow}</span>
          </p>
          <h2 id="about-title" className="font-display text-4xl font-bold leading-tight tracking-[-0.04em] text-text-primary sm:text-5xl lg:text-6xl">
            <span data-rh="about-home.title">{copy.title}</span>
          </h2>
          <p className="mt-6 max-w-xl text-base leading-7 text-text-secondary sm:text-lg sm:leading-8">
            <span data-rh="about-home.intro">{copy.intro}</span>
          </p>
          <p className="mt-4 max-w-xl text-[15px] leading-7 text-text-muted sm:text-base">
            <span data-rh="about-home.details">{copy.details}</span>
          </p>
        </div>

        <div className="border-t border-border">
          {signals.map((signal, index) => (
            <a
              key={signal.id}
              href={signal.href} data-rh-link={`about-signals.${index}.href`}
              target="_blank"
              rel="noopener noreferrer"
              className="group grid gap-3 border-b border-border py-6 sm:grid-cols-[3rem_1fr_auto] sm:gap-5 sm:py-7"
            >
              <span className="font-mono text-xs text-text-muted">
                {String(index + 1).padStart(2, '0')}
              </span>
              <span>
                <span className="font-mono text-[11px] uppercase tracking-[0.16em] text-primary">
                  <span data-rh={`about-signals.${index}.label`}>{signal.label}</span>
                </span>
                <span className="mt-2 block text-lg font-semibold leading-7 text-text-primary group-hover:text-primary sm:text-xl">
                  <span data-rh={`about-signals.${index}.title`}>{signal.title}</span>
                </span>
                <span className="mt-2 block text-[15px] leading-7 text-text-secondary sm:text-base">
                  <span data-rh={`about-signals.${index}.description`}>{signal.description}</span>
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
