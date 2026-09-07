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

## Status (Session 1)

Only structural discoverability is currently checked
(`task structure:validate` / `scripts/validate-structure.sh`): every `SKILL.md`
above exists with required `name`/`description` frontmatter, and points at a
Task interface that is itself defined (even if the underlying component isn't
implemented yet, in which case the target fails loudly with a clear message
rather than silently succeeding).

Full end-to-end skill execution evidence will be added in Sessions 2–4 as each
skill's underlying component becomes real.
