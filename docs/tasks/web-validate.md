# `task web:validate`

## What it does

Runs, in order:

```
generate → generated-freshness check → lint → test → build
```

## Why it exists

The single command that proves the web app is correctly
generated-from-contract, lint-clean, tested, and buildable — the `web`
boundary of `task validate`.

## How to use it

```bash
task web:validate
```

## Related

* [Task Reference index](../tasks-reference.md)
* [`task web:generate`](./web-generate.md)
* [`task web:lint`](./web-lint.md)
* [`task web:test`](./web-test.md)
* [`task web:build`](./web-build.md)
* [`task validate`](./validate.md)
