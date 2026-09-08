# `task compose:up`

## What it does

Runs `docker compose up -d --build --wait` against
`deploy/compose/docker-compose.yml`, building and starting `pulse-service`,
`pulse-web`, and three deterministic fake HTTP targets
(healthy/degraded/down), waiting for all Compose healthchecks to pass.

## Why it exists

The one command to stand up a fully local, deterministic integration
environment — no external network dependency, no real internet endpoints.

## How to use it

```bash
task compose:up
curl http://localhost:18080/api/v1/services
open http://localhost:18000
```

## Related

* [Task Reference index](../tasks-reference.md)
* [`task compose:down`](./compose-down.md)
* [`task compose:test`](./compose-test.md)
