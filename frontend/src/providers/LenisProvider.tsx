import { createContext, useContext, useEffect, useRef, useCallback, ReactNode, useState } from 'react'
import { useLocation } from 'react-router-dom'
import Lenis from 'lenis'
import { gsap } from 'gsap'
import { ScrollTrigger } from 'gsap/ScrollTrigger'
import { setLenisInstance, resetScrollLock } from '../utils/scrollLock'

gsap.registerPlugin(ScrollTrigger)

interface LenisContextValue {
  lenis: Lenis | null
  scrollTo: (target: string | number | HTMLElement, options?: { immediate?: boolean }) => void
}

const LenisContext = createContext<LenisContextValue>({
  lenis: null,
  scrollTo: () => {}
})

export const useLenis = () => useContext(LenisContext)

interface LenisProviderProps {
  children: ReactNode
}

function clearScrollLockState() {
  document.documentElement.classList.remove('lenis-stopped', 'scroll-unlocked', 'scroll-locked')
  document.body.style.overflow = ''
  document.body.style.position = ''
  document.body.style.top = ''
  document.body.style.width = ''
  resetScrollLock()
}

export default function LenisProvider({ children }: LenisProviderProps) {
  const lenisRef = useRef<Lenis | null>(null)
  const rafIdRef = useRef<number | null>(null)
  const [isReady, setIsReady] = useState(false)
  const location = useLocation()
  const isLandingPage = location.pathname === '/'

  useEffect(() => {
    clearScrollLockState()

    const prefersReducedMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches

    if (prefersReducedMotion) {
      setIsReady(true)
      return
    }

    if (!isLandingPage) {
      lenisRef.current = null
      setLenisInstance(null)
      setIsReady(true)
      return
    }

    const lenis = new Lenis({
      duration: 1.2,
      easing: (t) => Math.min(1, 1.001 - Math.pow(2, -10 * t)),
      orientation: 'vertical',
      gestureOrientation: 'vertical',
      smoothWheel: true,
      wheelMultiplier: 1,
      touchMultiplier: 2,
      infinite: false,
      autoResize: true,
    })

    lenisRef.current = lenis
    setLenisInstance(lenis)

    lenis.on('scroll', ScrollTrigger.update)

    const raf = (time: number) => {
      lenis.raf(time)
      rafIdRef.current = requestAnimationFrame(raf)
    }

    rafIdRef.current = requestAnimationFrame(raf)

    lenis.stop()
    document.documentElement.classList.add('scroll-locked')

    const observer = new MutationObserver(() => {
      if (document.documentElement.classList.contains('scroll-unlocked')) {
        document.documentElement.classList.remove('scroll-locked')
        lenis.scrollTo(0, { immediate: true })
        lenis.start()
        ScrollTrigger.refresh()
        setIsReady(true)
        observer.disconnect()
      }
    })

    observer.observe(document.documentElement, {
      attributes: true,
      attributeFilter: ['class'],
    })

    return () => {
      observer.disconnect()
      if (rafIdRef.current) {
        cancelAnimationFrame(rafIdRef.current)
      }
      lenis.destroy()
      clearScrollLockState()
      lenisRef.current = null
      setLenisInstance(null)
    }
  }, [isLandingPage])

  const scrollTo = useCallback((target: string | number | HTMLElement, options?: { immediate?: boolean }) => {
    if (lenisRef.current) {
      lenisRef.current.scrollTo(target, {
        offset: 0,
        immediate: options?.immediate ?? false,
        duration: options?.immediate ? 0 : 1.2,
        easing: (t) => Math.min(1, 1.001 - Math.pow(2, -10 * t))
      })
    } else {
      if (typeof target === 'string') {
        document.querySelector(target)?.scrollIntoView({ behavior: 'smooth' })
      } else if (typeof target === 'number') {
        window.scrollTo({ top: target, behavior: 'smooth' })
      } else if (target instanceof HTMLElement) {
        target.scrollIntoView({ behavior: 'smooth' })
      }
    }
  }, [])

  return (
    <LenisContext.Provider value={{ lenis: lenisRef.current, scrollTo }}>
      {children}
    </LenisContext.Provider>
  )
}
