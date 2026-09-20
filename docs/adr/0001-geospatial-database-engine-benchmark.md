# ADR 0001: Geospatial Database Engine Benchmark

- Status: Proposed
- Date: 2026-09-20
- Decision owners: Architecture team
- Related issue: #3

## Context

The map platform needs a spatial data layer for bulk OpenStreetMap ingestion, spatial predicates, joins, and high-concurrency reads. PostgreSQL/PostGIS is the current local implementation. The alternatives named in issue #3 are ScyllaDB, ClickHouse, and Uber H3, but they do not expose the same storage and query model, so they must be compared through workload adapters rather than a single SQL benchmark.

## Benchmark protocol

The executable baseline is `scripts/benchmark_postgis.sh`. It runs inside the local PostGIS container and records:

- OSM-shaped roads, POIs, and polygon bulk insert elapsed time and rows per second
- `ST_DWithin`, `ST_Intersects`, and bounding-box query plans
- concurrent spatial query throughput plus average, p50, p95, and p99 latency using `pgbench`
- Docker container memory usage before and after a spatial self-join, plus buffer usage and execution details
- the PostGIS version and benchmark parameters

Run it with:

```bash
make benchmark-postgis
```

Optional controls:

```bash
BENCHMARK_SCALE=1000000 BENCHMARK_CLIENTS=16 BENCHMARK_DURATION=60 make benchmark-postgis
```

Reports are written to `benchmarks/results/` and should not be committed. The ClickHouse, ScyllaDB, and H3 adapters are available through `make benchmark-clickhouse`, `make benchmark-scylladb`, and `make benchmark-h3`. Compare engines using the same dataset shape, concurrency, warm/cold-cache policy, hardware limits, and correctness checks.

## PostGIS baseline evidence

The verified local run used 100,000 rows per point/line workload and 1,000 polygons:

| Measurement | Result |
| --- | ---: |
| Bulk rows | 301,000 |
| Bulk insert throughput | 251,882.85 rows/s |
| `ST_DWithin` p95 / p99 | 6.449 ms / 7.414 ms |
| `ST_Intersects` p95 / p99 | 2.309 ms / 2.795 ms |
| Bounding-box p95 / p99 | 1.074 ms / 1.297 ms |
| Spatial join execution time | 58.107 ms |
| Concurrent query failures | 0 |

The report is intentionally generated locally because benchmark results depend on host resources. Cross-engine measurements are not yet recorded; the final choice must not be marked accepted until equivalent adapters have been run for the candidate engines.

## Decision

Pending benchmark results. The existing application scaffold may continue using PostGIS for local development, but production engine selection remains provisional until the benchmark table is populated.

## Follow-up work

1. Run the ClickHouse adapter against the shared point workload.
2. Run the ScyllaDB adapter with a driver-based load test; native spatial predicates are not equivalent to PostGIS.
3. Run the H3 adapter as an indexing baseline; H3 is not a transactional database.
4. Capture p50/p95/p99 latency, throughput, container memory, and correctness results in a comparison table.
5. Publish the final comparison and update this ADR to `Accepted` or `Rejected`.
