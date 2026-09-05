# OTel-Shop Lab

Laboratory project untuk mempelajari **OpenTelemetry distributed tracing** secara
end-to-end: 3 microservices Golang (Order, Inventory, Payment) di-deploy ke
**k3d Kubernetes**, dengan OpenTelemetry Collector dan Jaeger sebagai backend
tracing.

## Arsitektur

```
                    ┌──────────────────── k3d cluster (otel-shop) ────────────────┐
                    │                                                             │
host :18080 ───────▶│ order-service ──HTTP──▶ inventory-service ──SQL──▶ Postgres │
host :18081         │   │  otelhttp handler ▲ otelhttp server     otelsql         │
host :18082         │   │  otelhttp client  │                                     │
host :14317/14318   │   ▼                   │                                     │
host :16686         │ payment-service (chaos delay/error %)                       │
                    │   │                                                         │
                    │   ▼ OTLP gRPC 14317                                         │
                    │ otel-collector ──OTLP──▶ jaeger ──UI :16686──▶ host         │
                    └─────────────────────────────────────────────────────────────┘
```

Satu request `POST /checkout` menghasilkan **satu trace ID** lintas ketiga
service: `POST /checkout` → `validate-order` → `check-inventory` (+ DB span
`sql.*`) → `process-payment` → `POST /pay` (events + baggage `order.id`).

## Struktur Repo

```
services/          # 3 service (module Go terpisah)
  order/           #   POST /checkout — orchestrator
  inventory/       #   GET  /inventory/{id} — baca PostgreSQL
  payment/         #   POST /pay — chaos configurable
pkg/telemetry/     # shared OTel SDK setup (resource, exporter, sampler, propagator)
deploy/            # manifest K8s per komponen
db/init.sql        # schema + seed produk
scripts/           # cluster.sh build.sh deploy.sh test.sh chaos-test.sh check.sh
docs/              # progress-tracker, daily-log, tracing-examples, experiments
```

## Prerequisites

| Tool | Versi | Catatan |
|---|---|---|
| Go | 1.27 | workspace `go.work`, 4 module |
| Docker | daemon aktif | runtime k3d (k3d tidak support podman) |
| Podman | 6.x | build image (PRD D14) |
| k3d | 5.8+ | cluster lokal |
| kubectl | 1.28+ | context `k3d-otel-shop` |
| golangci-lint | 2.x | untuk `check.sh` (opsional) |

## Ports

| Port | Tujuan |
|---|---|
| 18080 | order-service (`POST /checkout`, `GET /health`) |
| 18081 | inventory-service (`GET /inventory/{id}`, `GET /health`) |
| 18082 | payment-service (`POST /pay`, `GET /health`) |
| 15432 | PostgreSQL (NodePort → 5432, db `oteldb`) |
| 14317 / 14318 | OTel Collector OTLP gRPC / HTTP |
| 16686 | Jaeger UI |

## Quickstart

```bash
# 1. Cluster + infra (Postgres, Collector, Jaeger)
./scripts/cluster.sh
./scripts/deploy.sh

# 2. Build image (podman → docker → k3d import) + deploy 3 service
./scripts/build.sh
./scripts/deploy.sh

# 3. Verifikasi
./scripts/test.sh          # E2E: health + checkout + error scenarios (9 checks)
./scripts/chaos-test.sh    # chaos: forced error + 2s delay di cluster
./scripts/check.sh         # gofmt/vet/lint/test + coverage gate >=70%
```

Lalu coba:

```bash
curl -X POST http://localhost:18080/checkout \
  -H "Content-Type: application/json" \
  -d '{"item_id":"A123","qty":1}'
# → {"order_id":"O-xxxx","status":"paid"}
```

Buka `http://localhost:16686` → service `order-service` → operation
`POST /checkout` untuk melihat distributed trace lengkap.

## Konfigurasi (env)

| Env | Default | Fungsi |
|---|---|---|
| `OTEL_EXPORTER_OTLP_ENDPOINT` | `otel-collector:14317` | tujuan OTLP exporter |
| `OTEL_TRACES_SAMPLER_ARG` | `1.0` | rasio sampling root trace (0.0–1.0) |
| `PAYMENT_DELAY_PERCENT` / `PAYMENT_DELAY_MS` | `20` / `2000` | chaos delay |
| `PAYMENT_ERROR_PERCENT` | `10` | chaos error |
| `INVENTORY_URL` / `PAYMENT_URL` | internal DNS | target downstream order |

## Dokumentasi

- [PRD](prd.md) — requirements & design
- [Implementation Plan](implementation_plan.md) — breakdown task
- [Progress Tracker](docs/progress-tracker.md) · [Daily Log](docs/daily-log.md)
- [Contoh Trace](docs/tracing-examples.md) — normal / slow / error
- [Eksperimen](docs/experiments.md) — 6 eksperimen inti observability
- [Glosarium](docs/glossary.md) — kamus istilah untuk pemula (ELI5)

## Teardown

```bash
k3d cluster delete otel-shop
```
