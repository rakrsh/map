# Makefile for local development and service orchestration

.PHONY: dev-up dev-down logs build test fmt lint

dev-up:
	docker compose up --build -d

dev-down:
	docker compose down -v

up: dev-up
down: dev-down

logs:
	docker compose logs -f

build:
	docker compose build

test:
	@echo "Running Go tests..."
	@cd services/routing-service && go test ./...
	@cd services/tile-service && go test ./...
	@cd services/geocoding-service && python -m compileall app

fmt:
	@echo "Formatting Go and Python sources..."
	@gofmt -w services/routing-service/**/*.go 2>/dev/null || true
	@gofmt -w services/tile-service/**/*.go 2>/dev/null || true
	@cd services/geocoding-service && python -m compileall app

lint:
	@echo "Running static checks..."
	@cd services/routing-service && go vet ./...
	@cd services/tile-service && go vet ./...
	@cd services/geocoding-service && python -m compileall app
