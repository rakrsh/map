---
name: generate-unit-tests
description: Create focused unit tests for changed behavior with edge cases and repository-native runners.
---

# Generate Unit Tests

## Workflow

1. Read the implementation, neighboring tests, and issue acceptance criteria.
2. State the behavior hypothesis and the smallest discriminating test.
3. Cover the happy path, boundary conditions, invalid input, and error behavior.
4. Prefer real domain logic and small deterministic fixtures over broad mocks.
5. Use the existing test framework and naming conventions.
6. Run the narrow test first, then the relevant package or service suite.
7. Report commands, results, and remaining coverage gaps.

## Required Context

- The changed implementation and its callers
- Existing tests in the same package
- `AGENTS.md`
- `.github/copilot-instructions.md`
- The issue acceptance criteria

## Quality Rules

Tests must fail for the relevant defect before the fix when practical, must not add test-only production APIs, and must avoid assertions that only verify mock wiring.
