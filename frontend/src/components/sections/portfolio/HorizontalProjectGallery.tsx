import { useSiteContent } from '../../../lib/site-content'
import { useEffect, useRef, useState } from 'react'
import type { CSSProperties, KeyboardEvent } from 'react'
import {
  ArrowDownRight,
  ArrowRight,
  ExternalLink,
  Github,
  Play,
} from 'lucide-react'
import { Link } from 'react-router-dom'
import { projectFilms } from '../../../data/projectFilms'
import ProjectFilmPlayer, { filmTime } from './ProjectFilmPlayer'
import './project-films.css'

export default function HorizontalProjectGallery() {
  const siteContent = useSiteContent()
  const copy = siteContent['projects-home']
  const descriptions = {
    'raizhost': siteContent['film-raizhost'],
    'cloudcostmcp': siteContent['film-cloudcostmcp'],
    'tickethacker': siteContent['film-tickethacker'],
    'llm-lint': siteContent['film-llm-lint'],
    'sre-reference-app': siteContent['film-sre-reference-app'],
    'sre-landing-zone': siteContent['film-sre-landing-zone'],
  }
  const films = projectFilms.map((film) => ({ ...film, ...descriptions[film.id as keyof typeof descriptions] }))
  const settings = siteContent['project-settings']
  const defaultIndex = Math.max(0, films.findIndex(film => film.id === settings.defaultFilm))
  const [selected, setSelected] = useState(defaultIndex)
  useEffect(() => { setSelected(defaultIndex) }, [defaultIndex])
  const tabs = useRef<(HTMLButtonElement | null)[]>([])
  const film = films[selected]
  const navigate = (event: KeyboardEvent<HTMLButtonElement>, index: number) => {
    const keys: Record<string, number> = {
      ArrowRight: (index + 1) % films.length,
      ArrowDown: (index + 1) % films.length,
      ArrowLeft: (index + films.length - 1) % films.length,
      ArrowUp: (index + films.length - 1) % films.length,
      Home: 0,
      End: films.length - 1,
    }
    if (!(event.key in keys)) return
    event.preventDefault()
    setSelected(keys[event.key])
    tabs.current[keys[event.key]]?.focus({ preventScroll: true })
  }
  return (
    <section
      id="projects"
      aria-labelledby="projects-title"
      className="project-cinema"
    >
      <div className="portfolio-container">
        <div className="cinema-heading">
          <div>
            <p className="cinema-eyebrow">
              Selected work <span>01 — 06</span>
            </p>
            <h2 id="projects-title"><span data-rh="projects-home.title">{copy.title}</span></h2>
          </div>
          <p>
            <span data-rh="projects-home.intro">{copy.intro}</span>
            <br />
            <span data-rh="projects-home.detail">{copy.detail}</span>
          </p>
        </div>
        <div
          className="cinema-layout"
          data-rh-control="project-settings.defaultFilm"
          style={{ '--film-accent': film.color } as CSSProperties}
        >
          <div className="cinema-selection">
            <p className="cinema-list-label">
              <span data-rh="projects-home.selection">{copy.selection}</span>{' '}
              <ArrowDownRight size={15} aria-hidden="true" />
            </p>
            <div
              className="cinema-tabs"
              role="tablist"
              aria-label="Project walkthroughs"
            >
              {films.map((project, index) => (
                <button
                  type="button"
                  key={project.id}
                  ref={(element) => {
                    tabs.current[index] = element
                  }}
                  id={`film-tab-${project.id}`}
                  role="tab"
                  aria-selected={index === selected}
                  aria-controls={`film-panel-${project.id}`}
                  tabIndex={index === selected ? 0 : -1}
                  onClick={() => setSelected(index)}
                  onKeyDown={(event) => navigate(event, index)}
                  className={
                    index === selected ? 'cinema-tab is-selected' : 'cinema-tab'
                  }
                >
                  <span className="cinema-tab-index">
                    {String(index + 1).padStart(2, '0')}
                  </span>
                  <span className="cinema-tab-copy">
                    <strong data-rh={`film-${project.id}.title`}>{project.title}</strong>
                    <span data-rh={`film-${project.id}.category`}>{project.category}</span>
                  </span>
                  <span className="cinema-tab-time">
                    {index === selected ? (
                      <Play size={13} fill="currentColor" aria-hidden="true" />
                    ) : (
                      filmTime(project.duration)
                    )}
                  </span>
                </button>
              ))}
            </div>
            <p className="cinema-note">
              Short, captioned, and sound-free.
              <br />
              Press play when you’re ready.
            </p>
          </div>
          {films.map((project, index) => (
            <div
              key={project.id}
              id={`film-panel-${project.id}`}
              role="tabpanel"
              aria-labelledby={`film-tab-${project.id}`}
              hidden={index !== selected}
              tabIndex={0}
              className="cinema-panel"
            >
              {index === selected && (
                <>
                  <ProjectFilmPlayer key={project.id} film={project} />
                  <div className="film-details">
                    <div>
                      <p className="film-category" data-rh={`film-${project.id}.category`}>{project.category}</p>
                      <h3 data-rh={`film-${project.id}.headline`}>{project.headline}</h3>
                      <p id={`film-summary-${project.id}`} data-rh={`film-${project.id}.summary`}>{project.summary}</p>
                      {settings.showTags && <ul
                        data-rh-control="project-settings.showTags"
                        className="film-tags"
                        aria-label={`${project.title} technologies`}
                      >
                        {project.tags.map((tag) => (
                          <li key={tag}>{tag}</li>
                        ))}
                      </ul>}
                    </div>
                    <div className="film-links">
                      <a
                        href={project.evidenceUrl}
                        target="_blank"
                        rel="noopener noreferrer"
                      >
                        <Github size={16} aria-hidden="true" />
                        {project.evidenceLabel}
                        <ArrowRight size={16} aria-hidden="true" />
                      </a>
                      {project.liveUrl && (
                        <a
                          href={project.liveUrl}
                          target="_blank"
                          rel="noopener noreferrer"
                        >
                          <ExternalLink size={16} aria-hidden="true" />
                          {project.liveLabel}
                        </a>
                      )}
                      <a
                        href={`mailto:contact@jadenrazo.dev?subject=${encodeURIComponent(`Let’s talk about ${project.title}`)}`}
                        className="film-inquiry"
                      >
                        Discuss this project{' '}
                        <ArrowRight size={15} aria-hidden="true" />
                      </a>
                    </div>
                  </div>
                </>
              )}
            </div>
          ))}
        </div>
        <div className="cinema-footer">
          <span><span data-rh="projects-home.footer">{copy.footer}</span></span>
          <Link to="/projects">
            <span data-rh="projects-home.allProjects">{copy.allProjects}</span> <ArrowRight size={16} aria-hidden="true" />
          </Link>
        </div>
      </div>
    </section>
  )
}
