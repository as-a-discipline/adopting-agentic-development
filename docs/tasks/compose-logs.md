# `task compose:logs`

## What it does

Runs `docker compose logs` against the stack defined in
`deploy/compose/docker-compose.yml`.

## Why it exists

Quick diagnostic access to every service's logs without remembering
individual container names.

## How to use it

```bash
task compose:logs
```

## Related

* [Task Reference index](../tasks-reference.md)
* [`task compose:up`](./compose-up.md)
