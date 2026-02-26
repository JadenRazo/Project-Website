import { useState, useEffect } from 'react'
import { Code2 } from 'lucide-react'
import { motion } from 'framer-motion'
import { fetchGitHubLOC } from '../../utils/codeStats'

function AnimatedNumber({ value, duration = 2000 }: { value: number; duration?: number }) {
  const [displayValue, setDisplayValue] = useState(0)
  const [hasAnimated, setHasAnimated] = useState(false)

  useEffect(() => {
    if (hasAnimated) return

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
  }, [value, duration, hasAnimated])

  return <span className="tabular-nums">{displayValue.toLocaleString()}</span>
}

export default function CodeCounter() {
  const [totalLines, setTotalLines] = useState<number | null>(null)

  useEffect(() => {
    fetchGitHubLOC()
      .then((data) => setTotalLines(data.totalLines))
      .catch(() => {})
  }, [])

  if (!totalLines) return null

  return (
    <motion.div
      initial={{ opacity: 0, scale: 0.9 }}
      animate={{ opacity: 1, scale: 1 }}
      transition={{ duration: 0.5, delay: 0.2 }}
      className="flex items-center gap-2 px-2.5 sm:px-3 py-1 sm:py-1.5 glass rounded-full"
    >
      <div className="w-6 h-6 rounded-lg bg-primary/20 flex items-center justify-center">
        <Code2 className="w-3.5 h-3.5 text-primary" />
      </div>
      <div className="flex flex-col">
        <span className="text-xs font-bold gradient-text leading-none">
          <AnimatedNumber value={totalLines} />
        </span>
        <span className="text-[10px] text-text-muted leading-none mt-0.5">lines</span>
      </div>
    </motion.div>
  )
}
