import { useEffect, useState, useRef } from 'react'
import { motion, useSpring, useMotionValue, useTransform } from 'framer-motion'

interface CursorState {
  isHovering: boolean
  isHidden: boolean
}

export default function CustomCursor() {
  const [cursorState, setCursorState] = useState<CursorState>({
    isHovering: false,
    isHidden: false
  })
  const [isTouchDevice, setIsTouchDevice] = useState(false)
  const cursorOuterRef = useRef<HTMLDivElement>(null)
  const cursorInnerRef = useRef<HTMLDivElement>(null)

  const mouseX = useMotionValue(0)
  const mouseY = useMotionValue(0)

  const springConfig = { stiffness: 150, damping: 15, mass: 0.1 }
  const outerX = useSpring(mouseX, springConfig)
  const outerY = useSpring(mouseY, springConfig)

  const innerX = useSpring(mouseX, { stiffness: 500, damping: 28, mass: 0.1 })
  const innerY = useSpring(mouseY, { stiffness: 500, damping: 28, mass: 0.1 })

  const outerScale = useTransform(
    [outerX, outerY],
    () => cursorState.isHovering ? 2 : 1
  )

  useEffect(() => {
    const checkTouchDevice = () => {
      setIsTouchDevice(
        'ontouchstart' in window ||
        navigator.maxTouchPoints > 0 ||
        window.matchMedia('(pointer: coarse)').matches
      )
    }

    checkTouchDevice()

    const prefersReducedMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches
    if (prefersReducedMotion || isTouchDevice) {
      return
    }

    const handleMouseMove = (e: MouseEvent) => {
      mouseX.set(e.clientX)
      mouseY.set(e.clientY)
    }

    const handleMouseEnter = () => {
      setCursorState(prev => ({ ...prev, isHidden: false }))
    }

    const handleMouseLeave = () => {
      setCursorState(prev => ({ ...prev, isHidden: true }))
    }

    const interactiveElements = document.querySelectorAll(
      'a, button, [role="button"], input, textarea, select, [data-cursor-hover]'
    )

    const handleElementEnter = () => {
      setCursorState(prev => ({ ...prev, isHovering: true }))
    }

    const handleElementLeave = () => {
      setCursorState(prev => ({ ...prev, isHovering: false }))
    }

    window.addEventListener('mousemove', handleMouseMove)
    document.body.addEventListener('mouseenter', handleMouseEnter)
    document.body.addEventListener('mouseleave', handleMouseLeave)

    interactiveElements.forEach(el => {
      el.addEventListener('mouseenter', handleElementEnter)
      el.addEventListener('mouseleave', handleElementLeave)
    })

    const observer = new MutationObserver(() => {
      const newElements = document.querySelectorAll(
        'a:not([data-cursor-observed]), button:not([data-cursor-observed]), [role="button"]:not([data-cursor-observed])'
      )
      newElements.forEach(el => {
        el.setAttribute('data-cursor-observed', 'true')
        el.addEventListener('mouseenter', handleElementEnter)
        el.addEventListener('mouseleave', handleElementLeave)
      })
    })

    observer.observe(document.body, { childList: true, subtree: true })

    return () => {
      window.removeEventListener('mousemove', handleMouseMove)
      document.body.removeEventListener('mouseenter', handleMouseEnter)
      document.body.removeEventListener('mouseleave', handleMouseLeave)
      interactiveElements.forEach(el => {
        el.removeEventListener('mouseenter', handleElementEnter)
        el.removeEventListener('mouseleave', handleElementLeave)
      })
      observer.disconnect()
    }
  }, [mouseX, mouseY, isTouchDevice])

  if (isTouchDevice) {
    return null
  }

  return (
    <>
      <motion.div
        ref={cursorOuterRef}
        className="fixed pointer-events-none z-[9999] mix-blend-difference"
        style={{
          x: outerX,
          y: outerY,
          translateX: '-50%',
          translateY: '-50%',
        }}
        animate={{
          opacity: cursorState.isHidden ? 0 : 1,
          scale: cursorState.isHovering ? 2 : 1,
        }}
        transition={{
          opacity: { duration: 0.2 },
          scale: { type: 'spring', stiffness: 400, damping: 30 }
        }}
      >
        <div
          className="w-10 h-10 rounded-full border-2 border-white transition-colors duration-200"
          style={{
            borderColor: cursorState.isHovering ? 'rgba(255,255,255,0.5)' : 'white'
          }}
        />
      </motion.div>
      <motion.div
        ref={cursorInnerRef}
        className="fixed pointer-events-none z-[9999] mix-blend-difference"
        style={{
          x: innerX,
          y: innerY,
          translateX: '-50%',
          translateY: '-50%',
        }}
        animate={{
          opacity: cursorState.isHidden ? 0 : 1,
          scale: cursorState.isHovering ? 0.5 : 1,
        }}
        transition={{
          opacity: { duration: 0.2 },
          scale: { type: 'spring', stiffness: 400, damping: 30 }
        }}
      >
        <div className="w-2 h-2 rounded-full bg-white" />
      </motion.div>
    </>
  )
}
