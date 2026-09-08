# `task api:validate`

## What it does

Runs, in order:

```
lint → compat → breaking-changes → generate → generated-freshness check
```

The freshness check (`scripts/check-generated-fresh.sh`) regenerates to a
temp location and diffs against the committed `service/generated/`, failing
on any difference.

## Why it exists

The single command that proves the API contract is valid, compatible with
both the stored baseline and mainline/last-pushed-commit, and that
generated code hasn't drifted — the `api` boundary of `task validate`.

## How to use it

```bash
task api:validate
```

## Related

* [Task Reference index](../tasks-reference.md)
* [`task api:lint`](./api-lint.md)
* [`task api:compat`](./api-compat.md)
* [`task api:breaking-changes`](./api-breaking-changes.md)
* [`task api:generate`](./api-generate.md)
* [`task validate`](./validate.md)
