import { useSyncExternalStore } from 'react';
import initial from '../../../raizhost/content.json';
import map from '../../../raizhost/content-map.json';
import type { WebsiteDocument } from '../../../raizhost/editor/website-kit';

type SiteContent = Omit<typeof initial, keyof WebsiteDocument> & WebsiteDocument;
type Entry = Exclude<keyof SiteContent, 'version' | 'updatedAt'>;
let current: SiteContent = initial as SiteContent;
const listeners = new Set<() => void>();

function onContentUpdate(event: Event) {
  const detail = (event as CustomEvent<{ path?: unknown; value?: unknown }>).detail;
  if (!detail || typeof detail.path !== 'string') return;
  const [entry, key, extra] = detail.path.split('.');
  if (extra) return;
  const contract = map.entries.find(item => item.id === entry);
  if (!key && contract && 'type' in contract && contract.type === 'list' && Array.isArray(detail.value)) {
    current = { ...current, [entry]: detail.value };
    listeners.forEach(listener => listener());
    return;
  }
  const field = contract?.fields.find(item => item.key === key) as { type: string; min?: number; max?: number; step?: number; options?: { value: string }[] } | undefined;
  if (!field || !Object.hasOwn(current, entry)) return;
  const previous = current[entry as Entry] as unknown as Record<string, unknown>;
  let value = detail.value;
  if (field.type === 'image') {
    if (value !== null && (!value || typeof value !== 'object' || !('src' in value) || typeof value.src !== 'string')) return;
  } else if (field.type === 'number') {
    if (typeof value !== 'number' && typeof value !== 'string') return;
    if (typeof value === 'string' && value.trim() === '') return;
    value = Number(value);
    if (!Number.isFinite(value) || (field.min != null && (value as number) < field.min)
      || (field.max != null && (value as number) > field.max)) return;
  } else if (field.type === 'boolean') {
    if (typeof value !== 'boolean') return;
  } else if (field.type === 'select') {
    if (typeof value !== 'string' || !field.options?.some(option => option.value === value)) return;
  } else if (typeof value !== 'string' || (field.max != null && value.length > field.max)) return;
  if (Object.is(previous[key], value)) return;
  current = { ...current, [entry]: { ...previous, [key]: value } };
  listeners.forEach(listener => listener());
}

function onDocumentUpdate(event: Event) {
  const doc = (event as CustomEvent<Record<string, unknown>>).detail;
  if (!doc || doc.version !== 1 || typeof doc !== 'object') return;
  const next = { ...current };
  for (const entry of map.entries) {
    const value = doc[entry.id];
    if (value && typeof value === 'object') Object.assign(next, { [entry.id]: value });
  }
  // Stable section IDs are reconciled by React. Never replace its DOM through
  // the editor's scalar marker bridge.
  current = next;
  listeners.forEach(listener => listener());
}

function subscribe(listener: () => void) {
  if (listeners.size === 0) {
    window.addEventListener('raizhost:content-update', onContentUpdate);
    window.addEventListener('raizhost:document-update', onDocumentUpdate);
    queueMicrotask(() => window.dispatchEvent(new CustomEvent('raizhost:content-ready')));
  }
  listeners.add(listener);
  return () => {
    listeners.delete(listener);
    if (listeners.size === 0) {
      window.removeEventListener('raizhost:content-update', onContentUpdate);
      window.removeEventListener('raizhost:document-update', onDocumentUpdate);
    }
  };
}

/** React owns its nested spans and state; preview messages update values only. */
export function useSiteContent() {
  return useSyncExternalStore(subscribe, () => current, () => initial);
}
