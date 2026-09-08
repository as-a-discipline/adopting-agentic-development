# `task service:validate`

## What it does

Runs, in order:

```
generate → generated-freshness check → fmt → lint → test → build
```

## Why it exists

The single command that proves the Go service is correctly
generated-from-contract, formatted, lint-clean, tested, and buildable — the
`service` boundary of `task validate`.

## How to use it

```bash
task service:validate
```

## Related

* [Task Reference index](../tasks-reference.md)
* [`task service:generate`](./service-generate.md)
* [`task service:fmt`](./service-fmt.md)
* [`task service:lint`](./service-lint.md)
* [`task service:test`](./service-test.md)
* [`task service:build`](./service-build.md)
* [`task validate`](./validate.md)
