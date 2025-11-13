import { productClient } from './client'
import { handleApiError } from './client'

export interface Product {
  id: number
  name: string
  description?: string
  price: number
  sku: string
  stock: number
  category?: string
  image_url?: string
  created_by: number
  created_at: string
  updated_at: string
}

export interface CreateProductRequest {
  name: string
  description?: string
  price: number
  sku: string
  stock?: number
  category?: string
  image_url?: string
}

export interface UpdateProductRequest {
  name?: string
  description?: string
  price?: number
  sku?: string
  stock?: number
  category?: string
  image_url?: string
}

export interface ProductListResponse {
  products: Product[]
  offset: number
  limit: number
}

export const productApi = {
  getById: async (id: number): Promise<Product> => {
    try {
      const response = await productClient.get(`/api/v1/products/${id}`)
      return response.data
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },

  list: async (offset = 0, limit = 20): Promise<ProductListResponse> => {
    try {
      const response = await productClient.get('/api/v1/products', {
        params: { offset, limit },
      })
      return response.data
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },

  search: async (query: string, limit = 20): Promise<{ products: Product[]; query: string }> => {
    try {
      const response = await productClient.get('/api/v1/products/search', {
        params: { q: query, limit },
      })
      return response.data
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },

  getByCategory: async (
    category: string,
    offset = 0,
    limit = 20
  ): Promise<{ products: Product[]; category: string; offset: number; limit: number }> => {
    try {
      const response = await productClient.get(`/api/v1/products/category/${category}`, {
        params: { offset, limit },
      })
      return response.data
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },

  create: async (data: CreateProductRequest): Promise<Product> => {
    try {
      const response = await productClient.post('/api/v1/products', data)
      return response.data
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },

  update: async (id: number, data: UpdateProductRequest): Promise<Product> => {
    try {
      const response = await productClient.put(`/api/v1/products/${id}`, data)
      return response.data
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },

  delete: async (id: number): Promise<void> => {
    try {
      await productClient.delete(`/api/v1/products/${id}`)
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },
}

