.PHONY: test test-coverage test-race build clean fmt vet lint run-collector help

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Test targets
test: ## Run all tests
	go test -v ./...

test-coverage: ## Run tests with coverage
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-race: ## Run tests with race detection
	go test -v -race ./...

test-clean: ## Clean test cache and coverage files
	go clean -testcache
	rm -f coverage.out coverage.html

# Build targets
build: ## Build all binaries
	go build -o bin/collector cmd/collector/main.go

build-all: ## Build all binaries for multiple platforms
	GOOS=linux GOARCH=amd64 go build -o bin/collector-linux-amd64 cmd/collector/main.go
	GOOS=darwin GOARCH=amd64 go build -o bin/collector-darwin-amd64 cmd/collector/main.go
	GOOS=windows GOARCH=amd64 go build -o bin/collector-windows-amd64.exe cmd/collector/main.go

# Code quality targets
fmt: ## Format code
	go fmt ./...

vet: ## Run go vet
	go vet ./...

lint: ## Run linter (requires golangci-lint)
	golangci-lint run

# Runtime targets
run-collector: ## Run the feed collector
	go run cmd/collector/main.go

# Cleanup targets
clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.out coverage.html
	go clean -cache
	go clean -testcache

# Development targets
deps: ## Download dependencies
	go mod download
	go mod verify

deps-update: ## Update dependencies
	go mod tidy
	go get -u ./...

# CI targets
ci-test: test-race test-coverage ## Run CI tests (race detection + coverage)

install-tools: ## Install development tools
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	