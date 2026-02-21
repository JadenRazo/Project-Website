/**
 * API client with retry logic, error handling, and request/response interceptors
 */

export class ApiError extends Error {
  constructor(
    message: string,
    public status: number,
    public statusText: string,
    public data?: unknown
  ) {
    super(message)
    this.name = 'ApiError'
  }
}

type RequestInterceptor = (config: RequestInit) => RequestInit | Promise<RequestInit>
type ResponseInterceptor = (response: Response) => Response | Promise<Response>

interface ApiClientConfig {
  baseUrl?: string
  timeout?: number
  maxRetries?: number
  baseDelay?: number
}

interface RequestOptions extends RequestInit {
  timeout?: number
  retries?: number
  retryDelay?: number
}

class ApiClient {
  private baseUrl: string
  private timeout: number
  private maxRetries: number
  private baseDelay: number
  private requestInterceptors: RequestInterceptor[] = []
  private responseInterceptors: ResponseInterceptor[] = []

  constructor(config: ApiClientConfig = {}) {
    this.baseUrl = config.baseUrl ?? ''
    this.timeout = config.timeout ?? 30000
    this.maxRetries = config.maxRetries ?? 3
    this.baseDelay = config.baseDelay ?? 1000
  }

  /**
   * Add a request interceptor
   */
  addRequestInterceptor(interceptor: RequestInterceptor): () => void {
    this.requestInterceptors.push(interceptor)
    return () => {
      const index = this.requestInterceptors.indexOf(interceptor)
      if (index > -1) {
        this.requestInterceptors.splice(index, 1)
      }
    }
  }

  /**
   * Add a response interceptor
   */
  addResponseInterceptor(interceptor: ResponseInterceptor): () => void {
    this.responseInterceptors.push(interceptor)
    return () => {
      const index = this.responseInterceptors.indexOf(interceptor)
      if (index > -1) {
        this.responseInterceptors.splice(index, 1)
      }
    }
  }

  /**
   * Sleep for a specified duration
   */
  private sleep(ms: number): Promise<void> {
    return new Promise(resolve => setTimeout(resolve, ms))
  }

  /**
   * Calculate exponential backoff delay
   */
  private getRetryDelay(attempt: number, baseDelay: number): number {
    return Math.min(baseDelay * Math.pow(2, attempt), 30000)
  }

  /**
   * Check if error is retryable
   */
  private isRetryable(status: number): boolean {
    return status >= 500 || status === 408 || status === 429
  }

  /**
   * Execute request with timeout
   */
  private async fetchWithTimeout(
    url: string,
    options: RequestInit,
    timeout: number
  ): Promise<Response> {
    const controller = new AbortController()
    const timeoutId = setTimeout(() => controller.abort(), timeout)

    try {
      const response = await fetch(url, {
        ...options,
        signal: controller.signal
      })
      return response
    } finally {
      clearTimeout(timeoutId)
    }
  }

  /**
   * Core request method with retry logic
   */
  async request<T>(endpoint: string, options: RequestOptions = {}): Promise<T> {
    const {
      timeout = this.timeout,
      retries = this.maxRetries,
      retryDelay = this.baseDelay,
      ...fetchOptions
    } = options

    const url = this.baseUrl ? `${this.baseUrl}${endpoint}` : endpoint

    // Apply request interceptors
    let config: RequestInit = { ...fetchOptions }
    for (const interceptor of this.requestInterceptors) {
      config = await interceptor(config)
    }

    let lastError: Error | null = null
    let attempt = 0

    while (attempt < retries) {
      try {
        let response = await this.fetchWithTimeout(url, config, timeout)

        // Apply response interceptors
        for (const interceptor of this.responseInterceptors) {
          response = await interceptor(response)
        }

        if (!response.ok) {
          // Check if we should retry this error
          if (this.isRetryable(response.status) && attempt < retries - 1) {
            const delay = this.getRetryDelay(attempt, retryDelay)
            await this.sleep(delay)
            attempt++
            continue
          }

          // Parse error response body
          let errorData: unknown
          try {
            errorData = await response.json()
          } catch {
            errorData = null
          }

          throw new ApiError(
            `Request failed: ${response.status} ${response.statusText}`,
            response.status,
            response.statusText,
            errorData
          )
        }

        // Handle empty responses
        const contentType = response.headers.get('content-type')
        if (!contentType || !contentType.includes('application/json')) {
          return {} as T
        }

        const text = await response.text()
        if (!text) {
          return {} as T
        }

        return JSON.parse(text) as T
      } catch (error) {
        if (error instanceof ApiError) {
          throw error
        }

        // Handle abort/timeout errors
        if (error instanceof DOMException && error.name === 'AbortError') {
          lastError = new ApiError('Request timeout', 408, 'Request Timeout')
          if (attempt < retries - 1) {
            const delay = this.getRetryDelay(attempt, retryDelay)
            await this.sleep(delay)
            attempt++
            continue
          }
          throw lastError
        }

        // Network errors - retry
        lastError = error instanceof Error ? error : new Error(String(error))
        if (attempt < retries - 1) {
          const delay = this.getRetryDelay(attempt, retryDelay)
          await this.sleep(delay)
          attempt++
          continue
        }
        throw new ApiError(
          lastError.message || 'Network error',
          0,
          'Network Error'
        )
      }
    }

    throw lastError ?? new ApiError('Max retries exceeded', 0, 'Max Retries')
  }

  /**
   * GET request
   */
  async get<T>(endpoint: string, options: RequestOptions = {}): Promise<T> {
    return this.request<T>(endpoint, {
      ...options,
      method: 'GET'
    })
  }

  /**
   * POST request
   */
  async post<T>(endpoint: string, data?: unknown, options: RequestOptions = {}): Promise<T> {
    return this.request<T>(endpoint, {
      ...options,
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        ...options.headers
      },
      body: data ? JSON.stringify(data) : undefined
    })
  }

  /**
   * PUT request
   */
  async put<T>(endpoint: string, data?: unknown, options: RequestOptions = {}): Promise<T> {
    return this.request<T>(endpoint, {
      ...options,
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        ...options.headers
      },
      body: data ? JSON.stringify(data) : undefined
    })
  }

  /**
   * DELETE request
   */
  async delete<T>(endpoint: string, options: RequestOptions = {}): Promise<T> {
    return this.request<T>(endpoint, {
      ...options,
      method: 'DELETE'
    })
  }

  /**
   * Set the base URL
   */
  setBaseUrl(url: string): void {
    this.baseUrl = url
  }
}

// Export singleton instance
export const apiClient = new ApiClient()

// Export class for custom instances
export { ApiClient }
