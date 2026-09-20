# Map Navigation Engine

A modern map and navigation platform focused on geocoding, routing, and vector map delivery. The repository is structured as a multi-service system that reflects the product brief and technical statement of work for a scalable navigation application.

## Overview

The system is designed to support:
- accurate address search and reverse geocoding
- multi-modal route planning for cars, cyclists, pedestrians, and transit
- vector tile rendering for fast map visualization
- spatial storage with PostGIS
- search indexing with Elasticsearch
- caching and transient state management with Redis
- low-stress and custom route preferences aligned to the user experience requirements

## Architecture

```text
map/
├── config/
│   ├── elasticsearch/
│   │   └── mappings.json
│   └── postgres/
│       └── init.sql
├── scripts/
│   └── osm_ingest.sh
├── services/
│   ├── geocoding-service/
│   │   ├── app/
│   │   ├── Dockerfile
│   │   └── pyproject.toml
│   ├── routing-service/
│   │   ├── cmd/
│   │   ├── internal/
│   │   ├── Dockerfile
│   │   └── go.mod
│   └── tile-service/
│       ├── cmd/
│       ├── Dockerfile
│       └── go.mod
├── .env.example
├── .gitignore
├── docker-compose.yml
├── Makefile
├── README.md
├── project_bootstrap.md
├── sow.txt
├── technical_details.txt
└── LICENSE
```

## Services

### 1) Routing Service
Location: `services/routing-service`

Technology: Go

Responsibilities:
- route generation and shortest-path calculations
- low-stress routing heuristics
- route preference handling for tolls, highways, dirt roads, and similar constraints
- health and route API endpoints

### 2) Geocoding Service
Location: `services/geocoding-service`

Technology: Python + FastAPI

Responsibilities:
- address autocomplete and search
- reverse geocoding
- geospatial indexing with Elasticsearch
- requests for location lookup and place metadata

### 3) Tile Service
Location: `services/tile-service`

Technology: Go

Responsibilities:
- serving map tiles for frontend rendering
- vector tile endpoints
- map-layer integration support

## Data Layer

- PostgreSQL + PostGIS for persistent geospatial data
- Elasticsearch for search, autocomplete, and place documents
- Redis for ephemeral caching and coordination

## Prerequisites

Before running the project locally, install:
- Docker
- Docker Compose
- Go 1.22+
- Python 3.11+
- Git

## Quick Start

1. Copy the environment template:

```bash
cp .env.example .env
```

2. Start the local stack:

```bash
docker compose up --build
```

3. Verify the services:

- Geocoding API: http://localhost:8000/health
- Routing Service: http://localhost:8081/health
- Tile Service: http://localhost:8082/health
- PostgreSQL: localhost:5432
- Elasticsearch: http://localhost:9200
- Redis: localhost:6379

The backing services expose Docker healthchecks. Use `docker compose ps` and wait for `healthy` before testing dependent services.

## Environment Variables

The project uses values defined in `.env.example`:

```env
POSTGRES_DB=mapdb
POSTGRES_USER=mapuser
POSTGRES_PASSWORD=map_password
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
REDIS_ADDR=localhost:6379
ELASTICSEARCH_URL=http://localhost:9200
ROUTING_SERVICE_PORT=8081
GEOCODING_SERVICE_PORT=8000
TILE_SERVICE_PORT=8082
```

## Local Development Commands

From the repository root:

```bash
# Start the full stack
docker compose up --build

# Stop the stack
docker compose down -v

# View logs
docker compose logs -f

# Run the Go service directly
cd services/routing-service
go run ./cmd/server

# Run the Python service directly
cd services/geocoding-service
python -m uvicorn app.main:app --host 0.0.0.0 --port 8000 --reload
```

## Geospatial Benchmarking

Issue #3 includes a reproducible PostGIS baseline and optional comparison adapters:

```bash
# Start PostGIS
docker compose up -d postgres

# Run the PostGIS workload and write a local report
make benchmark-postgis

# Optional comparison adapters
make benchmark-clickhouse
make benchmark-scylladb
make benchmark-h3
```

The PostGIS benchmark generates OSM-shaped roads, POIs, and polygons, then reports bulk throughput, spatial query plans, concurrent p50/p95/p99 latency, and container memory snapshots. Reports are written to `benchmarks/results/` and are ignored by Git. See `docs/adr/0001-geospatial-database-engine-benchmark.md` for the protocol and current baseline.

Issue #4 includes a CH/CCH navigation graph update benchmark:

```bash
make benchmark-graph
```

The current prototype selects CCH for dynamic traffic updates. See `docs/adr/0002-navigation-graph-architecture.md` for the measured comparison and production follow-up work.

## Issue and PR Linking

Use recognized GitHub closing keywords in the PR description:

```text
Closes #123
Fixes #456
```

When the PR is merged into the repository's default branch, GitHub closes the referenced issues automatically. `Related #123` and `Depends on #456` create links but do not close issues. Project-board status changes require GitHub Projects workflow automation to be configured in the project settings; they are not controlled by repository files.

## Current Status

This repository is in its bootstrap stage. The structure is set up to support the planned navigation application, but the core business logic, spatial indexes, and API contracts are still being implemented.

## Planned Next Steps

1. define concrete routing API payloads and response models
2. implement geocoding/search indexing pipelines
3. wire PostGIS and Elasticsearch integration for place and route data
4. build a client frontend or map interface
5. add automated tests and CI checks

## References

- `project_bootstrap.md` — repository architecture and scaffold guidance
- `sow.txt` — product and delivery statement of work
- `technical_details.txt` — deeper technical requirements and constraints
