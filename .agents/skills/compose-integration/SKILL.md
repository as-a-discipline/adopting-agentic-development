---
name: compose-integration
description: Use when starting, testing, or debugging the Docker Compose integration environment for Pulse (deploy/compose/). Trigger on requests to run end-to-end/integration checks, or to investigate a compose:test failure.
---

# Compose Integration Skill

## When to use

* An end-to-end check of the API + web + deterministic fake targets is needed.
* `compose:test` is failing and the failure needs to be isolated.

## Workflow

1. Read `deploy/compose/AGENTS.md` before changing the environment.
2. Start the environment: `task compose:up`.
3. Wait for health checks to report healthy (the compose file defines explicit
   health checks — do not add ad hoc sleep/polling logic outside the task).
4. Run `task compose:test` to verify API and UI behavior against the deterministic
   healthy/degraded/down fake targets.
5. On failure, collect logs with `task compose:logs` before making changes.
6. Always tear down with `task compose:down`, even after a failure.
7. Run `task compose:validate` for the full deterministic check.

## Evidence of success

* `task compose:validate` exits 0.
* `task compose:down` leaves no dangling containers/networks/volumes.
