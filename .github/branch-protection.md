# Main Branch Protection Checklist

GitHub repository administrators must configure these settings for `main`; they cannot be enforced by repository files alone.

- Require a pull request before merging.
- Require at least one approving review.
- Require approval of the most recent push.
- Require the `Verify` status check from `Trunk CI`.
- Require the `SAST` status check from `Trunk CI`.
- Require the AI governance validation check.
- Require branches to be up to date before merging.
- Require linear history.
- Allow squash merge and rebase merge; disable merge commits.
- Dismiss stale approvals when new commits are pushed.
- Restrict direct pushes, force pushes, and branch deletion, including for administrators where organizational policy permits.
- Enable automatic head-branch deletion after merge.
- Use a ruleset or branch protection rule targeting `main`.

## Feature Flags

Incomplete work should be merged behind a configuration-controlled feature flag. Flags must default to the safe or disabled state, avoid secrets, and include an owner and removal issue. A flag should be removed when the feature is fully released.

## Branch Lifecycle

Use short-lived branches named `feat/<scope>`, `fix/<scope>`, `docs/<scope>`, or `chore/<scope>`. Keep branches synchronized with `main` and aim to merge within 24-48 hours.
