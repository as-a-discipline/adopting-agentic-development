#!/usr/bin/env bash
# Deterministic, intentionally simple check for forbidden cross-component
# source dependencies (see /adr/0001-component-boundaries-in-the-monorepo.md
# and root AGENTS.md's Repository Architecture section).
#
# This is a transparency mechanism, not an exhaustive static-analysis
# platform: it looks for source-level references that reach across a
# top-level component boundary (e.g. service/ Go code importing/reading from
# web/, or web/ source importing from service/). It does not understand full
# language semantics — it is a grep-based scan a reader can verify by eye.
set -euo pipefail

cd "$(dirname "$0")/.."

fail=0

# Pairs of (component whose source we scan, forbidden-reference patterns).
# Only handwritten source is scanned — generated/ and bin/ are derived
# artifacts, not hand-authored cross-component dependencies.
check_component() {
  local component="$1"
  shift
  local forbidden=("$@")

  local files
  files=$(find "$component" -type f \
    \( -name '*.go' -o -name '*.ts' -o -name '*.mjs' -o -name '*.js' \) \
    -not -path "*/generated/*" \
    -not -path "*/bin/*" \
    -not -path "*/dist/*" \
    2>/dev/null || true)

  [ -z "$files" ] && return 0

  for pattern in "${forbidden[@]}"; do
    local hits
    hits=$(grep -lE "$pattern" $files 2>/dev/null || true)
    if [ -n "$hits" ]; then
      echo "FORBIDDEN: $component source references '$pattern':"
      echo "$hits" | sed 's/^/  /'
      fail=1
    fi
  done
}

echo "== Component boundary scan =="

# service/ (Go) must never reference web/ paths.
check_component service '(^|[^A-Za-z0-9_./-])\.\./web(/|$)' '"web/'

# web/ (TypeScript/JS) must never reference service/ Go source paths.
check_component web '(^|[^A-Za-z0-9_./-])\.\./service(/|$)' '"service/'

if [ "$fail" -eq 0 ]; then
  echo "OK: no forbidden cross-component source references found."
else
  echo "FAILED: forbidden cross-component source references found (see above)."
fi

exit "$fail"
