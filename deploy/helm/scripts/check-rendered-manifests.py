#!/usr/bin/env python3
"""Static validation of rendered Pulse Helm manifests.

Run against the output of `helm template` (see deploy/helm/Taskfile.yml,
task template:check). This is intentionally lightweight: it does not require
a live cluster, only that the rendered YAML is well-formed and follows the
rules in deploy/helm/AGENTS.md (probes on every deployment, no secrets).
"""
import sys
import yaml


def fail(msg: str) -> None:
    print(f"FAIL: {msg}")
    sys.exit(1)


def main() -> None:
    if len(sys.argv) != 2:
        fail("usage: check-rendered-manifests.py <rendered.yaml>")

    with open(sys.argv[1]) as f:
        docs = [d for d in yaml.safe_load_all(f) if d]

    if not docs:
        fail("no manifests rendered")

    kinds = {}
    for doc in docs:
        kind = doc.get("kind")
        kinds.setdefault(kind, []).append(doc)

        if kind == "Secret":
            fail(f"chart must not commit secrets, found Secret: {doc.get('metadata', {}).get('name')}")

        if kind == "Deployment":
            name = doc["metadata"]["name"]
            containers = doc["spec"]["template"]["spec"]["containers"]
            for c in containers:
                if "readinessProbe" not in c:
                    fail(f"Deployment {name} container {c['name']} missing readinessProbe")
                if "livenessProbe" not in c:
                    fail(f"Deployment {name} container {c['name']} missing livenessProbe")
                if not c.get("image"):
                    fail(f"Deployment {name} container {c['name']} missing image")

    required = {"Deployment", "Service"}
    missing = required - kinds.keys()
    if missing:
        fail(f"expected kinds missing from rendered output: {sorted(missing)}")

    if len(kinds.get("Deployment", [])) < 2:
        fail("expected at least 2 Deployments (service + web)")
    if len(kinds.get("Service", [])) < 2:
        fail("expected at least 2 Services (service + web)")

    print(f"OK: rendered manifests valid — {sum(len(v) for v in kinds.values())} resources, "
          f"no secrets, all deployments have readiness/liveness probes.")


if __name__ == "__main__":
    main()
