import type { ReactNode } from 'react'

import PortfolioNavbar from './PortfolioNavbar'

interface PortfolioLayoutProps {
  children: ReactNode
}

export default function PortfolioLayout({ children }: PortfolioLayoutProps) {
  return (
    <div className="portfolio-shell relative min-h-screen overflow-x-hidden">
      <a
        href="#main-content"
        className="skip fixed left-4 top-4 z-[100] -translate-y-24 rounded-md bg-primary px-4 py-3 font-semibold text-white transition-transform focus:translate-y-0"
      >
        Skip to content
      </a>
      <PortfolioNavbar />
      <main id="main-content" className="relative z-10">
        {children}
      </main>
    </div>
  )
}
