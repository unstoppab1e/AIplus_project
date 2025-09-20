.PHONY: build run test test-coverage clean docker-build docker-run docker-stop docker-logs dev-db dev-db-stop deps lint help

# Переменные
APP_NAME=Aiplus_project
DOCKER_IMAGE=employee-service:latest
DATABASE_URL=postgres://postgres:password@localhost:5432/employees?sslmode=disable

# Установка зависимостей
deps:
	@echo "📦 Installing dependencies..."
	go mod download
	go mod tidy
	@echo "✅ Dependencies installed"

# Сборка приложения
build:
	@echo "🔨 Building application..."
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/$(APP_NAME) cmd/server/main.go
	@echo "✅ Build completed: bin/$(APP_NAME)"

# Запуск приложения локально
run:
	@echo "🚀 Starting application locally..."
	go run cmd/server/main.go

# Запуск unit тестов
test:
	@echo "🧪 Running unit tests..."
	go test ./tests/... -v -race

# Запуск тестов с покрытием
test-coverage:
	@echo "📊 Running tests with coverage..."
	go test ./tests/... -v -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"
	@go tool cover -func=coverage.out | grep total

# Линтер (требует golangci-lint)
lint:
	@echo "🔍 Running linter..."
	golangci-lint run

# Очистка артефактов
clean:
	@echo "🧹 Cleaning up..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	@echo "✅ Cleanup completed"

# Docker команды
docker-build:
	@echo "🐳 Building Docker image..."
	docker build -t $(DOCKER_IMAGE) .
	@echo "✅ Docker image built: $(DOCKER_IMAGE)"

docker-run:
	@echo "🚀 Starting services with Docker Compose..."
	docker-compose up --build

docker-stop:
	@echo "🛑 Stopping Docker services..."
	docker-compose down

docker-logs:
	@echo "📋 Showing application logs..."
	docker-compose logs -f app

# Тестирование в Docker
docker-test:
	@echo "🧪 Running Docker integration tests..."
	@if [ -f test-docker.sh ]; then \
		chmod +x test-docker.sh && ./test-docker.sh; \
	else \
		echo "❌ test-docker.sh not found"; \
	fi

# Локальная БД для разработки
dev-db:
	@echo "🗄️ Starting PostgreSQL for development..."
	docker run --name employee-postgres \
		-e POSTGRES_DB=employees \
		-e POSTGRES_USER=postgres \
		-e POSTGRES_PASSWORD=password \
		-p 5432:5432 \
		-d postgres:15
	@echo "✅ PostgreSQL started on port 5432"

dev-db-stop:
	@echo "🛑 Stopping development database..."
	docker stop employee-postgres && docker rm employee-postgres
	@echo "✅ Development database stopped"

# Полная очистка Docker
docker-clean:
	@echo "🧹 Cleaning Docker resources..."
	docker-compose down -v
	docker system prune -f
	@echo "✅ Docker cleanup completed"

# Быстрые тесты
quick-test:
	@echo "⚡ Running quick tests..."
	go test ./tests/... -short

# Проверка готовности к коммиту
check: deps lint test
	@echo "✅ All checks passed - ready to commit!"

# Полный цикл разработки
dev: deps docker-run

# Production сборка
production-build:
	@echo "🏭 Building for production..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
		-ldflags="-w -s" \
		-a -installsuffix cgo \
		-o bin/$(APP_NAME) cmd/server/main.go
	@echo "✅ Production build completed"

# Benchmark тесты
benchmark:
	@echo "⚡ Running benchmarks..."
	go test ./tests/... -bench=. -benchmem

# Показать статус
status:
	@echo "📊 Project Status:"
	@echo "Go version: $(shell go version)"
	@echo "Docker: $(shell docker --version 2>/dev/null || echo 'not installed')"
	@echo "Docker Compose: $(shell docker-compose --version 2>/dev/null || echo 'not installed')"
	@echo "Project: $(APP_NAME)"
	@echo "Database URL: $(DATABASE_URL)"

# Помощь
help:
	@echo "🚀 Employee Service - Available Commands:"
	@echo ""
	@echo "📦 Development:"
	@echo "  deps            - Install Go dependencies"
	@echo "  build           - Build the application"
	@echo "  run             - Run application locally"
	@echo "  dev             - Full development setup"
	@echo ""
	@echo "🧪 Testing:"
	@echo "  test            - Run unit tests"
	@echo "  test-coverage   - Run tests with coverage report"
	@echo "  quick-test      - Run short tests only"
	@echo "  docker-test     - Run Docker integration tests"
	@echo "  benchmark       - Run performance benchmarks"
	@echo ""
	@echo "🐳 Docker:"
	@echo "  docker-build    - Build Docker image"
	@echo "  docker-run      - Start services with Docker Compose"
	@echo "  docker-stop     - Stop Docker services"
	@echo "  docker-logs     - Show application logs"
	@echo "  docker-clean    - Clean Docker resources"
	@echo ""
	@echo "🗄️ Database:"
	@echo "  dev-db          - Start local PostgreSQL"
	@echo "  dev-db-stop     - Stop local PostgreSQL"
	@echo ""
	@echo "🔧 Utilities:"
	@echo "  lint            - Run code linter"
	@echo "  clean           - Clean build artifacts"
	@echo "  check           - Run all checks (lint + test)"
	@echo "  status          - Show project status"
	@echo "  production-build - Build optimized binary"
	@echo ""
	@echo "💡 Quick start: make dev"