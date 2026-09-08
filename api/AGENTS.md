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

See [`adr/`](./adr/) for API-specific decisions (e.g. contract-first
generation, the plugin-type discovery endpoint in
[`adr/0002-plugin-type-discovery-endpoint.md`](./adr/0002-plugin-type-discovery-endpoint.md)).
Cross-cutting decisions are in the root [`/adr/`](../adr/).

## Plugin-type discovery (`GET /plugin-types`)

`service`'s check mechanisms are implemented as plugins (see
[`service/adr/0002-factory-based-plugin-model-for-checks.md`](../service/adr/0002-factory-based-plugin-model-for-checks.md)).
The API surfaces this as a `type` field on `MonitoredService` (which plugin
implements a given service's check) and a `GET /plugin-types` endpoint that
lists every registered plugin's `type`, `inputSchema`, and `outputSchema` —
so plugin capabilities are discoverable through the API contract itself, not
just in Go source. When adding a new plugin type in `service/`, no contract
change is required unless the new type needs its own request/response
shape beyond the generic `{type, inputSchema, outputSchema}` descriptor.

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
pinned containers, no host install required. It is a real step in
`task api:validate` / `task validate` (in addition to, not a replacement
for, `task api:compat`'s static-baseline check) and is also run directly by
the local pre-commit hook — a breaking change blocks both. See
[`/docs/api-breaking-changes.md`](../docs/api-breaking-changes.md) for
details and captured proof.
