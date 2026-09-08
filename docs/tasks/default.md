# `task` (no arguments) / `task default`

Per [`/adr/0002-task-as-the-execution-contract.md`](../../adr/0002-task-as-the-execution-contract.md),
`task <name>` is the only supported entry point for building, linting,
testing, generating, or validating anything in this repository. Running
`task` with no arguments is the friendly discovery command for that
contract.

## What it does

Runs `task --list-all`, printing every task defined across the root and
every included component `Taskfile.yml` (`api`, `service`, `web`,
`policies`, `compose`, `helm`), each with its one-line description.

## Why it exists

A no-argument invocation of `task` should never fail or do nothing — it
should show you what you can do. This is the live, drift-proof source of
truth for "what tasks exist" (it reads directly from the Taskfiles, so it
can never disagree with [`/docs/tasks-reference.md`](../tasks-reference.md)
for long without being noticed).

## How to use it

```bash
task
```

## Related

* [Task Reference index](../tasks-reference.md)
* [`task validate`](./validate.md) — the canonical full validation gate
