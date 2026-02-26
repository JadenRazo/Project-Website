import { useEffect, useRef, useState, useMemo } from 'react'
import { motion } from 'framer-motion'
import { gsap } from 'gsap'
import { ArrowDown } from 'lucide-react'
import GlassIcosahedron from '../../3d/GlassIcosahedron'
import { useLenis } from '../../../providers/LenisProvider'

interface SplitTextProps {
  text: string
  className?: string
  delay?: number
}

function SplitText({ text, className = '', delay = 0 }: SplitTextProps) {
  const containerRef = useRef<HTMLSpanElement>(null)
  const charsRef = useRef<(HTMLSpanElement | null)[]>([])

  useEffect(() => {
    const prefersReducedMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches
    if (prefersReducedMotion) return

    const chars = charsRef.current.filter(Boolean)

    gsap.set(chars, {
      y: '100%',
      opacity: 0,
      rotationX: -80,
    })

    gsap.to(chars, {
      y: '0%',
      opacity: 1,
      rotationX: 0,
      duration: 0.8,
      ease: 'power3.out',
      stagger: 0.025,
      delay: delay,
    })
  }, [delay])

  const chars = useMemo(() => text.split(''), [text])

  return (
    <span ref={containerRef} className={`inline-block overflow-hidden perspective-1000 ${className}`}>
      {chars.map((char, i) => (
        <span
          key={i}
          ref={(el) => (charsRef.current[i] = el)}
          className="inline-block preserve-3d"
          style={{ transformOrigin: 'center bottom' }}
        >
          {char === ' ' ? '\u00A0' : char}
        </span>
      ))}
    </span>
  )
}

const roles = ['websites', 'experiences', 'interfaces', 'solutions']

function AnimatedRole() {
  const [currentIndex, setCurrentIndex] = useState(0)
  const [isAnimating, setIsAnimating] = useState(false)
  const textRef = useRef<HTMLSpanElement>(null)

  useEffect(() => {
    const interval = setInterval(() => {
      setIsAnimating(true)

      setTimeout(() => {
        setCurrentIndex((prev) => (prev + 1) % roles.length)
        setIsAnimating(false)
      }, 400)
    }, 3000)

    return () => clearInterval(interval)
  }, [])

  return (
    <span className="inline-block relative overflow-hidden">
      <motion.span
        ref={textRef}
        className="inline-block gradient-text-hero"
        animate={{
          y: isAnimating ? '-100%' : '0%',
          opacity: isAnimating ? 0 : 1,
        }}
        transition={{
          duration: 0.4,
          ease: [0.22, 1, 0.36, 1],
        }}
      >
        {roles[currentIndex]}
      </motion.span>
    </span>
  )
}

export default function Hero() {
  const containerRef = useRef<HTMLDivElement>(null)
  const { scrollTo } = useLenis()

  const handleScrollToAbout = () => {
    scrollTo('#about')
  }

  const handleScrollToProjects = () => {
    scrollTo('#projects')
  }

  const handleScrollToContact = () => {
    scrollTo('#contact')
  }

  return (
    <section
      id="home"
      ref={containerRef}
      className="relative w-full min-h-[85vh] sm:min-h-[88vh] overflow-hidden"
    >
      {/* 3D Elements - desktop only */}
      <div className="absolute right-[3%] top-1/2 -translate-y-1/2 hidden xl:block opacity-50 pointer-events-none">
        <GlassIcosahedron size={320} />
      </div>

      {/* Main content - centered vertically */}
      <div className="relative z-10 w-full flex flex-col items-center pt-[28vh] sm:pt-[30vh] lg:pt-[32vh] pb-16 sm:pb-20 px-4 sm:px-6 lg:px-8">
        {/* Content wrapper */}
        <div className="w-full max-w-4xl text-center">
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.6, delay: 0.2 }}
              className="mb-2 sm:mb-4"
            >
              <span className="inline-block px-4 py-2 glass-enhanced rounded-full text-xs sm:text-sm font-mono uppercase tracking-widest text-primary">
                Available for freelance work
              </span>
            </motion.div>

            <h1 className="text-editorial mb-3 sm:mb-5">
              <motion.div
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                transition={{ duration: 0.3, delay: 0.4 }}
                className="overflow-hidden"
              >
                <SplitText
                  text="I craft beautiful"
                  className="text-[2.5rem] sm:text-5xl md:text-6xl lg:text-7xl text-text-primary block leading-[1.1]"
                  delay={0.5}
                />
              </motion.div>
              <motion.div
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                transition={{ duration: 0.3, delay: 0.8 }}
                className="overflow-hidden"
              >
                <span className="text-[2.5rem] sm:text-5xl md:text-6xl lg:text-7xl leading-[1.1]">
                  <AnimatedRole />
                </span>
              </motion.div>
            </h1>

            <div className="w-full flex justify-center">
              <motion.p
                initial={{ opacity: 0, y: 30 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.6, delay: 1.2, ease: [0.22, 1, 0.36, 1] }}
                className="text-base sm:text-lg md:text-xl text-text-secondary max-w-lg mb-4 sm:mb-6 leading-relaxed text-center"
              >
                Full-stack developer specializing in stunning, high-performance web applications with modern technologies.
              </motion.p>
            </div>

            <motion.div
              initial={{ opacity: 0, y: 30 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.6, delay: 1.4, ease: [0.22, 1, 0.36, 1] }}
              className="flex flex-wrap gap-3 justify-center"
            >
              <button
                className="btn-primary"
                onClick={handleScrollToProjects}
              >
                <span>View My Work</span>
                <ArrowDown size={16} />
              </button>
              <button
                className="btn-secondary"
                onClick={handleScrollToContact}
              >
                <span>Get In Touch</span>
              </button>
            </motion.div>

            <motion.div
              initial={{ opacity: 0, y: 30 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.6, delay: 1.6, ease: [0.22, 1, 0.36, 1] }}
              className="flex gap-4 justify-center mt-4 sm:mt-6 items-center"
            >
              <a
                href="https://github.com/JadenRazo"
                target="_blank"
                rel="noopener noreferrer"
                className="group flex items-center gap-2 text-text-secondary hover:text-primary transition-colors duration-300"
                aria-label="GitHub Profile"
              >
                <svg
                  className="w-5 h-5 transition-transform duration-300 group-hover:scale-110"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                >
                  <path d="M9 19c-5 1.5-5-2.5-7-3m14 6v-3.87a3.37 3.37 0 0 0-.94-2.61c3.14-.35 6.44-1.54 6.44-7A5.44 5.44 0 0 0 20 4.77 5.07 5.07 0 0 0 19.91 1S18.73.65 16 2.48a13.38 13.38 0 0 0-7 0C6.27.65 5.09 1 5.09 1A5.07 5.07 0 0 0 5 4.77a5.44 5.44 0 0 0-1.5 3.78c0 5.42 3.3 6.61 6.44 7A3.37 3.37 0 0 0 9 18.13V22" />
                </svg>
                <span className="text-xs sm:text-sm font-medium">GitHub</span>
              </a>
              <div className="w-px h-5 bg-border" />
              <a
                href="https://jadenrazo.dev/s/linkedin"
                target="_blank"
                rel="noopener noreferrer"
                className="group flex items-center gap-2 text-text-secondary hover:text-primary transition-colors duration-300"
                aria-label="LinkedIn Profile"
              >
                <svg
                  className="w-5 h-5 transition-transform duration-300 group-hover:scale-110"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                >
                  <path d="M16 8a6 6 0 0 1 6 6v7h-4v-7a2 2 0 0 0-2-2 2 2 0 0 0-2 2v7h-4v-7a6 6 0 0 1 6-6z" />
                  <rect x="2" y="9" width="4" height="12" />
                  <circle cx="4" cy="4" r="2" />
                </svg>
                <span className="text-xs sm:text-sm font-medium">LinkedIn</span>
              </a>
              <div className="w-px h-5 bg-border" />
              <a
                href="mailto:contact@jadenrazo.dev"
                className="group flex items-center gap-2 text-text-secondary hover:text-primary transition-colors duration-300"
                aria-label="Email Contact"
              >
                <svg
                  className="w-5 h-5 transition-transform duration-300 group-hover:scale-110"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                >
                  <path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z" />
                  <polyline points="22,6 12,13 2,6" />
                </svg>
                <span className="text-xs sm:text-sm font-medium">Email</span>
              </a>
            </motion.div>

            {/* Stats row */}
            <motion.div
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              transition={{ delay: 2, duration: 0.6 }}
              className="mt-4 sm:mt-8 flex items-center justify-center gap-5 sm:gap-6"
            >
              <div className="flex flex-col">
                <span className="text-xl sm:text-2xl font-bold gradient-text">50+</span>
                <span className="text-[10px] sm:text-xs text-text-muted uppercase tracking-wider">Projects</span>
              </div>
              <div className="w-px h-8 bg-border" />
              <div className="flex flex-col">
                <span className="text-xl sm:text-2xl font-bold gradient-text">5+</span>
                <span className="text-[10px] sm:text-xs text-text-muted uppercase tracking-wider">Years</span>
              </div>
              <div className="w-px h-8 bg-border" />
              <div className="flex flex-col">
                <span className="text-xl sm:text-2xl font-bold gradient-text">100%</span>
                <span className="text-[10px] sm:text-xs text-text-muted uppercase tracking-wider">Quality</span>
              </div>
            </motion.div>
          </div>

          {/* Scroll indicator */}
          <motion.button
            onClick={handleScrollToAbout}
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 2, duration: 0.6 }}
            className="mt-8 sm:mt-10 flex flex-col items-center gap-2 group cursor-pointer z-20"
            whileHover={{ scale: 1.05 }}
            whileTap={{ scale: 0.95 }}
          >
            <span className="text-[10px] sm:text-xs font-mono uppercase tracking-widest text-text-muted group-hover:text-primary transition-colors">
              Scroll
            </span>
            <motion.div
              animate={{ y: [0, 8, 0] }}
              transition={{ duration: 1.5, repeat: Infinity, ease: 'easeInOut' }}
              className="p-2 glass-enhanced rounded-full"
            >
              <ArrowDown size={20} className="text-primary group-hover:text-accent transition-colors" />
            </motion.div>
          </motion.button>

        </div>

      {/* Bottom gradient */}
      <div className="absolute bottom-0 left-0 right-0 h-24 bg-gradient-to-t from-background to-transparent pointer-events-none" />
    </section>
  )
}
