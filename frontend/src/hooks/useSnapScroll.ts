import { useEffect, useRef, useState, useCallback } from 'react'

interface UseSnapScrollOptions {
  sectionSelector?: string
  transitionDuration?: number
  swipeThreshold?: number
}

export function useSnapScroll(
  containerRef: React.RefObject<HTMLElement | null>,
  options: UseSnapScrollOptions = {}
) {
  const {
    sectionSelector = '.snap-section',
    transitionDuration = 800,
    swipeThreshold = 50
  } = options
  const [currentIndex, setCurrentIndex] = useState(0)
  const isScrolling = useRef(false)
  const touchStartY = useRef(0)
  const touchStartX = useRef(0)
  const touchCurrentY = useRef(0)
  const touchCurrentX = useRef(0)
  const isVerticalSwipe = useRef<boolean | null>(null)
  const sectionsCount = useRef(0)

  const scrollToSection = useCallback((index: number, sections: NodeListOf<Element>) => {
    if (index < 0 || index >= sections.length) return false
    if (isScrolling.current) return false

    isScrolling.current = true
    setCurrentIndex(index)

    const section = sections[index] as HTMLElement
    section.scrollIntoView({ behavior: 'smooth', block: 'start' })

    setTimeout(() => {
      isScrolling.current = false
    }, transitionDuration)

    return true
  }, [transitionDuration])

  useEffect(() => {
    const container = containerRef.current
    if (!container) return

    const sections = container.querySelectorAll(sectionSelector)
    if (sections.length === 0) return
    sectionsCount.current = sections.length

    const handleWheel = (e: WheelEvent) => {
      if (isScrolling.current) {
        e.preventDefault()
        return
      }

      const atFirstSection = currentIndex === 0
      const atLastSection = currentIndex === sections.length - 1
      const scrollingUp = e.deltaY < 0
      const scrollingDown = e.deltaY > 0

      if ((atFirstSection && scrollingUp) || (atLastSection && scrollingDown)) {
        return
      }

      e.preventDefault()

      if (scrollingDown) {
        scrollToSection(currentIndex + 1, sections)
      } else if (scrollingUp) {
        scrollToSection(currentIndex - 1, sections)
      }
    }

    const handleTouchStart = (e: TouchEvent) => {
      touchStartY.current = e.touches[0].clientY
      touchStartX.current = e.touches[0].clientX
      touchCurrentY.current = e.touches[0].clientY
      touchCurrentX.current = e.touches[0].clientX
      isVerticalSwipe.current = null
    }

    const handleTouchMove = (e: TouchEvent) => {
      touchCurrentY.current = e.touches[0].clientY
      touchCurrentX.current = e.touches[0].clientX

      const diffY = Math.abs(touchCurrentY.current - touchStartY.current)
      const diffX = Math.abs(touchCurrentX.current - touchStartX.current)

      if (isVerticalSwipe.current === null && (diffX > 10 || diffY > 10)) {
        isVerticalSwipe.current = diffY > diffX
      }

      if (isVerticalSwipe.current === true && !isScrolling.current) {
        const atFirstSection = currentIndex === 0
        const atLastSection = currentIndex === sections.length - 1
        const swipingUp = touchCurrentY.current < touchStartY.current
        const swipingDown = touchCurrentY.current > touchStartY.current

        if ((atFirstSection && swipingDown) || (atLastSection && swipingUp)) {
          return
        }

        e.preventDefault()
      }
    }

    const handleTouchEnd = () => {
      if (isScrolling.current) return
      if (isVerticalSwipe.current !== true) return

      const diffY = touchStartY.current - touchCurrentY.current

      if (Math.abs(diffY) > swipeThreshold) {
        const atFirstSection = currentIndex === 0
        const atLastSection = currentIndex === sections.length - 1

        if (diffY > 0 && !atLastSection) {
          scrollToSection(currentIndex + 1, sections)
        } else if (diffY < 0 && !atFirstSection) {
          scrollToSection(currentIndex - 1, sections)
        }
      }

      isVerticalSwipe.current = null
    }

    const handleKeyDown = (e: KeyboardEvent) => {
      if (isScrolling.current) return

      if (e.key === 'ArrowDown' || e.key === 'PageDown') {
        e.preventDefault()
        scrollToSection(currentIndex + 1, sections)
      } else if (e.key === 'ArrowUp' || e.key === 'PageUp') {
        e.preventDefault()
        scrollToSection(currentIndex - 1, sections)
      } else if (e.key === 'Home') {
        e.preventDefault()
        scrollToSection(0, sections)
      } else if (e.key === 'End') {
        e.preventDefault()
        scrollToSection(sections.length - 1, sections)
      }
    }

    container.addEventListener('wheel', handleWheel, { passive: false })
    container.addEventListener('touchstart', handleTouchStart, { passive: true })
    container.addEventListener('touchmove', handleTouchMove, { passive: false })
    container.addEventListener('touchend', handleTouchEnd, { passive: true })
    window.addEventListener('keydown', handleKeyDown)

    return () => {
      container.removeEventListener('wheel', handleWheel)
      container.removeEventListener('touchstart', handleTouchStart)
      container.removeEventListener('touchmove', handleTouchMove)
      container.removeEventListener('touchend', handleTouchEnd)
      window.removeEventListener('keydown', handleKeyDown)
    }
  }, [containerRef, currentIndex, sectionSelector, swipeThreshold, scrollToSection])

  return { currentIndex, setCurrentIndex, totalSections: sectionsCount.current }
}

export default useSnapScroll
