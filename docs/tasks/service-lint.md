# `task service:lint`

## What it does

Runs `golangci-lint run ./...` via a pinned
`golangci/golangci-lint:v1.64.8` container.

## Why it exists

Standard Go static-analysis aggregation (unused code, common bug patterns,
style) without a host toolchain.

## How to use it

```bash
task service:lint
```

## Related

* [Task Reference index](../tasks-reference.md)
* [`task service:fmt`](./service-fmt.md)
* [`task service:validate`](./service-validate.md)
