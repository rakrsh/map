# ADR 0001: Geospatial Database Engine Benchmark

- Status: Proposed
- Date: 2026-09-20
- Decision owners: Architecture team
- Related issue: #3

## Context

The map platform needs a spatial data layer for bulk OpenStreetMap ingestion, spatial predicates, joins, and high-concurrency reads. PostgreSQL/PostGIS is the current local implementation. The alternatives named in issue #3 are ScyllaDB, ClickHouse, and Uber H3, but they do not expose the same storage and query model, so they must be compared through workload adapters rather than a single SQL benchmark.

## Benchmark protocol

The first executable baseline is `scripts/benchmark_postgis.sh`. It runs inside the local PostGIS container and records:

- bulk point insert elapsed time and approximate throughput
- `ST_DWithin`, `ST_Intersects`, and bounding-box query plans
- concurrent spatial query throughput and latency using `pgbench`
- buffer usage and execution details for a spatial self-join
- the PostGIS version and benchmark parameters

Run it with:

```bash
make benchmark-postgis
```

Optional controls:

```bash
BENCHMARK_SCALE=1000000 BENCHMARK_CLIENTS=16 BENCHMARK_DURATION=60 make benchmark-postgis
```

Reports are written to `benchmarks/results/` and should not be committed. Compare engines using the same dataset shape, concurrency, warm/cold-cache policy, hardware limits, and correctness checks.

## Current evidence

No cross-engine measurements are recorded yet. The PostGIS script is the reproducible baseline; the final choice must not be marked accepted until equivalent adapters have been run for the candidate engines.

## Decision

Pending benchmark results. The existing application scaffold may continue using PostGIS for local development, but production engine selection remains provisional until the benchmark table is populated.

## Follow-up work

1. Add a dataset generator for representative OSM roads, POIs, and polygons.
2. Add adapters for ClickHouse, ScyllaDB, and H3 where the workload is semantically equivalent.
3. Capture p50/p95/p99 latency, throughput, RSS/peak memory, and correctness results.
4. Publish the final comparison and update this ADR to `Accepted` or `Rejected`.
