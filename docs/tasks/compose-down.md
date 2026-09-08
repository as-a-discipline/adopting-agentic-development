# `task compose:down`

## What it does

Runs `docker compose down --volumes --remove-orphans`.

## Why it exists

Guarantees a clean teardown (containers, network, volumes) — no leftover
`pulse-*` resources between runs.

## How to use it

```bash
task compose:down
```

## Related

* [Task Reference index](../tasks-reference.md)
* [`task compose:up`](./compose-up.md)
* [`task compose:test`](./compose-test.md)
