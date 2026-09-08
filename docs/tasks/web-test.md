# `task web:test`

## What it does

Runs Node's built-in test runner
(`node --experimental-transform-types --test src/**/*.test.ts`) in the
pinned `node:22-alpine` container.

## Why it exists

Runs unit tests for pure logic (status mapping, summary computation) with
zero test-framework dependency.

## How to use it

```bash
task web:test
```

## Related

* [Task Reference index](../tasks-reference.md)
* [`task web:build`](./web-build.md)
* [`task web:validate`](./web-validate.md)
