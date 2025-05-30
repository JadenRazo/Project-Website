export class DevCacheManager {
  private static instance: DevCacheManager;
  private isDevEnvironment: boolean;

  private constructor() {
    this.isDevEnvironment = process.env.NODE_ENV === 'development';
  }

  static getInstance(): DevCacheManager {
    if (!DevCacheManager.instance) {
      DevCacheManager.instance = new DevCacheManager();
    }
    return DevCacheManager.instance;
  }

  async clearAllCaches(): Promise<void> {
    if (!this.isDevEnvironment) return;

    await this.unregisterServiceWorkers();
    await this.clearBrowserCaches();
    this.clearLocalStorage();
    this.clearSessionStorage();
  }

  private async unregisterServiceWorkers(): Promise<void> {
    if ('serviceWorker' in navigator) {
      try {
        const registrations = await navigator.serviceWorker.getRegistrations();
        await Promise.all(
          registrations.map(registration => registration.unregister())
        );
        console.log('[DevCacheManager] All service workers unregistered');
      } catch (error) {
        console.error('[DevCacheManager] Error unregistering service workers:', error);
      }
    }
  }

  private async clearBrowserCaches(): Promise<void> {
    if ('caches' in window) {
      try {
        const cacheNames = await caches.keys();
        await Promise.all(
          cacheNames.map(cacheName => caches.delete(cacheName))
        );
        console.log('[DevCacheManager] All caches cleared');
      } catch (error) {
        console.error('[DevCacheManager] Error clearing caches:', error);
      }
    }
  }

  private clearLocalStorage(): void {
    try {
      localStorage.clear();
      console.log('[DevCacheManager] Local storage cleared');
    } catch (error) {
      console.error('[DevCacheManager] Error clearing local storage:', error);
    }
  }

  private clearSessionStorage(): void {
    try {
      sessionStorage.clear();
      console.log('[DevCacheManager] Session storage cleared');
    } catch (error) {
      console.error('[DevCacheManager] Error clearing session storage:', error);
    }
  }

  setupDevTools(): void {
    if (!this.isDevEnvironment) return;

    if (typeof window !== 'undefined') {
      (window as any).__clearAllCaches = () => this.clearAllCaches();
      console.log('[DevCacheManager] Dev tools setup complete. Use window.__clearAllCaches() to clear all caches.');
    }

    this.detectAndWarnAboutServiceWorkers();
  }

  private async detectAndWarnAboutServiceWorkers(): Promise<void> {
    if ('serviceWorker' in navigator) {
      const registrations = await navigator.serviceWorker.getRegistrations();
      if (registrations.length > 0) {
        console.warn(
          '[DevCacheManager] Service workers detected in development!',
          'This may cause caching issues. Running automatic cleanup...'
        );
        await this.unregisterServiceWorkers();
      }
    }
  }
}

export const devCacheManager = DevCacheManager.getInstance();