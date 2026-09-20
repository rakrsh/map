#!/usr/bin/env bash
# Copyright (c) 2026 Ravi Sharma

set -euo pipefail

SCYLLA_HOST="${SCYLLA_HOST:-localhost}"
SCYLLA_PORT="${SCYLLA_PORT:-9042}"
SCALE="${BENCHMARK_SCALE:-100000}"
OUTPUT_DIR="${BENCHMARK_OUTPUT_DIR:-benchmarks/results}"
mkdir -p "$OUTPUT_DIR"
REPORT="$OUTPUT_DIR/scylladb-$(date -u +%Y%m%dT%H%M%SZ).txt"

if ! command -v cqlsh >/dev/null 2>&1; then
  echo "cqlsh is required for the ScyllaDB adapter." >&2
  exit 1
fi

{
  echo "ScyllaDB benchmark"
  echo "host=$SCYLLA_HOST"
  echo "scale=$SCALE"
  cqlsh "$SCYLLA_HOST" "$SCYLLA_PORT" <<CQL
CREATE KEYSPACE IF NOT EXISTS benchmark WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};
CREATE TABLE IF NOT EXISTS benchmark.points (id bigint PRIMARY KEY, lon double, lat double);
CQL
  echo "The Scylla adapter provisions the workload schema. Use a driver-based load test for bulk writes and token/range queries; native spatial predicates are not equivalent to PostGIS."
} | tee "$REPORT"

echo "Benchmark report written to $REPORT"
