import { useState, useEffect } from 'react'
import { Link, useLocation, useNavigate } from 'react-router-dom'
import { motion, AnimatePresence, useAnimation } from 'framer-motion'
import { Menu, X } from 'lucide-react'
import CodeCounter from '../ui/CodeCounter'
import { useIntroComplete } from '../../context/IntroContext'

const navLinks = [
  { name: 'About', href: '/about' },
  { name: 'Contact', href: '/contact' },
]

export default function PortfolioNavbar() {
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false)
  const location = useLocation()
  const navigate = useNavigate()
  const introComplete = useIntroComplete()
  const controls = useAnimation()

  useEffect(() => {
    setIsMobileMenuOpen(false)
  }, [location])

  useEffect(() => {
    if (introComplete) {
      controls.start({ y: 0 })
    }
  }, [introComplete, controls])

  const handleNavClick = (href: string) => {
    navigate(href)
    setIsMobileMenuOpen(false)
  }

  return (
    <motion.nav
      initial={{ y: -100 }}
      animate={controls}
      transition={{ duration: 0.6, ease: 'easeOut' }}
      className="fixed top-0 left-0 right-0 z-50 py-3 sm:py-4 lg:py-5 bg-background/70 backdrop-blur-xl border-b border-border/50"
    >
      <div className="container flex items-center justify-between">
        <div className="flex items-center gap-2 sm:gap-4">
          <Link to="/" className="relative group">
            <span className="text-xl sm:text-2xl font-bold gradient-text">Portfolio</span>
            <span className="absolute -bottom-1 left-0 w-0 h-0.5 bg-gradient-to-r from-primary to-accent transition-all duration-300 group-hover:w-full" />
          </Link>
          <CodeCounter />
        </div>

        <button
          onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
          className="lg:hidden p-3 min-w-[44px] min-h-[44px] flex items-center justify-center text-text-primary rounded-full hover:bg-surface-hover active:bg-surface-hover transition-colors"
          aria-label={isMobileMenuOpen ? "Close menu" : "Open menu"}
        >
          {isMobileMenuOpen ? <X size={22} /> : <Menu size={22} />}
        </button>

        <div className="hidden lg:flex items-center gap-10">
          <Link
            to="/blog"
            className="relative text-sm font-medium text-text-secondary hover:text-text-primary transition-colors duration-300 group tracking-wide"
          >
            Blog
            <span className="absolute -bottom-1.5 left-0 w-0 h-[2px] bg-gradient-to-r from-primary to-accent transition-all duration-300 group-hover:w-full" />
          </Link>
          {navLinks.map((link) => (
            <button
              key={link.name}
              onClick={() => handleNavClick(link.href)}
              className="relative text-sm font-medium text-text-secondary hover:text-text-primary transition-colors duration-300 group tracking-wide"
            >
              {link.name}
              <span className="absolute -bottom-1.5 left-0 w-0 h-[2px] bg-gradient-to-r from-primary to-accent transition-all duration-300 group-hover:w-full" />
            </button>
          ))}
          <button className="btn-primary" onClick={() => handleNavClick('/contact')}>
            <span>Hire Me</span>
          </button>
        </div>
      </div>

      <AnimatePresence>
        {isMobileMenuOpen && (
          <motion.div
            initial={{ opacity: 0, height: 0 }}
            animate={{ opacity: 1, height: 'auto' }}
            exit={{ opacity: 0, height: 0 }}
            transition={{ duration: 0.3 }}
            className="lg:hidden glass mt-4 mx-4 rounded-2xl overflow-hidden"
          >
            <div className="flex flex-col p-5 gap-1">
              <Link
                to="/blog"
                onClick={() => setIsMobileMenuOpen(false)}
                className="text-left text-base text-text-secondary hover:text-text-primary active:text-text-primary transition-colors duration-300 py-3 min-h-[44px] flex items-center"
              >
                Blog
              </Link>
              {navLinks.map((link) => (
                <button
                  key={link.name}
                  onClick={() => handleNavClick(link.href)}
                  className="text-left text-base text-text-secondary hover:text-text-primary active:text-text-primary transition-colors duration-300 py-3 min-h-[44px] flex items-center"
                >
                  {link.name}
                </button>
              ))}
              <button className="btn-primary mt-3 w-full" onClick={() => handleNavClick('/contact')}>
                <span>Hire Me</span>
              </button>
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    </motion.nav>
  )
}
