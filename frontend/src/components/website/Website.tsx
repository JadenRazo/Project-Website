import type { ReactNode } from 'react';
import { Helmet } from 'react-helmet-async';
import { useSiteContent } from '../../lib/site-content';
import { WEBSITE_ENTRIES, isWebsiteLink, websiteDesignVariables, type WebsitePage } from '../../../../raizhost/editor/website-kit';
import '../../../../raizhost/editor/website.css';

export function WebsiteDesign() {
  const document = useSiteContent();
  const vars = websiteDesignVariables(document[WEBSITE_ENTRIES.design]);
  const native: Record<string, string> = {};
  if (vars['--rh-accent']) for (const name of ['--color-primary', '--color-primary-light', '--primary-color']) native[name] = vars['--rh-accent'];
  if (vars['--rh-background']) for (const name of ['--color-background', '--background']) native[name] = vars['--rh-background'];
  if (vars['--rh-foreground']) for (const name of ['--color-text-primary', '--text']) native[name] = vars['--rh-foreground'];
  let css = `:root{${Object.entries({ ...vars, ...native }).map(([key, value]) => `${key}:${value}`).join(';')}}`;
  if (vars['--rh-font']) css += 'body,h1,h2,h3,h4,p,a,button,input,textarea{font-family:var(--rh-font)!important}';
  if (vars['--rh-content-width']) css += '[data-rh-native] .container,[data-rh-native] .max-w-7xl,[data-rh-native] .portfolio-container{max-width:var(--rh-content-width)}';
  if (vars['--rh-section-space']) css += '[data-rh-native]>section{padding-block:var(--rh-section-space)}';
  if (vars['--rh-radius']) css += '[data-rh-native] .project-card,[data-rh-native] .btn-primary,[data-rh-native] input,[data-rh-native] textarea{border-radius:var(--rh-radius)}';
  return <style data-rh-website-design>{css}</style>;
}

export function WebsiteNavigation({ belowFixedHeader = false }: { belowFixedHeader?: boolean }) {
  const document = useSiteContent();
  const pages = document[WEBSITE_ENTRIES.pages].filter(page => page.enabled && page.inNavigation);
  if (!pages.length) return null;
  return <nav aria-label="Website pages" data-rh-pages className={`mx-auto max-w-7xl px-5 py-3 ${belowFixedHeader ? 'mt-20' : ''}`}>
    {pages.map(page => <a key={page.id} href={page.path}>{page.title}</a>)}
  </nav>;
}

export function WebsiteSections({ page, native = {} }: { page: WebsitePage; native?: Record<string, ReactNode> }) {
  const document = useSiteContent();
  if (!page.enabled) return <><Helmet><meta name="robots" content="noindex, nofollow" /></Helmet><p className="mx-auto max-w-7xl p-8">This page is unavailable.</p></>;
  return <div data-rh-controlled data-rh-website-root>
    {document[WEBSITE_ENTRIES.sections].map((section, index) => {
      if (section.pageId !== page.id || section.hidden) return null;
      const field = (key: string) => `${WEBSITE_ENTRIES.sections}.${index}.${key}`;
      if (section.component.startsWith('authored:')) {
        const key = section.component.slice(9);
        return <div key={section.id} data-rh-native={key} style={{ display: 'contents' }}>{native[key] ?? null}</div>;
      }
      return <section key={section.id} className="rh-owner-section" data-rh-section={section.id} data-tone={section.tone} style={{ textAlign: section.alignment === 'center' ? 'center' : 'left' }}>
        <div className="rh-owner-section-content">
          {section.title && <h2 data-rh={field('title')}>{section.title}</h2>}
          {section.body && <p data-rh={field('body')}>{section.body}</p>}
          {section.image?.src && <img data-rh-img={field('image')} src={section.image.src} alt={section.image.alt ?? ''} loading="lazy" />}
          {section.link && section.linkLabel && isWebsiteLink(section.link) && <a className="rh-owner-button" href={section.link} data-rh={field('linkLabel')}>{section.linkLabel}</a>}
        </div>
      </section>;
    })}
    <Helmet><title>{page.title}</title>{page.description && <meta name="description" content={page.description} />}</Helmet>
  </div>;
}
