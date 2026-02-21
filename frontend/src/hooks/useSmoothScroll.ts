type EasingFunction = (t: number) => number

const easeInOutCubic: EasingFunction = (t) => {
  return t < 0.5 ? 4 * t * t * t : 1 - Math.pow(-2 * t + 2, 3) / 2
}

const isSafari = (): boolean => {
  if (typeof window === 'undefined') return false
  return /^((?!chrome|android).)*safari/i.test(navigator.userAgent)
}

const supportsNativeSmoothScroll = (): boolean => {
  if (typeof window === 'undefined') return false
  return 'scrollBehavior' in document.documentElement.style
}

export const smoothScrollTo = (
  target: HTMLElement | number,
  duration: number = 800,
  offset: number = 0
): Promise<void> => {
  return new Promise((resolve) => {
    const startPosition = window.pageYOffset
    const targetPosition = typeof target === 'number'
      ? target
      : target.getBoundingClientRect().top + startPosition - offset
    const distance = targetPosition - startPosition

    if (distance === 0) {
      resolve()
      return
    }

    if (!isSafari() && supportsNativeSmoothScroll()) {
      window.scrollTo({
        top: targetPosition,
        behavior: 'smooth'
      })
      setTimeout(resolve, duration)
      return
    }

    let startTime: number | null = null

    const animation = (currentTime: number) => {
      if (startTime === null) startTime = currentTime
      const elapsed = currentTime - startTime
      const progress = Math.min(elapsed / duration, 1)
      const easeProgress = easeInOutCubic(progress)

      window.scrollTo(0, startPosition + distance * easeProgress)

      if (elapsed < duration) {
        requestAnimationFrame(animation)
      } else {
        resolve()
      }
    }

    requestAnimationFrame(animation)
  })
}

export const scrollToElement = (
  selector: string,
  duration: number = 800,
  offset: number = 0
): Promise<void> => {
  const element = document.querySelector(selector) as HTMLElement
  if (!element) {
    return Promise.resolve()
  }
  return smoothScrollTo(element, duration, offset)
}

export const scrollToTop = (duration: number = 800): Promise<void> => {
  return smoothScrollTo(0, duration)
}

export const useSmoothScroll = () => {
  return {
    scrollTo: smoothScrollTo,
    scrollToElement,
    scrollToTop,
    isSafari: isSafari(),
    supportsNativeSmooth: supportsNativeSmoothScroll()
  }
}

export default useSmoothScroll
