---
name: web-development
description: Use when implementing or fixing the Pulse web application (web/) — UI, or regenerating the typed API client from the contract. Trigger on requests to change what the UI shows, fix a failing web test, or wire up a new API operation on the frontend.
---

# Web Development Skill

## When to use

* A UI change is needed for service list/status/summary display, or a check action.
* A `web:test`, `web:lint`, or `web:build` failure needs a root-cause fix.
* The API contract changed and the frontend client needs regeneration.

## Workflow

1. Read `web/AGENTS.md` before editing.
2. If the API contract changed, run `task web:generate` first — do not hand-write
   duplicate models of types that already exist under `web/generated/`.
3. Modify UI code under `web/src/` (or equivalent) only.
4. Run `task web:lint` and `task web:test`.
5. Run `task web:build` to confirm the app still builds.
6. Run `task web:validate` to confirm the component as a whole is valid.

## Evidence of success

* `task web:validate` exits 0.
* No handwritten type duplicates a generated API type.
