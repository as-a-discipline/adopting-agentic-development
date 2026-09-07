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
task api:breaking-changes
```

`task api:breaking-changes` compares the current working tree against the
current branch's last pushed commit (or the target mainline branch if
nothing has been pushed yet) using either `oasdiff` (default,
`TOOL=oasdiff`) or `pb33f/openapi-changes` (`TOOL=openapi-changes`) — both
pinned containers, no host install required. This is a developer/reviewer
tool for catching breaking changes before merge; it is separate from
`task api:compat`'s static-baseline check, which is part of the always-on
`task validate` gate.
