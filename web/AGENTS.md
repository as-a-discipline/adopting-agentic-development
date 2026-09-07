# AGENTS.md — web/

Consumes the generated public API client. Does not implement business logic beyond
presenting monitored-service status.

## Toolchain

No npm dependencies are installed for this component. The app is built and tested
using only Node's built-in capabilities:

* **Generation** — `web:generate` runs `openapi-generator-cli` (pinned container)
  to produce a TypeScript client under `generated/`, then a small fixup script
  (`tools/fix-generated-imports.mjs`) normalizes its relative import specifiers to
  explicit `.ts` extensions (required by Node's strict ESM resolver).
* **Check (`web:lint`)** — `tools/check.mjs` verifies every source file strips
  cleanly via Node's built-in TypeScript type-stripping (`node:module`
  `stripTypeScriptTypes`) and that every relative import resolves to a real file.
  This is a syntax + import-resolution safety net, not full type checking — no
  TypeScript compiler binary is installed.
* **Test (`web:test`)** — Node's built-in test runner (`node --test`,
  `node:assert`), executing `.test.ts` files directly via
  `--experimental-transform-types`. No test framework dependency.
* **Build (`web:build`)** — `tools/build.mjs` compiles `src/` and `generated/`
  TypeScript to plain JavaScript (type-stripped, `.ts` import specifiers rewritten
  to `.js`) into `dist/`, which is directly loadable by browsers via native ES
  modules — no bundler.

This keeps the dependency surface minimal and the build fully deterministic and
offline-capable, consistent with the containerized-tooling ADR (`/adr/0003`).

## Rules

* Consume the generated API client under `generated/`; do not hand-write duplicate
  models of API types that already exist as generated types.
* No assumptions about Go service internals — the only contract is the generated
  client derived from `api/openapi.yaml`.
* Keep the UI deliberately small: header, summary line, service list/cards, and a
  per-service check action.
* Basic accessibility and responsive behavior are required; a large design system
  is not.
* Use `task web:*` targets — do not invoke `node`/frontend tooling directly outside
  of them.
* `generated/` and hand-written source files under `tools/` are not the same
  thing: `tools/*.mjs` are hand-written build tooling and may be edited freely;
  `generated/**/*.ts` must never be hand-edited (see root `AGENTS.md`'s
  Generated-Code Rule).

## Task Interface

```text
task web:generate
task web:lint
task web:test
task web:build
task web:validate
```
