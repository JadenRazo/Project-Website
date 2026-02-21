import { ref, type Ref } from 'vue'

interface UseRetryOptions {
  maxRetries?: number
  baseDelay?: number
  maxDelay?: number
}

interface UseRetryReturn<T> {
  execute: (fn: () => Promise<T>) => Promise<T>
  attempts: Ref<number>
  error: Ref<Error | null>
  isRetrying: Ref<boolean>
  reset: () => void
}

/**
 * Composable for handling retry logic with exponential backoff
 */
export function useRetry<T = unknown>(options: UseRetryOptions = {}): UseRetryReturn<T> {
  const { maxRetries = 3, baseDelay = 1000, maxDelay = 30000 } = options

  const attempts = ref(0)
  const error = ref<Error | null>(null)
  const isRetrying = ref(false)

  /**
   * Sleep for a specified duration
   */
  function sleep(ms: number): Promise<void> {
    return new Promise(resolve => setTimeout(resolve, ms))
  }

  /**
   * Calculate delay with exponential backoff
   */
  function getDelay(attempt: number): number {
    const delay = baseDelay * Math.pow(2, attempt)
    return Math.min(delay, maxDelay)
  }

  /**
   * Execute a function with retry logic
   */
  async function execute(fn: () => Promise<T>): Promise<T> {
    attempts.value = 0
    error.value = null
    isRetrying.value = false

    while (attempts.value < maxRetries) {
      try {
        attempts.value++
        const result = await fn()
        isRetrying.value = false
        return result
      } catch (e) {
        error.value = e instanceof Error ? e : new Error(String(e))

        if (attempts.value < maxRetries) {
          isRetrying.value = true
          const delay = getDelay(attempts.value - 1)
          await sleep(delay)
        } else {
          isRetrying.value = false
          throw error.value
        }
      }
    }

    throw error.value ?? new Error('Max retries exceeded')
  }

  /**
   * Reset the retry state
   */
  function reset(): void {
    attempts.value = 0
    error.value = null
    isRetrying.value = false
  }

  return {
    execute,
    attempts,
    error,
    isRetrying,
    reset
  }
}
