# `task service:test`

## What it does

Runs `go test ./... -v -cover` via the pinned `golang:1.23-alpine`
container.

## Why it exists

Runs the Go unit/handler test suite (status normalization, HTTP handler
behavior, edge cases) with coverage reporting, with no host Go required.

## How to use it

```bash
task service:test
```

## Related

* [Task Reference index](../tasks-reference.md)
* [`task service:build`](./service-build.md)
* [`task service:validate`](./service-validate.md)
