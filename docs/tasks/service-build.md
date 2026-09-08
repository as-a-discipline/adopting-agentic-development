# `task service:build`

## What it does

Builds the `pulse` binary from `./cmd/pulse` into `service/bin/pulse`, via
the pinned Go container.

## Why it exists

Proves the service actually compiles as a real binary (not just "tests
pass") — also the same binary `service/Dockerfile` packages for
Compose/Helm.

## How to use it

```bash
task service:build
```

## Related

* [Task Reference index](../tasks-reference.md)
* [`task service:test`](./service-test.md)
* [`task service:validate`](./service-validate.md)
