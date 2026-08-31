import { createContext, type ReactNode, useCallback, useContext } from 'react'

interface ScrollContextValue {
  scrollTo: (target: string | number | HTMLElement, options?: { immediate?: boolean }) => void
}

const ScrollContext = createContext<ScrollContextValue>({
  scrollTo: () => {},
})

export const useScroll = () => useContext(ScrollContext)

interface ScrollProviderProps {
  children: ReactNode
}

export default function ScrollProvider({ children }: ScrollProviderProps) {
  const scrollTo = useCallback((target: string | number | HTMLElement, options?: { immediate?: boolean }) => {
    const reduceMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches
    const behavior: ScrollBehavior = options?.immediate || reduceMotion ? 'auto' : 'smooth'

    if (typeof target === 'string') {
      document.querySelector(target)?.scrollIntoView({ behavior })
    } else if (typeof target === 'number') {
      window.scrollTo({ top: target, behavior })
    } else if (target instanceof HTMLElement) {
      target.scrollIntoView({ behavior })
    }
  }, [])

  return (
    <ScrollContext.Provider value={{ scrollTo }}>
      {children}
    </ScrollContext.Provider>
  )
}
