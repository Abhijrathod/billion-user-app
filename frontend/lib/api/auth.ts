import { authClient } from './client'
import { handleApiError } from './client'

export interface RegisterRequest {
  email: string
  username: string
  password: string
}

export interface LoginRequest {
  email: string
  password: string
}

export interface AuthResponse {
  access_token: string
  refresh_token: string
  token_type: string
}

export interface User {
  id: number
  email: string
  username: string
  is_active: boolean
}

export const authApi = {
  register: async (data: RegisterRequest): Promise<{ message: string; user: User }> => {
    try {
      const response = await authClient.post('/api/v1/register', data)
      return response.data
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },

  login: async (data: LoginRequest): Promise<AuthResponse> => {
    try {
      const response = await authClient.post('/api/v1/login', data)
      return response.data
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },

  refreshToken: async (refreshToken: string): Promise<AuthResponse> => {
    try {
      const response = await authClient.post('/api/v1/refresh', {
        refresh_token: refreshToken,
      })
      return response.data
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },

  logout: async (refreshToken: string): Promise<void> => {
    try {
      await authClient.post('/api/v1/logout', {
        refresh_token: refreshToken,
      })
    } catch (error) {
      // Log error but don't throw - logout should always succeed
      console.error('Logout error:', error)
    }
  },

  getProfile: async (): Promise<User> => {
    try {
      const response = await authClient.get('/api/v1/auth/profile')
      return response.data
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },
}

