#!/usr/bin/env bash
set -euo pipefail

# Build service images with Podman (PRD D14), load into the Docker daemon
# used by k3d, and import into the cluster.
cd "$(dirname "$0")/.."

CLUSTER="${K3D_CLUSTER_NAME:-otel-shop}"
SERVICES=(order inventory payment)

if ! command -v podman >/dev/null 2>&1; then echo "ERROR: podman tidak ditemukan" >&2; exit 1; fi

for svc in "${SERVICES[@]}"; do
  echo "==> podman build otel-shop/${svc}:local (context: repo root)"
  podman build -t "otel-shop/${svc}:local" -f "services/${svc}/Dockerfile" .
done

echo "==> load images ke docker daemon (runtime k3d)"
mkdir -p /tmp/otel-shop-images
for svc in "${SERVICES[@]}"; do
  tar="/tmp/otel-shop-images/${svc}.tar"
  podman save "otel-shop/${svc}:local" -o "$tar" --format docker-archive
  docker load -i "$tar"
  # podman menandai image dengan prefix localhost/ — samakan ke nama tanpa prefix
  if docker image inspect "localhost/otel-shop/${svc}:local" >/dev/null 2>&1; then
    docker tag "localhost/otel-shop/${svc}:local" "otel-shop/${svc}:local"
  fi
done

echo "==> k3d image import ke cluster '$CLUSTER'"
k3d image import otel-shop/order:local otel-shop/inventory:local otel-shop/payment:local -c "$CLUSTER"

echo "Build & import selesai."