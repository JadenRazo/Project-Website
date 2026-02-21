import { defineStore } from 'pinia'
import { ref } from 'vue'
import { useToastStore } from './toastStore'

export type AppealStatus = 'pending' | 'under_review' | 'approved' | 'denied'

export interface Appeal {
  id: string
  username: string
  email: string
  banReason: string
  appealText: string
  status: AppealStatus
  createdAt: string
  updatedAt: string
  staffResponse?: string
}

export interface AppealSubmission {
  username: string
  email: string
  banReason: string
  appealText: string
}

export const useAppealStore = defineStore('appeal', () => {
  const toastStore = useToastStore()

  const appeals = ref<Appeal[]>([])
  const currentAppeal = ref<Appeal | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  function generateAppealId(): string {
    const timestamp = Date.now().toString(36).toUpperCase()
    const random = Math.random().toString(36).substring(2, 6).toUpperCase()
    return `WS-${timestamp}-${random}`
  }

  async function submitAppeal(data: AppealSubmission): Promise<Appeal | null> {
    loading.value = true
    error.value = null

    try {
      // Simulate API call delay
      await new Promise(resolve => setTimeout(resolve, 1500))

      const now = new Date().toISOString()
      const newAppeal: Appeal = {
        id: generateAppealId(),
        username: data.username,
        email: data.email,
        banReason: data.banReason,
        appealText: data.appealText,
        status: 'pending',
        createdAt: now,
        updatedAt: now
      }

      // Store in localStorage for demo purposes
      const storedAppeals = JSON.parse(localStorage.getItem('weenie_appeals') || '[]')
      storedAppeals.push(newAppeal)
      localStorage.setItem('weenie_appeals', JSON.stringify(storedAppeals))

      appeals.value.push(newAppeal)
      currentAppeal.value = newAppeal

      toastStore.success('Appeal submitted successfully!')
      return newAppeal
    } catch (e) {
      const errorMsg = e instanceof Error ? e.message : 'Failed to submit appeal'
      error.value = errorMsg
      toastStore.error(errorMsg)
      return null
    } finally {
      loading.value = false
    }
  }

  async function checkAppealStatus(id: string, email: string): Promise<Appeal | null> {
    loading.value = true
    error.value = null

    try {
      // Simulate API call delay
      await new Promise(resolve => setTimeout(resolve, 1000))

      // Check localStorage for demo purposes
      const storedAppeals: Appeal[] = JSON.parse(localStorage.getItem('weenie_appeals') || '[]')
      const appeal = storedAppeals.find(
        a => a.id.toLowerCase() === id.toLowerCase() && a.email.toLowerCase() === email.toLowerCase()
      )

      if (!appeal) {
        error.value = 'No appeal found with that ID and email combination'
        toastStore.error('Appeal not found. Please check your appeal ID and email.')
        return null
      }

      currentAppeal.value = appeal
      return appeal
    } catch (e) {
      const errorMsg = e instanceof Error ? e.message : 'Failed to check appeal status'
      error.value = errorMsg
      toastStore.error(errorMsg)
      return null
    } finally {
      loading.value = false
    }
  }

  function clearCurrentAppeal() {
    currentAppeal.value = null
    error.value = null
  }

  return {
    appeals,
    currentAppeal,
    loading,
    error,
    submitAppeal,
    checkAppealStatus,
    clearCurrentAppeal
  }
})
