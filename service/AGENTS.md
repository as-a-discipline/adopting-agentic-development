# AGENTS.md — service/

Implements the OpenAPI contract in Go.

## Rules

* Generated interfaces/types under `generated/` are not handwritten — implement the
  generated server interface from `internal/api/`, don't edit generated files.
* HTTP handlers stay thin; business behavior belongs in the appropriate `internal/`
  package (e.g. `internal/monitoring`).
* No arbitrary execution of operating-system commands from application code.
* No persistence — see [`adr/0001-no-persistence-in-baseline.md`](./adr/0001-no-persistence-in-baseline.md).
* Any outbound HTTP call (health checks against monitored services) must use an
  explicit timeout and honor request context.
* Use `task service:*` targets — do not invoke `go` directly outside of them.

## Architecture Decisions

See [`adr/`](./adr/) for service-specific decisions. Cross-cutting decisions are in
the root [`/adr/`](../adr/).

## Task Interface

```text
task service:generate
task service:fmt
task service:lint
task service:test
task service:build
task service:validate
```
