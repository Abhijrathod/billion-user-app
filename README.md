# Billion User App - Production-Ready Go Microservices Architecture

A scalable, production-ready microservices architecture built with Go, designed to support up to 1 billion users. This project implements a complete backend system with authentication, user management, products, tasks, and media services.

## 🏗️ Architecture Overview

This project follows a microservices architecture with the following components:

### Services

1. **Auth Service** (Port 3001)
   - User registration and authentication
   - JWT token generation and validation
   - Refresh token management
   - Password hashing with bcrypt

2. **User Service** (Port 3002)
   - User profile management
   - User search and listing
   - Profile updates

3. **Product Service** (Port 3003)
   - Product CRUD operations
   - Product search and categorization
   - Inventory management

4. **Task Service** (Port 3004)
   - Task management (create, update, delete)
   - Task status tracking
   - User-specific task filtering

5. **Media Service** (Port 3005)
   - Media file metadata management
   - Presigned URL generation (S3-ready)
   - Media listing and deletion

### Shared Packages (`pkg/`)

- **config**: Environment configuration management
- **database**: GORM database connection utilities
- **jwtutils**: JWT token generation and validation
- **kafkaclient**: Kafka event publishing client
- **logger**: Structured logging with zerolog

## 🚀 Quick Start

### Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose
- PostgreSQL 15+ (or use Docker Compose)
- Redis (or use Docker Compose)
- Kafka (or use Docker Compose)

### Local Development Setup

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd billion-user-app
   ```

2. **Set up environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. **Start infrastructure services**
   ```bash
   docker-compose up -d
   ```

   This starts:
   - PostgreSQL (with multiple databases)
   - Redis
   - Kafka + Zookeeper

4. **Run database migrations**
   Migrations run automatically on service startup via GORM AutoMigrate.

5. **Run services locally**

   In separate terminals:
   ```bash
   # Auth Service
   cd services/auth-service
   go run cmd/api/main.go

   # User Service
   cd services/user-service
   go run cmd/api/main.go

   # Product Service
   cd services/product-service
   go run cmd/api/main.go

   # Task Service
   cd services/task-service
   go run cmd/api/main.go

   # Media Service
   cd services/media-service
   go run cmd/api/main.go
   ```

## 📦 Building with Docker

### Build individual services

```bash
# Build auth service
docker build -t auth-service:latest -f services/auth-service/Dockerfile .

# Build all services
for service in auth-service user-service product-service task-service media-service; do
  docker build -t $service:latest -f services/$service/Dockerfile .
done
```

## 🔌 API Endpoints

### Auth Service (Port 3001)

- `POST /api/v1/register` - Register a new user
- `POST /api/v1/login` - Login and get tokens
- `POST /api/v1/refresh` - Refresh access token
- `POST /api/v1/logout` - Logout (invalidate refresh token)
- `GET /api/v1/auth/profile` - Get current user profile (protected)

### User Service (Port 3002)

- `GET /api/v1/users/:id` - Get user by ID
- `GET /api/v1/users/username/:username` - Get user by username
- `GET /api/v1/users` - List users (paginated)
- `GET /api/v1/users/search?q=query` - Search users
- `POST /api/v1/users` - Create user (protected)
- `PUT /api/v1/users/:id` - Update user (protected)
- `DELETE /api/v1/users/:id` - Delete user (protected)

### Product Service (Port 3003)

- `GET /api/v1/products/:id` - Get product by ID
- `GET /api/v1/products` - List products (paginated)
- `GET /api/v1/products/search?q=query` - Search products
- `GET /api/v1/products/category/:category` - Get products by category
- `POST /api/v1/products` - Create product (protected)
- `PUT /api/v1/products/:id` - Update product (protected)
- `DELETE /api/v1/products/:id` - Delete product (protected)

### Task Service (Port 3004)

- `GET /api/v1/tasks` - Get my tasks (protected)
- `GET /api/v1/tasks/status/:status` - Get tasks by status (protected)
- `GET /api/v1/tasks/:id` - Get task by ID (protected)
- `POST /api/v1/tasks` - Create task (protected)
- `PUT /api/v1/tasks/:id` - Update task (protected)
- `DELETE /api/v1/tasks/:id` - Delete task (protected)

### Media Service (Port 3005)

- `GET /api/v1/media` - Get my media (protected)
- `GET /api/v1/media/:id` - Get media by ID (protected)
- `POST /api/v1/media` - Create media record (protected)
- `POST /api/v1/media/presigned-url` - Get presigned URL for upload (protected)
- `DELETE /api/v1/media/:id` - Delete media (protected)

## 🔐 Authentication

All protected endpoints require a JWT token in the Authorization header:

```
Authorization: Bearer <access_token>
```

### Example: Register and Login

```bash
# Register
curl -X POST http://localhost:3001/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "username": "johndoe",
    "password": "securepassword123"
  }'

# Login
curl -X POST http://localhost:3001/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "securepassword123"
  }'

# Use the access_token in subsequent requests
curl -X GET http://localhost:3002/api/v1/users/1 \
  -H "Authorization: Bearer <access_token>"
```

## 🏭 Production Considerations

### Scaling Strategy

1. **Horizontal Scaling**: All services are stateless and can be scaled horizontally
2. **Database Sharding**: Use CockroachDB, Vitess, or Citus for distributed SQL
3. **Caching**: Redis for session storage and hot data caching
4. **Message Queue**: Kafka for event streaming and async processing
5. **CDN**: Use CloudFront/Cloudflare for static assets and media

### Multi-Region Deployment

- Use geo-routing (Route53 latency-based)
- Multi-region database replication (CockroachDB or managed Postgres)
- Cross-region Kafka replication
- Redis Cluster with regional replicas

### Monitoring & Observability

- **Metrics**: Prometheus + Grafana
- **Tracing**: Jaeger for distributed tracing
- **Logging**: Centralized logging with Loki or ELK stack
- **Health Checks**: `/health` endpoint on each service

## 📊 Database Schema

Each service has its own database:
- `auth_db` - Authentication and refresh tokens
- `user_db` - User profiles
- `product_db` - Products
- `task_db` - Tasks
- `media_db` - Media metadata

## 🔄 Event-Driven Architecture

Services publish events to Kafka topics:
- `user.created` - When a user is created
- `user.updated` - When a user is updated
- `product.created` - When a product is created
- `task.created` - When a task is created

Consumers can subscribe to these events for analytics, notifications, or other processing.

## 🧪 Testing

```bash
# Run tests for a specific service
cd services/auth-service
go test ./...

# Run all tests
go test ./...
```

## 📝 Project Structure

```
billion-user-app/
├── services/
│   ├── auth-service/
│   │   ├── cmd/api/main.go
│   │   ├── internal/
│   │   │   ├── domain/
│   │   │   ├── handler/
│   │   │   ├── repository/
│   │   │   └── service/
│   │   └── Dockerfile
│   ├── user-service/
│   ├── product-service/
│   ├── task-service/
│   └── media-service/
├── pkg/
│   ├── config/
│   ├── database/
│   ├── jwtutils/
│   ├── kafkaclient/
│   └── logger/
├── k8s/              # Kubernetes manifests
├── infra/            # Terraform infrastructure
├── docker-compose.yml
└── README.md
```

## 🚧 Roadmap

### Phase 1 (Current)
- ✅ Core microservices implementation
- ✅ JWT authentication
- ✅ Database integration
- ✅ Kafka event publishing
- ✅ Docker setup

### Phase 2 (Next)
- [ ] Kubernetes deployment manifests
- [ ] Terraform infrastructure as code
- [ ] CI/CD pipelines
- [ ] Comprehensive testing
- [ ] API documentation (OpenAPI/Swagger)

### Phase 3 (Future)
- [ ] S3/MinIO integration for media
- [ ] Elasticsearch for advanced search
- [ ] Redis caching layer
- [ ] Rate limiting middleware
- [ ] Distributed tracing
- [ ] Monitoring dashboards

## 📄 License

[Your License Here]

## 🤝 Contributing

[Contributing Guidelines]

## 📧 Contact

[Contact Information]
