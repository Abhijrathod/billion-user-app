# Quick Start Guide

Get up and running with the Billion User App in 5 minutes!

## Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose
- Make (optional, but recommended)

## Step 1: Clone and Setup

```bash
git clone <repository-url>
cd billion-user-app
```

## Step 2: Start Infrastructure

```bash
# Option A: Use the setup script (Linux/Mac)
chmod +x scripts/setup.sh
./scripts/setup.sh

# Option B: Manual setup
docker-compose up -d
```

This starts:
- PostgreSQL (port 5432)
- Redis (port 6379)
- Kafka + Zookeeper (port 9092)

## Step 3: Configure Environment

Create a `.env` file in the root directory:

```bash
# Copy the example
cp .env.example .env

# Or create manually with these values:
APP_ENV=development
PORT=3000

DB_HOST=localhost
DB_PORT=5432
DB_USER=admin
DB_PASSWORD=secret
DB_NAME=postgres
DB_SSLMODE=disable

REDIS_ADDRESS=localhost:6379
KAFKA_BROKERS=localhost:9092
JWT_SECRET=super-secret-key-change-in-production
```

## Step 4: Install Dependencies

```bash
# Install all Go dependencies
make install-deps

# Or manually for each service:
cd services/auth-service && go mod download
cd ../user-service && go mod download
# ... etc
```

## Step 5: Run Services

Open 5 terminal windows and run each service:

**Terminal 1 - Auth Service:**
```bash
cd services/auth-service
go run cmd/api/main.go
```

**Terminal 2 - User Service:**
```bash
cd services/user-service
go run cmd/api/main.go
```

**Terminal 3 - Product Service:**
```bash
cd services/product-service
go run cmd/api/main.go
```

**Terminal 4 - Task Service:**
```bash
cd services/task-service
go run cmd/api/main.go
```

**Terminal 5 - Media Service:**
```bash
cd services/media-service
go run cmd/api/main.go
```

Or use Make commands:
```bash
make run-auth    # In one terminal
make run-user    # In another terminal
# ... etc
```

## Step 6: Test the API

### Register a User

```bash
curl -X POST http://localhost:3001/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "username": "testuser",
    "password": "password123"
  }'
```

### Login

```bash
curl -X POST http://localhost:3001/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

Save the `access_token` from the response.

### Get User Profile (Protected)

```bash
curl -X GET http://localhost:3002/api/v1/users/1 \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### Create a Product (Protected)

```bash
curl -X POST http://localhost:3003/api/v1/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "name": "Test Product",
    "description": "A test product",
    "price": 99.99,
    "sku": "TEST-001",
    "stock": 100,
    "category": "electronics"
  }'
```

### Create a Task (Protected)

```bash
curl -X POST http://localhost:3004/api/v1/tasks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "title": "Complete project",
    "description": "Finish the billion user app",
    "status": "in_progress",
    "priority": 2
  }'
```

## Health Checks

Check if services are running:

```bash
curl http://localhost:3001/health  # Auth Service
curl http://localhost:3002/health  # User Service
curl http://localhost:3003/health  # Product Service
curl http://localhost:3004/health  # Task Service
curl http://localhost:3005/health  # Media Service
```

## Troubleshooting

### Services won't start

1. **Check if ports are available:**
   ```bash
   # Check if ports are in use
   lsof -i :3001
   lsof -i :5432
   ```

2. **Check Docker containers:**
   ```bash
   docker-compose ps
   docker-compose logs
   ```

3. **Verify database connection:**
   ```bash
   docker exec -it postgres_db psql -U admin -d auth_db -c "SELECT 1;"
   ```

### Database connection errors

- Ensure PostgreSQL is running: `docker-compose ps`
- Check `.env` file has correct credentials
- Verify database was created: `docker exec postgres_db psql -U admin -l`

### Kafka connection errors

- Kafka takes ~30 seconds to start
- Check logs: `docker-compose logs kafka`
- Services will continue without Kafka (events just won't be published)

## Next Steps

- Read [README.md](README.md) for full documentation
- Check [ARCHITECTURE.md](ARCHITECTURE.md) for system design
- Explore the API endpoints in the README
- Set up monitoring and observability
- Deploy to Kubernetes (see `k8s/` directory)

## Development Tips

1. **Hot Reload**: Use [Air](https://github.com/cosmtrek/air) for hot reloading:
   ```bash
   go install github.com/cosmtrek/air@latest
   air
   ```

2. **Database Migrations**: Migrations run automatically on startup. For manual migrations, use GORM's migration tools.

3. **Testing**: Run tests with `make test` or `go test ./...`

4. **Logging**: Logs are structured JSON. Use `jq` to pretty-print:
   ```bash
   docker-compose logs -f | jq
   ```

## Need Help?

- Check the [README.md](README.md) for detailed documentation
- Review [ARCHITECTURE.md](ARCHITECTURE.md) for system design
- Open an issue on GitHub

Happy coding! 🚀

