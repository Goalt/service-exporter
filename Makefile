# Makefile for service-exporter Go project

# Variables
BINARY_NAME=service-exporter
BUILD_DIR=.
COVERAGE_FILE=coverage.out
COVERAGE_HTML=coverage.html

# Default target
.PHONY: help
help: ## Show this help message
	@printf '\033[1;36m📖 Service Exporter - Available Commands\033[0m\n'
	@printf '\033[1;34m========================================\033[0m\n'
	@echo ''
	@printf '\033[1;32mUsage:\033[0m make [target]\n'
	@echo ''
	@printf '\033[1;32mTargets:\033[0m\n'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[33m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: lint
lint: ## Run linting (go vet and go fmt check)
	@printf '\033[1;34m🔍 Running go vet...\033[0m\n'
	go vet ./...
	@printf '\033[1;34m🔍 Checking go fmt...\033[0m\n'
	@if [ -n "$$(go fmt ./...)" ]; then \
		printf '\033[1;31m❌ Code is not properly formatted. Run '\''make fmt'\'' to fix.\033[0m\n'; \
		exit 1; \
	else \
		printf '\033[1;32m✅ Code is properly formatted.\033[0m\n'; \
	fi

.PHONY: fmt
fmt: ## Format Go code
	@printf '\033[1;35m✨ Formatting Go code...\033[0m\n'
	go fmt ./...
	@printf '\033[1;32m✅ Code formatting complete!\033[0m\n'

.PHONY: test
test: ## Run tests with coverage
	@printf '\033[1;34m🧪 Running tests with coverage...\033[0m\n'
	go test -cover ./...

.PHONY: test-coverage
test-coverage: ## Run tests and generate HTML coverage report
	@printf '\033[1;34m📊 Running tests and generating coverage report...\033[0m\n'
	go test -coverprofile=$(COVERAGE_FILE) ./...
	go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@printf '\033[1;32m✅ Coverage report generated: $(COVERAGE_HTML)\033[0m\n'

.PHONY: test-verbose
test-verbose: ## Run tests with verbose output and coverage
	@printf '\033[1;34m🔬 Running tests with verbose output and coverage...\033[0m\n'
	go test -v -cover ./...

.PHONY: build
build: ## Build the application
	@printf '\033[1;33m🏗️  Building $(BINARY_NAME)...\033[0m\n'
	go build -o $(BINARY_NAME) ./cmd
	@printf '\033[1;32m✅ Build complete: $(BINARY_NAME)\033[0m\n'

.PHONY: run
run: build ## Build and run the application
	@printf '\033[1;32m🚀 Running $(BINARY_NAME)...\033[0m\n'
	./$(BINARY_NAME)

.PHONY: clean
clean: ## Clean build artifacts and coverage files
	@printf '\033[1;33m🧹 Cleaning up...\033[0m\n'
	rm -f $(BINARY_NAME)
	rm -f $(COVERAGE_FILE)
	rm -f $(COVERAGE_HTML)
	@printf '\033[1;32m✅ Cleanup complete!\033[0m\n'

.PHONY: deps
deps: ## Download and tidy dependencies
	@printf '\033[1;34m📦 Downloading dependencies...\033[0m\n'
	go mod download
	go mod tidy
	@printf '\033[1;32m✅ Dependencies updated successfully!\033[0m\n'

.PHONY: all
all: deps lint test build ## Run all targets (deps, lint, test, build)

.PHONY: ci
ci: deps lint test ## Run CI pipeline (deps, lint, test)