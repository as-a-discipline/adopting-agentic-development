# `task helm:lint`

## What it does

Runs `helm lint pulse` via a pinned `alpine/helm:3.21.4` container against
the chart in `deploy/helm/pulse/`.

## Why it exists

Standard Helm chart structural/syntax validation.

## How to use it

```bash
task helm:lint
```

## Related

* [Task Reference index](../tasks-reference.md)
* [`task helm:validate`](./helm-validate.md)
