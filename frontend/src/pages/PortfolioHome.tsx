import { Hero, About, Services, Contact } from '../components/sections/portfolio'
import HorizontalProjectGallery from '../components/sections/portfolio/HorizontalProjectGallery'
import SEO from '../components/common/SEO'

export default function PortfolioHome() {
  return (
    <>
      <SEO
        title="Jaden Razo | Full Stack Developer - Building scalable web applications"
        description="Full Stack Developer specializing in React, TypeScript, Go, and cloud technologies. Building scalable web applications and microservices with modern best practices."
        path="/"
      />
      <div className="relative">
        <Hero />
        <About />
        <HorizontalProjectGallery />
        <Services />
        <Contact />
      </div>
    </>
  )
}
