# Copilot Instructions

## Project

Map Navigation Engine is a service-oriented map platform with Go routing and tile services, a Python/FastAPI geocoding service, PostgreSQL/PostGIS, Elasticsearch, and Redis.

## Repository Layout

- `services/routing-service`: Go routing service and graph experiments
- `services/tile-service`: Go vector tile service
- `services/geocoding-service`: Python/FastAPI geocoding service
- `tests/e2e`: Playwright E2E browser and API test suite
- `config`: database and search configuration
- `scripts`: local ingestion, benchmark, and test runner scripts
- `docs/adr`: architecture decisions
- `.github/skills`: reusable AI workflows

## Engineering Rules

- Follow trunk-based development with short-lived branches and small PRs.
- Preserve existing public APIs unless the issue requires a contract change.
- Write or update focused unit tests for production behavior changes, plus integration, end-to-end, or performance tests when the change requires them.
- For every feature update, review and update all affected skills, agent instructions, documentation, tests, README, CONTRIBUTING.md, configuration, API contracts, and other necessary files in the same change.
- Prefer existing project patterns and standard libraries over new abstractions.
- Do not add secrets, credentials, tokens, local paths, or machine-specific URLs.
- Keep configuration in environment variables and update `.env.example` when needed.
- Use versioned database migrations for schema changes; do not edit deployed schemas manually.
- Use ASCII by default and add comments only when they clarify non-obvious logic.
- Include the approved project copyright header in newly authored files when repository policy requires it; never invent legal text. Do not rewrite existing or third-party headers.
- Avoid deprecated libraries and patterns; justify any exception in the PR description.
- Use the latest stable releases of reusable CI actions, verified from each action's official release page before updating workflows.
- Pin every workflow `uses:` reference to the full 40-character commit SHA for the verified release and retain a version comment, for example `actions/checkout@<full-sha> # v7.0.1`. Never use shortened SHAs, mutable tags, or branches for reusable CI actions.

## Feature Update Checklist

Every feature change must assess and update, when applicable:

- `.github/skills/` and agent instructions
- documentation and architecture decision records
- unit, integration, end-to-end, and performance tests
- `README.md` and `CONTRIBUTING.md`
- configuration, migrations, API contracts, workflows, and deployment files

Do not mark a feature complete until the affected supporting files are updated or the PR records why they are not applicable.

## Required Verification

Before proposing a PR, run the narrowest relevant checks and report the commands and results:

- Go: `cd services/routing-service && go test ./... && go vet ./...`
- Go tiles: `cd services/tile-service && go test ./... && go vet ./...`
- Python: `python -m compileall services/geocoding-service/app`
- Compose: `docker compose --env-file .env.example config --quiet`
- Repository checks: `make test`, `make lint`, and `make scan-sast`

Do not claim a check passed without fresh output. If a dependency is unavailable, report the blocker explicitly.

## Change Boundaries

- Do not modify unrelated user changes.
- Do not commit or push unless explicitly requested.
- Do not use destructive Git commands such as `reset --hard` or force-push.
- Ask before changing public API contracts, data retention, security policy, or deployment behavior.
