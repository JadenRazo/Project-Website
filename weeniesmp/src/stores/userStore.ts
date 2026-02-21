import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useToastStore } from './toastStore'

export interface Purchase {
  id: string
  packageName: string
  price: number
  date: string
  status: 'completed' | 'pending' | 'refunded'
}

export interface User {
  username: string
  uuid: string | null
  email: string | null
  purchases: Purchase[]
  joinDate: string | null
}

const MOJANG_API = 'https://api.mojang.com/users/profiles/minecraft'

export const useUserStore = defineStore('user', () => {
  const toastStore = useToastStore()

  const user = ref<User | null>(null)
  const purchases = ref<Purchase[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  const isAuthenticated = computed(() => user.value !== null)
  const totalSpent = computed(() =>
    purchases.value
      .filter(p => p.status === 'completed')
      .reduce((sum, p) => sum + p.price, 0)
  )
  const purchaseCount = computed(() =>
    purchases.value.filter(p => p.status === 'completed').length
  )

  function isValidMinecraftUsername(name: string): boolean {
    return /^[a-zA-Z0-9_]{3,16}$/.test(name)
  }

  async function login(username: string): Promise<boolean> {
    if (!isValidMinecraftUsername(username)) {
      error.value = 'Invalid username. Must be 3-16 characters (letters, numbers, underscores only).'
      toastStore.error(error.value)
      return false
    }

    loading.value = true
    error.value = null

    try {
      // Validate username against Mojang API
      const response = await fetch(`${MOJANG_API}/${encodeURIComponent(username)}`)

      if (!response.ok) {
        if (response.status === 404) {
          error.value = 'Minecraft account not found. Please check your username.'
          toastStore.error(error.value)
          return false
        }
        throw new Error('Failed to validate username')
      }

      const data = await response.json()

      user.value = {
        username: data.name,
        uuid: data.id,
        email: null,
        purchases: [],
        joinDate: null
      }

      localStorage.setItem('weenie_user', JSON.stringify({
        username: data.name,
        uuid: data.id
      }))

      toastStore.success(`Welcome, ${data.name}!`)

      // Fetch purchases after login
      await fetchPurchases()

      return true
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to login'
      toastStore.error(error.value)
      return false
    } finally {
      loading.value = false
    }
  }

  async function fetchPurchases(): Promise<void> {
    if (!user.value) return

    loading.value = true
    error.value = null

    try {
      // Tebex does not provide a public API for purchase history lookup by username
      // This would require a backend proxy with the secret key
      // For now, we'll use mock data or skip if no backend is available

      // Mock implementation - in production, this would call your backend
      // which would then call Tebex's authenticated API
      purchases.value = []

      // Update user's purchases reference
      if (user.value) {
        user.value.purchases = purchases.value
      }
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to fetch purchases'
      console.error('Failed to fetch purchases:', e)
    } finally {
      loading.value = false
    }
  }

  function logout(): void {
    user.value = null
    purchases.value = []
    error.value = null
    localStorage.removeItem('weenie_user')
    toastStore.info('Logged out successfully')
  }

  function restoreSession(): void {
    const saved = localStorage.getItem('weenie_user')
    if (saved) {
      try {
        const data = JSON.parse(saved)
        user.value = {
          username: data.username,
          uuid: data.uuid,
          email: null,
          purchases: [],
          joinDate: null
        }
        // Fetch purchases in background
        fetchPurchases()
      } catch (e) {
        console.error('Failed to restore session:', e)
        localStorage.removeItem('weenie_user')
      }
    }
  }

  // Auto-restore session on store creation (client-side only)
  if (typeof window !== 'undefined') {
    restoreSession()
  }

  return {
    user,
    purchases,
    loading,
    error,
    isAuthenticated,
    totalSpent,
    purchaseCount,
    isValidMinecraftUsername,
    login,
    fetchPurchases,
    logout,
    restoreSession
  }
})
