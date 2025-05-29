.PHONY: help build run test clean docker-build docker-up docker-down docker-logs migrate

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application
	@echo "Building..."
	go build -o bin/server ./cmd/server/main.go

run: ## Run the application locally
	@echo "Running..."
	go run ./cmd/server/main.go

test: ## Run tests
	@echo "Running tests..."
	go test -v ./...

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -rf bin/
	rm -f server

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker-compose build

docker-up: ## Start Docker containers
	@echo "Starting Docker containers..."
	docker-compose up -d

docker-down: ## Stop Docker containers
	@echo "Stopping Docker containers..."
	docker-compose down

docker-logs: ## Show Docker logs
	docker-compose logs -f

docker-restart: docker-down docker-up ## Restart Docker containers

fmt: ## Format code
	@echo "Formatting code..."
	go fmt ./...

lint: ## Run linter
	@echo "Running linter..."
	golangci-lint run

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

.DEFAULT_GOAL := help
