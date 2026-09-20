#!/usr/bin/env bash
# Copyright (c) 2026 Ravi Sharma

set -euo pipefail

CLICKHOUSE_URL="${CLICKHOUSE_URL:-http://localhost:8123}"
SCALE="${BENCHMARK_SCALE:-100000}"
OUTPUT_DIR="${BENCHMARK_OUTPUT_DIR:-benchmarks/results}"
mkdir -p "$OUTPUT_DIR"
REPORT="$OUTPUT_DIR/clickhouse-$(date -u +%Y%m%dT%H%M%SZ).txt"

if ! curl --fail --silent "$CLICKHOUSE_URL/ping" >/dev/null; then
  echo "ClickHouse is unavailable at $CLICKHOUSE_URL. Start the benchmark adapter first." >&2
  exit 1
fi

run_query() {
  curl --fail --silent --data-binary "$1" "$CLICKHOUSE_URL/" 
}

{
  echo "ClickHouse benchmark"
  echo "scale=$SCALE"
  run_query "CREATE DATABASE IF NOT EXISTS benchmark"
  run_query "DROP TABLE IF EXISTS benchmark.points"
  run_query "CREATE TABLE benchmark.points (id UInt64, lon Float64, lat Float64) ENGINE = MergeTree ORDER BY id"
  run_query "INSERT INTO benchmark.points SELECT number, (randUniform(-180, 180)), (randUniform(-90, 90)) FROM numbers($SCALE)"
  echo "bulk_rows=$SCALE"
  for query in \
    "SELECT count() FROM benchmark.points WHERE greatCircleDistance(lon, lat, 0, 0) < 100000" \
    "SELECT count() FROM benchmark.points WHERE lon BETWEEN -10 AND 10 AND lat BETWEEN -10 AND 10"; do
    echo "query=$query"
    run_query "$query"
  done
} | tee "$REPORT"

echo "Benchmark report written to $REPORT"
