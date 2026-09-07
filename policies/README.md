# Policies

Deterministic engineering constraints applied across the Pulse repository. These
are intentionally few and transparent — the goal is to demonstrate executable
engineering constraints, not to build an enterprise policy framework.

Each rule below states its current enforcement status. Rules that can't yet be
checked (because the component they apply to doesn't exist yet) are documented now
and wired into `task validate` as the relevant component is built.

## 1. Generated Code Integrity

**Rule:** Directories under a component's `generated/` path may not be manually
maintained. They must exactly match what regenerating from the source contract
produces.

**Mechanism:** regenerate, then diff against the committed tree; fail if non-empty.

**Status:** documented now; enforced starting in Session 2 (`service/generated/`)
and Session 3 (`web/generated/`) once those directories exist.

## 2. API Contract Must Lint and Stay Structurally Valid

**Rule:** `api/openapi.yaml` must pass linting and remain a structurally valid
OpenAPI document before generation or merge.

**Mechanism:** `task api:lint`, run from a pinned container.

**Status:** documented now; enforced starting in Session 2 once `openapi.yaml`
exists.

## 3. Component Boundaries

**Rule:** No component may import another component's implementation internals
directly (e.g. `web` importing Go source, `service` importing frontend build
output). See [`/adr/0001-component-boundaries-in-the-monorepo.md`](../adr/0001-component-boundaries-in-the-monorepo.md).

**Mechanism:** a simple, documented check (e.g. grep-based import-path scan) run as
part of `task structure:validate` / `task validate`, flagging obvious violations.
This is a transparency mechanism, not an exhaustive static-analysis platform.

**Status:** basic structural check active now (see `scripts/validate-structure.sh`);
deeper source-import scanning is added once `service/` and `web/` have real source
trees.

## 4. Dependency Governance (stub)

**Rule:** New third-party dependencies should be recorded somewhere reviewable
(e.g. `go.mod` / `package.json` diffs are visible in the change) rather than
silently vendored or hidden. This starts as a transparency convention, not an
automated license/supportability scanner.

**Status:** convention documented now; a lightweight, understandable check (e.g.
a dependency-manifest diff summary in `task validate` output) may be added once
`service`/`web` have real manifests. No enterprise dependency platform will be
introduced.

## 5. Container/Deployment Quality (basic)

**Rule:** Docker images and the Helm chart must pass their respective lint/template
validation, use non-root runtime users where practical, and commit no secrets.

**Mechanism:** `task helm:lint` / `task helm:template`; Dockerfile review as part of
component validation.

**Status:** documented now; enforced starting in Session 3 once Dockerfiles and the
Helm chart exist.
