import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useToastStore } from './toastStore'
import { ApiClient, ApiError } from '@/utils/api'

const TEBEX_TOKEN = import.meta.env.VITE_TEBEX_PUBLIC_TOKEN
if (!TEBEX_TOKEN) {
  console.error('Missing VITE_TEBEX_PUBLIC_TOKEN environment variable')
}
const API_BASE = `https://headless.tebex.io/api/accounts/${TEBEX_TOKEN}`
const BASKET_API = 'https://headless.tebex.io/api/baskets'

// Create dedicated API client for Tebex with retry disabled for basket operations
const tebexApi = new ApiClient({ timeout: 30000, maxRetries: 1 })
const tebexApiWithRetry = new ApiClient({ timeout: 30000, maxRetries: 3 })

export interface TebexPackage {
  id: number
  name: string
  description: string
  image: string | null
  base_price: number
  sales_price: number
  total_price: number
  currency: string
  category: {
    id: number
    name: string
  }
  type: string
}

export interface TebexCategory {
  id: number
  name: string
  packages: TebexPackage[]
}

export interface TebexBasketPackage {
  id: number
  name: string
  description?: string
  image?: string | null
  // Tebex API nests quantity/price inside in_basket
  in_basket: {
    quantity: number
    price: number
    gift_username_valid?: boolean
  }
}

export interface TebexBasket {
  ident: string
  complete: boolean
  email: string | null
  username: string | null
  base_price: number
  sales_tax: number
  total_price: number
  currency: string
  packages: TebexBasketPackage[]
  coupons?: Array<{
    id: string
    code: string
  }>
  creator_code?: string
  links: {
    checkout: string
  }
}

export const useTebexStore = defineStore('tebex', () => {
  const toastStore = useToastStore()

  const packages = ref<TebexPackage[]>([])
  const categories = ref<TebexCategory[]>([])
  const basket = ref<TebexBasket | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)
  const username = ref<string | null>(localStorage.getItem('tebex_username'))
  const basketInitialized = ref(false)
  const basketCreating = ref(false) // Mutex for basket creation

  const basketTotal = computed(() => basket.value?.total_price ?? 0)
  const basketItemCount = computed(() =>
    basket.value?.packages.reduce((sum, p) => sum + (p.in_basket?.quantity ?? 0), 0) ?? 0
  )

  // Initialize basket from localStorage on store creation
  async function initBasket() {
    if (basketInitialized.value) return
    basketInitialized.value = true

    const savedIdent = localStorage.getItem('tebex_basket')
    if (!savedIdent) return

    try {
      const data = await tebexApi.get<{ data: TebexBasket }>(`${BASKET_API}/${savedIdent}`)
      if (!data.data.complete) {
        basket.value = data.data
      } else {
        // Basket was completed, clear it
        localStorage.removeItem('tebex_basket')
      }
    } catch (e) {
      console.error('Failed to restore basket:', e)
      localStorage.removeItem('tebex_basket')
    }
  }

  // Auto-init basket when store is created (only on client side)
  if (typeof window !== 'undefined') {
    initBasket()
  }

  async function fetchCategories() {
    loading.value = true
    error.value = null
    try {
      const data = await tebexApiWithRetry.get<{ data: TebexCategory[] }>(
        `${API_BASE}/categories?includePackages=1`
      )
      categories.value = data.data || []
      packages.value = categories.value.flatMap(c => c.packages || [])
    } catch (e) {
      error.value = e instanceof ApiError ? e.message : (e instanceof Error ? e.message : 'Unknown error')
      console.error('Tebex fetch error:', e)
    } finally {
      loading.value = false
    }
  }

  async function fetchPackages() {
    loading.value = true
    error.value = null
    try {
      const data = await tebexApiWithRetry.get<{ data: TebexPackage[] }>(`${API_BASE}/packages`)
      packages.value = data.data || []
    } catch (e) {
      error.value = e instanceof ApiError ? e.message : (e instanceof Error ? e.message : 'Unknown error')
      console.error('Tebex fetch error:', e)
    } finally {
      loading.value = false
    }
  }

  function isValidMinecraftUsername(name: string): boolean {
    return /^[a-zA-Z0-9_]{3,16}$/.test(name)
  }

  function setUsername(name: string): boolean {
    if (!isValidMinecraftUsername(name)) {
      toastStore.error('Invalid username. Must be 3-16 characters (letters, numbers, underscores only).')
      return false
    }
    username.value = name
    localStorage.setItem('tebex_username', name)
    // Clear old basket when username changes
    basket.value = null
    localStorage.removeItem('tebex_basket')
    toastStore.success(`Welcome, ${name}!`)
    return true
  }

  function logout() {
    username.value = null
    localStorage.removeItem('tebex_username')
    basket.value = null
    localStorage.removeItem('tebex_basket')
  }

  async function createBasket() {
    if (!username.value) {
      error.value = 'Please enter your Minecraft username first'
      toastStore.error('Please enter your Minecraft username first')
      return null
    }

    // Prevent concurrent basket creation
    if (basketCreating.value) {
      await new Promise(resolve => setTimeout(resolve, 100))
      return basket.value
    }

    basketCreating.value = true
    try {
      const data = await tebexApi.post<{ data: TebexBasket }>(
        `${API_BASE}/baskets`,
        { username: username.value }
      )
      basket.value = data.data
      localStorage.setItem('tebex_basket', data.data.ident)
      return data.data
    } catch (e) {
      error.value = e instanceof ApiError ? e.message : (e instanceof Error ? e.message : 'Unknown error')
      console.error('Tebex basket error:', e)
      return null
    } finally {
      basketCreating.value = false
    }
  }

  async function getOrCreateBasket() {
    // Return existing basket if available
    if (basket.value) {
      return basket.value
    }

    // Wait for any in-progress basket creation
    if (basketCreating.value) {
      // Wait and retry
      await new Promise(resolve => setTimeout(resolve, 100))
      return getOrCreateBasket()
    }

    const savedIdent = localStorage.getItem('tebex_basket')
    if (savedIdent) {
      try {
        const data = await tebexApi.get<{ data: TebexBasket }>(`${BASKET_API}/${savedIdent}`)
        if (!data.data.complete) {
          basket.value = data.data
          return data.data
        }
      } catch (e) {
        console.error('Failed to restore basket:', e)
      }
      // Clear invalid basket
      localStorage.removeItem('tebex_basket')
    }
    return await createBasket()
  }

  async function addToBasket(packageId: number, quantity = 1, retried = false): Promise<boolean> {
    const b = await getOrCreateBasket()
    if (!b) return false

    // Find package for toast message and type info
    const pkg = packages.value.find(p => p.id === packageId)
    const pkgName = pkg?.name ?? 'Item'

    // Build request body - include type for packages that support both single and subscription
    const requestBody: Record<string, unknown> = { package_id: packageId, quantity }
    if (pkg?.type === 'both' || pkg?.type === 'single' || pkg?.type === 'subscription') {
      // For 'both' type packages, default to 'single' (one-time purchase)
      // For explicit single/subscription types, pass them through
      requestBody.type = pkg.type === 'both' ? 'single' : pkg.type
    }

    try {
      const data = await tebexApi.post<{ data: TebexBasket }>(
        `${BASKET_API}/${b.ident}/packages`,
        requestBody
      )
      basket.value = data.data
      toastStore.success(`Added ${pkgName} to cart`)
      return true
    } catch (e) {
      // Handle ApiError with detailed error info
      if (e instanceof ApiError) {
        const errorData = e.data as Record<string, string> | null

        // If basket is invalid/expired and we haven't retried, create new basket and retry
        if ((e.status === 400 || e.status === 404) && !retried) {
          const errorLower = (errorData?.message || errorData?.error || '').toLowerCase()
          if (errorLower.includes('basket') || errorLower.includes('not found') || errorLower.includes('invalid') || errorLower.includes('expired')) {
            console.log('Basket appears invalid, creating new basket and retrying...')
            clearBasket()
            return addToBasket(packageId, quantity, true)
          }
        }

        let errorMessage = 'Failed to add to basket'

        if (errorData?.message) {
          errorMessage = errorData.message
        } else if (errorData?.error) {
          errorMessage = errorData.error
        } else if (errorData?.error_message) {
          errorMessage = errorData.error_message
        } else if (e.status === 400) {
          errorMessage = 'This item cannot be added. It may already be in your cart or not available for purchase.'
        } else if (e.status === 404) {
          errorMessage = 'Package not found'
        } else if (e.status === 422) {
          errorMessage = 'Invalid request. Please try again.'
        }

        error.value = errorMessage
        console.error('Tebex add to basket error:', e)
        toastStore.error(errorMessage)
        return false
      }

      const errorMsg = e instanceof Error ? e.message : 'Unknown error'
      error.value = errorMsg
      console.error('Tebex add to basket error:', e)
      toastStore.error(errorMsg)
      return false
    }
  }

  async function removeFromBasket(packageId: number) {
    if (!basket.value) return false

    try {
      const data = await tebexApi.delete<{ data: TebexBasket }>(
        `${BASKET_API}/${basket.value.ident}/packages/${packageId}`
      )
      basket.value = data.data
      toastStore.info('Item removed from cart')
      return true
    } catch (e) {
      error.value = e instanceof ApiError ? e.message : (e instanceof Error ? e.message : 'Unknown error')
      console.error('Tebex remove error:', e)
      toastStore.error('Failed to remove item')
      return false
    }
  }

  async function updateBasketQuantity(packageId: number, quantity: number) {
    if (!basket.value) return false

    try {
      const data = await tebexApi.put<{ data: TebexBasket }>(
        `${BASKET_API}/${basket.value.ident}/packages/${packageId}`,
        { quantity }
      )
      basket.value = data.data
      return true
    } catch (e) {
      error.value = e instanceof ApiError ? e.message : (e instanceof Error ? e.message : 'Unknown error')
      console.error('Tebex update error:', e)
      toastStore.error('Failed to update quantity')
      return false
    }
  }

  async function applyCoupon(code: string) {
    if (!basket.value) {
      toastStore.error('Please add items to cart first')
      return false
    }

    try {
      const data = await tebexApi.post<{ data: TebexBasket }>(
        `${BASKET_API}/${basket.value.ident}/coupons`,
        { coupon_code: code }
      )
      basket.value = data.data
      toastStore.success('Coupon applied!')
      return true
    } catch (e) {
      let msg = 'Failed to apply coupon'
      if (e instanceof ApiError) {
        const errorData = e.data as Record<string, string> | null
        msg = errorData?.message || 'Invalid coupon code'
      } else if (e instanceof Error) {
        msg = e.message
      }
      error.value = msg
      toastStore.error(msg)
      return false
    }
  }

  async function removeCoupon(couponId: string) {
    if (!basket.value) return false

    try {
      const data = await tebexApi.delete<{ data: TebexBasket }>(
        `${BASKET_API}/${basket.value.ident}/coupons/${encodeURIComponent(couponId)}`
      )
      basket.value = data.data
      toastStore.info('Coupon removed')
      return true
    } catch (e) {
      error.value = e instanceof ApiError ? e.message : (e instanceof Error ? e.message : 'Failed to remove coupon')
      toastStore.error('Failed to remove coupon')
      return false
    }
  }

  async function applyCreatorCode(code: string) {
    if (!basket.value) {
      toastStore.error('Please add items to cart first')
      return false
    }

    try {
      const data = await tebexApi.post<{ data: TebexBasket }>(
        `${BASKET_API}/${basket.value.ident}/creator-codes`,
        { creator_code: code }
      )
      basket.value = data.data
      toastStore.success('Creator code applied!')
      return true
    } catch (e) {
      let msg = 'Failed to apply creator code'
      if (e instanceof ApiError) {
        const errorData = e.data as Record<string, string> | null
        msg = errorData?.message || 'Invalid creator code'
      } else if (e instanceof Error) {
        msg = e.message
      }
      error.value = msg
      toastStore.error(msg)
      return false
    }
  }

  async function removeCreatorCode() {
    if (!basket.value) return false

    try {
      const data = await tebexApi.delete<{ data: TebexBasket }>(
        `${BASKET_API}/${basket.value.ident}/creator-codes`
      )
      basket.value = data.data
      toastStore.info('Creator code removed')
      return true
    } catch (e) {
      error.value = e instanceof ApiError ? e.message : (e instanceof Error ? e.message : 'Failed to remove creator code')
      toastStore.error('Failed to remove creator code')
      return false
    }
  }

  function getCheckoutUrl() {
    return basket.value?.links.checkout ?? null
  }

  function checkout() {
    const url = getCheckoutUrl()
    if (url) {
      window.location.href = url
    }
  }

  function clearBasket() {
    basket.value = null
    localStorage.removeItem('tebex_basket')
  }

  return {
    packages,
    categories,
    basket,
    loading,
    error,
    username,
    basketTotal,
    basketItemCount,
    fetchCategories,
    fetchPackages,
    setUsername,
    logout,
    createBasket,
    getOrCreateBasket,
    addToBasket,
    removeFromBasket,
    updateBasketQuantity,
    applyCoupon,
    removeCoupon,
    applyCreatorCode,
    removeCreatorCode,
    getCheckoutUrl,
    checkout,
    clearBasket
  }
})
