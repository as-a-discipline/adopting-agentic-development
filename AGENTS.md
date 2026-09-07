# AGENTS.md

Root engineering instructions for the Pulse repository. Read this before making any
change. Read the nearest component `AGENTS.md` before touching a specific component —
do not load every component's instructions just because they exist.

## Repository Architecture

This is a monorepo, but each top-level component is **independently governed**:

* `api/` — owns the public OpenAPI contract. This is the single source of truth for
  the HTTP interface.
* `service/` — implements the contract in Go. Consumes generated server
  interfaces/types; never redefines the contract itself.
* `web/` — consumes only the generated public API client. Never assumes anything
  about Go internals.
* `deploy/` — owns local/runtime packaging (Docker Compose, Helm). Consumes built
  artifacts and configuration; does not own application behavior.
* `policies/` — deterministic engineering constraints applied across components.

Repository locality does not grant dependency access: `web` may not reach into
`service` internals, `service` may not reach into `web` internals, and `deploy` does
not implement business logic.

## Execution Rule

Once Task is available, all project actions (build, test, lint, generate, validate)
must go through documented `task` targets. Agents must not invent alternate shell
commands unless modifying the Task abstraction is itself the assigned work.

`task validate` is the canonical repository-level validation entry point.

## Generated-Code Rule

Files under a component's `generated/` directory are produced by tooling and must
never be hand-edited. To change generated output, change:

* the source contract (`api/openapi.yaml`),
* generator configuration,
* explicitly-allowed templates, or
* the handwritten implementation that consumes the generated interface —

then regenerate via the relevant `task *:generate` target.

## Architecture Decision Records

* Cross-cutting decisions that affect more than one component live in root `/adr/`.
* Component-specific decisions live in that component's own `adr/` folder
  (e.g. `api/adr/`, `service/adr/`). Check the local component's `adr/` folder in
  addition to the root before assuming a decision hasn't been made.
* ADRs use the [MADR](https://adr.github.io/madr/) format.

## Forbidden Actions

Agents may not:

* change CI/CD configuration
* deploy to any environment
* disable or weaken tests
* weaken policy checks
* bypass API compatibility validation
* introduce persistence (no database)
* add authentication
* access unrelated repository areas without cause

## Context Rule

Read the nearest applicable `AGENTS.md` before changing a component. Root
instructions contain global invariants; component instructions contain local
engineering knowledge; skills (`.agents/skills/`) contain workflows loaded only
when relevant to the task at hand.

## Completion Rule

A change is complete only when:

* the intended requirement is satisfied,
* unrelated behavior was not materially changed,
* the component's own validation passes, and
* required integration validation passes.
