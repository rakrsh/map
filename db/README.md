Database migrations
-------------------

This folder contains SQL migrations used to version the PostgreSQL schema.

We use the `golang-migrate` CLI (via Docker image `ghcr.io/golang-migrate/migrate`) to apply migrations.

Common commands (from project root):

```sh
# Apply all pending migrations (uses DATABASE_URL env var or local postgres default)
make db-migrate

# Rollback last migration
make db-rollback

# Create new migration stub (provide NAME)
make db-create NAME=add_new_table
```

Guidelines:
- Keep migrations additive and backward-compatible where possible.
- Provide both `.up.sql` and `.down.sql` when creating non-trivial changes.
- Follow the `V<timestamp>__description.sql` naming convention used by the tooling.
