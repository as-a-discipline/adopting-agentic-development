# `task policies:validate`

## What it does

Runs, in order:

```
boundaries → dependencies → containers
```

Rules 1 and 2 (generated-code integrity and API lint/compat) are already
enforced by `api:validate`/`service:validate`/`web:validate` and are
deliberately not duplicated here.

## Why it exists

The single command for "all currently-enforceable policy rules"; also the
check run by the local pre-commit hook (alongside `structure:validate` and
`api:breaking-changes`).

## How to use it

```bash
task policies:validate
```

## Related

* [Task Reference index](../tasks-reference.md)
* [`task policies:boundaries`](./policies-boundaries.md)
* [`task policies:dependencies`](./policies-dependencies.md)
* [`task policies:containers`](./policies-containers.md)
* [`policies/README.md`](../../policies/README.md)
* [`task validate`](./validate.md)
