# `task helm:template`

## What it does

Runs `helm template pulse pulse`, printing the fully rendered Kubernetes
manifests to stdout.

## Why it exists

Lets you inspect exactly what Kubernetes objects the chart would create,
without a live cluster.

## How to use it

```bash
task helm:template
```

## Related

* [Task Reference index](../tasks-reference.md)
* [`task helm:template:check`](./helm-template-check.md)
