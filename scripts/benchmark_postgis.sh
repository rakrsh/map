#!/usr/bin/env bash
# Copyright (c) 2026 Ravi Sharma

set -euo pipefail

CONTAINER_NAME="${POSTGIS_CONTAINER:-map-postgres}"
DATABASE="${POSTGIS_DATABASE:-mapdb}"
USER_NAME="${POSTGIS_USER:-mapuser}"
SCALE="${BENCHMARK_SCALE:-100000}"
CLIENTS="${BENCHMARK_CLIENTS:-8}"
DURATION="${BENCHMARK_DURATION:-30}"
OUTPUT_DIR="${BENCHMARK_OUTPUT_DIR:-benchmarks/results}"
PYTHON_BIN="${PYTHON_BIN:-python3}"

mkdir -p "$OUTPUT_DIR"
TIMESTAMP="$(date -u +%Y%m%dT%H%M%SZ)"
REPORT="$OUTPUT_DIR/postgis-$TIMESTAMP.txt"
TEMP_DIR="$(mktemp -d)"

cleanup() {
  psql_exec -c "DROP SCHEMA IF EXISTS benchmark CASCADE;" >/dev/null 2>&1 || true
  rm -rf "$TEMP_DIR"
}
trap cleanup EXIT

psql_exec() {
  docker exec "$CONTAINER_NAME" psql -v ON_ERROR_STOP=1 -U "$USER_NAME" -d "$DATABASE" "$@"
}

if ! docker inspect "$CONTAINER_NAME" >/dev/null 2>&1; then
  echo "PostGIS container '$CONTAINER_NAME' is not available. Start it with: docker compose up -d postgres" >&2
  exit 1
fi

{
  echo "PostGIS benchmark"
  echo "timestamp_utc=$TIMESTAMP"
  echo "database=$DATABASE"
  echo "scale=$SCALE"
  echo "clients=$CLIENTS"
  echo "duration_seconds=$DURATION"
  echo
  echo "## PostGIS version"
  psql_exec -Atc "SELECT postgis_full_version();"

  echo
  echo "## Bulk insert setup and throughput"
  psql_exec -c "DROP SCHEMA IF EXISTS benchmark CASCADE; CREATE SCHEMA benchmark; CREATE TABLE benchmark.points (id BIGINT PRIMARY KEY, geom geometry(Point, 4326) NOT NULL); CREATE TABLE benchmark.roads (id BIGINT PRIMARY KEY, geom geometry(LineString, 4326) NOT NULL, highway TEXT NOT NULL); CREATE TABLE benchmark.pois (id BIGINT PRIMARY KEY, geom geometry(Point, 4326) NOT NULL, category TEXT NOT NULL); CREATE TABLE benchmark.areas (id BIGINT PRIMARY KEY, geom geometry(Polygon, 4326) NOT NULL);"
  start_ms="$(date +%s%3N)"
  psql_exec -c "INSERT INTO benchmark.points (id, geom) SELECT id, ST_SetSRID(ST_MakePoint(-180 + random() * 360, -90 + random() * 180), 4326) FROM generate_series(1, $SCALE) AS id; INSERT INTO benchmark.roads (id, geom, highway) SELECT id, ST_SetSRID(ST_MakeLine(ST_MakePoint(-180 + random() * 360, -90 + random() * 180), ST_MakePoint(-180 + random() * 360, -90 + random() * 180)), 4326), CASE WHEN id % 4 = 0 THEN 'primary' ELSE 'residential' END FROM generate_series(1, $SCALE) AS id; INSERT INTO benchmark.pois (id, geom, category) SELECT id, ST_SetSRID(ST_MakePoint(-180 + random() * 360, -90 + random() * 180), 4326), CASE WHEN id % 3 = 0 THEN 'restaurant' ELSE 'landmark' END FROM generate_series(1, $SCALE) AS id; INSERT INTO benchmark.areas (id, geom) SELECT id, ST_SetSRID(ST_MakeEnvelope(-180 + random() * 350, -90 + random() * 170, -180 + random() * 350 + 10, -90 + random() * 170 + 10), 4326) FROM generate_series(1, greatest(1, $SCALE / 100)) AS id;"
  end_ms="$(date +%s%3N)"
  elapsed_ms="$((end_ms - start_ms))"
  awk -v rows="$((SCALE * 3 + (SCALE / 100)))" -v elapsed="$elapsed_ms" 'BEGIN { printf "rows_inserted=%d\nelapsed_seconds=%.3f\nrows_per_second=%.2f\n", rows, elapsed / 1000, rows / (elapsed / 1000) }'
  psql_exec -c "CREATE INDEX benchmark_points_geom_gist ON benchmark.points USING GIST (geom); CREATE INDEX benchmark_roads_geom_gist ON benchmark.roads USING GIST (geom); CREATE INDEX benchmark_pois_geom_gist ON benchmark.pois USING GIST (geom); CREATE INDEX benchmark_areas_geom_gist ON benchmark.areas USING GIST (geom); ANALYZE benchmark.points; ANALYZE benchmark.roads; ANALYZE benchmark.pois; ANALYZE benchmark.areas;"

  echo
  echo "## Spatial operator plans"
  echo "### ST_DWithin"
  psql_exec -P pager=off -c "EXPLAIN (ANALYZE, BUFFERS, FORMAT TEXT) SELECT count(*) FROM benchmark.points WHERE ST_DWithin(geom, ST_SetSRID(ST_MakePoint(0, 0), 4326), 1);"
  echo "### ST_Intersects"
  psql_exec -P pager=off -c "EXPLAIN (ANALYZE, BUFFERS, FORMAT TEXT) SELECT count(*) FROM benchmark.points WHERE ST_Intersects(geom, ST_MakeEnvelope(-10, -10, 10, 10, 4326));"
  echo "### Bounding box"
  psql_exec -P pager=off -c "EXPLAIN (ANALYZE, BUFFERS, FORMAT TEXT) SELECT count(*) FROM benchmark.points WHERE geom && ST_MakeEnvelope(-10, -10, 10, 10, 4326);"

  echo
  echo "## Concurrent query benchmark"
  psql_exec -c "CREATE TABLE benchmark.query (id integer PRIMARY KEY, geom geometry(Point, 4326) NOT NULL); INSERT INTO benchmark.query VALUES (1, ST_SetSRID(ST_MakePoint(0, 0), 4326));"
  for operator in dwithin intersects bbox; do
    case "$operator" in
      dwithin) predicate="ST_DWithin(points.geom, query.geom, 1)" ;;
      intersects) predicate="ST_Intersects(points.geom, ST_MakeEnvelope(-10, -10, 10, 10, 4326))" ;;
      bbox) predicate="points.geom && ST_MakeEnvelope(-10, -10, 10, 10, 4326)" ;;
    esac
    cat <<SQL | docker exec -i "$CONTAINER_NAME" psql -v ON_ERROR_STOP=1 -U "$USER_NAME" -d "$DATABASE" >/dev/null
CREATE OR REPLACE FUNCTION benchmark.$operator() RETURNS void LANGUAGE plpgsql AS \$\$
BEGIN
  PERFORM 1 FROM benchmark.points, benchmark.query WHERE $predicate;
END;
\$\$;
SQL
    printf 'SELECT benchmark.%s();\n' "$operator" > "$TEMP_DIR/pgbench-$operator.sql"
    docker cp "$TEMP_DIR/pgbench-$operator.sql" "$CONTAINER_NAME:/tmp/pgbench-$operator.sql"
    log_prefix="/tmp/pgbench-$operator-$TIMESTAMP"
    docker exec "$CONTAINER_NAME" pgbench -n -M prepared -U "$USER_NAME" -c "$CLIENTS" -T "$DURATION" -l --log-prefix="$log_prefix" -f "/tmp/pgbench-$operator.sql" "$DATABASE"
    docker exec "$CONTAINER_NAME" sh -c "cat ${log_prefix}.*" > "$TEMP_DIR/pgbench-$operator.log"
    echo "### $operator latency percentiles"
    "$PYTHON_BIN" scripts/benchmark_metrics.py "$TEMP_DIR/pgbench-$operator.log"
  done

  echo
  echo "## Memory-sensitive join plan"
  echo "### PostgreSQL memory before join"
  docker stats --no-stream --format "container_memory_usage={{.MemUsage}}" "$CONTAINER_NAME"
  echo "### Join plan and buffer consumption"
  psql_exec -P pager=off -c "EXPLAIN (ANALYZE, BUFFERS, FORMAT TEXT) SELECT count(*) FROM benchmark.points p1 JOIN benchmark.points p2 ON ST_DWithin(p1.geom, p2.geom, 0.01) WHERE p1.id <= 1000;"
  echo "### PostgreSQL memory after join"
  docker stats --no-stream --format "container_memory_usage={{.MemUsage}}" "$CONTAINER_NAME"
} | tee "$REPORT"
echo "Benchmark report written to $REPORT"
