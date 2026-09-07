# Policies

Deterministic engineering constraints applied across the Pulse repository. These
are intentionally few and transparent — the goal is to demonstrate executable
engineering constraints, not to build an enterprise policy framework.

Each rule below states its current enforcement status and the exact command that
enforces it. All five rules are enforced as of Session 4, wired into
`task validate` (rules 1-2 via `api:validate`/`service:validate`/`web:validate`;
rules 3-5 via `task policies:validate`, see `policies/Taskfile.yml`).

## 1. Generated Code Integrity

**Rule:** Directories under a component's `generated/` path may not be manually
maintained. They must exactly match what regenerating from the source contract
produces.

**Mechanism:** `scripts/check-generated-fresh.sh` regenerates to a temp location
and diffs against the committed tree, failing on any difference.

**Status:** enforced. Run directly via `task api:validate` (checks
`service/generated/`) and `task web:validate` (checks `web/generated/`); both are
part of `task validate`.

## 2. API Contract Must Lint and Stay Structurally Valid

**Rule:** `api/openapi.yaml` must pass linting and remain a structurally valid
OpenAPI document before generation or merge, and must not introduce breaking
changes against the stored baseline without an explicit, reviewed baseline
update.

**Mechanism:** `task api:lint` (Redocly CLI, pinned container) plus
`task api:compat` (oasdiff, pinned container, `breaking --fail-on ERR`
against `api/openapi.baseline.yaml`). A related but separate developer tool,
`task api:breaking-changes`, compares the working tree against the last
pushed commit or target mainline branch (via `oasdiff` or
`pb33f/openapi-changes`) — useful before a baseline update is even proposed,
but not itself part of the `task validate` gate.

**Status:** enforced. Part of `task api:validate` / `task validate`. Note:
`task api:compat`'s underlying `oasdiff breaking` call does not fail on its
own — oasdiff only returns a non-zero exit code when `--fail-on ERR` (or
`WARN`) is passed; this was a real gap found and fixed while adding
`task api:breaking-changes` (verified by injecting a breaking change and
confirming `task api:compat` failed only after adding the flag).

## 3. Component Boundaries

**Rule:** No component may import another component's implementation internals
directly (e.g. `web` importing Go source, `service` importing frontend build
output). See [`/adr/0001-component-boundaries-in-the-monorepo.md`](../adr/0001-component-boundaries-in-the-monorepo.md).

**Mechanism:** `scripts/check-component-boundaries.sh` — a transparent,
grep-based scan of handwritten source (excluding `generated/`/`bin/`/`dist/`) for
relative-import references that reach across a top-level component boundary
(e.g. `service/` source referencing `../web`). This is a transparency mechanism,
not an exhaustive static-analysis platform — a reader can verify exactly what it
checks by reading the script.

**Status:** enforced. Run via `task policies:boundaries` (part of
`task policies:validate` / `task validate`).

## 4. Dependency Governance

**Rule:** New direct third-party dependencies must be recorded in a small,
reviewable allow-list rather than added silently.

**Mechanism:** `scripts/check-dependency-governance.py` lists each component's
direct dependencies (`service/go.mod` non-indirect requires,
`web/package.json` `dependencies`) and confirms every one is present in
[`allowed-dependencies.txt`](./allowed-dependencies.txt). Adding a new direct
dependency requires an explicit, visible diff to that file. This is a
transparency convention, not an automated license/supportability scanner — no
enterprise dependency-management platform is introduced.

**Status:** enforced. Run via `task policies:dependencies` (part of
`task policies:validate` / `task validate`).

## 5. Container/Deployment Quality (basic)

**Rule:** Docker images and the Helm chart must pass their respective lint/
template validation, run as a non-root user, use pinned (non-`latest`) base
image tags, and commit no secrets.

**Mechanism:** `task helm:lint` / `task helm:template:check` (Helm side); and
`scripts/check-container-quality.sh`, which statically checks both
`service/Dockerfile` and `web/Dockerfile` for a non-root final `USER`, an
explicit (non-`latest`) base image tag, and scans Compose/Helm files for
secret-like literal values.

**Status:** enforced. Run via `task helm:validate` and
`task policies:containers` (both part of `task validate`). This check found and
led to fixing a real issue during Session 4: `web/Dockerfile`'s runtime image
(`nginx:1.27.5-alpine`) ran as root; switched to
`nginxinc/nginx-unprivileged:1.27.5-alpine` (listening on 8080, since
unprivileged users can't bind port 80) with an explicit `USER nginx`.

