'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import { useAuthStore } from '@/lib/store/authStore'
import { mediaApi } from '@/lib/api/media'
import Layout from '@/components/Layout'
import toast from 'react-hot-toast'

interface Media {
  id: number
  user_id: number
  file_name: string
  file_type: string
  file_size: number
  url: string
  thumbnail_url?: string
  created_at: string
}

export default function MediaPage() {
  const router = useRouter()
  const { isAuthenticated } = useAuthStore()
  const [media, setMedia] = useState<Media[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (!isAuthenticated) {
      router.push('/login')
      return
    }
    loadMedia()
  }, [isAuthenticated, router])

  const loadMedia = async () => {
    try {
      setLoading(true)
      const response = await mediaApi.list()
      setMedia(response.media)
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Failed to load media')
    } finally {
      setLoading(false)
    }
  }

  const handleDelete = async (id: number) => {
    if (!confirm('Are you sure you want to delete this media?')) return
    try {
      await mediaApi.delete(id)
      toast.success('Media deleted')
      loadMedia()
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Failed to delete media')
    }
  }

  const formatFileSize = (bytes: number): string => {
    if (bytes === 0) return '0 Bytes'
    const k = 1024
    const sizes = ['Bytes', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
  }

  return (
    <Layout>
      <div>
        <h1 className="text-3xl font-bold text-gray-900 mb-6">Media</h1>

        {loading ? (
          <div className="text-center py-12">
            <p className="text-gray-600">Loading media...</p>
          </div>
        ) : media.length === 0 ? (
          <div className="text-center py-12">
            <p className="text-gray-600">No media files found</p>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {media.map((item) => (
              <div key={item.id} className="bg-white rounded-lg shadow-md overflow-hidden">
                {item.thumbnail_url ? (
                  <img
                    src={item.thumbnail_url}
                    alt={item.file_name}
                    className="w-full h-48 object-cover"
                  />
                ) : (
                  <div className="w-full h-48 bg-gray-200 flex items-center justify-center">
                    <span className="text-gray-400">No preview</span>
                  </div>
                )}
                <div className="p-4">
                  <h3 className="font-semibold mb-2 truncate">{item.file_name}</h3>
                  <p className="text-sm text-gray-600 mb-2">
                    {item.file_type} • {formatFileSize(item.file_size)}
                  </p>
                  <div className="flex space-x-2">
                    <a
                      href={item.url}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="flex-1 text-center px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 transition-colors"
                    >
                      View
                    </a>
                    <button
                      onClick={() => handleDelete(item.id)}
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
      </div>
    </Layout>
  )
}

