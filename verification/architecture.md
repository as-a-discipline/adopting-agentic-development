# Architecture Verification

How to manually verify the repository's architectural guardrails. This document
covers cross-cutting checks. Component-specific verification evidence lives with
the component once it exists (e.g. `service` build/test evidence in Session 2,
`web`/`deploy` evidence in Session 3).

## 1. Generated files are identifiable

* Generated directories are named consistently (`service/generated/`,
  `web/generated/`) and, once populated, will carry a generated-file notice.
* `AGENTS.md` (root and component) states the generated-code rule explicitly.

Check: `grep -r "generated/" AGENTS.md */AGENTS.md`

## 2. Component boundaries exist

* `api/`, `service/`, `web/`, `deploy/`, `policies/` each have their own
  `AGENTS.md` describing their boundary.
* `AGENTS.md`/ADRs state that repository locality does not grant dependency
  access (see `/adr/0001-component-boundaries-in-the-monorepo.md`).

Check: `ls */AGENTS.md`

## 3. Root Taskfile delegates to component Taskfiles

* `Taskfile.yml` at root uses `includes:` for `api`, `service`, `web`, `compose`,
  `helm` — it does not reimplement their logic.

Check: `task --list-all` and inspect that each namespaced task (e.g. `api:lint`)
resolves to the corresponding component `Taskfile.yml`.

## 4. No host Go/Node dependency for generation/build

* `make check-prereqs` only checks for `git`, `docker`, `task`.
* Component Taskfiles are expected to invoke pinned container images for
  generation/build/lint once those components exist (see
  `/adr/0003-containerized-tooling.md`).

Check: `make check-prereqs` succeeds without Go/Node installed on the host.

## 5. Instructions are scoped hierarchically

* Root `AGENTS.md` contains only global invariants.
* Each component `AGENTS.md` contains only that component's local rules.
* `.github/copilot-instructions.md` is a small pointer, not a duplicate.

Check: read root `AGENTS.md` — it should be short enough to serve as a table of
contents, not an exhaustive manual.

## 6. Deterministic repository-structure validation

Run:

```text
task structure:validate
```

This validates required `AGENTS.md`/`SKILL.md` presence, skill frontmatter, and
required Taskfiles. It does **not** validate semantic correctness of any
instruction or skill — that requires the manual checks above plus, later, actual
usage evidence recorded in `verification/skills.md`.
