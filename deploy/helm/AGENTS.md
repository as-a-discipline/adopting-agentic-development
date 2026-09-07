# AGENTS.md — deploy/helm/

Kubernetes packaging for Pulse.

## Rules

* Chart must produce valid Kubernetes resources without requiring a live cluster to
  validate.
* Sensible readiness/liveness probes for every deployment.
* Configuration (image names/tags, replicas, app config) goes through `values.yaml`
  — no hardcoded environment-specific values.
* No secrets committed to the chart.
* `helm:validate` must run lint, template rendering, and reasonable static manifest
  validation.

## Task Interface

```text
task helm:lint
task helm:template
task helm:validate
```
