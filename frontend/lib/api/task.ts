import { taskClient } from './client'
import { handleApiError } from './client'

export type TaskStatus = 'pending' | 'in_progress' | 'completed' | 'cancelled'

export interface Task {
  id: number
  user_id: number
  title: string
  description?: string
  status: TaskStatus
  priority: number
  due_date?: string
  created_at: string
  updated_at: string
}

export interface CreateTaskRequest {
  title: string
  description?: string
  status?: TaskStatus
  priority?: number
  due_date?: string
}

export interface UpdateTaskRequest {
  title?: string
  description?: string
  status?: TaskStatus
  priority?: number
  due_date?: string
}

export interface TaskListResponse {
  tasks: Task[]
  offset: number
  limit: number
}

export const taskApi = {
  getById: async (id: number): Promise<Task> => {
    try {
      const response = await taskClient.get(`/api/v1/tasks/${id}`)
      return response.data
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },

  list: async (offset = 0, limit = 20): Promise<TaskListResponse> => {
    try {
      const response = await taskClient.get('/api/v1/tasks', {
        params: { offset, limit },
      })
      return response.data
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },

  getByStatus: async (
    status: TaskStatus,
    offset = 0,
    limit = 20
  ): Promise<{ tasks: Task[]; status: TaskStatus; offset: number; limit: number }> => {
    try {
      const response = await taskClient.get(`/api/v1/tasks/status/${status}`, {
        params: { offset, limit },
      })
      return response.data
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },

  create: async (data: CreateTaskRequest): Promise<Task> => {
    try {
      const response = await taskClient.post('/api/v1/tasks', data)
      return response.data
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },

  update: async (id: number, data: UpdateTaskRequest): Promise<Task> => {
    try {
      const response = await taskClient.put(`/api/v1/tasks/${id}`, data)
      return response.data
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },

  delete: async (id: number): Promise<void> => {
    try {
      await taskClient.delete(`/api/v1/tasks/${id}`)
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },
}

