import { useSiteContent } from "../../../lib/site-content";
import { useState } from "react";
import type { CSSProperties } from "react";
import { flushSync } from "react-dom";
import { ArrowRight, ExternalLink, Github } from "lucide-react";
import { Link } from "react-router-dom";
import { projectFilms } from "../../../data/projectFilms";
import ProjectFilmPlayer from "./ProjectFilmPlayer";
import "./project-films.css";

export default function HorizontalProjectGallery() {
  const siteContent = useSiteContent();
  const copy = siteContent['projects-home'];
  const settings = siteContent['project-settings'];
  const descriptions = {
    'raizhost': siteContent['film-raizhost'],
    'cloudcostmcp': siteContent['film-cloudcostmcp'],
    'tickethacker': siteContent['film-tickethacker'],
    'llm-lint': siteContent['film-llm-lint'],
    'sre-reference-app': siteContent['film-sre-reference-app'],
    'sre-landing-zone': siteContent['film-sre-landing-zone'],
  };
  const films = projectFilms.map(film => ({ ...film, ...descriptions[film.id as keyof typeof descriptions] }));
  const [activeFilm, setActiveFilm] = useState<string | null>(null);

  return (
    <section
      id="projects"
      aria-labelledby="projects-title"
      className="project-gallery"
    >
      <div className="portfolio-container">
        <div className="project-gallery-heading">
          <p data-rh="projects-home.eyebrow">{copy.eyebrow}</p>
          <h2 id="projects-title">
            <span data-rh="projects-home.title">{copy.title}</span>{" "}<span className="gradient-text" data-rh="projects-home.accent">{copy.accent}</span>
          </h2>
        </div>
        <div className="project-cards">
          {films.map((film, index) => (
            <article
              key={film.id}
              className="project-card"
              data-project={film.id}
              aria-labelledby={`project-title-${film.id}`}
              style={{ "--film-accent": film.color } as CSSProperties}
            >
              <div className="project-card-media">
                <ProjectFilmPlayer
                  film={film}
                  active={activeFilm === film.id}
                  onActivate={() => {
                    // Mount inside the play gesture to preserve mobile activation.
                    flushSync(() => setActiveFilm(film.id));
                  }}
                />
              </div>
              <div className="project-card-copy">
                <p className="project-card-number">
                  Project {String(index + 1).padStart(2, "0")}
                </p>
                <h3 id={`project-title-${film.id}`} data-rh={`film-${film.id}.title`}>{film.title}</h3>
                {settings.showCategories && <p className="project-card-category" data-rh={`film-${film.id}.category`}>{film.category}</p>}
                <p id={`film-summary-${film.id}`} data-rh={`film-${film.id}.summary`}>{film.summary}</p>
                {settings.showTags && <ul
                  data-rh-control="project-settings.showTags"
                  className="film-tags"
                  aria-label={`${film.title} technologies`}
                >
                  {film.tags.map((tag) => (
                    <li key={tag}>{tag}</li>
                  ))}
                </ul>}
                <div className="project-card-links">
                  {film.liveUrl && (
                    <a
                      href={film.liveUrl}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="btn-primary"
                    >
                      {film.liveLabel}{" "}
                      <ExternalLink size={15} aria-hidden="true" />
                    </a>
                  )}
                  <a
                    href={film.evidenceUrl}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="btn-secondary"
                  >
                    <Github size={15} aria-hidden="true" /> {film.evidenceLabel}
                  </a>
                </div>
                <a
                  href={`mailto:contact@jadenrazo.dev?subject=${encodeURIComponent(`Let’s talk about ${film.title}`)}`}
                  className="project-card-inquiry"
                >
                  <span data-rh="projects-home.inquiry">{copy.inquiry}</span>{" "}
                  <ArrowRight size={15} aria-hidden="true" />
                </a>
              </div>
            </article>
          ))}
        </div>
        <div className="project-gallery-footer">
          <Link to="/projects" className="btn-secondary">
            <span data-rh="projects-home.allProjects">{copy.allProjects}</span> <ArrowRight size={17} aria-hidden="true" />
          </Link>
        </div>
      </div>
    </section>
  );
}
