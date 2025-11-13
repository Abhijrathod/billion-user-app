/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  env: {
    NEXT_PUBLIC_API_AUTH_URL: process.env.NEXT_PUBLIC_API_AUTH_URL || 'http://localhost:3001',
    NEXT_PUBLIC_API_USER_URL: process.env.NEXT_PUBLIC_API_USER_URL || 'http://localhost:3002',
    NEXT_PUBLIC_API_PRODUCT_URL: process.env.NEXT_PUBLIC_API_PRODUCT_URL || 'http://localhost:3003',
    NEXT_PUBLIC_API_TASK_URL: process.env.NEXT_PUBLIC_API_TASK_URL || 'http://localhost:3004',
    NEXT_PUBLIC_API_MEDIA_URL: process.env.NEXT_PUBLIC_API_MEDIA_URL || 'http://localhost:3005',
  },
}

module.exports = nextConfig

