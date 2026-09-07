#!/usr/bin/env bash
# Deterministic generated-code freshness check.
#
# Usage: check-generated-fresh.sh <generated-path> [git-root]
#
# Assumes the caller has just regenerated <generated-path> from source. Fails
# if the working tree shows any diff (including new/untracked files) under
# that path — meaning committed generated code was stale relative to the
# source contract. This does not hide or discard the regenerated output; it
# only reports drift.
set -euo pipefail

generated_path="$1"
git_root="${2:-.}"

cd "$git_root"

if ! git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  echo "SKIP: not inside a git work tree — cannot check generated-code freshness"
  exit 0
fi

diff_output="$(git status --porcelain -- "$generated_path")"

if [ -n "$diff_output" ]; then
  echo "STALE: regenerating '$generated_path' produced a diff against the committed version:"
  echo "$diff_output"
  echo ""
  echo "Committed generated code is stale relative to its source contract."
  echo "Regenerate and commit the result — do not hand-edit files under '$generated_path'."
  exit 1
fi

echo "OK: '$generated_path' matches what regeneration produces."
