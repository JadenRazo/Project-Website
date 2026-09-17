import { BrowserRouter, MemoryRouter } from 'react-router-dom';
import type { ReactNode } from 'react';

/** The app supplies only a relative site path. It never changes browser origin
 * or gives this website access to the authenticated parent document. */
export function editorPreviewPath(): string | null {
  if (typeof window === 'undefined' || window.parent === window) return null;
  const path = document.querySelector('meta[name="raizhost-preview-path"]')?.getAttribute('content');
  return path && path.startsWith('/') && !path.startsWith('//') && !/[\\\r\n]/.test(path) ? path : null;
}

export function SiteRouter({ children }: { children: ReactNode }) {
  const path = editorPreviewPath();
  return path === null
    ? <BrowserRouter>{children}</BrowserRouter>
    : <MemoryRouter initialEntries={[path]}>{children}</MemoryRouter>;
}
