import { useState, useEffect, useRef } from 'react'
import { motion, AnimatePresence } from 'framer-motion'

interface IntroAnimationProps {
  onComplete: () => void
}

const words = ['Build', 'Design', 'Create', 'Innovate']

export default function IntroAnimation({ onComplete }: IntroAnimationProps) {
  const [currentIndex, setCurrentIndex] = useState(0)
  const [isComplete, setIsComplete] = useState(false)
  const hasCompleted = useRef(false)

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

  useEffect(() => {
    if (isComplete && !hasCompleted.current) {
      const safetyTimeout = setTimeout(() => {
        if (!hasCompleted.current) {
          hasCompleted.current = true
          onComplete()
        }
      }, 1000)
      return () => clearTimeout(safetyTimeout)
    }
  }, [isComplete, onComplete])

  const handleAnimationComplete = () => {
    if (!hasCompleted.current) {
      hasCompleted.current = true
      onComplete()
    }
  }

  return (
    <AnimatePresence mode="wait" onExitComplete={handleAnimationComplete}>
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
