#!/usr/bin/env bash
# E2E smoke test (F12.4): health + checkout scenarios via NodePort.
# Usage: ./scripts/test.sh  (override bases via env vars)
set -euo pipefail

ORDER_BASE="${ORDER_BASE:-http://localhost:18080}"
INVENTORY_BASE="${INVENTORY_BASE:-http://localhost:18081}"
PAYMENT_BASE="${PAYMENT_BASE:-http://localhost:18082}"

pass=0
fail=0

check() {
  local name="$1"; local expect="$2"; local actual="$3"
  if [ "$actual" = "$expect" ]; then
    echo "PASS  $name (got: $actual)"
    pass=$((pass+1))
  else
    echo "FAIL  $name (expected: $expect, got: $actual)"
    fail=$((fail+1))
  fi
}

code() { curl -s -o /dev/null -w '%{http_code}' "$@"; }
body() { curl -s "$@"; }

echo "==> health checks"
check "order /health"      "200" "$(code "$ORDER_BASE/health")"
check "inventory /health"   "200" "$(code "$INVENTORY_BASE/health")"
check "payment /health"     "200" "$(code "$PAYMENT_BASE/health")"

echo "==> inventory"
check "inventory A123 stock" "10"  "$(body "$INVENTORY_BASE/inventory/A123" | python3 -c 'import sys,json;print(json.load(sys.stdin)["stock"])')"
check "inventory UNKNOWN"      "404" "$(code "$INVENTORY_BASE/inventory/UNKNOWN")"

echo "==> checkout"
check "checkout A123 qty1"        "200" "$(code -X POST "$ORDER_BASE/checkout" -H 'Content-Type: application/json' -d '{"item_id":"A123","qty":1}')"
check "checkout UNKNOWN"          "500" "$(code -X POST "$ORDER_BASE/checkout" -H 'Content-Type: application/json' -d '{"item_id":"UNKNOWN","qty":1}')"
check "checkout C123 qty999"      "400" "$(code -X POST "$ORDER_BASE/checkout" -H 'Content-Type: application/json' -d '{"item_id":"C123","qty":999}')"
check "checkout invalid qty0"     "400" "$(code -X POST "$ORDER_BASE/checkout" -H 'Content-Type: application/json' -d '{"item_id":"A123","qty":0}')"

echo ""
echo "Results: $pass passed, $fail failed"
[ "$fail" -eq 0 ]
