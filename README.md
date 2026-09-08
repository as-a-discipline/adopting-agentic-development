# Pulse

Pulse is a tiny service-monitoring application: it checks a few HTTP services and
shows whether they are healthy.

This repository is a small public demonstration built for a conference
presentation about agentic software engineering. **The repository itself is the
demonstration** — the application stays intentionally boring so the engineering
system around it (instructions, skills, contracts, generated-code boundaries, task
execution, guardrails) is what's worth looking at.

## Why This Repository Exists

It demonstrates an engineered environment for agentic development: contract-first
API design, independently governed components sharing one repository, hierarchical
agent instructions, reusable agent skills, staged context, deterministic task
interfaces, generated-code boundaries, deterministic policy enforcement, and
containerized tooling — all backed by real, runnable validation rather than prose
promises. Someone inspecting this repository should be able to understand exactly
what an agent was taught, what skills were available, what commands it could use,
what constraints existed, how generated artifacts were produced, and how
correctness was evaluated — with nothing hidden behind tribal knowledge.

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
independently governed: repository locality does not grant dependency access
(enforced by a deterministic source-boundary scan, not just convention — see
[Guardrails](#guardrails)).

## Repository Structure

```text
/
├── AGENTS.md              Root engineering instructions
├── Makefile                Host bootstrap only (not the build system)
├── Taskfile.yml             Root task composition; `task validate` is canonical
├── adr/                     Cross-cutting architecture decisions (MADR format)
├── verification/            Cross-cutting verification (skills, architecture)
├── scripts/                 Repo-wide validation scripts + git-hooks/ source
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
verification evidence lives with that component; only the cross-cutting
skill/architecture verification lives in root `/verification/`.

## Prerequisites

Only host tools: **Git**, **Docker**, **Make**, **Task**. Go, Node, generators, and
linters run from pinned containers and are not required on the host.

## Bootstrap

```bash
make check-prereqs
make bootstrap
```

Optionally install the local-only git pre-commit hook (see
[Guardrails](#guardrails)):

```bash
task hooks:install
```

## Common Tasks

```bash
task api:lint
task api:validate
task api:breaking-changes   # compare working tree vs last pushed commit / mainline

task service:test
task service:validate

task web:test
task web:validate

task policies:validate

task compose:test
task helm:validate
```

## Run Locally

```bash
task compose:up      # build + start service, web, and deterministic fake targets
curl http://localhost:18080/api/v1/services   # API directly
open http://localhost:18000                    # web UI (or curl it)
task compose:down    # stop and fully clean up
```

## Validate Everything

```bash
task validate
```

This is the canonical repository-level validation entry point; it orchestrates
structure, API, service, web, policy, Compose, and Helm validation — in that
order, failing clearly at whichever boundary breaks — without duplicating any
component's own implementation.

## Agent Instructions

Instructions are scoped hierarchically: root [`AGENTS.md`](./AGENTS.md) holds
global invariants; each component has its own `AGENTS.md` with local rules
(`api/AGENTS.md`, `service/AGENTS.md`, `web/AGENTS.md`, `deploy/AGENTS.md`,
`deploy/compose/AGENTS.md`, `deploy/helm/AGENTS.md`, `policies/AGENTS.md`).
[`.github/copilot-instructions.md`](./.github/copilot-instructions.md) is a small
pointer into this hierarchy, not a duplicate of it. An agent never needs to load
every file — only the root (always) plus whichever component's `AGENTS.md` the
current task actually touches. See
[`verification/architecture.md`](./verification/architecture.md#7-instruction-hierarchy-in-practice--a-worked-example)
for a concrete worked example of this composition.

## Agent Skills

Reusable workflows live in [`.agents/skills/`](./.agents/skills/): `api-contract`,
`service-development`, `web-development`, `compose-integration`,
`helm-validation`, `repository-validation`. Each skill documents its trigger, the
Task interface it uses, and what evidence indicates success — and each has been
verified end-to-end (not just structurally) in
[`verification/skills.md`](./verification/skills.md), which also documents how to
verify skill discovery with both Copilot CLI and OpenCode.

## Generated Code

The OpenAPI contract in `api/openapi.yaml` is the source of truth. Go server
types (`service/generated/`) and the TypeScript client (`web/generated/`) are
generated from it and must never be hand-edited — change the contract or
generator configuration, then regenerate via `task *:generate`. A deterministic
freshness check (regenerate to a temp location, diff against the committed tree)
gates both `api:validate`/`service:validate` and `web:validate`.

## Breaking-Change Detection

`task api:breaking-changes` compares the current working tree against the
current branch's last pushed commit (or the target mainline branch, when
nothing has been pushed yet) and fails if it finds a breaking change. It is
now a real step in `task api:validate` / `task validate` (between `compat`
and `generate`), and it's also run by the local pre-commit hook — a breaking
API change blocks both the validation gate and the commit itself, without
`--no-verify`. Two pinned, containerized tools are supported, selected with
`TOOL=`:

```bash
task api:breaking-changes                              # TOOL=oasdiff (default)
task api:breaking-changes TOOL=openapi-changes          # pb33f/openapi-changes
task api:breaking-changes -- main                       # explicit ref override
```

See [`docs/api-breaking-changes.md`](./docs/api-breaking-changes.md) for how
it works in detail, plus real captured proof (a genuine breaking change
introduced, detected by both tools, and blocking `task validate` and a real
`git commit`).

## Guardrails

* **Execution interface** — all project actions go through `task` targets; Make
  is bootstrap-only.
* **Engineering validation** — lint, generation-freshness checks, tests, and
  static analysis gate every component.
* **Component boundaries** — `web` never assumes Go internals; `service` never
  depends on `web`; `deploy` owns no application behavior. Enforced by a
  grep-based source scan (`task policies:boundaries`), not just convention.
* **Dependency governance** — every direct third-party dependency must appear in
  a small, reviewable allow-list (`policies/allowed-dependencies.txt`), checked
  by `task policies:dependencies`.
* **Container/deployment quality** — Docker images run as non-root users with
  pinned (non-`latest`) base image tags; Compose/Helm files are scanned for
  secret-like literal values; the Helm chart is lint/template validated with no
  live cluster required. Checked by `task policies:containers` and
  `task helm:validate`.
* **Integration validation** — deterministic Docker Compose tests validate the
  system as deployed (healthy/degraded/down behavior end-to-end), not just in
  isolation.
* **Local pre-commit hook** (opt-in, local-only) — `task hooks:install` installs
  a fast guardrail subset (`task structure:validate` + `task policies:validate`
  + `task api:breaking-changes`) as `.git/hooks/pre-commit`. It is never
  tracked by git or pushed; it is not a CI/CD mechanism, just a fast local
  check before the full `task validate` gate — a breaking API change is
  blocked at commit time, not just at `task validate` time.

See [`policies/README.md`](./policies/README.md) for the full rule set, each
rule's exact enforcement command, and its status.

## Architecture Decisions

See [`/adr/`](./adr/) for cross-cutting decisions and each component's own
`adr/` folder (e.g. [`api/adr/`](./api/adr/), [`service/adr/`](./service/adr/))
for component-specific decisions. Records use the
[MADR](https://adr.github.io/madr/) format.

## Verification

See [`/verification/`](./verification/) for the skill verification matrix
(including how to verify skill discovery with Copilot CLI and OpenCode) and
manual architecture-verification steps (including a worked example of the
instruction hierarchy in practice).

## Scope

The baseline intentionally omits: persistence/database, authentication, alerting,
historical metrics, external integrations, and advanced diagnostics (e.g.
traceroute). Future work must not introduce these without an explicit requirement
change captured in a new ADR.
