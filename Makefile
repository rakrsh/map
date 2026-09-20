# Makefile for local development and service orchestration

.PHONY: dev-up dev-down logs build test fmt lint validate-config benchmark-postgis benchmark-clickhouse benchmark-scylladb benchmark-h3 benchmark-graph

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
	@$(MAKE) validate-config
	@cd services/routing-service && go vet ./...
	@cd services/tile-service && go vet ./...
	@cd services/geocoding-service && python -m compileall app
	@test -z "$$(gofmt -l services/routing-service services/tile-service)"

validate-config:
	@docker compose --env-file .env.example config --quiet

benchmark-postgis:
	@bash scripts/benchmark_postgis.sh

benchmark-clickhouse:
	@bash scripts/benchmark_clickhouse.sh

benchmark-scylladb:
	@bash scripts/benchmark_scylladb.sh

benchmark-h3:
	@python3 scripts/benchmark_h3.py

benchmark-graph:
	@cd services/routing-service && go run ./cmd/graph-benchmark
