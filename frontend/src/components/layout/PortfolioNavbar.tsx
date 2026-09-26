import { useSiteContent } from '../../lib/site-content';
import SiteText from '../website/SiteText';
import { useState } from 'react'
import { Github, Menu, X } from 'lucide-react'
import { Link } from 'react-router-dom'

import { useScroll } from '../../providers/ScrollProvider'



export default function PortfolioNavbar() {
  const content = useSiteContent();
  const sectionLinks = content['navigation-links'];
  const blogEnabled = content['website-pages'].find(page => page.path === '/blog')?.enabled !== false;
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false)
  const { scrollTo } = useScroll()

  const goToSection = (event: React.MouseEvent<HTMLAnchorElement>, href: string) => {
    if (!href.startsWith('#')) return
    event.preventDefault()
    setIsMobileMenuOpen(false)
    scrollTo(href)
  }

  return (
    <nav className="fixed inset-x-0 top-0 z-50 border-b border-border bg-background/95 py-3 backdrop-blur-md sm:py-4" aria-label="Primary navigation">
      <div className="portfolio-container flex items-center justify-between gap-4">
        <Link to="/" className="flex min-h-11 items-center gap-3 text-text-primary">
          <span className="font-display text-lg font-bold tracking-[-0.02em] sm:text-xl"><SiteText entry="copy-portfolio-navbar-1" name="textf2888132b55cd3b7" /></span>
          <span className="hidden border-l border-border pl-3 font-mono text-[10px] uppercase tracking-[0.16em] text-text-muted sm:block"><SiteText entry="copy-portfolio-navbar-1" name="text4d1eaa3fda45be1c" /></span>
        </Link>

        <button
          type="button"
          onClick={() => setIsMobileMenuOpen((open) => !open)}
          className="flex min-h-11 min-w-11 items-center justify-center border border-border text-text-primary hover:border-primary hover:text-primary lg:hidden"
          aria-label={isMobileMenuOpen ? 'Close menu' : 'Open menu'}
          aria-expanded={isMobileMenuOpen}
          aria-controls="portfolio-mobile-menu"
        >
          {isMobileMenuOpen ? <X size={21} /> : <Menu size={21} />}
        </button>

        <div className="hidden items-center gap-8 lg:flex">
          {sectionLinks.map((link) => (
            <a
              key={link.id}
              href={link.href}
              onClick={(event) => goToSection(event, link.href)}
              className="inline-flex min-h-11 items-center text-sm font-medium text-text-secondary hover:text-primary"
            >
              {link.name}
            </a>
          ))}
          {blogEnabled && <Link to="/blog" className="inline-flex min-h-11 items-center text-sm font-medium text-text-secondary hover:text-primary"><SiteText entry="copy-portfolio-navbar-1" name="text8c6bc099534a0251" /></Link>}
          <a
            href="https://github.com/JadenRazo"
            target="_blank"
            rel="noopener noreferrer"
            className="inline-flex min-h-11 items-center gap-2 text-sm font-medium text-text-secondary hover:text-primary"
          >
            <Github size={16} aria-hidden="true" /><SiteText entry="copy-portfolio-navbar-1" name="textf911e414cf6bdfc5" /></a>
          <a href="mailto:contact@jadenrazo.dev" className="inline-flex min-h-11 items-center border border-primary px-4 text-sm font-semibold text-primary hover:bg-primary hover:text-white"><SiteText entry="copy-portfolio-navbar-1" name="text2b5c3d26721ae9c3" /></a>
        </div>
      </div>

      {isMobileMenuOpen && (
        <div id="portfolio-mobile-menu" className="portfolio-container pt-3 lg:hidden">
          <div className="border border-border bg-background-secondary p-3">
            {sectionLinks.map((link) => (
              <a
                key={link.id}
                href={link.href}
                onClick={(event) => goToSection(event, link.href)}
                className="flex min-h-12 items-center border-b border-border px-2 text-base text-text-secondary last:border-b-0 hover:text-primary"
              >
                {link.name}
              </a>
            ))}
            {blogEnabled && <Link to="/blog" onClick={() => setIsMobileMenuOpen(false)} className="flex min-h-12 items-center border-b border-border px-2 text-base text-text-secondary hover:text-primary"><SiteText entry="copy-portfolio-navbar-1" name="text8c6bc099534a0251" /></Link>}
            <a
              href="https://github.com/JadenRazo"
              target="_blank"
              rel="noopener noreferrer"
              className="flex min-h-12 items-center gap-2 border-b border-border px-2 text-base text-text-secondary hover:text-primary"
            >
              <Github size={17} aria-hidden="true" /><SiteText entry="copy-portfolio-navbar-1" name="textf911e414cf6bdfc5" /></a>
            <a href="mailto:contact@jadenrazo.dev" className="mt-3 flex min-h-12 items-center justify-center bg-primary px-4 font-semibold text-white"><SiteText entry="copy-portfolio-navbar-1" name="text9d4575ccb2a697bc" /></a>
          </div>
        </div>
      )}
    </nav>
  )
}
