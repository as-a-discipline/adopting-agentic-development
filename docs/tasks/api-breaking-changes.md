# `task api:breaking-changes`

## What it does

Compares the current working tree against the current branch's last pushed
commit (its upstream tracking branch), or the target mainline branch
(`origin/main` → `origin/master` → local `main` → local `master`, whichever
resolves first) if nothing has been pushed yet. Uses either `oasdiff`
(default) or `pb33f/openapi-changes`, both pinned containers, selected with
`TOOL=`.

## Why it exists

Answers a different question than [`api:compat`](./api-compat.md): "did we
break compatibility with what's already pushed / with mainline?" — useful
before a baseline update is even proposed. It's enforced in both
`task api:validate` and the local pre-commit hook, so a breaking change is
caught before it's even committed, not just before it's merged.

## How to use it

```bash
task api:breaking-changes                          # TOOL=oasdiff (default), auto-detected ref
task api:breaking-changes TOOL=openapi-changes      # pb33f/openapi-changes instead
task api:breaking-changes -- main                   # explicit ref override
task api:breaking-changes TOOL=openapi-changes -- v1.2.0
```

## Full guide and captured proof

This task has a dedicated, in-depth guide with real captured evidence
(including a genuine breaking change being detected and a real `git commit`
being blocked): **[`/docs/api-breaking-changes.md`](../api-breaking-changes.md)**.
It covers:

* the full ref-resolution algorithm,
* a comparison table of the two supported tools,
* why a custom JSON parser was needed for `pb33f/openapi-changes`,
* exactly where this task is wired in (`api:validate`, pre-commit hook),
* five captured proof scenarios with raw logs.

## Related

* [Task Reference index](../tasks-reference.md)
* [`task api:compat`](./api-compat.md)
* [`task api:validate`](./api-validate.md)
* [`task hooks:install`](./hooks-install.md)
