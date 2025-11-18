.PHONY: help build run test test-cover clean deps lint docker-compose migrate

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application
	@echo "Building..."
	@go build -o bin/api cmd/api/main.go

run: ## Run the application
	@echo "Running..."
	@go run cmd/api/main.go

test: ## Run tests
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

test-cover: test ## Run tests with coverage report
	@go tool cover -func=coverage.out

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

lint: ## Run linter
	@echo "Running linter..."
	@echo "If linter was not found, install it on your device: https://golangci-lint.run/docs/welcome/install/"
	@golangci-lint run

docker-compose: ## Build Docker image
	@echo "Building Docker compose..."
	@docker-compose up --build -d

migrate: ## Run database migrations
	@echo "Running migrations..."
	@go run cmd/api/main.go migrate
