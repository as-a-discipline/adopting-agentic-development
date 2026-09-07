# 3. Containerized Tooling

* Status: accepted
* Date: 2026-09-07

## Context and Problem Statement

Pulse's build involves an OpenAPI linter, an OpenAPI-to-Go generator, an
OpenAPI-to-TypeScript generator, Go toolchain, Node toolchain, and Helm. Requiring
every contributor (and every agent's execution environment) to install and keep all
of these in sync on the host is fragile and violates the goal of minimal host
prerequisites.

## Decision Drivers

* Host prerequisites should be limited to Git, Docker, Make, and Task.
* Tool versions must be reproducible across machines and over time.
* Agents should not need to reason about host-specific tool installation issues.

## Considered Options

1. Require host installation of Go, Node, the OpenAPI generator, and linters.
2. Use containerized, pinned versions of these tools, invoked from Task targets.
3. A mix: some tools containerized, others required on host.

## Decision Outcome

Chosen option: **2 — containerized, pinned tooling where practical.**

* Generators and major build/lint tooling run inside pinned container images,
  invoked via `docker run` from the relevant component's `Taskfile.yml`.
* Host prerequisites remain: Git, Docker, Make, Task.
* Where a tool is genuinely lightweight and has no meaningful version-drift risk,
  a host fallback may be documented, but containerized execution is the default.

## Consequences

* Good: `make bootstrap && task validate` works the same on any contributor's or
  agent's machine, regardless of locally installed Go/Node versions.
* Good: pinned image tags make generation and lint output reproducible over time.
* Bad: first-run tasks are slower (image pulls) and require Docker to be usable —
  this is why `make check-prereqs` explicitly verifies Docker is running.
