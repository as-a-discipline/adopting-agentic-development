# AGENTS.md — web/

Consumes the generated public API client. Does not implement business logic beyond
presenting monitored-service status.

## Rules

* Consume the generated API client under `generated/`; do not hand-write duplicate
  models of API types that already exist as generated types.
* No assumptions about Go service internals — the only contract is the generated
  client derived from `api/openapi.yaml`.
* Keep the UI deliberately small: header, summary line, service list/cards, and a
  per-service check action.
* Basic accessibility and responsive behavior are required; a large design system
  is not.
* Use `task web:*` targets — do not invoke `npm`/frontend tooling directly outside
  of them.

## Task Interface

```text
task web:generate
task web:lint
task web:test
task web:build
task web:validate
```
