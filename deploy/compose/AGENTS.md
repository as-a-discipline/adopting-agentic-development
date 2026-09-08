# AGENTS.md — deploy/compose/

Deterministic local Docker Compose integration environment.

## Rules

* No external network dependencies for test success — use deterministic, trivial
  fake HTTP targets (e.g. `healthy-target`, `degraded-target`, `down-target`)
  instead of the public internet.
* Explicit health checks for every service in the Compose file.
* Runtime state is disposable; `compose:down` must clean up fully.
* `compose:test` must build/start the stack, wait deterministically for readiness,
  verify representative healthy/degraded/down behavior, and tear down cleanly even
  on failure where practical.

## Task Interface

```text
task compose:up
task compose:down
task compose:test
task compose:logs
task compose:validate
```
