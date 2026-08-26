#!/usr/bin/env bash
set -euo pipefail

NS=otel-shop
KUBECTL="kubectl"
CLUSTER="${K3D_CLUSTER_NAME:-otel-shop}"
CLUSTER_NAME_NS="k3d-${CLUSTER}"

cd "$(dirname "$0")/.."

check_context() {
  local ctx
  ctx=$($KUBECTL config current-context 2>/dev/null || echo "")
  if [ "$ctx" != "$CLUSTER_NAME_NS" ]; then
    echo "ERROR: context saat ini '$ctx' bukan '$CLUSTER_NAME_NS'."
    echo "Jalankan: k3d kubeconfig merge '$CLUSTER' && kubectl config use-context '$CLUSTER_NAME_NS'"
    exit 1
  fi
}

apply() {
  echo "==> apply $1"
  $KUBECTL apply -f "$1" -n "$NS" --dry-run=client >/dev/null && $KUBECTL apply -f "$1"
}

check_context

# Urutan: namespace -> secret/config -> postgres -> jaeger -> collector
$KUBECTL apply -f deploy/namespace.yaml
apply deploy/postgres/secret.yaml
apply deploy/postgres/configmap.yaml
apply deploy/postgres/deployment.yaml
apply deploy/postgres/service.yaml
apply deploy/jaeger/deployment.yaml
apply deploy/jaeger/service.yaml
apply deploy/otel-collector/configmap.yaml
apply deploy/otel-collector/deployment.yaml
apply deploy/otel-collector/service.yaml

echo "==> Menunggu pods ready..."
$KUBECTL rollout status deploy/postgres -n "$NS" --timeout 120s
$KUBECTL rollout status deploy/jaeger -n "$NS" --timeout 120s
$KUBECTL rollout status deploy/otel-collector -n "$NS" --timeout 120s

echo "==> Status:"
$KUBECTL get pods,svc -n "$NS"