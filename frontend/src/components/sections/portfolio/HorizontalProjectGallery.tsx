import { useState } from "react";
import type { CSSProperties } from "react";
import { flushSync } from "react-dom";
import { ArrowRight, ExternalLink, Github } from "lucide-react";
import { Link } from "react-router-dom";
import { projectFilms } from "../../../data/projectFilms";
import ProjectFilmPlayer from "./ProjectFilmPlayer";
import "./project-films.css";

export default function HorizontalProjectGallery() {
  const [activeFilm, setActiveFilm] = useState<string | null>(null);

  return (
    <section
      id="projects"
      aria-labelledby="projects-title"
      className="project-gallery"
    >
      <div className="portfolio-container">
        <div className="project-gallery-heading">
          <p>Selected Works</p>
          <h2 id="projects-title">
            Featured <span className="gradient-text">Projects</span>
          </h2>
        </div>
        <div className="project-cards">
          {projectFilms.map((film, index) => (
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
                <h3 id={`project-title-${film.id}`}>{film.title}</h3>
                <p className="project-card-category">{film.category}</p>
                <p id={`film-summary-${film.id}`}>{film.summary}</p>
                <ul
                  className="film-tags"
                  aria-label={`${film.title} technologies`}
                >
                  {film.tags.map((tag) => (
                    <li key={tag}>{tag}</li>
                  ))}
                </ul>
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
                  Discuss this project{" "}
                  <ArrowRight size={15} aria-hidden="true" />
                </a>
              </div>
            </article>
          ))}
        </div>
        <div className="project-gallery-footer">
          <Link to="/projects" className="btn-secondary">
            View All Projects <ArrowRight size={17} aria-hidden="true" />
          </Link>
        </div>
      </div>
    </section>
  );
}
