#!/usr/bin/env bash
set -euo pipefail

CLUSTER="${K3D_CLUSTER_NAME:-otel-shop}"

# Cek port bentrok di host
for p in 18080 18081 18082 15432 14317 14318 16686; do
  if lsof -nP -iTCP:"$p" -sTCP:LISTEN >/dev/null 2>&1; then
    echo "WARN: port $p sudah dipakai di host — k3d mapping akan gagal." >&2
  fi
done

if k3d cluster list | grep -q "^${CLUSTER}[[:space:]]"; then
  echo "Cluster '$CLUSTER' sudah ada — skip create."
  k3d cluster start "$CLUSTER"
  exit 0
fi

k3d cluster create "$CLUSTER" \
  --servers 1 --agents 1 \
  --k3s-arg '--service-node-port-range=10000-32767@server:0' \
  --port '18080:18080@server:0' \
  --port '18081:18081@server:0' \
  --port '18082:18082@server:0' \
  --port '15432:15432@server:0' \
  --port '14317:14317@server:0' \
  --port '14318:14318@server:0' \
  --port '16686:16686@server:0' \
  --wait

echo "Cluster '$CLUSTER' siap. kubectl context: k3d-$CLUSTER"