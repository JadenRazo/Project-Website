/** Portable website composition contract. No framework, network, or server dependencies.
 * The site renderer and editor use this same contract; content never contains code. */
export const WEBSITE_ENTRIES = { pages: "website-pages", sections: "website-sections", design: "website-design" } as const;
export const SECTION_TEMPLATES = ["text", "image", "call-to-action"] as const;
export interface AuthoredPage { id: string; path: string; title: string }
export interface AuthoredSection { id: string; pageId: string; label: string }
export interface WebsiteContract {
  version: 1;
  /** Native shell includes the component registry, even for unpublished routes. */
  previewPath: string;
  pages: AuthoredPage[];
  sections: AuthoredSection[];
  /** Paths and path prefixes owned by the application, never owner-created pages. */
  reservedPaths: string[];
}
export interface WebsitePage { id: string; path: string; title: string; description: string; inNavigation: boolean; enabled: boolean }
export interface WebsiteSection {
  id: string; pageId: string; component: string; title: string; body: string;
  image: { src: string; alt?: string } | null; link: string; linkLabel: string;
  alignment: string; tone: string; hidden: boolean;
}
export interface WebsiteDesign {
  accent: string; background: string; foreground: string; font: string;
  contentWidth: number | null; cornerRadius: number | null; spacing: string;
}
export interface WebsiteDocument {
  "website-pages": WebsitePage[];
  "website-sections": WebsiteSection[];
  "website-design": WebsiteDesign;
}
export interface WebsiteProblem { path: string; message: string; code: "format" | "required" }

/** Deliberately canonical routes: no decoding ambiguity, traversal, query or hash. */
export function isWebsitePath(value: unknown): value is string {
  return typeof value === "string" && value.length <= 180
    && /^(?:\/|(?:\/[a-z0-9]+(?:-[a-z0-9]+)*)+)$/.test(value);
}
export function isWebsiteLink(value: unknown): value is string {
  if (typeof value !== "string" || value.length > 500 || /[\s<>"'\\\u0000-\u001f]/.test(value)) return false;
  if (value === "" || value === "#" || /^#[a-zA-Z][\w-]*$/.test(value) || /^mailto:[^@]+@[^@]+\.[^@]+$/.test(value) || /^tel:\+?[0-9()-]{7,20}$/.test(value)) return true;
  try {
    if (/^https?:\/\//i.test(value)) { const url = new URL(value); return Boolean(url.hostname) && !url.username && !url.password; }
    if (!value.startsWith("/") || value.startsWith("//")) return false;
    const decoded = decodeURIComponent(value);
    if (/[\\\u0000-\u0020]/.test(decoded) || decoded.startsWith("//") || decoded.split(/[/?#]/).includes("..")) return false;
    return new URL(value, "https://website.invalid").origin === "https://website.invalid";
  } catch { return false; }
}
export const isWebsiteColor = (value: unknown): value is string => typeof value === "string" && (value === "" || /^#[0-9a-f]{6}$/i.test(value));
const record = (value: unknown): Record<string, unknown> => value && typeof value === "object" && !Array.isArray(value) ? value as Record<string, unknown> : {};

export function websiteDefaults(contract: WebsiteContract): WebsiteDocument {
  return {
    "website-pages": contract.pages.map(page => ({ ...page, description: "", inNavigation: false, enabled: true })),
    "website-sections": contract.sections.map(section => ({
      id: section.id, pageId: section.pageId, component: `authored:${section.id}`,
      title: section.label, body: "", image: null, link: "", linkLabel: "",
      alignment: "left", tone: "inherit", hidden: false,
    })),
    "website-design": { accent: "", background: "", foreground: "", font: "inherit", contentWidth: null, cornerRadius: null, spacing: "inherit" },
  };
}

/** Runs after normal content validation on every save, preview, publish and restore. */
export function validateWebsite(contract: WebsiteContract, document: unknown): WebsiteProblem[] {
  const doc = record(document);
  const pages = Array.isArray(doc[WEBSITE_ENTRIES.pages]) ? doc[WEBSITE_ENTRIES.pages] as Record<string, unknown>[] : [];
  const sections = Array.isArray(doc[WEBSITE_ENTRIES.sections]) ? doc[WEBSITE_ENTRIES.sections] as Record<string, unknown>[] : [];
  const errors: WebsiteProblem[] = [];
  const paths = new Set<string>();
  const ids = new Set<string>();
  pages.forEach((page, index) => {
    const base = `${WEBSITE_ENTRIES.pages}.${index}`;
    const authored = contract.pages.find(candidate => candidate.id === page.id);
    if (!isWebsitePath(page.path)) errors.push({ path: `${base}.path`, message: "Use a web address such as /about-us, with lowercase letters, numbers and hyphens.", code: "format" });
    else {
      if (paths.has(page.path)) errors.push({ path: `${base}.path`, message: "Another page already uses this web address.", code: "format" });
      paths.add(page.path);
      if (authored && authored.path !== page.path) errors.push({ path: `${base}.path`, message: "This existing page keeps its web address so its links and integrations continue to work.", code: "format" });
      if (!authored && (contract.pages.some(candidate => candidate.path === page.path)
        || contract.reservedPaths.some(prefix => page.path === prefix || (page.path as string).startsWith(prefix.replace(/\/$/, "") + "/")))) {
        errors.push({ path: `${base}.path`, message: "This web address belongs to an existing website feature. Choose another address.", code: "format" });
      }
    }
    if (typeof page.id === "string") ids.add(page.id);
  });
  // Authored routes stay represented; owners may hide their content, but cannot
  // remove an app route from the contract or accidentally replace the homepage.
  for (const page of contract.pages) {
    if (!pages.some(candidate => candidate.id === page.id)) errors.push({ path: WEBSITE_ENTRIES.pages, message: `Keep ${page.title} in the page list; turn it off to hide its content.`, code: "required" });
  }
  if (!pages.some(page => page.path === "/" && page.enabled === true)) errors.push({ path: WEBSITE_ENTRIES.pages, message: "Keep the homepage enabled.", code: "required" });
  const native = new Set<string>();
  sections.forEach((section, index) => {
    const base = `${WEBSITE_ENTRIES.sections}.${index}`;
    if (!ids.has(String(section.pageId))) errors.push({ path: `${base}.pageId`, message: "Choose an existing page for this section, or remove the section with its page.", code: "format" });
    if (typeof section.component === "string" && section.component.startsWith("authored:")) {
      const id = section.component.slice(9);
      const authored = contract.sections.find(candidate => candidate.id === id);
      if (!authored) errors.push({ path: `${base}.component`, message: "Choose an available section.", code: "format" });
      else if (authored.pageId !== section.pageId) errors.push({ path: `${base}.pageId`, message: "This custom section belongs to its existing page. New text, photo and button sections can move between pages.", code: "format" });
      if (native.has(id)) errors.push({ path: `${base}.component`, message: "This custom section is already on the website. Move the existing section instead.", code: "format" });
      native.add(id);
    } else if (!SECTION_TEMPLATES.includes(section.component as typeof SECTION_TEMPLATES[number])) {
      errors.push({ path: `${base}.component`, message: "Choose an available section.", code: "format" });
    }
  });
  return errors;
}

/** Safe CSS values only. Empty settings preserve the site's authored defaults. */
export function websiteDesignVariables(value: unknown): Record<string, string> {
  const design = record(value);
  const result: Record<string, string> = {};
  for (const key of ["accent", "background", "foreground"] as const) {
    if (isWebsiteColor(design[key]) && design[key]) result[`--rh-${key}`] = design[key];
  }
  if (design.font === "sans") result["--rh-font"] = 'Inter, ui-sans-serif, system-ui, sans-serif';
  if (design.font === "serif") result["--rh-font"] = 'Georgia, ui-serif, serif';
  if (typeof design.contentWidth === "number" && design.contentWidth >= 640 && design.contentWidth <= 1600) result["--rh-content-width"] = `${design.contentWidth}px`;
  if (typeof design.cornerRadius === "number" && design.cornerRadius >= 0 && design.cornerRadius <= 48) result["--rh-radius"] = `${design.cornerRadius}px`;
  const spacing: Record<string, string> = { compact: "2rem", comfortable: "4rem", spacious: "6rem" };
  if (typeof design.spacing === "string" && spacing[design.spacing]) result["--rh-section-space"] = spacing[design.spacing];
  return result;
}

/** Authored as ordinary content entries so existing drafts, history and CAS
 * remain the only write path. `reference` choices come from the current draft. */
export function websiteEntries(contract: WebsiteContract) {
  const identity = { key: "id", type: "text", label: "Identity", access: "managed", required: true, editorHidden: true, max: 100 };
  const options = (values: string[]) => values.map(value => ({ value, label: value === "inherit" ? "Original design" : value.charAt(0).toUpperCase() + value.slice(1) }));
  return [
    { id: WEBSITE_ENTRIES.pages, label: "Pages & navigation", group: "Website", access: "owner", type: "list", min: 1, max: 100, identityKey: "id", itemLabel: "{title}",
      help: "Add pages, choose their web addresses, and arrange their order in navigation. Existing web addresses stay stable.",
      create: { label: "page", defaults: { title: "New page", path: "/new-page", description: "", inNavigation: true, enabled: true } },
      fields: [identity, { key: "title", type: "text", label: "Page title", required: true, max: 120 },
        { key: "path", type: "route", label: "Web address", required: true },
        { key: "description", type: "textarea", label: "Search description", max: 300 },
        { key: "inNavigation", type: "boolean", label: "Show in navigation" },
        { key: "enabled", type: "boolean", label: "Page enabled" }],
    },
    { id: WEBSITE_ENTRIES.sections, label: "Sections & layout", group: "Website", access: "owner", type: "list", min: 0, max: 200, identityKey: "id", itemLabel: "{title}",
      help: "Add, move, hide or remove sections. Their order here is their order on each page.",
      create: { label: "section", defaults: { pageId: contract.pages[0].id, component: "text", title: "New section", body: "", image: null, link: "", linkLabel: "", alignment: "left", tone: "inherit", hidden: false } },
      fields: [identity, { key: "pageId", type: "reference", collection: WEBSITE_ENTRIES.pages, label: "Page", required: true },
        { key: "component", type: "select", label: "Section type", required: true, options: [...options([...SECTION_TEMPLATES]), ...contract.sections.map(section => ({ value: `authored:${section.id}`, label: section.label }))] },
        { key: "title", type: "text", label: "Heading", max: 180 },
        { key: "body", type: "textarea", label: "Text", max: 10000 },
        { key: "image", type: "image", label: "Photo", alt: true, maxWidth: 2400 },
        { key: "link", type: "link", label: "Button link" }, { key: "linkLabel", type: "text", label: "Button text", max: 100 },
        { key: "alignment", type: "select", label: "Text alignment", options: options(["left", "center"]) },
        { key: "tone", type: "select", label: "Section background", options: options(["inherit", "accent", "muted"]) },
        { key: "hidden", type: "boolean", label: "Hide section" }],
    },
    { id: WEBSITE_ENTRIES.design, label: "Colors, type & spacing", group: "Website", access: "owner", fields: [
      { key: "accent", type: "color", label: "Accent color", help: "Leave blank to keep the original design." },
      { key: "background", type: "color", label: "Background color" }, { key: "foreground", type: "color", label: "Text color" },
      { key: "font", type: "select", label: "Typeface", options: options(["inherit", "sans", "serif"]) },
      { key: "contentWidth", type: "number", label: "Content width", min: 640, max: 1600, step: 20, help: "Pixels on wide screens. Leave blank for the original width." },
      { key: "cornerRadius", type: "number", label: "Corner rounding", min: 0, max: 48, step: 1 },
      { key: "spacing", type: "select", label: "Section spacing", options: options(["inherit", "compact", "comfortable", "spacious"]) },
    ] },
  ];
}
