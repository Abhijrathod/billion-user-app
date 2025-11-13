#!/bin/bash

# Setup script for Billion User App
# This script helps set up the development environment

set -e

echo "🚀 Setting up Billion User App..."

# Check if .env exists
if [ ! -f .env ]; then
    echo "📝 Creating .env file from .env.example..."
    cp .env.example .env
    echo "✅ Created .env file. Please update it with your configuration."
else
    echo "✅ .env file already exists"
fi

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker and try again."
    exit 1
fi

echo "🐳 Starting infrastructure services..."
docker-compose up -d

echo "⏳ Waiting for services to be ready..."
sleep 10

# Check if PostgreSQL is ready
until docker exec postgres_db pg_isready -U admin > /dev/null 2>&1; do
    echo "⏳ Waiting for PostgreSQL..."
    sleep 2
done

echo "✅ PostgreSQL is ready!"

# Check if Redis is ready
until docker exec redis_cache redis-cli ping > /dev/null 2>&1; do
    echo "⏳ Waiting for Redis..."
    sleep 2
done

echo "✅ Redis is ready!"

# Check if Kafka is ready (this might take longer)
echo "⏳ Waiting for Kafka (this may take a minute)..."
sleep 30

echo ""
echo "✅ Setup complete!"
echo ""
echo "Next steps:"
echo "1. Update .env file with your configuration if needed"
echo "2. Run services locally:"
echo "   - Auth Service:    make run-auth"
echo "   - User Service:    make run-user"
echo "   - Product Service: make run-product"
echo "   - Task Service:    make run-task"
echo "   - Media Service:   make run-media"
echo ""
echo "Or use Docker Compose to run all services together."
echo ""
echo "To stop infrastructure: docker-compose down"

