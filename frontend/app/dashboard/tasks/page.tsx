'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import { useAuthStore } from '@/lib/store/authStore'
import { taskApi, Task, TaskStatus } from '@/lib/api/task'
import Layout from '@/components/Layout'
import toast from 'react-hot-toast'
import Link from 'next/link'
import { Plus } from 'lucide-react'
import { format } from 'date-fns'

export default function TasksPage() {
  const router = useRouter()
  const { isAuthenticated } = useAuthStore()
  const [tasks, setTasks] = useState<Task[]>([])
  const [loading, setLoading] = useState(true)
  const [filter, setFilter] = useState<TaskStatus | 'all'>('all')

  useEffect(() => {
    if (!isAuthenticated) {
      router.push('/login')
      return
    }
    loadTasks()
  }, [isAuthenticated, router, filter])

  const loadTasks = async () => {
    try {
      setLoading(true)
      const response =
        filter === 'all'
          ? await taskApi.list()
          : await taskApi.getByStatus(filter)
      setTasks(response.tasks)
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Failed to load tasks')
    } finally {
      setLoading(false)
    }
  }

  const handleDelete = async (id: number) => {
    if (!confirm('Are you sure you want to delete this task?')) return
    try {
      await taskApi.delete(id)
      toast.success('Task deleted')
      loadTasks()
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Failed to delete task')
    }
  }

  const getStatusColor = (status: TaskStatus) => {
    switch (status) {
      case 'completed':
        return 'bg-green-100 text-green-800'
      case 'in_progress':
        return 'bg-blue-100 text-blue-800'
      case 'pending':
        return 'bg-yellow-100 text-yellow-800'
      case 'cancelled':
        return 'bg-red-100 text-red-800'
      default:
        return 'bg-gray-100 text-gray-800'
    }
  }

  const getPriorityColor = (priority: number) => {
    if (priority >= 2) return 'text-red-600 font-bold'
    if (priority === 1) return 'text-yellow-600 font-semibold'
    return 'text-gray-600'
  }

  return (
    <Layout>
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold text-gray-900">Tasks</h1>
        <Link
          href="/dashboard/tasks/new"
          className="flex items-center space-x-2 bg-primary-600 text-white px-4 py-2 rounded-lg hover:bg-primary-700 transition-colors"
        >
          <Plus size={20} />
          <span>New Task</span>
        </Link>
      </div>

      <div className="mb-6 flex space-x-2">
        {(['all', 'pending', 'in_progress', 'completed', 'cancelled'] as const).map((status) => (
          <button
            key={status}
            onClick={() => setFilter(status)}
            className={`px-4 py-2 rounded-lg transition-colors ${
              filter === status
                ? 'bg-primary-600 text-white'
                : 'bg-white text-gray-700 hover:bg-gray-100'
            }`}
          >
            {status.charAt(0).toUpperCase() + status.slice(1).replace('_', ' ')}
          </button>
        ))}
      </div>

      {loading ? (
        <div className="text-center py-12">
          <p className="text-gray-600">Loading tasks...</p>
        </div>
      ) : tasks.length === 0 ? (
        <div className="text-center py-12">
          <p className="text-gray-600">No tasks found</p>
        </div>
      ) : (
        <div className="space-y-4">
          {tasks.map((task) => (
            <div key={task.id} className="bg-white rounded-lg shadow-md p-6">
              <div className="flex justify-between items-start mb-4">
                <div className="flex-1">
                  <h3 className="text-xl font-semibold mb-2">{task.title}</h3>
                  {task.description && (
                    <p className="text-gray-600 mb-4">{task.description}</p>
                  )}
                  <div className="flex items-center space-x-4 text-sm">
                    <span className={`px-2 py-1 rounded ${getStatusColor(task.status)}`}>
                      {task.status.replace('_', ' ')}
                    </span>
                    <span className={getPriorityColor(task.priority)}>
                      Priority: {task.priority}
                    </span>
                    {task.due_date && (
                      <span className="text-gray-600">
                        Due: {format(new Date(task.due_date), 'MMM dd, yyyy')}
                      </span>
                    )}
                  </div>
                </div>
                <div className="flex space-x-2">
                  <Link
                    href={`/dashboard/tasks/${task.id}`}
                    className="px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 transition-colors"
                  >
                    Edit
                  </Link>
                  <button
                    onClick={() => handleDelete(task.id)}
                    className="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors"
                  >
                    Delete
                  </button>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </Layout>
  )
}

