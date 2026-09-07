#!/usr/bin/env python3
"""Dependency governance summary + allow-list check (Session 4).

Deliberately simple: for each component, list DIRECT third-party
dependencies and confirm each is present in policies/allowed-dependencies.txt.
This is a transparency mechanism, not a license/supportability scanner --
see policies/README.md rule 4 and policies/allowed-dependencies.txt.
"""
import json
import re
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]


def load_allowlist() -> dict[str, set[str]]:
    allowed: dict[str, set[str]] = {}
    path = ROOT / "policies" / "allowed-dependencies.txt"
    for line in path.read_text().splitlines():
        line = line.strip()
        if not line or line.startswith("#"):
            continue
        component, _, dep = line.partition(":")
        allowed.setdefault(component, set()).add(dep)
    return allowed


def go_mod_direct_deps(go_mod: Path) -> list[str]:
    if not go_mod.exists():
        return []
    text = go_mod.read_text()
    deps = []
    in_require_block = False
    for line in text.splitlines():
        stripped = line.strip()
        if stripped == "require (":
            in_require_block = True
            continue
        if in_require_block:
            if stripped == ")":
                in_require_block = False
                continue
            if stripped.endswith("// indirect"):
                continue
            m = re.match(r"(\S+)\s+v", stripped)
            if m:
                deps.append(m.group(1))
        else:
            m = re.match(r"require\s+(\S+)\s+v", stripped)
            if m:
                deps.append(m.group(1))
    return deps


def package_json_deps(package_json: Path) -> list[str]:
    if not package_json.exists():
        return []
    data = json.loads(package_json.read_text())
    return sorted(data.get("dependencies", {}).keys())


def main() -> int:
    allowed = load_allowlist()
    fail = False

    components = {
        "service": go_mod_direct_deps(ROOT / "service" / "go.mod"),
        "web": package_json_deps(ROOT / "web" / "package.json"),
    }

    print("== Dependency governance summary ==")
    for component, deps in components.items():
        if not deps:
            print(f"{component}: (no direct third-party dependencies)")
            continue
        print(f"{component}:")
        allowed_for_component = allowed.get(component, set())
        for dep in deps:
            marker = "OK" if dep in allowed_for_component else "NOT ALLOW-LISTED"
            print(f"  - {dep}  [{marker}]")
            if dep not in allowed_for_component:
                fail = True

    if fail:
        print(
            "\nFAILED: one or more direct dependencies are not in "
            "policies/allowed-dependencies.txt. Add the new dependency there "
            "explicitly (a reviewable diff) if it's intentional."
        )
        return 1

    print("\nOK: all direct dependencies are allow-listed.")
    return 0


if __name__ == "__main__":
    sys.exit(main())
