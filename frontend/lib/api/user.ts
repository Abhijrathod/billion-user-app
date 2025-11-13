import { userClient } from './client'
import { handleApiError } from './client'

export interface User {
  id: number
  email: string
  username: string
  display_name?: string
  bio?: string
  avatar_url?: string
  region?: string
  created_at: string
  updated_at: string
}

export interface CreateUserRequest {
  email: string
  username: string
  display_name?: string
  bio?: string
  region?: string
}

export interface UpdateUserRequest {
  display_name?: string
  bio?: string
  avatar_url?: string
  region?: string
  metadata?: string
}

export interface UserListResponse {
  users: User[]
  offset: number
  limit: number
}

export const userApi = {
  getById: async (id: number): Promise<User> => {
    try {
      const response = await userClient.get(`/api/v1/users/${id}`)
      return response.data
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },

  getByUsername: async (username: string): Promise<User> => {
    try {
      const response = await userClient.get(`/api/v1/users/username/${username}`)
      return response.data
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },

  list: async (offset = 0, limit = 20): Promise<UserListResponse> => {
    try {
      const response = await userClient.get('/api/v1/users', {
        params: { offset, limit },
      })
      return response.data
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },

  search: async (query: string, limit = 20): Promise<{ users: User[]; query: string }> => {
    try {
      const response = await userClient.get('/api/v1/users/search', {
        params: { q: query, limit },
      })
      return response.data
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },

  create: async (data: CreateUserRequest): Promise<User> => {
    try {
      const response = await userClient.post('/api/v1/users', data)
      return response.data
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },

  update: async (id: number, data: UpdateUserRequest): Promise<User> => {
    try {
      const response = await userClient.put(`/api/v1/users/${id}`, data)
      return response.data
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },

  delete: async (id: number): Promise<void> => {
    try {
      await userClient.delete(`/api/v1/users/${id}`)
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },
}

