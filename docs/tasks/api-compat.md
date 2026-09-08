# `task api:compat`

## What it does

Compares [`api/openapi.yaml`](../../api/openapi.yaml) against the committed
snapshot `api/openapi.baseline.yaml` using `oasdiff breaking --fail-on ERR`
(`tufin/oasdiff:v1.31.0`, pinned container). If no baseline file exists yet,
it bootstraps one from the current contract instead of failing, with an
explicit message that this establishes, not proves, compatibility.

## Why it exists

A static, always-available compatibility gate that doesn't depend on git
history or network access — it answers "did we break compatibility with the
last explicitly-agreed-good contract?". Updating the baseline (i.e.
accepting a breaking change on purpose) is a deliberate, reviewable diff to
`openapi.baseline.yaml`, not a silent bypass.

`--fail-on ERR` matters: oasdiff's `breaking` subcommand does **not** fail
on its own regardless of detected breaking changes unless `--fail-on ERR`
(or `WARN`) is explicitly passed — this was a real, silent bug fixed
alongside adding [`task api:breaking-changes`](./api-breaking-changes.md).
See [`docs/api-breaking-changes.md`](../api-breaking-changes.md) for the
full story.

## How to use it

```bash
task api:compat
```

To accept an intentional breaking change:

```bash
cp api/openapi.yaml api/openapi.baseline.yaml
```

...as part of the same commit, with the reason visible in the diff/PR.

## Related

* [Task Reference index](../tasks-reference.md)
* [`task api:breaking-changes`](./api-breaking-changes.md) — a related but
  distinct check (compares against last-pushed-commit/mainline, not a
  stored baseline)
* [`task api:validate`](./api-validate.md)
