# `task policies:containers`

## What it does

Runs `scripts/check-container-quality.sh`, which checks that
`service/Dockerfile` and `web/Dockerfile` each end with a non-root `USER`,
use an explicit (non-`latest`) base image tag, and scans Compose/Helm files
for secret-like literal values.

## Why it exists

A basic, transparent container/deployment quality gate — this check is what
found the real `web/Dockerfile` root-user bug fixed in Session 4 (see
[`policies/README.md`](../../policies/README.md) rule 5).

## How to use it

```bash
task policies:containers
```

## Related

* [Task Reference index](../tasks-reference.md)
* [`task policies:validate`](./policies-validate.md)
