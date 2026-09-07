# AGENTS.md — deploy/

Owns local/runtime packaging: Docker Compose integration environment and the Helm
chart. Consumes artifacts and configuration produced by `api`/`service`/`web`; does
not implement application behavior.

## Subcomponents

* [`compose/`](./compose/AGENTS.md) — deterministic local integration environment.
* [`helm/`](./helm/AGENTS.md) — Kubernetes packaging (no live cluster required for
  validation).

## Task Interface

```text
task compose:up
task compose:down
task compose:test
task compose:validate

task helm:lint
task helm:template
task helm:validate
```
