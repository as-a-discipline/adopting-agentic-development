# 1. Contract-First API

* Status: accepted
* Date: 2026-09-07

## Context and Problem Statement

Pulse exposes an HTTP API consumed by both a Go service (as the implementer) and a
TypeScript web app (as a consumer). If either side hand-writes its own model of the
contract, the two will drift, and the API stops being an actual boundary.

## Decision Drivers

* `service` and `web` must never invent their own, potentially divergent, views of
  the API shape.
* Changes to the contract should be reviewable in one place.
* Generated code must be trivially identifiable and never hand-edited.

## Considered Options

1. Hand-write Go structs and TypeScript interfaces independently in each component.
2. Write the API contract by hand in each component and reconcile manually.
3. Define `api/openapi.yaml` as the single source of truth, and generate both the Go
   server interfaces/types and the TypeScript client from it.

## Decision Outcome

Chosen option: **3 — contract-first with generated consumers.**

* `api/openapi.yaml` is the source of truth for the HTTP contract.
* `service/generated/` contains Go server interfaces/types generated from the
  contract; `service/internal/...` implements those interfaces (handwritten code
  never redefines the contract).
* `web/generated/` contains a TypeScript client generated from the same contract.
* Generated code is never manually edited; changes go through the contract (or
  generator configuration), then regeneration via `task api:generate` /
  `task service:generate` / `task web:generate`.
* API changes must pass lint and a compatibility check before being considered valid.

## Consequences

* Good: one edit to `openapi.yaml`, followed by regeneration, keeps both consumers
  in sync by construction.
* Good: generated-code freshness is mechanically checkable (regenerate and diff).
* Bad: adds a generation step to the development loop; mitigated by containerized,
  pinned generators (see the containerized-tooling ADR) and clear `task *:generate`
  targets.
