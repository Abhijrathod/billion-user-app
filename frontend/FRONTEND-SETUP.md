# Frontend Setup Guide

## Quick Start

### 1. Install Dependencies

```bash
cd frontend
npm install
```

### 2. Start Development Server

```bash
npm run dev
```

The frontend will be available at [http://localhost:3000](http://localhost:3000)

### 3. Make Sure Backend Services Are Running

Before using the frontend, ensure all backend services are running:

```bash
# From project root
docker-compose up -d  # Start infrastructure
# Then start each service (see main README)
```

## Features

✅ **Authentication**
- Login/Register pages
- JWT token management
- Automatic token refresh
- Protected routes

✅ **Dashboard**
- Overview with statistics
- Quick actions
- Navigation to all sections

✅ **User Management**
- View all users
- Search users
- User profiles

✅ **Product Management**
- List products
- Create new products
- Search products
- Delete products

✅ **Task Management**
- List tasks
- Create tasks
- Filter by status
- Update and delete tasks

✅ **Media Management**
- View media files
- Upload media (presigned URLs)
- Delete media

## Project Structure

```
frontend/
├── app/                    # Next.js pages
│   ├── dashboard/         # Dashboard pages
│   ├── login/             # Login page
│   └── register/          # Register page
├── components/            # Reusable components
│   └── Layout.tsx         # Main layout with navigation
├── lib/
│   ├── api/              # API clients
│   │   ├── auth.ts       # Auth API
│   │   ├── user.ts       # User API
│   │   ├── product.ts    # Product API
│   │   ├── task.ts       # Task API
│   │   ├── media.ts      # Media API
│   │   └── client.ts     # Axios clients
│   └── store/            # State management
│       └── authStore.ts  # Auth state (Zustand)
└── public/               # Static assets
```

## Environment Variables

Optional - defaults work for local development:

```env
NEXT_PUBLIC_API_AUTH_URL=http://localhost:3001
NEXT_PUBLIC_API_USER_URL=http://localhost:3002
NEXT_PUBLIC_API_PRODUCT_URL=http://localhost:3003
NEXT_PUBLIC_API_TASK_URL=http://localhost:3004
NEXT_PUBLIC_API_MEDIA_URL=http://localhost:3005
```

## Usage

1. **Start the frontend**: `npm run dev`
2. **Register a new account** at `/register`
3. **Login** at `/login`
4. **Access dashboard** at `/dashboard`
5. **Navigate** using the top navigation bar

## Tech Stack

- **Next.js 14** - React framework with App Router
- **TypeScript** - Type safety
- **Tailwind CSS** - Utility-first CSS
- **Zustand** - Lightweight state management
- **React Hook Form** - Form handling
- **Axios** - HTTP client
- **React Hot Toast** - Toast notifications
- **Lucide React** - Icon library

## Build for Production

```bash
npm run build
npm start
```

## Troubleshooting

### CORS Errors
- Ensure backend services have CORS enabled (already configured)
- Check that services are running on correct ports

### Authentication Issues
- Clear browser localStorage
- Check browser console for errors
- Verify backend auth service is running

### API Connection Errors
- Verify all backend services are running
- Check environment variables
- Ensure network connectivity

