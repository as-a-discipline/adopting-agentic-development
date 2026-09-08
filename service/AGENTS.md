# AGENTS.md — service/

Implements the OpenAPI contract in Go.

## Rules

* Generated interfaces/types under `generated/` are not handwritten — implement the
  generated server interface from `internal/api/`, don't edit generated files.
* HTTP handlers stay thin; business behavior belongs in the appropriate `internal/`
  package (e.g. `internal/monitoring`).
* No arbitrary execution of operating-system commands from application code.
* No persistence — see [`adr/0001-no-persistence-in-baseline.md`](./adr/0001-no-persistence-in-baseline.md).
* Any outbound HTTP call (health checks against monitored services) must use an
  explicit timeout and honor request context.
* Use `task service:*` targets — do not invoke `go` directly outside of them.
* **All check mechanisms are plugins** — see
  [`adr/0002-factory-based-plugin-model-for-checks.md`](./adr/0002-factory-based-plugin-model-for-checks.md).
  A new check type (e.g. `tcp`, `dns`) is a new package under
  `internal/plugins/<name>/` that fully implements `plugins.Plugin`
  (`Type`, `InputSchema`, `OutputSchema`, `Check`) from `internal/plugins/`.
  Both schemas are required and must actually describe the plugin's I/O —
  `internal/monitoring.Checker` validates every plugin's input before
  calling `Check` and its output before trusting the result
  (`internal/plugins.Validate`), so an inaccurate schema breaks real checks,
  not just documentation.
* **No `init()`-based plugin self-registration.** Every plugin is registered
  explicitly by name in `cmd/pulse/main.go`
  (`registry.Register(httpcheck.Type, httpcheck.Factory)`) — consistent with
  this repo's "no hidden magic" convention. A plugin package must never
  register itself as a side effect of being imported.
* A plugin's `Check` output is expected to map onto the generic
  `{success, elapsedMs, statusCode, message}` shape that
  `internal/monitoring.Checker` understands for status normalization — see
  the ADR's Consequences section for this constraint's rationale and
  limitations.
* **Adding a new third-party Go dependency:** write the code that actually
  imports the package first, *then* run `go get <module>@<version>` and
  `go mod tidy` inside the pinned container (see `service:generate`'s
  container pattern) — running `go mod tidy` before anything imports the
  package will silently remove it from `go.mod` again, since tidy only
  keeps what's actually referenced. Record the new direct dependency in
  [`policies/allowed-dependencies.txt`](../policies/allowed-dependencies.txt)
  per policies rule 4 as part of the same change.

## Architecture Decisions

See [`adr/`](./adr/) for service-specific decisions. Cross-cutting decisions are in
the root [`/adr/`](../adr/).

## Task Interface

```text
task service:generate
task service:fmt
task service:lint
task service:test
task service:build
task service:validate
```
