import { ArrowRight, ExternalLink, Github } from 'lucide-react'
import { Link } from 'react-router-dom'

import { evidenceProjects, type EvidenceProject } from '../../../data/projects'

function ProjectRow({ project, index }: { project: EvidenceProject; index: number }) {
  return (
    <article className="grid gap-5 border-b border-border py-8 sm:py-10 lg:grid-cols-[8rem_minmax(0,1fr)_12rem] lg:gap-10 lg:py-12">
      <div>
        <span className="font-mono text-xs text-text-muted">
          {String(index + 1).padStart(2, '0')}
        </span>
        <p className="mt-2 font-mono text-[11px] uppercase tracking-[0.16em] text-primary">
          {project.discipline}
        </p>
      </div>

      <div className="max-w-3xl">
        <a
          href={project.repo_url}
          target="_blank"
          rel="noopener noreferrer"
          className="group inline-flex items-start gap-2"
        >
          <h3 className="font-display text-2xl font-bold tracking-[-0.03em] text-text-primary group-hover:text-primary sm:text-3xl lg:text-4xl">
            {project.name}
          </h3>
          <ExternalLink className="mt-1 h-4 w-4 text-text-muted group-hover:text-primary" aria-hidden="true" />
        </a>
        <p className="mt-4 text-[15px] leading-7 text-text-secondary sm:text-base sm:leading-8">
          {project.description}
        </p>
        <p className="mt-5 border-l-2 border-primary pl-4 text-[15px] leading-7 text-text-primary sm:text-base">
          <span className="font-semibold">Inspectable evidence:</span>{' '}
          {project.evidence}
        </p>
        <ul aria-label={`${project.name} technologies`} className="mt-5 flex flex-wrap gap-x-4 gap-y-2">
          {project.tags.map((tag) => (
            <li key={tag} className="font-mono text-[11px] uppercase tracking-[0.12em] text-text-muted">
              {tag}
            </li>
          ))}
        </ul>
      </div>

      <div className="flex items-start gap-3 lg:flex-col lg:items-stretch">
        <a
          href={project.repo_url}
          target="_blank"
          rel="noopener noreferrer"
          className="inline-flex min-h-11 flex-1 items-center justify-center gap-2 border border-border px-4 py-2 text-sm font-semibold text-text-primary hover:border-primary hover:text-primary"
        >
          <Github size={16} aria-hidden="true" />
          Repository
        </a>
        {project.live_url && (
          <a
            href={project.live_url}
            target="_blank"
            rel="noopener noreferrer"
            className="inline-flex min-h-11 flex-1 items-center justify-center gap-2 border border-border px-4 py-2 text-sm font-semibold text-text-secondary hover:border-primary hover:text-primary"
          >
            <ExternalLink size={16} aria-hidden="true" />
            Artifact / live
          </a>
        )}
      </div>
    </article>
  )
}

export default function HorizontalProjectGallery() {
  return (
    <section id="projects" aria-labelledby="projects-title" className="relative w-full border-b border-border py-16 sm:py-20 lg:py-28">
      <div className="portfolio-container">
        <div className="grid gap-6 pb-10 lg:grid-cols-[0.7fr_1.3fr] lg:items-end lg:pb-14">
          <div>
            <p className="mb-4 font-mono text-xs uppercase tracking-[0.2em] text-primary">
              Selected repositories
            </p>
            <h2 id="projects-title" className="font-display text-4xl font-bold tracking-[-0.04em] text-text-primary sm:text-5xl lg:text-6xl">
              Work you can inspect.
            </h2>
          </div>
          <p className="max-w-2xl text-base leading-7 text-text-secondary lg:justify-self-end lg:text-lg lg:leading-8">
            These projects lead with the engineering artifact: a failing case,
            a measured recovery, a gated write path, a signed release, or a
            documented constraint. Each link goes to the evidence rather than a mockup.
          </p>
        </div>

        <div className="border-t border-border">
          {evidenceProjects.map((project, index) => (
            <ProjectRow key={project.id} project={project} index={index} />
          ))}
        </div>

        <div className="mt-10 flex justify-end">
          <Link to="/projects" className="inline-flex min-h-11 items-center gap-2 text-sm font-semibold text-text-primary hover:text-primary">
            Browse the full repository index
            <ArrowRight size={17} aria-hidden="true" />
          </Link>
        </div>
      </div>
    </section>
  )
}
