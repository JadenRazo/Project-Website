import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import {
  getCachedLeaderboard,
  setCachedLeaderboard,
  shouldRefreshCache,
  isCacheStale,
  getCacheAge
} from '@/utils/leaderboardCache'

// Use the nginx proxy path for API calls to avoid CORS issues
const API_BASE = '/api/mc'

export interface LeaderboardEntry {
  rank: number
  uuid: string
  username: string
  value: number
  metadata?: {
    formatted?: string
    job?: string
    [key: string]: any
  }
}

export type LeaderboardCategory = 'playtime' | 'economy' | 'claims' | 'chunks' | 'jobs' | 'voting'

export const useLeaderboardStore = defineStore('leaderboard', () => {
  const entries = ref<LeaderboardEntry[]>([])
  const loading = ref(false)
  const initialLoading = ref(true)
  const updating = ref(false)
  const error = ref<string | null>(null)
  const currentCategory = ref<LeaderboardCategory>('playtime')
  const lastUpdated = ref<number>(0)

  // Computed properties
  const isStale = computed(() => isCacheStale(currentCategory.value))
  const cacheAge = computed(() => getCacheAge(currentCategory.value))
  const hasCache = computed(() => entries.value.length > 0)

  /**
   * Fetch leaderboard with cache-first strategy
   * 1. Load from cache instantly if available
   * 2. Fetch fresh data in background
   * 3. Update UI smoothly when new data arrives
   */
  async function fetchLeaderboard(category: LeaderboardCategory, forceRefresh = false) {
    currentCategory.value = category
    error.value = null

    // Step 1: Load from cache instantly
    const cached = getCachedLeaderboard(category)
    if (cached && !forceRefresh) {
      entries.value = cached.entries
      lastUpdated.value = cached.lastUpdated
      initialLoading.value = false
      loading.value = false
    } else {
      // No cache available, show loading state
      loading.value = true
      initialLoading.value = true
    }

    // Step 2: Check if we need to fetch fresh data
    if (!forceRefresh && !shouldRefreshCache(category)) {
      // Cache is still fresh, no need to fetch
      return
    }

    // Step 3: Fetch fresh data in background
    updating.value = true

    try {
      const res = await fetch(`${API_BASE}/leaderboards?type=${category}`)
      if (!res.ok) throw new Error('Failed to fetch leaderboard')

      const data = await res.json()
      const newEntries = data.entries || []
      const newLastUpdated = data.lastUpdated || Date.now()

      // Step 4: Update cache
      setCachedLeaderboard(category, newEntries, newLastUpdated)

      // Step 5: Update UI smoothly
      entries.value = newEntries
      lastUpdated.value = newLastUpdated
      error.value = null
    } catch (e) {
      // Only show error if we don't have cached data
      if (!cached) {
        error.value = e instanceof Error ? e.message : 'Unknown error'
        entries.value = []
      }
      console.error('Leaderboard fetch error:', e)
    } finally {
      loading.value = false
      initialLoading.value = false
      updating.value = false
    }
  }

  /**
   * Force refresh (bypass cache)
   */
  async function refresh() {
    await fetchLeaderboard(currentCategory.value, true)
  }

  /**
   * Prefetch data for a category without switching to it
   */
  async function prefetchCategory(category: LeaderboardCategory) {
    // Only fetch if cache is missing or expired
    if (shouldRefreshCache(category)) {
      try {
        const res = await fetch(`${API_BASE}/leaderboards?type=${category}`)
        if (res.ok) {
          const data = await res.json()
          setCachedLeaderboard(category, data.entries || [], data.lastUpdated || Date.now())
        }
      } catch (error) {
        console.error(`Failed to prefetch ${category}:`, error)
      }
    }
  }

  return {
    entries,
    loading,
    initialLoading,
    updating,
    error,
    currentCategory,
    lastUpdated,
    isStale,
    cacheAge,
    hasCache,
    fetchLeaderboard,
    refresh,
    prefetchCategory
  }
})
