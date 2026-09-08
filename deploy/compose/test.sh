#!/usr/bin/env bash
# Deterministic end-to-end integration test for the Compose stack. Assumes
# `task compose:up` has already started the environment (services are healthy
# per their own Compose healthchecks) — this script verifies actual monitoring
# *behavior* against the deterministic fake targets: healthy stays healthy,
# the deliberately-slow target reports degraded, and the deliberately-failing
# target reports down. No external network dependency is involved.
set -euo pipefail

API="http://localhost:18080/api/v1"
WEB="http://localhost:18000"

fail() {
  echo "FAIL: $1" >&2
  exit 1
}

echo "== compose:test — verifying pulse-service API directly =="

services_json="$(curl -sf "$API/services")" || fail "GET $API/services did not respond"
echo "$services_json" | python3 -m json.tool >/dev/null || fail "GET /services did not return valid JSON"

for id in healthy-target degraded-target down-target; do
  echo "-- checking $id --"
  result="$(curl -sf -X POST "$API/services/$id/check")" || fail "POST /services/$id/check did not respond"
  status="$(echo "$result" | python3 -c 'import json,sys; print(json.load(sys.stdin)["status"])')"
  echo "   status=$status"
  case "$id" in
    healthy-target)
      [ "$status" = "healthy" ] || fail "expected healthy-target to report 'healthy', got '$status'"
      ;;
    degraded-target)
      [ "$status" = "degraded" ] || fail "expected degraded-target to report 'degraded', got '$status'"
      ;;
    down-target)
      [ "$status" = "down" ] || fail "expected down-target to report 'down', got '$status'"
      ;;
  esac
done

echo "== compose:test — verifying pulse-web serves the app and proxies the API =="

curl -sf "$WEB/" | grep -q "<title>Pulse</title>" || fail "pulse-web did not serve the expected index.html"
curl -sf "$WEB/api/v1/services" >/dev/null || fail "pulse-web did not proxy /api/v1/services to pulse-service"

echo "OK: compose integration behaves deterministically (healthy/degraded/down all verified)."
