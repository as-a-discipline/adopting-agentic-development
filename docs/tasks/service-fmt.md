# `task service:fmt`

## What it does

Runs `gofmt -l -w` over `service/cmd` and `service/internal` inside a pinned
`golang:1.23-alpine` container, then fails if anything still needed
formatting — i.e. it both fixes and verifies in one step.

## Why it exists

Keeps Go source canonically formatted without requiring a host Go install.

## How to use it

```bash
task service:fmt
```

## Related

* [Task Reference index](../tasks-reference.md)
* [`task service:lint`](./service-lint.md)
* [`task service:validate`](./service-validate.md)
