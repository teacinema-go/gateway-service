.PHONY: help
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: build
build: ## Build the application
	@echo "Building application..."
	go build -o bin/gateway-service cmd/gateway-service/main.go
	@echo "Build complete: bin/gateway-service"

.PHONY: run
run: ## Run the application
	@echo "Starting application..."
	go run cmd/gateway-service/main.go