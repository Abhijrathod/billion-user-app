'use client'

import { useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { useAuthStore } from '@/lib/store/authStore'
import Layout from '@/components/Layout'
import { User, Package, CheckSquare, Image } from 'lucide-react'
import Link from 'next/link'

export default function DashboardPage() {
  const router = useRouter()
  const { isAuthenticated, user } = useAuthStore()

  useEffect(() => {
    if (!isAuthenticated) {
      router.push('/login')
    }
  }, [isAuthenticated, router])

  if (!isAuthenticated) {
    return null
  }

  const stats = [
    { label: 'Users', value: '0', icon: User, href: '/dashboard/users', color: 'bg-blue-500' },
    { label: 'Products', value: '0', icon: Package, href: '/dashboard/products', color: 'bg-green-500' },
    { label: 'Tasks', value: '0', icon: CheckSquare, href: '/dashboard/tasks', color: 'bg-yellow-500' },
    { label: 'Media', value: '0', icon: Image, href: '/dashboard/media', color: 'bg-purple-500' },
  ]

  return (
    <Layout>
      <div>
        <h1 className="text-3xl font-bold text-gray-900 mb-2">Dashboard</h1>
        <p className="text-gray-600 mb-8">
          Welcome back, {user?.username || 'User'}!
        </p>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
          {stats.map((stat) => {
            const Icon = stat.icon
            return (
              <Link
                key={stat.label}
                href={stat.href}
                className="bg-white rounded-lg shadow-md p-6 hover:shadow-lg transition-shadow"
              >
                <div className="flex items-center justify-between">
                  <div>
                    <p className="text-sm font-medium text-gray-600">{stat.label}</p>
                    <p className="text-2xl font-bold text-gray-900 mt-2">{stat.value}</p>
                  </div>
                  <div className={`${stat.color} p-3 rounded-lg`}>
                    <Icon className="text-white" size={24} />
                  </div>
                </div>
              </Link>
            )
          })}
        </div>

        <div className="bg-white rounded-lg shadow-md p-6">
          <h2 className="text-xl font-semibold mb-4">Quick Actions</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <Link
              href="/dashboard/products/new"
              className="p-4 border-2 border-dashed border-gray-300 rounded-lg hover:border-primary-500 hover:bg-primary-50 transition-colors text-center"
            >
              <Package className="mx-auto mb-2 text-gray-400" size={32} />
              <p className="font-medium">Create Product</p>
            </Link>
            <Link
              href="/dashboard/tasks/new"
              className="p-4 border-2 border-dashed border-gray-300 rounded-lg hover:border-primary-500 hover:bg-primary-50 transition-colors text-center"
            >
              <CheckSquare className="mx-auto mb-2 text-gray-400" size={32} />
              <p className="font-medium">Create Task</p>
            </Link>
          </div>
        </div>
      </div>
    </Layout>
  )
}

