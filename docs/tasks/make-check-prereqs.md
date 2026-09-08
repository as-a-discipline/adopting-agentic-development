# `make check-prereqs`

Per
[`/adr/0002-task-as-the-execution-contract.md`](../../adr/0002-task-as-the-execution-contract.md),
`make` targets intentionally do none of the actual build/test/lint/generate/
validate work — they exist purely to get `task` itself onto a fresh machine.

## What it does

Verifies `git`, `docker` (and that the Docker daemon is actually running),
and `task` are installed; exits non-zero listing anything missing.

## Why it exists

The first thing to run on a fresh clone/machine, before any `task` command
can work at all.

## How to use it

```bash
make check-prereqs
```

## Related

* [Task Reference index](../tasks-reference.md)
* [`make bootstrap`](./make-bootstrap.md)
