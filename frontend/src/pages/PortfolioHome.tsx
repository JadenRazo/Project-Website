import { Hero, About, Services, Contact } from '../components/sections/portfolio'
import HorizontalProjectGallery from '../components/sections/portfolio/HorizontalProjectGallery'
import SEO from '../components/common/SEO'

export default function PortfolioHome() {
  return (
    <>
      <SEO
        title="Jaden Razo | AWS Cloud & DevOps Engineer"
        description="AWS cloud and DevOps engineer building reliable, secure, cost-aware systems with Terraform, Go, TypeScript, SRE practices, and inspectable operational evidence."
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
