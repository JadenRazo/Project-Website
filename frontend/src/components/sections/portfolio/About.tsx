import { useEffect, useRef, useState } from 'react'
import { motion, useInView } from 'framer-motion'

const stats = [
  { value: 5, suffix: '+', label: 'Years Experience' },
  { value: 30, suffix: '+', label: 'Projects Delivered' },
  { value: 15, suffix: '+', label: 'Clients Served' },
  { value: 100, suffix: '%', label: 'Satisfaction Rate' },
]

function AnimatedStat({ value, suffix, label, animate }: { value: number; suffix: string; label: string; animate: boolean }) {
  const [displayValue, setDisplayValue] = useState(0)
  const hasAnimated = useRef(false)

  useEffect(() => {
    if (animate && !hasAnimated.current) {
      hasAnimated.current = true
      const duration = 2000
      const startTime = Date.now()

      const animateValue = () => {
        const elapsed = Date.now() - startTime
        const progress = Math.min(elapsed / duration, 1)
        const easeOut = 1 - Math.pow(1 - progress, 3)
        setDisplayValue(Math.floor(value * easeOut))

        if (progress < 1) {
          requestAnimationFrame(animateValue)
        }
      }

      requestAnimationFrame(animateValue)
    }
  }, [animate, value])

  return (
    <div className="text-center p-2 sm:p-3">
      <div className="text-xl sm:text-2xl md:text-3xl lg:text-4xl font-bold gradient-text mb-1">
        {displayValue}{suffix}
      </div>
      <div className="text-text-secondary text-[10px] sm:text-xs md:text-sm">{label}</div>
    </div>
  )
}

export default function About() {
  const sectionRef = useRef<HTMLElement>(null)
  const isInView = useInView(sectionRef, { once: true, margin: '-100px' })

  return (
    <section ref={sectionRef} id="about" className="relative w-full pt-0 pb-12 sm:pb-16 md:pb-20">
      <div className="w-full flex flex-col items-center px-4 sm:px-6 lg:px-8">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={isInView ? { opacity: 1, y: 0 } : {}}
          transition={{ duration: 0.6 }}
          className="mb-8 sm:mb-10 text-center max-w-3xl"
        >
          <h2 className="text-2xl sm:text-3xl md:text-4xl font-bold text-text-primary mb-3">
            About <span className="gradient-text">Me</span>
          </h2>
          <p className="text-sm sm:text-base md:text-lg text-text-secondary leading-relaxed">
            Full-stack developer passionate about crafting exceptional digital experiences.
            I specialize in building high-performance web applications with modern technologies,
            seamless user experiences, and scalable architecture.
          </p>
        </motion.div>

        <motion.div
          initial={{ opacity: 0, y: 40 }}
          animate={isInView ? { opacity: 1, y: 0 } : {}}
          transition={{ duration: 0.6, delay: 0.3 }}
          className="grid grid-cols-2 md:grid-cols-4 gap-2 sm:gap-4 w-full max-w-4xl"
        >
          {stats.map((stat) => (
            <AnimatedStat
              key={stat.label}
              value={stat.value}
              suffix={stat.suffix}
              label={stat.label}
              animate={isInView}
            />
          ))}
        </motion.div>
      </div>
    </section>
  )
}
