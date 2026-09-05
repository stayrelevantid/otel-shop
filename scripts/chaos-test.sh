#!/usr/bin/env bash
# Chaos test (F12.5): verify payment chaos knobs against the cluster.
# Scenario 1: ERROR_PERCENT=100   -> checkout must fail (500)
# Scenario 2: DELAY_PERCENT=100   -> checkout must succeed but take ~2s
# Then revert to deployment defaults. Exit non-zero on any failure.
set -uo pipefail

NS=otel-shop
ORDER_URL="${ORDER_URL:-http://localhost:18080}"
FAILED=0

set_env() { kubectl -n "$NS" set env deploy/payment-service "$@" >/dev/null; }
wait_rollout() { kubectl -n "$NS" rollout status deploy/payment-service --timeout 90s >/dev/null; }
checkout_code() {
  curl -s -o /dev/null -w '%{http_code}' \
    -X POST "$ORDER_URL/checkout" \
    -H 'Content-Type: application/json' \
    -d '{"item_id":"A123","qty":1}'
}

echo "==> baseline"
set_env PAYMENT_ERROR_PERCENT=10 PAYMENT_DELAY_PERCENT=20 PAYMENT_DELAY_MS=2000
wait_rollout

echo "==> scenario 1: ERROR_PERCENT=100 -> checkout harus 500"
set_env PAYMENT_ERROR_PERCENT=100
wait_rollout
sleep 1
code=$(checkout_code)
if [ "$code" = "500" ]; then
  echo "PASS  forced error (got $code)"
else
  echo "FAIL  expected 500, got $code"
  FAILED=1
fi

echo "==> scenario 2: DELAY_PERCENT=100, DELAY_MS=2000 -> 200 dan ~2s"
set_env PAYMENT_ERROR_PERCENT=0 PAYMENT_DELAY_PERCENT=100 PAYMENT_DELAY_MS=2000
wait_rollout
sleep 1
start=$(date +%s%N)
code=$(checkout_code)
ms=$(( ($(date +%s%N) - start) / 1000000 ))
if [ "$code" = "200" ]; then
  echo "PASS  delayed checkout 200 (${ms}ms)"
else
  echo "FAIL  expected 200, got $code (${ms}ms)"
  FAILED=1
fi
if [ "$ms" -ge 1900 ]; then
  echo "PASS  delay terukur >= 1900ms (${ms}ms)"
else
  echo "FAIL  delay hanya ${ms}ms (expected >= 1900ms)"
  FAILED=1
fi

echo "==> revert ke default (error 10%, delay 20%)"
set_env PAYMENT_ERROR_PERCENT=10 PAYMENT_DELAY_PERCENT=20 PAYMENT_DELAY_MS=2000
wait_rollout
code=$(checkout_code)
if [ "$code" = "200" ]; then
  echo "PASS  revert checkout 200"
else
  echo "FAIL  revert got $code"
  FAILED=1
fi

if [ "$FAILED" -eq 0 ]; then
  echo "CHAOS TEST PASSED"
else
  echo "CHAOS TEST FAILED"
fi
exit "$FAILED"
