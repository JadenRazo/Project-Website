import { useEffect, useRef, useState, useMemo } from 'react'
import { motion, useAnimation } from 'framer-motion'
import { gsap } from 'gsap'
import { ArrowDown } from 'lucide-react'
import GlassIcosahedron from '../../3d/GlassIcosahedron'
import { useNavigate } from 'react-router-dom'
import { useLenis } from '../../../providers/LenisProvider'
import { useIntroComplete } from '../../../context/IntroContext'

interface SplitTextProps {
  text: string
  className?: string
  delay?: number
  ready?: boolean
}

function SplitText({ text, className = '', delay = 0, ready = false }: SplitTextProps) {
  const containerRef = useRef<HTMLSpanElement>(null)
  const charsRef = useRef<(HTMLSpanElement | null)[]>([])
  const hasAnimated = useRef(false)
  const tweenRef = useRef<gsap.core.Tween | null>(null)

  useEffect(() => {
    if (!ready || hasAnimated.current) return
    hasAnimated.current = true

    const prefersReducedMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches
    const chars = charsRef.current.filter(Boolean)

    if (prefersReducedMotion) {
      gsap.set(chars, { y: '0%', opacity: 1, rotationX: 0 })
      return
    }

    gsap.set(chars, {
      y: '100%',
      opacity: 0,
      rotationX: -80,
    })

    tweenRef.current = gsap.to(chars, {
      y: '0%',
      opacity: 1,
      rotationX: 0,
      duration: 0.5,
      ease: 'power3.out',
      stagger: 0.012,
      delay: delay,
    })

    return () => {
      tweenRef.current?.kill()
      tweenRef.current = null
    }
  }, [delay, ready])

  useEffect(() => {
    if (ready) return
    const chars = charsRef.current.filter(Boolean)
    gsap.set(chars, { y: '100%', opacity: 0, rotationX: -80 })
    return () => {
      gsap.killTweensOf(chars)
    }
  }, [ready])

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

function AnimatedRole({ ready }: { ready: boolean }) {
  const [currentIndex, setCurrentIndex] = useState(0)
  const [isAnimating, setIsAnimating] = useState(false)

  useEffect(() => {
    if (!ready) return
    let pendingTimeout: ReturnType<typeof setTimeout> | null = null

    const interval = setInterval(() => {
      setIsAnimating(true)

      pendingTimeout = setTimeout(() => {
        setCurrentIndex((prev) => (prev + 1) % roles.length)
        setIsAnimating(false)
        pendingTimeout = null
      }, 400)
    }, 3000)

    return () => {
      clearInterval(interval)
      if (pendingTimeout !== null) {
        clearTimeout(pendingTimeout)
      }
    }
  }, [ready])

  return (
    <span className="inline-block relative overflow-hidden">
      <motion.span
        className="inline-block gradient-text-hero"
        initial={{ opacity: 0 }}
        animate={ready ? {
          y: isAnimating ? '-100%' : '0%',
          opacity: isAnimating ? 0 : 1,
        } : { opacity: 0 }}
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
  const navigate = useNavigate()
  const { scrollTo } = useLenis()
  const introComplete = useIntroComplete()
  const controls = useAnimation()

  useEffect(() => {
    if (introComplete) {
      controls.start('visible')
    }
  }, [introComplete, controls])

  const handleScrollToAbout = () => {
    navigate('/about')
  }

  const handleScrollToProjects = () => {
    scrollTo('#projects')
  }

  const handleScrollToContact = () => {
    navigate('/contact')
  }

  const item = (delay: number) => ({
    hidden: { opacity: 0, y: 30 },
    visible: {
      opacity: 1,
      y: 0,
      transition: { duration: 0.4, delay, ease: [0.22, 1, 0.36, 1] },
    },
  })

  const fade = (delay: number) => ({
    hidden: { opacity: 0 },
    visible: {
      opacity: 1,
      transition: { duration: 0.2, delay },
    },
  })

  return (
    <section
      id="home"
      ref={containerRef}
      className="relative w-full h-[100dvh] overflow-hidden"
    >
      <div className="absolute right-[5%] top-1/2 -translate-y-1/2 hidden xl:block xl:opacity-40 2xl:opacity-50 pointer-events-none transition-opacity duration-700">
        <GlassIcosahedron size={380} />
      </div>

      <div className="relative z-10 w-full h-full flex flex-col items-center justify-center px-4 sm:px-6 lg:px-8 xl:px-12 pb-16">
        <div className="w-full max-w-4xl xl:max-w-5xl text-center">
            <motion.div
              variants={item(0.0)}
              initial="hidden"
              animate={controls}
              className="mb-2 sm:mb-4"
            >
              <span className="inline-block px-4 py-2 glass-enhanced rounded-full text-xs sm:text-sm font-mono uppercase tracking-widest text-primary">
                Available for freelance work
              </span>
            </motion.div>

            <h1 className="text-editorial mb-3 sm:mb-5">
              <motion.div
                variants={fade(0.1)}
                initial="hidden"
                animate={controls}
                className="overflow-hidden"
              >
                <SplitText
                  text="I craft beautiful"
                  className="text-[2.5rem] sm:text-5xl md:text-6xl lg:text-7xl xl:text-[5.5rem] 2xl:text-[6rem] text-text-primary block leading-[1.05]"
                  delay={0.15}
                  ready={introComplete}
                />
              </motion.div>
              <motion.div
                variants={fade(0.3)}
                initial="hidden"
                animate={controls}
                className="overflow-hidden"
              >
                <span className="text-[2.5rem] sm:text-5xl md:text-6xl lg:text-7xl xl:text-[5.5rem] 2xl:text-[6rem] leading-[1.05]">
                  <AnimatedRole ready={introComplete} />
                </span>
              </motion.div>
            </h1>

            <div className="w-full flex justify-center">
              <motion.p
                variants={item(0.5)}
                initial="hidden"
                animate={controls}
                className="text-base sm:text-lg md:text-xl lg:text-[1.35rem] text-text-secondary max-w-lg lg:max-w-2xl mb-4 sm:mb-6 lg:mb-8 leading-relaxed text-center"
              >
                Full-stack developer specializing in stunning, high-performance web applications with modern technologies.
              </motion.p>
            </div>

            <motion.div
              variants={item(0.6)}
              initial="hidden"
              animate={controls}
              className="flex flex-wrap gap-3 lg:gap-4 justify-center"
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
              variants={item(0.7)}
              initial="hidden"
              animate={controls}
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

            <motion.div
              variants={fade(0.8)}
              initial="hidden"
              animate={controls}
              className="mt-4 sm:mt-8 lg:mt-10 flex items-center justify-center gap-5 sm:gap-6 lg:gap-10"
            >
              <div className="flex flex-col">
                <span className="text-xl sm:text-2xl lg:text-3xl font-bold gradient-text tabular-nums">50+</span>
                <span className="text-[10px] sm:text-xs lg:text-sm text-text-muted uppercase tracking-wider">Projects</span>
              </div>
              <div className="w-px h-8 lg:h-10 bg-border" />
              <div className="flex flex-col">
                <span className="text-xl sm:text-2xl lg:text-3xl font-bold gradient-text tabular-nums">5+</span>
                <span className="text-[10px] sm:text-xs lg:text-sm text-text-muted uppercase tracking-wider">Years</span>
              </div>
              <div className="w-px h-8 lg:h-10 bg-border" />
              <div className="flex flex-col">
                <span className="text-xl sm:text-2xl lg:text-3xl font-bold gradient-text tabular-nums">100%</span>
                <span className="text-[10px] sm:text-xs lg:text-sm text-text-muted uppercase tracking-wider">Quality</span>
              </div>
            </motion.div>
          </div>

        </div>

      <motion.button
        onClick={handleScrollToAbout}
        variants={fade(0.9)}
        initial="hidden"
        animate={controls}
        className="absolute bottom-6 sm:bottom-8 left-1/2 -translate-x-1/2 flex flex-col items-center gap-2 group cursor-pointer z-20"
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

    </section>
  )
}
