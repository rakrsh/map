# Contributing

## Development Workflow

Use a short-lived branch from `main`, keep commits focused, and open a PR using the repository template. Use Conventional Commit-style titles such as `feat:`, `fix:`, `docs:`, `test:`, or `chore:`.

## Trunk-Based Development

`main` is the single source of truth. Work on short-lived branches named `feat/<scope>`, `fix/<scope>`, `docs/<scope>`, or `chore/<scope>` and keep them synchronized with `main`. Aim to merge focused changes within 24-48 hours using squash or rebase merge; merge commits are disabled by repository policy.

All PRs must pass Trunk CI, AI governance validation, required reviews, and the configured security gates before merge. PR size labels are advisory; `size/XL` is reserved for changes over 1,500 lines and receives a warning, but is not blocked.

Incomplete work must be protected by a disabled-by-default feature flag with an owner and removal issue. Do not use feature flags for secrets or security controls.

Repository administrators must configure the settings in [.github/branch-protection.md](.github/branch-protection.md).

For every feature update, review the complete change surface and update affected skills, agent instructions, documentation, tests, README, this guide, configuration, API contracts, migrations, workflows, and deployment files. Record any intentionally unchanged area as not applicable in the PR.

## AI-Assisted Development

AI agents must follow:

- [AGENTS.md](AGENTS.md)
- [.github/copilot-instructions.md](.github/copilot-instructions.md)
- reusable workflows under `.github/skills/`

In VS Code, repository instructions are loaded automatically by supported Copilot features. Ask the agent to use a skill by naming its directory and provide the relevant issue or file scope. Review generated changes as if they were human-authored.

To activate Agent Mode in VS Code, open Copilot Chat, choose **Agent** from the mode selector, and provide the issue number, requested scope, and expected verification. In other supported IDE or CLI integrations, select the equivalent agent workflow and apply the same repository instructions.

## Verification

Before opening a PR, run:

```bash
make test
make lint
make scan-sast
docker compose --env-file .env.example config --quiet
```

Document any unavailable dependency, skipped check, or post-merge action in the PR.

## Issue Linking

Use `Closes #123` or `Fixes #123` in the PR description when the PR fully resolves an issue. Use `Related #123` or `Depends on #123` for non-closing relationships.
