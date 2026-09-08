# Detecting Breaking API Changes: `task api:breaking-changes`

This document explains what `task api:breaking-changes` does, how it's wired
into `task validate` and the local git pre-commit hook, and provides real,
captured proof that it detects a genuine breaking change and blocks both
validation and a commit.

> Note on placement: most cross-cutting docs in this repository live in root
> `/verification/` (see that directory's `AGENTS.md`-adjacent convention).
> This guide lives under `/docs/` at the requester's explicit direction, as a
> standalone, task-specific how-to rather than an architecture-verification
> document.

## What it does

`task api:breaking-changes` compares `api/openapi.yaml` between:

* the **current working tree** (right side — including uncommitted changes), and
* a **comparison reference** (left side), resolved as:
  1. an explicit ref you pass yourself (`task api:breaking-changes -- <ref>`), or
  2. the current branch's **upstream tracking branch** — i.e. its last pushed
     commit (`git rev-parse --abbrev-ref --symbolic-full-name @{u}`), or
  3. if there's no upstream (e.g. the branch has never been pushed), the
     **target mainline branch**: `origin/main`, then `origin/master`, then
     local `main`, then local `master` — whichever resolves first.

If the resolved ref doesn't have `api/openapi.yaml` at all (e.g. a brand new
contract, or the ref predates the API), the task prints an explanatory
message and passes — there's nothing to compare against yet.

## Choosing a detection tool

Two pinned, containerized tools are supported, selected with `TOOL=`:

```bash
task api:breaking-changes                          # TOOL=oasdiff (default)
task api:breaking-changes TOOL=openapi-changes      # pb33f/openapi-changes
task api:breaking-changes -- main                   # explicit ref override
task api:breaking-changes TOOL=openapi-changes -- main
```

| Tool | Image | How it decides "breaking" |
|---|---|---|
| [`oasdiff`](https://github.com/oasdiff/oasdiff) (default) | `tufin/oasdiff:v1.31.0` | Runs `oasdiff breaking --fail-on ERR`; oasdiff's own built-in breaking-change rules; the task fails if oasdiff reports any `ERR`-level finding. |
| [`pb33f/openapi-changes`](https://github.com/pb33f/openapi-changes) | `pb33f/openapi-changes:v0.2.11` | Runs `openapi-changes report --reproducible`, which classifies every detected change with a `breaking: true/false` field (via [libopenapi's configurable breaking-change rules](https://pb33f.io/libopenapi/what-changed/#configurable-breaking-change-rules)); `scripts/parse-openapi-changes-report.py` reads that JSON and fails if any change has `breaking: true`. |

Both are pinned, run as containers (no host install required), and are
genuinely independent implementations — useful for cross-checking a
suspicious result.

### Why a custom JSON parser for `openapi-changes`?

`openapi-changes` does not (as of `v0.2.11`) document a CLI flag that gates
the exit code on **breaking-only** changes. Its only documented gating flag,
`--error-on-diff` (on the `summary` command), fails on **any** diff, breaking
or not — too strict for our purpose. Instead, `check-api-breaking-changes.sh`
runs the machine-readable `report` command and
`scripts/parse-openapi-changes-report.py` applies the breaking-only gate
itself, based on the JSON shape verified by running the tool directly against
this repository (see [Evidence](#evidence) below):

```json
{
  "reportSummary": {"paths": {"totalChanges": 1, "breakingChanges": 1}, "...": "..."},
  "changes": [
    {"breaking": true, "path": "$.paths./services/{id}/check", "changeText": "object_removed", "...": "..."}
  ]
}
```

## Where it's wired in

* **`task api:validate`** (and therefore the canonical **`task validate`**)
  now runs `breaking-changes` as a real step, between `compat` and
  `generate`. A breaking change vs. the last pushed commit / mainline branch
  fails `task validate` — the same gate that structure, lint, tests, and
  policies go through.
* **The local git pre-commit hook** (`scripts/git-hooks/pre-commit`, opt-in
  via `task hooks:install` — see the root
  [README's Guardrails section](../README.md#guardrails)) now runs
  `task api:breaking-changes` directly, in addition to its existing fast
  subset (`task structure:validate` + `task policies:validate`). This check
  is deliberately included in the pre-commit hook (unlike the rest of
  `api`/`service`/`web`/`compose`/`helm` validation, which is too slow for a
  pre-commit gate) because it's fast — no container build, no test suite, no
  Compose stack — just a git-history read plus a single containerized diff
  tool run. **A breaking API change blocks the commit** unless you
  deliberately bypass with `git commit --no-verify`.

This means a breaking API change is now caught at three points, in
increasing order of "too late to easily fix":

1. **Before you commit** — the pre-commit hook.
2. **Any time you run `task validate`** — locally, before pushing.
3. **Explicitly, on demand** — `task api:breaking-changes` itself, useful for
   checking a change against a specific ref before you've even finished it.

## Note on `task api:compat` (the pre-existing static-baseline check)

`task api:compat` (added in Session 2, part of `task api:validate`) checks
the contract against a *stored baseline snapshot*
(`api/openapi.baseline.yaml`) rather than git history. It answers a
different question ("did we break compatibility with the last time we
explicitly agreed the contract was good?") than `api:breaking-changes`
("did we break compatibility with what's already pushed / mainline?"). Both
now run as part of `task api:validate`, and — as the evidence below shows —
both correctly failed on the same real breaking change, from two
independent code paths.

**A real, pre-existing bug was found and fixed while building
`api:breaking-changes`:** `oasdiff`'s `breaking` subcommand does not fail on
its own — it only returns a non-zero exit code when `--fail-on ERR` (or
`WARN`) is explicitly passed. `task api:compat` had never passed this flag,
meaning it has silently never failed the build on a breaking change since
Session 2, despite `policies/README.md` documenting it as "enforced". This
has been fixed (`--fail-on ERR` added to `api:compat`'s oasdiff invocation)
and is included in the same evidence run below.

## Evidence

All logs below are real command output, captured in this repository, proving
the mechanism actually works — not narrated. Full logs are in
[`api-breaking-changes-evidence/`](./api-breaking-changes-evidence/).

### Test setup

1. **A real breaking change** was introduced: the entire
   `POST /services/{id}/check` operation was deleted from
   `api/openapi.yaml` (removing a public endpoint without deprecation — an
   unambiguous, realistic breaking change).
2. Because this repository's `in-progress` branch had never been pushed at
   the time of this test (no upstream tracking branch, and `origin/main` /
   local `main` are still just the initial README+LICENSE commit — they
   predate the API contract entirely), the default reference-resolution
   logic would otherwise report "no API contract found at ref — nothing to
   compare" (a real, correctly-handled case, but not useful for proving
   breaking-change *detection*). To exercise the actual default-resolution
   code path end-to-end, the local record of `refs/remotes/origin/main` was
   temporarily pointed (via `git update-ref`, a purely local operation — no
   `git push`, no network call, nothing sent anywhere) at commit `811d8cb`
   (the commit that originally introduced `api/openapi.yaml`, still real
   history in this repository), then restored to its real value
   (`ccaebc8`) immediately after capturing this evidence. This means the
   "comparing 'origin/main' ... -> working tree" lines below reflect exactly
   what would happen for real once this branch (or any branch) is actually
   pushed and diverges from a mainline that already has the contract.

### Scenario 1 — `task api:breaking-changes` (default: `TOOL=oasdiff`)

[`01-default-oasdiff.log`](./api-breaking-changes-evidence/01-default-oasdiff.log)

```text
$ task api:breaking-changes
task: [api:breaking-changes] ...
== api:breaking-changes — comparing 'origin/main' (target mainline branch (no upstream configured for this branch)) -> working tree ==
tool: oasdiff
1 changes: 1 error, 0 warning, 0 info
error	[api-path-removed-without-deprecation] at /old/old-openapi.yaml
	in API POST /services/{id}/check
		api path removed without deprecation

task: Failed to run task "api:breaking-changes": exit status 1
EXIT=201
```

`oasdiff` correctly identifies the removed path as breaking and the task
fails (Task wraps the underlying `exit 1` as its own `201`).

### Scenario 2 — `task api:breaking-changes TOOL=openapi-changes`

[`02-default-openapi-changes.log`](./api-breaking-changes-evidence/02-default-openapi-changes.log)

```text
$ task api:breaking-changes TOOL=openapi-changes
== api:breaking-changes — comparing 'origin/main' ... -> working tree ==
tool: openapi-changes
Total changes: 1; breaking: 1

Breaking changes:
  - $.paths./services/{id}/check: object_removed ('/services/{id}/check' -> None)

FAIL: 1 breaking change(s) detected.
task: Failed to run task "api:breaking-changes": exit status 1
EXIT=201
```

The independent `pb33f/openapi-changes` engine detects the same removal and
the task fails, via the custom JSON-parsing gate described above.

### Scenario 3 — `task api:validate` fails

[`03-api-validate-fails.log`](./api-breaking-changes-evidence/03-api-validate-fails.log)

```text
$ task api:validate
task: [api:lint] ...
Woohoo! Your API description is valid. 🎉

task: [api:compat] ...
1 changes: 1 error, 0 warning, 0 info
error	[api-path-removed-without-deprecation] at /spec/openapi.baseline.yaml
	in API POST /services/{id}/check
		api path removed without deprecation

task: Failed to run task "api:validate": task: Failed to run task "api:compat": exit status 1
EXIT=201
```

`api:lint` still passes (the contract is still structurally valid YAML/
OpenAPI — removing an operation doesn't make the document invalid), but
`api:compat` (the static-baseline check, now fixed with `--fail-on ERR`)
correctly fails the chain before `breaking-changes` even runs — proving the
`--fail-on ERR` fix works, and that `api:validate` stops at the first real
problem it finds.

### Scenario 4 — full `task validate` fails

[`04-full-validate-fails.log`](./api-breaking-changes-evidence/04-full-validate-fails.log)

```text
$ task validate
Repository structure: OK
task: [api:lint] ...
Woohoo! Your API description is valid. 🎉

task: [api:compat] ...
1 changes: 1 error, 0 warning, 0 info
error	[api-path-removed-without-deprecation] ...

task: Failed to run task "validate": task: Failed to run task "api:validate": task: Failed to run task "api:compat": exit status 1
EXIT=201
```

The canonical, repository-wide `task validate` gate — the same one described
in the root README — stops at this same boundary. A breaking API change
never reaches service/web/policies/compose/helm validation; it fails fast.

### Scenario 5 — a real `git commit` is blocked

[`05-commit-blocked.log`](./api-breaking-changes-evidence/05-commit-blocked.log)

```text
$ git add api/openapi.yaml
$ git commit -m 'test: attempt to commit a breaking API change (should be blocked)'
== pre-commit: task structure:validate ==
Repository structure: OK

== pre-commit: task policies:validate ==
OK: container/deployment quality checks passed.

== pre-commit: task api:breaking-changes ==
tool: oasdiff
1 changes: 1 error, 0 warning, 0 info
error	[api-path-removed-without-deprecation] ...

task: Failed to run task "api:breaking-changes": exit status 1
EXIT=1
```

This is a **real `git commit` invocation** (not a simulation) against the
actual installed `.git/hooks/pre-commit` — the commit was refused, and
`git log` immediately afterward confirmed `HEAD` had not moved: the breaking
change was never committed.

### Cleanup

After capturing this evidence, the breaking change was reverted
(`api/openapi.yaml` restored to its original content, confirmed via `diff`
against a pre-test backup) and the temporarily-modified
`refs/remotes/origin/main` local record was restored to its real value
(confirmed via `git rev-parse origin/main`). No commit containing the
breaking change exists in this repository's history; no remote was ever
contacted.

## Summary

* `task api:breaking-changes` works with both supported tools and correctly
  fails on a genuine breaking change (a removed, non-deprecated operation).
* It is now a real step in `task api:validate` / `task validate`.
* It is now enforced by the local pre-commit hook — a real `git commit`
  containing a breaking API change is refused before it ever reaches
  history, without `--no-verify`.
