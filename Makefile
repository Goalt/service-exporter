# Makefile for service-exporter Go project

# Variables
BINARY_NAME=service-exporter
BUILD_DIR=.
COVERAGE_FILE=coverage.out
COVERAGE_HTML=coverage.html

# Default target
.PHONY: help
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: lint
lint: ## Run linting (go vet and go fmt check)
	@echo "Running go vet..."
	go vet ./...
	@echo "Checking go fmt..."
	@if [ -n "$$(go fmt ./...)" ]; then \
		echo "Code is not properly formatted. Run 'make fmt' to fix."; \
		exit 1; \
	else \
		echo "Code is properly formatted."; \
	fi

.PHONY: fmt
fmt: ## Format Go code
	@echo "Formatting Go code..."
	go fmt ./...

.PHONY: test
test: ## Run tests with coverage
	@echo "Running tests with coverage..."
	go test -cover ./...

.PHONY: test-coverage
test-coverage: ## Run tests and generate HTML coverage report
	@echo "Running tests and generating coverage report..."
	go test -coverprofile=$(COVERAGE_FILE) ./...
	go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "Coverage report generated: $(COVERAGE_HTML)"

.PHONY: test-verbose
test-verbose: ## Run tests with verbose output and coverage
	@echo "Running tests with verbose output and coverage..."
	go test -v -cover ./...

.PHONY: build
build: ## Build the application
	@echo "Building $(BINARY_NAME)..."
	go build -o $(BINARY_NAME) ./cmd

.PHONY: run
run: build ## Build and run the application
	@echo "Running $(BINARY_NAME)..."
	./$(BINARY_NAME)

.PHONY: clean
clean: ## Clean build artifacts and coverage files
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)
	rm -f $(COVERAGE_FILE)
	rm -f $(COVERAGE_HTML)

.PHONY: deps
deps: ## Download and tidy dependencies
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

.PHONY: all
all: deps lint test build ## Run all targets (deps, lint, test, build)

.PHONY: ci
ci: deps lint test ## Run CI pipeline (deps, lint, test)