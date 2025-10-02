.PHONY: help build run test clean migrate-up migrate-down docker-build docker-run

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application
	go build -o bin/job-service cmd/main.go

run: ## Run the application
	go run cmd/main.go

test: ## Run tests
	go test -v ./...

clean: ## Clean build files
	rm -rf bin/

migrate-up: ## Run database migrations up
	migrate -path ./migrations -database "postgres://jobfair_user:jobfair_pass@localhost:5435/jobfair_jobs?sslmode=disable" up

migrate-down: ## Run database migrations down
	migrate -path ./migrations -database "postgres://jobfair_user:jobfair_pass@localhost:5435/jobfair_jobs?sslmode=disable" down

migrate-create: ## Create a new migration file (usage: make migrate-create name=create_table_name)
	migrate create -ext sql -dir ./migrations -seq $(name)

docker-build: ## Build docker image
	docker build -t jobfair-job-service:latest .

docker-run: ## Run docker container
	docker run -p 8082:8082 --env-file .env jobfair-job-service:latest

deps: ## Download dependencies
	go mod download
	go mod tidy

lint: ## Run linter
	golangci-lint run

dev: ## Run in development mode with hot reload (requires air)
	air
