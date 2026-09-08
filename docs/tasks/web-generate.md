# `task web:generate`

## What it does

Runs `openapitools/openapi-generator-cli:v7.11.0` (pinned container,
`typescript-fetch` generator) against
[`api/openapi.yaml`](../../api/openapi.yaml), writes output to
`web/generated/`, then runs a post-generation fixup script
(`tools/fix-generated-imports.mjs`) that adds explicit `.ts` import
extensions, required by Node's strict ESM resolver.

## Why it exists

`web/generated/` must never be hand-written, same rule as
`service/generated/` — this is the one mechanism that produces it.

## How to use it

```bash
task web:generate
```

## Related

* [Task Reference index](../tasks-reference.md)
* [`task web:validate`](./web-validate.md)
* [`task api:generate`](./api-generate.md)
