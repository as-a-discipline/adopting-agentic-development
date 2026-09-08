# `task policies:boundaries`

## What it does

Runs `scripts/check-component-boundaries.sh` — a grep-based scan of
handwritten source (excluding `generated/`, `bin/`, `dist/`) for
relative-import references that cross a top-level component boundary (e.g.
`service/` code referencing `../web`).

## Why it exists

Enforces
[`/adr/0001-component-boundaries-in-the-monorepo.md`](../../adr/0001-component-boundaries-in-the-monorepo.md)
mechanically rather than by convention alone. Full rule description and
status live in [`policies/README.md`](../../policies/README.md).

## How to use it

```bash
task policies:boundaries
```

## Related

* [Task Reference index](../tasks-reference.md)
* [`task policies:validate`](./policies-validate.md)
