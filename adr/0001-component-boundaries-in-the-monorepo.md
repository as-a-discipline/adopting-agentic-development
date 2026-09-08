# 1. Component Boundaries in the Monorepo

* Status: accepted
* Date: 2026-09-07

## Context and Problem Statement

Pulse is organized as a single repository containing `api`, `service`, `web`, and
`deploy`. A shared repository makes it easy — and tempting — for one component to
reach directly into another component's internals (e.g. `web` importing a Go
package, or `service` depending on frontend build output). We need a clear rule for
what repository locality does and does not grant.

## Decision Drivers

* The repository must demonstrate independently governed components despite shared
  storage.
* Agents (and humans) need an unambiguous rule to check changes against.
* The OpenAPI contract must remain the only interface between `api`/`service` and
  `web`.

## Considered Options

1. Treat the repo as one codebase with no enforced boundaries.
2. Split into separate repositories (polyrepo).
3. Keep a single repository, but treat each top-level component as independently
   governed, communicating only through defined contracts/artifacts.

## Decision Outcome

Chosen option: **3 — single repository, independently governed components.**

* `web` communicates with `service` only through the generated client for the public
  OpenAPI contract. It must never assume Go implementation details.
* `service` implements the OpenAPI contract; it does not depend on anything in `web`.
* `deploy` consumes built artifacts and configuration from `api`/`service`/`web`; it
  does not implement application behavior.
* Internal packages of `service` (e.g. `service/internal/...`) are not shared with
  `web` or `deploy`.

This preserves most of the ergonomic benefits of a monorepo (atomic commits across
components, shared CI-equivalent validation) while keeping the components as
replaceable as if they were separate repositories.

## Consequences

* Good: cross-component dependencies are limited to the contract, making it possible
  to reason about and later split components without a rewrite.
* Good: a "no cross-component source import" check is mechanically checkable and can
  be added to `policies/`.
* Bad: some duplication (e.g. re-fetching config) may occur instead of directly
  importing a sibling component's code; this is accepted as the cost of the boundary.
