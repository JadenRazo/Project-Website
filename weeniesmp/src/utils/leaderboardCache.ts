import type { LeaderboardEntry, LeaderboardCategory } from '@/stores/leaderboardStore'

export interface CachedLeaderboardData {
  entries: LeaderboardEntry[]
  timestamp: number
  lastUpdated: number
  category: LeaderboardCategory
}

export interface LeaderboardCache {
  [key: string]: CachedLeaderboardData
}

// Cache duration: 5 minutes (matching server-side cache)
const CACHE_DURATION_MS = 5 * 60 * 1000

// Stale threshold: 30 seconds (data is considered "updating soon")
const STALE_THRESHOLD_MS = 30 * 1000

const STORAGE_KEY = 'weeniesmp_leaderboards_cache'

/**
 * Get cached leaderboard data from localStorage
 */
export function getCachedLeaderboard(category: LeaderboardCategory): CachedLeaderboardData | null {
  try {
    const cacheStr = localStorage.getItem(STORAGE_KEY)
    if (!cacheStr) return null

    const cache: LeaderboardCache = JSON.parse(cacheStr)
    const cached = cache[category]

    if (!cached) return null

    // Check if cache is expired
    const now = Date.now()
    if (now - cached.timestamp > CACHE_DURATION_MS) {
      return null
    }

    return cached
  } catch (error) {
    console.error('Error reading leaderboard cache:', error)
    return null
  }
}

/**
 * Save leaderboard data to localStorage cache
 */
export function setCachedLeaderboard(
  category: LeaderboardCategory,
  entries: LeaderboardEntry[],
  lastUpdated: number
): void {
  try {
    const cacheStr = localStorage.getItem(STORAGE_KEY)
    const cache: LeaderboardCache = cacheStr ? JSON.parse(cacheStr) : {}

    cache[category] = {
      entries,
      timestamp: Date.now(),
      lastUpdated,
      category
    }

    localStorage.setItem(STORAGE_KEY, JSON.stringify(cache))
  } catch (error) {
    console.error('Error saving leaderboard cache:', error)
    // If localStorage is full, clear old entries
    if (error instanceof DOMException && error.name === 'QuotaExceededError') {
      try {
        localStorage.removeItem(STORAGE_KEY)
        // Retry saving just this category
        const cache: LeaderboardCache = {
          [category]: {
            entries,
            timestamp: Date.now(),
            lastUpdated,
            category
          }
        }
        localStorage.setItem(STORAGE_KEY, JSON.stringify(cache))
      } catch (retryError) {
        console.error('Failed to save cache after clearing:', retryError)
      }
    }
  }
}

/**
 * Check if cached data is stale (older than threshold but not expired)
 */
export function isCacheStale(category: LeaderboardCategory): boolean {
  const cached = getCachedLeaderboard(category)
  if (!cached) return false

  const now = Date.now()
  const age = now - cached.timestamp

  return age > STALE_THRESHOLD_MS && age < CACHE_DURATION_MS
}

/**
 * Check if cache needs refresh (should fetch new data)
 */
export function shouldRefreshCache(category: LeaderboardCategory): boolean {
  const cached = getCachedLeaderboard(category)
  if (!cached) return true

  const now = Date.now()
  return now - cached.timestamp > CACHE_DURATION_MS
}

/**
 * Clear all cached leaderboard data
 */
export function clearLeaderboardCache(): void {
  try {
    localStorage.removeItem(STORAGE_KEY)
  } catch (error) {
    console.error('Error clearing leaderboard cache:', error)
  }
}

/**
 * Clear cache for a specific category
 */
export function clearCategoryCache(category: LeaderboardCategory): void {
  try {
    const cacheStr = localStorage.getItem(STORAGE_KEY)
    if (!cacheStr) return

    const cache: LeaderboardCache = JSON.parse(cacheStr)
    delete cache[category]

    if (Object.keys(cache).length === 0) {
      localStorage.removeItem(STORAGE_KEY)
    } else {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(cache))
    }
  } catch (error) {
    console.error('Error clearing category cache:', error)
  }
}

/**
 * Get cache age in seconds
 */
export function getCacheAge(category: LeaderboardCategory): number | null {
  const cached = getCachedLeaderboard(category)
  if (!cached) return null

  return Math.floor((Date.now() - cached.timestamp) / 1000)
}

/**
 * Preload cache for all categories (useful for background prefetching)
 */
export async function preloadAllCategories(
  categories: LeaderboardCategory[],
  apiBase: string
): Promise<void> {
  const fetchPromises = categories.map(async (category) => {
    // Only fetch if cache is missing or expired
    if (shouldRefreshCache(category)) {
      try {
        const res = await fetch(`${apiBase}/leaderboards?type=${category}`)
        if (res.ok) {
          const data = await res.json()
          setCachedLeaderboard(category, data.entries || [], data.lastUpdated || Date.now())
        }
      } catch (error) {
        console.error(`Failed to preload ${category}:`, error)
      }
    }
  })

  await Promise.all(fetchPromises)
}
