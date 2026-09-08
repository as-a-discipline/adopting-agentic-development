# Vestigial: `deploy/Taskfile.yml`'s `validate`

## What it does

`cd deploy && task validate` exists but only prints a pointer message
("Use root Taskfile: task compose:validate / task helm:validate").

## Why it's not part of the reachable task set

It is **not** included in the root `Taskfile.yml`'s `includes:`, so it is
not reachable as `task deploy:validate` from the repo root. It predates
`compose:*` and `helm:*` being split into their own namespaces early in this
repository's history and is not part of the `task validate` chain.

## How to use the real equivalent

```bash
task compose:validate
task helm:validate
```

## Related

* [Task Reference index](../tasks-reference.md)
* [`task compose:validate`](./compose-validate.md)
* [`task helm:validate`](./helm-validate.md)
