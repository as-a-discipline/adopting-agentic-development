# `task helm:validate`

## What it does

Runs, in order:

```
lint → template:check
```

## Why it exists

The single command that proves the Helm chart is valid and its rendered
output meets baseline deployment-quality expectations — the `helm` boundary
of `task validate`.

## How to use it

```bash
task helm:validate
```

## Related

* [Task Reference index](../tasks-reference.md)
* [`task helm:lint`](./helm-lint.md)
* [`task helm:template:check`](./helm-template-check.md)
* [`task validate`](./validate.md)
