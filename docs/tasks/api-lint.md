# `task api:lint`

## What it does

Runs Redocly CLI (`redocly/cli:1.25.3`, pinned container) against
[`api/openapi.yaml`](../../api/openapi.yaml).

## Why it exists

The contract must stay structurally valid OpenAPI and pass sensible-default
linting before anything downstream (generation, compatibility checks) is
trusted. See [`policies/README.md`](../../policies/README.md) rule 2.

## How to use it

```bash
task api:lint
```

## Related

* [Task Reference index](../tasks-reference.md)
* [`task api:validate`](./api-validate.md)
* [`task api:generate`](./api-generate.md)
