#!/usr/bin/env bash
# Deterministic, deliberately basic container/deployment quality checks
# (see policies/README.md rule 5). Not a full security scanner: it checks
# for a small number of concrete, explainable things any reviewer can verify
# by reading the same files this script reads.
set -euo pipefail

cd "$(dirname "$0")/.."

fail=0

echo "== Container/deployment quality checks =="

check_dockerfile() {
  local dockerfile="$1"

  # 1. No unpinned "latest" base images.
  local latest_hits
  latest_hits=$(grep -nE '^FROM\s+\S+:latest(\s|$)' "$dockerfile" || true)
  if [ -n "$latest_hits" ]; then
    echo "FAIL ($dockerfile): base image pinned to 'latest':"
    echo "$latest_hits" | sed 's/^/  /'
    fail=1
  fi

  # 2. Every FROM line must specify an explicit tag (not just a bare image
  #    name, which implicitly means 'latest').
  local unpinned_hits
  unpinned_hits=$(grep -nE '^FROM\s+[^:[:space:]]+(\s|$)' "$dockerfile" | grep -vE '\bAS\b' || true)
  # (FROM lines that reference a previous build stage by name, e.g.
  # "FROM build", are excluded further below by name-matching.)
  if [ -n "$unpinned_hits" ]; then
    # Filter out references to earlier stage names declared via "AS <name>"
    stage_names=$(grep -oE '\bAS\s+\S+' "$dockerfile" | awk '{print $2}')
    real_hits=""
    while IFS= read -r hit; do
      [ -z "$hit" ] && continue
      image=$(echo "$hit" | sed -E 's/^[0-9]+:FROM\s+//' | awk '{print $1}')
      is_stage_ref=0
      for s in $stage_names; do
        [ "$image" = "$s" ] && is_stage_ref=1
      done
      [ "$is_stage_ref" -eq 0 ] && real_hits="$real_hits$hit\n"
    done <<< "$unpinned_hits"
    if [ -n "$real_hits" ]; then
      echo "FAIL ($dockerfile): base image missing an explicit tag:"
      echo -e "$real_hits" | sed '/^$/d;s/^/  /'
      fail=1
    fi
  fi

  # 3. Final stage must run as a non-root user.
  local last_user
  last_user=$(grep -E '^USER\s+' "$dockerfile" | tail -1 | awk '{print $2}')
  if [ -z "$last_user" ]; then
    echo "FAIL ($dockerfile): no USER instruction — runtime stage defaults to root."
    fail=1
  elif [ "$last_user" = "root" ] || [ "$last_user" = "0" ]; then
    echo "FAIL ($dockerfile): final USER is '$last_user' (root)."
    fail=1
  else
    echo "OK ($dockerfile): runs as non-root user '$last_user'."
  fi
}

for df in service/Dockerfile web/Dockerfile; do
  if [ -f "$df" ]; then
    check_dockerfile "$df"
  fi
done

echo ""
echo "-- secrets scan (Compose + Helm) --"
secret_pattern='(password|secret|api[_-]?key|access[_-]?token|private[_-]?key)\s*[:=]\s*["'"'"']?[A-Za-z0-9]'
secret_hits=$(grep -rniE "$secret_pattern" deploy/compose/docker-compose.yml deploy/helm/pulse/values.yaml deploy/helm/pulse/templates/ 2>/dev/null || true)
if [ -n "$secret_hits" ]; then
  echo "FAIL: possible committed secret-like values found:"
  echo "$secret_hits" | sed 's/^/  /'
  fail=1
else
  echo "OK: no secret-like literal values found in Compose/Helm files."
fi

echo ""
if [ "$fail" -eq 0 ]; then
  echo "OK: container/deployment quality checks passed."
else
  echo "FAILED: container/deployment quality checks found issues (see above)."
fi

exit "$fail"
