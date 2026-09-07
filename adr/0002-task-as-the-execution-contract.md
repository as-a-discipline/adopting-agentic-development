# 2. Task as the Execution Contract

* Status: accepted
* Date: 2026-09-07

## Context and Problem Statement

Pulse involves multiple toolchains (OpenAPI linting/generation, Go, a TypeScript
frontend, Docker Compose, Helm). If humans, agents, and (eventually) CI each learn
raw per-tool commands, the actual command needed to build/test/lint/validate any
given component becomes tribal knowledge, and it's easy for an agent to invent a
plausible-looking but incorrect command.

## Decision Drivers

* Agents must have a small, deterministic, discoverable set of verbs per component.
* The same commands must work identically for a human and an agent, locally and
  (later) in CI.
* Adding a new component or toolchain should not require re-teaching agents new raw
  commands.

## Considered Options

1. Document raw shell commands per tool in each `AGENTS.md`.
2. Use Make as the full build system.
3. Use [Task](https://taskfile.dev) as a composable task runner, with each component
   owning its own `Taskfile.yml`, composed by a root `Taskfile.yml`.

## Decision Outcome

Chosen option: **3 — Task as the execution contract.**

* Once Task is available, all build/test/lint/generate/validate actions go through
  documented `task` targets (e.g. `task service:test`, `task api:validate`).
* Agents must not invent alternate build/test/lint/generation commands unless
  modifying the Task abstraction is itself the assigned work.
* `task validate` is the canonical, root-level, whole-repository validation command.
* Make is reserved for host bootstrap only (`make check-prereqs`, `make bootstrap`)
  and is never a second build system.

## Consequences

* Good: one mental model ("what's the task target?") replaces per-tool command
  memorization.
* Good: swapping an underlying tool (e.g. a different OpenAPI generator) only
  requires updating the Taskfile that wraps it — the interface (`task api:generate`)
  stays stable.
* Bad: Task itself becomes a required host dependency (mitigated by `make bootstrap`
  installing it, and by containerizing the tools it wraps rather than Task itself).
