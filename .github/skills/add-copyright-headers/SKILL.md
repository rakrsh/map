---
name: add-copyright-headers
description: Audit files for the repository copyright policy and add a header only when explicitly required.
---

# Add Copyright Headers

## Workflow

1. Identify the requested scope and file types.
2. Inspect existing headers and repository policy.
3. Do not invent an owner, year, license, or legal text.
4. Ask for the approved copyright text if it is not already defined.
5. Apply the smallest change and preserve generated or third-party files.
6. Run the relevant formatter and tests.
7. Report files changed and any files intentionally skipped.

## Required Context

- `LICENSE`
- `AGENTS.md`
- `.github/copilot-instructions.md`
- The issue or PR containing the explicit header requirement

## Safety

Never add legal headers based only on a guess. Never modify vendored, generated, or third-party files without explicit scope.
