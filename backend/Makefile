.DEFAULT_GOAL := help
.PHONY: help run dev build tidy test fmt vet

help: ## List all available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-8s\033[0m %s\n", $$1, $$2}'

run: ## Run the API
	go run ./cmd/api

dev: ## Run the API with hot reload (requires air)
	air

build: ## Build the API binary to ./bin/api
	go build -o ./bin/api ./cmd/api

tidy: ## Sync module dependencies
	go mod tidy

test: ## Run all tests
	go test ./...

fmt: ## Format all Go source
	go fmt ./...

vet: ## Run go vet
	go vet ./...
