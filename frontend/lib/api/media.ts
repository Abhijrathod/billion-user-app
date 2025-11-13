import { mediaClient } from './client'
import { handleApiError } from './client'

export interface Media {
  id: number
  user_id: number
  file_name: string
  file_type: string
  file_size: number
  url: string
  thumbnail_url?: string
  metadata?: string
  created_at: string
  updated_at: string
}

export interface CreateMediaRequest {
  file_name: string
  file_type: string
  file_size: number
  url: string
  thumbnail_url?: string
  metadata?: string
}

export interface PresignedURLRequest {
  file_name: string
  file_type: string
  file_size: number
  expires_in?: number
}

export interface MediaListResponse {
  media: Media[]
  offset: number
  limit: number
}

export const mediaApi = {
  getById: async (id: number): Promise<Media> => {
    try {
      const response = await mediaClient.get(`/api/v1/media/${id}`)
      return response.data
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },

  list: async (offset = 0, limit = 20): Promise<MediaListResponse> => {
    try {
      const response = await mediaClient.get('/api/v1/media', {
        params: { offset, limit },
      })
      return response.data
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },

  create: async (data: CreateMediaRequest): Promise<Media> => {
    try {
      const response = await mediaClient.post('/api/v1/media', data)
      return response.data
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },

  delete: async (id: number): Promise<void> => {
    try {
      await mediaClient.delete(`/api/v1/media/${id}`)
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },

  getPresignedURL: async (data: PresignedURLRequest): Promise<{
    url: string
    key: string
    expires_in: number
  }> => {
    try {
      const response = await mediaClient.post('/api/v1/media/presigned-url', data)
      return response.data
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },
}

