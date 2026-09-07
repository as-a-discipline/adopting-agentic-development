# Skill Verification Matrix

Cross-cutting record of every skill defined in `.agents/skills/`. Component-specific
verification evidence (e.g. actual `task service:test` output once the service
exists) is recorded alongside that component as it's built in later sessions; this
matrix tracks discoverability and the intended trigger/interface for the whole
skill set.

| Skill | Intended trigger | Task interface used | Evidence of success |
| ----- | ----------------- | -------------------- | -------------------- |
| `api-contract` | Changing/adding an API operation or schema; fixing an API lint/compat failure | `task api:lint`, `task api:generate`, `task api:validate` | `task api:validate` exits 0; generated files match the contract |
| `service-development` | Implementing/fixing Go service behavior against an existing contract | `task service:generate`, `task service:fmt`, `task service:lint`, `task service:test`, `task service:build`, `task service:validate` | `task service:validate` exits 0; new tests include a negative case |
| `web-development` | Implementing/fixing the web UI or wiring a new API operation | `task web:generate`, `task web:lint`, `task web:test`, `task web:build`, `task web:validate` | `task web:validate` exits 0; no hand-duplicated generated types |
| `compose-integration` | Running/debugging end-to-end integration checks | `task compose:up`, `task compose:test`, `task compose:logs`, `task compose:down`, `task compose:validate` | `task compose:validate` exits 0; clean teardown |
| `helm-validation` | Changing/validating the Helm chart | `task helm:lint`, `task helm:template`, `task helm:validate` | `task helm:validate` exits 0; valid probes, no secrets |
| `repository-validation` | Running full validation or triaging a `task validate` failure | `task validate`, `task structure:validate` | `task validate` exits 0; fix scoped to the failing component |

## Status (Session 4): full verification evidence

For each skill below, six things were checked:

1. **Structurally discoverable** — `SKILL.md` exists with valid `name`/
   `description` frontmatter (enforced by `task structure:validate`).
2. **Trigger clarity** — the `description` states concretely when to use the
   skill (not just a vague topic label).
3. **Points to correct scoped instructions** — the skill's own text names the
   right component `AGENTS.md` to read first.
4. **Task interface named matches reality** — every `task <x>` command the
   skill names actually exists (cross-checked against the real `Taskfile.yml`
   files, not just the skill's own claims).
5. **Commands actually exist and run** — executed directly, not just grepped.
6. **Workflow executed manually** — the documented sequence was actually run
   (in some cases as part of the real work performed in Sessions 2–4, in
   others as a fresh check during Session 4), with results recorded below.

| Skill | 1 | 2 | 3 | 4 | 5/6 — manual execution evidence |
| ----- | - | - | - | - | -------------------------------- |
| `api-contract` | ✅ | ✅ | ✅ (`api/AGENTS.md`, `api/adr/`) | ✅ | Used throughout Session 2 to author `api/openapi.yaml`; `task api:lint`, `task api:generate`, `task api:validate` all run repeatedly and pass (latest: `task validate` output, Session 4, `api:validate` → "OK: 'service/generated' matches what regeneration produces."). |
| `service-development` | ✅ | ✅ | ✅ (`service/AGENTS.md`, `service/adr/`) | ✅ | Used in Session 2 to build the Go service and again in Session 3 to add the `PULSE_SEED_SERVICES_JSON` config override (with a new negative-case test: falls back to default seed on invalid JSON). `task service:validate` passes (fmt/lint/test/build all green, 96.4% coverage on `internal/monitoring`). |
| `web-development` | ✅ | ✅ | ✅ (`web/AGENTS.md`) | ✅ | Used in Session 3 to build the whole web app + generated-client toolchain. `task web:validate` passes (generate/freshness/lint/test/build all green); verified with a real smoke test against the live Go service. |
| `compose-integration` | ✅ | ✅ | ✅ (`deploy/compose/AGENTS.md`) | ✅ | Used in Session 3 to build and debug the Compose stack — including a real failure (pulse-web healthcheck) diagnosed via `task compose:logs` + manual `docker compose exec`, fixed, then re-verified with `task compose:test`. Re-run again in Session 4 after switching the web runtime image; `task compose:down` confirmed to leave zero dangling containers/networks each time. |
| `helm-validation` | ✅ | ✅ | ✅ (`deploy/helm/AGENTS.md`) | ✅ | Used in Session 3 to build the chart and in Session 4 to update the web port after the non-root nginx fix. `task helm:validate` passes (`helm lint` clean; rendered manifests statically checked for probes/no-secrets via `deploy/helm/scripts/check-rendered-manifests.py`). |
| `repository-validation` | ✅ | ✅ | ✅ (root `AGENTS.md`) | ✅ | Used at the start of every session (fresh `task validate` baseline check) and to triage the Session 3 `pulse-web` failure (identified the failing boundary from `task validate` output, then loaded only `deploy/compose/AGENTS.md` + the `compose-integration` skill to fix it — not every component's context). Session 4: `task validate` passes end-to-end through all seven boundaries (structure → api → service → web → policies → compose → helm). |

### What this does and does not prove

* It proves: every skill is discoverable, names a real and currently-working
  Task interface, and was actually exercised at least once to do real work (not
  hypothetical instructions that were never run).
* It does **not** prove that a model will always autonomously select the
  correct skill for a given request. Skill *discovery* (the file exists and is
  syntactically valid) and skill *instruction correctness* (the documented
  workflow is accurate and its commands work) are both deterministic and
  checked here. *Model skill-selection behavior* — whether an agent picks the
  right skill unprompted for an arbitrary phrasing of a request — is a
  probabilistic property of the model and prompt, not something a structural
  test can guarantee. Treat this matrix as proof of a well-formed, working
  skill library, not as proof of perfect autonomous routing.

## How to verify with Copilot CLI

Useful interactive checks, run from within this repository:

* **List available skills**: run `/skills` to browse/manage the skills Copilot
  CLI has discovered, including every entry under `.agents/skills/`.
* **Inspect what's currently loaded**: run `/env` — it reports loaded
  instructions, skills, agents, and other environment context, so you can
  confirm a specific skill (e.g. `service-development`) was actually
  discovered for this repository.
* **Encourage a specific skill's use**: reference the skill's trigger
  scenario directly in a prompt (e.g. "the compose:test is failing, debug it"
  should surface `compose-integration`); you can also explicitly name the
  skill in your prompt (e.g. "use the helm-validation skill to check the
  chart") to make the intent unambiguous. Whether the model chooses to load a
  skill unprompted is model behavior, not a deterministic guarantee (see
  above).

## How to verify with OpenCode

OpenCode natively recognizes the same `.agents/skills/<name>/SKILL.md` layout
this repository already uses (no duplication or shim required) — per OpenCode's
own skill-discovery docs, project-local skills are walked up from the current
working directory to the git worktree root, matching `.agents/skills/*/SKILL.md`
among its other supported locations.

To verify with OpenCode:

* Start an OpenCode session from the repository root and ask what skills are
  available — OpenCode surfaces discovered skills (name + description) to the
  model via its built-in `skill` tool, so the model can report back what it
  sees.
* Confirm each `SKILL.md`'s frontmatter satisfies OpenCode's naming rules
  (lowercase alphanumeric with single-hyphen separators, 1–64 characters,
  `name` matching its containing directory) — the same rules `task
  structure:validate` already checks a subset of (presence of `name`/
  `description`); all six skills in this repo already conform.
* If a skill doesn't appear, OpenCode's own troubleshooting steps apply:
  confirm the file is named exactly `SKILL.md`, confirm required frontmatter
  fields are present, and confirm no naming collision exists across the
  supported skill directories.

