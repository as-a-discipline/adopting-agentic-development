# `task helm:template:check`

## What it does

Renders the chart (as [`helm:template`](./helm-template.md) does) to a temp
file, then runs `deploy/helm/scripts/check-rendered-manifests.py` against
it — a static check that readiness/liveness probes are present on every
Deployment, no `Secret` resources exist, and the expected resource
kinds/counts are present.

## Why it exists

Proves the rendered output is actually correct, not just that
`helm template` didn't crash — still with no live cluster required.

## How to use it

```bash
task helm:template:check
```

## Related

* [Task Reference index](../tasks-reference.md)
* [`task helm:template`](./helm-template.md)
* [`task helm:validate`](./helm-validate.md)
