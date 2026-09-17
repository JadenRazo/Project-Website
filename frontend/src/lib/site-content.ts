import { useSyncExternalStore } from 'react';
import initial from '../../../raizhost/content.json';
import map from '../../../raizhost/content-map.json';

type SiteContent = typeof initial;
type Entry = Exclude<keyof SiteContent, 'version' | 'updatedAt'>;
let current = initial;
const listeners = new Set<() => void>();

function onContentUpdate(event: Event) {
  const detail = (event as CustomEvent<{ path?: unknown; value?: unknown }>).detail;
  if (!detail || typeof detail.path !== 'string') return;
  const [entry, key, extra] = detail.path.split('.');
  if (extra) return;
  const contract = map.entries.find(item => item.id === entry);
  const field = contract?.fields.find(item => item.key === key) as { type: string; min?: number; max?: number; step?: number; options?: { value: string }[] } | undefined;
  if (!field || !Object.hasOwn(current, entry)) return;
  const previous = current[entry as Entry] as Record<string, unknown>;
  let value = detail.value;
  if (field.type === 'number') {
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

function subscribe(listener: () => void) {
  if (listeners.size === 0) window.addEventListener('raizhost:content-update', onContentUpdate);
  listeners.add(listener);
  return () => {
    listeners.delete(listener);
    if (listeners.size === 0) window.removeEventListener('raizhost:content-update', onContentUpdate);
  };
}

/** React owns its nested spans and state; preview messages update values only. */
export function useSiteContent() {
  return useSyncExternalStore(subscribe, () => current, () => initial);
}
