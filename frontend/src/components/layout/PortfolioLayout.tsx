import { useState } from 'react'
import type { ReactNode } from 'react'
import { motion } from 'framer-motion'
import PortfolioNavbar from './PortfolioNavbar'
import WebGLBackground from '../animations/WebGLBackground'
import IntroAnimation from '../animations/IntroAnimation'

interface PortfolioLayoutProps {
  children: ReactNode
}

export default function PortfolioLayout({ children }: PortfolioLayoutProps) {
  const [introComplete, setIntroComplete] = useState(false)

  return (
    <div className="relative min-h-screen overflow-x-hidden noise-overlay">
      <IntroAnimation onComplete={() => setIntroComplete(true)} />
      <WebGLBackground />
      <motion.div
        className="site-content relative z-10"
        initial={{ y: -20, opacity: 0 }}
        animate={{
          y: introComplete ? 0 : -20,
          opacity: introComplete ? 1 : 0
        }}
        transition={{
          duration: 0.8,
          ease: [0.22, 1, 0.36, 1]
        }}
      >
        <PortfolioNavbar />
        <main className="relative">
          {children}
        </main>
      </motion.div>
    </div>
  )
}
