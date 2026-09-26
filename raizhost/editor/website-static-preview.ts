/** Native adapter for static HTML/Astro. React islands use their own state
 * adapter; never clone or rewrite a framework-owned subtree here. */
import { WEBSITE_ENTRIES, isWebsiteLink, websiteDesignVariables, type WebsiteDocument } from "./website-kit";

export function installStaticWebsitePreview(aliases: Record<string, string[]> = {}) {
  const meta = document.querySelector<HTMLMetaElement>('meta[name="raizhost-preview-path"]');
  if (!meta || window.parent === window) return () => {};
  const root = document.querySelector<HTMLElement>("[data-rh-website-root]");
  const parking = document.querySelector<HTMLElement>("[data-rh-website-library]");
  if (!root || !parking) return () => {};
  const originals = new Map([...document.querySelectorAll<HTMLElement>("[data-rh-native]")].map(node => [node.dataset.rhNative!, node]));
  const generated = new Map([...root.querySelectorAll<HTMLElement>('[data-rh-section]')].map(node => [node.dataset.rhSection!, node]));
  let previousDesign = [...Array.from(document.documentElement.style).filter(key => key.startsWith('--rh-')), ...Object.values(aliases).flat()];
  const description = document.querySelector<HTMLMetaElement>('meta[name="description"]');
  const originalDescription = description?.dataset.rhDefault ?? description?.content ?? '';
  const apply = (doc: WebsiteDocument) => {
    if (!Array.isArray(doc[WEBSITE_ENTRIES.pages]) || !Array.isArray(doc[WEBSITE_ENTRIES.sections])) return;
    const path = meta.content.replace(/\/$/, "") || "/";
    const page = doc[WEBSITE_ENTRIES.pages].find(page => page.path === path);
    const style = document.documentElement.style;
    for (const key of previousDesign) style.removeProperty(key);
    const design = websiteDesignVariables(doc[WEBSITE_ENTRIES.design]);
    for (const [key, targets] of Object.entries(aliases)) if (design[key]) for (const target of targets) design[target] = design[key];
    for (const [key, value] of Object.entries(design)) style.setProperty(key, value);
    previousDesign = Object.keys(design);
    const desired: HTMLElement[] = [];
    const used = new Set<string>();
    const rows = doc[WEBSITE_ENTRIES.sections];
    rows.forEach((row, index) => {
      if (!page?.enabled || row.pageId !== page.id || row.hidden || used.has(row.id)) return;
      used.add(row.id);
      if (typeof row.component === "string" && row.component.startsWith("authored:")) {
        const node = originals.get(row.component.slice(9));
        if (node) desired.push(node);
        return;
      }
      if (!["text", "image", "call-to-action"].includes(row.component)) return;
      const section = generated.get(row.id) ?? document.createElement("section");
      generated.set(row.id, section);
      section.className = "rh-owner-section";
      section.dataset.rhSection = row.id;
      section.dataset.tone = ["accent", "muted"].includes(row.tone) ? row.tone : "inherit";
      section.style.textAlign = row.alignment === "center" ? "center" : "left";
      section.replaceChildren();
      const content = document.createElement("div"); content.className = "rh-owner-section-content";
      const text = (tag: string, key: "title" | "body" | "linkLabel", value: unknown) => {
        const node = document.createElement(tag); node.textContent = typeof value === "string" ? value : "";
        node.dataset.rh = `${WEBSITE_ENTRIES.sections}.${index}.${key}`; content.append(node); return node;
      };
      if (row.title) text("h2", "title", row.title);
      if (row.body) text("p", "body", row.body);
      if (row.image?.src && /^(?:\/(?!\/)|data:image\/|https?:\/\/)/.test(row.image.src)) {
        const image = document.createElement("img"); image.src = row.image.src; image.alt = row.image.alt ?? "";
        image.dataset.rhImg = `${WEBSITE_ENTRIES.sections}.${index}.image`; image.loading = "lazy"; content.append(image);
      }
      if (row.link && row.linkLabel && isWebsiteLink(row.link)) {
        const link = text("a", "linkLabel", row.linkLabel) as HTMLAnchorElement; link.href = row.link; link.className = "rh-owner-button";
      }
      section.append(content); desired.push(section);
    });
    // Keep unchanged native elements mounted so form/animation state survives
    // ordinary copy changes. Reordering moves the same elements and listeners.
    for (const node of originals.values()) if (!desired.includes(node) && node.parentElement !== parking) parking.append(node);
    for (const node of generated.values()) if (!desired.includes(node)) node.remove();
    desired.forEach((node, index) => { if (root.children[index] !== node) root.insertBefore(node, root.children[index] ?? null); });
    for (const [id, node] of generated) if (!used.has(id)) { node.remove(); generated.delete(id); }
    for (const nav of document.querySelectorAll<HTMLElement>("[data-rh-pages]")) {
      nav.replaceChildren();
      for (const item of doc[WEBSITE_ENTRIES.pages]) if (item.enabled && item.inNavigation && isWebsiteLink(item.path)) {
        const link = document.createElement("a"); link.href = item.path; link.textContent = item.title; nav.append(link);
      }
    }
    for (const link of document.querySelectorAll<HTMLElement>('[data-rh-page-link]')) {
      const target = doc[WEBSITE_ENTRIES.pages].find(page => page.path === link.dataset.rhPageLink);
      link.hidden = target?.enabled === false;
    }
    if (page?.title) document.title = page.title;
    if (description) description.content = page?.description || originalDescription;
  };
  const update = (event: Event) => apply((event as CustomEvent<WebsiteDocument>).detail);
  window.addEventListener("raizhost:document-update", update);
  window.dispatchEvent(new CustomEvent("raizhost:content-ready"));
  return () => window.removeEventListener("raizhost:document-update", update);
}
