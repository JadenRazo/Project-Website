import { useState, useEffect, useRef, useCallback } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { ExternalLink, Github, ChevronLeft, ChevronRight, Code, FileCode } from 'lucide-react'

interface Project {
  id: number
  title: string
  category: string
  description: string
  liveUrl: string
  githubUrl: string
  images: string[]
  technologies: string[]
  linesOfCode: number
  files: number
  mediaType?: 'image' | 'video'
  videoUrl?: string
}

const projects: Project[] = [
  {
    id: 1,
    title: 'Project Website',
    category: 'Full-Stack Portfolio',
    description: 'Personal portfolio and project showcase with Go microservices backend, React frontend, real-time messaging, and comprehensive admin panel.',
    liveUrl: 'https://jadenrazo.dev',
    githubUrl: 'https://github.com/jadenrazo/Project-Website',
    images: [
      'https://images.unsplash.com/photo-1507238691740-187a5b1d37b8?w=800&h=600&auto=format&fit=crop&crop=bottom&q=80',
    ],
    technologies: ['React', 'TypeScript', 'Go', 'PostgreSQL', 'Redis'],
    linesOfCode: 114000,
    files: 847,
  },
  {
    id: 2,
    title: 'Showers Auto Detail',
    category: 'Business Website',
    description: 'Professional auto detailing business website with service booking, gallery showcase, and customer management system.',
    liveUrl: 'https://showersautodetail.com',
    githubUrl: 'https://github.com/jadenrazo/showersautodetail',
    images: [
      'https://images.unsplash.com/photo-1520340356584-f9917d1eea6f?w=800&h=600&auto=format&fit=crop&crop=bottom&q=80',
      'https://images.unsplash.com/photo-1507136566006-cfc505b114fc?w=800&h=600&auto=format&fit=crop&crop=bottom&q=80',
      'https://images.unsplash.com/photo-1600880292203-757bb62b4baf?w=800&h=600&auto=format&fit=crop&crop=bottom&q=80',
      'https://images.unsplash.com/photo-1605559424843-9e4c228bf1c2?w=800&h=600&auto=format&fit=crop&crop=bottom&q=80',
    ],
    technologies: ['Next.js', 'TailwindCSS', 'Node.js', 'Stripe'],
    linesOfCode: 18500,
    files: 156,
  },
  {
    id: 3,
    title: 'Quiz Bot',
    category: 'Discord Bot',
    description: 'AI-powered Discord bot featuring interactive quizzes, trivia games, and educational content with real-time multiplayer support.',
    liveUrl: '#',
    githubUrl: 'https://github.com/jadenrazo/Quiz-Bot',
    images: [
      'https://images.unsplash.com/photo-1614680376739-414d95ff43df?w=800&h=600&auto=format&fit=crop&crop=bottom&q=80',
    ],
    technologies: ['Python', 'Discord.py', 'OpenAI', 'MongoDB'],
    linesOfCode: 8200,
    files: 42,
  },
  {
    id: 4,
    title: 'WeenieSMP',
    category: 'Full-Stack Minecraft Ecosystem',
    description: 'Production-grade Minecraft server ecosystem with Vue 3 e-commerce site, Go microservices for real-time stats, and comprehensive Docker infrastructure serving 12,000+ players.',
    liveUrl: 'https://weeniesmp.net',
    githubUrl: 'https://github.com/jadenrazo/Project-Website/tree/main/weeniesmp',
    images: ['/videos/weeniesmp_gambling_demo_optimized.mp4'],
    technologies: ['Vue 3', 'Go', 'MariaDB', 'Redis', 'Docker', 'Nginx'],
    linesOfCode: 45000,
    files: 238,
    mediaType: 'video',
    videoUrl: '/videos/weeniesmp_gambling_demo_optimized.mp4',
  },
]

function AnimatedCounter({ value, duration = 2000 }: { value: number; duration?: number }) {
  const [displayValue, setDisplayValue] = useState(0)
  const [hasAnimated, setHasAnimated] = useState(false)
  const ref = useRef<HTMLSpanElement>(null)

  useEffect(() => {
    const observer = new IntersectionObserver(
      (entries) => {
        if (entries[0].isIntersecting && !hasAnimated) {
          setHasAnimated(true)
          const startTime = Date.now()
          const animate = () => {
            const elapsed = Date.now() - startTime
            const progress = Math.min(elapsed / duration, 1)
            const easeOut = 1 - Math.pow(1 - progress, 3)
            setDisplayValue(Math.floor(value * easeOut))
            if (progress < 1) {
              requestAnimationFrame(animate)
            }
          }
          requestAnimationFrame(animate)
        }
      },
      { threshold: 0.5 }
    )

    if (ref.current) {
      observer.observe(ref.current)
    }

    return () => observer.disconnect()
  }, [value, duration, hasAnimated])

  useEffect(() => {
    setHasAnimated(false)
    setDisplayValue(0)
  }, [value])

  return (
    <span ref={ref} className="tabular-nums">
      {displayValue.toLocaleString()}
    </span>
  )
}

export default function Projects() {
  const [currentIndex, setCurrentIndex] = useState(0)
  const touchStartX = useRef(0)
  const touchStartY = useRef(0)
  const touchMoveX = useRef(0)
  const isHorizontal = useRef<boolean | null>(null)

  const currentProject = projects[currentIndex]

  const nextProject = useCallback(() => {
    setCurrentIndex((prev) => (prev + 1) % projects.length)
  }, [])

  const prevProject = useCallback(() => {
    setCurrentIndex((prev) => (prev - 1 + projects.length) % projects.length)
  }, [])

  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === 'ArrowRight') nextProject()
      if (e.key === 'ArrowLeft') prevProject()
    }
    window.addEventListener('keydown', handleKeyDown)
    return () => window.removeEventListener('keydown', handleKeyDown)
  }, [nextProject, prevProject])

  const handleTouchStart = (e: React.TouchEvent) => {
    touchStartX.current = e.touches[0].clientX
    touchStartY.current = e.touches[0].clientY
    touchMoveX.current = e.touches[0].clientX
    isHorizontal.current = null
  }

  const handleTouchMove = (e: React.TouchEvent) => {
    const currentX = e.touches[0].clientX
    const currentY = e.touches[0].clientY
    touchMoveX.current = currentX

    const diffX = Math.abs(currentX - touchStartX.current)
    const diffY = Math.abs(currentY - touchStartY.current)

    // Determine swipe direction after 15px movement
    if (isHorizontal.current === null && (diffX > 15 || diffY > 15)) {
      isHorizontal.current = diffX > diffY
    }
  }

  const handleTouchEnd = () => {
    // Only handle horizontal swipes
    if (isHorizontal.current !== true) {
      isHorizontal.current = null
      return
    }

    const diff = touchStartX.current - touchMoveX.current
    const threshold = 50

    if (Math.abs(diff) > threshold) {
      if (diff > 0) {
        nextProject()
      } else {
        prevProject()
      }
    }

    isHorizontal.current = null
  }

  return (
    <section id="projects" className="snap-section relative w-full bg-background">
      <div className="w-full h-full flex flex-col">
        {/* Section Header */}
        <div className="text-center pt-6 sm:pt-8 md:pt-10 pb-3 sm:pb-4 px-4 flex-shrink-0">
          <motion.h2
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            className="text-2xl sm:text-3xl md:text-4xl font-bold text-text-primary mb-2"
          >
            Featured <span className="gradient-text">Projects</span>
          </motion.h2>
          <motion.p
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            transition={{ delay: 0.1 }}
            className="text-text-secondary text-sm md:text-base text-center"
          >
            A selection of my recent work
          </motion.p>
        </div>

        {/* Project Content */}
        <div
          className="flex-1 relative min-h-0"
          onTouchStart={handleTouchStart}
          onTouchMove={handleTouchMove}
          onTouchEnd={handleTouchEnd}
        >
          <AnimatePresence mode="wait">
            <motion.div
              key={currentIndex}
              initial={{ opacity: 0, x: 60 }}
              animate={{ opacity: 1, x: 0 }}
              exit={{ opacity: 0, x: -60 }}
              transition={{ duration: 0.35, ease: 'easeOut' }}
              className="absolute inset-0 flex flex-col lg:flex-row"
            >
              {/* Image/Video Section */}
              <div className="relative w-full lg:w-1/2 aspect-video sm:aspect-[16/10] lg:h-full flex-shrink-0">
                {currentProject.mediaType === 'video' && currentProject.videoUrl ? (
                  <video
                    src={currentProject.videoUrl}
                    poster={currentProject.videoUrl.replace('.mp4', '_poster.jpg').replace('_optimized.mp4', '_poster.jpg')}
                    autoPlay
                    loop
                    muted
                    playsInline
                    preload="none"
                    className="w-full h-full object-cover"
                    style={{
                      WebkitBackfaceVisibility: 'hidden',
                      backfaceVisibility: 'hidden',
                      transform: 'translateZ(0)',
                    }}
                    onCanPlay={(e) => {
                      e.currentTarget.style.opacity = '1';
                    }}
                    aria-label={`${currentProject.title} demo video`}
                  />
                ) : (
                  <picture>
                    <source
                      srcSet={currentProject.images[0].replace('.jpg', '.webp').replace('.png', '.webp')}
                      type="image/webp"
                    />
                    <img
                      src={currentProject.images[0]}
                      alt={currentProject.title}
                      loading="lazy"
                      className="w-full h-full object-cover"
                    />
                  </picture>
                )}
                <div className="absolute inset-0 bg-gradient-to-t from-background/80 via-background/10 to-transparent lg:bg-gradient-to-r lg:from-transparent lg:via-background/10 lg:to-background/80" />

                {/* Project Counter */}
                <div className="absolute top-3 left-3 sm:top-4 sm:left-4 px-2.5 py-1 sm:px-3 sm:py-1.5 glass rounded-full text-xs sm:text-sm text-text-primary font-medium">
                  {currentIndex + 1} / {projects.length}
                </div>

                {/* Navigation Arrows on Image */}
                <button
                  onClick={prevProject}
                  className="absolute left-2 sm:left-3 top-1/2 -translate-y-1/2 p-2 sm:p-2.5 glass-card rounded-full z-10 active:scale-95 transition-transform"
                  aria-label="Previous project"
                >
                  <ChevronLeft className="w-4 h-4 sm:w-5 sm:h-5 text-text-primary" />
                </button>
                <button
                  onClick={nextProject}
                  className="absolute right-2 sm:right-3 top-1/2 -translate-y-1/2 p-2 sm:p-2.5 glass-card rounded-full z-10 active:scale-95 transition-transform"
                  aria-label="Next project"
                >
                  <ChevronRight className="w-4 h-4 sm:w-5 sm:h-5 text-text-primary" />
                </button>
              </div>

              {/* Content Section */}
              <div className="flex-1 lg:w-1/2 flex flex-col justify-center px-4 sm:px-6 lg:px-10 py-3 sm:py-4 lg:py-6 overflow-y-auto">
                <div className="max-w-lg mx-auto lg:mx-0 w-full">
                  {/* Category */}
                  <span className="inline-block px-2.5 py-1 glass rounded-full text-xs text-primary mb-2 sm:mb-3">
                    {currentProject.category}
                  </span>

                  {/* Title */}
                  <h3 className="text-lg sm:text-xl lg:text-3xl xl:text-4xl font-bold text-text-primary mb-1.5 sm:mb-2 lg:mb-3">
                    {currentProject.title}
                  </h3>

                  {/* Description */}
                  <p className="text-text-secondary text-xs sm:text-sm lg:text-base mb-3 sm:mb-4 leading-relaxed line-clamp-2 lg:line-clamp-none">
                    {currentProject.description}
                  </p>

                  {/* Stats */}
                  <div className="flex gap-2 sm:gap-3 mb-3 sm:mb-4">
                    <div className="glass-card px-2.5 py-1.5 sm:px-3 sm:py-2 lg:px-4 lg:py-3 flex items-center gap-2">
                      <div className="w-7 h-7 sm:w-8 sm:h-8 lg:w-10 lg:h-10 rounded-lg bg-primary/20 flex items-center justify-center">
                        <Code className="w-3.5 h-3.5 sm:w-4 sm:h-4 lg:w-5 lg:h-5 text-primary" />
                      </div>
                      <div>
                        <div className="text-sm sm:text-base lg:text-lg font-bold gradient-text">
                          <AnimatedCounter value={currentProject.linesOfCode} />
                        </div>
                        <div className="text-[10px] sm:text-xs text-text-muted">Lines</div>
                      </div>
                    </div>
                    <div className="glass-card px-2.5 py-1.5 sm:px-3 sm:py-2 lg:px-4 lg:py-3 flex items-center gap-2">
                      <div className="w-7 h-7 sm:w-8 sm:h-8 lg:w-10 lg:h-10 rounded-lg bg-accent/20 flex items-center justify-center">
                        <FileCode className="w-3.5 h-3.5 sm:w-4 sm:h-4 lg:w-5 lg:h-5 text-accent" />
                      </div>
                      <div>
                        <div className="text-sm sm:text-base lg:text-lg font-bold gradient-text">
                          <AnimatedCounter value={currentProject.files} />
                        </div>
                        <div className="text-[10px] sm:text-xs text-text-muted">Files</div>
                      </div>
                    </div>
                  </div>

                  {/* Technologies */}
                  <div className="flex flex-wrap gap-1 sm:gap-1.5 mb-3 sm:mb-4">
                    {currentProject.technologies.map((tech) => (
                      <span
                        key={tech}
                        className="px-2 py-0.5 sm:px-2.5 sm:py-1 text-[10px] sm:text-xs glass rounded-md sm:rounded-lg text-text-secondary"
                      >
                        {tech}
                      </span>
                    ))}
                  </div>

                  {/* Action Buttons */}
                  <div className="flex gap-2 sm:gap-3">
                    {currentProject.liveUrl !== '#' && (
                      <a
                        href={currentProject.liveUrl}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="btn-primary py-2 px-3 sm:py-2.5 sm:px-4 lg:py-3 lg:px-5 text-xs sm:text-sm"
                      >
                        <ExternalLink className="w-3.5 h-3.5 sm:w-4 sm:h-4" />
                        <span>View Live</span>
                      </a>
                    )}
                    <a
                      href={currentProject.githubUrl}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="btn-secondary py-2 px-3 sm:py-2.5 sm:px-4 lg:py-3 lg:px-5 text-xs sm:text-sm"
                    >
                      <Github className="w-3.5 h-3.5 sm:w-4 sm:h-4" />
                      <span>Source</span>
                    </a>
                  </div>
                </div>
              </div>
            </motion.div>
          </AnimatePresence>
        </div>

        {/* Pagination Dots */}
        <div className="flex justify-center items-center gap-2 sm:gap-3 py-3 sm:py-4 flex-shrink-0">
          {projects.map((_, index) => (
            <button
              key={index}
              onClick={() => setCurrentIndex(index)}
              className="p-1.5 sm:p-2"
              aria-label={`Go to project ${index + 1}`}
            >
              <span
                className={`block h-1.5 sm:h-2 rounded-full transition-all duration-300 ${
                  index === currentIndex
                    ? 'w-6 sm:w-8 bg-primary'
                    : 'w-1.5 sm:w-2 bg-surface hover:bg-surface-hover'
                }`}
              />
            </button>
          ))}
        </div>
      </div>
    </section>
  )
}
