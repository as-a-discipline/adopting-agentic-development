# Task Reference: Every `task` Command in This Repository

This is the single reference for every task defined across this repository's
Taskfiles: what it does, why it exists, and how to run it. Per
[`/adr/0002-task-as-the-execution-contract.md`](../adr/0002-task-as-the-execution-contract.md),
`task <name>` is the **only** supported way to build, lint, test, generate, or
validate anything here — there is no other documented entry point, and this
document is meant to make that contract fully transparent in one place.

Run `task --list-all` from the repo root at any time to see this same list
live, with one-line descriptions, generated directly from the Taskfiles below
(so it can never drift from what's actually runnable).

## How this document is organized

Tasks are grouped by namespace, matching the repository structure
(`api:*`, `service:*`, `web:*`, `policies:*`, `compose:*`, `helm:*`, plus
root-level tasks with no namespace). Each entry has:

* **Command** — exactly what to type.
* **What it does** — the real mechanism (which container, which script).
* **Why it exists** — the guardrail or workflow need it satisfies.
* **How to use it** — example invocations, including any variables/flags.

For the two-sentence version of the whole validation pipeline, see the root
[README's "Validate Everything" section](../README.md#validate-everything).

---

## Root-level tasks

### `task` (no arguments) / `task default`

* **What it does:** Runs `task --list-all`, printing every available task
  with its one-line description.
* **Why it exists:** A friendly no-argument entry point — running bare
  `task` shows you what you can do instead of failing or doing nothing.
* **How to use it:**
  ```bash
  task
  ```

### `task validate`

* **What it does:** Runs, in order: `structure:validate` → `api:validate` →
  `service:validate` → `web:validate` → `policies:validate` →
  `compose:validate` → `helm:validate`. Stops at the first failure — it does
  not continue past a broken boundary and does not reimplement any
  component's own logic (it only composes/delegates).
* **Why it exists:** The single canonical "is this repository correct right
  now" gate — the same one a contributor runs before pushing, and the same
  one this repository's own documentation points to as proof of engineering
  soundness. See
  [`/docs/api-breaking-changes.md`](./api-breaking-changes.md) for a real,
  captured example of this task failing on a genuine breaking API change.
* **How to use it:**
  ```bash
  task validate
  ```
  No variables. Takes a few minutes on first run (builds containers, starts
  a Compose stack for `compose:validate`); subsequent runs are faster due to
  Docker layer caching.

### `task structure:validate`

* **What it does:** Runs `scripts/validate-structure.sh`, which checks that
  required root files, every component's `AGENTS.md`, every component's
  `Taskfile.yml`, all 6 skills (with valid `SKILL.md` frontmatter), the
  required root and component ADR directories, and the root verification
  docs all exist.
* **Why it exists:** A deterministic check that the repository's own
  documented shape (instructions, skills, task interfaces, ADRs) hasn't
  silently rotted — the first, fastest thing `task validate` checks.
* **How to use it:**
  ```bash
  task structure:validate
  ```

### `task hooks:install`

* **What it does:** Copies `scripts/git-hooks/pre-commit` to
  `.git/hooks/pre-commit` and makes it executable.
* **Why it exists:** Installs the **local-only** pre-commit guardrail (see
  [Guardrails in the README](../README.md#guardrails)). `.git/hooks/` is
  never tracked by git, so this hook only exists on a machine after someone
  explicitly runs this task — it is never silently applied to a fresh clone
  and it is not a CI/CD mechanism.
* **How to use it:**
  ```bash
  task hooks:install
  ```
  Re-run this any time `scripts/git-hooks/pre-commit` changes — the
  installed copy in `.git/hooks/` is not auto-synced.

---

## `api:*` — OpenAPI contract (`api/`)

### `task api:lint`

* **What it does:** Runs Redocly CLI (`redocly/cli:1.25.3`, pinned
  container) against `api/openapi.yaml`.
* **Why it exists:** The contract must stay structurally valid OpenAPI and
  pass sensible-default linting before anything downstream (generation,
  compatibility checks) is trusted. See
  [`policies/README.md`](../policies/README.md) rule 2.
* **How to use it:**
  ```bash
  task api:lint
  ```

### `task api:compat`

* **What it does:** Compares `api/openapi.yaml` against the committed
  snapshot `api/openapi.baseline.yaml` using
  `oasdiff breaking --fail-on ERR` (`tufin/oasdiff:v1.31.0`, pinned
  container). If no baseline file exists yet, it bootstraps one from the
  current contract instead of failing (with an explicit message that this
  establishes, not proves, compatibility).
* **Why it exists:** A static, always-available compatibility gate that
  doesn't depend on git history or network access — it answers "did we
  break compatibility with the last explicitly-agreed-good contract?".
  Updating the baseline (i.e. accepting a breaking change on purpose) is a
  deliberate, reviewable diff to `openapi.baseline.yaml`.
* **How to use it:**
  ```bash
  task api:compat
  ```
  To accept an intentional breaking change: update
  `api/openapi.baseline.yaml` (e.g. `cp api/openapi.yaml api/openapi.baseline.yaml`)
  as part of the same commit, with the reason visible in the diff/PR.

### `task api:breaking-changes`

* **What it does:** Compares the current working tree against the current
  branch's last pushed commit (its upstream tracking branch), or the target
  mainline branch (`origin/main` → `origin/master` → local `main` → local
  `master`, whichever resolves first) if nothing has been pushed yet. Uses
  either `oasdiff` (default) or `pb33f/openapi-changes`, both pinned
  containers.
* **Why it exists:** Answers a different question than `api:compat`: "did we
  break compatibility with what's already pushed / with mainline?" — useful
  before a baseline update is even proposed, and enforced both in
  `task api:validate` and the local pre-commit hook. Full details, tool
  comparison, and captured proof:
  [`/docs/api-breaking-changes.md`](./api-breaking-changes.md).
* **How to use it:**
  ```bash
  task api:breaking-changes                          # TOOL=oasdiff (default), auto-detected ref
  task api:breaking-changes TOOL=openapi-changes      # pb33f/openapi-changes instead
  task api:breaking-changes -- main                   # explicit ref override
  task api:breaking-changes TOOL=openapi-changes -- v1.2.0
  ```

### `task api:generate`

* **What it does:** Runs `oapi-codegen` (via a pinned `golang:1.23-alpine`
  container, version-pinned with `@v2.4.1`) against `api/openapi.yaml`,
  producing `service/generated/api.gen.go` (Go types, `std-http` server
  interfaces, and the embedded spec). Depends on `api:lint` (won't generate
  from an invalid contract).
* **Why it exists:** `service/generated/` must never be hand-written — this
  is the one mechanism that produces it, per
  [`api/adr/0001-contract-first-api.md`](../api/adr/0001-contract-first-api.md).
* **How to use it:**
  ```bash
  task api:generate
  ```
  You will not normally run this directly — `task api:validate` and
  `task service:generate` both call it. Run it directly right after editing
  `api/openapi.yaml`, to see the generated diff before running full
  validation.

### `task api:validate`

* **What it does:** Runs `lint` → `compat` → `breaking-changes` → `generate`
  → a freshness check (`scripts/check-generated-fresh.sh`, which
  regenerates to a temp location and diffs against the committed
  `service/generated/`, failing on any difference).
* **Why it exists:** The single command that proves the API contract is
  valid, compatible, and that generated code hasn't drifted — the `api`
  boundary of `task validate`.
* **How to use it:**
  ```bash
  task api:validate
  ```

---

## `service:*` — Go implementation (`service/`)

### `task service:generate`

* **What it does:** Delegates directly to `api:generate` (via
  `task --taskfile ../api/Taskfile.yml --dir ../api generate`) — it does not
  duplicate the generation logic.
* **Why it exists:** Lets you run generation from the `service` namespace
  without needing to know it actually lives in `api` — while keeping exactly
  one implementation of "how do we generate", per the containerized-tooling
  and component-boundary ADRs.
* **How to use it:**
  ```bash
  task service:generate
  ```

### `task service:fmt`

* **What it does:** Runs `gofmt -l -w` over `service/cmd` and
  `service/internal` inside a pinned `golang:1.23-alpine` container, then
  fails if anything still needed formatting (i.e. it both fixes and
  verifies in one step).
* **Why it exists:** Keeps Go source canonically formatted without
  requiring a host Go install.
* **How to use it:**
  ```bash
  task service:fmt
  ```

### `task service:lint`

* **What it does:** Runs `golangci-lint run ./...` via a pinned
  `golangci/golangci-lint:v1.64.8` container.
* **Why it exists:** Standard Go static-analysis aggregation (unused code,
  common bug patterns, style) without a host toolchain.
* **How to use it:**
  ```bash
  task service:lint
  ```

### `task service:test`

* **What it does:** Runs `go test ./... -v -cover` via the pinned
  `golang:1.23-alpine` container.
* **Why it exists:** Runs the Go unit/handler test suite (status
  normalization, HTTP handler behavior, edge cases) with coverage reporting,
  with no host Go required.
* **How to use it:**
  ```bash
  task service:test
  ```

### `task service:build`

* **What it does:** Builds the `pulse` binary from `./cmd/pulse` into
  `service/bin/pulse`, via the pinned Go container.
* **Why it exists:** Proves the service actually compiles as a real binary
  (not just "tests pass") — also the same binary `service/Dockerfile`
  packages for Compose/Helm.
* **How to use it:**
  ```bash
  task service:build
  ```

### `task service:validate`

* **What it does:** Runs `generate` → generated-freshness check → `fmt` →
  `lint` → `test` → `build`, in that order.
* **Why it exists:** The single command that proves the Go service is
  correctly generated-from-contract, formatted, lint-clean, tested, and
  buildable — the `service` boundary of `task validate`.
* **How to use it:**
  ```bash
  task service:validate
  ```

---

## `web:*` — Web application (`web/`)

### `task web:generate`

* **What it does:** Runs `openapitools/openapi-generator-cli:v7.11.0`
  (pinned container, `typescript-fetch` generator) against
  `api/openapi.yaml`, writes output to `web/generated/`, then runs a
  post-generation fixup script (`tools/fix-generated-imports.mjs`) that adds
  explicit `.ts` import extensions, required by Node's strict ESM resolver.
* **Why it exists:** `web/generated/` must never be hand-written, same rule
  as `service/generated/` — this is the one mechanism that produces it.
* **How to use it:**
  ```bash
  task web:generate
  ```

### `task web:lint`

* **What it does:** Runs `tools/check.mjs` in a pinned `node:22-alpine`
  container — a syntax + import-resolution check (not full TypeScript type
  checking, since no `tsc` binary is installed; see `web/AGENTS.md` for why).
* **Why it exists:** Catches syntax errors and broken import paths in web
  source without adding an npm dependency (this repo's web toolchain is
  deliberately dependency-free — see `web/AGENTS.md`).
* **How to use it:**
  ```bash
  task web:lint
  ```

### `task web:test`

* **What it does:** Runs Node's built-in test runner
  (`node --experimental-transform-types --test src/**/*.test.ts`) in the
  pinned `node:22-alpine` container.
* **Why it exists:** Runs unit tests for pure logic (status mapping, summary
  computation) with zero test-framework dependency.
* **How to use it:**
  ```bash
  task web:test
  ```

### `task web:build`

* **What it does:** Runs `tools/build.mjs` (Node's built-in
  `stripTypeScriptTypes`) in the pinned container to compile `.ts` sources
  to static, browser-loadable JavaScript.
* **Why it exists:** Produces the actual deployable web asset — the same
  build `web/Dockerfile` packages into the nginx image used by Compose/Helm.
* **How to use it:**
  ```bash
  task web:build
  ```

### `task web:validate`

* **What it does:** Runs `generate` → generated-freshness check → `lint` →
  `test` → `build`.
* **Why it exists:** The single command that proves the web app is
  correctly generated-from-contract, lint-clean, tested, and buildable —
  the `web` boundary of `task validate`.
* **How to use it:**
  ```bash
  task web:validate
  ```

---

## `policies:*` — Deterministic engineering constraints (`policies/`)

Full rule descriptions and status live in
[`policies/README.md`](../policies/README.md); this section covers the task
mechanics only.

### `task policies:boundaries`

* **What it does:** Runs `scripts/check-component-boundaries.sh` — a
  grep-based scan of handwritten source (excluding `generated/`, `bin/`,
  `dist/`) for relative-import references that cross a top-level component
  boundary (e.g. `service/` code referencing `../web`).
* **Why it exists:** Enforces
  [`/adr/0001-component-boundaries-in-the-monorepo.md`](../adr/0001-component-boundaries-in-the-monorepo.md)
  mechanically rather than by convention alone.
* **How to use it:**
  ```bash
  task policies:boundaries
  ```

### `task policies:dependencies`

* **What it does:** Runs `scripts/check-dependency-governance.py`, which
  lists `service/go.mod`'s direct (non-`// indirect`) requires and
  `web/package.json`'s `dependencies`, and fails if any of them is missing
  from [`policies/allowed-dependencies.txt`](../policies/allowed-dependencies.txt).
* **Why it exists:** Makes adding a new direct third-party dependency a
  visible, reviewable diff instead of a silent addition — without
  introducing a license/supportability scanning platform.
* **How to use it:**
  ```bash
  task policies:dependencies
  ```
  Adding a real new dependency: add the line to `allowed-dependencies.txt`
  in the same commit that introduces the dependency.

### `task policies:containers`

* **What it does:** Runs `scripts/check-container-quality.sh`, which checks
  that `service/Dockerfile` and `web/Dockerfile` each end with a non-root
  `USER`, use an explicit (non-`latest`) base image tag, and scans
  Compose/Helm files for secret-like literal values.
* **Why it exists:** A basic, transparent container/deployment quality gate
  — this check is what found the real `web/Dockerfile` root-user bug fixed
  in Session 4 (see `policies/README.md` rule 5).
* **How to use it:**
  ```bash
  task policies:containers
  ```

### `task policies:validate`

* **What it does:** Runs `boundaries` → `dependencies` → `containers`.
  (Rules 1 and 2 — generated-code integrity and API lint/compat — are
  already enforced by `api:validate`/`service:validate`/`web:validate` and
  are deliberately not duplicated here.)
* **Why it exists:** The single command for "all currently-enforceable
  policy rules"; also the check run by the local pre-commit hook (alongside
  `structure:validate` and `api:breaking-changes`).
* **How to use it:**
  ```bash
  task policies:validate
  ```

---

## `compose:*` — Local integration environment (`deploy/compose/`)

### `task compose:up`

* **What it does:** Runs
  `docker compose up -d --build --wait` against
  `deploy/compose/docker-compose.yml`, building and starting
  `pulse-service`, `pulse-web`, and three deterministic fake HTTP targets
  (healthy/degraded/down), waiting for all Compose healthchecks to pass.
* **Why it exists:** The one command to stand up a fully local, deterministic
  integration environment — no external network dependency, no real
  internet endpoints.
* **How to use it:**
  ```bash
  task compose:up
  curl http://localhost:18080/api/v1/services
  open http://localhost:18000
  ```

### `task compose:down`

* **What it does:** Runs `docker compose down --volumes --remove-orphans`.
* **Why it exists:** Guarantees a clean teardown (containers, network,
  volumes) — no leftover `pulse-*` resources between runs.
* **How to use it:**
  ```bash
  task compose:down
  ```

### `task compose:logs`

* **What it does:** Runs `docker compose logs` against the stack.
* **Why it exists:** Quick diagnostic access to every service's logs
  without remembering individual container names.
* **How to use it:**
  ```bash
  task compose:logs
  ```

### `task compose:test`

* **What it does:** Runs `up`, then `deploy/compose/test.sh` (which verifies
  healthy/degraded/down status both directly against the API and via the
  `pulse-web` reverse-proxy), with `down` registered via Task's `defer:` so
  teardown always runs — even if the test script fails.
* **Why it exists:** Deterministic, automated proof that the whole system
  (service + web + fakes) behaves correctly when actually deployed together,
  not just in isolated unit tests.
* **How to use it:**
  ```bash
  task compose:test
  ```

### `task compose:validate`

* **What it does:** Runs `test`.
* **Why it exists:** Naming consistency with every other component's
  `*:validate` task — the `compose` boundary of `task validate`.
* **How to use it:**
  ```bash
  task compose:validate
  ```

---

## `helm:*` — Kubernetes packaging (`deploy/helm/`)

### `task helm:lint`

* **What it does:** Runs `helm lint pulse` via a pinned
  `alpine/helm:3.21.4` container against the chart in
  `deploy/helm/pulse/`.
* **Why it exists:** Standard Helm chart structural/syntax validation.
* **How to use it:**
  ```bash
  task helm:lint
  ```

### `task helm:template`

* **What it does:** Runs `helm template pulse pulse`, printing the fully
  rendered Kubernetes manifests to stdout.
* **Why it exists:** Lets you inspect exactly what Kubernetes objects the
  chart would create, without a live cluster.
* **How to use it:**
  ```bash
  task helm:template
  ```

### `task helm:template:check`

* **What it does:** Renders the chart (as above) to a temp file, then runs
  `deploy/helm/scripts/check-rendered-manifests.py` against it — a static
  check that readiness/liveness probes are present on every Deployment, no
  `Secret` resources exist, and the expected resource kinds/counts are
  present.
* **Why it exists:** Proves the rendered output is actually correct, not
  just that `helm template` didn't crash — still with no live cluster
  required.
* **How to use it:**
  ```bash
  task helm:template:check
  ```

### `task helm:validate`

* **What it does:** Runs `lint` → `template:check`.
* **Why it exists:** The single command that proves the Helm chart is valid
  and its rendered output meets baseline deployment-quality expectations —
  the `helm` boundary of `task validate`.
* **How to use it:**
  ```bash
  task helm:validate
  ```

---

## Vestigial: `deploy/Taskfile.yml`

`cd deploy && task validate` exists but only prints a pointer message
("Use root Taskfile: task compose:validate / task helm:validate") — it is
not included in the root `Taskfile.yml`'s `includes:`, so it is not reachable
as `task deploy:validate` from the repo root. It predates `compose:*` and
`helm:*` being split into their own namespaces early in this repository's
history and is not part of the `task validate` chain. Use
`task compose:validate` and `task helm:validate` directly instead.

---

## `make` targets (host bootstrap only — not part of `task validate`)

These exist purely to get `task` itself onto a fresh machine; per
[`/adr/0002-task-as-the-execution-contract.md`](../adr/0002-task-as-the-execution-contract.md)
they intentionally do none of the actual build/test/lint/generate/validate
work.

* **`make check-prereqs`** — verifies `git`, `docker` (and that the Docker
  daemon is actually running), and `task` are installed; exits non-zero
  listing anything missing.
* **`make bootstrap`** — installs `task` if missing (via Homebrew, then
  `go install`, whichever is available), then runs `check-prereqs`.
* **`make help`** — prints the two commands above with a one-line
  description each.
