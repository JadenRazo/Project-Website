import { useRef, useState, useEffect } from 'react'
import { motion, useInView } from 'framer-motion'
import { ExternalLink, Github, ArrowRight } from 'lucide-react'
import { mockProjects, ProjectData } from '../../../data/projects'
import CICDPipelineAnimation from '../../animations/CICDPipelineAnimation'

const componentMap: Record<string, React.FC> = {
  'cicd-pipeline': CICDPipelineAnimation,
}

interface ProjectCardProps {
  project: ProjectData
  index: number
}

function ProjectCard({ project, index }: ProjectCardProps) {
  const [isHovered, setIsHovered] = useState(false)
  const [isMobile, setIsMobile] = useState(false)
  const cardRef = useRef<HTMLDivElement>(null)
  const isInView = useInView(cardRef, { once: true, margin: '-50px' })

  useEffect(() => {
    const checkMobile = () => {
      setIsMobile(window.matchMedia('(max-width: 1024px)').matches || 'ontouchstart' in window)
    }
    checkMobile()
    window.addEventListener('resize', checkMobile)
    return () => window.removeEventListener('resize', checkMobile)
  }, [])

  const isEven = index % 2 === 0

  return (
    <motion.div
      ref={cardRef}
      className={`flex flex-col ${isEven ? 'lg:flex-row' : 'lg:flex-row-reverse'} gap-4 sm:gap-6 lg:gap-10 items-center`}
      initial={{ opacity: 0, y: 60 }}
      animate={isInView ? { opacity: 1, y: 0 } : {}}
      transition={{ duration: 0.8, delay: 0.1, ease: [0.22, 1, 0.36, 1] }}
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
    >
      <div className="w-full lg:w-1/2 relative group">
        <motion.div
          className="relative rounded-xl overflow-hidden glass-card"
          whileHover={{ scale: 1.02 }}
          transition={{ duration: 0.4, ease: [0.22, 1, 0.36, 1] }}
        >
          <div className="aspect-[16/10] w-full bg-surface overflow-hidden relative">
            {project.mediaType === 'component' && project.mediaUrl && componentMap[project.mediaUrl] ? (
              <div className="absolute inset-0">
                {(() => { const Comp = componentMap[project.mediaUrl!]; return <Comp />; })()}
              </div>
            ) : project.mediaType === 'video' ? (
              <video
                src={project.mediaUrl}
                autoPlay
                muted
                loop
                playsInline
                className="absolute inset-0 w-full h-full object-cover"
              />
            ) : (
              <div className="absolute inset-0 flex items-center justify-center p-6 bg-surface">
                <img
                  src={project.mediaUrl || '/images/projects/placeholder.jpg'}
                  alt={project.name}
                  className="max-w-full max-h-full object-contain transition-transform duration-700 group-hover:scale-105"
                />
              </div>
            )}
          </div>
          <div className={`absolute inset-0 bg-gradient-to-t from-background/80 via-background/20 to-transparent transition-opacity duration-300 ${isMobile ? 'opacity-60' : 'opacity-0 group-hover:opacity-60'}`} />

          <motion.div
            className={`absolute inset-0 flex items-end justify-center pb-4 ${isMobile ? '' : 'lg:items-center lg:pb-0'}`}
            initial={{ opacity: 0 }}
            animate={{ opacity: isMobile || isHovered ? 1 : 0 }}
            transition={{ duration: 0.3 }}
          >
            <div className="flex gap-2">
              {project.live_url && (
                <a
                  href={project.live_url}
                  target="_blank"
                  rel="noopener noreferrer"
                  aria-label={`View ${project.name} live demo`}
                  className="p-2.5 sm:p-3 rounded-full bg-primary text-white hover:bg-primary-light transition-colors active:scale-95"
                >
                  <ExternalLink size={18} />
                </a>
              )}
              {project.repo_url && (
                <a
                  href={project.repo_url}
                  target="_blank"
                  rel="noopener noreferrer"
                  aria-label={`View ${project.name} source code on GitHub`}
                  className="p-2.5 sm:p-3 rounded-full bg-surface/90 backdrop-blur-sm border border-border text-text-primary hover:border-primary transition-colors active:scale-95"
                >
                  <Github size={18} />
                </a>
              )}
            </div>
          </motion.div>
        </motion.div>

        <motion.div
          className="absolute -bottom-2 -right-2 sm:-bottom-3 sm:-right-3 w-10 h-10 sm:w-12 sm:h-12 rounded-full glass-enhanced flex items-center justify-center text-sm sm:text-base font-bold text-primary font-mono z-10"
          initial={{ scale: 0, rotate: -180 }}
          animate={isInView ? { scale: 1, rotate: 0 } : {}}
          transition={{ duration: 0.6, delay: 0.3, type: 'spring' }}
        >
          {String(index + 1).padStart(2, '0')}
        </motion.div>
      </div>

      <div className="w-full lg:w-1/2 space-y-3 sm:space-y-4">
        <motion.div
          initial={{ opacity: 0, x: isEven ? -30 : 30 }}
          animate={isInView ? { opacity: 1, x: 0 } : {}}
          transition={{ duration: 0.6, delay: 0.2 }}
        >
          <span className="text-[10px] sm:text-xs font-mono uppercase tracking-widest text-primary mb-2 block">
            Project {String(index + 1).padStart(2, '0')}
          </span>
          <h3 className="text-lg sm:text-xl md:text-2xl lg:text-3xl font-bold text-text-primary tracking-tight leading-tight">
            {project.name}
          </h3>
          {project.badges && project.badges.length > 0 && (
            <div className="flex flex-wrap gap-1.5 mt-2">
              {project.badges.map((badge) => (
                <span
                  key={badge}
                  className={`text-[10px] sm:text-xs px-2 py-0.5 rounded-full font-mono uppercase tracking-wider border ${
                    badge === 'client'
                      ? 'bg-emerald-500/20 text-emerald-400 border-emerald-500/30'
                      : badge === 'live'
                        ? 'bg-primary/20 text-primary border-primary/30'
                        : badge === 'demo'
                          ? 'bg-amber-500/20 text-amber-400 border-amber-500/30'
                          : 'bg-slate-500/20 text-slate-400 border-slate-500/30'
                  }`}
                >
                  {badge === 'client' ? 'Client Work' : badge === 'live' ? 'Live' : badge === 'demo' ? 'Demo' : 'Internal'}
                </span>
              ))}
            </div>
          )}
        </motion.div>

        <motion.p
          className="text-text-secondary text-xs sm:text-sm md:text-base leading-relaxed"
          initial={{ opacity: 0, x: isEven ? -30 : 30 }}
          animate={isInView ? { opacity: 1, x: 0 } : {}}
          transition={{ duration: 0.6, delay: 0.3 }}
        >
          {project.description}
        </motion.p>

        <motion.div
          className="flex flex-wrap gap-1.5 sm:gap-2"
          initial={{ opacity: 0, x: isEven ? -30 : 30 }}
          animate={isInView ? { opacity: 1, x: 0 } : {}}
          transition={{ duration: 0.6, delay: 0.4 }}
        >
          {project.tags.slice(0, 5).map((tag) => (
            <span
              key={tag}
              className="px-2 py-0.5 sm:px-2.5 sm:py-1 text-[10px] sm:text-xs font-mono bg-surface/80 backdrop-blur-sm rounded-full text-text-secondary border border-border hover:border-primary/50 transition-colors"
            >
              {tag}
            </span>
          ))}
          {project.tags.length > 5 && (
            <span className="px-2 py-0.5 sm:px-2.5 sm:py-1 text-[10px] sm:text-xs font-mono bg-surface/80 backdrop-blur-sm rounded-full text-text-muted border border-border">
              +{project.tags.length - 5}
            </span>
          )}
        </motion.div>

        <motion.div
          className="flex flex-wrap gap-2 pt-1"
          initial={{ opacity: 0, x: isEven ? -30 : 30 }}
          animate={isInView ? { opacity: 1, x: 0 } : {}}
          transition={{ duration: 0.6, delay: 0.5 }}
        >
          {project.live_url && (
            <a
              href={project.live_url}
              target="_blank"
              rel="noopener noreferrer"
              className="btn-primary text-sm py-2 px-4"
            >
              <span>View Live</span>
              <ExternalLink size={14} />
            </a>
          )}
          {project.repo_url && (
            <a
              href={project.repo_url}
              target="_blank"
              rel="noopener noreferrer"
              className="btn-secondary text-sm py-2 px-4"
            >
              <span>Source Code</span>
              <Github size={14} />
            </a>
          )}
        </motion.div>
      </div>
    </motion.div>
  )
}

export default function HorizontalProjectGallery() {
  const sectionRef = useRef<HTMLElement>(null)
  const headerRef = useRef<HTMLDivElement>(null)
  const isInView = useInView(headerRef, { once: true, margin: '-50px' })
  const projects = mockProjects.slice(0, 6)

  return (
    <section
      id="projects"
      ref={sectionRef}
      className="relative w-full py-12 sm:py-16 md:py-20"
    >
      <div className="w-full flex flex-col items-center px-4 sm:px-6 lg:px-8">
        {/* Header */}
        <div ref={headerRef} className="text-center mb-8 sm:mb-10 md:mb-12">
          <motion.span
            className="block text-[10px] sm:text-xs font-mono uppercase tracking-[0.2em] text-primary mb-2"
            initial={{ opacity: 0, y: 20 }}
            animate={isInView ? { opacity: 1, y: 0 } : {}}
            transition={{ duration: 0.6 }}
          >
            Selected Works
          </motion.span>
          <motion.h2
            className="text-2xl sm:text-3xl md:text-4xl lg:text-5xl font-bold text-text-primary tracking-tight"
            initial={{ opacity: 0, y: 30 }}
            animate={isInView ? { opacity: 1, y: 0 } : {}}
            transition={{ duration: 0.6, delay: 0.1 }}
          >
            Featured <span className="gradient-text">Projects</span>
          </motion.h2>
        </div>

        {/* Project cards */}
        <div className="w-full max-w-7xl">
          <div className="space-y-10 sm:space-y-12 md:space-y-16">
            {projects.map((project, index) => (
              <ProjectCard
                key={project.id}
                project={project}
                index={index}
              />
            ))}
          </div>
        </div>

        {/* View all button */}
        <motion.div
          className="mt-8 sm:mt-10 md:mt-12"
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          transition={{ duration: 0.6 }}
        >
          <a href="/projects" className="btn-secondary">
            <span>View All Projects</span>
            <ArrowRight size={18} />
          </a>
        </motion.div>
      </div>
    </section>
  )
}
