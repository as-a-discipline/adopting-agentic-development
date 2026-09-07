---
name: api-contract
description: Use when inspecting or changing the Pulse OpenAPI contract (api/openapi.yaml), or when downstream Go/TypeScript artifacts need to be regenerated from it. Trigger on requests to add/change an API operation or schema, or to fix an API lint/compatibility failure.
---

# API Contract Skill

## When to use

* A request changes or adds an HTTP operation, request/response schema, or error
  shape for Pulse.
* An API lint or compatibility check is failing and needs to be fixed at the
  source.
* Downstream generated code (`service/generated/`, `web/generated/`) is stale
  relative to `api/openapi.yaml`.

## Workflow

1. Read `api/AGENTS.md` and `api/adr/` before editing.
2. Edit `api/openapi.yaml` only — never edit generated output directly.
3. Run `task api:lint`. Fix reported issues before proceeding.
4. Run `task api:generate` to regenerate downstream Go/TypeScript artifacts.
5. Run `task api:validate` to confirm lint, compatibility, and generation
   freshness all pass.
6. If compatibility checking flags a breaking change, treat that as a design
   question, not an obstacle to silence.

## Evidence of success

* `task api:validate` exits 0.
* `git status` shows only the contract change plus the corresponding regenerated
  files under `service/generated/` and `web/generated/` — no stray manual edits to
  generated paths.
