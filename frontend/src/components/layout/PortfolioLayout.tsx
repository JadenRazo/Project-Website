import { useState, useEffect } from 'react'
import type { ReactNode } from 'react'
import { motion } from 'framer-motion'
import PortfolioNavbar from './PortfolioNavbar'
import WebGLBackground from '../animations/WebGLBackground'
import IntroAnimation from '../animations/IntroAnimation'
import { useLenis } from '../../providers/LenisProvider'
import { IntroContext } from '../../context/IntroContext'

interface PortfolioLayoutProps {
  children: ReactNode
}

export default function PortfolioLayout({ children }: PortfolioLayoutProps) {
  const [introComplete, setIntroComplete] = useState(false)
  const { scrollTo } = useLenis()

  useEffect(() => {
    if (introComplete) {
      scrollTo(0, { immediate: true })
    }
  }, [introComplete, scrollTo])

  return (
    <IntroContext.Provider value={introComplete}>
      <div className="relative min-h-screen overflow-x-hidden noise-overlay">
        <IntroAnimation onComplete={() => setIntroComplete(true)} />
        <WebGLBackground />
        <motion.div
          className="site-content relative z-10"
          initial={{ opacity: 0 }}
          animate={introComplete ? { opacity: 1 } : { opacity: 0 }}
          transition={{
            duration: 0.4,
            ease: [0.22, 1, 0.36, 1]
          }}
          style={{ pointerEvents: introComplete ? 'auto' : 'none' }}
        >
          <PortfolioNavbar />
          <main className="relative">
            {children}
          </main>
        </motion.div>
      </div>
    </IntroContext.Provider>
  )
}
