# Pulse

Pulse is a tiny service-monitoring application: it checks a few HTTP services and
shows whether they are healthy.

> Pulse checks a few HTTP services and shows whether they are healthy.

This repository is a small public demonstration built for a conference
presentation about agentic software engineering. **The repository itself is the
demonstration** — the application stays intentionally boring so the engineering
system around it (instructions, skills, contracts, generated-code boundaries, task
execution, guardrails) is what's worth looking at.

## Why This Repository Exists

It demonstrates an engineered environment for agentic development: contract-first
API design, independently governed components sharing one repository, hierarchical
agent instructions, reusable agent skills, staged context, deterministic task
interfaces, generated-code boundaries, and containerized tooling — all backed by
real, runnable validation rather than prose promises.

## Architecture

```text
                    OpenAPI Contract
                          │
              ┌───────────┴───────────┐
              │                       │
       generated Go            generated TS
       interfaces                 client
              │                       │
              ▼                       ▼
        Go Service ────────────── Web App
              │                       │
              └───────────┬───────────┘
                          │
                   Docker Compose
                          │
                      Helm Chart
```

`api` owns the OpenAPI contract. `service` (Go) implements it. `web` (TypeScript)
consumes only the generated client. `deploy` packages built artifacts for local
Compose integration and Kubernetes (Helm) — it owns none of the application
behavior. Although this is one repository, each component is treated as
independently governed: repository locality does not grant dependency access.

## Repository Structure

```text
/
├── AGENTS.md              Root engineering instructions
├── Makefile                Host bootstrap only (not the build system)
├── Taskfile.yml             Root task composition; `task validate` is canonical
├── adr/                     Cross-cutting architecture decisions (MADR format)
├── verification/            Cross-cutting verification (skills, architecture)
├── .github/
│   └── copilot-instructions.md   Small pointer to AGENTS.md/skills/Task
├── .agents/skills/           Reusable agent workflows (SKILL.md per skill)
├── api/                      OpenAPI contract + api/adr/ (API-specific decisions)
├── service/                  Go implementation + service/adr/
├── web/                      Web application (generated client consumer)
├── deploy/
│   ├── compose/              Deterministic local integration environment
│   └── helm/                 Kubernetes packaging
└── policies/                 Deterministic engineering constraints
```

Architecture decision records are colocated with the component they concern
(`api/adr/`, `service/adr/`); only genuinely cross-cutting decisions live in root
`/adr/`. The same principle applies to verification docs: component-specific
verification evidence will live with that component as it's built; only the
cross-cutting skill/architecture verification lives in root `/verification/`.

## Prerequisites

Only host tools: **Git**, **Docker**, **Make**, **Task**. Go, Node, generators, and
linters run from pinned containers and are not required on the host.

## Bootstrap

```bash
make check-prereqs
make bootstrap
```

## Common Tasks

```bash
task api:lint
task api:validate

task service:test
task service:validate

task web:test
task web:validate

task compose:test
task helm:validate
```

## Validate Everything

```bash
task validate
```

This is the canonical repository-level validation entry point; it orchestrates
structure, API, service, web, Compose, and Helm validation without duplicating
any component's implementation.

## Agent Instructions

Instructions are scoped hierarchically: root [`AGENTS.md`](./AGENTS.md) holds
global invariants; each component has its own `AGENTS.md` with local rules
(`api/AGENTS.md`, `service/AGENTS.md`, `web/AGENTS.md`, `deploy/AGENTS.md`,
`deploy/compose/AGENTS.md`, `deploy/helm/AGENTS.md`, `policies/AGENTS.md`).
[`.github/copilot-instructions.md`](./.github/copilot-instructions.md) is a small
pointer into this hierarchy, not a duplicate of it.

## Agent Skills

Reusable workflows live in [`.agents/skills/`](./.agents/skills/): `api-contract`,
`service-development`, `web-development`, `compose-integration`,
`helm-validation`, `repository-validation`. Each skill documents its trigger, the
Task interface it uses, and what evidence indicates success.

## Generated Code

The OpenAPI contract in `api/openapi.yaml` (Session 2) is the source of truth.
Go server types (`service/generated/`) and the TypeScript client
(`web/generated/`) are generated from it and must never be hand-edited — change
the contract or generator configuration, then regenerate via `task *:generate`.

## Guardrails

* **Execution interface** — all project actions go through `task` targets; Make
  is bootstrap-only.
* **Engineering validation** — lint, generation-freshness checks, tests, and
  static analysis gate every component.
* **Component boundaries** — `web` never assumes Go internals; `service` never
  depends on `web`; `deploy` owns no application behavior.
* **Integration validation** — deterministic Docker Compose tests and Helm
  lint/template checks validate the system as deployed, not just in isolation.

## Architecture Decisions

See [`/adr/`](./adr/) for cross-cutting decisions and each component's own
`adr/` folder (e.g. [`api/adr/`](./api/adr/), [`service/adr/`](./service/adr/))
for component-specific decisions. Records use the
[MADR](https://adr.github.io/madr/) format.

## Verification

See [`/verification/`](./verification/) for the skill verification matrix and
manual architecture-verification steps.

## Current Implementation Status

This is **Session 1**: repository skeleton only (instructions, skills, Task
composition, ADRs, policy placeholders, verification docs). No OpenAPI contract,
Go service, web application, Compose environment, or Helm chart exists yet —
those are built in later, separately-run sessions. Component `task` targets for
not-yet-built parts fail with an explicit "not implemented yet" message rather
than a fake success.

## Scope

The baseline intentionally omits: persistence/database, authentication, alerting,
historical metrics, external integrations, and advanced diagnostics (e.g.
traceroute). Future work must not introduce these without an explicit requirement
change captured in a new ADR.