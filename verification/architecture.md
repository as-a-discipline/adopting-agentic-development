# Architecture Verification

How to manually verify the repository's architectural guardrails. This document
covers cross-cutting checks. Component-specific verification evidence lives with
the component (e.g. `service` build/test evidence, `web`/`deploy` evidence) —
see each component's own `AGENTS.md` and Taskfile output.

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

* `Taskfile.yml` at root uses `includes:` for `api`, `service`, `web`,
  `policies`, `compose`, `helm` — it does not reimplement their logic.

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
instruction or skill — that requires the manual checks above plus actual usage
evidence recorded in `verification/skills.md`.

## 7. Instruction hierarchy in practice — a worked example

The design intent (see root `AGENTS.md`'s Context Rule) is that an agent never
needs one giant static instructions file. Instead, context is composed from
three layers plus the specific request:

```text
root AGENTS.md
      +
service/AGENTS.md
      +
service-development skill
      +
specific story/task
```

Here is what that composition actually looks like for a concrete task —
**"Add a `tags: []string` field to the `MonitoredService` API schema and thread
it through the Go service and the web UI"** — showing what each layer
contributes and, critically, what it does *not* need to contain because a
different layer already owns it:

* **Root `AGENTS.md` contributes** (global invariants, loaded once,
  independent of which component the task touches):
  * The repository architecture map: `api` owns the contract, `service`/`web`
    only consume generated artifacts, so the agent knows the schema change
    starts in `api/openapi.yaml`, not in `service/generated/` or
    `web/generated/` directly.
  * The Generated-Code Rule: never hand-edit `generated/` — regenerate via
    `task *:generate` after the contract changes.
  * The Execution Rule: use `task` targets, don't invent ad hoc shell
    commands.
  * The Forbidden Actions list: e.g. don't introduce persistence just because
    a new field is being added.
  * It does **not** need to know anything about Go syntax, TypeScript syntax,
    or this specific field — that's out of scope for a global-invariants file.

* **`service/AGENTS.md` contributes** (loaded only because this task touches
  `service/`):
  * That generated interfaces live under `generated/` and must be implemented
    from `internal/api/`, not edited directly.
  * That business behavior belongs in `internal/monitoring`, not the HTTP
    handler layer.
  * The exact `task service:*` interface to use.
  * It does **not** repeat the root's architecture map or forbidden-actions
    list — those are already in scope from the root layer.

* **The `service-development` skill contributes** (loaded on-demand because
  the request matches its trigger — "add a Go-side feature that already has an
  API contract"):
  * The concrete step order: contract changed → run `task service:generate`
    first, then implement in handwritten code, then test/lint/validate in that
    order, then interpret failures rather than iterating blindly.
  * The evidence-of-success bar: `task service:validate` exits 0, and new
    behavior has a test including a negative case.
  * It does **not** contain the architecture rules or the Task command
    reference itself — it composes with, rather than duplicates, the two
    layers above (and would separately trigger `api-contract` for the schema
    edit, and `web-development` for the frontend half of this same task).

* **The specific task** contributes the only thing none of the above layers
  could: that this particular change is a `tags` field, on this particular
  schema, right now.

This is what "context is scoped" (root `AGENTS.md`'s design principle) means in
practice: each layer is loaded because it's relevant, contains only what it
uniquely knows, and none of the three static layers had to be rewritten to
accommodate this specific task — only the task description itself was new.

**To verify this yourself:** read root `AGENTS.md`, then `service/AGENTS.md`,
then `.agents/skills/service-development/SKILL.md`, in that order, and confirm
none of them restate content already covered by an earlier layer (e.g.
`service/AGENTS.md` doesn't re-explain what `service/` is for at a
whole-repository level; the skill doesn't redefine what `task service:test`
does versus just naming it).
