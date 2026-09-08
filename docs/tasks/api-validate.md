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
both the stored baseline and mainline/last-pushed-commit, and that the
**Go** generated code (`service/generated/`) hasn't drifted — the `api`
boundary of `task validate`.

**Note:** this does not check `web/generated/` (the TypeScript client) —
that freshness check lives in `task web:validate` instead, since the two
generators are invoked and validated independently. After any contract
change, run both `task api:validate` and `task web:validate` (or the full
`task validate`, which runs both).

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
* [`task web:validate`](./web-validate.md) — the companion check for
  `web/generated/` freshness, not covered by this task
* [`task validate`](./validate.md)
