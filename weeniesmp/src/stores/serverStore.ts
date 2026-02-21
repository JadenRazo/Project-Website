import { defineStore } from 'pinia'
import { ref } from 'vue'
import { apiClient, ApiError } from '@/utils/api'

interface McStatusResponse {
  online: boolean
  version?: {
    name_clean?: string
  }
  players?: {
    online?: number
    max?: number
  }
}

interface ServerStatus {
  online: boolean
  players: {
    online: number
    max: number
  }
  version: string
}

export const useServerStore = defineStore('server', () => {
  const status = ref<ServerStatus>({
    online: false,
    players: { online: 0, max: 100 },
    version: '1.21.8'
  })
  const loading = ref(true)
  const error = ref<string | null>(null)

  async function fetchStatus() {
    loading.value = true
    error.value = null

    try {
      const data = await apiClient.get<McStatusResponse>(
        'https://api.mcstatus.io/v2/status/java/play.weeniesmp.net',
        { retries: 3 }
      )

      // Clean up version string (remove "Velocity " prefix if present)
      let versionStr = data.version?.name_clean ?? '1.21.8'
      versionStr = versionStr.replace(/^Velocity\s*/i, '')

      status.value = {
        online: data.online ?? false,
        players: {
          online: data.players?.online ?? 0,
          max: data.players?.max ?? 100
        },
        version: versionStr
      }
    } catch (e) {
      if (e instanceof ApiError) {
        error.value = `Failed to fetch server status: ${e.message}`
      } else {
        error.value = 'Failed to fetch server status'
      }
      status.value.online = false
      console.error('Server status fetch error:', e)
    } finally {
      loading.value = false
    }
  }

  return {
    status,
    loading,
    error,
    fetchStatus
  }
})
