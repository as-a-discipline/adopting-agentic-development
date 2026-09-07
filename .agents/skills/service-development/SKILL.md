---
name: service-development
description: Use when implementing or fixing the Go Pulse service (service/) — handlers, monitoring behavior, or regenerating server scaffolding from the API contract. Trigger on requests to change service behavior, fix a failing service test, or add a Go-side feature that already has an API contract.
---

# Service Development Skill

## When to use

* The API contract already defines the operation; a Go implementation is needed or
  needs to change.
* A `service:test`, `service:lint`, or `service:build` failure needs a root-cause
  fix.
* Service-side monitoring/status-normalization behavior needs to change.

## Workflow

1. Read `service/AGENTS.md` and `service/adr/` before editing.
2. If the API contract changed, run `task service:generate` first — do not hand-edit
   `service/generated/`.
3. Modify only handwritten code (`service/internal/...`, `service/cmd/...`) to
   implement the generated interface.
4. Run targeted tests for the area you changed before the full suite:
   `task service:test`.
5. Run `task service:lint` and `task service:fmt`.
6. Run `task service:validate` to confirm the component as a whole is valid.
7. Interpret failures (read the actual error) before changing code — do not
   iterate blindly.

## Evidence of success

* `task service:validate` exits 0.
* New/changed behavior has corresponding unit tests, including at least one
  negative case.
