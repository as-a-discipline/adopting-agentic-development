# Feature: Add TCP Monitoring as a New Pluggable Check Type

## Summary

Extend Pulse so that monitored services are not limited to HTTP health checks.

Introduce a pluggable monitoring architecture that allows Pulse to support multiple check types, and deliver **TCP port monitoring** as the first additional implementation alongside the existing HTTP check.

A user should be able to configure a monitored service to use either an HTTP check or a TCP check. Pulse must execute the appropriate check implementation, normalize the result into the existing service-status model, and display the result consistently through the API and web application.

This Epic establishes the extension pattern that future monitoring types can follow without requiring the monitoring core to be redesigned.

---

# Business Outcome

Pulse currently answers service-health questions using HTTP checks only.

That limits monitoring to services that expose an HTTP endpoint.

Many infrastructure services are better validated by confirming that a TCP endpoint is reachable, for example:

* databases
* message brokers
* SSH endpoints
* proxies
* application ports
* infrastructure appliances

After this Epic, Pulse should support monitoring both:

* HTTP endpoints
* TCP endpoints

through a common monitoring model.

The user should not need to understand which internal implementation performs the check.

---

# User Value

A user can configure a monitored service with the check type appropriate to that service.

Examples:

**HTTP**

`https://payments.example/health`

**TCP**

`database.example:5432`

Pulse should present both through the same status experience:

* healthy
* degraded
* down
* unknown

The addition of another check type should not require separate monitoring screens or a separate product workflow.

---

# Epic Scope

This Epic includes:

* introducing an explicit monitoring-check abstraction
* preserving the existing HTTP monitoring capability through that abstraction
* implementing TCP port monitoring as a second check type
* extending the API contract to represent check type and check configuration
* generating updated Go and TypeScript API artifacts
* updating the Go service to select the appropriate check implementation
* updating the web UI to represent the selected monitoring type
* adding deterministic TCP test targets to Docker Compose
* updating Helm configuration where required
* unit, contract, integration, and regression testing
* documentation of the plugin/check-type extension pattern

---

# Out of Scope

This Epic does not include:

* ping / ICMP checks
* DNS checks
* TLS/certificate checks
* traceroute
* UDP monitoring
* arbitrary command execution
* custom user-uploaded plugins
* dynamically loaded binary plugins
* third-party plugin marketplaces
* scripting support
* persistence
* historical metrics
* alerting
* authentication or authorization changes
* automatic discovery of monitoring types

The term **plugin** refers to an internal extension pattern in the monitoring service, not arbitrary executable third-party code.

---

# Functional Requirements

## FR1 — Monitoring Type

Each monitored service must explicitly identify the type of check used to evaluate its health.

Initially supported values are:

* `http`
* `tcp`

The API contract must define the supported values.

Existing HTTP-based services must continue to function after the change.

Where practical, compatibility with existing service definitions should be preserved.

---

## FR2 — Monitoring Check Abstraction

The monitoring service must no longer contain HTTP-specific behavior as the only implementation path.

Introduce an internal abstraction for a monitoring check.

Conceptually:

`Check(ctx, configuration) → Result`

or equivalent.

Each check implementation must return a normalized result understood by the monitoring domain.

The abstraction should support:

* HTTP check implementation
* TCP check implementation

without requiring monitoring-domain callers to know implementation details.

---

## FR3 — HTTP Check Plugin

The existing HTTP monitoring behavior must be migrated behind the new check abstraction.

Existing HTTP behavior must remain functionally equivalent unless the current implementation contains behavior that must change to satisfy the new architecture.

Existing HTTP regression tests must continue to pass.

---

## FR4 — TCP Check Plugin

Pulse must support checking whether a TCP connection can be established to a configured host and port.

A TCP monitor must support configuration equivalent to:

* hostname or IP address
* TCP port

The check must:

1. attempt to establish a TCP connection
2. use a bounded timeout
3. honor cancellation through context where applicable
4. close the connection cleanly after a successful check
5. normalize the result into the common monitoring result

A successful TCP connection should result in a healthy service.

A timeout, refusal, or connection failure should result in an appropriate unhealthy/down result.

---

## FR5 — Common Check Result

HTTP and TCP check implementations must produce a shared normalized result.

The result should support the information required by the existing monitoring domain, including:

* status
* check timestamp
* response/duration measurement
* human-readable message where appropriate

Transport-specific details should not leak unnecessarily into the common monitoring model.

---

## FR6 — API Contract

The OpenAPI contract must represent the monitoring type and the configuration required by that type.

The contract should make invalid configurations difficult to represent.

For example:

* an HTTP check requires a URL
* a TCP check requires a host and port

Prefer explicit typed schemas over a generic unstructured configuration map.

The exact OpenAPI modeling approach is an implementation decision, but it must remain understandable and generate usable Go and TypeScript types.

---

## FR7 — API Compatibility

Existing API behavior must remain backward compatible where practical.

If existing services are represented as HTTP monitors implicitly, the implementation should preserve a reasonable compatibility path rather than forcing unrelated breaking API changes.

Any unavoidable contract change must be explicitly identified during planning and must pass the repository's API compatibility process.

The agent or implementation team must not silently introduce a breaking API change for architectural convenience.

---

## FR8 — Web UI

The Pulse UI must display the check type used by each monitored service.

Users must be able to distinguish, at minimum:

* HTTP monitored service
* TCP monitored service

The UI should continue to display the normalized health information consistently regardless of check type.

The UI does not need a complex configuration editor for this Epic unless one already exists.

Seed/configured TCP services are sufficient to demonstrate the capability.

---

## FR9 — Deterministic Integration Environment

Docker Compose must include deterministic targets for TCP monitoring.

At minimum, provide:

* one reachable TCP target
* one unavailable/refused target or deterministic equivalent

Tests must not depend on external network resources.

The integration environment must allow validation of both HTTP and TCP monitoring from a clean local environment.

---

# Non-Functional Requirements

## NFR1 — Extensibility

Adding a future check type should not require significant modification to the monitoring orchestration logic.

A future implementation should primarily require:

* a new check implementation
* check-specific configuration
* registration or selection logic
* appropriate tests
* API/UI changes required by that check

Do not over-generalize the architecture into a dynamic plugin framework.

The goal is a clean extension seam, not a plugin platform.

---

## NFR2 — Bounded Execution

Every check implementation must have bounded execution.

HTTP and TCP checks must use explicit timeouts.

A failed target must not cause the monitoring loop or request to hang indefinitely.

---

## NFR3 — No Arbitrary Execution

Monitoring plugins must not execute user-defined shell commands.

The plugin abstraction must represent typed application capabilities, not arbitrary executable behavior.

No new shell execution mechanism should be introduced for TCP monitoring.

---

## NFR4 — Isolation of Transport Logic

HTTP-specific and TCP-specific behavior must remain inside their respective implementations.

The monitoring orchestration layer should depend on the common check abstraction.

Transport-specific connection logic must not spread into API handlers or unrelated domain packages.

---

## NFR5 — Testability

Each check implementation must be independently testable.

Tests must cover both positive and negative behavior.

The monitoring orchestration logic must be testable without requiring real external network services.

---

## NFR6 — Performance

A TCP health check should add minimal overhead beyond the connection attempt itself.

No unnecessary retries, polling loops, or background workers should be introduced.

Timeout values should be reasonable and configurable through the existing configuration pattern where appropriate.

---

## NFR7 — Portability

The TCP implementation should use portable Go networking APIs.

It must not depend on host utilities such as:

* `nc`
* `telnet`
* `nmap`
* shell scripts

The functionality should operate consistently in the service container.

---

## NFR8 — Observability

Check execution should produce sufficient structured logging to understand:

* monitored service
* check type
* duration
* success/failure

Logs should not expose unnecessary sensitive endpoint information beyond what is already part of the monitored-service configuration.

---

# Architecture Constraints

The implementation must preserve the existing repository architecture.

Specifically:

1. OpenAPI remains the source of truth for HTTP contracts.
2. Generated code must not be manually edited.
3. HTTP and TCP checks must implement a common monitoring abstraction.
4. Monitoring orchestration must not depend directly on HTTP-specific behavior.
5. The web application must consume generated API types/client behavior.
6. Project operations must use documented Task targets.
7. No persistence may be introduced.
8. No CI/CD changes are required for this Epic.
9. The implementation must not create a dynamic arbitrary-code plugin system.

---

# Preferred Conceptual Architecture

The implementation should resemble:

```text
                  Monitoring Service
                         |
                         v
                  Check Dispatcher
                         |
              +----------+----------+
              |                     |
              v                     v
        HTTP Check              TCP Check
              |                     |
              +----------+----------+
                         |
                         v
                Normalized Result
                         |
                         v
               Service Health Model
```

The exact package and interface names are implementation decisions.

The architectural requirement is the dependency direction.

The monitoring core selects and invokes a check implementation.

Individual check implementations do not redefine the monitoring domain.

---

# Check Registration / Selection

The architecture must provide an explicit way to map a configured check type to its implementation.

This may be:

* a registry
* a constructor/factory
* dependency injection
* another simple explicit mechanism

Avoid:

* reflection-heavy plugin discovery
* runtime-loaded binaries
* arbitrary scripting
* unnecessary framework dependencies

The mechanism should be obvious to an engineer examining the repository.

---

# Guardrail Expectations

The completed Epic must pass all existing repository validation and extend it where necessary.

## API Guardrails

* OpenAPI lint
* schema validation
* generated-code freshness
* compatibility checking

## Service Guardrails

* formatting
* static analysis
* unit tests
* HTTP regression tests
* TCP positive tests
* TCP negative tests
* timeout/cancellation tests
* common abstraction tests where useful

## Web Guardrails

* generated-client freshness
* lint
* build
* tests covering monitoring-type rendering

## Integration Guardrails

Docker Compose integration must demonstrate:

* healthy HTTP target
* degraded/down HTTP target as already supported
* reachable TCP target
* unavailable TCP target

`task validate` must include the applicable integration validation.

## Deployment Guardrails

* container builds
* Helm lint
* Helm rendering/static validation

---

# Required Negative Tests

At minimum, TCP monitoring must test:

* reachable endpoint
* refused connection
* unreachable endpoint or deterministic failure equivalent
* timeout
* invalid port
* missing required TCP configuration
* cancellation where applicable

The broader architecture should test:

* unsupported check type
* HTTP monitor still functions
* one check implementation cannot accidentally bypass the normalized result path

---

# Definition of Complete

This Epic is complete when:

1. Pulse supports both HTTP and TCP monitoring.
2. Each monitored service explicitly uses an appropriate check type.
3. HTTP monitoring operates through the new common monitoring abstraction.
4. TCP monitoring is implemented through the same abstraction.
5. TCP checks can successfully identify reachable and unavailable endpoints.
6. Both check types produce the existing normalized service-health representation.
7. OpenAPI accurately represents the new monitoring configuration.
8. Generated Go artifacts are current.
9. Generated TypeScript artifacts are current.
10. The web application represents both HTTP and TCP monitored services.
11. Existing HTTP behavior remains regression-tested.
12. TCP positive, negative, timeout, and invalid-configuration cases are tested.
13. Docker Compose provides deterministic TCP integration targets.
14. No external network dependency is required for validation.
15. No arbitrary shell-command or dynamic-code execution mechanism has been introduced.
16. No database or other out-of-scope infrastructure has been introduced.
17. Relevant architecture documentation is updated.
18. `task validate` succeeds from a clean working tree.
19. The capability is demonstrable end-to-end using the local Compose environment.

---

# Expected Demonstration

When complete, Pulse should show multiple monitored services using different monitoring types.

For example:

| Service    | Check Type | Target    | Status   |
| ---------- | ---------- | --------- | -------- |
| Web API    | HTTP       | `/health` | Healthy  |
| Legacy App | HTTP       | `/status` | Degraded |
| PostgreSQL | TCP        | `:5432`   | Healthy  |
| Broker     | TCP        | `:5672`   | Down     |

A user should see one coherent monitoring experience even though different implementations are responsible for determining health.

---

# Success Measures

The primary success criterion is not simply that TCP connectivity can be tested.

Success means the monitoring architecture now has a proven extension pattern.

Evaluate:

* ability to add the second check type without duplicating monitoring orchestration
* preservation of existing HTTP behavior
* API compatibility
* deterministic validation
* clear component boundaries
* negative-path coverage
* future check types can follow the established seam
* no unnecessary generic plugin framework was introduced

The desired architectural result is:

> **HTTP is no longer the monitoring system. HTTP is one implementation of the monitoring system.**

---

# Future Extensibility — Not Part of This Epic

The architecture established here may later support:

* DNS checks
* ping
* TLS checks
* UDP checks
* application-specific protocol checks

These are not required now.

Do not implement speculative support for them.

Use the simplest architecture that cleanly supports the two real implementations required today:

**HTTP and TCP.**
