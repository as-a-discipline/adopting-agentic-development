#!/usr/bin/env bash
# Deterministic repository-structure validation.
#
# This validates DISCOVERABILITY and STRUCTURE only — required AGENTS.md/SKILL.md
# files exist, skill frontmatter has required fields, and required component
# Taskfiles exist. It does NOT validate the semantic correctness of any agent
# skill or instruction content.
set -euo pipefail

cd "$(dirname "$0")/.."

fail=0

require_file() {
  local path="$1"
  if [ ! -f "$path" ]; then
    echo "MISSING: $path"
    fail=1
  fi
}

echo "== Required root files =="
require_file "AGENTS.md"
require_file "README.md"
require_file "Makefile"
require_file "Taskfile.yml"
require_file ".github/copilot-instructions.md"

echo "== Required component AGENTS.md =="
for d in api service web deploy deploy/compose deploy/helm policies; do
  require_file "$d/AGENTS.md"
done

echo "== Required component Taskfile.yml =="
for d in api service web deploy/compose deploy/helm; do
  require_file "$d/Taskfile.yml"
done

echo "== Required skills =="
for skill in api-contract service-development web-development compose-integration helm-validation repository-validation; do
  skill_file=".agents/skills/$skill/SKILL.md"
  require_file "$skill_file"
  if [ -f "$skill_file" ]; then
    if ! grep -q "^name:" "$skill_file" || ! grep -q "^description:" "$skill_file"; then
      echo "INVALID FRONTMATTER: $skill_file (missing name/description)"
      fail=1
    fi
  fi
done

echo "== Required root ADRs =="
for f in adr/0001-component-boundaries-in-the-monorepo.md adr/0002-task-as-the-execution-contract.md adr/0003-containerized-tooling.md; do
  require_file "$f"
done

echo "== Required component ADRs =="
require_file "api/adr/0001-contract-first-api.md"
require_file "service/adr/0001-no-persistence-in-baseline.md"

echo "== Required verification docs =="
require_file "verification/skills.md"
require_file "verification/architecture.md"

if [ "$fail" -eq 0 ]; then
  echo ""
  echo "Repository structure: OK"
else
  echo ""
  echo "Repository structure: FAILED"
fi

exit "$fail"
