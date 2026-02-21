import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface PlayerStats {
  kills: number
  deaths: number
  blocksPlaced: number
  blocksBroken: number
  mobsKilled: number
  distanceTraveled: number
}

export interface Achievement {
  id: string
  name: string
  description: string
  icon: string
  unlockedAt: string | null
}

export interface PlayerProfile {
  uuid: string
  username: string
  firstJoin: string
  lastSeen: string
  playtime: number
  stats: PlayerStats
  achievements: Achievement[]
  isOnline: boolean
}

export interface PlayerSearchResult {
  uuid: string
  username: string
  isOnline: boolean
  lastSeen: string
}

// Mock data for development - replace with actual API calls
const MOCK_PROFILES: Record<string, PlayerProfile> = {
  'notch': {
    uuid: '069a79f4-44e9-4726-a5be-fca90e38aaf5',
    username: 'Notch',
    firstJoin: '2024-01-15T12:00:00Z',
    lastSeen: '2025-01-17T08:30:00Z',
    playtime: 12450,
    isOnline: false,
    stats: {
      kills: 342,
      deaths: 156,
      blocksPlaced: 48923,
      blocksBroken: 52341,
      mobsKilled: 8934,
      distanceTraveled: 1234567
    },
    achievements: [
      { id: 'first_join', name: 'Welcome!', description: 'Join the server for the first time', icon: 'star', unlockedAt: '2024-01-15T12:00:00Z' },
      { id: 'first_kill', name: 'First Blood', description: 'Get your first player kill', icon: 'sword', unlockedAt: '2024-01-16T14:30:00Z' },
      { id: 'mine_diamond', name: 'Diamonds!', description: 'Mine your first diamond', icon: 'gem', unlockedAt: '2024-01-15T18:45:00Z' },
      { id: 'playtime_100', name: 'Dedicated', description: 'Play for 100 hours', icon: 'clock', unlockedAt: '2024-03-20T10:00:00Z' },
      { id: 'build_master', name: 'Build Master', description: 'Place 50,000 blocks', icon: 'hammer', unlockedAt: null },
    ]
  },
  'jeb_': {
    uuid: '853c80ef-3c37-49fd-aa49-938b674adae6',
    username: 'jeb_',
    firstJoin: '2024-02-20T09:00:00Z',
    lastSeen: '2025-01-17T10:15:00Z',
    playtime: 8760,
    isOnline: true,
    stats: {
      kills: 89,
      deaths: 234,
      blocksPlaced: 123456,
      blocksBroken: 98765,
      mobsKilled: 12345,
      distanceTraveled: 2345678
    },
    achievements: [
      { id: 'first_join', name: 'Welcome!', description: 'Join the server for the first time', icon: 'star', unlockedAt: '2024-02-20T09:00:00Z' },
      { id: 'mine_diamond', name: 'Diamonds!', description: 'Mine your first diamond', icon: 'gem', unlockedAt: '2024-02-20T15:30:00Z' },
      { id: 'builder', name: 'Builder', description: 'Place 10,000 blocks', icon: 'blocks', unlockedAt: '2024-04-10T12:00:00Z' },
    ]
  }
}

const MOCK_PLAYERS: PlayerSearchResult[] = [
  { uuid: '069a79f4-44e9-4726-a5be-fca90e38aaf5', username: 'Notch', isOnline: false, lastSeen: '2025-01-17T08:30:00Z' },
  { uuid: '853c80ef-3c37-49fd-aa49-938b674adae6', username: 'jeb_', isOnline: true, lastSeen: '2025-01-17T10:15:00Z' },
  { uuid: 'f498513c-e8c8-4773-be26-ecfc7ed5185d', username: 'Dinnerbone', isOnline: false, lastSeen: '2025-01-16T22:00:00Z' },
  { uuid: '7125ba8b-1c86-4508-b92b-b5c042ccfe2b', username: 'Grumm', isOnline: true, lastSeen: '2025-01-17T09:45:00Z' },
]

export const usePlayerStore = defineStore('player', () => {
  const profile = ref<PlayerProfile | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)
  const searchResults = ref<PlayerSearchResult[]>([])
  const searchLoading = ref(false)
  const recentPlayers = ref<PlayerSearchResult[]>([])
  const popularPlayers = ref<PlayerSearchResult[]>([])

  async function fetchProfile(username: string) {
    loading.value = true
    error.value = null
    profile.value = null

    try {
      // Simulate API delay
      await new Promise(resolve => setTimeout(resolve, 500))

      // Mock implementation - replace with actual API call
      const normalizedUsername = username.toLowerCase()
      const mockProfile = MOCK_PROFILES[normalizedUsername]

      if (mockProfile) {
        profile.value = mockProfile
      } else {
        // Try to create a placeholder profile using Mojang API for UUID
        try {
          const mojangRes = await fetch(`https://api.mojang.com/users/profiles/minecraft/${username}`)
          if (mojangRes.ok) {
            const mojangData = await mojangRes.json()
            // Player exists in Mojang but not on our server
            error.value = `Player "${mojangData.name}" has not joined WeenieSMP yet`
          } else {
            error.value = `Player "${username}" not found`
          }
        } catch {
          error.value = `Player "${username}" not found`
        }
      }

      // TODO: Replace with actual API call
      // const response = await fetch(`/api/players/${username}`)
      // if (!response.ok) throw new Error('Player not found')
      // profile.value = await response.json()
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to fetch player profile'
    } finally {
      loading.value = false
    }
  }

  async function searchPlayers(query: string) {
    if (!query.trim()) {
      searchResults.value = []
      return
    }

    searchLoading.value = true

    try {
      // Simulate API delay
      await new Promise(resolve => setTimeout(resolve, 300))

      // Mock implementation - filter mock players
      const normalizedQuery = query.toLowerCase()
      searchResults.value = MOCK_PLAYERS.filter(p =>
        p.username.toLowerCase().includes(normalizedQuery)
      )

      // TODO: Replace with actual API call
      // const response = await fetch(`/api/players/search?q=${encodeURIComponent(query)}`)
      // searchResults.value = await response.json()
    } catch (e) {
      console.error('Search error:', e)
      searchResults.value = []
    } finally {
      searchLoading.value = false
    }
  }

  async function fetchRecentPlayers() {
    try {
      // Mock implementation
      recentPlayers.value = MOCK_PLAYERS.slice(0, 4)

      // TODO: Replace with actual API call
      // const response = await fetch('/api/players/recent')
      // recentPlayers.value = await response.json()
    } catch (e) {
      console.error('Failed to fetch recent players:', e)
    }
  }

  async function fetchPopularPlayers() {
    try {
      // Mock implementation
      popularPlayers.value = MOCK_PLAYERS.slice(0, 4)

      // TODO: Replace with actual API call
      // const response = await fetch('/api/players/popular')
      // popularPlayers.value = await response.json()
    } catch (e) {
      console.error('Failed to fetch popular players:', e)
    }
  }

  function clearProfile() {
    profile.value = null
    error.value = null
  }

  return {
    profile,
    loading,
    error,
    searchResults,
    searchLoading,
    recentPlayers,
    popularPlayers,
    fetchProfile,
    searchPlayers,
    fetchRecentPlayers,
    fetchPopularPlayers,
    clearProfile
  }
})
