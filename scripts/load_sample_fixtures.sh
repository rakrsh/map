#!/usr/bin/env bash
set -euo pipefail

# Load sample fixtures into the configured database. Uses DATABASE_URL or defaults to local postgres.
: ${DATABASE_URL:=postgres://postgres:postgres@127.0.0.1:5432/postgres}

if [ ! -f db/fixtures/sample_pois.sql ]; then
  echo "No fixture found at db/fixtures/sample_pois.sql" >&2
  exit 1
fi

echo "Loading sample fixtures into ${DATABASE_URL}"
psql "${DATABASE_URL}" -f db/fixtures/sample_pois.sql
echo "Fixtures loaded."
