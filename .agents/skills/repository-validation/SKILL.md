---
name: repository-validation
description: Use when running full repository validation, or when a validation failure needs to be triaged to the right component. Trigger on requests to "validate everything", "check the repo is healthy", or when task validate fails and the cause is unclear.
---

# Repository Validation Skill

## When to use

* A broad "is the repository valid?" check is requested.
* `task validate` failed and it's unclear which boundary (structure/api/service/
  web/compose/helm) caused it.

## Workflow

1. Run `task validate` (or `task structure:validate` alone for just the
   structural check) from the repository root.
2. Read the output to identify which task in the chain failed first.
3. Load only the applicable component's `AGENTS.md` (and matching skill) for the
   failing boundary — do not load every component's context.
4. Fix the root cause in that component.
5. Re-run the scoped task first (e.g. `task service:validate`), then re-run
   `task validate` for the full picture.

## Evidence of success

* `task validate` exits 0.
* The fix touched only the component that actually owned the failure.
