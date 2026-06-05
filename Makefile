.PHONY: run dev build test tidy docker-up docker-down swag

# ── Development ────────────────────────────────────────────────────────────────

## run: Run the app directly
run:
	go run ./cmd/api

## dev: Run with hot reload (requires: go install github.com/air-verse/air@latest)
dev:
	air

## build: Build production binary
build:
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o bin/portfolio-cms ./cmd/api

## tidy: Tidy go modules
tidy:
	go mod tidy

## test: Run all tests
test:
	go test ./... -v -race -coverprofile=coverage.out

## test-cover: Show test coverage in browser
test-cover: test
	go tool cover -html=coverage.out

# ── Swagger ────────────────────────────────────────────────────────────────────

## swag: Generate Swagger docs (requires: go install github.com/swaggo/swag/cmd/swag@latest)
swag:
	swag init -g cmd/api/main.go -o docs

# ── Docker ─────────────────────────────────────────────────────────────────────

## docker-up: Start all services with Docker Compose
docker-up:
	docker compose up -d --build

## docker-down: Stop all services
docker-down:
	docker compose down

## docker-logs: Follow app logs
docker-logs:
	docker compose logs -f app

# ── Database ───────────────────────────────────────────────────────────────────

## db-reset: Drop and recreate database (development only!)
db-reset:
	docker compose exec postgres psql -U postgres -c "DROP DATABASE IF EXISTS portfolio_cms;"
	docker compose exec postgres psql -U postgres -c "CREATE DATABASE portfolio_cms;"

# ── Helpers ────────────────────────────────────────────────────────────────────

## help: Show this help message
help:
	@echo "Available commands:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'
