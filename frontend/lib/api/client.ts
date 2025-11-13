import axios, { AxiosInstance, AxiosError } from 'axios'

const API_BASE_URLS = {
  auth: process.env.NEXT_PUBLIC_API_AUTH_URL || 'http://localhost:3001',
  user: process.env.NEXT_PUBLIC_API_USER_URL || 'http://localhost:3002',
  product: process.env.NEXT_PUBLIC_API_PRODUCT_URL || 'http://localhost:3003',
  task: process.env.NEXT_PUBLIC_API_TASK_URL || 'http://localhost:3004',
  media: process.env.NEXT_PUBLIC_API_MEDIA_URL || 'http://localhost:3005',
}

// Create axios instances for each service
export const authClient: AxiosInstance = axios.create({
  baseURL: API_BASE_URLS.auth,
  headers: {
    'Content-Type': 'application/json',
  },
})

export const userClient: AxiosInstance = axios.create({
  baseURL: API_BASE_URLS.user,
  headers: {
    'Content-Type': 'application/json',
  },
})

export const productClient: AxiosInstance = axios.create({
  baseURL: API_BASE_URLS.product,
  headers: {
    'Content-Type': 'application/json',
  },
})

export const taskClient: AxiosInstance = axios.create({
  baseURL: API_BASE_URLS.task,
  headers: {
    'Content-Type': 'application/json',
  },
})

export const mediaClient: AxiosInstance = axios.create({
  baseURL: API_BASE_URLS.media,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Add token interceptor for authenticated requests
const addAuthToken = (config: any) => {
  if (typeof window !== 'undefined') {
    const token = localStorage.getItem('access_token')
    if (token && token.trim() !== '') {
      config.headers = config.headers || {}
      config.headers.Authorization = `Bearer ${token}`
    }
  }
  return config
}

// Apply interceptors to all clients (including auth for protected endpoints like /auth/profile)
[authClient, userClient, productClient, taskClient, mediaClient].forEach((client) => {
  client.interceptors.request.use(addAuthToken)
  client.interceptors.response.use(
    (response) => response,
    async (error: AxiosError) => {
      if (error.response?.status === 401) {
        // Token expired, try to refresh
        const refreshToken = localStorage.getItem('refresh_token')
        if (refreshToken) {
          try {
            const response = await authClient.post('/api/v1/refresh', {
              refresh_token: refreshToken,
            })
            const { access_token, refresh_token: newRefreshToken } = response.data
            localStorage.setItem('access_token', access_token)
            if (newRefreshToken) {
              localStorage.setItem('refresh_token', newRefreshToken)
            }
            // Retry original request
            if (error.config) {
              error.config.headers.Authorization = `Bearer ${access_token}`
              return axios.request(error.config)
            }
          } catch (refreshError) {
            // Refresh failed, redirect to login
            localStorage.removeItem('access_token')
            localStorage.removeItem('refresh_token')
            if (typeof window !== 'undefined') {
              window.location.href = '/login'
            }
          }
        } else {
          // No refresh token, redirect to login
          if (typeof window !== 'undefined') {
            window.location.href = '/login'
          }
        }
      }
      return Promise.reject(error)
    }
  )
})

export const handleApiError = (error: unknown): string => {
  if (axios.isAxiosError(error)) {
    const axiosError = error as AxiosError<{ error?: string; message?: string }>
    return (
      axiosError.response?.data?.error ||
      axiosError.response?.data?.message ||
      axiosError.message ||
      'An error occurred'
    )
  }
  return 'An unexpected error occurred'
}

