# How to Run the Billion User App

## Quick Start (5 Steps)

### Step 1: Start Infrastructure Services

Open a terminal in the project root and run:

```bash
docker-compose up -d
```

This starts:
- PostgreSQL (port 5432)
- Redis (port 6379)  
- Kafka + Zookeeper (port 9092)

**Wait 30-60 seconds** for all services to be ready.

### Step 2: Create Environment File

Create a `.env` file in the root directory with this content:

```env
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

### Step 3: Install Go Dependencies

```bash
make install-deps
```

Or manually:
```bash
cd services/auth-service && go mod download && cd ../..
cd services/user-service && go mod download && cd ../..
cd services/product-service && go mod download && cd ../..
cd services/task-service && go mod download && cd ../..
cd services/media-service && go mod download && cd ../..
```

### Step 4: Run Services

You need **5 terminal windows** (one for each service):

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

**OR use Make commands** (each in separate terminals):
```bash
make run-auth
make run-user
make run-product
make run-task
make run-media
```

### Step 5: Test It Works

Open a new terminal and test:

```bash
# Check health
curl http://localhost:3001/health

# Register a user
curl -X POST http://localhost:3001/api/v1/register ^
  -H "Content-Type: application/json" ^
  -d "{\"email\":\"test@example.com\",\"username\":\"testuser\",\"password\":\"password123\"}"

# Login (save the access_token from response)
curl -X POST http://localhost:3001/api/v1/login ^
  -H "Content-Type: application/json" ^
  -d "{\"email\":\"test@example.com\",\"password\":\"password123\"}"
```

## Service Ports

- Auth Service: `http://localhost:3001`
- User Service: `http://localhost:3002`
- Product Service: `http://localhost:3003`
- Task Service: `http://localhost:3004`
- Media Service: `http://localhost:3005`

## Troubleshooting

### Port Already in Use

If you get "port already in use" errors:

**Windows:**
```powershell
# Find what's using the port
netstat -ano | findstr :3001

# Kill the process (replace PID with actual process ID)
taskkill /PID <PID> /F
```

**Linux/Mac:**
```bash
lsof -ti:3001 | xargs kill -9
```

### Database Connection Error

1. Check if PostgreSQL is running:
   ```bash
   docker-compose ps
   ```

2. Check PostgreSQL logs:
   ```bash
   docker-compose logs db
   ```

3. Verify databases were created:
   ```bash
   docker exec -it postgres_db psql -U admin -l
   ```

### Service Won't Start

1. Check if `.env` file exists in root directory
2. Verify Go version: `go version` (needs 1.21+)
3. Check service logs for errors
4. Make sure infrastructure is running: `docker-compose ps`

### Kafka Connection Warnings

Kafka takes 30-60 seconds to start. If you see Kafka connection warnings, wait a bit and restart the service. Services will work without Kafka (events just won't be published).

## Stop Everything

```bash
# Stop all services (Ctrl+C in each terminal)
# Then stop infrastructure:
docker-compose down

# Or remove everything including data:
docker-compose down -v
```

## Next Steps

- See [QUICKSTART.md](QUICKSTART.md) for detailed API examples
- Check [README.md](README.md) for full documentation
- Review [ARCHITECTURE.md](ARCHITECTURE.md) for system design

