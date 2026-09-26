import {
  Hero,
  About,
  Services,
  Contact,
} from "../components/sections/portfolio";
import HorizontalProjectGallery from "../components/sections/portfolio/HorizontalProjectGallery";
import SEO from "../components/common/SEO";
import { useSiteContent } from "../lib/site-content";
import { WebsiteSections } from "../components/website/Website";

export default function PortfolioHome() {
  const content = useSiteContent();
  const page = content['website-pages'].find(page => page.id === 'home')!;
  return (
    <>
      <SEO
        title={page.title}
        description={page.description || "AWS cloud and DevOps engineer building reliable, secure, cost-aware systems with Terraform, Go, TypeScript, SRE practices, and inspectable operational evidence."}
        path="/"
      />
      <div className="relative">
        <WebsiteSections page={page} native={{ hero: <Hero />, about: <About />, projects: <HorizontalProjectGallery />, services: <Services />, contact: <Contact /> }} />
      </div>
    </>
  );
}
