# Billion User App - Frontend

Modern React frontend built with Next.js 14, TypeScript, and Tailwind CSS.

## Features

- 🔐 Authentication (Login/Register)
- 👥 User Management
- 📦 Product Management
- ✅ Task Management
- 🖼️ Media Management
- 📱 Responsive Design
- 🎨 Modern UI with Tailwind CSS

## Getting Started

### Prerequisites

- Node.js 18+ and npm/yarn
- Backend services running (see main README)

### Installation

```bash
cd frontend
npm install
```

### Environment Variables

Create a `.env.local` file (optional - defaults work for local dev):

```env
NEXT_PUBLIC_API_AUTH_URL=http://localhost:3001
NEXT_PUBLIC_API_USER_URL=http://localhost:3002
NEXT_PUBLIC_API_PRODUCT_URL=http://localhost:3003
NEXT_PUBLIC_API_TASK_URL=http://localhost:3004
NEXT_PUBLIC_API_MEDIA_URL=http://localhost:3005
```

### Run Development Server

```bash
npm run dev
```

Open [http://localhost:3000](http://localhost:3000) in your browser.

### Build for Production

```bash
npm run build
npm start
```

## Project Structure

```
frontend/
├── app/                    # Next.js App Router pages
│   ├── dashboard/          # Dashboard pages
│   ├── login/             # Login page
│   ├── register/          # Register page
│   └── layout.tsx         # Root layout
├── components/            # Reusable components
├── lib/                   # Utilities and API clients
│   ├── api/              # API client functions
│   └── store/            # State management (Zustand)
└── public/               # Static assets
```

## API Integration

The frontend communicates with 5 microservices:

- **Auth Service** (3001): Authentication and user management
- **User Service** (3002): User profiles
- **Product Service** (3003): Products
- **Task Service** (3004): Tasks
- **Media Service** (3005): Media files

## Features

### Authentication

- JWT token-based authentication
- Automatic token refresh
- Protected routes
- Persistent login state

### State Management

- Zustand for global state
- Local storage persistence
- React Hook Form for forms

### UI Components

- Responsive design
- Tailwind CSS styling
- Toast notifications
- Loading states
- Error handling

## Development

### Adding a New Page

1. Create a new file in `app/` directory
2. Use the `Layout` component for consistent navigation
3. Add route to navigation in `components/Layout.tsx`

### Adding a New API Endpoint

1. Add function to appropriate API file in `lib/api/`
2. Use the service-specific client (authClient, userClient, etc.)
3. Handle errors with `handleApiError` utility

## Troubleshooting

### CORS Errors

Make sure backend services have CORS enabled (already configured in services).

### Authentication Issues

- Check if tokens are stored in localStorage
- Verify backend services are running
- Check browser console for errors

### API Connection Errors

- Verify backend services are running on correct ports
- Check `.env.local` configuration
- Ensure network connectivity

## Tech Stack

- **Next.js 14** - React framework
- **TypeScript** - Type safety
- **Tailwind CSS** - Styling
- **Zustand** - State management
- **React Hook Form** - Form handling
- **Axios** - HTTP client
- **React Hot Toast** - Notifications
- **Lucide React** - Icons

