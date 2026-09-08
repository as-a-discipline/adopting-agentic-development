# Task Reference: Every `task` Command in This Repository

This is the index for every task defined across this repository's
Taskfiles. Per
[`/adr/0002-task-as-the-execution-contract.md`](../adr/0002-task-as-the-execution-contract.md),
`task <name>` is the **only** supported way to build, lint, test, generate,
or validate anything here. Each task below has its own standalone document
(what it does, why it exists, how to use it) — the same format used for
[`/docs/api-breaking-changes.md`](./api-breaking-changes.md).

Run `task --list-all` from the repo root at any time to see this same list
live, generated directly from the Taskfiles (so it can never drift from
what's actually runnable).

For the two-sentence version of the whole validation pipeline, see the root
[README's "Validate Everything" section](../README.md#validate-everything).

## Root-level tasks

* [`task` (no arguments) / `task default`](./tasks/default.md)
* [`task validate`](./tasks/validate.md)
* [`task structure:validate`](./tasks/structure-validate.md)
* [`task hooks:install`](./tasks/hooks-install.md)

## `api:*` — OpenAPI contract (`api/`)

* [`task api:lint`](./tasks/api-lint.md)
* [`task api:compat`](./tasks/api-compat.md)
* [`task api:breaking-changes`](./tasks/api-breaking-changes.md) — see also
  the full guide: [`/docs/api-breaking-changes.md`](./api-breaking-changes.md)
* [`task api:generate`](./tasks/api-generate.md)
* [`task api:validate`](./tasks/api-validate.md)

## `service:*` — Go implementation (`service/`)

* [`task service:generate`](./tasks/service-generate.md)
* [`task service:fmt`](./tasks/service-fmt.md)
* [`task service:lint`](./tasks/service-lint.md)
* [`task service:test`](./tasks/service-test.md)
* [`task service:build`](./tasks/service-build.md)
* [`task service:validate`](./tasks/service-validate.md)

## `web:*` — Web application (`web/`)

* [`task web:generate`](./tasks/web-generate.md)
* [`task web:lint`](./tasks/web-lint.md)
* [`task web:test`](./tasks/web-test.md)
* [`task web:build`](./tasks/web-build.md)
* [`task web:validate`](./tasks/web-validate.md)

## `policies:*` — Deterministic engineering constraints (`policies/`)

Full rule descriptions and status live in
[`policies/README.md`](../policies/README.md); the linked docs below cover
task mechanics only.

* [`task policies:boundaries`](./tasks/policies-boundaries.md)
* [`task policies:dependencies`](./tasks/policies-dependencies.md)
* [`task policies:containers`](./tasks/policies-containers.md)
* [`task policies:validate`](./tasks/policies-validate.md)

## `compose:*` — Local integration environment (`deploy/compose/`)

* [`task compose:up`](./tasks/compose-up.md)
* [`task compose:down`](./tasks/compose-down.md)
* [`task compose:logs`](./tasks/compose-logs.md)
* [`task compose:test`](./tasks/compose-test.md)
* [`task compose:validate`](./tasks/compose-validate.md)

## `helm:*` — Kubernetes packaging (`deploy/helm/`)

* [`task helm:lint`](./tasks/helm-lint.md)
* [`task helm:template`](./tasks/helm-template.md)
* [`task helm:template:check`](./tasks/helm-template-check.md)
* [`task helm:validate`](./tasks/helm-validate.md)

## Vestigial

* [`deploy/Taskfile.yml`'s `validate`](./tasks/deploy-validate-vestigial.md) —
  not reachable from the repo root, not part of `task validate`

## `make` targets (host bootstrap only — not part of `task validate`)

These exist purely to get `task` itself onto a fresh machine; they
intentionally do none of the actual build/test/lint/generate/validate work.

* [`make check-prereqs`](./tasks/make-check-prereqs.md)
* [`make bootstrap`](./tasks/make-bootstrap.md)
* [`make help`](./tasks/make-help.md)
