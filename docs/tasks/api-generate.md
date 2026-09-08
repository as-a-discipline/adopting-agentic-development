# `task api:generate`

## What it does

Runs `oapi-codegen` (via a pinned `golang:1.23-alpine` container,
version-pinned with `@v2.4.1`) against
[`api/openapi.yaml`](../../api/openapi.yaml), producing
`service/generated/api.gen.go` (Go types, `std-http` server interfaces, and
the embedded spec). Depends on [`api:lint`](./api-lint.md) — it won't
generate from an invalid contract.

## Why it exists

`service/generated/` must never be hand-written — this is the one mechanism
that produces it, per
[`api/adr/0001-contract-first-api.md`](../../api/adr/0001-contract-first-api.md).

## How to use it

```bash
task api:generate
```

You will not normally run this directly — `task api:validate` and
`task service:generate` both call it. Run it directly right after editing
`api/openapi.yaml` to see the generated diff before running full
validation.

**Note:** this only regenerates the **Go** server artifacts
(`service/generated/`). It does not touch the TypeScript client — after any
contract change, also run `task web:generate` (see
[`docs/tasks/web-generate.md`](./web-generate.md)) to keep `web/generated/`
in sync.

## Related

* [Task Reference index](../tasks-reference.md)
* [`task api:validate`](./api-validate.md)
* [`task service:generate`](./service-generate.md)
* [`task web:generate`](./web-generate.md) — the companion regeneration step
  for the TypeScript client; not called by this task
