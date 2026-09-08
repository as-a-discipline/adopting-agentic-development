# `task compose:test`

## What it does

Runs [`up`](./compose-up.md), then `deploy/compose/test.sh` (which verifies
healthy/degraded/down status both directly against the API and via the
`pulse-web` reverse-proxy), with [`down`](./compose-down.md) registered via
Task's `defer:` so teardown always runs, even if the test script fails.

## Why it exists

Deterministic, automated proof that the whole system (service + web +
fakes) behaves correctly when actually deployed together, not just in
isolated unit tests.

## How to use it

```bash
task compose:test
```

## Related

* [Task Reference index](../tasks-reference.md)
* [`task compose:up`](./compose-up.md)
* [`task compose:down`](./compose-down.md)
* [`task compose:validate`](./compose-validate.md)
