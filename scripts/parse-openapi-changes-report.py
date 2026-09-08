#!/usr/bin/env python3
"""Parse a `pb33f/openapi-changes report` JSON file and fail (exit 1) if any
detected change is marked breaking.

`report`'s own process exit code only signals tool errors (e.g. an
unresolvable git ref) — it exits 0 whether or not breaking changes were
found, and there is no documented flag to gate on breaking-only changes
(`--error-on-diff` gates on *any* diff, breaking or not). This script reads
the JSON `changes` array and applies the breaking-only gate ourselves,
matching the fields verified by running the tool directly against this
repository:

  {"reportSummary": {"<section>": {"totalChanges": N, "breakingChanges": N}}},
  "changes": [{"breaking": true|false, "path": "...", "property": "...", ...}]}

or, when there is no diff at all: {"message": "No changes found..."}.
"""
import json
import sys


def main() -> int:
    if len(sys.argv) != 2:
        print("usage: parse-openapi-changes-report.py <report.json>", file=sys.stderr)
        return 2

    with open(sys.argv[1], encoding="utf-8") as f:
        report = json.load(f)

    if "message" in report and "changes" not in report:
        print(f"OK: {report['message']} — no breaking changes.")
        return 0

    changes = report.get("changes", [])
    breaking = [c for c in changes if c.get("breaking")]

    print(f"Total changes: {len(changes)}; breaking: {len(breaking)}")

    if not breaking:
        print("OK: no breaking changes detected.")
        return 0

    print("\nBreaking changes:")
    for c in breaking:
        path = c.get("path", "?")
        prop = c.get("property", "")
        change_text = c.get("changeText", "changed")
        original = c.get("original")
        new = c.get("new")
        detail = f"{path}" + (f".{prop}" if prop else "")
        if original is not None or new is not None:
            print(f"  - {detail}: {change_text} ({original!r} -> {new!r})")
        else:
            print(f"  - {detail}: {change_text}")

    print(f"\nFAIL: {len(breaking)} breaking change(s) detected.")
    return 1


if __name__ == "__main__":
    sys.exit(main())
