# 2. Plugin-Type Discovery Endpoint

* Status: accepted
* Date: 2026-09-08

## Context and Problem Statement

[`service/adr/0002-factory-based-plugin-model-for-checks.md`](../../service/adr/0002-factory-based-plugin-model-for-checks.md)
introduces a factory-based plugin model on the service side: a monitored
service's health check is now implemented by one of several possible
check-plugin types, each with its own self-describing, JSON-Schema-validated
input and output shape. This has two direct consequences for the public API
contract:

1. A `MonitoredService` now has a check **type**, not just a URL — clients
   (in particular the web app) may want to know which kind of check produced
   a given result.
2. Plugin capabilities (which types exist, and what configuration/result
   shape each expects) now exist mechanically on the service side (as Go
   code + embedded JSON Schemas), but nothing about them was previously
   visible through the API at all.

The question this ADR resolves: **how, if at all, should plugin-type
information be exposed through the API contract?**

## Decision Drivers

* Consumers should be able to discover what check-plugin types exist and
  what their input/output shapes are without reading Go source.
* The contract change should be additive (non-breaking) — this repository's
  `task api:breaking-changes` and `task api:compat` gates must continue to
  pass, per [`docs/api-breaking-changes.md`](../../docs/api-breaking-changes.md).
* Avoid over-fitting the `MonitoredService` resource shape to internal
  plugin-configuration details that a client has no need to see or edit
  (Pulse has no service-creation/editing API — services are seeded
  server-side).

## Considered Options

1. **Type discriminator + embedded schema fields on `MonitoredService`**:
   add `type`, `inputSchema`, and `outputSchema` directly to every service
   resource. Rejected — repeats the same (potentially large) schema
   documents on every list item, and conflates "this service's current
   status" with "this plugin type's general capabilities", which are
   logically distinct and change independently.
2. **A schema-discovery endpoint** (`GET /plugin-types`) listing every
   registered plugin type with its input and output JSON Schemas, plus a
   `type` discriminator field on `MonitoredService` pointing at one of those
   types.
3. **No API change at all** — keep plugin structure as a purely internal,
   service-side implementation detail. Rejected — this defeats the stated
   goal of the underlying decision (self-describable, inspectable check
   behavior); a client would have no way to know what `type` values even
   mean or what a plugin expects/returns.

## Decision Outcome

Chosen option: **2 — schema-discovery endpoint plus a `type` discriminator.**

* `MonitoredService` gains a required `type` field (e.g. `"http"`) — the
  discriminator into `GET /plugin-types`, present in every service response.
  This is an additive change to a response schema (a new server-provided
  field on data the server already fully owns); existing consumers reading
  only fields they know about are unaffected.
* A new endpoint, `GET /plugin-types`, returns every registered plugin's
  `type`, `inputSchema`, and `outputSchema` — the schemas are represented
  generically (OpenAPI `type: object` with `additionalProperties: true`)
  since actual JSON Schema keywords vary per plugin and this contract should
  not need to change every time a plugin's schema does.
* This keeps per-service responses small (just a type name) while making
  full plugin capability information available in exactly one place,
  fetched once and cached by a client rather than repeated on every service.
* Verified via `task api:breaking-changes` (comparing against the prior
  contract) and `task api:compat` (against the stored baseline) that both
  changes are additive/non-breaking.

## Consequences

* Good: a client (e.g. the web app) can resolve `service.type` against
  `GET /plugin-types` to show what kind of check produced a result, and
  even render/validate configuration using the published JSON Schemas, all
  without any Go-source knowledge of the service implementation.
* Good: additive-only change — no existing field or endpoint's meaning
  changes, confirmed by the breaking-change detection tasks.
* Neutral: `GET /plugin-types` is a capability/metadata endpoint, not a
  per-service resource — it does not follow the `/services/{id}` pattern
  because it describes plugin *types*, which are process-wide, not
  per-service.
* Future agents adding a new check-plugin type on the service side (per the
  companion service ADR) must ensure it becomes visible through
  `GET /plugin-types` automatically (by registering it with the plugin
  registry) — no separate manual step should be required to keep this
  endpoint accurate.
