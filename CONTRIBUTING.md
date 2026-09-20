# Contributing

## Development Workflow

Use a short-lived branch from `main`, keep commits focused, and open a PR using the repository template. Use Conventional Commit-style titles such as `feat:`, `fix:`, `docs:`, `test:`, or `chore:`.

For every feature update, review the complete change surface and update affected skills, agent instructions, documentation, tests, README, this guide, configuration, API contracts, migrations, workflows, and deployment files. Record any intentionally unchanged area as not applicable in the PR.

## AI-Assisted Development

AI agents must follow:

- [AGENTS.md](AGENTS.md)
- [.github/copilot-instructions.md](.github/copilot-instructions.md)
- reusable workflows under `.github/skills/`

In VS Code, repository instructions are loaded automatically by supported Copilot features. Ask the agent to use a skill by naming its directory and provide the relevant issue or file scope. Review generated changes as if they were human-authored.

## Verification

Before opening a PR, run:

```bash
make test
make lint
docker compose --env-file .env.example config --quiet
```

Document any unavailable dependency, skipped check, or post-merge action in the PR.

## Issue Linking

Use `Closes #123` or `Fixes #123` in the PR description when the PR fully resolves an issue. Use `Related #123` or `Depends on #123` for non-closing relationships.
