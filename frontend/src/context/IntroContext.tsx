import { createContext, useContext } from 'react'

export const IntroContext = createContext(false)

export function useIntroComplete() {
  return useContext(IntroContext)
}
