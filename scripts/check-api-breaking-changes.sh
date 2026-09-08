#!/usr/bin/env bash
# Detect breaking changes in api/openapi.yaml between a git reference (the
# current branch's last pushed commit, or a target mainline branch) and the
# current working tree — so a breaking change is visible before it's pushed
# or merged, not only against the static compatibility baseline checked by
# `task api:compat`.
#
# Usage: check-api-breaking-changes.sh [ref]
#
# Env vars:
#   TOOL          "oasdiff" (default) or "openapi-changes" — which pinned,
#                 containerized tool performs the comparison.
#   OASDIFF_IMAGE          pinned oasdiff image (passed in by the Taskfile)
#   OPENAPI_CHANGES_IMAGE  pinned pb33f/openapi-changes image (passed in by
#                          the Taskfile)
#
# Reference resolution (when [ref] is not given):
#   1. The current branch's upstream tracking branch (i.e. its last pushed
#      commit) — `git rev-parse --abbrev-ref --symbolic-full-name @{u}`.
#   2. Otherwise, the target mainline branch: origin/main, origin/master,
#      then local main/master, whichever resolves first.
# This mirrors "compare against what's already pushed, or against mainline if
# nothing has been pushed yet" — a reviewer's actual question before merging.
set -euo pipefail

TOOL="${TOOL:-oasdiff}"
OASDIFF_IMAGE="${OASDIFF_IMAGE:-tufin/oasdiff:v1.31.0}"
OPENAPI_CHANGES_IMAGE="${OPENAPI_CHANGES_IMAGE:-pb33f/openapi-changes:v0.2.11}"
SPEC_PATH="api/openapi.yaml"

repo_root="$(git rev-parse --show-toplevel)"
cd "$repo_root"

ref="${1:-}"

if [ -z "$ref" ]; then
  if upstream="$(git rev-parse --abbrev-ref --symbolic-full-name '@{u}' 2>/dev/null)"; then
    ref="$upstream"
    ref_reason="upstream tracking branch (last pushed commit)"
  else
    for candidate in origin/main origin/master main master; do
      if git rev-parse --verify --quiet "$candidate" >/dev/null; then
        ref="$candidate"
        ref_reason="target mainline branch (no upstream configured for this branch)"
        break
      fi
    done
  fi
fi

if [ -z "$ref" ]; then
  echo "ERROR: could not resolve a reference to compare against." >&2
  echo "No upstream tracking branch is configured and none of origin/main," >&2
  echo "origin/master, main, master could be resolved. Pass one explicitly:" >&2
  echo "  task api:breaking-changes -- <ref>" >&2
  exit 1
fi

echo "== api:breaking-changes — comparing '$ref' (${ref_reason:-explicit ref}) -> working tree =="
echo "tool: $TOOL"

if ! git cat-file -e "${ref}:${SPEC_PATH}" 2>/dev/null; then
  echo "No '$SPEC_PATH' found at '$ref' — nothing to compare (new contract, or"
  echo "'$ref' predates the API). Treating as pass."
  exit 0
fi

case "$TOOL" in
  oasdiff)
    tmpdir="$(mktemp -d)"
    trap 'rm -rf "$tmpdir"' EXIT
    git show "${ref}:${SPEC_PATH}" > "$tmpdir/old-openapi.yaml"
    docker run --rm \
      -v "$tmpdir":/old \
      -v "$repo_root/api":/new \
      "$OASDIFF_IMAGE" \
      breaking --fail-on ERR /old/old-openapi.yaml /new/openapi.yaml
    ;;

  openapi-changes)
    report_file="$(mktemp)"
    trap 'rm -f "$report_file"' EXIT
    docker run --rm -v "$repo_root":/repo -w /repo \
      "$OPENAPI_CHANGES_IMAGE" \
      report --reproducible --no-logo "${ref}:${SPEC_PATH}" "./${SPEC_PATH}" \
      > "$report_file"
    python3 "$(dirname "$0")/parse-openapi-changes-report.py" "$report_file"
    ;;

  *)
    echo "ERROR: unknown TOOL '$TOOL' (expected 'oasdiff' or 'openapi-changes')" >&2
    exit 1
    ;;
esac
