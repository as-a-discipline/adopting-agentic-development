# `task policies:dependencies`

## What it does

Runs `scripts/check-dependency-governance.py`, which lists
`service/go.mod`'s direct (non-`// indirect`) requires and
`web/package.json`'s `dependencies`, and fails if any of them is missing
from [`policies/allowed-dependencies.txt`](../../policies/allowed-dependencies.txt).

## Why it exists

Makes adding a new direct third-party dependency a visible, reviewable diff
instead of a silent addition — without introducing a license/supportability
scanning platform. Full rule description and status live in
[`policies/README.md`](../../policies/README.md).

## How to use it

```bash
task policies:dependencies
```

Adding a real new dependency: add the line to `allowed-dependencies.txt` in
the same commit that introduces the dependency.

## Related

* [Task Reference index](../tasks-reference.md)
* [`task policies:validate`](./policies-validate.md)
