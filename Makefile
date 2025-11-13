.PHONY: help build up down logs clean test migrate

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build all Docker images
	docker build -t auth-service:latest -f services/auth-service/Dockerfile .
	docker build -t user-service:latest -f services/user-service/Dockerfile .
	docker build -t product-service:latest -f services/product-service/Dockerfile .
	docker build -t task-service:latest -f services/task-service/Dockerfile .
	docker build -t media-service:latest -f services/media-service/Dockerfile .

up: ## Start all infrastructure services
	docker-compose up -d

down: ## Stop all infrastructure services
	docker-compose down

logs: ## Show logs from all services
	docker-compose logs -f

clean: ## Remove all containers and volumes
	docker-compose down -v
	docker system prune -f

test: ## Run tests for all services
	@echo "Running tests..."
	@cd services/auth-service && go test ./... || true
	@cd services/user-service && go test ./... || true
	@cd services/product-service && go test ./... || true
	@cd services/task-service && go test ./... || true
	@cd services/media-service && go test ./... || true

migrate: ## Run database migrations (auto-migrates on service start)
	@echo "Migrations run automatically on service startup"

run-auth: ## Run auth service locally
	cd services/auth-service && go run cmd/api/main.go

run-user: ## Run user service locally
	cd services/user-service && go run cmd/api/main.go

run-product: ## Run product service locally
	cd services/product-service && go run cmd/api/main.go

run-task: ## Run task service locally
	cd services/task-service && go run cmd/api/main.go

run-media: ## Run media service locally
	cd services/media-service && go run cmd/api/main.go

install-deps: ## Install Go dependencies for all services
	@echo "Installing dependencies..."
	@cd pkg/config && go mod download || true
	@cd pkg/database && go mod download || true
	@cd pkg/jwtutils && go mod download || true
	@cd pkg/kafkaclient && go mod download || true
	@cd pkg/logger && go mod download || true
	@cd services/auth-service && go mod download || true
	@cd services/user-service && go mod download || true
	@cd services/product-service && go mod download || true
	@cd services/task-service && go mod download || true
	@cd services/media-service && go mod download || true

