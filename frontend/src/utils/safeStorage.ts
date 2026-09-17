/** Storage can be unavailable in private browsing or an opaque editor iframe.
 * Keep a tab-local fallback; never request access to the embedding app's data. */
export function safeStorage(kind: 'localStorage' | 'sessionStorage'): Storage {
  const fallback = new Map<string, string>();
  let disabled = false;
  const storage = (): Storage | null => {
    if (disabled) return null;
    try { return typeof window === 'undefined' ? null : window[kind]; }
    catch { disabled = true; return null; }
  };
  return {
    get length() { try { return storage()?.length ?? fallback.size; } catch { disabled = true; return fallback.size; } },
    key(index) {
      try { const target = storage(); return target ? target.key(index) : [...fallback.keys()][index] ?? null; }
      catch { disabled = true; return [...fallback.keys()][index] ?? null; }
    },
    getItem(key) {
      try {
        const target = storage();
        if (!target) return fallback.get(key) ?? null;
        const value = target.getItem(key);
        if (value === null) fallback.delete(key); else fallback.set(key, value);
        return value;
      } catch { disabled = true; return fallback.get(key) ?? null; }
    },
    setItem(key, value) { fallback.set(key, String(value)); try { storage()?.setItem(key, String(value)); } catch { disabled = true; } },
    removeItem(key) { fallback.delete(key); try { storage()?.removeItem(key); } catch { disabled = true; } },
    clear() { fallback.clear(); try { storage()?.clear(); } catch { disabled = true; } },
  };
}

export const local = safeStorage('localStorage');
export const session = safeStorage('sessionStorage');
