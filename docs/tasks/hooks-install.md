# `task hooks:install`

## What it does

Copies [`scripts/git-hooks/pre-commit`](../../scripts/git-hooks/pre-commit)
to `.git/hooks/pre-commit` and makes it executable.

## Why it exists

Installs the **local-only** pre-commit guardrail described in the
[README's Guardrails section](../../README.md#guardrails). `.git/hooks/` is
never tracked by git, so this hook only exists on a machine after someone
explicitly runs this task — it is never silently applied to a fresh clone,
never pushed, and it is not a CI/CD mechanism. The hook runs
`structure:validate`, `policies:validate`, and `api:breaking-changes` — a
fast subset, not the full `task validate` (see
[`docs/api-breaking-changes.md`](../api-breaking-changes.md) for the
breaking-changes half of that wiring, with captured proof of a real commit
being blocked).

## How to use it

```bash
task hooks:install
```

Re-run this any time `scripts/git-hooks/pre-commit` changes — the installed
copy in `.git/hooks/` is not auto-synced. To bypass the hook deliberately
(standard git behavior, not something this repo fights):

```bash
git commit --no-verify
```

## Related

* [Task Reference index](../tasks-reference.md)
* [`task structure:validate`](./structure-validate.md)
* [`task policies:validate`](./policies-validate.md)
* [`task api:breaking-changes`](./api-breaking-changes.md)
