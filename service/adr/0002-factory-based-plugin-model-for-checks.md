# 2. Factory-Based Plugin Model for Service Checks

* Status: accepted
* Date: 2026-09-08

## Context and Problem Statement

Pulse currently performs exactly one kind of health check: an HTTP GET
against a service's URL, hardcoded inside `internal/monitoring.Checker`. As
soon as a second check mechanism is needed (e.g. TCP dial, DNS lookup, a
script-based probe), that hardcoded implementation has nowhere natural to
grow — either the `Checker` grows a type-switch that mixes check-specific
logic with status-normalization logic, or every new check type requires
editing code that has nothing to do with it.

Separately, a monitored service's check configuration (today: just a URL)
and a check's result (today: success/elapsed time/message, implicit in Go
struct fields) are not self-describing anywhere outside the Go source. There
is no mechanical way to answer "what configuration does this check type
need?" or "what does its result look like?" without reading
`monitoring.go`.

## Decision Drivers

* New check types must be addable without modifying the core check-running
  loop or its status-normalization logic (open/closed principle).
* A check type's configuration and result shape should be inspectable and
  validatable mechanically, not just documented in code comments.
* Keep the mechanism itself small and legible — one interface, one registry,
  one real implementation — not a general-purpose plugin framework.
* No new business functionality is required by this decision alone; the
  existing HTTP check must keep behaving identically once ported to the new
  model.

## Considered Options

1. **Status quo**: a single hardcoded HTTP check inside `Checker`.
2. **Inline type-switch**: add a `switch svc.Type` inside `Checker` with a
   case per check kind, all still living in `monitoring.go`.
3. **Factory-based plugin model**: define a `Plugin` interface
   (`Type() string`, `InputSchema() []byte`, `OutputSchema() []byte`,
   `Check(ctx, input) (output, error)`), a `Registry` that maps a type name
   to a `Factory` (constructor) and constructs `Plugin` instances by name,
   and require every plugin to ship JSON Schema documents describing its
   input and output shapes. Inputs are validated against `InputSchema()`
   before `Check()` runs; outputs are validated against `OutputSchema()`
   before the result is trusted.

## Decision Outcome

Chosen option: **3 — factory-based plugin model with self-describing,
validated I/O.**

* `internal/plugins` defines the `Plugin` interface, `Factory` type, and a
  thread-safe `Registry` (`Register`, `New`, `List`). The registry is the
  "factory" in the plugin-model sense: it holds named constructors and
  builds a fresh `Plugin` instance by type name.
* Every plugin type lives in its own package under `internal/plugins/<name>/`
  and must implement `Plugin` in full, including both JSON Schema methods.
  There is no default/fallback schema — an unschematized plugin is a defect.
* JSON Schema validation is real and enforced (via
  `github.com/santhosh-tekuri/jsonschema/v5`), not just documentation: a
  malformed input is rejected before `Check()` runs, and a malformed output
  is rejected before it reaches the caller.
* Registration is **explicit**, done once in `cmd/pulse/main.go`
  (`registry.Register("http", httpcheck.Factory)`) — no `init()`-based
  self-registration, consistent with this repository's "no hidden magic"
  convention (see e.g. the local pre-commit hook install step).
* `internal/monitoring.Checker` no longer performs HTTP requests itself. It
  resolves each service's `Type` to a `Plugin` via the registry, builds that
  plugin's input, calls `Check`, and applies Pulse's status-normalization
  rules (healthy/degraded/down thresholds) to the plugin's generic
  `{success, elapsedMs, message}` output. Status normalization stays
  centralized in the `Checker` — it is a Pulse-level concept (what counts as
  "healthy"), not a property of any one check mechanism.
* The first (and, as of this decision, only) plugin is `http`
  (`internal/plugins/httpcheck`), behaviorally identical to the check it
  replaces: a GET request with an explicit timeout, returning success,
  elapsed time, and an optional message.
* `config.Service` gains an optional `Type` (defaults to `"http"` when
  empty) and an optional `Config` (raw JSON; when empty, the `http` plugin's
  input is built from the existing `URL` field) — so the existing seed
  format used by `PULSE_SEED_SERVICES_JSON` and the Compose fake targets
  keeps working unchanged.

## Consequences

* Good: adding a new check type (e.g. `tcp`, `dns`) means adding a new
  package and one registration line — the `Checker`'s core loop and status
  rules never change.
* Good: a check type's configuration and result shape are mechanically
  inspectable (see the companion API ADR,
  [`api/adr/0002-plugin-type-discovery-endpoint.md`](../../api/adr/0002-plugin-type-discovery-endpoint.md))
  and mechanically validated at runtime, not just documented.
* Good: the `http` plugin's behavior is unchanged from the reader's/operator's
  perspective — this ADR changes internal structure, not observable check
  behavior, for the one check type that exists today.
* Bad: adds a small amount of indirection (registry, interface, schema
  validation) for what is, today, a single check type — accepted as the
  cost of making the extension point real rather than aspirational.
* Bad: adds one new third-party dependency
  (`github.com/santhosh-tekuri/jsonschema/v5`) — recorded as a visible,
  reviewable addition to `policies/allowed-dependencies.txt` per policies
  rule 4, not a silent addition.
* Future agents adding a new check type must implement the full `Plugin`
  interface (including both schemas) and register it explicitly; they must
  not special-case a new type inside `Checker`.
