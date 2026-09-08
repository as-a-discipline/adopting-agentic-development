# `task validate`

The single canonical "is this repository correct right now" gate.

## What it does

Runs, in order:

```
structure:validate → api:validate → service:validate → web:validate
  → policies:validate → compose:validate → helm:validate
```

It stops at the first failure and does not continue past a broken boundary.
It also does not reimplement any component's own logic — it only composes
and delegates to each component's `Taskfile.yml`, per
[`/adr/0001-component-boundaries-in-the-monorepo.md`](../../adr/0001-component-boundaries-in-the-monorepo.md).

## Why it exists

This is the same command a contributor runs before pushing, and the same
one this repository's own documentation points to as proof of engineering
soundness — every guardrail (contract linting/compatibility, generated-code
freshness, service/web build+test, policy enforcement, integration and
deployment validation) is reachable from this one entry point. See
[`/docs/api-breaking-changes.md`](../api-breaking-changes.md) for a real,
captured example of this task failing on a genuine breaking API change.

## How to use it

```bash
task validate
```

No variables. Takes a few minutes on a first run (it builds containers and
starts a Compose stack for `compose:validate`); subsequent runs are faster
due to Docker layer caching.

## Related

* [Task Reference index](../tasks-reference.md)
* [`task structure:validate`](./structure-validate.md) — first stage
* [`task api:validate`](./api-validate.md), [`task service:validate`](./service-validate.md),
  [`task web:validate`](./web-validate.md), [`task policies:validate`](./policies-validate.md),
  [`task compose:validate`](./compose-validate.md), [`task helm:validate`](./helm-validate.md)
