import { useEffect, useRef, useState } from "react";
import { Play, RotateCcw, VolumeX } from "lucide-react";
import type { ProjectFilm } from "../../../data/projectFilms";

export const filmTime = (seconds: number) =>
  `${Math.floor(seconds / 60)}:${String(Math.floor(seconds % 60)).padStart(2, "0")}`;

export default function ProjectFilmPlayer({
  film,
  active,
  onActivate,
}: {
  film: ProjectFilm;
  active: boolean;
  onActivate: () => void;
}) {
  const videoRef = useRef<HTMLVideoElement>(null);
  const frameRef = useRef<HTMLDivElement>(null);
  const requested = useRef(false);
  const pendingSeek = useRef<number | null>(null);
  const attempt = useRef(0);
  const playbackAllowed = useRef(false);
  const [mode, setMode] = useState<
    "poster" | "loading" | "playing" | "paused" | "ended" | "error"
  >("poster");
  const [time, setTime] = useState(0);
  const [slow, setSlow] = useState(false);
  const [composition, setComposition] = useState<"portrait" | "wide">();
  const [source, setSource] = useState(film.src);

  useEffect(() => {
    if (!active) {
      requested.current = false;
      pendingSeek.current = null;
      playbackAllowed.current = false;
      setMode("poster");
      setTime(0);
      setComposition(undefined);
      setSource(film.src);
      return;
    }
    const video = videoRef.current;
    const frame = frameRef.current;
    if (!video || !frame) return;
    let settledFrame = 0;
    const allowNativeControlsAfterPaint = () => {
      cancelAnimationFrame(settledFrame);
      // WebKit can issue automatic play during re-entry. Preserve the pause
      // through that paint, then accept native controls even when their events
      // stay inside the browser's own shadow tree.
      settledFrame = requestAnimationFrame(() => {
        settledFrame = requestAnimationFrame(() => {
          const bounds = frame.getBoundingClientRect();
          if (
            !document.hidden &&
            bounds.bottom > 0 &&
            bounds.top < window.innerHeight
          )
            playbackAllowed.current = true;
        });
      });
    };
    const pause = () => {
      cancelAnimationFrame(settledFrame);
      playbackAllowed.current = false;
      attempt.current += 1;
      video.pause();
      if (requested.current)
        setMode((current) =>
          current === "error" || current === "ended" ? current : "paused",
        );
    };
    const onVisibility = () => {
      if (document.hidden) pause();
      else allowNativeControlsAfterPaint();
    };
    const observer = new IntersectionObserver(([entry]) => {
      if (!entry.isIntersecting) pause();
      else allowNativeControlsAfterPaint();
    });
    observer.observe(frame);
    document.addEventListener("visibilitychange", onVisibility);
    return () => {
      observer.disconnect();
      cancelAnimationFrame(settledFrame);
      document.removeEventListener("visibilitychange", onVisibility);
      attempt.current += 1;
      video.pause();
      video.removeAttribute("src");
      video.load();
    };
  }, [active, film.src]);

  useEffect(() => {
    if (mode !== "loading") {
      setSlow(false);
      return;
    }
    const timer = window.setTimeout(() => setSlow(true), 10000);
    return () => window.clearTimeout(timer);
  }, [mode]);

  const play = (start?: number, transferFocus = false) => {
    onActivate();
    const video = videoRef.current;
    if (!video) return;
    const bounds = frameRef.current?.getBoundingClientRect();
    if (bounds && (bounds.bottom <= 0 || bounds.top >= window.innerHeight)) {
      frameRef.current?.scrollIntoView({
        block: "nearest",
        behavior: "instant",
      });
    }
    const currentAttempt = ++attempt.current;
    playbackAllowed.current = true;
    pendingSeek.current = start ?? (mode === "ended" ? 0 : null);
    if (!requested.current || mode === "error") {
      // Attach one composition in the user gesture; no video requests on load.
      const portrait = window.matchMedia("(max-width: 639px)").matches;
      const mp4 = video.canPlayType('video/mp4; codecs="avc1.640028"') !== "";
      const selectedSource = portrait
        ? mp4
          ? film.mobileSrc
          : film.mobileWebmSrc
        : mp4
          ? film.src
          : film.webmSrc;
      video.src = selectedSource;
      setComposition(portrait ? "portrait" : "wide");
      setSource(selectedSource);
      requested.current = true;
      video.load();
    }
    if (video.readyState >= 1 && pendingSeek.current !== null) {
      video.currentTime = pendingSeek.current;
      pendingSeek.current = null;
    }
    setMode("loading");
    // The play/retry overlay disappears. Keep keyboard focus on the controls.
    if (transferFocus)
      requestAnimationFrame(() => video.focus({ preventScroll: true }));
    void video.play().catch((error: DOMException) => {
      if (attempt.current !== currentAttempt) return;
      setMode(
        error.name === "NotAllowedError" || error.name === "AbortError"
          ? "paused"
          : "error",
      );
    });
  };
  const chapterIndex = Math.max(
    0,
    film.chapters.reduce(
      (last, chapter, index) => (time >= chapter.at ? index : last),
      0,
    ),
  );
  const guardPlayback = (video: HTMLVideoElement) => {
    const bounds = frameRef.current?.getBoundingClientRect();
    if (
      !playbackAllowed.current ||
      document.hidden ||
      (bounds && (bounds.bottom <= 0 || bounds.top >= window.innerHeight))
    ) {
      playbackAllowed.current = false;
      video.pause();
      return false;
    }
    return true;
  };

  return (
    <div className="project-film-player">
      <div className="film-toolbar">
        <span>{film.kind}</span>
        <span>
          <VolumeX size={14} aria-hidden="true" /> Sound-free ·{" "}
          {filmTime(film.duration)}
        </span>
      </div>
      <div className="film-frame" ref={frameRef} data-composition={composition}>
        {active && (
          <video
            ref={videoRef}
            className={
              mode === "poster" ? "film-video is-unstarted" : "film-video"
            }
            controls={mode !== "poster" && mode !== "error"}
            preload="none"
            playsInline
            muted
            tabIndex={0}
            width="1600"
            height="1000"
            aria-label={`${film.title} walkthrough`}
            aria-describedby={`film-summary-${film.id}`}
            onLoadedMetadata={() => {
              if (pendingSeek.current !== null && videoRef.current) {
                videoRef.current.currentTime = pendingSeek.current;
                pendingSeek.current = null;
              }
            }}
            onTimeUpdate={(event) => {
              const video = event.currentTarget;
              setTime(video.currentTime);
              // Some WebKit seeks advance without another `playing` event.
              if (!video.paused && video.readyState >= 3)
                setMode((current) =>
                  current === "loading" ? "playing" : current,
                );
            }}
            onPointerDownCapture={() => {
              playbackAllowed.current = true;
            }}
            onKeyDownCapture={() => {
              playbackAllowed.current = true;
            }}
            onPlay={(event) => {
              guardPlayback(event.currentTarget);
            }}
            onPlaying={(event) => {
              if (guardPlayback(event.currentTarget)) setMode("playing");
            }}
            onPause={() => {
              const bounds = frameRef.current?.getBoundingClientRect();
              if (
                document.hidden ||
                (bounds &&
                  (bounds.bottom <= 0 || bounds.top >= window.innerHeight))
              )
                playbackAllowed.current = false;
              if (requested.current)
                setMode((current) =>
                  current === "error" || current === "ended"
                    ? current
                    : "paused",
                );
            }}
            onWaiting={() => {
              if (requested.current && !videoRef.current?.paused)
                setMode("loading");
            }}
            onEnded={() => setMode("ended")}
            onError={() => {
              if (requested.current) setMode("error");
            }}
          >
            <track
              kind="captions"
              src={film.captions}
              srcLang="en"
              label="English descriptions"
            />
            Your browser cannot play this video. Read the walkthrough below.
          </video>
        )}
        {(!active || mode === "poster") && (
          <>
            <picture className="film-poster">
              <source media="(max-width: 639px)" srcSet={film.mobilePoster} />
              <img
                src={film.poster}
                alt=""
                width="1600"
                height="1000"
                loading="lazy"
                decoding="async"
              />
            </picture>
            <button
              className="film-play"
              onClick={() => play(undefined, true)}
              aria-label={`Play ${film.title} walkthrough, ${film.duration} seconds`}
            >
              <span className="film-play-icon">
                <Play size={21} fill="currentColor" aria-hidden="true" />
              </span>
              <span>
                Watch walkthrough{" "}
                <span className="film-play-time">
                  {filmTime(film.duration)}
                </span>
              </span>
            </button>
          </>
        )}
        {mode === "error" && (
          <div className="film-recovery" role="status">
            <p>The video couldn’t load.</p>
            <p>You can retry, open the video, or read the walkthrough below.</p>
            <button onClick={() => play(time, true)}>
              <RotateCcw size={16} aria-hidden="true" /> Try again
            </button>
            <a href={source} target="_blank" rel="noopener noreferrer">
              Open video
            </a>
          </div>
        )}
        {mode === "ended" && (
          <button
            className="film-play film-replay"
            onClick={() => play(0, true)}
          >
            <RotateCcw size={19} aria-hidden="true" /> Watch again
          </button>
        )}
        {mode === "loading" && (
          <div className="film-loading" role="status">
            {slow
              ? "Taking a moment. The text walkthrough is available below."
              : "Loading video…"}
          </div>
        )}
      </div>
      <div
        className="film-chapters"
        aria-label={`${film.title} video chapters`}
      >
        {film.chapters.map((chapter, index) => (
          <button
            key={chapter.title}
            onClick={() => play(chapter.at)}
            className={index === chapterIndex ? "is-current" : ""}
            aria-label={`Play chapter ${index + 1}: ${chapter.title}, at ${filmTime(chapter.at)}`}
            aria-current={
              index === chapterIndex && mode !== "poster" ? "step" : undefined
            }
          >
            <span>{filmTime(chapter.at)}</span>
            <strong>{chapter.title}</strong>
          </button>
        ))}
      </div>
      <details className="film-transcript">
        <summary>Read the walkthrough</summary>
        <ol>
          {film.chapters.map((chapter) => (
            <li key={chapter.title}>
              <strong>{chapter.title}.</strong> {chapter.description}
            </li>
          ))}
        </ol>
        <p>{film.sourceNote}</p>
      </details>
    </div>
  );
}
