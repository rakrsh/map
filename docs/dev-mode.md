Dev-mode Map (minimal)
----------------------

This document describes the minimal developer workflow to run a dev-mode map locally.

Steps:

1. Start Postgres/PostGIS via Docker compose:

   docker compose up -d postgres

2. Start the tile service (dev server):

   cd services/tile-service
   go run ./cmd/server

3. Serve the frontend and open the map in a browser:

   # from repo root
   make web-dev
   # open http://localhost:3000

4. For a full dev loop, run the geocoding and routing services locally or via `docker compose up`.

Notes:
- The `dev-map` Makefile target prints quick start instructions.
- To enable quick integration tests, add a small fixture dataset under `db/fixtures/` and load it via `psql`.
