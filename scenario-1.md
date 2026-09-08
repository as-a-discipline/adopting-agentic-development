# Feature: Interactive Service Diagnostics

## Summary

Extend Pulse beyond passive service-health reporting by allowing users to perform targeted diagnostics against monitored services when a problem is detected.

The first diagnostic capability delivered under this Epic will be **network path tracing**. A user should be able to select a monitored service, initiate a traceroute-style diagnostic, and view the resulting network hops and timing information through the Pulse web application.

The capability must be implemented as a distinct Diagnostics domain and API surface rather than being embedded directly into the existing monitoring implementation.

This Epic establishes the architectural pattern that future diagnostic capabilities can follow without requiring the initial release to implement additional tools such as DNS lookup, port testing, TLS inspection, or ping.

---

# Business Outcome

When Pulse reports that a monitored service is degraded or unavailable, a user should be able to gather basic diagnostic information without leaving the application.

Today, Pulse can answer:

> Is the service healthy?

This Epic begins allowing Pulse to answer:

> What can I inspect to understand why it may not be healthy?

The delivered capability should reduce the gap between detection and initial investigation while establishing a reusable architecture for additional diagnostics.

---

# User Value

A user viewing a monitored service can initiate a network-path diagnostic against that service and inspect the resulting route without needing direct access to the Pulse runtime environment or a shell.

The diagnostic should be:

* simple to initiate
* scoped to an existing monitored service
* safe to execute
* deterministic in how it is invoked
* represented through a stable API contract
* understandable through both the API and web interface

---

# Epic Scope

This Epic includes the architecture and product capability necessary to support service diagnostics, with **network path tracing as the first supported diagnostic**.

The Epic includes:

* a new Diagnostics API domain
* an API operation for initiating a network path trace
* a diagnostic request and response contract
* generated Go API scaffolding
* generated frontend client support
* a distinct diagnostics implementation boundary within the Go service
* a safe abstraction for executing network path diagnostics
* deterministic testing through a simulated diagnostic implementation
* web UI for initiating and viewing a diagnostic
* component and integration tests
* documentation
* container and deployment changes required to support the capability
* applicable security, architecture, and policy guardrails

---

# Out of Scope

The following are explicitly outside this Epic unless required to satisfy an acceptance criterion:

* persistent diagnostic history
* database storage
* scheduled diagnostics
* alerting
* notification integrations
* automatic diagnostics triggered by health-state changes
* arbitrary command execution
* user-supplied shell commands or command-line options
* diagnostic execution against arbitrary targets not represented by monitored services
* DNS diagnostics
* ICMP ping
* TCP port diagnostics
* TLS/certificate inspection
* packet capture
* remote command execution
* Kubernetes diagnostics
* authentication or authorization changes
* production deployment automation

Future diagnostic types should be able to use the architecture established by this Epic, but they are not required for completion.

---

# Functional Requirements

## FR1 — Diagnostics API

Pulse must expose a new Diagnostics API surface that is separate from the existing monitoring API responsibilities.

The API must support initiating a network path diagnostic for an existing monitored service.

A representative API shape may resemble:

`POST /api/v1/services/{serviceId}/diagnostics/traceroute`

The final contract must follow the repository's existing OpenAPI conventions.

The request must not permit arbitrary shell arguments.

---

## FR2 — Diagnostic Target

Diagnostics must operate against an existing monitored service.

The service's configured endpoint or hostname must determine the diagnostic target.

A caller must not be able to provide an unrestricted arbitrary hostname, IP address, shell command, or command-line argument through this API.

If the referenced monitored service does not exist, the API must return the repository's standard not-found response.

---

## FR3 — Network Path Result

A successful diagnostic response must provide structured data representing the observed network path.

At minimum, each returned hop must support:

* hop number
* address or hostname when available
* observed latency when available
* indication that the hop did not respond, where applicable

The overall result must identify:

* the monitored service
* the diagnostic target
* whether the diagnostic completed
* the ordered collection of hops

The API response must be structured data and must not simply return raw command-line output.

---

## FR4 — Diagnostic Execution Abstraction

The HTTP/API layer must not execute operating-system commands directly.

The diagnostics capability must use an internal abstraction representing network-path diagnostics.

For example, conceptually:

`Tracer.Trace(context, target) → TraceResult`

The exact interface is an implementation decision.

The architecture must allow:

* a real runtime implementation
* a deterministic simulated/test implementation

without changing API behavior.

---

## FR5 — Deterministic Simulation

The repository must contain a deterministic way to test network-path diagnostics without depending on:

* public internet connectivity
* conference/hotel networking
* external DNS
* ICMP availability
* host networking behavior
* external cloud services

The existing Docker Compose integration environment should be extended as required to provide predictable diagnostic results.

Integration tests must not require access to the public internet.

---

## FR6 — Web User Experience

The Pulse web application must allow a user to initiate the diagnostic from the context of an existing monitored service.

The user should not be required to manually enter a target hostname or IP address.

The UI must display the resulting path in an understandable ordered format.

At minimum, display:

* hop number
* address/hostname
* latency or timeout/unavailable indication

The UI must clearly represent:

* diagnostic in progress
* successful completion
* diagnostic failure
* unavailable/non-responsive hops

The implementation should remain visually consistent with the existing Pulse application and should not introduce a new UI framework solely for this feature.

---

## FR7 — API Generation

Changes to the OpenAPI contract must continue to generate:

* applicable Go server types/interfaces
* applicable TypeScript client/types

Generated files must not be modified manually.

The existing generated-code freshness checks must continue to pass.

---

# Non-Functional Requirements

## NFR1 — Security

The diagnostic feature must not create a general-purpose remote command execution mechanism.

User-controlled values must never be concatenated into shell command strings.

The implementation must not expose arbitrary command options through the API.

The application must validate the diagnostic target derived from the monitored service before execution.

Any OS-level diagnostic implementation must use the narrowest practical execution mechanism.

---

## NFR2 — Bounded Execution

A diagnostic must have a finite execution duration.

The implementation must honor request cancellation and/or context timeout.

A failed or unreachable destination must not allow the diagnostic process to run indefinitely.

Reasonable maximum hop and execution limits must be defined by the implementation or configuration.

---

## NFR3 — API Compatibility

Existing Pulse API operations and schemas must remain backward compatible.

Introducing Diagnostics must not materially change the behavior of:

* service listing
* service details
* existing health checks

Any OpenAPI compatibility guardrails must continue to pass.

---

## NFR4 — Architectural Isolation

Diagnostics must be implemented as a separate domain responsibility from monitoring.

Conceptually:

`API → Diagnostics Domain → Diagnostic Provider`

The monitoring domain must not become responsible for traceroute implementation details.

The API handler must not become the location for operating-system execution logic.

---

## NFR5 — Testability

The diagnostic execution implementation must be replaceable for testing.

Unit and integration tests must not require a functioning host traceroute implementation unless specifically testing that adapter.

Business and API behavior must be testable through deterministic substitutes.

---

## NFR6 — Portability

The baseline application must continue to operate through its documented containerized development and validation workflows.

The feature must not introduce an undocumented host prerequisite.

If an operating-system utility is required at runtime, it must be provided by the applicable service container image rather than assumed to exist on the developer workstation.

---

## NFR7 — Observability

Diagnostic execution must produce sufficient structured logging to identify:

* requested monitored service
* diagnostic type
* completion/failure
* duration

Logs must not expose unsafe command strings or unrelated sensitive runtime information.

---

## NFR8 — Deployment

The existing Docker Compose environment and Helm chart must remain valid.

Any new container capability or runtime requirement introduced for diagnostics must be explicitly represented in the appropriate deployment artifacts.

The feature must not require a live Kubernetes cluster for local validation.

---

# Architecture Constraints

The implementation must comply with the repository's existing architecture decisions.

Specifically:

1. **OpenAPI remains the source of truth for HTTP contracts.**
2. **Generated code is not manually modified.**
3. **The web application interacts with the service only through the API contract.**
4. **Component-specific logic remains inside its owning component/domain.**
5. **Project operations use the defined Task interfaces.**
6. **Build and generation tooling remains containerized where practical.**
7. **No persistence is introduced.**
8. **No CI/CD modifications are required for this Epic unless separately approved.**

A new ADR should be created only if this Epic introduces an architectural decision not already represented by an existing ADR.

---

# Guardrail Expectations

The completed implementation must pass all existing repository validation.

In addition, this Epic should add or extend deterministic checks where necessary to cover the new architectural and security constraints.

Relevant checks should include:

### API Guardrails

* OpenAPI lint
* schema validation
* generated-code freshness
* compatibility validation

### Service Guardrails

* formatting
* static analysis
* unit tests
* generated-interface compliance
* bounded execution tests
* negative security cases

### Web Guardrails

* generated-client freshness
* lint
* tests
* build

### Integration Guardrails

* deterministic diagnostic integration test
* existing monitoring behavior regression tests
* Docker Compose validation

### Deployment Guardrails

* container build
* Helm lint
* Helm template/static validation

The feature is not complete merely because the traceroute implementation itself works.

---

# Negative Conditions That Must Be Tested

At minimum, validation must exercise:

* unknown service ID
* unreachable target
* diagnostic timeout
* one or more non-responsive intermediate hops
* diagnostic execution failure
* cancelled request
* malformed or unexpected diagnostic result from the implementation boundary
* attempts to bypass target restrictions where applicable

Tests must verify that these conditions result in controlled application behavior rather than crashes, hangs, or arbitrary command execution.

---

# Definition of Complete

This Epic is complete when:

1. A user can initiate a network-path diagnostic for an existing monitored service through the Pulse API.
2. The same capability is available through the web application.
3. Results are represented as structured diagnostic data.
4. The implementation uses a dedicated Diagnostics domain and execution abstraction.
5. Arbitrary command execution is not exposed.
6. Diagnostic execution is bounded.
7. A deterministic simulated implementation supports automated testing.
8. OpenAPI-generated Go and TypeScript artifacts are current.
9. Existing monitoring functionality remains behaviorally compatible.
10. Component validation passes.
11. Docker Compose integration validation passes.
12. Helm validation passes.
13. Negative and failure conditions are tested.
14. Relevant architecture and usage documentation is updated.
15. `task validate` completes successfully from a clean working tree.
16. The feature is demonstrable end-to-end without relying on the public internet.
17. No database, authentication system, unrelated diagnostic capability, or other out-of-scope architecture has been introduced.

---

# Expected User Demonstration

When complete, a user should be able to:

1. Open Pulse.
2. See a monitored service that is degraded or down.
3. Open that service.
4. Choose the network-path diagnostic.
5. Start the diagnostic.
6. See that the diagnostic is running.
7. Receive an ordered list of network hops and latency information.
8. See a clear representation of unreachable hops or failures.
9. Return to the normal monitoring view without affecting existing monitoring behavior.

This complete workflow must be demonstrable locally using the deterministic development environment.

---

# Success Measures

This Epic should be evaluated primarily on engineering completeness rather than the amount of code produced.

Relevant measures include:

* all acceptance criteria satisfied
* no regression in existing monitoring behavior
* deterministic end-to-end validation
* no new uncontrolled execution capability
* no manual modification of generated artifacts
* new API remains consistent with existing API conventions
* clear separation between monitoring and diagnostics
* successful repository-wide validation
* ability for a future diagnostic type to follow the established architecture without redesigning the HTTP/API and domain boundaries

---

# Future Extensibility — Architectural Intent Only

This Epic should establish a pattern capable of later supporting additional diagnostic types such as:

* DNS lookup
* TCP connectivity
* ping
* TLS inspection

These are **not part of the Epic's required implementation**.

Do not build speculative abstractions simply to support every imaginable diagnostic.

The design should create a clean extension point while implementing only what the current requirements require.

The desired principle is:

> Design the seam; implement the need.
