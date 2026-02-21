import axios, { AxiosError } from 'axios'

export const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || '/api/v1',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Request interceptor
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response interceptor
api.interceptors.response.use(
  (response) => {
    return response
  },
  (error: AxiosError) => {
    // Don't redirect if there's no response (backend not running)
    if (!error.response) {
      console.warn('API not available, running in offline mode')
      return Promise.reject(error)
    }
    
    if (error.response.status === 401) {
      // Only redirect if we have a token (not demo mode)
      const token = localStorage.getItem('token')
      if (token && !token.startsWith('demo-')) {
        localStorage.removeItem('token')
        window.location.href = '/login'
      }
    }
    return Promise.reject(error)
  }
)

// API response types
export interface ApiResponse<T = unknown> {
  code: number
  message: string
  data: T
}

export interface PaginatedResponse<T = unknown> {
  items: T[]
  total: number
  page: number
  page_size: number
}

// Helper functions
export function getErrorMessage(error: unknown): string {
  if (axios.isAxiosError(error)) {
    const axiosError = error as AxiosError<ApiResponse>
    
    // Network error (backend not running)
    if (!error.response) {
      return '无法连接到服务器，请确保后端服务已启动'
    }
    
    if (axiosError.response?.data?.message) {
      return axiosError.response.data.message
    }
    
    // HTTP status error messages
    const statusMessages: Record<number, string> = {
      400: '请求参数错误',
      401: '登录已过期，请重新登录',
      403: '没有权限访问',
      404: '请求的资源不存在',
      500: '服务器内部错误',
    }
    
    const status = axiosError.response?.status
    if (status && statusMessages[status]) {
      return statusMessages[status]
    }
    
    return axiosError.message
  }
  return 'An unexpected error occurred'
}

export function isSuccessResponse(response: ApiResponse): boolean {
  return response.code === 0
}

export default api
