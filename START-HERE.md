# 🚀 START HERE - How to Run

## Simple 3-Step Guide

### 1️⃣ Start Infrastructure (PostgreSQL, Redis, Kafka)

```bash
docker-compose up -d
```

Wait 30-60 seconds for services to start.

### 2️⃣ Install Dependencies

```bash
make install-deps
```

### 3️⃣ Run All Services

Open **5 separate terminal windows** and run:

**Terminal 1:**
```bash
cd services/auth-service
go run cmd/api/main.go
```

**Terminal 2:**
```bash
cd services/user-service
go run cmd/api/main.go
```

**Terminal 3:**
```bash
cd services/product-service
go run cmd/api/main.go
```

**Terminal 4:**
```bash
cd services/task-service
go run cmd/api/main.go
```

**Terminal 5:**
```bash
cd services/media-service
go run cmd/api/main.go
```

## ✅ Verify It's Working

Open a new terminal and test:

```bash
# Health check
curl http://localhost:3001/health

# Register a user
curl -X POST http://localhost:3001/api/v1/register -H "Content-Type: application/json" -d "{\"email\":\"test@example.com\",\"username\":\"testuser\",\"password\":\"password123\"}"
```

## 📝 Quick Reference

| Service | Port | URL |
|---------|------|-----|
| Auth | 3001 | http://localhost:3001 |
| User | 3002 | http://localhost:3002 |
| Product | 3003 | http://localhost:3003 |
| Task | 3004 | http://localhost:3004 |
| Media | 3005 | http://localhost:3005 |

## 🛑 Stop Everything

Press `Ctrl+C` in each service terminal, then:

```bash
docker-compose down
```

## 📚 More Help

- **Detailed guide**: See [RUN.md](RUN.md)
- **API examples**: See [QUICKSTART.md](QUICKSTART.md)
- **Full docs**: See [README.md](README.md)

## ⚠️ Troubleshooting

**Port in use?**
- Windows: `netstat -ano | findstr :3001` then `taskkill /PID <PID> /F`
- Linux/Mac: `lsof -ti:3001 | xargs kill -9`

**Database error?**
- Check: `docker-compose ps`
- Logs: `docker-compose logs db`

**Service won't start?**
- Make sure `.env` file exists in root
- Check Go version: `go version` (needs 1.21+)
- Verify infrastructure is running: `docker-compose ps`

