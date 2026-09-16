import { useRef, useState } from 'react'
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
  const [selected, setSelected] = useState(0)
  const tabs = useRef<(HTMLButtonElement | null)[]>([])
  const film = projectFilms[selected]
  const navigate = (event: KeyboardEvent<HTMLButtonElement>, index: number) => {
    const keys: Record<string, number> = {
      ArrowRight: (index + 1) % projectFilms.length,
      ArrowDown: (index + 1) % projectFilms.length,
      ArrowLeft: (index + projectFilms.length - 1) % projectFilms.length,
      ArrowUp: (index + projectFilms.length - 1) % projectFilms.length,
      Home: 0,
      End: projectFilms.length - 1,
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
            <h2 id="projects-title">Behind the build.</h2>
          </div>
          <p>
            Six projects. A closer look.
            <br />
            See how each one works in under a minute, then explore the code and
            decisions behind it.
          </p>
        </div>
        <div
          className="cinema-layout"
          style={{ '--film-accent': film.color } as CSSProperties}
        >
          <div className="cinema-selection">
            <p className="cinema-list-label">
              Choose a walkthrough{' '}
              <ArrowDownRight size={15} aria-hidden="true" />
            </p>
            <div
              className="cinema-tabs"
              role="tablist"
              aria-label="Project walkthroughs"
            >
              {projectFilms.map((project, index) => (
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
                    <strong>{project.title}</strong>
                    <span>{project.category}</span>
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
          {projectFilms.map((project, index) => (
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
                      <p className="film-category">{project.category}</p>
                      <h3>{project.headline}</h3>
                      <p id={`film-summary-${project.id}`}>{project.summary}</p>
                      <ul
                        className="film-tags"
                        aria-label={`${project.title} technologies`}
                      >
                        {project.tags.map((tag) => (
                          <li key={tag}>{tag}</li>
                        ))}
                      </ul>
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
          <span>Built by Jaden Razo. Open for a closer look.</span>
          <Link to="/projects">
            All repositories <ArrowRight size={16} aria-hidden="true" />
          </Link>
        </div>
      </div>
    </section>
  )
}
