#!/usr/bin/env bash
# Quality pipeline (F12.1): gofmt, vet, lint, tests, coverage gate >= 70%
# on ./internal/... per service. Exit non-zero on any failure.
set -uo pipefail

cd "$(dirname "$0")/.."

COVERAGE_MIN="${COVERAGE_MIN:-70}"
MODULES=(services/order services/inventory services/payment pkg/telemetry)
COVERAGE_MODULES=(services/order services/inventory services/payment)
FAILED=0

for m in "${MODULES[@]}"; do
  echo "==> [$m] gofmt"
  unformatted=$(cd "$m" && gofmt -l .)
  if [ -n "$unformatted" ]; then
    echo "FAIL: file tidak terformat:"
    echo "$unformatted"
    FAILED=1
  else
    echo "clean"
  fi

  echo "==> [$m] go vet"
  if ! (cd "$m" && go vet ./...); then
    FAILED=1
  fi

  echo "==> [$m] golangci-lint"
  if command -v golangci-lint >/dev/null 2>&1; then
    if ! (cd "$m" && golangci-lint run); then
      FAILED=1
    fi
  else
    echo "SKIP: golangci-lint tidak terinstall"
  fi

  echo "==> [$m] go test"
  if ! (cd "$m" && go test ./...); then
    FAILED=1
  fi
  echo ""
done

echo "==> Coverage gate (>= ${COVERAGE_MIN}% per package, ./internal/...)"
for m in "${COVERAGE_MODULES[@]}"; do
  cover_out=$(cd "$m" && go test ./internal/... -cover 2>/dev/null)
  min=$(echo "$cover_out" | python3 -c "
import sys, re
vals = [float(m.group(1)) for line in sys.stdin if (m := re.search(r'coverage: ([\d.]+)%', line))]
print(min(vals) if vals else 0)
")
  echo "$cover_out" | sed 's/^/  /'
  if python3 -c "import sys; sys.exit(0 if float('$min') >= float('$COVERAGE_MIN') else 1)"; then
    echo "PASS  [$m] coverage min ${min}% >= ${COVERAGE_MIN}%"
  else
    echo "FAIL  [$m] coverage min ${min}% < ${COVERAGE_MIN}%"
    FAILED=1
  fi
  echo ""
done

if [ "$FAILED" -eq 0 ]; then
  echo "ALL CHECKS PASSED"
else
  echo "CHECKS FAILED"
fi
exit "$FAILED"
