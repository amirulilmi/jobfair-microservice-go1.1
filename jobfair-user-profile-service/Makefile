.PHONY: help run build test clean migrate-up migrate-down migrate-create

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

run: ## Run the application
	@echo "🚀 Starting User Profile Service..."
	@go run cmd/main.go

build: ## Build the application
	@echo "🔨 Building User Profile Service..."
	@go build -o bin/user-profile-service cmd/main.go
	@echo "✅ Build completed: bin/user-profile-service"

test: ## Run tests
	@echo "🧪 Running tests..."
	@go test -v ./...

clean: ## Clean build artifacts
	@echo "🧹 Cleaning..."
	@rm -rf bin/
	@rm -rf uploads/
	@echo "✅ Cleaned"

deps: ## Install dependencies
	@echo "📦 Installing dependencies..."
	@go mod download
	@go mod tidy
	@echo "✅ Dependencies installed"

dev: ## Run in development mode with auto-reload (requires air)
	@echo "🔥 Starting development server with auto-reload..."
	@air

install-air: ## Install air for hot reload
	@echo "📦 Installing air..."
	@go install github.com/cosmtrek/air@latest
	@echo "✅ Air installed"

docker-build: ## Build Docker image
	@echo "🐳 Building Docker image..."
	@docker build -t jobfair-user-profile-service:latest .

docker-run: ## Run Docker container
	@echo "🐳 Running Docker container..."
	@docker run -p 8083:8083 --env-file .env jobfair-user-profile-service:latest

# Database migration commands (requires golang-migrate)
migrate-install: ## Install golang-migrate tool
	@echo "📦 Installing golang-migrate..."
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@echo "✅ golang-migrate installed"

migrate-up: ## Run database migrations
	@echo "🔼 Running migrations..."
	@migrate -path migrations -database "$(shell grep DATABASE_URL .env | cut -d '=' -f2)" up
	@echo "✅ Migrations completed"

migrate-down: ## Rollback database migrations
	@echo "🔽 Rolling back migrations..."
	@migrate -path migrations -database "$(shell grep DATABASE_URL .env | cut -d '=' -f2)" down
	@echo "✅ Rollback completed"

migrate-force: ## Force migration version (use: make migrate-force VERSION=1)
	@echo "⚠️  Forcing migration version $(VERSION)..."
	@migrate -path migrations -database "$(shell grep DATABASE_URL .env | cut -d '=' -f2)" force $(VERSION)

migrate-create: ## Create new migration (use: make migrate-create NAME=migration_name)
	@echo "📝 Creating new migration: $(NAME)..."
	@migrate create -ext sql -dir migrations -seq $(NAME)
	@echo "✅ Migration files created"

# Database management
db-create: ## Create database
	@echo "🗄️  Creating database..."
	@psql -U postgres -c "CREATE DATABASE jobfair_user_profile;"
	@echo "✅ Database created"

db-drop: ## Drop database
	@echo "⚠️  Dropping database..."
	@psql -U postgres -c "DROP DATABASE IF EXISTS jobfair_user_profile;"
	@echo "✅ Database dropped"

db-reset: db-drop db-create migrate-up ## Reset database (drop, create, migrate)
	@echo "✅ Database reset completed"

# Helper commands
check-env: ## Check if .env file exists
	@if [ ! -f .env ]; then \
		echo "❌ .env file not found! Copying from .env.example..."; \
		cp .env.example .env; \
		echo "✅ .env file created. Please update with your settings."; \
	else \
		echo "✅ .env file exists"; \
	fi

setup: check-env deps migrate-install ## Setup project (check env, install deps, install migrate tool)
	@echo "✅ Setup completed! You can now run: make migrate-up && make run"

