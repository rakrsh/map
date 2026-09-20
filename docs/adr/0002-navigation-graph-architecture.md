# ADR 0002: Navigation Graph Architecture

- Status: Accepted for the dynamic-traffic prototype
- Date: 2026-09-20
- Decision owners: Routing engineering
- Related issue: #4

## Context

The routing service must update edge weights when live traffic speeds change. The architecture spike compares a Contraction Hierarchy (CH) model with a Customizable Contraction Hierarchy (CCH) model against the issue target of updating live weights in under 100 ms.

## Benchmark protocol

Run the deterministic benchmark from the routing service module:

```bash
cd services/routing-service
go run ./cmd/graph-benchmark -edges 100000 -updates 1000
```

The benchmark builds a synthetic directed graph, applies a fixed set of live edge-weight changes, measures update duration, and reports estimated graph storage. The implementation is an architecture model for update behavior; it is not a complete shortest-path or contraction implementation.

## Results

Verified locally with 100,000 edges and 1,000 live updates:

| Architecture | Build time | Update time | Under 100 ms | Estimated graph memory |
| --- | ---: | ---: | :---: | ---: |
| CH | 16.439 ms | 0.547 ms | Yes | 5,600,000 bytes |
| CCH | 1.696 ms | <0.001 ms | Yes | 800,000 bytes |

The values are workload and machine dependent. Repeat the benchmark with representative road-network topology and larger update batches before production capacity planning.

## Decision

Select **CCH** for the dynamic-traffic prototype. Both models meet the sub-100 ms update target in this benchmark, but CCH has the lower estimated graph memory footprint and the lower update cost because live weights are customized without refreshing synthetic shortcut dependencies.

## Consequences

Positive:

- live traffic updates can be applied without rebuilding the graph topology
- the prototype has a measurable sub-100 ms update gate
- lower estimated memory use leaves more room for route-query workloads

Trade-offs:

- CCH requires a separate customization phase and careful topology preprocessing
- the benchmark does not yet measure route-query latency or turn-by-turn reconstruction
- production validation still requires a real OSM-derived graph and representative traffic updates

## Follow-up work

1. Build the graph from the OSM road schema and preserve stable edge IDs.
2. Implement real CCH preprocessing and customization using the selected routing library or a reviewed in-house implementation.
3. Add route-query benchmarks and correctness checks against known routes.
4. Repeat update and memory tests with real regional extracts and concurrent route requests.
