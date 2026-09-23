# Makefile for local development and service orchestration
# Copyright (c) 2026 Ravi Sharma

.PHONY: dev-up dev-down logs build test test-unit test-coverage test-integration test-e2e test-all fmt lint scan-sast scan-secrets scan-dependencies scan-licenses scan-dast scan-security validate-config benchmark-postgis benchmark-clickhouse benchmark-scylladb benchmark-h3 benchmark-graph

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

test-unit:
	@echo "Running unit tests in parallel with isolated doubles..."
	@cd services/routing-service && go test -v -p 4 ./internal/...
	@cd services/tile-service && go test -v -p 4 ./internal/...
	@cd services/geocoding-service && python -m pytest tests/unit -v -n auto

test-coverage:
	@echo "Enforcing code coverage thresholds (>=80%)..."
	@python scripts/check_coverage.py --threshold 80.0

test-integration:
	@echo "Running integration tests with ephemeral containers..."
	@python scripts/run_integration_tests.py

test-e2e:
	@echo "Running Playwright E2E automation with artifact capture..."
	@python scripts/run_e2e_tests.py

test-all: test-unit test-coverage test-integration test-e2e

test: test-unit test-coverage

fmt:
	@echo "Formatting Go and Python sources..."
	@gofmt -w services/routing-service/**/*.go 2>/dev/null || true
	@gofmt -w services/tile-service/**/*.go 2>/dev/null || true
	@cd services/geocoding-service && python -m compileall app

.PHONY: web-dev
web-dev:
	@echo "Opening dev web app (web/index.html). Ensure tile-service is running on :8082"
	@python -m http.server --directory web 3000

lint:
	@echo "Running static checks..."
	@$(MAKE) validate-config
	@cd services/routing-service && go vet ./...
	@cd services/tile-service && go vet ./...
	@cd services/geocoding-service && python -m compileall app
	@test -z "$$(gofmt -l services/routing-service services/tile-service)"

scan-sast:
	@command -v semgrep >/dev/null 2>&1 || (echo "semgrep is required; install it with: python -m pip install semgrep" >&2; exit 1)
	@semgrep scan --config auto --error .

scan-secrets:
	@command -v gitleaks >/dev/null 2>&1 || (echo "gitleaks is required; install it from https://github.com/gitleaks/gitleaks" >&2; exit 1)
	@gitleaks detect --source . --redact --no-banner

scan-dependencies:
	@docker run --rm -v "$$(pwd):/repo" aquasec/trivy:0.70.0 fs --scanners vuln,misconfig,secret --severity HIGH,CRITICAL --ignore-unfixed --exit-code 1 /repo

scan-licenses:
	@docker run --rm -v "$$(pwd):/repo" aquasec/trivy:0.70.0 fs --scanners license --severity HIGH,CRITICAL --exit-code 1 /repo

scan-dast:
	@docker run --rm --network host -t ghcr.io/zaproxy/zaproxy:stable zap-baseline.py -t "$${TARGET:-http://localhost:8000/health}" -I

scan-security: scan-sast scan-secrets scan-dependencies scan-licenses

validate-config:
	@docker compose --env-file .env.example config --quiet
	@docker compose -f docker-compose.test.yml config --quiet

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
