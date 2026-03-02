import { useState, useEffect, useRef, useCallback } from 'react'
import { motion, AnimatePresence } from 'framer-motion'

interface IntroAnimationProps {
  onComplete: () => void
}

const words = ['Build', 'Design', 'Create', 'Innovate']

function lockScroll() {
  document.documentElement.style.overflow = 'hidden'
  document.body.style.overflow = 'hidden'
  document.documentElement.style.position = 'fixed'
  document.documentElement.style.width = '100%'
  document.documentElement.style.height = '100%'
  document.body.style.position = 'fixed'
  document.body.style.width = '100%'
  document.body.style.height = '100%'
  document.body.style.top = '0'
}

function unlockScroll() {
  document.documentElement.style.overflow = ''
  document.body.style.overflow = ''
  document.documentElement.style.position = ''
  document.documentElement.style.width = ''
  document.documentElement.style.height = ''
  document.body.style.position = ''
  document.body.style.width = ''
  document.body.style.height = ''
  document.body.style.top = ''
  window.scrollTo(0, 0)
}

export default function IntroAnimation({ onComplete }: IntroAnimationProps) {
  const [currentIndex, setCurrentIndex] = useState(0)
  const [isComplete, setIsComplete] = useState(false)
  const hasCompleted = useRef(false)

  useEffect(() => {
    lockScroll()

    const preventDefault = (e: Event) => e.preventDefault()
    document.addEventListener('touchmove', preventDefault, { passive: false })
    document.addEventListener('wheel', preventDefault, { passive: false })

    return () => {
      document.removeEventListener('touchmove', preventDefault)
      document.removeEventListener('wheel', preventDefault)
      unlockScroll()
    }
  }, [])

  useEffect(() => {
    if (currentIndex < words.length - 1) {
      const timer = setTimeout(() => {
        setCurrentIndex((prev) => prev + 1)
      }, 1100)
      return () => clearTimeout(timer)
    } else {
      const completeTimer = setTimeout(() => {
        setIsComplete(true)
      }, 1200)
      return () => clearTimeout(completeTimer)
    }
  }, [currentIndex])

  const finishIntro = useCallback(() => {
    if (!hasCompleted.current) {
      hasCompleted.current = true
      unlockScroll()
      onComplete()
    }
  }, [onComplete])

  useEffect(() => {
    if (isComplete && !hasCompleted.current) {
      const safetyTimeout = setTimeout(finishIntro, 1000)
      return () => clearTimeout(safetyTimeout)
    }
  }, [isComplete, finishIntro])

  return (
    <AnimatePresence mode="wait" onExitComplete={finishIntro}>
      {!isComplete && (
        <motion.div
          key="intro-overlay"
          className="intro-overlay"
          initial={{ opacity: 1 }}
          exit={{
            opacity: 0,
            transition: {
              duration: 0.6,
              ease: 'easeInOut'
            }
          }}
        >
          <div className="intro-content">
            <AnimatePresence mode="wait">
              <motion.span
                key={words[currentIndex]}
                className="intro-word gradient-text"
                initial={{
                  opacity: 0,
                  y: 40
                }}
                animate={{
                  opacity: 1,
                  y: 0
                }}
                exit={{
                  opacity: 0,
                  y: -40
                }}
                transition={{
                  duration: 0.5,
                  ease: [0.22, 1, 0.36, 1]
                }}
              >
                {words[currentIndex]}
              </motion.span>
            </AnimatePresence>

            <motion.div
              className="intro-progress"
              initial={{ scaleX: 0 }}
              animate={{ scaleX: (currentIndex + 1) / words.length }}
              transition={{ duration: 0.3, ease: 'easeOut' }}
            />
          </div>
        </motion.div>
      )}
    </AnimatePresence>
  )
}
