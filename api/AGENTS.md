# AGENTS.md — api/

Owns the public OpenAPI contract (`openapi.yaml`). This is the single source of
truth for the Pulse HTTP interface.

## Rules

* The OpenAPI contract is the source of truth. `service` and `web` consume
  generated artifacts derived from it; they do not define their own view of the API
  shape.
* Maintain backward compatibility unless a requirement explicitly permits a
  breaking change.
* Lint before generating: `task api:lint` must pass before `task api:generate`.
* Generated consumers (`service/generated/`, `web/generated/`) are downstream
  artifacts of this contract — do not hand-edit them from here or elsewhere.
* Examples and schemas in the contract must remain valid against the spec.
* Do not change `service` or `web` implementation code as part of API work; change
  the contract, then let the owning component regenerate and adapt.

## Architecture Decisions

See [`adr/`](./adr/) for API-specific decisions (e.g. contract-first generation).
Cross-cutting decisions are in the root [`/adr/`](../adr/).

## Task Interface

```text
task api:lint
task api:validate
task api:generate
```
