import { handleApiError } from './errorHandler';

interface ApiConfig {
  baseUrl: string;
  timeout: number;
  headers: Record<string, string>;
}

interface ApiRequestOptions extends RequestInit {
  timeout?: number;
  skipAuth?: boolean;
  customHeaders?: Record<string, string>;
}

class ApiClient {
  private static instance: ApiClient;
  private config: ApiConfig;
  private authTokenKey = 'auth_token';

  private constructor() {
    this.config = this.getDefaultConfig();
  }

  static getInstance(): ApiClient {
    if (!ApiClient.instance) {
      ApiClient.instance = new ApiClient();
    }
    return ApiClient.instance;
  }

  private getDefaultConfig(): ApiConfig {
    // Check if we're running in development or production
    const isDevelopment = process.env.NODE_ENV === 'development';
    
    // In development, use proxy (configured in setupProxy.js)
    // In production, use relative URLs
    const baseUrl = '';

    return {
      baseUrl,
      timeout: 30000, // 30 seconds
      headers: {
        'Content-Type': 'application/json',
      },
    };
  }

  private getAuthToken(): string | null {
    return localStorage.getItem(this.authTokenKey);
  }

  setAuthToken(token: string): void {
    localStorage.setItem(this.authTokenKey, token);
  }

  clearAuthToken(): void {
    localStorage.removeItem(this.authTokenKey);
  }

  private async fetchWithTimeout(
    url: string,
    options: RequestInit,
    timeout: number
  ): Promise<Response> {
    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), timeout);

    try {
      const response = await fetch(url, {
        ...options,
        signal: controller.signal,
      });
      clearTimeout(timeoutId);
      return response;
    } catch (error) {
      clearTimeout(timeoutId);
      if (error instanceof Error && error.name === 'AbortError') {
        throw new Error('Request timeout');
      }
      throw error;
    }
  }

  private buildUrl(endpoint: string): string {
    // Remove leading slash if present
    const cleanEndpoint = endpoint.startsWith('/') ? endpoint : `/${endpoint}`;
    return `${this.config.baseUrl}${cleanEndpoint}`;
  }

  private buildHeaders(options: ApiRequestOptions): Headers {
    const headers = new Headers(this.config.headers);

    // Add auth token if available and not skipped
    if (!options.skipAuth) {
      const token = this.getAuthToken();
      if (token) {
        headers.set('Authorization', `Bearer ${token}`);
      }
    }

    // Add custom headers
    if (options.customHeaders) {
      Object.entries(options.customHeaders).forEach(([key, value]) => {
        headers.set(key, value);
      });
    }

    // Add existing headers from options
    if (options.headers) {
      const optionHeaders = options.headers instanceof Headers 
        ? options.headers 
        : new Headers(options.headers as HeadersInit);
      
      optionHeaders.forEach((value, key) => {
        headers.set(key, value);
      });
    }

    return headers;
  }

  async request<T = any>(
    endpoint: string,
    options: ApiRequestOptions = {}
  ): Promise<T> {
    const url = this.buildUrl(endpoint);
    const timeout = options.timeout || this.config.timeout;

    try {
      const response = await this.fetchWithTimeout(
        url,
        {
          ...options,
          headers: this.buildHeaders(options),
        },
        timeout
      );

      // Handle non-2xx responses
      if (!response.ok) {
        const errorData = await response.json().catch(() => ({
          error: `HTTP ${response.status}: ${response.statusText}`,
        }));
        
        const error = new Error(errorData.error || errorData.message || 'Request failed');
        (error as any).status = response.status;
        (error as any).data = errorData;
        throw error;
      }

      // Handle empty responses
      if (response.status === 204 || response.headers.get('content-length') === '0') {
        return {} as T;
      }

      // Parse JSON response
      const data = await response.json();
      return data as T;
    } catch (error) {
      handleApiError(
        error,
        endpoint,
        options.method || 'GET',
        'Failed to complete request'
      );
      throw error;
    }
  }

  // Convenience methods
  async get<T = any>(endpoint: string, options?: ApiRequestOptions): Promise<T> {
    return this.request<T>(endpoint, { ...options, method: 'GET' });
  }

  async post<T = any>(
    endpoint: string,
    body?: any,
    options?: ApiRequestOptions
  ): Promise<T> {
    return this.request<T>(endpoint, {
      ...options,
      method: 'POST',
      body: body ? JSON.stringify(body) : undefined,
    });
  }

  async put<T = any>(
    endpoint: string,
    body?: any,
    options?: ApiRequestOptions
  ): Promise<T> {
    return this.request<T>(endpoint, {
      ...options,
      method: 'PUT',
      body: body ? JSON.stringify(body) : undefined,
    });
  }

  async patch<T = any>(
    endpoint: string,
    body?: any,
    options?: ApiRequestOptions
  ): Promise<T> {
    return this.request<T>(endpoint, {
      ...options,
      method: 'PATCH',
      body: body ? JSON.stringify(body) : undefined,
    });
  }

  async delete<T = any>(endpoint: string, options?: ApiRequestOptions): Promise<T> {
    return this.request<T>(endpoint, { ...options, method: 'DELETE' });
  }
}

// Export singleton instance
export const apiClient = ApiClient.getInstance();

// Export convenience functions
export const api = {
  get: <T = any>(endpoint: string, options?: ApiRequestOptions) => 
    apiClient.get<T>(endpoint, options),
  post: <T = any>(endpoint: string, body?: any, options?: ApiRequestOptions) => 
    apiClient.post<T>(endpoint, body, options),
  put: <T = any>(endpoint: string, body?: any, options?: ApiRequestOptions) => 
    apiClient.put<T>(endpoint, body, options),
  patch: <T = any>(endpoint: string, body?: any, options?: ApiRequestOptions) => 
    apiClient.patch<T>(endpoint, body, options),
  delete: <T = any>(endpoint: string, options?: ApiRequestOptions) => 
    apiClient.delete<T>(endpoint, options),
  setAuthToken: (token: string) => apiClient.setAuthToken(token),
  clearAuthToken: () => apiClient.clearAuthToken(),
};