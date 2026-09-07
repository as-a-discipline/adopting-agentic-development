---
name: helm-validation
description: Use when linting, templating, or validating the Pulse Helm chart (deploy/helm/). Trigger on requests to change chart values/templates, or to fix a helm:validate failure.
---

# Helm Validation Skill

## When to use

* The Helm chart's deployments/services/values need to change.
* `helm:lint`, `helm:template`, or `helm:validate` is failing.

## Workflow

1. Read `deploy/helm/AGENTS.md` before editing chart files.
2. Run `task helm:lint` and fix reported issues first.
3. Run `task helm:template` and inspect rendered manifests for correctness (probes,
   image references, config).
4. Distinguish chart errors (bad template logic, missing values) from application
   errors (the app itself misbehaving) — this skill only owns the former.
5. Run `task helm:validate` for the combined check. No live cluster is required.

## Evidence of success

* `task helm:validate` exits 0.
* Rendered manifests contain sensible readiness/liveness probes and no committed
  secrets.
