# 1. No Persistence in Baseline

* Status: accepted
* Date: 2026-09-07

## Context and Problem Statement

Pulse checks a small set of HTTP services and reports their current status. It would
be easy to "helpfully" add a database to store check history, but the baseline spec
explicitly defines Pulse as reporting current state only, with no history.

## Decision Drivers

* Keep the baseline application boring and small — the repository/engineering system
  is the demonstration, not the business logic.
* Avoid the operational and architectural surface area a database adds (migrations,
  connection handling, backup/restore, additional deploy dependencies).
* A restart resetting runtime state is an acceptable, even desirable, simplicity.

## Considered Options

1. Add a database (e.g. SQLite/Postgres) to persist check history.
2. Keep monitored-service configuration and last-check results entirely in memory,
   seeded deterministically at startup.

## Decision Outcome

Chosen option: **2 — no persistence.**

* Monitored services are deterministic, seeded/configured at startup (suitable for
  Docker Compose integration tests).
* `GET /api/v1/services` and `GET /api/v1/services/{id}` reflect only the current
  in-memory state.
* `POST /api/v1/services/{id}/check` performs an immediate check and returns the
  resulting status; it does not write to any persistent store.
* A process restart may reset current runtime state — this is expected, not a bug.

## Consequences

* Good: no database dependency, no migrations, no persistence-related deploy
  complexity.
* Good: integration tests are simpler — state is fully determined by process
  lifecycle and configuration.
* Bad: no historical trend data is available; this is an explicit, accepted scope
  limit for the baseline. Future agents must not introduce persistence unless a
  requirement explicitly supersedes this decision (and does so via a new ADR).
